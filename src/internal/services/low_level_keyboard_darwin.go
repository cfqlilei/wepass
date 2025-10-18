//go:build darwin
// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework Foundation -framework AppKit -framework Carbon

#include <ApplicationServices/ApplicationServices.h>
#include <Foundation/Foundation.h>
#include <AppKit/AppKit.h>
#include <Carbon/Carbon.h>
#include <unistd.h>
#include <stdlib.h>

// 函数声明
int mapCharToVirtualKeyWithUSLayout(unichar character, int* virtualKey, int* needsShift);

// 发送底层键盘事件
int sendLowLevelKeyEventMac(int virtualKey, int isKeyDown) {
    CGEventRef keyEvent = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)virtualKey, isKeyDown);
    if (keyEvent == NULL) {
        return -1;
    }

    // 设置事件源为系统级别，提高兼容性
    CGEventSetSource(keyEvent, NULL);

    // 发送事件到系统
    CGEventPost(kCGHIDEventTap, keyEvent);

    CFRelease(keyEvent);
    return 0;
}

// 发送带修饰键的键盘事件
int sendKeyEventWithModifiers(int virtualKey, int isKeyDown, int modifiers) {
    CGEventRef keyEvent = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)virtualKey, isKeyDown);
    if (keyEvent == NULL) {
        return -1;
    }

    // 设置修饰键标志
    CGEventFlags flags = 0;
    if (modifiers & 1) flags |= kCGEventFlagMaskShift;
    if (modifiers & 2) flags |= kCGEventFlagMaskControl;
    if (modifiers & 4) flags |= kCGEventFlagMaskAlternate;
    if (modifiers & 8) flags |= kCGEventFlagMaskCommand;

    CGEventSetFlags(keyEvent, flags);
    CGEventSetSource(keyEvent, NULL);

    CGEventPost(kCGHIDEventTap, keyEvent);
    CFRelease(keyEvent);
    return 0;
}

// 获取键盘布局标识符
char* getKeyboardLayoutMac() {
    @autoreleasepool {
        TISInputSourceRef currentKeyboard = TISCopyCurrentKeyboardInputSource();
        if (currentKeyboard == NULL) {
            return strdup("US");
        }

        CFStringRef layoutID = (CFStringRef)TISGetInputSourceProperty(currentKeyboard, kTISPropertyInputSourceID);
        if (layoutID == NULL) {
            CFRelease(currentKeyboard);
            return strdup("US");
        }

        const char* layoutCStr = CFStringGetCStringPtr(layoutID, kCFStringEncodingUTF8);
        char* result;
        if (layoutCStr != NULL) {
            result = strdup(layoutCStr);
        } else {
            // 如果无法直接获取C字符串，使用CFStringGetCString
            CFIndex length = CFStringGetLength(layoutID);
            CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
            result = malloc(maxSize);
            if (!CFStringGetCString(layoutID, result, maxSize, kCFStringEncodingUTF8)) {
                free(result);
                result = strdup("US");
            }
        }

        CFRelease(currentKeyboard);
        return result;
    }
}

