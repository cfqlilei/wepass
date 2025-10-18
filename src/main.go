package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"wepassword/internal/app"
	"wepassword/internal/config"
	"wepassword/internal/database"
	"wepassword/internal/logger"
	"wepassword/internal/version"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// getLogDirectory 获取日志目录路径，自动检测正确的工作目录
// 20251002 陈凤庆 解决双击启动时工作目录不正确的问题
// 20251003 陈凤庆 优化Windows平台日志目录获取逻辑，存放到exe同目录
// 20251003 陈凤庆 修复Windows平台日志权限和路径问题
// 20251004 陈凤庆 优化开发模式日志目录优先级，确保开发环境日志存放在工作区
func getLogDirectory() string {
	// 20251004 陈凤庆 最高优先级：检查是否在开发环境（wails dev模式）
	// 如果是开发模式，直接使用项目根目录的logs，不再尝试其他策略
	workDir, _ := os.Getwd()
	isDevMode := false

	// 方法1：检查可执行文件路径是否包含wails相关路径
	if execPath, err := os.Executable(); err == nil {
		if strings.Contains(execPath, "wails") || strings.Contains(execPath, "build/bin") {
			isDevMode = true
		}
	}

	// 方法2：检查环境变量
	if os.Getenv("WAILS_DEV") != "" {
		isDevMode = true
	}

	// 方法3：检查是否存在 wails.json 文件（开发环境标志）
	if _, err := os.Stat(filepath.Join(workDir, "wails.json")); err == nil {
		isDevMode = true
	}

	// 如果是开发模式，优先使用项目根目录的logs
	if isDevMode {
		var projectRoot string

		// 检查当前工作目录是否是 src 目录
		if filepath.Base(workDir) == "src" {
			// 如果在 src 目录，回到项目根目录
			projectRoot = filepath.Dir(workDir)
		} else {
			srcDir := filepath.Join(workDir, "src")
			if _, err := os.Stat(srcDir); err == nil {
				// 当前目录是项目根目录
				projectRoot = workDir
			} else {
				// 尝试查找项目根目录（向上查找包含 src 目录的父目录）
				currentDir := workDir
				for i := 0; i < 5; i++ { // 最多向上查找5级
					srcDir := filepath.Join(currentDir, "src")
					if _, err := os.Stat(srcDir); err == nil {
						// 找到项目根目录
						projectRoot = currentDir
						break
					}

					// 向上一级目录
					parentDir := filepath.Dir(currentDir)
					if parentDir == currentDir {
						// 已经到达根目录
						break
					}
					currentDir = parentDir
				}
			}
		}

		// 如果找到了项目根目录，使用其logs子目录
		if projectRoot != "" {
			logDir := filepath.Join(projectRoot, "logs")
			// 确保日志目录存在
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] 🚀 开发环境模式: %s\n", logDir)
					return logDir
				}
			}
		}

		// 如果无法使用项目根目录，在当前目录创建 logs
		logDir := filepath.Join(workDir, "logs")
		if err := os.MkdirAll(logDir, 0755); err == nil {
			// 测试写入权限
			testFile := filepath.Join(logDir, "test-write.tmp")
			if file, err := os.Create(testFile); err == nil {
				file.WriteString("test")
				file.Close()
				os.Remove(testFile) // 清理测试文件
				fmt.Printf("[日志目录] 🚀 开发环境模式(当前目录): %s\n", logDir)
				return logDir
			}
		}
	}

	// 20251003 陈凤庆 Windows平台使用多种策略确保日志目录可用
	if runtime.GOOS == "windows" {
		// 策略1：尝试使用exe同目录下的logs文件夹
		execPath, err := os.Executable()
		if err == nil {
			// 获取可执行文件所在目录
			execDir := filepath.Dir(execPath)
			logDir := filepath.Join(execDir, "logs")

			// 尝试创建日志目录并测试写入权限
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] Windows策略1成功: %s\n", logDir)
					return logDir
				}
			}
		}

		// 策略2：使用用户文档目录下的wepass文件夹
		userHomeDir, err := os.UserHomeDir()
		if err == nil {
			logDir := filepath.Join(userHomeDir, "Documents", "wepass", "logs")
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] Windows策略2成功: %s\n", logDir)
					return logDir
				}
			}
		}

		// 策略3：使用AppData目录
		appDataDir := os.Getenv("APPDATA")
		if appDataDir != "" {
			logDir := filepath.Join(appDataDir, "wepass", "logs")
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] Windows策略3成功: %s\n", logDir)
					return logDir
				}
			}
		}

		// 策略4：使用当前工作目录
		workDir, err := os.Getwd()
		if err == nil {
			logDir := filepath.Join(workDir, "logs")
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] Windows策略4成功: %s\n", logDir)
					return logDir
				}
			}
		}

		// 策略5：使用临时目录作为最后的备选
		tempDir := os.TempDir()
		logDir := filepath.Join(tempDir, "wepass-logs")
		if err := os.MkdirAll(logDir, 0755); err == nil {
			fmt.Printf("[日志目录] Windows策略5(临时目录)成功: %s\n", logDir)
			return logDir
		}

		fmt.Printf("[日志目录] ❌ Windows所有策略都失败，使用默认临时目录\n")
		return filepath.Join(os.TempDir(), "wepass-logs")
	}

	// 20251003 陈凤庆 Mac平台使用多种策略确保日志目录可用
	if runtime.GOOS == "darwin" {
		// 策略1：使用用户Library/Logs目录（Mac标准做法）
		userHomeDir, err := os.UserHomeDir()
		if err == nil {
			logDir := filepath.Join(userHomeDir, "Library", "Logs", "wepass")
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] Mac策略1成功: %s\n", logDir)
					return logDir
				}
			}
		}

		// 策略2：使用用户Documents目录
		if userHomeDir != "" {
			logDir := filepath.Join(userHomeDir, "Documents", "wepass", "logs")
			if err := os.MkdirAll(logDir, 0755); err == nil {
				// 测试写入权限
				testFile := filepath.Join(logDir, "test-write.tmp")
				if file, err := os.Create(testFile); err == nil {
					file.WriteString("test")
					file.Close()
					os.Remove(testFile) // 清理测试文件
					fmt.Printf("[日志目录] Mac策略2成功: %s\n", logDir)
					return logDir
				}
			}
		}

		// 策略3：尝试使用app bundle同目录（如果是从app运行）
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			// 检查是否在app bundle内
			if strings.Contains(execDir, ".app/Contents/MacOS") {
				// 在app bundle内，使用app同级目录
				appDir := filepath.Dir(filepath.Dir(filepath.Dir(execDir))) // 回到.app的父目录
				logDir := filepath.Join(appDir, "wepass-logs")
				if err := os.MkdirAll(logDir, 0755); err == nil {
					// 测试写入权限
					testFile := filepath.Join(logDir, "test-write.tmp")
					if file, err := os.Create(testFile); err == nil {
						file.WriteString("test")
						file.Close()
						os.Remove(testFile) // 清理测试文件
						fmt.Printf("[日志目录] Mac策略3成功: %s\n", logDir)
						return logDir
					}
				}
			} else {
				// 不在app bundle内，使用可执行文件同目录
				logDir := filepath.Join(execDir, "logs")
				if err := os.MkdirAll(logDir, 0755); err == nil {
					// 测试写入权限
					testFile := filepath.Join(logDir, "test-write.tmp")
					if file, err := os.Create(testFile); err == nil {
						file.WriteString("test")
						file.Close()
						os.Remove(testFile) // 清理测试文件
						fmt.Printf("[日志目录] Mac策略3b成功: %s\n", logDir)
						return logDir
					}
				}
			}
		}
	}

	// 通用策略：适用于所有平台（非开发环境）
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "."
	}

	// 如果无法找到项目根目录，在当前目录创建 logs
	logDir := filepath.Join(workDir, "logs")

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// 如果无法创建，使用临时目录
		fmt.Printf("[日志目录] 所有策略失败，使用临时目录: %s\n", filepath.Join(os.TempDir(), "wepass-logs"))
		return filepath.Join(os.TempDir(), "wepass-logs")
	}

	fmt.Printf("[日志目录] 通用策略成功: %s\n", logDir)
	return logDir
}

