package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"time"

	"wepassword/internal/config"
	"wepassword/internal/crypto"
	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/models"
	"wepassword/internal/services"
	"wepassword/internal/utils"
	"wepassword/internal/version"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

/**
 * App 应用核心结构体
 * @author 陈凤庆
 * @description 管理应用的生命周期和核心服务
 */
type App struct {
	ctx                   context.Context
	configManager         *config.ConfigManager
	dbManager             *database.DatabaseManager
	vaultService          *services.VaultService
	accountService        *services.AccountService // 20251002 陈凤庆 passwordService改名为accountService
	groupService          *services.GroupService
	typeService           *services.TypeService           // 20251002 陈凤庆 tabService改名为typeService
	exportService         *services.ExportService         // 20251003 陈凤庆 导出服务
	importService         *services.ImportService         // 20251003 陈凤庆 导入服务
	lockService           *services.LockService           // 20251004 陈凤庆 锁定服务
	keyboardService       interface{}                     // 20251003 陈凤庆 平台特定的键盘服务
	keyboardHelperService *services.KeyboardHelperService // 20251005 陈凤庆 键盘助手服务
	remoteInputService    *services.RemoteInputService    // 20251005 陈凤庆 远程输入服务
	globalHotkeyService   *services.GlobalHotkeyService   // 20251014 陈凤庆 全局快捷键服务
	passwordRuleApp       *PasswordRuleApp                // 20251017 陈凤庆 密码规则应用服务
	usernameHistoryApp    *UsernameHistoryApp             // 20251017 陈凤庆 用户名历史记录应用服务
}

/**
 * NewApp 创建新的应用实例
 * @param configManager 配置管理器
 * @param dbManager 数据库管理器
 * @return *App 应用实例
 */
func NewApp(configManager *config.ConfigManager, dbManager *database.DatabaseManager) *App {
	return &App{
		configManager: configManager,
		dbManager:     dbManager,
	}
}

/**
 * Startup 应用启动时调用
 * @param ctx 应用上下文
 */
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	log.Printf("========================================")
	log.Printf("%s 启动中...", version.GetFullVersion())
	log.Printf("构建日期: %s", version.BuildDate)
	log.Printf("作者: %s", version.Author)
	log.Printf("========================================")

	// 初始化服务
	log.Println("[启动] 初始化服务...")
	a.vaultService = services.NewVaultService(a.dbManager, a.configManager)
	// 20251002 陈凤庆 passwordService改名为accountService，tabService改名为typeService
	a.accountService = services.NewAccountService(a.dbManager)
	a.groupService = services.NewGroupService(a.dbManager)
	a.typeService = services.NewTypeService(a.dbManager)
	// 20251003 陈凤庆 初始化导出导入服务
	a.exportService = services.NewExportService(a.dbManager, a.accountService, a.groupService, a.typeService)
	a.importService = services.NewImportService(a.dbManager, a.accountService, a.groupService, a.typeService, nil) // cryptoManager稍后设置
	// 20251004 陈凤庆 初始化锁定服务
	a.lockService = services.NewLockService(a.configManager)
	// 20251003 陈凤庆 平台特定的键盘服务初始化
	a.initKeyboardService()
	// 20251014 陈凤庆 初始化全局快捷键服务
	a.initGlobalHotkeyService()
	// 20251017 陈凤庆 初始化密码规则应用服务
	a.initPasswordRuleApp()
	// 20251017 陈凤庆 初始化用户名历史记录应用服务
	a.initUsernameHistoryApp()
	log.Println("[启动] 服务初始化完成")

	log.Printf("[启动] %s 启动完成\n", version.GetAppName())
}

/**
 * handleLock 处理锁定回调
 * @return error 锁定处理错误
 * @author 陈凤庆
 * @date 20251004
 * @description 当定时器触发锁定时调用，关闭密码库并清理状态
 */
func (a *App) handleLock() error {
	logger.Info("[锁定] 开始执行锁定操作")

	// 关闭密码库
	if a.vaultService != nil && a.vaultService.IsOpened() {
		a.vaultService.CloseVault()
		logger.Info("[锁定] 密码库已关闭")
	}

	// 清理加密管理器
	if a.accountService != nil {
		a.accountService.SetCryptoManager(nil)
		logger.Info("[锁定] 加密管理器已清理")
	}

	// 停止锁定服务（避免重复锁定）
	if a.lockService != nil {
		a.lockService.StopLockService()
		logger.Info("[锁定] 锁定服务已停止")
	}

	logger.Info("[锁定] 锁定操作完成")
	return nil
}

/**
 * GetLockConfig 获取锁定配置
 * @return models.LockConfig 锁定配置
 * @author 陈凤庆
 * @date 20251004
 */
func (a *App) GetLockConfig() models.LockConfig {
	return a.configManager.GetLockConfig()
}

/**
 * IsLockTriggered 检查是否触发了锁定
 * @return bool 是否触发锁定
 * @author 陈凤庆
 * @date 20251004
 */
func (a *App) IsLockTriggered() bool {
	if a.lockService != nil {
		return a.lockService.IsLockTriggered()
	}
	return false
}

/**
 * GetTimerStatus 获取定时器状态（调试用）
 * @return map[string]interface{} 定时器状态信息
 * @author 陈凤庆
 * @date 20251004
 */
func (a *App) GetTimerStatus() map[string]interface{} {
	if a.lockService != nil {
		return a.lockService.GetTimerStatus()
	}
	return map[string]interface{}{
		"error": "锁定服务未初始化",
	}
}

/**
 * SetLockConfig 设置锁定配置
 * @param lockConfig 锁定配置
 * @return error 设置错误
 * @author 陈凤庆
 * @date 20251004
 */
func (a *App) SetLockConfig(lockConfig models.LockConfig) error {
	logger.Info("[锁定配置] 开始更新锁定配置")

	// 保存到配置文件
	if err := a.configManager.SetLockConfig(lockConfig); err != nil {
		logger.Error("[锁定配置] 保存配置失败: %v", err)
		return fmt.Errorf("保存锁定配置失败: %w", err)
	}

	// 更新锁定服务配置
	if a.lockService != nil {
		// 转换为锁定服务的配置格式
		serviceConfig := &services.LockConfig{
			EnableAutoLock:     lockConfig.EnableAutoLock,
			EnableTimerLock:    lockConfig.EnableTimerLock,
			EnableMinimizeLock: lockConfig.EnableMinimizeLock,
			LockTimeMinutes:    lockConfig.LockTimeMinutes,
			EnableSystemLock:   lockConfig.EnableSystemLock,
			SystemLockMinutes:  lockConfig.SystemLockMinutes,
		}

		if err := a.lockService.UpdateLockConfig(serviceConfig); err != nil {
			logger.Error("[锁定配置] 更新锁定服务配置失败: %v", err)
			return fmt.Errorf("更新锁定服务配置失败: %w", err)
		}
	}

	logger.Info("[锁定配置] 锁定配置更新成功")
	return nil
}

/**
 * GetHotkeyConfig 获取快捷键配置
 * @return models.HotkeyConfig 快捷键配置
 * @author 陈凤庆
 * @date 20251014
 */
func (a *App) GetHotkeyConfig() models.HotkeyConfig {
	return a.configManager.GetHotkeyConfig()
}

/**
 * SetHotkeyConfig 设置快捷键配置
 * @param hotkeyConfig 快捷键配置
 * @return error 设置错误
 * @author 陈凤庆
 * @date 20251014
 */
func (a *App) SetHotkeyConfig(hotkeyConfig models.HotkeyConfig) error {
	logger.Info("[快捷键配置] 开始更新快捷键配置")

	// 验证快捷键格式
	if a.globalHotkeyService != nil && hotkeyConfig.ShowHideHotkey != "" {
		if err := a.globalHotkeyService.ValidateHotkey(hotkeyConfig.ShowHideHotkey); err != nil {
			logger.Error("[快捷键配置] 快捷键格式验证失败: %v", err)
			return fmt.Errorf("快捷键格式无效: %w", err)
		}
	}

	// 更新全局快捷键服务配置
	if a.globalHotkeyService != nil {
		if err := a.globalHotkeyService.UpdateConfig(hotkeyConfig); err != nil {
			logger.Error("[快捷键配置] 更新全局快捷键服务配置失败: %v", err)
			return fmt.Errorf("更新快捷键服务配置失败: %w", err)
		}
	} else {
		// 如果服务未初始化，只保存配置
		if err := a.configManager.SetHotkeyConfig(hotkeyConfig); err != nil {
			logger.Error("[快捷键配置] 保存配置失败: %v", err)
			return fmt.Errorf("保存快捷键配置失败: %w", err)
		}
	}

	logger.Info("[快捷键配置] 快捷键配置更新成功")
	return nil
}

/**
 * UpdateUserActivity 更新用户活动
 * @author 陈凤庆
 * @date 20251004
 * @description 当用户有操作时调用，重置定时锁定计时器
 */
func (a *App) UpdateUserActivity() {
	if a.lockService != nil {
		a.lockService.UpdateActivity()
	}
}

/**
 * TriggerLock 手动触发锁定
 * @return error 锁定错误
 * @author 陈凤庆
 * @date 20251004
 * @description 手动锁定密码库
 */
func (a *App) TriggerLock() error {
	logger.Info("[手动锁定] 用户手动触发锁定")
	return a.handleLock()
}

/**
 * OnWindowMinimize 窗口最小化事件处理
 * @author 陈凤庆
 * @date 20251004
 * @description 当窗口最小化时触发锁定检查
 */
func (a *App) OnWindowMinimize() {
	logger.Info("[窗口事件] 检测到窗口最小化")
	if a.lockService != nil {
		a.lockService.OnMinimize()
	}
}

/**
 * OnWindowFocus 窗口获得焦点事件处理
 * @author 陈凤庆
 * @date 20251004
 * @description 当窗口获得焦点时更新用户活动
 */
