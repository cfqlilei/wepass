package services

import (
	"errors"
	"fmt"
	"runtime"
	"time"
	"wepassword/internal/logger"
)

/**
 * RemoteInputService 远程输入服务
 * @author 陈凤庆
 * @description 专为ToDesk等远程桌面环境设计的输入服务，支持逐字符输入和窗口切换
 * @date 20251005
 */
type RemoteInputService struct {
	windowMonitor *WindowMonitorService // 窗口监控服务
}

/**
 * NewRemoteInputService 创建远程输入服务实例
 * @param windowMonitor 窗口监控服务
 * @return *RemoteInputService 远程输入服务实例
 */
func NewRemoteInputService(windowMonitor *WindowMonitorService) *RemoteInputService {
	return &RemoteInputService{
		windowMonitor: windowMonitor,
	}
}

/**
 * SimulateText 远程输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 逐字符输入，每个字符都进行窗口切换，确保ToDesk能正确捕获
 */
func (ris *RemoteInputService) SimulateText(text string) error {
	if ris.windowMonitor == nil {
		logger.Error("远程输入: 窗口监控服务未初始化")
		return errors.New("窗口监控服务未初始化")
	}

	if text == "" {
		logger.Error("远程输入: 输入文本为空")
		return errors.New("输入文本不能为空")
	}

	logger.Info("远程输入: ========== 开始远程输入文本 ==========")
	logger.Info("远程输入: 服务: RemoteInputService")
	logger.Info("远程输入: 函数: SimulateText")
	logger.Info("远程输入: 输入内容长度: %d", len(text))
	logger.Info("远程输入: 输入策略: 逐字符输入，每字符都切换窗口")
	logger.Info("远程输入: 底层API: 硬件键码 + CGEventPost (macOS)")
	logger.Info("远程输入: 延迟优化: 50ms字符间隔 (ToDesk专用)")

	// 获取目标窗口信息
	lastWindow := ris.windowMonitor.GetLastWindow()
	if lastWindow == nil {
		logger.Error("远程输入: 未找到目标窗口")
		return errors.New("未找到目标窗口")
	}

	logger.Info("远程输入: 目标窗口: %s (PID: %d)", lastWindow.Title, lastWindow.PID)

	// 一次性切换到目标窗口，保持活动状态
	logger.Info("远程输入: 切换到目标窗口并保持活动状态")
	if !ris.windowMonitor.SwitchToLastWindow() {
		logger.Error("远程输入: 切换到目标窗口失败")
		return errors.New("切换到目标窗口失败")
	}

	// 等待窗口切换完成
	logger.Info("远程输入: 等待窗口切换完成 (200ms)")
	time.Sleep(200 * time.Millisecond)

	// 逐字符输入，不再切换窗口
	runes := []rune(text)
	logger.Info("远程输入: 开始逐字符输入，共 %d 个字符", len(runes))
	for i, char := range runes {
		logger.Info("远程输入: 步骤 %d/%d: 输入字符 '%c' (Unicode: U+%04X)", i+1, len(runes), char, char)

		// 直接输入字符，不切换窗口
		logger.Info("远程输入: 调用底层API输入字符 '%c'", char)
		if err := ris.inputSingleCharacter(char); err != nil {
			logger.Error("远程输入: 底层API输入失败: %v", err)
			return fmt.Errorf("输入字符 '%c' 失败: %w", char, err)
		}

		// 等待字符输入完成 - 增加延迟确保ToDesk和目标窗口完全处理
		logger.Info("远程输入: 等待字符输入完成 (100ms - ToDesk处理时间)")
		time.Sleep(100 * time.Millisecond)

		logger.Info("远程输入: 字符 '%c' 输入完成", char)
	}

	// 输入完成后等待一段时间再切换回密码管理器
	logger.Info("远程输入: 输入完成，等待目标窗口完全处理 (600ms)...")
	time.Sleep(600 * time.Millisecond) // 等待600ms确保输入被目标窗口完全处理
	logger.Info("远程输入: 切换回密码管理器")
	if !ris.windowMonitor.SwitchToPasswordManager() {
		logger.Error("远程输入: 切换回密码管理器失败")
	}

	logger.Info("远程输入: 远程输入完成，共输入 %d 个字符", len(runes))
	logger.Info("远程输入: ========== 远程输入结束 ==========")
	return nil
}

