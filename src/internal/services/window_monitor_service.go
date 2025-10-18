//go:build darwin
// +build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework Foundation -framework AppKit
#include <ApplicationServices/ApplicationServices.h>
#include <Foundation/Foundation.h>
#include <AppKit/AppKit.h>
#include <unistd.h>
#include <stdlib.h>

// 获取当前进程PID
int getCurrentProcessPID() {
    return getpid();
}

// 获取当前活动窗口的PID
int getCurrentActiveWindowPID() {
    @autoreleasepool {
        NSRunningApplication *frontApp = [[NSWorkspace sharedWorkspace] frontmostApplication];
        if (frontApp) {
            return (int)[frontApp processIdentifier];
        }
        return 0;
    }
}

// 获取当前活动窗口的标题
const char* getCurrentActiveWindowTitle() {
    @autoreleasepool {
        NSRunningApplication *frontApp = [[NSWorkspace sharedWorkspace] frontmostApplication];
        if (frontApp) {
            NSString *appName = [frontApp localizedName];
            if (appName) {
                const char *utf8String = [appName UTF8String];
                return strdup(utf8String);
            }
        }
        return strdup("未知");
    }
}

// 激活指定PID的窗口
bool activateWindowByPID(int pid) {
    @autoreleasepool {
        NSArray *runningApps = [[NSWorkspace sharedWorkspace] runningApplications];
        for (NSRunningApplication *app in runningApps) {
            if ([app processIdentifier] == pid) {
                return [app activateWithOptions:NSApplicationActivateIgnoringOtherApps];
            }
        }
        return false;
    }
}

// 检查PID是否存在
bool checkPIDExists(int pid) {
    @autoreleasepool {
        NSArray *runningApps = [[NSWorkspace sharedWorkspace] runningApplications];
        for (NSRunningApplication *app in runningApps) {
            if ([app processIdentifier] == pid) {
                return true;
            }
        }
        return false;
    }
}

// 窗口边界结构体
typedef struct {
    int x;
    int y;
    int width;
    int height;
    bool found;
} WindowRect;

// 获取指定PID进程的主窗口位置和大小
WindowRect getWindowRectByPID(int pid) {
    WindowRect rect = {0, 0, 0, 0, false};
    @autoreleasepool {
        CFArrayRef windowList = CGWindowListCopyWindowInfo(kCGWindowListOptionOnScreenOnly, kCGNullWindowID);
        if (windowList) {
            CFIndex count = CFArrayGetCount(windowList);
            for (CFIndex i = 0; i < count; i++) {
                CFDictionaryRef window = (CFDictionaryRef)CFArrayGetValueAtIndex(windowList, i);

                CFNumberRef pidRef = (CFNumberRef)CFDictionaryGetValue(window, kCGWindowOwnerPID);
                int windowPID = 0;
                if (pidRef) {
                    CFNumberGetValue(pidRef, kCFNumberIntType, &windowPID);
                }

                if (windowPID == pid) {
                    CFDictionaryRef bounds = (CFDictionaryRef)CFDictionaryGetValue(window, kCGWindowBounds);
                    if (bounds) {
                        CGRect windowBounds;
                        if (CGRectMakeWithDictionaryRepresentation(bounds, &windowBounds)) {
                            rect.x = (int)windowBounds.origin.x;
                            rect.y = (int)windowBounds.origin.y;
                            rect.width = (int)windowBounds.size.width;
                            rect.height = (int)windowBounds.size.height;
                            rect.found = true;
                            break;
                        }
                    }
                }
            }
            CFRelease(windowList);
        }
    }
    return rect;
}

// 获取当前鼠标位置
typedef struct {
    int x;
    int y;
} MousePosition;

MousePosition getCurrentMousePosition() {
    MousePosition pos = {0, 0};
    @autoreleasepool {
        NSPoint mouseLocation = [NSEvent mouseLocation];
        NSScreen *mainScreen = [NSScreen mainScreen];
        if (mainScreen) {
            NSRect screenFrame = [mainScreen frame];
            pos.x = (int)mouseLocation.x;
            pos.y = (int)(screenFrame.size.height - mouseLocation.y);
        }
    }
    return pos;
}

