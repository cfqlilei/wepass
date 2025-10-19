#!/bin/bash

# wepass 密码管理工具 - 统一构建脚本
# 作者：陈凤庆
# 创建时间：2025年9月28日
# 修改时间：2025年10月2日
# @modify 陈凤庆 整合多个构建脚本为统一脚本，支持更多平台和选项

set -e

# 版本信息
VERSION="1.0.7"
APP_NAME="wepass"
BUILD_DATE=$(date +%Y%m%d)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 全局变量
COMPRESS=""
WITH_DATE=""
CUSTOM_VERSION=""
WORKING_DIR=""

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 命令未找到，请先安装 $1"
        exit 1
    fi
}

# 检查 Wails CLI
check_wails() {
    if ! command -v wails &> /dev/null; then
        if [ -f "$(go env GOPATH)/bin/wails" ]; then
            export PATH="$PATH:$(go env GOPATH)/bin"
            log_info "已将 Go bin 目录添加到 PATH"
        else
            log_error "Wails CLI 未找到，请运行: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
            exit 1
        fi
    fi
}

# 显示帮助信息
show_help() {
    echo "wepass 密码管理工具 - 统一构建脚本"
    echo ""
    echo "用法: $0 [平台] [选项]"
    echo ""
    echo "平台选项:"
    echo "  all              编译所有支持的平台版本"
    echo "  mac              编译 macOS 版本 (包含 X64 和 arm64)"
    echo "  mac-amd64        编译 macOS Intel 版本"
    echo "  mac-arm64        编译 macOS Apple Silicon 版本"
    echo "  mac-universal    编译 macOS 通用版本"
    echo "  win              编译 Windows 版本 (默认 amd64)"
    echo "  win-amd64        编译 Windows 64位版本"
    echo "  win-386          编译 Windows 32位版本"
    echo "  linux            编译 Linux 版本 (默认 amd64)"
    echo "  linux-amd64      编译 Linux 64位版本"
    echo ""
    echo "其他选项:"
    echo "  -h, --help       显示此帮助信息"
    echo "  -c, --clean      清理构建目录"
    echo "  -d, --dev        启动开发模式"
    echo "  -t, --test       运行测试"
    echo "  --compress       启用压缩"
    echo "  --version VER    指定版本号 (默认: $VERSION)"
    echo "  --with-date      在输出文件名中包含日期"
    echo ""
    echo "示例:"
    echo "  $0 all                    # 编译所有平台"
    echo "  $0 mac                    # 编译 macOS 所有架构"
    echo "  $0 mac-amd64             # 仅编译 macOS Intel 版本"
    echo "  $0 win --compress         # 编译压缩的 Windows 版本"
    echo "  $0 all --with-date        # 编译所有平台并在文件名中包含日期"
    echo "  $0 linux --version 1.0.2  # 编译 Linux 版本并指定版本号"
}

# 清理构建目录
clean_build() {
    log_info "清理构建目录..."

    # 清理根目录的构建目录
    if [ -d "build" ]; then
        rm -rf build
        log_success "构建目录已清理"
    else
        log_info "构建目录不存在，跳过清理"
    fi
}

# 进入工作目录
enter_working_dir() {
    if [ -d "src" ]; then
        log_info "进入源码目录..."
        cd src
        WORKING_DIR="src"
    fi
}

# 返回原目录
exit_working_dir() {
    if [ "$WORKING_DIR" = "src" ]; then
        cd ..
        WORKING_DIR=""
    fi
}

# 安装依赖
install_dependencies() {
    log_info "检查并安装依赖..."

    # 进入工作目录
    enter_working_dir

    # 安装 Go 依赖
    log_info "安装 Go 模块依赖..."
    go mod tidy

    # 安装前端依赖
    if [ -d "frontend" ]; then
        log_info "安装前端依赖..."
        cd frontend
        npm install
        cd ..
    fi

    log_success "依赖安装完成"
}

