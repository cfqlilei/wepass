package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"wepassword/internal/crypto"
	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/models"
	"wepassword/internal/utils"
)

/**
 * 账号服务
 * @author 陈凤庆
 * @description 管理账号的增删改查操作
 * @modify 20251002 陈凤庆 password_service.go改名为account_service.go，对应accounts表
 */

/**
 * AccountService 账号服务
 */
type AccountService struct {
	dbManager     *database.DatabaseManager
	cryptoManager *crypto.CryptoManager
}

/**
 * NewAccountService 创建新的账号服务
 * @param dbManager 数据库管理器
 * @return *AccountService 账号服务实例
 * @modify 20251002 陈凤庆 NewPasswordService改名为NewAccountService
 */
func NewAccountService(dbManager *database.DatabaseManager) *AccountService {
	return &AccountService{
		dbManager: dbManager,
	}
}

/**
 * SetCryptoManager 设置加密管理器
 * @param cryptoManager 加密管理器
 * @modify 20251002 陈凤庆 添加日志记录
 */
func (as *AccountService) SetCryptoManager(cryptoManager *crypto.CryptoManager) {
	as.cryptoManager = cryptoManager
	logger.Debug("[账号服务] 加密管理器已设置")
}

/**
 * IsCryptoManagerSet 检查加密管理器是否已设置
 * @return bool 是否已设置
 */
func (as *AccountService) IsCryptoManagerSet() bool {
	return as.cryptoManager != nil
}

/**
 * GetAccountsByGroup 根据分组ID获取账号列表
 * @param groupID 分组ID
 * @return []models.AccountDecrypted 解密后的账号列表
 * @return error 错误信息
 * @modify 20251002 陈凤庆 GetPasswordItemsByGroup改名为GetAccountsByGroup
 */
func (as *AccountService) GetAccountsByGroup(groupID string) ([]models.AccountDecrypted, error) {
	// 构建查询条件JSON
	conditions := fmt.Sprintf(`{"group_id":"%s"}`, groupID)
	return as.GetAccountsByConditions(conditions)
}

// decryptField 解密单个字段
func (as *AccountService) decryptField(encryptedField string) string {
	if encryptedField == "" {
		return ""
	}
	decrypted, err := as.cryptoManager.Decrypt(encryptedField)
	if err != nil {
		logger.Error("[账号服务] 解密字段失败: %v", err)
		return ""
	}
	return decrypted
}

// maskUsername 脱敏用户名，只显示前2位和后2位，中间用*代替
func (as *AccountService) maskUsername(username string) string {
	if username == "" {
		return ""
	}

	// 如果用户名长度小于等于4，只显示第一个字符
	if len(username) <= 4 {
		return string(username[0]) + "***"
	}

	// 显示前2位和后2位，中间用*代替
	return username[:2] + "***" + username[len(username)-2:]
}

/**
 * GetAccountsByConditions 根据查询条件获取账号列表
 * @param conditions 查询条件JSON字符串，格式：{"group_id":"xxx","type_id":"xxx"}
 * @return []models.AccountDecrypted 解密后的账号列表
 * @return error 错误信息
 * @author 20251003 陈凤庆 统一账号查询方法，支持多种查询条件
 */
func (as *AccountService) GetAccountsByConditions(conditions string) ([]models.AccountDecrypted, error) {
	logger.Debug("[账号服务] GetAccountsByConditions 被调用，条件: %s", conditions)

	if !as.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	if as.cryptoManager == nil {
		return nil, fmt.Errorf("加密管理器未设置")
	}

	// 解析查询条件
	var conditionsMap map[string]interface{}
	if err := json.Unmarshal([]byte(conditions), &conditionsMap); err != nil {
		return nil, fmt.Errorf("解析查询条件失败: %w", err)
	}

	db := as.dbManager.GetDB()
	// 构建基础查询语句，包含地址字段用于右键菜单功能，并关联分组和类型表用于排序
	sqlQuery := `
		SELECT a.id, a.title, a.username, a.url, a.typeid, a.input_method, t.group_id
		FROM accounts a
		INNER JOIN types t ON a.typeid = t.id
		INNER JOIN groups g ON t.group_id = g.id
		WHERE 1=1`

	var args []interface{}

	// 根据条件动态添加WHERE子句
	if groupID, exists := conditionsMap["group_id"]; exists && groupID != "" {
		sqlQuery += " AND t.group_id = ?"
		args = append(args, groupID)
	}

	if typeID, exists := conditionsMap["type_id"]; exists && typeID != "" {
		sqlQuery += " AND a.typeid = ?"
		args = append(args, typeID)
	}

	// 按照分组、类别、标题进行排序
	sqlQuery += " ORDER BY g.sort_order ASC, g.name ASC, t.sort_order ASC, t.name ASC, a.title ASC"

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("查询账号失败: %w", err)
	}
	defer rows.Close()

	// 初始化为非nil的空切片，确保前端接收到[]而不是null
	accounts := make([]models.AccountDecrypted, 0)
	for rows.Next() {
		var account models.Account
		var groupID string

		err := rows.Scan(
			&account.ID, &account.Title, &account.Username, &account.URL, &account.TypeID, &account.InputMethod, &groupID,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描账号数据失败: %w", err)
		}

		// 解密用户名
		decryptedUsername, err := as.cryptoManager.Decrypt(account.Username)
		if err != nil {
			logger.Error("[解密] 解密用户名失败，账号ID: %s, 原始数据: %s, 错误: %v", account.ID, account.Username, err)
			logger.Info("[解密] 跳过损坏的账号数据，账号ID: %s, 标题: %s", account.ID, account.Title)
			continue // 跳过解密失败的账号
		}

		// 解密地址
		decryptedURL, err := as.cryptoManager.Decrypt(account.URL)
		if err != nil {
			logger.Error("[解密] 解密地址失败，账号ID: %s, 原始数据: %s, 错误: %v", account.ID, account.URL, err)
			continue // 跳过解密失败的账号
		}

		// 创建脱敏用户名
		maskedUsername := as.maskUsername(decryptedUsername)

		// 构建返回的账号数据（包含地址字段用于右键菜单功能）
		decryptedAccount := models.AccountDecrypted{
			ID:             account.ID,
			Title:          account.Title,
			URL:            decryptedURL,
			TypeID:         account.TypeID,
			InputMethod:    account.InputMethod,
			GroupID:        groupID,
			MaskedUsername: maskedUsername,
		}

		accounts = append(accounts, decryptedAccount)
	}

	return accounts, nil
}

