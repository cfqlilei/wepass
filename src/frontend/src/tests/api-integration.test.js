/**
 * API集成测试
 * @author 陈凤庆
 * @date 2025-09-28
 * @description 测试前端与Go后端的API集成功能
 */

import { describe, it, expect, beforeEach } from 'vitest'
import { ApiService } from '../services/api.js'

describe('API集成测试', () => {
  let apiService

  beforeEach(() => {
    apiService = new ApiService()
  })

  describe('最近使用密码库状态检查', () => {
    it('应该能够检查最近使用的密码库状态', async () => {
      // 20250928 陈凤庆 测试checkRecentVaultStatus方法
      const result = await apiService.checkRecentVaultStatus()
      
      expect(result).toBeDefined()
      expect(result).toHaveProperty('hasValidVault')
      expect(result).toHaveProperty('vaultPath')
      expect(result).toHaveProperty('isSimplified')
      
      // 验证返回值类型
      expect(typeof result.hasValidVault).toBe('boolean')
      expect(typeof result.vaultPath).toBe('string')
      expect(typeof result.isSimplified).toBe('boolean')
    })

    it('在没有最近使用文件时应该返回正确的状态', async () => {
      // 20250928 陈凤庆 测试无最近使用文件的情况
      const result = await apiService.checkRecentVaultStatus()
      
      // 在开发环境下，应该返回模拟数据
      if (!window.go?.app?.App) {
        expect(result.hasValidVault).toBe(false)
        expect(result.vaultPath).toBe('')
        expect(result.isSimplified).toBe(false)
      }
    })
  })

  describe('密码库操作', () => {
    it('应该能够检查密码库文件是否存在', async () => {
      // 20250928 陈凤庆 测试checkVaultExists方法
      const testPath = '/test/path/vault.db'
      const result = await apiService.checkVaultExists(testPath)
      
      expect(typeof result).toBe('boolean')
    })

    it('应该能够获取最近使用的密码库列表', async () => {
      // 20250928 陈凤庆 测试getRecentVaults方法
      const result = await apiService.getRecentVaults()
      
      expect(Array.isArray(result)).toBe(true)
    })
  })

  describe('错误处理', () => {
    it('在API调用失败时应该正确处理错误', async () => {
      // 20250928 陈凤庆 测试错误处理机制
      try {
        // 模拟一个会失败的API调用
        await apiService.openVault('invalid-path', 'invalid-password')
      } catch (error) {
        expect(error).toBeDefined()
        expect(error.message).toBeDefined()
      }
    })
  })
})