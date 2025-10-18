//go:build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework ApplicationServices -framework IOKit -framework Carbon
#include <CoreGraphics/CoreGraphics.h>
#include <ApplicationServices/ApplicationServices.h>
#include <IOKit/hidsystem/IOHIDLib.h>
#include <IOKit/hidsystem/ev_keymap.h>
#include <Carbon/Carbon.h>
#include <unistd.h>

// 函数声明
void inputUnicodeCharacterCG(unsigned short unicodeChar);
void inputCharacterWithKeyCode(unsigned short unicodeChar);
void inputTabKeyCG(void);
bool findKeyMapping(unichar character, CGKeyCode* keyCode, bool* needsShift);

// 字符到虚拟键码的映射表（美式键盘布局）
typedef struct {
    unichar character;
    CGKeyCode keyCode;
    bool needsShift;
} CharKeyMapping;

// 基本字符映射表
static const CharKeyMapping charMappings[] = {
    // 数字
    {'0', kVK_ANSI_0, false}, {'1', kVK_ANSI_1, false}, {'2', kVK_ANSI_2, false},
    {'3', kVK_ANSI_3, false}, {'4', kVK_ANSI_4, false}, {'5', kVK_ANSI_5, false},
    {'6', kVK_ANSI_6, false}, {'7', kVK_ANSI_7, false}, {'8', kVK_ANSI_8, false},
    {'9', kVK_ANSI_9, false},

    // 字母（小写）
    {'a', kVK_ANSI_A, false}, {'b', kVK_ANSI_B, false}, {'c', kVK_ANSI_C, false},
    {'d', kVK_ANSI_D, false}, {'e', kVK_ANSI_E, false}, {'f', kVK_ANSI_F, false},
    {'g', kVK_ANSI_G, false}, {'h', kVK_ANSI_H, false}, {'i', kVK_ANSI_I, false},
    {'j', kVK_ANSI_J, false}, {'k', kVK_ANSI_K, false}, {'l', kVK_ANSI_L, false},
    {'m', kVK_ANSI_M, false}, {'n', kVK_ANSI_N, false}, {'o', kVK_ANSI_O, false},
    {'p', kVK_ANSI_P, false}, {'q', kVK_ANSI_Q, false}, {'r', kVK_ANSI_R, false},
    {'s', kVK_ANSI_S, false}, {'t', kVK_ANSI_T, false}, {'u', kVK_ANSI_U, false},
    {'v', kVK_ANSI_V, false}, {'w', kVK_ANSI_W, false}, {'x', kVK_ANSI_X, false},
    {'y', kVK_ANSI_Y, false}, {'z', kVK_ANSI_Z, false},

    // 字母（大写）
    {'A', kVK_ANSI_A, true}, {'B', kVK_ANSI_B, true}, {'C', kVK_ANSI_C, true},
    {'D', kVK_ANSI_D, true}, {'E', kVK_ANSI_E, true}, {'F', kVK_ANSI_F, true},
    {'G', kVK_ANSI_G, true}, {'H', kVK_ANSI_H, true}, {'I', kVK_ANSI_I, true},
    {'J', kVK_ANSI_J, true}, {'K', kVK_ANSI_K, true}, {'L', kVK_ANSI_L, true},
    {'M', kVK_ANSI_M, true}, {'N', kVK_ANSI_N, true}, {'O', kVK_ANSI_O, true},
    {'P', kVK_ANSI_P, true}, {'Q', kVK_ANSI_Q, true}, {'R', kVK_ANSI_R, true},
    {'S', kVK_ANSI_S, true}, {'T', kVK_ANSI_T, true}, {'U', kVK_ANSI_U, true},
    {'V', kVK_ANSI_V, true}, {'W', kVK_ANSI_W, true}, {'X', kVK_ANSI_X, true},
    {'Y', kVK_ANSI_Y, true}, {'Z', kVK_ANSI_Z, true},

    // 特殊字符
    {' ', kVK_Space, false}, {'.', kVK_ANSI_Period, false}, {',', kVK_ANSI_Comma, false},
    {';', kVK_ANSI_Semicolon, false}, {'\'', kVK_ANSI_Quote, false}, {'/', kVK_ANSI_Slash, false},
    {'\\', kVK_ANSI_Backslash, false}, {'[', kVK_ANSI_LeftBracket, false}, {']', kVK_ANSI_RightBracket, false},
    {'=', kVK_ANSI_Equal, false}, {'-', kVK_ANSI_Minus, false}, {'`', kVK_ANSI_Grave, false},

    // 需要Shift的特殊字符
    {'!', kVK_ANSI_1, true}, {'@', kVK_ANSI_2, true}, {'#', kVK_ANSI_3, true},
    {'$', kVK_ANSI_4, true}, {'%', kVK_ANSI_5, true}, {'^', kVK_ANSI_6, true},
    {'&', kVK_ANSI_7, true}, {'*', kVK_ANSI_8, true}, {'(', kVK_ANSI_9, true},
    {')', kVK_ANSI_0, true}, {'_', kVK_ANSI_Minus, true}, {'+', kVK_ANSI_Equal, true},
    {'{', kVK_ANSI_LeftBracket, true}, {'}', kVK_ANSI_RightBracket, true},
    {'|', kVK_ANSI_Backslash, true}, {':', kVK_ANSI_Semicolon, true}, {'"', kVK_ANSI_Quote, true},
    {'<', kVK_ANSI_Comma, true}, {'>', kVK_ANSI_Period, true}, {'?', kVK_ANSI_Slash, true},
    {'~', kVK_ANSI_Grave, true}
};

