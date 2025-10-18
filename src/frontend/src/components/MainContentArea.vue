<template>
  <div class="main-content">
    <!-- 左侧页签栏（搜索状态下隐藏） -->
    <TabsSidebar
      v-if="showTabs"
      :tabs="tabs"
      :current-tab-id="currentTabId"
      :width="sidebarWidth"
      @select-tab="handleSelectTab"
      @create-tab="handleCreateTab"
      @show-tab-menu="handleShowTabMenu"
      @drag-reorder="handleTabDragReorder"
    />

    <!-- 可拖动的分隔栏（搜索状态下隐藏） -->
    <Resizer v-if="showTabs" @start-resize="handleStartResize" />

    <!-- 右侧账号列表 -->
    <AccountListPanel
      ref="accountListPanelRef"
      :current-tab="currentTab"
      :filtered-accounts="filteredAccounts"
      :class="{ 'full-width': !showTabs }"
      @create-account="handleCreateAccount"
      @show-account-detail="handleShowAccountDetail"
      @show-account-context-menu="handleShowAccountContextMenu"
      @input-username-and-password="handleInputUsernameAndPassword"
      @input-username="handleInputUsername"
      @input-password="handleInputPassword"
    />
  </div>
</template>

<script>
import { ref } from 'vue'
import TabsSidebar from './TabsSidebar.vue'
import Resizer from './Resizer.vue'
import AccountListPanel from './AccountListPanel.vue'

/**
 * 主内容区域组件
 * @author 陈凤庆
 * @date 20251001
 * @modify 20251002 陈凤庆 密码列表改为账号列表，保持命名一致性
 * @description 包含左侧标签栏、可拖动分隔符和右侧账号列表
 */
export default {
  name: 'MainContentArea',
  components: {
    TabsSidebar,
    Resizer,
    AccountListPanel
  },
  props: {
    tabs: {
      type: Array,
      default: () => []
    },
    currentTabId: {
      type: Number,
      default: null
    },
    currentTab: {
      type: Object,
      default: null
    },
    filteredAccounts: {
      type: Array,
      default: () => []
    },
    showTabs: {
      type: Boolean,
      default: true
    }
  },
  emits: [
    'select-tab',
    'create-tab',
    'show-tab-menu',
    'drag-reorder',
    'create-account',
    'show-account-detail',
    'show-account-context-menu',
    'input-username-and-password',
    'input-username',
    'input-password'
  ],
  setup(props, { emit, expose }) {
    const sidebarWidth = ref(120)
    const isResizing = ref(false)
    const accountListPanelRef = ref(null)

    /**
     * 处理选择标签
     */
    const handleSelectTab = (tabId) => {
      emit('select-tab', tabId)
    }

    /**
     * 处理创建标签
     */
    const handleCreateTab = () => {
      emit('create-tab')
    }

    /**
     * 处理显示标签菜单
     */
    const handleShowTabMenu = (event, tab) => {
      emit('show-tab-menu', event, tab)
    }

    /**
     * 处理创建账号
     */
    const handleCreateAccount = () => {
      emit('create-account')
    }

    /**
     * 处理显示账号详情
     */
    const handleShowAccountDetail = (account) => {
      emit('show-account-detail', account)
    }

    /**
     * 处理显示账号右键菜单
     */
    const handleShowAccountContextMenu = (event, account) => {
      emit('show-account-context-menu', event, account)
    }

    /**
     * 处理输入用户名和密码
     */
    const handleInputUsernameAndPassword = (password) => {
      emit('input-username-and-password', password)
    }

    /**
     * 处理输入用户名
     */
    const handleInputUsername = (password) => {
      emit('input-username', password)
    }

    /**
     * 处理输入密码
     */
    const handleInputPassword = (password) => {
      emit('input-password', password)
    }

    /**
     * 处理标签拖拽排序
     */
    const handleTabDragReorder = (sourceTabId, targetTabId) => {
      emit('drag-reorder', sourceTabId, targetTabId)
    }

    /**
     * 开始拖动调整侧边栏宽度
     * @param {MouseEvent} event 鼠标事件
     */
    const handleStartResize = (event) => {
      isResizing.value = true
      const startX = event.clientX
      const startWidth = sidebarWidth.value

      const handleMouseMove = (e) => {
        if (!isResizing.value) return

        const deltaX = e.clientX - startX
        const newWidth = startWidth + deltaX

        // 限制最小和最大宽度
        if (newWidth >= 80 && newWidth <= 300) {
          sidebarWidth.value = newWidth
        }
      }

      const handleMouseUp = () => {
        isResizing.value = false
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
        document.body.style.cursor = 'default'
        document.body.style.userSelect = 'auto'
      }

      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
      document.body.style.cursor = 'col-resize'
      document.body.style.userSelect = 'none'
    }

    // 暴露 AccountListPanel 的方法给父组件
    expose({
      navigateToAccountPage: (accountId) => {
        return accountListPanelRef.value?.navigateToAccountPage(accountId)
      }
    })

    return {
      sidebarWidth,
      isResizing,
      accountListPanelRef,
      handleSelectTab,
      handleCreateTab,
      handleShowTabMenu,
      handleTabDragReorder,
      handleCreateAccount,
      handleShowAccountDetail,
      handleShowAccountContextMenu,
      handleInputUsernameAndPassword,
      handleInputUsername,
      handleInputPassword,
      handleStartResize
    }
  }
}
</script>

<style scoped>
.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
  min-height: 0;
}

/* 搜索状态下账号列表全宽显示 */
:deep(.full-width) {
  flex: 1;
  width: 100%;
}
</style>