// 将字符映射到虚拟键码
int mapCharToVirtualKeyMac(unichar character, int* virtualKey, int* needsShift) {
    @autoreleasepool {
        // 首先尝试使用固定的US键盘布局映射
        // 这样可以确保英文字符和数字的正确映射
        int result = mapCharToVirtualKeyWithUSLayout(character, virtualKey, needsShift);
        if (result == 0) {
            return 0;
        }

        // 如果US布局映射失败，再尝试当前键盘布局
        TISInputSourceRef currentKeyboard = TISCopyCurrentKeyboardInputSource();
        if (currentKeyboard == NULL) {
            return -1;
        }

        CFDataRef layoutData = (CFDataRef)TISGetInputSourceProperty(currentKeyboard, kTISPropertyUnicodeKeyLayoutData);
        if (layoutData == NULL) {
            CFRelease(currentKeyboard);
            return -1;
        }

        const UCKeyboardLayout* keyboardLayout = (const UCKeyboardLayout*)CFDataGetBytePtr(layoutData);

        // 遍历所有可能的虚拟键码
        for (int vk = 0; vk < 128; vk++) {
            UInt32 deadKeyState = 0;
            UniChar unicodeString[4];
            UniCharCount actualStringLength = 0;

            // 不带修饰键的情况
            OSStatus status = UCKeyTranslate(keyboardLayout, vk, kUCKeyActionDown, 0,
                                           LMGetKbdType(), kUCKeyTranslateNoDeadKeysBit,
                                           &deadKeyState, 4, &actualStringLength, unicodeString);

            if (status == noErr && actualStringLength > 0 && unicodeString[0] == character) {
                *virtualKey = vk;
                *needsShift = 0;
                CFRelease(currentKeyboard);
                return 0;
            }

            // 带Shift修饰键的情况
            deadKeyState = 0;
            status = UCKeyTranslate(keyboardLayout, vk, kUCKeyActionDown, shiftKey >> 8,
                                  LMGetKbdType(), kUCKeyTranslateNoDeadKeysBit,
                                  &deadKeyState, 4, &actualStringLength, unicodeString);

            if (status == noErr && actualStringLength > 0 && unicodeString[0] == character) {
                *virtualKey = vk;
                *needsShift = 1;
                CFRelease(currentKeyboard);
                return 0;
            }
        }

        CFRelease(currentKeyboard);
        return -1;
    }
}

// 使用固定的US键盘布局映射字符
int mapCharToVirtualKeyWithUSLayout(unichar character, int* virtualKey, int* needsShift) {
    // 基于US键盘布局的固定映射表
    // 这确保了英文字符和数字的正确映射，不受当前键盘布局影响

    // 小写字母 a-z (虚拟键码 0-25)
    if (character >= 'a' && character <= 'z') {
        *virtualKey = character - 'a';
        *needsShift = 0;
        return 0;
    }

    // 大写字母 A-Z (虚拟键码 0-25 + Shift)
    if (character >= 'A' && character <= 'Z') {
        *virtualKey = character - 'A';
        *needsShift = 1;
        return 0;
    }

    // 数字和特殊字符
    switch (character) {
        // 数字行 (不带Shift)
        case '1': *virtualKey = 18; *needsShift = 0; return 0;
        case '2': *virtualKey = 19; *needsShift = 0; return 0;
        case '3': *virtualKey = 20; *needsShift = 0; return 0;
        case '4': *virtualKey = 21; *needsShift = 0; return 0;
        case '5': *virtualKey = 23; *needsShift = 0; return 0;
        case '6': *virtualKey = 22; *needsShift = 0; return 0;
        case '7': *virtualKey = 26; *needsShift = 0; return 0;
        case '8': *virtualKey = 28; *needsShift = 0; return 0;
        case '9': *virtualKey = 25; *needsShift = 0; return 0;
        case '0': *virtualKey = 29; *needsShift = 0; return 0;

        // 数字行 (带Shift)
        case '!': *virtualKey = 18; *needsShift = 1; return 0;
        case '@': *virtualKey = 19; *needsShift = 1; return 0;
        case '#': *virtualKey = 20; *needsShift = 1; return 0;
        case '$': *virtualKey = 21; *needsShift = 1; return 0;
        case '%': *virtualKey = 23; *needsShift = 1; return 0;
        case '^': *virtualKey = 22; *needsShift = 1; return 0;
        case '&': *virtualKey = 26; *needsShift = 1; return 0;
        case '*': *virtualKey = 28; *needsShift = 1; return 0;
        case '(': *virtualKey = 25; *needsShift = 1; return 0;
        case ')': *virtualKey = 29; *needsShift = 1; return 0;

        // 其他常用字符
        case ' ': *virtualKey = 49; *needsShift = 0; return 0; // 空格
        case '-': *virtualKey = 27; *needsShift = 0; return 0;
        case '_': *virtualKey = 27; *needsShift = 1; return 0;
        case '=': *virtualKey = 24; *needsShift = 0; return 0;
        case '+': *virtualKey = 24; *needsShift = 1; return 0;
        case '[': *virtualKey = 33; *needsShift = 0; return 0;
        case '{': *virtualKey = 33; *needsShift = 1; return 0;
        case ']': *virtualKey = 30; *needsShift = 0; return 0;
        case '}': *virtualKey = 30; *needsShift = 1; return 0;
        case '\\': *virtualKey = 42; *needsShift = 0; return 0;
        case '|': *virtualKey = 42; *needsShift = 1; return 0;
        case ';': *virtualKey = 41; *needsShift = 0; return 0;
        case ':': *virtualKey = 41; *needsShift = 1; return 0;
        case '\'': *virtualKey = 39; *needsShift = 0; return 0;
        case '"': *virtualKey = 39; *needsShift = 1; return 0;
        case ',': *virtualKey = 43; *needsShift = 0; return 0;
        case '<': *virtualKey = 43; *needsShift = 1; return 0;
        case '.': *virtualKey = 47; *needsShift = 0; return 0;
        case '>': *virtualKey = 47; *needsShift = 1; return 0;
        case '/': *virtualKey = 44; *needsShift = 0; return 0;
        case '?': *virtualKey = 44; *needsShift = 1; return 0;
        case '`': *virtualKey = 50; *needsShift = 0; return 0;
        case '~': *virtualKey = 50; *needsShift = 1; return 0;

        default:
            return -1; // 未找到映射
    }
}