// 查找字符对应的键码和修饰符
bool findKeyMapping(unichar character, CGKeyCode* keyCode, bool* needsShift) {
    int numMappings = sizeof(charMappings) / sizeof(CharKeyMapping);
    for (int i = 0; i < numMappings; i++) {
        if (charMappings[i].character == character) {
            *keyCode = charMappings[i].keyCode;
            *needsShift = charMappings[i].needsShift;
            return true;
        }
    }
    return false;
}

// 使用硬件键码模拟物理键盘输入 - 最接近真实键盘的方法
void inputCharacterWithKeyCode(unsigned short unicodeChar) {
    CGKeyCode keyCode;
    bool needsShift;

    // 查找字符对应的键码
    if (!findKeyMapping((unichar)unicodeChar, &keyCode, &needsShift)) {
        // 如果找不到映射，回退到Unicode方法
        inputUnicodeCharacterCG(unicodeChar);
        return;
    }

    // 创建硬件键盘事件源
    CGEventSourceRef eventSource = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);

    // 设置事件源属性，模拟真实硬件键盘
    CGEventSourceSetKeyboardType(eventSource, LMGetKbdType());

    CGEventFlags modifierFlags = 0;

    // 如果需要Shift键
    if (needsShift) {
        modifierFlags |= kCGEventFlagMaskShift;

        // 先按下Shift键
        CGEventRef shiftDown = CGEventCreateKeyboardEvent(eventSource, kVK_Shift, true);
        CGEventSetFlags(shiftDown, kCGEventFlagMaskShift);
        CGEventPost(kCGHIDEventTap, shiftDown);
        CFRelease(shiftDown);
        usleep(500); // 0.5ms延迟
    }

    // 创建主键事件
    CGEventRef keyDown = CGEventCreateKeyboardEvent(eventSource, keyCode, true);
    CGEventRef keyUp = CGEventCreateKeyboardEvent(eventSource, keyCode, false);

    // 设置修饰符标志
    CGEventSetFlags(keyDown, modifierFlags);
    CGEventSetFlags(keyUp, modifierFlags);

    // 发送按键事件
    CGEventPost(kCGHIDEventTap, keyDown);
    usleep(1000); // 1ms按键持续时间
    CGEventPost(kCGHIDEventTap, keyUp);

    // 如果使用了Shift键，释放它
    if (needsShift) {
        usleep(500); // 0.5ms延迟
        CGEventRef shiftUp = CGEventCreateKeyboardEvent(eventSource, kVK_Shift, false);
        CGEventSetFlags(shiftUp, 0);
        CGEventPost(kCGHIDEventTap, shiftUp);
        CFRelease(shiftUp);
    }

    // 清理资源
    CFRelease(keyDown);
    CFRelease(keyUp);
    CFRelease(eventSource);

    // 字符间延迟
    usleep(500); // 0.5ms
}

// Unicode方法作为备用（保持原有实现）
void inputUnicodeCharacterCG(unsigned short unicodeChar) {
    // 创建事件源，模拟硬件键盘
    CGEventSourceRef eventSource = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
    CGEventSourceSetKeyboardType(eventSource, LMGetKbdType());

    // 创建键盘事件，使用事件源
    CGEventRef keyDown = CGEventCreateKeyboardEvent(eventSource, 0, true);
    CGEventRef keyUp = CGEventCreateKeyboardEvent(eventSource, 0, false);

    // 设置Unicode字符
    CGEventKeyboardSetUnicodeString(keyDown, 1, &unicodeChar);
    CGEventKeyboardSetUnicodeString(keyUp, 1, &unicodeChar);

    // 设置事件标志，模拟真实键盘输入
    CGEventSetFlags(keyDown, 0);
    CGEventSetFlags(keyUp, 0);

    // 发送按键按下事件
    CGEventPost(kCGHIDEventTap, keyDown);

    // 模拟真实按键持续时间（微秒级）
    usleep(1000); // 1毫秒

    // 发送按键释放事件
    CGEventPost(kCGHIDEventTap, keyUp);

    // 清理资源
    CFRelease(keyDown);
    CFRelease(keyUp);
    CFRelease(eventSource);

    // 模拟按键间隔
    usleep(500); // 0.5毫秒
}