/**
 * GetAccountsByTab 根据标签ID获取账号列表（保持向后兼容）
 * @param tabID 标签ID
 * @return []models.AccountDecrypted 解密后的账号列表
 * @return error 错误信息
 * @author 20251002 陈凤庆 新增按标签ID查询账号的方法
 * @modify 20251003 陈凤庆 重构为调用统一的GetAccountsByConditions方法
 */
func (as *AccountService) GetAccountsByTab(tabID string) ([]models.AccountDecrypted, error) {
	// 构建查询条件JSON
	conditions := fmt.Sprintf(`{"type_id":"%s"}`, tabID)
	return as.GetAccountsByConditions(conditions)
}

/**
 * GetAllAccounts 获取所有账号
 * @return []models.AccountDecrypted 解密后的账号列表
 * @return error 错误信息
 * @modify 20251002 陈凤庆 GetAllPasswordItems改名为GetAllAccounts
 */
func (as *AccountService) GetAllAccounts() ([]models.AccountDecrypted, error) {
	if !as.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	if as.cryptoManager == nil {
		return nil, fmt.Errorf("加密管理器未设置")
	}

	db := as.dbManager.GetDB()
	// 20251002 陈凤庆 查询accounts表，删除group_id字段
	rows, err := db.Query(`
		SELECT id, title, username, password, url, typeid, notes, icon,
			   is_favorite, use_count, last_used_at, created_at, updated_at
		FROM accounts
		ORDER BY is_favorite DESC, use_count DESC, title
	`)
	if err != nil {
		return nil, fmt.Errorf("查询账号失败: %w", err)
	}
	defer rows.Close()

	// 20251002 陈凤庆 初始化为非nil的空切片，确保前端接收到[]而不是null
	accounts := make([]models.AccountDecrypted, 0)
	for rows.Next() {
		var account models.Account
		// 20251002 陈凤庆 删除group_id和tab_id字段的扫描
		err := rows.Scan(
			&account.ID, &account.Title, &account.Username, &account.Password, &account.URL, &account.TypeID, &account.Notes, &account.Icon,
			&account.IsFavorite, &account.UseCount, &account.LastUsedAt, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描账号数据失败: %w", err)
		}

		// 解密账号
		decryptedAccount, err := as.decryptAccount(account)
		if err != nil {
			logger.Error("[账号服务] 解密账号失败，账号ID: %s, 标题: %s, 错误: %v", account.ID, account.Title, err)

			// 检查是否是明文数据，如果是则跳过并记录，稍后修复
			if !as.isBase64(account.Username) || !as.isBase64(account.Password) {
				logger.Info("[账号服务] 检测到明文数据，账号ID: %s，将在后台修复", account.ID)
				// 启动后台修复任务
				go as.fixPlaintextAccountAsync(account.ID)
			} else {
				logger.Info("[账号服务] 跳过损坏的账号数据，账号ID: %s, 标题: %s", account.ID, account.Title)
			}
			continue // 跳过解密失败的账号
		}

		accounts = append(accounts, decryptedAccount)
	}

	logger.Debug("[账号服务] GetAccountsByConditions 完成，返回 %d 个有效账号", len(accounts))
	return accounts, nil
}

/**
 * CreateAccount 创建账号
 * @param title 标题
 * @param username 用户名
 * @param password 密码
 * @param url 网址
 * @param typeID 类型ID
 * @param notes 备注
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
 * @return models.AccountDecrypted 创建的账号（解密后）
 * @return error 错误信息
 * @modify 20251002 陈凤庆 CreatePasswordItem改名为CreateAccount
 * @modify 20251003 陈凤庆 添加inputMethod参数
 * @modify 20251005 陈凤庆 支持第5种输入方式（键盘助手输入）
 */