// 检查鼠标是否在指定窗口区域内
bool isMouseInWindowRect(int windowX, int windowY, int windowWidth, int windowHeight) {
    MousePosition mousePos = getCurrentMousePosition();
    return (mousePos.x >= windowX && mousePos.x <= windowX + windowWidth &&
            mousePos.y >= windowY && mousePos.y <= windowY + windowHeight);
}

// 检查鼠标左键是否按下
bool isMouseButtonPressed() {
    @autoreleasepool {
        NSUInteger pressedButtons = [NSEvent pressedMouseButtons];
        return (pressedButtons & 1) != 0; // 检查左键是否按下
    }
}
*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

/**
 * WindowInfo 窗口信息结构
 * @author 陈凤庆
 * @description 存储窗口的基本信息
 */
type WindowInfo struct {
	PID       int    `json:"pid"`       // 进程ID
	Title     string `json:"title"`     // 窗口标题
	Timestamp int64  `json:"timestamp"` // 记录时间戳
}

/**
 * WindowMonitorService 窗口监控服务
 * @author 陈凤庆
 * @description 提供后台窗口监控、PID管理、自动切换活动窗口等功能
 */
type WindowMonitorService struct {
	selfPID             int           // 密码管理器自身的PID
	lastPID             int           // 最后活动的非密码管理器窗口PID
	lastWindow          *WindowInfo   // 最后活动窗口的详细信息
	running             bool          // 监控服务是否运行中
	stopChan            chan struct{} // 停止监控的信号通道
	mu                  sync.RWMutex  // 读写锁，保护共享数据
	pidFilePath         string        // PID文件存储路径
	mouseInSelfWindow   bool          // 20251002 陈凤庆 鼠标是否在当前窗口内
	lastActiveWindowPID int           // 20251002 陈凤庆 上次检查的活动窗口PID，用于检测窗口切换
}

/**
 * NewWindowMonitorService 创建新的窗口监控服务实例
 * @return *WindowMonitorService 窗口监控服务实例
 * @description 20251002 陈凤庆 初始化窗口监控服务，获取自身PID并设置文件路径
 */
func NewWindowMonitorService() *WindowMonitorService {
	// 20251002 陈凤庆 获取当前进程PID（密码管理器的PID）
	selfPID := int(C.getCurrentProcessPID())

	// 20251002 陈凤庆 设置PID文件存储路径
	pidFilePath := filepath.Join("logs", "lastPID.txt")

	fmt.Printf("[窗口监控] 初始化 WindowMonitorService，selfPID = %d（密码管理器进程）\n", selfPID)

	service := &WindowMonitorService{
		selfPID:     selfPID,
		lastPID:     0,
		running:     false,
		stopChan:    make(chan struct{}),
		pidFilePath: pidFilePath,
	}

	// 20251002 陈凤庆 启动时尝试从文件加载lastPID
	service.loadLastPIDFromFile()

	return service
}

/**
 * StartMonitoring 启动窗口监控
 * @return error 启动错误信息
 * @description 20251002 陈凤庆 启动后台监控服务，监控全局活动窗口变化
 */
func (wms *WindowMonitorService) StartMonitoring() error {
	wms.mu.Lock()
	defer wms.mu.Unlock()

	if wms.running {
		return fmt.Errorf("窗口监控服务已在运行")
	}

	fmt.Println("[窗口监控] 启动窗口监控服务...")
	wms.running = true

	// 20251002 陈凤庆 启动后台监控协程
	go wms.monitorLoop()

	fmt.Println("[窗口监控] 窗口监控服务启动完成")
	return nil
}

/**
 * StopMonitoring 停止窗口监控
 * @return error 停止错误信息
 * @description 20251002 陈凤庆 停止后台监控服务
 */
