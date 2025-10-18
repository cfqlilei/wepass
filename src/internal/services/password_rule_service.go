package services

import (
	"encoding/json"
	"fmt"
	"time"

	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/models"
	"wepassword/internal/utils"
)

/**
 * 密码规则服务
 * @author 陈凤庆
 * @description 管理密码生成规则的业务逻辑
 */

/**
 * PasswordRuleService 密码规则服务
 */
type PasswordRuleService struct {
	dbManager *database.DatabaseManager
}

/**
 * NewPasswordRuleService 创建密码规则服务实例
 * @param dbManager 数据库管理器
 * @return *PasswordRuleService 密码规则服务实例
 */
func NewPasswordRuleService(dbManager *database.DatabaseManager) *PasswordRuleService {
	return &PasswordRuleService{
		dbManager: dbManager,
	}
}

/**
 * CreateRule 创建密码规则
 * @param name 规则名称
 * @param description 规则描述
 * @param ruleType 规则类型
 * @param config 规则配置
 * @return models.PasswordRule 创建的密码规则
 * @return error 错误信息
 */
func (prs *PasswordRuleService) CreateRule(name, description, ruleType string, config interface{}) (models.PasswordRule, error) {
	if name == "" {
		return models.PasswordRule{}, fmt.Errorf("规则名称不能为空")
	}

	if ruleType != "general" && ruleType != "custom" {
		return models.PasswordRule{}, fmt.Errorf("无效的规则类型: %s", ruleType)
	}

	// 检查规则名称是否已存在
	exists, err := prs.checkRuleNameExists(name)
	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("检查规则名称失败: %w", err)
	}
	if exists {
		return models.PasswordRule{}, fmt.Errorf("规则名称已存在: %s", name)
	}

	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("序列化规则配置失败: %w", err)
	}

	// 创建规则
	rule := models.PasswordRule{
		ID:          utils.GenerateGUID(),
		Name:        name,
		Description: description,
		RuleType:    ruleType,
		Config:      string(configJSON),
		IsDefault:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	err = prs.saveRule(rule)
	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("保存密码规则失败: %w", err)
	}

	logger.Info("[密码规则服务] 创建密码规则成功: %s", name)
	return rule, nil
}

/**
 * UpdateRule 更新密码规则
 * @param id 规则ID
 * @param name 规则名称
 * @param description 规则描述
 * @param config 规则配置
 * @return models.PasswordRule 更新后的密码规则
 * @return error 错误信息
 */
func (prs *PasswordRuleService) UpdateRule(id, name, description string, config interface{}) (models.PasswordRule, error) {
	if id == "" {
		return models.PasswordRule{}, fmt.Errorf("规则ID不能为空")
	}

	if name == "" {
		return models.PasswordRule{}, fmt.Errorf("规则名称不能为空")
	}

	// 获取现有规则
	rule, err := prs.GetRuleByID(id)
	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("获取密码规则失败: %w", err)
	}

	// 检查规则名称是否已被其他规则使用
	if rule.Name != name {
		exists, err := prs.checkRuleNameExists(name)
		if err != nil {
			return models.PasswordRule{}, fmt.Errorf("检查规则名称失败: %w", err)
		}
		if exists {
			return models.PasswordRule{}, fmt.Errorf("规则名称已存在: %s", name)
		}
	}

	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("序列化规则配置失败: %w", err)
	}

	// 更新规则
	rule.Name = name
	rule.Description = description
	rule.Config = string(configJSON)
	rule.UpdatedAt = time.Now()

	// 保存到数据库
	err = prs.saveRule(rule)
	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("更新密码规则失败: %w", err)
	}

	logger.Info("[密码规则服务] 更新密码规则成功: %s", name)
	return rule, nil
}

/**
 * DeleteRule 删除密码规则
 * @param id 规则ID
 * @return error 错误信息
 */
func (prs *PasswordRuleService) DeleteRule(id string) error {
	if id == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	// 获取规则信息用于日志记录
	rule, err := prs.GetRuleByID(id)
	if err != nil {
		return fmt.Errorf("获取密码规则失败: %w", err)
	}

	db := prs.dbManager.GetDB()
	_, err = db.Exec("DELETE FROM password_rules WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("删除密码规则失败: %w", err)
	}

	logger.Info("[密码规则服务] 删除密码规则成功: %s", rule.Name)
	return nil
}

