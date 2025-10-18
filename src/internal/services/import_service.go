package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"wepassword/internal/crypto"
	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/models"

	"github.com/alexmullins/zip"
)

/**
 * 导入服务
 * @author 陈凤庆
 * @date 2025-10-03
 * @description 处理密码库导入功能，包括ZIP解压、JSON反序列化、数据验证和导入
 */

/**
 * ImportService 导入服务
 */
type ImportService struct {
	dbManager      *database.DatabaseManager
	accountService *AccountService
	groupService   *GroupService
	typeService    *TypeService
	cryptoManager  *crypto.CryptoManager // 当前密码库的加密管理器
}

/**
 * ImportOptions 导入选项
 */
type ImportOptions struct {
	ImportPath     string `json:"import_path"`     // 导入文件路径
	BackupPassword string `json:"backup_password"` // 备份密码（解压密码）
}

/**
 * ImportResult 导入结果
 */
type ImportResult struct {
	Success               bool                 `json:"success"`                 // 是否成功
	TotalAccounts         int                  `json:"total_accounts"`          // 总账号数
	ImportedAccounts      int                  `json:"imported_accounts"`       // 成功导入的账号数
	SkippedAccounts       int                  `json:"skipped_accounts"`        // 跳过的账号数
	ErrorAccounts         int                  `json:"error_accounts"`          // 错误的账号数
	TotalGroups           int                  `json:"total_groups"`            // 总分组数
	ImportedGroups        int                  `json:"imported_groups"`         // 成功导入的分组数
	SkippedGroups         int                  `json:"skipped_groups"`          // 跳过的分组数
	TotalTypes            int                  `json:"total_types"`             // 总类型数
	ImportedTypes         int                  `json:"imported_types"`          // 成功导入的类型数
	SkippedTypes          int                  `json:"skipped_types"`           // 跳过的类型数
	SkippedAccountDetails []SkippedAccountInfo `json:"skipped_account_details"` // 跳过的账号详情
	ErrorMessage          string               `json:"error_message"`           // 错误信息
}

/**
 * SkippedAccountInfo 跳过的账号信息
 */
type SkippedAccountInfo struct {
	ID    string `json:"id"`    // 账号ID
	Title string `json:"title"` // 账号标题
	Name  string `json:"name"`  // 账号名称（用户名）
}

/**
 * NewImportService 创建导入服务
 * @param dbManager 数据库管理器
 * @param accountService 账号服务
 * @param groupService 分组服务
 * @param typeService 类型服务
 * @param cryptoManager 当前密码库的加密管理器
 * @return *ImportService 导入服务实例
 */
func NewImportService(dbManager *database.DatabaseManager, accountService *AccountService, groupService *GroupService, typeService *TypeService, cryptoManager *crypto.CryptoManager) *ImportService {
	return &ImportService{
		dbManager:      dbManager,
		accountService: accountService,
		groupService:   groupService,
		typeService:    typeService,
		cryptoManager:  cryptoManager,
	}
}

/**
 * ImportVault 导入密码库
 * @param options 导入选项
 * @return ImportResult 导入结果
 * @return error 错误信息
 */
