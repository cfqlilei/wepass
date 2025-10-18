package services

import (
	"errors"
	"fmt"
	"time"
)

/**
 * KeyboardHelperService 键盘助手输入服务
 * @author 陈凤庆
 * @description 提供跨平台的键盘助手输入功能，通过调用系统键盘助手实现字符输入
 * @date 20251005
 */
type KeyboardHelperService struct {
	windowMonitor *WindowMonitorService  // 窗口监控服务
	helperPID     int                    // 键盘助手进程PID
	targetPID     int                    // 目标输入窗口PID
	platformImpl  KeyboardHelperPlatform // 平台特定实现
}

/**
 * KeyboardHelperPlatform 键盘助手平台接口
 * @description 定义各平台需要实现的键盘助手功能
 */
type KeyboardHelperPlatform interface {
	// LaunchHelper 启动键盘助手程序
	// @return int 键盘助手进程PID
	// @return error 错误信息
	LaunchHelper() (int, error)

	// CloseHelper 关闭键盘助手程序
	// @param helperPID 键盘助手进程PID
	// @return error 错误信息
	CloseHelper(helperPID int) error

	// SimulateCharWithHelper 使用键盘助手输入单个字符
	// @param char 要输入的字符
	// @param helperPID 键盘助手进程PID
	// @return error 错误信息
	SimulateCharWithHelper(char rune, helperPID int) error

	// CheckHelperAvailable 检查键盘助手是否可用
	// @return bool 是否可用
	CheckHelperAvailable() bool

	// GetHelperName 获取键盘助手名称
	// @return string 键盘助手名称
	GetHelperName() string
}

/**
 * NewKeyboardHelperService 创建键盘助手服务实例
 * @param windowMonitor 窗口监控服务
 * @return *KeyboardHelperService 服务实例
 * @return error 错误信息
 */
func NewKeyboardHelperService(windowMonitor *WindowMonitorService) (*KeyboardHelperService, error) {
	// 创建平台特定实现
	platformImpl, err := newKeyboardHelperPlatformImpl()
	if err != nil {
		return nil, fmt.Errorf("创建平台实现失败: %w", err)
	}

	// 检查键盘助手是否可用
	if !platformImpl.CheckHelperAvailable() {
		return nil, fmt.Errorf("键盘助手不可用: %s", platformImpl.GetHelperName())
	}

	service := &KeyboardHelperService{
		windowMonitor: windowMonitor,
		helperPID:     0,
		targetPID:     0,
		platformImpl:  platformImpl,
	}

	fmt.Printf("[键盘助手] 服务初始化成功，使用: %s\n", platformImpl.GetHelperName())
	return service, nil
}

/**
 * SimulateText 使用键盘助手模拟输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 */
func (khs *KeyboardHelperService) SimulateText(text string) error {
	if text == "" {
		return errors.New("输入文本不能为空")
	}

	if khs.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	fmt.Printf("[键盘助手] ========== 开始键盘助手输入 ==========\n")
	fmt.Printf("[键盘助手] 输入内容: %s (长度: %d)\n", text, len(text))

	// 步骤1: 从lastPID获取目标窗口PID
	lastWindow := khs.windowMonitor.GetLastWindow()
	if lastWindow == nil || lastWindow.PID <= 0 {
		return errors.New("无法获取目标窗口PID")
	}
	khs.targetPID = lastWindow.PID
	fmt.Printf("[键盘助手] 步骤1: 获取目标窗口PID: %d (%s)\n", khs.targetPID, lastWindow.Title)

	// 步骤2: 启动键盘助手程序
	fmt.Printf("[键盘助手] 步骤2: 启动键盘助手程序\n")
	helperPID, err := khs.platformImpl.LaunchHelper()
	if err != nil {
		return fmt.Errorf("启动键盘助手失败: %w", err)
	}
	khs.helperPID = helperPID
	fmt.Printf("[键盘助手] ✅ 键盘助手已启动，PID: %d\n", helperPID)

	// 确保函数退出时关闭键盘助手
	defer func() {
		if khs.helperPID > 0 {
			fmt.Printf("[键盘助手] 关闭键盘助手，PID: %d\n", khs.helperPID)
			if err := khs.platformImpl.CloseHelper(khs.helperPID); err != nil {
				fmt.Printf("[键盘助手] ⚠️ 关闭键盘助手失败: %v\n", err)
			}
			khs.helperPID = 0
		}
	}()

	// 等待键盘助手启动完成
	time.Sleep(500 * time.Millisecond)

	// 步骤3: 一次性切换到目标窗口并保持活动状态
	fmt.Printf("[键盘助手] 步骤3: 切换到目标窗口并保持活动状态\n")
	if !khs.windowMonitor.SwitchToLastWindow() {
		return errors.New("切换到目标窗口失败")
	}
	time.Sleep(200 * time.Millisecond)

	// 步骤4: 逐字符输入，不再切换窗口
	fmt.Printf("[键盘助手] 步骤4: 开始逐字符输入\n")
	for i, char := range text {
		// 直接使用键盘助手输入字符，不切换窗口
		fmt.Printf("[键盘助手] 输入字符 %d/%d: '%c'\n", i+1, len(text), char)
		if err := khs.platformImpl.SimulateCharWithHelper(char, khs.helperPID); err != nil {
			return fmt.Errorf("输入字符'%c'失败: %w", char, err)
		}

		// 等待字符输入完成（增加等待时间确保字符被正确接收）
		time.Sleep(100 * time.Millisecond)
	}

	// 输入完成后等待一段时间再切换回密码管理器
	fmt.Printf("[键盘助手] 输入完成，等待目标窗口完全处理...\n")
	time.Sleep(600 * time.Millisecond) // 等待600ms确保输入被目标窗口完全处理
	fmt.Printf("[键盘助手] 切换回密码管理器\n")
	if !khs.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[键盘助手] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[键盘助手] ✅ 文本输入完成\n")
	fmt.Printf("[键盘助手] ========== 键盘助手输入结束 ==========\n")

	return nil
}