func (wms *WindowMonitorService) StopMonitoring() error {
	wms.mu.Lock()
	defer wms.mu.Unlock()

	if !wms.running {
		return fmt.Errorf("窗口监控服务未运行")
	}

	fmt.Println("[窗口监控] 停止窗口监控服务...")
	wms.running = false
	close(wms.stopChan)

	fmt.Println("[窗口监控] 窗口监控服务已停止")
	return nil
}

/**
 * monitorLoop 监控循环
 * @description 20251002 陈凤庆 后台监控循环，每100ms检查一次活动窗口变化和鼠标状态
 */
func (wms *WindowMonitorService) monitorLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	fmt.Println("[窗口监控] 开始监控活动窗口变化和鼠标交互...")

	for {
		select {
		case <-wms.stopChan:
			fmt.Println("[窗口监控] 监控循环退出")
			return
		case <-ticker.C:
			wms.checkWindowChange()
			wms.checkMouseInteraction()
		}
	}
}

/**
 * checkWindowChange 检查窗口变化
 * @description 20251002 陈凤庆 检查当前活动窗口，只在窗口真正发生变化时记录日志，并处理鼠标点击状态
 */
func (wms *WindowMonitorService) checkWindowChange() {
	// 20251002 陈凤庆 获取当前活动窗口信息
	currentPID := int(C.getCurrentActiveWindowPID())
	if currentPID == 0 {
		return // 无法获取当前窗口PID
	}

	cTitle := C.getCurrentActiveWindowTitle()
	defer C.free(unsafe.Pointer(cTitle))
	currentTitle := C.GoString(cTitle)

	// 20251002 陈凤庆 检查窗口是否真正发生变化
	wms.mu.RLock()
	lastRecordedPID := wms.lastPID
	var lastRecordedTitle string
	if wms.lastWindow != nil {
		lastRecordedTitle = wms.lastWindow.Title
	}
	lastActiveWindowPID := wms.lastActiveWindowPID
	wms.mu.RUnlock()

	// 20251002 陈凤庆 检测活动窗口是否发生变化
	if currentPID != lastActiveWindowPID {
		wms.mu.Lock()
		wms.lastActiveWindowPID = currentPID
		wms.mu.Unlock()
	}

	// 20251002 陈凤庆 判断是否应该记录为lastPID
	if wms.shouldRecordAsLastPID(currentPID, currentTitle) {
		// 20251002 陈凤庆 只有当窗口真正发生变化时才记录日志和更新
		if currentPID != lastRecordedPID || currentTitle != lastRecordedTitle {
			wms.mu.Lock()
			wms.lastPID = currentPID
			wms.lastWindow = &WindowInfo{
				PID:       currentPID,
				Title:     currentTitle,
				Timestamp: time.Now().Unix(),
			}
			wms.mu.Unlock()

			// 20251002 陈凤庆 保存到文件
			wms.saveLastPIDToFile()

			// 20251002 陈凤庆 只在窗口变化时记录日志
			if lastRecordedPID != 0 && lastRecordedTitle != "" {
				fmt.Printf("[窗口监控] 窗口切换: %s (PID: %d) → %s (PID: %d, selfPID: %d, lastPID: %d)\n",
					lastRecordedTitle, lastRecordedPID, currentTitle, currentPID, wms.selfPID, wms.lastPID)
			} else {
				fmt.Printf("[窗口监控] 初始记录 lastPID: %s (PID: %d, selfPID: %d, lastPID: %d)\n",
					currentTitle, currentPID, wms.selfPID, wms.lastPID)
			}
		}
	}
}

/**
 * shouldRecordAsLastPID 判断是否应该记录为lastPID
 * @param pid 进程ID
 * @param title 窗口标题
 * @return bool 是否应该记录
 * @description 20251002 陈凤庆 判断当前窗口是否应该被记录为lastPID
 * @modify 20251005 陈凤庆 增强过滤逻辑，排除macOS系统界面和登录界面
 */
