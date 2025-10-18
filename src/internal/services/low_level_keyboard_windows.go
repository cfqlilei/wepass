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

// 发送底层键盘事件
// 20251005 陈凤庆 修复按键映射问题,只使用虚拟键码,不使用扫描码
int sendLowLevelKeyEvent(int virtualKey, int scanCode, int isKeyDown, int modifiers) {
    INPUT input = {0};
    input.type = INPUT_KEYBOARD;
    input.ki.wVk = virtualKey;
    input.ki.wScan = 0; // 20251005 陈凤庆 不使用扫描码,避免按键映射错误
    input.ki.dwFlags = isKeyDown ? 0 : KEYEVENTF_KEYUP;

    // 20251005 陈凤庆 不使用KEYEVENTF_SCANCODE标志,只使用虚拟键码
    // 这样可以确保按键映射正确,避免cfq变成dgy的问题

    input.ki.time = 0;
    input.ki.dwExtraInfo = 0;

    UINT result = SendInput(1, &input, sizeof(INPUT));
    return result == 1 ? 0 : -1;
}

// 获取键盘布局标识符
char* getKeyboardLayoutName() {
    HKL hkl = GetKeyboardLayout(0);
    static char layoutName[32];
    sprintf(layoutName, "%08X", (unsigned int)(uintptr_t)hkl);
    return layoutName;
}

// 将字符映射到虚拟键码和扫描码
int mapCharToVirtualKey(wchar_t wc, int* virtualKey, int* scanCode, int* needsShift) {
    // 使用VkKeyScanW获取虚拟键码
    SHORT vkResult = VkKeyScanW(wc);

    if (vkResult == -1) {
        return -1; // 无法映射
    }

    *virtualKey = LOBYTE(vkResult);
    *needsShift = (HIBYTE(vkResult) & 1) ? 1 : 0;

    // 使用MapVirtualKeyW获取扫描码
    *scanCode = MapVirtualKeyW(*virtualKey, MAPVK_VK_TO_VSC);

    return 0;
}

// 检查是否有必要的权限
int checkKeyboardPermissions() {
    // Windows通常不需要特殊权限，但可以检查是否在安全桌面上
    HDESK hDesk = GetThreadDesktop(GetCurrentThreadId());
    if (hDesk == NULL) {
        return 0;
    }

    char deskName[256];
    DWORD needed;
    if (!GetUserObjectInformationA(hDesk, UOI_NAME, deskName, sizeof(deskName), &needed)) {
        return 0;
    }

    // 检查是否在默认桌面上
    return strcmp(deskName, "Default") == 0 ? 1 : 0;
}
*/
import "C"

import (
	"fmt"
)

/**
 * WindowsLowLevelKeyboard Windows平台底层键盘实现
 * @author 陈凤庆
 * @description 使用Windows SendInput API实现底层键盘输入
 * @date 20251004
 */
type WindowsLowLevelKeyboard struct {
	keyboardLayout string
}

/**
 * NewWindowsLowLevelKeyboard 创建Windows底层键盘实例
 * @return *WindowsLowLevelKeyboard 实例
 * @return error 错误信息
 */
func NewWindowsLowLevelKeyboard() (*WindowsLowLevelKeyboard, error) {
	return &WindowsLowLevelKeyboard{}, nil
}

/**
 * newPlatformImpl 创建平台特定的底层键盘实现 (Windows版本)
 * @return LowLevelKeyboardPlatform 平台实现
 * @return error 错误信息
 */
func newPlatformImpl() (LowLevelKeyboardPlatform, error) {
	return NewWindowsLowLevelKeyboard()
}

/**
 * SendKeyEvent 发送键盘事件
 * @param virtualKey 虚拟键码
 * @param scanCode 扫描码
 * @param isKeyDown 是否为按下事件
 * @param modifiers 修饰键列表
 * @return error 错误信息
 */
