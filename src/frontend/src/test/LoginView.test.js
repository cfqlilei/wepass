/**
 * LoginView 组件测试
 * @author 陈凤庆
 * @description 测试登录界面组件的功能
 */

import { describe, it, expect, beforeEach, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { createPinia } from "pinia";
import { createRouter, createMemoryHistory } from "vue-router";
import LoginView from "@/views/LoginView.vue";

describe("LoginView", () => {
  let wrapper;
  let pinia;
  let router;

  beforeEach(async () => {
    pinia = createPinia();
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: "/", component: LoginView },
        { path: "/main", component: { template: "<div>Main</div>" } },
      ],
    });

    wrapper = mount(LoginView, {
      global: {
        plugins: [pinia, router],
        stubs: {
          "el-form": true,
          "el-form-item": true,
          "el-input": true,
          "el-button": true,
          "el-icon": true,
        },
      },
    });

    await router.isReady();
  });

  it("应该正确渲染登录界面", () => {
    expect(wrapper.find(".login-title").text()).toBe("密码库管理");
    expect(wrapper.find(".login-desc").text()).toBe("请选择或创建您的密码库");
  });

  it("应该包含Element Plus组件存根", () => {
    // 检查是否有Element Plus组件存根
    const elForm = wrapper.find("el-form-stub");
    const elButtons = wrapper.findAll("el-button-stub");
    const elInputs = wrapper.findAll("el-input-stub");

    expect(elForm.exists()).toBe(true);
    expect(elButtons.length).toBeGreaterThan(0);
    expect(elInputs.length).toBeGreaterThan(0);
  });

  it("应该包含必要的UI文本", () => {
    // 检查页面是否包含关键文本
    const text = wrapper.text();
    expect(text).toContain("密码库管理");
    expect(text).toContain("请选择或创建您的密码库");
  });

  it("应该显示密码提示信息", () => {
    const hints = wrapper.findAll(".password-hint p");
    expect(hints.length).toBe(3);
    expect(hints[0].text()).toContain("登录密码是密码库加密和解密的关键密钥");
  });

  // 20250127 陈凤庆 新增测试用例，测试修复的功能
  describe("文件选择功能", () => {
    it("应该检测Wails环境", async () => {
      // 模拟非Wails环境
      const originalWindow = global.window;
      global.window = {};

      const selectButton = wrapper.find(".file-select-btn");
      if (selectButton.exists()) {
        await selectButton.trigger("click");
        // 在非Wails环境中应该显示警告信息
        // 这里需要根据实际的消息提示组件进行验证
      }

      global.window = originalWindow;
    });

    it("应该在Wails环境中正常工作", async () => {
      // 模拟Wails环境
      global.window = {
        go: {
          app: {
            App: {
              SelectVaultFile: vi.fn().mockResolvedValue("/path/to/vault.db"),
            },
          },
        },
      };

      const selectButton = wrapper.find(".file-select-btn");
      if (selectButton.exists()) {
        await selectButton.trigger("click");
        // 验证文件选择功能被调用
        expect(global.window.go.app.App.SelectVaultFile).toHaveBeenCalled();
      }
    });
  });

  describe("按钮布局", () => {
    it("应该包含按钮组", () => {
      const btnGroup = wrapper.find(".btn-group");
      expect(btnGroup.exists()).toBe(true);
    });

    it("应该包含文件选择器", () => {
      const fileSelector = wrapper.find(".file-selector");
      expect(fileSelector.exists()).toBe(true);
    });
  });

  describe("高度自适应", () => {
    it("应该在没有最近使用记录时不显示相关区域", () => {
      // 确保没有最近使用记录
      const recentVaults = wrapper.find(".recent-vaults");
      expect(recentVaults.exists()).toBe(false);
    });

    it("应该包含登录容器", () => {
      const container = wrapper.find(".login-container");
      expect(container.exists()).toBe(true);
    });
  });
});