// 检查辅助功能权限
int checkAccessibilityPermissionsMac() {
    // 检查是否有辅助功能权限
    NSDictionary* options = @{(id)kAXTrustedCheckOptionPrompt: @NO};
    Boolean hasPermission = AXIsProcessTrustedWithOptions((CFDictionaryRef)options);
    return hasPermission ? 1 : 0;
}

// 请求辅助功能权限
void requestAccessibilityPermissionsMac() {
    NSDictionary* options = @{(id)kAXTrustedCheckOptionPrompt: @YES};
    AXIsProcessTrustedWithOptions((CFDictionaryRef)options);
}
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

/**
 * DarwinLowLevelKeyboard macOS平台底层键盘实现
 * @author 陈凤庆
 * @description 使用macOS CGEventPost API实现底层键盘输入
 * @date 20251004
 */
type DarwinLowLevelKeyboard struct {
	keyboardLayout string
}

/**
 * NewDarwinLowLevelKeyboard 创建macOS底层键盘实例
 * @return *DarwinLowLevelKeyboard 实例
 * @return error 错误信息
 */
func NewDarwinLowLevelKeyboard() (*DarwinLowLevelKeyboard, error) {
	return &DarwinLowLevelKeyboard{}, nil
}

/**
 * newPlatformImpl 创建平台特定的底层键盘实现 (macOS版本)
 * @return LowLevelKeyboardPlatform 平台实现
 * @return error 错误信息
 */
func newPlatformImpl() (LowLevelKeyboardPlatform, error) {
	return NewDarwinLowLevelKeyboard()
}

/**
 * SendKeyEvent 发送键盘事件
 * @param virtualKey 虚拟键码
 * @param scanCode 扫描码（macOS中不使用）
 * @param isKeyDown 是否为按下事件
 * @param modifiers 修饰键列表
 * @return error 错误信息
 */
func (d *DarwinLowLevelKeyboard) SendKeyEvent(virtualKey int, scanCode int, isKeyDown bool, modifiers []string) error {
	keyDown := 0
	if isKeyDown {
		keyDown = 1
	}

	// 构建修饰键标志位
	modifierFlags := 0
	for _, modifier := range modifiers {
		switch modifier {
		case "shift":
			modifierFlags |= 1
		case "ctrl":
			modifierFlags |= 2
		case "alt":
			modifierFlags |= 4
		case "cmd":
			modifierFlags |= 8
		}
	}

	var result C.int
	if len(modifiers) > 0 {
		result = C.sendKeyEventWithModifiers(C.int(virtualKey), C.int(keyDown), C.int(modifierFlags))
	} else {
		result = C.sendLowLevelKeyEventMac(C.int(virtualKey), C.int(keyDown))
	}

	if result != 0 {
		return fmt.Errorf("CGEventPost失败，错误码: %d", result)
	}

	return nil
}

