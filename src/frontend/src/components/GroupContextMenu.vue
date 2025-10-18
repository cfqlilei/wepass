<template>
  <div
    v-if="visible"
    class="group-context-menu"
    :style="{ left: position.x + 'px', top: position.y + 'px' }"
    @click.stop
  >
    <div class="menu-item" @click="handleRename">
      <el-icon><Edit /></el-icon>
      <span>{{ t('group.rename') }}</span>
    </div>
    <div class="menu-item" @click="handleDelete">
      <el-icon><Delete /></el-icon>
      <span>{{ t('group.deleteGroup') }}</span>
    </div>
    <div class="menu-divider"></div>
    <div class="menu-item" @click="handleCreateGroup">
      <el-icon><Plus /></el-icon>
      <span>{{ t('group.createGroup') }}</span>
    </div>
    <div class="menu-divider"></div>
    <div class="menu-item" @click="handleMoveLeft">
      <el-icon><ArrowLeft /></el-icon>
      <span>{{ t('group.moveLeft') }}</span>
      <span class="shortcut">Shift+←</span>
    </div>
    <div class="menu-item" @click="handleMoveRight">
      <el-icon><ArrowRight /></el-icon>
      <span>{{ t('group.moveRight') }}</span>
      <span class="shortcut">Shift+→</span>
    </div>
  </div>
</template>

<script>
import { ref, nextTick } from 'vue'
import { Edit, Delete, ArrowLeft, ArrowRight, Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

/**
 * 分组右键菜单组件
 * @author 陈凤庆
 * @date 20251002
 * @description 分组右键菜单，包含重命名、删除、左移、右移功能
 */
export default {
  name: 'GroupContextMenu',
  components: {
    Edit,
    Delete,
    ArrowLeft,
    ArrowRight,
    Plus
  },
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    position: {
      type: Object,
      default: () => ({ x: 0, y: 0 })
    },
    group: {
      type: Object,
      default: null
    }
  },
  emits: ['close', 'rename', 'delete', 'move-left', 'move-right', 'create-group'],
  setup(props, { emit }) {
    const { t } = useI18n()
    /**
     * 处理重命名
     */
    const handleRename = async () => {
      if (!props.group) return
      
      try {
        const { value: newName } = await ElMessageBox.prompt(
          t('group.newGroupName'),
          t('group.renameGroup'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            inputPlaceholder: t('group.newGroupPlaceholder'),
            inputValue: props.group.name,
            inputValidator: (value) => {
              if (!value || value.trim() === '') {
                return t('group.groupNameCannotBeEmpty')
              }
              if (value.trim() === props.group.name) {
                return t('group.groupNameNotChanged')
              }
              return true
            }
          }
        )

        if (newName && newName.trim() !== props.group.name) {
          emit('rename', props.group.id, newName.trim())
        }
      } catch (error) {
        // 用户取消操作
        console.log('用户取消重命名操作')
      }
      
      emit('close')
    }

    /**
     * 处理删除分组
     */
    const handleDelete = async () => {
      if (!props.group) return

      // 检查是否是默认分组
      if (props.group.name === t('group.defaultGroupName')) {
        ElMessage.warning(t('group.defaultGroupCannotBeDeleted'))
        emit('close')
        return
      }

      try {
        await ElMessageBox.confirm(
          t('group.confirmDeleteGroup', { groupName: props.group.name }),
          t('group.deleteGroupTitle'),
          {
            confirmButtonText: t('group.confirmDelete'),
            cancelButtonText: t('common.cancel'),
            type: 'warning',
            dangerouslyUseHTMLString: false
          }
        )

        emit('delete', props.group.id)
      } catch (error) {
        // 用户取消操作
        console.log('用户取消删除操作')
      }

      emit('close')
    }

    /**
     * 处理左移
     */
    const handleMoveLeft = () => {
      if (!props.group) return
      
      emit('move-left', props.group.id)
      emit('close')
    }

    /**
     * 处理右移
     */
    const handleMoveRight = () => {
      if (!props.group) return

      emit('move-right', props.group.id)
      emit('close')
    }

    /**
     * 处理新建分组
     */
    const handleCreateGroup = async () => {
      try {
        const { value: groupName } = await ElMessageBox.prompt(
          t('group.newGroupName'),
          t('group.createGroup'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            inputPlaceholder: t('group.newGroupPlaceholder'),
            inputValidator: (value) => {
              if (!value || value.trim() === '') {
                return t('group.groupNameCannotBeEmpty')
              }
              return true
            }
          }
        )

        if (groupName && groupName.trim()) {
          emit('create-group', groupName.trim())
        }
      } catch (error) {
        // 用户取消操作
        console.log('用户取消新建分组操作')
      }

      emit('close')
    }

    return {
      t,
      handleRename,
      handleDelete,
      handleMoveLeft,
      handleMoveRight,
      handleCreateGroup
    }
  }
}
</script>

<style scoped>
.group-context-menu {
  position: fixed;
  z-index: 9999;
  background: #ffffff;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 4px 0;
  min-width: 140px;
  font-size: 14px;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  cursor: pointer;
  color: #606266;
  transition: background-color 0.3s;
  gap: 8px;
}

.menu-item:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.menu-item .el-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.menu-item span:first-of-type {
  flex: 1;
}

.shortcut {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.menu-divider {
  height: 1px;
  background-color: #e4e7ed;
  margin: 4px 0;
}

/* 删除项特殊样式 */
.menu-item:nth-child(2):hover {
  background-color: #fef0f0;
  color: #f56c6c;
}
</style>
