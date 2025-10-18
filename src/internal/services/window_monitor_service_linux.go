//go:build linux
// +build linux

package services

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

/**
 * WindowInfo 窗口信息结构体
 * @author 陈凤庆
 * @date 20251002
 */
type WindowInfo struct {
	PID       int    `json:"pid"`       // 进程ID
	Title     string `json:"title"`     // 窗口标题
	Timestamp int64  `json:"timestamp"` // 记录时间戳
}

/**
 * WindowMonitorService Linux 版本的窗口监控服务
 * @author 陈凤庆
 * @date 20251002
 */
type WindowMonitorService struct {
	selfPID             int          // 密码管理器自身的PID
	lastPID             int          // 最后活动的非密码管理器窗口PID
	lastWindow          *WindowInfo  // 最后活动窗口的详细信息
	running             bool         // 监控服务是否运行中
	stopChan            chan bool    // 停止信号通道
	mu                  sync.RWMutex // 读写锁，保护共享数据
	pidFilePath         string       // PID文件存储路径
	mouseInSelfWindow   bool         // 鼠标是否在当前窗口内
	lastActiveWindowPID int          // 上次检查的活动窗口PID，用于检测窗口切换
}

/**
 * NewWindowMonitorService 创建新的窗口监控服务实例
 * @return *WindowMonitorService 窗口监控服务实例
 */
func NewWindowMonitorService() *WindowMonitorService {
	log.Println("[WindowMonitorService] 初始化 Linux 窗口监控服务...")

	service := &WindowMonitorService{
		selfPID:  os.Getpid(),
		stopChan: make(chan bool, 1),
		running:  false,
	}

	log.Printf("[WindowMonitorService] Linux 窗口监控服务初始化完成，自身PID: %d", service.selfPID)
	return service
}

/**
 * StartMonitoring 开始窗口监控（Linux 版本暂不支持）
 * @return error 错误信息
 */
func (wms *WindowMonitorService) StartMonitoring() error {
	wms.mu.Lock()
	defer wms.mu.Unlock()

	if wms.running {
		return fmt.Errorf("窗口监控服务已在运行中")
	}

	log.Println("[WindowMonitorService] Linux 版本暂不支持窗口监控")
	wms.running = true
	return nil
}

/**
 * StopMonitoring 停止窗口监控
 * @return error 错误信息
 */
func (wms *WindowMonitorService) StopMonitoring() error {
	wms.mu.Lock()
	defer wms.mu.Unlock()

	if !wms.running {
		return nil
	}

	wms.running = false
	select {
	case wms.stopChan <- true:
	default:
	}

	log.Println("[WindowMonitorService] Linux 窗口监控服务已停止")
	return nil
}

/**
 * Stop 停止窗口监控服务
 * @return error 错误信息
 */
func (wms *WindowMonitorService) Stop() error {
	return wms.StopMonitoring()
}

/**
 * SwitchToLastWindow 切换到最后记录的窗口（Linux 版本暂不支持）
 * @return bool 是否成功
 */
func (wms *WindowMonitorService) SwitchToLastWindow() bool {
	log.Println("[WindowMonitorService] Linux 版本暂不支持切换到最后记录的窗口")
	return false
}

/**
 * SwitchToPasswordManager 切换到密码管理器窗口（Linux 版本暂不支持）
 * @return bool 是否成功
 */
func (wms *WindowMonitorService) SwitchToPasswordManager() bool {
	log.Println("[WindowMonitorService] Linux 版本暂不支持切换到密码管理器窗口")
	return false
}

/**
 * RecordCurrentWindow 记录当前窗口信息（Linux 版本暂不支持）
 * @return error 错误信息
 */
func (wms *WindowMonitorService) RecordCurrentWindow() error {
	wms.mu.Lock()
	defer wms.mu.Unlock()

	// Linux 版本暂时创建一个虚拟的窗口信息
	wms.lastWindow = &WindowInfo{
		PID:       0,
		Title:     "Linux 版本暂不支持",
		Timestamp: time.Now().Unix(),
	}

	log.Println("[WindowMonitorService] Linux 版本暂不支持记录当前窗口信息")
	return nil
}

/**
 * GetLastWindowInfo 获取最后记录的窗口信息
 * @return *WindowInfo 窗口信息
 */
func (wms *WindowMonitorService) GetLastWindowInfo() *WindowInfo {
	wms.mu.RLock()
	defer wms.mu.RUnlock()

	if wms.lastWindow == nil {
		return &WindowInfo{
			PID:       0,
			Title:     "Linux 版本暂不支持",
			Timestamp: time.Now().Unix(),
		}
	}

	return wms.lastWindow
}

/**
 * RecordCurrentAsLastPID 记录当前窗口为最后活动窗口（Linux 版本暂不支持）
 * @return bool 是否成功
 */
func (wms *WindowMonitorService) RecordCurrentAsLastPID() bool {
	log.Println("[WindowMonitorService] Linux 版本暂不支持记录当前窗口为最后活动窗口")
	return false
}

/**
 * GetLastWindow 获取最后活动的窗口信息
 * @return *WindowInfo 窗口信息
 */
func (wms *WindowMonitorService) GetLastWindow() *WindowInfo {
	return wms.GetLastWindowInfo()
}

/**
 * GetSelfPID 获取自身进程ID
 * @return int 进程ID
 */
func (wms *WindowMonitorService) GetSelfPID() int {
	return wms.selfPID
}

/**
 * IsRunning 检查监控服务是否运行中
 * @return bool 是否运行中
 */
func (wms *WindowMonitorService) IsRunning() bool {
	wms.mu.RLock()
	defer wms.mu.RUnlock()
	return wms.running
}