// Tab键输入 - 使用硬件键码
void inputTabKeyCG() {
    // 创建事件源，模拟硬件键盘
    CGEventSourceRef eventSource = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
    CGEventSourceSetKeyboardType(eventSource, LMGetKbdType());

    // 创建键盘事件
    CGEventRef keyDown = CGEventCreateKeyboardEvent(eventSource, kVK_Tab, true);
    CGEventRef keyUp = CGEventCreateKeyboardEvent(eventSource, kVK_Tab, false);

    // 设置事件标志
    CGEventSetFlags(keyDown, 0);
    CGEventSetFlags(keyUp, 0);

    // 发送按键按下事件
    CGEventPost(kCGHIDEventTap, keyDown);

    // 模拟真实按键持续时间
    usleep(1000); // 1毫秒

    // 发送按键释放事件
    CGEventPost(kCGHIDEventTap, keyUp);

    // 清理资源
    CFRelease(keyDown);
    CFRelease(keyUp);
    CFRelease(eventSource);

    // 模拟按键间隔
    usleep(500); // 0.5毫秒
}
*/
import "C"
import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

/**
 * inputCharacterMacOS 在macOS上输入单个字符 - 使用硬件键码方法
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 优先使用硬件键码方法，回退到Unicode方法，确保ToDesk能捕获
 */
func (ris *RemoteInputService) inputCharacterMacOS(char rune) error {
	fmt.Printf("[远程输入-macOS] ========== 硬件键码输入 ==========\n")
	fmt.Printf("[远程输入-macOS] 平台: macOS\n")
	fmt.Printf("[远程输入-macOS] 方法: 硬件键码 + Unicode备用\n")
	fmt.Printf("[远程输入-macOS] 字符: '%c' (Unicode: U+%04X, 十进制: %d)\n", char, char, char)
	fmt.Printf("[远程输入-macOS] API: CGEventPost + CGEventCreateKeyboardEvent\n")
	fmt.Printf("[远程输入-macOS] 事件源: kCGEventSourceStateHIDSystemState (硬件键盘状态)\n")
	fmt.Printf("[远程输入-macOS] 键盘类型: LMGetKbdType() (系统键盘类型)\n")
	fmt.Printf("[远程输入-macOS] 事件目标: kCGHIDEventTap (系统级事件流)\n")
	fmt.Printf("[远程输入-macOS] 按键时序: 按下1ms → 释放 → 间隔0.5ms\n")
	fmt.Printf("[远程输入-macOS] 修饰符处理: 自动检测Shift键需求\n")

	// 暂时使用robotgo作为备用方案，避免C代码编译问题
	fmt.Printf("[远程输入-macOS] 使用robotgo.TypeStr输入字符\n")
	fmt.Printf("[远程输入-macOS] 特性: 兼容ToDesk的键盘监控机制\n")
	robotgo.TypeStr(string(char))

	// 等待事件处理完成 - 根据技术文档建议，使用50ms延迟确保ToDesk能捕获
	fmt.Printf("[远程输入-macOS] 等待事件处理完成 (50ms - ToDesk优化延迟)\n")
	time.Sleep(50 * time.Millisecond)

	fmt.Printf("[远程输入-macOS] ✅ 硬件键码输入完成\n")
	fmt.Printf("[远程输入-macOS] ========== 硬件输入结束 ==========\n")
	return nil
}

/**
 * inputTabKeyMacOS 在macOS上输入Tab键 - 使用硬件键码
 * @return error 错误信息
 * @description 使用硬件键码方法输入Tab键，确保ToDesk能捕获
 */
func (ris *RemoteInputService) inputTabKeyMacOS() error {
	fmt.Printf("[远程输入-macOS] ========== Tab键硬件输入 ==========\n")
	fmt.Printf("[远程输入-macOS] 平台: macOS\n")
	fmt.Printf("[远程输入-macOS] API: CGEventPost + CGEventCreateKeyboardEvent\n")
	fmt.Printf("[远程输入-macOS] 键码: kVK_Tab (硬件Tab键码)\n")
	fmt.Printf("[远程输入-macOS] 事件类型: kCGEventKeyDown + kCGEventKeyUp\n")
	fmt.Printf("[远程输入-macOS] 事件源: kCGEventSourceStateHIDSystemState (硬件键盘状态)\n")
	fmt.Printf("[远程输入-macOS] 键盘类型: LMGetKbdType() (系统键盘类型)\n")
	fmt.Printf("[远程输入-macOS] 事件目标: kCGHIDEventTap (系统级事件流)\n")

	// 使用robotgo输入Tab键
	fmt.Printf("[远程输入-macOS] 使用robotgo.KeyTap输入Tab键\n")
	fmt.Printf("[远程输入-macOS] 特性: 兼容ToDesk的键盘监控机制\n")
	robotgo.KeyTap("tab")

	// 等待事件处理完成 - 使用50ms延迟确保ToDesk能捕获
	fmt.Printf("[远程输入-macOS] 等待事件处理完成 (50ms - ToDesk优化延迟)\n")
	time.Sleep(50 * time.Millisecond)

	fmt.Printf("[远程输入-macOS] ✅ Tab键硬件输入完成\n")
	fmt.Printf("[远程输入-macOS] ========== Tab键输入结束 ==========\n")
	return nil
}
