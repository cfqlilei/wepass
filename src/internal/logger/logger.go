package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

/**
 * 日志管理器
 * @author 陈凤庆
 * @date 20251001
 * @description 统一的日志记录管理，支持文件和控制台输出
 */

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger

	// 20251003 陈凤庆 添加文件句柄，用于强制刷新
	infoFile  *os.File
	errorFile *os.File
	debugFile *os.File
	logMutex  sync.Mutex

	// 20251003 陈凤庆 添加日志配置控制
	logConfig LogConfig

	// 20251005 陈凤庆 添加环境检测变量
	isDevelopmentEnv bool
)

/**
 * getModuleName 根据文件名获取模块名称
 * @param filename 文件名
 * @return string 模块名称
 * @description 20251009 陈凤庆 为Windows平台日志格式提供模块名称映射
 */
func getModuleName(filename string) string {
	// 移除文件扩展名
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 根据文件名映射到模块名称
	switch name {
	case "window_monitor_service_windows":
		return "WindowManager"
	case "account_service":
		return "AccountService"
	case "keyboard_service_windows":
		return "KeyboardService"
	case "app_windows", "app_darwin", "app_linux", "app":
		return "Application"
	case "logger":
		return "Logger"
	case "main":
		return "Main"
	default:
		// 将下划线转换为驼峰命名
		parts := strings.Split(name, "_")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
		return strings.Join(parts, "")
	}
}

/**
 * LogConfig 日志配置结构
 * @author 陈凤庆
 * @date 20251003
 */
type LogConfig struct {
	EnableInfoLog  bool
	EnableDebugLog bool
}

/**
 * SetLogConfig 设置日志配置
 * @param config 日志配置
 * @author 陈凤庆
 * @date 20251003
 */
func SetLogConfig(config LogConfig) {
	logMutex.Lock()
	defer logMutex.Unlock()
	logConfig = config
}

/**
 * detectEnvironment 检测运行环境
 * @return bool 是否为开发环境
 * @author 陈凤庆
 * @date 20251005
 * @description 检测当前是开发环境还是生产环境，影响默认日志配置
 */
func detectEnvironment() bool {
	// 方法1：检查环境变量（最高优先级）
	if os.Getenv("WAILS_DEV") != "" {
		return true
	}

	// 方法2：检查可执行文件路径是否包含wails相关路径
	if execPath, err := os.Executable(); err == nil {
		// 如果可执行文件路径包含wails临时目录，则为开发环境
		if strings.Contains(execPath, "wailsbindings") || strings.Contains(execPath, "/tmp/") {
			return true
		}

		// 如果在build/bin目录下，但同时存在源码目录结构，则为开发环境
		if strings.Contains(execPath, "build/bin") {
			// 检查是否存在源码目录结构
			execDir := filepath.Dir(execPath)
			// 向上查找，看是否能找到src目录
			for i := 0; i < 5; i++ {
				execDir = filepath.Dir(execDir)
				if _, err := os.Stat(filepath.Join(execDir, "src", "wails.json")); err == nil {
					return true
				}
				if _, err := os.Stat(filepath.Join(execDir, "src", "go.mod")); err == nil {
					return true
				}
			}
		}
	}

	// 方法3：检查当前工作目录是否为开发环境
	if workDir, err := os.Getwd(); err == nil {
		// 如果当前目录或父目录存在 wails.json，则为开发环境
		checkDir := workDir
		for i := 0; i < 5; i++ {
			if _, err := os.Stat(filepath.Join(checkDir, "wails.json")); err == nil {
				return true
			}
			if _, err := os.Stat(filepath.Join(checkDir, "src", "wails.json")); err == nil {
				return true
			}
			parentDir := filepath.Dir(checkDir)
			if parentDir == checkDir {
				break // 已到根目录
			}
			checkDir = parentDir
		}
	}

	// 默认为生产环境
	return false
}

