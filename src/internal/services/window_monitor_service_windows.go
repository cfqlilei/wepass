//go:build windows
// +build windows

package services

/*
#cgo CFLAGS: -I. -DUNICODE -D_UNICODE
#cgo LDFLAGS: -luser32 -lkernel32 -ladvapi32

#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <tlhelp32.h>

// 窗口信息结构体
typedef struct {
    HWND hwnd;
    DWORD pid;
    wchar_t title[256];  // 20251002 陈凤庆 改为wchar_t支持Unicode，避免中文乱码
    BOOL found;
} WindowSearchInfo;

// 鼠标位置结构体
typedef struct {
    int x;
    int y;
} MousePosition;

// 窗口矩形结构体
typedef struct {
    int x;
    int y;
    int width;
    int height;
    BOOL found;
} WindowRect;

// 获取当前进程PID
DWORD getCurrentProcessPID() {
    return GetCurrentProcessId();
}

// 获取当前活动窗口PID
DWORD getCurrentActiveWindowPID() {
    HWND hwnd = GetForegroundWindow();
    if (hwnd) {
        DWORD pid = 0;
        GetWindowThreadProcessId(hwnd, &pid);
        return pid;
    }
    return 0;
}

// 窗口枚举回调函数
BOOL CALLBACK WepassEnumWindowsProc(HWND hwnd, LPARAM lParam) {
    WindowSearchInfo* info = (WindowSearchInfo*)lParam;
    DWORD pid = 0;

    GetWindowThreadProcessId(hwnd, &pid);
    if (pid == info->pid) {
        // 20251005 陈凤庆 改进窗口选择逻辑，优先选择可见的主窗口

        // 检查是否为有效窗口
        if (IsWindow(hwnd)) {
            // 获取窗口标题
            int len = GetWindowTextW(hwnd, info->title, 255);

            // 优先选择可见且非最小化的窗口
            if (IsWindowVisible(hwnd) && !IsIconic(hwnd) && len > 0) {
                info->hwnd = hwnd;
                info->found = TRUE;
                return FALSE; // 找到最佳窗口，停止枚举
            }

            // 如果还没找到窗口，接受可见窗口（即使最小化）
            if (!info->found && IsWindowVisible(hwnd)) {
                info->hwnd = hwnd;
                info->found = TRUE;
                // 继续枚举，寻找更好的窗口
            }

            // 如果还没找到任何窗口，接受任何有标题的窗口
            if (!info->found && len > 0) {
                info->hwnd = hwnd;
                info->found = TRUE;
                // 继续枚举，寻找更好的窗口
            }
        }
    }
    return TRUE; // 继续枚举
}

// 激活指定PID的窗口
BOOL activateWindowByPID(DWORD pid) {
    WindowSearchInfo info = {0};
    info.pid = pid;
    info.found = FALSE;

    // 枚举所有窗口，找到匹配PID的主窗口
    EnumWindows(WepassEnumWindowsProc, (LPARAM)&info);

    if (info.found && info.hwnd) {
        // 20251005 陈凤庆 改进窗口激活逻辑，参考Mac平台的强制激活方式

        // 步骤1：如果窗口最小化，先恢复
        if (IsIconic(info.hwnd)) {
            ShowWindow(info.hwnd, SW_RESTORE);
            Sleep(50); // 等待窗口恢复
        }

        // 步骤2：确保窗口可见
        ShowWindow(info.hwnd, SW_SHOW);

        // 步骤3：将窗口置于前台
        BringWindowToTop(info.hwnd);

        // 步骤4：强制激活窗口 - 使用更强制的方法
        HWND currentForeground = GetForegroundWindow();
        DWORD currentThreadId = GetCurrentThreadId();
        DWORD targetThreadId = GetWindowThreadProcessId(info.hwnd, NULL);

        // 如果目标窗口在不同线程，需要附加输入
        if (currentThreadId != targetThreadId) {
            AttachThreadInput(currentThreadId, targetThreadId, TRUE);
        }

        // 尝试多种激活方法
        BOOL result = FALSE;

        // 方法1：直接设置前台窗口
        result = SetForegroundWindow(info.hwnd);

        if (!result) {
            // 方法2：使用SetActiveWindow
            SetActiveWindow(info.hwnd);
            result = TRUE;
        }

        if (!result) {
            // 方法3：使用SwitchToThisWindow（更强制）
            SwitchToThisWindow(info.hwnd, TRUE);
            result = TRUE;
        }

        // 步骤5：设置焦点
        SetFocus(info.hwnd);

        // 清理线程输入附加
        if (currentThreadId != targetThreadId) {
            AttachThreadInput(currentThreadId, targetThreadId, FALSE);
        }

        // 步骤6：验证激活结果
        Sleep(100); // 等待激活完成
        HWND newForeground = GetForegroundWindow();
        if (newForeground == info.hwnd) {
            result = TRUE;
        }

        return result;
    }
    return FALSE;
}

// 检查PID是否存在
BOOL checkPIDExists(DWORD pid) {
    HANDLE hProcess = OpenProcess(PROCESS_QUERY_INFORMATION, FALSE, pid);
    if (hProcess) {
        CloseHandle(hProcess);
        return TRUE;
    }
    return FALSE;
}

// 20251002 陈凤庆 获取当前活动窗口标题（Unicode版本，避免中文乱码）
const wchar_t* getCurrentActiveWindowTitleW() {
    static wchar_t title[256];
    HWND hwnd = GetForegroundWindow();
    if (hwnd) {
        int len = GetWindowTextW(hwnd, title, 255);
        if (len > 0) {
            return title;
        }
    }
    wcscpy(title, L"未知");
    return title;
}

// 获取当前鼠标位置
MousePosition getCurrentMousePosition() {
    MousePosition pos = {0, 0};
    POINT point;
    if (GetCursorPos(&point)) {
        pos.x = point.x;
        pos.y = point.y;
    }
    return pos;
}

// 获取指定PID进程的主窗口位置和大小
WindowRect getWindowRectByPID(DWORD pid) {
    WindowRect rect = {0, 0, 0, 0, FALSE};
    WindowSearchInfo info = {0};
    info.pid = pid;
    info.found = FALSE;

    EnumWindows(WepassEnumWindowsProc, (LPARAM)&info);

    if (info.found && info.hwnd) {
        RECT windowRect;
        if (GetWindowRect(info.hwnd, &windowRect)) {
            rect.x = windowRect.left;
            rect.y = windowRect.top;
            rect.width = windowRect.right - windowRect.left;
            rect.height = windowRect.bottom - windowRect.top;
            rect.found = TRUE;
        }
    }
    return rect;
}

// 检查鼠标是否在指定窗口区域内
BOOL isMouseInWindowRect(int windowX, int windowY, int windowWidth, int windowHeight) {
    MousePosition mousePos = getCurrentMousePosition();
    return (mousePos.x >= windowX && mousePos.x <= windowX + windowWidth &&
            mousePos.y >= windowY && mousePos.y <= windowY + windowHeight);
}

// 检查鼠标左键是否按下
BOOL isMouseButtonPressed() {
    return (GetAsyncKeyState(VK_LBUTTON) & 0x8000) != 0;
}

// 获取最后错误信息
const char* getLastErrorMessage() {
    static char buffer[256];
    DWORD error = GetLastError();
    FormatMessageA(
        FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS,
        NULL,
        error,
        MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
        buffer,
        sizeof(buffer),
        NULL
    );
    return buffer;
}

// 20251006 陈凤庆 根据PID获取窗口标题
char* getWindowTitleByPID(DWORD pid) {
    WindowSearchInfo info = {0};
    info.pid = pid;
    info.found = FALSE;

    // 枚举所有窗口，找到匹配PID的主窗口
    EnumWindows(WepassEnumWindowsProc, (LPARAM)&info);

    if (info.found && info.hwnd) {
        // 将wchar_t转换为char*
        int len = WideCharToMultiByte(CP_UTF8, 0, info.title, -1, NULL, 0, NULL, NULL);
        if (len > 0) {
            char* result = (char*)malloc(len);
            if (result) {
                WideCharToMultiByte(CP_UTF8, 0, info.title, -1, result, len, NULL, NULL);
                return result;
            }
        }
    }

    // 如果没有找到窗口或转换失败，返回默认值
    char* result = (char*)malloc(8);
    if (result) {
        strcpy(result, "未知");
    }
    return result;
}
*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unicode/utf16"
	"unsafe"
	"wepassword/internal/logger"
)

// 20251002 陈凤庆 辅助函数：将 wchar_t* 转换为 Go string（支持Unicode）
func wcharToGoString(wstr *C.wchar_t) string {
	if wstr == nil {
		return ""
	}

	// 计算字符串长度
	length := 0
	ptr := uintptr(unsafe.Pointer(wstr))
	for {
		ch := *(*uint16)(unsafe.Pointer(ptr))
		if ch == 0 {
			break
		}
		length++
		ptr += 2 // wchar_t 在 Windows 上是 2 字节
	}

	if length == 0 {
		return ""
	}

	// 转换为 []uint16
	slice := make([]uint16, length)
	ptr = uintptr(unsafe.Pointer(wstr))
	for i := 0; i < length; i++ {
		slice[i] = *(*uint16)(unsafe.Pointer(ptr + uintptr(i)*2))
	}

	// 使用 utf16.Decode 转换为 Go string
	return string(utf16.Decode(slice))
}

/**
 * WindowInfo Windows版本窗口信息结构
 * @author 陈凤庆
 * @description Windows平台的窗口信息
 */
