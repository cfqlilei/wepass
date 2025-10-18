<template>
  <el-dialog
    v-model="visible"
    :title="t('settings.title')"
    :width="dialogWidth"
    top="5vh"
    :before-close="handleClose"
    :modal="true"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    :destroy-on-close="false"
    :append-to-body="true"
    :lock-scroll="true"
    class="settings-dialog"
    draggable
  >
    <div class="settings-container">
      <!-- 左侧侧边栏 -->
      <div class="settings-sidebar">
        <el-menu
          :key="generalSettings.language"
          :default-active="activeTab"
          @select="handleTabSelect"
          class="settings-menu"
        >
          <el-menu-item index="log">
            <el-icon><Document /></el-icon>
            <span>{{ t('settings.log') }}</span>
          </el-menu-item>
          <el-menu-item index="general">
            <el-icon><Setting /></el-icon>
            <span>{{ t('settings.general') }}</span>
          </el-menu-item>
          <el-menu-item index="lock">
            <el-icon><Lock /></el-icon>
            <span>{{ t('settings.lock') }}</span>
          </el-menu-item>
          <el-menu-item index="hotkey">
            <el-icon><Operation /></el-icon>
            <span>{{ t('settings.hotkey') }}</span>
          </el-menu-item>
        </el-menu>
      </div>

      <!-- 右侧内容区域 -->
      <div class="settings-content">
        <!-- 日志设置 -->
        <div v-if="activeTab === 'log'" class="settings-panel">
          <h3>{{ t('settings.logSettings') }}</h3>
          <div class="setting-item">
            <el-checkbox
              v-model="logSettings.enableInfoLog"
              @change="handleLogSettingChange"
            >
              {{ t('settings.recordInfoLog') }}
            </el-checkbox>
            <p class="setting-description">
              {{ t('settings.infoLogDescription') }}
            </p>
          </div>
          <div class="setting-item">
            <el-checkbox
              v-model="logSettings.enableDebugLog"
              @change="handleLogSettingChange"
            >
              {{ t('settings.recordDebugLog') }}
            </el-checkbox>
            <p class="setting-description">
              {{ t('settings.debugLogDescription') }}
            </p>
          </div>
        </div>

        <!-- 通用设置 -->
        <div v-if="activeTab === 'general'" class="settings-panel">
          <h3>{{ t('settings.general') }}</h3>
          <div class="setting-item">
            <label>{{ t('settings.theme') }}</label>
            <el-select v-model="generalSettings.theme" :placeholder="t('settings.selectTheme')">
              <el-option :label="t('settings.lightTheme')" value="light" />
              <el-option :label="t('settings.darkTheme')" value="dark" />
            </el-select>
          </div>
          <div class="setting-item">
            <label>{{ t('settings.language') }}</label>
            <el-select v-model="generalSettings.language" :placeholder="t('settings.selectLanguage')">
              <el-option :label="t('settings.simplifiedChinese')" value="zh-CN" />
              <el-option :label="t('settings.traditionalChinese')" value="zh-TW" />
              <el-option :label="t('settings.english')" value="en-US" />
              <el-option :label="t('settings.russian')" value="ru-RU" />
              <el-option :label="t('settings.spanish')" value="es-ES" />
              <el-option :label="t('settings.french')" value="fr-FR" />
              <el-option :label="t('settings.german')" value="de-DE" />
            </el-select>
          </div>
        </div>

        <!-- 锁定设置 -->
        <div v-if="activeTab === 'lock'" class="settings-panel">
          <h3>{{ t('settings.lockSettings') }}</h3>

          <!-- 自动锁定开关 -->
          <div class="setting-item">
            <el-checkbox
              v-model="lockSettings.enableAutoLock"
              @change="handleLockSettingChange"
            >
              {{ t('settings.enableAutoLock') }}
            </el-checkbox>
            <p class="setting-description">
              {{ t('settings.autoLockDescription') }}
            </p>
          </div>

          <!-- 锁定方式选择 -->
          <div v-if="lockSettings.enableAutoLock" class="setting-item">
            <label>{{ t('settings.lockMethod') }}</label>
            <div class="lock-type-options">
              <el-checkbox
                v-model="lockSettings.enableTimerLock"
                @change="handleLockSettingChange"
              >
                {{ t('settings.timerLock') }}
              </el-checkbox>
              <el-checkbox
                v-model="lockSettings.enableMinimizeLock"
                @change="handleLockSettingChange"
              >
                {{ t('settings.minimizeLock') }}
              </el-checkbox>
            </div>
            <p class="setting-description">
              {{ t('settings.lockMethodDescription') }}
            </p>
          </div>

          <!-- 定时锁定时间设置 -->
          <div v-if="lockSettings.enableAutoLock && lockSettings.enableTimerLock" class="setting-item">
            <label>{{ t('settings.lockTime') }}</label>
            <div class="time-input-group">
              <el-input-number
                v-model="lockSettings.lockTimeMinutes"
                :min="1"
                :max="1440"
                :step="1"
                @change="handleLockSettingChange"
              />
              <span class="time-unit">{{ t('settings.minutesAfterLock') }}</span>
            </div>
            <p class="setting-description">
              {{ t('settings.lockTimeDescription') }}
            </p>
          </div>

          <!-- 系统固定锁定提示 -->
          <div class="setting-item system-lock-info">
            <el-alert
              :title="t('settings.systemFixedLock')"
              type="info"
              :closable="false"
              show-icon
            >
              <template #default>
                <p>{{ t('settings.systemLockDescription') }}</p>
              </template>
            </el-alert>
          </div>
        </div>

        <!-- 快捷键设置 -->
        <div v-if="activeTab === 'hotkey'" class="settings-panel">
          <h3>{{ t('settings.hotkey') }}</h3>

          <!-- 启用全局快捷键 -->
          <div class="setting-item">
            <div class="setting-row">
              <el-switch
                v-model="hotkeySettings.enableGlobalHotkey"
                @change="handleHotkeySettingChange"
              />
              <label>{{ t('settings.enableGlobalHotkey') }}</label>
            </div>
            <p class="setting-description">
              {{ t('settings.enableGlobalHotkeyDescription') }}
            </p>
          </div>

          <!-- 显示/隐藏快捷键 -->
          <div v-if="hotkeySettings.enableGlobalHotkey" class="setting-item">
            <div class="setting-row">
              <el-switch
                v-model="hotkeySettings.enableShowHideHotkey"
                @change="handleHotkeySettingChange"
              />
              <label>{{ t('settings.enableShowHideHotkey') }}</label>
            </div>
            <div v-if="hotkeySettings.enableShowHideHotkey" class="hotkey-config">
              <label>{{ t('settings.showHideHotkey') }}</label>
              <div class="hotkey-input-group">
                <el-input
                  v-model="hotkeySettings.showHideHotkey"
                  :placeholder="t('settings.showHideHotkeyPlaceholder')"
                  @keydown="handleHotkeyInput"
                  @blur="handleHotkeySettingChange"
                  readonly
                  class="hotkey-input"
                />
                <el-button @click="clearHotkey" size="small">
                  {{ t('common.clear') }}
                </el-button>
              </div>
              <p class="setting-description">
                {{ t('settings.showHideHotkeyDescription') }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script>
/**
 * 设置对话框组件
 * @author 陈凤庆
 * @date 2025-10-03
 * @description 应用设置界面，包含日志设置等配置选项
 */
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Setting, Lock, Operation } from '@element-plus/icons-vue'
import { apiService } from '@/services/api'
import { useI18n } from 'vue-i18n'
import { switchLanguage } from '@/utils/language'

