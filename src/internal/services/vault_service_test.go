package services

import (
	"os"
	"path/filepath"
	"testing"

	"wepassword/internal/config"
	"wepassword/internal/database"
)

/**
 * 密码库服务集成测试
 * @author 陈凤庆
 * @description 测试密码库的创建、打开和验证功能
 */

func TestVaultService_CreateAndOpenVault(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.db")

	// 创建配置管理器（使用临时目录）
	configManager := config.NewConfigManager()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager()

	// 创建密码库服务
	vaultService := NewVaultService(dbManager, configManager)

	// 测试创建密码库
	password := "Test246!Asd"
	err := vaultService.CreateVault(vaultPath, password)
	if err != nil {
		t.Fatalf("创建密码库失败: %v", err)
	}

	// 检查文件是否存在
	if !vaultService.CheckVaultExists(vaultPath) {
		t.Error("密码库文件未创建")
	}

	// 关闭当前连接
	vaultService.CloseVault()

	// 测试打开密码库
	err = vaultService.OpenVault(vaultPath, password)
	if err != nil {
		t.Fatalf("打开密码库失败: %v", err)
	}

	// 检查状态
	if !vaultService.IsVaultOpened() {
		t.Error("密码库应该处于打开状态")
	}

	// 测试错误密码
	vaultService.CloseVault()
	err = vaultService.OpenVault(vaultPath, "wrongpassword")
	if err == nil {
		t.Error("错误密码应该导致打开失败")
	}

	// 测试不存在的文件
	err = vaultService.OpenVault("/nonexistent/path.db", password)
	if err == nil {
		t.Error("不存在的文件应该导致打开失败")
	}

	// 清理
	vaultService.CloseVault()
}

func TestVaultService_CheckVaultExists(t *testing.T) {
	vaultService := NewVaultService(nil, nil)

	// 测试空路径
	if vaultService.CheckVaultExists("") {
		t.Error("空路径应该返回 false")
	}

	// 测试不存在的文件
	if vaultService.CheckVaultExists("/nonexistent/path.db") {
		t.Error("不存在的文件应该返回 false")
	}

	// 创建临时文件
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.db")
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	file.Close()

	// 测试存在的文件
	if !vaultService.CheckVaultExists(tempFile) {
		t.Error("存在的文件应该返回 true")
	}
}
