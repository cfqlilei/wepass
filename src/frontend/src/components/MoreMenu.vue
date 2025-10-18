<template>
  <div class="more-menu-container">
    <el-dropdown
      ref="dropdownRef"
      trigger="click"
      placement="bottom-end"
      @visible-change="handleVisibleChange"
      @command="handleCommand"
    >
      <el-button
        type="text"
        class="action-btn"
        :title="$t('moreMenu.more')"
      >
        <el-icon><MoreFilled /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <!-- 20251017 陈凤庆 将锁定密码库移动到第一项并增加分隔栏 -->
          <el-dropdown-item command="lockVault">
            <el-icon><Lock /></el-icon>
            {{ $t('moreMenu.lockVault') }}
          </el-dropdown-item>
          <el-dropdown-item divided command="openFile">
            <el-icon><FolderOpened /></el-icon>
            {{ $t('moreMenu.selectNewVault') }}
          </el-dropdown-item>
          <el-dropdown-item command="openVaultDirectory">
            <el-icon><Folder /></el-icon>
            {{ $t('moreMenu.openVaultDirectory') }}
          </el-dropdown-item>
          <!-- 20251017 陈凤庆 删除生成密码菜单项，该功能在其他地方已有提供 -->
          <el-dropdown-item divided command="setPasswordRules">
            <el-icon><Setting /></el-icon>
            {{ $t('moreMenu.setPasswordRules') }}
          </el-dropdown-item>

          <el-dropdown-item divided command="changeLoginPassword">
            <el-icon><Lock /></el-icon>
            {{ $t('moreMenu.changeLoginPassword') }}
          </el-dropdown-item>
          <el-dropdown-item divided command="exportVault">
            <el-icon><Download /></el-icon>
            {{ $t('moreMenu.exportVault') }}
          </el-dropdown-item>
          <el-dropdown-item command="importVault">
            <el-icon><Upload /></el-icon>
            {{ $t('moreMenu.importVault') }}
          </el-dropdown-item>
          <el-dropdown-item divided command="showChangeLog">
            <el-icon><Document /></el-icon>
            {{ $t('moreMenu.changeLog') }}
          </el-dropdown-item>
          <el-dropdown-item command="openSettings">
            <el-icon><Tools /></el-icon>
            {{ $t('moreMenu.settings') }}
          </el-dropdown-item>
          <!-- <el-dropdown-item @click="openTestDialog">
            <el-icon><Tools /></el-icon>
            测试对话框
          </el-dropdown-item> -->
          
          <el-dropdown-item command="showHelp">
            <el-icon><QuestionFilled /></el-icon>
            {{ $t('moreMenu.help') }}
          </el-dropdown-item>
          <el-dropdown-item command="showAbout">
            <el-icon><InfoFilled /></el-icon>
            {{ $t('moreMenu.about') }}
          </el-dropdown-item>
          <!-- 20251017 陈凤庆 添加退出程序功能 -->
          <el-dropdown-item divided command="exitApp">
            <el-icon><SwitchButton /></el-icon>
            {{ $t('moreMenu.exit') }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 密码生成器对话框 -->
    <PasswordGenerator
      v-model="showPasswordGenerator"
      @generated="handlePasswordGenerated"
    />

    <!-- 密码规则设置对话框 -->
    <PasswordRuleSettings
      v-model="showPasswordRuleSettings"
    />

    <!-- 更新日志对话框 -->
    <ChangeLogDialog
      v-model="showChangeLogDialog"
    />

    <!-- 关于对话框 -->
    <AboutDialog
      v-model="showAboutDialog"
    />

    <!-- 帮助对话框 -->
    <HelpDialog
      v-model="showHelpDialog"
    />

    <!-- 设置对话框 -->
    <SettingsDialog
      v-model="showSettingsDialog"
    />

    <!-- 导出密码库对话框 -->
    <ExportVaultDialog
      v-model="showExportDialog"
    />

    <!-- 导入密码库对话框 -->
    <ImportVaultDialog
      v-model="showImportDialog"
      @import-success="handleImportSuccess"
    /> 

    <!-- 修改登录密码对话框 -->
    <el-dialog
      v-model="showChangePasswordDialog"
      :title="$t('moreMenu.changeLoginPassword')"
      width="500px"
      :before-close="handleChangePasswordDialogClose"
    >
      <div class="change-password-content">
        <el-form
          ref="changePasswordFormRef"
          :model="changePasswordForm"
          :rules="changePasswordRules"
          label-width="120px"
        >
          <el-form-item :label="$t('moreMenu.oldLoginPasswordLabel')" prop="oldPassword">
            <el-input
              v-model="changePasswordForm.oldPassword"
              type="password"
              :placeholder="$t('moreMenu.oldLoginPasswordPlaceholder')"
              show-password
            />
          </el-form-item>
          <el-form-item :label="$t('moreMenu.newLoginPasswordLabel')" prop="newPassword">
            <el-input
              v-model="changePasswordForm.newPassword"
              type="password"
              :placeholder="$t('moreMenu.newLoginPasswordPlaceholder')"
              show-password
            />
          </el-form-item>
          <el-form-item :label="$t('moreMenu.confirmNewPasswordLabel')" prop="confirmPassword">
            <el-input
              v-model="changePasswordForm.confirmPassword"
              type="password"
              :placeholder="$t('moreMenu.confirmNewPasswordPlaceholder')"
              show-password
            />
          </el-form-item>
        </el-form>

        <!-- 进度显示 -->
        <div v-if="passwordChangeProgress.show" class="progress-section">
          <div class="progress-info">
            <p>{{ passwordChangeProgress.currentStep || $t('status.updatingAccountData') }}</p>
            <p v-if="passwordChangeProgress.total > 0">{{ $t('status.totalAccounts') }}：{{ passwordChangeProgress.total }}</p>
            <p v-if="passwordChangeProgress.total > 0">{{ $t('status.processed') }}：{{ passwordChangeProgress.current }}</p>
            <p v-if="passwordChangeProgress.total > 0">{{ $t('status.success') }}：{{ passwordChangeProgress.success }}</p>
            <p v-if="passwordChangeProgress.total > 0">{{ $t('status.failed') }}：{{ passwordChangeProgress.error }}</p>
          </div>
          <el-progress
            :percentage="passwordChangeProgress.percentage"
            :status="passwordChangeProgress.status"
          />
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleChangePasswordDialogClose" :disabled="passwordChangeProgress.isChangingPassword">{{ $t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            @click="handleChangePassword"
            :loading="passwordChangeProgress.show"
            :disabled="passwordChangeProgress.isChangingPassword"
          >
            {{ passwordChangeProgress.show ? (passwordChangeProgress.isChangingPassword ? $t('status.changingPassword') : $t('status.verifying')) : $t('common.confirm') }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 文件选择对话框 -->
    <el-dialog
      v-model="showFileDialog"
      :title="$t('moreMenu.selectVaultFile')"
      width="500px"
      :before-close="handleFileDialogClose"
    >
      <div class="file-dialog-content">
        <p>{{ $t('moreMenu.selectVaultFilePrompt') }}</p>
        <el-input
          v-model="selectedFilePath"
          :placeholder="$t('moreMenu.selectVaultFilePlaceholder')"
          readonly
          @click="selectFile"
        >
          <template #append>
            <el-button @click="selectFile">{{ $t('moreMenu.browse') }}</el-button>
          </template>
        </el-input>
        <div class="file-tips">
          <p>{{ $t('moreMenu.supportedFileFormats') }}</p>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleFileDialogClose">{{ $t('common.cancel') }}</el-button>
          <el-button
            type="primary"
            @click="handleOpenFile"
            :disabled="!selectedFilePath"
          >
            {{ $t('common.open') }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
/**
 * 更多菜单组件
 * @author 陈凤庆
 * @modify 陈凤庆 添加文件选择和帮助对话框功能
 * @modify 20251002 陈凤庆 添加打开密码库所在目录功能
 */
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { useVaultStore } from '@/stores/vault'
import { useI18n } from 'vue-i18n'
import {
  MoreFilled,
  FolderOpened,
  Folder,
  Key,
  Setting,
  Tools,
  Lock,
  InfoFilled,
  QuestionFilled,
  SwitchButton,
  Document,
  Download,
  Upload
} from '@element-plus/icons-vue'
import PasswordGenerator from './PasswordGenerator.vue'
import PasswordRuleSettings from './PasswordRuleSettings.vue'
import ChangeLogDialog from './ChangeLogDialog.vue'
import AboutDialog from './AboutDialog.vue'
import HelpDialog from './HelpDialog.vue'
import SettingsDialog from './SettingsDialog.vue'
import TestDialog from './TestDialog.vue'
import ExportVaultDialog from './ExportVaultDialog.vue'
import ImportVaultDialog from './ImportVaultDialog.vue'
import { SelectVaultFile } from '../../wailsjs/go/app/App'
import { apiService } from '@/services/api'

export default {
  name: 'MoreMenu',
  emits: ['refresh-data'],
  components: {
    PasswordGenerator,
    PasswordRuleSettings,
    ChangeLogDialog,
    AboutDialog,
    HelpDialog,
    SettingsDialog,
    TestDialog,
    ExportVaultDialog,
    ImportVaultDialog
  },
  setup(props, { emit }) {
    const router = useRouter()
    const vaultStore = useVaultStore()
    const { t } = useI18n()
    
    const dropdownRef = ref(null)
    const showPasswordGenerator = ref(false)
    const showPasswordRuleSettings = ref(false)
    const showChangeLogDialog = ref(false)
    const showAboutDialog = ref(false)
    const showHelpDialog = ref(false)
    const showSettingsDialog = ref(false)
    const showTestDialog = ref(false)
    const showExportDialog = ref(false)
    const showImportDialog = ref(false)
    const showFileDialog = ref(false)
    const selectedFilePath = ref('')
    const showChangePasswordDialog = ref(false)
    const changePasswordFormRef = ref(null)

    // 修改密码表单数据
    const changePasswordForm = reactive({
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    })

    // 密码修改进度
    const passwordChangeProgress = reactive({
      show: false,
      total: 0,
      current: 0,
      success: 0,
      error: 0,
      percentage: 0,
      status: '',
      currentStep: '', // 20251003 陈凤庆 当前步骤描述
      isChangingPassword: false // 20251003 陈凤庆 是否正在修改密码（区分验证阶段）
    })

    // 表单验证规则
    const changePasswordRules = {
      oldPassword: [
        { required: true, message: t('moreMenu.oldLoginPasswordRequired'), trigger: 'blur' }
      ],
      newPassword: [
        { required: true, message: t('moreMenu.newLoginPasswordRequired'), trigger: 'blur' },
        { min: 8, message: t('moreMenu.newLoginPasswordMinLength'), trigger: 'blur' },
        {
          validator: (_rule, value, callback) => {
            if (!value) {
              callback()
              return
            }
            // 密码强度验证
            const hasUpper = /[A-Z]/.test(value)
            const hasLower = /[a-z]/.test(value)
            const hasNumber = /\d/.test(value)
            const hasSpecial = /[!@#$%^&*()_+\-=\[\]{}|;:',.<>/?]/.test(value)

            if (!hasUpper || !hasLower || !hasNumber || !hasSpecial) {
              callback(new Error(t('moreMenu.newLoginPasswordStrength')))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ],
      confirmPassword: [
        { required: true, message: t('moreMenu.confirmNewLoginPasswordRequired'), trigger: 'blur' },
        {
          validator: (_rule, value, callback) => {
            if (value !== changePasswordForm.newPassword) {
              callback(new Error(t('moreMenu.passwordsDoNotMatch')))
            } else {
              callback()
            }
          },
          trigger: 'blur'
        }
      ]
    }

    /**
     * 处理下拉菜单显示状态变化
     * @param {boolean} visible - 是否显示
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const handleVisibleChange = (visible) => {
      // 可以在这里添加下拉菜单状态变化的处理逻辑
    }

    /**
     * 处理下拉菜单命令
     * @author 陈凤庆
     * @date 2025-10-05
     */
    const handleCommand = (command) => {
      console.log('[MoreMenu] handleCommand 被调用，命令:', command)

      switch (command) {
        case 'openFile':
          openFile()
          break
        case 'openVaultDirectory':
          openVaultDirectory()
          break
        // 20251017 陈凤庆 删除生成密码命令处理，该功能在其他地方已有提供
        case 'setPasswordRules':
          setPasswordRules()
          break
        case 'changeLoginPassword':
          changeLoginPassword()
          break
        case 'exportVault':
          exportVault()
          break
        case 'importVault':
          importVault()
          break
        case 'showChangeLog':
          showChangeLog()
          break
        case 'openSettings':
          openSettings()
          break
        case 'showHelp':
          showHelp()
          break
        case 'showAbout':
          showAbout()
          break
        case 'lockVault':
          lockVault()
          break
        case 'exitApp':
          exitApp()
          break
        default:
          console.warn('[MoreMenu] 未知命令:', command)
      }
    }

    /**
     * 选择新的密码框
     * @author 陈凤庆
     * @date 2025-10-03
     * @description 返回到登录窗口的完全界面，用于选择新的密码框
     */
    const openFile = async () => {
      try {
        await ElMessageBox.confirm(
          t('moreMenu.selectNewVaultConfirm'),
          t('moreMenu.selectNewVault'),
          {
            confirmButtonText: t('common.continue'),
            cancelButtonText: t('common.cancel'),
            type: 'info'
          }
        )

        // 调用后端关闭密码库
        try {
          await apiService.closeVault()
          console.log('[选择新密码框] 后端密码库已关闭')
        } catch (error) {
          console.error('[选择新密码框] 关闭后端密码库失败:', error)
          // 即使后端关闭失败，也继续前端清理流程
        }

        // 清除前端vault状态
        vaultStore.closeVault()

        // 导航到登录页面，并确保使用完全模式
        await router.push('/login?mode=full')

        ElMessage.success(t('moreMenu.selectNewVaultSuccess'))
      } catch {
        // 用户取消操作
      }
    }

    /**
     * 选择文件
     * @author 陈凤庆
     * @date 2025-01-27
     * @modify 陈凤庆 修复文件选择功能，使用Wails的SelectVaultFile API
     */
    const selectFile = async () => {
      try {
        // 20251001 陈凤庆 使用 Wails 的文件选择对话框API
        const selectedPath = await SelectVaultFile()
        if (selectedPath) {
          selectedFilePath.value = selectedPath
          ElMessage.success(t('moreMenu.fileSelectionSuccess'))
        }
      } catch (error) {
        console.error('选择文件失败:', error)
        ElMessage.error(t('moreMenu.fileSelectionFailed'))
      }
    }

    /**
     * 处理打开文件
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const handleOpenFile = async () => {
      if (!selectedFilePath.value) {
        ElMessage.warning(t('moreMenu.pleaseSelectFile'))
        return
      }

      try {
        // TODO: 实现打开密码库文件的逻辑
        ElMessage.success(t('moreMenu.openingFile', { filePath: selectedFilePath.value }))
        handleFileDialogClose()
      } catch (error) {
        console.error('打开文件失败:', error)
        ElMessage.error(t('moreMenu.openFileFailed'))
      }
    }

    /**
     * 关闭文件选择对话框
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const handleFileDialogClose = () => {
      showFileDialog.value = false
      selectedFilePath.value = ''
    }

    // 20251017 陈凤庆 删除生成密码函数，该功能在其他地方已有提供

    /**
     * 处理密码生成完成
     * @param {string} password - 生成的密码
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const handlePasswordGenerated = (password) => {
      // 可以在这里处理生成的密码，比如复制到剪贴板
      if (navigator.clipboard) {
        navigator.clipboard.writeText(password).then(() => {
          ElMessage.success(t('moreMenu.passwordCopied'))
        }).catch(() => {
          ElMessage.warning(t('moreMenu.copyFailedManual'))
        })
      }
    }

    /**
     * 设置密码规则
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const setPasswordRules = () => {
      showPasswordRuleSettings.value = true
    }

    /**
     * 修改登录密码
     * @author 20251003 陈凤庆
     */
    const changeLoginPassword = () => {
      showChangePasswordDialog.value = true
    }

    /**
     * 处理修改密码对话框关闭
     */
    const handleChangePasswordDialogClose = () => {
      if (passwordChangeProgress.isChangingPassword) {
        ElMessage.warning(t('moreMenu.passwordChangeInProgress'))
        return
      }
      showChangePasswordDialog.value = false
      // 重置表单
      changePasswordForm.oldPassword = ''
      changePasswordForm.newPassword = ''
      changePasswordForm.confirmPassword = ''
      // 重置进度
      passwordChangeProgress.show = false
      passwordChangeProgress.total = 0
      passwordChangeProgress.current = 0
      passwordChangeProgress.success = 0
      passwordChangeProgress.error = 0
      passwordChangeProgress.percentage = 0
      passwordChangeProgress.status = ''
      passwordChangeProgress.currentStep = ''
      passwordChangeProgress.isChangingPassword = false
    }

    /**
     * 处理修改密码
     */
    const handleChangePassword = async () => {
      try {
        // 表单验证
        if (!changePasswordFormRef.value) {
          ElMessage.error(t('moreMenu.formRefNotFound'))
          return
        }

        const valid = await changePasswordFormRef.value.validate()
        if (!valid) {
          return
        }

        // 20251003 陈凤庆 先验证旧密码，不显示验证进度（避免一闪而过影响体验）
        try {
          // 先验证旧密码，不显示进度
          await apiService.verifyOldPassword(changePasswordForm.oldPassword)

          // 验证成功，不显示成功提示，直接进入确认对话框

        } catch (error) {
          // 旧密码验证失败，直接显示错误消息，不显示进度
          console.error('旧密码验证失败:', error)
          ElMessage.error(error.message || t('moreMenu.oldPasswordIncorrect'))
          return
        }

        // 确认对话框
        await ElMessageBox.confirm(
          t('moreMenu.changePasswordConfirm'),
          t('moreMenu.changeLoginPassword'),
          {
            confirmButtonText: t('common.continue'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
        )

        // 开始修改密码
        passwordChangeProgress.show = true
        passwordChangeProgress.status = 'active'
        passwordChangeProgress.percentage = 0
        passwordChangeProgress.currentStep = t('status.startingChangeLoginPassword')
        passwordChangeProgress.isChangingPassword = true // 修改阶段，按钮禁用

        try {
          // 调用后端API
          await apiService.changeLoginPassword(
            changePasswordForm.oldPassword,
            changePasswordForm.newPassword
          )

          passwordChangeProgress.status = 'success'
          passwordChangeProgress.percentage = 100
          passwordChangeProgress.currentStep = '登录密码修改完成'
          passwordChangeProgress.isChangingPassword = false

          ElMessage.success(t('moreMenu.passwordChangeSuccess'))

          // 延迟关闭对话框
          setTimeout(() => {
            handleChangePasswordDialogClose()
          }, 2000)

        } catch (error) {
          passwordChangeProgress.status = 'exception'
          passwordChangeProgress.currentStep = t('status.changeFailed')
          passwordChangeProgress.isChangingPassword = false
          console.error('修改登录密码失败:', error)
          ElMessage.error(error.message || t('moreMenu.passwordChangeFailed'))
        }
      } catch {
        // 用户取消操作，重置进度状态
        passwordChangeProgress.show = false
        passwordChangeProgress.status = 'active'
        passwordChangeProgress.percentage = 0
        passwordChangeProgress.currentStep = ''
        passwordChangeProgress.isChangingPassword = false
      }
    }

    /**
     * 导出密码库
     * @author 20251003 陈凤庆
     */
    const exportVault = () => {
      showExportDialog.value = true
    }

    /**
     * 导入密码库
     * @author 20251003 陈凤庆
     */
    const importVault = () => {
      showImportDialog.value = true
    }

    /**
     * 处理导入成功事件
     * @author 20251003 陈凤庆
     * @modify 20251004 陈凤庆 添加数据重新加载逻辑
     */
    const handleImportSuccess = async () => {
      try {
        // 通知父组件刷新数据
        emit('refresh-data')
        ElMessage.success(t('moreMenu.importSuccessDataUpdated'))
      } catch (error) {
        console.error('刷新数据失败:', error)
        ElMessage.error(t('moreMenu.importSuccessDataRefreshFailed'))
      }
    }

    /**
     * 显示更新日志对话框
     * @author 陈凤庆
     * @date 2025-10-01
     */
    const showChangeLog = () => {
      showChangeLogDialog.value = true
    }

    /**
     * 显示关于对话框
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const showAbout = () => {
      console.log('[MoreMenu] showAbout 被调用')
      showAboutDialog.value = true
      console.log('[MoreMenu] showAboutDialog 设置为:', showAboutDialog.value)
    }

    /**
     * 显示帮助对话框
     * @author 陈凤庆
     * @date 2025-01-27
     */
    const showHelp = () => {
      showHelpDialog.value = true
    }

    /**
     * 打开设置对话框
     * @author 陈凤庆
     * @date 2025-10-03
     */
    const openSettings = () => {
      console.log('Opening settings dialog...')
      showSettingsDialog.value = true
      console.log('showSettingsDialog.value:', showSettingsDialog.value)
    }

    /**
     * 打开测试对话框
     * @author 陈凤庆
     * @date 2025-10-03
     */
    const openTestDialog = () => {
      console.log('Opening test dialog...')
      showTestDialog.value = true
      console.log('showTestDialog.value:', showTestDialog.value)
    }

    /**
     * 打开密码库所在目录
     * @author 陈凤庆
     * @date 2025-10-02
     * @description 打开当前密码库文件所在的目录，支持跨平台
     */
     const openVaultDirectory = async () => {
      try {
        await apiService.openVaultDirectory()
        ElMessage.success(t('moreMenu.vaultDirectoryOpened'))
      } catch (error) {
        console.error('打开目录失败:', error)
        ElMessage.error(error.message || t('moreMenu.openDirectoryFailed'))
      }
    }

    /**
     * 锁定密码库
     * @author 陈凤庆
     * @date 2025-10-04
     * @description 手动锁定密码库，返回登录界面
     */
    const lockVault = async () => {
      try {
        await ElMessageBox.confirm(
          t('moreMenu.lockVaultConfirm'),
          t('moreMenu.lockVault'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
        )

        // 调用后端锁定密码库
        try {
          await apiService.triggerLock()
          console.log('[锁定密码库] 后端锁定成功')
        } catch (error) {
          console.error('[锁定密码库] 后端锁定失败:', error)
          // 即使后端锁定失败，也继续前端锁定流程
        }

        // 更新前端状态
        vaultStore.setVaultOpened(false, '')
        console.log('[锁定密码库] 前端状态已更新')

        // 20251017 陈凤庆 清理sessionStorage中的敏感数据
        // 注意：登录密码现在存储在后端内存中，会在密码库关闭时自动清理
        sessionStorage.clear()
        console.log('[锁定密码库] 已清理sessionStorage中的敏感数据')

        // 跳转到登录页面
        await router.push('/login')
        console.log('[锁定密码库] 已跳转到登录页面')

        ElMessage.success(t('moreMenu.vaultLocked'))
      } catch (error) {
        if (error !== 'cancel') {
          console.error('锁定密码库失败:', error)
          ElMessage.error(t('moreMenu.lockVaultFailed'))
        }
      }
    }

    /**
     * 退出程序
     * @author 陈凤庆
     * @date 2025-10-17
     * @description 关闭整个应用程序
     */
    const exitApp = async () => {
      try {
        await ElMessageBox.confirm(
          t('moreMenu.exitConfirm'),
          t('moreMenu.exit'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
        )

        // 调用Wails运行时退出程序
        if (window.runtime && window.runtime.Quit) {
          window.runtime.Quit()
        } else {
          console.warn('[退出程序] Wails运行时不可用，尝试关闭窗口')
          window.close()
        }
      } catch {
        // 用户取消操作
      }
    }

    return {
      // 图标
      MoreFilled,
      FolderOpened,
      Folder,
      Key,
      Setting,
      Tools,
      Lock,
      InfoFilled,
      QuestionFilled,
      SwitchButton,
      Document,
      Download,
      Upload,

      // 引用
      dropdownRef,

      // 响应式数据
      showPasswordGenerator,
      showPasswordRuleSettings,
      showChangeLogDialog,
      showAboutDialog,
      showHelpDialog,
      showSettingsDialog,
      showTestDialog,
      showExportDialog,
      showImportDialog,
      showFileDialog,
      selectedFilePath,
      showChangePasswordDialog,
      changePasswordFormRef,
      changePasswordForm,
      passwordChangeProgress,
      changePasswordRules,

      // 方法
      handleVisibleChange,
      handleCommand,
      openFile,
      selectFile,
      handleOpenFile,
      handleFileDialogClose,
      openVaultDirectory,
      // 20251017 陈凤庆 删除生成密码函数导出
      handlePasswordGenerated,
      setPasswordRules,
      openSettings,
      openTestDialog,
      changeLoginPassword,
      handleChangePasswordDialogClose,
      handleChangePassword,
      exportVault,
      importVault,
      handleImportSuccess,
      showChangeLog,
      showAbout,
      showHelp,
      lockVault,
      exitApp
    }
  }
}
</script>

<style scoped>
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

.el-dropdown-menu__item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.el-dropdown-menu__item .el-icon {
  font-size: 14px;
}

/* 文件对话框样式 */
.file-dialog-content {
  padding: 16px 0;
}

.file-dialog-content p {
  margin-bottom: 16px;
  color: #606266;
}

.file-tips {
  margin-top: 12px;
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border-left: 3px solid #409eff;
}

.file-tips p {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 修改密码对话框样式 */
.change-password-content {
  padding: 16px 0;
}

.progress-section {
  margin-top: 24px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border-left: 3px solid #409eff;
}

.progress-info {
  margin-bottom: 16px;
}

.progress-info p {
  margin: 4px 0;
  font-size: 14px;
  color: #606266;
}
</style>
