package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"wepassword/internal/logger"
	"wepassword/internal/models"
	"wepassword/internal/utils"
)

/**
 * 配置管理模块
 * @author 陈凤庆
 * @description 管理应用配置文件的读写操作
 */

/**
 * ConfigManager 配置管理器
 */
type ConfigManager struct {
	configPath string
	config     *models.AppConfig
}

/**
 * NewConfigManager 创建新的配置管理器
 * @return *ConfigManager 配置管理器实例
 * @modify 20251002 陈凤庆 使用跨平台路径工具函数获取配置文件路径
 */
func NewConfigManager() *ConfigManager {
	// 20251002 陈凤庆 使用跨平台路径工具函数获取应用数据目录
	appDataDir, err := utils.GetAppDataDir()
	if err != nil {
		// 如果获取失败，回退到用户主目录下的 .wepass
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			homeDir = "."
		}
		appDataDir = filepath.Join(homeDir, ".wepass")
	}

	configPath := filepath.Join(appDataDir, "config.json")

	// 20251005 陈凤庆 获取基于环境的默认日志配置
	defaultLogConfig := logger.GetDefaultLogConfig()

	cm := &ConfigManager{
		configPath: configPath,
		config: &models.AppConfig{
			CurrentVaultPath: "",
			RecentVaults:     make([]string, 0),
			WindowWidth:      550,
			WindowHeight:     800,
			Theme:            "light",
			Language:         "zh-CN",
			// 20251005 陈凤庆 使用基于环境的默认日志配置
			LogConfig: models.LogConfig{
				EnableInfoLog:  defaultLogConfig.EnableInfoLog,
				EnableDebugLog: defaultLogConfig.EnableDebugLog,
			},
			// 20251004 陈凤庆 添加锁定配置默认值
			LockConfig: models.LockConfig{
				EnableAutoLock:     false, // 默认不启用自动锁定
				EnableTimerLock:    false, // 默认不启用定时锁定
				EnableMinimizeLock: false, // 默认不启用最小化锁定
				LockTimeMinutes:    10,    // 默认10分钟
				EnableSystemLock:   true,  // 系统锁定固定启用
				SystemLockMinutes:  120,   // 固定2小时
			},
			// 20251014 陈凤庆 添加快捷键配置默认值
			HotkeyConfig: models.HotkeyConfig{
				EnableGlobalHotkey:   true,         // 默认启用全局快捷键
				ShowHideHotkey:       "Ctrl+Alt+H", // 默认快捷键为Ctrl+Alt+H
				EnableShowHideHotkey: true,         // 默认启用显示/隐藏快捷键
			},
		},
	}

	// 加载配置
	cm.LoadConfig()
	return cm
}

/**
 * LoadConfig 加载配置文件
 * @return error 错误信息
 */
