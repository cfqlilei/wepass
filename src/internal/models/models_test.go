package models

import (
	"encoding/json"
	"testing"
	"time"
)

/**
 * 数据模型单元测试
 * @author 陈凤庆
 * @description 测试数据模型的序列化和反序列化功能
 */

func TestGroup_JSONSerialization(t *testing.T) {
	// 创建测试分组
	// 20251001 陈凤庆 删除created_by和updated_by字段，ID改为string类型
	now := time.Now()
	// 20251002 陈凤庆 删除ParentID字段
	group := Group{
		ID:        "00000000-0000-4000-8000-000000000001",
		Name:      "测试分组",
		Icon:      "folder",
		SortOrder: 1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 测试序列化
	jsonData, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("分组序列化失败: %v", err)
	}

	// 测试反序列化
	var deserializedGroup Group
	err = json.Unmarshal(jsonData, &deserializedGroup)
	if err != nil {
		t.Fatalf("分组反序列化失败: %v", err)
	}

	// 验证数据一致性
	// 20251001 陈凤庆 ID改为string类型，格式化符号改为%s
	if deserializedGroup.ID != group.ID {
		t.Errorf("分组ID不一致，期望: %s, 实际: %s", group.ID, deserializedGroup.ID)
	}

	if deserializedGroup.Name != group.Name {
		t.Errorf("分组名称不一致，期望: %s, 实际: %s", group.Name, deserializedGroup.Name)
	}

	// 20251002 陈凤庆 删除ParentID字段检查
}

