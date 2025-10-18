<template>
  <el-dialog
    v-model="visible"
    :title="$t('export.title')"
    :width="dialogWidth"
    top="5vh"
    :before-close="handleClose"
    :modal="true"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    :destroy-on-close="false"
    :append-to-body="true"
    :lock-scroll="true"
    class="export-vault-dialog"
    draggable
  >
    <div class="export-content">
      <!-- 步骤指示器 -->
      <el-steps :active="currentStep" finish-status="success" align-center>
        <el-step :title="$t('export.steps.verifyPassword')" />
        <el-step :title="$t('export.steps.selectAccounts')" />
        <el-step :title="$t('export.steps.setBackup')" />
        <el-step :title="$t('export.steps.exportComplete')" />
      </el-steps>

      <!-- 步骤1: 验证登录密码 -->
      <div v-if="currentStep === 0" class="step-content">
        <h3>{{ $t('export.verifyPasswordTitle') }}</h3>
        <p>{{ $t('export.verifyPasswordDesc') }}</p>
        <el-form :model="formData" :rules="rules" ref="passwordFormRef" label-width="100px">
          <el-form-item :label="$t('export.loginPassword')" prop="loginPassword">
            <el-input
              v-model="formData.loginPassword"
              type="password"
              :placeholder="$t('export.loginPasswordPlaceholder')"
              show-password
              @keyup.enter="verifyPassword"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤2: 选择要导出的账号 -->
      <div v-if="currentStep === 1" class="step-content">
        <h3>{{ $t('export.selectAccountsTitle') }}</h3>
        <div class="account-selection">
          <!-- 选择模式 -->
          <el-radio-group v-model="selectionMode" @change="handleSelectionModeChange">
            <el-radio label="all">{{ $t('export.exportAll') }}</el-radio>
            <el-radio label="group">{{ $t('export.exportByGroup') }}</el-radio>
            <el-radio label="type">{{ $t('export.exportByType') }}</el-radio>
            <el-radio label="manual">{{ $t('export.exportSelected') }}</el-radio>
          </el-radio-group>

          <!-- 按分组选择 -->
          <div v-if="selectionMode === 'group'" class="selection-content">
            <div class="selection-header">
              <h4>{{ $t('export.selectGroups') }}</h4>
              <div class="selection-actions">
                <el-button size="small" @click="selectAllGroups">{{ $t('export.selectAll') }}</el-button>
                <el-button size="small" @click="clearAllGroups">{{ $t('export.clearAll') }}</el-button>
              </div>
            </div>
            <el-checkbox-group v-model="selectedGroups" @change="handleGroupSelectionChange">
              <el-checkbox
                v-for="group in groups"
                :key="group.id"
                :label="group.id"
              >
                {{ group.name }}
              </el-checkbox>
            </el-checkbox-group>
            <div class="selection-summary">
              {{ $t('export.groupSelectionSummary', { count: selectedGroups.length, accountCount: getGroupAccountCount() }) }}
            </div>
          </div>

          <!-- 按类别选择 -->
          <div v-if="selectionMode === 'type'" class="selection-content">
            <div class="selection-header">
              <h4>{{ $t('export.selectTypes') }}</h4>
              <div class="selection-actions">
                <el-button size="small" @click="selectAllTypes">{{ $t('export.selectAll') }}</el-button>
                <el-button size="small" @click="clearAllTypes">{{ $t('export.clearAll') }}</el-button>
              </div>
            </div>
            <el-checkbox-group v-model="selectedTypes" @change="handleTypeSelectionChange">
              <el-checkbox
                v-for="type in types"
                :key="type.id"
                :label="type.id"
              >
                {{ type.name }}
              </el-checkbox>
            </el-checkbox-group>
            <div class="selection-summary">
              {{ $t('export.typeSelectionSummary', { count: selectedTypes.length, accountCount: getTypeAccountCount() }) }}
            </div>
          </div>

          <!-- 手动选择账号 -->
          <div v-if="selectionMode === 'manual'" class="selection-content">
            <h4>{{ $t('export.selectAccounts') }}</h4>
            <div v-if="isLoadingData" class="loading-container">
              <el-loading-directive v-loading="true" :element-loading-text="$t('export.loadingAccounts')">
                <div style="height: 100px;"></div>
              </el-loading-directive>
            </div>
            <div v-else class="account-list">
              <el-checkbox
                v-model="selectAllAccounts"
                :indeterminate="isIndeterminate"
                @change="handleSelectAllChange"
              >
                {{ $t('export.selectAll') }}
              </el-checkbox>
              <el-checkbox-group v-model="selectedAccounts" @change="handleAccountSelectionChange">
                <div
                  v-for="account in accounts"
                  :key="account.id"
                  class="account-item"
                >
                  <el-checkbox :label="account.id">
                    <div class="account-info">
                      <span class="account-title">{{ account.title }}</span>
                      <span class="account-username">{{ account.masked_username }}</span>
                    </div>
                  </el-checkbox>
                </div>
              </el-checkbox-group>
              <div v-if="accounts.length === 0" class="no-accounts">
                <el-empty :description="$t('export.noAccounts')" />
              </div>
            </div>
          </div>

          <!-- 选择统计 -->
          <div class="selection-summary">
            <el-alert
              :title="`已选择 ${getSelectedAccountCount()} 个账号`"
              type="info"
              :closable="false"
            />
          </div>
        </div>
      </div>

      <!-- 步骤3: 设置备份密码和导出路径 -->
      <div v-if="currentStep === 2" class="step-content">
        <h3>{{ $t('export.setBackupTitle') }}</h3>
        <el-form :model="formData" :rules="backupRules" ref="backupFormRef" label-width="120px">
          <el-form-item :label="$t('export.backupPassword')" prop="backupPassword">
            <el-input
              v-model="formData.backupPassword"
              type="password"
              :placeholder="$t('export.backupPasswordPlaceholder')"
              show-password
            >
              <template #append>
                <el-button @click="generateBackupPassword">{{ $t('export.generate') }}</el-button>
              </template>
            </el-input>
            <div class="form-tip">{{ $t('export.backupPasswordTip') }}</div>
          </el-form-item>

          <el-form-item :label="$t('export.exportPath')" prop="exportPath">
            <el-input
              v-model="formData.exportPath"
              :placeholder="$t('export.exportPathPlaceholder')"
              readonly
            >
              <template #append>
                <el-button @click="selectExportPath">{{ $t('export.browse') }}</el-button>
              </template>
            </el-input>
          </el-form-item>
        </el-form>
      </div>

      <!-- 步骤4: 导出进度和结果 -->
      <div v-if="currentStep === 3" class="step-content">
        <div v-if="exportProgress.isExporting" class="export-progress">
          <h3>{{ $t('export.exporting') }}</h3>
          <el-progress
            :percentage="exportProgress.percentage"
            :status="exportProgress.status"
          />
          <p class="progress-text">{{ exportProgress.currentStep }}</p>
        </div>

        <div v-else-if="exportProgress.completed" class="export-result">
          <div v-if="exportProgress.success" class="success-result">
            <el-result
              icon="success"
              :title="$t('export.exportSuccessTitle')"
              :sub-title="$t('export.exportSuccessSubTitle', { path: formData.exportPath })"
            />
          </div>

          <div v-else class="error-result">
            <el-result
              icon="error"
              :title="$t('export.exportFailedTitle')"
              :sub-title="exportProgress.errorMessage"
            />
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="exportProgress.isExporting">
          {{ currentStep === 3 && exportProgress.completed ? $t('common.close') : $t('common.cancel') }}
        </el-button>
        
        <el-button
          v-if="currentStep > 0 && currentStep < 3"
          @click="previousStep"
          :disabled="exportProgress.isExporting"
        >
          {{ $t('common.previous') }}
        </el-button>
        
        <el-button
          v-if="currentStep < 3"
          type="primary"
          @click="nextStep"
          :disabled="exportProgress.isExporting || !canProceedToNextStep()"
        >
          {{ currentStep === 2 ? $t('export.startExport') : $t('common.next') }}
        </el-button>

        <el-button
          v-if="currentStep === 3 && exportProgress.completed && exportProgress.success"
          type="primary"
          @click="openExportFolder"
        >
          {{ $t('export.openFolder') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { apiService } from '@/services/api'

/**
 * 导出密码库对话框组件
 * @author 陈凤庆
 * @date 2025-10-03
 * @description 提供密码库导出功能的完整界面，包括密码验证、账号选择、备份设置等
 */

// 国际化
const { t } = useI18n()

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['update:modelValue'])

// 响应式数据
const visible = ref(false)
const currentStep = ref(0)
const windowWidth = ref(window.innerWidth)
const isLoadingData = ref(false)

// 表单数据
const formData = reactive({
  loginPassword: '',
  backupPassword: '',
  exportPath: ''
})

// 选择模式
const selectionMode = ref('all')
const selectedGroups = ref([])
const selectedTypes = ref([])
const selectedAccounts = ref([])
const selectAllAccounts = ref(false)

// 数据
const groups = ref([])
const types = ref([])
const accounts = ref([])

// 导出进度
const exportProgress = reactive({
  isExporting: false,
  completed: false,
  success: false,
  percentage: 0,
  status: 'active',
  currentStep: '',
  errorMessage: ''
})

// 表单引用
const passwordFormRef = ref()
const backupFormRef = ref()

// 响应式对话框宽度 - 符合需求：最小300px，最大600px，左右保持边距
const dialogWidth = computed(() => {
  const maxWidth = 600  // 最大600px
  const minWidth = 300  // 最小300px
  const margin = 80     // 左右边距各40px，总共80px
  const availableWidth = windowWidth.value - margin
  const finalWidth = Math.min(maxWidth, Math.max(minWidth, availableWidth))
  return `${finalWidth}px`
})

// 计算属性
const isIndeterminate = computed(() => {
  const selectedCount = selectedAccounts.value.length
  const totalCount = accounts.value.length
  return selectedCount > 0 && selectedCount < totalCount
})

// 表单验证规则
const rules = {
  loginPassword: [
    { required: true, message: t('export.loginPasswordRequired'), trigger: 'blur' }
  ]
}

const backupRules = {
  backupPassword: [
    { required: true, message: t('export.backupPasswordRequired'), trigger: 'blur' },
    { min: 6, message: t('export.backupPasswordMinLength'), trigger: 'blur' }
  ],
  exportPath: [
    { required: true, message: t('export.exportPathRequired'), trigger: 'blur' }
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
    // 默认加载全部导出模式的数据
    loadDataByMode('all')
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
  formData.loginPassword = ''
  formData.backupPassword = ''
  formData.exportPath = ''
  selectionMode.value = 'all'
  selectedGroups.value = []
  selectedTypes.value = []
  selectedAccounts.value = []
  selectAllAccounts.value = false
  
  exportProgress.isExporting = false
  exportProgress.completed = false
  exportProgress.success = false
  exportProgress.percentage = 0
  exportProgress.status = 'active'
  exportProgress.currentStep = ''
  exportProgress.errorMessage = ''
}

/**
 * 按需加载数据
 */
const loadDataByMode = async (mode) => {
  try {
    isLoadingData.value = true
    console.log('按需加载数据，模式:', mode)

    if (mode === 'group') {
      // 加载分组数据
      if (groups.value.length === 0) {
        const groupsData = await apiService.getGroups()
        groups.value = groupsData || []
        console.log('- 分组数量:', groups.value.length)
      }

      // 加载类型数据（用于计算账号数量）
      if (types.value.length === 0) {
        const typesData = await apiService.getTypes()
        types.value = typesData || []
      }

      // 加载账号数据（用于计算数量）
      if (accounts.value.length === 0) {
        const accountsData = await apiService.getAllAccounts()
        accounts.value = accountsData || []
      }
    } else if (mode === 'type') {
      // 加载类型数据
      if (types.value.length === 0) {
        const typesData = await apiService.getTypes()
        types.value = typesData || []
        console.log('- 类型数量:', types.value.length)
      }

      // 加载账号数据（用于计算数量）
      if (accounts.value.length === 0) {
        const accountsData = await apiService.getAllAccounts()
        accounts.value = accountsData || []
      }
    } else if (mode === 'manual') {
      // 加载账号数据
      if (accounts.value.length === 0) {
        const accountsData = await apiService.getAllAccounts()
        accounts.value = accountsData || []
        console.log('- 账号数量:', accounts.value.length)
      }
    } else if (mode === 'all') {
      // 全部导出模式，只需要知道总账号数量
      if (accounts.value.length === 0) {
        const accountsData = await apiService.getAllAccounts()
        accounts.value = accountsData || []
        console.log('- 账号数量:', accounts.value.length)
      }
    }

    console.log('数据加载完成')
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败: ' + error.message)
  } finally {
    isLoadingData.value = false
  }
}

/**
 * 验证登录密码
 */
const verifyPassword = async () => {
  if (!passwordFormRef.value) return false

  try {
    await passwordFormRef.value.validate()

    // 调用后端验证密码
    await apiService.verifyOldPassword(formData.loginPassword)
    ElMessage.success('密码验证成功')
    return true
  } catch (error) {
    console.error('密码验证失败:', error)
    ElMessage.error(error.message || '密码验证失败')
    return false
  }
}

/**
 * 下一步
 */
const nextStep = async () => {
  if (currentStep.value === 0) {
    // 验证登录密码
    if (await verifyPassword()) {
      currentStep.value++
    }
  } else if (currentStep.value === 1) {
    // 检查是否选择了账号
    if (getSelectedAccountCount() === 0) {
      ElMessage.warning('请至少选择一个账号')
      return
    }
    currentStep.value++
  } else if (currentStep.value === 2) {
    // 验证备份设置并开始导出
    if (!backupFormRef.value) return

    try {
      await backupFormRef.value.validate()
      currentStep.value++
      await startExport()
    } catch (error) {
      console.error('表单验证失败:', error)
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
      return formData.loginPassword.trim() !== ''
    case 1:
      const canProceed = getSelectedAccountCount() > 0
      console.log(`步骤1是否可以继续: ${canProceed}, 选中数量: ${getSelectedAccountCount()}`)
      return canProceed
    case 2:
      return formData.backupPassword.trim() !== '' && formData.exportPath.trim() !== ''
    default:
      return false
  }
}

/**
 * 处理选择模式变化
 */
const handleSelectionModeChange = async () => {
  console.log('选择模式变化:', selectionMode.value)
  selectedGroups.value = []
  selectedTypes.value = []
  selectedAccounts.value = []

  // 如果是全部导出模式，设置全选状态
  if (selectionMode.value === 'all') {
    selectAllAccounts.value = true
    console.log('切换到全部导出模式，设置全选状态')
  } else {
    selectAllAccounts.value = false
    console.log('切换到其他模式，取消全选状态')
  }

  // 按需加载数据
  await loadDataByMode(selectionMode.value)
}

/**
 * 处理分组选择变化
 */
const handleGroupSelectionChange = () => {
  // 根据选中的分组获取对应的账号
  const accountsInSelectedGroups = accounts.value.filter(account => {
    const accountType = types.value.find(type => type.id === account.typeid)
    return accountType && selectedGroups.value.includes(accountType.group_id)
  })

  selectedAccounts.value = accountsInSelectedGroups.map(account => account.id)
}

/**
 * 处理类型选择变化
 */
const handleTypeSelectionChange = () => {
  // 根据选中的类型获取对应的账号
  const accountsInSelectedTypes = accounts.value.filter(account =>
    selectedTypes.value.includes(account.typeid)
  )

  selectedAccounts.value = accountsInSelectedTypes.map(account => account.id)
}

/**
 * 处理账号选择变化
 */
const handleAccountSelectionChange = (value) => {
  selectAllAccounts.value = value.length === accounts.value.length
}

/**
 * 处理全选变化
 */
const handleSelectAllChange = (value) => {
  selectedAccounts.value = value ? accounts.value.map(account => account.id) : []
}

/**
 * 全选分组
 */
const selectAllGroups = () => {
  selectedGroups.value = groups.value.map(group => group.id)
  handleGroupSelectionChange()
}

/**
 * 取消全选分组
 */
const clearAllGroups = () => {
  selectedGroups.value = []
  handleGroupSelectionChange()
}

/**
 * 全选类别
 */
const selectAllTypes = () => {
  selectedTypes.value = types.value.map(type => type.id)
  handleTypeSelectionChange()
}

/**
 * 取消全选类别
 */
const clearAllTypes = () => {
  selectedTypes.value = []
  handleTypeSelectionChange()
}

/**
 * 获取分组对应的账号数量
 */
const getGroupAccountCount = () => {
  if (selectedGroups.value.length === 0) return 0

  const accountsInSelectedGroups = accounts.value.filter(account => {
    const accountType = types.value.find(type => type.id === account.typeid)
    return accountType && selectedGroups.value.includes(accountType.group_id)
  })

  return accountsInSelectedGroups.length
}

/**
 * 获取类别对应的账号数量
 */
const getTypeAccountCount = () => {
  if (selectedTypes.value.length === 0) return 0

  const accountsInSelectedTypes = accounts.value.filter(account =>
    selectedTypes.value.includes(account.typeid)
  )

  return accountsInSelectedTypes.length
}

/**
 * 获取选中的账号数量
 */
const getSelectedAccountCount = () => {
  let count = 0

  if (selectionMode.value === 'all') {
    count = accounts.value.length
  } else if (selectionMode.value === 'group') {
    count = getGroupAccountCount()
  } else if (selectionMode.value === 'type') {
    count = getTypeAccountCount()
  } else if (selectionMode.value === 'manual') {
    count = selectedAccounts.value.length
  }

  console.log(`获取选中账号数量: 模式=${selectionMode.value}, 数量=${count}`)
  return count
}

/**
 * 生成备份密码
 */
const generateBackupPassword = () => {
  // 生成一个随机密码
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
  let password = ''
  for (let i = 0; i < 16; i++) {
    password += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  formData.backupPassword = password
  ElMessage.success('备份密码已生成')
}

/**
 * 选择导出路径
 */
const selectExportPath = async () => {
  try {
    const path = await apiService.selectExportPath()
    if (path) {
      formData.exportPath = path
      ElMessage.success('导出路径选择成功')
    }
  } catch (error) {
    console.error('选择导出路径失败:', error)
    ElMessage.error('选择导出路径失败')
  }
}

/**
 * 开始导出
 */
const startExport = async () => {
  exportProgress.isExporting = true
  exportProgress.completed = false
  exportProgress.success = false
  exportProgress.percentage = 0
  exportProgress.status = 'active'
  exportProgress.currentStep = '准备导出...'

  try {
    // 确定导出参数
    let accountIDs = []
    let groupIDs = []
    let typeIDs = []
    let exportAll = false

    if (selectionMode.value === 'all') {
      exportAll = true
      console.log('[导出] 全部导出模式')
    } else if (selectionMode.value === 'group') {
      groupIDs = selectedGroups.value
      console.log('[导出] 按分组导出，分组数量:', groupIDs.length)
    } else if (selectionMode.value === 'type') {
      typeIDs = selectedTypes.value
      console.log('[导出] 按类别导出，类别数量:', typeIDs.length)
    } else if (selectionMode.value === 'manual') {
      accountIDs = selectedAccounts.value
      console.log('[导出] 手动选择导出，账号数量:', accountIDs.length)
    }

    // 更新进度
    exportProgress.percentage = 20
    exportProgress.currentStep = '验证参数...'

    // 调用后端导出API
    await apiService.exportVault(
      formData.loginPassword,
      formData.backupPassword,
      formData.exportPath,
      accountIDs,
      groupIDs,
      typeIDs,
      exportAll
    )

    // 导出成功
    exportProgress.percentage = 100
    exportProgress.status = 'success'
    exportProgress.currentStep = '导出完成'
    exportProgress.success = true
    exportProgress.completed = true

    ElMessage.success('密码库导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    exportProgress.percentage = 100
    exportProgress.status = 'exception'
    exportProgress.currentStep = '导出失败'
    exportProgress.success = false
    exportProgress.completed = true
    exportProgress.errorMessage = error.message || '导出过程中发生未知错误'

    ElMessage.error(error.message || '导出失败')
  } finally {
    exportProgress.isExporting = false
  }
}

/**
 * 打开导出文件夹
 */
const openExportFolder = () => {
  // 这里可以调用系统API打开文件夹
  ElMessage.info('请在文件管理器中查看导出的文件')
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  if (exportProgress.isExporting) {
    ElMessage.warning('导出正在进行中，请稍候...')
    return
  }
  visible.value = false
}
</script>

<style scoped>
.export-vault-dialog {
  --el-dialog-padding-primary: 20px;
}

.export-content {
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

.account-selection {
  margin-top: 20px;
}

.selection-content {
  margin-top: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.selection-content h4 {
  margin-bottom: 15px;
  color: #303133;
}

.selection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.selection-header h4 {
  margin: 0;
  color: #303133;
}

.selection-actions {
  display: flex;
  gap: 8px;
}

.selection-summary {
  margin-top: 15px;
  padding: 10px;
  background-color: #f0f9ff;
  border: 1px solid #b3d8ff;
  border-radius: 4px;
  color: #409eff;
  font-size: 14px;
}

.account-list {
  max-height: 300px;
  overflow-y: auto;
}

.account-item {
  margin-bottom: 10px;
  padding: 8px;
  background-color: white;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.account-info {
  display: flex;
  flex-direction: column;
  margin-left: 8px;
}

.account-title {
  font-weight: 500;
  color: #303133;
}

.account-username {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.selection-summary {
  margin-top: 20px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.export-progress {
  text-align: center;
  padding: 40px 20px;
}

.export-progress h3 {
  margin-bottom: 20px;
}

.progress-text {
  margin-top: 10px;
  color: #606266;
}

.export-result {
  padding: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
