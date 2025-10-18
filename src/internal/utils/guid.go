package utils

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

/**
 * GUID生成器
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 基于UUID v4算法的GUID生成器，替代雪花ID解决JavaScript精度丢失问题
 */

/**
 * GenerateGUID 生成GUID
 * @return string GUID字符串，格式：xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
 * @description 生成符合UUID v4标准的GUID，确保全局唯一性
 */
func GenerateGUID() string {
	// 生成16字节随机数
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备用方案
		return GenerateSimpleGUID()
	}

	// 设置版本号(4)和变体位
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // 版本4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // 变体位

	// 格式化为标准GUID格式
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16])
}

/**
 * GenerateSimpleGUID 生成简化GUID（备用方案）
 * @return string 简化的GUID字符串
 * @description 当随机数生成失败时的备用方案，基于时间戳和简单随机数
 */
func GenerateSimpleGUID() string {
	// 使用当前时间戳的纳秒部分作为基础
	now := time.Now().UnixNano()

	// 生成简单的随机后缀
	suffix := now % 999999

	return fmt.Sprintf("%016x-%06x-4000-8000-000000000000", now, suffix)
}

/**
 * GenerateShortGUID 生成短GUID
 * @return string 短GUID字符串，去掉连字符
 * @description 生成不带连字符的32位GUID字符串
 */
func GenerateShortGUID() string {
	guid := GenerateGUID()
	return strings.ReplaceAll(guid, "-", "")
}

/**
 * IsValidGUID 验证GUID格式
 * @param guid GUID字符串
 * @return bool 是否为有效的GUID格式
 * @description 验证字符串是否符合标准GUID格式
 */
func IsValidGUID(guid string) bool {
	if len(guid) != 36 {
		return false
	}

	// 检查连字符位置
	if guid[8] != '-' || guid[13] != '-' || guid[18] != '-' || guid[23] != '-' {
		return false
	}

	// 检查字符是否都是十六进制
	hexChars := "0123456789abcdefABCDEF"
	for i, char := range guid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			continue // 跳过连字符
		}

		found := false
		for _, hexChar := range hexChars {
			if char == hexChar {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

/**
 * IsValidShortGUID 验证短GUID格式
 * @param guid 短GUID字符串
 * @return bool 是否为有效的短GUID格式
 * @description 验证字符串是否符合32位短GUID格式
 */
func IsValidShortGUID(guid string) bool {
	if len(guid) != 32 {
		return false
	}

	// 检查字符是否都是十六进制
	hexChars := "0123456789abcdefABCDEF"
	for _, char := range guid {
		found := false
		for _, hexChar := range hexChars {
			if char == hexChar {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

/**
 * ConvertGUIDToShort 将标准GUID转换为短GUID
 * @param guid 标准GUID字符串
 * @return string 短GUID字符串
 * @return error 错误信息
 */
func ConvertGUIDToShort(guid string) (string, error) {
	if !IsValidGUID(guid) {
		return "", fmt.Errorf("无效的GUID格式: %s", guid)
	}

	return strings.ReplaceAll(guid, "-", ""), nil
}

/**
 * ConvertShortGUIDToStandard 将短GUID转换为标准GUID
 * @param shortGuid 短GUID字符串
 * @return string 标准GUID字符串
 * @return error 错误信息
 */
func ConvertShortGUIDToStandard(shortGuid string) (string, error) {
	if !IsValidShortGUID(shortGuid) {
		return "", fmt.Errorf("无效的短GUID格式: %s", shortGuid)
	}

	// 插入连字符
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		shortGuid[0:8],
		shortGuid[8:12],
		shortGuid[12:16],
		shortGuid[16:20],
		shortGuid[20:32]), nil
}

/**
 * GenerateGUIDWithPrefix 生成带前缀的GUID
 * @param prefix 前缀字符串
 * @return string 带前缀的GUID字符串
 * @description 生成格式为 prefix_xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx 的GUID
 */
func GenerateGUIDWithPrefix(prefix string) string {
	guid := GenerateGUID()
	if prefix == "" {
		return guid
	}
	return fmt.Sprintf("%s_%s", prefix, guid)
}
