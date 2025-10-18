package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

/**
 * 系统信息管理
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 管理系统配置信息,包括数据库版本号等
 */

const (
	// KeyDatabaseVersion 数据库版本号的键名
	KeyDatabaseVersion = "databaseversion"
)

/**
 * SysInfoManager 系统信息管理器
 */
type SysInfoManager struct {
	db *sql.DB
}

/**
 * NewSysInfoManager 创建系统信息管理器
 * @param db 数据库连接
 * @return *SysInfoManager 系统信息管理器实例
 */
func NewSysInfoManager(db *sql.DB) *SysInfoManager {
	return &SysInfoManager{db: db}
}

/**
 * InitSysInfoTable 初始化系统信息表
 * @return error 错误信息
 */
func (s *SysInfoManager) InitSysInfoTable() error {
	// 创建sysinfo表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sysinfo (
		id INTEGER PRIMARY KEY,
		keyname TEXT NOT NULL UNIQUE,
		keyvalue TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建sysinfo表失败: %w", err)
	}

	log.Println("sysinfo表初始化成功")
	return nil
}

/**
 * GetValue 获取配置值
 * @param keyname 键名
 * @return string 键值
 * @return error 错误信息
 */
func (s *SysInfoManager) GetValue(keyname string) (string, error) {
	var keyvalue string
	err := s.db.QueryRow("SELECT keyvalue FROM sysinfo WHERE keyname = ?", keyname).Scan(&keyvalue)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // 键不存在,返回空字符串
		}
		return "", fmt.Errorf("获取配置值失败: %w", err)
	}
	return keyvalue, nil
}

/**
 * SetValue 设置配置值
 * @param keyname 键名
 * @param keyvalue 键值
 * @return error 错误信息
 */
func (s *SysInfoManager) SetValue(keyname string, keyvalue string) error {
	// 检查键是否存在
	existingValue, err := s.GetValue(keyname)
	if err != nil {
		return err
	}

	now := time.Now()

	if existingValue == "" {
		// 插入新记录,使用自增ID
		_, err = s.db.Exec(`
			INSERT INTO sysinfo (keyname, keyvalue, created_at, updated_at)
			VALUES (?, ?, ?, ?)
		`, keyname, keyvalue, now, now)
		if err != nil {
			return fmt.Errorf("插入配置值失败: %w", err)
		}
		log.Printf("配置 %s = %s 插入成功", keyname, keyvalue)
	} else {
		// 更新现有记录
		_, err = s.db.Exec(`
			UPDATE sysinfo
			SET keyvalue = ?, updated_at = ?
			WHERE keyname = ?
		`, keyvalue, now, keyname)
		if err != nil {
			return fmt.Errorf("更新配置值失败: %w", err)
		}
		log.Printf("配置 %s = %s 更新成功", keyname, keyvalue)
	}

	return nil
}

/**
 * GetDatabaseVersion 获取数据库版本号
 * @return int 版本号,0表示未初始化
 * @return error 错误信息
 */
func (s *SysInfoManager) GetDatabaseVersion() (int, error) {
	versionStr, err := s.GetValue(KeyDatabaseVersion)
	if err != nil {
		return 0, err
	}

	if versionStr == "" {
		return 0, nil // 未初始化
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, fmt.Errorf("版本号格式错误: %w", err)
	}

	return version, nil
}

/**
 * SetDatabaseVersion 设置数据库版本号
 * @param version 版本号
 * @return error 错误信息
 */
func (s *SysInfoManager) SetDatabaseVersion(version int) error {
	return s.SetValue(KeyDatabaseVersion, strconv.Itoa(version))
}

/**
 * InitDatabaseVersion 初始化数据库版本号
 * @return error 错误信息
 */
func (s *SysInfoManager) InitDatabaseVersion() error {
	version, err := s.GetDatabaseVersion()
	if err != nil {
		return err
	}

	if version == 0 {
		// 初始化为版本1
		err = s.SetDatabaseVersion(1)
		if err != nil {
			return err
		}
		log.Println("数据库版本号初始化为 1")
	}

	return nil
}

/**
 * KeyExists 检查键是否存在
 * @param keyname 键名
 * @return bool 是否存在
 * @return error 错误信息
 */
func (s *SysInfoManager) KeyExists(keyname string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM sysinfo WHERE keyname = ?", keyname).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("检查键是否存在失败: %w", err)
	}
	return count > 0, nil
}

/**
 * DeleteKey 删除键
 * @param keyname 键名
 * @return error 错误信息
 */
func (s *SysInfoManager) DeleteKey(keyname string) error {
	_, err := s.db.Exec("DELETE FROM sysinfo WHERE keyname = ?", keyname)
	if err != nil {
		return fmt.Errorf("删除键失败: %w", err)
	}
	log.Printf("配置 %s 删除成功", keyname)
	return nil
}

/**
 * GetAllKeys 获取所有键值对
 * @return map[string]string 所有键值对
 * @return error 错误信息
 */
func (s *SysInfoManager) GetAllKeys() (map[string]string, error) {
	rows, err := s.db.Query("SELECT keyname, keyvalue FROM sysinfo")
	if err != nil {
		return nil, fmt.Errorf("查询所有键值对失败: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var keyname, keyvalue string
		err := rows.Scan(&keyname, &keyvalue)
		if err != nil {
			return nil, fmt.Errorf("扫描键值对失败: %w", err)
		}
		result[keyname] = keyvalue
	}

	return result, nil
}
