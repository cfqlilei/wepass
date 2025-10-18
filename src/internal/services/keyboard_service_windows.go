//go:build windows
// +build windows

package services

/*
#cgo CFLAGS: -I. -DUNICODE -D_UNICODE
#cgo LDFLAGS: -luser32 -lkernel32

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 20251003 陈凤庆 移除了输入法切换相关代码，因为使用Unicode方式输入不依赖输入法

// 20251002 陈凤庆 使用Unicode方式模拟键盘输入文本（不依赖输入法）
void simulateTextInputUnicode(const WCHAR* text) {
    int len = wcslen(text);
    INPUT* inputs = (INPUT*)malloc(len * 2 * sizeof(INPUT)); // 每个字符需要按下和释放

    for (int i = 0; i < len; i++) {
        WCHAR wc = text[i];
        int inputIndex = i * 2;

        // 按下字符键（使用KEYEVENTF_UNICODE标志）
        inputs[inputIndex].type = INPUT_KEYBOARD;
        inputs[inputIndex].ki.wVk = 0;  // 使用Unicode时VK必须为0
        inputs[inputIndex].ki.wScan = wc;  // Unicode字符
        inputs[inputIndex].ki.dwFlags = KEYEVENTF_UNICODE;
        inputs[inputIndex].ki.time = 0;
        inputs[inputIndex].ki.dwExtraInfo = 0;

        // 释放字符键
        inputs[inputIndex + 1].type = INPUT_KEYBOARD;
        inputs[inputIndex + 1].ki.wVk = 0;
        inputs[inputIndex + 1].ki.wScan = wc;
        inputs[inputIndex + 1].ki.dwFlags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP;
        inputs[inputIndex + 1].ki.time = 0;
        inputs[inputIndex + 1].ki.dwExtraInfo = 0;

        SendInput(2, &inputs[inputIndex], sizeof(INPUT));
        Sleep(10); // 字符间延迟
    }

    free(inputs);
}

// 20251002 陈凤庆 模拟键盘输入文本（VK方式，保留用于兼容）
void simulateTextInput(const char* text) {
    int len = strlen(text);
    INPUT* inputs = (INPUT*)malloc(len * 2 * sizeof(INPUT)); // 每个字符需要按下和释放

    for (int i = 0; i < len; i++) {
        char c = text[i];
        SHORT vk = VkKeyScanA(c);
        BYTE virtualKey = LOBYTE(vk);
        BYTE shiftState = HIBYTE(vk);

        int inputIndex = i * 2;

        // 如果需要Shift键
        if (shiftState & 1) {
            // 按下Shift
            inputs[inputIndex].type = INPUT_KEYBOARD;
            inputs[inputIndex].ki.wVk = VK_SHIFT;
            inputs[inputIndex].ki.dwFlags = 0;
            inputs[inputIndex].ki.time = 0;
            inputs[inputIndex].ki.dwExtraInfo = 0;

            // 按下字符键
            inputs[inputIndex + 1].type = INPUT_KEYBOARD;
            inputs[inputIndex + 1].ki.wVk = virtualKey;
            inputs[inputIndex + 1].ki.dwFlags = 0;
            inputs[inputIndex + 1].ki.time = 0;
            inputs[inputIndex + 1].ki.dwExtraInfo = 0;

            SendInput(2, &inputs[inputIndex], sizeof(INPUT));
            Sleep(10); // 短暂延迟

            // 释放字符键
            inputs[inputIndex].type = INPUT_KEYBOARD;
            inputs[inputIndex].ki.wVk = virtualKey;
            inputs[inputIndex].ki.dwFlags = KEYEVENTF_KEYUP;
            inputs[inputIndex].ki.time = 0;
            inputs[inputIndex].ki.dwExtraInfo = 0;

            // 释放Shift
            inputs[inputIndex + 1].type = INPUT_KEYBOARD;
            inputs[inputIndex + 1].ki.wVk = VK_SHIFT;
            inputs[inputIndex + 1].ki.dwFlags = KEYEVENTF_KEYUP;
            inputs[inputIndex + 1].ki.time = 0;
            inputs[inputIndex + 1].ki.dwExtraInfo = 0;

            SendInput(2, &inputs[inputIndex], sizeof(INPUT));
        } else {
            // 按下字符键
            inputs[inputIndex].type = INPUT_KEYBOARD;
            inputs[inputIndex].ki.wVk = virtualKey;
            inputs[inputIndex].ki.dwFlags = 0;
            inputs[inputIndex].ki.time = 0;
            inputs[inputIndex].ki.dwExtraInfo = 0;

            // 释放字符键
            inputs[inputIndex + 1].type = INPUT_KEYBOARD;
            inputs[inputIndex + 1].ki.wVk = virtualKey;
            inputs[inputIndex + 1].ki.dwFlags = KEYEVENTF_KEYUP;
            inputs[inputIndex + 1].ki.time = 0;
            inputs[inputIndex + 1].ki.dwExtraInfo = 0;

            SendInput(2, &inputs[inputIndex], sizeof(INPUT));
        }

        Sleep(10); // 字符间延迟
    }

    free(inputs);
}

// 模拟Tab键
void simulateTabKey() {
    INPUT inputs[2];

    // 按下Tab键
    inputs[0].type = INPUT_KEYBOARD;
    inputs[0].ki.wVk = VK_TAB;
    inputs[0].ki.dwFlags = 0;
    inputs[0].ki.time = 0;
    inputs[0].ki.dwExtraInfo = 0;

    // 释放Tab键
    inputs[1].type = INPUT_KEYBOARD;
    inputs[1].ki.wVk = VK_TAB;
    inputs[1].ki.dwFlags = KEYEVENTF_KEYUP;
    inputs[1].ki.time = 0;
    inputs[1].ki.dwExtraInfo = 0;

    SendInput(2, inputs, sizeof(INPUT));
}
*/
import "C"
import (
	"errors"
	"fmt"
	"time"
	"unicode/utf16"
	"unsafe"
)