/**
 * SimulateUsernameAndPassword 远程输入用户名和密码
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 先输入用户名，然后Tab键，再输入密码
 */
func (ris *RemoteInputService) SimulateUsernameAndPassword(username, password string) error {
	if ris.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	fmt.Printf("[远程输入] ========== 开始远程输入用户名和密码 ==========\n")
	fmt.Printf("[远程输入] 服务: RemoteInputService\n")
	fmt.Printf("[远程输入] 函数: SimulateUsernameAndPassword\n")
	fmt.Printf("[远程输入] 用户名: %s (长度: %d)\n", username, len(username))
	fmt.Printf("[远程输入] 密码: *** (长度: %d)\n", len(password))
	fmt.Printf("[远程输入] 输入策略: 用户名 → Tab键 → 密码（整个过程保持在目标窗口）\n")

	// 步骤1: 一次性切换到目标窗口并保持活动状态
	fmt.Printf("[远程输入] 步骤1: 切换到目标窗口并保持活动状态\n")
	if !ris.windowMonitor.SwitchToLastWindow() {
		return errors.New("切换到目标窗口失败")
	}
	// 等待窗口切换完成
	fmt.Printf("[远程输入] 等待窗口切换完成 (200ms)\n")
	time.Sleep(200 * time.Millisecond)

	// 步骤2: 输入用户名（不切换回来）
	if username != "" {
		fmt.Printf("[远程输入] 步骤2: 输入用户名\n")
		if err := ris.simulateTextWithoutSwitching(username); err != nil {
			return fmt.Errorf("输入用户名失败: %w", err)
		}
		// 用户名输入完成后等待更长时间，确保目标窗口完全接收
		fmt.Printf("[远程输入] 用户名输入完成，等待目标窗口处理 (500ms)...\n")
		time.Sleep(500 * time.Millisecond)
	} else {
		fmt.Printf("[远程输入] 步骤2: 跳过用户名输入（为空）\n")
	}

	// 步骤3: 按Tab键切换到密码输入框（不切换回来）
	fmt.Printf("[远程输入] 步骤3: 按Tab键切换到密码输入框\n")
	if err := ris.simulateTabWithoutSwitching(); err != nil {
		return fmt.Errorf("Tab键输入失败: %w", err)
	}
	// Tab键后等待更长时间，确保焦点正确切换到密码输入框
	fmt.Printf("[远程输入] Tab键输入完成，等待焦点切换 (400ms)...\n")
	time.Sleep(400 * time.Millisecond)

	// 步骤4: 输入密码（不切换回来）
	if password != "" {
		fmt.Printf("[远程输入] 步骤4: 输入密码\n")
		if err := ris.simulateTextWithoutSwitching(password); err != nil {
			return fmt.Errorf("输入密码失败: %w", err)
		}
		// 密码输入完成后等待更长时间，确保目标窗口完全接收
		fmt.Printf("[远程输入] 密码输入完成，等待目标窗口处理 (500ms)...\n")
		time.Sleep(500 * time.Millisecond)
	} else {
		fmt.Printf("[远程输入] 步骤4: 跳过密码输入（为空）\n")
	}

	// 步骤5: 全部输入完成后等待一段时间再切换回密码管理器
	fmt.Printf("[远程输入] 步骤5: 全部输入完成，等待目标窗口完全处理 (800ms)...\n")
	time.Sleep(800 * time.Millisecond) // 等待800ms确保所有输入都被目标窗口完全处理
	fmt.Printf("[远程输入] 步骤6: 切换回密码管理器\n")
	if !ris.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[远程输入] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[远程输入] ✅ 用户名和密码输入完成\n")
	return nil
}

/**
 * simulateTextWithoutSwitching 输入文本但不切换窗口
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 用于用户名密码输入，避免中途切换窗口
 */
