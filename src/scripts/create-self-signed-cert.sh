#!/bin/bash

# wepass 自签名证书创建脚本
# 作者：陈凤庆
# 创建时间：2025年10月3日
# 用途：创建自签名证书用于应用签名，解决双击启动问题

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 证书配置
CERT_NAME="wepass Developer Certificate"
KEYCHAIN_NAME="wepass-dev.keychain"
KEYCHAIN_PASSWORD="wepass123456"

# 检查是否在 macOS 上运行
check_macos() {
    if [[ "$OSTYPE" != "darwin"* ]]; then
        log_error "此脚本只能在 macOS 上运行"
        exit 1
    fi
}

# 检查证书是否已存在
check_existing_cert() {
    if security find-identity -v -p codesigning | grep -q "$CERT_NAME"; then
        log_warning "证书 '$CERT_NAME' 已存在"
        read -p "是否要重新创建？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "使用现有证书"
            return 1
        fi
        
        # 删除现有证书
        log_info "删除现有证书..."
        security delete-identity -c "$CERT_NAME" 2>/dev/null || true
    fi
    return 0
}

# 创建自签名证书
create_certificate() {
    log_info "创建自签名证书..."
    
    # 创建证书请求配置
    local cert_config="/tmp/cert.conf"
    cat > "$cert_config" << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C=CN
ST=Beijing
L=Beijing
O=wepass
OU=Development
CN=$CERT_NAME

[v3_req]
keyUsage = digitalSignature
extendedKeyUsage = codeSigning
EOF

    # 生成私钥
    local private_key="/tmp/private.key"
    openssl genrsa -out "$private_key" 2048
    
    # 生成证书请求
    local cert_request="/tmp/cert.csr"
    openssl req -new -key "$private_key" -out "$cert_request" -config "$cert_config"
    
    # 生成自签名证书
    local certificate="/tmp/cert.crt"
    openssl x509 -req -days 3650 -in "$cert_request" -signkey "$private_key" -out "$certificate" -extensions v3_req -extfile "$cert_config"
    
    # 转换为 p12 格式（使用密码）
    local p12_file="/tmp/cert.p12"
    openssl pkcs12 -export -out "$p12_file" -inkey "$private_key" -in "$certificate" -name "$CERT_NAME" -passout pass:$KEYCHAIN_PASSWORD
    
    # 导入到钥匙串
    security import "$p12_file" -k "$KEYCHAIN_NAME" -P "$KEYCHAIN_PASSWORD" -A
    
    # 清理临时文件
    rm -f "$cert_config" "$private_key" "$cert_request" "$certificate" "$p12_file"
    
    log_success "证书创建完成"
    return 0
}

# 创建专用钥匙串
create_keychain() {
    log_info "创建专用钥匙串..."
    
    # 删除现有钥匙串（如果存在）
    security delete-keychain "$KEYCHAIN_NAME" 2>/dev/null || true
    
    # 创建新钥匙串
    security create-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_NAME"
    
    # 设置钥匙串属性
    security set-keychain-settings -t 86400 -l "$KEYCHAIN_NAME"
    
    # 解锁钥匙串
    security unlock-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN_NAME"
    
    # 添加到搜索列表
    security list-keychains -d user -s "$KEYCHAIN_NAME" $(security list-keychains -d user | sed s/\"//g)
    
    log_success "钥匙串创建完成: $KEYCHAIN_NAME"
}

# 设置证书信任
setup_certificate_trust() {
    log_info "设置证书信任..."
    
    # 设置证书信任
    security set-key-partition-list -S apple-tool:,apple: -s -k "$KEYCHAIN_PASSWORD" "$KEYCHAIN_NAME"
    
    log_success "证书信任设置完成"
}

# 验证证书
verify_certificate() {
    log_info "验证证书..."
    
    # 检查证书是否可用
    if security find-identity -v -p codesigning | grep -q "$CERT_NAME"; then
        log_success "证书验证成功"
        
        # 显示证书信息
        log_info "可用的代码签名证书:"
        security find-identity -v -p codesigning | grep "$CERT_NAME"
        
        return 0
    else
        log_error "证书验证失败"
        return 1
    fi
}

# 清理临时文件
cleanup() {
    log_info "清理临时文件..."
    rm -f /tmp/wepass-cert.*
}

# 显示使用说明
show_usage() {
    echo ""
    log_info "=========================================="
    log_info "证书创建完成！"
    log_info "=========================================="
    echo ""
    log_info "现在可以使用以下命令重新编译应用："
    echo "  ./build.sh mac-arm64"
    echo ""
    log_info "或者手动签名现有应用："
    echo "  codesign --force --deep --sign \"$CERT_NAME\" build/bin/wepass-1.0.1-arm64.app"
    echo ""
    log_info "证书信息："
    echo "  名称: $CERT_NAME"
    echo "  钥匙串: $KEYCHAIN_NAME"
    echo "  有效期: 10年"
    echo ""
    log_warning "注意："
    echo "  - 这是自签名证书，仅用于本地开发"
    echo "  - 首次运行时仍需要用户确认信任"
    echo "  - 但之后可以正常双击启动"
}

# 主函数
main() {
    log_info "=========================================="
    log_info "wepass 自签名证书创建工具"
    log_info "时间: $(date)"
    log_info "=========================================="
    
    # 检查运行环境
    check_macos
    
    # 检查现有证书
    if ! check_existing_cert; then
        show_usage
        exit 0
    fi
    
    # 创建钥匙串
    create_keychain
    
    # 创建证书
    if create_certificate; then
        log_success "证书创建成功"
    else
        log_error "证书创建失败"
        exit 1
    fi
    
    # 设置证书信任
    setup_certificate_trust
    
    # 验证证书
    if verify_certificate; then
        show_usage
    else
        log_error "证书验证失败，请检查配置"
        exit 1
    fi
    
    # 清理临时文件
    cleanup
    
    log_success "自签名证书设置完成！"
}

# 运行主函数
main "$@"
