package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"wepassword/internal/logger"
	"wepassword/internal/models"
)

/**
 * 密码生成器
 * @author 陈凤庆
 * @description 实现各种密码生成算法
 */

// 字符集定义
const (
	LowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers          = "0123456789"
	SpecialChars     = "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// 自定义规则字符集
	LowerAlphanumeric = "abcdefghijklmnopqrstuvwxyz0123456789"                           // a = 小写字母数字
	MixedAlphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" // A = 混合字母数字
	UpperAlphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"                           // U = 大写字母数字
	Digits            = "0123456789"
	LowerHex          = "0123456789abcdef"
	UpperHex          = "0123456789ABCDEF"
	LowerLettersOnly  = "abcdefghijklmnopqrstuvwxyz"
	MixedLetters      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	UpperLettersOnly  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerVowels       = "aeiou"
	MixedVowels       = "AEIOUaeiou"
	UpperVowels       = "AEIOU"
	LowerConsonants   = "bcdfghjklmnpqrstvwxyz"
	MixedConsonants   = "BCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz"
	UpperConsonants   = "BCDFGHJKLMNPQRSTVWXYZ"
	Punctuation       = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Brackets          = "()[]{}<>"
	SpecialCharsSet   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	PrintableASCII    = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	Latin1Supplement  = "¡¢£¤¥¦§¨©ª«¬­®¯°±²³´µ¶·¸¹º»¼½¾¿ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖ×ØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõö÷øùúûüýþÿ"
)

/**
 * generateGeneralPassword 生成通用规则密码
 * @param configJSON 配置JSON字符串
 * @return string 生成的密码
 * @return error 错误信息
 */
func (prs *PasswordRuleService) generateGeneralPassword(configJSON string) (string, error) {
	var config models.GeneralRuleConfig
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return "", fmt.Errorf("解析通用规则配置失败: %w", err)
	}

	// 验证配置
	if err := prs.validateGeneralConfig(config); err != nil {
		return "", err
	}

	// 构建字符集和最小要求
	var charSets []string
	var minCounts []int
	var labels []string

	if config.IncludeUppercase {
		charSets = append(charSets, UppercaseLetters)
		minCounts = append(minCounts, config.MinUppercase)
		labels = append(labels, "大写字母")
	}

	if config.IncludeLowercase {
		charSets = append(charSets, LowercaseLetters)
		minCounts = append(minCounts, config.MinLowercase)
		labels = append(labels, "小写字母")
	}

	if config.IncludeNumbers {
		charSets = append(charSets, Numbers)
		minCounts = append(minCounts, config.MinNumbers)
		labels = append(labels, "数字")
	}

	if config.IncludeSpecialChars {
		charSets = append(charSets, SpecialChars)
		minCounts = append(minCounts, config.MinSpecialChars)
		labels = append(labels, "特殊字符")
	}

	if config.IncludeCustomChars && config.CustomSpecialChars != "" {
		charSets = append(charSets, config.CustomSpecialChars)
		minCounts = append(minCounts, config.MinCustomChars)
		labels = append(labels, "自定义字符")
	}

	if len(charSets) == 0 {
		return "", fmt.Errorf("至少需要选择一种字符类型")
	}

	// 生成密码
	password, err := prs.generatePasswordWithConstraints(charSets, minCounts, config.Length)
	if err != nil {
		return "", fmt.Errorf("生成密码失败: %w", err)
	}

	// 验证生成的密码是否符合要求
	if err := prs.validateGeneratedPassword(password, config); err != nil {
		logger.Error("[密码生成器] 生成的密码不符合要求，重新生成: %v", err)
		// 重试一次
		password, err = prs.generatePasswordWithConstraints(charSets, minCounts, config.Length)
		if err != nil {
			return "", fmt.Errorf("重新生成密码失败: %w", err)
		}
	}

	return password, nil
}

/**
 * generateCustomPassword 生成自定义规则密码
 * @param configJSON 配置JSON字符串
 * @return string 生成的密码
 * @return error 错误信息
 */
func (prs *PasswordRuleService) generateCustomPassword(configJSON string) (string, error) {
	var config models.CustomRuleConfig
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return "", fmt.Errorf("解析自定义规则配置失败: %w", err)
	}

	if config.Pattern == "" {
		return "", fmt.Errorf("自定义规则模式不能为空")
	}

	password, err := prs.parseCustomPattern(config.Pattern)
	if err != nil {
		return "", fmt.Errorf("解析自定义规则失败: %w", err)
	}

	return password, nil
}

/**
 * validateGeneralConfig 验证通用规则配置
 * @param config 通用规则配置
 * @return error 错误信息
 */