func (a *App) OnWindowFocus() {
	logger.Debug("[窗口事件] 窗口获得焦点")
	// 20251005 陈凤庆 添加防抖机制，避免频繁的用户活动更新
	if a.lockService != nil {
		// 只有在真正的用户操作时才更新活动时间，避免系统窗口切换触发
		go func() {
			time.Sleep(100 * time.Millisecond) // 短暂延迟，过滤掉快速的窗口切换
			a.lockService.UpdateActivity()
		}()
	}
}

/**
 * OnWindowBlur 窗口失去焦点事件处理
 * @author 陈凤庆
 * @date 20251004
 * @description 当窗口失去焦点时记录日志
 */
func (a *App) OnWindowBlur() {
	logger.Debug("[窗口事件] 窗口失去焦点")
}

/**
 * ShowWindow 显示窗口
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251014
 * @description 实现WindowController接口，显示应用窗口
 */
func (a *App) ShowWindow() error {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
		logger.Info("[窗口控制] 窗口已显示")
		return nil
	}
	return fmt.Errorf("应用上下文未初始化")
}

/**
 * HideWindow 隐藏窗口
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251014
 * @description 实现WindowController接口，隐藏应用窗口
 */
func (a *App) HideWindow() error {
	if a.ctx != nil {
		runtime.WindowHide(a.ctx)
		logger.Info("[窗口控制] 窗口已隐藏")
		return nil
	}
	return fmt.Errorf("应用上下文未初始化")
}

/**
 * IsWindowVisible 检查窗口是否可见
 * @return bool 是否可见
 * @author 陈凤庆
 * @date 20251014
 * @description 实现WindowController接口，检查窗口可见性
 */
func (a *App) IsWindowVisible() bool {
	if a.ctx != nil {
		// Wails runtime没有直接的IsVisible方法，这里使用IsMinimised的反向逻辑
		isMinimised := runtime.WindowIsMinimised(a.ctx)
		return !isMinimised
	}
	return false
}

/**
 * ToggleWindow 切换窗口显示状态
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251014
 * @description 实现WindowController接口，切换窗口显示/隐藏状态
 */
func (a *App) ToggleWindow() error {
	if a.ctx == nil {
		return fmt.Errorf("应用上下文未初始化")
	}

	// 获取窗口状态信息
	isMinimised := runtime.WindowIsMinimised(a.ctx)
	isMaximised := runtime.WindowIsMaximised(a.ctx)
	isFullscreen := runtime.WindowIsFullscreen(a.ctx)

	logger.Info("[窗口控制] 当前窗口状态 - 最小化: %v, 最大化: %v, 全屏: %v",
		isMinimised, isMaximised, isFullscreen)

	if isMinimised {
		// 窗口已最小化，恢复显示
		logger.Info("[窗口控制] 检测到窗口已最小化，正在恢复显示...")

		// 先显示窗口
		runtime.WindowShow(a.ctx)
		// 然后取消最小化
		runtime.WindowUnminimise(a.ctx)
		// 最后将窗口置于前台
		runtime.WindowSetAlwaysOnTop(a.ctx, true)
		runtime.WindowSetAlwaysOnTop(a.ctx, false) // 立即取消置顶，只是为了激活窗口

		logger.Info("[窗口控制] ✅ 窗口已从最小化状态恢复显示")
	} else {
		// 窗口可见，最小化
		logger.Info("[窗口控制] 检测到窗口可见，正在最小化...")

		runtime.WindowMinimise(a.ctx)

		logger.Info("[窗口控制] ✅ 窗口已最小化")
	}

	return nil
}

/**
 * TestToggleWindow 测试窗口切换功能
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251015
 * @description 手动测试窗口显示/隐藏切换功能
 */
func (a *App) TestToggleWindow() error {
	logger.Info("[测试] 手动触发窗口切换测试")
	return a.ToggleWindow()
}

/**
 * DomReady DOM 准备就绪时调用
 * @param ctx 应用上下文
 */
func (a *App) DomReady(ctx context.Context) {
	log.Println("[启动] DOM 准备就绪，前端界面已加载")
}

/**
 * BeforeClose 应用关闭前调用
 * @param ctx 应用上下文
 * @return bool 是否允许关闭
 */
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	log.Println("应用即将关闭")
	return false
}

/**
 * CheckAndRepairCorruptedAccounts 检查并修复损坏的账号数据
 * @return map[string]interface{} 修复结果
 * @author 20251004 陈凤庆 新增数据修复API
 */
func (a *App) CheckAndRepairCorruptedAccounts() map[string]interface{} {
	result := map[string]interface{}{
		"success":        false,
		"repaired_count": 0,
		"message":        "",
	}

	if a.accountService == nil {
		result["message"] = "账号服务未初始化"
		return result
	}

	repairedCount, err := a.accountService.CheckAndRepairCorruptedAccounts()
	if err != nil {
		result["message"] = fmt.Sprintf("检查修复失败: %v", err)
		return result
	}

	result["success"] = true
	result["repaired_count"] = repairedCount
	result["message"] = fmt.Sprintf("检查完成，修复了 %d 个损坏的账号", repairedCount)

	return result
}

/**
 * Shutdown 应用关闭时调用
 * @param ctx 应用上下文
 */
func (a *App) Shutdown(ctx context.Context) {
	log.Println("应用正在关闭...")

	// 20251003 陈凤庆 关闭密码库并清理状态，解决重启后需要输入两次密码的问题
	if a.vaultService != nil {
		log.Println("[关闭] 正在关闭密码库服务...")
		a.vaultService.CloseVault()
		log.Println("[关闭] 密码库服务已关闭")
	}

	// 20251003 陈凤庆 清理配置文件中的当前密码库路径，确保下次启动时状态干净
	if a.configManager != nil {
		log.Println("[关闭] 正在清理配置文件中的当前密码库路径...")
		if err := a.configManager.SetCurrentVaultPath(""); err != nil {
			log.Printf("[关闭] ⚠️ 清理配置文件失败: %v", err)
		} else {
			log.Println("[关闭] ✅ 配置文件已清理")
		}
	}

	// 20251014 陈凤庆 停止全局快捷键服务
	if a.globalHotkeyService != nil {
		log.Println("[关闭] 正在停止全局快捷键服务...")
		if err := a.globalHotkeyService.Stop(); err != nil {
			log.Printf("[关闭] ⚠️ 停止全局快捷键服务失败: %v", err)
		} else {
			log.Println("[关闭] 全局快捷键服务已停止")
		}
	}

	// 关闭数据库连接
	if a.dbManager != nil {
		log.Println("[关闭] 正在关闭数据库连接...")
		a.dbManager.Close()
		log.Println("[关闭] 数据库连接已关闭")
	}

	log.Println("应用已关闭")
}

// API 方法 - 供前端调用

/**
 * GetAppInfo 获取应用信息
 * @return map[string]interface{} 应用信息
 */
func (a *App) GetAppInfo() map[string]interface{} {
	log.Println("[API] 获取应用信息")
	return map[string]interface{}{
		"name":      version.GetAppName(),
		"version":   version.GetVersion(),
		"author":    version.Author,
		"buildDate": version.BuildDate,
		"time":      time.Now().Format("2006-01-02 15:04:05"),
	}
}

/**
 * GetChangeLog 获取更新日志
 * @return []version.ChangeLogEntry 更新日志列表
 */
func (a *App) GetChangeLog() []version.ChangeLogEntry {
	log.Println("[API] 获取更新日志")
	return version.GetChangeLog()
}

/**
 * GetLogConfig 获取日志配置
 * @return models.LogConfig 日志配置
 * @author 陈凤庆
 * @date 2025-10-03
 */
func (a *App) GetLogConfig() models.LogConfig {
	logger.Info("[API] 获取日志配置")
	return a.configManager.GetLogConfig()
}

/**
 * SetLogConfig 设置日志配置
 * @param config 日志配置
 * @return error 错误信息
 * @author 陈凤庆
 * @date 2025-10-03
 */
func (a *App) SetLogConfig(config models.LogConfig) error {
	logger.Info("[API] 设置日志配置: Info=%v, Debug=%v", config.EnableInfoLog, config.EnableDebugLog)

	// 保存到配置文件
	if err := a.configManager.SetLogConfig(config); err != nil {
		logger.Error("[API] 保存日志配置失败: %v", err)
		return err
	}

	// 更新日志记录器配置
	logger.SetLogConfig(logger.LogConfig{
		EnableInfoLog:  config.EnableInfoLog,
		EnableDebugLog: config.EnableDebugLog,
	})

	logger.Info("[API] 日志配置已更新并生效")
	return nil
}

/**
 * GetAppConfig 获取应用配置
 * @return models.AppConfig 应用配置
 * @author 陈凤庆
 * @date 2025-10-03
 */
func (a *App) GetAppConfig() models.AppConfig {
	logger.Info("[API] 获取应用配置")
	return *a.configManager.GetConfig()
}

/**
 * SetAppConfig 设置应用配置
 * @param config 应用配置的部分字段
 * @return error 错误信息
 * @author 陈凤庆
 * @date 2025-10-03
 */
func (a *App) SetAppConfig(config map[string]interface{}) error {
	logger.Info("[API] 设置应用配置: %+v", config)

	// 获取当前配置
	currentConfig := a.configManager.GetConfig()

	// 更新指定字段
	if theme, ok := config["theme"].(string); ok {
		if err := a.configManager.SetTheme(theme); err != nil {
			logger.Error("[API] 设置主题失败: %v", err)
			return err
		}
	}

	if language, ok := config["language"].(string); ok {
		if err := a.configManager.SetLanguage(language); err != nil {
			logger.Error("[API] 设置语言失败: %v", err)
			return err
		}
	}

	logger.Info("[API] 应用配置已更新: 主题=%s, 语言=%s", currentConfig.Theme, currentConfig.Language)
	return nil
}

/**
 * CheckVaultExists 检查密码库是否存在
 * @param vaultPath 密码库文件路径
 * @return bool 是否存在
 */
func (a *App) CheckVaultExists(vaultPath string) bool {
	return a.vaultService.CheckVaultExists(vaultPath)
}

