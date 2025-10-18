<template>
  <div class="search-results">
    <div class="search-header">
      <h3>{{ $t('main.searchResults') }}</h3>
      <el-button type="text" :icon="Close" @click="closeSearch" />
    </div>
    
    <!-- 分组搜索结果 -->
    <div v-if="groupResults.length > 0" class="result-section">
      <h4 class="section-title">
        <el-icon><Folder /></el-icon>
        {{ $t('searchResults.groups') }} ({{ groupResults.length }})
      </h4>
      <div class="result-list">
        <div
          v-for="group in groupResults"
          :key="`group-${group.id}`"
          class="result-item group-item"
          @click="selectGroup(group)"
        >
          <el-icon class="result-icon"><Folder /></el-icon>
          <div class="result-content">
            <div class="result-title">{{ group.name }}</div>
            <div class="result-subtitle">{{ $t('main.group') }}</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 账号搜索结果 -->
    <div v-if="passwordResults.length > 0" class="result-section">
      <h4 class="section-title">
        <el-icon><Key /></el-icon>
        {{ $t('searchResults.accounts') }} ({{ passwordResults.length }})
      </h4>
      <div class="result-list">
        <div
          v-for="password in passwordResults"
          :key="`password-${password.id}`"
          class="result-item password-item"
          @click="selectPassword(password)"
        >
          <el-icon class="result-icon"><component :is="getIconComponent(password.icon)" /></el-icon>
          <div class="result-content">
            <div class="result-title">{{ password.title }}</div>
            <div class="result-subtitle">{{ password.url || $t('searchResults.noUrl') }}</div>
            <div class="result-group">{{ $t('searchResults.belongsToGroup') }}: {{ getGroupName(password.group_id) }}</div>
          </div>
          <div class="result-actions">
            <el-button
              type="text"
              :icon="Key"
              size="small"
              :title="$t('main.inputUsernameAndPassword')"
              @click.stop="inputUsernameAndPassword(password)"
            />
            <el-button
              type="text"
              :icon="User"
              size="small"
              :title="$t('main.inputUsername')"
              @click.stop="inputUsername(password)"
            />
            <el-button
              type="text"
              :icon="Lock"
              size="small"
              :title="$t('main.inputPassword')"
              @click.stop="inputPassword(password)"
            />
          </div>
        </div>
      </div>
    </div>
    
    <!-- 无结果提示 -->
    <div v-if="groupResults.length === 0 && passwordResults.length === 0" class="no-results">
      <el-icon class="no-results-icon"><Search /></el-icon>
      <p>{{ $t('searchResults.noSearchResults') }}</p>
      <p class="no-results-tip">{{ $t('searchResults.tryOtherKeywords') }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import {
  Close, Folder, Key, Link, User, Lock, Search
} from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'

/**
 * 搜索结果组件
 * @author 陈凤庆
 * @description 显示搜索结果，支持分组和账号的搜索
 */

// Props
const props = defineProps({
  groupResults: {
    type: Array,
    default: () => []
  },
  passwordResults: {
    type: Array,
    default: () => []
  }
})

// Emits
const emit = defineEmits(['close', 'selectGroup', 'selectPassword'])

const appStore = useAppStore()
const { t } = useI18n()

// 计算属性
const groups = computed(() => appStore.groups)

/**
 * 获取分组名称
 * @param {number} groupId 分组ID
 * @returns {string} 分组名称
 */
const getGroupName = (groupId) => {
  const group = groups.value.find(g => g.id === groupId)
  return group ? group.name : t('searchResults.unknownGroup')
}

/**
 * 获取图标组件
 * @param {string} iconName 图标名称
 * @returns {Object} 图标组件
 */
const getIconComponent = (iconName) => {
  // 简化的图标映射
  return Key // 默认使用 Key 图标
}

/**
 * 关闭搜索
 */
const closeSearch = () => {
  emit('close')
}

/**
 * 选择分组
 * @param {Object} group 分组对象
 */
const selectGroup = (group) => {
  emit('selectGroup', group)
}

/**
 * 选择账号
 * @param {Object} password 密码对象
 */
const selectPassword = (password) => {
  emit('selectPassword', password)
}

/**
 * 输入用户名和密码
 * @param {Object} password 密码对象
 * @author 陈凤庆
 * @date 2025-01-27
 * @modify 20251002 陈凤庆 添加窗口聚焦管理，确保输入到正确的目标窗口
 * @modify 20251002 陈凤庆 优化提示信息，移除敏感信息显示
 */
const inputUsernameAndPassword = async (password) => {
      try {
        // 检查辅助功能权限
        const hasPermission = await window.go.app.App.CheckAccessibilityPermission()
        if (!hasPermission) {
          ElMessage.error(t('error.accessibilityPermissionRequired'))
          return
        }
        
        // 20250127 陈凤庆 获取当前存储的目标应用程序名称
        const targetAppName = await window.go.app.App.GetPreviousFocusedAppName()
        
        // 模拟输入用户名和密码到其他应用程序
        await window.go.app.App.SimulateUsernameAndPassword(password.username, password.password)
        ElMessage.success(t('success.autofillUsernameAndPassword', { appName: targetAppName }))
      } catch (error) {
        console.error('自动填充失败:', error)
        ElMessage.error(t('error.autofillFailed', { error: error.message || error }))
      }
    }

/**
 * 输入用户名
 * @param {Object} password 密码对象
 * @author 陈凤庆
 * @date 2025-01-27
 * @modify 20251002 陈凤庆 添加窗口聚焦管理，确保输入到正确的目标窗口
 * @modify 20251002 陈凤庆 优化提示信息，移除敏感信息显示
 */
const inputUsername = async (password) => {
  try {
    // 检查辅助功能权限
    const hasPermission = await window.go.app.App.CheckAccessibilityPermission()
    if (!hasPermission) {
      ElMessage.error(t('error.accessibilityPermissionRequired'))
      return
    }
    
    // 20250127 陈凤庆 获取当前存储的目标应用程序名称
    const targetAppName = await window.go.app.App.GetPreviousFocusedAppName()
    
    // 模拟输入用户名到其他应用程序
        await window.go.app.App.SimulateUsername(password.username)
        ElMessage.success(t('success.autofillUsername', { appName: targetAppName }))
  } catch (error) {
    console.error('自动填充用户名失败:', error)
    ElMessage.error(t('error.autofillUsernameFailed', { error: error.message || error }))
  }
}

/**
 * 输入密码
 * @param {Object} password 密码对象
 * @author 陈凤庆
 * @date 2025-01-27
 * @modify 20251002 陈凤庆 添加窗口聚焦管理，确保输入到正确的目标窗口
 * @modify 20251002 陈凤庆 优化提示信息，移除敏感信息显示
 */
const inputPassword = async (password) => {
  try {
    // 检查辅助功能权限
    const hasPermission = await window.go.app.App.CheckAccessibilityPermission()
    if (!hasPermission) {
      ElMessage.error(t('error.accessibilityPermissionRequired'))
      return
    }
    
    // 20250127 陈凤庆 获取当前存储的目标应用程序名称
    const targetAppName = await window.go.app.App.GetPreviousFocusedAppName()
    
    // 模拟输入密码到其他应用程序
        await window.go.app.App.SimulatePassword(password.password)
        ElMessage.success(t('success.autofillPassword', { appName: targetAppName }))
  } catch (error) {
    console.error('自动填充密码失败:', error)
    ElMessage.error(t('error.autofillPasswordFailed', { error: error.message || error }))
  }
}
</script>

<style scoped>
.search-results {
  position: absolute;
  top: 76px; /* 标题栏 + 搜索栏高度 */
  left: 0;
  right: 0;
  bottom: 30px; /* 状态栏高度 */
  background: #ffffff;
  border-top: 1px solid #dee2e6;
  z-index: 1000;
  overflow-y: auto;
}

.search-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.search-header h3 {
  font-size: 16px;
  color: #333;
  margin: 0;
}

.result-section {
  margin-bottom: 20px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #f0f0f0;
  font-size: 14px;
  color: #666;
  margin: 0;
}

.result-list {
  padding: 0 16px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.2s;
}

.result-item:hover {
  background-color: rgba(0, 123, 255, 0.02);
}

.result-item:last-child {
  border-bottom: none;
}

.result-icon {
  font-size: 18px;
  color: #909399;
}

.result-content {
  flex: 1;
}

.result-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.result-subtitle {
  font-size: 12px;
  color: #666;
}

.result-group {
  font-size: 11px;
  color: #999;
  margin-top: 2px;
}

.result-actions {
  display: flex;
  gap: 4px;
}

.result-actions .el-button {
  width: 24px !important;
  height: 24px !important;
  padding: 0 !important;
}

.no-results {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.no-results-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.no-results p {
  margin: 8px 0;
}

.no-results-tip {
  font-size: 12px;
  color: #ccc;
}
</style>