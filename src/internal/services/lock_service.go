package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"wepassword/internal/config"
	"wepassword/internal/logger"
	"wepassword/internal/models"
)

/**
 * LockService 简化锁定服务
 * @author 陈凤庆
 * @date 20251004
 * @description 简化的锁定服务，只负责定时器管理，不管理锁定状态
 * 锁定逻辑由前端处理：检测到锁定条件时直接返回登录界面
 */

// LockConfig 锁定配置
type LockConfig struct {
	EnableAutoLock     bool `json:"enable_auto_lock"`     // 是否启用自动锁定
	EnableTimerLock    bool `json:"enable_timer_lock"`    // 是否启用定时锁定
	EnableMinimizeLock bool `json:"enable_minimize_lock"` // 是否启用最小化锁定
	LockTimeMinutes    int  `json:"lock_time_minutes"`    // 定时锁定时间（分钟）
	EnableSystemLock   bool `json:"enable_system_lock"`   // 是否启用系统固定锁定（固定为true）
	SystemLockMinutes  int  `json:"system_lock_minutes"`  // 系统锁定时间（固定120分钟）
}

// LockService 简化锁定服务结构体
type LockService struct {
	mu               sync.RWMutex
	configManager    *config.ConfigManager
	config           *LockConfig
	userTimer        *time.Timer // 用户设置的定时器
	systemTimer      *time.Timer // 系统固定定时器
	ctx              context.Context
	cancel           context.CancelFunc
	running          bool
	lastActivityTime time.Time // 最后活动时间
	lockTriggered    bool      // 锁定触发标志（供前端查询）
}

// 默认配置常量
const (
	DefaultLockTimeMinutes = 10  // 默认10分钟锁定
	SystemLockMinutes      = 120 // 系统固定2小时锁定
)

/**
 * NewLockService 创建新的锁定服务
 * @param configManager 配置管理器
 * @return *LockService 锁定服务实例
 */
func NewLockService(configManager *config.ConfigManager) *LockService {
	log.Println("[锁定服务] 开始创建锁定服务")
	logger.Info("[锁定服务] 开始创建锁定服务")

	ctx, cancel := context.WithCancel(context.Background())

	service := &LockService{
		configManager:    configManager,
		ctx:              ctx,
		cancel:           cancel,
		lastActivityTime: time.Now(),
		config: &LockConfig{
			EnableAutoLock:     false,                  // 默认不启用自动锁定
			EnableTimerLock:    false,                  // 默认不启用定时锁定
			EnableMinimizeLock: false,                  // 默认不启用最小化锁定
			LockTimeMinutes:    DefaultLockTimeMinutes, // 默认10分钟
			EnableSystemLock:   true,                   // 系统锁定固定启用
			SystemLockMinutes:  SystemLockMinutes,      // 固定2小时
		},
	}

	log.Println("[锁定服务] 锁定服务结构体已创建")
	logger.Info("[锁定服务] 锁定服务结构体已创建")

	// 加载配置
	service.loadConfig()

	log.Println("[锁定服务] 锁定服务初始化完成")
	logger.Info("[锁定服务] 锁定服务初始化完成")
	return service
}

/**
 * IsLockTriggered 检查是否触发了锁定
 * @return bool 是否触发锁定
 */
func (ls *LockService) IsLockTriggered() bool {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.lockTriggered
}

/**
 * ResetLockTrigger 重置锁定触发标志（登录后调用）
 */
func (ls *LockService) ResetLockTrigger() {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.lockTriggered = false
	log.Println("[锁定服务] 锁定触发标志已重置")
}

/**
 * GetTimerStatus 获取定时器状态（用于调试）
 */
func (ls *LockService) GetTimerStatus() map[string]interface{} {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	status := map[string]interface{}{
		"running":           ls.running,
		"lockTriggered":     ls.lockTriggered,
		"lastActivityTime":  ls.lastActivityTime,
		"userTimerActive":   ls.userTimer != nil,
		"systemTimerActive": ls.systemTimer != nil,
		"config": map[string]interface{}{
			"enableAutoLock":     ls.config.EnableAutoLock,
			"enableTimerLock":    ls.config.EnableTimerLock,
			"enableMinimizeLock": ls.config.EnableMinimizeLock,
			"lockTimeMinutes":    ls.config.LockTimeMinutes,
			"systemLockMinutes":  ls.config.SystemLockMinutes,
		},
	}

	return status
}

