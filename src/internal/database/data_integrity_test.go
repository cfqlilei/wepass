package database

import (
	"os"
	"testing"
)

/**
 * 数据完整性检查测试
 * @author 陈凤庆
 * @description 测试EnsureDataIntegrity方法
 */

/**
 * TestEnsureDataIntegrity_EmptyDatabase 测试空数据库的数据完整性检查
 */
func TestEnsureDataIntegrity_EmptyDatabase(t *testing.T) {
	// 创建临时数据库文件
	dbPath := "./test_integrity_empty.db"
	defer os.Remove(dbPath)

	// 创建数据库管理器
	dbManager := NewDatabaseManager()

	// 打开数据库
	if err := dbManager.OpenDatabase(dbPath); err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dbManager.Close()

	// 创建表结构
	if err := dbManager.CreateTables(); err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	// 删除默认数据（模拟空数据库）
	// 20251002 陈凤庆 tabs表改为types表
	db := dbManager.GetDB()
	if _, err := db.Exec("DELETE FROM types"); err != nil {
		t.Fatalf("删除types数据失败: %v", err)
	}
	if _, err := db.Exec("DELETE FROM groups"); err != nil {
		t.Fatalf("删除groups数据失败: %v", err)
	}

	// 执行数据完整性检查
	if err := dbManager.EnsureDataIntegrity(); err != nil {
		t.Fatalf("数据完整性检查失败: %v", err)
	}

	// 验证结果：应该有1个默认分组
	var groupCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&groupCount); err != nil {
		t.Fatalf("查询分组数量失败: %v", err)
	}
	if groupCount != 1 {
		t.Errorf("期望有1个分组，实际有%d个", groupCount)
	}

	// 验证默认分组的属性
	var groupID int64
	var groupName string
	if err := db.QueryRow("SELECT id, name FROM groups WHERE id = 1").Scan(&groupID, &groupName); err != nil {
		t.Fatalf("查询默认分组失败: %v", err)
	}
	if groupID != 1 {
		t.Errorf("期望分组ID为1，实际为%d", groupID)
	}
	if groupName != "默认" {
		t.Errorf("期望分组名称为'默认'，实际为'%s'", groupName)
	}

	// 验证结果：应该有1个默认标签
	var tabCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM tabs WHERE group_id = 1").Scan(&tabCount); err != nil {
		t.Fatalf("查询标签数量失败: %v", err)
	}
	if tabCount != 1 {
		t.Errorf("期望有1个标签，实际有%d个", tabCount)
	}

	// 验证默认标签的属性
	var tabName string
	var tabType string
	if err := db.QueryRow("SELECT name, type FROM tabs WHERE group_id = 1").Scan(&tabName, &tabType); err != nil {
		t.Fatalf("查询默认标签失败: %v", err)
	}
	if tabName != "网站账号" {
		t.Errorf("期望标签名称为'网站账号'，实际为'%s'", tabName)
	}
	if tabType != "default" {
		t.Errorf("期望标签类型为'default'，实际为'%s'", tabType)
	}
}

/**
 * TestEnsureDataIntegrity_GroupWithoutTabs 测试有分组但没有标签的情况
 */
func TestEnsureDataIntegrity_GroupWithoutTabs(t *testing.T) {
	// 创建临时数据库文件
	dbPath := "./test_integrity_no_tabs.db"
	defer os.Remove(dbPath)

	// 创建数据库管理器
	dbManager := NewDatabaseManager()

	// 打开数据库
	if err := dbManager.OpenDatabase(dbPath); err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dbManager.Close()

	// 创建表结构
	if err := dbManager.CreateTables(); err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	db := dbManager.GetDB()

	// 删除所有类型（模拟有分组但没有类型的情况）
	// 20251002 陈凤庆 tabs表改为types表
	if _, err := db.Exec("DELETE FROM types"); err != nil {
		t.Fatalf("删除types数据失败: %v", err)
	}

	// 添加一个测试分组
	// 20251002 陈凤庆 删除parent_id、created_by、updated_by字段
	if _, err := db.Exec(`
		INSERT INTO groups (id, name, icon, sort_order)
		VALUES ('00000000-0000-4000-8000-000000000002', '测试分组', 'fa-folder', 1)
	`); err != nil {
		t.Fatalf("插入测试分组失败: %v", err)
	}

	// 执行数据完整性检查
	if err := dbManager.EnsureDataIntegrity(); err != nil {
		t.Fatalf("数据完整性检查失败: %v", err)
	}

	// 验证结果：默认分组应该有标签
	var tabCount1 int
	if err := db.QueryRow("SELECT COUNT(*) FROM tabs WHERE group_id = 1").Scan(&tabCount1); err != nil {
		t.Fatalf("查询默认分组标签数量失败: %v", err)
	}
	if tabCount1 != 1 {
		t.Errorf("期望默认分组有1个标签，实际有%d个", tabCount1)
	}

	// 验证结果：测试分组应该有标签
	var tabCount2 int
	if err := db.QueryRow("SELECT COUNT(*) FROM tabs WHERE group_id = 2").Scan(&tabCount2); err != nil {
		t.Fatalf("查询测试分组标签数量失败: %v", err)
	}
	if tabCount2 != 1 {
		t.Errorf("期望测试分组有1个标签，实际有%d个", tabCount2)
	}
}

