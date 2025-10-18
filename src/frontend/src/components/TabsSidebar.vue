<template>
  <div class="tabs-sidebar" :style="{ width: width + 'px' }">
    <!-- 标签列表 -->
    <TabItem
      v-for="tab in tabs"
      :key="tab.id"
      :tab="tab"
      :active="currentTabId === tab.id"
      @click="handleSelectTab"
      @contextmenu="handleTabContextMenu"
      @drag-reorder="handleTabDragReorder"
    />

    <!-- 空状态提示 -->
    <div v-if="tabs.length === 0" class="empty-tabs-hint">
      <el-icon><Document /></el-icon>
      <p>{{ $t('tabsSidebar.emptyTabs') }}</p>
      <p class="hint-text">{{ $t('tabsSidebar.createTabHint') }}</p>
    </div>

    <!-- 新建标签按钮 -->
    <div class="add-tab-btn" @click="handleCreateTab">
      <el-icon><Plus /></el-icon>
      <span>{{ $t('tabsSidebar.newTab') }}</span>
    </div>
  </div>
</template>

<script setup>
import { Plus, Document } from '@element-plus/icons-vue'
import TabItem from './TabItem.vue'
import { useI18n } from 'vue-i18n'

/**
 * 标签侧边栏组件
 * @author 陈凤庆
 * @date 20251001
 * @description 显示标签列表，处理标签选择和创建
 */

const props = defineProps({
  tabs: {
    type: Array,
    default: () => []
  },
  currentTabId: {
    type: Number,
    default: null
  },
  width: {
    type: Number,
    default: 120
  }
})

const emit = defineEmits(['select-tab', 'create-tab', 'show-tab-menu', 'drag-reorder'])

const { t } = useI18n()

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
 * 处理标签右键菜单
 */
const handleTabContextMenu = (event, tab) => {
  emit('show-tab-menu', event, tab)
}

/**
 * 处理标签拖拽排序
 */
const handleTabDragReorder = (sourceTabId, targetTabId) => {
  emit('drag-reorder', sourceTabId, targetTabId)
}
</script>

<style scoped>
.tabs-sidebar {
  min-width: 80px;
  max-width: 300px;
  border-right: 1px solid #dee2e6;
  background-color: #f8f9fa;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.add-tab-btn {
  margin: 8px;
  padding: 8px;
  text-align: center;
  border: 1px dashed #dee2e6;
  border-radius: 4px;
  color: #6c757d;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.add-tab-btn:hover {
  border-color: #007bff;
  color: #007bff;
}

.empty-tabs-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #909399;
  text-align: center;
  flex: 1;
  min-height: 200px;
}

.empty-tabs-hint .el-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-tabs-hint p {
  margin: 8px 0;
  font-size: 14px;
}

.empty-tabs-hint .hint-text {
  font-size: 12px;
  color: #c0c4cc;
}
</style>
