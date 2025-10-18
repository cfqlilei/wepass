package utils

import (
	"strings"
	"testing"
)

/**
 * GUID生成器测试
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 测试GUID生成器的各项功能
 */

/**
 * TestGenerateGUID 测试GUID生成
 */
func TestGenerateGUID(t *testing.T) {
	// 生成多个GUID，确保格式正确且不重复
	guids := make(map[string]bool)
	
	for i := 0; i < 1000; i++ {
		guid := GenerateGUID()
		
		// 检查格式
		if !IsValidGUID(guid) {
			t.Errorf("生成的GUID格式无效: %s", guid)
		}
		
		// 检查长度
		if len(guid) != 36 {
			t.Errorf("GUID长度不正确，期望36，实际%d: %s", len(guid), guid)
		}
		
		// 检查连字符位置
		if guid[8] != '-' || guid[13] != '-' || guid[18] != '-' || guid[23] != '-' {
			t.Errorf("GUID连字符位置不正确: %s", guid)
		}
		
		// 检查版本号（第13位应该是4）
		if guid[14] != '4' {
			t.Errorf("GUID版本号不正确，期望4，实际%c: %s", guid[14], guid)
		}
		
		// 检查变体位（第17位应该是8、9、a、b之一）
		variantChar := guid[19]
		if variantChar != '8' && variantChar != '9' && variantChar != 'a' && variantChar != 'b' &&
		   variantChar != 'A' && variantChar != 'B' {
			t.Errorf("GUID变体位不正确，期望8/9/a/b，实际%c: %s", variantChar, guid)
		}
		
		// 检查唯一性
		if guids[guid] {
			t.Errorf("生成了重复的GUID: %s", guid)
		}
		guids[guid] = true
	}
	
	t.Logf("成功生成并验证了%d个唯一GUID", len(guids))
}

/**
 * TestGenerateShortGUID 测试短GUID生成
 */
func TestGenerateShortGUID(t *testing.T) {
	// 生成多个短GUID
	shortGuids := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		shortGuid := GenerateShortGUID()
		
		// 检查格式
		if !IsValidShortGUID(shortGuid) {
			t.Errorf("生成的短GUID格式无效: %s", shortGuid)
		}
		
		// 检查长度
		if len(shortGuid) != 32 {
			t.Errorf("短GUID长度不正确，期望32，实际%d: %s", len(shortGuid), shortGuid)
		}
		
		// 检查唯一性
		if shortGuids[shortGuid] {
			t.Errorf("生成了重复的短GUID: %s", shortGuid)
		}
		shortGuids[shortGuid] = true
	}
	
	t.Logf("成功生成并验证了%d个唯一短GUID", len(shortGuids))
}

/**
 * TestIsValidGUID 测试GUID格式验证
 */
func TestIsValidGUID(t *testing.T) {
	// 有效的GUID
	validGuids := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6ba7b811-9dad-11d1-80b4-00c04fd430c8",
		"00000000-0000-4000-8000-000000000000",
	}
	
	for _, guid := range validGuids {
		if !IsValidGUID(guid) {
			t.Errorf("有效GUID被判断为无效: %s", guid)
		}
	}
	
	// 无效的GUID
	invalidGuids := []string{
		"",                                      // 空字符串
		"550e8400-e29b-41d4-a716",              // 长度不够
		"550e8400-e29b-41d4-a716-446655440000-extra", // 长度过长
		"550e8400_e29b_41d4_a716_446655440000", // 错误的分隔符
		"550e8400-e29b-41d4-a716-44665544000g", // 包含非十六进制字符
		"550e8400-e29b-41d4-a716-4466554400",   // 最后一段长度不够
	}
	
	for _, guid := range invalidGuids {
		if IsValidGUID(guid) {
			t.Errorf("无效GUID被判断为有效: %s", guid)
		}
	}
}

/**
 * TestIsValidShortGUID 测试短GUID格式验证
 */
