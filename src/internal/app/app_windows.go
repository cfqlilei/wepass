//go:build windows
// +build windows

package app

import (
	"errors"
	"fmt"
	"time"
	"wepassword/internal/logger"
	"wepassword/internal/services"

	"github.com/go-vgo/robotgo"
)

/**
 * safeRobotgoWriteAll 安全的剪贴板写入操作
 * @param text 要写入的文本
 * @return error 错误信息
 * @description 20251004 陈凤庆 包装robotgo.WriteAll，添加异常处理
 */
func safeRobotgoWriteAll(text string) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] WriteAll操作发生异常: %v", r)
		}
	}()

	return robotgo.WriteAll(text)
}

/**
 * safeRobotgoReadAll 安全的剪贴板读取操作
 * @return string 剪贴板内容
 * @return error 错误信息
 * @description 20251004 陈凤庆 包装robotgo.ReadAll，添加异常处理
 */
func safeRobotgoReadAll() (string, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] ReadAll操作发生异常: %v", r)
		}
	}()

	return robotgo.ReadAll()
}

/**
 * safeRobotgoPasteStr 安全的粘贴操作
 * @param text 要粘贴的文本
 * @return error 错误信息
 * @description 20251004 陈凤庆 包装robotgo.PasteStr，添加异常处理
 */
func safeRobotgoPasteStr(text string) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] PasteStr操作发生异常: %v", r)
		}
	}()

	return robotgo.PasteStr(text)
}

/**
 * safeRobotgoTypeStr 安全的文本输入操作
 * @param text 要输入的文本
 * @description 20251004 陈凤庆 包装robotgo.TypeStr，添加异常处理
 */
func safeRobotgoTypeStr(text string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] TypeStr操作发生异常: %v", r)
		}
	}()

	robotgo.TypeStr(text)
}

/**
 * safeRobotgoKeyTap 安全的按键操作
 * @param key 要按的键
 * @description 20251004 陈凤庆 包装robotgo.KeyTap，添加异常处理
 */
func safeRobotgoKeyTap(key string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] KeyTap操作发生异常: %v", r)
		}
	}()

	robotgo.KeyTap(key)
}

/**
 * initKeyboardService 初始化键盘服务 (Windows版本)
 * @description 20251003 陈凤庆 Windows平台初始化键盘服务
 * @modify 20251003 陈凤庆 启用Windows平台键盘服务，移除反射机制
 * @modify 20251005 陈凤庆 添加键盘助手服务初始化
 */
func (a *App) initKeyboardService() {
	// 20251003 陈凤庆 Windows平台创建键盘服务实例
	a.keyboardService = services.NewKeyboardService()

	// 20251005 陈凤庆 初始化键盘助手服务和远程输入服务
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		windowMonitor := ks.GetWindowMonitor()
		if windowMonitor != nil {
			// 初始化键盘助手服务
			helperService, err := services.NewKeyboardHelperService(windowMonitor)
			if err != nil {
				fmt.Printf("[App] ⚠️ 键盘助手服务初始化失败: %v\n", err)
			} else {
				a.keyboardHelperService = helperService
				fmt.Printf("[App] ✅ 键盘助手服务初始化成功\n")
			}

			// 初始化远程输入服务
			a.remoteInputService = services.NewRemoteInputService(windowMonitor)
			fmt.Printf("[App] ✅ 远程输入服务初始化成功\n")
		}
	}
}

/**
 * CheckAccessibilityPermission 检查是否有辅助功能权限
 * @return bool 是否有权限
 * @description 20250127 陈凤庆 Windows平台总是返回true
 */
func (a *App) CheckAccessibilityPermission() bool {
	// Windows平台不需要特殊权限
	return true
}

/**
 * StorePreviousFocusedApp 存储当前聚焦的应用程序
 * @description 20250127 陈凤庆 在执行自动填充前调用，记录当前聚焦的应用程序
 * @modify 20251003 陈凤庆 使用interface{}类型避免编译时类型检查问题
 */
func (a *App) StorePreviousFocusedApp() {
	if a.keyboardService != nil {
		// 20251003 陈凤庆 使用类型断言调用Windows版本的KeyboardService方法
		if ks, ok := a.keyboardService.(interface{ StorePreviousFocusedApp() }); ok {
			ks.StorePreviousFocusedApp()
		}
	}
}

