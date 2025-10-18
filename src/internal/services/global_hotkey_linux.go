//go:build linux

package services

import (
	"context"
	"fmt"
	"strings"

	"wepassword/internal/logger"
)

/**
 * GlobalHotkeyLinux Linux平台全局快捷键实现
 * @author 陈凤庆
 * @date 20251014
 * @description Linux平台的全局快捷键实现，目前为简化版本
 */
type GlobalHotkeyLinux struct {
	registeredHotkeys map[string]func() // 已注册的快捷键回调
}

/**
 * newGlobalHotkeyLinux 创建Linux平台实现
 * @return GlobalHotkeyPlatformImpl 平台实现
 * @return error 错误信息
 */
func newGlobalHotkeyLinux() (GlobalHotkeyPlatformImpl, error) {
	return &GlobalHotkeyLinux{
		registeredHotkeys: make(map[string]func()),
	}, nil
}

/**
 * RegisterHotkey 注册全局快捷键
 * @param hotkey 快捷键字符串
 * @param callback 回调函数
 * @return error 错误信息
 */
func (ghl *GlobalHotkeyLinux) RegisterHotkey(hotkey string, callback func()) error {
	logger.Info("[全局快捷键-Linux] 注册快捷键: %s", hotkey)

	// 解析快捷键
	modifiers, key, err := ghl.ParseHotkey(hotkey)
	if err != nil {
		return fmt.Errorf("解析快捷键失败: %w", err)
	}

	logger.Info("[全局快捷键-Linux] 解析结果 - 修饰键: %v, 主键: %s", modifiers, key)

	// 保存回调函数
	ghl.registeredHotkeys[hotkey] = callback

	// TODO: 实现Linux平台的实际快捷键注册
	// 可以使用X11库或者其他Linux特定的API
	logger.Warn("[全局快捷键-Linux] Linux平台快捷键注册暂未完全实现，仅保存配置")

	logger.Info("[全局快捷键-Linux] 快捷键注册成功: %s", hotkey)
	return nil
}

/**
 * UnregisterHotkey 取消注册全局快捷键
 * @param hotkey 快捷键字符串
 * @return error 错误信息
 */
func (ghl *GlobalHotkeyLinux) UnregisterHotkey(hotkey string) error {
	logger.Info("[全局快捷键-Linux] 取消注册快捷键: %s", hotkey)

	// 从映射中删除
	delete(ghl.registeredHotkeys, hotkey)

	// TODO: 实现实际的取消注册逻辑
	logger.Info("[全局快捷键-Linux] 快捷键取消注册成功: %s", hotkey)
	return nil
}

/**
 * StartListening 开始监听快捷键事件
 * @param ctx 上下文
 * @return error 错误信息
 */
func (ghl *GlobalHotkeyLinux) StartListening(ctx context.Context) error {
	logger.Info("[全局快捷键-Linux] 开始监听快捷键事件")

	// TODO: 实现Linux平台的事件监听
	// 可以使用X11事件循环或者其他Linux特定的事件机制

	// 目前只是等待上下文取消
	<-ctx.Done()

	logger.Info("[全局快捷键-Linux] 停止监听快捷键事件")
	return nil
}

/**
 * StopListening 停止监听快捷键事件
 * @return error 错误信息
 */
func (ghl *GlobalHotkeyLinux) StopListening() error {
	logger.Info("[全局快捷键-Linux] 停止监听")
	return nil
}

/**
 * IsSupported 检查当前平台是否支持全局快捷键
 * @return bool 是否支持
 */
func (ghl *GlobalHotkeyLinux) IsSupported() bool {
	// Linux平台支持有限，返回false表示功能不完整
	logger.Warn("[全局快捷键-Linux] Linux平台快捷键支持有限")
	return false
}

/**
 * ParseHotkey 解析快捷键字符串
 * @param hotkey 快捷键字符串，如"Ctrl+Space"
 * @return []string 修饰键列表
 * @return string 主键
 * @return error 错误信息
 */
func (ghl *GlobalHotkeyLinux) ParseHotkey(hotkey string) ([]string, string, error) {
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
		case "super", "win", "windows":
			modifiers = append(modifiers, "Super")
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
 * GetRegisteredHotkeys 获取已注册的快捷键列表（用于调试）
 * @return map[string]func() 已注册的快捷键映射
 */
func (ghl *GlobalHotkeyLinux) GetRegisteredHotkeys() map[string]func() {
	return ghl.registeredHotkeys
}
