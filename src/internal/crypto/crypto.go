package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

/**
 * 加密模块
 * @author 陈凤庆
 * @description 提供 AES-256 加密解密功能，用于密码数据保护
 */

const (
	// AES-256 密钥长度
	KeyLength = 32
	// PBKDF2 迭代次数
	PBKDF2Iterations = 100000
	// 盐值长度
	SaltLength = 16
)

/**
 * CryptoManager 加密管理器
 */
type CryptoManager struct {
	masterKey []byte // 主密钥，从登录密码派生
}

/**
 * NewCryptoManager 创建新的加密管理器
 * @return *CryptoManager 加密管理器实例
 */
func NewCryptoManager() *CryptoManager {
	return &CryptoManager{}
}

/**
 * SetMasterPassword 设置主密码
 * @param password 登录密码
 * @param salt 盐值
 */
func (cm *CryptoManager) SetMasterPassword(password string, salt []byte) {
	// 使用 PBKDF2 从密码派生密钥
	cm.masterKey = pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeyLength, sha256.New)
}

/**
 * GenerateSalt 生成随机盐值
 * @return []byte 盐值
 * @return error 错误信息
 */
func (cm *CryptoManager) GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("生成盐值失败: %w", err)
	}
	return salt, nil
}

/**
 * HashPassword 对密码进行哈希处理（用于存储登录密码）
 * @param password 原始密码
 * @param salt 盐值
 * @return string 哈希后的密码
 */
func (cm *CryptoManager) HashPassword(password string, salt []byte) string {
	hash := pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeyLength, sha256.New)
	return base64.StdEncoding.EncodeToString(hash)
}

/**
 * VerifyPassword 验证密码
 * @param password 输入的密码
 * @param hashedPassword 存储的哈希密码
 * @param salt 盐值
 * @return bool 是否匹配
 */
func (cm *CryptoManager) VerifyPassword(password string, hashedPassword string, salt []byte) bool {
	hash := cm.HashPassword(password, salt)
	return hash == hashedPassword
}

/**
 * Encrypt 加密数据
 * @param plaintext 明文数据
 * @return string 加密后的数据（Base64编码）
 * @return error 错误信息
 */
func (cm *CryptoManager) Encrypt(plaintext string) (string, error) {
	if cm.masterKey == nil {
		return "", errors.New("主密钥未设置")
	}

	if plaintext == "" {
		return "", nil
	}

	// 创建 AES 加密器
	block, err := aes.NewCipher(cm.masterKey)
	if err != nil {
		return "", fmt.Errorf("创建 AES 加密器失败: %w", err)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 模式失败: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成 nonce 失败: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// 返回 Base64 编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

/**
 * Decrypt 解密数据
 * @param ciphertext 加密数据（Base64编码）
 * @return string 解密后的明文数据
 * @return error 错误信息
 */
func (cm *CryptoManager) Decrypt(ciphertext string) (string, error) {
	if cm.masterKey == nil {
		return "", errors.New("主密钥未设置")
	}

	if ciphertext == "" {
		return "", nil
	}

	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %w", err)
	}

	// 创建 AES 解密器
	block, err := aes.NewCipher(cm.masterKey)
	if err != nil {
		return "", fmt.Errorf("创建 AES 解密器失败: %w", err)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建 GCM 模式失败: %w", err)
	}

	// 检查数据长度
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("加密数据格式错误")
	}

	// 提取 nonce 和密文
	nonce, cipherData := data[:nonceSize], data[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}

	return string(plaintext), nil
}

/**
 * MaskString 对字符串进行脱敏处理
 * @param str 原始字符串
 * @return string 脱敏后的字符串
 */
func MaskString(str string) string {
	if str == "" {
		return ""
	}
	
	length := len([]rune(str))
	if length <= 2 {
		return strings.Repeat("*", length)
	}
	
	runes := []rune(str)
	if length <= 6 {
		// 短字符串：显示首尾各1个字符
		return string(runes[0]) + strings.Repeat("*", length-2) + string(runes[length-1])
	} else {
		// 长字符串：显示首尾各2个字符
		return string(runes[:2]) + strings.Repeat("*", length-4) + string(runes[length-2:])
	}
}

/**
 * GenerateRandomPassword 生成随机密码
 * @param length 密码长度
 * @param includeUppercase 是否包含大写字母
 * @param includeLowercase 是否包含小写字母
 * @param includeNumbers 是否包含数字
 * @param includeSymbols 是否包含特殊符号
 * @return string 生成的密码
 * @return error 错误信息
 */
func GenerateRandomPassword(length int, includeUppercase, includeLowercase, includeNumbers, includeSymbols bool) (string, error) {
	if length <= 0 {
		return "", errors.New("密码长度必须大于0")
	}

	var charset string
	if includeUppercase {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if includeLowercase {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}
	if includeNumbers {
		charset += "0123456789"
	}
	if includeSymbols {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	if charset == "" {
		return "", errors.New("至少需要选择一种字符类型")
	}

	password := make([]byte, length)
	for i := range password {
		randomIndex := make([]byte, 1)
		_, err := rand.Read(randomIndex)
		if err != nil {
			return "", fmt.Errorf("生成随机密码失败: %w", err)
		}
		password[i] = charset[int(randomIndex[0])%len(charset)]
	}

	return string(password), nil
}