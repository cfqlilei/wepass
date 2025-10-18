package models

import (
	"time"
)

/**
 * 数据模型定义
 * @author 陈凤庆
 * @description 定义密码管理工具的核心数据结构
 */

/**
 * Group 分组模型
 * @modify 20251002 陈凤庆 ID字段改为string类型，避免Wails传输时JavaScript精度丢失
 * @modify 20251002 陈凤庆 删除parent_id字段，不需要层级结构
 */
type Group struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Icon      string    `json:"icon" db:"icon"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

/**
 * Account 账号模型（原PasswordItem）
 * @modify 20251002 陈凤庆 删除TabID和GroupID字段，通过TypeID关联
 * @modify 20251003 陈凤庆 添加InputMethod字段，支持三种输入方式
 */
type Account struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`               // 标题（明文）
	Username    string    `json:"username" db:"username"`         // 用户名（加密）
	Password    string    `json:"password" db:"password"`         // 密码（加密）
	URL         string    `json:"url" db:"url"`                   // 地址（加密）
	TypeID      string    `json:"typeid" db:"typeid"`             // 类型ID（明文）
	Notes       string    `json:"notes" db:"notes"`               // 备注（加密）
	Icon        string    `json:"icon" db:"icon"`                 // 图标
	IsFavorite  bool      `json:"is_favorite" db:"is_favorite"`   // 是否收藏
	UseCount    int       `json:"use_count" db:"use_count"`       // 使用次数
	LastUsedAt  time.Time `json:"last_used_at" db:"last_used_at"` // 最后使用时间
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	InputMethod int       `json:"input_method" db:"input_method"` // 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
}

/**
 * PasswordItem 账号模型（为了兼容性保留的别名）
 * @deprecated 请使用Account模型
 */
type PasswordItem = Account

/**
 * AccountDecrypted 解密后的账号模型（用于前端显示）
 * @modify 20250101 陈凤庆 添加TabID字段，用于关联标签
 * @modify 20251001 陈凤庆 ID字段改为string类型，避免JavaScript精度丢失
 * @modify 20251001 陈凤庆 PasswordItemDecrypted改名为AccountDecrypted，Type字段改为TypeID，删除created_by和updated_by字段
 * @modify 20251003 陈凤庆 添加InputMethod字段，支持三种输入方式
 */
type AccountDecrypted struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Username       string    `json:"username"` // 解密后的用户名
	Password       string    `json:"password"` // 解密后的密码
	URL            string    `json:"url"`      // 解密后的地址
	TypeID         string    `json:"typeid"`   // 类型ID
	GroupID        string    `json:"group_id"` // 20251002 陈凤庆 添加分组ID字段，用于前端加载类型列表
	Notes          string    `json:"notes"`    // 解密后的备注
	Icon           string    `json:"icon"`
	IsFavorite     bool      `json:"is_favorite"`
	UseCount       int       `json:"use_count"`
	LastUsedAt     time.Time `json:"last_used_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	InputMethod    int       `json:"input_method"`    // 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
	MaskedUsername string    `json:"masked_username"` // 脱敏用户名，用于列表显示
	MaskedPassword string    `json:"masked_password"` // 脱敏密码，用于列表显示
}

/**
 * PasswordItemDecrypted 解密后的账号模型（为了兼容性保留的别名）
 * @deprecated 请使用AccountDecrypted模型
 */
type PasswordItemDecrypted = AccountDecrypted

/**
 * VaultConfig 密码库配置模型
 * @modify 20251001 陈凤庆 ID字段改为string类型，避免JavaScript精度丢失
 */
type VaultConfig struct {
	ID           string    `json:"id" db:"id"`
	PasswordHash string    `json:"password_hash" db:"password_hash"` // 登录密码哈希
	Salt         string    `json:"salt" db:"salt"`                   // 密码盐值
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

/**
 * Type 类型模型（原Tab）
 * @modify 20251001 陈凤庆 ID字段改为string类型，避免JavaScript精度丢失
 * @modify 20251001 陈凤庆 Tab改名为Type，删除Type字段，删除created_by和updated_by字段
 */
type Type struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Icon      string    `json:"icon" db:"icon"`
	Filter    string    `json:"filter" db:"filter"`     // 筛选规则（JSON格式）
	GroupID   string    `json:"group_id" db:"group_id"` // 所属分组ID
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

/**
 * Tab 页签模型（为了兼容性保留的别名）
 * @deprecated 请使用Type模型
 */
type Tab = Type

/**
 * PasswordRule 密码规则模型
 * @author 陈凤庆
 * @description 存储密码生成规则的配置信息
 */
type PasswordRule struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`                 // 规则名称
	Description string    `json:"description" db:"description"`   // 规则描述
	RuleType    string    `json:"rule_type" db:"rule_type"`       // 规则类型：general(通用规则)、custom(自定义规则)
	Config      string    `json:"config" db:"config"`             // 规则配置（JSON格式）
	IsDefault   bool      `json:"is_default" db:"is_default"`     // 是否为默认规则
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

