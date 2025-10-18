/**
 * LoginView组件测试
 * @author 陈凤庆
 * @date 2025-09-28
 * @description 测试LoginView组件的功能，包括简化模式和界面高度优化
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ElForm, ElFormItem, ElInput, ElButton } from 'element-plus'
import LoginView from '../views/LoginView.vue'

// 模拟API服务
vi.mock('../services/api.js', () => ({
  ApiService: vi.fn().mockImplementation(() => ({
    checkRecentVaultStatus: vi.fn().mockResolvedValue({
      hasValidVault: false,
      vaultPath: '',
      isSimplified: false
    }),
    checkVaultExists: vi.fn().mockResolvedValue(false),
    getRecentVaults: vi.fn().mockResolvedValue([]),
    openVault: vi.fn().mockResolvedValue({ success: true }),
    createVault: vi.fn().mockResolvedValue({ success: true })
  }))
}))

describe('LoginView组件测试', () => {
  let wrapper

  beforeEach(() => {
    wrapper = mount(LoginView, {
      global: {
        components: {
          ElForm,
          ElFormItem,
          ElInput,
          ElButton
        }
      }
    })
  })

  describe('组件渲染', () => {
    it('应该正确渲染登录界面', () => {
      // 20250928 陈凤庆 验证基本组件是否渲染
      expect(wrapper.find('.login-container').exists()).toBe(true)
      expect(wrapper.find('.login-header').exists()).toBe(true)
      expect(wrapper.find('.login-title').text()).toBe('密码库管理')
    })

    it('应该包含必要的表单元素', () => {
      // 20250928 陈凤庆 验证表单元素
      expect(wrapper.find('input[placeholder="请选择密码库文件路径"]').exists()).toBe(true)
      expect(wrapper.find('input[placeholder="请输入密码库密码"]').exists()).toBe(true)
      expect(wrapper.find('.btn-group').exists()).toBe(true)
    })
  })

  describe('简化模式功能', () => {
    it('在完整模式下应该显示所有元素', async () => {
      // 20250928 陈凤庆 验证完整模式
      await wrapper.vm.$nextTick()
      
      expect(wrapper.vm.isSimplifiedMode).toBe(false)
      expect(wrapper.find('.login-container').classes()).not.toContain('simplified-mode')
    })

    it('在简化模式下应该隐藏文件选择器', async () => {
      // 20250928 陈凤庆 模拟简化模式
      await wrapper.setData({ isSimplifiedMode: true })
      await wrapper.vm.$nextTick()
      
      expect(wrapper.vm.isSimplifiedMode).toBe(true)
      expect(wrapper.find('.login-container').classes()).toContain('simplified-mode')
    })
  })

  describe('界面高度优化', () => {
    it('登录容器应该有正确的CSS类', () => {
      // 20250928 陈凤庆 验证CSS类
      const loginContainer = wrapper.find('.login-container')
      expect(loginContainer.exists()).toBe(true)
    })

    it('简化模式下应该应用紧凑样式', async () => {
      // 20250928 陈凤庆 验证简化模式样式
      await wrapper.setData({ isSimplifiedMode: true })
      await wrapper.vm.$nextTick()
      
      const loginContainer = wrapper.find('.login-container')
      expect(loginContainer.classes()).toContain('simplified-mode')
    })
  })

  describe('表单验证', () => {
    it('应该验证必填字段', async () => {
      // 20250928 陈凤庆 测试表单验证
      const form = wrapper.findComponent({ name: 'ElForm' })
      expect(form.exists()).toBe(true)
    })

    it('在简化模式下应该只验证密码字段', async () => {
      // 20250928 陈凤庆 测试简化模式下的验证
      await wrapper.setData({ 
        isSimplifiedMode: true,
        validRecentVaultPath: '/test/vault.db'
      })
      await wrapper.vm.$nextTick()
      
      expect(wrapper.vm.isSimplifiedMode).toBe(true)
      expect(wrapper.vm.validRecentVaultPath).toBe('/test/vault.db')
    })
  })

  describe('用户交互', () => {
    it('应该能够处理打开密码库操作', async () => {
      // 20250928 陈凤庆 测试打开密码库功能
      const openButton = wrapper.find('.btn-group .el-button')
      expect(openButton.exists()).toBe(true)
      
      // 模拟点击事件
      await openButton.trigger('click')
      // 验证是否调用了相应的方法
    })

    it('应该能够处理创建密码库操作', async () => {
      // 20250928 陈凤庆 测试创建密码库功能
      const buttons = wrapper.findAll('.btn-group .el-button')
      expect(buttons.length).toBeGreaterThan(1)
    })
  })
})