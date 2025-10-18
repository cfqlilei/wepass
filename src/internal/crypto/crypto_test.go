package crypto

import (
	"testing"
)

/**
 * 加密模块单元测试
 * @author 陈凤庆
 * @description 测试加密解密功能的正确性和安全性
 */

func TestCryptoManager_GenerateSalt(t *testing.T) {
	cm := NewCryptoManager()

	salt1, err := cm.GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	salt2, err := cm.GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	// 检查盐值长度
	if len(salt1) != SaltLength {
		t.Errorf("盐值长度不正确，期望 %d，实际 %d", SaltLength, len(salt1))
	}

	// 检查两次生成的盐值不同
	if string(salt1) == string(salt2) {
		t.Error("两次生成的盐值相同，随机性不足")
	}
}

func TestCryptoManager_HashPassword(t *testing.T) {
	cm := NewCryptoManager()

	salt, err := cm.GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	password := "Test246!Asd"
	hash1 := cm.HashPassword(password, salt)
	hash2 := cm.HashPassword(password, salt)

	// 相同密码和盐值应该产生相同的哈希
	if hash1 != hash2 {
		t.Error("相同密码和盐值产生了不同的哈希")
	}

	// 不同盐值应该产生不同的哈希
	salt2, _ := cm.GenerateSalt()
	hash3 := cm.HashPassword(password, salt2)
	if hash1 == hash3 {
		t.Error("不同盐值产生了相同的哈希")
	}
}

func TestCryptoManager_VerifyPassword(t *testing.T) {
	cm := NewCryptoManager()

	salt, err := cm.GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	password := "Test246!Asd"
	hash := cm.HashPassword(password, salt)

	// 正确密码应该验证通过
	if !cm.VerifyPassword(password, hash, salt) {
		t.Error("正确密码验证失败")
	}

	// 错误密码应该验证失败
	if cm.VerifyPassword("wrongpassword", hash, salt) {
		t.Error("错误密码验证通过")
	}
}

func TestCryptoManager_EncryptDecrypt(t *testing.T) {
	cm := NewCryptoManager()

	// 设置主密钥
	salt, err := cm.GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}
	cm.SetMasterPassword("Test246!Asd", salt)

	testCases := []string{
		"",
		"简单文本",
		"包含特殊字符!@#$%^&*()",
		"很长的文本内容，包含中文和英文以及数字123456789",
		"emoji测试🔐🔑🛡️",
	}

	for _, plaintext := range testCases {
		// 加密
		ciphertext, err := cm.Encrypt(plaintext)
		if err != nil {
			t.Errorf("加密失败，明文: %s, 错误: %v", plaintext, err)
			continue
		}

		// 解密
		decrypted, err := cm.Decrypt(ciphertext)
		if err != nil {
			t.Errorf("解密失败，密文: %s, 错误: %v", ciphertext, err)
			continue
		}

		// 验证结果
		if decrypted != plaintext {
			t.Errorf("加密解密结果不匹配，原文: %s, 解密: %s", plaintext, decrypted)
		}
	}
}

func TestMaskString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "*"},
		{"ab", "**"},
		{"abc", "a*c"},
		{"abcd", "a**d"},
		{"abcde", "a***e"},
		{"abcdef", "a****f"},
		{"github_user", "gi*******er"},
		{"email@gmail.com", "em***********om"},
	}

	for _, tc := range testCases {
		result := MaskString(tc.input)
		if result != tc.expected {
			t.Errorf("脱敏结果不正确，输入: %s, 期望: %s, 实际: %s", tc.input, tc.expected, result)
		}
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	testCases := []struct {
		length           int
		includeUppercase bool
		includeLowercase bool
		includeNumbers   bool
		includeSymbols   bool
		shouldError      bool
	}{
		{8, true, true, true, false, false},
		{12, true, true, true, true, false},
		{0, true, true, true, true, true},     // 长度为0应该报错
		{6, false, false, false, false, true}, // 没有字符类型应该报错
	}

	for _, tc := range testCases {
		password, err := GenerateRandomPassword(
			tc.length, tc.includeUppercase, tc.includeLowercase,
			tc.includeNumbers, tc.includeSymbols,
		)

		if tc.shouldError {
			if err == nil {
				t.Errorf("期望出错但没有出错，参数: %+v", tc)
			}
		} else {
			if err != nil {
				t.Errorf("生成密码失败: %v, 参数: %+v", err, tc)
				continue
			}

			if len(password) != tc.length {
				t.Errorf("密码长度不正确，期望: %d, 实际: %d", tc.length, len(password))
			}
		}
	}
}
