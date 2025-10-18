<template>
  <div
    class="group-tag"
    :class="{ active: active, dragging: isDragging }"
    draggable="true"
    @click="handleClick"
    @contextmenu.prevent="handleContextMenu"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
    @dragover.prevent="handleDragOver"
    @drop.prevent="handleDrop"
  >
    <el-icon class="group-icon">
      <component :is="iconComponent" />
    </el-icon>
    <span>{{ group.name }}</span>
  </div>
</template>

<script>
import { computed, ref } from 'vue'
import { FolderOpened, Folder, Document, Search } from '@element-plus/icons-vue'

/**
 * 分组标签组件
 * @author 陈凤庆
 * @date 20251001
 * @description 单个分组标签的显示和交互
 */
export default {
  name: 'GroupTag',
  props: {
    group: {
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
        'fa-envelope': Document,
        'fa-search': Search
      }
      return iconMap[props.group.icon] || Document
    })

    /**
     * 处理点击事件
     */
    const handleClick = () => {
      emit('click', props.group.id)
    }

    /**
     * 处理右键菜单事件
     */
    const handleContextMenu = (event) => {
      // 搜索结果分组不显示右键菜单
      if (props.group.isSearchResult) {
        return
      }
      console.log('[GroupTag] 右键菜单事件触发:', props.group.name, event)
      emit('contextmenu', event, props.group)
    }

    // 拖拽状态
    const isDragging = ref(false)

    /**
     * 处理拖拽开始
     */
    const handleDragStart = (event) => {
      console.log('[GroupTag] 开始拖拽分组:', props.group.name)
      isDragging.value = true
      event.dataTransfer.setData('text/plain', props.group.id)
      event.dataTransfer.effectAllowed = 'move'
    }

    /**
     * 处理拖拽结束
     */
    const handleDragEnd = () => {
      console.log('[GroupTag] 拖拽结束:', props.group.name)
      isDragging.value = false
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
      const sourceGroupId = event.dataTransfer.getData('text/plain')
      const targetGroupId = props.group.id

      if (sourceGroupId !== targetGroupId) {
        console.log('[GroupTag] 拖拽排序:', sourceGroupId, '->', targetGroupId)
        emit('drag-reorder', sourceGroupId, targetGroupId)
      }
    }

    return {
      iconComponent,
      handleClick,
      handleContextMenu,
      isDragging,
      handleDragStart,
      handleDragEnd,
      handleDragOver,
      handleDrop
    }
  }
}
</script>

<style scoped>
.group-tag {
  padding: 6px 12px;
  background-color: #e9ecef;
  border: 1px solid #dee2e6;
  border-bottom-color: transparent;
  border-radius: 4px 4px 0 0;
  font-size: 13px;
  font-weight: 500;
  color: #495057;
  cursor: pointer;
  white-space: nowrap;
  height: 35px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.group-tag.active {
  background-color: #ffffff;
  color: #007bff;
  border-color: #dee2e6 #dee2e6 #ffffff;
  border-top: 2px solid #007bff;
  margin-bottom: -1px;
}

.group-tag:hover:not(.active) {
  background-color: #dee2e6;
  color: #212529;
}

.group-icon {
  font-size: 16px;
}

/* 拖拽状态样式 */
.group-tag.dragging {
  opacity: 0.5;
  transform: rotate(5deg);
  cursor: grabbing;
}

.group-tag[draggable="true"] {
  cursor: grab;
}

.group-tag[draggable="true"]:active {
  cursor: grabbing;
}
</style>

