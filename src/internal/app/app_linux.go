//go:build linux
// +build linux

package app

import (
	"errors"
	"fmt"
	"time"
	"wepassword/internal/services"

	"github.com/go-vgo/robotgo"
)

/**
 * safeRobotgoWriteAll 安全的剪贴板写入操作 (Linux版本)
 * @param text 要写入的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 包装robotgo.WriteAll，添加异常处理
 */
func safeRobotgoWriteAll(text string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[robotgo] WriteAll操作发生异常: %v\n", r)
		}
	}()

	return robotgo.WriteAll(text)
}

/**
 * safeRobotgoReadAll 安全的剪贴板读取操作 (Linux版本)
 * @return string 剪贴板内容
 * @return error 错误信息
 * @description 20251005 陈凤庆 包装robotgo.ReadAll，添加异常处理
 */
func safeRobotgoReadAll() (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[robotgo] ReadAll操作发生异常: %v\n", r)
		}
	}()

	return robotgo.ReadAll()
}

/**
 * safeRobotgoPasteStr 安全的粘贴操作 (Linux版本)
 * @param text 要粘贴的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 包装robotgo.PasteStr，添加异常处理
 */
func safeRobotgoPasteStr(text string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[robotgo] PasteStr操作发生异常: %v\n", r)
		}
	}()

	return robotgo.PasteStr(text)
}

/**
 * safeRobotgoTypeStr 安全的逐键输入操作 (Linux版本)
 * @param text 要输入的文本
 * @description 20251005 陈凤庆 包装robotgo.TypeStr，添加异常处理
 */
func safeRobotgoTypeStr(text string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[robotgo] TypeStr操作发生异常: %v\n", r)
		}
	}()

	robotgo.TypeStr(text)
}

/**
 * safeRobotgoKeyTap 安全的按键操作 (Linux版本)
 * @param key 要按的键
 * @description 20251005 陈凤庆 包装robotgo.KeyTap，添加异常处理
 */
func safeRobotgoKeyTap(key string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[robotgo] KeyTap操作发生异常: %v\n", r)
		}
	}()

	robotgo.KeyTap(key)
}

/**
 * initKeyboardService 初始化键盘服务 (Linux版本)
 * @description 20251005 陈凤庆 Linux平台使用robotgo实现基本键盘服务
 */
func (a *App) initKeyboardService() {
	// Linux平台使用基本的键盘服务
	keyboardService := services.NewKeyboardService()
	if keyboardService != nil {
		a.keyboardService = keyboardService
		fmt.Printf("[Linux键盘服务] ✅ 键盘服务初始化成功\n")
	} else {
		a.keyboardService = nil
		fmt.Printf("[Linux键盘服务] ⚠️ 键盘服务初始化失败\n")
	}

	// 初始化键盘助手服务
	keyboardHelperService := services.NewKeyboardHelperService()
	if keyboardHelperService != nil {
		a.keyboardHelperService = keyboardHelperService
		fmt.Printf("[Linux键盘助手] ✅ 键盘助手服务初始化成功\n")
	} else {
		a.keyboardHelperService = nil
		fmt.Printf("[Linux键盘助手] ⚠️ 键盘助手服务初始化失败\n")
	}
}

/**
 * CheckAccessibilityPermission 检查是否有辅助功能权限
 * @return bool 是否有权限
 * @description 20250127 陈凤庆 Linux平台总是返回true
 */
func (a *App) CheckAccessibilityPermission() bool {
	// Linux平台不需要特殊权限
	return true
}

/**
 * StorePreviousFocusedApp 存储当前聚焦的应用程序
 * @description 20250127 陈凤庆 Linux平台暂不支持
 */
func (a *App) StorePreviousFocusedApp() {
	// Linux平台暂不支持
}

/**
 * GetPreviousFocusedAppName 获取上一次聚焦应用程序的名称
 * @return string 应用程序名称
 * @description 20250127 陈凤庆 Linux平台暂不支持
 */
func (a *App) GetPreviousFocusedAppName() string {
	return ""
}

/**
 * SimulateUsernameAndPassword 模拟输入用户名和密码
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 Linux平台暂不支持
 * @modify 20251003 陈凤庆 修改方法签名，只接收accountId参数
 */
func (a *App) SimulateUsernameAndPassword(accountId string) error {
	return nil
}

/**
 * SimulateUsername 模拟输入用户名
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 Linux平台暂不支持
 * @modify 20251003 陈凤庆 修改方法签名，只接收accountId参数
 */
func (a *App) SimulateUsername(accountId string) error {
	return nil
}