/**
 * KeyboardService Windows版本键盘输入服务
 * @author 陈凤庆
 * @description Windows平台的键盘输入模拟功能
 * @modify 陈凤庆 实现完整的Windows平台键盘模拟功能
 */
type KeyboardService struct {
	windowMonitor *WindowMonitorService // 窗口监控服务
}

/**
 * NewKeyboardService 创建新的键盘服务实例（Windows版本）
 * @return *KeyboardService 键盘服务实例
 * @description 20251002 陈凤庆 初始化Windows版本键盘服务，集成窗口监控功能
 */
func NewKeyboardService() *KeyboardService {
	fmt.Printf("[键盘服务] Windows版本键盘服务初始化\n")

	// 20251002 陈凤庆 创建窗口监控服务
	windowMonitor := NewWindowMonitorService()

	// 20251002 陈凤庆 启动窗口监控
	if err := windowMonitor.StartMonitoring(); err != nil {
		fmt.Printf("[键盘服务] ⚠️ 启动窗口监控失败: %v\n", err)
	}

	return &KeyboardService{
		windowMonitor: windowMonitor,
	}
}

/**
 * GetWindowMonitor 获取窗口监控服务（Windows版本）
 * @return *WindowMonitorService 窗口监控服务实例
 * @description 20251005 陈凤庆 提供窗口监控服务访问接口
 */
func (ks *KeyboardService) GetWindowMonitor() *WindowMonitorService {
	return ks.windowMonitor
}

/**
 * CheckAccessibilityPermission 检查是否有辅助功能权限（Windows版本）
 * @return bool 是否有权限
 * @description 20251002 陈凤庆 Windows版本检查键盘输入权限
 */
func (ks *KeyboardService) CheckAccessibilityPermission() bool {
	// 20251002 陈凤庆 Windows版本通常不需要特殊权限，但可能需要管理员权限
	return true
}

/**
 * StorePreviousFocusedApp 存储当前聚焦的应用程序（Windows版本）
 * @description 20251002 陈凤庆 通过窗口监控服务记录当前活动窗口
 */
func (ks *KeyboardService) StorePreviousFocusedApp() {
	if ks.windowMonitor != nil {
		ks.windowMonitor.RecordCurrentAsLastPID()
	}
}

/**
 * RestorePreviousFocusedApp 恢复到上一次聚焦的应用程序（Windows版本）
 * @return error 错误信息
 * @description 20251002 陈凤庆 切换到上一次记录的活动窗口
 */