/**
 * IsVaultOpened 检查密码库是否已打开
 * @return bool 是否已打开
 * @description 20251003 陈凤庆 用于前端检查登录状态，避免重复登录
 */
func (a *App) IsVaultOpened() bool {
	if a.vaultService == nil {
		log.Printf("[状态检查] VaultService为空，返回false")
		return false
	}

	// 20251003 陈凤庆 详细的状态检查，帮助调试登录两次密码问题
	isOpened := a.vaultService.IsOpened()
	dbOpened := a.dbManager != nil && a.dbManager.IsOpened()
	currentPath := a.vaultService.GetCurrentVaultPath()
	configPath := ""
	if a.configManager != nil {
		configPath = a.configManager.GetCurrentVaultPath()
	}

	log.Printf("[状态检查] VaultService.IsOpened=%t, DB.IsOpened=%t", isOpened, dbOpened)
	log.Printf("[状态检查] VaultService.CurrentPath=%s", currentPath)
	log.Printf("[状态检查] Config.CurrentPath=%s", configPath)

	// 20251003 陈凤庆 如果配置文件中有路径但VaultService状态为false，说明状态不一致
	if configPath != "" && !isOpened {
		log.Printf("[状态检查] ⚠️ 检测到状态不一致：配置文件有路径但密码库未打开")
		log.Printf("[状态检查] 这可能是导致需要输入两次密码的原因")

		// 清理不一致的配置状态
		if err := a.configManager.SetCurrentVaultPath(""); err != nil {
			log.Printf("[状态检查] ❌ 清理配置路径失败: %v", err)
		} else {
			log.Printf("[状态检查] ✅ 已清理不一致的配置路径")
		}
	}

	return isOpened
}

/**
 * CloseVault 关闭密码库
 * @description 20251003 陈凤庆 关闭密码库并清理状态，用于退出登录
 */
func (a *App) CloseVault() {
	logger.Info("[密码库] 开始关闭密码库")

	if a.vaultService != nil {
		a.vaultService.CloseVault()
		logger.Info("[密码库] 密码库服务已关闭")
	}

	// 清理配置文件中的当前密码库路径
	if a.configManager != nil {
		if err := a.configManager.SetCurrentVaultPath(""); err != nil {
			logger.Error("[密码库] 清理配置文件失败: %v", err)
		} else {
			logger.Info("[密码库] 配置文件已清理")
		}
	}

	logger.Info("[密码库] 密码库关闭完成")
}

/**
 * GetCurrentVaultPath 获取当前密码库路径
 * @return string 当前密码库路径
 * @description 20251017 陈凤庆 获取当前打开的密码库文件路径
 */
func (a *App) GetCurrentVaultPath() string {
	if a.vaultService == nil {
		return ""
	}
	return a.vaultService.GetCurrentVaultPath()
}

/**
 * CreateVault 创建新密码库
 * @param vaultName 密码库名称（无需后缀）
 * @param password 登录密码
 * @param language 语言代码（用于初始化多语言数据）
 * @return string 创建的密码库完整路径
 * @return error 错误信息
 * @modify 20251002 陈凤庆 使用跨平台路径工具函数，根据操作系统标准存储密码库文件
 * @modify 20251005 陈凤庆 添加语言参数，支持多语言初始数据
 */
func (a *App) CreateVault(vaultName string, password string, language string) (string, error) {
	log.Printf("[密码库] 开始创建密码库: %s", vaultName)

	// 20251002 陈凤庆 使用跨平台路径工具函数获取默认密码库路径
	vaultPath, err := utils.GetDefaultVaultPath(vaultName)
	if err != nil {
		log.Printf("[密码库] 获取密码库路径失败: %v", err)
		return "", fmt.Errorf("获取密码库路径失败: %w", err)
	}
	log.Printf("[密码库] 密码库路径: %s", vaultPath)

	err = a.vaultService.CreateVault(vaultPath, password, language)
	if err != nil {
		log.Printf("[密码库] 创建失败: %v", err)
		return "", err
	}
	log.Printf("[密码库] 密码库创建成功: %s", vaultPath)

	// 设置账号服务的加密管理器
	// 20251002 陈凤庆 passwordService改名为accountService
	a.accountService.SetCryptoManager(a.vaultService.GetCryptoManager())
	return vaultPath, nil
}

/**
 * OpenVault 打开密码库
 * @param vaultPath 密码库文件路径
 * @param password 登录密码
 * @return error 错误信息
 * @modify 20251001 陈凤庆 添加详细日志记录，确保加密管理器正确设置
 */
func (a *App) OpenVault(vaultPath string, password string) error {
	logger.Info("[密码库] 开始打开密码库: %s", vaultPath)

	// 20251001 陈凤庆 检查密码库服务是否已初始化
	if a.vaultService == nil {
		logger.Error("[密码库] 密码库服务未初始化")
		return fmt.Errorf("密码库服务未初始化")
	}

	err := a.vaultService.OpenVault(vaultPath, password)
	if err != nil {
		logger.Error("[密码库] 打开失败: %v", err)
		return err
	}
	logger.Info("[密码库] 密码库打开成功: %s", vaultPath)

	// 20251001 陈凤庆 设置密码服务的加密管理器
	cryptoManager := a.vaultService.GetCryptoManager()
	if cryptoManager == nil {
		logger.Error("[密码库] 获取加密管理器失败")
		return fmt.Errorf("获取加密管理器失败")
	}
	logger.Debug("[密码库] 加密管理器获取成功")

	// 20251002 陈凤庆 检查账号服务是否已初始化，passwordService改名为accountService
	if a.accountService == nil {
		logger.Error("[密码库] 账号服务未初始化")
		return fmt.Errorf("账号服务未初始化")
	}

	a.accountService.SetCryptoManager(cryptoManager)
	logger.Info("[密码库] 加密管理器设置完成")

	// 20251003 陈凤庆 设置导入服务的加密管理器
	if a.importService != nil {
		a.importService.SetCryptoManager(cryptoManager)
		logger.Info("[密码库] 导入服务加密管理器设置完成")
	}

	// 20251002 陈凤庆 验证加密管理器是否正确设置
	if !a.accountService.IsCryptoManagerSet() {
		logger.Error("[密码库] 加密管理器设置验证失败")
		return fmt.Errorf("加密管理器设置失败")
	}
	logger.Info("[密码库] 加密管理器设置验证成功")

	// 20251004 陈凤庆 启动锁定服务
	if a.lockService != nil {
		// 重置锁定触发标志（登录后重置）
		a.lockService.ResetLockTrigger()

		if err := a.lockService.StartLockService(); err != nil {
			logger.Error("[密码库] 启动锁定服务失败: %v", err)
		} else {
			logger.Info("[密码库] 锁定服务已启动")
		}
	}

	return nil
}

/**
 * GetRecentVaults 获取最近使用的密码库列表
 * @return []string 最近使用的密码库路径列表
 */
func (a *App) GetRecentVaults() []string {
	return a.configManager.GetRecentVaults()
}

/**
 * CheckRecentVaultStatus 检查最近使用的密码库状态
 * @return map[string]interface{} 包含存在的文件路径和简化模式状态
 * @author 陈凤庆
 * @description 20250928 检查最近使用的密码库文件是否存在，决定是否使用简化模式
 */
func (a *App) CheckRecentVaultStatus() map[string]interface{} {
	recentVaults := a.configManager.GetRecentVaults()
	result := map[string]interface{}{
		"hasValidVault": false,
		"vaultPath":     "",
		"isSimplified":  true,
	}

	// 检查最近使用的密码库文件是否存在
	for _, vaultPath := range recentVaults {
		if a.vaultService.CheckVaultExists(vaultPath) {
			result["hasValidVault"] = true
			result["vaultPath"] = vaultPath
			result["isSimplified"] = true // 有有效文件时使用简化模式
			break
		}
	}

	// 如果没有有效的最近使用文件，使用完整模式
	if !result["hasValidVault"].(bool) {
		result["isSimplified"] = false
	}

	return result
}

/**
 * GetGroups 获取所有分组
 * @return []models.Group 分组列表
 * @modify 20251002 陈凤庆 添加调试日志，排查分组ID问题
 */
func (a *App) GetGroups() ([]models.Group, error) {
	groups, err := a.groupService.GetAllGroups()
	if err != nil {
		log.Printf("[GetGroups] 获取分组失败: %v", err)
		return nil, err
	}

	log.Printf("[GetGroups] 获取到 %d 个分组:", len(groups))
	for i, group := range groups {
		// 20251001 陈凤庆 ID改为string类型，格式化符号改为%s
		log.Printf("[GetGroups] 分组 %d: ID=%s, Name='%s', Icon='%s'", i, group.ID, group.Name, group.Icon)
	}

	return groups, nil
}

/**
 * CreateGroup 创建新分组
 * @param name 分组名称
 * @return models.Group 创建的分组
 * @modify 20251002 陈凤庆 删除parentID参数，不需要层级结构
 */
func (a *App) CreateGroup(name string) (models.Group, error) {
	return a.groupService.CreateGroup(name)
}

/**
 * RenameGroup 重命名分组
 * @param id 分组ID
 * @param newName 新的分组名称
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组重命名API
 */
func (a *App) RenameGroup(id string, newName string) error {
	return a.groupService.RenameGroup(id, newName)
}

/**
 * DeleteGroup 删除分组
 * @param id 分组ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组删除API
 */
func (a *App) DeleteGroup(id string) error {
	return a.groupService.DeleteGroup(id)
}

/**
 * MoveGroupLeft 将分组向左移动一位
 * @param id 分组ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组左移API
 */
func (a *App) MoveGroupLeft(id string) error {
	return a.groupService.MoveGroupLeft(id)
}

/**
 * MoveGroupRight 将分组向右移动一位
 * @param id 分组ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组右移API
 */
func (a *App) MoveGroupRight(id string) error {
	return a.groupService.MoveGroupRight(id)
}

/**
 * UpdateGroupSortOrder 更新分组排序
 * @param groupID 分组ID
 * @param newSortOrder 新的排序号
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增分组排序更新API，用于拖拽排序
 */