/**
 * StartLockService 启动锁定服务
 * @return error 启动错误
 */
func (ls *LockService) StartLockService() error {
	log.Println("[锁定服务] 开始启动锁定服务")
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// 20251004 陈凤庆 每次登录都重新启动定时器，确保计时重新开始
	if ls.running {
		log.Println("[锁定服务] 服务已在运行，重新启动定时器")
		// 停止现有定时器
		if ls.userTimer != nil {
			ls.userTimer.Stop()
			ls.userTimer = nil
			log.Println("[锁定服务] 已停止现有用户定时器")
		}
		if ls.systemTimer != nil {
			ls.systemTimer.Stop()
			ls.systemTimer = nil
			log.Println("[锁定服务] 已停止现有系统定时器")
		}
	}

	ls.running = true
	ls.lastActivityTime = time.Now()
	log.Printf("[锁定服务] 设置运行状态为true，最后活动时间: %v", ls.lastActivityTime)

	// 启动系统固定锁定定时器（始终启用）
	ls.startSystemTimer()
	log.Println("[锁定服务] 系统定时器已启动")

	// 根据配置启动用户定时器
	if ls.config.EnableAutoLock && ls.config.EnableTimerLock {
		ls.startUserTimer()
		log.Println("[锁定服务] 用户定时器已启动")
	} else {
		log.Printf("[锁定服务] 用户定时器未启动 - 启用=%t, 定时锁定=%t", ls.config.EnableAutoLock, ls.config.EnableTimerLock)
	}

	log.Println("[锁定服务] 锁定服务启动完成")
	logger.Info("[锁定服务] 锁定服务已启动")
	logger.Info("[锁定服务] 用户锁定配置: 启用=%t, 定时锁定=%t, 最小化锁定=%t, 时间=%d分钟",
		ls.config.EnableAutoLock, ls.config.EnableTimerLock, ls.config.EnableMinimizeLock, ls.config.LockTimeMinutes)
	logger.Info("[锁定服务] 系统锁定配置: 启用=%t, 时间=%d分钟",
		ls.config.EnableSystemLock, ls.config.SystemLockMinutes)

	return nil
}

/**
 * StopLockService 停止锁定服务
 */
func (ls *LockService) StopLockService() {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if !ls.running {
		return
	}

	ls.running = false

	// 停止定时器
	if ls.userTimer != nil {
		ls.userTimer.Stop()
		ls.userTimer = nil
	}

	if ls.systemTimer != nil {
		ls.systemTimer.Stop()
		ls.systemTimer = nil
	}

	// 取消上下文
	ls.cancel()

	logger.Info("[锁定服务] 锁定服务已停止")
}

/**
 * UpdateActivity 更新用户活动时间
 * @description 当用户有操作时调用此方法重置定时器
 */
func (ls *LockService) UpdateActivity() {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.lastActivityTime = time.Now()

	// 重启用户定时器
	if ls.running && ls.config.EnableAutoLock && ls.config.EnableTimerLock {
		ls.restartUserTimer()
	}

	logger.Debug("[锁定服务] 用户活动时间已更新")
}

/**
 * OnMinimize 处理窗口最小化事件
 * @description 当窗口最小化时调用
 */
func (ls *LockService) OnMinimize() {
	log.Println("[锁定服务] 收到窗口最小化事件")
	ls.mu.RLock()
	running := ls.running
	enableAutoLock := ls.config.EnableAutoLock
	enableMinimizeLock := ls.config.EnableMinimizeLock
	shouldLock := running && enableAutoLock && enableMinimizeLock
	ls.mu.RUnlock()

	log.Printf("[锁定服务] 最小化检查: running=%t, enableAutoLock=%t, enableMinimizeLock=%t, shouldLock=%t",
		running, enableAutoLock, enableMinimizeLock, shouldLock)

	if shouldLock {
		log.Println("[锁定服务] 检测到窗口最小化，触发锁定")
		logger.Info("[锁定服务] 检测到窗口最小化，触发锁定")
		ls.triggerLock("最小化锁定")
	} else {
		log.Println("[锁定服务] 窗口最小化但不满足锁定条件，跳过锁定")
		logger.Info("[锁定服务] 窗口最小化但不满足锁定条件，跳过锁定")
	}
}