func (prs *PasswordRuleService) validateGeneralConfig(config models.GeneralRuleConfig) error {
	if config.Length <= 0 {
		return fmt.Errorf("密码长度必须大于0")
	}

	if config.Length > 128 {
		return fmt.Errorf("密码长度不能超过128")
	}

	// 计算最小位数总和
	minTotal := 0
	if config.IncludeUppercase {
		minTotal += config.MinUppercase
	}
	if config.IncludeLowercase {
		minTotal += config.MinLowercase
	}
	if config.IncludeNumbers {
		minTotal += config.MinNumbers
	}
	if config.IncludeSpecialChars {
		minTotal += config.MinSpecialChars
	}
	if config.IncludeCustomChars {
		minTotal += config.MinCustomChars
	}

	if minTotal > config.Length {
		return fmt.Errorf("最小位数总和(%d)不能大于密码长度(%d)", minTotal, config.Length)
	}

	return nil
}

/**
 * generatePasswordWithConstraints 根据约束生成密码
 * @param charSets 字符集列表
 * @param minCounts 最小数量列表
 * @param totalLength 总长度
 * @return string 生成的密码
 * @return error 错误信息
 */
func (prs *PasswordRuleService) generatePasswordWithConstraints(charSets []string, minCounts []int, totalLength int) (string, error) {
	if len(charSets) != len(minCounts) {
		return "", fmt.Errorf("字符集和最小数量列表长度不匹配")
	}

	var password []rune
	allChars := ""

	// 首先满足最小要求
	for i, charset := range charSets {
		minCount := minCounts[i]
		for j := 0; j < minCount; j++ {
			char, err := prs.getRandomChar(charset)
			if err != nil {
				return "", err
			}
			password = append(password, char)
		}
		allChars += charset
	}

	// 填充剩余长度
	remaining := totalLength - len(password)
	for i := 0; i < remaining; i++ {
		char, err := prs.getRandomChar(allChars)
		if err != nil {
			return "", err
		}
		password = append(password, char)
	}

	// 打乱密码顺序
	for i := len(password) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

/**
 * getRandomChar 从字符集中获取随机字符
 * @param charset 字符集
 * @return rune 随机字符
 * @return error 错误信息
 */
func (prs *PasswordRuleService) getRandomChar(charset string) (rune, error) {
	if charset == "" {
		return 0, fmt.Errorf("字符集不能为空")
	}

	chars := []rune(charset)
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, fmt.Errorf("生成随机索引失败: %w", err)
	}

	return chars[index.Int64()], nil
}

/**
 * validateGeneratedPassword 验证生成的密码是否符合要求
 * @param password 生成的密码
 * @param config 通用规则配置
 * @return error 错误信息
 */
func (prs *PasswordRuleService) validateGeneratedPassword(password string, config models.GeneralRuleConfig) error {
	if len(password) != config.Length {
		return fmt.Errorf("密码长度不符合要求: 期望%d，实际%d", config.Length, len(password))
	}

	// 统计各类字符数量
	uppercaseCount := 0
	lowercaseCount := 0
	numberCount := 0
	specialCount := 0
	customCount := 0

	for _, char := range password {
		if strings.ContainsRune(UppercaseLetters, char) {
			uppercaseCount++
		} else if strings.ContainsRune(LowercaseLetters, char) {
			lowercaseCount++
		} else if strings.ContainsRune(Numbers, char) {
			numberCount++
		} else if strings.ContainsRune(SpecialChars, char) {
			specialCount++
		} else if config.CustomSpecialChars != "" && strings.ContainsRune(config.CustomSpecialChars, char) {
			customCount++
		}
	}

	// 验证最小要求
	if config.IncludeUppercase && uppercaseCount < config.MinUppercase {
		return fmt.Errorf("大写字母数量不足: 需要%d，实际%d", config.MinUppercase, uppercaseCount)
	}

	if config.IncludeLowercase && lowercaseCount < config.MinLowercase {
		return fmt.Errorf("小写字母数量不足: 需要%d，实际%d", config.MinLowercase, lowercaseCount)
	}

	if config.IncludeNumbers && numberCount < config.MinNumbers {
		return fmt.Errorf("数字数量不足: 需要%d，实际%d", config.MinNumbers, numberCount)
	}

	if config.IncludeSpecialChars && specialCount < config.MinSpecialChars {
		return fmt.Errorf("特殊字符数量不足: 需要%d，实际%d", config.MinSpecialChars, specialCount)
	}

	if config.IncludeCustomChars && customCount < config.MinCustomChars {
		return fmt.Errorf("自定义字符数量不足: 需要%d，实际%d", config.MinCustomChars, customCount)
	}

	return nil
}

/**
 * parseCustomPattern 解析自定义规则模式
 * @param pattern 自定义规则模式
 * @return string 生成的密码
 * @return error 错误信息
 */
