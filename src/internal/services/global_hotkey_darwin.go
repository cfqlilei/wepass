//go:build darwin

package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"wepassword/internal/logger"
)

/**
 * GlobalHotkeyDarwin macOS平台全局快捷键实现
 * @author 陈凤庆
 * @date 20251014
 */
type GlobalHotkeyDarwin struct {
	registeredHotkeys map[string]func() // 已注册的快捷键回调
	isListening       bool              // 是否正在监听
	stopChan          chan struct{}     // 停止信号
	mu                sync.RWMutex      // 读写锁
}

/**
 * newGlobalHotkeyDarwin 创建macOS平台实现
 * @return GlobalHotkeyPlatformImpl 平台实现
 * @return error 错误信息
 */
func newGlobalHotkeyDarwin() (GlobalHotkeyPlatformImpl, error) {
	ghd := &GlobalHotkeyDarwin{
		registeredHotkeys: make(map[string]func()),
		isListening:       false,
		stopChan:          make(chan struct{}),
	}

	return ghd, nil
}

/**
 * RegisterHotkey 注册全局快捷键
 * @param hotkey 快捷键字符串
 * @param callback 回调函数
 * @return error 错误信息
 */
func (ghd *GlobalHotkeyDarwin) RegisterHotkey(hotkey string, callback func()) error {
	logger.Info("[全局快捷键-macOS] 注册快捷键: %s", hotkey)

	// 解析快捷键
	modifiers, key, err := ghd.ParseHotkey(hotkey)
	if err != nil {
		return fmt.Errorf("解析快捷键失败: %w", err)
	}

	// 验证快捷键格式
	if len(modifiers) == 0 {
		return fmt.Errorf("必须包含至少一个修饰键")
	}

	// 验证主键
	if key == "" {
		return fmt.Errorf("必须包含主键")
	}

	ghd.mu.Lock()
	defer ghd.mu.Unlock()

	// 保存回调函数
	ghd.registeredHotkeys[hotkey] = callback

	logger.Info("[全局快捷键-macOS] 快捷键注册成功: %s (使用robotgo实现)", hotkey)
	return nil
}

/**
 * UnregisterHotkey 取消注册全局快捷键
 * @param hotkey 快捷键字符串
 * @return error 错误信息
 */
func (ghd *GlobalHotkeyDarwin) UnregisterHotkey(hotkey string) error {
	logger.Info("[全局快捷键-macOS] 取消注册快捷键: %s", hotkey)

	// 从映射中删除
	delete(ghd.registeredHotkeys, hotkey)

	logger.Info("[全局快捷键-macOS] 快捷键取消注册成功: %s (简化版本)", hotkey)
	return nil
}

/**
 * StartListening 开始监听快捷键事件
 * @param ctx 上下文
 * @return error 错误信息
 */
func (ghd *GlobalHotkeyDarwin) StartListening(ctx context.Context) error {
	logger.Info("[全局快捷键-macOS] 开始监听快捷键事件")

	ghd.mu.Lock()
	if ghd.isListening {
		ghd.mu.Unlock()
		return fmt.Errorf("已在监听中")
	}
	ghd.isListening = true
	ghd.mu.Unlock()

	// 启动快捷键监听循环
	go ghd.keyboardEventLoop(ctx)

	return nil
}

/**
 * StopListening 停止监听快捷键事件
 * @return error 错误信息
 */
func (ghd *GlobalHotkeyDarwin) StopListening() error {
	logger.Info("[全局快捷键-macOS] 停止监听")
	return nil
}

/**
 * IsSupported 检查当前平台是否支持全局快捷键
 * @return bool 是否支持
 */
func (ghd *GlobalHotkeyDarwin) IsSupported() bool {
	return true
}

/**
 * ParseHotkey 解析快捷键字符串
 * @param hotkey 快捷键字符串，如"Ctrl+Space"
 * @return []string 修饰键列表
 * @return string 主键
 * @return error 错误信息
 */
