<template>
  <!-- 右键菜单遮罩层 -->
  <div
    v-if="show"
    class="context-menu-overlay"
    @click="handleClose"
  >
    <!-- 账号右键菜单 -->
    <div
      class="password-context-menu"
      :style="{ left: position.x + 'px', top: position.y + 'px' }"
      @click.stop
    >
      <div class="context-menu-item" @click="handleOpenUrl">
        <el-icon><Link /></el-icon>
        <span>{{ t('contextMenu.openUrl') }}</span>
      </div>
      <div class="context-menu-item" @click="handleShowPassword">
        <el-icon><View /></el-icon>
        <span>{{ t('contextMenu.showPassword') }}</span>
      </div>
      <div class="context-menu-item" @click="handleDuplicate">
        <el-icon><CopyDocument /></el-icon>
        <span>{{ t('contextMenu.duplicate') }}</span>
      </div>
      <div class="context-menu-divider"></div>
      <div class="context-menu-item" @click="handleView">
        <el-icon><View /></el-icon>
        <span>{{ t('contextMenu.view') }}</span>
      </div>
      <div class="context-menu-item" @click="handleEdit">
        <el-icon><Edit /></el-icon>
        <span>{{ t('contextMenu.edit') }}</span>
      </div>
      <div class="context-menu-item" @click="handleChangeGroup">
        <el-icon><FolderOpened /></el-icon>
        <span>{{ t('contextMenu.changeGroup') }}</span>
      </div>
      <div class="context-menu-divider"></div>
      <div class="context-menu-item" @click="handleCopyUsername">
        <el-icon><User /></el-icon>
        <span>{{ t('contextMenu.copyUsername') }}</span>
      </div>
      <div class="context-menu-item" @click="handleCopyPassword">
        <el-icon><Lock /></el-icon>
        <span>{{ t('contextMenu.copyPassword') }}</span>
      </div>
      <div class="context-menu-item" @click="handleCopyUsernameAndPassword">
        <el-icon><CopyDocument /></el-icon>
        <span>{{ t('contextMenu.copyUsernameAndPassword') }}</span>
      </div>
      <div class="context-menu-item" @click="handleCopyUrl">
        <el-icon><Link /></el-icon>
        <span>{{ t('contextMenu.copyUrl') }}</span>
      </div>
      <div class="context-menu-item" @click="handleCopyTitle">
        <el-icon><Document /></el-icon>
        <span>{{ t('contextMenu.copyTitle') }}</span>
      </div>
      <div class="context-menu-item" @click="handleCopyNotes">
        <el-icon><Memo /></el-icon>
        <span>{{ t('contextMenu.copyNotes') }}</span>
      </div>
      <div class="context-menu-divider"></div>
      <div class="context-menu-item danger" @click="handleDelete">
        <el-icon><Delete /></el-icon>
        <span>{{ t('contextMenu.delete') }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import { Key, Link, View, Edit, Lock, User, Document, Memo, Delete, CopyDocument, FolderOpened } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

/**
 * 账号右键菜单组件
 * @author 陈凤庆
 * @date 20251001
 * @description 封装账号右键菜单的 UI 和逻辑
 */
export default {
  name: 'PasswordContextMenu',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    position: {
      type: Object,
      default: () => ({ x: 0, y: 0 })
    },
    account: {
      type: Object,
      default: null
    }
  },
  emits: [
    'close',
    'open-url',
    'view',
    'edit',
    'change-group',
    'duplicate',
    'copy-password',
    'show-password',
    'copy-username',
    'copy-username-and-password',
    'copy-url',
    'copy-title',
    'copy-notes',
    'delete'
  ],
  setup(props, { emit }) {
    const { t } = useI18n()
    /**
     * 处理关闭菜单
     */
    const handleClose = () => {
      emit('close')
    }

    /**
     * 处理打开URL
     */
    const handleOpenUrl = () => {
      emit('open-url', props.account)
      handleClose()
    }

    /**
     * 处理查看
     */
    const handleView = () => {
      emit('view', props.account)
      handleClose()
    }

    /**
     * 处理编辑
     */
    const handleEdit = () => {
      emit('edit', props.account)
      handleClose()
    }

    /**
     * 处理更改分组
     */
    const handleChangeGroup = () => {
      emit('change-group', props.account)
      handleClose()
    }

    /**
     * 处理生成副本
     */
    const handleDuplicate = () => {
      emit('duplicate', props.account)
      handleClose()
    }

    /**
     * 处理复制密码
     */
    const handleCopyPassword = () => {
      emit('copy-password', props.account)
      handleClose()
    }

    /**
     * 处理显示密码
     */
    const handleShowPassword = () => {
      emit('show-password', props.account)
      handleClose()
    }

    /**
     * 处理复制用户名
     */
    const handleCopyUsername = () => {
      emit('copy-username', props.account)
      handleClose()
    }

    /**
     * 处理复制用户名和密码
     */
    const handleCopyUsernameAndPassword = () => {
      emit('copy-username-and-password', props.account)
      handleClose()
    }

    /**
     * 处理复制URL
     */
    const handleCopyUrl = () => {
      emit('copy-url', props.account)
      handleClose()
    }

    /**
     * 处理复制标题
     */
    const handleCopyTitle = () => {
      emit('copy-title', props.account)
      handleClose()
    }

    /**
     * 处理复制备注
     */
    const handleCopyNotes = () => {
      emit('copy-notes', props.account)
      handleClose()
    }

    /**
     * 处理删除
     */
    const handleDelete = () => {
      emit('delete', props.account)
      handleClose()
    }

    return {
      t,
      Key,
      Link,
      View,
      Edit,
      Lock,
      User,
      Document,
      Memo,
      Delete,
      CopyDocument,
      FolderOpened,
      handleClose,
      handleOpenUrl,
      handleView,
      handleEdit,
      handleChangeGroup,
      handleDuplicate,
      handleCopyPassword,
      handleShowPassword,
      handleCopyUsername,
      handleCopyUsernameAndPassword,
      handleCopyUrl,
      handleCopyTitle,
      handleCopyNotes,
      handleDelete
    }
  }
}
</script>

<style scoped>
.context-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000;
  pointer-events: auto;
}

.password-context-menu {
  position: fixed;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 4px 0;
  min-width: 160px;
  z-index: 1001;
  pointer-events: auto;
}

.context-menu-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
  transition: background-color 0.2s;
}

.context-menu-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.context-menu-item .el-icon {
  margin-right: 8px;
  font-size: 16px;
}

.context-menu-item.danger {
  color: #f56c6c;
}

.context-menu-item.danger:hover {
  background-color: #fef0f0;
  color: #f56c6c;
}

.context-menu-divider {
  height: 1px;
  background-color: #e4e7ed;
  margin: 4px 0;
}
</style>
