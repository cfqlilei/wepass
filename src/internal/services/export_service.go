package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"wepassword/internal/crypto"
	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/models"

	"github.com/alexmullins/zip"
)

/**
 * 导出服务
 * @author 陈凤庆
 * @date 2025-10-03
 * @description 处理密码库导出功能，包括数据查询、加密解密、JSON序列化和ZIP压缩
 */

/**
 * ExportService 导出服务
 */
type ExportService struct {
	dbManager      *database.DatabaseManager
	accountService *AccountService
	groupService   *GroupService
	typeService    *TypeService
}

/**
 * ExportData 导出数据结构
 */
type ExportData struct {
	Version  string          `json:"version"`   // 导出格式版本
	ExportAt time.Time       `json:"export_at"` // 导出时间
	Groups   []models.Group  `json:"groups"`    // 分组数据
	Types    []models.Type   `json:"types"`     // 类型数据
	Accounts []ExportAccount `json:"accounts"`  // 账号数据（用备份密码加密）
	Metadata ExportMetadata  `json:"metadata"`  // 导出元数据
}

/**
 * ExportAccount 导出账号结构（用备份密码加密）
 */
type ExportAccount struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`    // 标题不加密
	Username    string    `json:"username"` // 用备份密码加密
	Password    string    `json:"password"` // 用备份密码加密
	URL         string    `json:"url"`      // 用备份密码加密
	TypeID      string    `json:"typeid"`   // 类型ID不加密
	Notes       string    `json:"notes"`    // 用备份密码加密
	Icon        string    `json:"icon"`
	IsFavorite  bool      `json:"is_favorite"`
	UseCount    int       `json:"use_count"`
	LastUsedAt  time.Time `json:"last_used_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	InputMethod int       `json:"input_method"`
}

/**
 * ExportMetadata 导出元数据
 */
type ExportMetadata struct {
	TotalAccounts int      `json:"total_accounts"` // 导出账号总数
	TotalGroups   int      `json:"total_groups"`   // 导出分组总数
	TotalTypes    int      `json:"total_types"`    // 导出类型总数
	AccountIDs    []string `json:"account_ids"`    // 导出的账号ID列表
	BackupSalt    string   `json:"backup_salt"`    // 备份密码盐值（Base64编码）
}

/**
 * ExportOptions 导出选项
 */
type ExportOptions struct {
	LoginPassword  string   `json:"login_password"`  // 登录密码
	BackupPassword string   `json:"backup_password"` // 备份密码
	ExportPath     string   `json:"export_path"`     // 导出路径
	AccountIDs     []string `json:"account_ids"`     // 要导出的账号ID列表
	GroupIDs       []string `json:"group_ids"`       // 按分组过滤（可选）
	TypeIDs        []string `json:"type_ids"`        // 按类型过滤（可选）
	ExportAll      bool     `json:"export_all"`      // 是否导出所有账号
}

/**
 * NewExportService 创建导出服务
 * @param dbManager 数据库管理器
 * @param accountService 账号服务
 * @param groupService 分组服务
 * @param typeService 类型服务
 * @return *ExportService 导出服务实例
 */
func NewExportService(dbManager *database.DatabaseManager, accountService *AccountService, groupService *GroupService, typeService *TypeService) *ExportService {
	return &ExportService{
		dbManager:      dbManager,
		accountService: accountService,
		groupService:   groupService,
		typeService:    typeService,
	}
}

/**
 * ExportVault 导出密码库
 * @param options 导出选项
 * @return error 错误信息
 */