func (a *App) UpdateGroupSortOrder(groupID string, newSortOrder int) error {
	return a.groupService.UpdateGroupSortOrder(groupID, newSortOrder)
}

/**
 * GetTypesByGroup 根据分组ID获取类型列表
 * @param groupID 分组ID（字符串格式）
 * @return []models.Type 类型列表
 * @modify 20251002 陈凤庆 GetTabsByGroup改名为GetTypesByGroup，Tab改为Type
 */
func (a *App) GetTypesByGroup(groupID string) ([]models.Type, error) {
	// 直接传递字符串GUID给服务层
	return a.typeService.GetTypesByGroup(groupID)
}

/**
 * GetAllTypes 获取所有类型
 * @return []models.Type 所有类型列表
 * @return error 错误信息
 * @author 20251004 陈凤庆 新增获取所有类型的API，用于导出功能
 */
func (a *App) GetAllTypes() ([]models.Type, error) {
	return a.typeService.GetAllTypes()
}

/**
 * CreateType 创建新类型
 * @param name 类型名称
 * @param groupID 所属分组ID
 * @param icon 图标
 * @return models.Type 创建的类型
 * @modify 20251002 陈凤庆 CreateTab改名为CreateType，Tab改为Type
 */
func (a *App) CreateType(name string, groupID string, icon string) (models.Type, error) {
	return a.typeService.CreateType(name, groupID, icon)
}

/**
 * UpdateType 更新类型
 * @param typeItem 类型信息
 * @return error 错误信息
 * @modify 20251002 陈凤庆 UpdateTab改名为UpdateType，Tab改为Type
 */
func (a *App) UpdateType(typeItem models.Type) error {
	return a.typeService.UpdateType(typeItem)
}

/**
 * DeleteType 删除类型
 * @param id 类型ID
 * @return error 错误信息
 * @modify 20251002 陈凤庆 DeleteTab改名为DeleteType，Tab改为Type
 */
func (a *App) DeleteType(id string) error {
	return a.typeService.DeleteType(id)
}

/**
 * MoveTypeUp 将类型向上移动一位
 * @param id 类型ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增类型上移API
 */
func (a *App) MoveTypeUp(id string) error {
	return a.typeService.MoveTypeUp(id)
}

/**
 * MoveTypeDown 将类型向下移动一位
 * @param id 类型ID
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增类型下移API
 */
func (a *App) MoveTypeDown(id string) error {
	return a.typeService.MoveTypeDown(id)
}

/**
 * UpdateTypeSortOrder 更新类型排序
 * @param typeID 类型ID
 * @param newSortOrder 新的排序号
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增类型排序更新API，用于拖拽排序
 */
func (a *App) UpdateTypeSortOrder(typeID string, newSortOrder int) error {
	return a.typeService.UpdateTypeSortOrder(typeID, newSortOrder)
}

/**
 * InsertTypeAfter 在指定标签后插入新标签
 * @param name 新标签名称
 * @param groupID 分组ID
 * @param icon 图标
 * @param afterTypeID 在此标签后插入，如果为空则插入到最后
 * @return models.Type 创建的类型
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增在指定标签后插入新标签API
 */
func (a *App) InsertTypeAfter(name string, groupID string, icon string, afterTypeID string) (models.Type, error) {
	return a.typeService.InsertTypeAfter(name, groupID, icon, afterTypeID)
}

/**
 * GetAccountsByGroup 根据分组获取账号
 * @param groupID 分组ID（字符串格式）
 * @return []models.AccountDecrypted 账号列表
 * @modify 20251002 陈凤庆 GetPasswordsByGroup改名为GetAccountsByGroup，账号改为账号
 */
func (a *App) GetAccountsByGroup(groupID string) ([]models.AccountDecrypted, error) {
	logger.LogAPICall("GetAccountsByGroup", fmt.Sprintf("groupID=%s", groupID), "开始处理")

	// 20251002 陈凤庆 检查账号服务是否已初始化，passwordService改名为accountService
	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，分组ID: %s", groupID)
		return nil, fmt.Errorf("账号服务未初始化")
	}

	// 直接传递字符串GUID给服务层
	result, err := a.accountService.GetAccountsByGroup(groupID)
	if err != nil {
		logger.Error("[API] GetAccountsByGroup失败，分组ID: %s, 错误: %v", groupID, err)
		logger.LogAPICall("GetAccountsByGroup", fmt.Sprintf("groupID=%s", groupID), fmt.Sprintf("失败: %v", err))
		return nil, err
	}

	logger.LogAPICall("GetAccountsByGroup", fmt.Sprintf("groupID=%s", groupID), fmt.Sprintf("成功，返回%d个账号", len(result)))
	return result, nil
}

/**
 * GetAccountsByTab 根据标签获取账号
 * @param tabID 标签ID（字符串格式）
 * @return []models.AccountDecrypted 账号列表
 * @author 20251002 陈凤庆 新增按标签ID查询账号的API
 */
func (a *App) GetAccountsByTab(tabID string) ([]models.AccountDecrypted, error) {
	logger.LogAPICall("GetAccountsByTab", fmt.Sprintf("tabID=%s", tabID), "开始处理")

	// 20251002 陈凤庆 检查账号服务是否已初始化
	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，标签ID: %s", tabID)
		return nil, fmt.Errorf("账号服务未初始化")
	}

	// 直接传递字符串GUID给服务层
	result, err := a.accountService.GetAccountsByTab(tabID)
	if err != nil {
		logger.Error("[API] GetAccountsByTab失败，标签ID: %s, 错误: %v", tabID, err)
		logger.LogAPICall("GetAccountsByTab", fmt.Sprintf("tabID=%s", tabID), fmt.Sprintf("失败: %v", err))
		return nil, err
	}

	logger.LogAPICall("GetAccountsByTab", fmt.Sprintf("tabID=%s", tabID), fmt.Sprintf("成功，返回%d个账号", len(result)))
	return result, nil
}

/**
 * GetAccountsByConditions 根据查询条件获取账号列表
 * @param conditions 查询条件JSON字符串，格式：{"group_id":"xxx","type_id":"xxx"}
 * @return []models.AccountDecrypted 账号列表
 * @author 20251003 陈凤庆 新增统一的账号查询API
 */
func (a *App) GetAccountsByConditions(conditions string) ([]models.AccountDecrypted, error) {
	logger.LogAPICall("GetAccountsByConditions", fmt.Sprintf("conditions=%s", conditions), "开始处理")

	// 检查账号服务是否已初始化
	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，查询条件: %s", conditions)
		return nil, fmt.Errorf("账号服务未初始化")
	}

	result, err := a.accountService.GetAccountsByConditions(conditions)
	if err != nil {
		logger.Error("[API] GetAccountsByConditions失败，查询条件: %s, 错误: %v", conditions, err)
		logger.LogAPICall("GetAccountsByConditions", fmt.Sprintf("conditions=%s", conditions), fmt.Sprintf("失败: %v", err))
		return nil, err
	}

	logger.LogAPICall("GetAccountsByConditions", fmt.Sprintf("conditions=%s", conditions), fmt.Sprintf("成功，返回%d个账号", len(result)))
	return result, nil
}

/**
 * CreateAccount 创建新账号
 * @param title 标题
 * @param username 用户名
 * @param password 密码
 * @param url 网址
 * @param typeID 类型ID
 * @param notes 备注
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
 * @return models.AccountDecrypted 创建的账号
 * @modify 20251002 陈凤庆 CreatePasswordItem改名为CreateAccount，账号改为账号
 * @modify 20251003 陈凤庆 添加inputMethod参数和详细日志记录
 */
func (a *App) CreateAccount(title, username, password, url, typeID, notes string, inputMethod int) (models.AccountDecrypted, error) {
	logger.LogAPICall("CreateAccount", fmt.Sprintf("title=%s, typeID=%s, inputMethod=%d", title, typeID, inputMethod), "开始处理")

	result, err := a.accountService.CreateAccount(title, username, password, url, typeID, notes, inputMethod)
	if err != nil {
		logger.Error("[API] CreateAccount失败，标题: %s, 错误: %v", title, err)
		logger.LogAPICall("CreateAccount", fmt.Sprintf("title=%s", title), fmt.Sprintf("失败: %v", err))
		return models.AccountDecrypted{}, err
	}

	logger.LogAPICall("CreateAccount", fmt.Sprintf("title=%s, accountID=%s", title, result.ID), "成功")
	return result, nil
}

/**
 * GetAccountCredentials 根据账号ID获取账号的用户名和密码（用于复制和输入操作）
 * @param accountID 账号ID
 * @return username 用户名
 * @return password 密码
 * @return inputMethod 输入方式
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 * @description 为安全考虑，只在需要时查询敏感信息，查询后立即返回，不在内存中长期保存
 * @modify 20251004 陈凤庆 添加数据验证和错误处理，防止返回损坏的数据
 */
