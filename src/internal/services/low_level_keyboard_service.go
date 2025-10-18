package services

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

/**
 * LowLevelKeyboardService 底层键盘API输入服务
 * @author 陈凤庆
 * @description 提供跨平台的底层键盘API输入功能，解决ToDesk等远程桌面工具中输入失效的问题
 * @date 20251004
 */
type LowLevelKeyboardService struct {
	keyboardLayout string                   // 当前键盘布局
	keyMap         map[rune]KeyMapping      // 字符到键码的映射
	platformImpl   LowLevelKeyboardPlatform // 平台特定实现
}

/**
 * KeyMapping 键码映射结构
 */
type KeyMapping struct {
	VirtualKey  int    // 虚拟键码
	ScanCode    int    // 扫描码
	NeedsShift  bool   // 是否需要Shift键
	NeedsAlt    bool   // 是否需要Alt键
	NeedsCtrl   bool   // 是否需要Ctrl键
	Description string // 键的描述
}

/**
 * LowLevelKeyboardPlatform 平台特定键盘接口
 */
type LowLevelKeyboardPlatform interface {
	// SendKeyEvent 发送键盘事件
	SendKeyEvent(virtualKey int, scanCode int, isKeyDown bool, modifiers []string) error

	// GetKeyboardLayout 获取当前键盘布局
	GetKeyboardLayout() (string, error)

	// MapCharToKey 将字符映射到键码
	MapCharToKey(char rune, layout string) (KeyMapping, error)

	// CheckPermissions 检查权限
	CheckPermissions() bool
}

/**
 * NewLowLevelKeyboardService 创建底层键盘服务实例
 * @return *LowLevelKeyboardService 服务实例
 * @return error 错误信息
 */
func NewLowLevelKeyboardService() (*LowLevelKeyboardService, error) {
	service := &LowLevelKeyboardService{
		keyMap: make(map[rune]KeyMapping),
	}

	// 根据平台创建对应的实现
	var err error
	service.platformImpl, err = newPlatformImpl()

	if err != nil {
		return nil, fmt.Errorf("创建平台实现失败: %w", err)
	}

	// 初始化键盘布局和键码映射
	if err := service.initializeKeyboardLayout(); err != nil {
		return nil, fmt.Errorf("初始化键盘布局失败: %w", err)
	}

	return service, nil
}

