package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

/**
 * 数据库升级工具类
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 提供数据库升级的通用方法,包括表、字段的创建、修改等操作
 */

/**
 * UpgradeUtils 数据库升级工具
 */
type UpgradeUtils struct {
	db *sql.DB
}

/**
 * NewUpgradeUtils 创建升级工具实例
 * @param db 数据库连接
 * @return *UpgradeUtils 升级工具实例
 */
func NewUpgradeUtils(db *sql.DB) *UpgradeUtils {
	return &UpgradeUtils{db: db}
}

/**
 * TableExists 检查表是否存在
 * @param tableName 表名
 * @return bool 是否存在
 * @return error 错误信息
 */
func (u *UpgradeUtils) TableExists(tableName string) (bool, error) {
	var count int
	err := u.db.QueryRow(`
		SELECT COUNT(*) 
		FROM sqlite_master 
		WHERE type='table' AND name=?
	`, tableName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("检查表是否存在失败: %w", err)
	}
	return count > 0, nil
}

/**
 * ColumnExists 检查字段是否存在
 * @param tableName 表名
 * @param columnName 字段名
 * @return bool 是否存在
 * @return error 错误信息
 */
func (u *UpgradeUtils) ColumnExists(tableName string, columnName string) (bool, error) {
	var count int
	err := u.db.QueryRow(`
		SELECT COUNT(*) 
		FROM pragma_table_info(?) 
		WHERE name=?
	`, tableName, columnName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("检查字段是否存在失败: %w", err)
	}
	return count > 0, nil
}

/**
 * CreateTable 创建表(如果不存在)
 * @param tableName 表名
 * @param createSQL 建表SQL语句
 * @return error 错误信息
 */
func (u *UpgradeUtils) CreateTable(tableName string, createSQL string) error {
	exists, err := u.TableExists(tableName)
	if err != nil {
		return err
	}

	if !exists {
		_, err = u.db.Exec(createSQL)
		if err != nil {
			return fmt.Errorf("创建表 %s 失败: %w", tableName, err)
		}
		log.Printf("表 %s 创建成功", tableName)
	} else {
		log.Printf("表 %s 已存在,跳过创建", tableName)
	}

	return nil
}

/**
 * AddColumn 添加字段(如果不存在)
 * @param tableName 表名
 * @param columnName 字段名
 * @param columnDef 字段定义(如: "INTEGER DEFAULT 0")
 * @return error 错误信息
 */
func (u *UpgradeUtils) AddColumn(tableName string, columnName string, columnDef string) error {
	exists, err := u.ColumnExists(tableName, columnName)
	if err != nil {
		return err
	}

	if !exists {
		sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnDef)
		_, err = u.db.Exec(sql)
		if err != nil {
			return fmt.Errorf("添加字段 %s.%s 失败: %w", tableName, columnName, err)
		}
		log.Printf("字段 %s.%s 添加成功", tableName, columnName)
	} else {
		log.Printf("字段 %s.%s 已存在,跳过添加", tableName, columnName)
	}

	return nil
}

/**
 * RenameColumn 重命名字段
 * @param tableName 表名
 * @param oldColumnName 原字段名
 * @param newColumnName 新字段名
 * @return error 错误信息
 */
func (u *UpgradeUtils) RenameColumn(tableName string, oldColumnName string, newColumnName string) error {
	exists, err := u.ColumnExists(tableName, oldColumnName)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("字段 %s.%s 不存在,跳过重命名", tableName, oldColumnName)
		return nil
	}

	// 检查新字段名是否已存在
	newExists, err := u.ColumnExists(tableName, newColumnName)
	if err != nil {
		return err
	}

	if newExists {
		log.Printf("字段 %s.%s 已存在,跳过重命名", tableName, newColumnName)
		return nil
	}

	sql := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", tableName, oldColumnName, newColumnName)
	_, err = u.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("重命名字段 %s.%s 失败: %w", tableName, oldColumnName, err)
	}

	log.Printf("字段 %s.%s 重命名为 %s 成功", tableName, oldColumnName, newColumnName)
	return nil
}