func (a *App) GetAccountCredentials(accountID string) (string, string, int, error) {
	// 验证账号ID
	if accountID == "" {
		return "", "", 0, fmt.Errorf("账号ID不能为空")
	}

	logger.Info("[获取凭据] 开始获取账号凭据，账号ID: %s", accountID)

	account, err := a.accountService.GetAccountByID(accountID)
	if err != nil {
		logger.Error("[获取凭据] 获取账号失败，账号ID: %s, 错误: %v", accountID, err)
		return "", "", 0, fmt.Errorf("获取账号失败: %w", err)
	}

	if account == nil {
		logger.Error("[获取凭据] 账号不存在，账号ID: %s", accountID)
		return "", "", 0, fmt.Errorf("账号不存在")
	}

	// 验证用户名和密码
	if account.Username == "" {
		logger.Error("[获取凭据] 用户名为空，账号ID: %s", accountID)
		return "", "", 0, fmt.Errorf("用户名为空")
	}

	if account.Password == "" {
		logger.Error("[获取凭据] 密码为空，账号ID: %s", accountID)
		return "", "", 0, fmt.Errorf("密码为空")
	}

	// 验证输入方式
	if account.InputMethod < 1 || account.InputMethod > 5 {
		logger.Info("[获取凭据] 输入方式无效: %d，重置为默认值1，账号ID: %s", account.InputMethod, accountID)
		account.InputMethod = 1 // 默认使用Unicode方式
	}

	// 检查用户名和密码是否包含异常重复字符
	if len(account.Username) > 2 {
		allSame := true
		firstChar := account.Username[0]
		for _, char := range account.Username {
			if char != rune(firstChar) {
				allSame = false
				break
			}
		}
		if allSame {
			logger.Error("[获取凭据] 检测到用户名异常重复字符: '%s'，账号ID: %s", account.Username, accountID)
			return "", "", 0, fmt.Errorf("用户名数据异常，包含重复字符: %s", account.Username)
		}
	}

	if len(account.Password) > 2 {
		allSame := true
		firstChar := account.Password[0]
		for _, char := range account.Password {
			if char != rune(firstChar) {
				allSame = false
				break
			}
		}
		if allSame {
			logger.Error("[获取凭据] 检测到密码异常重复字符，账号ID: %s", accountID)
			return "", "", 0, fmt.Errorf("密码数据异常，包含重复字符")
		}
	}

	logger.Info("[获取凭据] 成功获取账号凭据，账号ID: %s, 用户名长度: %d, 密码长度: %d, 输入方式: %d",
		accountID, len(account.Username), len(account.Password), account.InputMethod)

	return account.Username, account.Password, account.InputMethod, nil
}

/**
 * GetAccountByID 根据账号ID获取完整账号数据（用于编辑）
 * @param accountID 账号ID
 * @return *models.AccountDecrypted 完整账号数据（包括解密后的所有字段）
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 * @description 返回解密后的完整账号数据，包括用户名、密码等所有字段，用于编辑功能
 */
func (a *App) GetAccountByID(accountID string) (*models.AccountDecrypted, error) {
	logger.LogAPICall("GetAccountByID", fmt.Sprintf("accountID=%s", accountID), "开始处理")

	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，账号ID: %s", accountID)
		return nil, fmt.Errorf("账号服务未初始化")
	}

	result, err := a.accountService.GetAccountByID(accountID)
	if err != nil {
		logger.Error("[API] GetAccountByID失败，账号ID: %s, 错误: %v", accountID, err)
		logger.LogAPICall("GetAccountByID", fmt.Sprintf("accountID=%s", accountID), fmt.Sprintf("失败: %v", err))
		return nil, err
	}

	logger.LogAPICall("GetAccountByID", fmt.Sprintf("accountID=%s", accountID), "成功")
	return result, nil
}

/**
 * GetAccountDetail 根据账号ID获取账号详情（用于详情页面显示）
 * @param accountID 账号ID
 * @return models.AccountDecrypted 账号详情
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 * @description 返回解密后的用户名，密码不返回，备注返回脱敏版本
 */
func (a *App) GetAccountDetail(accountID string) (*models.AccountDecrypted, error) {
	logger.LogAPICall("GetAccountDetail", fmt.Sprintf("accountID=%s", accountID), "开始处理")

	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，账号ID: %s", accountID)
		return nil, fmt.Errorf("账号服务未初始化")
	}

	result, err := a.accountService.GetAccountDetail(accountID)
	if err != nil {
		logger.Error("[API] GetAccountDetail失败，账号ID: %s, 错误: %v", accountID, err)
		logger.LogAPICall("GetAccountDetail", fmt.Sprintf("accountID=%s", accountID), fmt.Sprintf("失败: %v", err))
		return nil, err
	}

	logger.LogAPICall("GetAccountDetail", fmt.Sprintf("accountID=%s", accountID), "成功")
	return result, nil
}

/**
 * GetAccountRaw 根据账号ID获取原始加密数据（用于特殊场景）
 * @param accountID 账号ID
 * @return *models.Account 数据库中的加密账号数据
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 * @description 返回数据库中的原始加密数据，不进行解密处理
 */
func (a *App) GetAccountRaw(accountID string) (*models.Account, error) {
	logger.LogAPICall("GetAccountRaw", fmt.Sprintf("accountID=%s", accountID), "开始处理")

	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，账号ID: %s", accountID)
		return nil, fmt.Errorf("账号服务未初始化")
	}

	result, err := a.accountService.GetAccountRaw(accountID)
	if err != nil {
		logger.Error("[API] GetAccountRaw失败，账号ID: %s, 错误: %v", accountID, err)
		logger.LogAPICall("GetAccountRaw", fmt.Sprintf("accountID=%s", accountID), fmt.Sprintf("失败: %v", err))
		return nil, err
	}

	logger.LogAPICall("GetAccountRaw", fmt.Sprintf("accountID=%s", accountID), "成功")
	return result, nil
}

/**
 * GetAccountPassword 根据账号ID获取解密后的密码（用于复制和显示操作）
 * @param accountID 账号ID
 * @return string 解密后的密码
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 * @description 为安全考虑，只在需要时查询密码，查询后立即返回，不在内存中长期保存
 */
func (a *App) GetAccountPassword(accountID string) (string, error) {
	logger.LogAPICall("GetAccountPassword", fmt.Sprintf("accountID=%s", accountID), "开始处理")

	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，账号ID: %s", accountID)
		return "", fmt.Errorf("账号服务未初始化")
	}

	account, err := a.accountService.GetAccountByID(accountID)
	if err != nil {
		logger.Error("[API] GetAccountPassword失败，账号ID: %s, 错误: %v", accountID, err)
		logger.LogAPICall("GetAccountPassword", fmt.Sprintf("accountID=%s", accountID), fmt.Sprintf("失败: %v", err))
		return "", err
	}

	logger.LogAPICall("GetAccountPassword", fmt.Sprintf("accountID=%s", accountID), "成功")
	return account.Password, nil
}

/**
 * CopyAccountUsername 复制账号用户名到剪贴板（10秒后自动清理）
 * @param accountID 账号ID
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 */
func (a *App) CopyAccountUsername(accountID string) error {
	username, _, _, err := a.GetAccountCredentials(accountID)
	if err != nil {
		return err
	}
	return a.copyToClipboardWithTimeout(username, "用户名")
}

/**
 * CopyAccountPassword 复制账号密码到剪贴板（10秒后自动清理）
 * @param accountID 账号ID
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 */
func (a *App) CopyAccountPassword(accountID string) error {
	_, password, _, err := a.GetAccountCredentials(accountID)
	if err != nil {
		return err
	}
	return a.copyToClipboardWithTimeout(password, "密码")
}

/**
 * CopyAccountUsernameAndPassword 复制账号用户名和密码到剪贴板（10秒后自动清理）
 * @param accountID 账号ID
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 */
func (a *App) CopyAccountUsernameAndPassword(accountID string) error {
	username, password, _, err := a.GetAccountCredentials(accountID)
	if err != nil {
		return err
	}
	combined := username + "\t" + password // 用Tab分隔用户名和密码
	return a.copyToClipboardWithTimeout(combined, "用户名和密码")
}

/**
 * CopyAccountNotes 复制账号备注到剪贴板（10秒后自动清理）
 * @param accountID 账号ID
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 */
func (a *App) CopyAccountNotes(accountID string) error {
	account, err := a.accountService.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	return a.copyToClipboardWithTimeout(account.Notes, "备注")
}

/**
 * GetAccountNotes 根据账号ID获取解密后的备注（用于显示操作）
 * @param accountID 账号ID
 * @return string 解密后的备注
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 * @description 为安全考虑，只在需要时查询备注，查询后立即返回，不在内存中长期保存
 */
func (a *App) GetAccountNotes(accountID string) (string, error) {
	logger.LogAPICall("GetAccountNotes", fmt.Sprintf("accountID=%s", accountID), "开始处理")

	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，账号ID: %s", accountID)
		return "", fmt.Errorf("账号服务未初始化")
	}

	account, err := a.accountService.GetAccountByID(accountID)
	if err != nil {
		logger.Error("[API] GetAccountNotes失败，账号ID: %s, 错误: %v", accountID, err)
		logger.LogAPICall("GetAccountNotes", fmt.Sprintf("accountID=%s", accountID), fmt.Sprintf("失败: %v", err))
		return "", err
	}

	logger.LogAPICall("GetAccountNotes", fmt.Sprintf("accountID=%s", accountID), "成功")
	return account.Notes, nil
}

/**
 * copyToClipboardWithTimeout 复制内容到剪贴板，并在指定时间后自动清理
 * @param content 要复制的内容
 * @param description 内容描述（用于日志）
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 */
func (a *App) copyToClipboardWithTimeout(content, description string) error {
	// 使用安全的robotgo复制到剪贴板
	err := a.safeWriteToClipboard(content)
	if err != nil {
		logger.Error("[剪贴板] 复制失败: %v", err)
		return err
	}

	logger.Info("[剪贴板] %s已复制到剪贴板，10秒后自动清理", description)

	// 启动协程，10秒后清理剪贴板
	go func() {
		time.Sleep(10 * time.Second)
		// 清理剪贴板，写入空字符串
		err := a.safeWriteToClipboard("")
		if err != nil {
			logger.Error("[剪贴板] 清理失败: %v", err)
		} else {
			logger.Info("[剪贴板] %s已从剪贴板清理", description)
		}
	}()

	return nil
}

/**
 * safeWriteToClipboard 安全的剪贴板写入操作
 * @param content 要写入的内容
 * @return error 错误信息
 * @description 20251004 陈凤庆 安全的剪贴板写入，添加异常处理，防止程序崩溃
 */
func (a *App) safeWriteToClipboard(content string) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[剪贴板] WriteAll操作发生异常: %v", r)
		}
	}()

	return robotgo.WriteAll(content)
}

/**
 * UpdateAccount 更新账号
 * @param account 账号信息（解密后）
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增更新账号的App层方法
 */
