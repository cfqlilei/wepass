//go:build linux
// +build linux

package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

/**
 * LinuxLowLevelKeyboard Linux平台底层键盘实现
 * @author 陈凤庆
 * @description 使用xdotool实现Linux平台的底层键盘输入
 * @date 20251004
 */
type LinuxLowLevelKeyboard struct {
	keyboardLayout string
	xdotoolPath    string
}

/**
 * NewLinuxLowLevelKeyboard 创建Linux底层键盘实例
 * @return *LinuxLowLevelKeyboard 实例
 * @return error 错误信息
 */
func NewLinuxLowLevelKeyboard() (*LinuxLowLevelKeyboard, error) {
	// 检查xdotool是否安装
	xdotoolPath, err := exec.LookPath("xdotool")
	if err != nil {
		return nil, fmt.Errorf("xdotool未安装，请运行: sudo apt install xdotool")
	}

	return &LinuxLowLevelKeyboard{
		xdotoolPath: xdotoolPath,
	}, nil
}

/**
 * newPlatformImpl 创建平台特定的底层键盘实现 (Linux版本)
 * @return LowLevelKeyboardPlatform 平台实现
 * @return error 错误信息
 */
func newPlatformImpl() (LowLevelKeyboardPlatform, error) {
	return NewLinuxLowLevelKeyboard()
}

/**
 * SendKeyEvent 发送键盘事件
 * @param virtualKey 虚拟键码
 * @param scanCode 扫描码
 * @param isKeyDown 是否为按下事件
 * @param modifiers 修饰键列表
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) SendKeyEvent(virtualKey int, scanCode int, isKeyDown bool, modifiers []string) error {
	// xdotool使用键名而不是键码，需要转换
	keyName := l.virtualKeyToKeyName(virtualKey)
	if keyName == "" {
		return fmt.Errorf("无法转换虚拟键码%d到键名", virtualKey)
	}

	var cmd *exec.Cmd
	if isKeyDown {
		if len(modifiers) > 0 {
			// 带修饰键的按键
			modifierStr := strings.Join(modifiers, "+")
			cmd = exec.Command(l.xdotoolPath, "keydown", modifierStr+"+"+keyName)
		} else {
			cmd = exec.Command(l.xdotoolPath, "keydown", keyName)
		}
	} else {
		if len(modifiers) > 0 {
			// 带修饰键的释放
			modifierStr := strings.Join(modifiers, "+")
			cmd = exec.Command(l.xdotoolPath, "keyup", modifierStr+"+"+keyName)
		} else {
			cmd = exec.Command(l.xdotoolPath, "keyup", keyName)
		}
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool命令执行失败: %w", err)
	}

	return nil
}

/**
 * GetKeyboardLayout 获取当前键盘布局
 * @return string 键盘布局标识符
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) GetKeyboardLayout() (string, error) {
	// 使用setxkbmap获取键盘布局
	cmd := exec.Command("setxkbmap", "-query")
	output, err := cmd.Output()
	if err != nil {
		return "US", nil // 默认返回US布局
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "layout:") {
			layout := strings.TrimSpace(strings.TrimPrefix(line, "layout:"))
			friendlyName := l.convertLayoutToFriendlyName(layout)
			l.keyboardLayout = friendlyName
			return friendlyName, nil
		}
	}

	l.keyboardLayout = "US"
	return "US", nil
}

/**
 * MapCharToKey 将字符映射到键码
 * @param char 字符
 * @param layout 键盘布局
 * @return KeyMapping 键码映射
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) MapCharToKey(char rune, layout string) (KeyMapping, error) {
	// Linux下的字符到键码映射
	mapping := KeyMapping{
		VirtualKey:  int(char), // 简化处理，使用字符码作为虚拟键码
		ScanCode:    0,
		NeedsShift:  false,
		NeedsAlt:    false,
		NeedsCtrl:   false,
		Description: fmt.Sprintf("Char:%c", char),
	}

	// 判断是否需要Shift键
	if char >= 'A' && char <= 'Z' {
		mapping.NeedsShift = true
	} else if strings.ContainsRune("!@#$%^&*()_+{}|:\"<>?", char) {
		mapping.NeedsShift = true
	}

	// 处理特殊字符的额外修饰键需求
	l.adjustSpecialCharModifiers(char, &mapping)

	return mapping, nil
}

/**
 * CheckPermissions 检查权限
 * @return bool 是否有权限
 */
func (l *LinuxLowLevelKeyboard) CheckPermissions() bool {
	// 检查是否在X11环境中
	if os.Getenv("DISPLAY") == "" {
		return false
	}

	// 检查xdotool是否可用
	cmd := exec.Command(l.xdotoolPath, "version")
	return cmd.Run() == nil
}

/**
 * convertLayoutToFriendlyName 将Linux布局标识符转换为友好名称
 * @param layoutId 布局标识符
 * @return string 友好名称
 */