/**
 * GetKeyboardLayout 获取当前键盘布局
 * @return string 键盘布局标识符
 * @return error 错误信息
 */
func (d *DarwinLowLevelKeyboard) GetKeyboardLayout() (string, error) {
	layoutPtr := C.getKeyboardLayoutMac()
	defer C.free(unsafe.Pointer(layoutPtr))

	layout := C.GoString(layoutPtr)

	// 将macOS键盘布局标识符转换为友好名称
	friendlyName := d.convertLayoutToFriendlyName(layout)
	d.keyboardLayout = friendlyName

	return friendlyName, nil
}

/**
 * MapCharToKey 将字符映射到键码
 * @param char 字符
 * @param layout 键盘布局
 * @return KeyMapping 键码映射
 * @return error 错误信息
 */
func (d *DarwinLowLevelKeyboard) MapCharToKey(char rune, layout string) (KeyMapping, error) {
	// 将Go rune转换为macOS unichar
	unichar := C.unichar(char)

	var virtualKey, needsShift C.int

	result := C.mapCharToVirtualKeyMac(unichar, &virtualKey, &needsShift)
	if result != 0 {
		return KeyMapping{}, fmt.Errorf("无法映射字符'%c'到虚拟键码", char)
	}

	mapping := KeyMapping{
		VirtualKey:  int(virtualKey),
		ScanCode:    0, // macOS不使用扫描码
		NeedsShift:  needsShift != 0,
		NeedsAlt:    false,
		NeedsCtrl:   false,
		Description: fmt.Sprintf("VK:%d", int(virtualKey)),
	}

	// 处理特殊字符的额外修饰键需求
	d.adjustSpecialCharModifiers(char, &mapping)

	return mapping, nil
}

/**
 * CheckPermissions 检查权限
 * @return bool 是否有权限
 */
func (d *DarwinLowLevelKeyboard) CheckPermissions() bool {
	result := C.checkAccessibilityPermissionsMac()
	return result != 0
}

/**
 * RequestPermissions 请求权限
 */
func (d *DarwinLowLevelKeyboard) RequestPermissions() {
	C.requestAccessibilityPermissionsMac()
}

/**
 * convertLayoutToFriendlyName 将macOS布局标识符转换为友好名称
 * @param layoutId 布局标识符
 * @return string 友好名称
 */
func (d *DarwinLowLevelKeyboard) convertLayoutToFriendlyName(layoutId string) string {
	// macOS键盘布局标识符映射
	layoutMap := map[string]string{
		"com.apple.keylayout.US":                 "US",
		"com.apple.keylayout.ABC":                "US",
		"com.apple.keylayout.USInternational-PC": "US-Intl",
		"com.apple.keylayout.German":             "DE",
		"com.apple.keylayout.French":             "FR",
		"com.apple.keylayout.Italian":            "IT",
		"com.apple.keylayout.Spanish":            "ES",
		"com.apple.keylayout.Russian":            "RU",
		"com.apple.keylayout.Japanese":           "JP",
		"com.apple.keylayout.Korean":             "KR",
		"com.apple.inputmethod.SCIM.ITABC":       "CN",
		"com.apple.inputmethod.TCIM.Cangjie":     "TW",
	}

	// 精确匹配
	if friendlyName, exists := layoutMap[layoutId]; exists {
		return friendlyName
	}

	// 模糊匹配
	for key, value := range layoutMap {
		if strings.Contains(layoutId, key) || strings.Contains(key, layoutId) {
			return value
		}
	}

	// 根据关键词匹配
	if strings.Contains(strings.ToLower(layoutId), "chinese") || strings.Contains(strings.ToLower(layoutId), "pinyin") {
		return "CN"
	}
	if strings.Contains(strings.ToLower(layoutId), "japanese") {
		return "JP"
	}
	if strings.Contains(strings.ToLower(layoutId), "korean") {
		return "KR"
	}

	// 默认返回US布局
	return "US"
}