func (wms *WindowMonitorService) shouldRecordAsLastPID(pid int, title string) bool {
	// 20251002 陈凤庆 PID必须有效且不能是密码管理器自身
	if pid <= 0 || pid == wms.selfPID {
		return false
	}

	// 20251005 陈凤庆 排除密码管理器相关的窗口标题和进程
	if title == "密码管理器" || title == "wails-demo" || title == "" {
		return false
	}

	// 20251005 陈凤庆 排除包含密码管理器关键词的窗口
	if strings.Contains(title, "wepass") || strings.Contains(title, "密码管理") {
		fmt.Printf("[窗口监控] 🚫 排除密码管理器相关窗口: %s (PID: %d)\n", title, pid)
		return false
	}

	// 20251005 陈凤庆 排除macOS系统界面和登录相关界面
	systemTitles := []string{
		"登录窗口",                // macOS登录界面
		"Login Window",        // macOS登录界面（英文）
		"ScreenSaverEngine",   // 屏幕保护程序
		"Screen Saver",        // 屏幕保护程序
		"锁定屏幕",                // 锁定屏幕
		"Lock Screen",         // 锁定屏幕（英文）
		"用户切换",                // 用户切换界面
		"User Switcher",       // 用户切换界面（英文）
		"Dock",                // macOS Dock
		"Finder",              // 当Finder是唯一活动窗口时通常表示桌面
		"桌面",                  // 桌面
		"Desktop",             // 桌面（英文）
		"Mission Control",     // Mission Control
		"Spotlight",           // Spotlight搜索
		"Notification Center", // 通知中心
		"Control Center",      // 控制中心
		"Siri",                // Siri界面
		"SystemUIServer",      // 系统UI服务器
		"WindowServer",        // 窗口服务器
	}

	for _, systemTitle := range systemTitles {
		if title == systemTitle {
			fmt.Printf("[窗口监控] 🚫 排除系统界面: %s (PID: %d)\n", title, pid)
			return false
		}
	}

	// 20251005 陈凤庆 排除以特定前缀开头的系统进程
	systemPrefixes := []string{
		"com.apple.",        // Apple系统进程
		"loginwindow",       // 登录窗口进程
		"screensaverengine", // 屏幕保护程序
		"SystemUIServer",    // 系统UI服务器
		"WindowServer",      // 窗口服务器
	}

	titleLower := strings.ToLower(title)
	for _, prefix := range systemPrefixes {
		if strings.HasPrefix(titleLower, strings.ToLower(prefix)) {
			fmt.Printf("[窗口监控] 🚫 排除系统进程: %s (PID: %d)\n", title, pid)
			return false
		}
	}

	return true
}

/**
 * isPIDValid 验证PID是否仍然有效
 * @param pid 进程ID
 * @return bool PID是否有效
 * @description 20251005 陈凤庆 检查指定PID的进程是否仍然存在
 */
func (wms *WindowMonitorService) isPIDValid(pid int) bool {
	if pid <= 0 {
		return false
	}
	return bool(C.checkPIDExists(C.int(pid)))
}

/**
 * clearInvalidLastPID 清理无效的lastPID
 * @description 20251005 陈凤庆 当检测到lastPID无效时，清理相关数据
 */
func (wms *WindowMonitorService) clearInvalidLastPID() {
	wms.mu.Lock()
	oldPID := wms.lastPID
	oldTitle := ""
	if wms.lastWindow != nil {
		oldTitle = wms.lastWindow.Title
	}

	wms.lastPID = 0
	wms.lastWindow = nil
	wms.mu.Unlock()

	// 清理文件中的记录
	wms.saveLastPIDToFile()

	fmt.Printf("[窗口监控] 🧹 已清理无效的lastPID: %d (%s)\n", oldPID, oldTitle)

	// 20251005 陈凤庆 尝试自动恢复：查找当前活动的有效窗口作为新的lastPID
	wms.tryRecoverLastPID()
}

