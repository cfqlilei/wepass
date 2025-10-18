/**
 * 测试环境设置
 * @author 陈凤庆
 * @description 配置前端测试的全局设置和模拟
 */

import { config } from "@vue/test-utils";

// 配置全局组件存根
config.global.stubs = {
  "el-form": true,
  "el-form-item": true,
  "el-input": true,
  "el-button": true,
  "el-icon": true,
  "el-divider": true,
};

// 模拟 Wails API
global.window = global.window || {};
global.window.go = {
  app: {
    App: {
      CheckVaultExists: vi.fn().mockResolvedValue(false),
      CreateVault: vi.fn().mockResolvedValue(undefined),
      OpenVault: vi.fn().mockResolvedValue(undefined),
      GetRecentVaults: vi.fn().mockResolvedValue([]),
      GetGroups: vi
        .fn()
        .mockResolvedValue([
          {
            id: 1,
            name: "默认",
            parent_id: 0,
            icon: "fa-folder-open",
            sort_order: 0,
          },
        ]),
      CreateGroup: vi.fn().mockResolvedValue({ id: 2, name: "新分组" }),
      GetPasswordsByGroup: vi.fn().mockResolvedValue([]),
      CreatePasswordItem: vi.fn().mockResolvedValue({ id: 1 }),
      SearchAccounts: vi.fn().mockResolvedValue([]),
      GetAppInfo: vi.fn().mockResolvedValue({
        name: "wepass",
        version: "1.0.0",
        author: "陈凤庆",
      }),
    },
  },
};

// 模拟剪贴板 API
Object.assign(navigator, {
  clipboard: {
    writeText: vi.fn().mockResolvedValue(undefined),
    readText: vi.fn().mockResolvedValue(""),
  },
});

// 模拟 Element Plus 消息组件
global.ElMessage = {
  success: vi.fn(),
  error: vi.fn(),
  warning: vi.fn(),
  info: vi.fn(),
};

global.ElMessageBox = {
  prompt: vi.fn().mockResolvedValue({ value: "test" }),
  confirm: vi.fn().mockResolvedValue(undefined),
};