func (prs *PasswordRuleService) parseCustomPattern(pattern string) (string, error) {
	var result strings.Builder
	i := 0

	for i < len(pattern) {
		char := pattern[i]

		// 处理转义字符
		if char == '\\' && i+1 < len(pattern) {
			nextChar := pattern[i+1]
			switch nextChar {
			case 'a':
				result.WriteRune('a')
			case 'b':
				result.WriteRune('b')
			case 'c':
				result.WriteRune('c')
			default:
				result.WriteByte(nextChar)
			}
			i += 2
			continue
		}

		// 处理自定义字符集 [...]
		if char == '[' {
			end := strings.Index(pattern[i:], "]")
			if end == -1 {
				return "", fmt.Errorf("未找到匹配的']'")
			}

			customSet := pattern[i+1 : i+end]
			if customSet == "" {
				return "", fmt.Errorf("自定义字符集不能为空")
			}

			// 检查是否有重复次数 {n}
			nextPos := i + end + 1
			if nextPos < len(pattern) && pattern[nextPos] == '{' {
				// 查找重复次数
				countEnd := strings.Index(pattern[nextPos:], "}")
				if countEnd == -1 {
					return "", fmt.Errorf("未找到匹配的'}'")
				}

				countStr := pattern[nextPos+1 : nextPos+countEnd]
				count, err := strconv.Atoi(countStr)
				if err != nil {
					return "", fmt.Errorf("无效的重复次数: %s", countStr)
				}

				if count <= 0 {
					return "", fmt.Errorf("重复次数必须大于0")
				}

				// 生成指定次数的字符
				for j := 0; j < count; j++ {
					randomChar, err := prs.getRandomChar(customSet)
					if err != nil {
						return "", fmt.Errorf("从自定义字符集生成字符失败: %w", err)
					}
					result.WriteRune(randomChar)
				}

				i = nextPos + countEnd + 1
			} else {
				// 没有重复次数，生成一个字符
				randomChar, err := prs.getRandomChar(customSet)
				if err != nil {
					return "", fmt.Errorf("从自定义字符集生成字符失败: %w", err)
				}
				result.WriteRune(randomChar)
				i += end + 1
			}
			continue
		}

		// 处理重复次数 {n}
		if i+1 < len(pattern) && pattern[i+1] == '{' {
			// 查找重复次数
			end := strings.Index(pattern[i+1:], "}")
			if end == -1 {
				return "", fmt.Errorf("未找到匹配的'}'")
			}

			countStr := pattern[i+2 : i+1+end]
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return "", fmt.Errorf("无效的重复次数: %s", countStr)
			}

			if count <= 0 {
				return "", fmt.Errorf("重复次数必须大于0")
			}

			// 生成指定次数的字符
			charset, err := prs.getCharsetForIdentifier(char)
			if err != nil {
				return "", err
			}

			for j := 0; j < count; j++ {
				randomChar, err := prs.getRandomChar(charset)
				if err != nil {
					return "", fmt.Errorf("生成字符失败: %w", err)
				}
				result.WriteRune(randomChar)
			}

			i += end + 2
			continue
		}

		// 处理单个字符标识符
		charset, err := prs.getCharsetForIdentifier(char)
		if err != nil {
			// 如果不是标识符，直接添加字符
			result.WriteByte(char)
		} else {
			randomChar, err := prs.getRandomChar(charset)
			if err != nil {
				return "", fmt.Errorf("生成字符失败: %w", err)
			}
			result.WriteRune(randomChar)
		}

		i++
	}

	return result.String(), nil
}

/**
 * getCharsetForIdentifier 根据标识符获取字符集
 * @param identifier 字符标识符
 * @return string 对应的字符集
 * @return error 错误信息
 */
func (prs *PasswordRuleService) getCharsetForIdentifier(identifier byte) (string, error) {
	switch identifier {
	case 'a':
		return LowerLettersOnly, nil // 小写字母（不含数字）
	case 'A':
		return UpperLettersOnly, nil // 大写字母（不含数字）
	case 'U':
		return UpperAlphanumeric, nil // 大写字母数字
	case 'd':
		return Digits, nil
	case 'h':
		return LowerHex, nil
	case 'H':
		return UpperHex, nil
	case 'l':
		return LowerLettersOnly, nil
	case 'L':
		return MixedLetters, nil
	case 'u':
		return UpperLettersOnly, nil
	case 'v':
		return LowerVowels, nil
	case 'V':
		return MixedVowels, nil
	case 'Z':
		return UpperVowels, nil
	case 'c':
		return LowerConsonants, nil
	case 'C':
		return MixedConsonants, nil
	case 'z':
		return UpperConsonants, nil
	case 'p':
		return Punctuation, nil
	case 'b':
		return Brackets, nil
	case 's':
		return SpecialCharsSet, nil
	case 'S':
		return PrintableASCII, nil
	case 'x':
		return Latin1Supplement, nil
	default:
		return "", fmt.Errorf("未知的字符标识符: %c", identifier)
	}
}
