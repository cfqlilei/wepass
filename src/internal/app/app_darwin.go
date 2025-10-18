//go:build darwin
// +build darwin

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -framework ApplicationServices

#import <ApplicationServices/ApplicationServices.h>

// 使用CGEvent发送Tab键到活动窗口
void sendTabKeyWithCGEvent() {
    CGEventRef keyDownEvent = CGEventCreateKeyboardEvent(NULL, 48, true); // Tab键的虚拟键码是48
    CGEventRef keyUpEvent = CGEventCreateKeyboardEvent(NULL, 48, false);

    CGEventPost(kCGHIDEventTap, keyDownEvent);
    usleep(20000); // 20ms延迟
    CGEventPost(kCGHIDEventTap, keyUpEvent);

    CFRelease(keyDownEvent);
    CFRelease(keyUpEvent);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"wepassword/internal/logger"
	"wepassword/internal/services"

	"github.com/go-vgo/robotgo"
)

/**
 * safeRobotgoWriteAll 安全的剪贴板写入操作 (Mac版本)
 * @param text 要写入的文本
 * @return error 错误信息
 * @description 20251004 陈凤庆 包装robotgo.WriteAll，添加异常处理，防止程序崩溃
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
 * safeRobotgoReadAll 安全的剪贴板读取操作 (Mac版本)
 * @return string 剪贴板内容
 * @return error 错误信息
 * @description 20251004 陈凤庆 包装robotgo.ReadAll，添加异常处理，防止程序崩溃
 */
func safeRobotgoReadAll() (string, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] ReadAll操作发生异常: %v", r)
		}
	}()

	return robotgo.ReadAll()
}

// safeRobotgoPasteStr 已移除，因为robotgo.PasteStr在Mac上会导致SIGTRAP崩溃
// 现在使用 safeRobotgoWriteAll + Cmd+V 组合键的方式替代

/**
 * sendTabKeyWithCGEventWrapper 使用CGEvent发送Tab键
 * @description 20251005 陈凤庆 使用CGEvent发送Tab键，确保Tab键正常输入
 */
func sendTabKeyWithCGEventWrapper() {
	fmt.Printf("[CGEvent] 使用CGEvent发送Tab键\n")
	C.sendTabKeyWithCGEvent()
	fmt.Printf("[CGEvent] ✅ Tab键发送完成\n")
}

/**
 * safeRobotgoTypeStr 安全的文本输入操作 (Mac版本)
 * @param text 要输入的文本
 * @description 20251004 陈凤庆 包装robotgo.TypeStr，添加异常处理，防止程序崩溃
 */
func safeRobotgoTypeStr(text string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[robotgo] TypeStr操作发生异常: %v", r)
		}
	}()

	robotgo.TypeStr(text)
}

// 已移除所有robotgo键盘模拟函数，因为它们都会调用keyCodeForChar导致SIGTRAP崩溃
// 包括：safeRobotgoKeyTap, safeRobotgoKeyDown, safeRobotgoKeyUp
// 现在使用AppleScript实现自动粘贴功能

/**
 * executeAppleScriptPaste 使用AppleScript执行Cmd+V粘贴操作
 * @return error 执行错误
 * @description 20251004 陈凤庆 使用AppleScript替代robotgo键盘模拟，避免SIGTRAP崩溃
 */
