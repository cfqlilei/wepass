package app

import (
	"context"
	"fmt"

	"wepassword/internal/logger"
	"wepassword/internal/services"
)

/**
 * UsernameHistoryApp 用户名历史记录应用服务
 * @author 陈凤庆
 * @date 2025-10-17
 * @description 用户名历史记录应用层服务，提供用户名历史记录管理功能
 */
type UsernameHistoryApp struct {
	usernameHistoryService *services.UsernameHistoryService
}

/**
 * NewUsernameHistoryApp 创建用户名历史记录应用服务
 * @param usernameHistoryService 用户名历史记录服务
 * @return *UsernameHistoryApp 用户名历史记录应用服务实例
 */
func NewUsernameHistoryApp(usernameHistoryService *services.UsernameHistoryService) *UsernameHistoryApp {
	return &UsernameHistoryApp{
		usernameHistoryService: usernameHistoryService,
	}
}

/**
 * GetUsernameHistory 获取用户名历史记录
 * @param ctx 上下文
 * @param password 登录密码，用于解密
 * @return []string 用户名列表
 * @return error 错误信息
 */
func (uha *UsernameHistoryApp) GetUsernameHistory(ctx context.Context, password string) ([]string, error) {
	logger.Info("[用户名历史应用] 获取用户名历史记录")

	usernames, err := uha.usernameHistoryService.GetUsernameHistory(password)
	if err != nil {
		logger.Error("[用户名历史应用] 获取用户名历史记录失败: %v", err)
		return nil, fmt.Errorf("获取用户名历史记录失败: %w", err)
	}

	logger.Info("[用户名历史应用] 成功获取 %d 个历史用户名", len(usernames))
	return usernames, nil
}

/**
 * SaveUsernameToHistory 保存用户名到历史记录
 * @param ctx 上下文
 * @param username 用户名
 * @param password 登录密码，用于加密
 * @return error 错误信息
 */
func (uha *UsernameHistoryApp) SaveUsernameToHistory(ctx context.Context, username, password string) error {
	logger.Info("[用户名历史应用] 保存用户名到历史记录: %s", username)

	err := uha.usernameHistoryService.SaveUsernameToHistory(username, password)
	if err != nil {
		logger.Error("[用户名历史应用] 保存用户名到历史记录失败: %v", err)
		return fmt.Errorf("保存用户名到历史记录失败: %w", err)
	}

	logger.Info("[用户名历史应用] 成功保存用户名到历史记录: %s", username)
	return nil
}

/**
 * ClearUsernameHistory 清空用户名历史记录
 * @param ctx 上下文
 * @return error 错误信息
 */
func (uha *UsernameHistoryApp) ClearUsernameHistory(ctx context.Context) error {
	logger.Info("[用户名历史应用] 清空用户名历史记录")

	err := uha.usernameHistoryService.ClearUsernameHistory()
	if err != nil {
		logger.Error("[用户名历史应用] 清空用户名历史记录失败: %v", err)
		return fmt.Errorf("清空用户名历史记录失败: %w", err)
	}

	logger.Info("[用户名历史应用] 成功清空用户名历史记录")
	return nil
}
