package database

import (
	"os"
	"path/filepath"
	"testing"
)

/**
 * 数据库升级测试
 * @author 陈凤庆
 * @description 测试数据库版本管理和升级功能
 */

/**
 * TestDatabaseUpgrade_NewDatabase 测试新数据库创建
 */
func TestDatabaseUpgrade_NewDatabase(t *testing.T) {
	// 创建临时数据库文件
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_new.db")

	dm := NewDatabaseManager()

	// 打开数据库
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	// 创建表
	err = dm.CreateTables()
	if err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	// 检查版本号
	db := dm.GetDB()
	sysInfoMgr := NewSysInfoManager(db)
	version, err := sysInfoMgr.GetDatabaseVersion()
	if err != nil {
		t.Fatalf("获取版本号失败: %v", err)
	}

	if version != CurrentDatabaseVersion {
		t.Errorf("版本号不正确，期望: %d, 实际: %d", CurrentDatabaseVersion, version)
	}

	// 检查types表是否有group_id字段
	// 20251002 陈凤庆 tabs表改为types表
	var columnExists bool
	err = db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM pragma_table_info('types')
		WHERE name='group_id'
	`).Scan(&columnExists)
	if err != nil {
		t.Fatalf("检查字段失败: %v", err)
	}

	if !columnExists {
		t.Error("types表缺少group_id字段")
	}

	// 检查默认分组的类型
	var typeCount int
	err = db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = '00000000-0000-4000-8000-000000000001'").Scan(&typeCount)
	if err != nil {
		t.Fatalf("查询类型数量失败: %v", err)
	}

	// 20251002 陈凤庆 检查默认类型数量
	if typeCount != 3 {
		t.Errorf("默认分组类型数量不正确，期望: 3, 实际: %d", typeCount)
	}
}

/**
 * TestDatabaseUpgrade_OldDatabase 测试旧数据库升级
 */
func TestDatabaseUpgrade_OldDatabase(t *testing.T) {
	// 创建临时数据库文件
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_old.db")

	dm := NewDatabaseManager()

	// 打开数据库
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	db := dm.GetDB()

	// 创建旧版本的表结构(没有group_id字段)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			parent_id INTEGER DEFAULT 0,
			icon TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT DEFAULT '陈凤庆',
			updated_by TEXT DEFAULT '陈凤庆'
		)
	`)
	if err != nil {
		t.Fatalf("创建groups表失败: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tabs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			icon TEXT DEFAULT '',
			type TEXT DEFAULT 'default',
			filter TEXT DEFAULT '',
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by TEXT DEFAULT '陈凤庆',
			updated_by TEXT DEFAULT '陈凤庆'
		)
	`)
	if err != nil {
		t.Fatalf("创建tabs表失败: %v", err)
	}

	// 插入默认分组
	_, err = db.Exec(`
		INSERT INTO groups (id, name, parent_id, icon, sort_order) 
		VALUES (1, '默认', 0, 'fa-folder-open', 0)
	`)
	if err != nil {
		t.Fatalf("插入默认分组失败: %v", err)
	}

	// 插入一些旧版本的页签(没有group_id)
	_, err = db.Exec(`
		INSERT INTO tabs (name, icon, type, sort_order) 
		VALUES ('网站账号', 'fa-globe', 'default', 0)
	`)
	if err != nil {
		t.Fatalf("插入页签失败: %v", err)
	}

	// 关闭数据库
	dm.Close()

	// 重新打开数据库并执行升级
	dm2 := NewDatabaseManager()
	err = dm2.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("重新打开数据库失败: %v", err)
	}
	defer dm2.Close()

	err = dm2.CreateTables()
	if err != nil {
		t.Fatalf("升级数据库失败: %v", err)
	}

	// 检查版本号
	db2 := dm2.GetDB()
	sysInfoMgr := NewSysInfoManager(db2)
	version, err := sysInfoMgr.GetDatabaseVersion()
	if err != nil {
		t.Fatalf("获取版本号失败: %v", err)
	}

	if version != CurrentDatabaseVersion {
		t.Errorf("版本号不正确，期望: %d, 实际: %d", CurrentDatabaseVersion, version)
	}

	// 检查tabs表是否有group_id字段
	var columnExists bool
	err = db2.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('tabs') 
		WHERE name='group_id'
	`).Scan(&columnExists)
	if err != nil {
		t.Fatalf("检查字段失败: %v", err)
	}

	if !columnExists {
		t.Error("升级后tabs表缺少group_id字段")
	}

	// 检查旧页签是否正确关联到默认分组
	var groupID int64
	err = db2.QueryRow("SELECT group_id FROM tabs WHERE name = '网站账号'").Scan(&groupID)
	if err != nil {
		t.Fatalf("查询页签group_id失败: %v", err)
	}

	if groupID != 1 {
		t.Errorf("旧页签的group_id不正确，期望: 1, 实际: %d", groupID)
	}
}

/**
 * TestCreateDefaultTabForGroup 测试为新分组创建默认页签
 */
func TestCreateDefaultTabForGroup(t *testing.T) {
	// 创建临时数据库文件
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_default_tab.db")

	dm := NewDatabaseManager()

	// 打开数据库
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	// 创建表
	err = dm.CreateTables()
	if err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	db := dm.GetDB()

	// 创建一个新分组
	// 20251002 陈凤庆 删除parent_id字段
	testGroupID := "00000000-0000-4000-8000-000000000002"
	_, err = db.Exec(`
		INSERT INTO groups (id, name, icon, sort_order, created_at, updated_at)
		VALUES (?, '测试分组', 'fa-folder', 1, datetime('now'), datetime('now'))
	`, testGroupID)
	if err != nil {
		t.Fatalf("创建分组失败: %v", err)
	}

	// 为新分组创建默认类型
	// 20251001 陈凤庆 CreateDefaultTabForGroup方法改为CreateDefaultTypeForGroup
	err = dm.CreateDefaultTypeForGroup(testGroupID)
	if err != nil {
		t.Fatalf("创建默认类型失败: %v", err)
	}

	// 检查类型是否创建成功
	// 20251001 陈凤庆 tabs表改为types
	var typeCount int
	err = db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = ?", testGroupID).Scan(&typeCount)
	if err != nil {
		t.Fatalf("查询类型数量失败: %v", err)
	}

	if typeCount != 1 {
		t.Errorf("新分组类型数量不正确，期望: 1, 实际: %d", typeCount)
	}

	// 检查类型名称
	var typeName string
	err = db.QueryRow("SELECT name FROM types WHERE group_id = ?", testGroupID).Scan(&typeName)
	if err != nil {
		t.Fatalf("查询类型名称失败: %v", err)
	}

	if typeName != "网站账号" {
		t.Errorf("类型名称不正确，期望: 网站账号, 实际: %s", typeName)
	}
}

/**
 * TestMain 测试入口
 */
func TestMain(m *testing.M) {
	// 运行测试
	code := m.Run()

	// 退出
	os.Exit(code)
}
