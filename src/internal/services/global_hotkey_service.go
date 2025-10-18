package services

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"wepassword/internal/config"
	"wepassword/internal/logger"
	"wepassword/internal/models"
)

/**
 * GlobalHotkeyService 全局快捷键服务
 * @author 陈凤庆
 * @date 20251014
 * @description 提供跨平台的全局快捷键监听和处理功能，支持显示/隐藏窗口
 */
type GlobalHotkeyService struct {
	configManager    *config.ConfigManager
	platformImpl     GlobalHotkeyPlatformImpl
	isRunning        bool
	mu               sync.RWMutex
	ctx              context.Context
	cancel           context.CancelFunc
	windowController WindowController // 窗口控制接口
}

/**
 * GlobalHotkeyPlatformImpl 平台特定实现接口
 */
type GlobalHotkeyPlatformImpl interface {
	// RegisterHotkey 注册全局快捷键
	RegisterHotkey(hotkey string, callback func()) error
	// UnregisterHotkey 取消注册全局快捷键
	UnregisterHotkey(hotkey string) error
	// StartListening 开始监听快捷键事件
	StartListening(ctx context.Context) error
	// StopListening 停止监听快捷键事件
	StopListening() error
	// IsSupported 检查当前平台是否支持全局快捷键
	IsSupported() bool
	// ParseHotkey 解析快捷键字符串
	ParseHotkey(hotkey string) (modifiers []string, key string, err error)
}

/**
 * WindowController 窗口控制接口
 */
type WindowController interface {
	// ShowWindow 显示窗口
	ShowWindow() error
	// HideWindow 隐藏窗口
	HideWindow() error
	// IsWindowVisible 检查窗口是否可见
	IsWindowVisible() bool
	// ToggleWindow 切换窗口显示状态
	ToggleWindow() error
}

/**
 * NewGlobalHotkeyService 创建全局快捷键服务实例
 * @param configManager 配置管理器
 * @param windowController 窗口控制器
 * @return *GlobalHotkeyService 服务实例
 * @return error 错误信息
 */
