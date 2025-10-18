<template>
  <div class="title-bar">
    <div class="app-title">
      <el-icon class="title-icon"><Key /></el-icon>
      <span>{{ displayTitle }}</span>
    </div>
    <div class="title-bar-actions">
      <el-button
        type="text"
        :icon="Search"
        class="action-btn"
        :title="$t('common.search')"
        @click="handleToggleSearch"
      />
      <MoreMenu @refresh-data="handleRefreshData" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Key, Search } from '@element-plus/icons-vue'
import MoreMenu from './MoreMenu.vue'

/**
 * 标题栏组件
 * @author 陈凤庆
 * @date 20251001
 * @description 显示应用标题、搜索按钮和更多菜单
 */

const { t } = useI18n()

const props = defineProps({
  appTitle: {
    type: String,
    default: ''
  }
})

/**
 * 显示的标题 - 如果没有传入appTitle，则使用国际化的默认值
 */
const displayTitle = computed(() => {
  return props.appTitle || t('titleBar.defaultAppTitle')
})

const emit = defineEmits(['toggle-search', 'refresh-data'])

/**
 * 处理切换搜索
 */
const handleToggleSearch = () => {
  emit('toggle-search')
}

/**
 * 处理刷新数据事件
 * @author 20251004 陈凤庆
 */
const handleRefreshData = () => {
  emit('refresh-data')
}
</script>

<style scoped>
.title-bar {
  height: 40px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
}

.app-title {
  font-size: 18px;
  font-weight: bold;
  color: #007bff;
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-icon {
  font-size: 16px;
}

.title-bar-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 28px !important;
  height: 28px !important;
  padding: 0 !important;
  border: none !important;
  background: transparent !important;
  color: #6c757d !important;
}

.action-btn:hover {
  background-color: rgba(0, 0, 0, 0.05) !important;
  color: #007bff !important;
}
</style>