func (a *App) UpdateAccount(account models.AccountDecrypted) error {
	logger.Info("[API] UpdateAccount - 参数: accountID=%s, title=%s - 结果: 开始处理", account.ID, account.Title)
	err := a.accountService.UpdateAccount(account)
	if err != nil {
		logger.Error("[API] UpdateAccount - 参数: accountID=%s - 结果: 失败，错误: %v", account.ID, err)
		return err
	}
	logger.Info("[API] UpdateAccount - 参数: accountID=%s - 结果: 成功", account.ID)
	return nil
}

/**
 * DeleteAccount 删除账号
 * @param accountID 账号ID
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增删除账号的App层方法
 */
func (a *App) DeleteAccount(accountID string) error {
	logger.LogAPICall("DeleteAccount", fmt.Sprintf("accountID=%s", accountID), "开始处理")

	if a.accountService == nil {
		logger.Error("[API] 账号服务未初始化，账号ID: %s", accountID)
		return fmt.Errorf("账号服务未初始化")
	}

	err := a.accountService.DeleteAccount(accountID)
	if err != nil {
		logger.Error("[API] DeleteAccount失败，账号ID: %s, 错误: %v", accountID, err)
		logger.LogAPICall("DeleteAccount", fmt.Sprintf("accountID=%s", accountID), fmt.Sprintf("失败: %v", err))
		return err
	}

	logger.LogAPICall("DeleteAccount", fmt.Sprintf("accountID=%s", accountID), "成功")
	return nil
}

/**
 * SearchAccounts 搜索账号
 * @param keyword 搜索关键词
 * @return []models.AccountDecrypted 搜索结果
 * @modify 20251002 陈凤庆 SearchPasswords改名为SearchAccounts，账号改为账号
 */
func (a *App) SearchAccounts(keyword string) ([]models.AccountDecrypted, error) {
	return a.accountService.SearchAccounts(keyword)
}

/**
 * UpdateAccountUsage 更新账号使用次数
 * @param accountId 账号ID
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增使用次数统计功能
 */
func (a *App) UpdateAccountUsage(accountId string) error {
	logger.Info("[API] UpdateAccountUsage - 参数: accountId=%s - 结果: 开始处理", accountId)
	err := a.accountService.UpdateAccountUsage(accountId)
	if err != nil {
		logger.Error("[API] UpdateAccountUsage - 参数: accountId=%s - 结果: 失败 - 错误: %v", accountId, err)
		return err
	}
	logger.Info("[API] UpdateAccountUsage - 参数: accountId=%s - 结果: 成功", accountId)
	return nil
}

/**
 * UpdateAccountGroup 更新账号分组
 * @param accountId 账号ID
 * @param typeId 新的类型ID
 * @return error 错误信息
 * @author 20251005 陈凤庆 新增更新账号分组的App层方法
 */
func (a *App) UpdateAccountGroup(accountId string, typeId string) error {
	logger.Info("[API] UpdateAccountGroup - 参数: accountId=%s, typeId=%s - 结果: 开始处理", accountId, typeId)
	err := a.accountService.UpdateAccountGroup(accountId, typeId)
	if err != nil {
		logger.Error("[API] UpdateAccountGroup - 参数: accountId=%s, typeId=%s - 结果: 失败 - 错误: %v", accountId, typeId, err)
		return err
	}
	logger.Info("[API] UpdateAccountGroup - 参数: accountId=%s, typeId=%s - 结果: 成功", accountId, typeId)
	return nil
}

/**
 * GetAllAccounts 获取所有账号
 * @return []models.AccountDecrypted 所有账号列表
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增使用次数统计功能
 */
func (a *App) GetAllAccounts() ([]models.AccountDecrypted, error) {
	logger.Info("[API] GetAllAccounts - 结果: 开始处理")
	accounts, err := a.accountService.GetAllAccounts()
	if err != nil {
		logger.Error("[API] GetAllAccounts - 结果: 失败 - 错误: %v", err)
		return nil, err
	}
	logger.Info("[API] GetAllAccounts - 结果: 成功 - 账号数量: %d", len(accounts))
	return accounts, nil
}

/**
 * UpdateAppUsage 更新应用使用统计
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增使用天数统计功能
 */
func (a *App) UpdateAppUsage() error {
	logger.Info("[API] UpdateAppUsage - 结果: 开始处理")

	// 获取当前时间
	now := time.Now()

	// 创建SysInfoManager
	sysInfoMgr := database.NewSysInfoManager(a.dbManager.GetDB())

	// 获取首次打开时间
	firstOpenTime, err := sysInfoMgr.GetValue("first_open_time")
	if err != nil || firstOpenTime == "" {
		// 首次打开，记录首次打开时间
		err = sysInfoMgr.SetValue("first_open_time", now.Format(time.RFC3339))
		if err != nil {
			logger.Error("[API] UpdateAppUsage - 设置首次打开时间失败: %v", err)
			return err
		}
		logger.Info("[API] UpdateAppUsage - 首次打开，记录首次打开时间: %s", now.Format(time.RFC3339))
	}

	// 更新最近打开时间
	err = sysInfoMgr.SetValue("last_open_time", now.Format(time.RFC3339))
	if err != nil {
		logger.Error("[API] UpdateAppUsage - 设置最近打开时间失败: %v", err)
		return err
	}

	logger.Info("[API] UpdateAppUsage - 结果: 成功")
	return nil
}

/**
 * GetUsageDays 获取使用天数
 * @return int 使用天数
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增使用天数统计功能
 */
func (a *App) GetUsageDays() (int, error) {
	logger.Info("[API] GetUsageDays - 结果: 开始处理")

	// 创建SysInfoManager
	sysInfoMgr := database.NewSysInfoManager(a.dbManager.GetDB())

	// 获取首次打开时间
	firstOpenTimeStr, err := sysInfoMgr.GetValue("first_open_time")
	if err != nil || firstOpenTimeStr == "" {
		logger.Info("[API] GetUsageDays - 首次打开时间不存在，返回0天")
		return 0, nil
	}

	// 解析首次打开时间
	firstOpenTime, err := time.Parse(time.RFC3339, firstOpenTimeStr)
	if err != nil {
		logger.Error("[API] GetUsageDays - 解析首次打开时间失败: %v", err)
		return 0, err
	}

	// 计算使用天数
	now := time.Now()
	duration := now.Sub(firstOpenTime)
	days := int(duration.Hours()/24) + 1 // 加1是因为当天也算一天

	logger.Info("[API] GetUsageDays - 结果: 成功 - 使用天数: %d", days)
	return days, nil
}

/**
 * SelectVaultFile 选择密码库文件
 * @return string 选择的文件路径，如果取消选择则返回空字符串
 */
func (a *App) SelectVaultFile() string {
	// 20250127 陈凤庆 添加文件选择对话框功能，解决前端无法选择文件的问题
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择密码库文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "密码库文件 (*.db)",
				Pattern:     "*.db",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		log.Printf("文件选择对话框错误: %v", err)
		return ""
	}

	return selection
}

/**
 * OpenVaultDirectory 打开密码库所在目录
 * @return error 错误信息
 * @description 20251002 陈凤庆 跨平台打开密码库文件所在的目录
 */
func (a *App) OpenVaultDirectory() error {
	logger.Info("[密码库] 准备打开密码库所在目录")

	// 20251002 陈凤庆 检查密码库服务是否已初始化
	if a.vaultService == nil {
		logger.Error("[密码库] 密码库服务未初始化")
		return fmt.Errorf("密码库服务未初始化")
	}

	// 20251002 陈凤庆 获取当前密码库路径
	vaultPath := a.configManager.GetCurrentVaultPath()
	if vaultPath == "" {
		logger.Error("[密码库] 当前没有打开的密码库")
		return fmt.Errorf("当前没有打开的密码库")
	}

	// 20251002 陈凤庆 检查文件是否存在
	if !a.vaultService.CheckVaultExists(vaultPath) {
		logger.Error("[密码库] 密码库文件不存在: %s", vaultPath)
		return fmt.Errorf("密码库文件不存在: %s", vaultPath)
	}

	// 20251002 陈凤庆 获取文件所在目录
	dirPath := filepath.Dir(vaultPath)
	logger.Info("[密码库] 密码库所在目录: %s", dirPath)

	// 20251002 陈凤庆 根据不同操作系统打开文件夹
	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "windows":
		// Windows: 使用 explorer 并选中文件
		cmd = exec.Command("explorer", "/select,", vaultPath)
		logger.Info("[密码库] 使用 Windows Explorer 打开目录")
	case "darwin":
		// macOS: 使用 open 命令并选中文件
		cmd = exec.Command("open", "-R", vaultPath)
		logger.Info("[密码库] 使用 macOS Finder 打开目录")
	case "linux":
		// Linux: 使用 xdg-open 打开目录（不能选中文件）
		cmd = exec.Command("xdg-open", dirPath)
		logger.Info("[密码库] 使用 xdg-open 打开目录")
	default:
		logger.Error("[密码库] 不支持的操作系统: %s", goruntime.GOOS)
		return fmt.Errorf("不支持的操作系统: %s", goruntime.GOOS)
	}

	// 20251002 陈凤庆 执行命令
	if err := cmd.Start(); err != nil {
		logger.Error("[密码库] 打开目录失败: %v", err)
		return fmt.Errorf("打开目录失败: %v", err)
	}

	logger.Info("[密码库] 成功打开密码库所在目录")
	return nil
}

/**
 * VerifyOldPassword 验证旧登录密码
 * @param oldPassword 旧登录密码
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增验证旧登录密码功能
 * @description 仅验证旧密码是否正确，不进行任何修改操作
 */
