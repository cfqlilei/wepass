//go:build darwin
// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework Foundation -framework AppKit
#include <ApplicationServices/ApplicationServices.h>
#include <Foundation/Foundation.h>
#include <AppKit/AppKit.h>
#include <unistd.h>
#include <stdlib.h>

// 全局变量存储上一次聚焦的应用程序
static NSRunningApplication *previousFocusedApp = nil;

// 获取当前聚焦的应用程序
NSRunningApplication* getCurrentFocusedApp() {
    @autoreleasepool {
        return [[NSWorkspace sharedWorkspace] frontmostApplication];
    }
}

// 存储当前聚焦的应用程序为上一次聚焦的应用程序
void storePreviousFocusedApp() {
    @autoreleasepool {
        NSRunningApplication *currentApp = getCurrentFocusedApp();
        if (currentApp) {
            // 安全释放之前的应用程序引用
            if (previousFocusedApp != nil && previousFocusedApp != currentApp) {
                [previousFocusedApp release];
                previousFocusedApp = nil;
            }

            // 只有当前应用程序不同时才更新
            if (previousFocusedApp != currentApp) {
                previousFocusedApp = [currentApp retain];
            }
        }
    }
}

// 恢复到上一次聚焦的应用程序
bool restorePreviousFocusedApp() {
    @autoreleasepool {
        if (previousFocusedApp != nil) {
            // 检查应用程序是否仍然在运行
            if ([previousFocusedApp isTerminated]) {
                // 应用程序已终止，清理引用
                [previousFocusedApp release];
                previousFocusedApp = nil;
                return false;
            }

            // 20250127 陈凤庆 使用activateWithOptions方法强制激活应用程序
            BOOL success = [previousFocusedApp activateWithOptions:NSApplicationActivateIgnoringOtherApps];
            if (success) {
                // 等待应用程序激活完成
                usleep(300000); // 300ms延迟，确保窗口切换完成
                return true;
            }
        }
        return false;
    }
}

// 获取上一次聚焦应用程序的名称（用于调试）
const char* getPreviousFocusedAppName() {
    @autoreleasepool {
        if (previousFocusedApp != nil) {
            // 检查应用程序是否仍然有效
            if ([previousFocusedApp isTerminated]) {
                // 应用程序已终止，清理引用
                [previousFocusedApp release];
                previousFocusedApp = nil;
                return strdup("应用程序已终止");
            }

            NSString *appName = [previousFocusedApp localizedName];
            if (appName != nil) {
                // 20250127 陈凤庆 使用strdup创建字符串副本，避免内存释放问题
                const char *utf8String = [appName UTF8String];
                if (utf8String != NULL) {
                    return strdup(utf8String);
                }
            }
        }
        return strdup("无");
    }
}

// 模拟键盘输入字符串
void simulateKeyboardInput(const char* text) {
    @autoreleasepool {
        NSString *nsText = [NSString stringWithUTF8String:text];

        for (NSUInteger i = 0; i < [nsText length]; i++) {
            unichar character = [nsText characterAtIndex:i];

            // 创建键盘事件
            CGEventRef keyDownEvent = CGEventCreateKeyboardEvent(NULL, 0, true);
            CGEventRef keyUpEvent = CGEventCreateKeyboardEvent(NULL, 0, false);

            // 设置Unicode字符
            CGEventKeyboardSetUnicodeString(keyDownEvent, 1, &character);
            CGEventKeyboardSetUnicodeString(keyUpEvent, 1, &character);

            // 发送事件
            CGEventPost(kCGHIDEventTap, keyDownEvent);
            CGEventPost(kCGHIDEventTap, keyUpEvent);

            // 释放事件
            CFRelease(keyDownEvent);
            CFRelease(keyUpEvent);

            // 短暂延迟，模拟真实输入
            usleep(10000); // 10ms
        }
    }
}

// 模拟Tab键
void simulateTabKey() {
    CGEventRef keyDownEvent = CGEventCreateKeyboardEvent(NULL, 48, true); // Tab键的虚拟键码是48
    CGEventRef keyUpEvent = CGEventCreateKeyboardEvent(NULL, 48, false);

    CGEventPost(kCGHIDEventTap, keyDownEvent);
    CGEventPost(kCGHIDEventTap, keyUpEvent);

    CFRelease(keyDownEvent);
    CFRelease(keyUpEvent);

    usleep(50000); // 50ms延迟
}