func (is *ImportService) ImportVault(options ImportOptions) (ImportResult, error) {
	logger.Info("[导入] 开始导入密码库，导入路径: %s", options.ImportPath)

	result := ImportResult{
		SkippedAccountDetails: make([]SkippedAccountInfo, 0),
	}

	// 1. 创建临时目录
	tempDir, err := os.MkdirTemp("", "wepass_import_*")
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("创建临时目录失败: %v", err)
		return result, fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir) // 清理临时目录

	// 2. 解压ZIP文件
	if err := is.extractZipFile(options.ImportPath, tempDir, options.BackupPassword); err != nil {
		result.ErrorMessage = fmt.Sprintf("解压ZIP文件失败: %v", err)
		return result, fmt.Errorf("解压ZIP文件失败: %w", err)
	}
	logger.Info("[导入] ✅ ZIP文件解压成功")

	// 3. 读取导出数据
	exportData, err := is.readExportData(tempDir)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("读取导出数据失败: %v", err)
		return result, fmt.Errorf("读取导出数据失败: %w", err)
	}
	logger.Info("[导入] ✅ 导出数据读取成功")

	// 4. 创建备份密码的加密管理器
	var backupCrypto *crypto.CryptoManager
	if exportData.Metadata.BackupSalt != "" {
		// 使用导出时保存的盐值（新版本导出文件）
		backupCrypto, err = is.createBackupCryptoManagerWithSalt(options.BackupPassword, exportData.Metadata.BackupSalt)
		if err != nil {
			result.ErrorMessage = fmt.Sprintf("创建备份加密管理器失败: %v", err)
			return result, fmt.Errorf("创建备份加密管理器失败: %w", err)
		}
		logger.Info("[导入] 使用导出文件中的盐值创建备份加密管理器")
	} else {
		// 回退到旧逻辑（旧版本导出文件）
		backupCrypto, err = is.createBackupCryptoManager(options.BackupPassword)
		if err != nil {
			result.ErrorMessage = fmt.Sprintf("创建备份加密管理器失败: %v", err)
			return result, fmt.Errorf("创建备份加密管理器失败: %w", err)
		}
		logger.Info("[导入] 导出文件无盐值信息，使用旧版本兼容模式")
	}

	// 5. 数据梳理：统一分组和分类ID
	err = is.normalizeImportData(&exportData)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("数据梳理失败: %v", err)
		return result, fmt.Errorf("数据梳理失败: %w", err)
	}
	logger.Info("[导入] ✅ 数据梳理完成")

	// 6. 导入分组数据
	result.TotalGroups = len(exportData.Groups)
	result.ImportedGroups, result.SkippedGroups = is.importGroups(exportData.Groups)
	logger.Info("[导入] 分组导入完成: 总数=%d, 导入=%d, 跳过=%d",
		result.TotalGroups, result.ImportedGroups, result.SkippedGroups)

	// 7. 导入类型数据
	result.TotalTypes = len(exportData.Types)
	result.ImportedTypes, result.SkippedTypes = is.importTypes(exportData.Types)
	logger.Info("[导入] 类型导入完成: 总数=%d, 导入=%d, 跳过=%d",
		result.TotalTypes, result.ImportedTypes, result.SkippedTypes)

	// 8. 导入账号数据
	result.TotalAccounts = len(exportData.Accounts)
	result.ImportedAccounts, result.SkippedAccounts, result.ErrorAccounts, result.SkippedAccountDetails =
		is.importAccounts(exportData.Accounts, backupCrypto)
	logger.Info("[导入] 账号导入完成: 总数=%d, 导入=%d, 跳过=%d, 错误=%d",
		result.TotalAccounts, result.ImportedAccounts, result.SkippedAccounts, result.ErrorAccounts)

	result.Success = true
	logger.Info("[导入] 🎉 密码库导入完成")
	return result, nil
}

/**
 * extractZipFile 解压带密码保护的ZIP文件
 * @param zipPath ZIP文件路径
 * @param destDir 目标目录
 * @param password 解压密码
 * @return error 错误信息
 */
func (is *ImportService) extractZipFile(zipPath, destDir, password string) error {
	logger.Info("[导入] 解压加密ZIP文件: %s", zipPath)

	// 打开ZIP文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer reader.Close()

	// 解压每个文件
	for _, file := range reader.File {
		// 构建目标文件路径
		destPath := filepath.Join(destDir, file.Name)

		// 确保目标目录存在
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("创建目标目录失败: %w", err)
		}

		// 跳过目录
		if file.FileInfo().IsDir() {
			continue
		}

		logger.Info("[导入] 解压文件: %s", file.Name)

		// 检查文件是否加密
		var fileReader io.ReadCloser
		if file.IsEncrypted() {
			// 使用密码打开加密文件
			file.SetPassword(password)
			fileReader, err = file.Open()
			if err != nil {
				return fmt.Errorf("使用密码打开加密文件失败 (%s): %w", file.Name, err)
			}
		} else {
			// 打开普通文件
			fileReader, err = file.Open()
			if err != nil {
				return fmt.Errorf("打开ZIP中的文件失败 (%s): %w", file.Name, err)
			}
		}

		// 创建目标文件
		destFile, err := os.Create(destPath)
		if err != nil {
			fileReader.Close()
			return fmt.Errorf("创建目标文件失败: %w", err)
		}

		// 复制文件内容
		_, err = io.Copy(destFile, fileReader)
		fileReader.Close()
		destFile.Close()

		if err != nil {
			return fmt.Errorf("复制文件内容失败 (%s): %w", file.Name, err)
		}
	}

	logger.Info("[导入] ZIP文件解压完成")
	return nil
}

/**
 * readExportData 读取导出数据
 * @param tempDir 临时目录
 * @return ExportData 导出数据
 * @return error 错误信息
 */
