package services

import (
	"encoding/base64"
	"fmt"
	"os"

	"wepassword/internal/config"
	"wepassword/internal/crypto"
	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/models"
)

/**
 * 密码库服务
 * @author 陈凤庆
 * @description 管理密码库的创建、打开、验证等操作
 */

/**
 * VaultService 密码库服务
 */
type VaultService struct {
	dbManager       *database.DatabaseManager
	configManager   *config.ConfigManager
	cryptoManager   *crypto.CryptoManager
	isOpened        bool   // 20251003 陈凤庆 密码库是否已打开
	currentPath     string // 20251003 陈凤庆 当前打开的密码库路径
	currentPassword string // 20251017 陈凤庆 当前登录密码，存储在内存中
}

/**
 * NewVaultService 创建新的密码库服务
 * @param dbManager 数据库管理器
 * @param configManager 配置管理器
 * @return *VaultService 密码库服务实例
 */
func NewVaultService(dbManager *database.DatabaseManager, configManager *config.ConfigManager) *VaultService {
	return &VaultService{
		dbManager:     dbManager,
		configManager: configManager,
		cryptoManager: crypto.NewCryptoManager(),
	}
}

/**
 * CheckVaultExists 检查密码库文件是否存在
 * @param vaultPath 密码库文件路径
 * @return bool 是否存在
 */
func (vs *VaultService) CheckVaultExists(vaultPath string) bool {
	if vaultPath == "" {
		return false
	}

	_, err := os.Stat(vaultPath)
	return !os.IsNotExist(err)
}

/**
 * CreateVault 创建新密码库
 * @param vaultPath 密码库文件路径
 * @param password 登录密码
 * @param language 语言代码（用于初始化多语言数据）
 * @return error 错误信息
 * @modify 20251005 陈凤庆 添加语言参数，支持多语言初始数据
 */
func (vs *VaultService) CreateVault(vaultPath string, password string, language string) error {
	// 检查文件是否已存在
	if vs.CheckVaultExists(vaultPath) {
		return fmt.Errorf("密码库文件已存在: %s", vaultPath)
	}

	// 打开数据库连接
	if err := vs.dbManager.OpenDatabase(vaultPath); err != nil {
		return fmt.Errorf("创建数据库失败: %w", err)
	}

	// 创建数据库表，传递语言参数用于初始化多语言数据
	if err := vs.dbManager.CreateTables(language); err != nil {
		return fmt.Errorf("创建数据库表失败: %w", err)
	}

	// 生成盐值
	salt, err := vs.cryptoManager.GenerateSalt()
	if err != nil {
		return fmt.Errorf("生成盐值失败: %w", err)
	}

	// 哈希密码
	hashedPassword := vs.cryptoManager.HashPassword(password, salt)

	// 保存密码库配置
	vaultConfig := &models.VaultConfig{
		PasswordHash: hashedPassword,
		Salt:         base64.StdEncoding.EncodeToString(salt),
	}

	if err := vs.dbManager.SaveVaultConfig(vaultConfig); err != nil {
		return fmt.Errorf("保存密码库配置失败: %w", err)
	}

	// 设置主密钥
	vs.cryptoManager.SetMasterPassword(password, salt)

	// 更新配置文件
	if err := vs.configManager.SetCurrentVaultPath(vaultPath); err != nil {
		return fmt.Errorf("更新配置文件失败: %w", err)
	}

	// 20251003 陈凤庆 标记密码库已打开
	vs.isOpened = true
	vs.currentPath = vaultPath
	logger.Info("[密码库] 密码库已创建并打开: %s", vaultPath)

	return nil
}

/**
 * OpenVault 打开密码库
 * @param vaultPath 密码库文件路径
 * @param password 登录密码
 * @return error 错误信息
 * @description 20251001 陈凤庆 添加数据库升级检查,确保打开旧数据库时自动升级
 * @modify 20250101 陈凤庆 添加数据完整性检查,确保有默认分组和标签
 * @modify 20251003 陈凤庆 添加详细的登录日志，便于跟踪Windows平台登录问题
 */