/**
 * GeneralRuleConfig 通用密码规则配置
 * @author 陈凤庆
 * @description 通用密码规则的具体配置参数
 */
type GeneralRuleConfig struct {
	IncludeUppercase     bool   `json:"include_uppercase"`      // 包含大写字母
	IncludeLowercase     bool   `json:"include_lowercase"`      // 包含小写字母
	IncludeNumbers       bool   `json:"include_numbers"`        // 包含数字
	IncludeSpecialChars  bool   `json:"include_special_chars"`  // 包含特殊字符
	IncludeCustomChars   bool   `json:"include_custom_chars"`   // 包含自定义特殊字符
	MinUppercase         int    `json:"min_uppercase"`          // 大写字母最小位数
	MinLowercase         int    `json:"min_lowercase"`          // 小写字母最小位数
	MinNumbers           int    `json:"min_numbers"`            // 数字最小位数
	MinSpecialChars      int    `json:"min_special_chars"`      // 特殊字符最小位数
	MinCustomChars       int    `json:"min_custom_chars"`       // 自定义特殊字符最小位数
	Length               int    `json:"length"`                 // 密码总长度
	CustomSpecialChars   string `json:"custom_special_chars"`   // 自定义特殊字符集
}

/**
 * CustomRuleConfig 自定义密码规则配置
 * @author 陈凤庆
 * @description 自定义密码规则的具体配置参数
 */
type CustomRuleConfig struct {
	Pattern     string `json:"pattern"`     // 自定义规则模式
	Description string `json:"description"` // 规则说明
}

/**
 * LogConfig 日志配置模型
 * @author 陈凤庆
 * @date 20251003
 * @description 日志系统的配置选项
 */
type LogConfig struct {
	EnableInfoLog  bool `json:"enable_info_log"`  // 是否启用Info日志
	EnableDebugLog bool `json:"enable_debug_log"` // 是否启用Debug日志
}

/**
 * LockConfig 锁定配置模型
 * @author 陈凤庆
 * @date 20251004
 * @description 定时锁定功能的配置选项
 */
type LockConfig struct {
	EnableAutoLock     bool `json:"enable_auto_lock"`     // 是否启用自动锁定
	EnableTimerLock    bool `json:"enable_timer_lock"`    // 是否启用定时锁定
	EnableMinimizeLock bool `json:"enable_minimize_lock"` // 是否启用最小化锁定
	LockTimeMinutes    int  `json:"lock_time_minutes"`    // 定时锁定时间（分钟）
	EnableSystemLock   bool `json:"enable_system_lock"`   // 是否启用系统固定锁定（固定为true）
	SystemLockMinutes  int  `json:"system_lock_minutes"`  // 系统锁定时间（固定120分钟）
}

/**
 * HotkeyConfig 快捷键配置模型
 * @author 陈凤庆
 * @date 20251014
 * @description 全局快捷键功能的配置选项
 */
type HotkeyConfig struct {
	EnableGlobalHotkey   bool   `json:"enable_global_hotkey"`    // 是否启用全局快捷键（全局开关）
	ShowHideHotkey       string `json:"show_hide_hotkey"`        // 显示/隐藏窗口的快捷键组合，默认"Ctrl+Alt+H"
	EnableShowHideHotkey bool   `json:"enable_show_hide_hotkey"` // 是否启用显示/隐藏快捷键（独立开关）
}

/**
 * AppConfig 应用配置模型
 * @modify 20251003 陈凤庆 添加日志配置
 * @modify 20251004 陈凤庆 添加锁定配置
 * @modify 20251014 陈凤庆 添加快捷键配置
 */
type AppConfig struct {
	CurrentVaultPath string       `json:"current_vault_path"`
	RecentVaults     []string     `json:"recent_vaults"`
	WindowWidth      int          `json:"window_width"`
	WindowHeight     int          `json:"window_height"`
	Theme            string       `json:"theme"`
	Language         string       `json:"language"`
	LogConfig        LogConfig    `json:"log_config"`    // 20251003 陈凤庆 添加日志配置
	LockConfig       LockConfig   `json:"lock_config"`   // 20251004 陈凤庆 添加锁定配置
	HotkeyConfig     HotkeyConfig `json:"hotkey_config"` // 20251014 陈凤庆 添加快捷键配置
}

/**
 * SearchResult 搜索结果模型
 * @modify 20251001 陈凤庆 PasswordItemDecrypted改为AccountDecrypted
 */
type SearchResult struct {
	Type     string             `json:"type"` // group 或 account
	Groups   []Group            `json:"groups"`
	Accounts []AccountDecrypted `json:"accounts"`
}
