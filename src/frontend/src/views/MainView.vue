<template>
  <div class="main-container">
    <!-- 顶部标题栏 -->
    <TitleBar
      :app-title="appTitle"
      @toggle-search="toggleSearch"
      @refresh-data="handleRefreshData"
    />

    <!-- 搜索栏 -->
    <SearchBar
      v-model="searchKeyword"
      :show="showSearchBar"
      @search="handleSearch"
      @clear="clearSearch"
    />

    <!-- 顶部分组标签 -->
    <GroupsBar
      :groups="groups"
      :current-group-id="currentGroupId"
      @select-group="selectGroup"
      @create-group="createGroup"
      @show-group-menu="showGroupMenu"
      @drag-reorder="handleGroupDragReorder"
    />

    <!-- 中部区域 -->
    <MainContentArea
      ref="mainContentAreaRef"
      :tabs="tabs"
      :current-tab-id="currentTabId"
      :current-tab="currentTab"
      :filtered-accounts="filteredAccounts"
      :show-tabs="!isSearching"
      @select-tab="selectTab"
      @create-tab="createTab"
      @show-tab-menu="showTabMenu"
      @drag-reorder="handleTabDragReorder"
      @create-account="createAccount"
      @show-account-detail="showAccountDetailDialog"
      @show-account-context-menu="showAccountContextMenu"
      @input-username-and-password="inputUsernameAndPassword"
      @input-username="inputUsername"
      @input-password="inputPassword"
    />

    <!-- 底部状态栏 -->
    <StatusBar
      :total-use-count="totalUseCount"
      :usage-days="usageDays"
    />

    <!-- 账号详情对话框 -->
    <AccountDetail
      v-model="showAccountDetail"
      :account-Data="selectedAccount"
      :is-new="isNewAccount"
      @save="handleSaveAccount"
      @delete="handleDeleteAccount"
    />

    <!-- 账号项右键菜单 -->
    <AccountContextMenu
      :show="showContextMenu"
      :position="contextMenuPosition"
      :account="contextMenuAccount"
      @close="hideContextMenu"
      @open-url="openAccountUrl"
      @view="showAccountDetailDialog"
      @edit="editAccount"
      @change-group="handleChangeGroup"
      @duplicate="duplicateAccount"
      @copy-password="() => copyAccountField(contextMenuAccount, 'password')"
      @show-password="showPasswordDialog"
      @copy-username="() => copyAccountField(contextMenuAccount, 'username')"
      @copy-username-and-password="() => copyAccountUsernameAndPassword(contextMenuAccount)"
      @copy-url="() => copyAccountField(contextMenuAccount, 'url')"
      @copy-title="() => copyAccountField(contextMenuAccount, 'title')"
      @copy-notes="() => copyAccountField(contextMenuAccount, 'notes')"
      @delete="deleteAccountWithConfirm"
    />

    <!-- 分组右键菜单 -->
    <GroupContextMenu
      :visible="showGroupContextMenu"
      :position="groupContextMenuPosition"
      :group="contextMenuGroup"
      @close="hideGroupContextMenu"
      @rename="handleRenameGroup"
      @delete="handleDeleteGroup"
      @move-left="handleMoveGroupLeft"
      @move-right="handleMoveGroupRight"
      @create-group="handleCreateGroupFromMenu"
    />

    <!-- 标签右键菜单 -->
    <TabContextMenu
      :visible="showTabContextMenu"
      :position="tabContextMenuPosition"
      :tab="contextMenuTab"
      @close="hideTabContextMenu"
      @rename="handleRenameTab"
      @delete="handleDeleteTab"
      @move-up="handleMoveTabUp"
      @move-down="handleMoveTabDown"
      @create-after="handleCreateTabAfter"
    />

    <!-- 密码显示对话框 -->
    <PasswordDisplayDialog
      v-model="showPasswordDisplayDialog"
      :account="passwordDisplayAccount"
    />

    <!-- 更改分组对话框 -->
    <ChangeGroupDialog
      :visible="showChangeGroupDialog"
      :account="changeGroupAccount"
      :groups="groups"
      @update:visible="showChangeGroupDialog = $event"
      @confirm="handleChangeGroupConfirm"
    />
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useVaultStore } from '@/stores/vault'
import { apiService } from '@/services/api'
import lockEventService from '@/services/lockEventService'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import { copyAccountPassword } from '@/utils/accountUtils'
import AccountDetail from '@/components/AccountDetail.vue'
import TitleBar from '@/components/TitleBar.vue'
import SearchBar from '@/components/SearchBar.vue'
import GroupsBar from '@/components/GroupsBar.vue'
import MainContentArea from '@/components/MainContentArea.vue'
import StatusBar from '@/components/StatusBar.vue'
import AccountContextMenu from '@/components/AccountContextMenu.vue'
import GroupContextMenu from '@/components/GroupContextMenu.vue'
import TabContextMenu from '@/components/TabContextMenu.vue'
import PasswordDisplayDialog from '@/components/PasswordDisplayDialog.vue'
import ChangeGroupDialog from '@/components/ChangeGroupDialog.vue'

/**
 * 主界面视图组件
 * @author 陈凤庆
 * @modify 20251001 陈凤庆 重构为多组件组合方式，提高代码可维护性
 * @description 密码管理工具的主界面，包含分组、页签和密码列表
 */
