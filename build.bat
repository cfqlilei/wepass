@echo off
REM wepass 密码管理工具 - 统一构建脚本 (Windows)
REM 作者：陈凤庆
REM 创建时间：2025年9月28日
REM 修改时间：2025年10月2日
REM @modify 陈凤庆 整合多个构建脚本为统一脚本，支持更多平台和选项

setlocal enabledelayedexpansion

REM 版本信息
set "VERSION=1.0.8"
set "APP_NAME=wepass"
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "BUILD_DATE=%dt:~0,8%"

REM 颜色定义（Windows 10+ 支持 ANSI 颜色）
set "RED=[31m"
set "GREEN=[32m"
set "YELLOW=[33m"
set "BLUE=[34m"
set "NC=[0m"

REM 全局变量
set "COMPRESS="
set "WITH_DATE="
set "CUSTOM_VERSION="
set "WORKING_DIR="

REM 日志函数
:log_info
echo %BLUE%[INFO]%NC% %~1
goto :eof

:log_success
echo %GREEN%[SUCCESS]%NC% %~1
goto :eof

:log_warning
echo %YELLOW%[WARNING]%NC% %~1
goto :eof

:log_error
echo %RED%[ERROR]%NC% %~1
goto :eof

REM 检查参数
if "%1"=="" goto show_help

REM 解析命令行参数
call :parse_arguments %*
goto main

:parse_arguments
set "PLATFORM_PARAM="
set "ACTION="

:parse_loop
if "%1"=="" goto :eof

if "%1"=="-h" goto show_help
if "%1"=="--help" goto show_help
if "%1"=="-c" goto clean_build
if "%1"=="--clean" goto clean_build
if "%1"=="-d" (
    set "ACTION=dev"
    shift
    goto parse_loop
)
if "%1"=="--dev" (
    set "ACTION=dev"
    shift
    goto parse_loop
)
if "%1"=="-t" (
    set "ACTION=test"
    shift
    goto parse_loop
)
if "%1"=="--test" (
    set "ACTION=test"
    shift
    goto parse_loop
)
if "%1"=="--compress" (
    set "COMPRESS=true"
    shift
    goto parse_loop
)
if "%1"=="--with-date" (
    set "WITH_DATE=true"
    shift
    goto parse_loop
)
if "%1"=="--version" (
    if "%2"=="" (
        call :log_error "--version 需要指定版本号"
        exit /b 1
    )
    set "CUSTOM_VERSION=%2"
    shift
    shift
    goto parse_loop
)

REM 检查是否为平台参数
echo %1 | findstr /r "^all$\|^mac\|^win\|^linux" >nul
if not errorlevel 1 (
    if defined PLATFORM_PARAM (
        call :log_error "只能指定一个平台参数"
        exit /b 1
    )
    set "PLATFORM_PARAM=%1"
    shift
    goto parse_loop
)

call :log_error "未知参数: %1"
goto show_help

:main
REM 检查环境
call :log_info "检查编译环境..."
call :check_environment
call :log_success "环境检查通过"

REM 执行操作
if "%ACTION%"=="dev" goto start_dev
if "%ACTION%"=="test" goto run_tests

REM 检查是否指定了平台
if "%PLATFORM_PARAM%"=="" (
    call :log_error "请指定要编译的平台"
    goto show_help
)

call :log_info "开始编译 %PLATFORM_PARAM% 平台..."

REM 进入工作目录
call :enter_working_dir

REM 安装依赖
call :install_dependencies

REM 构建前端
call :build_frontend

REM 解析并编译平台
call :build_platforms %PLATFORM_PARAM%

REM 返回原目录
call :exit_working_dir

REM 显示结果
call :show_build_results

goto end

:check_environment
where go >nul 2>&1
if errorlevel 1 (
    call :log_error "Go 未安装或不在 PATH 中"
    exit /b 1
)

where node >nul 2>&1
if errorlevel 1 (
    call :log_error "Node.js 未安装或不在 PATH 中"
    exit /b 1
)

where npm >nul 2>&1
if errorlevel 1 (
    call :log_error "npm 未安装或不在 PATH 中"
    exit /b 1
)

REM 检查 Wails
where wails >nul 2>&1
if errorlevel 1 (
    for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
    if exist "!GOPATH!\bin\wails.exe" (
        set "PATH=%PATH%;!GOPATH!\bin"
        call :log_info "已将 Go bin 目录添加到 PATH"
    ) else (
        call :log_error "Wails CLI 未找到，请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit /b 1
    )
)
goto :eof

:enter_working_dir
if exist "src" (
    call :log_info "进入源码目录..."
    cd src
    set "WORKING_DIR=src"
)
goto :eof

:exit_working_dir
if "%WORKING_DIR%"=="src" (
    cd ..
)
goto :eof

:install_dependencies
call :log_info "检查并安装依赖..."

call :log_info "安装 Go 模块依赖..."
go mod tidy
if errorlevel 1 (
    call :log_error "Go 依赖安装失败"
    exit /b 1
)

