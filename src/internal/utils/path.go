package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

/**
 * 跨平台路径工具函数
 * @author 陈凤庆
 * @date 2025-10-02
 * @description 提供跨平台的应用数据存储路径获取功能
 */

/**
 * GetAppDataDir 获取应用数据存储目录
 * @return string 应用数据存储目录路径
 * @return error 错误信息
 * @description 根据不同操作系统返回合适的应用数据存储路径：
 * - macOS: ~/Library/Application Support/wepass
 * - Windows: %APPDATA%/wepass
 * - Linux: ~/.config/wepass
 */
func GetAppDataDir() (string, error) {
	var appDataDir string

	switch runtime.GOOS {
	case "darwin": // macOS
		// 获取用户主目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		appDataDir = filepath.Join(homeDir, "Library", "Application Support", "wepass")

	case "windows": // Windows
		// 使用 APPDATA 环境变量
		appDataDir = os.Getenv("APPDATA")
		if appDataDir == "" {
			// 如果 APPDATA 不存在，使用用户主目录下的 AppData/Roaming
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appDataDir = filepath.Join(homeDir, "AppData", "Roaming")
		}
		appDataDir = filepath.Join(appDataDir, "wepass")

	default: // Linux 和其他 Unix 系统
		// 首先尝试使用 XDG_CONFIG_HOME 环境变量
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			// 如果 XDG_CONFIG_HOME 不存在，使用 ~/.config
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(homeDir, ".config")
		}
		appDataDir = filepath.Join(configDir, "wepass")
	}

	return appDataDir, nil
}

/**
 * GetVaultDataDir 获取密码库数据存储目录
 * @return string 密码库数据存储目录路径
 * @return error 错误信息
 * @description 在应用数据目录下创建 vaults 子目录用于存储密码库文件
 */
func GetVaultDataDir() (string, error) {
	appDataDir, err := GetAppDataDir()
	if err != nil {
		return "", err
	}

	vaultDataDir := filepath.Join(appDataDir, "vaults")
	return vaultDataDir, nil
}

/**
 * EnsureDir 确保目录存在
 * @param dirPath 目录路径
 * @return error 错误信息
 * @description 如果目录不存在则创建，包括所有必要的父目录
 */
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

/**
 * GetDefaultVaultPath 获取默认密码库文件路径
 * @param vaultName 密码库名称（无需后缀）
 * @return string 完整的密码库文件路径
 * @return error 错误信息
 * @description 在密码库数据目录下生成完整的密码库文件路径，自动添加.db后缀
 */
func GetDefaultVaultPath(vaultName string) (string, error) {
	vaultDataDir, err := GetVaultDataDir()
	if err != nil {
		return "", err
	}

	// 确保目录存在
	if err := EnsureDir(vaultDataDir); err != nil {
		return "", err
	}

	// 构建完整的密码库文件路径，自动添加.db后缀
	vaultPath := filepath.Join(vaultDataDir, vaultName+".db")
	return vaultPath, nil
}
