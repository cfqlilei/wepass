//go:build windows
// +build windows

package services

/*
#cgo CFLAGS: -I. -DUNICODE -D_UNICODE
#cgo LDFLAGS: -luser32 -lkernel32 -lshell32

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>

// 启动osk.exe键盘助手
DWORD launchOSK() {
    SHELLEXECUTEINFOW sei = {0};
    sei.cbSize = sizeof(SHELLEXECUTEINFOW);
    sei.fMask = SEE_MASK_NOCLOSEPROCESS;
    sei.lpVerb = L"open";
    sei.lpFile = L"osk.exe";
    sei.nShow = SW_SHOW;

    if (!ShellExecuteExW(&sei)) {
        return 0;
    }

    // 等待进程启动
    if (sei.hProcess) {
        WaitForInputIdle(sei.hProcess, 2000);
        DWORD pid = GetProcessId(sei.hProcess);
        CloseHandle(sei.hProcess);
        return pid;
    }

    return 0;
}

// 关闭进程
int terminateProcess(DWORD pid) {
    HANDLE hProcess = OpenProcess(PROCESS_TERMINATE, FALSE, pid);
    if (hProcess == NULL) {
        return -1;
    }

    BOOL result = TerminateProcess(hProcess, 0);
    CloseHandle(hProcess);

    return result ? 0 : -1;
}

// 检查进程是否存在
int checkProcessExists(DWORD pid) {
    HANDLE hProcess = OpenProcess(PROCESS_QUERY_INFORMATION, FALSE, pid);
    if (hProcess == NULL) {
        return 0;
    }

    DWORD exitCode;
    BOOL result = GetExitCodeProcess(hProcess, &exitCode);
    CloseHandle(hProcess);

    return (result && exitCode == STILL_ACTIVE) ? 1 : 0;
}

// 发送字符到活动窗口
void sendCharToActiveWindow(wchar_t ch) {
    INPUT inputs[2] = {0};
    
    // 按下键
    inputs[0].type = INPUT_KEYBOARD;
    inputs[0].ki.wVk = 0;
    inputs[0].ki.wScan = ch;
    inputs[0].ki.dwFlags = KEYEVENTF_UNICODE;
    
    // 释放键
    inputs[1].type = INPUT_KEYBOARD;
    inputs[1].ki.wVk = 0;
    inputs[1].ki.wScan = ch;
    inputs[1].ki.dwFlags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP;
    
    SendInput(2, inputs, sizeof(INPUT));
}

// 发送Tab键到活动窗口
void sendTabToActiveWindow() {
    INPUT inputs[2] = {0};
    
    // 按下Tab键
    inputs[0].type = INPUT_KEYBOARD;
    inputs[0].ki.wVk = VK_TAB;
    inputs[0].ki.dwFlags = 0;
    
    // 释放Tab键
    inputs[1].type = INPUT_KEYBOARD;
    inputs[1].ki.wVk = VK_TAB;
    inputs[1].ki.dwFlags = KEYEVENTF_KEYUP;
    
    SendInput(2, inputs, sizeof(INPUT));
}
*/
import "C"
import (
	"fmt"
	"os/exec"
	"time"
)

/**
 * WindowsKeyboardHelper Windows平台键盘助手实现
 * @author 陈凤庆
 * @description 使用osk.exe实现Windows平台的键盘助手功能
 * @date 20251005
 */
type WindowsKeyboardHelper struct {
	oskPath string // osk.exe路径
}

/**
 * NewWindowsKeyboardHelper 创建Windows键盘助手实例
 * @return *WindowsKeyboardHelper 实例
 * @return error 错误信息
 */
func NewWindowsKeyboardHelper() (*WindowsKeyboardHelper, error) {
	// 检查osk.exe是否存在
	oskPath, err := exec.LookPath("osk.exe")
	if err != nil {
		return nil, fmt.Errorf("未找到osk.exe: %w", err)
	}

	return &WindowsKeyboardHelper{
		oskPath: oskPath,
	}, nil
}

/**
 * newKeyboardHelperPlatformImpl 创建平台特定的键盘助手实现 (Windows版本)
 * @return KeyboardHelperPlatform 平台实现
 * @return error 错误信息
 */
func newKeyboardHelperPlatformImpl() (KeyboardHelperPlatform, error) {
	return NewWindowsKeyboardHelper()
}

/**
 * LaunchHelper 启动键盘助手程序
 * @return int 键盘助手进程PID
 * @return error 错误信息
 */
func (w *WindowsKeyboardHelper) LaunchHelper() (int, error) {
	fmt.Printf("[Windows键盘助手] 启动osk.exe\n")

	// 调用C函数启动osk.exe
	pid := C.launchOSK()
	if pid == 0 {
		return 0, fmt.Errorf("启动osk.exe失败")
	}

	fmt.Printf("[Windows键盘助手] ✅ osk.exe已启动，PID: %d\n", int(pid))
	return int(pid), nil
}

/**
 * CloseHelper 关闭键盘助手程序
 * @param helperPID 键盘助手进程PID
 * @return error 错误信息
 */
func (w *WindowsKeyboardHelper) CloseHelper(helperPID int) error {
	if helperPID <= 0 {
		return fmt.Errorf("无效的PID: %d", helperPID)
	}

	fmt.Printf("[Windows键盘助手] 关闭osk.exe，PID: %d\n", helperPID)

	// 检查进程是否存在
	if C.checkProcessExists(C.DWORD(helperPID)) == 0 {
		fmt.Printf("[Windows键盘助手] 进程已不存在，PID: %d\n", helperPID)
		return nil
	}

	// 调用C函数终止进程
	result := C.terminateProcess(C.DWORD(helperPID))
	if result != 0 {
		return fmt.Errorf("终止进程失败，PID: %d", helperPID)
	}

	fmt.Printf("[Windows键盘助手] ✅ osk.exe已关闭\n")
	return nil
}

/**
 * SimulateCharWithHelper 使用键盘助手输入单个字符
 * @param char 要输入的字符
 * @param helperPID 键盘助手进程PID
 * @return error 错误信息
 */
func (w *WindowsKeyboardHelper) SimulateCharWithHelper(char rune, helperPID int) error {
	// 检查键盘助手是否运行
	if C.checkProcessExists(C.DWORD(helperPID)) == 0 {
		return fmt.Errorf("键盘助手未运行，PID: %d", helperPID)
	}

	// 处理Tab键
	if char == '\t' {
		C.sendTabToActiveWindow()
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	// 发送字符到活动窗口
	C.sendCharToActiveWindow(C.wchar_t(char))
	time.Sleep(50 * time.Millisecond)

	return nil
}

/**
 * CheckHelperAvailable 检查键盘助手是否可用
 * @return bool 是否可用
 */
func (w *WindowsKeyboardHelper) CheckHelperAvailable() bool {
	// 检查osk.exe是否存在
	_, err := exec.LookPath("osk.exe")
	return err == nil
}

/**
 * GetHelperName 获取键盘助手名称
 * @return string 键盘助手名称
 */
func (w *WindowsKeyboardHelper) GetHelperName() string {
	return "Windows屏幕键盘 (osk.exe)"
}