if exist "frontend" (
    call :log_info "安装前端依赖..."
    cd frontend

    if not exist "node_modules" (
        call :log_info "安装前端依赖..."
        npm install
        if errorlevel 1 (
            call :log_error "前端依赖安装失败"
            cd ..
            exit /b 1
        )
    )

    cd ..
)

call :log_success "依赖安装完成"
goto :eof

:build_frontend
if exist "frontend" (
    call :log_info "构建前端资源..."
    cd frontend
    npm run build
    if errorlevel 1 (
        call :log_error "前端构建失败"
        cd ..
        exit /b 1
    )
    cd ..
    call :log_success "前端构建完成"
) else (
    call :log_warning "前端目录不存在，跳过前端构建"
)
goto :eof

:build_platforms
set "PLATFORM=%1"

if "%PLATFORM%"=="all" (
    call :build_platform_arch "darwin/amd64"
    call :build_platform_arch "darwin/arm64"
    call :build_platform_arch "windows/amd64"
    call :build_platform_arch "linux/amd64"
) else if "%PLATFORM%"=="mac" (
    call :build_platform_arch "darwin/amd64"
    call :build_platform_arch "darwin/arm64"
) else if "%PLATFORM%"=="mac-amd64" (
    call :build_platform_arch "darwin/amd64"
) else if "%PLATFORM%"=="mac-arm64" (
    call :build_platform_arch "darwin/arm64"
) else if "%PLATFORM%"=="mac-universal" (
    call :build_platform_arch "darwin/universal"
) else if "%PLATFORM%"=="win" (
    call :build_platform_arch "windows/amd64"
) else if "%PLATFORM%"=="win-amd64" (
    call :build_platform_arch "windows/amd64"
) else if "%PLATFORM%"=="win-386" (
    call :build_platform_arch "windows/386"
) else if "%PLATFORM%"=="linux" (
    call :build_platform_arch "linux/amd64"
) else if "%PLATFORM%"=="linux-amd64" (
    call :build_platform_arch "linux/amd64"
) else (
    call :log_error "不支持的平台: %PLATFORM%"
    exit /b 1
)
goto :eof

:build_platform_arch
set "PLATFORM_ARCH=%1"
set "COMPRESS_FLAG="

if "%COMPRESS%"=="true" (
    set "COMPRESS_FLAG=-upx"
)

call :log_info "编译 %PLATFORM_ARCH% 版本..."

wails build -platform %PLATFORM_ARCH% %COMPRESS_FLAG%
if errorlevel 1 (
    call :log_error "%PLATFORM_ARCH% 版本编译失败"
    exit /b 1
)

REM 移动文件到根目录的 build 目录并重命名
call :move_and_rename_output_file %PLATFORM_ARCH%

REM 显示成功信息
if "%PLATFORM_ARCH%"=="darwin/amd64" (
    call :log_success "macOS Intel 版本编译完成"
) else if "%PLATFORM_ARCH%"=="darwin/arm64" (
    call :log_success "macOS Apple Silicon 版本编译完成"
) else if "%PLATFORM_ARCH%"=="darwin/universal" (
    call :log_success "macOS 通用版本编译完成"
) else if "%PLATFORM_ARCH%"=="windows/amd64" (
    call :log_success "Windows 64位版本编译完成"
) else if "%PLATFORM_ARCH%"=="windows/386" (
    call :log_success "Windows 32位版本编译完成"
) else if "%PLATFORM_ARCH%"=="linux/amd64" (
    call :log_success "Linux 64位版本编译完成"
) else (
    call :log_success "%PLATFORM_ARCH% 版本编译完成"
)
goto :eof

:move_and_rename_output_file
set "PLATFORM_ARCH=%1"

REM 解析平台和架构
for /f "tokens=1,2 delims=/" %%a in ("%PLATFORM_ARCH%") do (
    set "PLATFORM=%%a"
    set "ARCH=%%b"
)

REM 确定源文件位置（Wails 默认输出位置）
set "SRC_BUILD_DIR=build\bin"

REM 确定目标目录（根目录的 build）
set "TARGET_BUILD_DIR=..\build\bin"
if "%WORKING_DIR%" NEQ "src" (
    set "TARGET_BUILD_DIR=build\bin"
)

REM 创建目标目录
if not exist "%TARGET_BUILD_DIR%" mkdir "%TARGET_BUILD_DIR%"

REM 确定原始文件名
if "%PLATFORM%"=="darwin" (
    set "ORIGINAL_FILE=%SRC_BUILD_DIR%\%APP_NAME%.app"
) else if "%PLATFORM%"=="windows" (
    set "ORIGINAL_FILE=%SRC_BUILD_DIR%\%APP_NAME%.exe"
) else (
    set "ORIGINAL_FILE=%SRC_BUILD_DIR%\%APP_NAME%"
)