/**
 * GetRuleByID 根据ID获取密码规则
 * @param id 规则ID
 * @return models.PasswordRule 密码规则
 * @return error 错误信息
 */
func (prs *PasswordRuleService) GetRuleByID(id string) (models.PasswordRule, error) {
	if id == "" {
		return models.PasswordRule{}, fmt.Errorf("规则ID不能为空")
	}

	db := prs.dbManager.GetDB()
	var rule models.PasswordRule

	err := db.QueryRow(`
		SELECT id, name, description, rule_type, config, is_default, created_at, updated_at
		FROM password_rules WHERE id = ?
	`, id).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.RuleType,
		&rule.Config, &rule.IsDefault, &rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("获取密码规则失败: %w", err)
	}

	return rule, nil
}

/**
 * GetAllRules 获取所有密码规则
 * @return []models.PasswordRule 密码规则列表
 * @return error 错误信息
 */
func (prs *PasswordRuleService) GetAllRules() ([]models.PasswordRule, error) {
	db := prs.dbManager.GetDB()
	rows, err := db.Query(`
		SELECT id, name, description, rule_type, config, is_default, created_at, updated_at
		FROM password_rules ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("查询密码规则失败: %w", err)
	}
	defer rows.Close()

	var rules []models.PasswordRule
	for rows.Next() {
		var rule models.PasswordRule
		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.Description, &rule.RuleType,
			&rule.Config, &rule.IsDefault, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描密码规则失败: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

/**
 * GetDefaultRule 获取默认密码规则
 * @return models.PasswordRule 默认密码规则
 * @return error 错误信息
 */
func (prs *PasswordRuleService) GetDefaultRule() (models.PasswordRule, error) {
	db := prs.dbManager.GetDB()
	var rule models.PasswordRule

	err := db.QueryRow(`
		SELECT id, name, description, rule_type, config, is_default, created_at, updated_at
		FROM password_rules WHERE is_default = 1 LIMIT 1
	`).Scan(
		&rule.ID, &rule.Name, &rule.Description, &rule.RuleType,
		&rule.Config, &rule.IsDefault, &rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return models.PasswordRule{}, fmt.Errorf("获取默认密码规则失败: %w", err)
	}

	return rule, nil
}

/**
 * checkRuleNameExists 检查规则名称是否已存在
 * @param name 规则名称
 * @return bool 是否存在
 * @return error 错误信息
 */
func (prs *PasswordRuleService) checkRuleNameExists(name string) (bool, error) {
	db := prs.dbManager.GetDB()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM password_rules WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

/**
 * saveRule 保存密码规则到数据库
 * @param rule 密码规则
 * @return error 错误信息
 */
func (prs *PasswordRuleService) saveRule(rule models.PasswordRule) error {
	db := prs.dbManager.GetDB()

	// 检查规则是否已存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM password_rules WHERE id = ?", rule.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查规则是否存在失败: %w", err)
	}

	if count > 0 {
		// 更新现有规则
		_, err = db.Exec(`
			UPDATE password_rules 
			SET name = ?, description = ?, rule_type = ?, config = ?, updated_at = ?
			WHERE id = ?
		`, rule.Name, rule.Description, rule.RuleType, rule.Config, rule.UpdatedAt, rule.ID)
	} else {
		// 插入新规则
		_, err = db.Exec(`
			INSERT INTO password_rules (id, name, description, rule_type, config, is_default, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, rule.ID, rule.Name, rule.Description, rule.RuleType, rule.Config, rule.IsDefault, rule.CreatedAt, rule.UpdatedAt)
	}

	return err
}

/**
 * GeneratePassword 根据规则生成密码
 * @param ruleID 规则ID
 * @return string 生成的密码
 * @return error 错误信息
 */
func (prs *PasswordRuleService) GeneratePassword(ruleID string) (string, error) {
	rule, err := prs.GetRuleByID(ruleID)
	if err != nil {
		return "", fmt.Errorf("获取密码规则失败: %w", err)
	}

	switch rule.RuleType {
	case "general":
		return prs.generateGeneralPassword(rule.Config)
	case "custom":
		return prs.generateCustomPassword(rule.Config)
	default:
		return "", fmt.Errorf("不支持的规则类型: %s", rule.RuleType)
	}
}