/**
 * SimulatePassword 模拟输入密码
 * @param accountId 账号ID
 * @return error 错误信息
 * @description 20250127 陈凤庆 Linux平台暂不支持
 * @modify 20251003 陈凤庆 修改方法签名，只接收accountId参数
 */
func (a *App) SimulatePassword(accountId string) error {
	return nil
}

/**
 * RecordLastWindow 记录当前活动窗口为lastPID
 * @return error 错误信息
 * @description 20251002 陈凤庆 Linux平台暂不支持
 */
func (a *App) RecordLastWindow() error {
	return nil
}

/**
 * GetLastWindowInfo 获取最后活动窗口信息
 * @return map[string]interface{} 窗口信息
 * @description 20251002 陈凤庆 Linux平台暂不支持
 */
func (a *App) GetLastWindowInfo() map[string]interface{} {
	return map[string]interface{}{
		"pid":       0,
		"title":     "",
		"timestamp": 0,
	}
}

/**
 * SwitchToLastWindow 切换到最后活动的窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 Linux平台暂不支持
 */
func (a *App) SwitchToLastWindow() bool {
	return false
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 Linux平台暂不支持
 */
func (a *App) SwitchToPasswordManager() bool {
	return false
}

/**
 * simulateTextByMethod 根据输入方式模拟输入文本 (Linux版本)
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入、3-复制粘贴输入、4-键盘助手输入
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 Linux平台使用robotgo实现输入功能
 */
func (a *App) simulateTextByMethod(inputMethod int, text string) error {
	if a.keyboardService != nil {
		switch inputMethod {
		case 1:
			// 默认方式：使用键盘服务的自动切换输入
			if ks, ok := a.keyboardService.(interface{ SimulateTextWithAutoSwitch(string) error }); ok {
				return ks.SimulateTextWithAutoSwitch(text)
			}
		case 2:
			// 模拟键盘输入：使用robotgo.TypeStr逐键输入
			return a.simulateTextWithRobotgoKeyTap(text)
		case 3:
			// 复制粘贴输入：使用robotgo.PasteStr
			return a.simulateTextWithRobotgoPaste(text)
		case 4:
			// 键盘助手输入：使用系统键盘助手
			if a.keyboardHelperService != nil {
				return a.keyboardHelperService.SimulateText(text)
			}
			return fmt.Errorf("键盘助手服务未初始化")
		default:
			// 未知方式，使用默认方式
			if ks, ok := a.keyboardService.(interface{ SimulateTextWithAutoSwitch(string) error }); ok {
				return ks.SimulateTextWithAutoSwitch(text)
			}
		}
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateUsernameAndPasswordByMethod 根据输入方式模拟输入用户名和密码 (Linux版本)
 * @param inputMethod 输入方式：1-默认方式(Unicode)、2-模拟键盘输入、3-复制粘贴输入、4-键盘助手输入
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251005 陈凤庆 Linux平台使用robotgo实现输入功能
 */
func (a *App) simulateUsernameAndPasswordByMethod(inputMethod int, username, password string) error {
	if a.keyboardService != nil {
		switch inputMethod {
		case 1:
			// 默认方式：使用键盘服务的自动切换输入
			if ks, ok := a.keyboardService.(interface{ SimulateUsernameAndPasswordWithAutoSwitch(string, string) error }); ok {
				return ks.SimulateUsernameAndPasswordWithAutoSwitch(username, password)
			}
		case 2:
			// 模拟键盘输入：使用robotgo.TypeStr逐键输入
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
		default:
			// 未知方式，使用默认方式
			if ks, ok := a.keyboardService.(interface{ SimulateUsernameAndPasswordWithAutoSwitch(string, string) error }); ok {
				return ks.SimulateUsernameAndPasswordWithAutoSwitch(username, password)
			}
		}
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateTextWithRobotgoKeyTap 使用robotgo.TypeStr逐键输入文本 (Linux版本)
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 使用robotgo逐键输入，兼容性好但速度较慢，不使用剪切板
 */
func (a *App) simulateTextWithRobotgoKeyTap(text string) error {
	if a.keyboardService != nil {
		fmt.Printf("[robotgo输入] 开始逐键输入文本: %s (长度: %d)\n", text, len(text))

		// 切换到目标窗口
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			if !ks.SwitchToLastWindow() {
				return errors.New("无法切换到目标窗口")
			}
		}

		time.Sleep(300 * time.Millisecond)

		// 准备输入文本
		fmt.Printf("[robotgo输入] 准备输入文本: %s\n", text)
		safeRobotgoTypeStr(text)
		fmt.Printf("[robotgo输入] ✅ 文本输入完成\n")

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
 * simulateTextWithRobotgoPaste 使用robotgo.PasteStr复制粘贴输入文本 (Linux版本)
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 使用robotgo复制粘贴方式输入文本，速度快且兼容性好
 */
func (a *App) simulateTextWithRobotgoPaste(text string) error {
	if a.keyboardService != nil {
		fmt.Printf("[robotgo粘贴] 开始复制粘贴输入文本: %s (长度: %d)\n", text, len(text))

		// 切换到目标窗口
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			if !ks.SwitchToLastWindow() {
				return errors.New("无法切换到目标窗口")
			}
		}

		time.Sleep(300 * time.Millisecond)

		// 复制到剪贴板
		err := safeRobotgoWriteAll(text)
		if err != nil {
			return err
		}
		fmt.Printf("[robotgo粘贴] ✅ 文本已复制到剪贴板: %s\n", text)

		// 粘贴文本
		err = safeRobotgoPasteStr(text)
		if err != nil {
			return err
		}
		fmt.Printf("[robotgo粘贴] ✅ 文本粘贴成功\n")

		// 自动销毁剪贴板内容
		go func() {
			time.Sleep(10 * time.Second)
			currentClipboard, _ := safeRobotgoReadAll()
			if currentClipboard == text {
				safeRobotgoWriteAll("")
				fmt.Printf("[robotgo粘贴] 🔒 剪贴板内容已自动销毁（10秒后）\n")
			}
		}()

		fmt.Printf("[robotgo粘贴] ✅ 文本输入完成\n")

		// 切换回密码管理器
		if ks, ok := a.keyboardService.(interface{ SwitchToPasswordManager() bool }); ok {
			if !ks.SwitchToPasswordManager() {
				fmt.Printf("[robotgo粘贴] ⚠️ 切换回密码管理器失败\n")
			}
		}

		return nil
	}
	return errors.New("键盘服务未初始化")
}

/**
 * simulateUsernameAndPasswordWithRobotgoKeyTap 使用robotgo.TypeStr逐键输入用户名和密码 (Linux版本)
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251005 陈凤庆 使用robotgo逐键输入用户名和密码，兼容性好但速度较慢，不使用剪切板
 */
func (a *App) simulateUsernameAndPasswordWithRobotgoKeyTap(username, password string) error {
	if a.keyboardService != nil {
		// 切换到目标窗口
		if ks, ok := a.keyboardService.(interface{ SwitchToLastWindow() bool }); ok {
			if !ks.SwitchToLastWindow() {
				return errors.New("无法切换到目标窗口")
			}
		}

		// 等待窗口切换完成
		time.Sleep(300 * time.Millisecond)

		// 输入用户名
		if username != "" {
			fmt.Printf("[robotgo输入] 开始逐键输入用户名: %s\n", username)
			safeRobotgoTypeStr(username)
			fmt.Printf("[robotgo输入] ✅ 用户名输入完成\n")
			time.Sleep(100 * time.Millisecond)
		}

		// 按Tab键切换到密码输入框
		fmt.Printf("[robotgo输入] 按Tab键切换到密码输入框\n")
		safeRobotgoKeyTap("tab")
		time.Sleep(200 * time.Millisecond)

		// 输入密码
		if password != "" {
			fmt.Printf("[robotgo输入] 开始输入密码\n")
			safeRobotgoTypeStr(password)
			fmt.Printf("[robotgo输入] ✅ 密码输入完成\n")
		}

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
 * simulateUsernameAndPasswordWithRobotgoPaste 使用robotgo.PasteStr复制粘贴输入用户名和密码 (Linux版本)
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251005 陈凤庆 使用robotgo复制粘贴方式输入用户名和密码，速度快且兼容性好
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
		time.Sleep(300 * time.Millisecond)

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
			time.Sleep(200 * time.Millisecond)
		}

		// 步骤5：在目标输入窗口发送Tab键
		fmt.Printf("[复制粘贴组合] 步骤5：在目标窗口发送Tab键\n")
		safeRobotgoKeyTap("tab")
		time.Sleep(200 * time.Millisecond)
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
			time.Sleep(200 * time.Millisecond)
			err = safeRobotgoPasteStr(password)
			if err != nil {
				fmt.Printf("[复制粘贴组合] ⚠️ 粘贴密码失败: %v\n", err)
				// 尝试切换回密码管理器
				ks.SwitchToPasswordManager()
				return err
			}
			fmt.Printf("[复制粘贴组合] ✅ 密码粘贴完成\n")
			time.Sleep(200 * time.Millisecond)
		}

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