/**
 * tryRecoverLastPID 尝试恢复lastPID
 * @description 20251005 陈凤庆 当lastPID无效时，尝试找到当前活动的有效窗口作为新的lastPID
 */
func (wms *WindowMonitorService) tryRecoverLastPID() {
	currentPID := int(C.getCurrentActiveWindowPID())
	if currentPID <= 0 {
		fmt.Printf("[窗口监控] 🔄 无法恢复lastPID：当前无活动窗口\n")
		return
	}

	cTitle := C.getCurrentActiveWindowTitle()
	defer C.free(unsafe.Pointer(cTitle))
	currentTitle := C.GoString(cTitle)

	if wms.shouldRecordAsLastPID(currentPID, currentTitle) {
		wms.mu.Lock()
		wms.lastPID = currentPID
		wms.lastWindow = &WindowInfo{
			PID:       currentPID,
			Title:     currentTitle,
			Timestamp: time.Now().Unix(),
		}
		wms.mu.Unlock()

		wms.saveLastPIDToFile()
		fmt.Printf("[窗口监控] 🔄 已恢复lastPID: %s (PID: %d)\n", currentTitle, currentPID)
	} else {
		fmt.Printf("[窗口监控] 🔄 当前窗口不适合作为lastPID: %s (PID: %d)\n", currentTitle, currentPID)
	}
}

/**
 * SwitchToLastWindow 切换到最后活动的窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活lastPID对应的窗口为活动窗口，记录切换前后状态
 */