# 构建前端
build_frontend() {
    if [ -d "frontend" ]; then
        log_info "构建前端资源..."
        cd frontend

        # 检查是否需要安装依赖
        if [ ! -d "node_modules" ]; then
            log_info "安装前端依赖..."
            npm install
        fi

        npm run build
        cd ..
        log_success "前端构建完成"
    else
        log_warning "前端目录不存在，跳过前端构建"
    fi
}

# 生成输出文件名
generate_output_name() {
    local platform=$1
    local arch=$2
    local base_name="$APP_NAME"

    # 使用自定义版本号或默认版本号
    local version_str="${CUSTOM_VERSION:-$VERSION}"

    # 根据平台确定文件扩展名
    case $platform in
        "darwin")
            local ext=".app"
            ;;
        "windows")
            local ext=".exe"
            ;;
        *)
            local ext=""
            ;;
    esac

    # 构建文件名：应用名称-版本号.后缀名
    if [ "$WITH_DATE" = "true" ]; then
        if [ -n "$arch" ] && [ "$arch" != "amd64" ]; then
            echo "${BUILD_DATE}-${base_name}-${version_str}-${arch}${ext}"
        else
            echo "${BUILD_DATE}-${base_name}-${version_str}${ext}"
        fi
    else
        if [ -n "$arch" ] && [ "$arch" != "amd64" ]; then
            echo "${base_name}-${version_str}-${arch}${ext}"
        else
            echo "${base_name}-${version_str}${ext}"
        fi
    fi
}

# 编译指定平台
build_platform() {
    local platform_arch=$1
    local compress_flag=""
    local build_flags=""

    if [ "$COMPRESS" = "true" ]; then
        compress_flag="-upx"
    fi

    # 解析平台和架构
    local platform=$(echo $platform_arch | cut -d'/' -f1)
    local arch=$(echo $platform_arch | cut -d'/' -f2)

    log_info "编译 $platform_arch 版本..."

    # 设置环境变量处理交叉编译
    export CGO_ENABLED=1
    case $platform_arch in
        "windows/amd64")
            # Windows 交叉编译需要特殊处理
            if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "win32" ]]; then
                log_warning "尝试在非 Windows 系统上编译 Windows 版本"
                log_info "检查是否安装了 MinGW-w64 交叉编译工具链..."

                # 检查 MinGW-w64 是否可用
                if command -v x86_64-w64-mingw32-gcc &> /dev/null; then
                    log_info "发现 MinGW-w64 工具链，尝试交叉编译"
                    export CC=x86_64-w64-mingw32-gcc
                    export CXX=x86_64-w64-mingw32-g++
                    export CGO_ENABLED=1
                else
                    log_error "Windows 版本需要 CGO 支持，但未找到 MinGW-w64 工具链"
                    log_error "请安装 MinGW-w64 或在 Windows 系统上编译"
                    log_info "安装命令："
                    log_info "  macOS: brew install mingw-w64"
                    log_info "  Ubuntu/Debian: sudo apt-get install gcc-mingw-w64-x86-64"
                    return 1
                fi
            fi
            ;;
        "linux/amd64")
            # Linux 交叉编译
            if [[ "$OSTYPE" != "linux-gnu"* ]]; then
                log_warning "从非 Linux 系统交叉编译到 Linux，禁用 CGO"
                export CGO_ENABLED=0
            fi
            ;;
    esac

    # 执行编译
    wails build -platform $platform_arch $compress_flag
    local build_result=$?

    if [ $build_result -eq 0 ]; then
        # 移动文件到根目录的 build 目录并重命名
        move_and_rename_output_file "$platform" "$arch"

        case $platform_arch in
            "darwin/amd64")
                log_success "macOS Intel 版本编译完成"
                ;;
            "darwin/arm64")
                log_success "macOS Apple Silicon 版本编译完成"
                ;;
            "darwin/universal")
                log_success "macOS 通用版本编译完成"
                ;;
            "windows/amd64")
                log_success "Windows 64位版本编译完成"
                ;;
            "windows/386")
                log_success "Windows 32位版本编译完成"
                ;;
            "linux/amd64")
                log_success "Linux 64位版本编译完成"
                ;;
            *)
                log_success "$platform_arch 版本编译完成"
                ;;
        esac
    else
        log_error "$platform_arch 版本编译失败"
        return 1
    fi
}