type WindowInfo struct {
	PID       int    `json:"pid"`       // 进程ID
	Title     string `json:"title"`     // 窗口标题
	Timestamp int64  `json:"timestamp"` // 记录时间戳
}

/**
 * WindowMonitorService Windows版本窗口监控服务
 * @author 陈凤庆
 * @description Windows平台的窗口监控功能
 * @modify 陈凤庆 实现完整的Windows平台窗口监控功能
 */
type WindowMonitorService struct {
	selfPID                int           // 密码管理器自身的PID
	lastPID                int           // 最后活动的非密码管理器窗口PID
	lastWindow             *WindowInfo   // 最后活动窗口的详细信息
	running                bool          // 监控服务是否运行中
	stopChan               chan struct{} // 停止监控的信号通道
	mu                     sync.RWMutex  // 读写锁，保护共享数据
	pidFilePath            string        // PID文件存储路径
	mouseClickedSelfWindow bool          // 20251002 陈凤庆 是否激活后鼠标点击当前窗口的变量
	mouseInSelfWindow      bool          // 20251002 陈凤庆 鼠标是否在当前窗口内
	lastActiveWindowPID    int           // 20251002 陈凤庆 上次检查的活动窗口PID，用于检测窗口切换
}

/**
 * NewWindowMonitorService 创建新的窗口监控服务实例（Windows版本）
 * @return *WindowMonitorService 窗口监控服务实例
 * @description 20251002 陈凤庆 初始化Windows版本窗口监控服务，获取自身PID并设置文件路径
 */