/**
 * GeneratePasswordByConfig 根据配置直接生成密码
 * @param ruleType 规则类型
 * @param config 规则配置
 * @return string 生成的密码
 * @return error 错误信息
 */
func (prs *PasswordRuleService) GeneratePasswordByConfig(ruleType string, config interface{}) (string, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %w", err)
	}

	switch ruleType {
	case "general":
		return prs.generateGeneralPassword(string(configJSON))
	case "custom":
		return prs.generateCustomPassword(string(configJSON))
	default:
		return "", fmt.Errorf("不支持的规则类型: %s", ruleType)
	}
}

/**
 * SetRuleAsDefault 设置密码规则为默认规则
 * @param ruleID 规则ID
 * @param isDefault 是否设为默认
 * @return error 错误信息
 */
func (prs *PasswordRuleService) SetRuleAsDefault(ruleID string, isDefault bool) error {
	if ruleID == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	db := prs.dbManager.GetDB()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer tx.Rollback()

	if isDefault {
		// 如果设为默认，先将所有其他规则设为非默认
		_, err = tx.Exec("UPDATE password_rules SET is_default = 0 WHERE id != ?", ruleID)
		if err != nil {
			return fmt.Errorf("清除其他默认规则失败: %w", err)
		}
	}

	// 更新指定规则的默认状态
	_, err = tx.Exec("UPDATE password_rules SET is_default = ?, updated_at = ? WHERE id = ?",
		isDefault, time.Now(), ruleID)
	if err != nil {
		return fmt.Errorf("更新规则默认状态失败: %w", err)
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	logger.Info("[密码规则服务] 设置规则默认状态成功: %s, 默认: %t", ruleID, isDefault)
	return nil
}

/**
 * InitializeDefaultRules 初始化默认密码规则
 * @return error 错误信息
 */
func (prs *PasswordRuleService) InitializeDefaultRules() error {
	return prs.ForceInitializeDefaultRules(false)
}

/**
 * ForceInitializeDefaultRules 强制初始化默认密码规则
 * @param force 是否强制重新创建
 * @return error 错误信息
 */
func (prs *PasswordRuleService) ForceInitializeDefaultRules(force bool) error {
	db := prs.dbManager.GetDB()

	if force {
		// 删除现有的默认规则
		_, err := db.Exec("DELETE FROM password_rules WHERE is_default = 1")
		if err != nil {
			return fmt.Errorf("删除现有默认规则失败: %w", err)
		}
		logger.Info("[密码规则服务] 已删除现有默认规则")
	} else {
		// 检查是否已存在默认规则
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM password_rules WHERE is_default = 1").Scan(&count)
		if err != nil {
			return fmt.Errorf("检查默认规则失败: %w", err)
		}

		if count > 0 {
			logger.Info("[密码规则服务] 默认规则已存在，跳过初始化")
			return nil
		}
	}

	// 创建通用规则
	generalConfig := models.GeneralRuleConfig{
		IncludeUppercase:    true,
		IncludeLowercase:    true,
		IncludeNumbers:      true,
		IncludeSpecialChars: false,
		IncludeCustomChars:  false,
		MinUppercase:        1,
		MinLowercase:        1,
		MinNumbers:          1,
		MinSpecialChars:     0,
		MinCustomChars:      0,
		Length:              12,
		CustomSpecialChars:  "",
	}

	configJSON, err := json.Marshal(generalConfig)
	if err != nil {
		return fmt.Errorf("序列化默认规则配置失败: %w", err)
	}

	defaultRule := models.PasswordRule{
		ID:          utils.GenerateGUID(),
		Name:        "通用规则",
		Description: "系统预置的通用密码规则",
		RuleType:    "general",
		Config:      string(configJSON),
		IsDefault:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = prs.saveRule(defaultRule)
	if err != nil {
		return fmt.Errorf("保存默认规则失败: %w", err)
	}

	logger.Info("[密码规则服务] 默认规则初始化成功")
	return nil
}