func (w *WindowsLowLevelKeyboard) SendKeyEvent(virtualKey int, scanCode int, isKeyDown bool, modifiers []string) error {
	keyDown := 0
	if isKeyDown {
		keyDown = 1
	}

	result := C.sendLowLevelKeyEvent(C.int(virtualKey), C.int(scanCode), C.int(keyDown), C.int(0))
	if result != 0 {
		return fmt.Errorf("SendInput失败，错误码: %d", result)
	}

	return nil
}

/**
 * GetKeyboardLayout 获取当前键盘布局
 * @return string 键盘布局标识符
 * @return error 错误信息
 */
func (w *WindowsLowLevelKeyboard) GetKeyboardLayout() (string, error) {
	layoutPtr := C.getKeyboardLayoutName()
	layout := C.GoString(layoutPtr)

	// 将Windows键盘布局标识符转换为友好名称
	friendlyName := w.convertLayoutToFriendlyName(layout)
	w.keyboardLayout = friendlyName

	return friendlyName, nil
}

/**
 * MapCharToKey 将字符映射到键码
 * @param char 字符
 * @param layout 键盘布局
 * @return KeyMapping 键码映射
 * @return error 错误信息
 */
func (w *WindowsLowLevelKeyboard) MapCharToKey(char rune, layout string) (KeyMapping, error) {
	// 将Go rune转换为Windows wchar_t
	wc := C.wchar_t(char)

	var virtualKey, scanCode, needsShift C.int

	result := C.mapCharToVirtualKey(wc, &virtualKey, &scanCode, &needsShift)
	if result != 0 {
		return KeyMapping{}, fmt.Errorf("无法映射字符'%c'到虚拟键码", char)
	}

	mapping := KeyMapping{
		VirtualKey:  int(virtualKey),
		ScanCode:    int(scanCode),
		NeedsShift:  needsShift != 0,
		NeedsAlt:    false,
		NeedsCtrl:   false,
		Description: fmt.Sprintf("VK:%d SC:%d", int(virtualKey), int(scanCode)),
	}

	// 处理特殊字符的额外修饰键需求
	w.adjustSpecialCharModifiers(char, &mapping)

	return mapping, nil
}

/**
 * CheckPermissions 检查权限
 * @return bool 是否有权限
 */
func (w *WindowsLowLevelKeyboard) CheckPermissions() bool {
	result := C.checkKeyboardPermissions()
	return result != 0
}

/**
 * convertLayoutToFriendlyName 将Windows布局标识符转换为友好名称
 * @param layoutId 布局标识符
 * @return string 友好名称
 */
func (w *WindowsLowLevelKeyboard) convertLayoutToFriendlyName(layoutId string) string {
	// Windows键盘布局标识符映射
	layoutMap := map[string]string{
		"00000409": "US", // 美式英语
		"00000804": "CN", // 中文(简体)
		"00000404": "TW", // 中文(繁体)
		"00000411": "JP", // 日语
		"00000412": "KR", // 韩语
		"00000407": "DE", // 德语
		"0000040C": "FR", // 法语
		"00000410": "IT", // 意大利语
		"0000040A": "ES", // 西班牙语
		"00000419": "RU", // 俄语
	}

	if friendlyName, exists := layoutMap[layoutId]; exists {
		return friendlyName
	}

	// 默认返回US布局
	return "US"
}

/**
 * adjustSpecialCharModifiers 调整特殊字符的修饰键需求
 * @param char 字符
 * @param mapping 键码映射指针
 */
func (w *WindowsLowLevelKeyboard) adjustSpecialCharModifiers(char rune, mapping *KeyMapping) {
	// 根据字符和键盘布局调整修饰键需求
	switch char {
	case '@':
		// 在某些布局中@需要AltGr键
		if w.keyboardLayout == "DE" || w.keyboardLayout == "FR" {
			mapping.NeedsAlt = true
			mapping.NeedsCtrl = true // AltGr = Ctrl + Alt
		}
	case '#':
		// 在英国布局中#需要Shift+3
		if w.keyboardLayout == "GB" {
			mapping.NeedsShift = true
		}
	case '€':
		// 欧元符号通常需要AltGr
		mapping.NeedsAlt = true
		mapping.NeedsCtrl = true
	}
}

