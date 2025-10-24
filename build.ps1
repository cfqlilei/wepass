# wepass 密码管理工具 - 统一构建脚本 (PowerShell)
# 作者：陈凤庆
# 创建时间：2025年10月10日
# 基于 build.sh 和 build.bat 的逻辑，提供更强大的 PowerShell 版本

param(
    [Parameter(Position=0)]
    [ValidateSet("all", "mac", "mac-amd64", "mac-arm64", "mac-universal", "win", "win-amd64", "win-386", "linux", "linux-amd64")]
    [string]$Platform,

    [switch]$Help,
    [switch]$Clean,
    [switch]$Dev,
    [switch]$Test,
    [switch]$Compress,
    [switch]$WithDate,
    [string]$Version
)

# 版本信息
$script:VERSION = "1.0.8"
$script:APP_NAME = "wepass"
$script:BUILD_DATE = Get-Date -Format "yyyyMMdd"
$script:WORKING_DIR = ""

# 日志函数
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# 检查命令是否存在
function Test-Command {
    param([string]$Command)

    $cmd = Get-Command $Command -ErrorAction SilentlyContinue
    if (-not $cmd) {
        Write-Error "$Command 命令未找到，请先安装 $Command"
        exit 1
    }
    return $true
}

# 检查 Wails CLI
function Test-Wails {
    $wails = Get-Command "wails" -ErrorAction SilentlyContinue
    if (-not $wails) {
        $goPath = go env GOPATH
        $wailsPath = Join-Path $goPath "bin\wails.exe"

        if (Test-Path $wailsPath) {
            $env:PATH += ";$(Join-Path $goPath 'bin')"
            Write-Info "已将 Go bin 目录添加到 PATH"
        } else {
            Write-Error "Wails CLI 未找到，请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
            exit 1
        }
    }
    return $true
}

# Show help information
function Show-Help {
    Write-Host "wepass Password Manager - Unified Build Script (PowerShell)"
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [platform] [options]"
    Write-Host ""
    Write-Host "Platform options:"
    Write-Host "  all              Build all supported platforms"
    Write-Host "  mac              Build macOS versions (X64 and arm64)"
    Write-Host "  mac-amd64        Build macOS Intel version"
    Write-Host "  mac-arm64        Build macOS Apple Silicon version"
    Write-Host "  mac-universal    Build macOS Universal version"
    Write-Host "  win              Build Windows version (default amd64)"
    Write-Host "  win-amd64        Build Windows 64-bit version"
    Write-Host "  win-386          Build Windows 32-bit version"
    Write-Host "  linux            Build Linux version (default amd64)"
    Write-Host "  linux-amd64      Build Linux 64-bit version"
    Write-Host ""
    Write-Host "Other options:"
    Write-Host "  -Help            Show this help information"
    Write-Host "  -Clean           Clean build directory"
    Write-Host "  -Dev             Start development mode"
    Write-Host "  -Test            Run tests"
    Write-Host "  -Compress        Enable compression"
    Write-Host "  -Version VER     Specify version number (default: $script:VERSION)"
    Write-Host "  -WithDate        Include date in output filename"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\build.ps1 all                    # Build all platforms"
    Write-Host "  .\build.ps1 mac                    # Build all macOS architectures"
    Write-Host "  .\build.ps1 mac-amd64             # Build macOS Intel version only"
    Write-Host "  .\build.ps1 win -Compress         # Build compressed Windows version"
    Write-Host "  .\build.ps1 all -WithDate         # Build all platforms with date"
    Write-Host "  .\build.ps1 linux -Version 1.0.2  # Build Linux version with version"
}

# 清理构建目录
function Clear-BuildDirectory {
    Write-Info "清理构建目录..."

    if (Test-Path "build") {
        Remove-Item "build" -Recurse -Force
        Write-Success "构建目录已清理"
    } else {
        Write-Info "构建目录不存在，跳过清理"
    }
}

# 进入工作目录
function Enter-WorkingDirectory {
    if (Test-Path "src") {
        Write-Info "进入源码目录..."
        Set-Location "src"
        $script:WORKING_DIR = "src"
    }
}

# 返回原目录
function Exit-WorkingDirectory {
    if ($script:WORKING_DIR -eq "src") {
        Set-Location ".."
        $script:WORKING_DIR = ""
    }
}

# 安装依赖
function Install-Dependencies {
    Write-Info "检查并安装依赖..."

    # 进入工作目录
    Enter-WorkingDirectory

    # 安装 Go 依赖
    Write-Info "安装 Go 模块依赖..."
    & go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Go 依赖安装失败"
        exit 1
    }

    # 安装前端依赖
    if (Test-Path "frontend") {
        Write-Info "安装前端依赖..."
        Push-Location "frontend"

        if (-not (Test-Path "node_modules")) {
            Write-Info "安装前端依赖..."
            & npm install
            if ($LASTEXITCODE -ne 0) {
                Write-Error "前端依赖安装失败"
                Pop-Location
                exit 1
            }
        }

        Pop-Location
    }

    Write-Success "依赖安装完成"
}