/**
 * GetPreviousFocusedAppName 获取上一次聚焦应用程序的名称
 * @return string 应用程序名称
 * @description 20250127 陈凤庆 用于调试和日志记录
 * @modify 20251003 陈凤庆 使用interface{}类型避免编译时类型检查问题
 */
func (a *App) GetPreviousFocusedAppName() string {
	if a.keyboardService != nil {
		// 20251003 陈凤庆 使用类型断言调用Windows版本的KeyboardService方法
		if ks, ok := a.keyboardService.(interface{ GetPreviousFocusedAppName() string }); ok {
			return ks.GetPreviousFocusedAppName()
		}
	}
	return ""
}

/**
 * SimulateUsernameAndPassword 模拟输入用户名和密码
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 实现自动填充功能，先输入用户名，按Tab键，再输入密码
 * @modify 20251003 陈凤庆 根据账号ID查询解密后的完整数据，然后根据输入方式选择不同的输入方法
 */
func (a *App) SimulateUsernameAndPassword(accountId string) error {
	if a.keyboardService != nil {
		logger.Info("账号列表操作: 开始输入用户名和密码，账号ID: %s", accountId)

		// 20251003 陈凤庆 根据账号ID获取解密后的完整账号信息
		username, password, inputMethod, err := a.GetAccountCredentials(accountId)
		if err != nil {
			logger.Error("账号列表操作: 获取账号信息失败，账号ID: %s, 错误: %v", accountId, err)
			return fmt.Errorf("获取账号信息失败: %w", err)
		}

		// 验证用户名和密码不为空
		if username == "" {
			logger.Error("账号列表操作: 用户名为空，账号ID: %s", accountId)
			return fmt.Errorf("用户名不能为空")
		}
		if password == "" {
			logger.Error("账号列表操作: 密码为空，账号ID: %s", accountId)
			return fmt.Errorf("密码不能为空")
		}

		logger.Info("账号列表操作: 准备输入用户名和密码，账号ID: %s, 输入方式: %d, 用户名长度: %d, 密码长度: %d",
			accountId, inputMethod, len(username), len(password))

		// 20251003 陈凤庆 根据账号的输入方式选择不同的输入方法
		err = a.simulateUsernameAndPasswordByMethod(inputMethod, username, password)
		if err != nil {
			logger.Error("[自动填充] 输入用户名和密码失败，账号ID: %s, 输入方式: %d, 错误: %v", accountId, inputMethod, err)
			return err
		}

		logger.Info("[自动填充] 成功输入用户名和密码，账号ID: %s", accountId)
		return nil
	}
	return fmt.Errorf("键盘服务未初始化")
}

/**
 * SimulateUsername 模拟输入用户名
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 实现用户名自动填充功能，包含窗口聚焦恢复
 * @modify 20251003 陈凤庆 根据账号ID查询解密后的用户名，然后根据输入方式选择不同的输入方法
 */
func (a *App) SimulateUsername(accountId string) error {
	if a.keyboardService != nil {
		logger.Info("账号列表操作: 开始输入用户名，账号ID: %s", accountId)

		// 20251003 陈凤庆 根据账号ID获取解密后的用户名
		username, _, inputMethod, err := a.GetAccountCredentials(accountId)
		if err != nil {
			logger.Error("账号列表操作: 获取账号信息失败，账号ID: %s, 错误: %v", accountId, err)
			return fmt.Errorf("获取账号信息失败: %w", err)
		}

		// 验证用户名不为空
		if username == "" {
			logger.Error("账号列表操作: 用户名为空，账号ID: %s", accountId)
			return fmt.Errorf("用户名不能为空")
		}

		logger.Info("账号列表操作: 准备输入用户名，账号ID: %s, 输入方式: %d, 用户名长度: %d", accountId, inputMethod, len(username))

		// 20251003 陈凤庆 根据账号的输入方式选择不同的输入方法
		return a.simulateTextByMethod(inputMethod, username)
	}
	return fmt.Errorf("键盘服务未初始化")
}

/**
 * SimulatePassword 模拟输入密码
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 实现密码自动填充功能，包含窗口聚焦恢复
 * @modify 20251003 陈凤庆 根据账号ID查询解密后的密码，然后根据输入方式选择不同的输入方法
 */
