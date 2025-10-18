<template>
  <el-dialog
    v-model="visible"
    :title="$t('importVault.title')"
    :width="dialogWidth"
    top="5vh"
    :before-close="handleClose"
    :modal="true"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    :destroy-on-close="false"
    :append-to-body="true"
    :lock-scroll="true"
    class="import-vault-dialog"
    draggable
  >
    <div class="import-content">
      <!-- 步骤指示器 -->
      <el-steps :active="currentStep" finish-status="success" align-center>
        <el-step :title="$t('importVault.step1')" />
        <el-step :title="$t('importVault.step2')" />
        <el-step :title="$t('importVault.step3')" />
      </el-steps>

      <!-- 步骤1: 选择导入文件 -->
      <div v-if="currentStep === 0" class="step-content">
        <h3>{{ $t('importVault.step1Title') }}</h3>
        <p>{{ $t('importVault.step1Description') }}</p>
        <el-form :model="formData" :rules="rules" ref="fileFormRef" label-width="100px">
          <el-form-item :label="$t('importVault.importFile')" prop="importPath">
            <el-input
              v-model="formData.importPath"
              :placeholder="$t('importVault.selectImportFilePlaceholder')"
              readonly
            >
              <template #append>
                <el-button @click="selectImportFile">{{ $t('importVault.browse') }}</el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤2: 输入解压密码 -->
      <div v-if="currentStep === 1" class="step-content">
        <h3>{{ $t('importVault.step2Title') }}</h3>
        <p>{{ $t('importVault.step2Description') }}</p>
        <el-form :model="formData" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
          <el-form-item :label="$t('importVault.backupPassword')" prop="backupPassword">
            <el-input
              v-model="formData.backupPassword"
              type="password"
              :placeholder="$t('importVault.enterBackupPasswordPlaceholder')"
              show-password
              @keyup.enter="startImport"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤3: 导入进度和结果 -->
      <div v-if="currentStep === 2" class="step-content">
        <div v-if="importProgress.isImporting" class="import-progress">
          <h3>{{ $t('importVault.importingVault') }}</h3>
          <el-progress
            :percentage="importProgress.percentage"
            :status="importProgress.status"
          />
          <p class="progress-text">{{ importProgress.currentStep }}</p>
        </div>

        <div v-else-if="importProgress.completed" class="import-result">
          <div v-if="importProgress.success" class="success-result">
            <el-result
              icon="success"
              :title="$t('importVault.importSuccess')"
              :sub-title="$t('importVault.vaultImportedSuccessfully')"
            />
            
            <!-- 导入报告 -->
            <div class="import-report">
              <h4>{{ $t('importVault.importReport') }}</h4>
              <el-descriptions :column="2" border>
                <el-descriptions-item :label="$t('importVault.totalAccounts')">
                  {{ importResult.total_accounts }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.successfullyImported')">
                  {{ importResult.imported_accounts }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.skippedAccounts')">
                  {{ importResult.skipped_accounts }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.errorAccounts')">
                  {{ importResult.error_accounts }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.totalGroups')">
                  {{ importResult.total_groups }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.importedGroups')">
                  {{ importResult.imported_groups }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.totalTypes')">
                  {{ importResult.total_types }}
                </el-descriptions-item>
                <el-descriptions-item :label="$t('importVault.importedTypes')">
                  {{ importResult.imported_types }}
                </el-descriptions-item>
              </el-descriptions>

              <!-- 跳过的账号详情 -->
              <div v-if="importResult.skipped_account_details && importResult.skipped_account_details.length > 0" class="skipped-accounts">
                <h5>{{ $t('importVault.skippedAccountDetails') }}</h5>
                <el-table
                  :data="importResult.skipped_account_details"
                  size="small"
                  max-height="200"
                >
                  <el-table-column prop="title" :label="$t('importVault.accountTitle')" />
                  <el-table-column prop="name" :label="$t('importVault.accountName')" />
                  <el-table-column prop="id" :label="$t('importVault.accountId')" width="200" />
                </el-table>
              </div>
            </div>
          </div>

          <div v-else class="error-result">
            <el-result
              icon="error"
              :title="$t('importVault.importFailed')"
              :sub-title="importProgress.errorMessage"
            />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="importProgress.isImporting">
          {{ currentStep === 2 && importProgress.completed ? $t('importVault.close') : $t('importVault.cancel') }}
        </el-button>
        
        <el-button
          v-if="currentStep > 0 && currentStep < 2"
          @click="previousStep"
          :disabled="importProgress.isImporting"
        >
          {{ $t('importVault.previousStep') }}
        </el-button>
        
        <el-button
          v-if="currentStep < 2"
          type="primary"
          @click="nextStep"
          :disabled="importProgress.isImporting || !canProceedToNextStep()"
        >
          {{ currentStep === 1 ? $t('importVault.startImport') : $t('importVault.nextStep') }}
        </el-button>

        <el-button
          v-if="currentStep === 2 && importProgress.completed && importProgress.success"
          type="primary"
          @click="refreshData"
        >
          {{ $t('importVault.refreshData') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { apiService } from '@/services/api'
import { useI18n } from 'vue-i18n'

/**
 * 导入密码库对话框组件
 * @author 陈凤庆
 * @date 2025-10-03
 * @description 提供密码库导入功能的完整界面，包括文件选择、密码验证、导入报告等
 */

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'import-success'])

// 响应式数据
const visible = ref(false)
const currentStep = ref(0)
const windowWidth = ref(window.innerWidth)

// 表单数据
const formData = reactive({
  importPath: '',
  backupPassword: ''
})

// 导入进度
const importProgress = reactive({
  isImporting: false,
  completed: false,
  success: false,
  percentage: 0,
  status: 'active',
  currentStep: '',
  errorMessage: ''
})

// 导入结果
const importResult = ref({})

// 表单引用
const fileFormRef = ref()
const passwordFormRef = ref()

// 响应式对话框宽度
const dialogWidth = computed(() => {
  const maxWidth = 600  // 最大600px
  const minWidth = 300  // 最小300px
  const margin = 40
  const availableWidth = windowWidth.value - margin
  const finalWidth = Math.min(maxWidth, Math.max(minWidth, availableWidth))
  return `${finalWidth}px`
})

// 表单验证规则
const { t } = useI18n()

const rules = {
  importPath: [
    { required: true, message: t('importVault.selectImportFileMessage'), trigger: 'blur' }
  ]
}

const passwordRules = {
  backupPassword: [
    { required: true, message: t('importVault.enterBackupPasswordMessage'), trigger: 'blur' }
  ]
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

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    resetDialog()
  }
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

/**
 * 重置对话框状态
 */
const resetDialog = () => {
  currentStep.value = 0
  formData.importPath = ''
  formData.backupPassword = ''
  
  importProgress.isImporting = false
  importProgress.completed = false
  importProgress.success = false
  importProgress.percentage = 0
  importProgress.status = 'active'
  importProgress.currentStep = ''
  importProgress.errorMessage = ''
  
  importResult.value = {}
}

/**
 * 下一步
 */
const nextStep = async () => {
  if (currentStep.value === 0) {
    // 验证文件选择
    if (!fileFormRef.value) return
    
    try {
      await fileFormRef.value.validate()
      currentStep.value++
    } catch (error) {
      console.error('文件选择验证失败:', error)
    }
  } else if (currentStep.value === 1) {
    // 验证密码并开始导入
    if (!passwordFormRef.value) return
    
    try {
      await passwordFormRef.value.validate()
      currentStep.value++
      await startImport()
    } catch (error) {
      console.error('密码验证失败:', error)
    }
  }
}

/**
 * 上一步
 */
const previousStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

/**
 * 检查是否可以进入下一步
 */
const canProceedToNextStep = () => {
  switch (currentStep.value) {
    case 0:
      return formData.importPath.trim() !== ''
    case 1:
      return formData.backupPassword.trim() !== ''
    default:
      return false
  }
}

/**
 * 选择导入文件
 */
const selectImportFile = async () => {
  try {
    const path = await apiService.selectImportFile()
    if (path) {
      formData.importPath = path
      ElMessage.success(t('importVault.selectFileSuccess'))
    }
  } catch (error) {
    console.error('选择导入文件失败:', error)
    ElMessage.error(t('importVault.selectFileFailed'))
  }
}

/**
 * 开始导入
 */
const startImport = async () => {
  importProgress.isImporting = true
  importProgress.completed = false
  importProgress.success = false
  importProgress.percentage = 0
  importProgress.status = 'active'
  importProgress.currentStep = t('importVault.preparingToImport')

  try {
    // 更新进度
    importProgress.percentage = 20
    importProgress.currentStep = t('importVault.validatingFileAndPassword')

    // 调用后端导入API
    const result = await apiService.importVault(
      formData.importPath,
      formData.backupPassword
    )

    // 导入成功
    importProgress.percentage = 100
    importProgress.status = 'success'
    importProgress.currentStep = t('importVault.importComplete')
    importProgress.success = true
    importProgress.completed = true
    importResult.value = result

    ElMessage.success(t('importVault.vaultImportSuccess'))

    // 通知父组件导入成功
    emit('import-success', result)
  } catch (error) {
    console.error('导入失败:', error)
    importProgress.percentage = 100
    importProgress.status = 'exception'
    importProgress.currentStep = t('importVault.importFailedStatus')
    importProgress.success = false
    importProgress.completed = true
    importProgress.errorMessage = error.message || t('importVault.unknownError')

    ElMessage.error(error.message || t('importVault.importFailed'))
  } finally {
    importProgress.isImporting = false
  }
}

/**
 * 刷新数据
 */
const refreshData = () => {
  // 通知父组件刷新数据
  emit('import-success')
  ElMessage.success(t('importVault.dataRefreshed'))
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  if (importProgress.isImporting) {
    ElMessage.warning(t('importVault.importInProgressWarning'))
    return
  }
  visible.value = false
}
</script>

<style scoped>
.import-vault-dialog {
  --el-dialog-padding-primary: 20px;
}

.import-content {
  min-height: 400px;
}

.step-content {
  margin-top: 30px;
  padding: 20px 0;
}

.step-content h3 {
  margin-bottom: 10px;
  color: #303133;
}

.step-content p {
  margin-bottom: 20px;
  color: #606266;
}

.import-progress {
  text-align: center;
  padding: 40px 20px;
}

.import-progress h3 {
  margin-bottom: 20px;
}

.progress-text {
  margin-top: 10px;
  color: #606266;
}

.import-result {
  padding: 20px;
}

.import-report {
  margin-top: 20px;
  text-align: left;
}

.import-report h4 {
  margin-bottom: 15px;
  color: #303133;
}

.skipped-accounts {
  margin-top: 20px;
}

.skipped-accounts h5 {
  margin-bottom: 10px;
  color: #303133;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
