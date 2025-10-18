<!--
/**
 * 帮助对话框组件
 * @author 陈凤庆
 * @date 2025-01-27
 * @description 显示应用程序的帮助信息和使用说明
 */
-->
<template>
  <el-dialog
    v-model="visible"
    :title="$t('help.title')"
    :width="dialogWidth"
    :before-close="handleClose"
    class="help-dialog"
  >
    <div class="help-content">
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 快速入门 -->
        <el-tab-pane :label="$t('help.quickStart')" name="quickstart">
          <div class="help-section">
            <h3>{{ $t('help.welcomeTitle') }}</h3>
            <p>{{ $t('help.welcomeDescription') }}</p>

            <h4>{{ $t('help.basicOperations') }}</h4>
            <ul>
              <li><strong>{{ $t('help.createVault') }}：</strong>{{ $t('help.createVaultDesc') }}</li>
              <li><strong>{{ $t('help.addPassword') }}：</strong>{{ $t('help.addPasswordDesc') }}</li>
              <li><strong>{{ $t('help.searchPassword') }}：</strong>{{ $t('help.searchPasswordDesc') }}</li>
              <li><strong>{{ $t('help.editPassword') }}：</strong>{{ $t('help.editPasswordDesc') }}</li>
              <li><strong>{{ $t('help.deletePassword') }}：</strong>{{ $t('help.deletePasswordDesc') }}</li>
            </ul>
          </div>
        </el-tab-pane>

        <!-- 功能说明 -->
        <el-tab-pane :label="$t('help.features')" name="features">
          <div class="help-section">
            <h4>{{ $t('help.mainFeatures') }}</h4>
            <ul>
              <li><strong>{{ $t('help.passwordGenerator') }}:</strong> {{ $t('help.passwordGeneratorDesc') }}</li>
              <li><strong>{{ $t('help.groupManagement') }}:</strong> {{ $t('help.groupManagementDesc') }}</li>
              <li><strong>{{ $t('help.secureEncryption') }}:</strong> {{ $t('help.secureEncryptionDesc') }}</li>
              <li><strong>{{ $t('help.importExport') }}:</strong> {{ $t('help.importExportDesc') }}</li>
              <li><strong>{{ $t('help.backupRestore') }}:</strong> {{ $t('help.backupRestoreDesc') }}</li>
            </ul>

            <h4>{{ $t('help.passwordGenerationRules') }}</h4>
            <p>{{ $t('help.passwordGenerationRulesDesc') }}</p>
            <ul>
              <li><code>a</code> - {{ $t('help.lowercaseLetters') }}</li>
              <li><code>A</code> - {{ $t('help.mixedCaseLetters') }}</li>
              <li><code>U</code> - {{ $t('help.uppercaseLetters') }}</li>
              <li><code>d</code> - {{ $t('help.digits') }}</li>
              <li><code>s</code> - {{ $t('help.specialCharacters') }}</li>
              <li><code>[...]</code> - {{ $t('help.customCharacterSet') }}</li>
            </ul>
          </div>
        </el-tab-pane>

        <!-- 安全提示 -->
        <el-tab-pane :label="$t('help.securityTips')" name="security">
          <div class="help-section">
            <h4>{{ $t('help.securitySuggestions') }}</h4>
            <ul>
              <li><strong>{{ $t('help.masterPassword') }}:</strong> {{ $t('help.masterPasswordDesc') }}</li>
              <li><strong>{{ $t('help.regularBackup') }}:</strong> {{ $t('help.regularBackupDesc') }}</li>
              <li><strong>{{ $t('help.timelyUpdate') }}:</strong> {{ $t('help.timelyUpdateDesc') }}</li>
              <li><strong>{{ $t('help.avoidRepetition') }}:</strong> {{ $t('help.avoidRepetitionDesc') }}</li>
              <li><strong>{{ $t('help.safeEnvironment') }}:</strong> {{ $t('help.safeEnvironmentDesc') }}</li>
            </ul>

            <h4>{{ $t('help.precautions') }}</h4>
            <ul>
              <li>{{ $t('help.precaution1') }}</li>
              <li>{{ $t('help.precaution2') }}</li>
              <li>{{ $t('help.precaution3') }}</li>
            </ul>
          </div>
        </el-tab-pane>

        <!-- 常见问题 -->
        <el-tab-pane :label="$t('help.faq')" name="faq">
          <div class="help-section">
            <h4>{{ $t('help.faqTitle') }}</h4>
            
            <div class="faq-item">
              <h5>Q: {{ $t('help.faq1_q') }}</h5>
              <p>A: {{ $t('help.faq1_a') }}</p>
            </div>

            <div class="faq-item">
              <h5>Q: {{ $t('help.faq2_q') }}</h5>
              <p>A: {{ $t('help.faq2_a') }}</p>
            </div>

            <div class="faq-item">
              <h5>Q: {{ $t('help.faq3_q') }}</h5>
              <p>A: {{ $t('help.faq3_a') }}</p>
            </div>

            <div class="faq-item">
              <h5>Q: {{ $t('help.faq4_q') }}</h5>
              <p>A: {{ $t('help.faq4_a') }}</p>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="handleClose">{{ $t('common.ok') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  name: 'HelpDialog',
  props: {
    modelValue: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const activeTab = ref('quickstart')

    // 窗口宽度响应式变量
    const windowWidth = ref(window.innerWidth)

    const visible = computed({
      get: () => props.modelValue,
      set: (value) => emit('update:modelValue', value)
    })

    // 响应式对话框宽度
    const dialogWidth = computed(() => {
      // 最大宽度700px，最小宽度350px（小于主界面最小宽度400px）
      const maxWidth = 700
      const minWidth = 350

      // 左右边距各30px，总共60px，保留适当边距
      const margin = 60

      // 计算可用宽度
      const availableWidth = windowWidth.value - margin

      // 返回合适的宽度，确保不超过最大宽度且不小于最小宽度，保留边距
      const finalWidth = Math.min(maxWidth, Math.max(minWidth, availableWidth))

      return `${finalWidth}px`
    })

    /**
     * 处理对话框关闭
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const handleClose = () => {
      visible.value = false
    }

    // 监听窗口大小变化
    const handleResize = () => {
      windowWidth.value = window.innerWidth
    }

    // 组件挂载时添加监听器
    onMounted(() => {
      window.addEventListener('resize', handleResize)
    })

    // 组件卸载时移除监听器
    onUnmounted(() => {
      window.removeEventListener('resize', handleResize)
    })

    return {
      activeTab,
      visible,
      dialogWidth,
      handleClose
    }
  }
}
</script>

<style scoped>
.help-dialog {
  --el-dialog-content-font-size: 14px;
}

.help-content {
  max-height: 500px;
  overflow-y: auto;
}

.help-section {
  padding: 20px;
  line-height: 1.6;
}

.help-section h3 {
  color: #409eff;
  margin-bottom: 16px;
  font-size: 18px;
}

.help-section h4 {
  color: #303133;
  margin: 20px 0 12px 0;
  font-size: 16px;
}

.help-section h5 {
  color: #606266;
  margin: 16px 0 8px 0;
  font-size: 14px;
  font-weight: 600;
}

.help-section ul {
  margin: 12px 0;
  padding-left: 20px;
}

.help-section li {
  margin: 8px 0;
  color: #606266;
}

.help-section code {
  background-color: #f5f7fa;
  color: #e6a23c;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.faq-item {
  margin-bottom: 20px;
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 6px;
  border-left: 4px solid #409eff;
}

.faq-item h5 {
  margin-top: 0;
  color: #409eff;
}

.faq-item p {
  margin-bottom: 0;
  color: #606266;
}

.dialog-footer {
  text-align: right;
}
</style>