func (a *App) SimulatePassword(accountId string) error {
	if a.keyboardService != nil {
		logger.Info("账号列表操作: 开始输入密码，账号ID: %s", accountId)

		// 20251003 陈凤庆 根据账号ID获取解密后的密码
		_, password, inputMethod, err := a.GetAccountCredentials(accountId)
		if err != nil {
			logger.Error("账号列表操作: 获取账号信息失败，账号ID: %s, 错误: %v", accountId, err)
			return fmt.Errorf("获取账号信息失败: %w", err)
		}

		// 验证密码不为空
		if password == "" {
			logger.Error("账号列表操作: 密码为空，账号ID: %s", accountId)
			return fmt.Errorf("密码不能为空")
		}

		logger.Info("账号列表操作: 准备输入密码，账号ID: %s, 输入方式: %d, 密码长度: %d", accountId, inputMethod, len(password))

		// 20251003 陈凤庆 根据账号的输入方式选择不同的输入方法
		err = a.simulateTextByMethod(inputMethod, password)
		if err != nil {
			logger.Error("[自动填充] 输入密码失败，账号ID: %s, 输入方式: %d, 错误: %v", accountId, inputMethod, err)
			return err
		}

		logger.Info("[自动填充] 成功输入密码，账号ID: %s", accountId)
		return nil
	}
	return fmt.Errorf("键盘服务未初始化")
}

/**
 * RecordLastWindow 记录当前活动窗口为lastPID
 * @return error 错误信息
 * @description 20251002 陈凤庆 手动记录当前活动窗口，用于鼠标悬停等场景
 * @modify 20251003 陈凤庆 使用interface{}类型避免编译时类型检查问题
 */
func (a *App) RecordLastWindow() error {
	if a.keyboardService != nil {
		// 20251003 陈凤庆 使用类型断言调用Windows版本的KeyboardService方法
		if ks, ok := a.keyboardService.(interface{ RecordCurrentWindow() error }); ok {
			return ks.RecordCurrentWindow()
		}
	}
	return nil
}

// WindowInfoInterface 定义窗口信息接口，避免直接依赖具体类型
type WindowInfoInterface interface {
	GetPID() int
	GetTitle() string
	GetTimestamp() int64
}

/**
 * GetLastWindowInfo 获取最后活动窗口信息
 * @return map[string]interface{} 窗口信息
 * @description 20251002 陈凤庆 获取lastPID对应的窗口信息
 * @modify 20251003 陈凤庆 使用interface{}类型避免编译时类型检查问题
 */
func (a *App) GetLastWindowInfo() map[string]interface{} {
	if a.keyboardService != nil {
		// 20251003 陈凤庆 使用类型断言调用Windows版本的KeyboardService方法
		// 由于WindowInfo结构体在不同平台有不同定义，这里简化处理
		if ks, ok := a.keyboardService.(interface{ GetLastWindowInfo() interface{} }); ok {
			windowInfo := ks.GetLastWindowInfo()
			if windowInfo != nil {
				// 尝试通过interface{}获取基本信息
				return map[string]interface{}{
					"pid":       0,      // 暂时返回默认值
					"title":     "未知窗口", // 暂时返回默认值
					"timestamp": 0,      // 暂时返回默认值
				}
			}
		}
	}

	return map[string]interface{}{
		"pid":       0,
		"title":     "",
		"timestamp": 0,
	}
}

/**
 * SwitchToLastWindow 切换到最后活动的窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活lastPID对应的窗口为活动窗口
 * @modify 20251003 陈凤庆 使用interface{}类型避免编译时类型检查问题
 */
func (a *App) SwitchToLastWindow() bool {
	if a.keyboardService != nil {
		// 20251003 陈凤庆 使用类型断言调用Windows版本的KeyboardService方法
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			return ks.SwitchToLastWindow()
		}
	}
	return false
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活密码管理器窗口为活动窗口
 * @modify 20251003 陈凤庆 使用interface{}类型避免编译时类型检查问题
 */
func (a *App) SwitchToPasswordManager() bool {
	if a.keyboardService != nil {
		// 20251003 陈凤庆 使用类型断言调用Windows版本的KeyboardService方法
		if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
			return ks.SwitchToPasswordManager()
		}
	}
	return false
}

/**
 * simulateTextByMethod 根据输入方式模拟输入文本 (Windows版本)
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251003 陈凤庆 根据不同的输入方式选择对应的输入方法
 * @modify 20251005 陈凤庆 添加第4种输入方式：键盘助手输入（删除原第4种底层键盘API）
 * @modify 20251005 陈凤庆 添加第5种输入方式：远程输入（专为ToDesk等远程桌面环境设计）
 */
