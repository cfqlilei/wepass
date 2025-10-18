package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"wepassword/internal/models"
)

/**
 * 数据库模块单元测试
 * @author 陈凤庆
 * @description 测试数据库连接、表创建和基本操作
 */

func TestDatabaseManager_OpenDatabase(t *testing.T) {
	// 创建临时数据库文件
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	dm := NewDatabaseManager()

	// 测试打开数据库
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}

	// 检查状态
	if !dm.IsOpened() {
		t.Error("数据库状态应该为已打开")
	}

	if dm.GetDatabasePath() != dbPath {
		t.Errorf("数据库路径不正确，期望: %s, 实际: %s", dbPath, dm.GetDatabasePath())
	}

	// 检查数据库文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("数据库文件未创建")
	}

	// 关闭数据库
	dm.Close()

	if dm.IsOpened() {
		t.Error("数据库状态应该为已关闭")
	}
}

func TestDatabaseManager_CreateTables(t *testing.T) {
	// 创建临时数据库文件
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

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

	// 检查表是否存在
	// 20251001 陈凤庆 password_items表改为accounts，tabs表改为types
	db := dm.GetDB()
	tables := []string{"vault_config", "groups", "accounts", "types"}

	for _, tableName := range tables {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
		if err != nil {
			t.Errorf("检查表 %s 失败: %v", tableName, err)
		}
		if count != 1 {
			t.Errorf("表 %s 不存在", tableName)
		}
	}

	// 检查默认数据是否插入
	// 20251002 陈凤庆 使用GUID检查默认分组
	var groupCount int
	err = db.QueryRow("SELECT COUNT(*) FROM groups WHERE id = '00000000-0000-4000-8000-000000000001'").Scan(&groupCount)
	if err != nil {
		t.Errorf("检查默认分组失败: %v", err)
	}
	if groupCount != 1 {
		t.Error("默认分组未创建")
	}
}

func TestDatabaseManager_VaultConfig(t *testing.T) {
	// 创建临时数据库文件
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	dm := NewDatabaseManager()

	// 打开数据库并创建表
	err := dm.OpenDatabase(dbPath)
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	defer dm.Close()

	err = dm.CreateTables()
	if err != nil {
		t.Fatalf("创建表失败: %v", err)
	}

	// 测试保存配置
	config := &models.VaultConfig{
		PasswordHash: "test_hash",
		Salt:         "test_salt",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = dm.SaveVaultConfig(config)
	if err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}

	// 测试获取配置
	savedConfig, err := dm.GetVaultConfig()
	if err != nil {
		t.Fatalf("获取配置失败: %v", err)
	}

	if savedConfig == nil {
		t.Fatal("获取的配置为空")
	}

	if savedConfig.PasswordHash != config.PasswordHash {
		t.Errorf("密码哈希不匹配，期望: %s, 实际: %s", config.PasswordHash, savedConfig.PasswordHash)
	}

	if savedConfig.Salt != config.Salt {
		t.Errorf("盐值不匹配，期望: %s, 实际: %s", config.Salt, savedConfig.Salt)
	}

	// 测试更新配置
	config.PasswordHash = "updated_hash"
	err = dm.SaveVaultConfig(config)
	if err != nil {
		t.Fatalf("更新配置失败: %v", err)
	}

	updatedConfig, err := dm.GetVaultConfig()
	if err != nil {
		t.Fatalf("获取更新后的配置失败: %v", err)
	}

	if updatedConfig.PasswordHash != "updated_hash" {
		t.Errorf("配置更新失败，期望: %s, 实际: %s", "updated_hash", updatedConfig.PasswordHash)
	}
}