/**
 * IsDevelopmentEnvironment 获取当前环境类型
 * @return bool 是否为开发环境
 * @author 陈凤庆
 * @date 20251005
 */
func IsDevelopmentEnvironment() bool {
	return isDevelopmentEnv
}

/**
 * GetDefaultLogConfig 获取默认日志配置
 * @return LogConfig 默认日志配置
 * @author 陈凤庆
 * @date 20251005
 * @description 根据环境返回不同的默认日志配置
 */
func GetDefaultLogConfig() LogConfig {
	if isDevelopmentEnv {
		// 开发环境：启用所有日志
		return LogConfig{
			EnableInfoLog:  true,
			EnableDebugLog: true,
		}
	} else {
		// 生产环境：默认只启用错误日志，关闭Info和Debug
		return LogConfig{
			EnableInfoLog:  false,
			EnableDebugLog: false,
		}
	}
}

/**
 * InfoStartup 记录程序启动相关的重要信息日志
 * @param format 格式字符串
 * @param v 参数
 * @author 陈凤庆
 * @date 20251005
 * @description 专门用于记录程序启动、关闭等重要提醒信息
 */
func InfoStartup(format string, v ...interface{}) {
	Info("[启动] "+format, v...)
}

/**
 * InfoOperation 记录操作相关的信息日志
 * @param format 格式字符串
 * @param v 参数
 * @author 陈凤庆
 * @date 20251005
 * @description 专门用于记录用户操作、系统状态变化等信息
 */
func InfoOperation(format string, v ...interface{}) {
	Info("[操作] "+format, v...)
}

/**
 * DebugDetail 记录详细的调试信息
 * @param format 格式字符串
 * @param v 参数
 * @author 陈凤庆
 * @date 20251005
 * @description 记录完整的调试信息，便于发现问题
 */
func DebugDetail(format string, v ...interface{}) {
	Debug("[详细] "+format, v...)
}

/**
 * ErrorSystem 记录系统错误
 * @param format 格式字符串
 * @param v 参数
 * @author 陈凤庆
 * @date 20251005
 * @description 只在程序出现错误时记录，避免冗余日志
 */
func ErrorSystem(format string, v ...interface{}) {
	Error("[系统错误] "+format, v...)
}

/**
 * GetLogConfig 获取日志配置
 * @return LogConfig 日志配置
 * @author 陈凤庆
 * @date 20251003
 */
func GetLogConfig() LogConfig {
	logMutex.Lock()
	defer logMutex.Unlock()
	return logConfig
}

/**
 * InitLogger 初始化日志系统
 * @param logDir 日志目录
 * @return error 错误信息
 */
