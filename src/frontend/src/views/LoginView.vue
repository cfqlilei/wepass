<template>
  <div class="login-container" :class="{ 'simplified-mode': isSimplifiedMode }">
    <!-- 登录头部 -->
    <div class="login-header">
      <!-- 返回按钮（仅在从简化模式切换到完整模式时显示） -->
      <div v-if="!isSimplifiedMode && showBackButton" class="back-button-container">
        <el-button type="text" @click="handleBackToSimplified" class="back-button">
          <el-icon><ArrowLeft /></el-icon>
          {{ $t('common.back') }}
        </el-button>
      </div>

      <!-- 登录标题 -->
      <h1 class="login-title">{{ appTitle }}</h1>
      <p v-if="!isSimplifiedMode" class="login-desc">{{ $t('login.title') }}</p>

      <div v-if="isSimplifiedMode" class="more-menu-container">
        <el-dropdown @command="handleMoreCommand" placement="bottom-end">
          <el-button type="text" class="more-button">
            <div class="hamburger-icon">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="select-vault">{{ $t('login.openVault') }}</el-dropdown-item>
              <el-dropdown-item command="create-vault">{{ $t('login.createVault') }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" label-width="100px" @submit.prevent="handleOpenVault">
      <!-- 20250127 陈凤庆 简化模式下隐藏密码库文件选择 -->
      <el-form-item v-if="!isSimplifiedMode" :label="$t('login.vaultFile')" prop="vaultPath">
        <div class="file-selector">
          <el-input
            v-model="loginForm.vaultPath"
            :placeholder="$t('login.vaultFilePlaceholder')"
            :disabled="loading"
          />
          <el-button
            type="default"
            :icon="FolderOpened"
            :disabled="loading"
            @click="selectVaultFile"
          >
            {{ $t('common.select') }}
          </el-button>
        </div>
      </el-form-item>

      <el-form-item :label="$t('login.password')" prop="password">
        <el-input
          v-model="loginForm.password"
          type="password"
          :placeholder="$t('login.passwordPlaceholder')"
          :disabled="loading"
          show-password
          @keydown.enter.prevent="handleOpenVault"
        />
      </el-form-item>

      <div class="btn-group">
        <el-button
          type="primary"
          :icon="Unlock"
          :loading="loading"
          @click="handleOpenVault"
          :class="{ 'full-width': isSimplifiedMode }"
        >
          {{ $t('login.login') }}
        </el-button>
        <!-- 20250127 陈凤庆 简化模式下隐藏创建密码库按钮 -->
        <el-button
          v-if="!isSimplifiedMode"
          type="default"
          :icon="Plus"
          :disabled="loading"
          @click="handleCreateVault"
        >
          {{ $t('login.createVault') }}
        </el-button>
      </div>
    </el-form>

    <!-- 20250127 陈凤庆 简化模式下隐藏最近使用的密码库 -->
    <div v-if="!isSimplifiedMode && recentVaults.length > 0" class="recent-vaults">
      <div class="divider"></div>
      <h3 class="recent-title">{{ $t('login.recentVaults') }}</h3>
      <div class="recent-list">
        <div
          v-for="vault in recentVaults"
          :key="vault"
          class="recent-item"
          @click="selectRecentVault(vault)"
        >
          <el-icon class="recent-icon"><Document /></el-icon>
          <span class="recent-path">{{ vault }}</span>
        </div>
      </div>
      <div class="divider"></div>
    </div>

    <!-- 20250127 陈凤庆 简化模式下隐藏密码提示信息 -->
    <div v-if="!isSimplifiedMode" class="password-hint">
      <p>• {{ $t('login.passwordHint1') }}</p>
      <p>• {{ $t('login.passwordHint2') }}</p>
      <p>• {{ $t('login.passwordHint3') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { FolderOpened, Unlock, Plus, Document, ArrowLeft } from '@element-plus/icons-vue'
import { useVaultStore } from '@/stores/vault'
import { SelectVaultFile } from '../../wailsjs/go/app/App'

// 20250928 陈凤庆 暂时回退到直接导入API服务，测试onMounted执行
import { apiService } from '@/services/api'

/**
 * 登录视图组件
 * @author 陈凤庆
 * @description 密码库登录界面，支持打开现有密码库或创建新密码库
 */

const { t } = useI18n()
const router = useRouter()
const vaultStore = useVaultStore()

// 响应式数据
const loginFormRef = ref()
const loading = ref(false)
const recentVaults = ref([])
// 20250127 陈凤庆 添加简化模式状态
// 20250928 陈凤庆 临时直接启用简化模式，绕过onMounted问题
const isSimplifiedMode = ref(true)
// 20250127 陈凤庆 存储最近使用的有效密码库文件路径
const validRecentVaultPath = ref('')
// 20251002 陈凤庆 添加应用标题
const appTitle = ref(t('main.title'))
// 20251017 陈凤庆 添加返回按钮显示状态
const showBackButton = ref(false)

const loginForm = reactive({
  vaultPath: '',
  password: ''
})

// 表单验证规则
const loginRules = {
  vaultPath: [
    { required: true, message: t('login.vaultPathRequired'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.passwordRequired'), trigger: 'blur' },
    { min: 6, message: t('login.passwordMinLength'), trigger: 'blur' }
  ]
}

/**
 * 检查最近使用的密码库状态
 * @author 陈凤庆
 * @description 20250928 检查最近使用的密码库文件是否存在，决定是否使用简化模式
 */
const checkRecentVaultStatus = async () => {
  try {
    console.log('开始检查最近使用的密码库状态...')
    
    // 使用新的API检查最近使用的密码库状态
    const status = await apiService.checkRecentVaultStatus()
    console.log('获取到的密码库状态:', status)
    
    if (status.hasValidVault) {
      // 有有效的最近使用密码库，启用简化模式
      console.log('有有效的最近使用密码库，启用简化模式')
      isSimplifiedMode.value = true
      validRecentVaultPath.value = status.vaultPath
      loginForm.vaultPath = status.vaultPath
      console.log('设置密码库路径为:', status.vaultPath)
    } else {
      // 没有有效的最近使用密码库，使用完整模式
      console.log('没有有效的最近使用密码库，使用完整模式')
      isSimplifiedMode.value = false
    }
    
    console.log('isSimplifiedMode设置为:', isSimplifiedMode.value)
  } catch (error) {
    console.error('检查最近使用的密码库状态失败:', error)
    // 出错时使用完整模式
    isSimplifiedMode.value = false
  }
}

/**
 * 组件挂载时执行
 */
onMounted(async () => {
  console.log('=== LoginView onMounted 开始执行 ===')

  try {
    console.log('apiService:', apiService)
    console.log('apiService类型:', typeof apiService)

    // 20251003 陈凤庆 检查URL参数，如果有mode=full参数，强制使用完整模式
    const urlParams = new URLSearchParams(window.location.search)
    const forceFullMode = urlParams.get('mode') === 'full'

    if (forceFullMode) {
      console.log('检测到mode=full参数，强制使用完整模式')
      isSimplifiedMode.value = false
      validRecentVaultPath.value = ''
      loginForm.vaultPath = ''
      loginForm.password = ''
    }

    // 20251002 陈凤庆 获取应用信息并设置标题
    try {
      const appInfo = await apiService.getAppInfo()
      appTitle.value = `${appInfo.name} v${appInfo.version}`
      // appTitle.value = `${appInfo.name}`
      console.log('登录界面标题已设置:', appTitle.value)
    } catch (error) {
      console.error('获取应用信息失败:', error)
      appTitle.value = t('main.title')
    }

    // 20250128 陈凤庆 修正简化模式逻辑，只有在有有效最近密码库时才启用简化模式
    // 20251003 陈凤庆 如果不是强制完整模式，才检查最近密码库状态
    if (!forceFullMode && apiService) {
      console.log('✅ apiService注入成功，开始检查最近使用的密码库状态')
      console.log('准备调用checkRecentVaultStatus...')
      await checkRecentVaultStatus()
      console.log('checkRecentVaultStatus执行完成')
    } else if (!forceFullMode) {
      console.log('❌ apiService注入失败，使用完整模式')
      // 20250128 陈凤庆 修正：当apiService不可用时，使用完整模式
      console.log('使用完整模式（备用逻辑）')
      isSimplifiedMode.value = false
    }
  } catch (error) {
    console.error('onMounted执行失败:', error)
    // 20250128 陈凤庆 修正：出错时使用完整模式
    console.log('出错时使用完整模式')
    isSimplifiedMode.value = false
  }

  console.log('=== LoginView onMounted 执行完成 ===')
})

/**
 * 选择密码库文件
 * @author 陈凤庆
 * @date 2025-01-27
 * @modify 陈凤庆 移除Wails环境检查，直接调用API，让后端处理错误
 */
const selectVaultFile = async () => {
  // 20251001 陈凤庆 直接使用Wails文件选择对话框API，不做前端环境检查
  // 文件选择只是为了获取路径，后续由后端验证文件是否存在
  try {
    const selectedPath = await SelectVaultFile()
    if (selectedPath) {
      loginForm.vaultPath = selectedPath
      ElMessage.success(t('common.success'))
    }
  } catch (error) {
    console.error('文件选择失败:', error)
    ElMessage.error(t('login.selectFileFailed'))
  }
}

/**
 * 选择最近使用的密码库
 * @param {string} vaultPath 密码库路径
 */
const selectRecentVault = (vaultPath) => {
  loginForm.vaultPath = vaultPath
}

/**
 * 检测当前运行平台
 * @returns {string} 平台名称：'windows', 'macos', 'linux', 'unknown'
 */
const detectPlatform = () => {
  if (typeof window !== "undefined") {
    const userAgent = window.navigator.userAgent.toLowerCase();
    if (userAgent.includes("win")) return "windows";
    if (userAgent.includes("mac")) return "macos";
    if (userAgent.includes("linux")) return "linux";
  }
  return "unknown";
}

/**
 * 验证强密码
 * @param {string} password 密码
 * @returns {object} 验证结果
 * @author 20251003 陈凤庆
 */
const validateStrongPassword = (password) => {
  if (!password) {
    return { isValid: false, message: t('login.passwordEmpty') }
  }

  if (password.length < 8) {
    return { isValid: false, message: t('login.passwordMinLength8') }
  }

  // 检查是否包含大写字母
  if (!/[A-Z]/.test(password)) {
    return { isValid: false, message: t('login.passwordNeedUppercase') }
  }

  // 检查是否包含小写字母
  if (!/[a-z]/.test(password)) {
    return { isValid: false, message: t('login.passwordNeedLowercase') }
  }

  // 检查是否包含数字
  if (!/[0-9]/.test(password)) {
    return { isValid: false, message: t('login.passwordNeedNumber') }
  }

  // 检查是否包含特殊字符
  if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~`]/.test(password)) {
    return { isValid: false, message: t('login.passwordNeedSpecialChar') }
  }

  return { isValid: true, message: t('login.passwordStrengthOk') }
}

/**
 * 处理打开密码库
 */
const handleOpenVault = async () => {
  if (!loginFormRef.value) return
  
  try {
    // 20250127 陈凤庆 简化模式下不需要验证文件路径，因为已经预设了
    if (!isSimplifiedMode.value) {
      // 验证表单
      await loginFormRef.value.validate()
    } else {
      // 简化模式下只验证密码
      if (!loginForm.password) {
        ElMessage.error(t('login.passwordRequired'))
        return
      }
      // 确保使用有效的最近密码库路径
      if (validRecentVaultPath.value) {
        loginForm.vaultPath = validRecentVaultPath.value
      }
    }
    
    loading.value = true
    
    // 检查密码库文件是否存在
    const exists = await apiService.checkVaultExists(loginForm.vaultPath)
    if (!exists) {
      ElMessage.error(t('login.vaultFileNotExists'))
      return
    }
    
    // 打开密码库
    console.log('[登录] 正在调用后端API打开密码库...')
    console.log('[登录] 密码库路径:', loginForm.vaultPath)
    console.log('[登录] 密码长度:', loginForm.password.length)

    await apiService.openVault(loginForm.vaultPath, loginForm.password)

    console.log('[登录] 🎉 后端API调用成功，密码库已打开')

    // 更新前端状态
    console.log('[登录] 正在更新前端store状态...')
    vaultStore.setVaultOpened(true, loginForm.vaultPath)

    console.log('[登录] ✅ 前端状态已更新，isVaultOpened =', vaultStore.isVaultOpened)
    console.log('[登录] 当前密码库路径:', vaultStore.currentVaultPath)

    // 20251004 陈凤庆 登录成功后获取锁定配置信息
    console.log('[登录] 正在获取锁定配置信息...')
    try {
      const lockConfig = await apiService.getLockConfig()
      console.log('[登录] ✅ 锁定配置获取成功:', lockConfig)

      // 将锁定配置存储到localStorage，供主界面使用
      localStorage.setItem('lock_config', JSON.stringify(lockConfig))
      console.log('[登录] 锁定配置已存储到localStorage')
    } catch (error) {
      console.error('[登录] ❌ 获取锁定配置失败:', error)
      // 不阻断登录流程，使用默认配置
    }

    // 20251003 陈凤庆 检测平台，Windows需要特殊处理
    const platform = detectPlatform()
    console.log('[登录] 检测到平台:', platform)

    console.log('[登录] 准备跳转到主界面')

    if (platform === 'windows') {
      // Windows平台使用Hash模式，需要确保状态完全同步后再跳转
      console.log('[登录] Windows平台，等待状态同步...')

      // 使用Vue的nextTick确保响应式更新完成
      await nextTick()

      // 再次确认状态
      console.log('[登录] Windows状态确认，isVaultOpened =', vaultStore.isVaultOpened)

      // 添加额外延迟确保状态稳定
      await new Promise(resolve => setTimeout(resolve, 100))

      console.log('[登录] Windows平台跳转到主界面')
      await router.push('/main')
    } else {
      // macOS/Linux平台正常跳转
      console.log('[登录] macOS/Linux平台跳转到主界面')
      await router.push('/main')
    }

    console.log('[登录] 路由跳转完成')

    // 20251005 陈凤庆 修复：只有在成功完成所有操作后才显示成功消息
    ElMessage.success(t('login.loginSuccess'))

  } catch (error) {
    console.error('打开密码库失败:', error)
    ElMessage.error(error.message || t('login.loginFailed'))
  } finally {
    loading.value = false
  }
}

/**
 * 处理创建密码库
 */
const handleCreateVault = async () => {
  try {
    // 提示用户输入新密码库名称
    const { value: vaultName } = await ElMessageBox.prompt(
      t('login.createVaultPrompt'),
      t('login.createVault'),
      {
        confirmButtonText: t('common.next'),
        cancelButtonText: t('common.cancel'),
        inputPlaceholder: t('login.vaultNamePlaceholder')
      }
    )

    if (!vaultName) return

    const { value: password } = await ElMessageBox.prompt(
      t('login.setPasswordPrompt'),
      t('login.setPassword'),
      {
        confirmButtonText: t('common.create'),
        cancelButtonText: t('common.cancel'),
        inputType: 'password',
        inputPlaceholder: t('login.passwordPlaceholder'),
        dangerouslyUseHTMLString: true,
        message: `
          <div style="text-align: left; line-height: 1.6; width: 100%; padding: 5px;">
            <p style="margin: 0 0 12px 0; font-size: 14px; color: #333;">${t('login.setPassword')}</p>
            <div style="background: #f8f9fa; padding: 12px; border-radius: 6px; border-left: 4px solid #409eff; margin-top: 8px;width: 100%;">
              <p style="margin: 0 0 8px 0; font-weight: 600; color: #333; font-size: 13px;">${t('login.passwordRequirements')}：</p>
              <ul style="margin: 0; padding-left: 16px; color: #666; font-size: 12px; list-style-type: disc;">
                <li style="margin-bottom: 4px;">${t('login.passwordReq1')}</li>
                <li style="margin-bottom: 4px;">${t('login.passwordReq2')}</li>
                <li style="margin-bottom: 4px;">${t('login.passwordReq3')}</li>
                <li style="margin-bottom: 4px;">${t('login.passwordReq4')}</li>
                <li style="margin-bottom: 0;">${t('login.passwordReq5')}</li>
              </ul>
            </div>
          </div>
        `
      }
    )

    // 20251003 陈凤庆 强密码验证：大小写+数字+特殊字符，至少8位
    if (!password) {
      return
    }

    const passwordValidation = validateStrongPassword(password)
    if (!passwordValidation.isValid) {
      ElMessage.error(passwordValidation.message)
      return
    }

    loading.value = true

    // 创建密码库，传递密码库名称和当前语言
    const currentLanguage = localStorage.getItem('app-language') || 'zh-CN'
    const createdVaultPath = await apiService.createVault(vaultName, password, currentLanguage)

    // 更新表单
    loginForm.vaultPath = createdVaultPath
    loginForm.password = password
    
    // 更新状态
    vaultStore.setVaultOpened(true, createdVaultPath)

    ElMessage.success(t('login.createVaultSuccess'))
    
    // 跳转到主界面
    router.push('/main')
    
  } catch (error) {
    console.error('创建密码库失败:', error)
    ElMessage.error(error.message || t('login.createVaultFailed'))
  } finally {
    loading.value = false
  }
}

/**
 * 处理更多菜单命令
 * @param {string} command 命令类型
 * @author 陈凤庆
 * @description 20250128 处理简化模式下更多菜单的命令，包括选择密码库和创建密码库
 */
const handleMoreCommand = async (command) => {
  try {
    switch (command) {
      case 'select-vault':
        // 切换到完整模式，允许用户选择密码库
        isSimplifiedMode.value = false
        showBackButton.value = true // 显示返回按钮
        validRecentVaultPath.value = ''

        // 尝试获取当前密码库路径并显示在输入框中
        try {
          const currentVaultPath = await apiService.getCurrentVaultPath()
          if (currentVaultPath) {
            loginForm.vaultPath = currentVaultPath
            console.log('[选择密码库] 已加载当前密码库路径:', currentVaultPath)
          } else {
            loginForm.vaultPath = ''
          }
        } catch (error) {
          console.warn('[选择密码库] 获取当前密码库路径失败:', error)
          loginForm.vaultPath = ''
        }

        loginForm.password = ''
        ElMessage.info(t('login.switchToFullMode'))
        break

      case 'create-vault':
        // 切换到完整模式并触发创建密码库流程
        isSimplifiedMode.value = false
        showBackButton.value = true // 显示返回按钮
        validRecentVaultPath.value = ''
        loginForm.vaultPath = ''
        loginForm.password = ''
        ElMessage.info(t('login.switchToCreateMode'))
        // 延迟一下再调用创建密码库，确保界面已经切换
        setTimeout(() => {
          handleCreateVault()
        }, 100)
        break

      default:
        console.warn('未知的更多菜单命令:', command)
    }
  } catch (error) {
    console.error('处理更多菜单命令失败:', error)
    ElMessage.error(t('error.operationFailed'))
  }
}

/**
 * 返回到简化模式
 * @author 陈凤庆
 * @date 2025-10-17
 * @description 从完整模式返回到简化模式
 */
const handleBackToSimplified = () => {
  isSimplifiedMode.value = true
  showBackButton.value = false
  loginForm.vaultPath = ''
  loginForm.password = ''
  ElMessage.info(t('login.backToSimplifiedMode'))
}
</script>

<style scoped>
.login-container {
  width: 100vw;
  height: 100vh; /* 20250928 陈凤庆 高度自适应内容 */
  background: #fff;
  /* 20250928 陈凤庆 完全取消边框、阴影和圆角，实现无边框设计 */
  padding: 20px 30px; 
  margin: 20px auto; 
  /* 20250928 陈凤庆 根据内容动态调整最小高度，确保界面紧凑 */
  min-height: fit-content;
  /* 20250928 陈凤庆 确保容器能够根据内容自动调整高度 */
  display: flex;
  flex-direction: column;
  /* 20250928 陈凤庆 添加过渡动画，使高度变化更平滑 */
  transition: height 0.3s ease-in-out;
}

/* 20250928 陈凤庆 简化模式下进一步优化间距，让界面更紧凑 */
.login-container.simplified-mode {
  padding: 10px 30px; /* 进一步减少上下内边距 */
  margin: 5px auto; /* 进一步减少上下外边距 */
  /* 20250928 陈凤庆 简化模式下设置最大高度，确保界面紧凑 */
  max-height: 300px;
}

.login-container.simplified-mode .login-header {
  margin-bottom: 10px; /* 进一步减少头部下边距 */
}

.login-container.simplified-mode .login-title {
  margin-bottom: 5px; /* 减少标题下边距 */
  font-size: 20px; /* 稍微减小标题字体 */
}

.login-container.simplified-mode .el-form-item {
  margin-bottom: 10px; /* 进一步减少表单项间距 */
}

.login-container.simplified-mode .btn-group {
  margin-top: 10px; /* 进一步减少按钮组上边距 */
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
  position: relative; /* 20251017 陈凤庆 为返回按钮定位 */
}

/* 20251017 陈凤庆 返回按钮样式 */
.back-button-container {
  position: absolute;
  left: 0;
  top: 0;
}

.back-button {
  color: #409eff;
  font-size: 14px;
  padding: 8px 12px;
}

.back-button:hover {
  color: #66b1ff;
}

.login-title {
  font-size: 22px;
  color: #333;
  margin-bottom: 8px;
}

.login-desc {
  font-size: 14px;
  color: #666;
}

/* 20250128 陈凤庆 更多菜单容器样式 */
.more-menu-container {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 10;
}

/* 20250128 陈凤庆 更多按钮样式 */
/* 20250928 陈凤庆 调整为三行横杆图标样式 */
.more-button {
  width: 32px !important;
  height: 32px !important;
  padding: 6px !important;
  border-radius: 4px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.more-button:hover {
  background-color: #f5f7fa;
}

/* 20250928 陈凤庆 三行横杆图标样式 */
.hamburger-icon {
  width: 18px;
  height: 14px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
}

.hamburger-icon span {
  width: 100%;
  height: 2px;
  background-color: #606266;
  border-radius: 1px;
  transition: all 0.3s;
}

.more-button:hover .hamburger-icon span {
  background-color: #409eff;
}

.file-selector {
  display: flex;
  gap: 10px;
  align-items: center;
  /* 20250928 陈凤庆 确保文件选择器总宽度与密码输入框宽度一致 */
  width: 100%;
}

.file-selector .el-input {
  /* 20250928 陈凤庆 输入框左对齐，自动适应宽度，减去选择按钮宽度和间距 */
  width: calc(100% - 80px - 10px); /* 总宽度减去按钮宽度80px和间距10px */
  min-width: 0;
}

.file-selector .el-button {
  /* 20250928 陈凤庆 选择按钮固定宽度80px，右对齐 */
  width: 80px;
  flex-shrink: 0;
}

.btn-group {
  display: flex;
  gap: 15px; /* 20250928 陈凤庆 按钮间距 */
  margin-top: 20px;
  /* 20250928 陈凤庆 确保按钮组与输入框对齐，左对齐到表单标签后 */
  margin-left: 100px;
}

/* 20250928 陈凤庆 简化模式下按钮组样式优化 */
.login-container.simplified-mode .btn-group {
  margin-left: 0;
  /* 20250928 陈凤庆 简化模式下按钮组与密码输入框对齐 */
  padding-left: 100px;
}

.btn-group .el-button {
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  /* 20250928 陈凤庆 按钮宽度自动计算，总宽度等于密码输入框宽度，减去中间间距 */
  width: calc(50% - 7.5px); /* 每个按钮宽度为50%减去一半间距 */
}

/* 20250127 陈凤庆 简化模式下按钮全宽样式 */
.btn-group .el-button.full-width {
  width: 100%;
  flex: none;
}

.recent-vaults {
  /* 20250928 陈凤庆 减少上边距，让界面更紧凑 */
  margin-top: 15px;
  /* 20250928 陈凤庆 确保在没有数据时完全不显示 */
  display: block;
}

.recent-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px; /* 20250928 陈凤庆 减少下边距 */
  margin-top: 8px; /* 20250928 陈凤庆 减少上边距 */
}

.recent-list {
  /* 20250928 陈凤庆 移除最大高度限制和滚动条，让内容自然显示 */
  overflow: visible;
  /* 20250928 陈凤庆 确保列表在没有内容时不占用额外空间 */
  min-height: 0;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.recent-item:hover {
  background-color: #f5f7fa;
}

.recent-icon {
  color: #909399;
  font-size: 16px;
}

.recent-path {
  font-size: 12px;
  color: #606266;
  word-break: break-all;
}

.divider {
  height: 1px;
  background: #eee;
  margin: 15px 0; /* 20250928 陈凤庆 减少分割线上下边距 */
}

.password-hint {
  font-size: 12px;
  color: #888;
  line-height: 1.6;
  margin-top: 10px; /* 20250928 陈凤庆 减少密码提示上边距 */
}

.password-hint p {
  margin-bottom: 4px;
}
</style>