func (es *ExportService) ExportVault(options ExportOptions) error {
	logger.Info("[导出] 开始导出密码库，导出路径: %s", options.ExportPath)

	// 1. 验证登录密码
	if err := es.verifyLoginPassword(options.LoginPassword); err != nil {
		return fmt.Errorf("登录密码验证失败: %w", err)
	}
	logger.Info("[导出] ✅ 登录密码验证成功")

	// 2. 获取要导出的账号列表
	accounts, err := es.getAccountsToExport(options)
	if err != nil {
		return fmt.Errorf("获取导出账号列表失败: %w", err)
	}
	logger.Info("[导出] 获取到 %d 个账号需要导出", len(accounts))

	// 3. 获取相关的分组和类型
	groups, types, err := es.getRelatedGroupsAndTypes(accounts)
	if err != nil {
		return fmt.Errorf("获取相关分组和类型失败: %w", err)
	}
	logger.Info("[导出] 获取到 %d 个分组，%d 个类型", len(groups), len(types))

	// 4. 创建备份密码的加密管理器
	backupCrypto, backupSalt, err := es.createBackupCryptoManager(options.BackupPassword)
	if err != nil {
		return fmt.Errorf("创建备份加密管理器失败: %w", err)
	}

	// 5. 转换账号数据（用备份密码重新加密）
	exportAccounts, err := es.convertAccountsForExport(accounts, backupCrypto)
	if err != nil {
		return fmt.Errorf("转换账号数据失败: %w", err)
	}

	// 6. 构建导出数据
	exportData := ExportData{
		Version:  "1.0",
		ExportAt: time.Now(),
		Groups:   groups,
		Types:    types,
		Accounts: exportAccounts,
		Metadata: ExportMetadata{
			TotalAccounts: len(exportAccounts),
			TotalGroups:   len(groups),
			TotalTypes:    len(types),
			AccountIDs:    es.extractAccountIDs(accounts),
			BackupSalt:    base64.StdEncoding.EncodeToString(backupSalt), // 保存备份密码盐值
		},
	}

	// 7. 创建临时目录
	tempDir, err := os.MkdirTemp("", "wepass_export_*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir) // 清理临时目录

	// 8. 导出JSON文件
	if err := es.exportJSONFiles(tempDir, exportData); err != nil {
		return fmt.Errorf("导出JSON文件失败: %w", err)
	}

	// 9. 创建ZIP压缩包
	if err := es.createZipArchive(tempDir, options.ExportPath, options.BackupPassword); err != nil {
		return fmt.Errorf("创建ZIP压缩包失败: %w", err)
	}

	logger.Info("[导出] 🎉 密码库导出完成: %s", options.ExportPath)
	return nil
}

/**
 * verifyLoginPassword 验证登录密码
 * @param loginPassword 登录密码
 * @return error 错误信息
 */
func (es *ExportService) verifyLoginPassword(loginPassword string) error {
	// 获取密码库配置
	vaultConfig, err := es.dbManager.GetVaultConfig()
	if err != nil {
		return fmt.Errorf("获取密码库配置失败: %w", err)
	}

	if vaultConfig == nil {
		return fmt.Errorf("密码库配置不存在")
	}

	// 创建临时加密管理器验证密码
	tempCrypto := crypto.NewCryptoManager()

	// 解码盐值
	salt, err := base64.StdEncoding.DecodeString(vaultConfig.Salt)
	if err != nil {
		return fmt.Errorf("解码盐值失败: %w", err)
	}

	// 验证密码
	if !tempCrypto.VerifyPassword(loginPassword, vaultConfig.PasswordHash, salt) {
		return fmt.Errorf("登录密码不正确")
	}

	return nil
}

/**
 * getAccountsToExport 获取要导出的账号列表
 * @param options 导出选项
 * @return []models.AccountDecrypted 账号列表
 * @return error 错误信息
 */
func (es *ExportService) getAccountsToExport(options ExportOptions) ([]models.AccountDecrypted, error) {
	if options.ExportAll {
		// 导出所有账号
		logger.Info("[导出] 导出所有账号")
		return es.accountService.GetAllAccounts()
	}

	// 按分组导出
	if len(options.GroupIDs) > 0 {
		logger.Info("[导出] 按分组导出，分组数量: %d", len(options.GroupIDs))
		return es.getAccountsByGroups(options.GroupIDs)
	}

	// 按类别导出
	if len(options.TypeIDs) > 0 {
		logger.Info("[导出] 按类别导出，类别数量: %d", len(options.TypeIDs))
		return es.getAccountsByTypes(options.TypeIDs)
	}

	// 手动选择账号导出
	if len(options.AccountIDs) > 0 {
		logger.Info("[导出] 手动选择账号导出，账号数量: %d", len(options.AccountIDs))
		return es.getAccountsByIDs(options.AccountIDs)
	}

	// 如果没有指定任何选择条件，返回空列表
	logger.Info("[导出] 没有指定任何导出条件，返回空列表")
	return []models.AccountDecrypted{}, nil
}

/**
 * getAccountsByGroups 根据分组ID列表获取账号
 * @param groupIDs 分组ID列表
 * @return []models.AccountDecrypted 账号列表
 * @return error 错误信息
 */
func (es *ExportService) getAccountsByGroups(groupIDs []string) ([]models.AccountDecrypted, error) {
	var allAccounts []models.AccountDecrypted

	for _, groupID := range groupIDs {
		logger.Info("[导出] 获取分组账号，分组ID: %s", groupID)
		accounts, err := es.accountService.GetAccountsByGroup(groupID)
		if err != nil {
			logger.Error("[导出] 获取分组账号失败，分组ID: %s, 错误: %v", groupID, err)
			continue // 跳过获取失败的分组，继续处理其他分组
		}
		allAccounts = append(allAccounts, accounts...)
	}

	// 去重处理（防止同一账号在多个分组中重复）
	return es.deduplicateAccounts(allAccounts), nil
}

/**
 * getAccountsByTypes 根据类别ID列表获取账号
 * @param typeIDs 类别ID列表
 * @return []models.AccountDecrypted 账号列表
 * @return error 错误信息
 */
func (es *ExportService) getAccountsByTypes(typeIDs []string) ([]models.AccountDecrypted, error) {
	var allAccounts []models.AccountDecrypted

	for _, typeID := range typeIDs {
		logger.Info("[导出] 获取类别账号，类别ID: %s", typeID)
		// 构建查询条件
		conditions := fmt.Sprintf(`{"type_id":"%s"}`, typeID)
		accounts, err := es.accountService.GetAccountsByConditions(conditions)
		if err != nil {
			logger.Error("[导出] 获取类别账号失败，类别ID: %s, 错误: %v", typeID, err)
			continue // 跳过获取失败的类别，继续处理其他类别
		}
		allAccounts = append(allAccounts, accounts...)
	}

	// 去重处理（防止同一账号在多个类别中重复）
	return es.deduplicateAccounts(allAccounts), nil
}

/**
 * getAccountsByIDs 根据账号ID列表获取账号
 * @param accountIDs 账号ID列表
 * @return []models.AccountDecrypted 账号列表
 * @return error 错误信息
 */
func (es *ExportService) getAccountsByIDs(accountIDs []string) ([]models.AccountDecrypted, error) {
	var accounts []models.AccountDecrypted

	for _, accountID := range accountIDs {
		account, err := es.accountService.GetAccountByID(accountID)
		if err != nil {
			logger.Error("[导出] 获取账号失败，ID: %s, 错误: %v", accountID, err)
			continue // 跳过获取失败的账号
		}
		if account != nil {
			accounts = append(accounts, *account)
		}
	}

	return accounts, nil
}

/**
 * deduplicateAccounts 账号去重
 * @param accounts 账号列表
 * @return []models.AccountDecrypted 去重后的账号列表
 */
func (es *ExportService) deduplicateAccounts(accounts []models.AccountDecrypted) []models.AccountDecrypted {
	seen := make(map[string]bool)
	var result []models.AccountDecrypted

	for _, account := range accounts {
		if !seen[account.ID] {
			seen[account.ID] = true
			result = append(result, account)
		}
	}

	return result
}

/**
 * getRelatedGroupsAndTypes 获取相关的分组和类型
 * @param accounts 账号列表
 * @return []models.Group 分组列表
 * @return []models.Type 类型列表
 * @return error 错误信息
 */
func (es *ExportService) getRelatedGroupsAndTypes(accounts []models.AccountDecrypted) ([]models.Group, []models.Type, error) {
	// 收集唯一的类型ID
	typeIDSet := make(map[string]bool)
	for _, account := range accounts {
		typeIDSet[account.TypeID] = true
	}

	// 获取类型信息
	var types []models.Type
	for typeID := range typeIDSet {
		typeInfo, err := es.typeService.GetTypeByID(typeID)
		if err != nil {
			logger.Error("[导出] 获取类型失败，ID: %s, 错误: %v", typeID, err)
			continue
		}
		if typeInfo != nil {
			types = append(types, *typeInfo)
		}
	}

	// 收集唯一的分组ID
	groupIDSet := make(map[string]bool)
	for _, typeInfo := range types {
		groupIDSet[typeInfo.GroupID] = true
	}

	// 获取分组信息
	var groups []models.Group
	for groupID := range groupIDSet {
		group, err := es.groupService.GetGroupByID(groupID)
		if err != nil {
			logger.Error("[导出] 获取分组失败，ID: %s, 错误: %v", groupID, err)
			continue
		}
		if group != nil {
			groups = append(groups, *group)
		}
	}

	return groups, types, nil
}

/**
 * createBackupCryptoManager 创建备份密码的加密管理器
 * @param backupPassword 备份密码
 * @return *crypto.CryptoManager 加密管理器
 * @return []byte 盐值
 * @return error 错误信息
 */
func (es *ExportService) createBackupCryptoManager(backupPassword string) (*crypto.CryptoManager, []byte, error) {
	backupCrypto := crypto.NewCryptoManager()

	// 生成新的盐值用于备份密码
	salt, err := backupCrypto.GenerateSalt()
	if err != nil {
		return nil, nil, fmt.Errorf("生成备份密码盐值失败: %w", err)
	}

	// 设置备份密码
	backupCrypto.SetMasterPassword(backupPassword, salt)

	return backupCrypto, salt, nil
}

/**
 * convertAccountsForExport 转换账号数据用于导出（用备份密码重新加密）
 * @param accounts 原始账号列表
 * @param backupCrypto 备份密码加密管理器
 * @return []ExportAccount 导出账号列表
 * @return error 错误信息
 */
func (es *ExportService) convertAccountsForExport(accounts []models.AccountDecrypted, backupCrypto *crypto.CryptoManager) ([]ExportAccount, error) {
	var exportAccounts []ExportAccount

	for _, account := range accounts {
		// 用备份密码重新加密敏感字段
		encryptedUsername, err := backupCrypto.Encrypt(account.Username)
		if err != nil {
			return nil, fmt.Errorf("加密用户名失败，账号ID: %s, 错误: %w", account.ID, err)
		}

		encryptedPassword, err := backupCrypto.Encrypt(account.Password)
		if err != nil {
			return nil, fmt.Errorf("加密密码失败，账号ID: %s, 错误: %w", account.ID, err)
		}

		encryptedURL, err := backupCrypto.Encrypt(account.URL)
		if err != nil {
			return nil, fmt.Errorf("加密地址失败，账号ID: %s, 错误: %w", account.ID, err)
		}

		encryptedNotes, err := backupCrypto.Encrypt(account.Notes)
		if err != nil {
			return nil, fmt.Errorf("加密备注失败，账号ID: %s, 错误: %w", account.ID, err)
		}

		exportAccount := ExportAccount{
			ID:          account.ID,
			Title:       account.Title, // 标题不加密
			Username:    encryptedUsername,
			Password:    encryptedPassword,
			URL:         encryptedURL,
			TypeID:      account.TypeID, // 类型ID不加密
			Notes:       encryptedNotes,
			Icon:        account.Icon,
			IsFavorite:  account.IsFavorite,
			UseCount:    account.UseCount,
			LastUsedAt:  account.LastUsedAt,
			CreatedAt:   account.CreatedAt,
			UpdatedAt:   account.UpdatedAt,
			InputMethod: account.InputMethod,
		}

		exportAccounts = append(exportAccounts, exportAccount)
	}

	return exportAccounts, nil
}

/**
 * extractAccountIDs 提取账号ID列表
 * @param accounts 账号列表
 * @return []string 账号ID列表
 */
func (es *ExportService) extractAccountIDs(accounts []models.AccountDecrypted) []string {
	var accountIDs []string
	for _, account := range accounts {
		accountIDs = append(accountIDs, account.ID)
	}
	return accountIDs
}

/**
 * exportJSONFiles 导出JSON文件到临时目录
 * @param tempDir 临时目录
 * @param exportData 导出数据
 * @return error 错误信息
 */
func (es *ExportService) exportJSONFiles(tempDir string, exportData ExportData) error {
	// 导出分组数据
	if err := es.writeJSONFile(filepath.Join(tempDir, "groups.json"), exportData.Groups); err != nil {
		return fmt.Errorf("导出分组数据失败: %w", err)
	}

	// 导出类型数据
	if err := es.writeJSONFile(filepath.Join(tempDir, "types.json"), exportData.Types); err != nil {
		return fmt.Errorf("导出类型数据失败: %w", err)
	}

	// 导出账号数据
	if err := es.writeJSONFile(filepath.Join(tempDir, "accounts.json"), exportData.Accounts); err != nil {
		return fmt.Errorf("导出账号数据失败: %w", err)
	}

	// 导出元数据
	if err := es.writeJSONFile(filepath.Join(tempDir, "metadata.json"), exportData.Metadata); err != nil {
		return fmt.Errorf("导出元数据失败: %w", err)
	}

	// 导出完整数据（包含所有信息的单个文件）
	if err := es.writeJSONFile(filepath.Join(tempDir, "export.json"), exportData); err != nil {
		return fmt.Errorf("导出完整数据失败: %w", err)
	}

	return nil
}

/**
 * writeJSONFile 写入JSON文件
 * @param filePath 文件路径
 * @param data 数据
 * @return error 错误信息
 */
func (es *ExportService) writeJSONFile(filePath string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

/**
 * createZipArchive 创建带密码保护的ZIP压缩包
 * @param sourceDir 源目录
 * @param zipPath ZIP文件路径
 * @param password 压缩密码
 * @return error 错误信息
 */
func (es *ExportService) createZipArchive(sourceDir, zipPath, password string) error {
	logger.Info("[导出] 创建加密ZIP压缩包: %s", zipPath)

	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("创建ZIP文件失败: %w", err)
	}
	defer zipFile.Close()

	// 创建支持密码保护的ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录中的所有文件
	return filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return fmt.Errorf("计算相对路径失败: %w", err)
		}

		logger.Info("[导出] 添加文件到ZIP: %s", relPath)

		// 在ZIP中创建带密码保护的文件
		zipFileWriter, err := zipWriter.Encrypt(relPath, password)
		if err != nil {
			return fmt.Errorf("在ZIP中创建加密文件失败: %w", err)
		}

		// 读取源文件
		sourceFile, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("打开源文件失败: %w", err)
		}
		defer sourceFile.Close()

		// 复制文件内容到ZIP
		_, err = io.Copy(zipFileWriter, sourceFile)
		if err != nil {
			return fmt.Errorf("复制文件内容失败: %w", err)
		}

		// 刷新写入器
		if flusher, ok := zipFileWriter.(interface{ Flush() error }); ok {
			if err := flusher.Flush(); err != nil {
				return fmt.Errorf("刷新ZIP写入器失败: %w", err)
			}
		}

		return nil
	})
}
