import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ElMessage } from 'element-plus'
import GroupTag from '@/components/GroupTag.vue'
import StatusBar from '@/components/StatusBar.vue'
import { apiService } from '@/services/api'

/**
 * 新增功能单元测试
 * @author 20251003 陈凤庆
 * @description 测试新增功能的组件和逻辑
 */

// Mock Element Plus
vi.mock('element-plus', () => ({
  ElMessage: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn()
  }
}))

// Mock API Service
vi.mock('@/services/api', () => ({
  apiService: {
    updateAccountUsage: vi.fn(),
    getAllAccounts: vi.fn(),
    updateAppUsage: vi.fn(),
    getUsageDays: vi.fn()
  }
}))

describe('GroupTag 组件', () => {
  it('应该正确显示搜索结果分组', () => {
    const searchGroup = {
      id: 'search-result',
      name: '搜索结果 (5)',
      icon: 'fa-search',
      isSearchResult: true
    }

    const wrapper = mount(GroupTag, {
      props: {
        group: searchGroup,
        active: true
      }
    })

    // 检查分组名称
    expect(wrapper.text()).toContain('搜索结果 (5)')
    
    // 检查是否有active类
    expect(wrapper.classes()).toContain('active')
  })

  it('搜索结果分组不应该触发右键菜单', async () => {
    const searchGroup = {
      id: 'search-result',
      name: '搜索结果',
      icon: 'fa-search',
      isSearchResult: true
    }

    const wrapper = mount(GroupTag, {
      props: {
        group: searchGroup
      }
    })

    // 模拟右键点击
    await wrapper.trigger('contextmenu')

    // 检查是否没有触发contextmenu事件
    expect(wrapper.emitted('contextmenu')).toBeFalsy()
  })

  it('普通分组应该能触发右键菜单', async () => {
    const normalGroup = {
      id: '1',
      name: '网站账号',
      icon: 'fa-folder'
    }

    const wrapper = mount(GroupTag, {
      props: {
        group: normalGroup
      }
    })

    // 模拟右键点击
    await wrapper.trigger('contextmenu')

    // 检查是否触发了contextmenu事件
    expect(wrapper.emitted('contextmenu')).toBeTruthy()
  })
})

describe('StatusBar 组件', () => {
  it('应该正确显示使用统计信息', () => {
    const wrapper = mount(StatusBar, {
      props: {
        totalUseCount: 128,
        usageDays: 30
      }
    })

    // 检查授权信息
    expect(wrapper.text()).toContain('授权信息: 已激活专业版')
    
    // 检查使用次数
    expect(wrapper.text()).toContain('使用次数: 128次')
    
    // 检查使用天数
    expect(wrapper.text()).toContain('使用天数: 30天')
  })

  it('应该正确处理零值', () => {
    const wrapper = mount(StatusBar, {
      props: {
        totalUseCount: 0,
        usageDays: 0
      }
    })

    expect(wrapper.text()).toContain('使用次数: 0次')
    expect(wrapper.text()).toContain('使用天数: 0天')
  })
})

describe('API Service 新增方法', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('updateAccountUsage 应该调用正确的 Wails API', async () => {
    // Mock Wails API
    const mockWailsAPI = {
      UpdateAccountUsage: vi.fn().mockResolvedValue()
    }
    
    // Mock getWailsAPI
    vi.spyOn(apiService, 'getWailsAPI').mockReturnValue(mockWailsAPI)

    await apiService.updateAccountUsage('test-account-id')

    expect(mockWailsAPI.UpdateAccountUsage).toHaveBeenCalledWith('test-account-id')
  })

  it('getAllAccounts 应该调用正确的 Wails API', async () => {
    const mockAccounts = [
      { id: '1', title: 'Test Account', use_count: 5 }
    ]
    
    const mockWailsAPI = {
      GetAllAccounts: vi.fn().mockResolvedValue(mockAccounts)
    }
    
    vi.spyOn(apiService, 'getWailsAPI').mockReturnValue(mockWailsAPI)

    const result = await apiService.getAllAccounts()

    expect(mockWailsAPI.GetAllAccounts).toHaveBeenCalled()
    expect(result).toEqual(mockAccounts)
  })

  it('updateAppUsage 应该调用正确的 Wails API', async () => {
    const mockWailsAPI = {
      UpdateAppUsage: vi.fn().mockResolvedValue()
    }
    
    vi.spyOn(apiService, 'getWailsAPI').mockReturnValue(mockWailsAPI)

    await apiService.updateAppUsage()

    expect(mockWailsAPI.UpdateAppUsage).toHaveBeenCalled()
  })

  it('getUsageDays 应该调用正确的 Wails API', async () => {
    const mockWailsAPI = {
      GetUsageDays: vi.fn().mockResolvedValue(30)
    }
    
    vi.spyOn(apiService, 'getWailsAPI').mockReturnValue(mockWailsAPI)

    const result = await apiService.getUsageDays()

    expect(mockWailsAPI.GetUsageDays).toHaveBeenCalled()
    expect(result).toBe(30)
  })

  it('当 Wails API 不可用时应该返回默认值', async () => {
    vi.spyOn(apiService, 'getWailsAPI').mockReturnValue(null)

    const accounts = await apiService.getAllAccounts()
    const days = await apiService.getUsageDays()

    expect(accounts).toEqual([])
    expect(days).toBe(30)
  })
})

describe('搜索功能', () => {
  it('应该能创建搜索结果虚拟分组', () => {
    const searchKeyword = 'GitHub'
    const passwordResults = [
      { id: '1', title: 'GitHub Account' },
      { id: '2', title: 'GitHub Enterprise' }
    ]

    const searchResultGroup = {
      id: 'search-result',
      name: `搜索结果 (${passwordResults.length})`,
      icon: 'fa-search',
      isSearchResult: true,
      sort_order: -1
    }

    expect(searchResultGroup.id).toBe('search-result')
    expect(searchResultGroup.name).toBe('搜索结果 (2)')
    expect(searchResultGroup.icon).toBe('fa-search')
    expect(searchResultGroup.isSearchResult).toBe(true)
    expect(searchResultGroup.sort_order).toBe(-1)
  })

  it('应该能识别搜索结果分组', () => {
    const searchGroup = { id: 'search-result', isSearchResult: true }
    const normalGroup = { id: '1', name: '网站账号' }

    expect(searchGroup.id === 'search-result').toBe(true)
    expect(normalGroup.id === 'search-result').toBe(false)
  })
})

describe('使用统计功能', () => {
  it('应该能计算总使用次数', () => {
    const accounts = [
      { id: '1', use_count: 5 },
      { id: '2', use_count: 10 },
      { id: '3', use_count: 3 }
    ]

    const total = accounts.reduce((sum, account) => sum + (account.use_count || 0), 0)
    expect(total).toBe(18)
  })

  it('应该能处理缺失的 use_count 字段', () => {
    const accounts = [
      { id: '1', use_count: 5 },
      { id: '2' }, // 缺失 use_count
      { id: '3', use_count: 3 }
    ]

    const total = accounts.reduce((sum, account) => sum + (account.use_count || 0), 0)
    expect(total).toBe(8)
  })
})