/**
 * GetLockConfig 获取锁定配置
 * @return *LockConfig 锁定配置
 */
func (ls *LockService) GetLockConfig() *LockConfig {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	// 返回配置副本
	return &LockConfig{
		EnableAutoLock:     ls.config.EnableAutoLock,
		EnableTimerLock:    ls.config.EnableTimerLock,
		EnableMinimizeLock: ls.config.EnableMinimizeLock,
		LockTimeMinutes:    ls.config.LockTimeMinutes,
		EnableSystemLock:   ls.config.EnableSystemLock,
		SystemLockMinutes:  ls.config.SystemLockMinutes,
	}
}

/**
 * UpdateLockConfig 更新锁定配置
 * @param config 新的锁定配置
 * @return error 更新错误
 */
func (ls *LockService) UpdateLockConfig(config *LockConfig) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// 验证配置
	if config.LockTimeMinutes < 1 || config.LockTimeMinutes > 1440 { // 1分钟到24小时
		return fmt.Errorf("锁定时间必须在1-1440分钟之间")
	}

	// 更新配置
	ls.config.EnableAutoLock = config.EnableAutoLock
	ls.config.EnableTimerLock = config.EnableTimerLock
	ls.config.EnableMinimizeLock = config.EnableMinimizeLock
	ls.config.LockTimeMinutes = config.LockTimeMinutes
	// 系统锁定配置不允许修改
	ls.config.EnableSystemLock = true
	ls.config.SystemLockMinutes = SystemLockMinutes

	// 保存配置
	if err := ls.saveConfig(); err != nil {
		return fmt.Errorf("保存锁定配置失败: %w", err)
	}

	// 重启定时器
	if ls.running {
		ls.restartTimers()
	}

	logger.Info("[锁定服务] 锁定配置已更新: 启用=%t, 定时锁定=%t, 最小化锁定=%t, 时间=%d分钟",
		ls.config.EnableAutoLock, ls.config.EnableTimerLock, ls.config.EnableMinimizeLock, ls.config.LockTimeMinutes)

	return nil
}

/**
 * startUserTimer 启动用户定时器
 */
func (ls *LockService) startUserTimer() {
	if ls.userTimer != nil {
		ls.userTimer.Stop()
		log.Println("[锁定服务] 停止现有用户定时器")
	}

	duration := time.Duration(ls.config.LockTimeMinutes) * time.Minute
	startTime := time.Now()

	ls.userTimer = time.AfterFunc(duration, func() {
		log.Printf("[锁定服务] 用户定时器触发！启动时间: %v, 当前时间: %v, 经过时间: %v",
			startTime, time.Now(), time.Since(startTime))
		ls.triggerLock("用户定时锁定")
	})

	log.Printf("[锁定服务] 用户定时器已启动，%d分钟后锁定，启动时间: %v", ls.config.LockTimeMinutes, startTime)
	logger.Info("[锁定服务] 用户定时器已启动，%d分钟后锁定", ls.config.LockTimeMinutes)
}

/**
 * startSystemTimer 启动系统定时器
 */
func (ls *LockService) startSystemTimer() {
	if ls.systemTimer != nil {
		ls.systemTimer.Stop()
		log.Println("[锁定服务] 停止现有系统定时器")
	}

	duration := time.Duration(ls.config.SystemLockMinutes) * time.Minute
	startTime := time.Now()

	ls.systemTimer = time.AfterFunc(duration, func() {
		log.Printf("[锁定服务] 系统定时器触发！启动时间: %v, 当前时间: %v, 经过时间: %v",
			startTime, time.Now(), time.Since(startTime))
		ls.triggerLock("系统固定锁定")
	})

	log.Printf("[锁定服务] 系统定时器已启动，%d分钟后锁定，启动时间: %v", ls.config.SystemLockMinutes, startTime)
	logger.Info("[锁定服务] 系统定时器已启动，%d分钟后锁定", ls.config.SystemLockMinutes)
}

/**
 * restartUserTimer 重启用户定时器
 */
