<!--
/**
 * 密码显示对话框组件
 * @author 陈凤庆
 * @date 2025-01-18
 * @description 显示密码的弹窗组件，支持密码分段显示、隐藏/显示切换等功能
 */
-->
<template>
  <el-dialog
    v-model="visible"
    :title="$t('passwordDisplay.title')"
    :width="dialogWidth"
    :before-close="handleClose"
    class="password-display-dialog"
    :close-on-click-modal="false"
  >
    <div class="password-display-content">
      <!-- 账号信息 -->
      <div class="account-info">
        <div class="account-title">
          <el-icon><Key /></el-icon>
          <span>{{ account?.title || '' }}</span>
        </div>
        <div class="account-username">
          <el-icon><User /></el-icon>
          <span>{{ account?.masked_username || account?.username || '' }}</span>
        </div>
      </div>

      <!-- 密码显示区域 -->
      <div class="password-section">
        <div class="password-controls">
          <!-- 隐藏密码开关 -->
          <div class="hide-password-control">
            <el-switch
              v-model="hidePassword"
              :active-text="$t('passwordDisplay.hidePassword')"
              :inactive-text="$t('passwordDisplay.showPassword')"
              size="small"
            />
          </div>
        </div>

        <!-- 密码显示按钮组 -->
        <div class="password-buttons">
          <div class="password-button-group">
            <el-button
              :type="currentPart === 1 ? 'primary' : 'default'"
              @click="showPart(1)"
              class="password-part-button"
            >
              {{ $t('passwordDisplay.firstPart') }}
            </el-button>
            <el-button
              :type="currentPart === 2 ? 'primary' : 'default'"
              @click="showPart(2)"
              class="password-part-button"
            >
              {{ $t('passwordDisplay.secondPart') }}
            </el-button>
          </div>
        </div>

        <!-- 密码显示框 -->
        <div class="password-display-box">
          <el-input
            :value="displayedPassword"
            readonly
            type="text"
            class="password-input"
            :placeholder="$t('passwordDisplay.passwordPlaceholder')"
          >
            <template #append>
              <el-button
                :icon="CopyDocument"
                @click="copyCurrentPassword"
                :title="$t('passwordDisplay.copyPassword')"
              />
            </template>
          </el-input>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">
          {{ $t('common.close') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { Key, User, CopyDocument } from '@element-plus/icons-vue'
import { apiService } from '@/services/api'

/**
 * 密码显示对话框组件
 */

// 国际化
const { t } = useI18n()

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  account: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['update:modelValue'])

// 响应式数据
const hidePassword = ref(false) // 隐藏密码开关
const currentPart = ref(1) // 当前显示的部分：1-第一部分，2-第二部分
const fullPassword = ref('') // 完整密码
const loading = ref(false) // 加载状态

// 窗口宽度响应式变量
const windowWidth = ref(window.innerWidth)

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 响应式对话框宽度
const dialogWidth = computed(() => {
  // 最大宽度500px，最小宽度350px
  const maxWidth = 500
  const minWidth = 350
  
  // 左右边距各30px，总共60px
  const margin = 60
  
  // 计算可用宽度
  const availableWidth = windowWidth.value - margin
  
  // 返回合适的宽度
  const finalWidth = Math.min(maxWidth, Math.max(minWidth, availableWidth))
  
  return `${finalWidth}px`
})

// 密码分段
const passwordParts = computed(() => {
  if (!fullPassword.value) {
    return { first: '', second: '' }
  }
  
  const password = fullPassword.value
  const length = password.length
  const midPoint = Math.ceil(length / 2)
  
  return {
    first: password.substring(0, midPoint),
    second: password.substring(midPoint)
  }
})

// 当前显示的密码
const displayedPassword = computed(() => {
  if (hidePassword.value) {
    return '*****' // 隐藏时显示5个*
  }
  
  if (currentPart.value === 1) {
    return passwordParts.value.first
  } else {
    return passwordParts.value.second
  }
})

/**
 * 处理对话框关闭
 */
const handleClose = () => {
  visible.value = false
  // 清理敏感数据
  fullPassword.value = ''
  currentPart.value = 1
  hidePassword.value = false
}

/**
 * 显示指定部分的密码
 * @param {number} part 部分编号：1-第一部分，2-第二部分
 */
const showPart = (part) => {
  currentPart.value = part
}

/**
 * 复制当前显示的密码
 */
const copyCurrentPassword = async () => {
  try {
    const passwordToCopy = hidePassword.value ? fullPassword.value : displayedPassword.value
    
    if (!passwordToCopy) {
      ElMessage.warning(t('passwordDisplay.noPasswordToCopy'))
      return
    }
    
    await navigator.clipboard.writeText(passwordToCopy)
    ElMessage.success(t('passwordDisplay.passwordCopied'))
    
    // 10秒后清理剪贴板
    setTimeout(async () => {
      try {
        await navigator.clipboard.writeText('')
      } catch (error) {
        console.warn('清理剪贴板失败:', error)
      }
    }, 10000)
  } catch (error) {
    console.error('复制密码失败:', error)
    ElMessage.error(t('passwordDisplay.copyFailed'))
  }
}

/**
 * 获取账号密码
 */
const fetchPassword = async () => {
  if (!props.account?.id) {
    console.error('账号ID为空，无法获取密码')
    return
  }
  
  loading.value = true
  try {
    console.log(`[PasswordDisplayDialog] 获取密码，账号ID: ${props.account.id}`)
    const password = await apiService.getAccountPassword(props.account.id)
    fullPassword.value = password
    console.log(`[PasswordDisplayDialog] 密码获取成功`)
  } catch (error) {
    console.error('[PasswordDisplayDialog] 获取密码失败:', error)
    ElMessage.error(t('passwordDisplay.getPasswordFailed'))
    handleClose()
  } finally {
    loading.value = false
  }
}

// 监听对话框显示状态
watch(visible, (newValue) => {
  if (newValue && props.account) {
    fetchPassword()
  } else if (!newValue) {
    // 对话框关闭时清理数据
    fullPassword.value = ''
    currentPart.value = 1
    hidePassword.value = false
  }
})

// 监听窗口大小变化
onMounted(() => {
  const handleResize = () => {
    windowWidth.value = window.innerWidth
  }
  
  window.addEventListener('resize', handleResize)
  
  // 组件卸载时移除监听器
  return () => {
    window.removeEventListener('resize', handleResize)
  }
})
</script>

<style scoped>
.password-display-dialog {
  --el-dialog-padding-primary: 20px;
}

.password-display-content {
  padding: 0;
}

.account-info {
  margin-bottom: 20px;
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}

.account-title,
.account-username {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
}

.account-title:last-child,
.account-username:last-child {
  margin-bottom: 0;
}

.account-title .el-icon,
.account-username .el-icon {
  margin-right: 8px;
  color: #606266;
}

.account-title {
  font-weight: 600;
  color: #303133;
}

.account-username {
  color: #606266;
}

.password-section {
  margin-top: 20px;
}

.password-controls {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.hide-password-control {
  display: flex;
  align-items: center;
}

.password-buttons {
  margin-bottom: 16px;
}

.password-button-group {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.password-part-button {
  min-width: 120px;
}

.password-display-box {
  margin-top: 16px;
}

.password-input {
  font-family: 'Courier New', monospace;
  font-size: 16px;
}

.password-input :deep(.el-input__inner) {
  text-align: center;
  letter-spacing: 2px;
  font-weight: 600;
}

.dialog-footer {
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .password-button-group {
    flex-direction: column;
    align-items: center;
  }
  
  .password-part-button {
    width: 100%;
    max-width: 200px;
  }
}
</style>