func InitLogger(logDir string) error {
	// 20251005 陈凤庆 检测运行环境
	isDevelopmentEnv = detectEnvironment()

	// 20251005 陈凤庆 根据环境设置默认日志配置
	logConfig = GetDefaultLogConfig()

	// 20251003 陈凤庆 创建一个临时的初始化日志收集器
	var initLogs []string

	// 添加日志记录函数
	addInitLog := func(format string, args ...interface{}) {
		// 20251003 陈凤庆 统一初始化日志格式：时间 [INIT] 消息
		now := time.Now()
		timeStr := now.Format("2006/01/02 15:04:05")
		msg := fmt.Sprintf(format, args...)
		formattedMsg := fmt.Sprintf("%s [INIT] %s", timeStr, msg)

		fmt.Printf("%s\n", formattedMsg)          // 输出到控制台
		initLogs = append(initLogs, formattedMsg) // 收集到数组中
	}

	// 20251003 陈凤庆 添加详细的调试信息，帮助排查Windows平台日志问题
	addInitLog("[日志初始化] 开始初始化日志系统，目标目录: %s", logDir)

	// 20251005 陈凤庆 添加环境检测信息
	envType := "生产环境"
	if isDevelopmentEnv {
		envType = "开发环境"
	}
	addInitLog("[日志初始化] 检测到运行环境: %s", envType)
	addInitLog("[日志初始化] 日志配置: Info=%v, Debug=%v", logConfig.EnableInfoLog, logConfig.EnableDebugLog)

	// 20251003 陈凤庆 添加系统环境信息
	addInitLog("[日志初始化] 运行环境信息:")
	addInitLog("  - 操作系统: %s", runtime.GOOS)
	addInitLog("  - 架构: %s", runtime.GOARCH)

	// 获取当前工作目录
	if workDir, err := os.Getwd(); err == nil {
		addInitLog("  - 当前工作目录: %s", workDir)
	}

	// 获取可执行文件路径
	if execPath, err := os.Executable(); err == nil {
		addInitLog("  - 可执行文件路径: %s", execPath)
		addInitLog("  - 可执行文件目录: %s", filepath.Dir(execPath))
	}

	// 获取用户主目录
	if homeDir, err := os.UserHomeDir(); err == nil {
		addInitLog("  - 用户主目录: %s", homeDir)
	}

	// 确保日志目录存在
	addInitLog("[日志初始化] 正在创建日志目录: %s", logDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		addInitLog("[日志初始化] ❌ 创建日志目录失败: %v", err)
		return fmt.Errorf("创建日志目录失败: %w", err)
	}
	addInitLog("[日志初始化] ✅ 日志目录创建成功: %s", logDir)

	// 检查目录权限
	if info, err := os.Stat(logDir); err == nil {
		addInitLog("[日志初始化] 目录权限: %s", info.Mode().String())
	}

	// 获取当前日期
	now := time.Now()
	dateStr := now.Format("20060102")
	addInitLog("[日志初始化] 使用日期字符串: %s", dateStr)

	// 创建日志文件路径
	infoLogFile := filepath.Join(logDir, fmt.Sprintf("%s-wepass-info.log", dateStr))
	errorLogFile := filepath.Join(logDir, fmt.Sprintf("%s-wepass-error.log", dateStr))
	debugLogFile := filepath.Join(logDir, fmt.Sprintf("%s-wepass-debug.log", dateStr))

	addInitLog("[日志初始化] Info日志文件: %s", infoLogFile)
	addInitLog("[日志初始化] Error日志文件: %s", errorLogFile)
	addInitLog("[日志初始化] Debug日志文件: %s", debugLogFile)

	// 打开日志文件
	addInitLog("[日志初始化] 正在打开Info日志文件: %s", infoLogFile)
	var err error
	infoFile, err = os.OpenFile(infoLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		addInitLog("[日志初始化] ❌ 创建info日志文件失败: %v", err)
		// 尝试检查父目录权限
		if parentDir := filepath.Dir(infoLogFile); parentDir != "" {
			if info, statErr := os.Stat(parentDir); statErr == nil {
				addInitLog("[日志初始化] 父目录权限: %s", info.Mode().String())
			}
		}
		return fmt.Errorf("创建info日志文件失败: %w", err)
	}
	addInitLog("[日志初始化] ✅ Info日志文件打开成功")

	// 检查文件权限
	if info, statErr := os.Stat(infoLogFile); statErr == nil {
		addInitLog("[日志初始化] Info文件权限: %s", info.Mode().String())
	}

	addInitLog("[日志初始化] 正在打开Error日志文件: %s", errorLogFile)
	errorFile, err = os.OpenFile(errorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		addInitLog("[日志初始化] ❌ 创建error日志文件失败: %v", err)
		infoFile.Close() // 清理已打开的文件
		return fmt.Errorf("创建error日志文件失败: %w", err)
	}
	addInitLog("[日志初始化] ✅ Error日志文件打开成功")

	// 检查文件权限
	if info, statErr := os.Stat(errorLogFile); statErr == nil {
		addInitLog("[日志初始化] Error文件权限: %s", info.Mode().String())
	}

	addInitLog("[日志初始化] 正在打开Debug日志文件: %s", debugLogFile)
	debugFile, err = os.OpenFile(debugLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		addInitLog("[日志初始化] ❌ 创建debug日志文件失败: %v", err)
		infoFile.Close()  // 清理已打开的文件
		errorFile.Close() // 清理已打开的文件
		return fmt.Errorf("创建debug日志文件失败: %w", err)
	}
	addInitLog("[日志初始化] ✅ Debug日志文件打开成功")

	// 检查文件权限
	if info, statErr := os.Stat(debugLogFile); statErr == nil {
		addInitLog("[日志初始化] Debug文件权限: %s", info.Mode().String())
	}

	// 测试文件写入
	addInitLog("[日志初始化] 正在测试文件写入...")
	testMsg := fmt.Sprintf("日志系统初始化测试 - %s\n", time.Now().Format("2006-01-02 15:04:05"))

	if _, err := infoFile.WriteString(testMsg); err != nil {
		addInitLog("[日志初始化] ❌ Info文件写入测试失败: %v", err)
	} else {
		addInitLog("[日志初始化] ✅ Info文件写入测试成功")
	}

	if _, err := errorFile.WriteString(testMsg); err != nil {
		addInitLog("[日志初始化] ❌ Error文件写入测试失败: %v", err)
	} else {
		addInitLog("[日志初始化] ✅ Error文件写入测试成功")
	}

	if _, err := debugFile.WriteString(testMsg); err != nil {
		addInitLog("[日志初始化] ❌ Debug文件写入测试失败: %v", err)
	} else {
		addInitLog("[日志初始化] ✅ Debug文件写入测试成功")
	}

	// 强制刷新文件缓冲区
	infoFile.Sync()
	errorFile.Sync()
	debugFile.Sync()
	addInitLog("[日志初始化] ✅ 文件缓冲区已刷新")

	// 20251003 陈凤庆 将所有初始化日志写入到文件中
	fmt.Printf("[日志初始化] 正在将初始化日志写入文件...\n")
	initLogContent := ""
	for _, logLine := range initLogs {
		initLogContent += logLine + "\n"
	}

	// 写入到所有日志文件
	if _, err := infoFile.WriteString(initLogContent); err == nil {
		infoFile.Sync()
		fmt.Printf("[日志初始化] ✅ 初始化日志已写入Info文件\n")
	} else {
		fmt.Printf("[日志初始化] ❌ 初始化日志写入Info文件失败: %v\n", err)
	}
	if _, err := errorFile.WriteString(initLogContent); err == nil {
		errorFile.Sync()
		fmt.Printf("[日志初始化] ✅ 初始化日志已写入Error文件\n")
	} else {
		fmt.Printf("[日志初始化] ❌ 初始化日志写入Error文件失败: %v\n", err)
	}
	if _, err := debugFile.WriteString(initLogContent); err == nil {
		debugFile.Sync()
		fmt.Printf("[日志初始化] ✅ 初始化日志已写入Debug文件\n")
	} else {
		fmt.Printf("[日志初始化] ❌ 初始化日志写入Debug文件失败: %v\n", err)
	}

	// 创建多重写入器（同时写入文件和控制台）
	infoWriter := io.MultiWriter(os.Stdout, infoFile)
	errorWriter := io.MultiWriter(os.Stderr, errorFile)
	debugWriter := io.MultiWriter(os.Stdout, debugFile)

	// 创建日志记录器
	// 20251009 陈凤庆 修改日志格式：不使用标准时间格式，改为自定义格式以符合Windows平台要求
	InfoLogger = log.New(infoWriter, "", 0)
	ErrorLogger = log.New(errorWriter, "", 0)
	DebugLogger = log.New(debugWriter, "", 0)

	addInitLog("[日志初始化] ✅ 日志记录器创建完成")

	// 测试日志记录器
	InfoLogger.Printf("日志系统初始化完成，日志目录: %s", logDir)
	addInitLog("[日志初始化] 🎉 日志系统初始化完成！")

	return nil
}

