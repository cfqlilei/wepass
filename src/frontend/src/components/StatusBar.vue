<template>
  <div class="status-bar">
    <div>{{ $t('statusBar.authInfo') }}</div>
    <div>{{ $t('statusBar.usageCount', { count: totalUseCount }) }} | {{ $t('statusBar.usageDays', { days: usageDays }) }}</div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

/**
 * 状态栏组件
 * @author 陈凤庆
 * @date 20251001
 * @description 显示底部状态信息
 */
const props = defineProps({
  totalUseCount: {
    type: Number,
    default: 0
  },
  usageDays: {
    type: Number,
    default: 0
  }
})

const { t } = useI18n()
</script>

<style scoped>
.status-bar {
  height: 30px;
  border-top: 1px solid #dee2e6;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 10px;
  color: #6c757d;
  background-color: #ffffff;
}
</style>
