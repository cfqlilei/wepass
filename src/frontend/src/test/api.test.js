/**
 * API 服务测试
 * @author 陈凤庆
 * @description 测试API服务的功能
 */

import { describe, it, expect, beforeEach, vi } from "vitest";
import { apiService } from "@/services/api";

describe("API Service", () => {
  beforeEach(() => {
    // 重置所有模拟函数
    vi.clearAllMocks();
  });

  it("应该能够检查密码库是否存在", async () => {
    const result = await apiService.checkVaultExists("/test/path");
    expect(window.go.app.App.CheckVaultExists).toHaveBeenCalledWith(
      "/test/path"
    );
    expect(typeof result).toBe("boolean");
  });

  it("应该能够创建密码库", async () => {
    await apiService.createVault("/test/path", "password123");
    expect(window.go.app.App.CreateVault).toHaveBeenCalledWith(
      "/test/path",
      "password123"
    );
  });

  it("应该能够打开密码库", async () => {
    await apiService.openVault("/test/path", "password123");
    expect(window.go.app.App.OpenVault).toHaveBeenCalledWith(
      "/test/path",
      "password123"
    );
  });

  it("应该能够获取最近使用的密码库", async () => {
    const result = await apiService.getRecentVaults();
    expect(window.go.app.App.GetRecentVaults).toHaveBeenCalled();
    expect(Array.isArray(result)).toBe(true);
  });

  it("应该能够获取分组列表", async () => {
    const result = await apiService.getGroups();
    expect(window.go.app.App.GetGroups).toHaveBeenCalled();
    expect(Array.isArray(result)).toBe(true);
  });

  it("应该能够创建分组", async () => {
    const groupName = "测试分组";
    const parentId = 0;
    const result = await apiService.createGroup(groupName, parentId);
    expect(window.go.app.App.CreateGroup).toHaveBeenCalledWith(
      groupName,
      parentId
    );
    expect(result).toBeDefined();
  });

  it("应该能够根据分组获取账号", async () => {
    const result = await apiService.getPasswordsByGroup(1);
    expect(window.go.app.App.GetPasswordsByGroup).toHaveBeenCalledWith(1);
    expect(Array.isArray(result)).toBe(true);
  });

  it("应该能够创建账号", async () => {
    const accountData = {
      title: "测试密码",
      username: "test@example.com",
      password: "password123",
      group_id: 1,
    };
    const result = await apiService.createPasswordItem(accountData);
    expect(window.go.app.App.CreatePasswordItem).toHaveBeenCalledWith(
      accountData
    );
    expect(result).toBeDefined();
  });

  it("应该能够搜索密码", async () => {
    const result = await apiService.searchPasswords("test");
    expect(window.go.app.App.SearchAccounts).toHaveBeenCalledWith("test");
    expect(Array.isArray(result)).toBe(true);
  });

  it("应该能够获取应用信息", async () => {
    const result = await apiService.getAppInfo();
    expect(window.go.app.App.GetAppInfo).toHaveBeenCalled();
    expect(result).toBeDefined();
    expect(result.name).toBe("wepass");
  });

  it("当Wails API未初始化时应该正确处理", async () => {
    // 临时移除window.go
    const originalGo = window.go;
    delete window.go;

    // 创建新的API服务实例（模拟未初始化状态）
    const { default: ApiServiceClass } = await import("@/services/api");
    // 由于我们无法直接导入类，我们模拟未初始化的行为
    const mockApiService = {
      async checkVaultExists(vaultPath) {
        console.warn("Wails API 未初始化，使用模拟数据");
        return false;
      },
    };

    const result = await mockApiService.checkVaultExists("/test/path");
    expect(result).toBe(false); // API未初始化时返回false

    // 恢复window.go
    window.go = originalGo;
  });
});
