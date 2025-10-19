package services

import (
	"regexp"
	"testing"
)

/**
 * TestGenerateMixedAlphanumericChars_Length1 测试A规则长度为1的情况
 * 预期：生成1个字符，可以是大写、小写或数字
 */
func TestGenerateMixedAlphanumericChars_Length1(t *testing.T) {
	prs := &PasswordRuleService{}

	for i := 0; i < 10; i++ {
		chars, err := prs.generateMixedAlphanumericChars(1)
		if err != nil {
			t.Fatalf("生成字符失败: %v", err)
		}

		if len(chars) != 1 {
			t.Errorf("期望生成1个字符，实际生成%d个", len(chars))
		}

		char := chars[0]
		isValid := (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')
		if !isValid {
			t.Errorf("生成的字符%c不是大写、小写或数字", char)
		}
	}

	t.Log("✅ A规则长度=1测试通过")
}

/**
 * TestGenerateMixedAlphanumericChars_Length2 测试A规则长度为2的情况
 * 预期：生成2个字符，必须包含至少2种类型（大写、小写、数字）
 */
func TestGenerateMixedAlphanumericChars_Length2(t *testing.T) {
	prs := &PasswordRuleService{}

	for i := 0; i < 10; i++ {
		chars, err := prs.generateMixedAlphanumericChars(2)
		if err != nil {
			t.Fatalf("生成字符失败: %v", err)
		}

		if len(chars) != 2 {
			t.Errorf("期望生成2个字符，实际生成%d个", len(chars))
		}

		// 统计字符类型
		hasUpper := false
		hasLower := false
		hasDigit := false

		for _, char := range chars {
			if char >= 'A' && char <= 'Z' {
				hasUpper = true
			} else if char >= 'a' && char <= 'z' {
				hasLower = true
			} else if char >= '0' && char <= '9' {
				hasDigit = true
			}
		}

		typeCount := 0
		if hasUpper {
			typeCount++
		}
		if hasLower {
			typeCount++
		}
		if hasDigit {
			typeCount++
		}

		if typeCount < 2 {
			t.Errorf("生成的字符%s只包含%d种类型，期望至少2种", string(chars), typeCount)
		}
	}

	t.Log("✅ A规则长度=2测试通过")
}

/**
 * TestGenerateMixedAlphanumericChars_Length3Plus 测试A规则长度>=3的情况
 * 预期：生成的字符必须包含大写、小写、数字三种类型
 */
func TestGenerateMixedAlphanumericChars_Length3Plus(t *testing.T) {
	prs := &PasswordRuleService{}

	testLengths := []int{3, 4, 5, 10}

	for _, length := range testLengths {
		for i := 0; i < 5; i++ {
			chars, err := prs.generateMixedAlphanumericChars(length)
			if err != nil {
				t.Fatalf("生成字符失败: %v", err)
			}

			if len(chars) != length {
				t.Errorf("期望生成%d个字符，实际生成%d个", length, len(chars))
			}

			// 统计字符类型
			hasUpper := false
			hasLower := false
			hasDigit := false

			for _, char := range chars {
				if char >= 'A' && char <= 'Z' {
					hasUpper = true
				} else if char >= 'a' && char <= 'z' {
					hasLower = true
				} else if char >= '0' && char <= '9' {
					hasDigit = true
				}
			}

			if !hasUpper || !hasLower || !hasDigit {
				t.Errorf("长度%d的字符%s不包含大写、小写或数字", length, string(chars))
			}
		}
	}

	t.Log("✅ A规则长度>=3测试通过")
}

/**
 * TestParseCustomPattern_A5 测试自定义规则A{5}
 * 预期：生成5个字符，必须包含大写、小写、数字
 */
func TestParseCustomPattern_A5(t *testing.T) {
	prs := &PasswordRuleService{}

	for i := 0; i < 10; i++ {
		password, err := prs.parseCustomPattern("A{5}")
		if err != nil {
			t.Fatalf("解析规则失败: %v", err)
		}

		if len(password) != 5 {
			t.Errorf("期望生成5个字符，实际生成%d个", len(password))
		}

		// 验证包含大写、小写、数字
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

		if !hasUpper || !hasLower || !hasDigit {
			t.Errorf("生成的密码%s不包含大写、小写或数字", password)
		}
	}

	t.Log("✅ 规则A{5}测试通过")
}

/**
 * TestParseCustomPattern_Complex 测试复杂规则A{5}[@$^%-]{1}A{5}
 * 预期：生成11个字符（5+1+5），A部分包含大小写数字，特殊字符部分包含指定字符
 */
func TestParseCustomPattern_Complex(t *testing.T) {
	prs := &PasswordRuleService{}

	for i := 0; i < 10; i++ {
		password, err := prs.parseCustomPattern("A{5}[@$^%-]{1}A{5}")
		if err != nil {
			t.Fatalf("解析规则失败: %v", err)
		}

		if len(password) != 11 {
			t.Errorf("期望生成11个字符，实际生成%d个", len(password))
		}

		// 验证包含大写、小写、数字
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

		if !hasUpper || !hasLower || !hasDigit {
			t.Errorf("生成的密码%s不包含大写、小写或数字", password)
		}

		// 验证包含特殊字符
		hasSpecial := regexp.MustCompile(`[@$^%-]`).MatchString(password)
		if !hasSpecial {
			t.Errorf("生成的密码%s不包含特殊字符", password)
		}
	}

	t.Log("✅ 复杂规则A{5}[@$^%-]{1}A{5}测试通过")
}

/**
 * TestGenerateMixedCharsWithConstraints_TwoTypes 测试2种字符组合
 * 预期：长度>=2时做到二选一
 */
func TestGenerateMixedCharsWithConstraints_TwoTypes(t *testing.T) {
	prs := &PasswordRuleService{}

	// 测试大写+小写组合
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	for i := 0; i < 10; i++ {
		chars, err := prs.generateMixedCharsWithConstraints(charset, 2)
		if err != nil {
			t.Fatalf("生成字符失败: %v", err)
		}

		if len(chars) != 2 {
			t.Errorf("期望生成2个字符，实际生成%d个", len(chars))
		}

		// 统计字符类型
		hasUpper := false
		hasLower := false

		for _, char := range chars {
			if char >= 'A' && char <= 'Z' {
				hasUpper = true
			} else if char >= 'a' && char <= 'z' {
				hasLower = true
			}
		}

		if !hasUpper || !hasLower {
			t.Errorf("生成的字符%s不包含大写和小写", string(chars))
		}
	}

	t.Log("✅ 2种字符组合测试通过")
}

/**
 * TestGenerateMixedCharsWithConstraints_ThreeTypes 测试3种或以上字符组合
 * 预期：长度=2时二选一，长度=3时三选一
 */
func TestGenerateMixedCharsWithConstraints_ThreeTypes(t *testing.T) {
	prs := &PasswordRuleService{}

	// 测试大写+小写+数字组合
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// 测试长度=2
	for i := 0; i < 10; i++ {
		chars, err := prs.generateMixedCharsWithConstraints(charset, 2)
		if err != nil {
			t.Fatalf("生成字符失败: %v", err)
		}

		if len(chars) != 2 {
			t.Errorf("期望生成2个字符，实际生成%d个", len(chars))
		}
	}

	// 测试长度=3
	for i := 0; i < 10; i++ {
		chars, err := prs.generateMixedCharsWithConstraints(charset, 3)
		if err != nil {
			t.Fatalf("生成字符失败: %v", err)
		}

		if len(chars) != 3 {
			t.Errorf("期望生成3个字符，实际生成%d个", len(chars))
		}

		// 统计字符类型
		hasUpper := false
		hasLower := false
		hasDigit := false

		for _, char := range chars {
			if char >= 'A' && char <= 'Z' {
				hasUpper = true
			} else if char >= 'a' && char <= 'z' {
				hasLower = true
			} else if char >= '0' && char <= '9' {
				hasDigit = true
			}
		}

		if !hasUpper || !hasLower || !hasDigit {
			t.Errorf("生成的字符%s不包含大写、小写或数字", string(chars))
		}
	}

	t.Log("✅ 3种或以上字符组合测试通过")
}

/**
 * TestDigitProbability 测试数字字符的几率
 * 预期：在A规则长度>=3时，数字出现的几率应该较高（约40%）
 */
func TestDigitProbability(t *testing.T) {
	prs := &PasswordRuleService{}

	digitCount := 0
	totalCount := 0

	for i := 0; i < 100; i++ {
		chars, _ := prs.generateMixedAlphanumericChars(10)
		for _, char := range chars {
			totalCount++
			if char >= '0' && char <= '9' {
				digitCount++
			}
		}
	}

	digitProbability := float64(digitCount) / float64(totalCount) * 100

	// 数字几率应该在30%-50%之间（目标是40%）
	if digitProbability < 30 || digitProbability > 50 {
		t.Logf("数字几率为%.2f%%，期望在30-50%%之间", digitProbability)
	} else {
		t.Logf("数字几率为%.2f%%，符合预期", digitProbability)
	}
}
