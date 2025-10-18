package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wepassword/internal/models"
	"wepassword/internal/utils"

	_ "modernc.org/sqlite"
)

// DefaultDataTranslations 默认数据的多语言翻译
type DefaultDataTranslations struct {
	DefaultGroupName string
	WebsiteType      string
	ApplicationType  string
	WiFiType         string
}

// getDefaultDataTranslations 获取默认数据的多语言翻译
func getDefaultDataTranslations(language string) DefaultDataTranslations {
	switch language {
	case "zh-TW":
		return DefaultDataTranslations{
			DefaultGroupName: "預設",
			WebsiteType:      "網站帳號",
			ApplicationType:  "應用程式",
			WiFiType:         "WiFi密碼",
		}
	case "en-US":
		return DefaultDataTranslations{
			DefaultGroupName: "Default",
			WebsiteType:      "Website Account",
			ApplicationType:  "Application",
			WiFiType:         "WiFi Password",
		}
	case "ru-RU":
		return DefaultDataTranslations{
			DefaultGroupName: "По умолчанию",
			WebsiteType:      "Веб-сайт",
			ApplicationType:  "Приложение",
			WiFiType:         "Пароль WiFi",
		}
	case "es-ES":
		return DefaultDataTranslations{
			DefaultGroupName: "Predeterminado",
			WebsiteType:      "Cuenta del sitio web",
			ApplicationType:  "Aplicación",
			WiFiType:         "Contraseña WiFi",
		}
	case "fr-FR":
		return DefaultDataTranslations{
			DefaultGroupName: "Par défaut",
			WebsiteType:      "Compte de site web",
			ApplicationType:  "Application",
			WiFiType:         "Mot de passe WiFi",
		}
	case "de-DE":
		return DefaultDataTranslations{
			DefaultGroupName: "Standard",
			WebsiteType:      "Website-Konto",
			ApplicationType:  "Anwendung",
			WiFiType:         "WiFi-Passwort",
		}
	default: // zh-CN 或其他语言默认使用中文
		return DefaultDataTranslations{
			DefaultGroupName: "默认",
			WebsiteType:      "网站账号",
			ApplicationType:  "应用程序",
			WiFiType:         "WiFi密码",
		}
	}
}

/**
 * 数据库管理模块
 * @author 陈凤庆
 * @description 管理 SQLite 数据库连接和数据操作
 */

const (
	// CurrentDatabaseVersion 当前数据库版本
	// 每次数据库结构变更时,需要递增此版本号
	// 20250101 陈凤庆 版本5: 为password_items表添加tab_id字段
	// 20251001 陈凤庆 版本6: 将所有ID字段从INTEGER改为TEXT类型(GUID)
	// 20251001 陈凤庆 版本7: 表结构重构 - password_items改为accounts，tabs改为types，type字段改为typeid，删除created_by和updated_by字段
	// 20251003 陈凤庆 版本8: 为accounts表添加input_method字段，支持三种输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)
	// 20251004 陈凤庆 版本9: 扩展input_method字段支持第四种输入方式：4-底层键盘API输入（已废弃）
	// 20251005 陈凤庆 版本10: 扩展input_method字段支持第四种输入方式：4-键盘助手输入（删除原第4种底层键盘API）
	// 20251017 陈凤庆 版本11: 添加password_rules表，支持密码规则管理
	// 20251017 陈凤庆 版本12: 添加username_history表，支持用户名历史记录管理
	CurrentDatabaseVersion = 12
)

/**
 * DatabaseManager 数据库管理器
 */
type DatabaseManager struct {
	db       *sql.DB
	dbPath   string
	isOpened bool
}

/**
 * NewDatabaseManager 创建新的数据库管理器
 * @return *DatabaseManager 数据库管理器实例
 */
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		isOpened: false,
	}
}

/**
 * OpenDatabase 打开数据库连接
 * @param dbPath 数据库文件路径
 * @return error 错误信息
 */
func (dm *DatabaseManager) OpenDatabase(dbPath string) error {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	dm.db = db
	dm.dbPath = dbPath
	dm.isOpened = true

	log.Printf("数据库已打开: %s", dbPath)
	return nil
}

/**
 * CreateTables 创建数据库表
 * @param language 语言代码（用于初始化多语言数据）
 * @return error 错误信息
 * @description 使用sysinfo表管理数据库版本,支持自动升级
 * @modify 20251005 陈凤庆 添加语言参数，支持多语言初始数据
 */