func (ris *RemoteInputService) simulateTextWithoutSwitching(text string) error {
	if text == "" {
		return errors.New("输入文本不能为空")
	}

	// 逐字符输入，不切换窗口
	runes := []rune(text)
	fmt.Printf("[远程输入] 开始逐字符输入，共 %d 个字符\n", len(runes))
	for i, char := range runes {
		fmt.Printf("[远程输入] 输入字符 %d/%d: '%c'\n", i+1, len(runes), char)

		// 使用底层系统API输入字符
		if err := ris.inputCharacterMacOS(char); err != nil {
			fmt.Printf("[远程输入] ❌ 字符 '%c' 输入失败: %v\n", char, err)
			return fmt.Errorf("输入字符 '%c' 失败: %w", char, err)
		}

		// 等待字符输入完成 - 增加延迟确保ToDesk和目标窗口完全处理
		fmt.Printf("[远程输入] 等待字符输入完成 (100ms - ToDesk处理时间)\n")
		time.Sleep(100 * time.Millisecond)

		fmt.Printf("[远程输入] ✅ 字符 '%c' 输入完成\n", char)
	}

	fmt.Printf("[远程输入] ✅ 文本输入完成\n")
	return nil
}

/**
 * simulateTabWithoutSwitching 输入Tab键但不切换窗口
 * @return error 错误信息
 * @description 20251005 陈凤庆 用于用户名密码输入，避免中途切换窗口
 */
func (ris *RemoteInputService) simulateTabWithoutSwitching() error {
	// 使用底层系统API输入Tab键
	if err := ris.inputTabKey(); err != nil {
		return fmt.Errorf("Tab键输入失败: %w", err)
	}
	fmt.Printf("[远程输入] 使用底层系统API输入Tab键\n")

	return nil
}

/**
 * SimulateTab 远程输入Tab键
 * @return error 错误信息
 * @description 使用robotgo.KeyTap模拟Tab键，适合远程桌面环境
 */
func (ris *RemoteInputService) SimulateTab() error {
	if ris.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	fmt.Printf("[远程输入] 开始远程输入Tab键\n")

	// 切换到目标窗口
	if !ris.windowMonitor.SwitchToLastWindow() {
		return errors.New("切换到目标窗口失败")
	}

	// 等待窗口切换完成
	time.Sleep(100 * time.Millisecond)

	// 输入Tab键，使用底层系统API
	if err := ris.inputTabKey(); err != nil {
		return fmt.Errorf("Tab键输入失败: %w", err)
	}
	fmt.Printf("[远程输入] 使用底层系统API输入Tab键\n")

	// 等待Tab键输入完成
	fmt.Printf("[远程输入] 等待Tab键处理完成 (200ms)...\n")
	time.Sleep(200 * time.Millisecond)

	// 切换回密码管理器
	if !ris.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[远程输入] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[远程输入] ✅ Tab键输入完成\n")
	return nil
}

/**
 * inputSingleCharacter 输入单个字符
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 使用底层系统API输入单个字符，确保ToDesk能捕获键盘事件
 */
func (ris *RemoteInputService) inputSingleCharacter(char rune) error {
	// 根据操作系统选择不同的底层API
	switch runtime.GOOS {
	case "darwin":
		// macOS: 使用CGEventPost API
		return ris.inputCharacterMacOS(char)
	case "windows":
		// Windows: 使用SendInput API (暂未实现)
		logger.Error("远程输入-Windows: Windows远程输入暂未实现")
		return fmt.Errorf("Windows远程输入暂未实现")
	case "linux":
		// Linux: 暂时使用robotgo作为备用方案
		logger.Error("远程输入-Linux: 使用robotgo.TypeStr输入字符: '%c'，但Linux远程输入暂未实现", char)
		// TODO: 实现xdotool
		return fmt.Errorf("Linux远程输入暂未实现")
	default:
		logger.Error("远程输入: 不支持的操作系统: %s", runtime.GOOS)
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

/**
 * inputTabKey 输入Tab键
 * @return error 错误信息
 * @description 根据操作系统使用不同的底层API输入Tab键
 */
func (ris *RemoteInputService) inputTabKey() error {
	switch runtime.GOOS {
	case "darwin":
		return ris.inputTabKeyMacOS()
	case "windows":
		// Windows: 使用SendInput API (暂未实现)
		logger.Error("远程输入-Windows: Windows Tab键输入暂未实现")
		return fmt.Errorf("Windows Tab键输入暂未实现")
	case "linux":
		// Linux: 暂时使用robotgo作为备用方案
		fmt.Printf("[远程输入-Linux] 使用robotgo.KeyTap输入Tab键\n")
		// TODO: 实现xdotool
		return fmt.Errorf("Linux远程输入Tab键暂未实现")
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

/**
 * IsServiceReady 检查服务是否就绪
 * @return bool 服务是否就绪
 */
func (ris *RemoteInputService) IsServiceReady() bool {
	return ris.windowMonitor != nil
}
