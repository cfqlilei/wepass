package services

import (
	"testing"
)

/**
 * 测试底层键盘API权限和字符映射
 * @author 陈凤庆
 * @date 20251004
 */

func TestLowLevelKeyboardDebug(t *testing.T) {
	service, err := NewLowLevelKeyboardService()
	if err != nil {
		t.Fatalf("创建底层键盘服务失败: %v", err)
	}

	// 测试权限
	t.Run("检查权限", func(t *testing.T) {
		impl := service.GetPlatformImpl()
		hasPermission := impl.CheckPermissions()
		t.Logf("辅助功能权限状态: %v", hasPermission)

		if !hasPermission {
			t.Log("⚠️ 没有辅助功能权限，这可能是字符映射失败的原因")
			t.Log("请在系统偏好设置 -> 安全性与隐私 -> 隐私 -> 辅助功能中添加此应用")
		}
	})

	// 测试键盘布局检测
	t.Run("检测键盘布局", func(t *testing.T) {
		impl := service.GetPlatformImpl()
		layout, err := impl.GetKeyboardLayout()
		if err != nil {
			t.Errorf("获取键盘布局失败: %v", err)
		} else {
			t.Logf("当前键盘布局: %s", layout)
		}
	})

	// 测试字符映射
	t.Run("测试字符映射", func(t *testing.T) {
		impl := service.GetPlatformImpl()

		// 测试几个基本字符
		testChars := []rune{'a', 'A', '1', '@', ' '}

		for _, char := range testChars {
			mapping, err := impl.MapCharToKey(char, "US")
			if err != nil {
				t.Logf("字符 '%c' 映射失败: %v", char, err)
			} else {
				t.Logf("字符 '%c' -> VK:%d, Shift:%v", char, mapping.VirtualKey, mapping.NeedsShift)
			}
		}
	})

	// 测试文本模拟
	t.Run("测试文本模拟", func(t *testing.T) {
		// 只测试不实际发送键盘事件，这里只是测试映射过程
		t.Log("跳过实际键盘事件发送，仅测试字符映射过程")
		t.Log("文本模拟测试完成")
	})
}
