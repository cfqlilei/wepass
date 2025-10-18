<template>
  <div
    v-if="visible"
    class="tab-context-menu"
    :style="{ left: position.x + 'px', top: position.y + 'px' }"
    @click.stop
  >
    <div class="menu-item" @click="handleRename">
      <el-icon><Edit /></el-icon>
      <span>{{ $t('tabContextMenu.rename') }}</span>
    </div>
    <div class="menu-item" @click="handleDelete">
      <el-icon><Delete /></el-icon>
      <span>{{ $t('tabContextMenu.deleteTab') }}</span>
    </div>
    <div class="menu-divider"></div>
    <div class="menu-item" @click="handleCreateAfter">
      <el-icon><Plus /></el-icon>
      <span>{{ $t('tabContextMenu.newTab') }}</span>
    </div>
    <div class="menu-divider"></div>
    <div class="menu-item" @click="handleMoveUp">
      <el-icon><ArrowUp /></el-icon>
      <span>{{ $t('tabContextMenu.moveUp') }}</span>
      <span class="shortcut">Shift+↑</span>
    </div>
    <div class="menu-item" @click="handleMoveDown">
      <el-icon><ArrowDown /></el-icon>
      <span>{{ $t('tabContextMenu.moveDown') }}</span>
      <span class="shortcut">Shift+↓</span>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { Edit, Delete, ArrowUp, ArrowDown, Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'

/**
 * 标签右键菜单组件
 * @author 陈凤庆
 * @date 20251002
 * @description 标签右键菜单，包含重命名、删除、上移、下移功能
 */
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  },
  tab: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'rename', 'delete', 'move-up', 'move-down', 'create-after'])

const { t } = useI18n()

/**
 * 处理重命名
 */
const handleRename = async () => {
  if (!props.tab) return
  
  try {
    const { value: newName } = await ElMessageBox.prompt(
      t('tabContextMenu.promptNewName'),
      t('tabContextMenu.renameTab'),
      {
        confirmButtonText: t('tabContextMenu.confirm'),
        cancelButtonText: t('tabContextMenu.cancel'),
        inputPlaceholder: t('tabContextMenu.tabName'),
        inputValue: props.tab.name,
        inputValidator: (value) => {
          if (!value || value.trim() === '') {
            return t('tabContextMenu.tabNameCannotBeEmpty')
          }
          if (value.trim() === props.tab.name.trim()) {
            return t('tabContextMenu.tabNameNotChanged')
          }
          return true
        }
      }
    )

    if (newName && newName.trim() !== props.tab.name) {
      emit('rename', props.tab.id, newName.trim())
    }
  } catch (error) {
    // 用户取消操作
    console.log(t('tabContextMenu.userCanceledRename'))
  }
  
  emit('close')
}

/**
 * 处理删除标签
 */
const handleDelete = async () => {
  if (!props.tab) return

  try {
    await ElMessageBox.confirm(
      t('tabContextMenu.confirmDeleteTab', { tabName: props.tab.name }),
      t('tabContextMenu.deleteTab'),
      {
        confirmButtonText: t('tabContextMenu.confirmDelete'),
        cancelButtonText: t('tabContextMenu.cancel'),
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    emit('delete', props.tab.id)
  } catch (error) {
    // 用户取消操作
    console.log(t('tabContextMenu.userCanceledDelete'))
  }

  emit('close')
}

/**
 * 处理上移
 */
const handleMoveUp = () => {
  if (!props.tab) return
  
  emit('move-up', props.tab.id)
  emit('close')
}

/**
 * 处理下移
 */
const handleMoveDown = () => {
  if (!props.tab) return

  emit('move-down', props.tab.id)
  emit('close')
}

/**
 * 处理在当前标签后新建标签
 */
const handleCreateAfter = async () => {
  if (!props.tab) return

  try {
    const { value: newName } = await ElMessageBox.prompt(
      t('tabContextMenu.promptNewTabName'),
      t('tabContextMenu.newTab'),
      {
        confirmButtonText: t('tabContextMenu.confirm'),
        cancelButtonText: t('tabContextMenu.cancel'),
        inputPlaceholder: t('tabContextMenu.tabName'),
        inputValidator: (value) => {
          if (!value || value.trim() === '') {
            return t('tabContextMenu.tabNameCannotBeEmpty')
          }
          return true
        }
      }
    )

    if (newName && newName.trim()) {
      emit('create-after', props.tab.id, newName.trim())
    }
  } catch (error) {
    // 用户取消操作
    console.log(t('tabContextMenu.userCanceledNewTab'))
  }

  emit('close')
}

</script>

<style scoped>
.tab-context-menu {
  position: fixed;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 2000;
  min-width: 120px;
  padding: 4px 0;
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
</style>