func (a *App) VerifyOldPassword(oldPassword string) error {
	logger.Info("[验证密码] 开始验证旧登录密码")

	// 检查密码库服务是否已初始化
	if a.vaultService == nil {
		logger.Error("[验证密码] 密码库服务未初始化")
		return fmt.Errorf("密码库服务未初始化")
	}

	// 验证旧密码
	logger.Info("[验证密码] 正在验证旧密码...")
	vaultConfig, err := a.dbManager.GetVaultConfig()
	if err != nil {
		logger.Error("[验证密码] 获取密码库配置失败: %v", err)
		return fmt.Errorf("获取密码库配置失败: %w", err)
	}

	if vaultConfig == nil {
		logger.Error("[验证密码] 密码库配置不存在")
		return fmt.Errorf("密码库配置不存在")
	}

	// 解码盐值
	salt, err := base64.StdEncoding.DecodeString(vaultConfig.Salt)
	if err != nil {
		logger.Error("[验证密码] 解码盐值失败: %v", err)
		return fmt.Errorf("解码盐值失败: %w", err)
	}

	// 验证旧密码
	cryptoManager := a.vaultService.GetCryptoManager()
	if !cryptoManager.VerifyPassword(oldPassword, vaultConfig.PasswordHash, salt) {
		logger.Error("[验证密码] 旧密码验证失败")
		return fmt.Errorf("旧密码不正确")
	}

	logger.Info("[验证密码] ✅ 旧密码验证成功")
	return nil
}

/**
 * ChangeLoginPassword 修改登录密码
 * @param oldPassword 旧登录密码
 * @param newPassword 新登录密码
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增修改登录密码功能
 * @description 验证旧密码，用新密码重新加密所有账号数据，更新密码库配置
 */
func (a *App) ChangeLoginPassword(oldPassword, newPassword string) error {
	logger.Info("[修改密码] 开始修改登录密码")

	// 检查密码库服务是否已初始化
	if a.vaultService == nil {
		logger.Error("[修改密码] 密码库服务未初始化")
		return fmt.Errorf("密码库服务未初始化")
	}

	// 检查账号服务是否已初始化
	if a.accountService == nil {
		logger.Error("[修改密码] 账号服务未初始化")
		return fmt.Errorf("账号服务未初始化")
	}

	// 验证旧密码
	logger.Info("[修改密码] 正在验证旧密码...")
	vaultConfig, err := a.dbManager.GetVaultConfig()
	if err != nil {
		logger.Error("[修改密码] 获取密码库配置失败: %v", err)
		return fmt.Errorf("获取密码库配置失败: %w", err)
	}

	if vaultConfig == nil {
		logger.Error("[修改密码] 密码库配置不存在")
		return fmt.Errorf("密码库配置不存在")
	}

	// 解码盐值
	salt, err := base64.StdEncoding.DecodeString(vaultConfig.Salt)
	if err != nil {
		logger.Error("[修改密码] 解码盐值失败: %v", err)
		return fmt.Errorf("解码盐值失败: %w", err)
	}

	// 验证旧密码
	cryptoManager := a.vaultService.GetCryptoManager()
	if !cryptoManager.VerifyPassword(oldPassword, vaultConfig.PasswordHash, salt) {
		logger.Error("[修改密码] 旧密码验证失败")
		return fmt.Errorf("旧密码不正确")
	}
	logger.Info("[修改密码] ✅ 旧密码验证成功")

	// 获取所有账号
	logger.Info("[修改密码] 正在获取所有账号...")
	accounts, err := a.accountService.GetAllAccounts()
	if err != nil {
		logger.Error("[修改密码] 获取账号列表失败: %v", err)
		return fmt.Errorf("获取账号列表失败: %w", err)
	}
	logger.Info("[修改密码] 获取到 %d 个账号", len(accounts))

	// 生成新的盐值
	logger.Info("[修改密码] 正在生成新的盐值...")
	newSalt, err := cryptoManager.GenerateSalt()
	if err != nil {
		logger.Error("[修改密码] 生成新盐值失败: %v", err)
		return fmt.Errorf("生成新盐值失败: %w", err)
	}

	// 哈希新密码
	newHashedPassword := cryptoManager.HashPassword(newPassword, newSalt)

	// 创建新的加密管理器
	newCryptoManager := crypto.NewCryptoManager()
	newCryptoManager.SetMasterPassword(newPassword, newSalt)

	// 开始事务
	logger.Info("[修改密码] 开始数据库事务...")
	tx, err := a.dbManager.GetDB().Begin()
	if err != nil {
		logger.Error("[修改密码] 开始事务失败: %v", err)
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	// 更新密码库配置
	logger.Info("[修改密码] 正在更新密码库配置...")
	newVaultConfig := &models.VaultConfig{
		PasswordHash: newHashedPassword,
		Salt:         base64.StdEncoding.EncodeToString(newSalt),
	}

	_, err = tx.Exec(`UPDATE vault_config SET password_hash = ?, salt = ?`,
		newVaultConfig.PasswordHash, newVaultConfig.Salt)
	if err != nil {
		logger.Error("[修改密码] 更新密码库配置失败: %v", err)
		return fmt.Errorf("更新密码库配置失败: %w", err)
	}

	// 重新加密所有账号
	logger.Info("[修改密码] 开始重新加密账号数据...")
	successCount := 0
	errorCount := 0

	for i, account := range accounts {
		logger.Info("[修改密码] 正在处理账号 %d/%d: %s", i+1, len(accounts), account.Title)

		// 用新密码重新加密账号数据
		newEncryptedUsername, err := newCryptoManager.Encrypt(account.Username)
		if err != nil {
			logger.Error("[修改密码] 重新加密用户名失败，账号: %s, 错误: %v", account.Title, err)
			errorCount++
			continue
		}

		newEncryptedPassword, err := newCryptoManager.Encrypt(account.Password)
		if err != nil {
			logger.Error("[修改密码] 重新加密密码失败，账号: %s, 错误: %v", account.Title, err)
			errorCount++
			continue
		}

		newEncryptedURL, err := newCryptoManager.Encrypt(account.URL)
		if err != nil {
			logger.Error("[修改密码] 重新加密URL失败，账号: %s, 错误: %v", account.Title, err)
			errorCount++
			continue
		}

		newEncryptedNotes, err := newCryptoManager.Encrypt(account.Notes)
		if err != nil {
			logger.Error("[修改密码] 重新加密备注失败，账号: %s, 错误: %v", account.Title, err)
			errorCount++
			continue
		}

		// 更新数据库中的账号数据
		_, err = tx.Exec(`UPDATE accounts SET username = ?, password = ?, url = ?, notes = ?, updated_at = ? WHERE id = ?`,
			newEncryptedUsername, newEncryptedPassword, newEncryptedURL, newEncryptedNotes, time.Now(), account.ID)
		if err != nil {
			logger.Error("[修改密码] 更新账号数据失败，账号: %s, 错误: %v", account.Title, err)
			errorCount++
			continue
		}

		successCount++
		logger.Info("[修改密码] ✅ 账号 %s 重新加密成功", account.Title)
	}

	// 提交事务
	logger.Info("[修改密码] 正在提交事务...")
	if err := tx.Commit(); err != nil {
		logger.Error("[修改密码] 提交事务失败: %v", err)
		return fmt.Errorf("提交事务失败: %w", err)
	}

	// 更新当前的加密管理器
	logger.Info("[修改密码] 正在更新加密管理器...")
	cryptoManager.SetMasterPassword(newPassword, newSalt)
	a.accountService.SetCryptoManager(cryptoManager)

	logger.Info("[修改密码] 🎉 登录密码修改完成！成功: %d, 失败: %d", successCount, errorCount)

	if errorCount > 0 {
		return fmt.Errorf("密码修改完成，但有 %d 个账号处理失败", errorCount)
	}

	return nil
}

/**
 * SelectExportPath 选择导出路径
 * @return string 选择的导出路径，如果取消选择则返回空字符串
 * @author 20251003 陈凤庆 新增导出路径选择功能
 */
func (a *App) SelectExportPath() string {
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "选择导出路径",
		DefaultFilename: "wepass_export.zip",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "ZIP压缩文件 (*.zip)",
				Pattern:     "*.zip",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		log.Printf("导出路径选择对话框错误: %v", err)
		return ""
	}

	return selection
}

/**
 * SelectImportFile 选择导入文件
 * @return string 选择的导入文件路径，如果取消选择则返回空字符串
 * @author 20251003 陈凤庆 新增导入文件选择功能
 */
func (a *App) SelectImportFile() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择导入文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "ZIP压缩文件 (*.zip)",
				Pattern:     "*.zip",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		log.Printf("导入文件选择对话框错误: %v", err)
		return ""
	}

	return selection
}

/**
 * ExportVault 导出密码库
 * @param loginPassword 登录密码
 * @param backupPassword 备份密码
 * @param exportPath 导出路径
 * @param accountIDs 要导出的账号ID列表（手动选择模式）
 * @param groupIDs 要导出的分组ID列表（按分组导出模式）
 * @param typeIDs 要导出的类别ID列表（按类别导出模式）
 * @param exportAll 是否导出所有账号
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增密码库导出功能
 * @modify 20251003 陈凤庆 支持按分组和按类别导出
 */
func (a *App) ExportVault(loginPassword, backupPassword, exportPath string, accountIDs []string, groupIDs []string, typeIDs []string, exportAll bool) error {
	logger.Info("[导出] 开始导出密码库")

	// 检查导出服务是否已初始化
	if a.exportService == nil {
		logger.Error("[导出] 导出服务未初始化")
		return fmt.Errorf("导出服务未初始化")
	}

	// 构建导出选项
	options := services.ExportOptions{
		LoginPassword:  loginPassword,
		BackupPassword: backupPassword,
		ExportPath:     exportPath,
		AccountIDs:     accountIDs,
		GroupIDs:       groupIDs,
		TypeIDs:        typeIDs,
		ExportAll:      exportAll,
	}

	// 执行导出
	if err := a.exportService.ExportVault(options); err != nil {
		logger.Error("[导出] 导出失败: %v", err)
		return err
	}

	logger.Info("[导出] 导出成功: %s", exportPath)
	return nil
}

/**
 * ImportVault 导入密码库
 * @param importPath 导入文件路径
 * @param backupPassword 备份密码（解压密码）
 * @return services.ImportResult 导入结果
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增密码库导入功能
 */
