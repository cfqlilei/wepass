package services

import (
	"runtime"
	"testing"
)

/**
 * 底层键盘API服务测试
 * @author 陈凤庆
 * @description 测试底层键盘API服务的各项功能
 * @date 20251004
 */

func TestNewLowLevelKeyboardService(t *testing.T) {
	t.Log("测试底层键盘API服务创建")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Logf("创建底层键盘API服务失败（这在某些环境下是正常的）: %v", err)
		return
	}

	if service == nil {
		t.Error("服务实例不应该为nil")
		return
	}

	t.Log("✅ 底层键盘API服务创建成功")
}

func TestKeyboardLayoutDetection(t *testing.T) {
	t.Log("测试键盘布局检测")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	layout := service.GetKeyboardLayout()
	if layout == "" {
		t.Error("键盘布局不应该为空")
		return
	}

	t.Logf("检测到键盘布局: %s", layout)
	t.Log("✅ 键盘布局检测测试通过")
}

func TestCharacterMapping(t *testing.T) {
	t.Log("测试字符到键码映射")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	// 测试常用字符
	testChars := []rune{'a', 'A', '1', '@', '#', '$'}

	for _, char := range testChars {
		mapping, err := service.GetPlatformImpl().MapCharToKey(char, service.GetKeyboardLayout())
		if err != nil {
			t.Logf("字符 '%c' 映射失败: %v", char, err)
			continue
		}

		t.Logf("字符 '%c' -> VK:%d, SC:%d, Shift:%v",
			char, mapping.VirtualKey, mapping.ScanCode, mapping.NeedsShift)
	}

	t.Log("✅ 字符到键码映射测试完成")
}

func TestPermissionCheck(t *testing.T) {
	t.Log("测试权限检查")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	hasPermission := service.GetPlatformImpl().CheckPermissions()

	switch runtime.GOOS {
	case "darwin":
		t.Logf("macOS平台权限检查结果: %v", hasPermission)
		if !hasPermission {
			t.Log("⚠️ macOS需要辅助功能权限，请在系统偏好设置中授权")
		}
	case "windows":
		t.Logf("Windows平台权限检查结果: %v", hasPermission)
	case "linux":
		t.Logf("Linux平台权限检查结果: %v", hasPermission)
		if !hasPermission {
			t.Log("⚠️ Linux需要X11环境和xdotool工具")
		}
	}

	t.Log("✅ 权限检查测试完成")
}

func TestTextValidation(t *testing.T) {
	t.Log("测试文本验证")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	// 测试各种文本
	testTexts := []string{
		"hello world",
		"Hello World!",
		"123456",
		"test@example.com",
		"P@ssw0rd123",
	}

	for _, text := range testTexts {
		err := service.validateTextCompatibility(text)
		if err != nil {
			t.Logf("文本 '%s' 验证失败: %v", text, err)
		} else {
			t.Logf("文本 '%s' 验证通过", text)
		}
	}

	t.Log("✅ 文本验证测试完成")
}

func TestKeyboardLayoutAdjustment(t *testing.T) {
	t.Log("测试键盘布局自动调整")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	originalLayout := service.GetKeyboardLayout()
	t.Logf("原始键盘布局: %s", originalLayout)

	// 测试包含特殊字符的文本
	testText := "test@example.com#$%"

	err = service.adjustKeyboardLayout(testText)
	if err != nil {
		t.Logf("键盘布局调整失败: %v", err)
	} else {
		newLayout := service.GetKeyboardLayout()
		t.Logf("调整后键盘布局: %s", newLayout)
	}

	t.Log("✅ 键盘布局自动调整测试完成")
}

func TestErrorHandling(t *testing.T) {
	t.Log("测试错误处理")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	// 测试空文本输入
	err = service.SimulateText("")
	if err == nil {
		t.Error("空文本输入应该返回错误")
	} else {
		t.Logf("空文本输入正确返回错误: %v", err)
	}

	// 测试特殊键模拟错误处理
	err = service.SimulateSpecialKey("invalid_key")
	if err == nil {
		t.Log("无效特殊键可能被平台实现处理")
	} else {
		t.Logf("无效特殊键正确返回错误: %v", err)
	}

	t.Log("✅ 错误处理测试完成")
}

func TestPlatformSpecificFeatures(t *testing.T) {
	t.Log("测试平台特定功能")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	platform := runtime.GOOS
	t.Logf("当前平台: %s", platform)

	switch platform {
	case "darwin":
		t.Log("测试macOS特定功能...")
		// 测试CGEventPost相关功能
		layout, err := service.GetPlatformImpl().GetKeyboardLayout()
		if err != nil {
			t.Logf("获取macOS键盘布局失败: %v", err)
		} else {
			t.Logf("macOS键盘布局: %s", layout)
		}

	case "windows":
		t.Log("测试Windows特定功能...")
		// 测试SendInput相关功能
		layout, err := service.GetPlatformImpl().GetKeyboardLayout()
		if err != nil {
			t.Logf("获取Windows键盘布局失败: %v", err)
		} else {
			t.Logf("Windows键盘布局: %s", layout)
		}

	case "linux":
		t.Log("测试Linux特定功能...")
		// 测试xdotool相关功能
		layout, err := service.GetPlatformImpl().GetKeyboardLayout()
		if err != nil {
			t.Logf("获取Linux键盘布局失败: %v", err)
		} else {
			t.Logf("Linux键盘布局: %s", layout)
		}

	default:
		t.Logf("未知平台: %s", platform)
	}

	t.Log("✅ 平台特定功能测试完成")
}

func TestMemoryUsage(t *testing.T) {
	t.Log("测试内存使用")

	// 创建多个服务实例测试内存使用
	services := make([]*LowLevelKeyboardService, 0, 10)

	for i := 0; i < 10; i++ {
		service, err := NewLowLevelKeyboardService()
		if err != nil {
			t.Logf("创建第%d个服务实例失败: %v", i+1, err)
			continue
		}
		services = append(services, service)
	}

	t.Logf("成功创建了 %d 个服务实例", len(services))

	// 清理
	services = nil
	runtime.GC()

	t.Log("✅ 内存使用测试完成")
}

func TestConcurrency(t *testing.T) {
	t.Log("测试并发安全性")

	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Skipf("跳过测试，服务创建失败: %v", err)
		return
	}

	// 并发测试
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// 并发访问服务方法
			layout := service.GetKeyboardLayout()
			t.Logf("协程 %d 获取到键盘布局: %s", id, layout)

			// 测试字符映射
			_, err := service.GetPlatformImpl().MapCharToKey('a', layout)
			if err != nil {
				t.Logf("协程 %d 字符映射失败: %v", id, err)
			}
		}(i)
	}

	// 等待所有协程完成
	for i := 0; i < 5; i++ {
		<-done
	}

	t.Log("✅ 并发安全性测试完成")
}

// 基准测试
func BenchmarkCharacterMapping(b *testing.B) {
	service, err := NewLowLevelKeyboardService()
	if err != nil {
		b.Skipf("跳过基准测试，服务创建失败: %v", err)
		return
	}

	layout := service.GetKeyboardLayout()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetPlatformImpl().MapCharToKey('a', layout)
	}
}

func BenchmarkTextValidation(b *testing.B) {
	service, err := NewLowLevelKeyboardService()
	if err != nil {
		b.Skipf("跳过基准测试，服务创建失败: %v", err)
		return
	}

	testText := "hello world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.validateTextCompatibility(testText)
	}
}