func (is *ImportService) readExportData(tempDir string) (ExportData, error) {
	var exportData ExportData

	// 尝试读取完整的导出文件
	exportFilePath := filepath.Join(tempDir, "export.json")
	if _, err := os.Stat(exportFilePath); err == nil {
		// 读取完整导出文件
		data, err := os.ReadFile(exportFilePath)
		if err != nil {
			return exportData, fmt.Errorf("读取导出文件失败: %w", err)
		}

		if err := json.Unmarshal(data, &exportData); err != nil {
			return exportData, fmt.Errorf("解析导出文件失败: %w", err)
		}

		return exportData, nil
	}

	// 如果没有完整导出文件，则分别读取各个文件
	// 读取分组数据
	if err := is.readJSONFile(filepath.Join(tempDir, "groups.json"), &exportData.Groups); err != nil {
		return exportData, fmt.Errorf("读取分组数据失败: %w", err)
	}

	// 读取类型数据
	if err := is.readJSONFile(filepath.Join(tempDir, "types.json"), &exportData.Types); err != nil {
		return exportData, fmt.Errorf("读取类型数据失败: %w", err)
	}

	// 读取账号数据
	if err := is.readJSONFile(filepath.Join(tempDir, "accounts.json"), &exportData.Accounts); err != nil {
		return exportData, fmt.Errorf("读取账号数据失败: %w", err)
	}

	// 读取元数据（可选）
	metadataPath := filepath.Join(tempDir, "metadata.json")
	if _, err := os.Stat(metadataPath); err == nil {
		if err := is.readJSONFile(metadataPath, &exportData.Metadata); err != nil {
			logger.Info("[导入] 读取元数据失败: %v", err)
		}
	}

	return exportData, nil
}

/**
 * readJSONFile 读取JSON文件
 * @param filePath 文件路径
 * @param target 目标对象
 * @return error 错误信息
 */
func (is *ImportService) readJSONFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	return nil
}

/**
 * createBackupCryptoManager 创建备份密码的加密管理器（已废弃，使用createBackupCryptoManagerWithSalt）
 * @param backupPassword 备份密码
 * @return *crypto.CryptoManager 加密管理器
 * @return error 错误信息
 */
func (is *ImportService) createBackupCryptoManager(backupPassword string) (*crypto.CryptoManager, error) {
	backupCrypto := crypto.NewCryptoManager()

	// 生成新的盐值用于备份密码
	salt, err := backupCrypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("生成备份密码盐值失败: %w", err)
	}

	// 设置备份密码
	backupCrypto.SetMasterPassword(backupPassword, salt)

	return backupCrypto, nil
}

/**
 * createBackupCryptoManagerWithSalt 使用指定盐值创建备份密码的加密管理器
 * @param backupPassword 备份密码
 * @param saltBase64 盐值（Base64编码）
 * @return *crypto.CryptoManager 加密管理器
 * @return error 错误信息
 */
func (is *ImportService) createBackupCryptoManagerWithSalt(backupPassword string, saltBase64 string) (*crypto.CryptoManager, error) {
	backupCrypto := crypto.NewCryptoManager()

	// 解码盐值
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return nil, fmt.Errorf("解码备份密码盐值失败: %w", err)
	}

	// 设置备份密码（使用导出时的盐值）
	backupCrypto.SetMasterPassword(backupPassword, salt)

	return backupCrypto, nil
}

/**
 * importGroups 导入分组数据
 * @param groups 分组列表
 * @return int 导入成功数量
 * @return int 跳过数量
 */
func (is *ImportService) importGroups(groups []models.Group) (int, int) {
	imported := 0
	skipped := 0

	for _, group := range groups {
		// 检查分组是否已存在
		existingGroup, err := is.groupService.GetGroupByID(group.ID)
		if err == nil && existingGroup != nil {
			logger.Info("[导入] 分组已存在，跳过: ID=%s, Name=%s", group.ID, group.Name)
			skipped++
			continue
		}

		// 创建分组（使用原有ID）
		if err := is.createGroupWithID(group); err != nil {
			logger.Error("[导入] 创建分组失败: ID=%s, Name=%s, 错误=%v", group.ID, group.Name, err)
			skipped++
			continue
		}

		logger.Info("[导入] 分组导入成功: ID=%s, Name=%s", group.ID, group.Name)
		imported++
	}

	return imported, skipped
}