func NewWindowMonitorService() *WindowMonitorService {
	// 20251002 陈凤庆 获取当前进程PID（密码管理器的PID）
	selfPID := int(C.getCurrentProcessPID())
	logger.Info("Windows版本窗口监控服务初始化，自身PID: %d", selfPID)

	// 20251003 陈凤庆 设置PID文件存储路径，Windows平台使用exe同目录
	var pidFilePath string
	if runtime.GOOS == "windows" {
		// Windows平台使用exe同目录下的logs文件夹
		execPath, err := os.Executable()
		if err == nil {
			// 获取可执行文件所在目录
			execDir := filepath.Dir(execPath)
			logDir := filepath.Join(execDir, "logs")
			os.MkdirAll(logDir, 0755) // 确保目录存在
			pidFilePath = filepath.Join(logDir, "lastPID.txt")
		} else {
			// 如果获取可执行文件路径失败，使用当前工作目录
			workDir, err := os.Getwd()
			if err == nil {
				logDir := filepath.Join(workDir, "logs")
				os.MkdirAll(logDir, 0755) // 确保目录存在
				pidFilePath = filepath.Join(logDir, "lastPID.txt")
			} else {
				// 最后回退到临时目录
				pidFilePath = filepath.Join(os.TempDir(), "wepass-lastPID.txt")
			}
		}
	} else {
		// 其他平台使用相对路径
		pidFilePath = filepath.Join("logs", "lastPID.txt")
	}
	logger.Info("PID文件路径: %s", pidFilePath)

	wms := &WindowMonitorService{
		selfPID:                selfPID,
		lastPID:                0,
		lastWindow:             nil,
		running:                false,
		stopChan:               make(chan struct{}),
		pidFilePath:            pidFilePath, // 保留字段但不使用文件操作
		mouseClickedSelfWindow: false,
		mouseInSelfWindow:      false,
		lastActiveWindowPID:    0,
	}

	// 20251009 陈凤庆 根据任务要求：不要把lastPID保存到文件，存放到内存即可，以便更高效地处理
	// 移除文件加载逻辑，只在内存中管理lastPID
	logger.Info("lastPID将仅在内存中管理，不使用文件持久化")

	return wms
}

/**
 * StartMonitoring 启动窗口监控（Windows版本）
 * @return error 错误信息
 * @description 20251002 陈凤庆 启动Windows版本窗口监控服务，监控活动窗口变化和鼠标交互
 */
func (wms *WindowMonitorService) StartMonitoring() error {
	wms.mu.Lock()
	if wms.running {
		wms.mu.Unlock()
		fmt.Printf("[窗口监控] Windows版本监控服务已在运行中\n")
		return nil
	}
	wms.running = true
	wms.mu.Unlock()

	logger.Info("Windows版本启动窗口监控服务")
	logger.Info("自身PID: %d", wms.selfPID)

	// 20251002 陈凤庆 启动后台监控协程
	go wms.monitorLoop()

	return nil
}

