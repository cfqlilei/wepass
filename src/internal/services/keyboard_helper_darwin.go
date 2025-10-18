//go:build darwin
// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -framework ApplicationServices

#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>
#import <ApplicationServices/ApplicationServices.h>

// 发送字符到活动窗口
void sendCharToActiveWindowMac(unichar ch) {
    CGEventRef keyDownEvent = CGEventCreateKeyboardEvent(NULL, 0, true);
    CGEventRef keyUpEvent = CGEventCreateKeyboardEvent(NULL, 0, false);

    // 使用Unicode字符
    CGEventKeyboardSetUnicodeString(keyDownEvent, 1, &ch);
    CGEventKeyboardSetUnicodeString(keyUpEvent, 1, &ch);

    CGEventPost(kCGHIDEventTap, keyDownEvent);
    usleep(20000); // 20ms延迟
    CGEventPost(kCGHIDEventTap, keyUpEvent);

    CFRelease(keyDownEvent);
    CFRelease(keyUpEvent);
}

// 发送Tab键到活动窗口
void sendTabToActiveWindowMac() {
    CGEventRef keyDownEvent = CGEventCreateKeyboardEvent(NULL, 48, true); // Tab键的虚拟键码是48
    CGEventRef keyUpEvent = CGEventCreateKeyboardEvent(NULL, 48, false);

    CGEventPost(kCGHIDEventTap, keyDownEvent);
    usleep(20000); // 20ms延迟
    CGEventPost(kCGHIDEventTap, keyUpEvent);

    CFRelease(keyDownEvent);
    CFRelease(keyUpEvent);
}

// 检查进程是否存在
int checkProcessExistsMac(int pid) {
    NSRunningApplication *app = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
    return (app != nil && ![app isTerminated]) ? 1 : 0;
}

// 终止进程
int terminateProcessMac(int pid) {
    NSRunningApplication *app = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
    if (app == nil) {
        return -1;
    }

    BOOL result = [app terminate];
    return result ? 0 : -1;
}
*/
import "C"
import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

/**
 * DarwinKeyboardHelper macOS平台键盘助手实现
 * @author 陈凤庆
 * @description 使用AppleScript实现macOS平台的键盘助手功能
 * @date 20251005
 */
type DarwinKeyboardHelper struct {
	helperScriptPID int // AppleScript进程PID
}

/**
 * NewDarwinKeyboardHelper 创建macOS键盘助手实例
 * @return *DarwinKeyboardHelper 实例
 * @return error 错误信息
 */
func NewDarwinKeyboardHelper() (*DarwinKeyboardHelper, error) {
	// 检查osascript是否可用
	_, err := exec.LookPath("osascript")
	if err != nil {
		return nil, fmt.Errorf("未找到osascript: %w", err)
	}

	return &DarwinKeyboardHelper{
		helperScriptPID: 0,
	}, nil
}

/**
 * newKeyboardHelperPlatformImpl 创建平台特定的键盘助手实现 (macOS版本)
 * @return KeyboardHelperPlatform 平台实现
 * @return error 错误信息
 */
func newKeyboardHelperPlatformImpl() (KeyboardHelperPlatform, error) {
	return NewDarwinKeyboardHelper()
}

/**
 * LaunchHelper 启动键盘助手程序
 * @return int 键盘助手进程PID
 * @return error 错误信息
 */
func (d *DarwinKeyboardHelper) LaunchHelper() (int, error) {
	fmt.Printf("[macOS键盘助手] 启动辅助功能键盘\n")

	// 使用AppleScript打开辅助功能键盘
	// 注意：这需要用户在系统偏好设置中启用辅助功能键盘
	script := `
		tell application "System Events"
			-- 尝试显示辅助功能键盘
			try
				do shell script "open -a 'Accessibility Keyboard'"
			on error
				-- 如果直接打开失败，尝试通过系统偏好设置启用
				tell application "System Preferences"
					reveal anchor "Keyboard" of pane id "com.apple.preference.universalaccess"
				end tell
				return "请在系统偏好设置中启用辅助功能键盘"
			end try
		end tell
	`

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil || strings.Contains(outputStr, "请在系统偏好设置") {
		fmt.Printf("[macOS键盘助手] ⚠️ 无法自动打开辅助功能键盘\n")
		fmt.Printf("[macOS键盘助手] 提示：请手动启用：系统偏好设置 -> 辅助功能 -> 键盘 -> 启用辅助功能键盘\n")
		fmt.Printf("[macOS键盘助手] 将使用CGEvent作为备用方案\n")
	} else {
		fmt.Printf("[macOS键盘助手] ✅ 辅助功能键盘已打开\n")
	}

	// 返回虚拟PID
	d.helperScriptPID = 99999
	return 99999, nil
}

/**
 * CloseHelper 关闭键盘助手程序
 * @param helperPID 键盘助手进程PID
 * @return error 错误信息
 */
func (d *DarwinKeyboardHelper) CloseHelper(helperPID int) error {
	fmt.Printf("[macOS键盘助手] 关闭键盘助手（无需操作）\n")
	// macOS使用AppleScript，不需要关闭独立进程
	d.helperScriptPID = 0
	return nil
}

/**
 * SimulateCharWithHelper 使用键盘助手输入单个字符
 * @param char 要输入的字符
 * @param helperPID 键盘助手进程PID
 * @return error 错误信息
 */
func (d *DarwinKeyboardHelper) SimulateCharWithHelper(char rune, helperPID int) error {
	// 处理Tab键
	if char == '\t' {
		C.sendTabToActiveWindowMac()
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	// 使用CGEvent发送字符
	C.sendCharToActiveWindowMac(C.unichar(char))
	time.Sleep(50 * time.Millisecond)

	return nil
}

/**
 * CheckHelperAvailable 检查键盘助手是否可用
 * @return bool 是否可用
 */
func (d *DarwinKeyboardHelper) CheckHelperAvailable() bool {
	// 检查osascript是否可用
	_, err := exec.LookPath("osascript")
	return err == nil
}

/**
 * GetHelperName 获取键盘助手名称
 * @return string 键盘助手名称
 */
func (d *DarwinKeyboardHelper) GetHelperName() string {
	return "macOS键盘助手 (AppleScript/CGEvent)"
}

/**
 * executeAppleScriptKeystroke 使用AppleScript执行按键操作（备用方法）
 * @param key 要按的键
 * @return error 执行错误
 */
func (d *DarwinKeyboardHelper) executeAppleScriptKeystroke(key string) error {
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
 * executeAppleScriptTab 使用AppleScript执行Tab键操作（备用方法）
 * @return error 执行错误
 */
func (d *DarwinKeyboardHelper) executeAppleScriptTab() error {
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
