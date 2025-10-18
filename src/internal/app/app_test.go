package app

import (
	"context"
	"path/filepath"
	"testing"

	"wepassword/internal/config"
	"wepassword/internal/database"
)

/**
 * App模块单元测试
 * @author 陈凤庆
 * @description 测试应用核心功能和生命周期管理
 */

func TestNewApp(t *testing.T) {
	// 创建配置管理器
	configManager := config.NewConfigManager()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager()

	// 创建应用实例
	app := NewApp(configManager, dbManager)

	// 检查应用实例是否正确创建
	if app == nil {
		t.Fatal("应用实例创建失败")
	}

	if app.configManager != configManager {
		t.Error("配置管理器设置不正确")
	}

	if app.dbManager != dbManager {
		t.Error("数据库管理器设置不正确")
	}
}

func TestApp_Startup(t *testing.T) {
	// 创建配置管理器
	configManager := config.NewConfigManager()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager()

	// 创建应用实例
	app := NewApp(configManager, dbManager)

	// 创建上下文
	ctx := context.Background()

	// 测试启动
	app.Startup(ctx)

	// 检查上下文是否设置
	if app.ctx != ctx {
		t.Error("应用上下文设置不正确")
	}

	// 检查服务是否初始化
	if app.vaultService == nil {
		t.Error("密码库服务未初始化")
	}

	// 20251002 陈凤庆 passwordService改名为accountService
	if app.accountService == nil {
		t.Error("账号服务未初始化")
	}

	if app.groupService == nil {
		t.Error("分组服务未初始化")
	}
}

func TestApp_GetAppInfo(t *testing.T) {
	// 创建配置管理器
	configManager := config.NewConfigManager()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager()

	// 创建应用实例
	app := NewApp(configManager, dbManager)

	// 获取应用信息
	info := app.GetAppInfo()

	// 检查应用信息
	if info == nil {
		t.Fatal("应用信息获取失败")
	}

	// 检查必要字段
	if name, ok := info["name"]; !ok || name != "wepass" {
		t.Error("应用名称不正确")
	}

	if version, ok := info["version"]; !ok || version == "" {
		t.Error("应用版本信息缺失")
	}

	if author, ok := info["author"]; !ok || author != "陈凤庆" {
		t.Error("应用作者信息不正确")
	}
}

func TestApp_VaultOperations(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.db")

	// 创建配置管理器
	configManager := config.NewConfigManager()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager()

	// 创建应用实例
	app := NewApp(configManager, dbManager)

	// 启动应用
	ctx := context.Background()
	app.Startup(ctx)

	// 测试检查密码库是否存在（不存在的情况）
	if app.CheckVaultExists(vaultPath) {
		t.Error("不存在的密码库检查结果应该为false")
	}

	// 测试创建密码库
	// 20251001 陈凤庆 CreateVault现在返回两个值
	password := "Test246!Asd"
	_, err := app.CreateVault(vaultPath, password)
	if err != nil {
		t.Fatalf("创建密码库失败: %v", err)
	}

	// 测试检查密码库是否存在（存在的情况）
	if !app.CheckVaultExists(vaultPath) {
		t.Error("存在的密码库检查结果应该为true")
	}

	// 关闭当前连接
	app.Shutdown(ctx)

	// 测试打开密码库
	app2 := NewApp(configManager, dbManager)
	app2.Startup(ctx)

	err = app2.OpenVault(vaultPath, password)
	if err != nil {
		t.Fatalf("打开密码库失败: %v", err)
	}

	// 清理
	app2.Shutdown(ctx)
}

func TestApp_Lifecycle(t *testing.T) {
	// 创建配置管理器
	configManager := config.NewConfigManager()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager()

	// 创建应用实例
	app := NewApp(configManager, dbManager)

	// 创建上下文
	ctx := context.Background()

	// 测试启动
	app.Startup(ctx)

	// 测试DomReady
	app.DomReady(ctx)

	// 测试BeforeClose
	prevent := app.BeforeClose(ctx)
	if prevent {
		t.Error("BeforeClose不应该阻止关闭")
	}

	// 测试Shutdown
	app.Shutdown(ctx)
}