export default {
  name: 'SettingsDialog',
  props: {
    modelValue: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const { t } = useI18n()
    const activeTab = ref('log')

    // 日志设置
    const logSettings = ref({
      enableInfoLog: true,
      enableDebugLog: true
    })

    // 通用设置
    const generalSettings = ref({
      theme: 'light',
      language: 'zh-CN'
    })

    // 锁定设置
    const lockSettings = ref({
      enableAutoLock: false,
      enableTimerLock: false,
      enableMinimizeLock: false,
      lockTimeMinutes: 10,
      enableSystemLock: true,
      systemLockMinutes: 120
    })

    // 快捷键设置
    const hotkeySettings = ref({
      enableGlobalHotkey: true,
      showHideHotkey: 'Ctrl+Alt+H',
      enableShowHideHotkey: true
    })

    // 控制对话框显示
    const visible = computed({
      get: () => props.modelValue,
      set: (value) => emit('update:modelValue', value)
    })

    // 窗口宽度响应式变量
    const windowWidth = ref(window.innerWidth)

    // 响应式对话框宽度
    const dialogWidth = computed(() => {
      // 最大宽度500px
      const maxWidth = 500

      // 左右边距各20px，总共40px
      const margin = 40

      // 计算可用宽度
      const availableWidth = windowWidth.value - margin

      // 返回较小的值，确保不超过最大宽度且保留边距
      const finalWidth = Math.min(maxWidth, availableWidth)

      return `${finalWidth}px`
    })

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

    /**
     * 加载设置
     */
    const loadSettings = async () => {
      try {
        // 加载日志配置
        const logConfig = await apiService.getLogConfig()
        if (logConfig) {
          logSettings.value = {
            enableInfoLog: logConfig.enable_info_log,
            enableDebugLog: logConfig.enable_debug_log
          }
        }

        // 加载通用配置
        const appConfig = await apiService.getAppConfig()
        if (appConfig) {
          generalSettings.value = {
            theme: appConfig.theme || 'light',
            language: appConfig.language || 'zh-CN'
          }
        }

        // 加载锁定配置
        const lockConfig = await apiService.getLockConfig()
        if (lockConfig) {
          lockSettings.value = {
            enableAutoLock: lockConfig.enable_auto_lock || false,
            enableTimerLock: lockConfig.enable_timer_lock || false,
            enableMinimizeLock: lockConfig.enable_minimize_lock || false,
            lockTimeMinutes: lockConfig.lock_time_minutes || 10,
            enableSystemLock: lockConfig.enable_system_lock || true,
            systemLockMinutes: lockConfig.system_lock_minutes || 120
          }
        }

        // 加载快捷键配置
        const hotkeyConfig = await apiService.getHotkeyConfig()
        if (hotkeyConfig) {
          hotkeySettings.value = {
            enableGlobalHotkey: hotkeyConfig.enable_global_hotkey !== false,
            showHideHotkey: hotkeyConfig.show_hide_hotkey || 'Ctrl+Alt+H',
            enableShowHideHotkey: hotkeyConfig.enable_show_hide_hotkey !== false
          }
        }
      } catch (error) {
        console.error('Failed to load settings:', error)
        ElMessage.error(t('settings.loadSettingsFailed'))
      }
    }

    /**
     * 处理标签页切换
     */
    const handleTabSelect = (key) => {
      activeTab.value = key
    }

    /**
     * 处理日志设置变化
     */
    const handleLogSettingChange = async () => {
      // 实时保存日志设置，不显示提示
      try {
        await saveLogSettings()
      } catch (error) {
        console.error('Failed to save log settings:', error)
        ElMessage.error(t('settings.saveSettingsFailed'))
      }
    }

    /**
     * 处理锁定设置变化
     */
    const handleLockSettingChange = async () => {
      // 实时保存锁定设置，不显示提示
      try {
        await saveLockSettings()
      } catch (error) {
        console.error('Failed to save lock settings:', error)
        ElMessage.error(t('settings.saveSettingsFailed'))
      }
    }

    /**
     * 保存日志设置
     */
    const saveLogSettings = async () => {
      try {
        await apiService.setLogConfig({
          enable_info_log: logSettings.value.enableInfoLog,
          enable_debug_log: logSettings.value.enableDebugLog
        })
        // 不在这里显示成功提示，由统一的保存方法显示
      } catch (error) {
        console.error('Failed to save log settings:', error)
        throw error // 抛出错误让上层处理
      }
    }

    /**
     * 保存锁定设置
     */
    const saveLockSettings = async () => {
      try {
        await apiService.setLockConfig({
          enable_auto_lock: lockSettings.value.enableAutoLock,
          enable_timer_lock: lockSettings.value.enableTimerLock,
          enable_minimize_lock: lockSettings.value.enableMinimizeLock,
          lock_time_minutes: lockSettings.value.lockTimeMinutes,
          enable_system_lock: lockSettings.value.enableSystemLock,
          system_lock_minutes: lockSettings.value.systemLockMinutes
        })
        // 不在这里显示成功提示，由统一的保存方法显示
      } catch (error) {
        console.error('Failed to save lock settings:', error)
        throw error // 抛出错误让上层处理
      }
    }

    /**
     * 保存快捷键设置
     */
    /**
     * 验证快捷键格式是否有效
     */
    const isValidHotkey = (hotkey) => {
      if (!hotkey || hotkey.trim() === '') {
        return true // 空快捷键是有效的（表示禁用）
      }

      const parts = hotkey.split('+')
      if (parts.length < 2) {
        return false // 必须至少有一个修饰键和一个主键
      }

      const modifiers = parts.slice(0, -1)
      const mainKey = parts[parts.length - 1]

      // 检查修饰键是否有效
      const validModifiers = ['Ctrl', 'Alt', 'Shift', 'Cmd', 'Command']
      for (const modifier of modifiers) {
        if (!validModifiers.includes(modifier)) {
          return false
        }
      }

      // 检查主键是否有效（不能是修饰键）
      if (validModifiers.includes(mainKey)) {
        return false
      }

      return true
    }

    const saveHotkeySettings = async () => {
      try {
        // 验证快捷键格式
        if (!isValidHotkey(hotkeySettings.value.showHideHotkey)) {
          ElMessage.error(t('settings.invalidHotkeyFormat'))
          throw new Error(t('settings.invalidHotkeyFormat'))
        }

        await apiService.setHotkeyConfig({
          enable_global_hotkey: hotkeySettings.value.enableGlobalHotkey,
          show_hide_hotkey: hotkeySettings.value.showHideHotkey,
          enable_show_hide_hotkey: hotkeySettings.value.enableShowHideHotkey
        })
        // 不在这里显示成功提示，由统一的保存方法显示
      } catch (error) {
        console.error('Failed to save hotkey settings:', error)
        throw error // 抛出错误让上层处理
      }
    }

    /**
     * 处理快捷键设置变化
     */
    const handleHotkeySettingChange = async () => {
      // 自动保存快捷键设置，不显示提示
      try {
        await saveHotkeySettings()
      } catch (error) {
        console.error('Failed to save hotkey settings:', error)
        ElMessage.error(t('settings.saveSettingsFailed'))
      }
    }

    /**
     * 处理快捷键输入
     */
    const handleHotkeyInput = (event) => {
      event.preventDefault()

      const modifierKeys = []
      let mainKey = ''

      // 检查修饰键
      if (event.ctrlKey || event.metaKey) modifierKeys.push('Ctrl')
      if (event.altKey) modifierKeys.push('Alt')
      if (event.shiftKey) modifierKeys.push('Shift')

      // 获取主键（非修饰键）
      if (!['Control', 'Alt', 'Shift', 'Meta'].includes(event.key)) {
        if (event.key === ' ') {
          mainKey = 'Space'
        } else if (event.key.length === 1) {
          mainKey = event.key.toUpperCase()
        } else {
          // 处理功能键，如F1, F2, Escape等
          mainKey = event.key
        }
      }

      // 必须有修饰键和主键才能构成有效的快捷键
      if (modifierKeys.length > 0 && mainKey) {
        const hotkey = [...modifierKeys, mainKey].join('+')
        hotkeySettings.value.showHideHotkey = hotkey
        handleHotkeySettingChange()
      }
    }

    /**
     * 清除快捷键
     */
    const clearHotkey = () => {
      hotkeySettings.value.showHideHotkey = ''
      handleHotkeySettingChange()
    }

    /**
     * 保存设置
     */
    const handleSave = async () => {
      try {
        // 获取当前语言设置，用于判断是否需要刷新页面
        const currentLanguage = generalSettings.value.language
        const previousLanguage = localStorage.getItem('app-language') || 'zh-CN'

        // 保存日志设置
        await saveLogSettings()

        // 保存锁定设置
        await saveLockSettings()

        // 保存快捷键设置
        await saveHotkeySettings()

        // 保存通用设置
        await apiService.setAppConfig({
          theme: generalSettings.value.theme,
          language: currentLanguage
        })

        // 显式更新 localStorage 中的语言设置
        localStorage.setItem('app-language', currentLanguage)

        // 只显示一次保存成功提示
        ElMessage.success(t('settings.settingsSaved'))

        // 如果语言发生变化，切换语言并刷新页面
        if (currentLanguage !== previousLanguage) {
          console.log('Language changed, attempting to switch language and reload page.');
          switchLanguage(currentLanguage, () => {
            // 延迟刷新页面，确保设置已保存
            setTimeout(() => {
              console.log('Reloading page...');
              window.location.reload()
            }, 500)
          })
        } else {
          console.log('Language not changed, closing dialog.');
          handleClose()
        }
      } catch (error) {
        console.error('Failed to save settings:', error)
        ElMessage.error(t('settings.saveSettingsFailed'))
      }
    }

    /**
     * 关闭对话框
     */
    const handleClose = () => {
      visible.value = false
    }

    // 监听对话框显示状态，加载设置
    watch(visible, (newValue) => {
      console.log('SettingsDialog visible changed:', newValue)
      if (newValue) {
        console.log('Loading settings...')
        loadSettings()
      }
    })

    // 监听props变化
    watch(() => props.modelValue, (newValue) => {
      console.log('SettingsDialog modelValue changed:', newValue)
    })

    return {
      t,
      visible,
      dialogWidth,
      activeTab,
      logSettings,
      generalSettings,
      lockSettings,
      hotkeySettings,
      Document,
      Setting,
      Lock,
      Operation,
      handleTabSelect,
      handleLogSettingChange,
      handleLockSettingChange,
      handleHotkeySettingChange,
      handleHotkeyInput,
      clearHotkey,
      handleSave,
      handleClose
    }
  }
}
</script>

