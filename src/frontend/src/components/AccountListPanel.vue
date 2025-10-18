<template>
  <div class="account-list-panel">
    <div class="list-header">
      <div class="list-title">{{ currentTab?.name || $t('main.allAccounts') }}</div>
      <el-button
        type="primary"
        :icon="Plus"
        size="small"
        @click="handleCreateAccount"
      >
        {{ $t('account.addAccount') }}
      </el-button>
    </div>

    <!-- 账号项表格 -->
    <div class="account-table">
      <div class="table-body-wrapper">
        <div class="table-body">
          <!-- 账号列表 -->
          <AccountRow
            v-for="account in filteredAccounts"
            :key="account.id"
            :account="account"
            @dblclick="handleShowAccountDetail"
            @contextmenu="handleShowAccountContextMenu"
            @input-username-and-password="handleInputUsernameAndPassword"
            @input-username="handleInputUsername"
            @input-password="handleInputPassword"
          />

          <!-- 空状态提示 -->
          <div v-if="filteredAccounts.length === 0" class="empty-accounts-hint">
            <el-icon><Key /></el-icon>
            <p>{{ $t('account.noAccounts') }}</p>
            <p class="hint-text">{{ $t('account.noAccountsHint') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Plus, Key } from '@element-plus/icons-vue'
import AccountRow from './AccountRow.vue'

/**
 * 账号列表面板组件
 * @author 陈凤庆
 * @date 20251001
 * @modify 20251002 陈凤庆 密码列表改为账号列表，保持命名一致性
 * @modify 20251003 陈凤庆 取消分页功能，采用全部显示的方式
 * @description 显示账号列表头部、账号项列表，处理账号项的创建、详情、操作和右键菜单
 */
export default {
  name: 'AccountListPanel',
  components: {
    AccountRow
  },
  props: {
    currentTab: {
      type: Object,
      default: null
    },
    filteredAccounts: {
      type: Array,
      default: () => []
    }
  },
  emits: [
    'create-account',
    'show-account-detail',
    'show-account-context-menu',
    'input-username-and-password',
    'input-username',
    'input-password'
  ],
  setup(props, { emit }) {
    /**
     * 处理创建账号
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     */
    const handleCreateAccount = () => {
      emit('create-account')
    }

    /**
     * 处理显示账号详情
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
     */
    const handleShowAccountDetail = (account) => {
      emit('show-account-detail', account)
    }

    /**
     * 处理显示账号右键菜单
     * @modify 20251002 陈凤庆 账号改为账号项，保持命名一致性
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
     * 导航到指定账号（无分页模式下无需实际操作）
     * @param {string} accountId 账号ID
     * @author 20251003 陈凤庆 取消分页后，此方法保留接口兼容性但无实际操作
     */
    const navigateToAccountPage = (accountId) => {
      if (!accountId) {
        console.warn('[AccountListPanel] navigateToAccountPage:', this.$t('account.accountIdEmpty'))
        return false
      }

      const allAccounts = props.filteredAccounts || []
      const targetIndex = allAccounts.findIndex(account => account.id === accountId)

      if (targetIndex === -1) {
        console.warn(`[AccountListPanel] navigateToAccountPage: ${this.$t('account.accountNotFound', { accountId })}`)
        return false
      }

      // 无分页模式下，账号已经在列表中显示，无需切换页面
      console.log(`[AccountListPanel] navigateToAccountPage: ${this.$t('account.accountAlreadyDisplayed', { accountId })}`)
      return true
    }

    return {
      Plus,
      Key,
      handleCreateAccount,
      handleShowAccountDetail,
      handleShowAccountContextMenu,
      handleInputUsernameAndPassword,
      handleInputUsername,
      handleInputPassword,
      navigateToAccountPage
    }
  }
}
</script>

<style scoped>
.account-list-panel {
  flex: 1;
  padding: 16px 0px 16px 16px; /* 右边距设置为 0px */
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-width: 0;
  /* 20251003 陈凤庆 为滚动条提供相对定位上下文 */
  position: relative;
}

.list-header {
  margin-bottom: 16px;
  margin-left: 16px; /* 左边距设置为 16px */
  margin-right: 16px; /* 右边距设置为 16px */
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-title {
  font-size: 16px;
  font-weight: bold;
  color: #212529;
}

.account-table {
  width: 100%;
  font-size: 13px;
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-body-wrapper {
  flex: 1;
  position: relative;
  overflow: hidden;
  /* 20251003 陈凤庆 为浮动滚动条提供定位上下文 */
}

.table-body {
  height: 100%;
  overflow-y: auto;
  /* 20251003 陈凤庆 滚动条样式优化：始终为滚动条预留空间，避免布局动态变化 */
  scrollbar-width: thin; /* Firefox - 始终显示细滚动条 */
  scrollbar-color: rgba(0, 0, 0, 0.3) transparent; /* Firefox 滚动条颜色 */
  /* 20251003 陈凤庆 添加右侧padding为滚动条预留空间，避免布局变化 */
  padding-right: 3px; /* 为滚动条预留3px空间 */
  margin-right: 0;
}

/* 20251003 陈凤庆 WebKit浏览器滚动条样式 - 宽度设置为3px */
.table-body::-webkit-scrollbar {
  width: 2px; /* 滚动条宽度设置为3px */
  opacity: 0; /* 始终显示滚动条，避免布局变化 */
}

.table-body::-webkit-scrollbar-track {
  background: transparent;
}

.table-body::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 1px; /* 调整圆角以适应3px宽度 */
  transition: background 0.3s ease;
}

.table-body::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.5);
}

/* 20251003 陈凤庆 确保滚动条始终占用空间，避免布局动态变化 */
.table-body {
  box-sizing: border-box; /* 改为border-box，确保padding计算在内 */
}

.empty-accounts-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #909399;
  text-align: center;
  min-height: 300px;
}

.empty-accounts-hint .el-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-accounts-hint p {
  margin: 8px 0;
  font-size: 14px;
}

.empty-accounts-hint .hint-text {
  font-size: 12px;
  color: #c0c4cc;
}


</style>