/**
 * adjustSpecialCharModifiers 调整特殊字符的修饰键需求
 * @param char 字符
 * @param mapping 键码映射指针
 */
func (d *DarwinLowLevelKeyboard) adjustSpecialCharModifiers(char rune, mapping *KeyMapping) {
	// 根据字符和键盘布局调整修饰键需求
	switch char {
	case '@':
		// 在某些欧洲布局中@需要Option键
		if d.keyboardLayout == "DE" || d.keyboardLayout == "FR" {
			mapping.NeedsAlt = true
		}
	case '#':
		// 在某些布局中#需要特殊组合
		if d.keyboardLayout == "GB" {
			mapping.NeedsAlt = true
		}
	case '€':
		// 欧元符号通常需要Option键
		mapping.NeedsAlt = true
	case '~':
		// 波浪号在某些布局中需要Option键
		if d.keyboardLayout != "US" {
			mapping.NeedsAlt = true
		}
	}
}

/**
 * GetSpecialKeyMapping 获取特殊键的映射
 * @param keyName 键名
 * @return KeyMapping 键码映射
 * @return error 错误信息
 */
func (d *DarwinLowLevelKeyboard) GetSpecialKeyMapping(keyName string) (KeyMapping, error) {
	// macOS特殊键映射（基于虚拟键码）
	specialKeys := map[string]KeyMapping{
		"tab":       {VirtualKey: 48, Description: "Tab"},
		"enter":     {VirtualKey: 36, Description: "Enter"},
		"return":    {VirtualKey: 36, Description: "Return"},
		"space":     {VirtualKey: 49, Description: "Space"},
		"backspace": {VirtualKey: 51, Description: "Backspace"},
		"delete":    {VirtualKey: 117, Description: "Delete"},
		"escape":    {VirtualKey: 53, Description: "Escape"},
		"f1":        {VirtualKey: 122, Description: "F1"},
		"f2":        {VirtualKey: 120, Description: "F2"},
		"f3":        {VirtualKey: 99, Description: "F3"},
		"f4":        {VirtualKey: 118, Description: "F4"},
		"f5":        {VirtualKey: 96, Description: "F5"},
		"f6":        {VirtualKey: 97, Description: "F6"},
		"f7":        {VirtualKey: 98, Description: "F7"},
		"f8":        {VirtualKey: 100, Description: "F8"},
		"f9":        {VirtualKey: 101, Description: "F9"},
		"f10":       {VirtualKey: 109, Description: "F10"},
		"f11":       {VirtualKey: 103, Description: "F11"},
		"f12":       {VirtualKey: 111, Description: "F12"},
		"up":        {VirtualKey: 126, Description: "Up Arrow"},
		"down":      {VirtualKey: 125, Description: "Down Arrow"},
		"left":      {VirtualKey: 123, Description: "Left Arrow"},
		"right":     {VirtualKey: 124, Description: "Right Arrow"},
		"home":      {VirtualKey: 115, Description: "Home"},
		"end":       {VirtualKey: 119, Description: "End"},
		"pageup":    {VirtualKey: 116, Description: "Page Up"},
		"pagedown":  {VirtualKey: 121, Description: "Page Down"},
		"insert":    {VirtualKey: 114, Description: "Insert"},
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
func (d *DarwinLowLevelKeyboard) SimulateSpecialKey(keyName string) error {
	mapping, err := d.GetSpecialKeyMapping(keyName)
	if err != nil {
		return err
	}

	// 按下键
	if err := d.SendKeyEvent(mapping.VirtualKey, 0, true, nil); err != nil {
		return fmt.Errorf("按下%s键失败: %w", keyName, err)
	}

	// 短暂延迟
	// time.Sleep(10 * time.Millisecond)

	// 释放键
	if err := d.SendKeyEvent(mapping.VirtualKey, 0, false, nil); err != nil {
		return fmt.Errorf("释放%s键失败: %w", keyName, err)
	}

	return nil
}