func (dm *DatabaseManager) CreateTables(language string) error {
	if !dm.isOpened {
		return errors.New("数据库未打开")
	}

	// 1. 初始化sysinfo表
	sysInfoMgr := NewSysInfoManager(dm.db)
	if err := sysInfoMgr.InitSysInfoTable(); err != nil {
		return fmt.Errorf("初始化sysinfo表失败: %w", err)
	}

	// 2. 获取当前数据库版本
	currentVersion, err := sysInfoMgr.GetDatabaseVersion()
	if err != nil {
		return fmt.Errorf("获取数据库版本失败: %w", err)
	}

	log.Printf("当前数据库版本: %d, 目标版本: %d", currentVersion, CurrentDatabaseVersion)

	// 3. 判断是否需要初始化或升级
	if currentVersion == 0 {
		// 检查是否是真正的新数据库
		var tableCount int
		err = dm.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tableCount)
		if err != nil {
			return fmt.Errorf("检查表数量失败: %w", err)
		}

		if tableCount <= 1 { // 只有sysinfo表或没有表
			// 新数据库,创建所有表
			log.Println("检测到新数据库,开始初始化...")
			if err := dm.createInitialSchema(language); err != nil {
				return err
			}
			// 设置版本号为当前版本
			if err := sysInfoMgr.SetDatabaseVersion(CurrentDatabaseVersion); err != nil {
				return err
			}
			log.Printf("数据库初始化完成,版本: %d", CurrentDatabaseVersion)
		} else {
			// 旧数据库(没有版本信息),从版本1开始升级
			log.Println("检测到旧数据库,开始升级...")
			if err := sysInfoMgr.SetDatabaseVersion(1); err != nil {
				return err
			}
			if err := dm.upgradeDatabase(1, CurrentDatabaseVersion); err != nil {
				return fmt.Errorf("数据库升级失败: %w", err)
			}
			log.Printf("数据库升级完成,版本: %d", CurrentDatabaseVersion)
		}
	} else if currentVersion < CurrentDatabaseVersion {
		// 需要升级
		log.Printf("开始升级数据库: 从版本 %d 到版本 %d", currentVersion, CurrentDatabaseVersion)
		if err := dm.upgradeDatabase(currentVersion, CurrentDatabaseVersion); err != nil {
			return fmt.Errorf("数据库升级失败: %w", err)
		}
		log.Printf("数据库升级完成,版本: %d", CurrentDatabaseVersion)
	} else {
		log.Printf("数据库版本已是最新: %d", currentVersion)
	}

	log.Println("数据库表创建完成")
	return nil
}

/**
 * createInitialSchema 创建初始数据库架构
 * @param language 语言代码（用于初始化多语言数据）
 * @return error 错误信息
 * @description 20251002 陈凤庆 版本7: 最新数据库结构，包含sysinfo表
 * @modify 20251005 陈凤庆 添加语言参数，支持多语言初始数据
 */