export default {
  name: 'MainView',
  components: {
    AccountDetail,
    TitleBar,
    SearchBar,
    GroupsBar,
    MainContentArea,
    StatusBar,
    AccountContextMenu,
    GroupContextMenu,
    TabContextMenu,
    PasswordDisplayDialog,
    ChangeGroupDialog
  },
  setup() {
    const { t } = useI18n()
    const router = useRouter()
    const appStore = useAppStore()
    const vaultStore = useVaultStore()

    // 组件引用
    const mainContentAreaRef = ref(null)

    // 响应式数据
    const appTitle = computed(() => t('main.title'))
    const showSearchBar = ref(false)
    const searchKeyword = ref('')
    const totalUseCount = ref(128)
    const usageDays = ref(30)
    const showAccountDetail = ref(false)
    const selectedAccount = ref({})
    const isNewAccount = ref(false)
    const searchResults = ref({ groups: [], accounts: [] })

    // 搜索请求控制器 - 用于取消过期的搜索请求
    let searchAbortController = null

    // 右键菜单相关状态
    const showContextMenu = ref(false)
    const contextMenuPosition = reactive({ x: 0, y: 0 })
    const contextMenuAccount = ref(null)

    // 分组右键菜单相关状态
    const showGroupContextMenu = ref(false)
    const groupContextMenuPosition = reactive({ x: 0, y: 0 })
    const contextMenuGroup = ref(null)

    // 标签右键菜单相关状态
    const showTabContextMenu = ref(false)
    const tabContextMenuPosition = reactive({ x: 0, y: 0 })
    const contextMenuTab = ref(null)

    // 密码显示对话框相关状态
    const showPasswordDisplayDialog = ref(false)
    const passwordDisplayAccount = ref(null)

    // 更改分组对话框相关状态
    const showChangeGroupDialog = ref(false)
    const changeGroupAccount = ref(null)

    // 计算属性
    const groups = computed(() => appStore.groups)
    const currentGroupId = computed(() => appStore.currentGroupId)
    const tabs = computed(() => appStore.tabs)
    const currentTabId = computed(() => appStore.currentTabId)
    const currentTab = computed(() => appStore.currentTab)
    const filteredAccounts = computed(() => appStore.filteredPasswords)
    const isSearching = computed(() => appStore.isSearching)
    
    /**
     * 组件挂载时执行
     * @modify 陈凤庆 添加窗口焦点管理，存储目标应用程序
     * @modify 20250101 陈凤庆 优化初始化逻辑，默认选中第一个分组和第一个标签，添加详细日志
     */
    onMounted(async () => {
      try {
        console.log('='.repeat(80))
        console.log('[onMounted] 🚀 开始初始化主界面...')
        console.log('[onMounted] 当前时间:', new Date().toISOString())
        console.log('[onMounted] 当前路由:', window.location.href)
        console.log('='.repeat(80))

        // 获取应用信息并设置标题
        console.log('[onMounted] 📱 获取应用信息...')
        const appInfo = await apiService.getAppInfo()
        console.log('[onMounted] ✅ 应用信息获取成功:', appInfo)
        appTitle.value = `${t('main.title')} v${appInfo.version}`
        console.log('[onMounted] 📝 标题已设置:', appTitle.value)

        // 20251019 陈凤庆 加载页签记忆
        console.log('[onMounted] 加载页签记忆...')
        appStore.loadGroupTabMemory()

        // 20250127 陈凤庆 存储当前聚焦的应用程序作为目标应用程序
        // 这样当用户从其他应用程序切换到密码管理器时，我们就知道目标应用程序是什么
        console.log('[onMounted] 存储目标应用程序...')
        await window.go.app.App.StorePreviousFocusedApp()

        // 加载分组数据
        console.log('[onMounted] 加载分组数据...')
        const groupsData = await apiService.getGroups()
        console.log('[onMounted] 获取到分组数据:', groupsData)
        appStore.setGroups(groupsData)

        // 初始化使用统计
        console.log('[onMounted] 初始化使用统计...')
        await apiService.updateAppUsage() // 更新应用使用统计
        await refreshTotalUseCount() // 刷新总使用次数
        await refreshUsageDays() // 刷新使用天数

        // 20250101 陈凤庆 如果有分组，默认选中第一个分组
        // 20251019 陈凤庆 使用selectGroup函数，支持页签记忆功能
        if (groupsData && groupsData.length > 0) {
          const firstGroupId = groupsData[0].id
          console.log(`[onMounted] 默认选中第一个分组: ${groupsData[0].name} (ID: ${firstGroupId})`)

          // 使用selectGroup函数，自动处理页签记忆和选择
          await selectGroup(firstGroupId)
        } else {
          console.warn('[onMounted] 没有分组数据')
          ElMessage.warning(t('warning.noGroupData'))
        }

        // 20250127 陈凤庆 添加窗口焦点事件监听，当应用程序失去焦点时存储目标应用程序
        window.addEventListener('blur', async () => {
          try {
            // 延迟一点时间，确保新的应用程序已经获得焦点
            setTimeout(async () => {
              await window.go.app.App.StorePreviousFocusedApp()
            }, 100)
          } catch (error) {
            console.error('存储目标应用程序失败:', error)
          }
        })

        // 20251004 陈凤庆 添加用户活动跟踪
        console.log('[onMounted] 🔒 初始化用户活动跟踪...')
        setupUserActivityTracking()
        console.log('[onMounted] ✅ 用户活动跟踪已启用')

        // 20251004 陈凤庆 初始化锁定事件服务
        console.log('[onMounted] 🔒 初始化锁定事件服务...')
        lockEventService.init(router, vaultStore)

        // 20251004 陈凤庆 获取并应用锁定配置
        try {
          console.log('[onMounted] 🔒 获取锁定配置...')

          // 优先从localStorage获取（登录时存储的配置）
          let lockConfig = null
          const storedConfig = localStorage.getItem('lock_config')
          if (storedConfig) {
            try {
              lockConfig = JSON.parse(storedConfig)
              console.log('[onMounted] 从localStorage获取锁定配置:', lockConfig)
            } catch (jsonError) {
              console.error('[onMounted] localStorage中的锁定配置格式错误:', jsonError)
              // 清理错误的配置数据
              localStorage.removeItem('lock_config')
              // 从后端重新获取
              lockConfig = await apiService.getLockConfig()
              console.log('[onMounted] 从后端重新获取锁定配置:', lockConfig)
            }
          } else {
            // 如果localStorage没有，从后端获取
            lockConfig = await apiService.getLockConfig()
            console.log('[onMounted] 从后端获取锁定配置:', lockConfig)
          }

          // 清理localStorage中的配置（避免过期配置）
          localStorage.removeItem('lock_config')

          // 20251004 陈凤庆 将锁定配置设置到锁定事件服务中
          lockEventService.setLockConfig(lockConfig)

          console.log('[onMounted] ✅ 锁定配置加载完成，启动监听服务')
        } catch (error) {
          console.error('[onMounted] ❌ 获取锁定配置失败:', error)
          // 继续启动服务，使用默认配置
          lockEventService.setLockConfig(null)
        }

        lockEventService.startMonitoring()
        console.log('[onMounted] ✅ 锁定事件服务已启动')

        console.log('='.repeat(80))
        console.log('[onMounted] ✅ 主界面初始化完成')
        console.log('[onMounted] 当前Store状态:', {
          groups: appStore.groups.length,
          currentGroupId: appStore.currentGroupId,
          tabs: appStore.tabs.length,
          currentTabId: appStore.currentTabId,
          passwords: appStore.passwords.length
        })
        console.log('='.repeat(80))

      } catch (error) {
        console.error('='.repeat(80))
        console.error('[onMounted] ❌ 初始化主界面失败:', error)
        console.error('[onMounted] 错误类型:', error.name)
        console.error('[onMounted] 错误消息:', error.message)
        console.error('[onMounted] 错误堆栈:', error.stack)
        console.error('='.repeat(80))
        ElMessage.error(`${t('error.dataLoadFailed')}: ${error.message || t('error.unknownError')}`)
      }
    })

    /**
     * 组件卸载时清理资源
     * @author 陈凤庆
     * @date 20251004
     */
    onUnmounted(() => {
      console.log('[onUnmounted] 清理资源...')

      // 停止锁定事件监听
      lockEventService.stopMonitoring()
      console.log('[onUnmounted] 锁定事件服务已停止')

      // 这里可以添加其他清理逻辑
      console.log('[onUnmounted] 资源清理完成')
    })

    /**
     * 根据分组加载页签
     * @param {number} groupId 分组ID
     * @param {boolean} autoSelectFirst 是否自动选中第一个标签，默认false
     * @modify 20250101 陈凤庆 加载页签后默认选中第一个，添加详细错误日志
     * @modify 20250101 陈凤庆 分组无标签时不报错，界面正常显示
     * @modify 20251001 陈凤庆 添加autoSelectFirst参数，控制是否自动选中第一个标签
     */
    const loadTabsByGroup = async (groupId, autoSelectFirst = false) => {
      try {
        console.log(`[loadTabsByGroup] 开始加载分组${groupId}的标签...`)
        const tabsData = await apiService.getTabsByGroup(groupId)
        console.log(`[loadTabsByGroup] 获取到标签数据:`, tabsData)

        appStore.setTabs(tabsData)

        // 20251001 陈凤庆 根据autoSelectFirst参数决定是否自动选中第一个标签
        if (autoSelectFirst && tabsData && tabsData.length > 0) {
          appStore.setCurrentTab(tabsData[0].id)
          console.log(`[loadTabsByGroup] 已加载分组${groupId}的${tabsData.length}个标签，默认选中第一个标签: ${tabsData[0].name}`)
        } else {
          // 20251001 陈凤庆 不自动选中标签，清空当前标签选中状态
          appStore.setCurrentTab(null)
          console.log(`[loadTabsByGroup] 已加载分组${groupId}的${tabsData?.length || 0}个标签，不自动选中`)
        }
      } catch (error) {
        console.error('[loadTabsByGroup] 加载页签失败:', error)
        console.error('[loadTabsByGroup] 错误详情:', error.message, error.stack)
        ElMessage.error(`${t('error.dataLoadFailed')}: ${error.message || t('error.unknownError')}`)
      }
    }

    /**
     * 根据分组加载账号
     * @param {number} groupId 分组ID
     * @modify 20251001 陈凤庆 添加详细日志记录，排查加载失败问题
     */
    const loadPasswordsByGroup = async (groupId) => {
      try {
        console.log('='.repeat(60))
        console.log('[loadPasswordsByGroup] 🔐 开始加载账号...')
        console.log('[loadPasswordsByGroup] 分组ID:', groupId)
        console.log('[loadPasswordsByGroup] 分组ID类型:', typeof groupId)
        console.log('[loadPasswordsByGroup] 当前时间:', new Date().toISOString())

        // 20251001 陈凤庆 检查API服务是否可用
        if (!apiService) {
          console.error('[loadPasswordsByGroup] ❌ API服务不可用')
          ElMessage.error(t('error.apiServiceUnavailable'))
          return
        }
        console.log('[loadPasswordsByGroup] ✅ API服务检查通过')

        // 20251001 陈凤庆 检查Wails API是否可用
        const wailsAPI = apiService.getWailsAPI()
        if (!wailsAPI) {
          console.error('[loadPasswordsByGroup] ❌ Wails API不可用')
          ElMessage.error(t('error.backendServiceUnavailable'))
          return
        }
        console.log('[loadPasswordsByGroup] ✅ Wails API检查通过')

        console.log('[loadPasswordsByGroup] 📡 开始调用后端API...')
        const passwords = await apiService.getPasswordsByGroup(groupId)
        console.log('[loadPasswordsByGroup] ✅ 后端API调用成功')
        console.log('[loadPasswordsByGroup] 返回的账号数量:', passwords ? passwords.length : 0)
        console.log('[loadPasswordsByGroup] 返回的数据类型:', typeof passwords)
        console.log('[loadPasswordsByGroup] 返回的数据:', passwords)

        // 20251001 陈凤庆 验证返回数据格式
        if (!Array.isArray(passwords)) {
          console.error('[loadPasswordsByGroup] ❌ 返回数据不是数组格式:', passwords)
          ElMessage.error(t('error.accountDataFormatError'))
          return
        }

        console.log('[loadPasswordsByGroup] 📝 更新应用状态...')
        appStore.setPasswords(passwords)
        console.log('[loadPasswordsByGroup] ✅ 应用状态更新完成')
        console.log('='.repeat(60))

      } catch (error) {
        console.error('='.repeat(60))
        console.error('[loadPasswordsByGroup] ❌ 加载账号失败')
        console.error('[loadPasswordsByGroup] 分组ID:', groupId)
        console.error('[loadPasswordsByGroup] 错误类型:', error.name)
        console.error('[loadPasswordsByGroup] 错误消息:', error.message)
        console.error('[loadPasswordsByGroup] 错误堆栈:', error.stack)
        console.error('[loadPasswordsByGroup] 完整错误对象:', error)
        console.error('='.repeat(60))

        // 20251001 陈凤庆 根据错误类型提供更具体的错误信息
        let errorMessage = '加载账号失败'
        if (error.message) {
          if (error.message.includes('数据库未打开')) {
            errorMessage = '数据库未打开，请重新打开密码库'
          } else if (error.message.includes('加密管理器未初始化')) {
            errorMessage = '加密服务未初始化，请重新打开密码库'
          } else if (error.message.includes('查询账号失败')) {
            errorMessage = '数据库查询失败，请检查密码库文件'
          } else if (error.message.includes('解密账号失败')) {
            errorMessage = '账号解密失败，可能密码库已损坏'
          } else {
            errorMessage = `加载账号失败: ${error.message}`
          }
        }

        ElMessage.error(errorMessage)
      }
    }

    /**
     * 选择分组
     * @param {number} groupId 分组ID
     * @modify 20250101 陈凤庆 切换分组时，先加载标签，再加载账号
     * @modify 20251001 陈凤庆 切换分组时不自动选中标签，只显示标签列表
     * @modify 20251003 陈凤庆 支持搜索结果分组
     * @modify 20251019 陈凤庆 修复问题 003-1：增加详细调试日志，跟踪分组切换过程
     * @modify 20251019 陈凤庆 新增分组页签记忆功能：自动选择记忆的页签或第一个页签
     */
    const selectGroup = async (groupId) => {
      console.log(`[selectGroup] ========== 开始切换分组 ==========`)
      console.log(`[selectGroup] 目标分组ID: ${groupId} (类型: ${typeof groupId})`)
      console.log(`[selectGroup] 切换前分组ID: ${currentGroupId.value}`)
      console.log(`[selectGroup] 切换前标签ID: ${currentTabId.value}`)

      // 设置当前分组
      appStore.setCurrentGroup(groupId)
      console.log(`[selectGroup] ✅ 已设置当前分组: ${appStore.currentGroupId}`)

      // 如果是搜索结果分组，直接显示搜索结果
      if (groupId === 'search-result') {
        console.log('[selectGroup] 切换到搜索结果分组')
        // 清空标签列表
        appStore.setTabs([])
        appStore.setCurrentTab(null)
        console.log('[selectGroup] 已清空标签列表和当前标签')
        // 显示搜索结果
        if (searchResults.value.passwords) {
          appStore.setPasswords(searchResults.value.passwords)
          console.log(`[selectGroup] 已显示${searchResults.value.passwords.length}个搜索结果`)
        }
        console.log(`[selectGroup] ========== 搜索结果分组切换完成 ==========`)
        return
      }

      console.log(`[selectGroup] 开始加载分组${groupId}的标签...`)
      // 20251019 陈凤庆 加载该分组的页签，并根据记忆选择页签
      await loadTabsByGroup(groupId, false)
      console.log(`[selectGroup] 标签加载完成，当前标签数量: ${appStore.tabs.length}`)

      // 20251019 陈凤庆 新增：分组页签记忆功能
      if (appStore.tabs.length > 0) {
        // 尝试获取记忆的页签
        const rememberedTabId = appStore.getRememberedTab(groupId)
        console.log(`[selectGroup] 分组${groupId}的记忆页签: ${rememberedTabId}`)

        // 检查记忆的页签是否存在于当前标签列表中
        const rememberedTab = rememberedTabId ? appStore.tabs.find(tab => tab.id === rememberedTabId) : null

        if (rememberedTab) {
          console.log(`[selectGroup] 找到记忆页签: ${rememberedTab.name} (ID: ${rememberedTab.id})`)
          await selectTab(rememberedTab.id)
        } else {
          // 如果没有记忆页签或记忆页签不存在，选择第一个页签
          const firstTab = appStore.tabs[0]
          console.log(`[selectGroup] 选择第一个页签: ${firstTab.name} (ID: ${firstTab.id})`)
          await selectTab(firstTab.id)
        }
      } else {
        console.log(`[selectGroup] 分组${groupId}没有页签，加载分组所有账号`)
        // 如果没有页签，加载该分组的所有账号
        await loadPasswordsByGroup(groupId)
      }

      console.log(`[selectGroup] ========== 分组${groupId}切换完成 ==========`)
      console.log(`[selectGroup] 最终状态检查:`)
      console.log(`  - 当前分组ID: ${currentGroupId.value}`)
      console.log(`  - 当前标签ID: ${currentTabId.value}`)
      console.log(`  - 标签数量: ${appStore.tabs.length}`)
      console.log(`  - 账号数量: ${appStore.passwords.length}`)
    }
    
    /**
     * 选择页签
     * @param {number} tabId 页签ID
     * @modify 20250101 陈凤庆 添加日志记录，标签切换时会自动触发账号筛选
     * @modify 20251002 陈凤庆 点击标签时调用后端API查询对应的账号列表
     * @modify 20251019 陈凤庆 修复问题 003-1：增加详细调试日志，跟踪标签切换过程
     */
    const selectTab = async (tabId) => {
      console.log(`[selectTab] ========== 开始切换标签 ==========`)
      console.log(`[selectTab] 目标标签ID: ${tabId} (类型: ${typeof tabId})`)
      console.log(`[selectTab] 切换前标签ID: ${currentTabId.value}`)
      console.log(`[selectTab] 当前分组ID: ${currentGroupId.value}`)

      appStore.setCurrentTab(tabId)
      console.log(`[selectTab] ✅ 已设置当前标签: ${appStore.currentTabId}`)

      // 20251002 陈凤庆 点击标签时，调用后端API查询该标签下的账号列表
      const tab = appStore.tabs.find(t => t.id === tabId)
      if (tab) {
        console.log(`[selectTab] 找到标签对象: ${tab.name} (ID: ${tab.id})`)

        try {
          // 调用后端API获取该标签下的账号列表
          console.log(`[selectTab] 开始加载标签${tabId}的账号列表...`)
          const accounts = await apiService.getAccountsByTab(tabId)
          console.log(`[selectTab] ✅ 获取到${accounts.length}个账号`)

          // 更新Store中的账号列表
          appStore.setPasswords(accounts)
          console.log(`[selectTab] ✅ 已更新Store中的账号列表`)

        } catch (error) {
          console.error('[selectTab] ❌ 加载标签账号列表失败:', error)
          ElMessage.error(`加载账号列表失败: ${error.message || '未知错误'}`)
        }
      } else {
        console.error(`[selectTab] ❌ 未找到标签ID为${tabId}的标签对象`)
        console.error(`[selectTab] 可用标签列表:`, appStore.tabs.map(t => ({id: t.id, name: t.name})))
      }

      console.log(`[selectTab] ========== 标签${tabId}切换完成 ==========`)
      console.log(`[selectTab] 最终状态检查:`)
      console.log(`  - 当前分组ID: ${currentGroupId.value}`)
      console.log(`  - 当前标签ID: ${currentTabId.value}`)
      console.log(`  - 账号数量: ${appStore.passwords.length}`)
    }
    
    /**
     * 切换搜索栏显示
     */
    const toggleSearch = () => {
      showSearchBar.value = !showSearchBar.value
      if (!showSearchBar.value) {
        clearSearch()
      }
    }

    /**
     * 处理数据刷新事件
     * @author 20251004 陈凤庆
     * @description 导入成功后刷新分组、分类、账号数据
     */
    const handleRefreshData = async () => {
      try {
        console.log('[数据刷新] 开始刷新数据...')

        // 1. 重新加载分组数据
        console.log('[数据刷新] 重新加载分组数据...')
        const groupsData = await apiService.getGroups()
        appStore.setGroups(groupsData)

        // 2. 如果有当前分组，重新加载分类和账号数据
        if (appStore.currentGroupId && appStore.currentGroupId !== 'search-result') {
          console.log(`[数据刷新] 重新加载分组${appStore.currentGroupId}的数据...`)

          // 重新加载分类（标签）
          await loadTabsByGroup(appStore.currentGroupId, false)

          // 重新加载账号数据
          if (appStore.currentTabId) {
            // 如果有选中的标签，加载该标签下的账号
            console.log(`[数据刷新] 重新加载标签${appStore.currentTabId}的账号...`)
            const accounts = await apiService.getAccountsByTab(appStore.currentTabId)
            appStore.setPasswords(accounts)
          } else {
            // 如果没有选中标签，加载整个分组的账号
            console.log(`[数据刷新] 重新加载分组${appStore.currentGroupId}的账号...`)
            await loadPasswordsByGroup(appStore.currentGroupId)
          }
        }

        console.log('[数据刷新] ✅ 数据刷新完成')
        ElMessage.success(t('success.dataRefreshed'))
      } catch (error) {
        console.error('[数据刷新] 刷新数据失败:', error)
        ElMessage.error('刷新数据失败: ' + (error.message || '未知错误'))
      }
    }

    /**
     * 处理搜索
     * @modify 20251018 陈凤庆 添加 AbortController 取消过期的搜索请求，解决快速输入时的时序问题
     */
    const handleSearch = async () => {
      if (searchKeyword.value.trim() === '') {
        clearSearch()
        return
      }

      try {
        // 取消之前的搜索请求
        if (searchAbortController) {
          console.log('[搜索] 取消之前的搜索请求')
          searchAbortController.abort()
        }

        // 创建新的 AbortController
        searchAbortController = new AbortController()
        const currentKeyword = searchKeyword.value

        console.log('[搜索] 开始搜索:', currentKeyword)

        // 搜索账号
        const passwords = await apiService.searchPasswords(currentKeyword)

        // 检查是否被取消
        if (searchAbortController.signal.aborted) {
          console.log('[搜索] 搜索请求已被取消，放弃处理结果')
          return
        }

        console.log('[搜索] 搜索到账号:', passwords.length)

        // 搜索分组
        const groups = await apiService.getGroups()

        // 再次检查是否被取消
        if (searchAbortController.signal.aborted) {
          console.log('[搜索] 搜索请求已被取消，放弃处理结果')
          return
        }

        const filteredGroups = groups.filter(g =>
          g.name.toLowerCase().includes(currentKeyword.toLowerCase())
        )
        console.log('[搜索] 搜索到分组:', filteredGroups.length)

        // 创建搜索结果虚拟分组
        const searchResultGroup = {
          id: 'search-result',
          name: `搜索结果 (${passwords.length})`,
          icon: 'fa-search',
          isSearchResult: true, // 标记为搜索结果，避免保存到数据库
          sort_order: 999 // 排在最后面
        }

        // 将搜索结果分组添加到分组列表的末尾
        const groupsWithSearch = [...groups, searchResultGroup]
        appStore.setGroups(groupsWithSearch)

        // 切换到搜索结果分组
        appStore.setCurrentGroup('search-result')

        // 设置搜索结果
        searchResults.value = {
          groups: filteredGroups,
          passwords: passwords
        }

        appStore.setSearchKeyword(currentKeyword)
        appStore.setSearchResults(passwords)

        console.log('[搜索] 搜索完成，切换到搜索结果分组')
      } catch (error) {
        // 如果是 AbortError，说明请求被取消，不需要显示错误
        if (error.name === 'AbortError') {
          console.log('[搜索] 搜索请求已被取消')
          return
        }
        console.error('搜索失败:', error)
        ElMessage.error('搜索失败')
      }
    }
    
    /**
     * 清除搜索
     */
    const clearSearch = async () => {
      console.log('[搜索] 清除搜索')

      searchKeyword.value = ''
      searchResults.value = { groups: [], accounts: [] }

      // 重新加载原始分组列表，移除搜索结果分组
      try {
        const originalGroups = await apiService.getGroups()
        appStore.setGroups(originalGroups)

        // 如果当前选中的是搜索结果分组，切换到第一个真实分组
        if (appStore.currentGroupId === 'search-result' && originalGroups.length > 0) {
          const firstGroupId = originalGroups[0].id
          console.log('[搜索] 切换回第一个分组:', firstGroupId)
          await selectGroup(firstGroupId)
        } else if (appStore.currentGroupId && appStore.currentGroupId !== 'search-result') {
          // 如果当前分组不是搜索结果分组，重新加载当前分组的标签和账号
          console.log('[搜索] 重新加载当前分组:', appStore.currentGroupId)
          await loadTabsByGroup(appStore.currentGroupId, false)
          await loadPasswordsByGroup(appStore.currentGroupId)
        }
      } catch (error) {
        console.error('[搜索] 清除搜索时重新加载分组失败:', error)
      }

      // 最后清除搜索状态，确保标签数据已加载
      appStore.clearSearch()

      console.log('[搜索] 搜索已清除')
    }
    
    /**
     * 显示分组菜单
     * @param {Event} event 右键事件
     * @param {Object} group 分组对象
     * @author 20251002 陈凤庆 实现分组右键菜单功能
     */
    const showGroupMenu = (event, group) => {
      console.log('[MainView] 显示分组菜单:', group, event)

      // 隐藏其他菜单
      hideContextMenu()

      // 设置菜单位置和分组信息
      groupContextMenuPosition.x = event.clientX
      groupContextMenuPosition.y = event.clientY
      contextMenuGroup.value = group
      showGroupContextMenu.value = true

      // 点击其他地方关闭菜单
      const handleClickOutside = () => {
        hideGroupContextMenu()
        document.removeEventListener('click', handleClickOutside)
      }

      // 延迟添加事件监听器，避免立即触发
      setTimeout(() => {
        document.addEventListener('click', handleClickOutside)
      }, 0)
    }

    /**
     * 隐藏分组右键菜单
     * @author 20251002 陈凤庆 隐藏分组右键菜单
     */
    const hideGroupContextMenu = () => {
      showGroupContextMenu.value = false
      contextMenuGroup.value = null
    }
    
    /**
     * 显示页签菜单
     * @param {Event} event 右键事件
     * @param {Object} tab 页签对象
     * @author 陈凤庆
     * @date 20250101
     * @modify 20251002 陈凤庆 改为使用TabContextMenu组件
     * @description 实现标签右键菜单：重命名、删除、上移、下移
     */
    const showTabMenu = (event, tab) => {
      console.log('[MainView] 显示标签菜单:', tab, event)

      // 隐藏其他菜单
      hideContextMenu()
      hideGroupContextMenu()

      // 设置菜单位置和标签信息
      tabContextMenuPosition.x = event.clientX
      tabContextMenuPosition.y = event.clientY
      contextMenuTab.value = tab
      showTabContextMenu.value = true

      // 点击其他地方关闭菜单
      const handleClickOutside = () => {
        hideTabContextMenu()
        document.removeEventListener('click', handleClickOutside)
      }

      // 延迟添加事件监听器，避免立即触发
      setTimeout(() => {
        document.addEventListener('click', handleClickOutside)
      }, 0)
    }

    /**
     * 隐藏标签右键菜单
     * @author 20251002 陈凤庆 隐藏标签右键菜单
     */
    const hideTabContextMenu = () => {
      showTabContextMenu.value = false
      contextMenuTab.value = null
    }

    /**
     * 处理标签重命名
     * @param {string} tabId 标签ID
     * @param {string} newName 新名称
     * @author 20251002 陈凤庆 处理标签重命名
     */
    const handleRenameTab = async (tabId, newName) => {
      try {
        console.log(`[handleRenameTab] 重命名标签: ${tabId} -> ${newName}`)

        // 找到要重命名的标签
        const tab = appStore.tabs.find(t => t.id === tabId)
        if (!tab) {
          ElMessage.error('标签不存在')
          return
        }

        // 更新标签
        const updatedTab = {
          ...tab,
          name: newName
        }
        await apiService.updateTab(updatedTab)

        // 重新加载标签列表
        await loadTabsByGroup(currentGroupId.value)

        ElMessage.success('标签重命名成功')
      } catch (error) {
        console.error('[handleRenameTab] 重命名标签失败:', error)
        ElMessage.error(`重命名标签失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 处理标签删除
     * @param {string} tabId 标签ID
     * @author 20251002 陈凤庆 处理标签删除
     * @modify 20251002 陈凤庆 修复错误信息显示，正确处理Wails错误对象
     */
    const handleDeleteTab = async (tabId) => {
      try {
        console.log(`[handleDeleteTab] 删除标签: ${tabId}`)

        // 调用API删除标签
        await apiService.deleteTab(tabId)

        // 重新加载标签列表
        await loadTabsByGroup(currentGroupId.value)

        ElMessage.success('标签删除成功')
      } catch (error) {
        console.error('[handleDeleteTab] 删除标签失败:', error)
        console.error('[handleDeleteTab] 错误类型:', typeof error)
        console.error('[handleDeleteTab] 错误对象:', error)

        // 20251002 陈凤庆 修复错误信息处理，支持Wails错误格式
        let errorMessage = '未知错误'
        if (typeof error === 'string') {
          errorMessage = error
        } else if (error && error.message) {
          errorMessage = error.message
        } else if (error && typeof error === 'object') {
          errorMessage = error.toString()
        }

        ElMessage.error(`删除标签失败: ${errorMessage}`)
      }
    }

    /**
     * 处理标签上移
     * @param {string} tabId 标签ID
     * @author 20251002 陈凤庆 处理标签上移
     */
    const handleMoveTabUp = async (tabId) => {
      try {
        console.log(`[handleMoveTabUp] 标签上移: ${tabId}`)
        await apiService.moveTabUp(tabId)

        // 重新加载标签列表
        await loadTabsByGroup(currentGroupId.value)

        ElMessage.success('标签上移成功')
      } catch (error) {
        console.error('[handleMoveTabUp] 标签上移失败:', error)
        ElMessage.error(`标签上移失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 处理标签下移
     * @param {string} tabId 标签ID
     * @author 20251002 陈凤庆 处理标签下移
     */
    const handleMoveTabDown = async (tabId) => {
      try {
        console.log(`[handleMoveTabDown] 标签下移: ${tabId}`)
        await apiService.moveTabDown(tabId)

        // 重新加载标签列表
        await loadTabsByGroup(currentGroupId.value)

        ElMessage.success('标签下移成功')
      } catch (error) {
        console.error('[handleMoveTabDown] 标签下移失败:', error)
        ElMessage.error(`标签下移失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 处理标签拖拽排序
     * @param {string} sourceTabId 源标签ID
     * @param {string} targetTabId 目标标签ID
     * @author 20251002 陈凤庆 处理标签拖拽排序
     */
    const handleTabDragReorder = async (sourceTabId, targetTabId) => {
      try {
        console.log(`[handleTabDragReorder] 拖拽排序: ${sourceTabId} -> ${targetTabId}`)

        // 找到源标签和目标标签
        const sourceTab = appStore.tabs.find(t => t.id === sourceTabId)
        const targetTab = appStore.tabs.find(t => t.id === targetTabId)

        if (!sourceTab || !targetTab) {
          console.error('[handleTabDragReorder] 找不到标签:', { sourceTabId, targetTabId })
          return
        }

        // 交换排序号
        const sourceSortOrder = sourceTab.sort_order
        const targetSortOrder = targetTab.sort_order

        // 更新源标签的排序号
        await apiService.updateTabSortOrder(sourceTabId, targetSortOrder)
        // 更新目标标签的排序号
        await apiService.updateTabSortOrder(targetTabId, sourceSortOrder)

        // 重新加载标签列表
        await loadTabsByGroup(currentGroupId.value)

        console.log('[handleTabDragReorder] 拖拽排序完成')
      } catch (error) {
        console.error('[handleTabDragReorder] 拖拽排序失败:', error)
        ElMessage.error(`拖拽排序失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 处理在指定标签后创建新标签
     * @param {string} afterTabId 在此标签后插入
     * @param {string} tabName 新标签名称
     * @author 20251002 陈凤庆 处理在指定标签后创建新标签
     */
    const handleCreateTabAfter = async (afterTabId, tabName) => {
      try {
        console.log(`[handleCreateTabAfter] 在标签 ${afterTabId} 后创建新标签: ${tabName}`)

        // 检查是否选中了分组
        if (!currentGroupId.value) {
          ElMessage.warning(t('warning.selectGroupFirst'))
          return
        }

        // 调用API在指定标签后创建新标签
        const newTab = await apiService.insertTabAfter(
          tabName,
          currentGroupId.value,
          'fa-tag', // 默认图标
          afterTabId
        )

        console.log(`[handleCreateTabAfter] 标签创建成功:`, newTab)

        // 重新加载当前分组的标签列表
        await loadTabsByGroup(currentGroupId.value)

        // 选中新创建的标签
        if (newTab && newTab.id) {
          appStore.setCurrentTab(newTab.id)
          console.log(`[handleCreateTabAfter] 已选中新创建的标签: ${newTab.name} (ID: ${newTab.id})`)
        }

        ElMessage.success(t('success.tabCreated'))
      } catch (error) {
        console.error('[handleCreateTabAfter] 创建标签失败:', error)
        ElMessage.error(`创建标签失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 处理分组拖拽排序
     * @param {string} sourceGroupId 源分组ID
     * @param {string} targetGroupId 目标分组ID
     * @author 20251002 陈凤庆 处理分组拖拽排序
     */
    const handleGroupDragReorder = async (sourceGroupId, targetGroupId) => {
      try {
        console.log(`[handleGroupDragReorder] 拖拽排序: ${sourceGroupId} -> ${targetGroupId}`)

        // 找到源分组和目标分组
        const sourceGroup = appStore.groups.find(g => g.id === sourceGroupId)
        const targetGroup = appStore.groups.find(g => g.id === targetGroupId)

        if (!sourceGroup || !targetGroup) {
          console.error('[handleGroupDragReorder] 找不到分组:', { sourceGroupId, targetGroupId })
          return
        }

        // 交换排序号
        const sourceSortOrder = sourceGroup.sort_order
        const targetSortOrder = targetGroup.sort_order

        // 更新源分组的排序号
        await apiService.updateGroupSortOrder(sourceGroupId, targetSortOrder)
        // 更新目标分组的排序号
        await apiService.updateGroupSortOrder(targetGroupId, sourceSortOrder)

        // 重新加载分组列表
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        console.log('[handleGroupDragReorder] 拖拽排序完成')
      } catch (error) {
        console.error('[handleGroupDragReorder] 拖拽排序失败:', error)
        ElMessage.error(`拖拽排序失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 处理右键菜单新建分组
     * @param {string} groupName 分组名称
     * @author 20251002 陈凤庆 处理右键菜单新建分组
     */
    const handleCreateGroupFromMenu = async (groupName) => {
      try {
        console.log(`[handleCreateGroupFromMenu] 创建分组: ${groupName}`)

        // 调用API创建分组
        const newGroup = await apiService.createGroup(groupName)
        console.log(`[handleCreateGroupFromMenu] 分组创建成功:`, newGroup)

        // 重新加载分组列表
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        // 切换到新创建的分组
        if (newGroup && newGroup.id) {
          await selectGroup(newGroup.id)
          console.log(`[handleCreateGroupFromMenu] 已切换到新创建的分组: ${newGroup.name} (ID: ${newGroup.id})`)
        }

        ElMessage.success('分组创建成功')
      } catch (error) {
        console.error('[handleCreateGroupFromMenu] 创建分组失败:', error)
        ElMessage.error(`创建分组失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 重命名标签
     * @param {Object} tab 标签对象
     * @author 陈凤庆
     * @date 20250101
     */
    const renameTab = async (tab) => {
      try {
        const { value: newName } = await ElMessageBox.prompt(
          t('tabContextMenu.promptNewName'),
          t('tabContextMenu.renameTab'),
          {
            confirmButtonText: t('tabContextMenu.confirm'),
            cancelButtonText: t('tabContextMenu.cancel'),
            inputPlaceholder: t('tabContextMenu.tabName'),
            inputValue: tab.name,
            inputValidator: (value) => {
              if (!value || value.trim() === '') {
                return t('tabContextMenu.tabNameCannotBeEmpty')
              }
              if (value.trim() === tab.name.trim()) {
                return t('tabContextMenu.tabNameNotChanged')
              }
              return true
            }
          }
        )

        if (newName && newName.trim() !== tab.name.trim()) {
          console.log(`[renameTab] 重命名标签: ${tab.name} -> ${newName}`)

          // 更新标签
          const updatedTab = {
            ...tab,
            name: newName.trim()
          }
          await apiService.updateTab(updatedTab)

          // 重新加载标签列表
          await loadTabsByGroup(currentGroupId.value)

          ElMessage.success(t('success.tabRenamed'))
        }
      } catch (error) {
        if (error === 'cancel') {
          return
        }
        console.error('[renameTab] 重命名标签失败:', error)
        ElMessage.error(`${t('error.renameTabFailed')}: ${error.message || t('error.unknownError')}`)
      }
    }

    /**
     * 删除标签
     * @param {Object} tab 标签对象
     * @author 陈凤庆
     * @date 20250101
     * @modify 20251002 陈凤庆 修复错误信息显示，正确处理Wails错误对象
     */
    const deleteTab = async (tab) => {
      try {
        await ElMessageBox.confirm(
          `确定要删除标签"${tab.name}"吗？`,
          '删除标签',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        console.log(`[deleteTab] 删除标签: ${tab.id} (${tab.name})`)

        // 调用API删除标签
        await apiService.deleteTab(tab.id)

        // 重新加载标签列表
        await loadTabsByGroup(currentGroupId.value)

        ElMessage.success('标签删除成功')
      } catch (error) {
        if (error === 'cancel') {
          return
        }
        console.error('[deleteTab] 删除标签失败:', error)
        console.error('[deleteTab] 错误类型:', typeof error)
        console.error('[deleteTab] 错误对象:', error)

        // 20251002 陈凤庆 修复错误信息处理，支持Wails错误格式
        let errorMessage = '未知错误'
        if (typeof error === 'string') {
          errorMessage = error
        } else if (error && error.message) {
          errorMessage = error.message
        } else if (error && typeof error === 'object') {
          errorMessage = error.toString()
        }

        ElMessage.error(`删除标签失败: ${errorMessage}`)
      }
    }
    
    /**
     * 创建分组
     * 20251017 陈凤庆 修复取消按钮提示问题
     */
    const createGroup = async () => {
      try {
        const { value: groupName } = await ElMessageBox.prompt(
          '请输入分组名称',
          '创建分组',
          {
            confirmButtonText: '创建',
            cancelButtonText: '取消',
            inputPlaceholder: '分组名称'
          }
        )

        if (!groupName) return

        // 20251002 陈凤庆 使用字符串ID，避免JavaScript精度丢失
        // 20251002 陈凤庆 删除parentID参数，后端不需要层级结构
        const newGroup = await apiService.createGroup(groupName)
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        // 切换到新创建的分组
        if (newGroup && newGroup.id) {
          await selectGroup(newGroup.id)
        }

        ElMessage.success('分组创建成功')
      } catch (error) {
        // 20251017 陈凤庆 修复取消按钮提示问题：用户取消操作时不显示错误提示
        if (error === 'cancel') {
          console.log('用户取消创建分组操作')
          return
        }
        console.error('创建分组失败:', error)
        ElMessage.error('创建分组失败')
      }
    }

    /**
     * 处理重命名分组
     * @param {string} groupId 分组ID
     * @param {string} newName 新的分组名称
     * @author 20251002 陈凤庆 处理分组重命名
     */
    const handleRenameGroup = async (groupId, newName) => {
      try {
        console.log(`重命名分组: ${groupId} -> ${newName}`)
        await apiService.renameGroup(groupId, newName)

        // 刷新分组列表
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        ElMessage.success('分组重命名成功')
      } catch (error) {
        console.error('重命名分组失败:', error)
        ElMessage.error('重命名分组失败: ' + error.message)
      }
    }

    /**
     * 处理删除分组
     * @param {string} groupId 分组ID
     * @author 20251002 陈凤庆 处理分组删除
     * @modify 20251002 陈凤庆 修复错误信息显示，正确处理Wails错误对象
     */
    const handleDeleteGroup = async (groupId) => {
      try {
        console.log(`删除分组: ${groupId}`)
        await apiService.deleteGroup(groupId)

        // 刷新分组列表
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        // 如果删除的是当前分组，切换到第一个分组
        if (currentGroupId.value === groupId && updatedGroups.length > 0) {
          await selectGroup(updatedGroups[0].id)
        }

        ElMessage.success('分组删除成功')
      } catch (error) {
        console.error('删除分组失败:', error)
        console.error('错误类型:', typeof error)
        console.error('错误对象:', error)

        // 20251002 陈凤庆 修复错误信息处理，支持Wails错误格式
        let errorMessage = '未知错误'
        if (typeof error === 'string') {
          errorMessage = error
        } else if (error && error.message) {
          errorMessage = error.message
        } else if (error && typeof error === 'object') {
          errorMessage = error.toString()
        }

        ElMessage.error('删除分组失败: ' + errorMessage)
      }
    }

    /**
     * 处理分组左移
     * @param {string} groupId 分组ID
     * @author 20251002 陈凤庆 处理分组左移
     */
    const handleMoveGroupLeft = async (groupId) => {
      try {
        console.log(`分组左移: ${groupId}`)
        await apiService.moveGroupLeft(groupId)

        // 刷新分组列表
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        ElMessage.success('分组左移成功')
      } catch (error) {
        console.error('分组左移失败:', error)
        ElMessage.error('分组左移失败: ' + error.message)
      }
    }

    /**
     * 处理分组右移
     * @param {string} groupId 分组ID
     * @author 20251002 陈凤庆 处理分组右移
     */
    const handleMoveGroupRight = async (groupId) => {
      try {
        console.log(`分组右移: ${groupId}`)
        await apiService.moveGroupRight(groupId)

        // 刷新分组列表
        const updatedGroups = await apiService.getGroups()
        appStore.setGroups(updatedGroups)

        ElMessage.success('分组右移成功')
      } catch (error) {
        console.error('分组右移失败:', error)
        ElMessage.error('分组右移失败: ' + error.message)
      }
    }

    /**
     * 创建页签
     * @author 陈凤庆
     * @date 20250101
     * @description 弹出对话框输入标签名称，创建标签并保存到数据库
     */
    const createTab = async () => {
      try {
        // 检查是否选中了分组
        if (!currentGroupId.value) {
          ElMessage.warning(t('warning.selectGroupFirst'))
          return
        }

        // 弹出输入框
        const { value: tabName } = await ElMessageBox.prompt(
          t('tabContextMenu.promptNewTabName'),
          t('tabContextMenu.newTab'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            inputPattern: /\S+/,
            inputErrorMessage: t('tabContextMenu.tabNameCannotBeEmpty')
          }
        )

        if (!tabName || !tabName.trim()) {
          return
        }

        console.log(`[createTab] 开始创建标签: ${tabName}`)

        // 调用API创建标签，插入到最后（afterTypeID为空）
        const newTab = await apiService.insertTabAfter(
          tabName.trim(),
          currentGroupId.value,
          'fa-tag', // 默认图标
          '' // 插入到最后
        )

        console.log(`[createTab] 标签创建成功:`, newTab)

        // 重新加载当前分组的标签列表
        await loadTabsByGroup(currentGroupId.value)

        // 选中新创建的标签
        if (newTab && newTab.id) {
          appStore.setCurrentTab(newTab.id)
          console.log(`[createTab] 已选中新创建的标签: ${newTab.name} (ID: ${newTab.id})`)
        }

        ElMessage.success(t('success.tabCreated'))
      } catch (error) {
        // 用户取消操作
        if (error === 'cancel') {
          console.log('[createTab] 用户取消创建标签')
          return
        }
        console.error('[createTab] 创建标签失败:', error)
        ElMessage.error(`${t('error.createTabFailed')}: ${error.message || t('error.unknownError')}`)
      }
    }
    
    /**
     * 创建账号项
     * @author 陈凤庆
     * @date 20250101
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     * @modify 20251019 陈凤庆 修复问题 002-1：切换标签后创建账号失败，确保在没有选中标签时提示用户先选择标签
     * @modify 20251019 陈凤庆 修复问题 003-1：增加详细调试日志，跟踪切换分组后的状态
     * @description 创建新账号项，传递当前分组ID和标签ID到编辑窗口
     */
    const createAccount = () => {
      console.log(`[createAccount] ========== 开始创建账号项 ==========`)
      console.log(`[createAccount] 当前分组ID: ${currentGroupId.value} (类型: ${typeof currentGroupId.value})`)
      console.log(`[createAccount] 当前标签ID: ${currentTabId.value} (类型: ${typeof currentTabId.value})`)
      console.log(`[createAccount] Store中的分组ID: ${appStore.currentGroupId} (类型: ${typeof appStore.currentGroupId})`)
      console.log(`[createAccount] Store中的标签ID: ${appStore.currentTabId} (类型: ${typeof appStore.currentTabId})`)
      console.log(`[createAccount] 当前分组对象:`, appStore.currentGroup)
      console.log(`[createAccount] 当前标签对象:`, appStore.currentTab)
      console.log(`[createAccount] 标签列表:`, appStore.tabs)

      // 20251019 陈凤庆 修复问题 002-1：检查是否选中了标签，如果没有选中标签则提示用户
      if (!currentTabId.value) {
        console.warn(`[createAccount] ⚠️ 未选中标签，无法创建账号`)
        ElMessage.warning(t('warning.selectTabFirst') || '请先选择一个标签')
        return
      }

      // 20251019 陈凤庆 修复问题 003-1：验证分组和标签的有效性
      if (!currentGroupId.value) {
        console.error(`[createAccount] ❌ 未选中分组，无法创建账号`)
        ElMessage.error('请先选择一个分组')
        return
      }

      // 验证标签是否存在于当前标签列表中
      const currentTabExists = appStore.tabs.find(tab => tab.id === currentTabId.value)
      if (!currentTabExists) {
        console.error(`[createAccount] ❌ 当前标签ID ${currentTabId.value} 不存在于标签列表中`)
        console.error(`[createAccount] 可用标签列表:`, appStore.tabs.map(t => ({id: t.id, name: t.name})))
        ElMessage.error('当前标签无效，请重新选择标签')
        return
      }

      console.log(`[createAccount] ✅ 验证通过，开始创建账号数据结构`)

      // 20251002 陈凤庆 确保所有ID字段都是字符串类型，新建时id为空字符串，使用typeid字段
      selectedAccount.value = {
        id: '', // 20251001 陈凤庆 新建时ID为空字符串，后端会自动生成GUID
        title: '',
        username: '',
        password: '',
        url: '',
        type: '',
        notes: '',
        group_id: currentGroupId.value || '', // 20251001 陈凤庆 确保为字符串
        typeid: currentTabId.value || '', // 20251002 陈凤庆 改为typeid，与后端模型保持一致
        icon: '',
        is_favorite: false,
        input_method: 1 // 20251003 陈凤庆 添加输入方式字段，默认为1（默认方式）
      }

      console.log(`[createAccount] 创建的账号数据结构:`, selectedAccount.value)
      console.log(`[createAccount] ========== 账号项创建完成，打开编辑窗口 ==========`)

      isNewAccount.value = true
      showAccountDetail.value = true
    }
    
    /**
     * 显示账号详情
     * @param {Object} account 账号对象
     * @modify 20251002 陈凤庆 密码详情改为账号详情，确保传递当前分组ID，用于加载对应的类型列表
     * @modify 20251003 陈凤庆 调用GetAccountDetail接口获取完整账号详情，密码字段默认显示5个*
     */
    const showAccountDetailDialog = async (account) => {
      try {
        console.log(`[showAccountDetailDialog] 开始获取账号详情，账号ID: ${account.id}`)

        // 调用新的GetAccountDetail接口获取账号详情
        const accountDetail = await apiService.getAccountDetail(account.id)
        console.log(`[showAccountDetailDialog] 获取到账号详情:`, accountDetail)

        // 20251002 陈凤庆 确保账号对象包含当前分组ID，用于类型加载，使用typeid字段
        selectedAccount.value = {
          ...accountDetail,
          password: "*****", // 密码字段默认显示5个*
          group_id: accountDetail.group_id || currentGroupId.value || '', // 确保有group_id
          typeid: accountDetail.typeid || currentTabId.value || '' // 20251002 陈凤庆 改为typeid字段
        }
        isNewAccount.value = false
        showAccountDetail.value = true

        console.log(`[showAccountDetailDialog] 显示账号详情，分组ID: ${selectedAccount.value.group_id}, 类型ID: ${selectedAccount.value.typeid}`)
      } catch (error) {
        console.error('[showAccountDetailDialog] 获取账号详情失败:', error)
        ElMessage.error(`获取账号详情失败: ${error.message || '未知错误'}`)
      }
    }
    
    /**
     * 输入用户名和密码
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 添加窗口聚焦管理，确保输入到正确的目标窗口
     * @modify 20251002 陈凤庆 优化提示信息，移除敏感信息显示
     * @modify 20251003 陈凤庆 修复逻辑：根据账号ID查询解密后的完整数据，然后根据输入方法进行输入
     */
    const inputUsernameAndPassword = async (account) => {
      try {
        // 检查辅助功能权限
        const hasPermission = await window.go.app.App.CheckAccessibilityPermission()
        if (!hasPermission) {
          ElMessage.error('需要辅助功能权限才能使用自动填充功能，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用')
          return
        }

        // 20250127 陈凤庆 获取当前存储的目标应用程序名称
        const targetAppName = await window.go.app.App.GetPreviousFocusedAppName()

        console.log(`[inputUsernameAndPassword] 开始输入用户名和密码，账号ID: ${account.id}`)

        // 20251003 陈凤庆 修复：后端会根据账号ID查询解密后的完整数据，然后根据输入方法进行输入
        // 前端只需要传递账号ID，不需要传递用户名和密码（避免使用脱敏数据）
        await window.go.app.App.SimulateUsernameAndPassword(account.id)

        // 更新账号使用次数
        await updateAccountUsageAndRefresh(account.id)

        ElMessage.success(`已向 ${targetAppName} 应用程序自动填充用户名和密码`)
      } catch (error) {
        console.error('自动填充失败:', error)
        ElMessage.error(`自动填充失败: ${error.message || error}`)
      }
    }

    /**
     * 输入用户名
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 添加窗口聚焦管理，确保输入到正确的目标窗口
     * @modify 20251002 陈凤庆 优化提示信息，移除敏感信息显示
     * @modify 20251003 陈凤庆 修复逻辑：根据账号ID查询解密后的用户名，然后根据输入方法进行输入
     */
    const inputUsername = async (account) => {
      try {
        // 检查辅助功能权限
        const hasPermission = await window.go.app.App.CheckAccessibilityPermission()
        if (!hasPermission) {
          ElMessage.error('需要辅助功能权限才能使用自动填充功能，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用')
          return
        }

        // 20250127 陈凤庆 获取当前存储的目标应用程序名称
        const targetAppName = await window.go.app.App.GetPreviousFocusedAppName()

        console.log(`[inputUsername] 开始输入用户名，账号ID: ${account.id}`)

        // 20251003 陈凤庆 修复：后端会根据账号ID查询解密后的用户名，然后根据输入方法进行输入
        // 前端只需要传递账号ID，不需要传递用户名（避免使用脱敏数据）
        await window.go.app.App.SimulateUsername(account.id)

        // 更新账号使用次数
        await updateAccountUsageAndRefresh(account.id)

        ElMessage.success(`已向 ${targetAppName} 应用程序自动填充用户名`)
      } catch (error) {
        console.error('自动填充用户名失败:', error)
        ElMessage.error(`自动填充用户名失败: ${error.message || error}`)
      }
    }
    
    /**
     * 输入密码
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 添加窗口聚焦管理，确保输入到正确的目标窗口
     * @modify 20251002 陈凤庆 优化提示信息，移除敏感信息显示
     * @modify 20251003 陈凤庆 修复逻辑：根据账号ID查询解密后的密码，然后根据输入方法进行输入
     */
    const inputPassword = async (account) => {
      try {
        // 检查辅助功能权限
        const hasPermission = await window.go.app.App.CheckAccessibilityPermission()
        if (!hasPermission) {
          ElMessage.error('需要辅助功能权限才能使用自动填充功能，请在系统偏好设置 > 安全性与隐私 > 隐私 > 辅助功能中添加此应用')
          return
        }

        // 20250127 陈凤庆 获取当前存储的目标应用程序名称
        const targetAppName = await window.go.app.App.GetPreviousFocusedAppName()

        console.log(`[inputPassword] 开始输入密码，账号ID: ${account.id}`)

        // 20251003 陈凤庆 修复：后端会根据账号ID查询解密后的密码，然后根据输入方法进行输入
        // 前端只需要传递账号ID，不需要传递密码（避免使用脱敏数据）
        await window.go.app.App.SimulatePassword(account.id)

        // 更新账号使用次数
        await updateAccountUsageAndRefresh(account.id)

        ElMessage.success(`已向 ${targetAppName} 应用程序自动填充密码`)
      } catch (error) {
        console.error('自动填充密码失败:', error)
        ElMessage.error(`自动填充密码失败: ${error.message || error}`)
      }
    }

    /**
     * 更新账号使用次数并刷新显示
     * @param {string} accountId 账号ID
     * @author 20251003 陈凤庆 新增使用次数统计功能
     */
    const updateAccountUsageAndRefresh = async (accountId) => {
      try {
        console.log('[使用统计] 更新账号使用次数:', accountId)

        // 调用后端API更新使用次数
        await apiService.updateAccountUsage(accountId)

        // 重新加载当前分组的账号列表以刷新使用次数显示
        if (appStore.currentGroupId && appStore.currentGroupId !== 'search-result') {
          await loadPasswordsByGroup(appStore.currentGroupId)
        }

        // 重新计算总使用次数
        await refreshTotalUseCount()

        console.log('[使用统计] 账号使用次数更新完成')
      } catch (error) {
        console.error('[使用统计] 更新账号使用次数失败:', error)
        // 不显示错误消息，避免影响用户体验
      }
    }

    /**
     * 刷新总使用次数
     * @author 20251003 陈凤庆 新增使用次数统计功能
     */
    const refreshTotalUseCount = async () => {
      try {
        // 获取所有账号的使用次数总和
        const accounts = await apiService.getAllAccounts()
        const total = accounts.reduce((sum, account) => sum + (account.use_count || 0), 0)
        totalUseCount.value = total
        console.log('[使用统计] 总使用次数已更新:', total)
      } catch (error) {
        console.error('[使用统计] 刷新总使用次数失败:', error)
      }
    }

    /**
     * 刷新使用天数
     * @author 20251003 陈凤庆 新增使用天数统计功能
     */
    const refreshUsageDays = async () => {
      try {
        const days = await apiService.getUsageDays()
        usageDays.value = days
        console.log('[使用统计] 使用天数已更新:', days)
      } catch (error) {
        console.error('[使用统计] 刷新使用天数失败:', error)
      }
    }

    /**
     * 处理账号项保存
     * @param {Object} accountData 账号项数据
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     * @modify 20251002 陈凤庆 实现更新账号功能，保存后刷新列表并定位到当前账号
     * @modify 20251019 陈凤庆 修复问题 003-1：增加详细调试日志，跟踪保存过程中的数据状态
     */
    const handleSaveAccount = async (accountData) => {
      try {
        console.log(`[handleSaveAccount] ========== 开始保存账号 ==========`)
        console.log(`[handleSaveAccount] 原始账号数据:`, accountData)
        console.log(`[handleSaveAccount] 是否为新建账号: ${isNewAccount.value}`)
        console.log(`[handleSaveAccount] 当前分组ID: ${currentGroupId.value}`)
        console.log(`[handleSaveAccount] 当前标签ID: ${currentTabId.value}`)

        // 20251001 陈凤庆 确保所有ID字段都是字符串类型
        const cleanedData = {
          ...accountData,
          id: accountData.id ? String(accountData.id) : '',
          group_id: accountData.group_id ? String(accountData.group_id) : '',
          typeid: accountData.typeid ? String(accountData.typeid) : '' // 20251002 陈凤庆 改为typeid字段
        }

        console.log('[handleSaveAccount] 清理后的账号数据:', cleanedData)

        // 20251019 陈凤庆 修复问题 003-1：验证关键字段
        if (!cleanedData.typeid) {
          console.error('[handleSaveAccount] ❌ typeid字段为空，这会导致保存失败')
          console.error('[handleSaveAccount] 当前数据状态检查:')
          console.error('  - cleanedData.typeid:', cleanedData.typeid)
          console.error('  - accountData.typeid:', accountData.typeid)
          console.error('  - currentTabId.value:', currentTabId.value)
          console.error('  - appStore.currentTabId:', appStore.currentTabId)
          ElMessage.error('标签信息丢失，请重新选择标签后再保存')
          return
        }

        if (!cleanedData.group_id) {
          console.error('[handleSaveAccount] ❌ group_id字段为空，这会导致保存失败')
          console.error('[handleSaveAccount] 当前数据状态检查:')
          console.error('  - cleanedData.group_id:', cleanedData.group_id)
          console.error('  - accountData.group_id:', accountData.group_id)
          console.error('  - currentGroupId.value:', currentGroupId.value)
          console.error('  - appStore.currentGroupId:', appStore.currentGroupId)
          ElMessage.error('分组信息丢失，请重新选择分组后再保存')
          return
        }

        console.log('[handleSaveAccount] ✅ 关键字段验证通过，开始调用API')

        let savedAccountId = cleanedData.id

        if (isNewAccount.value) {
          console.log('[handleSaveAccount] 创建新账号...')
          console.log('[handleSaveAccount] 调用API参数详情:')
          console.log('  - title:', cleanedData.title)
          console.log('  - username:', cleanedData.username)
          console.log('  - password:', cleanedData.password ? '***已设置***' : '未设置')
          console.log('  - url:', cleanedData.url)
          console.log('  - typeid:', cleanedData.typeid)
          console.log('  - notes:', cleanedData.notes)
          console.log('  - input_method:', cleanedData.input_method)

          const newAccount = await apiService.createAccount(cleanedData)
          console.log('[handleSaveAccount] ✅ 新账号创建成功:', newAccount)
          appStore.addPassword(newAccount)
          savedAccountId = newAccount.id
        } else {
          console.log('[handleSaveAccount] 更新现有账号...')
          // 20251002 陈凤庆 调用更新账号API
          await apiService.updatePasswordItem(cleanedData)
          console.log('[handleSaveAccount] ✅ 账号更新成功')
          appStore.updatePassword(cleanedData)
        }

        // 20251002 陈凤庆 根据当前标签刷新账号列表，而不是按分组刷新
        console.log('[handleSaveAccount] 刷新账号列表...')
        if (currentTabId.value) {
          // 如果有选中的标签，刷新该标签下的账号列表
          console.log(`[handleSaveAccount] 刷新标签${currentTabId.value}的账号列表`)
          const accounts = await apiService.getAccountsByTab(currentTabId.value)
          appStore.setPasswords(accounts)
        } else {
          // 如果没有选中标签，刷新整个分组的账号列表
          console.log(`[handleSaveAccount] 刷新分组${currentGroupId.value}的账号列表`)
          await loadPasswordsByGroup(currentGroupId.value)
        }

        // 20251002 陈凤庆 定位到保存的账号
        console.log(`[handleSaveAccount] 定位到账号ID: ${savedAccountId}`)
        await scrollToAccount(savedAccountId)

        ElMessage.success(isNewAccount.value ? '账号创建成功' : '账号更新成功')

      } catch (error) {
        console.error('保存账号失败:', error)
        ElMessage.error(`保存账号失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 定位到指定账号
     * @param {string} accountId 账号ID
     * @author 20251002 陈凤庆 新增定位到指定账号的功能
     * @modify 20251003 陈凤庆 优化分页定位逻辑，支持跨页面定位账号
     * @modify 20251003 陈凤庆 使用组件引用方式调用AccountListPanel的navigateToAccountPage方法，提高可靠性
     */
    const scrollToAccount = async (accountId) => {
      if (!accountId) {
        console.warn('[scrollToAccount] 账号ID为空，无法定位')
        return
      }

      try {
        console.log(`[scrollToAccount] 开始定位到账号ID: ${accountId}`)

        // 使用组件引用调用 AccountListPanel 的 navigateToAccountPage 方法
        const mainContentArea = mainContentAreaRef.value
        if (!mainContentArea) {
          console.warn('[scrollToAccount] MainContentArea 组件引用不存在')
          return
        }

        // 调用 MainContentArea 暴露的 navigateToAccountPage 方法
        const pageChanged = mainContentArea.navigateToAccountPage(accountId)
        console.log(`[scrollToAccount] 分页切换结果: ${pageChanged}`)

        // 如果进行了分页切换，等待页面更新
        if (pageChanged) {
          console.log('[scrollToAccount] 等待分页切换完成...')
          await new Promise(resolve => setTimeout(resolve, 300))
          await nextTick()
        }

        // 在当前页面中查找目标账号行并高亮显示
        const accountRows = document.querySelectorAll('.account-row')
        console.log(`[scrollToAccount] 当前页面找到${accountRows.length}个账号行`)

        for (let i = 0; i < accountRows.length; i++) {
          const row = accountRows[i]
          const rowAccountId = row.getAttribute('data-account-id')
          console.log(`[scrollToAccount] 检查行${i}，账号ID: ${rowAccountId}`)

          if (rowAccountId === accountId) {
            console.log(`[scrollToAccount] 找到目标账号行，开始滚动定位`)

            // 滚动到目标元素
            row.scrollIntoView({
              behavior: 'smooth',
              block: 'center'
            })

            // 高亮显示该行
            row.classList.add('highlight-account')

            // 2秒后移除高亮
            setTimeout(() => {
              row.classList.remove('highlight-account')
            }, 2000)

            console.log(`[scrollToAccount] 成功定位到账号: ${accountId}`)
            return
          }
        }

        console.warn(`[scrollToAccount] 未找到账号ID为${accountId}的行`)

      } catch (error) {
        console.error('[scrollToAccount] 定位账号失败:', error)
      }
    }
    
    /**
     * 处理账号项删除
     * @param {string} accountId 账号项ID
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     * @modify 20251003 陈凤庆 修复删除逻辑，添加API调用、列表刷新和定位功能
     */
    const handleDeleteAccount = async (accountId) => {
      try {
        console.log(`[handleDeleteAccount] 开始删除账号: ${accountId}`)

        // 获取当前账号列表，用于定位删除后的位置
        const currentAccounts = appStore.filteredPasswords
        const deletedIndex = currentAccounts.findIndex(account => account.id === accountId)

        console.log(`[handleDeleteAccount] 删除账号在列表中的索引: ${deletedIndex}`)

        // 调用后端API删除账号
        await apiService.deleteAccount(accountId)
        console.log(`[handleDeleteAccount] 后端删除成功`)

        // 从Store中移除账号
        appStore.removePassword(accountId)

        // 刷新账号列表
        console.log(`[handleDeleteAccount] 刷新账号列表...`)
        if (currentTabId.value) {
          // 如果有选中的标签，刷新该标签下的账号列表
          console.log(`[handleDeleteAccount] 刷新标签${currentTabId.value}的账号列表`)
          const accounts = await apiService.getAccountsByTab(currentTabId.value)
          appStore.setPasswords(accounts)
        } else {
          // 如果没有选中标签，刷新整个分组的账号列表
          console.log(`[handleDeleteAccount] 刷新分组${currentGroupId.value}的账号列表`)
          await loadPasswordsByGroup(currentGroupId.value)
        }

        // 定位到删除后的账号
        const updatedAccounts = appStore.filteredPasswords
        if (updatedAccounts.length > 0) {
          let targetAccountId = null

          if (deletedIndex < updatedAccounts.length) {
            // 如果删除的不是最后一个，定位到下一个账号
            targetAccountId = updatedAccounts[deletedIndex].id
            console.log(`[handleDeleteAccount] 定位到下一个账号: ${targetAccountId}`)
          } else if (updatedAccounts.length > 0) {
            // 如果删除的是最后一个，定位到新的最后一个账号
            targetAccountId = updatedAccounts[updatedAccounts.length - 1].id
            console.log(`[handleDeleteAccount] 定位到最后一个账号: ${targetAccountId}`)
          }

          if (targetAccountId) {
            await scrollToAccount(targetAccountId)
          }
        }

        ElMessage.success('账号删除成功')
        console.log(`[handleDeleteAccount] 账号删除完成`)

      } catch (error) {
        console.error('[handleDeleteAccount] 删除账号失败:', error)
        ElMessage.error(`删除账号失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 显示账号项右键菜单
     * @param {MouseEvent} event 鼠标事件
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     * @modify 20251003 陈凤庆 添加菜单位置智能计算，防止超出窗口边界
     * @modify 20251003 陈凤庆 优化菜单位置计算，根据实际菜单高度和窗口位置智能调整显示方向
     * @modify 20251003 陈凤庆 修正菜单高度计算，增强边界检测和位置调整逻辑
     * @modify 20251017 陈凤庆 优化菜单高度计算，添加动态检测和更精确的边界处理
     */
    const showAccountContextMenu = (event, account) => {
      contextMenuAccount.value = account

      // 先显示菜单以便获取实际尺寸
      contextMenuPosition.x = event.clientX
      contextMenuPosition.y = event.clientY
      showContextMenu.value = true

      // 等待下一帧，让菜单渲染完成后获取实际尺寸
      nextTick(() => {
        const menuElement = document.querySelector('.password-context-menu')
        if (!menuElement) {
          console.warn('[showAccountContextMenu] 未找到菜单元素，使用默认尺寸')
          adjustMenuPosition(event.clientX, event.clientY, 160, 477)
          return
        }

        // 获取菜单的实际尺寸
        const rect = menuElement.getBoundingClientRect()
        const actualWidth = rect.width
        const actualHeight = rect.height

        console.log(`[showAccountContextMenu] 菜单实际尺寸: ${actualWidth}x${actualHeight}`)

        // 重新调整位置
        adjustMenuPosition(event.clientX, event.clientY, actualWidth, actualHeight)
      })

      // 添加全局点击监听器来隐藏菜单
      document.addEventListener('click', hideContextMenu)
    }

    /**
     * 调整菜单位置
     * @param {number} mouseX 鼠标X坐标
     * @param {number} mouseY 鼠标Y坐标
     * @param {number} menuWidth 菜单宽度
     * @param {number} menuHeight 菜单高度
     */
    const adjustMenuPosition = (mouseX, mouseY, menuWidth, menuHeight) => {
      const windowWidth = window.innerWidth
      const windowHeight = window.innerHeight
      const margin = 10 // 边界安全距离

      let x = mouseX
      let y = mouseY

      console.log(`[adjustMenuPosition] 鼠标位置: (${x}, ${y}), 窗口尺寸: ${windowWidth}x${windowHeight}, 菜单尺寸: ${menuWidth}x${menuHeight}`)

      // 水平位置智能调整
      const spaceRight = windowWidth - x // 鼠标右侧剩余空间
      const spaceLeft = x // 鼠标左侧可用空间

      if (spaceRight >= menuWidth + margin) {
        // 右侧空间足够，向右展开（默认行为）
        console.log(`[adjustMenuPosition] 右侧空间足够(${spaceRight}px)，向右展开`)
      } else if (spaceLeft >= menuWidth + margin) {
        // 右侧空间不够但左侧空间足够，向左展开
        x = x - menuWidth
        console.log(`[adjustMenuPosition] 左侧空间足够(${spaceLeft}px)，向左展开，调整x坐标为: ${x}`)
      } else {
        // 左右空间都不够，选择空间较大的一侧
        if (spaceRight > spaceLeft) {
          // 右侧空间较大，贴右边界显示
          x = windowWidth - menuWidth - margin
          console.log(`[adjustMenuPosition] 右侧空间较大，贴右边界显示，调整x坐标为: ${x}`)
        } else {
          // 左侧空间较大，贴左边界显示
          x = margin
          console.log(`[adjustMenuPosition] 左侧空间较大，贴左边界显示，调整x坐标为: ${x}`)
        }
      }

      // 垂直位置智能调整
      const spaceBelow = windowHeight - y // 鼠标下方剩余空间
      const spaceAbove = y // 鼠标上方可用空间

      console.log(`[adjustMenuPosition] 下方空间: ${spaceBelow}px, 上方空间: ${spaceAbove}px`)

      if (spaceBelow >= menuHeight + margin) {
        // 下方空间足够，向下展开（默认行为）
        console.log(`[adjustMenuPosition] 下方空间足够，向下展开`)
      } else if (spaceAbove >= menuHeight + margin) {
        // 下方空间不够但上方空间足够，向上展开
        y = y - menuHeight
        console.log(`[adjustMenuPosition] 上方空间足够，向上展开，调整y坐标为: ${y}`)
      } else {
        // 上下空间都不够，选择空间较大的一侧并调整到边界
        if (spaceBelow > spaceAbove) {
          // 下方空间较大，贴底部显示
          y = windowHeight - menuHeight - margin
          console.log(`[adjustMenuPosition] 下方空间较大，贴底部显示，调整y坐标为: ${y}`)
        } else {
          // 上方空间较大，贴顶部显示
          y = margin
          console.log(`[adjustMenuPosition] 上方空间较大，贴顶部显示，调整y坐标为: ${y}`)
        }
      }

      // 最终边界检查，确保菜单完全在窗口内
      if (x < margin) {
        x = margin
        console.log(`[adjustMenuPosition] 最终边界检查：x坐标调整为: ${x}`)
      }
      if (y < margin) {
        y = margin
        console.log(`[adjustMenuPosition] 最终边界检查：y坐标调整为: ${y}`)
      }
      if (x + menuWidth > windowWidth - margin) {
        x = windowWidth - menuWidth - margin
        console.log(`[adjustMenuPosition] 最终边界检查：x坐标调整为: ${x}`)
      }
      if (y + menuHeight > windowHeight - margin) {
        y = windowHeight - menuHeight - margin
        console.log(`[adjustMenuPosition] 最终边界检查：y坐标调整为: ${y}`)
      }

      console.log(`[adjustMenuPosition] 最终菜单位置: (${x}, ${y})`)

      // 更新菜单位置
      contextMenuPosition.x = x
      contextMenuPosition.y = y
    }

    /**
     * 隐藏右键菜单
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const hideContextMenu = () => {
      showContextMenu.value = false
      contextMenuAccount.value = null
      document.removeEventListener('click', hideContextMenu)
    }

    /**
     * 打开账号的URL
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     * @modify 20251002 陈凤庆 完善地址为空的提示信息，符合需求要求
     * @modify 20251002 陈凤庆 使用Wails的BrowserOpenURL API替换window.open，实现真正的浏览器打开功能
     */
    const openAccountUrl = (account) => {
      hideContextMenu()

      // 20251002 陈凤庆 检查地址是否为空或只包含空白字符
      if (account.url && account.url.trim()) {
        try {
          // 确保URL格式正确，如果没有协议则添加https://
          let url = account.url.trim()
          if (!url.startsWith('http://') && !url.startsWith('https://')) {
            url = 'https://' + url
          }

          // 20251002 陈凤庆 使用Wails的BrowserOpenURL API打开系统默认浏览器
          BrowserOpenURL(url)
          ElMessage.success(`正在打开 ${account.title} 的地址`)
        } catch (error) {
          console.error('打开地址失败:', error)
          ElMessage.error('打开地址失败，请检查地址格式是否正确')
        }
      } else {
        ElMessage.warning('地址为空，无法打开')
      }
    }

    /**
     * 编辑账号
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     * @modify 20251003 陈凤庆 调用GetAccountByID接口获取完整的账号数据（包括解密后的所有字段）
     */
    const editAccount = async (account) => {
      hideContextMenu()

      try {
        console.log(`[editAccount] 开始获取完整账号数据，账号ID: ${account.id}`)

        // 调用GetAccountByID接口获取完整的账号数据（包括解密后的所有字段）
        const fullAccountData = await window.go.app.App.GetAccountByID(account.id)
        console.log(`[editAccount] 获取到完整账号数据:`, fullAccountData)

        selectedAccount.value = {
          ...fullAccountData,
          group_id: fullAccountData.group_id || currentGroupId.value || '',
          typeid: fullAccountData.typeid || currentTabId.value || ''
        }
        isNewAccount.value = false
        showAccountDetail.value = true

        console.log(`[editAccount] 编辑账号数据设置完成`)
      } catch (error) {
        console.error('[editAccount] 获取完整账号数据失败:', error)
        ElMessage.error(`获取账号数据失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 生成账号副本
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-10-03
     * @description 复制当前账号信息并打开创建账号界面，预填充账号信息供用户修改
     * @modify 2025-10-03 陈凤庆 添加group_id字段，确保创建账号窗口能正确加载类型列表
     * @modify 2025-10-03 陈凤庆 修复副本生成功能，先获取完整的解密账号数据
     */
    const duplicateAccount = async (account) => {
      hideContextMenu()

      try {
        console.log(`[duplicateAccount] 开始获取完整账号数据，账号ID: ${account.id}`)

        // 调用GetAccountByID接口获取完整的账号数据（包括解密后的所有字段）
        const fullAccountData = await window.go.app.App.GetAccountByID(account.id)
        console.log(`[duplicateAccount] 获取到完整账号数据:`, fullAccountData)

        // 创建账号副本，移除ID和时间戳字段，使用完整的解密数据
        const duplicatedAccount = {
          title: fullAccountData.title + ' - 副本',
          username: fullAccountData.username || '',
          password: fullAccountData.password || '',
          url: fullAccountData.url || '',
          notes: fullAccountData.notes || '',
          icon: fullAccountData.icon || 'fa-globe',
          typeid: fullAccountData.typeid || currentTab.value?.id || '',
          group_id: fullAccountData.group_id || currentGroupId.value || '', // 确保传递分组信息
          input_method: fullAccountData.input_method || 1, // 复制输入方式
          is_favorite: false, // 副本默认不是收藏
          use_count: 0 // 副本使用次数重置为0
        }

        console.log(`[duplicateAccount] 生成副本，分组ID: ${duplicatedAccount.group_id}, 类型ID: ${duplicatedAccount.typeid}`)

        selectedAccount.value = duplicatedAccount
        isNewAccount.value = true
        showAccountDetail.value = true
      } catch (error) {
        console.error('[duplicateAccount] 获取完整账号数据失败:', error)
        ElMessage.error(`获取账号数据失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 复制账号的指定字段
     * @param {Object} account 账号对象
     * @param {string} field 字段名称
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     */
    const copyAccountField = async (account, field) => {
      hideContextMenu()

      try {
        switch (field) {
          case 'password':
            // 使用标准复制密码函数
            await copyAccountPassword(account.id, account.title)
            break
          case 'username':
            // 使用后端安全复制方法，10秒后自动清理剪贴板
            await window.go.app.App.CopyAccountUsername(account.id)
            ElMessage.success(t('success.usernameCopied'))
            break
          case 'url':
            // URL不是敏感信息，可以直接复制
            await navigator.clipboard.writeText(account.url || '')
            ElMessage.success(t('success.urlCopied'))
            break
          case 'title':
            // 标题不是敏感信息，可以直接复制
            await navigator.clipboard.writeText(account.title || '')
            ElMessage.success(t('success.titleCopied'))
            break
          case 'notes':
            // 使用后端安全复制方法，10秒后自动清理剪贴板
            await window.go.app.App.CopyAccountNotes(account.id)
            ElMessage.success(t('success.notesCopied'))
            break
          default:
            ElMessage.warning(t('warning.unknownFieldType'))
            return
        }
      } catch (error) {
        console.error('复制失败:', error)
        ElMessage.error(t('error.copyFailed'))
      }
    }

    /**
     * 复制账号用户名和密码
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 20251003
     */
    const copyAccountUsernameAndPassword = async (account) => {
      hideContextMenu()

      try {
        // 使用后端安全复制方法，10秒后自动清理剪贴板
        await window.go.app.App.CopyAccountUsernameAndPassword(account.id)
        ElMessage.success(t('success.usernameAndPasswordCopied'))
      } catch (error) {
        console.error('复制失败:', error)
        ElMessage.error(t('error.copyFailed'))
      }
    }

    /**
     * 显示密码对话框
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-18
     * @description 显示密码分段显示对话框
     */
    const showPasswordDialog = (account) => {
      console.log('[MainView] 显示密码对话框:', account)
      passwordDisplayAccount.value = account
      showPasswordDisplayDialog.value = true
    }

    /**
     * 处理更改分组
     * @param {Object} account 账号对象
     * @author 20251005 陈凤庆 新增更改分组功能
     */
    const handleChangeGroup = (account) => {
      console.log('[MainView] 处理更改分组:', account)
      changeGroupAccount.value = account
      showChangeGroupDialog.value = true
    }

    /**
     * 处理更改分组确认
     * @param {Object} result 更改结果 {accountId, groupId, typeId}
     * @author 20251005 陈凤庆 新增更改分组确认处理
     */
    const handleChangeGroupConfirm = async (result) => {
      try {
        console.log('[MainView] 确认更改分组:', result)

        // 调用后端API更新账号分组
        await apiService.updateAccountGroup(result.accountId, result.typeId)

        ElMessage.success('分组更改成功')

        // 刷新账号列表
        if (currentTabId.value) {
          // 如果有选中的标签，刷新该标签下的账号列表
          console.log(`[handleChangeGroupConfirm] 刷新标签${currentTabId.value}的账号列表`)
          const accounts = await apiService.getAccountsByTab(currentTabId.value)
          appStore.setPasswords(accounts)
        } else {
          // 如果没有选中标签，刷新整个分组的账号列表
          console.log(`[handleChangeGroupConfirm] 刷新分组${currentGroupId.value}的账号列表`)
          await loadPasswordsByGroup(currentGroupId.value)
        }

        // TODO: 定位到下一个账号（需要实现账号定位逻辑）

      } catch (error) {
        console.error('[MainView] 更改分组失败:', error)
        ElMessage.error(`更改分组失败: ${error.message || '未知错误'}`)
      }
    }

    /**
     * 删除账号（带确认）
     * @param {Object} account 账号对象
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     */
    const deleteAccountWithConfirm = async (account) => {
      hideContextMenu()

      try {
        await ElMessageBox.confirm(
          `确定要删除账号"${account.title}"吗？此操作不可撤销。`,
          '删除确认',
          {
            confirmButtonText: '确定删除',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await handleDeleteAccount(account.id)
      } catch {
        // 用户取消操作
      }
    }

    /**
     * 设置用户活动跟踪
     * @author 陈凤庆
     * @date 20251004
     * @description 监听用户交互事件，定期更新活动状态
     */
    const setupUserActivityTracking = () => {
      let lastActivityTime = Date.now()
      let activityTimer = null

      // 节流函数，避免频繁调用API
      const throttledUpdateActivity = (() => {
        let lastCall = 0
        const delay = 30000 // 30秒节流

        return async () => {
          const now = Date.now()
          if (now - lastCall >= delay) {
            lastCall = now
            try {
              await apiService.updateUserActivity()
              console.log('[活动跟踪] 用户活动已更新')
            } catch (error) {
              console.error('[活动跟踪] 更新用户活动失败:', error)
            }
          }
        }
      })()

      // 用户活动事件列表
      const activityEvents = [
        'mousedown', 'mousemove', 'keydown', 'scroll',
        'click', 'touchstart', 'touchmove'
      ]

      // 活动事件处理函数
      const handleActivity = () => {
        lastActivityTime = Date.now()
        throttledUpdateActivity()
      }

      // 添加事件监听器
      activityEvents.forEach(event => {
        document.addEventListener(event, handleActivity, { passive: true })
      })

      // 定期检查活动状态（每分钟检查一次）
      activityTimer = setInterval(() => {
        const now = Date.now()
        const timeSinceLastActivity = now - lastActivityTime

        // 如果超过5分钟没有活动，停止更新
        if (timeSinceLastActivity < 5 * 60 * 1000) {
          throttledUpdateActivity()
        }
      }, 60000) // 每分钟检查一次

      // 窗口焦点事件
      window.addEventListener('focus', () => {
        console.log('[活动跟踪] 窗口获得焦点')
        console.log('[活动跟踪] 调用堆栈:', new Error().stack)
        handleActivity()
      })

      // 清理函数（组件卸载时调用）
      const cleanup = () => {
        activityEvents.forEach(event => {
          document.removeEventListener(event, handleActivity)
        })
        if (activityTimer) {
          clearInterval(activityTimer)
        }
      }

      // 返回清理函数
      return cleanup
    }

    return {
      // 组件引用
      mainContentAreaRef,

      // 响应式数据
      appTitle,
      showSearchBar,
      searchKeyword,
      totalUseCount,
      usageDays,
      showAccountDetail,
      selectedAccount,
      isNewAccount,
      searchResults,
      showContextMenu,
      contextMenuPosition,
      contextMenuAccount,
      showGroupContextMenu,
      groupContextMenuPosition,
      contextMenuGroup,
      showTabContextMenu,
      tabContextMenuPosition,
      contextMenuTab,
      showPasswordDisplayDialog,
      passwordDisplayAccount,
      showChangeGroupDialog,
      changeGroupAccount,

      // 计算属性
      groups,
      currentGroupId,
      tabs,
      currentTabId,
      currentTab,
      filteredAccounts,
      isSearching,

      // 方法
      loadTabsByGroup,
      loadPasswordsByGroup,
      selectGroup,
      selectTab,
      toggleSearch,
      handleRefreshData,
      handleSearch,
      clearSearch,
      showGroupMenu,
      hideGroupContextMenu,
      handleRenameGroup,
      handleDeleteGroup,
      handleMoveGroupLeft,
      handleMoveGroupRight,
      showTabMenu,
      hideTabContextMenu,
      handleRenameTab,
      handleDeleteTab,
      handleMoveTabUp,
      handleMoveTabDown,
      handleTabDragReorder,
      handleCreateTabAfter,
      handleGroupDragReorder,
      handleCreateGroupFromMenu,
      createGroup,
      createTab,
      renameTab,
      deleteTab,
      createAccount,
      showAccountDetailDialog,
      inputUsernameAndPassword,
      inputUsername,
      inputPassword,
      handleSaveAccount,
      handleDeleteAccount,
      showAccountContextMenu,
      hideContextMenu,
      openAccountUrl,
      editAccount,
      duplicateAccount,
      copyAccountField,
      copyAccountUsernameAndPassword,
      showPasswordDialog,
      handleChangeGroup,
      handleChangeGroupConfirm,
      deleteAccountWithConfirm,
      scrollToAccount,

      // 统计和活动跟踪函数
      updateAccountUsageAndRefresh,
      refreshTotalUseCount,
      refreshUsageDays,
      setupUserActivityTracking
    }
  }
}
</script>

<style scoped>
.main-container {
  width: 100vw;
  height: 100vh;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  position: relative;
  background: #ffffff;
  display: flex;
  flex-direction: column;
  min-width: 300px;
  min-height: 400px;
}

/* 20251002 陈凤庆 账号行高亮样式，用于保存后定位 */
:deep(.highlight-account) {
  background-color: rgba(64, 158, 255, 0.1) !important;
  border: 1px solid #409eff !important;
  border-radius: 4px;
  transition: all 0.3s ease;
}

:deep(.highlight-account:hover) {
  background-color: rgba(64, 158, 255, 0.15) !important;
}
</style>