package database

import (
	"path/filepath"
	"testing"
)

/**
 * 数据库升级工具测试
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 测试数据库升级工具的各项功能
 */

/**
 * TestUpgradeUtils_TableExists 测试检查表是否存在
 */
func TestUpgradeUtils_TableExists(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_utils.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	utils := NewUpgradeUtils(dm.GetDB())

	// 创建测试表
	_, err = dm.GetDB().Exec(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		t.Fatalf("创建测试表失败: %v", err)
	}

	// 检查表是否存在
	exists, err := utils.TableExists("test_table")
	if err != nil {
		t.Fatalf("检查表是否存在失败: %v", err)
	}

	if !exists {
		t.Error("test_table应该存在")
	}

	// 检查不存在的表
	exists, err = utils.TableExists("non_existent_table")
	if err != nil {
		t.Fatalf("检查表是否存在失败: %v", err)
	}

	if exists {
		t.Error("non_existent_table不应该存在")
	}
}

/**
 * TestUpgradeUtils_ColumnExists 测试检查字段是否存在
 */
func TestUpgradeUtils_ColumnExists(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_utils.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	utils := NewUpgradeUtils(dm.GetDB())

	// 创建测试表
	_, err = dm.GetDB().Exec(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		t.Fatalf("创建测试表失败: %v", err)
	}

	// 检查字段是否存在
	exists, err := utils.ColumnExists("test_table", "name")
	if err != nil {
		t.Fatalf("检查字段是否存在失败: %v", err)
	}

	if !exists {
		t.Error("name字段应该存在")
	}

	// 检查不存在的字段
	exists, err = utils.ColumnExists("test_table", "non_existent_column")
	if err != nil {
		t.Fatalf("检查字段是否存在失败: %v", err)
	}

	if exists {
		t.Error("non_existent_column不应该存在")
	}
}

/**
 * TestUpgradeUtils_CreateTable 测试创建表
 */
func TestUpgradeUtils_CreateTable(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_utils.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	utils := NewUpgradeUtils(dm.GetDB())

	// 创建表
	createSQL := `
	CREATE TABLE IF NOT EXISTS test_table (
		id INTEGER PRIMARY KEY,
		name TEXT
	)`

	err = utils.CreateTable("test_table", createSQL)
	if err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	// 检查表是否存在
	exists, err := utils.TableExists("test_table")
	if err != nil {
		t.Fatalf("检查表是否存在失败: %v", err)
	}

	if !exists {
		t.Error("test_table应该存在")
	}

	// 再次创建表(应该跳过)
	err = utils.CreateTable("test_table", createSQL)
	if err != nil {
		t.Fatalf("再次创建表失败: %v", err)
	}
}

/**
 * TestUpgradeUtils_AddColumn 测试添加字段
 */
func TestUpgradeUtils_AddColumn(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_utils.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	utils := NewUpgradeUtils(dm.GetDB())

	// 创建测试表
	_, err = dm.GetDB().Exec(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			name TEXT
		)
	`)
	if err != nil {
		t.Fatalf("创建测试表失败: %v", err)
	}

	// 添加字段
	err = utils.AddColumn("test_table", "age", "INTEGER DEFAULT 0")
	if err != nil {
		t.Fatalf("添加字段失败: %v", err)
	}

	// 检查字段是否存在
	exists, err := utils.ColumnExists("test_table", "age")
	if err != nil {
		t.Fatalf("检查字段是否存在失败: %v", err)
	}

	if !exists {
		t.Error("age字段应该存在")
	}

	// 再次添加字段(应该跳过)
	err = utils.AddColumn("test_table", "age", "INTEGER DEFAULT 0")
	if err != nil {
		t.Fatalf("再次添加字段失败: %v", err)
	}
}

/**
 * TestUpgradeUtils_RenameColumn 测试重命名字段
 */
func TestUpgradeUtils_RenameColumn(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_utils.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	utils := NewUpgradeUtils(dm.GetDB())

	// 创建测试表
	_, err = dm.GetDB().Exec(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			old_name TEXT
		)
	`)
	if err != nil {
		t.Fatalf("创建测试表失败: %v", err)
	}

	// 重命名字段
	err = utils.RenameColumn("test_table", "old_name", "new_name")
	if err != nil {
		t.Fatalf("重命名字段失败: %v", err)
	}

	// 检查新字段是否存在
	exists, err := utils.ColumnExists("test_table", "new_name")
	if err != nil {
		t.Fatalf("检查字段是否存在失败: %v", err)
	}

	if !exists {
		t.Error("new_name字段应该存在")
	}

	// 检查旧字段是否不存在
	exists, err = utils.ColumnExists("test_table", "old_name")
	if err != nil {
		t.Fatalf("检查字段是否存在失败: %v", err)
	}

	if exists {
		t.Error("old_name字段不应该存在")
	}
}

/**
 * TestUpgradeUtils_ChangeColumnType 测试修改字段类型
 */
func TestUpgradeUtils_ChangeColumnType(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_utils.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	utils := NewUpgradeUtils(dm.GetDB())

	// 创建测试表
	_, err = dm.GetDB().Exec(`
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			value TEXT
		)
	`)
	if err != nil {
		t.Fatalf("创建测试表失败: %v", err)
	}

	// 插入测试数据
	_, err = dm.GetDB().Exec("INSERT INTO test_table (id, value) VALUES (1, '123')")
	if err != nil {
		t.Fatalf("插入测试数据失败: %v", err)
	}

	// 修改字段类型
	err = utils.ChangeColumnType("test_table", "value", "INTEGER")
	if err != nil {
		t.Fatalf("修改字段类型失败: %v", err)
	}

	// 检查新字段是否存在
	exists, err := utils.ColumnExists("test_table", "value")
	if err != nil {
		t.Fatalf("检查字段是否存在失败: %v", err)
	}

	if !exists {
		t.Error("value字段应该存在")
	}

	// 获取字段类型
	columnType, err := utils.GetColumnType("test_table", "value")
	if err != nil {
		t.Fatalf("获取字段类型失败: %v", err)
	}

	if columnType != "INTEGER" {
		t.Errorf("字段类型不正确，期望: INTEGER, 实际: %s", columnType)
	}
}
