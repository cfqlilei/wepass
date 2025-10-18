//go:build linux
// +build linux

package services

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/**
 * LinuxKeyboardHelper Linux平台键盘助手实现
 * @author 陈凤庆
 * @description 使用xdotool实现Linux平台的键盘助手功能
 * @date 20251005
 */
type LinuxKeyboardHelper struct {
	xdotoolPath string // xdotool路径
	helperCmd   *exec.Cmd // 键盘助手命令
}

/**
 * NewLinuxKeyboardHelper 创建Linux键盘助手实例
 * @return *LinuxKeyboardHelper 实例
 * @return error 错误信息
 */
func NewLinuxKeyboardHelper() (*LinuxKeyboardHelper, error) {
	// 检查xdotool是否安装
	xdotoolPath, err := exec.LookPath("xdotool")
	if err != nil {
		return nil, fmt.Errorf("xdotool未安装，请运行: sudo apt install xdotool")
	}

	return &LinuxKeyboardHelper{
		xdotoolPath: xdotoolPath,
		helperCmd:   nil,
	}, nil
}

/**
 * newKeyboardHelperPlatformImpl 创建平台特定的键盘助手实现 (Linux版本)
 * @return KeyboardHelperPlatform 平台实现
 * @return error 错误信息
 */
func newKeyboardHelperPlatformImpl() (KeyboardHelperPlatform, error) {
	return NewLinuxKeyboardHelper()
}

/**
 * LaunchHelper 启动键盘助手程序
 * @return int 键盘助手进程PID
 * @return error 错误信息
 */
func (l *LinuxKeyboardHelper) LaunchHelper() (int, error) {
	fmt.Printf("[Linux键盘助手] 启动键盘助手（使用xdotool）\n")

	// Linux使用xdotool进行键盘输入，不需要启动独立的键盘助手程序
	// 这里返回一个虚拟PID，表示服务已就绪
	cmd := exec.Command("sh", "-c", "echo $$")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("获取进程PID失败: %w", err)
	}

	pidStr := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("解析PID失败: %w", err)
	}

	fmt.Printf("[Linux键盘助手] ✅ 键盘助手已就绪，虚拟PID: %d\n", pid)
	return pid, nil
}

/**
 * CloseHelper 关闭键盘助手程序
 * @param helperPID 键盘助手进程PID
 * @return error 错误信息
 */
func (l *LinuxKeyboardHelper) CloseHelper(helperPID int) error {
	fmt.Printf("[Linux键盘助手] 关闭键盘助手（无需操作）\n")
	// Linux使用xdotool，不需要关闭独立进程
	return nil
}

/**
 * SimulateCharWithHelper 使用键盘助手输入单个字符
 * @param char 要输入的字符
 * @param helperPID 键盘助手进程PID
 * @return error 错误信息
 */
func (l *LinuxKeyboardHelper) SimulateCharWithHelper(char rune, helperPID int) error {
	// 处理Tab键
	if char == '\t' {
		cmd := exec.Command(l.xdotoolPath, "key", "Tab")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("xdotool输入Tab键失败: %w", err)
		}
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	// 使用xdotool type命令输入字符
	cmd := exec.Command(l.xdotoolPath, "type", "--delay", "20", string(char))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xdotool输入字符失败: %w", err)
	}

	time.Sleep(50 * time.Millisecond)
	return nil
}

/**
 * CheckHelperAvailable 检查键盘助手是否可用
 * @return bool 是否可用
 */
func (l *LinuxKeyboardHelper) CheckHelperAvailable() bool {
	// 检查xdotool是否可用
	_, err := exec.LookPath("xdotool")
	if err != nil {
		return false
	}

	// 检查是否在X11环境中
	cmd := exec.Command(l.xdotoolPath, "version")
	return cmd.Run() == nil
}

/**
 * GetHelperName 获取键盘助手名称
 * @return string 键盘助手名称
 */
func (l *LinuxKeyboardHelper) GetHelperName() string {
	return "Linux键盘助手 (xdotool)"
}

/**
 * checkProcessExists 检查进程是否存在
 * @param pid 进程PID
 * @return bool 是否存在
 */
func (l *LinuxKeyboardHelper) checkProcessExists(pid int) bool {
	// 使用kill -0检查进程是否存在
	process, err := syscall.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

/**
 * terminateProcess 终止进程
 * @param pid 进程PID
 * @return error 错误信息
 */
func (l *LinuxKeyboardHelper) terminateProcess(pid int) error {
	process, err := syscall.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("查找进程失败: %w", err)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("终止进程失败: %w", err)
	}

	return nil
}

/**
 * LaunchOnboardHelper 启动Onboard屏幕键盘（可选实现）
 * @return int 进程PID
 * @return error 错误信息
 */
func (l *LinuxKeyboardHelper) LaunchOnboardHelper() (int, error) {
	// 检查onboard是否安装
	onboardPath, err := exec.LookPath("onboard")
	if err != nil {
		return 0, fmt.Errorf("onboard未安装，请运行: sudo apt install onboard")
	}

	// 启动onboard
	cmd := exec.Command(onboardPath)
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("启动onboard失败: %w", err)
	}

	pid := cmd.Process.Pid
	l.helperCmd = cmd

	fmt.Printf("[Linux键盘助手] ✅ Onboard已启动，PID: %d\n", pid)
	return pid, nil
}

/**
 * CloseOnboardHelper 关闭Onboard屏幕键盘（可选实现）
 * @return error 错误信息
 */
func (l *LinuxKeyboardHelper) CloseOnboardHelper() error {
	if l.helperCmd == nil || l.helperCmd.Process == nil {
		return nil
	}

	pid := l.helperCmd.Process.Pid
	fmt.Printf("[Linux键盘助手] 关闭Onboard，PID: %d\n", pid)

	if err := l.helperCmd.Process.Kill(); err != nil {
		return fmt.Errorf("关闭Onboard失败: %w", err)
	}

	l.helperCmd = nil
	return nil
}