func (a *App) ImportVault(importPath, backupPassword string) (services.ImportResult, error) {
	logger.Info("[导入] 开始导入密码库")

	var result services.ImportResult

	// 检查导入服务是否已初始化
	if a.importService == nil {
		logger.Error("[导入] 导入服务未初始化")
		return result, fmt.Errorf("导入服务未初始化")
	}

	// 构建导入选项
	options := services.ImportOptions{
		ImportPath:     importPath,
		BackupPassword: backupPassword,
	}

	// 执行导入
	result, err := a.importService.ImportVault(options)
	if err != nil {
		logger.Error("[导入] 导入失败: %v", err)
		return result, err
	}

	logger.Info("[导入] 导入完成: 成功=%d, 跳过=%d, 错误=%d",
		result.ImportedAccounts, result.SkippedAccounts, result.ErrorAccounts)
	return result, nil
}

/**
 * initGlobalHotkeyService 初始化全局快捷键服务
 * @author 陈凤庆
 * @date 20251014
 * @description 初始化全局快捷键服务，支持显示/隐藏窗口功能
 */
func (a *App) initGlobalHotkeyService() {
	logger.Info("[全局快捷键] 初始化全局快捷键服务...")

	// 创建全局快捷键服务
	hotkeyService, err := services.NewGlobalHotkeyService(a.configManager, a)
	if err != nil {
		logger.Error("[全局快捷键] 初始化失败: %v", err)
		return
	}

	a.globalHotkeyService = hotkeyService

	// 启动快捷键监听
	if err := a.globalHotkeyService.Start(); err != nil {
		logger.Error("[全局快捷键] 启动监听失败: %v", err)
		return
	}

	logger.Info("[全局快捷键] 全局快捷键服务初始化完成")
}

/**
 * initPasswordRuleApp 初始化密码规则应用服务
 * @author 陈凤庆
 * @date 20251017
 * @description 初始化密码规则应用服务，提供密码规则管理功能
 */
func (a *App) initPasswordRuleApp() {
	logger.Info("[密码规则] 初始化密码规则应用服务...")

	// 创建密码规则服务
	passwordRuleService := services.NewPasswordRuleService(a.dbManager)

	// 创建密码规则应用服务
	a.passwordRuleApp = NewPasswordRuleApp(passwordRuleService)

	logger.Info("[密码规则] 密码规则应用服务初始化完成")
}

/**
 * initUsernameHistoryApp 初始化用户名历史记录应用服务
 * @author 陈凤庆
 * @date 20251017
 * @description 初始化用户名历史记录应用服务，提供用户名历史记录管理功能
 */
func (a *App) initUsernameHistoryApp() {
	logger.Info("[用户名历史] 初始化用户名历史记录应用服务...")

	// 创建用户名历史记录服务
	usernameHistoryService := services.NewUsernameHistoryService(a.dbManager)

	// 创建用户名历史记录应用服务
	a.usernameHistoryApp = NewUsernameHistoryApp(usernameHistoryService)

	logger.Info("[用户名历史] 用户名历史记录应用服务初始化完成")
}

// 密码规则管理API

/**
 * GetAllPasswordRules 获取所有密码规则
 * @return []models.PasswordRule 密码规则列表
 * @return error 错误信息
 */
func (a *App) GetAllPasswordRules() ([]models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return nil, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.GetAllRules(a.ctx)
}

/**
 * GetPasswordRuleByID 根据ID获取密码规则
 * @param id 规则ID
 * @return models.PasswordRule 密码规则
 * @return error 错误信息
 */
func (a *App) GetPasswordRuleByID(id string) (models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return models.PasswordRule{}, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.GetRuleByID(a.ctx, id)
}

/**
 * GetDefaultPasswordRule 获取默认密码规则
 * @return models.PasswordRule 默认密码规则
 * @return error 错误信息
 */
func (a *App) GetDefaultPasswordRule() (models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return models.PasswordRule{}, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.GetDefaultRule(a.ctx)
}

/**
 * CreateGeneralPasswordRule 创建通用密码规则
 * @param name 规则名称
 * @param description 规则描述
 * @param config 通用规则配置
 * @return models.PasswordRule 创建的密码规则
 * @return error 错误信息
 */
func (a *App) CreateGeneralPasswordRule(name, description string, config models.GeneralRuleConfig) (models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return models.PasswordRule{}, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.CreateGeneralRule(a.ctx, name, description, config)
}

/**
 * CreateCustomPasswordRule 创建自定义密码规则
 * @param name 规则名称
 * @param description 规则描述
 * @param config 自定义规则配置
 * @return models.PasswordRule 创建的密码规则
 * @return error 错误信息
 */
func (a *App) CreateCustomPasswordRule(name, description string, config models.CustomRuleConfig) (models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return models.PasswordRule{}, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.CreateCustomRule(a.ctx, name, description, config)
}

/**
 * UpdateGeneralPasswordRule 更新通用密码规则
 * @param id 规则ID
 * @param name 规则名称
 * @param description 规则描述
 * @param config 通用规则配置
 * @return models.PasswordRule 更新后的密码规则
 * @return error 错误信息
 */
func (a *App) UpdateGeneralPasswordRule(id, name, description string, config models.GeneralRuleConfig) (models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return models.PasswordRule{}, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.UpdateGeneralRule(a.ctx, id, name, description, config)
}

/**
 * UpdateCustomPasswordRule 更新自定义密码规则
 * @param id 规则ID
 * @param name 规则名称
 * @param description 规则描述
 * @param config 自定义规则配置
 * @return models.PasswordRule 更新后的密码规则
 * @return error 错误信息
 */
func (a *App) UpdateCustomPasswordRule(id, name, description string, config models.CustomRuleConfig) (models.PasswordRule, error) {
	if a.passwordRuleApp == nil {
		return models.PasswordRule{}, fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.UpdateCustomRule(a.ctx, id, name, description, config)
}

/**
 * DeletePasswordRule 删除密码规则
 * @param id 规则ID
 * @return error 错误信息
 */
func (a *App) DeletePasswordRule(id string) error {
	if a.passwordRuleApp == nil {
		return fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.DeleteRule(a.ctx, id)
}

/**
 * GeneratePasswordByRule 根据规则生成密码
 * @param ruleID 规则ID
 * @return string 生成的密码
 * @return error 错误信息
 */
func (a *App) GeneratePasswordByRule(ruleID string) (string, error) {
	if a.passwordRuleApp == nil {
		return "", fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.GeneratePassword(a.ctx, ruleID)
}

/**
 * GeneratePasswordByGeneralConfig 根据通用配置生成密码
 * @param config 通用规则配置
 * @return string 生成的密码
 * @return error 错误信息
 */
func (a *App) GeneratePasswordByGeneralConfig(config models.GeneralRuleConfig) (string, error) {
	if a.passwordRuleApp == nil {
		return "", fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.GeneratePasswordByGeneralConfig(a.ctx, config)
}

/**
 * GeneratePasswordByCustomConfig 根据自定义配置生成密码
 * @param config 自定义规则配置
 * @return string 生成的密码
 * @return error 错误信息
 */
func (a *App) GeneratePasswordByCustomConfig(config models.CustomRuleConfig) (string, error) {
	if a.passwordRuleApp == nil {
		return "", fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.GeneratePasswordByCustomConfig(a.ctx, config)
}

/**
 * SetPasswordRuleAsDefault 设置密码规则为默认规则
 * @param ruleID 规则ID
 * @param isDefault 是否设为默认
 * @return error 错误信息
 */
func (a *App) SetPasswordRuleAsDefault(ruleID string, isDefault bool) error {
	if a.passwordRuleApp == nil {
		return fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.SetRuleAsDefault(a.ctx, ruleID, isDefault)
}

/**
 * ForceInitializeDefaultPasswordRules 强制初始化默认密码规则
 * @param force 是否强制重新创建
 * @return error 错误信息
 */
func (a *App) ForceInitializeDefaultPasswordRules(force bool) error {
	if a.passwordRuleApp == nil {
		return fmt.Errorf("密码规则应用服务未初始化")
	}
	return a.passwordRuleApp.ForceInitializeDefaultRules(a.ctx, force)
}

// 用户名历史记录管理API

/**
 * GetUsernameHistory 获取用户名历史记录
 * @param password 登录密码，用于解密
 * @return []string 用户名列表
 * @return error 错误信息
 */
func (a *App) GetUsernameHistory(password string) ([]string, error) {
	if a.usernameHistoryApp == nil {
		return nil, fmt.Errorf("用户名历史记录应用服务未初始化")
	}
	return a.usernameHistoryApp.GetUsernameHistory(a.ctx, password)
}

/**
 * SaveUsernameToHistory 保存用户名到历史记录
 * @param username 用户名
 * @param password 登录密码，用于加密
 * @return error 错误信息
 */
func (a *App) SaveUsernameToHistory(username, password string) error {
	if a.usernameHistoryApp == nil {
		return fmt.Errorf("用户名历史记录应用服务未初始化")
	}
	return a.usernameHistoryApp.SaveUsernameToHistory(a.ctx, username, password)
}

/**
 * ClearUsernameHistory 清空用户名历史记录
 * @return error 错误信息
 */
func (a *App) ClearUsernameHistory() error {
	if a.usernameHistoryApp == nil {
		return fmt.Errorf("用户名历史记录应用服务未初始化")
	}
	return a.usernameHistoryApp.ClearUsernameHistory(a.ctx)
}

/**
 * GetCurrentPassword 获取当前登录密码
 * @return string 当前登录密码
 * @description 20251017 陈凤庆 从后端内存中获取当前登录密码，用于用户名历史记录加密
 */
func (a *App) GetCurrentPassword() string {
	if a.vaultService == nil {
		logger.Error("[API] 密码库服务未初始化")
		return ""
	}
	return a.vaultService.GetCurrentPassword()
}

// 键盘输入模拟相关API

// 20251003 陈凤庆 键盘服务相关方法已移动到平台特定文件中
// app_darwin.go - macOS平台实现
// app_windows.go - Windows平台实现