func (a *App) simulateTextByMethod(inputMethod int, text string) error {
	if a.keyboardService != nil {
		logger.Info("模拟输入: 开始文本输入，输入方式: %d, 文本长度: %d", inputMethod, len(text))

		switch inputMethod {
		case 1:
			// 默认方式：使用Unicode发送按键码方式，速度快
			logger.Info("模拟输入: 使用默认方式(Unicode)输入文本")
			if ks, ok := a.keyboardService.(interface{ SimulateTextWithAutoSwitch(string) error }); ok {
				err := ks.SimulateTextWithAutoSwitch(text)
				if err != nil {
					logger.Error("模拟输入: 默认方式输入失败，错误: %v", err)
				} else {
					logger.Info("模拟输入: 默认方式输入成功")
				}
				return err
			}
		case 2:
			// 模拟键盘输入：使用robotgo.KeyTap逐键输入
			logger.Info("模拟输入: 使用模拟键盘输入(robotgo.KeyTap)方式")
			err := a.simulateTextWithRobotgoKeyTap(text)
			if err != nil {
				logger.Error("模拟输入: 模拟键盘输入失败，错误: %v", err)
			} else {
				logger.Info("模拟输入: 模拟键盘输入成功")
			}
			return err
		case 3:
			// 复制粘贴输入：使用robotgo.PasteStr
			logger.Info("模拟输入: 使用复制粘贴输入(robotgo.PasteStr)方式")
			err := a.simulateTextWithRobotgoPaste(text)
			if err != nil {
				logger.Error("模拟输入: 复制粘贴输入失败，错误: %v", err)
			} else {
				logger.Info("模拟输入: 复制粘贴输入成功")
			}
			return err
		case 4:
			// 键盘助手输入：使用系统键盘助手
			logger.Info("模拟输入: 使用键盘助手输入方式")
			if a.keyboardHelperService != nil {
				err := a.keyboardHelperService.SimulateText(text)
				if err != nil {
					logger.Error("模拟输入: 键盘助手输入失败，错误: %v", err)
				} else {
					logger.Info("模拟输入: 键盘助手输入成功")
				}
				return err
			}
			logger.Error("模拟输入: 键盘助手服务未初始化")
			return fmt.Errorf("键盘助手服务未初始化")
		case 5:
			// 远程输入：专为ToDesk等远程桌面环境设计，逐字符输入
			logger.Info("模拟输入: 使用远程输入方式(ToDesk等远程桌面)")
			if a.remoteInputService != nil {
				err := a.remoteInputService.SimulateText(text)
				if err != nil {
					logger.Error("模拟输入: 远程输入失败，错误: %v", err)
				} else {
					logger.Info("模拟输入: 远程输入成功")
				}
				return err
			}
			logger.Error("模拟输入: 远程输入服务未初始化")
			return fmt.Errorf("远程输入服务未初始化")
		default:
			// 未知方式，使用默认方式
			logger.Info("模拟输入: 未知输入方式 %d，使用默认方式", inputMethod)
			if ks, ok := a.keyboardService.(interface{ SimulateTextWithAutoSwitch(string) error }); ok {
				err := ks.SimulateTextWithAutoSwitch(text)
				if err != nil {
					logger.Error("模拟输入: 默认方式输入失败，错误: %v", err)
				} else {
					logger.Info("模拟输入: 默认方式输入成功")
				}
				return err
			}
		}
	}
	return nil
}

/**
 * simulateUsernameAndPasswordByMethod 根据输入方式模拟输入用户名和密码 (Windows版本)
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251003 陈凤庆 根据不同的输入方式选择对应的输入方法
 * @modify 20251005 陈凤庆 添加第4种输入方式：键盘助手输入（删除原第4种底层键盘API）
 * @modify 20251005 陈凤庆 添加第5种输入方式：远程输入（专为ToDesk等远程桌面环境设计）
 */