# 对 macOS 应用进行代码签名
sign_macos_app() {
    local app_path="$1"

    if [[ "$OSTYPE" == "darwin"* ]] && [ -d "$app_path" ]; then
        log_info "对 macOS 应用进行代码签名: $(basename "$app_path")"

        # 检查是否有开发者证书（包括自签名证书）
        local dev_cert=$(security find-identity -v -p codesigning | grep -E "(Developer ID Application|wepass Developer Certificate)" | head -1 | cut -d'"' -f2)

        if [ -n "$dev_cert" ]; then
            log_info "使用证书签名: $dev_cert"
            codesign --force --deep --sign "$dev_cert" "$app_path" 2>/dev/null
        else
            log_info "使用 adhoc 签名（本地开发）"
            codesign --force --deep --sign - "$app_path" 2>/dev/null
        fi

        # 验证签名
        if codesign --verify --verbose "$app_path" 2>/dev/null; then
            log_success "代码签名完成并验证成功"
        else
            log_warning "代码签名验证失败，但应用仍可能正常工作"
        fi

        # 移除隔离属性，允许直接启动
        log_info "移除应用隔离属性..."
        xattr -dr com.apple.quarantine "$app_path" 2>/dev/null || true
        xattr -d com.apple.quarantine "$app_path" 2>/dev/null || true

        # 检查 Gatekeeper 状态
        local gatekeeper_result=$(spctl --assess --verbose "$app_path" 2>&1)
        if echo "$gatekeeper_result" | grep -q "accepted"; then
            log_success "应用通过 Gatekeeper 验证，可以直接双击启动"
        else
            log_warning "应用未通过 Gatekeeper 验证"
            log_info "已移除隔离属性，应用应该可以直接双击启动"
            log_info "如果仍然无法启动，请右键点击选择'打开'"
        fi
    fi
}

# 创建 Gatekeeper 绕过脚本
create_gatekeeper_bypass_script() {
    local app_path="$1"
    local script_name="允许wepass运行.command"

    log_info "创建 Gatekeeper 绕过脚本: $script_name"

    cat > "$script_name" << EOF
#!/bin/bash
# wepass Gatekeeper 绕过脚本
# 自动生成于构建时间: $(date)

echo "正在为 wepass 添加 Gatekeeper 例外..."
echo "这将允许应用直接双击启动，无需每次确认"
echo ""

APP_PATH="$app_path"

if [ -d "\$APP_PATH" ]; then
    echo "应用路径: \$APP_PATH"
    echo ""
    echo "请输入管理员密码以添加 Gatekeeper 例外："

    if sudo spctl --add "\$APP_PATH"; then
        echo ""
        echo "✅ 成功添加 Gatekeeper 例外！"
        echo "现在可以直接双击启动 wepass 应用了"
        echo ""
        echo "如果仍然无法启动，请尝试："
        echo "1. 右键点击应用图标 → 选择'打开'"
        echo "2. 或运行启动器脚本：启动wepass.command"
    else
        echo ""
        echo "❌ 添加 Gatekeeper 例外失败"
        echo "请手动右键点击应用图标选择'打开'"
    fi
else
    echo "❌ 错误: 找不到应用文件"
    echo "请确保已经编译了应用: ./build.sh mac-arm64"
fi

echo ""
read -p "按任意键退出..."
EOF

    chmod +x "$script_name"
    log_success "Gatekeeper 绕过脚本已创建: $script_name"
}