<style scoped>
:deep(.el-dialog) {
  --el-dialog-padding-primary: 0;
  border-radius: 8px;
  box-shadow: 0 12px 32px 4px rgba(0, 0, 0, 0.36), 0 8px 20px rgba(0, 0, 0, 0.72);
  z-index: 9999 !important;
}

:deep(.el-overlay) {
  z-index: 9998 !important;
}

:deep(.el-dialog__header) {
  padding: 20px 20px 0 20px;
  border-bottom: 1px solid #e4e7ed;
  margin-right: 0;
}

:deep(.el-dialog__body) {
  padding: 0;
}

:deep(.el-dialog__headerbtn) {
  top: 15px;
  right: 15px;
}

.settings-container {
  display: flex;
  height: 400px;
  max-height: 70vh;
  min-height: 300px;
  background: #fff;
}

.settings-sidebar {
  width: 120px;
  min-width: 100px;
  border-right: 1px solid #e4e7ed;
  background-color: #fafafa;
}

.settings-menu {
  border: none;
  background-color: transparent;
}

.settings-menu .el-menu-item {
  height: 48px;
  line-height: 48px;
  padding-left: 16px;
  padding-right: 8px;
  font-size: 14px;
}

.settings-menu .el-menu-item .el-icon {
  margin-right: 8px;
}

