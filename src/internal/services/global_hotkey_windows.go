//go:build windows

package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"wepassword/internal/logger"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
)

// Windows虚拟键码常量
const (
	VK_SPACE  = 0x20
	VK_RETURN = 0x0D
	VK_TAB    = 0x09
	VK_ESCAPE = 0x1B
	VK_BACK   = 0x08
	VK_DELETE = 0x2E
	VK_F1     = 0x70
	VK_F2     = 0x71
	VK_F3     = 0x72
	VK_F4     = 0x73
	VK_F5     = 0x74
	VK_F6     = 0x75
	VK_F7     = 0x76
	VK_F8     = 0x77
	VK_F9     = 0x78
	VK_F10    = 0x79
	VK_F11    = 0x7A
	VK_F12    = 0x7B
)

// Windows修饰键常量
const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008
)

// Windows消息常量
const (
	WM_HOTKEY = 0x0312
	WM_QUIT   = 0x0012
)

// MSG 结构体
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct {
		X, Y int32
	}
}

/**
 * GlobalHotkeyWindows Windows平台全局快捷键实现
 * @author 陈凤庆
 * @date 20251014
 */
type GlobalHotkeyWindows struct {
	registeredHotkeys map[int]func() // 已注册的快捷键回调，键为快捷键ID
	hotkeyCounter     int            // 快捷键ID计数器
	mu                sync.RWMutex   // 读写锁
	isListening       bool           // 是否正在监听
	stopChan          chan struct{}  // 停止信号
}

/**
 * newGlobalHotkeyWindows 创建Windows平台实现
 * @return GlobalHotkeyPlatformImpl 平台实现
 * @return error 错误信息
 */
func newGlobalHotkeyWindows() (GlobalHotkeyPlatformImpl, error) {
	return &GlobalHotkeyWindows{
		registeredHotkeys: make(map[int]func()),
		hotkeyCounter:     1,
		stopChan:          make(chan struct{}),
	}, nil
}

/**
 * RegisterHotkey 注册全局快捷键
 * @param hotkey 快捷键字符串
 * @param callback 回调函数
 * @return error 错误信息
 */
func (ghw *GlobalHotkeyWindows) RegisterHotkey(hotkey string, callback func()) error {
	logger.Info("[全局快捷键-Windows] 注册快捷键: %s", hotkey)

	// 解析快捷键
	modifiers, key, err := ghw.ParseHotkey(hotkey)
	if err != nil {
		return fmt.Errorf("解析快捷键失败: %w", err)
	}

	// 转换修饰键
	winModifiers := ghw.convertModifiers(modifiers)

	// 转换键码
	vkCode, err := ghw.convertKeyCode(key)
	if err != nil {
		return fmt.Errorf("转换键码失败: %w", err)
	}

	ghw.mu.Lock()
	hotkeyID := ghw.hotkeyCounter
	ghw.hotkeyCounter++
	ghw.mu.Unlock()

	// 注册快捷键
	ret, _, err := procRegisterHotKey.Call(
		0,                     // hWnd (NULL for global)
		uintptr(hotkeyID),     // id
		uintptr(winModifiers), // fsModifiers
		uintptr(vkCode),       // vk
	)

	if ret == 0 {
		return fmt.Errorf("注册快捷键失败: %v", err)
	}

	// 保存回调函数
	ghw.mu.Lock()
	ghw.registeredHotkeys[hotkeyID] = callback
	ghw.mu.Unlock()

	logger.Info("[全局快捷键-Windows] 快捷键注册成功: %s (ID: %d)", hotkey, hotkeyID)
	return nil
}

/**
 * UnregisterHotkey 取消注册全局快捷键
 * @param hotkey 快捷键字符串
 * @return error 错误信息
 */
func (ghw *GlobalHotkeyWindows) UnregisterHotkey(hotkey string) error {
	logger.Info("[全局快捷键-Windows] 取消注册快捷键: %s", hotkey)

	ghw.mu.Lock()
	defer ghw.mu.Unlock()

	// 找到并取消注册所有匹配的快捷键
	for hotkeyID := range ghw.registeredHotkeys {
		ret, _, err := procUnregisterHotKey.Call(0, uintptr(hotkeyID))
		if ret == 0 {
			logger.Error("[全局快捷键-Windows] 取消注册快捷键失败 (ID: %d): %v", hotkeyID, err)
		}
		delete(ghw.registeredHotkeys, hotkeyID)
	}

	logger.Info("[全局快捷键-Windows] 快捷键取消注册成功: %s", hotkey)
	return nil
}

/**
 * StartListening 开始监听快捷键事件
 * @param ctx 上下文
 * @return error 错误信息
 */
func (ghw *GlobalHotkeyWindows) StartListening(ctx context.Context) error {
	logger.Info("[全局快捷键-Windows] 开始监听快捷键事件")

	ghw.mu.Lock()
	if ghw.isListening {
		ghw.mu.Unlock()
		return fmt.Errorf("已在监听中")
	}
	ghw.isListening = true
	ghw.mu.Unlock()

	// 启动消息循环
	go ghw.messageLoop(ctx)

	return nil
}