# 创建应用启动器
create_app_launcher() {
    local app_path="$1"
    local launcher_name="$2"

    if [ -d "$app_path" ]; then
        log_info "创建应用启动器: $launcher_name"

        # 获取应用的绝对路径
        local abs_app_path=$(realpath "$app_path")
        local executable_path="$abs_app_path/Contents/MacOS/wepass"

        cat > "$launcher_name" << EOF
#!/bin/bash
# wepass 应用启动器
# 自动生成于构建时间: $(date)

# 获取脚本所在目录
SCRIPT_DIR="\$(cd "\$(dirname "\${BASH_SOURCE[0]}")" && pwd)"
APP_PATH="\$SCRIPT_DIR/$app_path/Contents/MacOS/wepass"

echo "正在启动 wepass 密码管理器..."
echo "应用路径: \$APP_PATH"

if [ -f "\$APP_PATH" ]; then
    echo "启动应用..."
    "\$APP_PATH"
else
    echo "错误: 找不到应用文件"
    echo "请确保已经编译了应用: ./build.sh all"
    read -p "按任意键退出..."
fi
EOF

        chmod +x "$launcher_name"
        log_success "启动器已创建: $launcher_name"
    fi
}

# 移动文件到根目录的 build 目录并重命名
move_and_rename_output_file() {
    local platform=$1
    local arch=$2

    # 确定源文件位置（Wails 默认输出位置）
    local src_build_dir="build/bin"

    # 确定目标目录（根目录的 build）
    local target_build_dir="../build/bin"
    if [ "$WORKING_DIR" != "src" ]; then
        target_build_dir="build/bin"
    fi

    # 创建目标目录
    mkdir -p "$target_build_dir"

    # 确定原始文件名
    case $platform in
        "darwin")
            local original_file="$src_build_dir/${APP_NAME}.app"
            ;;
        "windows")
            local original_file="$src_build_dir/${APP_NAME}.exe"
            ;;
        *)
            local original_file="$src_build_dir/${APP_NAME}"
            ;;
    esac

    if [ -e "$original_file" ]; then
        local new_name=$(generate_output_name "$platform" "$arch")
        local target_file="$target_build_dir/$new_name"

        # 移动并重命名文件
        mv "$original_file" "$target_file"
        log_info "文件已移动并重命名为: $new_name"

        # 对 macOS 应用进行代码签名
        if [ "$platform" = "darwin" ]; then
            sign_macos_app "$target_file"

            # 创建启动器（仅为第一个 macOS 应用创建）
            if [ "$arch" = "arm64" ] || ([ "$arch" = "amd64" ] && [ ! -f "../启动wepass.command" ]); then
                local launcher_path="../启动wepass.command"
                if [ "$WORKING_DIR" != "src" ]; then
                    launcher_path="启动wepass.command"
                fi
                create_app_launcher "$target_file" "$launcher_path"
            fi
        fi

        # 清理空的源目录
        if [ -d "$src_build_dir" ] && [ -z "$(ls -A $src_build_dir)" ]; then
            rmdir "$src_build_dir"
        fi
        if [ -d "build" ] && [ -z "$(ls -A build)" ]; then
            rmdir "build"
        fi
    else
        log_warning "未找到构建输出文件: $original_file"
    fi
}

# 显示构建结果
show_build_results() {
    log_info "构建结果:"

    # 构建目录在根目录下
    local build_dir="build/bin"
    local abs_build_dir="$(pwd)/build/bin"

    if [ -d "$build_dir" ]; then
        echo ""
        echo "========================================"
        echo "构建完成！输出文件："
        echo "========================================"
        ls -lh "$build_dir/"
        echo ""

        # 计算总大小
        local total_size=$(du -sh "$build_dir/" | cut -f1)
        log_success "构建完成，总大小: $total_size"

        echo ""
        echo "输出目录: $abs_build_dir"
        echo "构建时间: $(date)"
        echo "========================================"
    else
        log_warning "构建目录不存在: $build_dir"
    fi
}