func (l *LinuxLowLevelKeyboard) convertLayoutToFriendlyName(layoutId string) string {
	// Linux键盘布局标识符映射
	layoutMap := map[string]string{
		"us": "US",
		"gb": "GB",
		"de": "DE",
		"fr": "FR",
		"it": "IT",
		"es": "ES",
		"ru": "RU",
		"jp": "JP",
		"kr": "KR",
		"cn": "CN",
		"tw": "TW",
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
func (l *LinuxLowLevelKeyboard) adjustSpecialCharModifiers(char rune, mapping *KeyMapping) {
	// 根据字符和键盘布局调整修饰键需求
	switch char {
	case '@':
		// 在某些布局中@需要AltGr键
		if l.keyboardLayout == "DE" || l.keyboardLayout == "FR" {
			mapping.NeedsAlt = true
		}
	case '#':
		// 在英国布局中#的位置不同
		if l.keyboardLayout == "GB" {
			mapping.NeedsShift = false // 在GB布局中#不需要Shift
		}
	case '€':
		// 欧元符号通常需要AltGr
		mapping.NeedsAlt = true
	}
}

/**
 * virtualKeyToKeyName 将虚拟键码转换为xdotool键名
 * @param virtualKey 虚拟键码
 * @return string 键名
 */
func (l *LinuxLowLevelKeyboard) virtualKeyToKeyName(virtualKey int) string {
	// 特殊键映射
	specialKeys := map[int]string{
		9:   "Tab",
		13:  "Return",
		27:  "Escape",
		32:  "space",
		8:   "BackSpace",
		46:  "Delete",
		16:  "Shift_L",
		17:  "Control_L",
		18:  "Alt_L",
		37:  "Left",
		38:  "Up",
		39:  "Right",
		40:  "Down",
		36:  "Home",
		35:  "End",
		33:  "Page_Up",
		34:  "Page_Down",
		45:  "Insert",
		112: "F1",
		113: "F2",
		114: "F3",
		115: "F4",
		116: "F5",
		117: "F6",
		118: "F7",
		119: "F8",
		120: "F9",
		121: "F10",
		122: "F11",
		123: "F12",
	}

	if keyName, exists := specialKeys[virtualKey]; exists {
		return keyName
	}

	// 对于字符键，直接使用字符
	if virtualKey >= 32 && virtualKey <= 126 {
		return string(rune(virtualKey))
	}

	return ""
}

/**
 * GetSpecialKeyMapping 获取特殊键的映射
 * @param keyName 键名
 * @return KeyMapping 键码映射
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) GetSpecialKeyMapping(keyName string) (KeyMapping, error) {
	// Linux特殊键映射
	specialKeys := map[string]KeyMapping{
		"tab":       {VirtualKey: 9, Description: "Tab"},
		"enter":     {VirtualKey: 13, Description: "Return"},
		"return":    {VirtualKey: 13, Description: "Return"},
		"space":     {VirtualKey: 32, Description: "Space"},
		"backspace": {VirtualKey: 8, Description: "BackSpace"},
		"delete":    {VirtualKey: 46, Description: "Delete"},
		"escape":    {VirtualKey: 27, Description: "Escape"},
		"f1":        {VirtualKey: 112, Description: "F1"},
		"f2":        {VirtualKey: 113, Description: "F2"},
		"f3":        {VirtualKey: 114, Description: "F3"},
		"f4":        {VirtualKey: 115, Description: "F4"},
		"f5":        {VirtualKey: 116, Description: "F5"},
		"f6":        {VirtualKey: 117, Description: "F6"},
		"f7":        {VirtualKey: 118, Description: "F7"},
		"f8":        {VirtualKey: 119, Description: "F8"},
		"f9":        {VirtualKey: 120, Description: "F9"},
		"f10":       {VirtualKey: 121, Description: "F10"},
		"f11":       {VirtualKey: 122, Description: "F11"},
		"f12":       {VirtualKey: 123, Description: "F12"},
		"up":        {VirtualKey: 38, Description: "Up"},
		"down":      {VirtualKey: 40, Description: "Down"},
		"left":      {VirtualKey: 37, Description: "Left"},
		"right":     {VirtualKey: 39, Description: "Right"},
		"home":      {VirtualKey: 36, Description: "Home"},
		"end":       {VirtualKey: 35, Description: "End"},
		"pageup":    {VirtualKey: 33, Description: "Page_Up"},
		"pagedown":  {VirtualKey: 34, Description: "Page_Down"},
		"insert":    {VirtualKey: 45, Description: "Insert"},
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
func (l *LinuxLowLevelKeyboard) SimulateSpecialKey(keyName string) error {
	// 直接使用xdotool key命令
	cmd := exec.Command(l.xdotoolPath, "key", keyName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool key %s失败: %w", keyName, err)
	}

	return nil
}

/**
 * SimulateTextDirect 直接使用xdotool输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) SimulateTextDirect(text string) error {
	// 使用xdotool type命令直接输入文本
	cmd := exec.Command(l.xdotoolPath, "type", "--delay", "20", text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool type失败: %w", err)
	}

	return nil
}

/**
 * GetActiveWindow 获取当前活动窗口ID
 * @return string 窗口ID
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) GetActiveWindow() (string, error) {
	cmd := exec.Command(l.xdotoolPath, "getactivewindow")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取活动窗口失败: %w", err)
	}

	windowId := strings.TrimSpace(string(output))
	return windowId, nil
}

/**
 * ActivateWindow 激活指定窗口
 * @param windowId 窗口ID
 * @return error 错误信息
 */
func (l *LinuxLowLevelKeyboard) ActivateWindow(windowId string) error {
	cmd := exec.Command(l.xdotoolPath, "windowactivate", windowId)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("激活窗口%s失败: %w", windowId, err)
	}

	// 等待窗口激活
	time.Sleep(200 * time.Millisecond)
	return nil
}