func executeAppleScriptPaste() error {
	// AppleScript命令：模拟按下Cmd+V
	script := `tell application "System Events" to keystroke "v" using command down`

	// 执行AppleScript
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("AppleScript执行失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

/**
 * executeAppleScriptKeyPress 使用AppleScript执行按键操作
 * @param key 要按的键
 * @return error 执行错误
 * @description 20251004 陈凤庆 使用AppleScript替代robotgo键盘模拟
 */
func executeAppleScriptKeyPress(key string) error {
	// AppleScript命令：按下指定键
	script := fmt.Sprintf(`tell application "System Events" to keystroke "%s"`, key)

	// 执行AppleScript
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("AppleScript按键执行失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

/**
 * executeAppleScriptTab 使用AppleScript执行Tab键操作
 * @return error 执行错误
 * @description 20251004 陈凤庆 使用AppleScript模拟Tab键切换
 */
func executeAppleScriptTab() error {
	// AppleScript命令：按下Tab键
	script := `tell application "System Events" to keystroke tab`

	// 执行AppleScript
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("AppleScript Tab键执行失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

/**
 * initKeyboardService 初始化键盘服务 (macOS版本)
 * @description 20251003 陈凤庆 macOS平台初始化键盘服务
 * @modify 20251005 陈凤庆 添加键盘助手服务初始化
 */
func (a *App) initKeyboardService() {
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
 * @description 20250127 陈凤庆 检查macOS辅助功能权限
 */
func (a *App) CheckAccessibilityPermission() bool {
	if a.keyboardService == nil {
		return false
	}
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok {
		return ks.CheckAccessibilityPermission()
	}
	return false
}

/**
 * StorePreviousFocusedApp 存储当前聚焦的应用程序
 * @description 20250127 陈凤庆 在执行自动填充前调用，记录当前聚焦的应用程序
 */
func (a *App) StorePreviousFocusedApp() {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		ks.StorePreviousFocusedApp()
	}
}

/**
 * GetPreviousFocusedAppName 获取上一次聚焦应用程序的名称
 * @return string 应用程序名称
 * @description 20250127 陈凤庆 用于调试和日志记录
 */
func (a *App) GetPreviousFocusedAppName() string {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		return ks.GetPreviousFocusedAppName()
	}
	return ""
}

/**
 * SimulateUsernameAndPassword 模拟输入用户名和密码
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 实现自动填充功能，先输入用户名，按Tab键，再输入密码
 * @modify 陈凤庆 使用新的自动切换功能
 * @modify 20251003 陈凤庆 根据账号ID查询解密后的完整数据，然后根据输入方式选择不同的输入方法
 */
func (a *App) SimulateUsernameAndPassword(accountId string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		// 20251003 陈凤庆 根据账号ID获取解密后的完整账号信息
		username, password, inputMethod, err := a.GetAccountCredentials(accountId)
		if err != nil {
			logger.Error("[自动填充] 获取账号凭据失败，账号ID: %s, 错误: %v", accountId, err)
			return fmt.Errorf("获取账号信息失败: %w", err)
		}

		// 验证用户名和密码不为空
		if username == "" {
			return fmt.Errorf("用户名不能为空")
		}
		if password == "" {
			return fmt.Errorf("密码不能为空")
		}

		logger.Info("[自动填充] 开始输入用户名和密码，账号ID: %s, 输入方式: %d, 用户名长度: %d, 密码长度: %d",
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
 * @modify 陈凤庆 使用新的自动切换功能
 * @modify 20251003 陈凤庆 根据账号ID查询解密后的用户名，然后根据输入方式选择不同的输入方法
 */
func (a *App) SimulateUsername(accountId string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		// 20251003 陈凤庆 根据账号ID获取解密后的用户名
		username, _, inputMethod, err := a.GetAccountCredentials(accountId)
		if err != nil {
			logger.Error("[自动填充] 获取账号凭据失败，账号ID: %s, 错误: %v", accountId, err)
			return fmt.Errorf("获取账号信息失败: %w", err)
		}

		// 验证用户名不为空
		if username == "" {
			return fmt.Errorf("用户名不能为空")
		}

		logger.Info("[自动填充] 开始输入用户名，账号ID: %s, 输入方式: %d", accountId, inputMethod)

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
 * @modify 陈凤庆 使用新的自动切换功能
 * @modify 20251003 陈凤庆 根据账号ID查询解密后的密码，然后根据输入方式选择不同的输入方法
 */
func (a *App) SimulatePassword(accountId string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		// 20251003 陈凤庆 根据账号ID获取解密后的密码
		_, password, inputMethod, err := a.GetAccountCredentials(accountId)
		if err != nil {
			logger.Error("[自动填充] 获取账号凭据失败，账号ID: %s, 错误: %v", accountId, err)
			return fmt.Errorf("获取账号信息失败: %w", err)
		}

		// 验证密码不为空
		if password == "" {
			return fmt.Errorf("密码不能为空")
		}

		logger.Info("[自动填充] 开始输入密码，账号ID: %s, 输入方式: %d, 密码长度: %d", accountId, inputMethod, len(password))

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
 */
func (a *App) RecordLastWindow() error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		return ks.RecordCurrentWindow()
	}
	return nil
}

/**
 * GetLastWindowInfo 获取最后活动窗口信息
 * @return map[string]interface{} 窗口信息
 * @description 20251002 陈凤庆 获取lastPID对应的窗口信息
 */
func (a *App) GetLastWindowInfo() map[string]interface{} {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		windowInfo := ks.GetLastWindowInfo()
		if windowInfo == nil {
			return map[string]interface{}{
				"pid":       0,
				"title":     "",
				"timestamp": 0,
			}
		}

		return map[string]interface{}{
			"pid":       windowInfo.PID,
			"title":     windowInfo.Title,
			"timestamp": windowInfo.Timestamp,
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
 */
func (a *App) SwitchToLastWindow() bool {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		return ks.SwitchToLastWindow()
	}
	return false
}

/**
 * getInputMethodName 获取输入方式名称
 * @param inputMethod 输入方式编号
 * @return string 输入方式名称
 * @description 20251005 陈凤庆 用于日志显示的输入方式名称映射
 */
func (a *App) getInputMethodName(inputMethod int) string {
	switch inputMethod {
	case 1:
		return "默认方式(Unicode)"
	case 2:
		return "模拟键盘输入(robotgo.KeyTap)"
	case 3:
		return "复制粘贴输入(robotgo.PasteStr)"
	case 4:
		return "键盘助手输入"
	case 5:
		return "远程输入(ToDesk专用)"
	default:
		return "未知方式"
	}
}

/**
 * simulateTextByMethod 根据输入方式模拟输入文本
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251003 陈凤庆 根据不同的输入方式选择对应的输入方法
 * @modify 20251005 陈凤庆 添加第4种输入方式：键盘助手输入（删除原第4种底层键盘API）
 * @modify 20251005 陈凤庆 添加第5种输入方式：远程输入（专为ToDesk等远程桌面环境设计）
 * @modify 20251005 陈凤庆 添加详细日志显示，便于调试和跟踪
 */
func (a *App) simulateTextByMethod(inputMethod int, text string) error {
	// 获取输入方式名称
	methodName := a.getInputMethodName(inputMethod)
	fmt.Printf("[输入方式] ========== 开始文本输入 ==========\n")
	fmt.Printf("[输入方式] 方式: %s (方式%d)\n", methodName, inputMethod)
	fmt.Printf("[输入方式] 函数: simulateTextByMethod\n")
	fmt.Printf("[输入方式] 内容: %s (长度: %d)\n", text, len(text))

	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		switch inputMethod {
		case 1:
			// 默认方式：使用Unicode发送按键码方式，速度快
			fmt.Printf("[输入方式] 调用函数: ks.SimulateTextWithAutoSwitch\n")
			return ks.SimulateTextWithAutoSwitch(text)
		case 2:
			// 模拟键盘输入：使用robotgo.KeyTap逐键输入
			fmt.Printf("[输入方式] 调用函数: a.simulateTextWithRobotgoKeyTap\n")
			return a.simulateTextWithRobotgoKeyTap(text)
		case 3:
			// 复制粘贴输入：使用robotgo.PasteStr
			fmt.Printf("[输入方式] 调用函数: a.simulateTextWithRobotgoPaste\n")
			return a.simulateTextWithRobotgoPaste(text)
		case 4:
			// 键盘助手输入：使用系统键盘助手
			fmt.Printf("[输入方式] 调用函数: a.keyboardHelperService.SimulateText\n")
			if a.keyboardHelperService != nil {
				return a.keyboardHelperService.SimulateText(text)
			}
			return fmt.Errorf("键盘助手服务未初始化")
		case 5:
			// 远程输入：专为ToDesk等远程桌面环境设计，逐字符输入
			fmt.Printf("[输入方式] 调用函数: a.remoteInputService.SimulateText\n")
			if a.remoteInputService != nil {
				return a.remoteInputService.SimulateText(text)
			}
			return fmt.Errorf("远程输入服务未初始化")
		default:
			// 未知方式，使用默认方式
			fmt.Printf("[输入方式] ⚠️ 未知输入方式，使用默认方式\n")
			fmt.Printf("[输入方式] 调用函数: ks.SimulateTextWithAutoSwitch (默认)\n")
			return ks.SimulateTextWithAutoSwitch(text)
		}
	}
	return nil
}

/**
 * simulateUsernameAndPasswordByMethod 根据输入方式模拟输入用户名和密码
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入(robotgo.KeyTap)、3-复制粘贴输入(robotgo.PasteStr)、4-键盘助手输入、5-远程输入
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251003 陈凤庆 根据不同的输入方式选择对应的输入方法
 * @modify 20251005 陈凤庆 添加第4种输入方式：键盘助手输入（删除原第4种底层键盘API）
 * @modify 20251005 陈凤庆 添加第5种输入方式：远程输入（专为ToDesk等远程桌面环境设计）
 */
func (a *App) simulateUsernameAndPasswordByMethod(inputMethod int, username, password string) error {
	// 获取输入方式名称
	methodName := a.getInputMethodName(inputMethod)
	fmt.Printf("[输入方式] ========== 开始用户名密码输入 ==========\n")
	fmt.Printf("[输入方式] 方式: %s (方式%d)\n", methodName, inputMethod)
	fmt.Printf("[输入方式] 函数: simulateUsernameAndPasswordByMethod\n")
	fmt.Printf("[输入方式] 用户名: %s (长度: %d)\n", username, len(username))
	fmt.Printf("[输入方式] 密码: *** (长度: %d)\n", len(password))

	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		switch inputMethod {
		case 1:
			// 默认方式：使用Unicode发送按键码方式，速度快
			fmt.Printf("[输入方式] 调用函数: ks.SimulateUsernameAndPasswordWithAutoSwitch\n")
			return ks.SimulateUsernameAndPasswordWithAutoSwitch(username, password)
		case 2:
			// 模拟键盘输入：使用robotgo.KeyTap逐键输入
			fmt.Printf("[输入方式] 调用函数: a.simulateUsernameAndPasswordWithRobotgoKeyTap\n")
			return a.simulateUsernameAndPasswordWithRobotgoKeyTap(username, password)
		case 3:
			// 复制粘贴输入：使用robotgo.PasteStr
			fmt.Printf("[输入方式] 调用函数: a.simulateUsernameAndPasswordWithRobotgoPaste\n")
			return a.simulateUsernameAndPasswordWithRobotgoPaste(username, password)
		case 4:
			// 键盘助手输入：使用系统键盘助手
			fmt.Printf("[输入方式] 调用函数: a.keyboardHelperService.SimulateUsernameAndPassword\n")
			if a.keyboardHelperService != nil {
				return a.keyboardHelperService.SimulateUsernameAndPassword(username, password)
			}
			return fmt.Errorf("键盘助手服务未初始化")
		case 5:
			// 远程输入：专为ToDesk等远程桌面环境设计，逐字符输入
			fmt.Printf("[输入方式] 调用函数: a.remoteInputService.SimulateUsernameAndPassword\n")
			if a.remoteInputService != nil {
				return a.remoteInputService.SimulateUsernameAndPassword(username, password)
			}
			return fmt.Errorf("远程输入服务未初始化")
		default:
			// 未知方式，使用默认方式
			fmt.Printf("[输入方式] ⚠️ 未知输入方式，使用默认方式\n")
			fmt.Printf("[输入方式] 调用函数: ks.SimulateUsernameAndPasswordWithAutoSwitch (默认)\n")
			return ks.SimulateUsernameAndPasswordWithAutoSwitch(username, password)
		}
	}
	return nil
}

/**
 * simulateTextWithRobotgoKeyTap 使用robotgo.KeyTap逐键输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo逐键输入，兼容性好但速度较慢
 * @modify 20251004 陈凤庆 添加输入验证和错误处理，防止输入错误字符
 */
func (a *App) simulateTextWithRobotgoKeyTap(text string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		// 验证输入文本
		if text == "" {
			return errors.New("输入文本不能为空")
		}

		// 验证文本内容，确保不是损坏的数据
		if len(text) < 2 && text != "" {
			fmt.Printf("[robotgo输入] ⚠️ 检测到可疑的短文本: '%s'，长度: %d\n", text, len(text))
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
				fmt.Printf("[robotgo输入] ❌ 检测到异常重复字符: '%s'，可能是数据损坏\n", text)
				return fmt.Errorf("检测到异常重复字符，可能是数据损坏: %s", text)
			}
		}

		fmt.Printf("[robotgo输入] 开始逐键输入文本: %s (长度: %d)\n", text, len(text))

		// 切换到目标窗口
		if !ks.SwitchToLastWindow() {
			return errors.New("无法切换到目标窗口")
		}

		// 等待窗口切换完成
		time.Sleep(300 * time.Millisecond)

		// 使用defer确保异常时也能切换回密码管理器
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[robotgo输入] ❌ 发生异常: %v\n", r)
				// 尝试切换回密码管理器
				if !ks.SwitchToPasswordManager() {
					fmt.Printf("[robotgo输入] ⚠️ 异常恢复时切换回密码管理器失败\n")
				}
			}
		}()

		// 输入文本，添加异常处理
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[robotgo输入] ❌ TypeStr操作发生异常: %v\n", r)
					// 不再重新抛出异常，避免程序崩溃
				}
			}()

			fmt.Printf("[robotgo输入] 准备输入文本: %s\n", text)
			safeRobotgoTypeStr(text)
		}()

		fmt.Printf("[robotgo输入] ✅ 文本输入完成\n")

		// 切换回密码管理器
		if !ks.SwitchToPasswordManager() {
			fmt.Printf("[robotgo输入] ⚠️ 切换回密码管理器失败\n")
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活密码管理器窗口为活动窗口
 */
func (a *App) SwitchToPasswordManager() bool {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		return ks.SwitchToPasswordManager()
	}
	return false
}

/**
 * simulateTextWithRobotgoPaste 使用robotgo.PasteStr复制粘贴输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo复制粘贴方式输入，速度快且兼容性好
 * @modify 20251004 陈凤庆 添加错误处理和异常恢复机制，防止程序崩溃
 */
func (a *App) simulateTextWithRobotgoPaste(text string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		// 验证输入文本
		if text == "" {
			return errors.New("输入文本不能为空")
		}

		fmt.Printf("[robotgo粘贴] 开始复制粘贴输入文本: %s (长度: %d)\n", text, len(text))

		// 切换到目标窗口
		if !ks.SwitchToLastWindow() {
			return errors.New("无法切换到目标窗口")
		}

		// 20251005 陈凤庆 等待窗口切换完成 - 增加等待时间确保目标窗口完全激活
		fmt.Printf("[robotgo粘贴] 等待窗口切换完成 (500ms)\n")
		time.Sleep(500 * time.Millisecond)

		// 使用defer确保异常时也能切换回密码管理器
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[robotgo粘贴] ❌ 发生异常: %v\n", r)
				// 尝试切换回密码管理器
				if !ks.SwitchToPasswordManager() {
					fmt.Printf("[robotgo粘贴] ⚠️ 异常恢复时切换回密码管理器失败\n")
				}
			}
		}()

		// 使用纯AppleScript粘贴方案，完全避免robotgo键盘模拟
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[AppleScript粘贴] ❌ 粘贴操作发生异常: %v\n", r)
				}
			}()

			// 将文本复制到剪贴板
			err := safeRobotgoWriteAll(text)
			if err != nil {
				fmt.Printf("[AppleScript粘贴] ⚠️ 写入剪贴板失败: %v\n", err)
				return
			}

			fmt.Printf("[AppleScript粘贴] ✅ 文本已复制到剪贴板: %s\n", text)

			// 启动10秒后自动销毁剪贴板内容的goroutine
			go func() {
				time.Sleep(10 * time.Second)
				currentClipboard, _ := safeRobotgoReadAll()
				if currentClipboard == text {
					// 如果剪贴板内容仍然是我们设置的文本，则清空它
					safeRobotgoWriteAll("")
					fmt.Printf("[AppleScript粘贴] 🔒 剪贴板内容已自动销毁（10秒后）\n")
				}
			}()

			// 20251005 陈凤庆 等待更长时间确保窗口切换完成
			fmt.Printf("[AppleScript粘贴] 等待窗口完全激活 (400ms)\n")
			time.Sleep(400 * time.Millisecond)

			// 使用AppleScript执行Cmd+V粘贴
			err = executeAppleScriptPaste()
			if err != nil {
				fmt.Printf("[AppleScript粘贴] ⚠️ AppleScript粘贴失败: %v，请手动按Cmd+V粘贴\n", err)
				// 给用户时间手动粘贴
				time.Sleep(2000 * time.Millisecond)
			} else {
				fmt.Printf("[AppleScript粘贴] ✅ AppleScript自动粘贴成功\n")
			}

			// 20251005 陈凤庆 等待粘贴操作完成 - 增加等待时间确保目标窗口完全处理
			fmt.Printf("[AppleScript粘贴] 等待粘贴操作完成 (500ms)\n")
			time.Sleep(500 * time.Millisecond)
		}()

		fmt.Printf("[robotgo粘贴] ✅ 文本输入完成\n")

		// 切换回密码管理器
		if !ks.SwitchToPasswordManager() {
			fmt.Printf("[robotgo粘贴] ⚠️ 切换回密码管理器失败\n")
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateUsernameAndPasswordWithRobotgoKeyTap 使用robotgo.KeyTap逐键输入用户名和密码
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251003 陈凤庆 使用robotgo逐键输入用户名和密码，兼容性好但速度较慢
 * @modify 20251005 陈凤庆 改为逐键输入，不使用剪切板
 */
func (a *App) simulateUsernameAndPasswordWithRobotgoKeyTap(username, password string) error {
	if ks, ok := a.keyboardService.(*services.KeyboardService); ok && ks != nil {
		// 切换到目标窗口
		if !ks.SwitchToLastWindow() {
			return errors.New("无法切换到目标窗口")
		}

		// 20251005 陈凤庆 等待窗口切换完成 - 增加等待时间确保目标窗口完全激活
		fmt.Printf("[robotgo输入] 等待窗口切换完成 (500ms)\n")
		time.Sleep(500 * time.Millisecond)

		// 输入用户名（逐键输入，不使用剪切板）
		if username != "" {
			fmt.Printf("[robotgo输入] 开始逐键输入用户名: %s (长度: %d)\n", username, len(username))
			safeRobotgoTypeStr(username)
			fmt.Printf("[robotgo输入] ✅ 用户名输入完成\n")
			// 20251005 陈凤庆 用户名输入完成后等待更长时间
			fmt.Printf("[robotgo输入] 等待用户名输入处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 按Tab键切换到密码输入框
		fmt.Printf("[robotgo输入] 按Tab键切换到密码输入框\n")
		err := executeAppleScriptTab()
		if err != nil {
			fmt.Printf("[AppleScript输入] ⚠️ AppleScript Tab键失败: %v\n", err)
		}
		// 20251005 陈凤庆 Tab键后等待更长时间，确保焦点切换完成
		fmt.Printf("[robotgo输入] 等待Tab键焦点切换完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)

		// 输入密码（逐键输入，不使用剪切板）
		if password != "" {
			fmt.Printf("[robotgo输入] 开始逐键输入密码 (长度: %d)\n", len(password))
			safeRobotgoTypeStr(password)
			fmt.Printf("[robotgo输入] ✅ 密码输入完成\n")
			// 20251005 陈凤庆 密码输入完成后等待更长时间
			fmt.Printf("[robotgo输入] 等待密码输入处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 20251005 陈凤庆 全部输入完成后等待一段时间再切换回密码管理器
		fmt.Printf("[robotgo输入] 全部输入完成，等待目标窗口完全处理 (600ms)\n")
		time.Sleep(600 * time.Millisecond)

		// 切换回密码管理器
		if !ks.SwitchToPasswordManager() {
			fmt.Printf("[robotgo输入] ⚠️ 切换回密码管理器失败\n")
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateUsernameAndPasswordWithRobotgoPaste 使用robotgo.WriteAll复制粘贴输入用户名和密码
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
		// 20251005 陈凤庆 等待窗口切换完成 - 增加等待时间确保目标窗口完全激活
		fmt.Printf("[复制粘贴组合] 等待窗口切换完成 (500ms)\n")
		time.Sleep(500 * time.Millisecond)

		// 步骤4：在目标窗口粘贴用户名
		if username != "" {
			fmt.Printf("[复制粘贴组合] 步骤4：在目标窗口粘贴用户名\n")
			err := executeAppleScriptPaste()
			if err != nil {
				fmt.Printf("[复制粘贴组合] ⚠️ 粘贴用户名失败: %v\n", err)
				// 尝试切换回密码管理器
				ks.SwitchToPasswordManager()
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 用户名粘贴完成\n")
			// 20251005 陈凤庆 用户名粘贴完成后等待更长时间
			fmt.Printf("[复制粘贴组合] 等待用户名粘贴处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 步骤5：在目标输入窗口发送Tab键
		fmt.Printf("[复制粘贴组合] 步骤5：在目标窗口发送Tab键\n")
		sendTabKeyWithCGEventWrapper()
		// 20251005 陈凤庆 Tab键后等待更长时间，确保焦点切换完成
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
			// 20251005 陈凤庆 粘贴密码前等待更长时间，确保剪贴板内容已更新
			fmt.Printf("[复制粘贴组合] 等待剪贴板内容更新 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
			err = executeAppleScriptPaste()
			if err != nil {
				fmt.Printf("[复制粘贴组合] ⚠️ 粘贴密码失败: %v\n", err)
				// 尝试切换回密码管理器
				ks.SwitchToPasswordManager()
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 密码粘贴完成\n")
			// 20251005 陈凤庆 密码粘贴完成后等待更长时间
			fmt.Printf("[复制粘贴组合] 等待密码粘贴处理完成 (400ms)\n")
			time.Sleep(400 * time.Millisecond)
		}

		// 20251005 陈凤庆 全部输入完成后等待一段时间再切换回密码管理器
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