func (wms *WindowMonitorService) SwitchToLastWindow() bool {
	// 20251002 陈凤庆 记录切换前的窗口状态
	beforePID := int(C.getCurrentActiveWindowPID())
	var beforeTitle string
	if beforePID > 0 {
		cTitle := C.getCurrentActiveWindowTitle()
		defer C.free(unsafe.Pointer(cTitle))
		beforeTitle = C.GoString(cTitle)
	}

	wms.mu.RLock()
	lastPID := wms.lastPID
	lastWindow := wms.lastWindow
	wms.mu.RUnlock()

	if lastPID <= 0 {
		fmt.Println("[窗口监控] ⚠️ lastPID无效，无法切换窗口")
		return false
	}

	targetTitle := "未知"
	if lastWindow != nil {
		targetTitle = lastWindow.Title
	}

	fmt.Printf("[窗口监控] ========== 开始窗口切换 ==========\n")
	fmt.Printf("[窗口监控] 切换前窗口: %s (PID: %d)\n", beforeTitle, beforePID)
	fmt.Printf("[窗口监控] 目标窗口: %s (PID: %d)\n", targetTitle, lastPID)

	// 20251005 陈凤庆 验证目标PID是否仍然有效
	if !wms.isPIDValid(lastPID) {
		fmt.Printf("[窗口监控] ⚠️ 目标PID已失效: %d (%s)，尝试清理并重新获取\n", lastPID, targetTitle)
		wms.clearInvalidLastPID()
		fmt.Printf("[窗口监控] ========== 切换失败 ==========\n")
		return false
	}

	// 20251002 陈凤庆 激活窗口
	if bool(C.activateWindowByPID(C.int(lastPID))) {
		// 20251002 陈凤庆 验证切换结果
		time.Sleep(100 * time.Millisecond) // 等待窗口激活
		afterPID := int(C.getCurrentActiveWindowPID())
		var afterTitle string
		if afterPID > 0 {
			cTitle := C.getCurrentActiveWindowTitle()
			defer C.free(unsafe.Pointer(cTitle))
			afterTitle = C.GoString(cTitle)
		}

		fmt.Printf("[窗口监控] 切换后窗口: %s (PID: %d)\n", afterTitle, afterPID)
		if afterPID == lastPID {
			fmt.Printf("[窗口监控] ✅ 窗口切换成功\n")
		} else {
			fmt.Printf("[窗口监控] ⚠️ 窗口切换验证失败，期望PID: %d，实际PID: %d\n", lastPID, afterPID)
			// 20251005 陈凤庆 如果切换后的窗口是系统界面，清理lastPID并提示用户
			if !wms.shouldRecordAsLastPID(afterPID, afterTitle) {
				fmt.Printf("[窗口监控] 🚫 切换后的窗口是系统界面，清理lastPID\n")
				wms.clearInvalidLastPID()
			}
		}
		fmt.Printf("[窗口监控] ========== 切换完成 ==========\n")
		return true
	}

	fmt.Printf("[窗口监控] ❌ 激活窗口失败: PID %d\n", lastPID)
	fmt.Printf("[窗口监控] ========== 切换失败 ==========\n")
	return false
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活密码管理器窗口为活动窗口，记录切换前后状态
 */
func (wms *WindowMonitorService) SwitchToPasswordManager() bool {
	// 20251002 陈凤庆 记录切换前的窗口状态
	beforePID := int(C.getCurrentActiveWindowPID())
	var beforeTitle string
	if beforePID > 0 {
		cTitle := C.getCurrentActiveWindowTitle()
		defer C.free(unsafe.Pointer(cTitle))
		beforeTitle = C.GoString(cTitle)
	}

	fmt.Printf("[窗口监控] ========== 切换回密码管理器 ==========\n")
	fmt.Printf("[窗口监控] 切换前窗口: %s (PID: %d)\n", beforeTitle, beforePID)
	fmt.Printf("[窗口监控] 目标窗口: 密码管理器 (PID: %d)\n", wms.selfPID)

	// 20251002 陈凤庆 激活密码管理器窗口
	if bool(C.activateWindowByPID(C.int(wms.selfPID))) {
		// 20251002 陈凤庆 验证切换结果
		time.Sleep(100 * time.Millisecond) // 等待窗口激活
		afterPID := int(C.getCurrentActiveWindowPID())
		var afterTitle string
		if afterPID > 0 {
			cTitle := C.getCurrentActiveWindowTitle()
			defer C.free(unsafe.Pointer(cTitle))
			afterTitle = C.GoString(cTitle)
		}

		fmt.Printf("[窗口监控] 切换后窗口: %s (PID: %d)\n", afterTitle, afterPID)
		if afterPID == wms.selfPID {
			fmt.Printf("[窗口监控] ✅ 成功切换回密码管理器\n")
		} else {
			fmt.Printf("[窗口监控] ⚠️ 切换验证失败，期望PID: %d，实际PID: %d\n", wms.selfPID, afterPID)
		}
		fmt.Printf("[窗口监控] ========== 切换完成 ==========\n")
		return true
	}

	fmt.Printf("[窗口监控] ❌ 激活密码管理器失败\n")
	fmt.Printf("[窗口监控] ========== 切换失败 ==========\n")
	return false
}

/**
 * RecordCurrentAsLastPID 手动记录当前活动窗口为lastPID
 * @return bool 记录是否成功
 * @description 20251002 陈凤庆 手动记录当前活动窗口，用于鼠标悬停等场景
 */
func (wms *WindowMonitorService) RecordCurrentAsLastPID() bool {
	currentPID := int(C.getCurrentActiveWindowPID())
	if currentPID == 0 {
		fmt.Println("[窗口监控] ⚠️ 无法获取当前窗口PID")
		return false
	}

	cTitle := C.getCurrentActiveWindowTitle()
	defer C.free(unsafe.Pointer(cTitle))
	currentTitle := C.GoString(cTitle)

	if wms.shouldRecordAsLastPID(currentPID, currentTitle) {
		wms.mu.Lock()
		wms.lastPID = currentPID
		wms.lastWindow = &WindowInfo{
			PID:       currentPID,
			Title:     currentTitle,
			Timestamp: time.Now().Unix(),
		}
		wms.mu.Unlock()

		// 20251002 陈凤庆 保存到文件
		wms.saveLastPIDToFile()

		fmt.Printf("[窗口监控] 手动记录 lastPID: %s (PID: %d)\n", currentTitle, currentPID)
		return true
	}

	fmt.Printf("[窗口监控] ⚠️ 当前窗口不符合记录条件: %s (PID: %d)\n", currentTitle, currentPID)
	return false
}

/**
 * GetLastWindow 获取最后活动窗口信息
 * @return *WindowInfo 窗口信息
 * @description 20251002 陈凤庆 获取lastPID对应的窗口信息
 */
func (wms *WindowMonitorService) GetLastWindow() *WindowInfo {
	wms.mu.RLock()
	defer wms.mu.RUnlock()

	if wms.lastWindow != nil {
		// 20251002 陈凤庆 返回副本，避免外部修改
		return &WindowInfo{
			PID:       wms.lastWindow.PID,
			Title:     wms.lastWindow.Title,
			Timestamp: wms.lastWindow.Timestamp,
		}
	}

	return nil
}

/**
 * GetSelfPID 获取密码管理器自身的PID
 * @return int 自身PID
 * @description 20251002 陈凤庆 获取密码管理器进程ID
 */
func (wms *WindowMonitorService) GetSelfPID() int {
	return wms.selfPID
}

/**
 * IsRunning 检查监控服务是否运行中
 * @return bool 是否运行中
 * @description 20251002 陈凤庆 检查监控服务状态
 */
func (wms *WindowMonitorService) IsRunning() bool {
	wms.mu.RLock()
	defer wms.mu.RUnlock()
	return wms.running
}

/**
 * saveLastPIDToFile 保存lastPID到文件
 * @description 20251002 陈凤庆 将lastPID保存到文件中，用于程序重启后恢复
 */
func (wms *WindowMonitorService) saveLastPIDToFile() {
	if wms.lastPID <= 0 {
		return
	}

	// 20251002 陈凤庆 确保logs目录存在
	if err := os.MkdirAll(filepath.Dir(wms.pidFilePath), 0755); err != nil {
		fmt.Printf("[窗口监控] ⚠️ 创建目录失败: %v\n", err)
		return
	}

	// 20251002 陈凤庆 写入PID到文件
	pidStr := strconv.Itoa(wms.lastPID)
	if err := os.WriteFile(wms.pidFilePath, []byte(pidStr), 0644); err != nil {
		fmt.Printf("[窗口监控] ⚠️ 保存PID到文件失败: %v\n", err)
	}
}

/**
 * loadLastPIDFromFile 从文件加载lastPID
 * @description 20251002 陈凤庆 从文件中加载lastPID，用于程序启动时恢复
 */
func (wms *WindowMonitorService) loadLastPIDFromFile() {
	if _, err := os.Stat(wms.pidFilePath); os.IsNotExist(err) {
		return // 文件不存在，跳过加载
	}

	data, err := os.ReadFile(wms.pidFilePath)
	if err != nil {
		fmt.Printf("[窗口监控] ⚠️ 读取PID文件失败: %v\n", err)
		return
	}

	pidStr := string(data)
	if pid, err := strconv.Atoi(pidStr); err == nil && pid > 0 {
		// 20251002 陈凤庆 验证进程是否仍然存在
		if bool(C.checkPIDExists(C.int(pid))) {
			wms.lastPID = pid
			fmt.Printf("[窗口监控] 从文件恢复 lastPID: %d\n", pid)
		} else {
			fmt.Printf("[窗口监控] 文件中的PID %d 已失效，清除文件\n", pid)
			os.Remove(wms.pidFilePath)
		}
	}
}

/**
 * checkMouseInteraction 检查鼠标交互
 * @description 20251002 陈凤庆 监控鼠标移动和点击，实现自动窗口切换逻辑
 */
func (wms *WindowMonitorService) checkMouseInteraction() {
	// 20251002 陈凤庆 获取当前活动窗口PID
	currentActivePID := int(C.getCurrentActiveWindowPID())

	wms.mu.RLock()
	selfPID := wms.selfPID
	lastPID := wms.lastPID
	mouseInSelfWindow := wms.mouseInSelfWindow
	wms.mu.RUnlock()

	// 20251002 陈凤庆 获取密码管理器窗口的位置和大小
	selfWindowRect := C.getWindowRectByPID(C.int(selfPID))

	// 20251002 陈凤庆 检查鼠标是否真正在密码管理器窗口内
	var mouseNowInSelfWindow bool
	if selfWindowRect.found {
		mouseNowInSelfWindow = bool(C.isMouseInWindowRect(
			selfWindowRect.x, selfWindowRect.y,
			selfWindowRect.width, selfWindowRect.height))
	} else {
		mouseNowInSelfWindow = false
	}

	// 20251002 陈凤庆 检测鼠标进入/离开密码管理器窗口
	if mouseNowInSelfWindow != mouseInSelfWindow {
		wms.mu.Lock()
		wms.mouseInSelfWindow = mouseNowInSelfWindow
		wms.mu.Unlock()

		if mouseNowInSelfWindow {
			fmt.Printf("[窗口监控] 鼠标进入密码管理器窗口范围 (PID: %d, selfPID: %d, lastPID: %d)\n",
				currentActivePID, selfPID, lastPID)

			// 20251005 陈凤庆 鼠标进入窗口时，自动激活窗口（恢复原逻辑）
			if currentActivePID != selfPID {
				// 20251002 陈凤庆 保存当前活动窗口到lastPID（如果不等于selfPID）
				if currentActivePID > 0 && currentActivePID != selfPID {
					cTitle := C.getCurrentActiveWindowTitle()
					defer C.free(unsafe.Pointer(cTitle))
					currentTitle := C.GoString(cTitle)

					wms.mu.Lock()
					wms.lastPID = currentActivePID
					wms.lastWindow = &WindowInfo{
						PID:       currentActivePID,
						Title:     currentTitle,
						Timestamp: time.Now().Unix(),
					}
					wms.mu.Unlock()

					// 20251002 陈凤庆 保存到文件
					wms.saveLastPIDToFile()

					fmt.Printf("[窗口监控] 保存当前活动窗口到lastPID: %s (PID: %d, selfPID: %d, lastPID: %d)\n",
						currentTitle, currentActivePID, selfPID, currentActivePID)
				}

				// 20251005 陈凤庆 激活密码管理器窗口
				fmt.Printf("[窗口监控] 激活密码管理器窗口 (当前活动窗口PID: %d, selfPID: %d, lastPID: %d)\n",
					currentActivePID, selfPID, wms.lastPID)
				success := bool(C.activateWindowByPID(C.int(selfPID)))
				if success {
					fmt.Printf("[窗口监控] ✅ 密码管理器窗口激活成功\n")
					fmt.Printf("[窗口监控] 鼠标进入窗口范围，等待用户点击操作\n")
				} else {
					fmt.Printf("[窗口监控] ❌ 密码管理器窗口激活失败\n")
				}
			} else {
				fmt.Printf("[窗口监控] 密码管理器已是活动窗口\n")
			}
		} else {
			fmt.Printf("[窗口监控] 鼠标离开密码管理器窗口范围 (PID: %d, selfPID: %d, lastPID: %d)\n",
				currentActivePID, selfPID, lastPID)

			// 20251005 陈凤庆 鼠标离开窗口时，立即切换回lastPID窗口（恢复原逻辑）
			if lastPID > 0 && lastPID != selfPID {
				fmt.Printf("[窗口监控] 鼠标离开窗口，自动切换到lastPID窗口 (PID: %d, selfPID: %d, lastPID: %d)\n",
					lastPID, selfPID, lastPID)
				wms.SwitchToLastWindow()
			} else {
				fmt.Printf("[窗口监控] lastPID无效或等于selfPID，不切换窗口 (PID: %d, selfPID: %d, lastPID: %d)\n",
					lastPID, selfPID, lastPID)
			}
		}
	}
}