/**
 * SimulateUsernameAndPassword 使用键盘助手输入用户名和密码
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251005 陈凤庆 修复Tab键问题：整个输入过程保持在目标窗口，最后统一切换回来
 */
func (khs *KeyboardHelperService) SimulateUsernameAndPassword(username, password string) error {
	fmt.Printf("[键盘助手] ========== 开始输入用户名和密码 ==========\n")

	// 步骤1: 启动键盘助手
	fmt.Printf("[键盘助手] 步骤1: 启动键盘助手\n")
	if err := khs.startHelper(); err != nil {
		return fmt.Errorf("启动键盘助手失败: %w", err)
	}

	// 确保在函数结束时关闭键盘助手
	defer func() {
		if khs.helperPID > 0 {
			fmt.Printf("[键盘助手] 关闭键盘助手 (PID: %d)\n", khs.helperPID)
			if err := khs.platformImpl.CloseHelper(khs.helperPID); err != nil {
				fmt.Printf("[键盘助手] ⚠️ 关闭键盘助手失败: %v\n", err)
			}
			khs.helperPID = 0
		}
	}()

	// 等待键盘助手启动完成
	time.Sleep(500 * time.Millisecond)

	// 步骤2: 一次性切换到目标窗口并保持活动状态
	fmt.Printf("[键盘助手] 步骤2: 切换到目标窗口并保持活动状态\n")
	if !khs.windowMonitor.SwitchToLastWindow() {
		return errors.New("切换到目标窗口失败")
	}
	time.Sleep(200 * time.Millisecond)

	// 步骤3: 输入用户名（不切换回来）
	if username != "" {
		fmt.Printf("[键盘助手] 步骤3: 输入用户名\n")
		if err := khs.simulateTextWithoutSwitching(username); err != nil {
			return fmt.Errorf("输入用户名失败: %w", err)
		}
		// 用户名输入完成后等待更长时间，确保目标窗口完全接收
		fmt.Printf("[键盘助手] 用户名输入完成，等待目标窗口处理...\n")
		time.Sleep(500 * time.Millisecond)
	}

	// 步骤4: 模拟Tab键切换到密码输入框（不切换回来）
	fmt.Printf("[键盘助手] 步骤4: 模拟Tab键\n")
	if err := khs.SimulateTab(); err != nil {
		return fmt.Errorf("模拟Tab键失败: %w", err)
	}
	// Tab键后等待更长时间，确保焦点正确切换到密码输入框
	fmt.Printf("[键盘助手] Tab键输入完成，等待焦点切换...\n")
	time.Sleep(400 * time.Millisecond)

	// 步骤5: 输入密码（不切换回来）
	if password != "" {
		fmt.Printf("[键盘助手] 步骤5: 输入密码\n")
		if err := khs.simulateTextWithoutSwitching(password); err != nil {
			return fmt.Errorf("输入密码失败: %w", err)
		}
		// 密码输入完成后等待更长时间，确保目标窗口完全接收
		fmt.Printf("[键盘助手] 密码输入完成，等待目标窗口处理...\n")
		time.Sleep(500 * time.Millisecond)
	}

	// 步骤6: 全部输入完成后等待一段时间再切换回密码管理器
	fmt.Printf("[键盘助手] 步骤6: 全部输入完成，等待目标窗口完全处理...\n")
	time.Sleep(800 * time.Millisecond) // 等待800ms确保所有输入都被目标窗口完全处理
	fmt.Printf("[键盘助手] 步骤7: 切换回密码管理器\n")
	if !khs.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[键盘助手] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[键盘助手] ========== 用户名和密码输入完成 ==========\n")
	return nil
}