if exist "%ORIGINAL_FILE%" (
    call :generate_output_name "%PLATFORM%" "%ARCH%"
    set "TARGET_FILE=%TARGET_BUILD_DIR%\%OUTPUT_NAME%"

    REM 移动并重命名文件
    move "%ORIGINAL_FILE%" "%TARGET_FILE%" >nul
    call :log_info "文件已移动并重命名为: %OUTPUT_NAME%"

    REM 清理空的源目录
    if exist "%SRC_BUILD_DIR%" (
        dir /b "%SRC_BUILD_DIR%" | findstr . >nul || rmdir "%SRC_BUILD_DIR%"
    )
    if exist "build" (
        dir /b "build" | findstr . >nul || rmdir "build"
    )
) else (
    call :log_warning "未找到构建输出文件: %ORIGINAL_FILE%"
)
goto :eof

:generate_output_name
set "PLATFORM=%1"
set "ARCH=%2"
set "BASE_NAME=%APP_NAME%"

REM 使用自定义版本号或默认版本号
if defined CUSTOM_VERSION (
    set "VERSION_STR=%CUSTOM_VERSION%"
) else (
    set "VERSION_STR=%VERSION%"
)

REM 根据平台确定文件扩展名
if "%PLATFORM%"=="darwin" (
    set "EXT=.app"
) else if "%PLATFORM%"=="windows" (
    set "EXT=.exe"
) else (
    set "EXT="
)

REM 构建文件名
if "%WITH_DATE%"=="true" (
    if not "%ARCH%"=="amd64" (
        set "OUTPUT_NAME=%BUILD_DATE%-%BASE_NAME%-%VERSION_STR%-%ARCH%%EXT%"
    ) else (
        set "OUTPUT_NAME=%BUILD_DATE%-%BASE_NAME%-%VERSION_STR%%EXT%"
    )
) else (
    if not "%ARCH%"=="amd64" (
        set "OUTPUT_NAME=%BASE_NAME%-%ARCH%%EXT%"
    ) else (
        set "OUTPUT_NAME=%BASE_NAME%%EXT%"
    )
)
goto :eof

:show_build_results
call :log_info "构建结果:"

REM 构建目录在根目录下
set "BUILD_DIR=..\build\bin"
if "%WORKING_DIR%" NEQ "src" (
    set "BUILD_DIR=build\bin"
)

if exist "%BUILD_DIR%" (
    echo.
    echo ========================================
    echo 构建完成！输出文件：
    echo ========================================
    dir "%BUILD_DIR%"
    echo.

    REM 计算总大小
    for /f "tokens=3" %%a in ('dir "%BUILD_DIR%" /s /-c ^| find "个文件"') do set "TOTAL_SIZE=%%a"
    call :log_success "构建完成，总大小: %TOTAL_SIZE% 字节"

    echo.
    REM 显示绝对路径
    if "%WORKING_DIR%"=="src" (
        for %%i in ("..\build\bin") do echo 输出目录: %%~fi
    ) else (
        for %%i in ("build\bin") do echo 输出目录: %%~fi
    )
    echo 构建时间: %date% %time%
    echo ========================================
) else (
    call :log_warning "构建目录不存在: %BUILD_DIR%"
)
goto :eof

:start_dev
call :log_info "启动开发模式..."
call :enter_working_dir
wails dev
goto end

:run_tests
call :log_info "运行 Go 测试..."
call :enter_working_dir
go test ./...
call :log_success "测试完成"
goto end

:clean_build
call :log_info "清理构建目录..."

REM 清理根目录的构建目录
if exist "build" (
    rmdir /s /q build
    call :log_success "构建目录已清理"
) else (
    call :log_info "构建目录不存在，跳过清理"
)
goto end

:show_help
echo wepass 密码管理工具 - 统一构建脚本 (Windows)
echo.
echo 用法: %0 [平台] [选项]
echo.
echo 平台选项:
echo   all              编译所有支持的平台版本
echo   mac              编译 macOS 版本 (包含 amd64 和 arm64)
echo   mac-amd64        编译 macOS Intel 版本
echo   mac-arm64        编译 macOS Apple Silicon 版本
echo   mac-universal    编译 macOS 通用版本
echo   win              编译 Windows 版本 (默认 amd64)
echo   win-amd64        编译 Windows 64位版本
echo   win-386          编译 Windows 32位版本
echo   linux            编译 Linux 版本 (默认 amd64)
echo   linux-amd64      编译 Linux 64位版本
echo.
echo 其他选项:
echo   -h, --help       显示此帮助信息
echo   -c, --clean      清理构建目录
echo   -d, --dev        启动开发模式
echo   -t, --test       运行测试
echo   --compress       启用压缩
echo   --version VER    指定版本号 (默认: %VERSION%)
echo   --with-date      在输出文件名中包含日期
echo.
echo 示例:
echo   %0 all                    # 编译所有平台
echo   %0 mac                    # 编译 macOS 所有架构
echo   %0 mac-amd64            # 仅编译 macOS Intel 版本
echo   %0 win --compress         # 编译压缩的 Windows 版本
echo   %0 all --with-date        # 编译所有平台并在文件名中包含日期
echo   %0 linux --version 1.0.2  # 编译 Linux 版本并指定版本号
echo.

:end
endlocal