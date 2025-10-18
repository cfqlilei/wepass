package config

import (
	"os"
	"path/filepath"
	"testing"

	"wepassword/internal/models"
)

/**
 * 配置管理模块单元测试
 * @author 陈凤庆
 * @description 测试配置文件的读写和管理功能
 */

// createTestConfigManager 创建用于测试的配置管理器
func createTestConfigManager(t *testing.T) *ConfigManager {
	tempDir := t.TempDir()
	return &ConfigManager{
		configPath: filepath.Join(tempDir, "test_config.json"),
		config: &models.AppConfig{
			CurrentVaultPath: "",
			RecentVaults:     make([]string, 0),
			WindowWidth:      550,
			WindowHeight:     800,
			Theme:            "light",
			Language:         "zh-CN",
		},
	}
}

func TestNewConfigManager(t *testing.T) {
	cm := createTestConfigManager(t)
	
	if cm == nil {
		t.Fatal("配置管理器创建失败")
	}
	
	if cm.config == nil {
		t.Fatal("默认配置未初始化")
	}
	
	// 检查默认配置值
	if cm.config.WindowWidth != 550 {
		t.Errorf("默认窗口宽度不正确，期望: 550, 实际: %d", cm.config.WindowWidth)
	}
	
	if cm.config.WindowHeight != 800 {
		t.Errorf("默认窗口高度不正确，期望: 800, 实际: %d", cm.config.WindowHeight)
	}
	
	if cm.config.Theme != "light" {
		t.Errorf("默认主题不正确，期望: light, 实际: %s", cm.config.Theme)
	}
	
	if cm.config.Language != "zh-CN" {
		t.Errorf("默认语言不正确，期望: zh-CN, 实际: %s", cm.config.Language)
	}
}

func TestConfigManager_VaultPath(t *testing.T) {
	cm := createTestConfigManager(t)
	
	// 测试初始状态
	if cm.GetCurrentVaultPath() != "" {
		t.Error("初始密码库路径应该为空")
	}
	
	// 测试设置密码库路径
	testPath := "/test/vault.db"
	err := cm.SetCurrentVaultPath(testPath)
	if err != nil {
		t.Fatalf("设置密码库路径失败: %v", err)
	}
	
	if cm.GetCurrentVaultPath() != testPath {
		t.Errorf("密码库路径设置不正确，期望: %s, 实际: %s", testPath, cm.GetCurrentVaultPath())
	}
}

func TestConfigManager_RecentVaults(t *testing.T) {
	cm := createTestConfigManager(t)
	
	// 测试初始状态
	recentVaults := cm.GetRecentVaults()
	if len(recentVaults) != 0 {
		t.Error("初始最近使用的密码库列表应该为空")
	}
	
	// 测试添加最近使用的密码库
	vault1 := "/test/vault1.db"
	vault2 := "/test/vault2.db"
	vault3 := "/test/vault3.db"
	
	cm.addToRecentVaults(vault1)
	recentVaults = cm.GetRecentVaults()
	if len(recentVaults) != 1 || recentVaults[0] != vault1 {
		t.Error("添加第一个密码库失败")
	}
	
	cm.addToRecentVaults(vault2)
	recentVaults = cm.GetRecentVaults()
	if len(recentVaults) != 2 || recentVaults[0] != vault2 || recentVaults[1] != vault1 {
		t.Error("添加第二个密码库失败，应该按最新使用顺序排列")
	}
	
	// 测试重复添加
	cm.addToRecentVaults(vault1)
	recentVaults = cm.GetRecentVaults()
	if len(recentVaults) != 2 || recentVaults[0] != vault1 || recentVaults[1] != vault2 {
		t.Error("重复添加密码库应该移动到最前面")
	}
	
	// 测试超过最大数量
	for i := 0; i < 10; i++ {
		cm.addToRecentVaults(vault3)
	}
	recentVaults = cm.GetRecentVaults()
	if len(recentVaults) > 5 {
		t.Errorf("最近使用的密码库数量不应超过5个，实际: %d", len(recentVaults))
	}
}

