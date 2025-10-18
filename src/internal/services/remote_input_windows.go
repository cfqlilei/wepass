//go:build windows

package services

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const (
	// Windows API常量
	INPUT_KEYBOARD    = 1
	KEYEVENTF_KEYUP   = 0x0002
	KEYEVENTF_UNICODE = 0x0004
	// VK_TAB 已在 global_hotkey_windows.go 中定义
)

// INPUT结构体
type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [8]byte // 填充到24字节
}

// KEYBDINPUT结构体
type KEYBDINPUT struct {
	Wvk         uint16
	Wscan       uint16
	Dwflags     uint32
	Time        uint32
	DwextraInfo uintptr
}

var (
	// user32 已在 global_hotkey_windows.go 中定义，这里使用不同的变量名
	user32RemoteInput = syscall.NewLazyDLL("user32.dll")
	procSendInput     = user32RemoteInput.NewProc("SendInput")
)

/**
 * inputCharacterWindows 在Windows上输入单个字符
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 使用SendInput API生成键盘事件，注入到系统输入队列
 */
func (ris *RemoteInputService) inputCharacterWindows(char rune) error {
	fmt.Printf("[远程输入-Windows] 使用SendInput输入字符: '%c' (Unicode: %d)\n", char, char)

	// 创建按键按下事件
	inputDown := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk:     0,
			Wscan:   uint16(char),
			Dwflags: KEYEVENTF_UNICODE,
			Time:    0,
		},
	}

	// 创建按键释放事件
	inputUp := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk:     0,
			Wscan:   uint16(char),
			Dwflags: KEYEVENTF_UNICODE | KEYEVENTF_KEYUP,
			Time:    0,
		},
	}

	// 发送按键按下事件
	ret1, _, err1 := procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputDown)),
		uintptr(unsafe.Sizeof(inputDown)),
	)

	if ret1 == 0 {
		return fmt.Errorf("SendInput按键按下失败: %v", err1)
	}

	// 等待一小段时间
	time.Sleep(5 * time.Millisecond)

	// 发送按键释放事件
	ret2, _, err2 := procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputUp)),
		uintptr(unsafe.Sizeof(inputUp)),
	)

	if ret2 == 0 {
		return fmt.Errorf("SendInput按键释放失败: %v", err2)
	}

	// 等待事件处理完成
	time.Sleep(10 * time.Millisecond)

	fmt.Printf("[远程输入-Windows] ✅ 字符输入完成\n")
	return nil
}

/**
 * inputTabKeyWindows 在Windows上输入Tab键
 * @return error 错误信息
 * @description 使用SendInput API输入Tab键
 */
func (ris *RemoteInputService) inputTabKeyWindows() error {
	fmt.Printf("[远程输入-Windows] 使用SendInput输入Tab键\n")

	// 创建Tab键按下事件
	inputDown := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk:     VK_TAB,
			Wscan:   0,
			Dwflags: 0,
			Time:    0,
		},
	}

	// 创建Tab键释放事件
	inputUp := INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk:     VK_TAB,
			Wscan:   0,
			Dwflags: KEYEVENTF_KEYUP,
			Time:    0,
		},
	}

	// 发送按键按下事件
	ret1, _, err1 := procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputDown)),
		uintptr(unsafe.Sizeof(inputDown)),
	)

	if ret1 == 0 {
		return fmt.Errorf("SendInput Tab按键按下失败: %v", err1)
	}

	// 等待一小段时间
	time.Sleep(5 * time.Millisecond)

	// 发送按键释放事件
	ret2, _, err2 := procSendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputUp)),
		uintptr(unsafe.Sizeof(inputUp)),
	)

	if ret2 == 0 {
		return fmt.Errorf("SendInput Tab按键释放失败: %v", err2)
	}

	// 等待事件处理完成
	time.Sleep(10 * time.Millisecond)

	fmt.Printf("[远程输入-Windows] ✅ Tab键输入完成\n")
	return nil
}

/**
 * inputCharacterMacOS Windows平台的别名方法
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 为了保持跨平台兼容性，提供与macOS相同的方法名
 * @modify 20251006 陈凤庆 添加跨平台兼容性别名方法
 */
func (ris *RemoteInputService) inputCharacterMacOS(char rune) error {
	return ris.inputCharacterWindows(char)
}

/**
 * inputTabKeyMacOS Windows平台的别名方法
 * @return error 错误信息
 * @description 为了保持跨平台兼容性，提供与macOS相同的方法名
 * @modify 20251006 陈凤庆 添加跨平台兼容性别名方法
 */
func (ris *RemoteInputService) inputTabKeyMacOS() error {
	return ris.inputTabKeyWindows()
}