/**
 * Info 记录信息日志
 * @param format 格式字符串
 * @param v 参数
 * @modify 20251003 陈凤庆 添加配置控制，可以关闭Info日志
 */
func Info(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 20251003 陈凤庆 检查是否启用Info日志
	if !logConfig.EnableInfoLog {
		return
	}

	if InfoLogger != nil {
		// 20251009 陈凤庆 修改为Windows平台要求的日志格式
		_, file, _, ok := runtime.Caller(1)
		if ok {
			file = filepath.Base(file) // 只保留文件名
		} else {
			file = "unknown"
		}

		// 20251009 陈凤庆 生成符合Windows平台要求的时间格式：2025-10-08T14:17:46.000Z
		now := time.Now().UTC()
		timeStr := now.Format("2006-01-02T15:04:05.000Z")

		// 20251009 陈凤庆 确定模块名称（基于文件名）
		moduleName := getModuleName(file)

		message := fmt.Sprintf(format, v...)

		// 20251009 陈凤庆 Windows平台日志格式：时间: xxx, 日志级别: xxx, 来源/模块: xxx, 消息: xxx
		logLine := fmt.Sprintf("时间: %s, 日志级别: INFO, 来源/模块: %s, 消息: %s", timeStr, moduleName, message)
		InfoLogger.Printf("%s", logLine)

		// 20251003 陈凤庆 强制刷新文件缓冲区，确保Windows平台日志立即写入
		if infoFile != nil {
			if err := infoFile.Sync(); err != nil {
				// 如果同步失败，输出到控制台作为备选
				log.Printf("[INFO-SYNC-ERROR] %s | 原始日志: %s", err.Error(), message)
			}
		}
	} else {
		log.Printf("[INFO] "+format, v...)
	}
}