/**
 * importTypes 导入类型数据
 * @param types 类型列表
 * @return int 导入成功数量
 * @return int 跳过数量
 */
func (is *ImportService) importTypes(types []models.Type) (int, int) {
	imported := 0
	skipped := 0

	for _, typeInfo := range types {
		// 检查类型是否已存在
		existingType, err := is.typeService.GetTypeByID(typeInfo.ID)
		if err == nil && existingType != nil {
			logger.Info("[导入] 类型已存在，跳过: ID=%s, Name=%s", typeInfo.ID, typeInfo.Name)
			skipped++
			continue
		}

		// 创建类型（使用原有ID）
		if err := is.createTypeWithID(typeInfo); err != nil {
			logger.Error("[导入] 创建类型失败: ID=%s, Name=%s, 错误=%v", typeInfo.ID, typeInfo.Name, err)
			skipped++
			continue
		}

		logger.Info("[导入] 类型导入成功: ID=%s, Name=%s", typeInfo.ID, typeInfo.Name)
		imported++
	}

	return imported, skipped
}

/**
 * importAccounts 导入账号数据
 * @param accounts 账号列表
 * @param backupCrypto 备份密码加密管理器
 * @return int 导入成功数量
 * @return int 跳过数量
 * @return int 错误数量
 * @return []SkippedAccountInfo 跳过的账号详情
 */
func (is *ImportService) importAccounts(accounts []ExportAccount, backupCrypto *crypto.CryptoManager) (int, int, int, []SkippedAccountInfo) {
	imported := 0
	skipped := 0
	errors := 0
	skippedDetails := make([]SkippedAccountInfo, 0)

	for _, account := range accounts {
		// 检查账号是否已存在
		existingAccount, err := is.accountService.GetAccountByID(account.ID)
		if err == nil && existingAccount != nil {
			logger.Info("[导入] 账号已存在，跳过: ID=%s, Title=%s", account.ID, account.Title)
			skipped++
			skippedDetails = append(skippedDetails, SkippedAccountInfo{
				ID:    account.ID,
				Title: account.Title,
				Name:  "", // 无法获取解密后的用户名，因为账号已存在
			})
			continue
		}

		// 解密账号数据并用当前登录密码重新加密
		decryptedAccount, err := is.convertAccountForImport(account, backupCrypto)
		if err != nil {
			logger.Error("[导入] 转换账号数据失败: ID=%s, Title=%s, 错误=%v", account.ID, account.Title, err)
			errors++
			continue
		}

		// 创建账号（使用原有ID）
		if err := is.createAccountWithID(decryptedAccount); err != nil {
			logger.Error("[导入] 创建账号失败: ID=%s, Title=%s, 错误=%v", account.ID, account.Title, err)
			errors++
			continue
		}

		logger.Info("[导入] 账号导入成功: ID=%s, Title=%s", account.ID, account.Title)
		imported++
	}

	return imported, skipped, errors, skippedDetails
}

/**
 * convertAccountForImport 转换账号数据用于导入（用当前登录密码重新加密）
 * @param exportAccount 导出账号
 * @param backupCrypto 备份密码加密管理器
 * @return models.AccountDecrypted 解密后的账号
 * @return error 错误信息
 */
func (is *ImportService) convertAccountForImport(exportAccount ExportAccount, backupCrypto *crypto.CryptoManager) (models.AccountDecrypted, error) {
	var account models.AccountDecrypted

	// 用备份密码解密敏感字段
	username, err := backupCrypto.Decrypt(exportAccount.Username)
	if err != nil {
		return account, fmt.Errorf("解密用户名失败: %w", err)
	}

	password, err := backupCrypto.Decrypt(exportAccount.Password)
	if err != nil {
		return account, fmt.Errorf("解密密码失败: %w", err)
	}

	url, err := backupCrypto.Decrypt(exportAccount.URL)
	if err != nil {
		return account, fmt.Errorf("解密地址失败: %w", err)
	}

	notes, err := backupCrypto.Decrypt(exportAccount.Notes)
	if err != nil {
		return account, fmt.Errorf("解密备注失败: %w", err)
	}

	// 构建解密后的账号数据
	account = models.AccountDecrypted{
		ID:          exportAccount.ID,
		Title:       exportAccount.Title,
		Username:    username,
		Password:    password,
		URL:         url,
		TypeID:      exportAccount.TypeID,
		Notes:       notes,
		Icon:        exportAccount.Icon,
		IsFavorite:  exportAccount.IsFavorite,
		UseCount:    exportAccount.UseCount,
		LastUsedAt:  exportAccount.LastUsedAt,
		CreatedAt:   exportAccount.CreatedAt,
		UpdatedAt:   exportAccount.UpdatedAt,
		InputMethod: exportAccount.InputMethod,
	}

	return account, nil
}

