<template>
  <div v-if="visible" class="test-dialog-overlay" @click="handleClose">
    <div class="test-dialog" @click.stop>
      <div class="test-dialog-header">
        <h3>{{ $t('testDialog.title') }}</h3>
        <button @click="handleClose" class="close-btn">×</button>
      </div>
      <div class="test-dialog-body">
        <p>{{ $t('testDialog.description1') }}</p>
        <p>{{ $t('testDialog.description2') }}</p>
      </div>
      <div class="test-dialog-footer">
        <button @click="handleClose" class="btn btn-primary">{{ $t('testDialog.close') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const handleClose = () => {
  visible.value = false
}
</script>

<style scoped>
.test-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.test-dialog {
  background: white;
  border-radius: 8px;
  box-shadow: 0 12px 32px 4px rgba(0, 0, 0, 0.36);
  width: 500px;
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
}

.test-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.test-dialog-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #909399;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #606266;
}

.test-dialog-body {
  padding: 20px;
}

.test-dialog-body p {
  margin: 0 0 16px 0;
  line-height: 1.6;
  color: #606266;
}

.test-dialog-footer {
  padding: 20px;
  border-top: 1px solid #e4e7ed;
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 8px 16px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn-primary {
  background: #409eff;
  border-color: #409eff;
  color: white;
}

.btn-primary:hover {
  background: #66b1ff;
  border-color: #66b1ff;
}
</style>
