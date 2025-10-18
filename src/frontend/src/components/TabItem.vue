<template>
  <div
    class="tab-item"
    :class="{ active: active, dragging: isDragging }"
    draggable="true"
    @click="handleClick"
    @contextmenu.prevent="handleContextMenu"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
    @dragover.prevent="handleDragOver"
    @drop.prevent="handleDrop"
  >
    <el-icon class="tab-icon">
      <component :is="iconComponent" />
    </el-icon>
    <span>{{ tab.name }}</span>
  </div>
</template>

<script>
import { computed, ref } from 'vue'
import { FolderOpened, Folder, Document } from '@element-plus/icons-vue'

/**
 * 标签项组件
 * @author 陈凤庆
 * @date 20251001
 * @modify 20251002 陈凤庆 添加拖拽排序功能
 * @description 单个标签项的显示和交互，支持拖拽排序
 */
export default {
  name: 'TabItem',
  props: {
    tab: {
      type: Object,
      required: true
    },
    active: {
      type: Boolean,
      default: false
    }
  },
  emits: ['click', 'contextmenu', 'drag-reorder'],
  setup(props, { emit }) {
    // 拖拽状态
    const isDragging = ref(false)
    /**
     * 获取图标组件
     */
    const iconComponent = computed(() => {
      const iconMap = {
        'fa-folder-open': FolderOpened,
        'fa-folder': Folder,
        'fa-globe': Document,
        'fa-mobile': Document,
        'fa-wifi': Document,
        'fa-credit-card': Document,
        'fa-github': Document,
        'fa-envelope': Document
      }
      return iconMap[props.tab.icon] || Document
    })

    /**
     * 处理点击事件
     */
    const handleClick = () => {
      emit('click', props.tab.id)
    }

    /**
     * 处理右键菜单事件
     */
    const handleContextMenu = (event) => {
      emit('contextmenu', event, props.tab)
    }

    /**
     * 处理拖拽开始
     */
    const handleDragStart = (event) => {
      isDragging.value = true
      event.dataTransfer.setData('text/plain', JSON.stringify({
        id: props.tab.id,
        name: props.tab.name,
        sortOrder: props.tab.sort_order
      }))
      event.dataTransfer.effectAllowed = 'move'
      console.log('[TabItem] 开始拖拽标签:', props.tab.name)
    }

    /**
     * 处理拖拽结束
     */
    const handleDragEnd = () => {
      isDragging.value = false
      console.log('[TabItem] 拖拽结束')
    }

    /**
     * 处理拖拽悬停
     */
    const handleDragOver = (event) => {
      event.preventDefault()
      event.dataTransfer.dropEffect = 'move'
    }

    /**
     * 处理拖拽放置
     */
    const handleDrop = (event) => {
      event.preventDefault()

      try {
        const dragData = JSON.parse(event.dataTransfer.getData('text/plain'))
        const sourceTabId = dragData.id
        const targetTabId = props.tab.id

        if (sourceTabId !== targetTabId) {
          console.log('[TabItem] 拖拽排序:', dragData.name, '->', props.tab.name)
          emit('drag-reorder', sourceTabId, targetTabId)
        }
      } catch (error) {
        console.error('[TabItem] 拖拽数据解析失败:', error)
      }
    }

    return {
      isDragging,
      iconComponent,
      handleClick,
      handleContextMenu,
      handleDragStart,
      handleDragEnd,
      handleDragOver,
      handleDrop
    }
  }
}
</script>

<style scoped>
.tab-item {
  padding: 10px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  color: #212529;
  cursor: pointer;
  border-left: 3px solid transparent;
  font-size: 13px;
}

.tab-item.active {
  background-color: #ffffff;
  border-left-color: #007bff;
  font-weight: 500;
}

.tab-item:hover {
  background-color: rgba(0, 123, 255, 0.05);
}

.tab-item.dragging {
  opacity: 0.5;
  transform: rotate(2deg);
}

.tab-item[draggable="true"] {
  cursor: move;
}

.tab-item[draggable="true"]:hover {
  background-color: rgba(0, 123, 255, 0.1);
}

.tab-icon {
  font-size: 16px;
  color: #6c757d;
}
</style>

