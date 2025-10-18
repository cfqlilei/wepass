<template>
  <div
    class="account-row"
    :data-account-id="account.id"
    @dblclick="handleDoubleClick"
    @contextmenu.prevent="handleContextMenu"
  >
    <div class="account-info">
      <div class="account-title-group">
        <div class="account-title-line">
          <el-icon class="account-icon">
            <component :is="iconComponent" />
          </el-icon>
          <span class="title-text">{{ account.title }}</span>
        </div>
        <div class="account-username-line">
          <span class="username-text">{{ account.masked_username }}</span>
        </div>
      </div>
    </div>

    <div class="account-actions">
      <el-tooltip :content="$t('main.inputUsername')" placement="top" :show-after="500">
        <el-button
          type="text"
          :icon="User"
          size="small"
          @click.stop="handleInputUsername"
        />
      </el-tooltip>
      <el-tooltip :content="$t('main.inputPassword')" placement="top" :show-after="500">
        <el-button
          type="text"
          :icon="Lock"
          size="small"
          @click.stop="handleInputPassword"
        />
      </el-tooltip>
      <el-tooltip :content="$t('main.inputUsernameAndPassword')" placement="top" :show-after="500">
        <el-button
          type="text"
          :icon="Key"
          size="small"
          @click.stop="handleInputUsernameAndPassword"
        />
      </el-tooltip>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { Key, User, Lock, Document, FolderOpened, Folder } from '@element-plus/icons-vue'

/**
 * 账号行组件
 * @author 陈凤庆
 * @date 20251002
 * @modify 20251002 陈凤庆 密码行改为账号行，保持命名一致性
 * @description 单个账号项的显示和交互
 */
export default {
  name: 'AccountRow',
  props: {
    account: {
      type: Object,
      required: true
    }
  },
  emits: [
    'dblclick',
    'contextmenu',
    'input-username-and-password',
    'input-username',
    'input-password'
  ],
  setup(props, { emit }) {
    /**
     * 获取图标组件
     */
    const iconComponent = computed(() => {
      const iconMap = {
        'fa-globe': Document,
        'fa-mobile': Document,
        'fa-wifi': Document,
        'fa-credit-card': Document,
        'fa-github': Document,
        'fa-envelope': Document,
        'fa-folder-open': FolderOpened,
        'fa-folder': Folder
      }
      return iconMap[props.account.icon] || Document
    })

    /**
     * 处理双击事件
     */
    const handleDoubleClick = () => {
      emit('dblclick', props.account)
    }

    /**
     * 处理右键菜单事件
     */
    const handleContextMenu = (event) => {
      emit('contextmenu', event, props.account)
    }

    /**
     * 处理输入用户名和密码
     */
    const handleInputUsernameAndPassword = () => {
      emit('input-username-and-password', props.account)
    }

    /**
     * 处理输入用户名
     */
    const handleInputUsername = () => {
      emit('input-username', props.account)
    }

    /**
     * 处理输入密码
     */
    const handleInputPassword = () => {
      emit('input-password', props.account)
    }

    return {
      Key,
      User,
      Lock,
      iconComponent,
      handleDoubleClick,
      handleContextMenu,
      handleInputUsernameAndPassword,
      handleInputUsername,
      handleInputPassword
    }
  }
}
</script>

<style scoped>
.account-row {
  display: flex;
  border-bottom: 1px solid #f1f1f1;
  cursor: pointer;
}

.account-row:hover {
  background-color: rgba(0, 123, 255, 0.02);
}

/* 20251003 陈凤庆 添加账号高亮样式，用于新建/编辑后的定位提示 */
.account-row.highlight-account {
  background-color: rgba(0, 123, 255, 0.1);
  border-left: 3px solid #007bff;
  animation: highlight-fade 2s ease-out;
}

@keyframes highlight-fade {
  0% {
    background-color: rgba(0, 123, 255, 0.2);
  }
  100% {
    background-color: rgba(0, 123, 255, 0.1);
  }
}

.account-info {
  flex: 1;
  padding: 10px 8px;
  display: flex;
  align-items: center;
}

.account-title-group {
  display: flex;
  flex-direction: column; /* 使内部元素垂直堆叠 */
  justify-content: center; /* 垂直居中 */
  height: 100%; /* 确保占据父容器高度 */
}

.account-title-line {
  display: flex;
  align-items: center;
  gap: 8px;
}

.account-username-line {
  margin-left: 24px; /* 根据图标宽度和gap调整，使脱敏账号与标题对齐 */
  margin-top: 2px; /* 标题和脱敏账号之间的间距 */
}

.account-icon {
  font-size: 16px;
  color: #6c757d;
}

.title-text {
  font-weight: 500;
  color: #212529;
}

.username-text {
  color: #6c757d;
  font-size: 12px;
}

.account-actions {
  width: 120px;
  padding: 10px 8px;
  display: flex;
  gap: 6px;
  justify-content: center;
  align-items: center;
}

.account-actions .el-button {
  width: 24px !important;
  height: 24px !important;
  padding: 0 !important;
  color: #6c757d !important;
}

.account-actions .el-button:hover {
  background-color: rgba(0, 0, 0, 0.05) !important;
  color: #007bff !important;
}
</style>