func (as *AccountService) CreateAccount(title, username, password, url, typeID, notes string, inputMethod int) (models.AccountDecrypted, error) {
	// 20251019 陈凤庆 修复问题 004：增加详细的参数验证和调试日志
	logger.Info("[账号服务] 🔍 CreateAccount 详细参数检查:")
	logger.Info("  - title: \"%s\" (长度: %d, 是否为空: %t)", title, len(title), title == "")
	logger.Info("  - username: \"%s\" (长度: %d, 是否为空: %t)", username, len(username), username == "")
	logger.Info("  - password: %s (长度: %d, 是否为空: %t)", func() string {
		if password != "" {
			return "***已设置***"
		}
		return "未设置"
	}(), len(password), password == "")
	logger.Info("  - url: \"%s\" (长度: %d, 是否为空: %t)", url, len(url), url == "")
	logger.Info("  - typeID: \"%s\" (长度: %d, 是否为空: %t)", typeID, len(typeID), typeID == "")
	logger.Info("  - notes: \"%s\" (长度: %d, 是否为空: %t)", notes, len(notes), notes == "")
	logger.Info("  - inputMethod: %d", inputMethod)

	if !as.dbManager.IsOpened() {
		logger.Error("[账号服务] ❌ 数据库未打开")
		return models.AccountDecrypted{}, fmt.Errorf("数据库未打开")
	}

	if as.cryptoManager == nil {
		logger.Error("[账号服务] ❌ 加密管理器未设置")
		return models.AccountDecrypted{}, fmt.Errorf("加密管理器未设置")
	}

	if title == "" {
		logger.Error("[账号服务] ❌ 标题不能为空")
		return models.AccountDecrypted{}, fmt.Errorf("标题不能为空")
	}

	// 20251019 陈凤庆 增强类型ID验证，增加详细日志
	logger.Info("[账号服务] 🔍 开始验证类型ID: %s", typeID)
	err := as.validateTypeID(typeID)
	if err != nil {
		logger.Error("[账号服务] ❌ 类型ID验证失败: %v", err)
		logger.Error("[账号服务] 失败的类型ID: \"%s\" (长度: %d)", typeID, len(typeID))
		return models.AccountDecrypted{}, err
	}
	logger.Info("[账号服务] ✅ 类型ID验证通过")

	db := as.dbManager.GetDB()
	now := time.Now()

	// 生成新的GUID
	newID := utils.GenerateGUID()

	// 20251003 陈凤庆 验证输入方式参数
	// 20251005 陈凤庆 支持第4种输入方式（键盘助手输入，删除原第4种底层键盘API）
	// 20251005 陈凤庆 支持第5种输入方式（远程输入）
	if inputMethod < 1 || inputMethod > 5 {
		logger.Info("[账号服务] 输入方式参数无效: %d，重置为默认值1", inputMethod)
		inputMethod = 1 // 默认使用Unicode方式
	}

	logger.Info("[账号服务] 创建账号，输入方式: %d", inputMethod)

	// 创建账号对象
	account := models.Account{
		ID:          newID,
		Title:       title,
		Username:    username,
		Password:    password,
		URL:         url,
		TypeID:      typeID,
		Notes:       notes,
		Icon:        "",
		IsFavorite:  false,
		UseCount:    0,
		LastUsedAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
		InputMethod: inputMethod, // 20251003 陈凤庆 添加输入方式字段
	}

	logger.Info("[账号服务] 创建的账号对象，InputMethod: %d", account.InputMethod)

	// 加密账号
	encryptedAccount, err := as.encryptAccount(account)
	if err != nil {
		return models.AccountDecrypted{}, fmt.Errorf("加密账号失败: %w", err)
	}

	logger.Info("[账号服务] 加密后的账号对象，InputMethod: %d", encryptedAccount.InputMethod)

	// 插入数据库
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型，插入accounts表，删除group_id字段
	// 20251003 陈凤庆 添加input_method字段
	var insertErr error
	_, insertErr = db.Exec(`
		INSERT INTO accounts (id, title, username, password, url, typeid, notes, icon, is_favorite, use_count, last_used_at, created_at, updated_at, input_method)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, encryptedAccount.ID, encryptedAccount.Title, encryptedAccount.Username, encryptedAccount.Password, encryptedAccount.URL, encryptedAccount.TypeID, encryptedAccount.Notes, encryptedAccount.Icon, encryptedAccount.IsFavorite, encryptedAccount.UseCount, encryptedAccount.LastUsedAt, encryptedAccount.CreatedAt, encryptedAccount.UpdatedAt, encryptedAccount.InputMethod)

	if insertErr != nil {
		return models.AccountDecrypted{}, fmt.Errorf("创建账号失败: %w", insertErr)
	}

	// 返回解密后的账号
	decryptedAccount, err := as.decryptAccount(encryptedAccount)
	if err != nil {
		return models.AccountDecrypted{}, fmt.Errorf("解密账号失败: %w", err)
	}

	// 20251003 陈凤庆 添加详细的创建成功日志，包含input_method字段
	logger.Info("[账号服务] 新账号创建成功: ID=%s, 标题=%s, 输入方式=%d", decryptedAccount.ID, decryptedAccount.Title, decryptedAccount.InputMethod)

	return decryptedAccount, nil
}

/**
 * UpdateAccount 更新账号
 * @param account 账号信息（解密后）
 * @return error 错误信息
 * @modify 20251002 陈凤庆 UpdatePasswordItem改名为UpdateAccount
 */
func (as *AccountService) UpdateAccount(account models.AccountDecrypted) error {
	if !as.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if as.cryptoManager == nil {
		return fmt.Errorf("加密管理器未设置")
	}

	// 验证类型ID
	err := as.validateTypeID(account.TypeID)
	if err != nil {
		return err
	}

	// 转换为加密的账号对象
	// 20251003 陈凤庆 添加InputMethod字段，修复更新时input_method字段丢失问题
	encryptedAccount := models.Account{
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
		UpdatedAt:   time.Now(),
		InputMethod: account.InputMethod, // 20251003 陈凤庆 添加输入方式字段，修复更新时丢失问题
	}

	// 加密账号
	encryptedAccount, err = as.encryptAccount(encryptedAccount)
	if err != nil {
		return fmt.Errorf("加密账号失败: %w", err)
	}

	db := as.dbManager.GetDB()

	// 更新数据库
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型，更新accounts表，删除group_id字段
	// 20251003 陈凤庆 添加input_method字段更新
	var updateErr error
	_, updateErr = db.Exec(`
		UPDATE accounts
		SET title = ?, username = ?, password = ?, url = ?, typeid = ?, notes = ?, icon = ?, is_favorite = ?, use_count = ?, last_used_at = ?, updated_at = ?, input_method = ?
		WHERE id = ?
	`, encryptedAccount.Title, encryptedAccount.Username, encryptedAccount.Password, encryptedAccount.URL, encryptedAccount.TypeID, encryptedAccount.Notes, encryptedAccount.Icon, encryptedAccount.IsFavorite, encryptedAccount.UseCount, encryptedAccount.LastUsedAt, encryptedAccount.UpdatedAt, encryptedAccount.InputMethod, encryptedAccount.ID)

	if updateErr != nil {
		return fmt.Errorf("更新账号失败: %w", updateErr)
	}

	return nil
}

/**
 * DeleteAccount 删除账号
 * @param id 账号ID
 * @return error 错误信息
 * @modify 20251002 陈凤庆 DeletePasswordItem改名为DeleteAccount
 */
func (as *AccountService) DeleteAccount(id string) error {
	if !as.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	db := as.dbManager.GetDB()
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型，删除accounts表数据
	var deleteErr error
	_, deleteErr = db.Exec("DELETE FROM accounts WHERE id = ?", id)
	if deleteErr != nil {
		return fmt.Errorf("删除账号失败: %w", deleteErr)
	}

	return nil
}

/**
 * SearchAccounts 搜索账号
 * @param keyword 搜索关键词
 * @return []models.AccountDecrypted 搜索结果
 * @return error 错误信息
 * @modify 20251002 陈凤庆 SearchPasswords改名为SearchAccounts
 */
func (as *AccountService) SearchAccounts(keyword string) ([]models.AccountDecrypted, error) {
	if !as.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := as.dbManager.GetDB()
	// 20251002 陈凤庆 查询accounts表，删除group_id字段
	// 20251003 陈凤庆 添加input_method字段查询
	rows, err := db.Query(`
		SELECT id, title, username, password, url, typeid, notes, icon,
			   is_favorite, use_count, last_used_at, created_at, updated_at, input_method
		FROM accounts
		WHERE title LIKE ? OR url LIKE ?
		ORDER BY is_favorite DESC, use_count DESC, title
	`, "%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("搜索账号失败: %w", err)
	}
	defer rows.Close()

	// 20251002 陈凤庆 初始化为非nil的空切片，确保前端接收到[]而不是null
	accounts := make([]models.AccountDecrypted, 0)
	for rows.Next() {
		var account models.Account
		// 20251003 陈凤庆 添加input_method字段的扫描
		err := rows.Scan(
			&account.ID, &account.Title, &account.Username, &account.Password, &account.URL, &account.TypeID, &account.Notes,
			&account.Icon, &account.IsFavorite, &account.UseCount, &account.LastUsedAt,
			&account.CreatedAt, &account.UpdatedAt, &account.InputMethod,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描账号数据失败: %w", err)
		}

		// 解密数据
		decryptedAccount, err := as.decryptAccount(account)
		if err != nil {
			logger.Error("[账号服务] 解密账号失败: %v", err)
			continue // 跳过解密失败的账号
		}

		accounts = append(accounts, decryptedAccount)
	}

	return accounts, nil
}

/**
 * encryptAccount 加密账号敏感数据
 * @param account 明文账号
 * @return models.Account 加密后的账号
 * @return error 错误信息
 * @modify 20251002 陈凤庆 encryptPasswordItem改名为encryptAccount
 */
func (as *AccountService) encryptAccount(account models.Account) (models.Account, error) {
	if as.cryptoManager == nil {
		return models.Account{}, fmt.Errorf("加密管理器未初始化")
	}

	// 20251002 陈凤庆 删除GroupID和TabID字段
	// 20251003 陈凤庆 添加InputMethod字段，修复input_method字段丢失问题
	encryptedAccount := models.Account{
		ID:          account.ID,
		Title:       account.Title, // 标题不加密
		Icon:        account.Icon,
		IsFavorite:  account.IsFavorite,
		UseCount:    account.UseCount,
		LastUsedAt:  account.LastUsedAt,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
		TypeID:      account.TypeID,      // TypeID不加密，保持明文存储
		InputMethod: account.InputMethod, // 20251003 陈凤庆 添加输入方式字段，修复丢失问题
	}

	// 加密敏感字段
	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型
	var err error
	encryptedAccount.Username, err = as.cryptoManager.Encrypt(account.Username)
	if err != nil {
		return models.Account{}, fmt.Errorf("加密用户名失败: %w", err)
	}

	encryptedAccount.Password, err = as.cryptoManager.Encrypt(account.Password)
	if err != nil {
		return models.Account{}, fmt.Errorf("加密密码失败: %w", err)
	}

	encryptedAccount.URL, err = as.cryptoManager.Encrypt(account.URL)
	if err != nil {
		return models.Account{}, fmt.Errorf("加密地址失败: %w", err)
	}

	encryptedAccount.Notes, err = as.cryptoManager.Encrypt(account.Notes)
	if err != nil {
		return models.Account{}, fmt.Errorf("加密备注失败: %w", err)
	}

	return encryptedAccount, nil
}

/**
 * decryptAccount 解密账号敏感数据
 * @param account 加密的账号
 * @return models.AccountDecrypted 解密后的账号
 * @return error 错误信息
 * @modify 20251002 陈凤庆 decryptPasswordItem改名为decryptAccount
 */
func (as *AccountService) decryptAccount(account models.Account) (models.AccountDecrypted, error) {
	if as.cryptoManager == nil {
		return models.AccountDecrypted{}, fmt.Errorf("加密管理器未初始化")
	}

	logger.Debug("[解密] 开始解密账号: ID=%s, Title=%s, InputMethod=%d", account.ID, account.Title, account.InputMethod)

	// 20251002 陈凤庆 删除GroupID和TabID字段
	// 20251003 陈凤庆 添加InputMethod字段，修复编辑和副本生成时input_method显示不正确的问题
	decryptedAccount := models.AccountDecrypted{
		ID:          account.ID,
		Title:       account.Title, // 标题不需要解密
		Icon:        account.Icon,
		IsFavorite:  account.IsFavorite,
		UseCount:    account.UseCount,
		LastUsedAt:  account.LastUsedAt,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
		TypeID:      account.TypeID,      // TypeID不需要解密
		InputMethod: account.InputMethod, // 20251003 陈凤庆 添加输入方式字段，修复编辑和副本生成显示问题
	}

	// 20251002 陈凤庆 解密敏感字段，不使用:=赋值，先声明变量类型
	var err error

	// 解密用户名
	logger.Debug("[解密] 解密用户名字段，原始数据长度: %d", len(account.Username))
	decryptedAccount.Username, err = as.cryptoManager.Decrypt(account.Username)
	if err != nil {
		logger.Error("[解密] 解密用户名失败，账号ID: %s, 原始数据: %s, 错误: %v", account.ID, account.Username, err)
		return models.AccountDecrypted{}, fmt.Errorf("解密用户名失败: %w", err)
	}

	// 验证解密后的用户名
	if decryptedAccount.Username == "" {
		logger.Error("[解密] 解密后用户名为空，账号ID: %s", account.ID)
		return models.AccountDecrypted{}, fmt.Errorf("解密后用户名为空")
	}

	// 检查用户名是否包含异常重复字符
	if len(decryptedAccount.Username) > 2 {
		allSame := true
		firstChar := decryptedAccount.Username[0]
		for _, char := range decryptedAccount.Username {
			if char != rune(firstChar) {
				allSame = false
				break
			}
		}
		if allSame {
			logger.Error("[解密] 解密后用户名包含异常重复字符: '%s'，账号ID: %s", decryptedAccount.Username, account.ID)
			return models.AccountDecrypted{}, fmt.Errorf("解密后用户名数据异常，包含重复字符: %s", decryptedAccount.Username)
		}
	}

	// 解密密码
	logger.Debug("[解密] 解密密码字段，原始数据长度: %d", len(account.Password))
	decryptedAccount.Password, err = as.cryptoManager.Decrypt(account.Password)
	if err != nil {
		logger.Error("[解密] 解密密码失败，账号ID: %s, 原始数据: %s, 错误: %v", account.ID, account.Password, err)
		return models.AccountDecrypted{}, fmt.Errorf("解密密码失败: %w", err)
	}

	// 验证解密后的密码
	if decryptedAccount.Password == "" {
		logger.Error("[解密] 解密后密码为空，账号ID: %s", account.ID)
		return models.AccountDecrypted{}, fmt.Errorf("解密后密码为空")
	}

	// 检查密码是否包含异常重复字符
	if len(decryptedAccount.Password) > 2 {
		allSame := true
		firstChar := decryptedAccount.Password[0]
		for _, char := range decryptedAccount.Password {
			if char != rune(firstChar) {
				allSame = false
				break
			}
		}
		if allSame {
			logger.Error("[解密] 解密后密码包含异常重复字符，账号ID: %s", account.ID)
			return models.AccountDecrypted{}, fmt.Errorf("解密后密码数据异常，包含重复字符")
		}
	}

	// 解密地址
	logger.Debug("[解密] 解密地址字段，原始数据长度: %d", len(account.URL))
	decryptedAccount.URL, err = as.cryptoManager.Decrypt(account.URL)
	if err != nil {
		logger.Error("[解密] 解密地址失败，账号ID: %s, 原始数据: %s, 错误: %v", account.ID, account.URL, err)
		return models.AccountDecrypted{}, fmt.Errorf("解密地址失败: %w", err)
	}

	// 解密备注
	logger.Debug("[解密] 解密备注字段，原始数据长度: %d", len(account.Notes))
	decryptedAccount.Notes, err = as.cryptoManager.Decrypt(account.Notes)
	if err != nil {
		logger.Error("[解密] 解密备注失败，账号ID: %s, 原始数据: %s, 错误: %v", account.ID, account.Notes, err)
		return models.AccountDecrypted{}, fmt.Errorf("解密备注失败: %w", err)
	}

	// 生成掩码数据
	// 20251002 陈凤庆 使用crypto包的MaskString方法
	decryptedAccount.MaskedUsername = crypto.MaskString(decryptedAccount.Username)
	decryptedAccount.MaskedPassword = crypto.MaskString(decryptedAccount.Password)

	logger.Debug("[解密] 账号解密完成: ID=%s, InputMethod=%d", account.ID, decryptedAccount.InputMethod)
	return decryptedAccount, nil
}

/**
 * GetAccountByID 根据ID获取账号
 * @param id 账号ID
 * @return *models.AccountDecrypted 解密后的账号
 * @return error 错误信息
 * @modify 20251002 陈凤庆 GetPasswordItemByID改名为GetAccountByID
 */
func (as *AccountService) GetAccountByID(id string) (*models.AccountDecrypted, error) {
	if !as.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := as.dbManager.GetDB()
	var account models.Account
	var groupID string // 20251002 陈凤庆 添加group_id变量
	// 20251002 陈凤庆 查询accounts表，关联types表获取group_id
	// 20251003 陈凤庆 添加input_method字段查询
	err := db.QueryRow(`
		SELECT a.id, a.title, a.username, a.password, a.url, a.typeid, a.notes, a.icon,
			   a.is_favorite, a.use_count, a.last_used_at, a.created_at, a.updated_at, a.input_method, t.group_id
		FROM accounts a
		INNER JOIN types t ON a.typeid = t.id
		WHERE a.id = ?
	`, id).Scan(
		&account.ID, &account.Title, &account.Username, &account.Password, &account.URL, &account.TypeID, &account.Notes,
		&account.Icon, &account.IsFavorite, &account.UseCount, &account.LastUsedAt,
		&account.CreatedAt, &account.UpdatedAt, &account.InputMethod, &groupID,
	)

	if err != nil {
		return nil, fmt.Errorf("查询账号失败: %w", err)
	}

	// 解密数据
	decryptedAccount, err := as.decryptAccount(account)
	if err != nil {
		return nil, fmt.Errorf("解密账号失败: %w", err)
	}

	// 20251002 陈凤庆 设置分组ID
	decryptedAccount.GroupID = groupID

	return &decryptedAccount, nil
}

/**
 * GetAccountDetail 根据ID获取账号详情（用于详情页面显示）
 * @param id 账号ID
 * @return *models.AccountDecrypted 解密后的账号详情
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增账号详情查询方法，返回解密后的用户名，密码不返回，备注返回脱敏版本
 */
func (as *AccountService) GetAccountDetail(id string) (*models.AccountDecrypted, error) {
	if !as.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	if as.cryptoManager == nil {
		return nil, fmt.Errorf("加密管理器未设置")
	}

	db := as.dbManager.GetDB()
	var account models.Account
	var groupID string
	// 查询账号基本信息
	err := db.QueryRow(`
		SELECT a.id, a.title, a.username, a.url, a.typeid, a.notes, a.icon,
			   a.is_favorite, a.use_count, a.last_used_at, a.created_at, a.updated_at, a.input_method, t.group_id
		FROM accounts a
		INNER JOIN types t ON a.typeid = t.id
		WHERE a.id = ?
	`, id).Scan(
		&account.ID, &account.Title, &account.Username, &account.URL, &account.TypeID, &account.Notes,
		&account.Icon, &account.IsFavorite, &account.UseCount, &account.LastUsedAt,
		&account.CreatedAt, &account.UpdatedAt, &account.InputMethod, &groupID,
	)

	if err != nil {
		return nil, fmt.Errorf("查询账号失败: %w", err)
	}

	// 解密必要字段
	decryptedAccount := models.AccountDecrypted{
		ID:          account.ID,
		Title:       account.Title,
		TypeID:      account.TypeID,
		Icon:        account.Icon,
		IsFavorite:  account.IsFavorite,
		UseCount:    account.UseCount,
		LastUsedAt:  account.LastUsedAt,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
		InputMethod: account.InputMethod,
		GroupID:     groupID,
	}

	// 解密用户名
	decryptedAccount.Username, err = as.cryptoManager.Decrypt(account.Username)
	if err != nil {
		logger.Error("[账号服务] 解密用户名失败: %v", err)
		return nil, fmt.Errorf("解密用户名失败: %w", err)
	}

	// 解密地址
	decryptedAccount.URL, err = as.cryptoManager.Decrypt(account.URL)
	if err != nil {
		logger.Error("[账号服务] 解密地址失败: %v", err)
		return nil, fmt.Errorf("解密地址失败: %w", err)
	}

	// 解密备注并脱敏
	decryptedNotes, err := as.cryptoManager.Decrypt(account.Notes)
	if err != nil {
		logger.Error("[账号服务] 解密备注失败: %v", err)
		return nil, fmt.Errorf("解密备注失败: %w", err)
	}
	decryptedAccount.Notes = as.maskNotes(decryptedNotes) // 脱敏备注

	// 密码不返回，保持空值
	decryptedAccount.Password = ""

	// 生成脱敏用户名
	decryptedAccount.MaskedUsername = as.maskUsername(decryptedAccount.Username)

	return &decryptedAccount, nil
}

/**
 * GetAccountRaw 根据ID获取账号原始数据（数据库中的加密数据）
 * @param id 账号ID
 * @return *models.Account 数据库中的加密账号数据
 * @return error 错误信息
 * @author 20251003 陈凤庆 新增获取原始加密数据的方法，用于特殊场景
 */
func (as *AccountService) GetAccountRaw(id string) (*models.Account, error) {
	if !as.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	db := as.dbManager.GetDB()
	var account models.Account
	// 查询账号原始数据（加密状态）
	err := db.QueryRow(`
		SELECT id, title, username, password, url, typeid, notes, icon,
			   is_favorite, use_count, last_used_at, created_at, updated_at, input_method
		FROM accounts
		WHERE id = ?
	`, id).Scan(
		&account.ID, &account.Title, &account.Username, &account.Password, &account.URL, &account.TypeID, &account.Notes,
		&account.Icon, &account.IsFavorite, &account.UseCount, &account.LastUsedAt,
		&account.CreatedAt, &account.UpdatedAt, &account.InputMethod,
	)

	if err != nil {
		return nil, fmt.Errorf("查询账号原始数据失败: %w", err)
	}

	return &account, nil
}

// maskNotes 脱敏备注，只显示前10个字符，后面用...代替
func (as *AccountService) maskNotes(notes string) string {
	if notes == "" {
		return ""
	}

	runes := []rune(notes)
	if len(runes) <= 10 {
		return notes
	}

	return string(runes[:10]) + "..."
}

/**
 * validateTypeID 验证类型ID是否存在
 * @param typeID 类型ID
 * @return error 错误信息
 * @modify 20251002 陈凤庆 替换validateGroupAndTab方法，只验证类型ID
 */
func (as *AccountService) validateTypeID(typeID string) error {
	// 20251019 陈凤庆 修复问题 004：增加详细的类型ID验证日志
	logger.Info("[账号服务] 🔍 validateTypeID 开始验证:")
	logger.Info("  - typeID: \"%s\" (长度: %d)", typeID, len(typeID))
	logger.Info("  - 数据库状态: %t", as.dbManager.IsOpened())

	if !as.dbManager.IsOpened() {
		logger.Error("[账号服务] ❌ 数据库未打开")
		return fmt.Errorf("数据库未打开")
	}

	if typeID == "" {
		logger.Error("[账号服务] ❌ 类型ID不能为空")
		return fmt.Errorf("类型ID不能为空")
	}

	db := as.dbManager.GetDB()
	logger.Info("[账号服务] 🔍 开始查询数据库中的类型ID")

	// 验证类型是否存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM types WHERE id = ?", typeID).Scan(&count)
	if err != nil {
		logger.Error("[账号服务] ❌ 查询类型失败: %v", err)
		logger.Error("[账号服务] SQL查询参数: typeID=\"%s\"", typeID)
		return fmt.Errorf("查询类型失败: %w", err)
	}

	logger.Info("[账号服务] 🔍 查询结果: count=%d", count)

	if count == 0 {
		logger.Error("[账号服务] ❌ 类型ID不存在: %s", typeID)

		// 20251019 陈凤庆 新增：查询所有可用的类型ID，帮助调试
		logger.Info("[账号服务] 🔍 查询所有可用的类型ID以供参考:")
		rows, queryErr := db.Query("SELECT id, name, group_id FROM types ORDER BY group_id, name")
		if queryErr == nil {
			defer rows.Close()
			for rows.Next() {
				var id, name, groupID string
				if scanErr := rows.Scan(&id, &name, &groupID); scanErr == nil {
					logger.Info("  - 可用类型: ID=\"%s\", Name=\"%s\", GroupID=\"%s\"", id, name, groupID)
				}
			}
		} else {
			logger.Error("[账号服务] 查询可用类型ID失败: %v", queryErr)
		}

		return fmt.Errorf("类型ID不存在: %s", typeID)
	}

	logger.Info("[账号服务] ✅ 类型ID验证通过: %s", typeID)
	return nil
}

/**
 * UpdateAccountUsage 更新账号使用次数和最后使用时间
 * @param id 账号ID
 * @return error 错误信息
 * @modify 20251002 陈凤庆 UpdatePasswordItemUsage改名为UpdateAccountUsage
 */
func (as *AccountService) UpdateAccountUsage(id string) error {
	if !as.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	db := as.dbManager.GetDB()
	now := time.Now()

	// 20251002 陈凤庆 不使用:=赋值，先声明变量类型，更新accounts表
	var updateErr error
	_, updateErr = db.Exec(`
		UPDATE accounts
		SET use_count = use_count + 1, last_used_at = ?, updated_at = ?
		WHERE id = ?
	`, now, now, id)

	if updateErr != nil {
		return fmt.Errorf("更新使用次数失败: %w", updateErr)
	}

	return nil
}

/**
 * UpdateAccountGroup 更新账号的分组（通过更新TypeID）
 * @param accountID 账号ID
 * @param typeID 新的类型ID
 * @return error 错误信息
 * @author 20251005 陈凤庆 新增更新账号分组的方法，用于更改分组功能
 */
func (as *AccountService) UpdateAccountGroup(accountID string, typeID string) error {
	if !as.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if accountID == "" {
		return fmt.Errorf("账号ID不能为空")
	}

	if typeID == "" {
		return fmt.Errorf("类型ID不能为空")
	}

	// 验证类型ID是否存在
	err := as.validateTypeID(typeID)
	if err != nil {
		return err
	}

	db := as.dbManager.GetDB()
	now := time.Now()

	// 更新账号的类型ID
	_, err = db.Exec(`
		UPDATE accounts
		SET typeid = ?, updated_at = ?
		WHERE id = ?
	`, typeID, now, accountID)

	if err != nil {
		return fmt.Errorf("更新账号分组失败: %w", err)
	}

	logger.Info("[账号服务] 账号分组更新成功，账号ID: %s, 新类型ID: %s", accountID, typeID)
	return nil
}

/**
 * tryFixPlaintextAccount 尝试修复明文账号数据
 * @param account 账号指针
 * @return bool 是否修复成功
 * @author 20251003 陈凤庆 修复明文数据问题
 */
func (as *AccountService) tryFixPlaintextAccount(account *models.Account) bool {
	if as.cryptoManager == nil {
		return false
	}

	// 检查是否是明文数据（不是Base64格式）
	if !as.isBase64(account.Username) || !as.isBase64(account.Password) {
		logger.Info("[账号修复] 检测到明文数据，账号ID: %s", account.ID)

		// 加密明文数据
		encryptedUsername, err := as.cryptoManager.Encrypt(account.Username)
		if err != nil {
			logger.Error("[账号修复] 加密用户名失败: %v", err)
			return false
		}

		encryptedPassword, err := as.cryptoManager.Encrypt(account.Password)
		if err != nil {
			logger.Error("[账号修复] 加密密码失败: %v", err)
			return false
		}

		// 更新数据库（使用事务确保数据一致性）
		db := as.dbManager.GetDB()
		tx, err := db.Begin()
		if err != nil {
			logger.Error("[账号修复] 开始事务失败: %v", err)
			return false
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
			UPDATE accounts
			SET username = ?, password = ?, updated_at = ?
			WHERE id = ?
		`, encryptedUsername, encryptedPassword, time.Now(), account.ID)

		if err != nil {
			logger.Error("[账号修复] 更新数据库失败: %v", err)
			return false
		}

		err = tx.Commit()
		if err != nil {
			logger.Error("[账号修复] 提交事务失败: %v", err)
			return false
		}

		// 更新内存中的数据
		account.Username = encryptedUsername
		account.Password = encryptedPassword

		logger.Info("[账号修复] 成功修复明文账号，账号ID: %s", account.ID)
		return true
	}

	return false
}