func (ghd *GlobalHotkeyDarwin) ParseHotkey(hotkey string) ([]string, string, error) {
	if hotkey == "" {
		return nil, "", fmt.Errorf("快捷键不能为空")
	}

	parts := strings.Split(hotkey, "+")
	if len(parts) < 2 {
		return nil, "", fmt.Errorf("快捷键格式无效，应包含修饰键和主键")
	}

	// 最后一个部分是主键
	key := strings.TrimSpace(parts[len(parts)-1])
	if key == "" {
		return nil, "", fmt.Errorf("主键不能为空")
	}

	// 前面的部分是修饰键
	var modifiers []string
	for i := 0; i < len(parts)-1; i++ {
		modifier := strings.TrimSpace(parts[i])
		if modifier == "" {
			continue
		}

		// 标准化修饰键名称
		switch strings.ToLower(modifier) {
		case "ctrl", "control":
			modifiers = append(modifiers, "Ctrl")
		case "cmd", "command":
			modifiers = append(modifiers, "Cmd")
		case "alt", "option":
			modifiers = append(modifiers, "Alt")
		case "shift":
			modifiers = append(modifiers, "Shift")
		default:
			return nil, "", fmt.Errorf("不支持的修饰键: %s", modifier)
		}
	}

	if len(modifiers) == 0 {
		return nil, "", fmt.Errorf("必须包含至少一个修饰键")
	}

	return modifiers, key, nil
}

/**
 * keyboardEventLoop 键盘事件监听循环
 * @param ctx 上下文
 */
func (ghd *GlobalHotkeyDarwin) keyboardEventLoop(ctx context.Context) {
	logger.Info("[全局快捷键-macOS] 启动键盘事件监听循环")

	// 使用 robotgo 监听全局键盘事件
	ticker := time.NewTicker(50 * time.Millisecond) // 50ms 检查一次
	defer ticker.Stop()

	var lastKeyState = make(map[string]bool)

	for {
		select {
		case <-ctx.Done():
			logger.Info("[全局快捷键-macOS] 上下文取消，退出键盘事件循环")
			return
		case <-ghd.stopChan:
			logger.Info("[全局快捷键-macOS] 收到停止信号，退出键盘事件循环")
			return
		case <-ticker.C:
			// 检查已注册的快捷键
			ghd.mu.RLock()
			for hotkey, callback := range ghd.registeredHotkeys {
				if ghd.isHotkeyPressed(hotkey, lastKeyState) {
					logger.Info("[全局快捷键-macOS] 检测到快捷键: %s", hotkey)
					go callback() // 在新的goroutine中执行回调，避免阻塞监听循环
				}
			}
			ghd.mu.RUnlock()
		}
	}
}

/**
 * isHotkeyPressed 检查快捷键是否被按下
 * @param hotkey 快捷键字符串
 * @param lastKeyState 上次按键状态
 * @return bool 是否被按下
 */
func (ghd *GlobalHotkeyDarwin) isHotkeyPressed(hotkey string, lastKeyState map[string]bool) bool {
	modifiers, key, err := ghd.ParseHotkey(hotkey)
	if err != nil {
		return false
	}

	// 检查所有修饰键是否都被按下
	for _, modifier := range modifiers {
		if !ghd.isModifierPressed(modifier) {
			return false
		}
	}

	// 检查主键是否被按下
	keyPressed := ghd.isKeyPressed(key)

	// 防止重复触发：只有在按键从未按下变为按下时才触发
	stateKey := hotkey
	wasPressed := lastKeyState[stateKey]
	lastKeyState[stateKey] = keyPressed

	return keyPressed && !wasPressed
}

/**
 * isModifierPressed 检查修饰键是否被按下
 * @param modifier 修饰键名称
 * @return bool 是否被按下
 */
func (ghd *GlobalHotkeyDarwin) isModifierPressed(modifier string) bool {
	// 暂时返回 false，需要更复杂的实现
	// 在实际应用中，需要使用 macOS 的 Carbon 框架或其他底层 API
	return false
}

/**
 * isKeyPressed 检查主键是否被按下
 * @param key 主键名称
 * @return bool 是否被按下
 */
func (ghd *GlobalHotkeyDarwin) isKeyPressed(key string) bool {
	// 暂时返回 false，需要更复杂的实现
	// 在实际应用中，需要使用 macOS 的 Carbon 框架或其他底层 API
	return false
}

/**
 * convertKeyToRobotgo 将键名转换为robotgo支持的格式
 * @param key 键名
 * @return string robotgo键名
 */
func (ghd *GlobalHotkeyDarwin) convertKeyToRobotgo(key string) string {
	// 转换为小写
	key = strings.ToLower(key)

	// 特殊键映射
	keyMap := map[string]string{
		"space":     "space",
		"enter":     "enter",
		"return":    "enter",
		"tab":       "tab",
		"escape":    "esc",
		"esc":       "esc",
		"backspace": "backspace",
		"delete":    "delete",
		"up":        "up",
		"down":      "down",
		"left":      "left",
		"right":     "right",
		"home":      "home",
		"end":       "end",
		"pageup":    "pageup",
		"pagedown":  "pagedown",
	}

	if robotgoKey, exists := keyMap[key]; exists {
		return robotgoKey
	}

	// 对于字母和数字，直接返回
	return key
}