/**
 * TestEnsureDataIntegrity_CompleteData 测试数据完整的情况
 */
func TestEnsureDataIntegrity_CompleteData(t *testing.T) {
	// 创建临时数据库文件
	dbPath := "./test_integrity_complete.db"
	defer os.Remove(dbPath)

	// 创建数据库管理器
	dbManager := NewDatabaseManager()

	// 打开数据库
	if err := dbManager.OpenDatabase(dbPath); err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dbManager.Close()

	// 创建表结构（会自动创建默认数据）
	if err := dbManager.CreateTables(); err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	db := dbManager.GetDB()

	// 记录执行前的数据量
	// 20251002 陈凤庆 tabs表改为types表
	var groupCountBefore, typeCountBefore int
	if err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&groupCountBefore); err != nil {
		t.Fatalf("查询分组数量失败: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM types").Scan(&typeCountBefore); err != nil {
		t.Fatalf("查询类型数量失败: %v", err)
	}

	// 执行数据完整性检查
	if err := dbManager.EnsureDataIntegrity(); err != nil {
		t.Fatalf("数据完整性检查失败: %v", err)
	}

	// 验证结果：数据量应该没有变化
	// 20251002 陈凤庆 tabs表改为types表
	var groupCountAfter, typeCountAfter int
	if err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&groupCountAfter); err != nil {
		t.Fatalf("查询分组数量失败: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM types").Scan(&typeCountAfter); err != nil {
		t.Fatalf("查询类型数量失败: %v", err)
	}

	if groupCountBefore != groupCountAfter {
		t.Errorf("分组数量发生变化：执行前%d，执行后%d", groupCountBefore, groupCountAfter)
	}
	if typeCountBefore != typeCountAfter {
		t.Errorf("类型数量发生变化：执行前%d，执行后%d", typeCountBefore, typeCountAfter)
	}
}

/**
 * TestEnsureDataIntegrity_MultipleGroups 测试多个分组的情况
 */
func TestEnsureDataIntegrity_MultipleGroups(t *testing.T) {
	// 创建临时数据库文件
	dbPath := "./test_integrity_multiple.db"
	defer os.Remove(dbPath)

	// 创建数据库管理器
	dbManager := NewDatabaseManager()

	// 打开数据库
	if err := dbManager.OpenDatabase(dbPath); err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dbManager.Close()

	// 创建表结构
	if err := dbManager.CreateTables(); err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	db := dbManager.GetDB()

	// 删除所有类型
	// 20251002 陈凤庆 tabs表改为types表
	if _, err := db.Exec("DELETE FROM types"); err != nil {
		t.Fatalf("删除types数据失败: %v", err)
	}

	// 添加多个测试分组
	groups := []struct {
		id   string
		name string
	}{
		{"00000000-0000-4000-8000-000000000002", "工作"},
		{"00000000-0000-4000-8000-000000000003", "个人"},
		{"00000000-0000-4000-8000-000000000004", "学习"},
	}

	// 20251002 陈凤庆 删除parent_id、created_by、updated_by字段
	for i, g := range groups {
		if _, err := db.Exec(`
			INSERT INTO groups (id, name, icon, sort_order)
			VALUES (?, ?, 'fa-folder', ?)
		`, g.id, g.name, i+2); err != nil {
			t.Fatalf("插入分组'%s'失败: %v", g.name, err)
		}
	}

	// 执行数据完整性检查
	if err := dbManager.EnsureDataIntegrity(); err != nil {
		t.Fatalf("数据完整性检查失败: %v", err)
	}

	// 验证结果：每个分组都应该有类型
	// 20251002 陈凤庆 修改为统一的string类型，tabs表改为types表
	allGroups := append([]struct {
		id   string
		name string
	}{{"00000000-0000-4000-8000-000000000001", "默认"}}, groups...)

	for _, g := range allGroups {
		var typeCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM types WHERE group_id = ?", g.id).Scan(&typeCount); err != nil {
			t.Fatalf("查询分组'%s'的类型数量失败: %v", g.name, err)
		}
		if typeCount < 1 {
			t.Errorf("期望分组'%s'至少有1个类型，实际有%d个", g.name, typeCount)
		}
	}
}
