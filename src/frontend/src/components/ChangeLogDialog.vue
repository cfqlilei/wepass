<template>
  <el-dialog
    v-model="visible"
    :title="$t('dialog.changeLogTitle')"
    :width="dialogWidth"
    :before-close="handleClose"
    class="changelog-dialog"
  >
    <div class="changelog-content">
      <div v-for="(entry, index) in changeLog" :key="index" class="changelog-entry">
        <div class="version-header">
          <el-tag type="primary" size="large">v{{ entry.Version }}</el-tag>
          <span class="release-date">{{ entry.Date }}</span>
        </div>
        <ul class="changes-list">
          <li v-for="(change, idx) in entry.Changes" :key="idx" class="change-item">
            <span :class="getChangeTypeClass(change)">{{ change }}</span>
          </li>
        </ul>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="handleClose">{{ $t('common.confirm') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { apiService } from '../services/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const visible = ref(props.modelValue)
const changeLog = ref([])
const windowWidth = ref(window.innerWidth)

// 计算对话框宽度，自适应主界面宽度，保留边距
const dialogWidth = computed(() => {
  // 左右外边距各20px，总共40px
  const horizontalMargin = 40
  // 计算可用宽度
  const availableWidth = windowWidth.value - horizontalMargin
  // 最大宽度1000px，最小宽度600px
  const maxWidth = 1000
  const minWidth = 600
  // 返回合适的宽度
  const finalWidth = Math.min(maxWidth, Math.max(minWidth, availableWidth))
  return finalWidth + 'px'
})

// 窗口大小变化监听
const handleResize = () => {
  windowWidth.value = window.innerWidth
}

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    loadChangeLog()
  }
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

// 生命周期钩子
onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

/**
 * 加载更新日志
 */
const loadChangeLog = async () => {
  try {
    const log = await apiService.getChangeLog()
    changeLog.value = log
  } catch (error) {
    console.error(t('error.loadChangeLogFailed'), error)
  }
}

/**
 * 获取变更类型的样式类
 */
const getChangeTypeClass = (change) => {
  if (change.startsWith(t('changeLog.new'))) {
    return 'change-new'
  } else if (change.startsWith(t('changeLog.optimize'))) {
    return 'change-improve'
  } else if (change.startsWith(t('changeLog.fix'))) {
    return 'change-fix'
  }
  return ''
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  visible.value = false
}
</script>

<style scoped>
/* 对话框样式 */
.changelog-dialog :deep(.el-dialog) {
  margin: 50px auto; /* 上下外边距50px */
  max-height: calc(100vh - 100px); /* 减去上下边距 */
}

.changelog-dialog :deep(.el-dialog__body) {
  padding: 20px;
  max-height: calc(100vh - 220px); /* 减去上下边距100px + 头部和底部的高度120px */
  overflow: hidden;
}

.changelog-content {
  max-height: calc(100vh - 260px); /* 自适应高度，保留所有边距 */
  overflow-y: auto;
  padding: 10px;
  margin: -5px; /* 抵消内边距 */
  padding: 5px;
}

.changelog-entry {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.changelog-entry:last-child {
  border-bottom: none;
}

.version-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.release-date {
  color: #909399;
  font-size: 14px;
}

.changes-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.change-item {
  padding: 8px 0;
  padding-left: 20px;
  position: relative;
  line-height: 1.6;
  font-size: 14px;
}

.change-item::before {
  content: '•';
  position: absolute;
  left: 5px;
  color: #409eff;
  font-weight: bold;
}

.change-new {
  color: #67c23a;
}

.change-improve {
  color: #409eff;
}

.change-fix {
  color: #f56c6c;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}

/* 滚动条样式 */
.changelog-content::-webkit-scrollbar {
  width: 6px;
}

.changelog-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.changelog-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.changelog-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