func (a *App) simulateUsernameAndPasswordByMethod(inputMethod int, username, password string) error {
	if a.keyboardService != nil {
		switch inputMethod {
		case 1:
			// 默认方式：使用Unicode发送按键码方式，速度快
			if ks, ok := a.keyboardService.(interface{ SimulateUsernameAndPasswordWithAutoSwitch(string, string) error }); ok {
				return ks.SimulateUsernameAndPasswordWithAutoSwitch(username, password)
			}
		case 2:
			// 模拟键盘输入：使用robotgo.KeyTap逐键输入
			return a.simulateUsernameAndPasswordWithRobotgoKeyTap(username, password)
		case 3:
			// 复制粘贴输入：使用robotgo.PasteStr
			return a.simulateUsernameAndPasswordWithRobotgoPaste(username, password)
		case 4:
			// 键盘助手输入：使用系统键盘助手
			if a.keyboardHelperService != nil {
				return a.keyboardHelperService.SimulateUsernameAndPassword(username, password)
			}
			return fmt.Errorf("键盘助手服务未初始化")
		case 5:
			// 远程输入：专为ToDesk等远程桌面环境设计，逐字符输入
			if a.remoteInputService != nil {
				return a.remoteInputService.SimulateUsernameAndPassword(username, password)
			}
			return fmt.Errorf("远程输入服务未初始化")
		default:
			// 未知方式，使用默认方式
			if ks, ok := a.keyboardService.(interface{ SimulateUsernameAndPasswordWithAutoSwitch(string, string) error }); ok {
				return ks.SimulateUsernameAndPasswordWithAutoSwitch(username, password)
			}
		}
	}
	return nil
}

/**
 * simulateTextWithRobotgoKeyTap 使用robotgo.KeyTap逐键输入文本 (Windows版本)
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo逐键输入，兼容性好但速度较慢
 * @modify 20251004 陈凤庆 添加输入验证和错误处理，防止输入错误字符
 */
func (a *App) simulateTextWithRobotgoKeyTap(text string) error {
	if a.keyboardService != nil {
		logger.Info("模拟键盘输入: 开始robotgo.KeyTap输入，文本长度: %d", len(text))

		// 验证输入文本
		if text == "" {
			logger.Error("模拟键盘输入: 输入文本为空")
			return errors.New("输入文本不能为空")
		}

		// 验证文本内容，确保不是损坏的数据
		if len(text) < 2 && text != "" {
			logger.Error("模拟键盘输入: 检测到可疑的短文本: '%s'，长度: %d", text, len(text))
		}

		// 检查是否包含重复字符（如"aaa"）
		if len(text) > 2 {
			allSame := true
			firstChar := text[0]
			for _, char := range text {
				if char != rune(firstChar) {
					allSame = false
					break
				}
			}
			if allSame {
				logger.Error("模拟键盘输入: 检测到全部相同字符的文本: '%s'，可能是输入错误", text)
				return fmt.Errorf("检测到异常重复字符，可能是数据损坏: %s", text)
			}
		}

		logger.Info("模拟键盘输入: 开始逐键输入文本，长度: %d", len(text))

		// 切换到目标窗口
		logger.Info("模拟键盘输入: 切换到目标窗口")
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			if !ks.SwitchToLastWindow() {
				logger.Error("模拟键盘输入: 无法切换到目标窗口")
				return errors.New("无法切换到目标窗口")
			}
		}

		// 等待窗口切换完成 - 参考Mac版本增加延迟时间
		logger.Info("模拟键盘输入: 等待窗口切换完成 (500ms)")
		time.Sleep(500 * time.Millisecond)

		// 使用defer确保异常时也能切换回密码管理器
		defer func() {
			if r := recover(); r != nil {
				logger.Error("模拟键盘输入: 发生异常: %v", r)
				// 尝试切换回密码管理器
				if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
					if !ks.SwitchToPasswordManager() {
						logger.Error("模拟键盘输入: 异常恢复时切换回密码管理器失败")
					}
				}
			}
		}()

		// 输入文本，添加异常处理
		func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("模拟键盘输入: TypeStr操作发生异常: %v", r)
					panic(r) // 重新抛出异常，让外层defer处理
				}
			}()

			logger.Info("模拟键盘输入: 准备输入文本，长度: %d", len(text))
			robotgo.TypeStr(text)
		}()

		logger.Info("模拟键盘输入: 文本输入完成")

		// 等待输入完成 - 参考Mac版本增加延迟时间
		logger.Info("模拟键盘输入: 等待输入处理完成 (400ms)")
		time.Sleep(400 * time.Millisecond)

		// 切换回密码管理器
		logger.Info("模拟键盘输入: 切换回密码管理器")
		if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
			if !ks.SwitchToPasswordManager() {
				logger.Error("模拟键盘输入: 切换回密码管理器失败")
			}
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateTextWithRobotgoPaste 使用robotgo.WriteAll复制粘贴输入文本 (Windows版本)
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo复制粘贴方式输入，速度快且兼容性好
 * @modify 20251004 陈凤庆 添加错误处理和异常恢复机制，防止程序崩溃
 * @modify 20251004 陈凤庆 改进剪贴板操作的安全性和稳定性
 */