func TestPasswordItem_JSONSerialization(t *testing.T) {
	// 创建测试账号
	now := time.Now()
	// 20251001 陈凤庆 Type字段改为TypeID，删除created_by和updated_by字段
	// 20251002 陈凤庆 删除GroupID和TabID字段
	passwordItem := PasswordItem{
		ID:         "test-guid-1",
		Title:      "测试网站",
		Username:   "encrypted_username",
		Password:   "encrypted_password",
		URL:        "encrypted_url",
		TypeID:     "encrypted_website",
		Notes:      "encrypted_notes",
		Icon:       "web",
		IsFavorite: true,
		UseCount:   5,
		LastUsedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 测试序列化
	jsonData, err := json.Marshal(passwordItem)
	if err != nil {
		t.Fatalf("账号序列化失败: %v", err)
	}

	// 测试反序列化
	var deserializedItem PasswordItem
	err = json.Unmarshal(jsonData, &deserializedItem)
	if err != nil {
		t.Fatalf("账号反序列化失败: %v", err)
	}

	// 验证数据一致性
	// 20251001 陈凤庆 ID改为string类型，格式化符号改为%s
	if deserializedItem.ID != passwordItem.ID {
		t.Errorf("账号ID不一致，期望: %s, 实际: %s", passwordItem.ID, deserializedItem.ID)
	}

	if deserializedItem.Title != passwordItem.Title {
		t.Errorf("标题不一致，期望: %s, 实际: %s", passwordItem.Title, deserializedItem.Title)
	}

	if deserializedItem.IsFavorite != passwordItem.IsFavorite {
		t.Errorf("收藏状态不一致，期望: %t, 实际: %t", passwordItem.IsFavorite, deserializedItem.IsFavorite)
	}

	if deserializedItem.UseCount != passwordItem.UseCount {
		t.Errorf("使用次数不一致，期望: %d, 实际: %d", passwordItem.UseCount, deserializedItem.UseCount)
	}
}

func TestPasswordItemDecrypted_JSONSerialization(t *testing.T) {
	// 创建测试解密账号
	now := time.Now()
	// 20251002 陈凤庆 删除GroupID和TabID字段
	decryptedItem := PasswordItemDecrypted{
		ID:             "test-guid-2",
		Title:          "测试网站",
		Username:       "testuser",
		Password:       "testpass123",
		URL:            "https://test.com",
		TypeID:         "网站",
		Notes:          "测试备注",
		Icon:           "web",
		IsFavorite:     true,
		UseCount:       5,
		LastUsedAt:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
		MaskedUsername: "test****",
		MaskedPassword: "****",
	}

	// 测试序列化
	jsonData, err := json.Marshal(decryptedItem)
	if err != nil {
		t.Fatalf("解密账号序列化失败: %v", err)
	}

	// 测试反序列化
	var deserializedItem PasswordItemDecrypted
	err = json.Unmarshal(jsonData, &deserializedItem)
	if err != nil {
		t.Fatalf("解密账号反序列化失败: %v", err)
	}

	// 验证数据一致性
	if deserializedItem.Username != decryptedItem.Username {
		t.Errorf("用户名不一致，期望: %s, 实际: %s", decryptedItem.Username, deserializedItem.Username)
	}

	if deserializedItem.Password != decryptedItem.Password {
		t.Errorf("密码不一致，期望: %s, 实际: %s", decryptedItem.Password, deserializedItem.Password)
	}

	if deserializedItem.MaskedUsername != decryptedItem.MaskedUsername {
		t.Errorf("脱敏用户名不一致，期望: %s, 实际: %s", decryptedItem.MaskedUsername, deserializedItem.MaskedUsername)
	}
}

func TestVaultConfig_JSONSerialization(t *testing.T) {
	// 创建测试密码库配置
	now := time.Now()
	vaultConfig := VaultConfig{
		ID:           "vault-guid-1",
		PasswordHash: "hashed_password",
		Salt:         "random_salt",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 测试序列化
	jsonData, err := json.Marshal(vaultConfig)
	if err != nil {
		t.Fatalf("密码库配置序列化失败: %v", err)
	}

	// 测试反序列化
	var deserializedConfig VaultConfig
	err = json.Unmarshal(jsonData, &deserializedConfig)
	if err != nil {
		t.Fatalf("密码库配置反序列化失败: %v", err)
	}

	// 验证数据一致性
	if deserializedConfig.PasswordHash != vaultConfig.PasswordHash {
		t.Errorf("密码哈希不一致，期望: %s, 实际: %s", vaultConfig.PasswordHash, deserializedConfig.PasswordHash)
	}

	if deserializedConfig.Salt != vaultConfig.Salt {
		t.Errorf("盐值不一致，期望: %s, 实际: %s", vaultConfig.Salt, deserializedConfig.Salt)
	}
}

func TestAppConfig_JSONSerialization(t *testing.T) {
	// 创建测试应用配置
	appConfig := AppConfig{
		CurrentVaultPath: "/path/to/vault.db",
		RecentVaults:     []string{"/path/to/vault1.db", "/path/to/vault2.db"},
		WindowWidth:      1024,
		WindowHeight:     768,
		Theme:            "dark",
		Language:         "zh-CN",
	}

	// 测试序列化
	jsonData, err := json.Marshal(appConfig)
	if err != nil {
		t.Fatalf("应用配置序列化失败: %v", err)
	}

	// 测试反序列化
	var deserializedConfig AppConfig
	err = json.Unmarshal(jsonData, &deserializedConfig)
	if err != nil {
		t.Fatalf("应用配置反序列化失败: %v", err)
	}

	// 验证数据一致性
	if deserializedConfig.CurrentVaultPath != appConfig.CurrentVaultPath {
		t.Errorf("当前密码库路径不一致，期望: %s, 实际: %s", appConfig.CurrentVaultPath, deserializedConfig.CurrentVaultPath)
	}

	if len(deserializedConfig.RecentVaults) != len(appConfig.RecentVaults) {
		t.Errorf("最近使用密码库数量不一致，期望: %d, 实际: %d", len(appConfig.RecentVaults), len(deserializedConfig.RecentVaults))
	}

	if deserializedConfig.WindowWidth != appConfig.WindowWidth {
		t.Errorf("窗口宽度不一致，期望: %d, 实际: %d", appConfig.WindowWidth, deserializedConfig.WindowWidth)
	}

	if deserializedConfig.Theme != appConfig.Theme {
		t.Errorf("主题不一致，期望: %s, 实际: %s", appConfig.Theme, deserializedConfig.Theme)
	}
}

func TestSearchResult_JSONSerialization(t *testing.T) {
	// 创建测试搜索结果
	now := time.Now()
	// 20251002 陈凤庆 删除ParentID、GroupID和TabID字段
	searchResult := SearchResult{
		Type: "mixed",
		Groups: []Group{
			{
				ID:        "group-guid-3",
				Name:      "测试分组",
				CreatedAt: now,
			},
		},
		Accounts: []AccountDecrypted{
			{
				ID:        "password-guid-3",
				Title:     "测试密码",
				Username:  "testuser",
				TypeID:    "type-guid-3",
				CreatedAt: now,
			},
		},
	}

	// 测试序列化
	jsonData, err := json.Marshal(searchResult)
	if err != nil {
		t.Fatalf("搜索结果序列化失败: %v", err)
	}

	// 测试反序列化
	var deserializedResult SearchResult
	err = json.Unmarshal(jsonData, &deserializedResult)
	if err != nil {
		t.Fatalf("搜索结果反序列化失败: %v", err)
	}

	// 验证数据一致性
	if deserializedResult.Type != searchResult.Type {
		t.Errorf("搜索类型不一致，期望: %s, 实际: %s", searchResult.Type, deserializedResult.Type)
	}

	if len(deserializedResult.Groups) != len(searchResult.Groups) {
		t.Errorf("分组数量不一致，期望: %d, 实际: %d", len(searchResult.Groups), len(deserializedResult.Groups))
	}

	// 20251001 陈凤庆 Passwords字段改为Accounts
	if len(deserializedResult.Accounts) != len(searchResult.Accounts) {
		t.Errorf("账号数量不一致，期望: %d, 实际: %d", len(searchResult.Accounts), len(deserializedResult.Accounts))
	}
}