func TestIsValidShortGUID(t *testing.T) {
	// 有效的短GUID
	validShortGuids := []string{
		"550e8400e29b41d4a716446655440000",
		"6ba7b8109dad11d180b400c04fd430c8",
		"00000000000040008000000000000000",
	}
	
	for _, guid := range validShortGuids {
		if !IsValidShortGUID(guid) {
			t.Errorf("有效短GUID被判断为无效: %s", guid)
		}
	}
	
	// 无效的短GUID
	invalidShortGuids := []string{
		"",                                // 空字符串
		"550e8400e29b41d4a716",           // 长度不够
		"550e8400e29b41d4a716446655440000extra", // 长度过长
		"550e8400e29b41d4a716446655440g0", // 包含非十六进制字符
	}
	
	for _, guid := range invalidShortGuids {
		if IsValidShortGUID(guid) {
			t.Errorf("无效短GUID被判断为有效: %s", guid)
		}
	}
}

/**
 * TestConvertGUIDToShort 测试GUID转短GUID
 */
func TestConvertGUIDToShort(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		hasError bool
	}{
		{
			input:    "550e8400-e29b-41d4-a716-446655440000",
			expected: "550e8400e29b41d4a716446655440000",
			hasError: false,
		},
		{
			input:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expected: "6ba7b8109dad11d180b400c04fd430c8",
			hasError: false,
		},
		{
			input:    "invalid-guid",
			expected: "",
			hasError: true,
		},
	}
	
	for _, tc := range testCases {
		result, err := ConvertGUIDToShort(tc.input)
		
		if tc.hasError {
			if err == nil {
				t.Errorf("期望出错但没有出错，输入: %s", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("不期望出错但出错了，输入: %s, 错误: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("转换结果不正确，输入: %s, 期望: %s, 实际: %s", tc.input, tc.expected, result)
			}
		}
	}
}

/**
 * TestConvertShortGUIDToStandard 测试短GUID转标准GUID
 */
func TestConvertShortGUIDToStandard(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		hasError bool
	}{
		{
			input:    "550e8400e29b41d4a716446655440000",
			expected: "550e8400-e29b-41d4-a716-446655440000",
			hasError: false,
		},
		{
			input:    "6ba7b8109dad11d180b400c04fd430c8",
			expected: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			hasError: false,
		},
		{
			input:    "invalid",
			expected: "",
			hasError: true,
		},
	}
	
	for _, tc := range testCases {
		result, err := ConvertShortGUIDToStandard(tc.input)
		
		if tc.hasError {
			if err == nil {
				t.Errorf("期望出错但没有出错，输入: %s", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("不期望出错但出错了，输入: %s, 错误: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Errorf("转换结果不正确，输入: %s, 期望: %s, 实际: %s", tc.input, tc.expected, result)
			}
		}
	}
}

/**
 * TestGenerateGUIDWithPrefix 测试带前缀的GUID生成
 */
func TestGenerateGUIDWithPrefix(t *testing.T) {
	prefix := "test"
	guid := GenerateGUIDWithPrefix(prefix)
	
	if !strings.HasPrefix(guid, prefix+"_") {
		t.Errorf("生成的GUID没有正确的前缀，期望前缀: %s_, 实际: %s", prefix, guid)
	}
	
	// 提取GUID部分并验证
	guidPart := strings.TrimPrefix(guid, prefix+"_")
	if !IsValidGUID(guidPart) {
		t.Errorf("提取的GUID部分格式无效: %s", guidPart)
	}
	
	// 测试空前缀
	emptyPrefixGuid := GenerateGUIDWithPrefix("")
	if !IsValidGUID(emptyPrefixGuid) {
		t.Errorf("空前缀生成的GUID格式无效: %s", emptyPrefixGuid)
	}
}

/**
 * TestCompatibilityFunctions 测试兼容性函数
 */
func TestCompatibilityFunctions(t *testing.T) {
	// 测试GenerateID兼容性函数
	id := GenerateID()
	if !IsValidGUID(id) {
		t.Errorf("GenerateID生成的ID格式无效: %s", id)
	}
	
	// 测试GenerateIDString兼容性函数
	idString := GenerateIDString()
	if !IsValidGUID(idString) {
		t.Errorf("GenerateIDString生成的ID格式无效: %s", idString)
	}
}