func (ls *LockService) restartUserTimer() {
	// 先停止现有定时器
	if ls.userTimer != nil {
		ls.userTimer.Stop()
		ls.userTimer = nil
	}

	// 如果启用了定时锁定，重新启动定时器
	if ls.config.EnableAutoLock && ls.config.EnableTimerLock {
		ls.startUserTimer()
		logger.Debug("[锁定服务] 用户定时器已重启")
	}
}

/**
 * restartTimers 重启所有定时器
 */
func (ls *LockService) restartTimers() {
	// 重启用户定时器
	if ls.config.EnableAutoLock && ls.config.EnableTimerLock {
		ls.startUserTimer()
	} else if ls.userTimer != nil {
		ls.userTimer.Stop()
		ls.userTimer = nil
	}

	// 系统定时器始终运行
	ls.startSystemTimer()
}

/**
 * triggerLock 触发锁定（简化版本）
 * @param reason 锁定原因
 */
func (ls *LockService) triggerLock(reason string) {
	log.Printf("[锁定服务] 触发锁定，原因: %s", reason)
	logger.Info("[锁定服务] 触发锁定，原因: %s", reason)

	// 设置锁定触发标志，供前端查询
	ls.mu.Lock()
	ls.lockTriggered = true

	// 停止所有定时器（锁定后不应该继续计时）
	if ls.userTimer != nil {
		ls.userTimer.Stop()
		ls.userTimer = nil
		log.Println("[锁定服务] 用户定时器已停止")
	}
	if ls.systemTimer != nil {
		ls.systemTimer.Stop()
		ls.systemTimer = nil
		log.Println("[锁定服务] 系统定时器已停止")
	}
	ls.running = false // 停止锁定服务
	ls.mu.Unlock()

	log.Println("[锁定服务] 锁定触发完成，前端将检测到此状态")
	logger.Info("[锁定服务] 锁定触发完成，前端将检测到此状态")
}

/**
 * loadConfig 加载配置
 */
func (ls *LockService) loadConfig() {
	if ls.configManager != nil {
		// 从配置管理器加载锁定配置
		modelConfig := ls.configManager.GetLockConfig()

		// 转换为服务配置格式
		ls.config.EnableAutoLock = modelConfig.EnableAutoLock
		ls.config.EnableTimerLock = modelConfig.EnableTimerLock
		ls.config.EnableMinimizeLock = modelConfig.EnableMinimizeLock
		ls.config.LockTimeMinutes = modelConfig.LockTimeMinutes
		ls.config.EnableSystemLock = modelConfig.EnableSystemLock
		ls.config.SystemLockMinutes = modelConfig.SystemLockMinutes

		log.Printf("[锁定服务] 已加载配置: 启用=%t, 定时锁定=%t, 最小化锁定=%t, 时间=%d分钟",
			ls.config.EnableAutoLock, ls.config.EnableTimerLock, ls.config.EnableMinimizeLock, ls.config.LockTimeMinutes)
		logger.Info("[锁定服务] 已从配置管理器加载锁定配置")
	} else {
		log.Println("[锁定服务] 配置管理器为空，使用默认配置")
		logger.Info("[锁定服务] 配置管理器为空，使用默认锁定配置")
	}
}

/**
 * saveConfig 保存配置
 * @return error 保存错误
 */
func (ls *LockService) saveConfig() error {
	if ls.configManager != nil {
		// 转换为模型配置格式
		modelConfig := models.LockConfig{
			EnableAutoLock:     ls.config.EnableAutoLock,
			EnableTimerLock:    ls.config.EnableTimerLock,
			EnableMinimizeLock: ls.config.EnableMinimizeLock,
			LockTimeMinutes:    ls.config.LockTimeMinutes,
			EnableSystemLock:   ls.config.EnableSystemLock,
			SystemLockMinutes:  ls.config.SystemLockMinutes,
		}

		// 保存到配置管理器
		if err := ls.configManager.SetLockConfig(modelConfig); err != nil {
			return fmt.Errorf("保存锁定配置到配置管理器失败: %w", err)
		}

		logger.Info("[锁定服务] 锁定配置已保存到配置文件")
		return nil
	}

	logger.Info("[锁定服务] 配置管理器为空，跳过保存")
	return nil
}