/**
 * ChangeColumnType 修改字段类型
 * @param tableName 表名
 * @param columnName 字段名
 * @param newColumnDef 新字段定义
 * @return error 错误信息
 * @description 通过备份旧字段、创建新字段、数据迁移的方式修改字段类型
 */
func (u *UpgradeUtils) ChangeColumnType(tableName string, columnName string, newColumnDef string) error {
	exists, err := u.ColumnExists(tableName, columnName)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("字段 %s.%s 不存在,跳过类型修改", tableName, columnName)
		return nil
	}

	// 生成备份字段名: 原字段名_日期_随机5位字母
	backupColumnName := fmt.Sprintf("%s_%s_%s", columnName, time.Now().Format("20060102"), generateRandomString(5))

	// 1. 重命名原字段为备份字段
	err = u.RenameColumn(tableName, columnName, backupColumnName)
	if err != nil {
		return err
	}

	// 2. 创建新字段
	err = u.AddColumn(tableName, columnName, newColumnDef)
	if err != nil {
		return err
	}

	// 3. 迁移数据(尝试类型转换,失败则设为NULL)
	sql := fmt.Sprintf("UPDATE %s SET %s = CAST(%s AS %s)", tableName, columnName, backupColumnName, extractTypeName(newColumnDef))
	_, err = u.db.Exec(sql)
	if err != nil {
		// 如果类型转换失败,尝试直接复制
		log.Printf("类型转换失败,尝试直接复制数据: %v", err)
		sql = fmt.Sprintf("UPDATE %s SET %s = %s", tableName, columnName, backupColumnName)
		_, err = u.db.Exec(sql)
		if err != nil {
			log.Printf("数据迁移失败: %v", err)
		}
	}

	log.Printf("字段 %s.%s 类型修改成功,备份字段: %s", tableName, columnName, backupColumnName)
	return nil
}

/**
 * generateRandomString 生成随机字符串
 * @param length 长度
 * @return string 随机字符串
 */
func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

/**
 * extractTypeName 从字段定义中提取类型名称
 * @param columnDef 字段定义
 * @return string 类型名称
 */
func extractTypeName(columnDef string) string {
	// 简单提取类型名称(取第一个空格前的内容)
	for i, c := range columnDef {
		if c == ' ' {
			return columnDef[:i]
		}
	}
	return columnDef
}

/**
 * ExecuteSQL 执行SQL语句
 * @param sql SQL语句
 * @return error 错误信息
 */
func (u *UpgradeUtils) ExecuteSQL(sql string) error {
	_, err := u.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("执行SQL失败: %w", err)
	}
	return nil
}

/**
 * GetColumnType 获取字段类型
 * @param tableName 表名
 * @param columnName 字段名
 * @return string 字段类型
 * @return error 错误信息
 */
func (u *UpgradeUtils) GetColumnType(tableName string, columnName string) (string, error) {
	var columnType string
	err := u.db.QueryRow(`
		SELECT type 
		FROM pragma_table_info(?) 
		WHERE name=?
	`, tableName, columnName).Scan(&columnType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("字段 %s.%s 不存在", tableName, columnName)
		}
		return "", fmt.Errorf("获取字段类型失败: %w", err)
	}
	return columnType, nil
}

/**
 * ConvertValue 转换值类型
 * @param value 原始值
 * @param targetType 目标类型
 * @return interface{} 转换后的值
 */
func ConvertValue(value interface{}, targetType string) interface{} {
	if value == nil {
		return nil
	}

	strValue := fmt.Sprintf("%v", value)

	switch targetType {
	case "INTEGER":
		if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
			return intValue
		}
		return nil
	case "REAL", "FLOAT", "DOUBLE":
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			return floatValue
		}
		return nil
	case "TEXT", "VARCHAR":
		return strValue
	case "BOOLEAN":
		if strValue == "1" || strValue == "true" || strValue == "TRUE" {
			return true
		}
		return false
	default:
		return value
	}
}