/**
 * Error 记录错误日志
 * @param format 格式字符串
 * @param v 参数
 */
func Error(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if ErrorLogger != nil {
		// 20251009 陈凤庆 修改为Windows平台要求的日志格式
		_, file, _, ok := runtime.Caller(1)
		if ok {
			file = filepath.Base(file) // 只保留文件名
		} else {
			file = "unknown"
		}

		// 20251009 陈凤庆 生成符合Windows平台要求的时间格式：2025-10-08T14:17:46.000Z
		now := time.Now().UTC()
		timeStr := now.Format("2006-01-02T15:04:05.000Z")

		// 20251009 陈凤庆 确定模块名称（基于文件名）
		moduleName := getModuleName(file)

		message := fmt.Sprintf(format, v...)

		// 20251009 陈凤庆 Windows平台日志格式：时间: xxx, 日志级别: xxx, 来源/模块: xxx, 消息: xxx
		logLine := fmt.Sprintf("时间: %s, 日志级别: ERROR, 来源/模块: %s, 消息: %s", timeStr, moduleName, message)
		ErrorLogger.Printf("%s", logLine)

		// 20251003 陈凤庆 强制刷新文件缓冲区，确保Windows平台日志立即写入
		if errorFile != nil {
			if err := errorFile.Sync(); err != nil {
				// 如果同步失败，输出到控制台作为备选
				log.Printf("[ERROR-SYNC-ERROR] %s | 原始日志: %s", err.Error(), message)
			}
		}
	} else {
		log.Printf("[ERROR] "+format, v...)
	}
}

/**
 * Debug 记录调试日志
 * @param format 格式字符串
 * @param v 参数
 * @modify 20251003 陈凤庆 添加配置控制，可以关闭Debug日志
 */
