//go:build darwin

package services

/**
 * newGlobalHotkeyPlatformImpl 创建macOS平台实现
 * @return GlobalHotkeyPlatformImpl 平台实现
 * @return error 错误信息
 */
func newGlobalHotkeyPlatformImpl() (GlobalHotkeyPlatformImpl, error) {
	return newGlobalHotkeyDarwin()
}