func (a *App) simulateTextWithRobotgoPaste(text string) error {
	if a.keyboardService != nil {
		logger.Info("复制粘贴输入: 开始复制粘贴输入，文本长度: %d", len(text))

		// 验证输入文本
		if text == "" {
			logger.Error("复制粘贴输入: 输入文本为空")
			return errors.New("输入文本不能为空")
		}

		// 验证文本长度，避免过长文本导致问题
		if len(text) > 1000 {
			logger.Error("复制粘贴输入: 文本过长 (%d 字符)，可能影响性能", len(text))
		}

		logger.Info("复制粘贴输入: 准备复制文本到剪贴板，长度: %d", len(text))

		// 注意：剪贴板内容将在10秒后自动清理，确保安全性

		// 步骤1：复制文本到剪贴板
		logger.Info("复制粘贴输入: 步骤1 - 复制文本到剪贴板")
		err := safeRobotgoWriteAll(text)
		if err != nil {
			logger.Error("复制粘贴输入: 复制文本到剪贴板失败: %v", err)
			return fmt.Errorf("复制文本到剪贴板失败: %w", err)
		}
		logger.Info("复制粘贴输入: ✅ 文本已复制到剪贴板")

		// 启动10秒后自动销毁剪贴板内容的goroutine
		go func() {
			time.Sleep(10 * time.Second)
			currentClipboard, _ := safeRobotgoReadAll()
			if currentClipboard == text {
				safeRobotgoWriteAll("")
				logger.Info("复制粘贴输入: 🔒 剪贴板内容已自动销毁")
			}
		}()

		// 步骤2：切换到目标窗口
		logger.Info("复制粘贴输入: 步骤2 - 切换到目标窗口")
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			if !ks.SwitchToLastWindow() {
				logger.Error("复制粘贴输入: 无法切换到目标窗口")
				return errors.New("无法切换到目标窗口")
			}
		}

		// 等待窗口切换完成 - 参考Mac版本增加延迟时间
		logger.Info("复制粘贴输入: 等待窗口切换完成 (500ms)")
		time.Sleep(500 * time.Millisecond)

		// 使用defer确保异常时也能切换回密码管理器和恢复剪贴板
		defer func() {
			if r := recover(); r != nil {
				logger.Error("复制粘贴输入: 发生异常: %v", r)
			}

			// 注意：剪贴板内容将在10秒后自动清理，不立即恢复原始内容
			logger.Info("复制粘贴输入: 剪贴板内容将在10秒后自动清理")

			// 尝试切换回密码管理器
			logger.Info("复制粘贴输入: 切换回密码管理器")
			if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
				if !ks.SwitchToPasswordManager() {
					logger.Error("复制粘贴输入: 切换回密码管理器失败")
				}
			}
		}()

		// 步骤3：执行粘贴操作，添加多重保护
		logger.Info("复制粘贴输入: 步骤3 - 开始粘贴操作")
		pasteErr := func() error {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("复制粘贴输入: 剪贴板操作发生异常: %v", r)
					panic(fmt.Errorf("剪贴板操作异常: %v", r))
				}
			}()

			// 验证剪贴板内容
			time.Sleep(100 * time.Millisecond)
			clipboardContent, err := safeRobotgoReadAll()
			if err != nil {
				logger.Error("复制粘贴输入: 读取剪贴板失败: %v", err)
				return fmt.Errorf("读取剪贴板失败: %w", err)
			}
			if clipboardContent != text {
				logger.Error("复制粘贴输入: 剪贴板内容不匹配，期望长度: %d，实际长度: %d", len(text), len(clipboardContent))
			} else {
				logger.Info("复制粘贴输入: 剪贴板内容验证成功")
			}

			logger.Info("复制粘贴输入: 步骤3 - 执行粘贴操作")
			// 使用PasteStr方法，但增加错误处理
			err = safeRobotgoPasteStr(text)
			if err != nil {
				logger.Error("复制粘贴输入: 粘贴操作失败: %v", err)
				return fmt.Errorf("粘贴操作失败: %w", err)
			}

			return nil
		}()

		if pasteErr != nil {
			logger.Error("复制粘贴输入: 粘贴操作失败: %v", pasteErr)
			return pasteErr
		}

		logger.Info("复制粘贴输入: 文本输入完成")

		// 步骤4：等待粘贴操作完成 - 参考Mac版本增加延迟时间
		logger.Info("复制粘贴输入: 步骤4 - 等待粘贴操作处理完成 (500ms)")
		time.Sleep(500 * time.Millisecond)

		// 步骤5：切换回密码管理器
		logger.Info("复制粘贴输入: 步骤5 - 切换回密码管理器")
		if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
			if !ks.SwitchToPasswordManager() {
				logger.Error("复制粘贴输入: 切换回密码管理器失败")
			} else {
				logger.Info("复制粘贴输入: ✅ 成功切换回密码管理器")
			}
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateUsernameAndPasswordWithRobotgoKeyTap 使用robotgo.KeyTap逐键输入用户名和密码 (Windows版本)
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo逐键输入用户名和密码，兼容性好但速度较慢
 * @modify 20251009 陈凤庆 还原为原始的逐键输入方式，这个功能原本是正确的
 */