/**
 * SimulateText 模拟输入文本
 * @param text 要输入的文本
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) SimulateText(text string) error {
	if !s.platformImpl.CheckPermissions() {
		return errors.New("缺少必要的系统权限")
	}

	if text == "" {
		return errors.New("输入文本不能为空")
	}

	fmt.Printf("[底层键盘] 开始输入文本: %s (长度: %d)\n", text, len(text))

	// 检测文本与键盘布局的兼容性
	if err := s.validateTextCompatibility(text); err != nil {
		fmt.Printf("[底层键盘] ⚠️ 文本兼容性检查失败: %v\n", err)
		// 尝试自动调整键盘布局
		if adjustErr := s.adjustKeyboardLayout(text); adjustErr != nil {
			return fmt.Errorf("文本兼容性检查失败且无法自动调整: %w", err)
		}
	}

	// 逐字符输入
	for i, char := range text {
		if err := s.simulateChar(char); err != nil {
			return fmt.Errorf("输入第%d个字符'%c'失败: %w", i+1, char, err)
		}

		// 字符间延迟，模拟真实输入
		time.Sleep(20 * time.Millisecond)
	}

	fmt.Printf("[底层键盘] ✅ 文本输入完成\n")
	return nil
}

/**
 * simulateChar 模拟输入单个字符
 * @param char 要输入的字符
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) simulateChar(char rune) error {
	// 从缓存中获取键码映射
	keyMapping, exists := s.keyMap[char]
	if !exists {
		// 动态映射字符到键码
		var err error
		keyMapping, err = s.platformImpl.MapCharToKey(char, s.keyboardLayout)
		if err != nil {
			return fmt.Errorf("无法映射字符'%c'到键码: %w", char, err)
		}
		// 缓存映射结果
		s.keyMap[char] = keyMapping
	}

	// 构建修饰键列表
	var modifiers []string
	if keyMapping.NeedsShift {
		modifiers = append(modifiers, "shift")
	}
	if keyMapping.NeedsAlt {
		modifiers = append(modifiers, "alt")
	}
	if keyMapping.NeedsCtrl {
		modifiers = append(modifiers, "ctrl")
	}

	// 按下修饰键
	for _, modifier := range modifiers {
		if err := s.sendModifierKey(modifier, true); err != nil {
			return fmt.Errorf("按下修饰键%s失败: %w", modifier, err)
		}
	}

	// 发送字符键事件（按下+释放）
	if err := s.platformImpl.SendKeyEvent(keyMapping.VirtualKey, keyMapping.ScanCode, true, modifiers); err != nil {
		return fmt.Errorf("发送按键事件失败: %w", err)
	}

	time.Sleep(10 * time.Millisecond) // 按键持续时间

	if err := s.platformImpl.SendKeyEvent(keyMapping.VirtualKey, keyMapping.ScanCode, false, modifiers); err != nil {
		return fmt.Errorf("发送释放键事件失败: %w", err)
	}

	// 释放修饰键（逆序）
	for i := len(modifiers) - 1; i >= 0; i-- {
		if err := s.sendModifierKey(modifiers[i], false); err != nil {
			return fmt.Errorf("释放修饰键%s失败: %w", modifiers[i], err)
		}
	}

	return nil
}

/**
 * sendModifierKey 发送修饰键事件
 * @param modifier 修饰键名称
 * @param isKeyDown 是否为按下事件
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) sendModifierKey(modifier string, isKeyDown bool) error {
	// 修饰键的虚拟键码映射
	modifierKeys := map[string]int{
		"shift": 16, // VK_SHIFT
		"ctrl":  17, // VK_CONTROL
		"alt":   18, // VK_MENU
	}

	virtualKey, exists := modifierKeys[modifier]
	if !exists {
		return fmt.Errorf("未知的修饰键: %s", modifier)
	}

	return s.platformImpl.SendKeyEvent(virtualKey, 0, isKeyDown, nil)
}

/**
 * initializeKeyboardLayout 初始化键盘布局
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) initializeKeyboardLayout() error {
	layout, err := s.platformImpl.GetKeyboardLayout()
	if err != nil {
		return fmt.Errorf("获取键盘布局失败: %w", err)
	}

	s.keyboardLayout = layout
	fmt.Printf("[底层键盘] 当前键盘布局: %s\n", layout)

	// 预加载常用字符的键码映射
	if err := s.preloadCommonKeyMappings(); err != nil {
		return fmt.Errorf("预加载键码映射失败: %w", err)
	}

	return nil
}

/**
 * preloadCommonKeyMappings 预加载常用字符的键码映射
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) preloadCommonKeyMappings() error {
	// 常用字符列表
	commonChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;':\",./<>?"

	for _, char := range commonChars {
		keyMapping, err := s.platformImpl.MapCharToKey(char, s.keyboardLayout)
		if err != nil {
			fmt.Printf("[底层键盘] ⚠️ 预加载字符'%c'映射失败: %v\n", char, err)
			continue
		}
		s.keyMap[char] = keyMapping
	}

	fmt.Printf("[底层键盘] ✅ 预加载了%d个字符的键码映射\n", len(s.keyMap))
	return nil
}

/**
 * validateTextCompatibility 验证文本与键盘布局的兼容性
 * @param text 要验证的文本
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) validateTextCompatibility(text string) error {
	incompatibleChars := []rune{}

	for _, char := range text {
		if _, err := s.platformImpl.MapCharToKey(char, s.keyboardLayout); err != nil {
			incompatibleChars = append(incompatibleChars, char)
		}
	}

	if len(incompatibleChars) > 0 {
		return fmt.Errorf("发现%d个不兼容字符: %s", len(incompatibleChars), string(incompatibleChars))
	}

	return nil
}

/**
 * adjustKeyboardLayout 自动调整键盘布局
 * @param text 参考文本
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) adjustKeyboardLayout(text string) error {
	// 检测文本中的字符特征，推断最适合的键盘布局
	if strings.ContainsAny(text, "@#$%^&*") {
		// 包含特殊字符，可能需要US布局
		if s.keyboardLayout != "US" {
			fmt.Printf("[底层键盘] 🔄 检测到特殊字符，尝试切换到US键盘布局\n")
			s.keyboardLayout = "US"
			s.keyMap = make(map[rune]KeyMapping) // 清空缓存
			return s.preloadCommonKeyMappings()
		}
	}

	return errors.New("无法自动调整键盘布局")
}

/**
 * GetKeyboardLayout 获取当前键盘布局
 * @return string 键盘布局
 */
func (s *LowLevelKeyboardService) GetKeyboardLayout() string {
	return s.keyboardLayout
}

/**
 * SetKeyboardLayout 设置键盘布局
 * @param layout 键盘布局
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) SetKeyboardLayout(layout string) error {
	s.keyboardLayout = layout
	s.keyMap = make(map[rune]KeyMapping) // 清空缓存
	return s.preloadCommonKeyMappings()
}

/**
 * GetPlatformImpl 获取平台特定实现
 * @return LowLevelKeyboardPlatform 平台实现接口
 */
func (s *LowLevelKeyboardService) GetPlatformImpl() LowLevelKeyboardPlatform {
	return s.platformImpl
}

/**
 * SimulateSpecialKey 模拟特殊键输入
 * @param keyName 键名
 * @return error 错误信息
 */
func (s *LowLevelKeyboardService) SimulateSpecialKey(keyName string) error {
	// 检查权限
	if !s.platformImpl.CheckPermissions() {
		return errors.New("缺少必要的系统权限")
	}

	// 根据平台类型调用相应的特殊键方法
	switch impl := s.platformImpl.(type) {
	case interface{ SimulateSpecialKey(string) error }:
		return impl.SimulateSpecialKey(keyName)
	default:
		return fmt.Errorf("当前平台不支持特殊键'%s'的模拟", keyName)
	}
}
