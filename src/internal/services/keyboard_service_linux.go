//go:build linux
// +build linux

package services

import (
	"errors"
	"fmt"
	"log"
)

/**
 * KeyboardService Linux 版本的键盘服务
 * @author 陈凤庆
 * @date 20251002
 */
type KeyboardService struct {
	windowMonitor *WindowMonitorService // 窗口监控服务
}

/**
 * NewKeyboardService 创建新的键盘服务实例
 * @return *KeyboardService 键盘服务实例
 */
func NewKeyboardService() *KeyboardService {
	log.Println("[KeyboardService] 初始化 Linux 键盘服务...")

	service := &KeyboardService{}

	// 初始化窗口监控服务
	service.windowMonitor = NewWindowMonitorService()

	log.Println("[KeyboardService] Linux 键盘服务初始化完成")
	return service
}

/**
 * CheckAccessibilityPermission 检查辅助功能权限（Linux 版本暂不支持）
 * @return bool 权限状态
 */
func (ks *KeyboardService) CheckAccessibilityPermission() bool {
	// Linux 版本暂时返回 true，实际使用时可能需要检查 X11 权限
	return true
}

/**
 * StorePreviousFocusedApp 存储上一次聚焦的应用程序（Linux 版本暂不支持）
 */
func (ks *KeyboardService) StorePreviousFocusedApp() {
	// Linux 版本暂不实现
	log.Println("[KeyboardService] Linux 版本暂不支持存储上一次聚焦的应用程序")
}

/**
 * RestorePreviousFocusedApp 恢复到上一次聚焦的应用程序（Linux 版本暂不支持）
 * @return error 错误信息
 */
func (ks *KeyboardService) RestorePreviousFocusedApp() error {
	// Linux 版本暂不实现
	return errors.New("Linux 版本暂不支持恢复到上一次聚焦的应用程序")
}

/**
 * GetPreviousFocusedAppName 获取上一次聚焦的应用程序名称（Linux 版本暂不支持）
 * @return string 应用程序名称
 */
func (ks *KeyboardService) GetPreviousFocusedAppName() string {
	// Linux 版本暂不实现
	return "Linux 版本暂不支持"
}

/**
 * GetWindowMonitor 获取窗口监控服务（Linux版本）
 * @return *WindowMonitorService 窗口监控服务实例
 * @description 20251005 陈凤庆 提供窗口监控服务访问接口
 */
func (ks *KeyboardService) GetWindowMonitor() *WindowMonitorService {
	return ks.windowMonitor
}

/**
 * SimulateText 模拟文本输入（Linux 版本暂不支持）
 * @param text 要输入的文本
 * @return error 错误信息
 */
func (ks *KeyboardService) SimulateText(text string) error {
	// Linux 版本暂不实现，可以考虑使用 xdotool 或其他工具
	return fmt.Errorf("Linux 版本暂不支持模拟文本输入: %s", text)
}

/**
 * SimulateTab 模拟 Tab 键输入（Linux 版本暂不支持）
 * @return error 错误信息
 */
func (ks *KeyboardService) SimulateTab() error {
	// Linux 版本暂不实现
	return errors.New("Linux 版本暂不支持模拟 Tab 键输入")
}

/**
 * SimulateUsernameAndPassword 模拟用户名和密码输入（Linux 版本暂不支持）
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 */
func (ks *KeyboardService) SimulateUsernameAndPassword(username, password string) error {
	// Linux 版本暂不实现
	return fmt.Errorf("Linux 版本暂不支持模拟用户名和密码输入: %s", username)
}

/**
 * SimulateTextWithAutoSwitch 模拟文本输入并自动切换窗口（Linux 版本暂不支持）
 * @param text 要输入的文本
 * @return error 错误信息
 */
func (ks *KeyboardService) SimulateTextWithAutoSwitch(text string) error {
	// Linux 版本暂不实现
	return fmt.Errorf("Linux 版本暂不支持模拟文本输入并自动切换窗口: %s", text)
}

/**
 * SimulateUsernameAndPasswordWithAutoSwitch 模拟用户名和密码输入并自动切换窗口（Linux 版本暂不支持）
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 */
func (ks *KeyboardService) SimulateUsernameAndPasswordWithAutoSwitch(username, password string) error {
	// Linux 版本暂不实现
	return fmt.Errorf("Linux 版本暂不支持模拟用户名和密码输入并自动切换窗口: %s", username)
}

/**
 * RecordCurrentWindow 记录当前窗口信息
 * @return error 错误信息
 */
func (ks *KeyboardService) RecordCurrentWindow() error {
	if ks.windowMonitor != nil {
		return ks.windowMonitor.RecordCurrentWindow()
	}
	return errors.New("窗口监控服务未初始化")
}

/**
 * GetLastWindowInfo 获取最后记录的窗口信息
 * @return *WindowInfo 窗口信息
 */
func (ks *KeyboardService) GetLastWindowInfo() *WindowInfo {
	if ks.windowMonitor != nil {
		return ks.windowMonitor.GetLastWindowInfo()
	}
	return nil
}

/**
 * StopWindowMonitoring 停止窗口监控
 * @return error 错误信息
 */
func (ks *KeyboardService) StopWindowMonitoring() error {
	if ks.windowMonitor != nil {
		return ks.windowMonitor.Stop()
	}
	return nil
}

/**
 * SwitchToLastWindow 切换到最后记录的窗口
 * @return bool 是否成功
 */
func (ks *KeyboardService) SwitchToLastWindow() bool {
	if ks.windowMonitor != nil {
		return ks.windowMonitor.SwitchToLastWindow()
	}
	return false
}

/**
 * SwitchToPasswordManager 切换到密码管理器窗口
 * @return bool 是否成功
 */
func (ks *KeyboardService) SwitchToPasswordManager() bool {
	if ks.windowMonitor != nil {
		return ks.windowMonitor.SwitchToPasswordManager()
	}
	return false
}