/**
 * StopMonitoring 停止窗口监控（Windows版本）
 * @return error 错误信息
 * @description 20251002 陈凤庆 停止Windows版本窗口监控服务
 */
func (wms *WindowMonitorService) StopMonitoring() error {
	wms.mu.Lock()
	if !wms.running {
		wms.mu.Unlock()
		fmt.Printf("[窗口监控] Windows版本监控服务未运行\n")
		return nil
	}
	wms.running = false
	wms.mu.Unlock()

	fmt.Printf("[窗口监控] Windows版本停止窗口监控服务\n")

	// 20251002 陈凤庆 发送停止信号
	select {
	case wms.stopChan <- struct{}{}:
	default:
	}

	return nil
}

/**
 * SwitchToLastWindow 切换到最后活动的窗口（Windows版本）
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活lastPID对应的窗口为活动窗口，记录切换前后状态
 */
func (wms *WindowMonitorService) SwitchToLastWindow() bool {
	// 20251002 陈凤庆 记录切换前的窗口状态
	beforePID := int(C.getCurrentActiveWindowPID())
	var beforeTitle string
	if beforePID > 0 {
		// 20251002 陈凤庆 使用Unicode版本获取窗口标题，避免中文乱码
		cTitle := C.getCurrentActiveWindowTitleW()
		beforeTitle = wcharToGoString(cTitle)
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

	logger.Info("========== 切换到目标窗口 ==========")
	logger.Info("切换前窗口: %s (PID: %d)", beforeTitle, beforePID)
	logger.Info("目标窗口: %s (PID: %d)", targetTitle, lastPID)

	// 20251002 陈凤庆 检查目标进程是否存在
	if !cBoolToBool(C.checkPIDExists(C.DWORD(lastPID))) {
		logger.Error("目标进程不存在: PID %d, 窗口标题: %s", lastPID, targetTitle)
		logger.Info("========== 切换失败 ==========")
		return false
	}

	// 20251005 陈凤庆 添加详细的激活日志
	logger.Info("开始激活目标窗口 (PID: %d, 标题: %s)", lastPID, targetTitle)

	// 20251002 陈凤庆 激活窗口
	if cBoolToBool(C.activateWindowByPID(C.DWORD(lastPID))) {
		// 20251005 陈凤庆 增加等待时间，确保窗口激活完成
		fmt.Printf("[窗口监控] ⏳ 等待窗口激活完成 (300ms)\n")
		time.Sleep(300 * time.Millisecond) // 等待窗口激活
		afterPID := int(C.getCurrentActiveWindowPID())
		var afterTitle string
		if afterPID > 0 {
			// 20251002 陈凤庆 使用Unicode版本获取窗口标题，避免中文乱码
			cTitle := C.getCurrentActiveWindowTitleW()
			afterTitle = wcharToGoString(cTitle)
		}

		logger.Info("切换后窗口: %s (PID: %d)", afterTitle, afterPID)
		if afterPID == lastPID {
			logger.Info("窗口切换成功，目标窗口: %s (PID: %d)", afterTitle, afterPID)
		} else {
			logger.Error("窗口切换验证失败，期望PID: %d (标题: %s)，实际PID: %d (标题: %s)", lastPID, targetTitle, afterPID, afterTitle)
		}
		logger.Info("========== 切换完成 ==========")
		return true
	}

	fmt.Printf("[窗口监控] ❌ 激活窗口失败: PID %d\n", lastPID)
	cErrorMsg := C.getLastErrorMessage()
	errorMsg := C.GoString(cErrorMsg)
	fmt.Printf("[窗口监控] 错误信息: %s\n", errorMsg)
	fmt.Printf("[窗口监控] ========== 切换失败 ==========\n")
	return false
}

/**
 * SwitchToPasswordManager 切换回密码管理器窗口（Windows版本）
 * @return bool 切换是否成功
 * @description 20251002 陈凤庆 激活密码管理器窗口为活动窗口，记录切换前后状态
 */
func (wms *WindowMonitorService) SwitchToPasswordManager() bool {
	// 20251002 陈凤庆 记录切换前的窗口状态
	beforePID := int(C.getCurrentActiveWindowPID())
	var beforeTitle string
	if beforePID > 0 {
		// 20251002 陈凤庆 使用Unicode版本获取窗口标题，避免中文乱码
		cTitle := C.getCurrentActiveWindowTitleW()
		beforeTitle = wcharToGoString(cTitle)
	}

	fmt.Printf("[窗口监控] ========== 切换回密码管理器 ==========\n")
	fmt.Printf("[窗口监控] 切换前窗口: %s (PID: %d)\n", beforeTitle, beforePID)
	fmt.Printf("[窗口监控] 目标窗口: 密码管理器 (PID: %d)\n", wms.selfPID)

	// 20251005 陈凤庆 添加详细的激活日志
	fmt.Printf("[窗口监控] 🔄 开始激活密码管理器窗口 (PID: %d)\n", wms.selfPID)

	// 20251002 陈凤庆 激活密码管理器窗口
	if cBoolToBool(C.activateWindowByPID(C.DWORD(wms.selfPID))) {
		// 20251005 陈凤庆 增加等待时间，确保窗口激活完成
		fmt.Printf("[窗口监控] ⏳ 等待密码管理器激活完成 (300ms)\n")
		time.Sleep(300 * time.Millisecond) // 等待窗口激活
		afterPID := int(C.getCurrentActiveWindowPID())
		var afterTitle string
		if afterPID > 0 {
			// 20251002 陈凤庆 使用Unicode版本获取窗口标题，避免中文乱码
			cTitle := C.getCurrentActiveWindowTitleW()
			afterTitle = wcharToGoString(cTitle)
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
	cErrorMsg := C.getLastErrorMessage()
	errorMsg := C.GoString(cErrorMsg)
	fmt.Printf("[窗口监控] 错误信息: %s\n", errorMsg)
	fmt.Printf("[窗口监控] ========== 切换失败 ==========\n")
	return false
}

/**
 * RecordCurrentAsLastPID 记录当前窗口为lastPID（Windows版本）
 * @return bool 记录是否成功
 * @description 20251002 陈凤庆 手动记录当前活动窗口为lastPID，用于特殊场景
 */
func (wms *WindowMonitorService) RecordCurrentAsLastPID() bool {
	currentPID := int(C.getCurrentActiveWindowPID())
	if currentPID <= 0 {
		fmt.Printf("[窗口监控] ⚠️ 获取当前活动窗口PID失败: %d\n", currentPID)
		return false
	}
	if currentPID == wms.selfPID {
		// 20251008 陈凤庆 获取当前窗口标题用于调试
		cTitle := C.getCurrentActiveWindowTitleW()
		currentTitle := wcharToGoString(cTitle)
		fmt.Printf("[窗口监控] ⚠️ 当前活动窗口是密码管理器自身，无法记录为lastPID (当前窗口: %s, PID: %d)\n", currentTitle, currentPID)
		return false
	}

	// 20251002 陈凤庆 使用Unicode版本获取窗口标题，避免中文乱码
	cTitle := C.getCurrentActiveWindowTitleW()
	title := wcharToGoString(cTitle)

	wms.mu.Lock()
	wms.lastPID = currentPID
	wms.lastWindow = &WindowInfo{
		PID:       currentPID,
		Title:     title,
		Timestamp: time.Now().Unix(),
	}
	wms.mu.Unlock()

	// 20251009 陈凤庆 根据任务要求：不保存到文件，只在内存中管理
	logger.Info("手动记录当前窗口: %s (PID: %d) [仅内存]", title, currentPID)
	return true
}

/**
 * GetLastWindow 获取最后活动窗口信息（Windows版本）
 * @return *WindowInfo 窗口信息
 * @description 20251002 陈凤庆 返回最后记录的窗口信息
 */
func (wms *WindowMonitorService) GetLastWindow() *WindowInfo {
	wms.mu.RLock()
	defer wms.mu.RUnlock()
	return wms.lastWindow
}

/**
 * IsRunning 检查窗口监控是否正在运行（Windows版本）
 * @return bool 是否正在运行
 * @description 20251002 陈凤庆 返回监控服务运行状态
 */
func (wms *WindowMonitorService) IsRunning() bool {
	wms.mu.RLock()
	defer wms.mu.RUnlock()
	return wms.running
}

/**
 * monitorLoop 窗口监控主循环（Windows版本）
 * @description 20251002 陈凤庆 后台监控活动窗口变化和鼠标交互
 */
func (wms *WindowMonitorService) monitorLoop() {
	fmt.Printf("[窗口监控] Windows版本监控循环启动\n")
	ticker := time.NewTicker(100 * time.Millisecond) // 100ms检查一次
	defer ticker.Stop()

	for {
		select {
		case <-wms.stopChan:
			fmt.Printf("[窗口监控] Windows版本监控循环停止\n")
			return
		case <-ticker.C:
			wms.checkWindowAndMouseState()
		}
	}
}

/**
 * checkWindowAndMouseState 检查窗口和鼠标状态（Windows版本）
 * @description 20251002 陈凤庆 检查活动窗口变化和鼠标交互状态
 * @modify 20251009 陈凤庆 增强调试日志，记录当前活动窗口的PID和标题
 */
func (wms *WindowMonitorService) checkWindowAndMouseState() {
	currentActiveWindowPID := int(C.getCurrentActiveWindowPID())
	if currentActiveWindowPID <= 0 {
		fmt.Printf("[窗口监控] ⚠️ 无效的当前窗口PID: %d\n", currentActiveWindowPID)
		return
	}

	// 20251009 陈凤庆 获取当前活动窗口标题用于调试
	cTitle := C.getCurrentActiveWindowTitleW()
	currentTitle := wcharToGoString(cTitle)

	// 20251002 陈凤庆 检查活动窗口是否发生变化
	wms.mu.Lock()
	lastActiveWindowPID := wms.lastActiveWindowPID
	wms.lastActiveWindowPID = currentActiveWindowPID
	selfPID := wms.selfPID
	lastPID := wms.lastPID
	wms.mu.Unlock()

	// 20251009 陈凤庆 增强调试日志：记录当前活动窗口信息
	if currentActiveWindowPID != lastActiveWindowPID {
		logger.Info("活动窗口变化: 上一个窗口 (PID: %d) -> %s (PID: %d)",
			lastActiveWindowPID, currentTitle, currentActiveWindowPID)
		logger.Info("状态信息: selfPID=%d, lastPID=%d, 当前PID=%d, 当前窗口标题=%s",
			selfPID, lastPID, currentActiveWindowPID, currentTitle)
	}

	// 20251002 陈凤庆 如果活动窗口发生变化，记录非自身窗口
	if currentActiveWindowPID != lastActiveWindowPID {
		if currentActiveWindowPID != wms.selfPID {
			wms.mu.Lock()
			wms.lastPID = currentActiveWindowPID
			wms.lastWindow = &WindowInfo{
				PID:       currentActiveWindowPID,
				Title:     currentTitle,
				Timestamp: time.Now().Unix(),
			}
			// 20251006 陈凤庆 重置鼠标点击状态（恢复e542405988358758f18c1bd7c15b7bd8c193ebfb版本逻辑）
			wms.mouseClickedSelfWindow = false

			wms.mu.Unlock()

			// 20251009 陈凤庆 根据任务要求：不保存到文件，只在内存中管理
			logger.Info("记录新的活动窗口: %s (PID: %d) [仅内存]", currentTitle, currentActiveWindowPID)
		} else {
			// 20251002 陈凤庆 切换到自身窗口，检查鼠标状态
			logger.Info("切换到密码管理器窗口: %s (PID: %d)", currentTitle, currentActiveWindowPID)
			wms.handleSelfWindowActivation()
		}
	}

	// 20251002 陈凤庆 检查鼠标状态
	wms.checkMouseState()
}

/**
 * handleSelfWindowActivation 处理自身窗口激活（Windows版本）
 * @description 20251002 陈凤庆 当密码管理器窗口被激活时的处理逻辑
 */
func (wms *WindowMonitorService) handleSelfWindowActivation() {
	// 20251002 陈凤庆 检查鼠标是否在自身窗口内
	mouseInWindow := wms.isMouseInSelfWindow()

	wms.mu.Lock()
	wms.mouseInSelfWindow = mouseInWindow
	if !mouseInWindow {
		// 20251002 陈凤庆 鼠标不在窗口内，可能是程序激活，重置点击状态
		wms.mouseClickedSelfWindow = false
	}
	wms.mu.Unlock()
}

/**
 * checkMouseState 检查鼠标状态（Windows版本）
 * @description 20251002 陈凤庆 检查鼠标位置和点击状态
 */
func (wms *WindowMonitorService) checkMouseState() {
	mouseInWindow := wms.isMouseInSelfWindow()
	mousePressed := cBoolToBool(C.isMouseButtonPressed())

	// 20251008 陈凤庆 添加详细的鼠标状态调试日志
	if mouseInWindow != wms.mouseInSelfWindow {
		fmt.Printf("[窗口监控] 🔍 鼠标状态变化: 当前在窗口内=%t, 之前在窗口内=%t\n", mouseInWindow, wms.mouseInSelfWindow)
	}

	wms.mu.Lock()
	defer wms.mu.Unlock()

	// 20251002 陈凤庆 鼠标进入自身窗口
	if mouseInWindow && !wms.mouseInSelfWindow {
		wms.mouseInSelfWindow = true

		// 20251009 陈凤庆 释放锁，避免在获取窗口信息时持有锁
		wms.mu.Unlock()

		// 20251008 陈凤庆 修复：鼠标移入时，必须先获取当前活动窗口PID并记录为lastPID
		currentActivePID := int(C.getCurrentActiveWindowPID())

		// 20251008 陈凤庆 获取当前活动窗口的标题用于调试
		var currentTitle string
		if currentActivePID > 0 {
			cTitle := C.getCurrentActiveWindowTitleW()
			currentTitle = wcharToGoString(cTitle)
		} else {
			currentTitle = "未知窗口"
		}

		logger.Info("鼠标进入密码管理器窗口范围")
		logger.Info("当前活动窗口: %s (PID: %d)", currentTitle, currentActivePID)
		logger.Info("密码管理器PID: %d", wms.selfPID)

		wms.mu.Lock() // 重新获取锁
		logger.Info("之前的lastPID: %d", wms.lastPID)

		// 20251008 陈凤庆 根据设计要求：无论当前活动窗口是什么，都要先记录为lastPID，然后激活密码管理器
		if currentActivePID > 0 {
			if currentActivePID != wms.selfPID {
				// 20251008 陈凤庆 当前活动窗口不是密码管理器，记录为lastPID
				wms.lastPID = currentActivePID
				wms.lastWindow = &WindowInfo{
					PID:       currentActivePID,
					Title:     currentTitle,
					Timestamp: time.Now().Unix(),
				}

				logger.Info("记录当前活动窗口为lastPID: %s (PID: %d) [仅内存]", currentTitle, currentActivePID)

				// 20251009 陈凤庆 释放锁后再进行窗口激活操作
				wms.mu.Unlock()

				// 20251008 陈凤庆 然后激活密码管理器窗口
				logger.Info("开始激活密码管理器窗口 (selfPID: %d)", wms.selfPID)
				success := cBoolToBool(C.activateWindowByPID(C.DWORD(wms.selfPID)))
				if success {
					fmt.Printf("[窗口监控] ✅ 密码管理器窗口激活成功\n")
				} else {
					fmt.Printf("[窗口监控] ❌ 密码管理器窗口激活失败\n")
				}

				wms.mu.Lock() // 重新获取锁以保持defer的一致性
			} else {
				// 20251008 陈凤庆 当前活动窗口已经是密码管理器，不需要切换
				fmt.Printf("[窗口监控] ℹ️ 密码管理器已是活动窗口，无需切换\n")
			}
		} else {
			fmt.Printf("[窗口监控] ⚠️ 获取当前活动窗口PID失败: %d\n", currentActivePID)
		}
	}

	// 20251002 陈凤庆 鼠标离开自身窗口
	if !mouseInWindow && wms.mouseInSelfWindow {
		wms.mouseInSelfWindow = false

		// 20251008 陈凤庆 根据设计要求：当鼠标离开密码管理器窗口时，自动将焦点切换回lastPID对应的窗口
		if wms.lastPID > 0 {
			// 20251008 陈凤庆 获取lastPID窗口的详细信息用于调试
			var lastWindowTitle string
			if wms.lastWindow != nil {
				lastWindowTitle = wms.lastWindow.Title
			} else {
				lastWindowTitle = "未知"
			}

			fmt.Printf("[窗口监控] 鼠标离开密码管理器窗口，准备切换回lastPID窗口\n")
			fmt.Printf("[窗口监控] 目标窗口: %s (PID: %d)\n", lastWindowTitle, wms.lastPID)

			// 20251008 陈凤庆 检查目标进程是否还存在
			if !cBoolToBool(C.checkPIDExists(C.DWORD(wms.lastPID))) {
				fmt.Printf("[窗口监控] ❌ 目标进程不存在: PID %d，无法切换\n", wms.lastPID)
				return
			}

			go func() {
				time.Sleep(50 * time.Millisecond) // 短暂延迟避免频繁切换

				// 20251008 陈凤庆 记录切换前的当前活动窗口
				beforePID := int(C.getCurrentActiveWindowPID())
				var beforeTitle string
				if beforePID > 0 {
					cTitle := C.getCurrentActiveWindowTitleW()
					beforeTitle = wcharToGoString(cTitle)
				}

				fmt.Printf("[窗口监控] 切换前活动窗口: %s (PID: %d)\n", beforeTitle, beforePID)
				fmt.Printf("[窗口监控] 开始激活目标窗口: %s (PID: %d)\n", lastWindowTitle, wms.lastPID)

				if cBoolToBool(C.activateWindowByPID(C.DWORD(wms.lastPID))) {
					// 20251008 陈凤庆 验证切换是否成功
					time.Sleep(100 * time.Millisecond) // 等待窗口激活完成
					afterPID := int(C.getCurrentActiveWindowPID())
					var afterTitle string
					if afterPID > 0 {
						cTitle := C.getCurrentActiveWindowTitleW()
						afterTitle = wcharToGoString(cTitle)
					}

					fmt.Printf("[窗口监控] 切换后活动窗口: %s (PID: %d)\n", afterTitle, afterPID)
					if afterPID == wms.lastPID {
						fmt.Printf("[窗口监控] ✅ 自动切换回lastPID窗口成功\n")
					} else {
						fmt.Printf("[窗口监控] ⚠️ 窗口切换验证失败，期望PID: %d，实际PID: %d\n", wms.lastPID, afterPID)
					}
				} else {
					fmt.Printf("[窗口监控] ❌ 切换回lastPID窗口失败 (PID: %d)\n", wms.lastPID)
					cErrorMsg := C.getLastErrorMessage()
					errorMsg := C.GoString(cErrorMsg)
					fmt.Printf("[窗口监控] 错误信息: %s\n", errorMsg)
				}
			}()
		} else {
			fmt.Printf("[窗口监控] ⚠️ lastPID无效，无法切换回上一个窗口\n")
		}
	}

	// 20251006 陈凤庆 检测鼠标点击（恢复e542405988358758f18c1bd7c15b7bd8c193ebfb版本逻辑）
	if mouseInWindow && mousePressed {
		// 鼠标在窗口内且按下，标记为已点击
		if !wms.mouseClickedSelfWindow {
			wms.mouseClickedSelfWindow = true
			fmt.Printf("[窗口监控] 检测到鼠标点击密码管理器窗口\n")
		}
	}

}

/**
 * isMouseInSelfWindow 检查鼠标是否在自身窗口内（Windows版本）
 * @return bool 鼠标是否在窗口内
 * @description 20251002 陈凤庆 检查鼠标位置是否在密码管理器窗口区域内
 */
func (wms *WindowMonitorService) isMouseInSelfWindow() bool {
	windowRect := C.getWindowRectByPID(C.DWORD(wms.selfPID))
	if !cBoolToBool(windowRect.found) {
		return false
	}

	return cBoolToBool(C.isMouseInWindowRect(
		C.int(windowRect.x),
		C.int(windowRect.y),
		C.int(windowRect.width),
		C.int(windowRect.height),
	))
}

/**
 * saveLastPIDToFile 保存lastPID到文件（Windows版本）
 * @description 20251002 陈凤庆 将lastPID保存到文件中，用于程序重启后恢复
 */
func (wms *WindowMonitorService) saveLastPIDToFile() {
	wms.mu.RLock()
	lastPID := wms.lastPID
	wms.mu.RUnlock()

	if lastPID <= 0 {
		return
	}

	// 20251002 陈凤庆 确保logs目录存在
	if err := os.MkdirAll(filepath.Dir(wms.pidFilePath), 0755); err != nil {
		fmt.Printf("[窗口监控] ⚠️ 创建logs目录失败: %v\n", err)
		return
	}

	// 20251002 陈凤庆 写入文件
	if err := os.WriteFile(wms.pidFilePath, []byte(strconv.Itoa(lastPID)), 0644); err != nil {
		fmt.Printf("[窗口监控] ⚠️ 保存lastPID到文件失败: %v\n", err)
	}
}

/**
 * loadLastPIDFromFile 从文件加载lastPID（Windows版本）
 * @description 20251002 陈凤庆 从文件中恢复lastPID，用于程序重启后恢复状态
 */
func (wms *WindowMonitorService) loadLastPIDFromFile() {
	data, err := os.ReadFile(wms.pidFilePath)
	if err != nil {
		fmt.Printf("[窗口监控] 无法读取lastPID文件: %v\n", err)
		return
	}

	lastPID, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Printf("[窗口监控] 解析lastPID失败: %v\n", err)
		return
	}

	// 20251002 陈凤庆 检查进程是否仍然存在
	if cBoolToBool(C.checkPIDExists(C.DWORD(lastPID))) {
		// 20251008 陈凤庆 修复：使用指定PID获取窗口标题，而不是当前活动窗口的标题
		cTitle := C.getWindowTitleByPID(C.DWORD(lastPID))
		defer C.free(unsafe.Pointer(cTitle))
		title := C.GoString(cTitle)

		wms.mu.Lock()
		wms.lastPID = lastPID
		wms.lastWindow = &WindowInfo{
			PID:       lastPID,
			Title:     title,
			Timestamp: time.Now().Unix(),
		}
		wms.mu.Unlock()

		fmt.Printf("[窗口监控] ✅ 从文件恢复lastPID: %d (%s)\n", lastPID, title)
	} else {
		fmt.Printf("[窗口监控] ⚠️ 文件中的lastPID进程不存在: %d\n", lastPID)
		// 20251002 陈凤庆 删除无效的PID文件
		os.Remove(wms.pidFilePath)
	}
}