/**
 * isBase64 检查字符串是否是有效的Base64格式
 * @param s 要检查的字符串
 * @return bool 是否是Base64格式
 */
func (as *AccountService) isBase64(s string) bool {
	if s == "" {
		return false
	}
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

/**
 * fixPlaintextAccountAsync 异步修复明文账号数据
 * @param accountID 账号ID
 * @author 20251003 陈凤庆 异步修复明文数据
 */
func (as *AccountService) fixPlaintextAccountAsync(accountID string) {
	// 等待一段时间，避免数据库忙碌
	time.Sleep(5 * time.Second)

	logger.Info("[账号修复] 开始异步修复明文账号，账号ID: %s", accountID)

	if as.cryptoManager == nil {
		logger.Error("[账号修复] CryptoManager未初始化，无法修复账号ID: %s", accountID)
		return
	}

	// 重新查询账号数据
	db := as.dbManager.GetDB()
	var account models.Account
	err := db.QueryRow(`
		SELECT id, title, username, password, url, notes, typeid,
		       is_favorite, use_count, created_at, updated_at
		FROM accounts WHERE id = ?
	`, accountID).Scan(
		&account.ID, &account.Title, &account.Username, &account.Password,
		&account.URL, &account.Notes, &account.TypeID,
		&account.IsFavorite, &account.UseCount, &account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		logger.Error("[账号修复] 查询账号失败，账号ID: %s, 错误: %v", accountID, err)
		return
	}

	// 检查是否仍然是明文数据
	if !as.isBase64(account.Username) || !as.isBase64(account.Password) {
		logger.Info("[账号修复] 确认为明文数据，开始加密，账号ID: %s", accountID)

		// 加密明文数据
		encryptedUsername, err := as.cryptoManager.Encrypt(account.Username)
		if err != nil {
			logger.Error("[账号修复] 加密用户名失败，账号ID: %s, 错误: %v", accountID, err)
			return
		}

		encryptedPassword, err := as.cryptoManager.Encrypt(account.Password)
		if err != nil {
			logger.Error("[账号修复] 加密密码失败，账号ID: %s, 错误: %v", accountID, err)
			return
		}

		// 使用事务更新数据库
		tx, err := db.Begin()
		if err != nil {
			logger.Error("[账号修复] 开始事务失败，账号ID: %s, 错误: %v", accountID, err)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
			UPDATE accounts
			SET username = ?, password = ?, updated_at = ?
			WHERE id = ?
		`, encryptedUsername, encryptedPassword, time.Now(), accountID)

		if err != nil {
			logger.Error("[账号修复] 更新数据库失败，账号ID: %s, 错误: %v", accountID, err)
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Error("[账号修复] 提交事务失败，账号ID: %s, 错误: %v", accountID, err)
			return
		}

		logger.Info("[账号修复] ✅ 成功修复明文账号，账号ID: %s", accountID)
	} else {
		logger.Info("[账号修复] 账号数据已经是加密格式，无需修复，账号ID: %s", accountID)
	}
}

/**
 * CheckAndRepairCorruptedAccounts 检查并修复损坏的账号数据
 * @return int 修复的账号数量
 * @return error 错误信息
 * @author 20251004 陈凤庆 新增数据修复功能
 */
func (as *AccountService) CheckAndRepairCorruptedAccounts() (int, error) {
	if !as.dbManager.IsOpened() {
		return 0, fmt.Errorf("数据库未打开")
	}

	if as.cryptoManager == nil {
		return 0, fmt.Errorf("加密管理器未设置")
	}

	logger.Info("[数据修复] 开始检查损坏的账号数据...")

	// 查询所有账号
	query := `SELECT id, title, username, password, url, notes FROM accounts`
	rows, err := as.dbManager.GetDB().Query(query)
	if err != nil {
		return 0, fmt.Errorf("查询账号数据失败: %w", err)
	}
	defer rows.Close()

	var corruptedCount int
	var repairedCount int

	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.ID, &account.Title, &account.Username, &account.Password, &account.URL, &account.Notes)
		if err != nil {
			logger.Error("[数据修复] 扫描账号数据失败: %v", err)
			continue
		}

		// 检查是否能正常解密
		isCorrupted := false

		// 检查用户名
		if account.Username != "" {
			if _, err := as.cryptoManager.Decrypt(account.Username); err != nil {
				logger.Error("[数据修复] 账号 %s (%s) 用户名解密失败: %v", account.ID, account.Title, err)
				isCorrupted = true
			}
		}

		// 检查密码
		if account.Password != "" {
			if _, err := as.cryptoManager.Decrypt(account.Password); err != nil {
				logger.Error("[数据修复] 账号 %s (%s) 密码解密失败: %v", account.ID, account.Title, err)
				isCorrupted = true
			}
		}

		// 检查URL
		if account.URL != "" {
			if _, err := as.cryptoManager.Decrypt(account.URL); err != nil {
				logger.Error("[数据修复] 账号 %s (%s) URL解密失败: %v", account.ID, account.Title, err)
				isCorrupted = true
			}
		}

		// 检查备注
		if account.Notes != "" {
			if _, err := as.cryptoManager.Decrypt(account.Notes); err != nil {
				logger.Error("[数据修复] 账号 %s (%s) 备注解密失败: %v", account.ID, account.Title, err)
				isCorrupted = true
			}
		}

		if isCorrupted {
			corruptedCount++
			logger.Info("[数据修复] 发现损坏的账号: %s (%s)", account.ID, account.Title)

			// 这里可以添加修复逻辑，比如删除损坏的账号或者尝试恢复
			// 目前只是记录，不进行实际修复
		}
	}

	logger.Info("[数据修复] 检查完成，发现 %d 个损坏的账号，修复了 %d 个", corruptedCount, repairedCount)
	return repairedCount, nil
}