func TestConfigManager_WindowSize(t *testing.T) {
	cm := createTestConfigManager(t)
	
	// 测试获取默认窗口大小
	width, height := cm.GetWindowSize()
	if width != 550 || height != 800 {
		t.Errorf("默认窗口大小不正确，期望: 550x800, 实际: %dx%d", width, height)
	}
	
	// 测试设置窗口大小
	newWidth, newHeight := 1024, 768
	err := cm.SetWindowSize(newWidth, newHeight)
	if err != nil {
		t.Fatalf("设置窗口大小失败: %v", err)
	}
	
	width, height = cm.GetWindowSize()
	if width != newWidth || height != newHeight {
		t.Errorf("窗口大小设置不正确，期望: %dx%d, 实际: %dx%d", newWidth, newHeight, width, height)
	}
}

func TestConfigManager_Theme(t *testing.T) {
	cm := createTestConfigManager(t)
	
	// 测试获取默认主题
	if cm.GetTheme() != "light" {
		t.Errorf("默认主题不正确，期望: light, 实际: %s", cm.GetTheme())
	}
	
	// 测试设置主题
	newTheme := "dark"
	err := cm.SetTheme(newTheme)
	if err != nil {
		t.Fatalf("设置主题失败: %v", err)
	}
	
	if cm.GetTheme() != newTheme {
		t.Errorf("主题设置不正确，期望: %s, 实际: %s", newTheme, cm.GetTheme())
	}
}

func TestConfigManager_Language(t *testing.T) {
	cm := createTestConfigManager(t)
	
	// 测试获取默认语言
	if cm.GetLanguage() != "zh-CN" {
		t.Errorf("默认语言不正确，期望: zh-CN, 实际: %s", cm.GetLanguage())
	}
	
	// 测试设置语言
	newLanguage := "en-US"
	err := cm.SetLanguage(newLanguage)
	if err != nil {
		t.Fatalf("设置语言失败: %v", err)
	}
	
	if cm.GetLanguage() != newLanguage {
		t.Errorf("语言设置不正确，期望: %s, 实际: %s", newLanguage, cm.GetLanguage())
	}
}

func TestConfigManager_SaveAndLoadConfig(t *testing.T) {
	// 创建配置管理器
	cm := createTestConfigManager(t)
	
	// 修改一些配置
	cm.SetWindowSize(1024, 768)
	cm.SetTheme("dark")
	cm.SetLanguage("en-US")
	
	// 通过SetCurrentVaultPath添加到最近使用列表
	cm.SetCurrentVaultPath("/test/vault1.db")
	cm.SetCurrentVaultPath("/test/vault2.db")
	cm.SetCurrentVaultPath("/test/vault.db")
	
	// 保存配置
	err := cm.SaveConfig()
	if err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}
	
	// 检查配置文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		t.Error("配置文件未创建")
	}
	
	// 创建新的配置管理器并加载配置
	cm2 := &ConfigManager{
		configPath: cm.configPath,
		config: &models.AppConfig{
			CurrentVaultPath: "",
			RecentVaults:     make([]string, 0),
			WindowWidth:      550,
			WindowHeight:     800,
			Theme:            "light",
			Language:         "zh-CN",
		},
	}
	
	err = cm2.LoadConfig()
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}
	
	// 验证配置是否正确加载
	if cm2.GetCurrentVaultPath() != "/test/vault.db" {
		t.Error("密码库路径加载不正确")
	}
	
	width, height := cm2.GetWindowSize()
	if width != 1024 || height != 768 {
		t.Error("窗口大小加载不正确")
	}
	
	if cm2.GetTheme() != "dark" {
		t.Error("主题加载不正确")
	}
	
	if cm2.GetLanguage() != "en-US" {
		t.Error("语言加载不正确")
	}
	
	recentVaults := cm2.GetRecentVaults()
	if len(recentVaults) != 3 || recentVaults[0] != "/test/vault.db" || recentVaults[1] != "/test/vault2.db" || recentVaults[2] != "/test/vault1.db" {
		t.Errorf("最近使用的密码库列表加载不正确，期望: [/test/vault.db, /test/vault2.db, /test/vault1.db]，实际: %v", recentVaults)
	}
}