func (dm *DatabaseManager) createInitialSchema(language string) error {
	// 1. 创建sysinfo表（如果不存在）
	sysInfoSQL := `
	CREATE TABLE IF NOT EXISTS sysinfo (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 2. 创建密码库配置表(使用GUID)
	vaultConfigSQL := `
	CREATE TABLE IF NOT EXISTS vault_config (
		id TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL,
		salt TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 3. 创建分组表(使用GUID)
	// 20251002 陈凤庆 删除parent_id字段，不需要层级结构
	groupsSQL := `
	CREATE TABLE IF NOT EXISTS groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		icon TEXT DEFAULT '',
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 4. 创建类型表(使用GUID)
	// 20251002 陈凤庆 group_id不设置默认值，必须由前端传递
	typesSQL := `
	CREATE TABLE IF NOT EXISTS types (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		icon TEXT DEFAULT '',
		filter TEXT DEFAULT '',
		group_id TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES groups(id)
	);`

	// 5. 创建账号表(使用GUID)
	// 20251002 陈凤庆 删除tab_id字段，删除group_id默认值，通过typeid关联
	// 20251003 陈凤庆 添加input_method字段，支持三种输入方式：1-默认方式、2-模拟键盘输入、3-复制粘贴输入
	accountsSQL := `
	CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		username TEXT DEFAULT '',
		password TEXT DEFAULT '',
		url TEXT DEFAULT '',
		typeid TEXT NOT NULL,
		notes TEXT DEFAULT '',
		icon TEXT DEFAULT '',
		is_favorite BOOLEAN DEFAULT FALSE,
		use_count INTEGER DEFAULT 0,
		last_used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		input_method INTEGER DEFAULT 1,
		FOREIGN KEY (typeid) REFERENCES types(id)
	);`

	// 6. 创建密码规则表(使用GUID)
	// 20251017 陈凤庆 添加密码规则表，支持密码规则管理
	passwordRulesSQL := `
	CREATE TABLE IF NOT EXISTS password_rules (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT DEFAULT '',
		rule_type TEXT NOT NULL CHECK (rule_type IN ('general', 'custom')),
		config TEXT NOT NULL,
		is_default BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 7. 创建用户名历史记录表
	// 20251017 陈凤庆 添加用户名历史记录表，支持加密存储历史用户名
	usernameHistorySQL := `
	CREATE TABLE IF NOT EXISTS username_history (
		id INTEGER PRIMARY KEY DEFAULT 1,
		encrypted_data TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 执行建表语句
	tables := []string{sysInfoSQL, vaultConfigSQL, groupsSQL, typesSQL, accountsSQL, passwordRulesSQL, usernameHistorySQL}
	for _, tableSQL := range tables {
		if _, err := dm.db.Exec(tableSQL); err != nil {
			return fmt.Errorf("创建数据库表失败: %w", err)
		}
	}

	// 插入默认数据，传递语言参数
	if err := dm.insertDefaultData(language); err != nil {
		return fmt.Errorf("插入默认数据失败: %w", err)
	}

	return nil
}

/**
 * upgradeDatabase 升级数据库
 * @param fromVersion 起始版本
 * @param toVersion 目标版本
 * @return error 错误信息
 * @description 按顺序执行各个版本的升级方法
 */
func (dm *DatabaseManager) upgradeDatabase(fromVersion int, toVersion int) error {
	sysInfoMgr := NewSysInfoManager(dm.db)
	upgradeUtils := NewUpgradeUtils(dm.db)

	// 按顺序执行升级
	for version := fromVersion + 1; version <= toVersion; version++ {
		log.Printf("执行升级: 版本 %d -> 版本 %d", version-1, version)

		var err error
		switch version {
		case 7:
			// 20251001 陈凤庆 版本7: 表结构重构
			err = dm.dbUpgrade_v7(upgradeUtils)
		case 8:
			// 20251003 陈凤庆 版本8: 为accounts表添加input_method字段
			err = dm.dbUpgrade_v8(upgradeUtils)
		case 9:
			// 20251004 陈凤庆 版本9: 扩展input_method字段支持第四种输入方式
			err = dm.dbUpgrade_v9(upgradeUtils)
		case 10:
			// 20251005 陈凤庆 版本10: 扩展input_method字段支持第五种输入方式
			err = dm.dbUpgrade_v10(upgradeUtils)
		case 11:
			// 20251017 陈凤庆 版本11: 添加password_rules表
			err = dm.dbUpgrade_v11(upgradeUtils)
		case 12:
			// 20251017 陈凤庆 版本12: 添加username_history表
			err = dm.dbUpgrade_v12(upgradeUtils)
		// 未来版本在这里添加
		// case 13:
		//     err = dm.dbUpgrade_v13(upgradeUtils)
		default:
			// 20251002 陈凤庆 不再支持v7之前的版本升级
			return fmt.Errorf("不支持从版本 %d 升级，请使用最新版本创建新的数据库", version-1)
		}

		if err != nil {
			return fmt.Errorf("升级到版本 %d 失败: %w", version, err)
		}

		// 更新版本号
		if err := sysInfoMgr.SetDatabaseVersion(version); err != nil {
			return fmt.Errorf("更新版本号失败: %w", err)
		}

		log.Printf("升级到版本 %d 完成", version)
	}

	return nil
}

// 20251002 陈凤庆 删除旧版本升级代码，不再支持v7之前的版本升级

/**
 * insertDefaultData 插入默认数据
 * @param language 语言代码（用于初始化多语言数据）
 * @return error 错误信息
 * @modify 20251002 陈凤庆 直接设置数据库版本为最新版本，无需升级
 * @modify 20251005 陈凤庆 添加语言参数，支持多语言初始数据
 */
func (dm *DatabaseManager) insertDefaultData(language string) error {
	// 20251002 陈凤庆 直接设置数据库版本为最新版本
	sysInfoMgr := NewSysInfoManager(dm.db)
	err := sysInfoMgr.SetDatabaseVersion(CurrentDatabaseVersion)
	if err != nil {
		return fmt.Errorf("设置数据库版本失败: %w", err)
	}
	log.Printf("数据库版本设置为最新版本: %d", CurrentDatabaseVersion)

	// 获取多语言翻译
	translations := getDefaultDataTranslations(language)

	// 生成默认分组ID(使用固定的GUID,便于识别)
	defaultGroupID := "00000000-0000-4000-8000-000000000001"

	// 检查是否已有默认分组
	var count int
	err = dm.db.QueryRow("SELECT COUNT(*) FROM groups WHERE id = ?", defaultGroupID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// 插入默认分组，使用多语言名称
		// 20251002 陈凤庆 删除parent_id字段
		_, err = dm.db.Exec(`
			INSERT INTO groups (id, name, icon, sort_order, created_at, updated_at)
			VALUES (?, ?, 'fa-folder-open', 0, datetime('now'), datetime('now'))
		`, defaultGroupID, translations.DefaultGroupName)
		if err != nil {
			return err
		}
	}

	// 检查默认分组是否已有类型
	// 20251001 陈凤庆 tabs表改名为types
	err = dm.db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = ?", defaultGroupID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// 为默认分组插入默认类型(使用GUID)，使用多语言名称
		// 20251002 陈凤庆 删除created_by和updated_by字段
		defaultTypes := []struct {
			name      string
			icon      string
			sortOrder int
		}{
			{translations.WebsiteType, "fa-globe", 0},
			{translations.ApplicationType, "fa-mobile", 1},
			{translations.WiFiType, "fa-wifi", 2},
		}

		var typeIDs []string
		for _, typeItem := range defaultTypes {
			typeID := utils.GenerateGUID()
			typeIDs = append(typeIDs, typeID)
			_, err = dm.db.Exec(`
				INSERT INTO types (id, name, icon, group_id, sort_order, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
			`, typeID, typeItem.name, typeItem.icon, defaultGroupID, typeItem.sortOrder)
			if err != nil {
				return err
			}
		}

		// 20251002 陈凤庆 添加演示账号
		if len(typeIDs) > 0 {
			// 使用第一个类型（网站账号）创建演示账号
			demoAccountID := utils.GenerateGUID()
			_, err = dm.db.Exec(`
				INSERT INTO accounts (id, title, username, password, url, typeid, notes, icon, is_favorite, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
			`, demoAccountID, "演示账号", "demo@example.com", "demo123456", "https://example.com", typeIDs[0], "这是一个演示账号，您可以修改或删除它", "fa-user", false)
			if err != nil {
				return err
			}
		}
	}

	// 20251017 陈凤庆 取消自动初始化默认密码规则，让用户自行决定是否创建密码规则
	// err = dm.initializeDefaultPasswordRules()
	// if err != nil {
	//     return fmt.Errorf("初始化默认密码规则失败: %w", err)
	// }

	return nil
}

/**
 * CreateDefaultTypeForGroup 为指定分组创建默认类型
 * @param groupID 分组ID
 * @return error 错误信息
 * @modify 20251001 陈凤庆 tabs表改名为types，删除created_by和updated_by字段
 */
func (dm *DatabaseManager) CreateDefaultTypeForGroup(groupID string) error {
	if !dm.isOpened {
		return errors.New("数据库未打开")
	}

	// 检查该分组是否已有类型
	var count int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = ?", groupID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// 为新分组创建默认的"网站账号"类型(使用GUID)
		typeID := utils.GenerateGUID()
		_, err = dm.db.Exec(`
			INSERT INTO types (id, name, icon, group_id, sort_order)
			VALUES (?, '网站账号', 'fa-globe', ?, 0)
		`, typeID, groupID)
		if err != nil {
			return fmt.Errorf("创建默认类型失败: %w", err)
		}
		log.Printf("已为分组 %s 创建默认类型", groupID)
	}

	return nil
}

/**
 * EnsureDataIntegrity 确保数据完整性
 * @return error 错误信息
 * @description 20250101 陈凤庆 检查并初始化默认分组和标签
 */
func (dm *DatabaseManager) EnsureDataIntegrity() error {
	if !dm.isOpened {
		return errors.New("数据库未打开")
	}

	log.Println("[数据完整性检查] 开始检查数据完整性...")

	// 1. 检查是否有分组
	var groupCount int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&groupCount)
	if err != nil {
		return fmt.Errorf("检查分组数量失败: %w", err)
	}

	log.Printf("[数据完整性检查] 当前分组数量: %d", groupCount)

	// 2. 如果没有分组，创建默认分组
	if groupCount == 0 {
		log.Println("[数据完整性检查] 没有分组，创建默认分组...")
		// 20251001 陈凤庆 使用固定的GUID作为默认分组ID
		defaultGroupID := "00000000-0000-4000-8000-000000000001"
		now := time.Now()

		// 20251002 陈凤庆 删除parent_id、created_by、updated_by字段
		_, err = dm.db.Exec(`
			INSERT INTO groups (id, name, icon, sort_order, created_at, updated_at)
			VALUES (?, '默认', 'fa-folder-open', 0, ?, ?)
		`, defaultGroupID, now, now)
		if err != nil {
			return fmt.Errorf("创建默认分组失败: %w", err)
		}
		log.Printf("[数据完整性检查] 已创建默认分组 (ID: %s)", defaultGroupID)

		// 为默认分组创建默认类型
		// 20251001 陈凤庆 方法名改为CreateDefaultTypeForGroup
		if err := dm.CreateDefaultTypeForGroup(defaultGroupID); err != nil {
			return fmt.Errorf("为默认分组创建类型失败: %w", err)
		}
	} else {
		// 3. 检查每个分组是否都有类型
		// 20251001 陈凤庆 标签改为类型，tabs表改为types
		log.Println("[数据完整性检查] 检查每个分组的类型...")
		rows, err := dm.db.Query("SELECT id, name FROM groups")
		if err != nil {
			return fmt.Errorf("查询分组失败: %w", err)
		}

		// 20250101 陈凤庆 先收集所有需要创建类型的分组，避免在遍历rows时执行INSERT导致数据库锁
		// 20251001 陈凤庆 修改为使用string类型的GUID
		type groupInfo struct {
			id   string
			name string
		}
		var groupsNeedTabs []groupInfo

		for rows.Next() {
			var groupID string
			var groupName string
			if err := rows.Scan(&groupID, &groupName); err != nil {
				rows.Close()
				return fmt.Errorf("扫描分组数据失败: %w", err)
			}

			// 检查该分组是否有类型
			var typeCount int
			err = dm.db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = ?", groupID).Scan(&typeCount)
			if err != nil {
				rows.Close()
				return fmt.Errorf("检查分组%s的类型数量失败: %w", groupID, err)
			}

			if typeCount == 0 {
				groupsNeedTabs = append(groupsNeedTabs, groupInfo{id: groupID, name: groupName})
			} else {
				log.Printf("[数据完整性检查] 分组 '%s' (ID: %s) 有 %d 个类型", groupName, groupID, typeCount)
			}
		}
		rows.Close()

		// 为需要类型的分组创建默认类型
		// 20251001 陈凤庆 标签改为类型，方法名改为CreateDefaultTypeForGroup
		for _, g := range groupsNeedTabs {
			log.Printf("[数据完整性检查] 分组 '%s' (ID: %s) 没有类型，创建默认类型...", g.name, g.id)
			if err := dm.CreateDefaultTypeForGroup(g.id); err != nil {
				return fmt.Errorf("为分组%s创建类型失败: %w", g.id, err)
			}
		}
	}

	log.Println("[数据完整性检查] 数据完整性检查完成")
	return nil
}

/**
 * Close 关闭数据库连接
 */
func (dm *DatabaseManager) Close() {
	if dm.db != nil {
		dm.db.Close()
		dm.isOpened = false
		log.Println("数据库连接已关闭")
	}
}

/**
 * GetDB 获取数据库连接
 * @return *sql.DB 数据库连接
 */
func (dm *DatabaseManager) GetDB() *sql.DB {
	return dm.db
}

/**
 * IsOpened 检查数据库是否已打开
 * @return bool 是否已打开
 */
func (dm *DatabaseManager) IsOpened() bool {
	return dm.isOpened
}

/**
 * GetDatabasePath 获取数据库文件路径
 * @return string 数据库文件路径
 */
func (dm *DatabaseManager) GetDatabasePath() string {
	return dm.dbPath
}

// 密码库配置相关操作

/**
 * SaveVaultConfig 保存密码库配置
 * @param config 密码库配置
 * @return error 错误信息
 */
func (dm *DatabaseManager) SaveVaultConfig(config *models.VaultConfig) error {
	if !dm.isOpened {
		return errors.New("数据库未打开")
	}

	// 检查是否已存在配置
	var count int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM vault_config").Scan(&count)
	if err != nil {
		return err
	}

	now := time.Now()
	if count == 0 {
		// 插入新配置
		// 20251001 陈凤庆 使用GUID作为主键
		configID := utils.GenerateGUID()
		_, err = dm.db.Exec(`
			INSERT INTO vault_config (id, password_hash, salt, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
		`, configID, config.PasswordHash, config.Salt, now, now)
	} else {
		// 更新现有配置
		// 20251001 陈凤庆 更新第一条记录（因为只有一条配置记录）
		_, err = dm.db.Exec(`
			UPDATE vault_config
			SET password_hash = ?, salt = ?, updated_at = ?
			WHERE id = (SELECT id FROM vault_config LIMIT 1)
		`, config.PasswordHash, config.Salt, now)
	}

	return err
}

/**
 * GetVaultConfig 获取密码库配置
 * @return *models.VaultConfig 密码库配置
 * @return error 错误信息
 */
func (dm *DatabaseManager) GetVaultConfig() (*models.VaultConfig, error) {
	if !dm.isOpened {
		return nil, errors.New("数据库未打开")
	}

	config := &models.VaultConfig{}
	err := dm.db.QueryRow(`
		SELECT id, password_hash, salt, created_at, updated_at 
		FROM vault_config 
		ORDER BY id LIMIT 1
	`).Scan(&config.ID, &config.PasswordHash, &config.Salt, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有配置
		}
		return nil, err
	}

	return config, nil
}

// 20251002 陈凤庆 删除所有旧版本升级相关的辅助方法
/**
 * dbUpgrade_v7 升级到版本7
 * @param utils 升级工具
 * @return error 错误信息
 * @description 表结构重构：password_items改为accounts，tabs改为types，type字段改为typeid，删除created_by和updated_by字段
 * @author 陈凤庆
 * @date 20251001
 */

/**
 * joinStrings 连接字符串数组
 * @param strs 字符串数组
 * @param sep 分隔符
 * @return string 连接后的字符串
 */
func joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

/**
 * dbUpgrade_v7 升级到版本7
 * @param utils 升级工具
 * @return error 错误信息
 * @description 表结构重构：password_items改为accounts，tabs改为types，type字段改为typeid，删除created_by和updated_by字段
 * @author 陈凤庆
 * @date 20251001
 */
func (dm *DatabaseManager) dbUpgrade_v7(utils *UpgradeUtils) error {
	log.Println("开始执行版本7升级: 表结构重构")

	// 1. 重命名表：password_items -> accounts
	if err := dm.renameTableWithDataMigration("password_items", "accounts", map[string]string{
		"type": "typeid", // 字段重命名
	}, []string{"created_by", "updated_by"}); err != nil {
		return fmt.Errorf("重命名password_items表失败: %w", err)
	}

	// 2. 重命名表：tabs -> types，删除type字段和created_by、updated_by字段
	if err := dm.renameTableWithDataMigration("tabs", "types", map[string]string{},
		[]string{"type", "created_by", "updated_by"}); err != nil {
		return fmt.Errorf("重命名tabs表失败: %w", err)
	}

	// 3. 删除groups表的created_by和updated_by字段
	if err := dm.removeColumnsFromTable("groups", []string{"created_by", "updated_by"}); err != nil {
		return fmt.Errorf("删除groups表字段失败: %w", err)
	}

	log.Println("版本7升级完成: 表结构重构完成")
	return nil
}

/**
 * dbUpgrade_v8 升级到版本8
 * @param utils 升级工具
 * @return error 错误信息
 * @description 为accounts表添加input_method字段，支持三种输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)
 * @author 陈凤庆
 * @date 20251003
 */
func (dm *DatabaseManager) dbUpgrade_v8(utils *UpgradeUtils) error {
	log.Println("开始执行版本8升级: 为accounts表添加input_method字段")

	// 为accounts表添加input_method字段，默认值为1（默认方式）
	if err := utils.AddColumn("accounts", "input_method", "INTEGER DEFAULT 1"); err != nil {
		return fmt.Errorf("添加input_method字段失败: %w", err)
	}

	// 更新现有账号的输入方式为默认方式（1）
	if err := utils.ExecuteSQL("UPDATE accounts SET input_method = 1 WHERE input_method IS NULL"); err != nil {
		return fmt.Errorf("更新现有账号输入方式失败: %w", err)
	}

	log.Println("版本8升级完成: input_method字段添加完成")
	return nil
}

/**
 * dbUpgrade_v9 升级到版本9
 * @param utils 升级工具
 * @return error 错误信息
 * @description 扩展input_method字段支持第四种输入方式：4-底层键盘API输入
 * @author 陈凤庆
 * @date 20251004
 */
func (dm *DatabaseManager) dbUpgrade_v9(utils *UpgradeUtils) error {
	log.Println("开始执行版本9升级: 扩展input_method字段支持第四种输入方式")

	// 版本9主要是扩展input_method字段的取值范围，不需要修改数据库结构
	// input_method字段已经是INTEGER类型，可以直接支持值4
	// 这里只需要更新文档说明，实际不需要执行SQL操作

	log.Println("版本9升级完成: input_method字段现在支持第四种输入方式（4-底层键盘API输入）")
	return nil
}

/**
 * dbUpgrade_v10 升级到版本10
 * @param utils 升级工具
 * @return error 错误信息
 * @description 扩展input_method字段支持第四种输入方式：4-键盘助手输入（删除原第4种底层键盘API）
 * @author 陈凤庆
 * @date 20251005
 */
func (dm *DatabaseManager) dbUpgrade_v10(utils *UpgradeUtils) error {
	log.Println("开始执行版本10升级: 扩展input_method字段支持第四种输入方式（键盘助手输入）")

	// 版本10主要是扩展input_method字段的取值范围，不需要修改数据库结构
	// input_method字段已经是INTEGER类型，可以直接支持值4
	// 这里只需要更新文档说明，实际不需要执行SQL操作

	log.Println("版本10升级完成: input_method字段现在支持第四种输入方式（4-键盘助手输入）")
	return nil
}

/**
 * dbUpgrade_v11 升级到版本11
 * @param utils 升级工具
 * @return error 错误信息
 * @description 添加password_rules表，支持密码规则管理
 * @author 陈凤庆
 * @date 20251017
 */
func (dm *DatabaseManager) dbUpgrade_v11(utils *UpgradeUtils) error {
	log.Println("开始执行版本11升级: 添加password_rules表")

	// 创建密码规则表
	passwordRulesSQL := `
	CREATE TABLE IF NOT EXISTS password_rules (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT DEFAULT '',
		rule_type TEXT NOT NULL CHECK (rule_type IN ('general', 'custom')),
		config TEXT NOT NULL,
		is_default BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := dm.db.Exec(passwordRulesSQL)
	if err != nil {
		return fmt.Errorf("创建password_rules表失败: %w", err)
	}

	log.Println("版本11升级完成: password_rules表创建成功")
	return nil
}

/**
 * dbUpgrade_v12 升级到版本12
 * @param utils 升级工具
 * @return error 错误信息
 * @description 添加username_history表，支持用户名历史记录管理
 * @author 陈凤庆
 * @date 20251017
 */
func (dm *DatabaseManager) dbUpgrade_v12(utils *UpgradeUtils) error {
	log.Println("开始执行版本12升级: 添加username_history表")

	// 创建用户名历史记录表
	usernameHistorySQL := `
	CREATE TABLE IF NOT EXISTS username_history (
		id INTEGER PRIMARY KEY DEFAULT 1,
		encrypted_data TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := dm.db.Exec(usernameHistorySQL)
	if err != nil {
		return fmt.Errorf("创建username_history表失败: %w", err)
	}

	log.Println("版本12升级完成: username_history表创建成功")
	return nil
}

/**
 * renameTableWithDataMigration 重命名表并进行数据迁移
 * @param oldTableName 旧表名
 * @param newTableName 新表名
 * @param fieldRenames 字段重命名映射 (旧字段名 -> 新字段名)
 * @param fieldsToRemove 要删除的字段列表
 * @return error 错误信息
 */
func (dm *DatabaseManager) renameTableWithDataMigration(oldTableName, newTableName string, fieldRenames map[string]string, fieldsToRemove []string) error {
	// 1. 检查旧表是否存在
	var count int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", oldTableName).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查表%s是否存在失败: %w", oldTableName, err)
	}
	if count == 0 {
		log.Printf("表%s不存在，跳过重命名", oldTableName)
		return nil
	}

	// 2. 获取旧表结构
	var createSQL string
	err = dm.db.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name=?", oldTableName).Scan(&createSQL)
	if err != nil {
		return fmt.Errorf("获取表%s结构失败: %w", oldTableName, err)
	}

	// 3. 读取所有数据
	dataRows, err := dm.db.Query(fmt.Sprintf("SELECT * FROM %s", oldTableName))
	if err != nil {
		return fmt.Errorf("读取表%s数据失败: %w", oldTableName, err)
	}

	columns, err := dataRows.Columns()
	if err != nil {
		dataRows.Close()
		return fmt.Errorf("获取列名失败: %w", err)
	}

	var allData []map[string]interface{}
	for dataRows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := dataRows.Scan(valuePtrs...); err != nil {
			dataRows.Close()
			return fmt.Errorf("扫描数据失败: %w", err)
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}
		allData = append(allData, rowData)
	}
	dataRows.Close()

	// 4. 修改建表SQL
	newCreateSQL := dm.modifyCreateSQL(createSQL, oldTableName, newTableName, fieldRenames, fieldsToRemove)

	// 5. 删除旧表
	if _, err := dm.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", oldTableName)); err != nil {
		return fmt.Errorf("删除旧表%s失败: %w", oldTableName, err)
	}

	// 6. 创建新表
	if _, err := dm.db.Exec(newCreateSQL); err != nil {
		return fmt.Errorf("创建新表%s失败: %w", newTableName, err)
	}

	// 7. 迁移数据
	if len(allData) > 0 {
		// 构建新的列名列表（排除要删除的字段，应用字段重命名）
		var newColumns []string
		for _, col := range columns {
			// 检查是否要删除此字段
			shouldRemove := false
			for _, removeField := range fieldsToRemove {
				if col == removeField {
					shouldRemove = true
					break
				}
			}
			if shouldRemove {
				continue
			}

			// 应用字段重命名
			newCol := col
			if newName, exists := fieldRenames[col]; exists {
				newCol = newName
			}
			newColumns = append(newColumns, newCol)
		}

		for _, rowData := range allData {
			var values []interface{}
			for _, col := range columns {
				// 跳过要删除的字段
				shouldRemove := false
				for _, removeField := range fieldsToRemove {
					if col == removeField {
						shouldRemove = true
						break
					}
				}
				if shouldRemove {
					continue
				}
				values = append(values, rowData[col])
			}

			placeholders := make([]string, len(newColumns))
			for i := range placeholders {
				placeholders[i] = "?"
			}

			insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
				newTableName,
				joinStrings(newColumns, ", "),
				joinStrings(placeholders, ", "))

			if _, err := dm.db.Exec(insertSQL, values...); err != nil {
				return fmt.Errorf("插入数据到新表%s失败: %w", newTableName, err)
			}
		}
	}

	log.Printf("表%s重命名为%s完成，共迁移%d条记录", oldTableName, newTableName, len(allData))
	return nil
}

/**
 * modifyCreateSQL 修改建表SQL
 * @param createSQL 原始建表SQL
 * @param oldTableName 旧表名
 * @param newTableName 新表名
 * @param fieldRenames 字段重命名映射
 * @param fieldsToRemove 要删除的字段列表
 * @return string 修改后的建表SQL
 */
func (dm *DatabaseManager) modifyCreateSQL(createSQL, oldTableName, newTableName string, fieldRenames map[string]string, fieldsToRemove []string) string {
	result := createSQL

	// 1. 替换表名
	result = strings.Replace(result, oldTableName, newTableName, 1)

	// 2. 应用字段重命名
	for oldField, newField := range fieldRenames {
		// 替换字段定义中的字段名
		result = strings.ReplaceAll(result, oldField+" TEXT", newField+" TEXT")
		result = strings.ReplaceAll(result, oldField+" INTEGER", newField+" INTEGER")
		result = strings.ReplaceAll(result, oldField+" BOOLEAN", newField+" BOOLEAN")
		result = strings.ReplaceAll(result, oldField+" DATETIME", newField+" DATETIME")
	}

	// 3. 删除指定字段
	for _, fieldToRemove := range fieldsToRemove {
		// 使用正则表达式删除整行字段定义
		lines := strings.Split(result, "\n")
		var newLines []string
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			// 检查是否是要删除的字段行
			if strings.HasPrefix(trimmedLine, fieldToRemove+" ") {
				continue // 跳过这一行
			}
			newLines = append(newLines, line)
		}
		result = strings.Join(newLines, "\n")
	}

	return result
}

/**
 * removeColumnsFromTable 从表中删除指定列
 * @param tableName 表名
 * @param columnsToRemove 要删除的列名列表
 * @return error 错误信息
 */
func (dm *DatabaseManager) removeColumnsFromTable(tableName string, columnsToRemove []string) error {
	// SQLite不支持直接删除列，需要通过重建表的方式
	return dm.renameTableWithDataMigration(tableName, tableName+"_temp", map[string]string{}, columnsToRemove)
}

/**
 * InitializeDefaultPasswordRules 初始化默认密码规则（公开方法）
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251017
 */
func (dm *DatabaseManager) InitializeDefaultPasswordRules() error {
	return dm.initializeDefaultPasswordRules()
}

/**
 * initializeDefaultPasswordRules 初始化默认密码规则
 * @return error 错误信息
 * @author 陈凤庆
 * @date 20251017
 */
func (dm *DatabaseManager) initializeDefaultPasswordRules() error {
	// 检查是否已存在默认密码规则
	var count int
	err := dm.db.QueryRow("SELECT COUNT(*) FROM password_rules WHERE is_default = 1").Scan(&count)
	if err != nil {
		return fmt.Errorf("检查默认密码规则失败: %w", err)
	}

	if count > 0 {
		log.Println("默认密码规则已存在，跳过初始化")
		return nil
	}

	// 创建默认通用规则配置
	defaultConfig := `{"include_uppercase":true,"include_lowercase":true,"include_numbers":true,"include_special_chars":false,"include_custom_chars":false,"min_uppercase":1,"min_lowercase":1,"min_numbers":1,"min_special_chars":0,"min_custom_chars":0,"length":12,"custom_special_chars":""}`

	// 插入默认密码规则
	defaultRuleID := utils.GenerateGUID()
	_, err = dm.db.Exec(`
		INSERT INTO password_rules (id, name, description, rule_type, config, is_default, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, defaultRuleID, "通用规则", "系统预置的通用密码规则", "general", defaultConfig, true)

	if err != nil {
		return fmt.Errorf("插入默认密码规则失败: %w", err)
	}

	log.Println("默认密码规则初始化成功")
	return nil
}
