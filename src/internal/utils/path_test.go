package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

/**
 * 跨平台路径工具函数测试
 * @author 陈凤庆
 * @date 2025-10-02
 * @description 测试跨平台路径工具函数在不同操作系统下的正确性
 */

func TestGetAppDataDir(t *testing.T) {
	appDataDir, err := GetAppDataDir()
	if err != nil {
		t.Fatalf("获取应用数据目录失败: %v", err)
	}

	if appDataDir == "" {
		t.Error("应用数据目录不能为空")
	}

	// 检查路径是否包含 wepass
	if !strings.Contains(appDataDir, "wepass") {
		t.Errorf("应用数据目录应该包含 'wepass'，实际路径: %s", appDataDir)
	}

	// 根据操作系统检查路径格式
	switch runtime.GOOS {
	case "darwin": // macOS
		expectedPattern := "Library/Application Support/wepass"
		if !strings.Contains(appDataDir, expectedPattern) {
			t.Errorf("macOS 应用数据目录应该包含 '%s'，实际路径: %s", expectedPattern, appDataDir)
		}
	case "windows": // Windows
		if !strings.Contains(appDataDir, "wepass") {
			t.Errorf("Windows 应用数据目录应该包含 'wepass'，实际路径: %s", appDataDir)
		}
		// Windows 路径应该包含 AppData 或使用 APPDATA 环境变量
		appData := os.Getenv("APPDATA")
		if appData != "" && !strings.Contains(appDataDir, appData) {
			t.Errorf("Windows 应用数据目录应该基于 APPDATA 环境变量，实际路径: %s", appDataDir)
		}
	default: // Linux 和其他 Unix 系统
		expectedPattern := ".config/wepass"
		if !strings.Contains(appDataDir, expectedPattern) {
			t.Errorf("Linux 应用数据目录应该包含 '%s'，实际路径: %s", expectedPattern, appDataDir)
		}
	}

	t.Logf("当前操作系统: %s, 应用数据目录: %s", runtime.GOOS, appDataDir)
}

func TestGetVaultDataDir(t *testing.T) {
	vaultDataDir, err := GetVaultDataDir()
	if err != nil {
		t.Fatalf("获取密码库数据目录失败: %v", err)
	}

	if vaultDataDir == "" {
		t.Error("密码库数据目录不能为空")
	}

	// 检查路径是否包含 wepass 和 vaults
	if !strings.Contains(vaultDataDir, "wepass") {
		t.Errorf("密码库数据目录应该包含 'wepass'，实际路径: %s", vaultDataDir)
	}

	if !strings.Contains(vaultDataDir, "vaults") {
		t.Errorf("密码库数据目录应该包含 'vaults'，实际路径: %s", vaultDataDir)
	}

	t.Logf("密码库数据目录: %s", vaultDataDir)
}

func TestGetDefaultVaultPath(t *testing.T) {
	testVaultName := "test_vault"
	
	vaultPath, err := GetDefaultVaultPath(testVaultName)
	if err != nil {
		t.Fatalf("获取默认密码库路径失败: %v", err)
	}

	if vaultPath == "" {
		t.Error("密码库路径不能为空")
	}

	// 检查路径是否包含密码库名称和 .db 后缀
	expectedFileName := testVaultName + ".db"
	if !strings.Contains(vaultPath, expectedFileName) {
		t.Errorf("密码库路径应该包含 '%s'，实际路径: %s", expectedFileName, vaultPath)
	}

	// 检查路径是否包含 vaults 目录
	if !strings.Contains(vaultPath, "vaults") {
		t.Errorf("密码库路径应该包含 'vaults' 目录，实际路径: %s", vaultPath)
	}

	// 检查文件扩展名
	if filepath.Ext(vaultPath) != ".db" {
		t.Errorf("密码库文件应该有 .db 扩展名，实际路径: %s", vaultPath)
	}

	t.Logf("默认密码库路径: %s", vaultPath)
}

func TestEnsureDir(t *testing.T) {
	// 创建临时目录进行测试
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test", "nested", "directory")

	// 测试创建嵌套目录
	err := EnsureDir(testDir)
	if err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("目录应该已经创建")
	}

	// 测试对已存在目录的处理
	err = EnsureDir(testDir)
	if err != nil {
		t.Errorf("对已存在目录调用 EnsureDir 不应该返回错误: %v", err)
	}

	t.Logf("测试目录: %s", testDir)
}

func TestPathIntegration(t *testing.T) {
	// 集成测试：测试完整的路径创建流程
	testVaultName := "integration_test_vault"
	
	// 获取默认密码库路径
	vaultPath, err := GetDefaultVaultPath(testVaultName)
	if err != nil {
		t.Fatalf("获取默认密码库路径失败: %v", err)
	}

	// 检查目录是否已创建
	vaultDir := filepath.Dir(vaultPath)
	if _, err := os.Stat(vaultDir); os.IsNotExist(err) {
		t.Error("密码库目录应该已经创建")
	}

	// 清理测试创建的目录（仅在测试环境中）
	if strings.Contains(vaultDir, "test") || strings.Contains(vaultDir, "Test") {
		os.RemoveAll(vaultDir)
		t.Logf("已清理测试目录: %s", vaultDir)
	} else {
		t.Logf("保留生产目录: %s", vaultDir)
	}
}

func TestCrossPlatformCompatibility(t *testing.T) {
	// 测试跨平台兼容性
	appDataDir, err := GetAppDataDir()
	if err != nil {
		t.Fatalf("获取应用数据目录失败: %v", err)
	}

	// 检查路径分隔符是否正确
	expectedSeparator := string(filepath.Separator)
	if !strings.Contains(appDataDir, expectedSeparator) && len(appDataDir) > 1 {
		t.Errorf("路径应该使用正确的分隔符 '%s'，实际路径: %s", expectedSeparator, appDataDir)
	}

	// 检查路径是否为绝对路径
	if !filepath.IsAbs(appDataDir) {
		t.Errorf("应用数据目录应该是绝对路径，实际路径: %s", appDataDir)
	}

	t.Logf("操作系统: %s, 路径分隔符: %s, 应用数据目录: %s", 
		runtime.GOOS, expectedSeparator, appDataDir)
}