func Debug(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 20251003 陈凤庆 检查是否启用Debug日志
	if !logConfig.EnableDebugLog {
		return
	}

	if DebugLogger != nil {
		// 20251009 陈凤庆 修改为Windows平台要求的日志格式
		_, file, _, ok := runtime.Caller(1)
		if ok {
			file = filepath.Base(file) // 只保留文件名
		} else {
			file = "unknown"
		}

		// 20251009 陈凤庆 生成符合Windows平台要求的时间格式：2025-10-08T14:17:46.000Z
		now := time.Now().UTC()
		timeStr := now.Format("2006-01-02T15:04:05.000Z")

		// 20251009 陈凤庆 确定模块名称（基于文件名）
		moduleName := getModuleName(file)

		message := fmt.Sprintf(format, v...)

		// 20251009 陈凤庆 Windows平台日志格式：时间: xxx, 日志级别: xxx, 来源/模块: xxx, 消息: xxx
		logLine := fmt.Sprintf("时间: %s, 日志级别: DEBUG, 来源/模块: %s, 消息: %s", timeStr, moduleName, message)
		DebugLogger.Printf("%s", logLine)

		// 20251003 陈凤庆 强制刷新文件缓冲区，确保Windows平台日志立即写入
		if debugFile != nil {
			if err := debugFile.Sync(); err != nil {
				// 如果同步失败，输出到控制台作为备选
				log.Printf("[DEBUG-SYNC-ERROR] %s | 原始日志: %s", err.Error(), message)
			}
		}
	} else {
		log.Printf("[DEBUG] "+format, v...)
	}
}

/**
 * LogPasswordOperation 记录密码操作日志
 * @param operation 操作类型
 * @param groupID 分组ID
 * @param details 详细信息
 */
func LogPasswordOperation(operation string, groupID string, details string) {
	Info("[密码操作] %s - 分组ID: %s - %s", operation, groupID, details)
}

/**
 * LogDatabaseOperation 记录数据库操作日志
 * @param operation 操作类型
 * @param table 表名
 * @param details 详细信息
 */
func LogDatabaseOperation(operation string, table string, details string) {
	Info("[数据库] %s - 表: %s - %s", operation, table, details)
}

/**
 * LogAPICall 记录API调用日志
 * @param method API方法名
 * @param params 参数
 * @param result 结果
 */
func LogAPICall(method string, params string, result string) {
	Info("[API] %s - 参数: %s - 结果: %s", method, params, result)
}

/**
 * CheckLoggerHealth 检查日志系统健康状态
 * @return map[string]interface{} 健康状态信息
 * @description 20251003 陈凤庆 添加日志系统健康检查，帮助诊断日志问题
 */
func CheckLoggerHealth() map[string]interface{} {
	health := make(map[string]interface{})

	// 检查日志记录器是否初始化
	health["info_logger_initialized"] = InfoLogger != nil
	health["error_logger_initialized"] = ErrorLogger != nil
	health["debug_logger_initialized"] = DebugLogger != nil

	// 检查文件句柄是否有效
	health["info_file_valid"] = infoFile != nil
	health["error_file_valid"] = errorFile != nil
	health["debug_file_valid"] = debugFile != nil

	// 尝试写入测试
	if infoFile != nil {
		testMsg := fmt.Sprintf("健康检查测试 - %s\n", time.Now().Format("2006-01-02 15:04:05"))
		if _, err := infoFile.WriteString(testMsg); err != nil {
			health["info_write_test"] = fmt.Sprintf("失败: %v", err)
		} else {
			health["info_write_test"] = "成功"
			infoFile.Sync() // 立即刷新
		}
	} else {
		health["info_write_test"] = "文件句柄无效"
	}

	return health
}

/**
 * LogHealthCheck 记录日志系统健康检查结果
 * @description 20251003 陈凤庆 定期检查日志系统状态
 */
func LogHealthCheck() {
	health := CheckLoggerHealth()
	fmt.Printf("[日志健康检查] 状态: %+v\n", health)
}