func (cm *ConfigManager) LoadConfig() error {
	// 检查配置文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置
		return cm.SaveConfig()
	}

	// 读取配置文件
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 JSON
	if err := json.Unmarshal(data, cm.config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

/**
 * SaveConfig 保存配置文件
 * @return error 错误信息
 */
func (cm *ConfigManager) SaveConfig() error {
	// 确保配置目录存在
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 序列化配置
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

/**
 * GetCurrentVaultPath 获取当前密码库路径
 * @return string 当前密码库路径
 */
func (cm *ConfigManager) GetCurrentVaultPath() string {
	return cm.config.CurrentVaultPath
}

/**
 * SetCurrentVaultPath 设置当前密码库路径
 * @param path 密码库路径
 * @return error 错误信息
 */
func (cm *ConfigManager) SetCurrentVaultPath(path string) error {
	cm.config.CurrentVaultPath = path

	// 更新最近使用列表
	cm.addToRecentVaults(path)

	return cm.SaveConfig()
}

/**
 * GetRecentVaults 获取最近使用的密码库列表
 * @return []string 最近使用的密码库路径列表
 */
func (cm *ConfigManager) GetRecentVaults() []string {
	return cm.config.RecentVaults
}

/**
 * addToRecentVaults 添加到最近使用列表
 * @param path 密码库路径
 */
func (cm *ConfigManager) addToRecentVaults(path string) {
	// 移除已存在的相同路径
	for i, vault := range cm.config.RecentVaults {
		if vault == path {
			cm.config.RecentVaults = append(cm.config.RecentVaults[:i], cm.config.RecentVaults[i+1:]...)
			break
		}
	}

	// 添加到列表开头
	cm.config.RecentVaults = append([]string{path}, cm.config.RecentVaults...)

	// 保持最多5条记录
	if len(cm.config.RecentVaults) > 5 {
		cm.config.RecentVaults = cm.config.RecentVaults[:5]
	}
}

/**
 * GetWindowSize 获取窗口大小
 * @return int, int 宽度和高度
 */
func (cm *ConfigManager) GetWindowSize() (int, int) {
	return cm.config.WindowWidth, cm.config.WindowHeight
}

/**
 * SetWindowSize 设置窗口大小
 * @param width 宽度
 * @param height 高度
 * @return error 错误信息
 */
func (cm *ConfigManager) SetWindowSize(width, height int) error {
	cm.config.WindowWidth = width
	cm.config.WindowHeight = height
	return cm.SaveConfig()
}

/**
 * GetTheme 获取主题
 * @return string 主题名称
 */
func (cm *ConfigManager) GetTheme() string {
	return cm.config.Theme
}

/**
 * SetTheme 设置主题
 * @param theme 主题名称
 * @return error 错误信息
 */
func (cm *ConfigManager) SetTheme(theme string) error {
	cm.config.Theme = theme
	return cm.SaveConfig()
}

/**
 * GetLanguage 获取语言
 * @return string 语言代码
 */
func (cm *ConfigManager) GetLanguage() string {
	return cm.config.Language
}

/**
 * SetLanguage 设置语言
 * @param language 语言代码
 * @return error 错误信息
 */
func (cm *ConfigManager) SetLanguage(language string) error {
	cm.config.Language = language
	return cm.SaveConfig()
}

/**
 * GetConfig 获取完整配置
 * @return *models.AppConfig 应用配置
 * @author 陈凤庆
 * @date 20251003
 */
func (cm *ConfigManager) GetConfig() *models.AppConfig {
	return cm.config
}

/**
 * GetLogConfig 获取日志配置
 * @return models.LogConfig 日志配置
 * @author 陈凤庆
 * @date 20251003
 */
func (cm *ConfigManager) GetLogConfig() models.LogConfig {
	return cm.config.LogConfig
}

/**
 * SetLogConfig 设置日志配置
 * @param logConfig 日志配置
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251003
 */
func (cm *ConfigManager) SetLogConfig(logConfig models.LogConfig) error {
	cm.config.LogConfig = logConfig
	return cm.SaveConfig()
}

/**
 * GetLockConfig 获取锁定配置
 * @return models.LockConfig 锁定配置
 * @author 陈凤庆
 * @date 20251004
 */
func (cm *ConfigManager) GetLockConfig() models.LockConfig {
	return cm.config.LockConfig
}

/**
 * SetLockConfig 设置锁定配置
 * @param lockConfig 锁定配置
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251004
 */
func (cm *ConfigManager) SetLockConfig(lockConfig models.LockConfig) error {
	cm.config.LockConfig = lockConfig
	return cm.SaveConfig()
}

/**
 * IsInfoLogEnabled 检查是否启用Info日志
 * @return bool 是否启用
 * @author 陈凤庆
 * @date 20251003
 */
func (cm *ConfigManager) IsInfoLogEnabled() bool {
	return cm.config.LogConfig.EnableInfoLog
}

/**
 * IsDebugLogEnabled 检查是否启用Debug日志
 * @return bool 是否启用
 * @author 陈凤庆
 * @date 20251003
 */
func (cm *ConfigManager) IsDebugLogEnabled() bool {
	return cm.config.LogConfig.EnableDebugLog
}

/**
 * GetHotkeyConfig 获取快捷键配置
 * @return models.HotkeyConfig 快捷键配置
 * @author 陈凤庆
 * @date 20251014
 */
func (cm *ConfigManager) GetHotkeyConfig() models.HotkeyConfig {
	return cm.config.HotkeyConfig
}

/**
 * SetHotkeyConfig 设置快捷键配置
 * @param config 快捷键配置
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251014
 */
func (cm *ConfigManager) SetHotkeyConfig(config models.HotkeyConfig) error {
	cm.config.HotkeyConfig = config
	return cm.SaveConfig()
}