func NewGlobalHotkeyService(configManager *config.ConfigManager, windowController WindowController) (*GlobalHotkeyService, error) {
	// 创建平台特定实现
	platformImpl, err := newGlobalHotkeyPlatformImpl()
	if err != nil {
		return nil, fmt.Errorf("创建平台实现失败: %w", err)
	}

	// 检查平台支持
	if !platformImpl.IsSupported() {
		return nil, fmt.Errorf("当前平台不支持全局快捷键功能")
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &GlobalHotkeyService{
		configManager:    configManager,
		platformImpl:     platformImpl,
		isRunning:        false,
		ctx:              ctx,
		cancel:           cancel,
		windowController: windowController,
	}

	logger.Info("[全局快捷键] 服务初始化成功，平台: %s", runtime.GOOS)
	return service, nil
}

/**
 * Start 启动全局快捷键监听
 * @return error 错误信息
 */
func (ghs *GlobalHotkeyService) Start() error {
	ghs.mu.Lock()
	defer ghs.mu.Unlock()

	if ghs.isRunning {
		return fmt.Errorf("全局快捷键服务已在运行")
	}

	// 获取快捷键配置
	hotkeyConfig := ghs.configManager.GetHotkeyConfig()
	if !hotkeyConfig.EnableGlobalHotkey {
		logger.Info("[全局快捷键] 全局快捷键已禁用，跳过启动")
		return nil
	}

	// 注册显示/隐藏快捷键（需要全局开关和独立开关都启用）
	if hotkeyConfig.ShowHideHotkey != "" && hotkeyConfig.EnableShowHideHotkey {
		err := ghs.platformImpl.RegisterHotkey(hotkeyConfig.ShowHideHotkey, ghs.handleShowHideHotkey)
		if err != nil {
			return fmt.Errorf("注册显示/隐藏快捷键失败: %w", err)
		}
		logger.Info("[全局快捷键] 已注册显示/隐藏快捷键: %s", hotkeyConfig.ShowHideHotkey)
	} else if hotkeyConfig.ShowHideHotkey != "" && !hotkeyConfig.EnableShowHideHotkey {
		logger.Info("[全局快捷键] 显示/隐藏快捷键已禁用: %s", hotkeyConfig.ShowHideHotkey)
	}

	// 启动监听
	go func() {
		if err := ghs.platformImpl.StartListening(ghs.ctx); err != nil {
			logger.Error("[全局快捷键] 监听失败: %v", err)
		}
	}()

	ghs.isRunning = true
	logger.Info("[全局快捷键] 服务已启动")
	return nil
}

/**
 * Stop 停止全局快捷键监听
 * @return error 错误信息
 */
func (ghs *GlobalHotkeyService) Stop() error {
	ghs.mu.Lock()
	defer ghs.mu.Unlock()

	if !ghs.isRunning {
		return nil
	}

	// 停止监听
	if err := ghs.platformImpl.StopListening(); err != nil {
		logger.Error("[全局快捷键] 停止监听失败: %v", err)
	}

	// 取消上下文
	ghs.cancel()

	// 取消注册所有快捷键
	hotkeyConfig := ghs.configManager.GetHotkeyConfig()
	if hotkeyConfig.ShowHideHotkey != "" && hotkeyConfig.EnableShowHideHotkey {
		if err := ghs.platformImpl.UnregisterHotkey(hotkeyConfig.ShowHideHotkey); err != nil {
			logger.Error("[全局快捷键] 取消注册显示/隐藏快捷键失败: %v", err)
		}
	}

	ghs.isRunning = false
	logger.Info("[全局快捷键] 服务已停止")
	return nil
}

/**
 * Restart 重启全局快捷键监听（用于配置更新后）
 * @return error 错误信息
 */
func (ghs *GlobalHotkeyService) Restart() error {
	logger.Info("[全局快捷键] 重启服务...")

	if err := ghs.Stop(); err != nil {
		return fmt.Errorf("停止服务失败: %w", err)
	}

	// 重新创建上下文
	ghs.ctx, ghs.cancel = context.WithCancel(context.Background())

	if err := ghs.Start(); err != nil {
		return fmt.Errorf("启动服务失败: %w", err)
	}

	return nil
}

/**
 * IsRunning 检查服务是否正在运行
 * @return bool 是否运行中
 */
func (ghs *GlobalHotkeyService) IsRunning() bool {
	ghs.mu.RLock()
	defer ghs.mu.RUnlock()
	return ghs.isRunning
}

/**
 * UpdateConfig 更新快捷键配置
 * @param config 新的快捷键配置
 * @return error 错误信息
 */
func (ghs *GlobalHotkeyService) UpdateConfig(config models.HotkeyConfig) error {
	// 保存配置
	if err := ghs.configManager.SetHotkeyConfig(config); err != nil {
		return fmt.Errorf("保存快捷键配置失败: %w", err)
	}

	// 如果服务正在运行，重启以应用新配置
	if ghs.IsRunning() {
		if err := ghs.Restart(); err != nil {
			return fmt.Errorf("重启服务以应用新配置失败: %w", err)
		}
	}

	logger.Info("[全局快捷键] 配置已更新")
	return nil
}

/**
 * ValidateHotkey 验证快捷键格式
 * @param hotkey 快捷键字符串
 * @return error 错误信息
 */
func (ghs *GlobalHotkeyService) ValidateHotkey(hotkey string) error {
	if hotkey == "" {
		return fmt.Errorf("快捷键不能为空")
	}

	// 使用平台实现解析快捷键
	modifiers, key, err := ghs.platformImpl.ParseHotkey(hotkey)
	if err != nil {
		return fmt.Errorf("快捷键格式无效: %w", err)
	}

	if key == "" {
		return fmt.Errorf("快捷键必须包含主键")
	}

	if len(modifiers) == 0 {
		return fmt.Errorf("快捷键必须包含修饰键（如Ctrl、Alt、Shift等）")
	}

	logger.Debug("[全局快捷键] 快捷键验证通过: %s (修饰键: %v, 主键: %s)", hotkey, modifiers, key)
	return nil
}

/**
 * handleShowHideHotkey 处理显示/隐藏快捷键
 */
func (ghs *GlobalHotkeyService) handleShowHideHotkey() {
	logger.Info("[全局快捷键] 触发显示/隐藏快捷键")

	if ghs.windowController == nil {
		logger.Error("[全局快捷键] 窗口控制器未设置")
		return
	}

	// 切换窗口显示状态
	if err := ghs.windowController.ToggleWindow(); err != nil {
		logger.Error("[全局快捷键] 切换窗口状态失败: %v", err)
	}
}