func (a *App) simulateUsernameAndPasswordWithRobotgoKeyTap(username, password string) error {
	if a.keyboardService != nil {
		// 切换到目标窗口
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			if !ks.SwitchToLastWindow() {
				return errors.New("无法切换到目标窗口")
			}
		}

		// 等待窗口切换完成 - 参考Mac版本增加延迟时间
		fmt.Printf("[robotgo输入] 等待窗口切换完成 (500ms)\n")
		time.Sleep(500 * time.Millisecond)

		// 输入用户名
		if username != "" {
			fmt.Printf("[robotgo输入] 开始逐键输入用户名: %s\n", username)
			safeRobotgoTypeStr(username)
			fmt.Printf("[robotgo输入] ✅ 用户名输入完成\n")
			// 参考Mac版本增加用户名输入后的延迟
			fmt.Printf("[robotgo输入] 等待用户名输入处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 按Tab键切换到密码输入框
		fmt.Printf("[robotgo输入] 按Tab键切换到密码输入框\n")
		safeRobotgoKeyTap("tab")
		// 参考Mac版本增加Tab键后的延迟
		fmt.Printf("[robotgo输入] 等待Tab键焦点切换完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)

		// 输入密码
		if password != "" {
			fmt.Printf("[robotgo输入] 开始输入密码\n")
			safeRobotgoTypeStr(password)
			fmt.Printf("[robotgo输入] ✅ 密码输入完成\n")
			// 参考Mac版本增加密码输入后的延迟
			fmt.Printf("[robotgo输入] 等待密码输入处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 全部输入完成后等待一段时间再切换回密码管理器 - 参考Mac版本
		fmt.Printf("[robotgo输入] 全部输入完成，等待目标窗口完全处理 (600ms)\n")
		time.Sleep(600 * time.Millisecond)

		// 切换回密码管理器
		if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
			if !ks.SwitchToPasswordManager() {
				fmt.Printf("[robotgo输入] ⚠️ 切换回密码管理器失败\n")
			}
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateUsernameAndPasswordWithRobotgoPaste 使用robotgo.WriteAll复制粘贴输入用户名和密码 (Windows版本)
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo复制粘贴方式输入用户名和密码，速度快且兼容性好
 * @modify 20251005 陈凤庆 按照需求文档重构：获取targetInputPID后保持在目标窗口完成所有输入
 */
func (a *App) simulateUsernameAndPasswordWithRobotgoPaste(username, password string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		fmt.Printf("[复制粘贴组合] ========== 开始输入用户名和密码 ==========\n")

		// 步骤1：从lastPID获取目标窗口PID，在整个输入过程保持不变
		wm := ks.GetWindowMonitor()
		if wm == nil {
			return fmt.Errorf("窗口监控服务未初始化")
		}
		lastWindow := wm.GetLastWindow()
		if lastWindow == nil {
			return fmt.Errorf("无法获取目标窗口信息")
		}
		targetInputPID := lastWindow.PID
		fmt.Printf("[复制粘贴组合] 步骤1：获取目标窗口PID: %d (窗口: %s)\n", targetInputPID, lastWindow.Title)

		// 步骤2：复制用户名到剪切板
		if username != "" {
			fmt.Printf("[复制粘贴组合] 步骤2：复制用户名到剪切板\n")
			err := safeRobotgoWriteAll(username)
			if err != nil {
				fmt.Printf("[复制粘贴组合] ❌ 用户名写入剪贴板失败: %v\n", err)
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 用户名已复制到剪贴板: %s\n", username)

			// 启动10秒后自动销毁剪贴板内容的goroutine
			go func() {
				time.Sleep(10 * time.Second)
				currentClipboard, _ := safeRobotgoReadAll()
				if currentClipboard == username {
					safeRobotgoWriteAll("")
					fmt.Printf("[复制粘贴组合] 🔒 用户名剪贴板内容已自动销毁\n")
				}
			}()
		}

		// 步骤3：定位到目标输入窗口
		fmt.Printf("[复制粘贴组合] 步骤3：切换到目标输入窗口 (PID: %d)\n", targetInputPID)
		if !ks.SwitchToLastWindow() {
			return errors.New("无法切换到目标窗口")
		}
		// 参考Mac版本增加窗口切换延迟
		fmt.Printf("[复制粘贴组合] 等待窗口切换完成 (500ms)\n")
		time.Sleep(500 * time.Millisecond)

		// 步骤4：在目标窗口粘贴用户名
		if username != "" {
			fmt.Printf("[复制粘贴组合] 步骤4：在目标窗口粘贴用户名\n")
			err := safeRobotgoPasteStr(username)
			if err != nil {
				fmt.Printf("[复制粘贴组合] ⚠️ 粘贴用户名失败: %v\n", err)
				// 尝试切换回密码管理器
				ks.SwitchToPasswordManager()
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 用户名粘贴完成\n")
			// 参考Mac版本增加用户名粘贴后的延迟
			fmt.Printf("[复制粘贴组合] 等待用户名粘贴处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 步骤5：在目标输入窗口发送Tab键
		fmt.Printf("[复制粘贴组合] 步骤5：在目标窗口发送Tab键\n")
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[复制粘贴组合] ⚠️ Tab键输入发生异常: %v\n", r)
				}
			}()
			robotgo.KeyTap("tab")
		}()
		// 参考Mac版本增加Tab键后的延迟
		fmt.Printf("[复制粘贴组合] 等待Tab键焦点切换完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("[复制粘贴组合] ✅ Tab键发送完成\n")

		// 步骤6：复制密码到剪切板
		if password != "" {
			fmt.Printf("[复制粘贴组合] 步骤6：复制密码到剪切板\n")
			err := safeRobotgoWriteAll(password)
			if err != nil {
				fmt.Printf("[复制粘贴组合] ❌ 密码写入剪贴板失败: %v\n", err)
				// 尝试切换回密码管理器
				ks.SwitchToPasswordManager()
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 密码已复制到剪贴板\n")

			// 启动10秒后自动销毁剪贴板内容的goroutine
			go func() {
				time.Sleep(10 * time.Second)
				currentClipboard, _ := safeRobotgoReadAll()
				if currentClipboard == password {
					safeRobotgoWriteAll("")
					fmt.Printf("[复制粘贴组合] 🔒 密码剪贴板内容已自动销毁\n")
				}
			}()

			// 步骤7：在目标窗口粘贴密码
			fmt.Printf("[复制粘贴组合] 步骤7：在目标窗口粘贴密码\n")
			// 参考Mac版本增加粘贴密码前的延迟
			fmt.Printf("[复制粘贴组合] 等待剪贴板内容更新 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
			err = safeRobotgoPasteStr(password)
			if err != nil {
				fmt.Printf("[复制粘贴组合] ⚠️ 粘贴密码失败: %v\n", err)
				// 尝试切换回密码管理器
				ks.SwitchToPasswordManager()
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 密码粘贴完成\n")
			// 参考Mac版本增加密码粘贴后的延迟
			fmt.Printf("[复制粘贴组合] 等待密码粘贴处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 全部输入完成后等待一段时间再切换回密码管理器 - 参考Mac版本
		fmt.Printf("[复制粘贴组合] 全部输入完成，等待目标窗口完全处理 (600ms)\n")
		time.Sleep(600 * time.Millisecond)

		// 步骤8：切换回当前程序主窗口
		fmt.Printf("[复制粘贴组合] 步骤8：切换回密码管理器主窗口\n")
		if !ks.SwitchToPasswordManager() {
			fmt.Printf("[复制粘贴组合] ⚠️ 切换回密码管理器失败\n")
		}

		fmt.Printf("[复制粘贴组合] ========== 用户名和密码输入完成 ==========\n")
		return nil
	}
	return errors.New("键盘服务未初始化")
}