.settings-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

/* 小屏幕适配 */
@media (max-width: 480px) {
  .settings-sidebar {
    width: 80px;
    min-width: 80px;
  }

  .settings-menu .el-menu-item {
    padding-left: 8px;
    padding-right: 4px;
    font-size: 12px;
  }

  .settings-menu .el-menu-item span {
    display: none;
  }

  .settings-content {
    padding: 12px;
  }

  .settings-container {
    height: 350px;
    min-height: 280px;
  }
}

.settings-panel h3 {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.setting-item {
  margin-bottom: 24px;
}

.setting-item label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
}

.setting-description {
  margin: 8px 0 0 24px;
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
}

/* 锁定设置样式 */
.time-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time-unit {
  font-size: 14px;
  color: #606266;
}

.system-lock-info {
  margin-top: 20px;
}

.system-lock-info :deep(.el-alert) {
  background-color: #f0f9ff;
  border-color: #b3d8ff;
}

.system-lock-info :deep(.el-alert__content) {
  color: #409eff;
}

.lock-type-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

.lock-type-options .el-checkbox {
  margin-right: 0;
}

/* 快捷键设置样式 */
.hotkey-config {
  margin-left: 24px;
  margin-top: 12px;
  padding-left: 16px;
  border-left: 2px solid var(--el-border-color-light);
}

.hotkey-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.hotkey-input {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-weight: bold;
}

.hotkey-input :deep(.el-input__inner) {
  background-color: #f5f7fa;
  border: 2px dashed #dcdfe6;
  text-align: center;
  cursor: pointer;
}

.hotkey-input :deep(.el-input__inner):focus {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