/**
 * wepass 密码管理工具主程序入口
 * @author 陈凤庆
 * @date 2025-01-27
 * @description 基于 Wails + Go + Vue.js 的跨平台密码管理工具
 */
func main() {
	// 20251002 陈凤庆 调试信息：记录启动状态到固定位置
	execPath, _ := os.Executable()
	workDir, _ := os.Getwd()
	debugLogPath := "/tmp/wepass-debug.log"
	debugFile, _ := os.Create(debugLogPath)
	if debugFile != nil {
		fmt.Fprintf(debugFile, "启动时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(debugFile, "可执行文件路径: %s\n", execPath)
		fmt.Fprintf(debugFile, "当前工作目录: %s\n", workDir)
		debugFile.Close()
	}

	// 初始化日志系统，自动检测正确的工作目录
	logDir := getLogDirectory()

	// 更新调试信息
	debugFile, _ = os.OpenFile(debugLogPath, os.O_APPEND|os.O_WRONLY, 0644)
	if debugFile != nil {
		fmt.Fprintf(debugFile, "检测到的日志目录: %s\n", logDir)
		debugFile.Close()
	}

	if err := logger.InitLogger(logDir); err != nil {
		// 记录错误到调试文件
		debugFile, _ = os.OpenFile(debugLogPath, os.O_APPEND|os.O_WRONLY, 0644)
		if debugFile != nil {
			fmt.Fprintf(debugFile, "日志系统初始化失败: %v\n", err)
			debugFile.Close()
		}
		log.Fatalf("初始化日志系统失败: %v", err)
	}

	// 20251003 陈凤庆 执行日志系统健康检查
	logger.LogHealthCheck()

	logger.Info("=== wepass 密码管理器启动 ===")

	// 初始化配置管理器
	configManager := config.NewConfigManager()
	logger.Info("配置管理器初始化完成")

	// 20251003 陈凤庆 从配置文件加载日志配置
	logConfig := configManager.GetLogConfig()
	logger.SetLogConfig(logger.LogConfig{
		EnableInfoLog:  logConfig.EnableInfoLog,
		EnableDebugLog: logConfig.EnableDebugLog,
	})
	logger.Info("日志配置已加载: Info=%v, Debug=%v", logConfig.EnableInfoLog, logConfig.EnableDebugLog)

	// 初始化数据库管理器
	dbManager := database.NewDatabaseManager()
	logger.Info("数据库管理器初始化完成")

	// 创建应用实例
	app := app.NewApp(configManager, dbManager)
	logger.Info("应用实例创建完成")

	// 创建应用选项
	// 20250127 陈凤庆 添加AlwaysOnTop选项实现全局浮动效果
	err := wails.Run(&options.App{
		Title:  version.GetFullVersion(),
		Width:  550, // 根据设计要求调整窗口宽度
		Height: 800, // 根据设计要求调整窗口高度
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnDomReady:       app.DomReady,
		OnBeforeClose:    app.BeforeClose,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Frameless:         false,
		MinWidth:          400,  // 20250928 陈凤庆 调整最小宽度，支持窗口调整
		MinHeight:         600,  // 20250928 陈凤庆 调整最小高度，支持窗口调整
		MaxWidth:          1200, // 20250928 陈凤庆 设置最大宽度，支持窗口调整
		MaxHeight:         0,    // 20251005 陈凤庆 设置为0表示不限制最大高度，允许使用显示器可视高度
		StartHidden:       false,
		HideWindowOnClose: false,
		AlwaysOnTop:       true, // 20250127 陈凤庆 设置窗口始终置顶，实现全局浮动效果
	})

	if err != nil {
		log.Fatal("启动应用失败:", err)
	}
}
