import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

/**
 * 密码库状态管理
 * @author 陈凤庆
 * @description 管理密码库的打开状态和相关信息
 */
export const useVaultStore = defineStore('vault', () => {
  // 状态
  const isVaultOpened = ref(false)
  const currentVaultPath = ref('')
  const recentVaults = ref([])
  
  // 计算属性
  const hasRecentVaults = computed(() => recentVaults.value.length > 0)
  
  // 动作
  /**
   * 设置密码库打开状态
   * @param {boolean} opened 是否已打开
   * @param {string} path 密码库路径
   */
  const setVaultOpened = (opened, path = '') => {
    isVaultOpened.value = opened
    currentVaultPath.value = path
  }
  
  /**
   * 设置最近使用的密码库列表
   * @param {Array<string>} vaults 密码库路径列表
   */
  const setRecentVaults = (vaults) => {
    recentVaults.value = vaults
  }
  
  /**
   * 关闭密码库
   */
  const closeVault = () => {
    isVaultOpened.value = false
    currentVaultPath.value = ''
  }
  
  return {
    // 状态
    isVaultOpened,
    currentVaultPath,
    recentVaults,
    
    // 计算属性
    hasRecentVaults,
    
    // 动作
    setVaultOpened,
    setRecentVaults,
    closeVault
  }
})