# 构建前端
function Build-Frontend {
    if (Test-Path "frontend") {
        Write-Info "构建前端资源..."
        Push-Location "frontend"

        # 检查是否需要安装依赖
        if (-not (Test-Path "node_modules")) {
            Write-Info "安装前端依赖..."
            & npm install
            if ($LASTEXITCODE -ne 0) {
                Write-Error "前端依赖安装失败"
                Pop-Location
                exit 1
            }
        }

        & npm run build
        if ($LASTEXITCODE -ne 0) {
            Write-Error "前端构建失败"
            Pop-Location
            exit 1
        }

        Pop-Location
        Write-Success "前端构建完成"
    } else {
        Write-Warning "前端目录不存在，跳过前端构建"
    }
}

# 生成输出文件名
function Get-OutputName {
    param(
        [string]$Platform,
        [string]$Arch
    )

    $baseName = $script:APP_NAME
    $versionStr = if ($Version) { $Version } else { $script:VERSION }

    # 根据平台确定文件扩展名
    $ext = switch ($Platform) {
        "darwin" { ".app" }
        "windows" { ".exe" }
        default { "" }
    }

    # 构建文件名
    if ($WithDate) {
        if ($Arch -and $Arch -ne "amd64") {
            return "$script:BUILD_DATE-$baseName-$versionStr-$Arch$ext"
        } else {
            return "$script:BUILD_DATE-$baseName-$versionStr$ext"
        }
    } else {
        if ($Arch -and $Arch -ne "amd64") {
            return "$baseName-$versionStr-$Arch$ext"
        } else {
            return "$baseName-$versionStr$ext"
        }
    }
}

# 移动文件到根目录的 build 目录并重命名
function Move-AndRenameOutputFile {
    param(
        [string]$Platform,
        [string]$Arch
    )

    # 确定源文件位置（Wails 默认输出位置）
    $srcBuildDir = "build\bin"

    # 确定目标目录（根目录的 build）
    $targetBuildDir = if ($script:WORKING_DIR -eq "src") { "..\build\bin" } else { "build\bin" }

    # 创建目标目录
    if (-not (Test-Path $targetBuildDir)) {
        New-Item -ItemType Directory -Path $targetBuildDir -Force | Out-Null
    }

    # 确定原始文件名
    $originalFile = switch ($Platform) {
        "darwin" { Join-Path $srcBuildDir "$script:APP_NAME.app" }
        "windows" { Join-Path $srcBuildDir "$script:APP_NAME.exe" }
        default { Join-Path $srcBuildDir $script:APP_NAME }
    }

    if (Test-Path $originalFile) {
        $newName = Get-OutputName $Platform $Arch
        $targetFile = Join-Path $targetBuildDir $newName

        # 移动并重命名文件
        Move-Item $originalFile $targetFile -Force
        Write-Info "文件已移动并重命名为: $newName"

        # 清理空的源目录
        if ((Test-Path $srcBuildDir) -and ((Get-ChildItem $srcBuildDir).Count -eq 0)) {
            Remove-Item $srcBuildDir -Force
        }
        if ((Test-Path "build") -and ((Get-ChildItem "build").Count -eq 0)) {
            Remove-Item "build" -Force
        }
    } else {
        Write-Warning "未找到构建输出文件: $originalFile"
    }
}

# 编译指定平台
function Invoke-PlatformBuild {
    param([string]$PlatformArch)

    $compressFlag = if ($Compress) { "-upx" } else { "" }

    # 解析平台和架构
    $parts = $PlatformArch -split "/"
    $platform = $parts[0]
    $arch = $parts[1]

    Write-Info "编译 $PlatformArch 版本..."

    # 设置环境变量处理交叉编译
    $env:CGO_ENABLED = "1"

    # 执行编译
    if ($compressFlag) {
        & wails build -platform $PlatformArch $compressFlag
    } else {
        & wails build -platform $PlatformArch
    }

    if ($LASTEXITCODE -eq 0) {
        # 移动文件到根目录的 build 目录并重命名
        Move-AndRenameOutputFile $platform $arch

        $successMessage = ""
        if ($PlatformArch -eq "darwin/amd64") {
            $successMessage = "macOS Intel version build completed"
        } elseif ($PlatformArch -eq "darwin/arm64") {
            $successMessage = "macOS Apple Silicon version build completed"
        } elseif ($PlatformArch -eq "darwin/universal") {
            $successMessage = "macOS Universal version build completed"
        } elseif ($PlatformArch -eq "windows/amd64") {
            $successMessage = "Windows 64-bit version build completed"
        } elseif ($PlatformArch -eq "windows/386") {
            $successMessage = "Windows 32-bit version build completed"
        } elseif ($PlatformArch -eq "linux/amd64") {
            $successMessage = "Linux 64-bit version build completed"
        } else {
            $successMessage = "$PlatformArch version build completed"
        }
        Write-Success $successMessage
        return $true
    } else {
        Write-Error "$PlatformArch version build failed"
        return $false
    }
}