/**
 * GetSpecialKeyMapping 获取特殊键的映射
 * @param keyName 键名
 * @return KeyMapping 键码映射
 * @return error 错误信息
 */
func (w *WindowsLowLevelKeyboard) GetSpecialKeyMapping(keyName string) (KeyMapping, error) {
	// Windows特殊键映射
	specialKeys := map[string]KeyMapping{
		"tab":       {VirtualKey: 9, ScanCode: 15, Description: "Tab"},
		"enter":     {VirtualKey: 13, ScanCode: 28, Description: "Enter"},
		"space":     {VirtualKey: 32, ScanCode: 57, Description: "Space"},
		"backspace": {VirtualKey: 8, ScanCode: 14, Description: "Backspace"},
		"delete":    {VirtualKey: 46, ScanCode: 83, Description: "Delete"},
		"escape":    {VirtualKey: 27, ScanCode: 1, Description: "Escape"},
		"f1":        {VirtualKey: 112, ScanCode: 59, Description: "F1"},
		"f2":        {VirtualKey: 113, ScanCode: 60, Description: "F2"},
		"f3":        {VirtualKey: 114, ScanCode: 61, Description: "F3"},
		"f4":        {VirtualKey: 115, ScanCode: 62, Description: "F4"},
		"f5":        {VirtualKey: 116, ScanCode: 63, Description: "F5"},
		"f6":        {VirtualKey: 117, ScanCode: 64, Description: "F6"},
		"f7":        {VirtualKey: 118, ScanCode: 65, Description: "F7"},
		"f8":        {VirtualKey: 119, ScanCode: 66, Description: "F8"},
		"f9":        {VirtualKey: 120, ScanCode: 67, Description: "F9"},
		"f10":       {VirtualKey: 121, ScanCode: 68, Description: "F10"},
		"f11":       {VirtualKey: 122, ScanCode: 87, Description: "F11"},
		"f12":       {VirtualKey: 123, ScanCode: 88, Description: "F12"},
		"up":        {VirtualKey: 38, ScanCode: 72, Description: "Up Arrow"},
		"down":      {VirtualKey: 40, ScanCode: 80, Description: "Down Arrow"},
		"left":      {VirtualKey: 37, ScanCode: 75, Description: "Left Arrow"},
		"right":     {VirtualKey: 39, ScanCode: 77, Description: "Right Arrow"},
		"home":      {VirtualKey: 36, ScanCode: 71, Description: "Home"},
		"end":       {VirtualKey: 35, ScanCode: 79, Description: "End"},
		"pageup":    {VirtualKey: 33, ScanCode: 73, Description: "Page Up"},
		"pagedown":  {VirtualKey: 34, ScanCode: 81, Description: "Page Down"},
		"insert":    {VirtualKey: 45, ScanCode: 82, Description: "Insert"},
	}

	if mapping, exists := specialKeys[keyName]; exists {
		return mapping, nil
	}

	return KeyMapping{}, fmt.Errorf("未知的特殊键: %s", keyName)
}

/**
 * SimulateSpecialKey 模拟特殊键输入
 * @param keyName 键名
 * @return error 错误信息
 */
func (w *WindowsLowLevelKeyboard) SimulateSpecialKey(keyName string) error {
	mapping, err := w.GetSpecialKeyMapping(keyName)
	if err != nil {
		return err
	}

	// 按下键
	if err := w.SendKeyEvent(mapping.VirtualKey, mapping.ScanCode, true, nil); err != nil {
		return fmt.Errorf("按下%s键失败: %w", keyName, err)
	}

	// 短暂延迟
	// time.Sleep(10 * time.Millisecond)

	// 释放键
	if err := w.SendKeyEvent(mapping.VirtualKey, mapping.ScanCode, false, nil); err != nil {
		return fmt.Errorf("释放%s键失败: %w", keyName, err)
	}

	return nil
}