func (ks *KeyboardService) RestorePreviousFocusedApp() error {
	if ks.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	if !ks.windowMonitor.SwitchToLastWindow() {
		return errors.New("无法切换到上一次聚焦的应用程序")
	}

	return nil
}

/**
 * GetPreviousFocusedAppName 获取上一次聚焦应用程序的名称（Windows版本）
 * @return string 应用程序名称
 * @description 20251002 陈凤庆 返回上一次记录的窗口标题
 */
func (ks *KeyboardService) GetPreviousFocusedAppName() string {
	if ks.windowMonitor == nil {
		return "未知应用程序"
	}

	lastWindow := ks.windowMonitor.GetLastWindow()
	if lastWindow != nil {
		return lastWindow.Title
	}

	return "未知应用程序"
}

/**
 * SimulateText 模拟输入文本（Windows版本）
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251002 陈凤庆 使用Windows API模拟键盘输入文本，使用Unicode方式不依赖输入法
 * @modify 20251002 陈凤庆 改用KEYEVENTF_UNICODE方式，完全绕过输入法问题
 */
func (ks *KeyboardService) SimulateText(text string) error {
	if text == "" {
		return errors.New("输入文本不能为空")
	}

	fmt.Printf("[键盘服务] Windows版本模拟输入文本: %s (长度: %d)\n", text, len(text))

	// 20251002 陈凤庆 将Go字符串转换为UTF-16（Windows宽字符）
	utf16Text := utf16.Encode([]rune(text))

	// 20251002 陈凤庆 添加null终止符
	utf16Text = append(utf16Text, 0)

	// 20251002 陈凤庆 转换为C宽字符串并调用C函数
	cWText := (*C.WCHAR)(unsafe.Pointer(&utf16Text[0]))

	fmt.Printf("[键盘服务] 使用Unicode方式输入（不依赖输入法）\n")
	C.simulateTextInputUnicode(cWText)

	// 20251002 陈凤庆 短暂延迟，确保输入完成
	time.Sleep(50 * time.Millisecond)

	fmt.Printf("[键盘服务] ✅ 文本输入完成\n")
	return nil
}

/**
 * SimulateTab 模拟Tab键（Windows版本）
 * @return error 错误信息
 * @description 20251002 陈凤庆 使用Windows API模拟Tab键输入
 */
func (ks *KeyboardService) SimulateTab() error {
	fmt.Printf("[键盘服务] Windows版本模拟Tab键输入\n")

	C.simulateTabKey()

	fmt.Printf("[键盘服务] ✅ Tab键输入完成\n")
	return nil
}

/**
 * SimulateUsernameAndPassword 模拟输入用户名和密码（Windows版本）
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251002 陈凤庆 先恢复到上一次聚焦的应用程序，然后输入用户名，按Tab键，再输入密码
 */
