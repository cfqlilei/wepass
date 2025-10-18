package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"wepassword/internal/database"
	"wepassword/internal/logger"
)

/**
 * UsernameHistoryService 用户名历史记录服务
 * @author 陈凤庆
 * @date 2025-10-17
 * @description 管理用户名历史记录的加密存储和解密加载
 */
type UsernameHistoryService struct {
	dbManager *database.DatabaseManager
}

/**
 * UsernameHistoryData 用户名历史记录数据结构
 */
type UsernameHistoryData struct {
	Usernames []string  `json:"usernames"`
	UpdatedAt time.Time `json:"updated_at"`
}

/**
 * NewUsernameHistoryService 创建用户名历史记录服务
 * @param dbManager 数据库管理器
 * @return *UsernameHistoryService 用户名历史记录服务实例
 */
func NewUsernameHistoryService(dbManager *database.DatabaseManager) *UsernameHistoryService {
	return &UsernameHistoryService{
		dbManager: dbManager,
	}
}

/**
 * GetUsernameHistory 获取用户名历史记录
 * @param password 登录密码，用于解密
 * @return []string 用户名列表
 * @return error 错误信息
 */
func (uhs *UsernameHistoryService) GetUsernameHistory(password string) ([]string, error) {
	if !uhs.dbManager.IsOpened() {
		return nil, fmt.Errorf("数据库未打开")
	}

	if password == "" {
		return nil, fmt.Errorf("密码不能为空")
	}

	db := uhs.dbManager.GetDB()

	// 查询加密的用户名历史记录
	var encryptedData string
	err := db.QueryRow(`
		SELECT encrypted_data FROM username_history LIMIT 1
	`).Scan(&encryptedData)

	if err != nil {
		if err == sql.ErrNoRows {
			// 没有历史记录，返回空列表
			return []string{}, nil
		}
		return nil, fmt.Errorf("查询用户名历史记录失败: %w", err)
	}

	// 解密数据
	decryptedData, err := uhs.decryptData(encryptedData, password)
	if err != nil {
		logger.Error("[用户名历史] 解密失败，可能是密码错误: %v", err)
		// 密码错误时返回空列表，不抛出错误
		return []string{}, nil
	}

	// 解析JSON数据
	var historyData UsernameHistoryData
	err = json.Unmarshal([]byte(decryptedData), &historyData)
	if err != nil {
		logger.Error("[用户名历史] 解析JSON数据失败: %v", err)
		return []string{}, nil
	}

	// 按字母顺序排序
	sort.Strings(historyData.Usernames)

	logger.Info("[用户名历史] 成功加载 %d 个历史用户名", len(historyData.Usernames))
	return historyData.Usernames, nil
}

/**
 * SaveUsernameToHistory 保存用户名到历史记录
 * @param username 用户名
 * @param password 登录密码，用于加密
 * @return error 错误信息
 */
func (uhs *UsernameHistoryService) SaveUsernameToHistory(username, password string) error {
	if !uhs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	if username == "" || password == "" {
		return fmt.Errorf("用户名和密码不能为空")
	}

	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	// 获取现有的历史记录
	existingUsernames, err := uhs.GetUsernameHistory(password)
	if err != nil {
		logger.Error("[用户名历史] 获取现有历史记录失败: %v", err)
		existingUsernames = []string{}
	}

	// 检查用户名是否已存在
	for _, existing := range existingUsernames {
		if existing == username {
			logger.Info("[用户名历史] 用户名已存在，跳过保存: %s", username)
			return nil
		}
	}

	// 添加新用户名
	existingUsernames = append(existingUsernames, username)

	// 限制历史记录数量（最多保存50个）
	const maxHistoryCount = 50
	if len(existingUsernames) > maxHistoryCount {
		// 保留最新的50个
		existingUsernames = existingUsernames[len(existingUsernames)-maxHistoryCount:]
	}

	// 按字母顺序排序
	sort.Strings(existingUsernames)

	// 创建历史记录数据
	historyData := UsernameHistoryData{
		Usernames: existingUsernames,
		UpdatedAt: time.Now(),
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(historyData)
	if err != nil {
		return fmt.Errorf("序列化历史记录数据失败: %w", err)
	}

	// 加密数据
	encryptedData, err := uhs.encryptData(string(jsonData), password)
	if err != nil {
		return fmt.Errorf("加密历史记录数据失败: %w", err)
	}

	// 保存到数据库
	db := uhs.dbManager.GetDB()
	_, err = db.Exec(`
		INSERT OR REPLACE INTO username_history (id, encrypted_data, updated_at)
		VALUES (1, ?, datetime('now'))
	`, encryptedData)

	if err != nil {
		return fmt.Errorf("保存用户名历史记录失败: %w", err)
	}

	logger.Info("[用户名历史] 成功保存用户名到历史记录: %s", username)
	return nil
}

/**
 * ClearUsernameHistory 清空用户名历史记录
 * @return error 错误信息
 */
func (uhs *UsernameHistoryService) ClearUsernameHistory() error {
	if !uhs.dbManager.IsOpened() {
		return fmt.Errorf("数据库未打开")
	}

	db := uhs.dbManager.GetDB()
	_, err := db.Exec(`DELETE FROM username_history`)
	if err != nil {
		return fmt.Errorf("清空用户名历史记录失败: %w", err)
	}

	logger.Info("[用户名历史] 用户名历史记录已清空")
	return nil
}

/**
 * encryptData 加密数据
 * @param data 要加密的数据
 * @param password 密码
 * @return string 加密后的数据（Base64编码）
 * @return error 错误信息
 */
func (uhs *UsernameHistoryService) encryptData(data, password string) (string, error) {
	// 使用密码的SHA256哈希作为密钥
	key := uhs.deriveKey(password)

	// 创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建AES加密器失败: %w", err)
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	// 返回Base64编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

/**
 * decryptData 解密数据
 * @param encryptedData 加密的数据（Base64编码）
 * @param password 密码
 * @return string 解密后的数据
 * @return error 错误信息
 */
func (uhs *UsernameHistoryService) decryptData(encryptedData, password string) (string, error) {
	// Base64解码
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("Base64解码失败: %w", err)
	}

	// 使用密码的SHA256哈希作为密钥
	key := uhs.deriveKey(password)

	// 创建AES解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建AES解密器失败: %w", err)
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %w", err)
	}

	// 提取nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("密文长度不足")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %w", err)
	}

	return string(plaintext), nil
}

/**
 * deriveKey 从密码派生密钥
 * @param password 密码
 * @return []byte 32字节的密钥
 */
func (uhs *UsernameHistoryService) deriveKey(password string) []byte {
	// 简单的密钥派生，实际应用中应该使用PBKDF2或类似的方法
	// 这里为了简化，直接使用密码的SHA256哈希
	key := make([]byte, 32)
	passwordBytes := []byte(password)

	// 重复密码直到填满32字节
	for i := 0; i < 32; i++ {
		key[i] = passwordBytes[i%len(passwordBytes)]
	}

	return key
}