# 显示构建结果
function Show-BuildResults {
    Write-Info "构建结果:"

    # 构建目录在根目录下
    $buildDir = "build\bin"

    if (Test-Path $buildDir) {
        Write-Host ""
        Write-Host "========================================"
        Write-Host "构建完成！输出文件："
        Write-Host "========================================"
        Get-ChildItem $buildDir | Format-Table Name, Length, LastWriteTime -AutoSize
        Write-Host ""

        # 计算总大小
        $totalSize = (Get-ChildItem $buildDir -Recurse | Measure-Object -Property Length -Sum).Sum
        $totalSizeMB = [math]::Round($totalSize / 1MB, 2)
        Write-Success "构建完成，总大小: $totalSizeMB MB"

        Write-Host ""
        $absBuildDir = (Resolve-Path $buildDir -ErrorAction SilentlyContinue).Path
        Write-Host "输出目录: $absBuildDir"
        Write-Host "构建时间: $(Get-Date)"
        Write-Host "========================================"
    } else {
        Write-Warning "构建目录不存在: $buildDir"
    }
}

# 运行测试
function Invoke-Tests {
    Write-Info "运行 Go 测试..."

    # 进入工作目录
    Enter-WorkingDirectory

    & go test ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "测试失败"
        exit 1
    }

    Write-Success "测试完成"
}

# 启动开发模式
function Start-Development {
    Write-Info "启动开发模式..."

    # 进入工作目录
    Enter-WorkingDirectory

    & wails dev
}

# 解析平台参数并返回平台架构列表
function Get-PlatformList {
    param([string]$Platform)

    switch ($Platform) {
        "all" { return @("darwin/amd64", "darwin/arm64", "windows/amd64", "linux/amd64") }
        "mac" { return @("darwin/amd64", "darwin/arm64") }
        "mac-amd64" { return @("darwin/amd64") }
        "mac-arm64" { return @("darwin/arm64") }
        "mac-universal" { return @("darwin/universal") }
        "win" { return @("windows/amd64") }
        "win-amd64" { return @("windows/amd64") }
        "win-386" { return @("windows/386") }
        "linux" { return @("linux/amd64") }
        "linux-amd64" { return @("linux/amd64") }
        default {
            Write-Error "不支持的平台: $Platform"
            Show-Help
            exit 1
        }
    }
}

# 主函数
function Main {
    # 处理帮助参数
    if ($Help -or (-not $Platform -and -not $Clean -and -not $Dev -and -not $Test)) {
        Show-Help
        return
    }

    # 处理清理参数
    if ($Clean) {
        Clear-BuildDirectory
        return
    }

    # 检查环境
    Write-Info "检查编译环境..."
    Test-Command "go"
    Test-Command "node"
    Test-Command "npm"
    Test-Wails
    Write-Success "环境检查通过"

    # 处理开发模式
    if ($Dev) {
        Start-Development
        return
    }

    # 处理测试
    if ($Test) {
        Invoke-Tests
        return
    }

    # 检查是否指定了平台
    if (-not $Platform) {
        Write-Error "请指定要编译的平台"
        Show-Help
        return
    }

    # 解析平台参数
    $platformsList = Get-PlatformList $Platform

    Write-Info "开始编译 $Platform 平台..."
    Write-Info "目标平台: $($platformsList -join ', ')"

    # 安装依赖
    Install-Dependencies

    # 构建前端
    Build-Frontend

    # 编译各平台
    $failedPlatforms = @()
    $successCount = 0

    foreach ($platformArch in $platformsList) {
        if (Invoke-PlatformBuild $platformArch) {
            $successCount++
        } else {
            $failedPlatforms += $platformArch
            Write-Warning "跳过 $platformArch 平台编译"
        }
    }

    # 显示编译结果摘要
    Write-Info "Build summary:"
    Write-Success "Successfully built $successCount platforms"
    if ($failedPlatforms.Count -gt 0) {
        Write-Warning "Skipped platforms: $($failedPlatforms -join ', ')"
    }

    # 返回原目录
    Exit-WorkingDirectory

    # 显示结果
    Show-BuildResults
}

# 运行主函数
Main