// 检查是否有辅助功能权限
bool checkAccessibilityPermission() {
    return AXIsProcessTrusted();
}
*/
import "C"
import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

/**
 * KeyboardService 键盘输入服务
 * @author 陈凤庆
 * @description 提供跨应用程序的键盘输入模拟功能，支持窗口聚焦管理
 * @modify 陈凤庆 集成窗口监控服务，实现自动切换活动窗口功能
 */
type KeyboardService struct {
	windowMonitor *WindowMonitorService // 20251002 陈凤庆 窗口监控服务
}

/**
 * NewKeyboardService 创建新的键盘服务实例
 * @return *KeyboardService 键盘服务实例
 * @modify 陈凤庆 初始化窗口监控服务并启动监控
 */
func NewKeyboardService() *KeyboardService {
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
 * GetWindowMonitor 获取窗口监控服务
 * @return *WindowMonitorService 窗口监控服务实例
 * @description 20251005 陈凤庆 提供窗口监控服务访问接口
 */
func (ks *KeyboardService) GetWindowMonitor() *WindowMonitorService {
	return ks.windowMonitor
}

/**
 * CheckAccessibilityPermission 检查是否有辅助功能权限
 * @return bool 是否有权限
 * @description 在macOS上需要辅助功能权限才能模拟键盘输入
 */
func (ks *KeyboardService) CheckAccessibilityPermission() bool {
	return bool(C.checkAccessibilityPermission())
}

/**
 * StorePreviousFocusedApp 存储当前聚焦的应用程序
 * @description 在执行自动填充前调用，记录当前聚焦的应用程序
 */
func (ks *KeyboardService) StorePreviousFocusedApp() {
	C.storePreviousFocusedApp()
}

/**
 * RestorePreviousFocusedApp 恢复到上一次聚焦的应用程序
 * @return error 错误信息
 * @description 在执行键盘输入前调用，确保输入到正确的目标窗口
 */
func (ks *KeyboardService) RestorePreviousFocusedApp() error {
	success := bool(C.restorePreviousFocusedApp())
	if !success {
		return errors.New("无法恢复到上一次聚焦的应用程序")
	}
	return nil
}

/**
 * GetPreviousFocusedAppName 获取上一次聚焦应用程序的名称
 * @return string 应用程序名称
 * @description 用于调试和日志记录
 * @modify 陈凤庆 修复内存管理问题，确保正确释放C字符串内存
 */
func (ks *KeyboardService) GetPreviousFocusedAppName() string {
	cName := C.getPreviousFocusedAppName()
	defer C.free(unsafe.Pointer(cName)) // 20250127 陈凤庆 确保释放C分配的内存
	return C.GoString(cName)
}

/**
 * SimulateText 模拟输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 支持大小写、数字、特殊字符的输入模拟
 */
func (ks *KeyboardService) SimulateText(text string) error {
	if !ks.CheckAccessibilityPermission() {
		return errors.New("需要辅助功能权限才能模拟键盘输入，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用")
	}

	if text == "" {
		return errors.New("输入文本不能为空")
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	C.simulateKeyboardInput(cText)
	return nil
}

/**
 * SimulateTab 模拟Tab键
 * @return error 错误信息
 * @description 模拟按下Tab键，用于在输入框之间切换
 */
func (ks *KeyboardService) SimulateTab() error {
	if !ks.CheckAccessibilityPermission() {
		return errors.New("需要辅助功能权限才能模拟键盘输入，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用")
	}

	C.simulateTabKey()
	return nil
}

/**
 * SimulateUsernameAndPassword 模拟输入用户名和密码
 * @param username 用户名
 * @param password 密码
 * @return error 错误信息
 * @description 先恢复到上一次聚焦的应用程序，然后输入用户名，按Tab键，再输入密码，适用于跨应用程序自动填充
 */
func (ks *KeyboardService) SimulateUsernameAndPassword(username, password string) error {
	if !ks.CheckAccessibilityPermission() {
		return errors.New("需要辅助功能权限才能模拟键盘输入，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用")
	}

	// 20250127 陈凤庆 恢复到上一次聚焦的应用程序，确保输入到正确的目标窗口
	if err := ks.RestorePreviousFocusedApp(); err != nil {
		return errors.New("无法切换到目标应用程序：" + err.Error())
	}

	// 20251005 陈凤庆 窗口切换后等待更长时间，确保应用程序完全激活
	fmt.Printf("[键盘服务] 等待窗口切换完成 (500ms)\n")
	time.Sleep(500 * time.Millisecond) // 从200ms增加到500ms

	// 输入用户名
	if username != "" {
		if err := ks.SimulateText(username); err != nil {
			return err
		}
		// 20251005 陈凤庆 用户名输入后等待更长时间，确保目标窗口完全处理
		fmt.Printf("[键盘服务] 等待用户名输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond) // 从100ms增加到400ms
	}

	// 按Tab键切换到密码输入框
	if err := ks.SimulateTab(); err != nil {
		return err
	}

	// 20251005 陈凤庆 Tab键后等待更长时间，确保焦点切换完成
	fmt.Printf("[键盘服务] 等待Tab键焦点切换完成 (400ms)\n")
	time.Sleep(400 * time.Millisecond) // 从200ms增加到400ms

	// 输入密码
	if password != "" {
		if err := ks.SimulateText(password); err != nil {
			return err
		}
		// 20251005 陈凤庆 密码输入完成后等待一段时间，确保目标窗口完全处理
		fmt.Printf("[键盘服务] 等待密码输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	}

	// 20251005 陈凤庆 全部输入完成后等待一段时间
	fmt.Printf("[键盘服务] 全部输入完成，等待目标窗口完全处理 (300ms)\n")
	time.Sleep(300 * time.Millisecond)

	return nil
}

/**
 * SimulateTextWithAutoSwitch 模拟输入文本（带自动窗口切换）
 * @param text 要输入的文本
 * @return error 错误信息
 * @description 20251002 陈凤庆 自动切换到lastPID窗口，输入文本，然后切换回密码管理器
 */
func (ks *KeyboardService) SimulateTextWithAutoSwitch(text string) error {
	fmt.Printf("[键盘服务] ========== 开始自动输入文本 ==========\n")
	fmt.Printf("[键盘服务] 输入内容: %s (长度: %d)\n", text, len(text))

	if !ks.CheckAccessibilityPermission() {
		fmt.Printf("[键盘服务] ❌ 缺少辅助功能权限\n")
		return errors.New("需要辅助功能权限才能模拟键盘输入，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用")
	}

	if text == "" {
		fmt.Printf("[键盘服务] ❌ 输入文本为空\n")
		return errors.New("输入文本不能为空")
	}

	if ks.windowMonitor == nil {
		fmt.Printf("[键盘服务] ❌ 窗口监控服务未初始化\n")
		return errors.New("窗口监控服务未初始化")
	}

	// 一次性切换到目标窗口并保持活动状态
	fmt.Printf("[键盘服务] 步骤1: 切换到目标窗口并保持活动状态\n")
	if !ks.windowMonitor.SwitchToLastWindow() {
		fmt.Printf("[键盘服务] ❌ 无法切换到目标窗口\n")
		return errors.New("无法切换到目标窗口")
	}

	// 等待窗口切换完成 - 增加等待时间确保目标窗口完全激活
	fmt.Printf("[键盘服务] 步骤2: 等待窗口切换完成 (500ms)\n")
	time.Sleep(500 * time.Millisecond)

	// 输入文本，不再切换窗口
	fmt.Printf("[键盘服务] 步骤3: 开始输入文本\n")
	if err := ks.SimulateText(text); err != nil {
		fmt.Printf("[键盘服务] ❌ 文本输入失败: %v\n", err)
		return err
	}
	fmt.Printf("[键盘服务] ✅ 文本输入完成\n")

	// 等待输入完成
	fmt.Printf("[键盘服务] 步骤4: 等待输入完成 (100ms)\n")
	time.Sleep(100 * time.Millisecond)

	// 输入完成后切换回密码管理器
	fmt.Printf("[键盘服务] 步骤5: 切换回密码管理器\n")
	if !ks.windowMonitor.SwitchToPasswordManager() {
		fmt.Printf("[键盘服务] ⚠️ 切换回密码管理器失败\n")
	}

	fmt.Printf("[键盘服务] ========== 自动输入完成 ==========\n")
	return nil
}

/**
 * SimulateUsernameAndPasswordWithAutoSwitch 模拟输入用户名和密码（带自动窗口切换）
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
		fmt.Printf("[键盘服务] ❌ 缺少辅助功能权限\n")
		return errors.New("需要辅助功能权限才能模拟键盘输入，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用")
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

	// 20251005 陈凤庆 等待窗口切换完成 - 增加等待时间确保目标窗口完全激活
	fmt.Printf("[键盘服务] 步骤2: 等待窗口切换完成 (500ms)\n")
	time.Sleep(500 * time.Millisecond)

	// 20251002 陈凤庆 输入用户名
	if username != "" {
		fmt.Printf("[键盘服务] 步骤3: 输入用户名\n")
		if err := ks.SimulateText(username); err != nil {
			fmt.Printf("[键盘服务] ❌ 用户名输入失败: %v\n", err)
			return err
		}
		fmt.Printf("[键盘服务] ✅ 用户名输入完成\n")
		// 20251005 陈凤庆 用户名输入完成后等待更长时间，确保目标窗口完全处理
		fmt.Printf("[键盘服务] 等待用户名输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	} else {
		fmt.Printf("[键盘服务] 步骤3: 跳过用户名输入（为空）\n")
	}

	// 20251002 陈凤庆 按Tab键切换到密码输入框
	fmt.Printf("[键盘服务] 步骤4: 按Tab键切换到密码输入框\n")
	if err := ks.SimulateTab(); err != nil {
		fmt.Printf("[键盘服务] ❌ Tab键输入失败: %v\n", err)
		return err
	}
	fmt.Printf("[键盘服务] ✅ Tab键输入完成\n")
	// 20251005 陈凤庆 Tab键后等待更长时间，确保焦点正确切换到密码输入框
	fmt.Printf("[键盘服务] 等待Tab键焦点切换完成 (400ms)\n")
	time.Sleep(400 * time.Millisecond)

	// 20251002 陈凤庆 输入密码
	if password != "" {
		fmt.Printf("[键盘服务] 步骤5: 输入密码\n")
		if err := ks.SimulateText(password); err != nil {
			fmt.Printf("[键盘服务] ❌ 密码输入失败: %v\n", err)
			return err
		}
		fmt.Printf("[键盘服务] ✅ 密码输入完成\n")
		// 20251005 陈凤庆 密码输入完成后等待更长时间，确保目标窗口完全处理
		fmt.Printf("[键盘服务] 等待密码输入处理完成 (400ms)\n")
		time.Sleep(400 * time.Millisecond)
	} else {
		fmt.Printf("[键盘服务] 步骤5: 跳过密码输入（为空）\n")
	}

	// 20251005 陈凤庆 全部输入完成后等待一段时间再切换回密码管理器
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
 * RecordCurrentWindow 记录当前活动窗口为lastPID
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
 * GetLastWindowInfo 获取最后活动窗口信息
 * @return *WindowInfo 窗口信息
 * @description 20251002 陈凤庆 获取lastPID对应的窗口信息
 */
func (ks *KeyboardService) GetLastWindowInfo() *WindowInfo {
	if ks.windowMonitor == nil {
		return nil
	}

	return ks.windowMonitor.GetLastWindow()
}

/**
 * StopWindowMonitoring 停止窗口监控
 * @return error 错误信息
 * @description 20251002 陈凤庆 停止窗口监控服务
 */
func (ks *KeyboardService) StopWindowMonitoring() error {
	if ks.windowMonitor == nil {
		return errors.New("窗口监控服务未初始化")
	}

	return ks.windowMonitor.StopMonitoring()
}

/**
 * SwitchToLastWindow 切换到最后活动的窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活lastPID对应的窗口为活动窗口
 */
func (ks *KeyboardService) SwitchToLastWindow() bool {
	if ks.windowMonitor == nil {
		return false
	}
	return ks.windowMonitor.SwitchToLastWindow()
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活密码管理器窗口为活动窗口
 */
func (ks *KeyboardService) SwitchToPasswordManager() bool {
	if ks.windowMonitor == nil {
		return false
	}
	return ks.windowMonitor.SwitchToPasswordManager()
}