# 运行测试
run_tests() {
    log_info "运行 Go 测试..."

    # 进入工作目录
    enter_working_dir

    go test ./...

    log_success "测试完成"
}

# 启动开发模式
start_dev() {
    log_info "启动开发模式..."

    # 进入工作目录
    enter_working_dir

    wails dev
}

# 解析平台参数并返回平台架构列表
parse_platform() {
    local platform=$1
    local platforms=()

    case $platform in
        "all")
            platforms=("darwin/amd64" "darwin/arm64" "windows/amd64" "linux/amd64")
            ;;
        "mac")
            platforms=("darwin/amd64" "darwin/arm64")
            ;;
        "mac-amd64")
            platforms=("darwin/amd64")
            ;;
        "mac-arm64")
            platforms=("darwin/arm64")
            ;;
        "mac-universal")
            platforms=("darwin/universal")
            ;;
        "win")
            platforms=("windows/amd64")
            ;;
        "win-amd64")
            platforms=("windows/amd64")
            ;;
        "win-386")
            platforms=("windows/386")
            ;;
        "linux")
            platforms=("linux/amd64")
            ;;
        "linux-amd64")
            platforms=("linux/amd64")
            ;;
        *)
            log_error "不支持的平台: $platform"
            show_help
            exit 1
            ;;
    esac

    echo "${platforms[@]}"
}

# 主函数
main() {
    local action=""
    local platform_param=""

    # 如果没有参数，显示帮助
    if [ $# -eq 0 ]; then
        show_help
        exit 1
    fi

    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                clean_build
                exit 0
                ;;
            -d|--dev)
                action="dev"
                shift
                ;;
            -t|--test)
                action="test"
                shift
                ;;
            --compress)
                COMPRESS="true"
                shift
                ;;
            --with-date)
                WITH_DATE="true"
                shift
                ;;
            --version)
                if [ -n "$2" ]; then
                    CUSTOM_VERSION="$2"
                    shift 2
                else
                    log_error "--version 需要指定版本号"
                    exit 1
                fi
                ;;
            all|mac|mac-*|win|win-*|linux|linux-*)
                if [ -z "$platform_param" ]; then
                    platform_param="$1"
                else
                    log_error "只能指定一个平台参数"
                    exit 1
                fi
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # 检查环境
    log_info "检查编译环境..."
    check_command go
    check_command node
    check_command npm
    check_wails
    log_success "环境检查通过"

    # 执行操作
    case $action in
        "dev")
            start_dev
            ;;
        "test")
            run_tests
            ;;
        *)
            # 检查是否指定了平台
            if [ -z "$platform_param" ]; then
                log_error "请指定要编译的平台"
                show_help
                exit 1
            fi

            # 解析平台参数
            local platforms_array=($(parse_platform "$platform_param"))

            log_info "开始编译 $platform_param 平台..."
            log_info "目标平台: ${platforms_array[*]}"

            # 安装依赖
            install_dependencies

            # 构建前端
            build_frontend

            # 编译各平台
            local failed_platforms=()
            local success_count=0
            for platform in "${platforms_array[@]}"; do
                if build_platform "$platform"; then
                    ((success_count++))
                else
                    failed_platforms+=("$platform")
                    log_warning "跳过 $platform 平台编译"
                fi
            done

            # 显示编译结果摘要
            log_info "编译完成摘要："
            log_success "成功编译 $success_count 个平台"
            if [ ${#failed_platforms[@]} -gt 0 ]; then
                log_warning "跳过的平台: ${failed_platforms[*]}"
            fi

            # 返回原目录
            exit_working_dir

            # 显示结果
            show_build_results
            ;;
    esac
}

# 运行主函数
main "$@"