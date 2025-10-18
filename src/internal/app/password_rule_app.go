package app

import (
	"context"
	"encoding/json"
	"fmt"

	"wepassword/internal/logger"
	"wepassword/internal/models"
	"wepassword/internal/services"
)

/**
 * 密码规则应用服务
 * @author 陈凤庆
 * @description 为前端提供密码规则管理的API接口
 */

/**
 * PasswordRuleApp 密码规则应用服务
 */
type PasswordRuleApp struct {
	passwordRuleService *services.PasswordRuleService
}

/**
 * NewPasswordRuleApp 创建密码规则应用服务实例
 * @param passwordRuleService 密码规则服务
 * @return *PasswordRuleApp 密码规则应用服务实例
 */
func NewPasswordRuleApp(passwordRuleService *services.PasswordRuleService) *PasswordRuleApp {
	return &PasswordRuleApp{
		passwordRuleService: passwordRuleService,
	}
}

/**
 * GetAllRules 获取所有密码规则
 * @param ctx 上下文
 * @return []models.PasswordRule 密码规则列表
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) GetAllRules(ctx context.Context) ([]models.PasswordRule, error) {
	logger.Info("[密码规则应用] 获取所有密码规则")

	rules, err := pra.passwordRuleService.GetAllRules()
	if err != nil {
		logger.Error("[密码规则应用] 获取密码规则失败: %v", err)
		return nil, fmt.Errorf("获取密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功获取 %d 个密码规则", len(rules))
	return rules, nil
}

/**
 * GetRuleByID 根据ID获取密码规则
 * @param ctx 上下文
 * @param id 规则ID
 * @return models.PasswordRule 密码规则
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) GetRuleByID(ctx context.Context, id string) (models.PasswordRule, error) {
	logger.Info("[密码规则应用] 获取密码规则: %s", id)

	if id == "" {
		return models.PasswordRule{}, fmt.Errorf("规则ID不能为空")
	}

	rule, err := pra.passwordRuleService.GetRuleByID(id)
	if err != nil {
		logger.Error("[密码规则应用] 获取密码规则失败: %v", err)
		return models.PasswordRule{}, fmt.Errorf("获取密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功获取密码规则: %s", rule.Name)
	return rule, nil
}

/**
 * GetDefaultRule 获取默认密码规则
 * @param ctx 上下文
 * @return models.PasswordRule 默认密码规则
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) GetDefaultRule(ctx context.Context) (models.PasswordRule, error) {
	logger.Info("[密码规则应用] 获取默认密码规则")

	rule, err := pra.passwordRuleService.GetDefaultRule()
	if err != nil {
		logger.Error("[密码规则应用] 获取默认密码规则失败: %v", err)
		return models.PasswordRule{}, fmt.Errorf("获取默认密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功获取默认密码规则: %s", rule.Name)
	return rule, nil
}

/**
 * CreateGeneralRule 创建通用密码规则
 * @param ctx 上下文
 * @param name 规则名称
 * @param description 规则描述
 * @param config 通用规则配置
 * @return models.PasswordRule 创建的密码规则
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) CreateGeneralRule(ctx context.Context, name, description string, config models.GeneralRuleConfig) (models.PasswordRule, error) {
	logger.Info("[密码规则应用] 创建通用密码规则: %s", name)

	rule, err := pra.passwordRuleService.CreateRule(name, description, "general", config)
	if err != nil {
		logger.Error("[密码规则应用] 创建通用密码规则失败: %v", err)
		return models.PasswordRule{}, fmt.Errorf("创建通用密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功创建通用密码规则: %s", rule.Name)
	return rule, nil
}

/**
 * CreateCustomRule 创建自定义密码规则
 * @param ctx 上下文
 * @param name 规则名称
 * @param description 规则描述
 * @param config 自定义规则配置
 * @return models.PasswordRule 创建的密码规则
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) CreateCustomRule(ctx context.Context, name, description string, config models.CustomRuleConfig) (models.PasswordRule, error) {
	logger.Info("[密码规则应用] 创建自定义密码规则: %s", name)

	rule, err := pra.passwordRuleService.CreateRule(name, description, "custom", config)
	if err != nil {
		logger.Error("[密码规则应用] 创建自定义密码规则失败: %v", err)
		return models.PasswordRule{}, fmt.Errorf("创建自定义密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功创建自定义密码规则: %s", rule.Name)
	return rule, nil
}

/**
 * UpdateGeneralRule 更新通用密码规则
 * @param ctx 上下文
 * @param id 规则ID
 * @param name 规则名称
 * @param description 规则描述
 * @param config 通用规则配置
 * @return models.PasswordRule 更新后的密码规则
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) UpdateGeneralRule(ctx context.Context, id, name, description string, config models.GeneralRuleConfig) (models.PasswordRule, error) {
	logger.Info("[密码规则应用] 更新通用密码规则: %s", id)

	rule, err := pra.passwordRuleService.UpdateRule(id, name, description, config)
	if err != nil {
		logger.Error("[密码规则应用] 更新通用密码规则失败: %v", err)
		return models.PasswordRule{}, fmt.Errorf("更新通用密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功更新通用密码规则: %s", rule.Name)
	return rule, nil
}

/**
 * UpdateCustomRule 更新自定义密码规则
 * @param ctx 上下文
 * @param id 规则ID
 * @param name 规则名称
 * @param description 规则描述
 * @param config 自定义规则配置
 * @return models.PasswordRule 更新后的密码规则
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) UpdateCustomRule(ctx context.Context, id, name, description string, config models.CustomRuleConfig) (models.PasswordRule, error) {
	logger.Info("[密码规则应用] 更新自定义密码规则: %s", id)

	rule, err := pra.passwordRuleService.UpdateRule(id, name, description, config)
	if err != nil {
		logger.Error("[密码规则应用] 更新自定义密码规则失败: %v", err)
		return models.PasswordRule{}, fmt.Errorf("更新自定义密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功更新自定义密码规则: %s", rule.Name)
	return rule, nil
}

/**
 * DeleteRule 删除密码规则
 * @param ctx 上下文
 * @param id 规则ID
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) DeleteRule(ctx context.Context, id string) error {
	logger.Info("[密码规则应用] 删除密码规则: %s", id)

	if id == "" {
		return fmt.Errorf("规则ID不能为空")
	}

	err := pra.passwordRuleService.DeleteRule(id)
	if err != nil {
		logger.Error("[密码规则应用] 删除密码规则失败: %v", err)
		return fmt.Errorf("删除密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功删除密码规则: %s", id)
	return nil
}

/**
 * GeneratePassword 根据规则生成密码
 * @param ctx 上下文
 * @param ruleID 规则ID
 * @return string 生成的密码
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) GeneratePassword(ctx context.Context, ruleID string) (string, error) {
	logger.Info("[密码规则应用] 根据规则生成密码: %s", ruleID)

	if ruleID == "" {
		return "", fmt.Errorf("规则ID不能为空")
	}

	password, err := pra.passwordRuleService.GeneratePassword(ruleID)
	if err != nil {
		logger.Error("[密码规则应用] 生成密码失败: %v", err)
		return "", fmt.Errorf("生成密码失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功生成密码，长度: %d", len(password))
	return password, nil
}

/**
 * GeneratePasswordByGeneralConfig 根据通用配置生成密码
 * @param ctx 上下文
 * @param config 通用规则配置
 * @return string 生成的密码
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) GeneratePasswordByGeneralConfig(ctx context.Context, config models.GeneralRuleConfig) (string, error) {
	logger.Info("[密码规则应用] 根据通用配置生成密码")

	password, err := pra.passwordRuleService.GeneratePasswordByConfig("general", config)
	if err != nil {
		logger.Error("[密码规则应用] 根据通用配置生成密码失败: %v", err)
		return "", fmt.Errorf("根据通用配置生成密码失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功根据通用配置生成密码，长度: %d", len(password))
	return password, nil
}

/**
 * GeneratePasswordByCustomConfig 根据自定义配置生成密码
 * @param ctx 上下文
 * @param config 自定义规则配置
 * @return string 生成的密码
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) GeneratePasswordByCustomConfig(ctx context.Context, config models.CustomRuleConfig) (string, error) {
	logger.Info("[密码规则应用] 根据自定义配置生成密码")

	password, err := pra.passwordRuleService.GeneratePasswordByConfig("custom", config)
	if err != nil {
		logger.Error("[密码规则应用] 根据自定义配置生成密码失败: %v", err)
		return "", fmt.Errorf("根据自定义配置生成密码失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功根据自定义配置生成密码，长度: %d", len(password))
	return password, nil
}

/**
 * ForceInitializeDefaultRules 强制初始化默认密码规则
 * @param ctx 上下文
 * @param force 是否强制重新创建
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) ForceInitializeDefaultRules(ctx context.Context, force bool) error {
	logger.Info("[密码规则应用] 强制初始化默认密码规则，force: %v", force)

	err := pra.passwordRuleService.ForceInitializeDefaultRules(force)
	if err != nil {
		logger.Error("[密码规则应用] 强制初始化默认密码规则失败: %v", err)
		return fmt.Errorf("强制初始化默认密码规则失败: %w", err)
	}

	logger.Info("[密码规则应用] 强制初始化默认密码规则成功")
	return nil
}

/**
 * SetRuleAsDefault 设置密码规则为默认规则
 * @param ctx 上下文
 * @param ruleID 规则ID
 * @param isDefault 是否设为默认
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) SetRuleAsDefault(ctx context.Context, ruleID string, isDefault bool) error {
	logger.Info("[密码规则应用] 设置密码规则为默认: %s, 默认: %t", ruleID, isDefault)

	err := pra.passwordRuleService.SetRuleAsDefault(ruleID, isDefault)
	if err != nil {
		logger.Error("[密码规则应用] 设置密码规则为默认失败: %v", err)
		return fmt.Errorf("设置密码规则为默认失败: %w", err)
	}

	logger.Info("[密码规则应用] 成功设置密码规则为默认: %s", ruleID)
	return nil
}

/**
 * ParseRuleConfig 解析规则配置
 * @param rule 密码规则
 * @return interface{} 解析后的配置
 * @return error 错误信息
 */
func (pra *PasswordRuleApp) ParseRuleConfig(rule models.PasswordRule) (interface{}, error) {
	switch rule.RuleType {
	case "general":
		var config models.GeneralRuleConfig
		err := json.Unmarshal([]byte(rule.Config), &config)
		if err != nil {
			return nil, fmt.Errorf("解析通用规则配置失败: %w", err)
		}
		return config, nil
	case "custom":
		var config models.CustomRuleConfig
		err := json.Unmarshal([]byte(rule.Config), &config)
		if err != nil {
			return nil, fmt.Errorf("解析自定义规则配置失败: %w", err)
		}
		return config, nil
	default:
		return nil, fmt.Errorf("不支持的规则类型: %s", rule.RuleType)
	}
}