func (vs *VaultService) OpenVault(vaultPath string, password string) error {
	logger.Info("[登录] 开始打开密码库: %s", vaultPath)

	// 检查文件是否存在
	if !vs.CheckVaultExists(vaultPath) {
		logger.Error("[登录] ❌ 密码库文件不存在: %s", vaultPath)
		return fmt.Errorf("密码库文件不存在: %s", vaultPath)
	}
	logger.Info("[登录] ✅ 密码库文件存在")

	// 打开数据库连接
	logger.Info("[登录] 正在打开数据库连接...")
	if err := vs.dbManager.OpenDatabase(vaultPath); err != nil {
		logger.Error("[登录] ❌ 打开数据库失败: %v", err)
		return fmt.Errorf("打开数据库失败: %w", err)
	}
	logger.Info("[登录] ✅ 数据库连接已建立")

	// 20251001 陈凤庆 检查并升级数据库结构
	// 这一步会自动检测数据库版本并执行必要的升级
	logger.Info("[登录] 正在检查数据库结构...")
	// 使用默认中文语言，因为这是打开现有密码库，不是创建新的
	if err := vs.dbManager.CreateTables("zh-CN"); err != nil {
		logger.Error("[登录] ❌ 检查数据库结构失败: %v", err)
		return fmt.Errorf("检查数据库结构失败: %w", err)
	}
	logger.Info("[登录] ✅ 数据库结构检查完成")

	// 20250101 陈凤庆 检查数据完整性，确保有默认分组和标签
	logger.Info("[登录] 正在检查数据完整性...")
	if err := vs.dbManager.EnsureDataIntegrity(); err != nil {
		logger.Error("[登录] ❌ 数据完整性检查失败: %v", err)
		return fmt.Errorf("数据完整性检查失败: %w", err)
	}
	logger.Info("[登录] ✅ 数据完整性检查完成")

	// 20251017 陈凤庆 取消自动初始化默认密码规则，让用户自行决定是否创建密码规则
	// logger.Info("[登录] 正在检查默认密码规则...")
	// if err := vs.dbManager.InitializeDefaultPasswordRules(); err != nil {
	//     logger.Error("[登录] ❌ 初始化默认密码规则失败: %v", err)
	//     return fmt.Errorf("初始化默认密码规则失败: %w", err)
	// }
	// logger.Info("[登录] ✅ 默认密码规则检查完成")

	// 获取密码库配置
	logger.Info("[登录] 正在获取密码库配置...")
	vaultConfig, err := vs.dbManager.GetVaultConfig()
	if err != nil {
		logger.Error("[登录] ❌ 获取密码库配置失败: %v", err)
		return fmt.Errorf("获取密码库配置失败: %w", err)
	}

	if vaultConfig == nil {
		logger.Error("[登录] ❌ 密码库配置不存在，可能不是有效的密码库文件")
		return fmt.Errorf("密码库配置不存在，可能不是有效的密码库文件")
	}
	logger.Info("[登录] ✅ 密码库配置获取成功")

	// 解码盐值
	logger.Info("[登录] 正在解码盐值...")
	salt, err := base64.StdEncoding.DecodeString(vaultConfig.Salt)
	if err != nil {
		logger.Error("[登录] ❌ 解码盐值失败: %v", err)
		return fmt.Errorf("解码盐值失败: %w", err)
	}
	logger.Info("[登录] ✅ 盐值解码成功")

	// 验证密码
	logger.Info("[登录] 正在验证登录密码...")
	if !vs.cryptoManager.VerifyPassword(password, vaultConfig.PasswordHash, salt) {
		logger.Error("[登录] ❌ 登录密码不正确")
		return fmt.Errorf("登录密码不正确")
	}
	logger.Info("[登录] ✅ 登录密码验证成功")

	// 设置主密钥
	logger.Info("[登录] 正在设置主密钥...")
	vs.cryptoManager.SetMasterPassword(password, salt)
	logger.Info("[登录] ✅ 主密钥设置完成")

	// 更新配置文件
	logger.Info("[登录] 正在更新配置文件...")
	if err := vs.configManager.SetCurrentVaultPath(vaultPath); err != nil {
		logger.Error("[登录] ❌ 更新配置文件失败: %v", err)
		return fmt.Errorf("更新配置文件失败: %w", err)
	}
	logger.Info("[登录] ✅ 配置文件更新完成")

	// 20251003 陈凤庆 标记密码库已打开
	vs.isOpened = true
	vs.currentPath = vaultPath
	// 20251017 陈凤庆 保存当前登录密码到内存
	vs.currentPassword = password
	logger.Info("[登录] 🎉 密码库登录成功: %s", vaultPath)
	logger.Info("[登录] 密码库状态: isOpened=%t, currentPath=%s", vs.isOpened, vs.currentPath)

	return nil
}

/**
 * GetCryptoManager 获取加密管理器
 * @return *crypto.CryptoManager 加密管理器
 */
func (vs *VaultService) GetCryptoManager() *crypto.CryptoManager {
	return vs.cryptoManager
}

/**
 * IsOpened 检查密码库是否已打开
 * @return bool 是否已打开
 * @description 20251003 陈凤庆 用于前端检查登录状态
 */
func (vs *VaultService) IsOpened() bool {
	return vs.isOpened && vs.dbManager.IsOpened()
}

/**
 * IsVaultOpened 检查密码库是否已打开（兼容旧方法名）
 * @return bool 是否已打开
 */
func (vs *VaultService) IsVaultOpened() bool {
	return vs.IsOpened()
}

/**
 * GetCurrentVaultPath 获取当前打开的密码库路径
 * @return string 密码库路径
 * @description 20251003 陈凤庆 返回当前打开的密码库文件路径
 */
func (vs *VaultService) GetCurrentVaultPath() string {
	return vs.currentPath
}

/**
 * CloseVault 关闭密码库
 * @description 20251003 陈凤庆 关闭密码库并清理状态
 */
func (vs *VaultService) CloseVault() {
	vs.dbManager.Close()
	vs.cryptoManager = crypto.NewCryptoManager()
	// 20251003 陈凤庆 清理状态
	vs.isOpened = false
	vs.currentPath = ""
	// 20251017 陈凤庆 清理内存中的密码
	vs.currentPassword = ""
	logger.Info("[密码库] 密码库已关闭")
}

/**
 * SetCurrentPassword 设置当前登录密码
 * @param password 登录密码
 * @description 20251017 陈凤庆 在内存中存储当前登录密码，用于用户名历史记录加密
 */
func (vs *VaultService) SetCurrentPassword(password string) {
	vs.currentPassword = password
	logger.Info("[密码库] 当前登录密码已设置")
}

/**
 * GetCurrentPassword 获取当前登录密码
 * @return string 当前登录密码
 * @description 20251017 陈凤庆 从内存中获取当前登录密码
 */
func (vs *VaultService) GetCurrentPassword() string {
	return vs.currentPassword
}

/**
 * ClearCurrentPassword 清除当前登录密码
 * @description 20251017 陈凤庆 清除内存中的当前登录密码
 */
func (vs *VaultService) ClearCurrentPassword() {
	vs.currentPassword = ""
	logger.Info("[密码库] 当前登录密码已清除")
}