/**
 * StopListening 停止监听快捷键事件
 * @return error 错误信息
 */
func (ghw *GlobalHotkeyWindows) StopListening() error {
	logger.Info("[全局快捷键-Windows] 停止监听")

	ghw.mu.Lock()
	if !ghw.isListening {
		ghw.mu.Unlock()
		return nil
	}
	ghw.isListening = false
	ghw.mu.Unlock()

	// 发送停止信号
	select {
	case ghw.stopChan <- struct{}{}:
	default:
	}

	// 发送退出消息
	procPostQuitMessage.Call(0)

	return nil
}

/**
 * IsSupported 检查当前平台是否支持全局快捷键
 * @return bool 是否支持
 */
func (ghw *GlobalHotkeyWindows) IsSupported() bool {
	return true
}

/**
 * ParseHotkey 解析快捷键字符串
 * @param hotkey 快捷键字符串，如"Ctrl+Space"
 * @return []string 修饰键列表
 * @return string 主键
 * @return error 错误信息
 */
func (ghw *GlobalHotkeyWindows) ParseHotkey(hotkey string) ([]string, string, error) {
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
		case "alt":
			modifiers = append(modifiers, "Alt")
		case "shift":
			modifiers = append(modifiers, "Shift")
		case "win", "windows":
			modifiers = append(modifiers, "Win")
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
 * convertModifiers 转换修饰键为Windows常量
 * @param modifiers 修饰键列表
 * @return uint32 Windows修饰键常量
 */
func (ghw *GlobalHotkeyWindows) convertModifiers(modifiers []string) uint32 {
	var winModifiers uint32
	for _, modifier := range modifiers {
		switch modifier {
		case "Ctrl":
			winModifiers |= MOD_CONTROL
		case "Alt":
			winModifiers |= MOD_ALT
		case "Shift":
			winModifiers |= MOD_SHIFT
		case "Win":
			winModifiers |= MOD_WIN
		}
	}
	return winModifiers
}

/**
 * convertKeyCode 转换键名为Windows虚拟键码
 * @param key 键名
 * @return uint32 虚拟键码
 * @return error 错误信息
 */
func (ghw *GlobalHotkeyWindows) convertKeyCode(key string) (uint32, error) {
	switch strings.ToUpper(key) {
	case "SPACE":
		return VK_SPACE, nil
	case "RETURN", "ENTER":
		return VK_RETURN, nil
	case "TAB":
		return VK_TAB, nil
	case "ESCAPE":
		return VK_ESCAPE, nil
	case "BACKSPACE":
		return VK_BACK, nil
	case "DELETE":
		return VK_DELETE, nil
	case "F1":
		return VK_F1, nil
	case "F2":
		return VK_F2, nil
	case "F3":
		return VK_F3, nil
	case "F4":
		return VK_F4, nil
	case "F5":
		return VK_F5, nil
	case "F6":
		return VK_F6, nil
	case "F7":
		return VK_F7, nil
	case "F8":
		return VK_F8, nil
	case "F9":
		return VK_F9, nil
	case "F10":
		return VK_F10, nil
	case "F11":
		return VK_F11, nil
	case "F12":
		return VK_F12, nil
	}

	// 数字键 0-9
	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		return uint32(key[0]), nil
	}

	// 字母键 A-Z
	if len(key) == 1 && ((key[0] >= 'A' && key[0] <= 'Z') || (key[0] >= 'a' && key[0] <= 'z')) {
		return uint32(strings.ToUpper(key)[0]), nil
	}

	return 0, fmt.Errorf("不支持的键: %s", key)
}

/**
 * messageLoop Windows消息循环
 * @param ctx 上下文
 */
func (ghw *GlobalHotkeyWindows) messageLoop(ctx context.Context) {
	logger.Info("[全局快捷键-Windows] 启动消息循环")

	for {
		select {
		case <-ctx.Done():
			logger.Info("[全局快捷键-Windows] 上下文取消，退出消息循环")
			return
		case <-ghw.stopChan:
			logger.Info("[全局快捷键-Windows] 收到停止信号，退出消息循环")
			return
		default:
			// 获取消息
			var msg MSG
			ret, _, _ := procGetMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				0, // hWnd
				0, // wMsgFilterMin
				0, // wMsgFilterMax
			)

			if ret == 0 { // WM_QUIT
				logger.Info("[全局快捷键-Windows] 收到WM_QUIT消息，退出消息循环")
				return
			}

			if ret == ^uintptr(0) { // -1, error
				logger.Error("[全局快捷键-Windows] GetMessage失败")
				continue
			}

			// 处理热键消息
			if msg.Message == WM_HOTKEY {
				hotkeyID := int(msg.WParam)
				ghw.mu.RLock()
				callback, exists := ghw.registeredHotkeys[hotkeyID]
				ghw.mu.RUnlock()

				if exists && callback != nil {
					logger.Info("[全局快捷键-Windows] 触发快捷键 (ID: %d)", hotkeyID)
					go callback() // 在新的goroutine中执行回调，避免阻塞消息循环
				}
			}

			// 翻译和分发消息
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
}
