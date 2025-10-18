<template>
  <div class="main-view-simple">
    <h1>简化版主界面</h1>
    <p>应用标题: {{ appTitle }}</p>
    <p>测试状态: {{ testStatus }}</p>
    <div v-if="error" class="error">
      错误: {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import lockEventService from '@/services/lockEventService'

const { t } = useI18n()
const apiService = inject('apiService')

const appTitle = ref('密码管理器')
const testStatus = ref('初始化中...')
const error = ref('')

onMounted(async () => {
  try {
    console.log('[简化版MainView] 开始初始化...')

    // 测试1: 获取应用信息
    testStatus.value = '获取应用信息...'
    const appInfo = await apiService.getAppInfo()
    appTitle.value = `${t('main.title')} v${appInfo.version}`
    console.log('[简化版MainView] ✅ 应用信息获取成功')

    // 测试2: 获取分组数据
    testStatus.value = '获取分组数据...'
    const groupsData = await apiService.getGroups()
    console.log('[简化版MainView] ✅ 分组数据获取成功:', groupsData.length)

    // 测试3: 获取锁定配置
    testStatus.value = '获取锁定配置...'
    const lockConfig = await apiService.getLockConfig()
    console.log('[简化版MainView] ✅ 锁定配置获取成功:', lockConfig)

    // 测试4: 初始化锁定服务
    testStatus.value = '初始化锁定服务...'
    lockEventService.setLockConfig(lockConfig)
    console.log('[简化版MainView] ✅ 锁定服务初始化成功')

    // 测试5: 逐个测试MoreMenu的子组件导入
    const componentsToTest = [
      'PasswordGenerator',
      'PasswordRuleSettings',
      'ChangeLogDialog',
      'AboutDialog',
      'HelpDialog',
      'SettingsDialog',
      'TestDialog',
      'ExportVaultDialog',
      'ImportVaultDialog'
    ]

    for (const componentName of componentsToTest) {
      testStatus.value = `测试${componentName}组件导入...`
      try {
        const component = await import(`@/components/${componentName}.vue`)
        console.log(`[简化版MainView] ✅ ${componentName}组件导入成功`)
      } catch (componentError) {
        console.error(`[简化版MainView] ❌ ${componentName}组件导入失败:`, componentError)
        throw new Error(`${componentName}组件导入失败: ${componentError.message}`)
      }
    }

    testStatus.value = '测试SearchBar组件导入...'
    try {
      const SearchBar = await import('@/components/SearchBar.vue')
      console.log('[简化版MainView] ✅ SearchBar组件导入成功')
    } catch (componentError) {
      console.error('[简化版MainView] ❌ SearchBar组件导入失败:', componentError)
      throw new Error(`SearchBar组件导入失败: ${componentError.message}`)
    }

    testStatus.value = '测试GroupsBar组件导入...'
    try {
      const GroupsBar = await import('@/components/GroupsBar.vue')
      console.log('[简化版MainView] ✅ GroupsBar组件导入成功')
    } catch (componentError) {
      console.error('[简化版MainView] ❌ GroupsBar组件导入失败:', componentError)
      throw new Error(`GroupsBar组件导入失败: ${componentError.message}`)
    }

    testStatus.value = '测试MainContentArea组件导入...'
    try {
      const MainContentArea = await import('@/components/MainContentArea.vue')
      console.log('[简化版MainView] ✅ MainContentArea组件导入成功')
    } catch (componentError) {
      console.error('[简化版MainView] ❌ MainContentArea组件导入失败:', componentError)
      throw new Error(`MainContentArea组件导入失败: ${componentError.message}`)
    }

    testStatus.value = '初始化完成'
    console.log('[简化版MainView] 🎉 所有测试通过')

  } catch (err) {
    console.error('[简化版MainView] ❌ 初始化失败:', err)
    error.value = err.message || '未知错误'
    testStatus.value = '初始化失败'
  }
})
</script>

<style scoped>
.main-view-simple {
  padding: 20px;
}

.error {
  color: red;
  background: #ffe6e6;
  padding: 10px;
  border-radius: 4px;
  margin-top: 10px;
}
</style>