func (ks *KeyboardService) SimulateUsernameAndPassword(username, password string) error {
	if !ks.CheckAccessibilityPermission() {
		return errors.New("需要相应权限才能模拟键盘输入")
	}

	// 20251002 陈凤庆 恢复到上一次聚焦的应用程序，确保输入到正确的目标窗口
	if err := ks.RestorePreviousFocusedApp(); err != nil {
		return errors.New("无法切换到目标应用程序：" + err.Error())
	}

	// 等待窗口切换完成 - 参考Mac版本增加延迟时间
	fmt.Printf("[键盘服务] 等待窗口切换完成 (500ms)\n")
	time.Sleep(500 * time.Millisecond)

	// 输入用户名
	if username != "" {
		if err := ks.SimulateText(username); err != nil {
			return err
		}
		// 参考Mac版本增加用户名输入后的延迟
		fmt.Printf("[键盘服务] 等待用户名输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	}

	// 按Tab键切换到密码输入框
	if err := ks.SimulateTab(); err != nil {
		return err
	}
	// 参考Mac版本增加Tab键后的延迟
	fmt.Printf("[键盘服务] 等待Tab键焦点切换完成 (400ms)\n")
	time.Sleep(400 * time.Millisecond)

	// 输入密码
	if password != "" {
		if err := ks.SimulateText(password); err != nil {
			return err
		}
		// 参考Mac版本增加密码输入后的延迟
		fmt.Printf("[键盘服务] 等待密码输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	}

	// 全部输入完成后等待一段时间 - 参考Mac版本
	fmt.Printf("[键盘服务] 全部输入完成，等待目标窗口完全处理 (300ms)\n")
	time.Sleep(300 * time.Millisecond)

	return nil
}

/**
 * SimulateTextWithAutoSwitch 模拟输入文本（带自动窗口切换）（Windows版本）
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251002 陈凤庆 自动切换到lastPID窗口，输入文本，然后切换回密码管理器
 */
func (ks *KeyboardService) SimulateTextWithAutoSwitch(text string) error {
	fmt.Printf("[键盘服务] ========== 开始自动输入文本 ==========\n")
	fmt.Printf("[键盘服务] 输入内容: %s (长度: %d)\n", text, len(text))

	if !ks.CheckAccessibilityPermission() {
		fmt.Printf("[键盘服务] ❌ 缺少相应权限\n")
		return errors.New("需要相应权限才能模拟键盘输入")
	}

	if text == "" {
		fmt.Printf("[键盘服务] ❌ 输入文本为空\n")
		return errors.New("输入文本不能为空")
	}

	if ks.windowMonitor == nil {
		fmt.Printf("[键盘服务] ❌ 窗口监控服务未初始化\n")
		return errors.New("窗口监控服务未初始化")
	}

	// 20251002 陈凤庆 切换到lastPID窗口
	fmt.Printf("[键盘服务] 步骤1: 切换到目标窗口\n")
	if !ks.windowMonitor.SwitchToLastWindow() {
		fmt.Printf("[键盘服务] ❌ 无法切换到目标窗口\n")
		return errors.New("无法切换到目标窗口")
	}

	// 等待窗口切换完成 - 参考Mac版本增加延迟时间
	fmt.Printf("[键盘服务] 步骤2: 等待窗口切换完成 (500ms)\n")
	time.Sleep(500 * time.Millisecond)

	// 输入文本
	fmt.Printf("[键盘服务] 步骤3: 开始输入文本\n")
	if err := ks.SimulateText(text); err != nil {
		fmt.Printf("[键盘服务] ❌ 文本输入失败: %v\n", err)
		return err
	}
	fmt.Printf("[键盘服务] ✅ 文本输入完成\n")

	// 等待输入完成 - 参考Mac版本增加延迟时间
	fmt.Printf("[键盘服务] 步骤4: 等待输入完成 (400ms)\n")
	time.Sleep(400 * time.Millisecond)

	// 20251002 陈凤庆 切换回密码管理器
	fmt.Printf("[键盘服务] 步骤5: 切换回密码管理器\n")
	if !ks.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[键盘服务] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[键盘服务] ========== 自动输入完成 ==========\n")
	return nil
}

/**
 * SimulateUsernameAndPasswordWithAutoSwitch 模拟输入用户名和密码（带自动窗口切换）（Windows版本）
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 20251002 陈凤庆 自动切换到lastPID窗口，输入用户名和密码，然后切换回密码管理器
 */
func (ks *KeyboardService) SimulateUsernameAndPasswordWithAutoSwitch(username, password string) error {
	fmt.Printf("[键盘服务] ========== 开始自动输入用户名和密码 ==========\n")
	fmt.Printf("[键盘服务] 用户名: %s (长度: %d)\n", username, len(username))
	fmt.Printf("[键盘服务] 密码: %s (长度: %d)\n", "***", len(password))

	if !ks.CheckAccessibilityPermission() {
		fmt.Printf("[键盘服务] ❌ 缺少相应权限\n")
		return errors.New("需要相应权限才能模拟键盘输入")
	}

	if ks.windowMonitor == nil {
		fmt.Printf("[键盘服务] ❌ 窗口监控服务未初始化\n")
		return errors.New("窗口监控服务未初始化")
	}

	// 20251002 陈凤庆 切换到lastPID窗口
	fmt.Printf("[键盘服务] 步骤1: 切换到目标窗口\n")
	if !ks.windowMonitor.SwitchToLastWindow() {
		fmt.Printf("[键盘服务] ❌ 无法切换到目标窗口\n")
		return errors.New("无法切换到目标窗口")
	}

	// 等待窗口切换完成 - 参考Mac版本增加延迟时间
	fmt.Printf("[键盘服务] 步骤2: 等待窗口切换完成 (500ms)\n")
	time.Sleep(500 * time.Millisecond)

	// 输入用户名
	if username != "" {
		fmt.Printf("[键盘服务] 步骤3: 输入用户名\n")
		if err := ks.SimulateText(username); err != nil {
			fmt.Printf("[键盘服务] ❌ 用户名输入失败: %v\n", err)
			return err
		}
		fmt.Printf("[键盘服务] ✅ 用户名输入完成\n")
		// 参考Mac版本增加用户名输入后的延迟
		fmt.Printf("[键盘服务] 等待用户名输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	} else {
		fmt.Printf("[键盘服务] 步骤3: 跳过用户名输入（为空）\n")
	}

	// 按Tab键切换到密码输入框
	fmt.Printf("[键盘服务] 步骤4: 按Tab键切换到密码输入框\n")
	if err := ks.SimulateTab(); err != nil {
		fmt.Printf("[键盘服务] ❌ Tab键输入失败: %v\n", err)
		return err
	}
	fmt.Printf("[键盘服务] ✅ Tab键输入完成\n")
	// 参考Mac版本增加Tab键后的延迟
	fmt.Printf("[键盘服务] 等待Tab键焦点切换完成 (400ms)\n")
	time.Sleep(400 * time.Millisecond)

	// 输入密码
	if password != "" {
		fmt.Printf("[键盘服务] 步骤5: 输入密码\n")
		if err := ks.SimulateText(password); err != nil {
			fmt.Printf("[键盘服务] ❌ 密码输入失败: %v\n", err)
			return err
		}
		fmt.Printf("[键盘服务] ✅ 密码输入完成\n")
		// 参考Mac版本增加密码输入后的延迟
		fmt.Printf("[键盘服务] 等待密码输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	} else {
		fmt.Printf("[键盘服务] 步骤5: 跳过密码输入（为空）\n")
	}

	// 全部输入完成后等待一段时间 - 参考Mac版本
	fmt.Printf("[键盘服务] 全部输入完成，等待目标窗口完全处理 (600ms)\n")
	time.Sleep(600 * time.Millisecond)

	// 20251002 陈凤庆 切换回密码管理器
	fmt.Printf("[键盘服务] 步骤6: 切换回密码管理器\n")
	if !ks.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[键盘服务] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[键盘服务] ========== 自动输入完成 ==========\n")
	return nil
}

/**
 * RecordCurrentWindow 记录当前活动窗口为lastPID（Windows版本）
 * @return error 错误信息
 * @description 20251002 陈凤庆 手动记录当前活动窗口，用于鼠标悬停等场景
 */
func (ks *KeyboardService) RecordCurrentWindow() error {
	if ks.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	if !ks.windowMonitor.RecordCurrentAsLastPID() {
		return errors.New("记录当前窗口失败")
	}

	return nil
}

/**
 * GetLastWindowInfo 获取最后活动窗口信息（Windows版本）
 * @return *WindowInfo 窗口信息
 * @description 20251002 陈凤庆 返回最后记录的窗口信息
 */
func (ks *KeyboardService) GetLastWindowInfo() *WindowInfo {
	if ks.windowMonitor == nil {
		return nil
	}

	return ks.windowMonitor.GetLastWindow()
}

/**
 * StopWindowMonitoring 停止窗口监控（Windows版本）
 * @return error 错误信息
 * @description 20251002 陈凤庆 停止窗口监控服务
 */
func (ks *KeyboardService) StopWindowMonitoring() error {
	if ks.windowMonitor == nil {
		return nil
	}

	return ks.windowMonitor.StopMonitoring()
}

/**
 * SwitchToLastWindow 切换到最后活动的窗口（Windows版本）
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 切换到最后记录的活动窗口
 */
func (ks *KeyboardService) SwitchToLastWindow() bool {
	if ks.windowMonitor == nil {
		return false
	}

	return ks.windowMonitor.SwitchToLastWindow()
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口（Windows版本）
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 切换回密码管理器窗口
 */
func (ks *KeyboardService) SwitchToPasswordManager() bool {
	if ks.windowMonitor == nil {
		return false
	}

	return ks.windowMonitor.SwitchToPasswordManager()
}