/**
 * createGroupWithID 创建分组（使用指定ID）
 * @param group 分组信息
 * @return error 错误信息
 */
func (is *ImportService) createGroupWithID(group models.Group) error {
	db := is.dbManager.GetDB()

	// 直接插入分组，使用原有ID
	_, err := db.Exec(`
		INSERT INTO groups (id, name, icon, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, group.ID, group.Name, group.Icon, group.SortOrder, group.CreatedAt, group.UpdatedAt)

	if err != nil {
		return fmt.Errorf("插入分组失败: %w", err)
	}

	return nil
}

/**
 * createTypeWithID 创建类型（使用指定ID）
 * @param typeInfo 类型信息
 * @return error 错误信息
 */
func (is *ImportService) createTypeWithID(typeInfo models.Type) error {
	db := is.dbManager.GetDB()

	// 直接插入类型，使用原有ID
	_, err := db.Exec(`
		INSERT INTO types (id, name, icon, group_id, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, typeInfo.ID, typeInfo.Name, typeInfo.Icon, typeInfo.GroupID, typeInfo.SortOrder, typeInfo.CreatedAt, typeInfo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("插入类型失败: %w", err)
	}

	return nil
}

/**
 * createAccountWithID 创建账号（使用指定ID）
 * @param account 账号信息
 * @return error 错误信息
 */
func (is *ImportService) createAccountWithID(account models.AccountDecrypted) error {
	// 将AccountDecrypted转换为Account类型
	accountToEncrypt := models.Account{
		ID:          account.ID,
		Title:       account.Title,
		Username:    account.Username,
		Password:    account.Password,
		URL:         account.URL,
		TypeID:      account.TypeID,
		Notes:       account.Notes,
		Icon:        account.Icon,
		IsFavorite:  account.IsFavorite,
		UseCount:    account.UseCount,
		LastUsedAt:  account.LastUsedAt,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
		InputMethod: account.InputMethod,
	}

	// 使用账号服务的加密功能来加密敏感字段
	encryptedAccount, err := is.accountService.encryptAccount(accountToEncrypt)
	if err != nil {
		return fmt.Errorf("加密账号数据失败: %w", err)
	}

	db := is.dbManager.GetDB()

	// 直接插入账号，使用原有ID
	_, err = db.Exec(`
		INSERT INTO accounts (id, title, username, password, url, typeid, notes, icon,
			is_favorite, use_count, last_used_at, created_at, updated_at, input_method)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, encryptedAccount.ID, encryptedAccount.Title, encryptedAccount.Username, encryptedAccount.Password,
		encryptedAccount.URL, encryptedAccount.TypeID, encryptedAccount.Notes, encryptedAccount.Icon,
		encryptedAccount.IsFavorite, encryptedAccount.UseCount, encryptedAccount.LastUsedAt,
		encryptedAccount.CreatedAt, encryptedAccount.UpdatedAt, encryptedAccount.InputMethod)

	if err != nil {
		return fmt.Errorf("插入账号失败: %w", err)
	}

	return nil
}

/**
 * normalizeImportData 数据梳理：统一分组和分类ID
 * @param exportData 导出数据指针
 * @return error 错误信息
 */
func (is *ImportService) normalizeImportData(exportData *ExportData) error {
	logger.Info("[数据梳理] 开始数据梳理，统一分组和分类ID")

	// 1. 统一分组ID
	if err := is.normalizeGroupIDs(exportData); err != nil {
		return fmt.Errorf("统一分组ID失败: %w", err)
	}

	// 2. 统一分类ID
	if err := is.normalizeTypeIDs(exportData); err != nil {
		return fmt.Errorf("统一分类ID失败: %w", err)
	}

	logger.Info("[数据梳理] 数据梳理完成")
	return nil
}

/**
 * normalizeGroupIDs 统一分组ID
 * @param exportData 导出数据指针
 * @return error 错误信息
 */
func (is *ImportService) normalizeGroupIDs(exportData *ExportData) error {
	logger.Info("[数据梳理] 开始统一分组ID")

	// 获取当前数据库中的所有分组
	dbGroups, err := is.groupService.GetAllGroups()
	if err != nil {
		return fmt.Errorf("获取数据库分组失败: %w", err)
	}

	// 创建分组名称到ID的映射
	dbGroupNameToID := make(map[string]string)
	for _, group := range dbGroups {
		dbGroupNameToID[group.Name] = group.ID
	}

	// 创建导入分组ID的映射关系（旧ID -> 新ID）
	groupIDMapping := make(map[string]string)

	// 遍历导入数据中的分组
	for i, importGroup := range exportData.Groups {
		if dbGroupID, exists := dbGroupNameToID[importGroup.Name]; exists {
			// 分组名称已存在，使用数据库中的分组ID
			oldID := importGroup.ID
			newID := dbGroupID
			groupIDMapping[oldID] = newID
			exportData.Groups[i].ID = newID
			logger.Info("[数据梳理] 分组名称 '%s' 已存在，ID映射: %s -> %s", importGroup.Name, oldID, newID)
		} else {
			// 分组名称不存在，保持原ID
			groupIDMapping[importGroup.ID] = importGroup.ID
		}
	}

	// 更新分类中的分组ID引用
	for i, importType := range exportData.Types {
		if newGroupID, exists := groupIDMapping[importType.GroupID]; exists {
			oldGroupID := importType.GroupID
			exportData.Types[i].GroupID = newGroupID
			logger.Info("[数据梳理] 分类 '%s' 的分组ID更新: %s -> %s", importType.Name, oldGroupID, newGroupID)
		}
	}

	// 更新账号中的分组ID引用（如果账号有分组ID字段的话）
	// 注意：根据当前的ExportAccount结构，账号没有直接的GroupID字段
	// 账号通过TypeID关联到分类，分类再关联到分组
	// 所以这里不需要更新账号的分组ID

	logger.Info("[数据梳理] 分组ID统一完成，共处理 %d 个分组", len(exportData.Groups))
	return nil
}

/**
 * normalizeTypeIDs 统一分类ID
 * @param exportData 导出数据指针
 * @return error 错误信息
 */
func (is *ImportService) normalizeTypeIDs(exportData *ExportData) error {
	logger.Info("[数据梳理] 开始统一分类ID")

	// 获取当前数据库中的所有分类
	dbTypes, err := is.typeService.GetAllTypes()
	if err != nil {
		return fmt.Errorf("获取数据库分类失败: %w", err)
	}

	// 创建分类名称到ID的映射（考虑分组ID）
	dbTypeNameToID := make(map[string]string) // key: "分组ID:分类名称", value: 分类ID
	for _, dbType := range dbTypes {
		key := fmt.Sprintf("%s:%s", dbType.GroupID, dbType.Name)
		dbTypeNameToID[key] = dbType.ID
	}

	// 创建导入分类ID的映射关系（旧ID -> 新ID）
	typeIDMapping := make(map[string]string)

	// 遍历导入数据中的分类
	for i, importType := range exportData.Types {
		key := fmt.Sprintf("%s:%s", importType.GroupID, importType.Name)
		if dbTypeID, exists := dbTypeNameToID[key]; exists {
			// 分类名称在相同分组下已存在，使用数据库中的分类ID
			oldID := importType.ID
			newID := dbTypeID
			typeIDMapping[oldID] = newID
			exportData.Types[i].ID = newID
			logger.Info("[数据梳理] 分类名称 '%s' (分组:%s) 已存在，ID映射: %s -> %s",
				importType.Name, importType.GroupID, oldID, newID)
		} else {
			// 分类名称不存在，保持原ID
			typeIDMapping[importType.ID] = importType.ID
		}
	}

	// 更新账号中的分类ID引用
	for i, importAccount := range exportData.Accounts {
		if newTypeID, exists := typeIDMapping[importAccount.TypeID]; exists {
			oldTypeID := importAccount.TypeID
			exportData.Accounts[i].TypeID = newTypeID
			logger.Info("[数据梳理] 账号 '%s' 的分类ID更新: %s -> %s", importAccount.Title, oldTypeID, newTypeID)
		}
	}

	logger.Info("[数据梳理] 分类ID统一完成，共处理 %d 个分类", len(exportData.Types))
	return nil
}

/**
 * SetCryptoManager 设置加密管理器
 * @param cryptoManager 加密管理器
 */
func (is *ImportService) SetCryptoManager(cryptoManager *crypto.CryptoManager) {
	is.cryptoManager = cryptoManager
}