/**
 * startHelper 启动键盘助手程序
 * @return error 错误信息
 * @description 20251005 陈凤庆 提取启动键盘助手的公共逻辑
 */
func (khs *KeyboardHelperService) startHelper() error {
	if khs.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	// 从lastPID获取目标窗口PID
	lastWindow := khs.windowMonitor.GetLastWindow()
	if lastWindow == nil || lastWindow.PID <= 0 {
		return errors.New("无法获取目标窗口PID")
	}
	khs.targetPID = lastWindow.PID
	fmt.Printf("[键盘助手] 获取目标窗口PID: %d (%s)\n", khs.targetPID, lastWindow.Title)

	// 启动键盘助手程序
	helperPID, err := khs.platformImpl.LaunchHelper()
	if err != nil {
		return fmt.Errorf("启动键盘助手失败: %w", err)
	}
	khs.helperPID = helperPID
	fmt.Printf("[键盘助手] ✅ 键盘助手已启动，PID: %d\n", helperPID)

	return nil
}

/**
 * simulateTextWithoutSwitching 输入文本但不切换窗口
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251005 陈凤庆 用于用户名密码输入，避免中途切换窗口
 */
func (khs *KeyboardHelperService) simulateTextWithoutSwitching(text string) error {
	if text == "" {
		return errors.New("输入文本不能为空")
	}

	if khs.helperPID <= 0 {
		return errors.New("键盘助手未启动")
	}

	fmt.Printf("[键盘助手] 开始逐字符输入: %s (长度: %d)\n", text, len(text))
	for i, char := range text {
		// 直接使用键盘助手输入字符，不切换窗口
		fmt.Printf("[键盘助手] 输入字符 %d/%d: '%c'\n", i+1, len(text), char)
		if err := khs.platformImpl.SimulateCharWithHelper(char, khs.helperPID); err != nil {
			return fmt.Errorf("输入字符'%c'失败: %w", char, err)
		}

		// 等待字符输入完成（增加等待时间确保字符被正确接收）
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("[键盘助手] ✅ 文本输入完成\n")
	return nil
}

/**
 * SimulateTab 模拟Tab键
 * @return error 错误信息
 */
func (khs *KeyboardHelperService) SimulateTab() error {
	if khs.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	// 直接使用键盘助手输入Tab字符，不切换窗口（假设目标窗口已经是活动窗口）
	if err := khs.platformImpl.SimulateCharWithHelper('\t', khs.helperPID); err != nil {
		return fmt.Errorf("输入Tab键失败: %w", err)
	}

	return nil
}

/**
 * IsHelperRunning 检查键盘助手是否正在运行
 * @return bool 是否运行中
 */
func (khs *KeyboardHelperService) IsHelperRunning() bool {
	return khs.helperPID > 0
}

/**
 * GetHelperPID 获取键盘助手进程PID
 * @return int 进程PID
 */
func (khs *KeyboardHelperService) GetHelperPID() int {
	return khs.helperPID
}

/**
 * GetTargetPID 获取目标输入窗口PID
 * @return int 窗口PID
 */
func (khs *KeyboardHelperService) GetTargetPID() int {
	return khs.targetPID
}

/**
 * Cleanup 清理资源
 * @description 关闭键盘助手，释放资源
 */
func (khs *KeyboardHelperService) Cleanup() {
	if khs.helperPID > 0 {
		fmt.Printf("[键盘助手] 清理资源，关闭键盘助手 PID: %d\n", khs.helperPID)
		if err := khs.platformImpl.CloseHelper(khs.helperPID); err != nil {
			fmt.Printf("[键盘助手] ⚠️ 关闭键盘助手失败: %v\n", err)
		}
		khs.helperPID = 0
	}
	khs.targetPID = 0
}
