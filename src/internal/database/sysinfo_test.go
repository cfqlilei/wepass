package database

import (
	"path/filepath"
	"testing"
)

/**
 * 系统信息管理测试
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 测试sysinfo表的管理功能
 */

/**
 * TestSysInfo_InitTable 测试初始化sysinfo表
 */
func TestSysInfo_InitTable(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_sysinfo.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	sysInfoMgr := NewSysInfoManager(dm.GetDB())

	// 初始化sysinfo表
	err = sysInfoMgr.InitSysInfoTable()
	if err != nil {
		t.Fatalf("初始化sysinfo表失败: %v", err)
	}

	// 检查表是否存在
	utils := NewUpgradeUtils(dm.GetDB())
	exists, err := utils.TableExists("sysinfo")
	if err != nil {
		t.Fatalf("检查表是否存在失败: %v", err)
	}

	if !exists {
		t.Error("sysinfo表不存在")
	}
}

/**
 * TestSysInfo_SetAndGetValue 测试设置和获取配置值
 */
func TestSysInfo_SetAndGetValue(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_sysinfo.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	sysInfoMgr := NewSysInfoManager(dm.GetDB())
	err = sysInfoMgr.InitSysInfoTable()
	if err != nil {
		t.Fatalf("初始化sysinfo表失败: %v", err)
	}

	// 设置配置值
	err = sysInfoMgr.SetValue("test_key", "test_value")
	if err != nil {
		t.Fatalf("设置配置值失败: %v", err)
	}

	// 获取配置值
	value, err := sysInfoMgr.GetValue("test_key")
	if err != nil {
		t.Fatalf("获取配置值失败: %v", err)
	}

	if value != "test_value" {
		t.Errorf("配置值不正确，期望: test_value, 实际: %s", value)
	}

	// 更新配置值
	err = sysInfoMgr.SetValue("test_key", "new_value")
	if err != nil {
		t.Fatalf("更新配置值失败: %v", err)
	}

	// 再次获取配置值
	value, err = sysInfoMgr.GetValue("test_key")
	if err != nil {
		t.Fatalf("获取配置值失败: %v", err)
	}

	if value != "new_value" {
		t.Errorf("配置值不正确，期望: new_value, 实际: %s", value)
	}
}

/**
 * TestSysInfo_DatabaseVersion 测试数据库版本管理
 */
func TestSysInfo_DatabaseVersion(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_sysinfo.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	sysInfoMgr := NewSysInfoManager(dm.GetDB())
	err = sysInfoMgr.InitSysInfoTable()
	if err != nil {
		t.Fatalf("初始化sysinfo表失败: %v", err)
	}

	// 获取初始版本号(应该为0)
	version, err := sysInfoMgr.GetDatabaseVersion()
	if err != nil {
		t.Fatalf("获取数据库版本失败: %v", err)
	}

	if version != 0 {
		t.Errorf("初始版本号不正确，期望: 0, 实际: %d", version)
	}

	// 初始化版本号
	err = sysInfoMgr.InitDatabaseVersion()
	if err != nil {
		t.Fatalf("初始化数据库版本失败: %v", err)
	}

	// 再次获取版本号(应该为1)
	version, err = sysInfoMgr.GetDatabaseVersion()
	if err != nil {
		t.Fatalf("获取数据库版本失败: %v", err)
	}

	if version != 1 {
		t.Errorf("版本号不正确，期望: 1, 实际: %d", version)
	}

	// 设置版本号为2
	err = sysInfoMgr.SetDatabaseVersion(2)
	if err != nil {
		t.Fatalf("设置数据库版本失败: %v", err)
	}

	// 获取版本号(应该为2)
	version, err = sysInfoMgr.GetDatabaseVersion()
	if err != nil {
		t.Fatalf("获取数据库版本失败: %v", err)
	}

	if version != 2 {
		t.Errorf("版本号不正确，期望: 2, 实际: %d", version)
	}
}

/**
 * TestSysInfo_KeyOperations 测试键操作
 */
func TestSysInfo_KeyOperations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_sysinfo.db")

	dm := NewDatabaseManager()
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	sysInfoMgr := NewSysInfoManager(dm.GetDB())
	err = sysInfoMgr.InitSysInfoTable()
	if err != nil {
		t.Fatalf("初始化sysinfo表失败: %v", err)
	}

	// 设置多个键值对
	err = sysInfoMgr.SetValue("key1", "value1")
	if err != nil {
		t.Fatalf("设置配置值失败: %v", err)
	}

	err = sysInfoMgr.SetValue("key2", "value2")
	if err != nil {
		t.Fatalf("设置配置值失败: %v", err)
	}

	// 检查键是否存在
	exists, err := sysInfoMgr.KeyExists("key1")
	if err != nil {
		t.Fatalf("检查键是否存在失败: %v", err)
	}

	if !exists {
		t.Error("key1应该存在")
	}

	// 获取所有键值对
	allKeys, err := sysInfoMgr.GetAllKeys()
	if err != nil {
		t.Fatalf("获取所有键值对失败: %v", err)
	}

	if len(allKeys) != 2 {
		t.Errorf("键值对数量不正确，期望: 2, 实际: %d", len(allKeys))
	}

	// 删除键
	err = sysInfoMgr.DeleteKey("key1")
	if err != nil {
		t.Fatalf("删除键失败: %v", err)
	}

	// 再次检查键是否存在
	exists, err = sysInfoMgr.KeyExists("key1")
	if err != nil {
		t.Fatalf("检查键是否存在失败: %v", err)
	}

	if exists {
		t.Error("key1不应该存在")
	}
}
