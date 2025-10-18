package version

/**
 * 版本信息管理
 * @author 陈凤庆
 * @date 2025-10-01
 * @description 统一管理应用版本号和更新日志
 */

const (
	// AppName 应用名称
	AppName = "密码管理器"

	// Version 当前版本号
	Version = "1.0.3"

	// BuildDate 构建日期
	BuildDate = "2025-10-18"

	// Author 作者
	Author = "陈凤庆"
)

// ChangeLog 更新日志
var ChangeLog = []ChangeLogEntry{
	{
		Version: "1.0.3",
		Date:    "2025-10-18",
		Changes: []string{
			"[新增] 账号右键菜单显示密码功能",
			"[新增] 密码显示对话框，支持密码分段显示",
			"[新增] 密码隐藏/显示切换功能",
			"[新增] 密码复制功能，10秒后自动清理剪贴板",
			"[优化] 移除账号右键菜单中的输入用户名和密码功能",
			"[优化] 调整右键菜单布局，显示密码移至打开地址后面",
		},
	},
	{
		Version: "1.0.2",
		Date:    "2025-10-03",
		Changes: []string{
			"[新增] 更多菜单中选择新密码框功能",
			"[新增] 账号详情中备注信息脱敏显示",
			"[优化] 账号详情中备注编辑框样式和对齐",
			"[优化] 备注框按钮位置和对齐方式",
			"[优化] 取消账号列表分页显示，采用全部显示",
			"[修复] 界面布局和交互体验优化",
		},
	},
	{
		Version: "1.0.1",
		Date:    "2025-10-02",
		Changes: []string{
			"[新增] 统一版本号管理系统",
			"[新增] 跨平台编译脚本和命令",
			"[优化] 版本号在各界面的统一显示",
			"[优化] 构建日期和时间规范",
		},
	},
	{
		Version: "1.0.0",
		Date:    "2025-10-01",
		Changes: []string{
			"[新增] 密码库创建界面优化，自动创建data目录",
			"[新增] 支持简单名称创建密码库，自动添加.db后缀",
			"[新增] 简化模式和完整模式自动切换",
			"[新增] 自动填充功能，支持用户名和密码自动填充",
			"[新增] 跨应用程序窗口焦点管理",
			"[新增] 密码生成器功能",
			"[新增] 分组和标签管理",
			"[新增] 搜索功能，支持模糊搜索",
			"[优化] 界面布局和交互体验",
			"[优化] 数据加密安全性",
			"[修复] 多个已知问题",
		},
	},
}

// ChangeLogEntry 更新日志条目
type ChangeLogEntry struct {
	Version string   // 版本号
	Date    string   // 发布日期
	Changes []string // 更新内容列表
}

// GetVersion 获取版本号
func GetVersion() string {
	return Version
}

// GetAppName 获取应用名称
func GetAppName() string {
	return AppName
}

// GetFullVersion 获取完整版本信息
func GetFullVersion() string {
	return AppName + " v" + Version
}

// GetChangeLog 获取更新日志
func GetChangeLog() []ChangeLogEntry {
	return ChangeLog
}
