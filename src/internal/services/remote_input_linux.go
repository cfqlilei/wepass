//go:build linux

package services

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

/**
 * inputCharacterLinux 在Linux上输入单个字符
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 使用xdotool命令生成键盘事件，ToDesk的X11或Wayland监听机制会捕获
 */
func (ris *RemoteInputService) inputCharacterLinux(char rune) error {
	fmt.Printf("[远程输入-Linux] 使用xdotool输入字符: '%c' (Unicode: %d)\n", char, char)
	
	// 首先检查xdotool是否可用
	if !ris.isXdotoolAvailable() {
		return fmt.Errorf("xdotool未安装，请运行: sudo apt install xdotool")
	}
	
	// 使用xdotool type命令输入字符
	charStr := string(char)
	cmd := exec.Command("xdotool", "type", charStr)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool输入字符失败: %w", err)
	}
	
	// 等待事件处理完成
	time.Sleep(10 * time.Millisecond)
	
	fmt.Printf("[远程输入-Linux] ✅ 字符输入完成\n")
	return nil
}

/**
 * inputTabKeyLinux 在Linux上输入Tab键
 * @return error 错误信息
 * @description 使用xdotool命令输入Tab键
 */
func (ris *RemoteInputService) inputTabKeyLinux() error {
	fmt.Printf("[远程输入-Linux] 使用xdotool输入Tab键\n")
	
	// 首先检查xdotool是否可用
	if !ris.isXdotoolAvailable() {
		return fmt.Errorf("xdotool未安装，请运行: sudo apt install xdotool")
	}
	
	// 使用xdotool key命令输入Tab键
	cmd := exec.Command("xdotool", "key", "Tab")
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool输入Tab键失败: %w", err)
	}
	
	// 等待事件处理完成
	time.Sleep(10 * time.Millisecond)
	
	fmt.Printf("[远程输入-Linux] ✅ Tab键输入完成\n")
	return nil
}

/**
 * isXdotoolAvailable 检查xdotool是否可用
 * @return bool xdotool是否可用
 */
func (ris *RemoteInputService) isXdotoolAvailable() bool {
	cmd := exec.Command("which", "xdotool")
	err := cmd.Run()
	return err == nil
}

/**
 * inputCharacterLinuxXTest 使用XTest扩展输入字符（备用方案）
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 如果xdotool不可用，可以尝试使用XTest扩展
 */
func (ris *RemoteInputService) inputCharacterLinuxXTest(char rune) error {
	fmt.Printf("[远程输入-Linux] 使用XTest输入字符: '%c' (Unicode: %d)\n", char, char)
	
	// 这里可以实现XTest的调用，需要CGO和X11库
	// 为了简化，暂时返回未实现错误
	return fmt.Errorf("XTest实现暂未完成，请使用xdotool")
}

/**
 * inputUnicodeLinux 使用Unicode方式输入字符
 * @param char 要输入的字符
 * @return error 错误信息
 * @description 使用Ctrl+Shift+U组合键输入Unicode字符
 */
func (ris *RemoteInputService) inputUnicodeLinux(char rune) error {
	fmt.Printf("[远程输入-Linux] 使用Unicode方式输入字符: '%c' (Unicode: %d)\n", char, char)
	
	if !ris.isXdotoolAvailable() {
		return fmt.Errorf("xdotool未安装，请运行: sudo apt install xdotool")
	}
	
	// 转换为十六进制Unicode
	unicodeHex := fmt.Sprintf("%x", char)
	
	// 按下Ctrl+Shift+U
	cmd1 := exec.Command("xdotool", "keydown", "ctrl+shift+u")
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("按下Ctrl+Shift+U失败: %w", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	
	// 输入Unicode十六进制值
	cmd2 := exec.Command("xdotool", "type", unicodeHex)
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("输入Unicode值失败: %w", err)
	}
	
	time.Sleep(50 * time.Millisecond)
	
	// 释放Ctrl+Shift+U并按Enter确认
	cmd3 := exec.Command("xdotool", "keyup", "ctrl+shift+u", "key", "Return")
	if err := cmd3.Run(); err != nil {
		return fmt.Errorf("确认Unicode输入失败: %w", err)
	}
	
	// 等待事件处理完成
	time.Sleep(10 * time.Millisecond)
	
	fmt.Printf("[远程输入-Linux] ✅ Unicode字符输入完成\n")
	return nil
}
