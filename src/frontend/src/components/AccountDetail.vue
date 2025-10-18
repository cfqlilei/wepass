<template>
  <div>
    <!-- 20251017 陈凤庆 增加对话框宽度以确保用户名输入框按钮有足够空间 -->
    <el-dialog
    v-model="visible"
    :title="isEdit ? (props.isNew ? t('account.createAccount') : t('account.editAccount')) : t('account.accountDetail')"
    width="550px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="isEdit ? formRules : {}"
      label-width="80px"
    >
      <el-form-item :label="t('account.title')" prop="title">
        <el-input
          v-model="formData.title"
          :placeholder="t('account.enterTitle')"
          :readonly="!isEdit"
        />
      </el-form-item>

      <el-form-item :label="t('account.username')" prop="username" class="username-form-item">
        <!-- 20251017 陈凤庆 改为与密码输入框相同的结构 -->
        <el-input
          v-model="formData.username"
          :placeholder="t('account.enterUsername')"
          :readonly="!isEdit"
          ref="usernameInputRef"
          @input="handleUsernameInput"
          @keydown="handleUsernameKeydown"
          @focus="handleUsernameFocus"
          @blur="handleUsernameBlur"
        >
          <template #append>
            <el-button-group>
              <!-- 20251017 陈凤庆 历史用户名下拉按钮始终显示，确保图标左右排列 -->
              <el-dropdown
                v-if="isEdit"
                @command="selectUsername"
                trigger="click"
                placement="bottom-end"
                class="username-history-dropdown"
                ref="usernameDropdownRef"
                :visible="showUsernameDropdown"
                @visible-change="handleDropdownVisibleChange"
                :disabled="filteredUsernameHistory.length === 0"
              >
                <template #default>
                  <el-button
                    :icon="ArrowDown"
                    :title="filteredUsernameHistory.length > 0 ? '历史用户名' : '暂无历史用户名'"
                    :disabled="filteredUsernameHistory.length === 0"
                  />
                </template>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item
                      v-for="(username, index) in filteredUsernameHistory"
                      :key="username"
                      :command="username"
                      :class="{ 'is-selected': selectedUsernameIndex === index }"
                    >
                      {{ username }}
                      <el-button
                        :icon="Delete"
                        size="small"
                        type="danger"
                        text
                        class="delete-username-btn"
                        @click.stop="handleDeleteUsername(username)"
                        :title="'删除用户名: ' + username"
                        style="margin-left: 8px;"
                      />
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-button :icon="CopyDocument" @click="copyUsername" :title="t('account.copyUsername')" />
            </el-button-group>
          </template>
        </el-input>

        <!-- 20251017 陈凤庆 用户名输入框下方显示历史用户名列表 -->
        <!-- 20251017 陈凤庆 优化历史用户名搜索结果显示，移除标题和数量 -->
        <div
          v-if="isEdit && showUsernameList && filteredUsernameHistory.length > 0"
          class="username-suggestions-list"
          ref="usernameSuggestionsRef"
        >
          <div class="suggestions-content">
            <div
              v-for="(username, index) in filteredUsernameHistory.slice(0, 8)"
              :key="username"
              :class="[
                'suggestion-item',
                { 'is-selected': selectedUsernameIndex === index }
              ]"
              @click="selectUsernameFromList(username)"
              @mouseenter="selectedUsernameIndex = index"
            >
              <span class="suggestion-text">{{ username }}</span>
              <el-button
                :icon="Delete"
                size="small"
                type="danger"
                text
                class="delete-suggestion-btn"
                @click.stop="handleDeleteUsername(username)"
                :title="'删除用户名: ' + username"
              />
            </div>
          </div>
        </div>
      </el-form-item>

      <el-form-item :label="t('account.password')" prop="password">
        <el-input
          v-model="formData.password"
          :type="showPassword ? 'text' : 'password'"
          :placeholder="t('account.enterPassword')"
          :readonly="!isEdit"
        >
          <template #append>
            <el-button-group>
              <el-button
                :icon="showPassword ? Hide : View"
                @click="togglePasswordVisibility"
                :title="t('account.togglePassword')"
              />
              <el-button :icon="CopyDocument" @click="copyPassword" :title="t('account.copyPassword')" />
              <el-button
                v-if="isEdit"
                :icon="Refresh"
                @click="generatePasswordWithDefaultRule"
                :title="t('account.generatePassword')"
              />
              <el-dropdown
                v-if="isEdit && passwordRules.length > 0"
                @command="handlePasswordRuleCommand"
                trigger="click"
                placement="bottom-end"
                class="password-rule-dropdown"
              >
                <template #default>
                  <el-button :icon="ArrowDown" :title="t('account.selectPasswordRule')" />
                </template>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item
                      v-for="rule in passwordRules"
                      :key="rule.id"
                      :command="rule.id"
                    >
                      {{ rule.name }}
                      <span v-if="rule.is_default" class="default-rule-tag">（默认）</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-button-group>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item :label="t('account.url')" prop="url">
        <el-input
          v-model="formData.url"
          :placeholder="t('account.enterUrl')"
          :readonly="!isEdit"
        >
          <template #append>
            <el-button :icon="Link" @click="openUrl" :title="t('account.openUrl')" />
          </template>
        </el-input>
      </el-form-item>
      
      <el-form-item :label="t('account.type')" prop="typeid">
        <el-select
          v-model="formData.typeid"
          :placeholder="t('account.selectType')"
          :disabled="!isEdit"
          :loading="tabListLoading"
        >
          <el-option
            v-for="tab in tabList"
            :key="tab.id"
            :label="tab.name"
            :value="tab.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('account.inputMethod')" prop="input_method">
        <el-select
          v-model="formData.input_method"
          :placeholder="t('account.selectInputMethod')"
          :disabled="!isEdit"
        >
          <el-option
            v-for="method in inputMethodOptions"
            :key="method.value"
            :label="method.label"
            :value="method.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('account.notes')">
        <div class="notes-field-container">
          <!-- 备注编辑框 -->
          <div class="notes-input-wrapper">
            <el-input
              v-model="formData.notes"
              type="textarea"
              :rows="3"
              :placeholder="t('account.enterNotes')"
              maxlength="500"
              show-word-limit
              :readonly="!isEdit"
            />
          </div>

          <!-- 备注操作按钮（仅在查看模式下显示） -->
          <div v-if="!isEdit" class="notes-actions">
            <el-button-group>
              <el-button
                :icon="showNotes ? Hide : View"
                @click="toggleNotesVisibility"
                title="显示/隐藏备注"
                size="small"
              />
              <el-button
                :icon="CopyDocument"
                @click="copyNotes"
                title="复制备注"
                size="small"
              />
            </el-button-group>
          </div>
        </div>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <div class="dialog-footer">
        <div class="left-footer-buttons">
          <el-button v-if="!isNew" type="danger" @click="handleDelete">{{ $t('common.delete') }}</el-button>
        </div>
        <div class="right-footer-buttons">
          <el-button @click="handleClose">{{ $t('common.cancel') }}</el-button>
          <el-button v-if="!isEdit" type="primary" @click="enableEdit">{{ $t('common.edit') }}</el-button>
          <el-button v-if="isEdit" type="primary" @click="handleSave" :loading="saving">{{ $t('common.save') }}</el-button>
        </div>
      </div>
    </template>
  </el-dialog>

  <!-- 密码规则选择对话框 -->
  <el-dialog
    v-model="showPasswordRuleDialog"
    :title="t('account.selectPasswordRule')"
    width="400px"
    :before-close="handlePasswordRuleDialogClose"
  >
    <el-form label-width="80px">
      <el-form-item :label="t('account.passwordRule')">
        <el-select
          v-model="selectedPasswordRuleId"
          :placeholder="t('account.selectRule')"
          :loading="passwordRulesLoading"
          style="width: 100%"
        >
          <el-option
            v-for="rule in passwordRules"
            :key="rule.id"
            :label="rule.name"
            :value="rule.id"
          />
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handlePasswordRuleDialogClose">{{ t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          @click="generatePasswordWithRule"
          :disabled="!selectedPasswordRuleId"
        >
          {{ t('account.generatePassword') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Link, CopyDocument, View, Hide, Refresh, ArrowDown, Delete
} from '@element-plus/icons-vue'
import { getAccountPassword, copyAccountPassword } from '@/utils/accountUtils'

const { t } = useI18n()

/**
 * 密码详情组件
 * @author 陈凤庆
 * @description 显示和编辑账号的详细信息
 * @modify 陈凤庆 移除收藏功能，修改编辑框为只读模式以允许文本选择和复制
 * @modify 20250101 陈凤庆 添加标签列表加载功能，支持根据分组ID加载标签
 */

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  accountData: {
    type: Object,
    default: () => ({})
  },
  isNew: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'save', 'delete'])

// 响应式数据
const formRef = ref()
const visible = ref(false)
const isEdit = ref(false)
const showPassword = ref(false)
const showNotes = ref(false) // 20251003 陈凤庆 添加备注显示/隐藏状态
const saving = ref(false)

// 20250101 陈凤庆 添加标签列表相关状态
const tabList = ref([])
const tabListLoading = ref(false)

// 密码规则相关状态
const showPasswordRuleDialog = ref(false)
const passwordRules = ref([])
const passwordRulesLoading = ref(false)
const selectedPasswordRuleId = ref('')
const defaultPasswordRuleId = ref('')

// 20251017 陈凤庆 用户名历史记录相关状态
const usernameHistory = ref([])
const usernameHistoryLoading = ref(false)
const usernameInputRef = ref()
const usernameDropdownRef = ref()
const usernameSuggestionsRef = ref()
const showUsernameDropdown = ref(false)
const showUsernameList = ref(false) // 20251017 陈凤庆 控制输入框下方历史用户名列表显示
const selectedUsernameIndex = ref(-1)
const usernameSearchQuery = ref('')

const formData = reactive({
  id: '', // 20251001 陈凤庆 改为字符串类型，使用GUID
  title: '',
  username: '',
  password: '',
  url: '',
  type: '',
  notes: '',
  group_id: '', // 20251001 陈凤庆 改为字符串类型，使用GUID
  typeid: '', // 20251002 陈凤庆 改为typeid，与后端模型保持一致
  icon: '',
  input_method: 1 // 20251003 陈凤庆 添加输入方式字段，默认为1（默认方式）
})

// 20251003 陈凤庆 输入方式选项
// 20251005 陈凤庆 添加第4种输入方式：键盘助手输入（原第4种底层键盘API已删除）
// 20251005 陈凤庆 添加第5种输入方式：远程输入（专为ToDesk等远程桌面环境设计）
const inputMethodOptions = [
  { value: 1, label: '默认方式 (发送按键码，速度快)' },
  { value: 2, label: '模拟键盘输入 (逐键输入，兼容性好)' },
  { value: 3, label: '复制粘贴输入 (先复制再粘贴)' },
  { value: 4, label: '键盘助手输入 (调用系统键盘助手)' },
  { value: 5, label: '远程输入 (专为ToDesk等远程桌面环境设计)' }
]

// 表单验证规则
const formRules = {
  title: [
    { required: true, message: t('validation.titleRequired'), trigger: 'blur' }
  ],
  username: [
    { required: true, message: t('validation.usernameRequired'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('validation.passwordRequired'), trigger: 'blur' }
  ],
  typeid: [
    { required: true, message: t('validation.typeRequired'), trigger: 'change' }
  ],
  input_method: [
    { required: true, message: t('validation.inputMethodRequired'), trigger: 'change' }
  ]
}

// 20251017 陈凤庆 计算属性：过滤和排序的历史记录
const filteredUsernameHistory = computed(() => {
  if (!usernameHistory.value || usernameHistory.value.length === 0) {
    return []
  }

  const query = usernameSearchQuery.value.toLowerCase()
  let filtered = usernameHistory.value

  // 如果有搜索查询，进行过滤
  if (query) {
    filtered = usernameHistory.value.filter(username =>
      username.toLowerCase().includes(query)
    )
  }

  // 按字母顺序排序
  return filtered.sort((a, b) => a.localeCompare(b))
})

// 监听 modelValue 变化
// 20250101 陈凤庆 修改为异步函数，支持加载标签列表
// 20251002 陈凤庆 增强日志记录，便于调试分析
watch(() => props.modelValue, async (newVal) => {
  console.log(`[PasswordDetail] 监听器触发，modelValue变化: ${newVal}`)

  visible.value = newVal
  if (newVal) {
    console.log(`[PasswordDetail] 对话框打开，开始初始化表单数据`)
    console.log(`[PasswordDetail] 传入的accountData:`, props.accountData)
    console.log(`[PasswordDetail] isNew标志:`, props.isNew)

    // 20251001 陈凤庆 重置表单数据，确保所有ID字段都是字符串类型
    Object.assign(formData, {
      id: '', // 20251001 陈凤庆 改为空字符串，新建时后端会自动生成GUID
      title: '',
      username: '',
      password: '',
      url: '',
      type: '',
      notes: '',
      group_id: '', // 20251001 陈凤庆 改为空字符串，避免硬编码
      typeid: '', // 20251002 陈凤庆 改为typeid，与后端模型保持一致
      icon: '',
      is_favorite: false,
      input_method: 1, // 20251003 陈凤庆 默认输入方式为1（默认方式）
      ...props.accountData // 20251001 陈凤庆 用传入的数据覆盖默认值
    })

    console.log(`[PasswordDetail] 表单数据初始化完成:`, formData)
    console.log(`[PasswordDetail] 分组ID: ${formData.group_id}, 类型ID: ${formData.typeid}`)

    // 加载密码规则列表
    await loadPasswordRulesForGeneration()

    // 20251017 陈凤庆 加载用户名历史记录
    await loadUsernameHistory()

    // 20250101 陈凤庆 如果有分组ID，加载该分组的标签列表
    if (formData.group_id) {
      console.log(`[PasswordDetail] 检测到分组ID，开始加载类型列表: ${formData.group_id}`)
      await loadTabsByGroup(formData.group_id)

      // 20250101 陈凤庆 如果有标签ID，定位到该标签
      if (formData.typeid) {
        console.log(`[PasswordDetail] 定位到类型ID: ${formData.typeid}`)
        // 20251002 陈凤庆 验证类型ID是否在列表中
        const foundType = tabList.value.find(tab => tab.id === formData.typeid)
        if (foundType) {
          console.log(`[PasswordDetail] 类型ID匹配成功，类型名称: ${foundType.name}`)
        } else {
          console.warn(`[PasswordDetail] 类型ID在列表中未找到: ${formData.typeid}`)
        }
      } else {
        console.log(`[PasswordDetail] 没有指定类型ID，将显示空选项`)
      }
    } else {
      console.warn(`[PasswordDetail] 没有分组ID，无法加载类型列表`)
      ElMessage.warning('没有分组信息，无法加载类型列表')
    }

    // 新建时默认进入编辑模式
    isEdit.value = props.isNew
    showPassword.value = false
    // 20251003 陈凤庆 备注显示状态：查看模式下默认显示脱敏备注，编辑模式下显示完整备注
    showNotes.value = isEdit.value
    console.log(`[PasswordDetail] 表单初始化完成，编辑模式: ${isEdit.value}，备注显示状态: ${showNotes.value}`)

    // 加载默认密码规则
    await loadDefaultPasswordRule()

    // 20251002 陈凤庆 清除表单验证状态，确保每次打开都是干净的状态
    if (formRef.value) {
      formRef.value.clearValidate()
    }
  } else {
    console.log(`[PasswordDetail] 对话框关闭`)
  }
})

/**
 * 关闭对话框
 * @modify 20251002 陈凤庆 添加表单验证状态清除
 */
const handleClose = () => {
  // 20251002 陈凤庆 清除表单验证状态，避免下次打开时显示验证错误
  if (formRef.value) {
    formRef.value.clearValidate()
  }

  emit('update:modelValue', false)
  isEdit.value = false
  showPassword.value = false
  showNotes.value = false // 20251003 陈凤庆 重置备注显示状态
}

/**
 * 启用编辑模式
 * @modify 20251002 陈凤庆 添加表单验证状态清除
 * @modify 20251003 陈凤庆 添加获取完整账号数据的功能
 */
const enableEdit = async () => {
  try {
    // 20251002 陈凤庆 清除之前的验证状态，确保编辑模式下重新开始验证
    if (formRef.value) {
      formRef.value.clearValidate()
    }

    // 20251003 陈凤庆 获取完整的解密账号数据用于编辑
    if (formData.id) {
      console.log('[AccountDetail] 获取完整账号数据用于编辑，账号ID:', formData.id)

      const wailsAPI = window.go?.app?.App
      if (wailsAPI) {
        const fullAccountData = await wailsAPI.GetAccountByID(formData.id)
        console.log('[AccountDetail] 获取到完整账号数据:', fullAccountData)

        // 更新表单数据为完整的解密数据
        Object.assign(formData, {
          ...fullAccountData,
          // 保持当前的一些状态
          id: formData.id,
          typeid: formData.typeid || fullAccountData.typeid,
          group_id: formData.group_id || fullAccountData.group_id
        })

        console.log('[AccountDetail] 表单数据已更新为完整数据')
      } else {
        console.warn('[AccountDetail] Wails API 未初始化，无法获取完整账号数据')
      }
    }

    isEdit.value = true
  } catch (error) {
    console.error('[AccountDetail] 获取完整账号数据失败:', error)
    ElMessage.error('获取账号详细信息失败: ' + (error.message || '未知错误'))
    // 即使获取失败，也允许进入编辑模式
    isEdit.value = true
  }
}

/**
 * 切换密码可见性
 * @modify 20251003 陈凤庆 添加显示密码功能，调用后端接口获取解密密码
 */
const togglePasswordVisibility = async () => {
  if (!showPassword.value) {
    // 显示密码：调用后端接口获取解密密码
    if (formData.id && !isEdit.value) {
      try {
        console.log(`[togglePasswordVisibility] 获取密码，账号ID: ${formData.id}`)
        const password = await getAccountPassword(formData.id)
        formData.password = password
        showPassword.value = true
        console.log(`[togglePasswordVisibility] 密码获取成功`)
      } catch (error) {
        console.error('[togglePasswordVisibility] 获取密码失败:', error)
        // 获取失败时不改变显示状态
      }
    } else {
      // 编辑模式或新建模式，直接切换显示状态
      showPassword.value = true
    }
  } else {
    // 隐藏密码：恢复默认显示
    if (!isEdit.value) {
      formData.password = "*****" // 恢复默认显示5个*
    }
    showPassword.value = false
  }
}

/**
 * 复制用户名
 */
const copyUsername = async () => {
  try {
    await navigator.clipboard.writeText(formData.username)
    ElMessage.success('用户名已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

/**
 * 加载用户名历史记录
 * 20251017 陈凤庆 简化加载机制，使用前端变量管理，如果为空则从后端加载
 */
const loadUsernameHistory = async () => {
  if (!isEdit.value) return // 只在编辑模式下加载

  console.log('开始加载用户名历史记录...')

  // 如果前端变量已有数据，直接使用
  if (usernameHistory.value && usernameHistory.value.length > 0) {
    console.log('使用前端变量的用户名历史记录:', usernameHistory.value.length, '个用户名', usernameHistory.value)
    return
  }

  // 前端变量为空，从后端加载
  try {
    usernameHistoryLoading.value = true
    console.log('从后端加载用户名历史记录...')

    // 获取当前登录密码（从后端内存获取）
    const currentPassword = await getCurrentPassword()
    if (!currentPassword) {
      console.log('未找到当前登录密码，跳过加载用户名历史记录')
      usernameHistory.value = []
      return
    }

    // 调用后端API获取历史用户名
    const history = await window.go.app.App.GetUsernameHistory(currentPassword)
    usernameHistory.value = history || []
    console.log('从后端加载用户名历史记录完成:', usernameHistory.value.length, '个用户名', usernameHistory.value)
  } catch (error) {
    console.error('加载用户名历史记录失败:', error)
    usernameHistory.value = []
  } finally {
    usernameHistoryLoading.value = false
  }
}

/**
 * 处理用户名输入事件
 * 20251017 陈凤庆 用户输入时自动搜索历史用户名，显示输入框下方的列表
 */
const handleUsernameInput = async (value) => {
  usernameSearchQuery.value = value
  selectedUsernameIndex.value = -1

  // 如果有输入内容且有历史记录，显示输入框下方的列表
  if (value && usernameHistory.value.length > 0) {
    await loadUsernameHistory()
    if (filteredUsernameHistory.value.length > 0) {
      showUsernameList.value = true
      // 不再自动显示下拉菜单，只有点击下拉按钮才显示
      showUsernameDropdown.value = false
    } else {
      showUsernameList.value = false
    }
  } else {
    showUsernameList.value = false
    showUsernameDropdown.value = false
  }
}

/**
 * 处理用户名键盘事件
 * 20251017 陈凤庆 支持键盘导航和确认
 */
const handleUsernameKeydown = (event) => {
  const hasVisibleSuggestions = showUsernameList.value && filteredUsernameHistory.value.length > 0
  const hasVisibleDropdown = showUsernameDropdown.value && filteredUsernameHistory.value.length > 0

  if (!hasVisibleSuggestions && !hasVisibleDropdown) {
    // 如果没有可见的建议列表或下拉菜单，Tab或Enter确认当前输入
    if (event.key === 'Tab' || event.key === 'Enter') {
      event.preventDefault()
      // 聚焦到下一个输入框（密码框）
      setTimeout(() => {
        const passwordInput = document.querySelector('input[placeholder*="密码"], input[placeholder*="password"]')
        if (passwordInput) {
          passwordInput.focus()
        }
      }, 50)
    }
    return
  }

  // 处理键盘导航
  switch (event.key) {
    case 'ArrowDown':
      event.preventDefault()
      selectedUsernameIndex.value = Math.min(
        selectedUsernameIndex.value + 1,
        Math.min(filteredUsernameHistory.value.length - 1, 7) // 最多显示8个，索引0-7
      )
      break
    case 'ArrowUp':
      event.preventDefault()
      selectedUsernameIndex.value = Math.max(selectedUsernameIndex.value - 1, -1)
      break
    case 'Enter':
    case 'Tab':
      event.preventDefault()
      if (selectedUsernameIndex.value >= 0) {
        // 选择高亮的用户名
        if (hasVisibleSuggestions) {
          selectUsernameFromList(filteredUsernameHistory.value[selectedUsernameIndex.value])
        } else {
          selectUsername(filteredUsernameHistory.value[selectedUsernameIndex.value])
        }
      } else {
        // 确认当前输入的用户名
        showUsernameDropdown.value = false
        showUsernameList.value = false
        // 聚焦到下一个输入框
        setTimeout(() => {
          const passwordInput = document.querySelector('input[placeholder*="密码"], input[placeholder*="password"]')
          if (passwordInput) {
            passwordInput.focus()
          }
        }, 50)
      }
      break
    case 'Escape':
      event.preventDefault()
      showUsernameDropdown.value = false
      showUsernameList.value = false
      selectedUsernameIndex.value = -1
      break
  }
}

/**
 * 处理用户名获得焦点事件
 * 20251017 陈凤庆 获得焦点时显示输入框下方的历史用户名列表
 */
const handleUsernameFocus = async () => {
  await loadUsernameHistory()
  usernameSearchQuery.value = formData.username || ''

  // 如果有历史记录且当前有输入内容，显示输入框下方的列表
  if (usernameHistory.value.length > 0 && formData.username) {
    if (filteredUsernameHistory.value.length > 0) {
      showUsernameList.value = true
    }
  }
}

/**
 * 处理用户名失去焦点事件
 * 20251017 陈凤庆 失去焦点时隐藏列表和下拉菜单
 */
const handleUsernameBlur = () => {
  // 延迟隐藏，以便点击列表项或下拉菜单项
  setTimeout(() => {
    if (!usernameDropdownRef.value?.visible && !isClickingOnSuggestions()) {
      showUsernameDropdown.value = false
      showUsernameList.value = false
      selectedUsernameIndex.value = -1
    }
  }, 200)
}

/**
 * 检查是否正在点击建议列表
 * 20251017 陈凤庆 防止点击建议列表时输入框失去焦点导致列表隐藏
 */
const isClickingOnSuggestions = () => {
  // 这里可以通过检查鼠标位置或其他方式来判断
  // 简单起见，我们通过延迟来处理
  return false
}

/**
 * 处理下拉菜单显示状态变化
 */
const handleDropdownVisibleChange = (visible) => {
  showUsernameDropdown.value = visible
  if (!visible) {
    selectedUsernameIndex.value = -1
  }
}

/**
 * 选择用户名
 * 20251017 陈凤庆 从历史用户名菜单中选择用户名
 */
const selectUsername = (username) => {
  formData.username = username
  usernameSearchQuery.value = username
  showUsernameDropdown.value = false
  selectedUsernameIndex.value = -1
  console.log('选择用户名:', username)

  // 聚焦到密码输入框
  setTimeout(() => {
    const passwordInput = document.querySelector('input[placeholder*="密码"], input[placeholder*="password"]')
    if (passwordInput) {
      passwordInput.focus()
    }
  }, 100)
}

/**
 * 从输入框下方的列表中选择用户名
 * 20251017 陈凤庆 从输入框下方的历史用户名列表中选择用户名
 */
const selectUsernameFromList = (username) => {
  formData.username = username
  usernameSearchQuery.value = username
  showUsernameList.value = false
  showUsernameDropdown.value = false
  selectedUsernameIndex.value = -1
  console.log('从列表选择用户名:', username)

  // 聚焦到密码输入框
  setTimeout(() => {
    const passwordInput = document.querySelector('input[placeholder*="密码"], input[placeholder*="password"]')
    if (passwordInput) {
      passwordInput.focus()
    }
  }, 100)
}



/**
 * 获取当前登录密码
 * 20251017 陈凤庆 从后端内存中获取当前登录密码，提高安全性
 */
const getCurrentPassword = async () => {
  try {
    // 从后端API获取当前登录密码
    const password = await window.go.app.App.GetCurrentPassword()
    return password || ''
  } catch (error) {
    console.error('获取当前登录密码失败:', error)
    return ''
  }
}

/**
 * 保存用户名到历史记录
 * 20251017 陈凤庆 保存账号时将用户名添加到历史记录，并同步更新前端变量
 */
const saveUsernameToHistory = async (username) => {
  if (!username || username.trim() === '') return

  try {
    // 获取当前登录密码
    const currentPassword = await getCurrentPassword()
    if (!currentPassword) {
      console.log('未找到当前登录密码，跳过保存用户名到历史记录')
      return
    }

    console.log('开始保存用户名到历史记录:', username.trim())

    // 调用后端API保存用户名到历史记录
    await window.go.app.App.SaveUsernameToHistory(username.trim(), currentPassword)
    console.log('用户名已保存到后端历史记录:', username.trim())

    // 同步更新前端变量：如果用户名不在历史记录中，则添加到开头
    const trimmedUsername = username.trim()
    if (!usernameHistory.value.includes(trimmedUsername)) {
      usernameHistory.value.unshift(trimmedUsername)
      console.log('用户名已添加到前端历史记录，当前历史记录数量:', usernameHistory.value.length)
    } else {
      console.log('用户名已存在于前端历史记录中')
    }
  } catch (error) {
    console.error('保存用户名到历史记录失败:', error)
    throw error // 重新抛出错误，让调用方知道保存失败
  }
}

/**
 * 刷新用户名历史记录
 * 20251017 陈凤庆 强制从后端重新加载用户名历史记录
 */
const refreshUsernameHistory = async () => {
  try {
    const currentPassword = await getCurrentPassword()
    if (!currentPassword) {
      console.log('未找到当前登录密码，无法刷新用户名历史记录')
      return
    }

    // 从后端重新查询历史记录
    const history = await window.go.app.App.GetUsernameHistory(currentPassword)
    usernameHistory.value = history || []
    console.log('用户名历史记录已刷新:', usernameHistory.value)
  } catch (error) {
    console.error('刷新用户名历史记录失败:', error)
  }
}

/**
 * 删除用户名历史记录
 * @param {string} username 要删除的用户名
 */
const handleDeleteUsername = async (username) => {
  try {
    // 弹出确认对话框
    await ElMessageBox.confirm(
      `确定要删除用户名"${username}"吗？删除后无法恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    // 用户确认删除，调用后端API删除
    const currentPassword = await getCurrentPassword()
    if (!currentPassword) {
      ElMessage.error('未找到当前登录密码，无法删除用户名')
      return
    }

    // 获取当前历史记录
    const currentHistory = await window.go.app.App.GetUsernameHistory(currentPassword)

    // 过滤掉要删除的用户名
    const updatedHistory = currentHistory.filter(name => name !== username)

    // 清空历史记录
    await window.go.app.App.ClearUsernameHistory()

    // 重新保存过滤后的历史记录
    for (const name of updatedHistory) {
      await window.go.app.App.SaveUsernameToHistory(name, currentPassword)
    }

    // 刷新前端缓存
    await refreshUsernameHistory()

    ElMessage.success(`用户名"${username}"删除成功`)
  } catch (error) {
    if (error === 'cancel') {
      // 用户取消删除
      return
    }

    console.error('删除用户名失败:', error)
    const errorMsg = error?.message || error?.toString() || '未知错误'
    ElMessage.error('删除用户名失败: ' + errorMsg)
  }
}

/**
 * 复制密码
 * @modify 20251003 陈凤庆 使用标准复制密码函数，调用后端安全复制接口
 */
const copyPassword = async () => {
  if (formData.id && !isEdit.value) {
    // 详情模式：使用标准复制密码函数
    try {
      await copyAccountPassword(formData.id, formData.title)
    } catch (error) {
      console.error('[copyPassword] 复制密码失败:', error)
    }
  } else {
    // 编辑模式或新建模式：直接复制当前表单中的密码
    try {
      await navigator.clipboard.writeText(formData.password)
      ElMessage.success('密码已复制到剪贴板')

      // 1分钟后清除剪贴板（安全考虑）
      setTimeout(async () => {
        try {
          await navigator.clipboard.writeText('')
        } catch (error) {
          console.warn('清除剪贴板失败:', error)
        }
      }, 60000)
    } catch (error) {
      ElMessage.error('复制失败')
    }
  }
}

/**
 * 切换备注可见性
 * @author 20251003 陈凤庆 添加备注显示/隐藏功能，支持脱敏显示
 */
const toggleNotesVisibility = async () => {
  if (!showNotes.value) {
    // 显示完整备注：调用后端接口获取完整解密备注
    if (formData.id && !isEdit.value) {
      try {
        console.log(`[toggleNotesVisibility] 获取完整备注，账号ID: ${formData.id}`)
        const wailsAPI = window.go?.app?.App
        if (wailsAPI) {
          const notes = await wailsAPI.GetAccountNotes(formData.id)
          formData.notes = notes
          showNotes.value = true
          console.log(`[toggleNotesVisibility] 完整备注获取成功`)
        } else {
          console.warn('[toggleNotesVisibility] Wails API 未初始化')
          ElMessage.warning('无法获取备注信息')
        }
      } catch (error) {
        console.error('[toggleNotesVisibility] 获取备注失败:', error)
        ElMessage.error('获取备注失败')
      }
    } else {
      // 编辑模式或新建模式，直接切换显示状态
      showNotes.value = true
    }
  } else {
    // 隐藏备注：恢复脱敏显示
    if (!isEdit.value && formData.id) {
      // 重新获取脱敏的备注数据
      try {
        console.log(`[toggleNotesVisibility] 恢复脱敏备注显示，账号ID: ${formData.id}`)
        const wailsAPI = window.go?.app?.App
        if (wailsAPI) {
          const accountDetail = await wailsAPI.GetAccountDetail(formData.id)
          formData.notes = accountDetail.notes // 这是脱敏的备注
          console.log(`[toggleNotesVisibility] 脱敏备注恢复成功`)
        }
      } catch (error) {
        console.error('[toggleNotesVisibility] 恢复脱敏备注失败:', error)
        // 如果获取失败，使用简单的脱敏逻辑
        const originalNotes = formData.notes
        if (originalNotes && originalNotes.length > 10) {
          formData.notes = originalNotes.substring(0, 10) + '...'
        }
      }
    }
    showNotes.value = false
  }
}

/**
 * 复制备注
 * @author 20251003 陈凤庆 添加备注复制功能
 */
const copyNotes = async () => {
  if (formData.id && !isEdit.value) {
    // 详情模式：使用后端安全复制方法
    try {
      const wailsAPI = window.go?.app?.App
      if (wailsAPI) {
        await wailsAPI.CopyAccountNotes(formData.id)
        ElMessage.success('备注已复制到剪贴板（10秒后自动清理）')
      } else {
        console.warn('[copyNotes] Wails API 未初始化')
        ElMessage.warning('无法复制备注')
      }
    } catch (error) {
      console.error('[copyNotes] 复制备注失败:', error)
      ElMessage.error('复制备注失败')
    }
  } else {
    // 编辑模式或新建模式：直接复制当前表单中的备注
    try {
      await navigator.clipboard.writeText(formData.notes || '')
      ElMessage.success('备注已复制到剪贴板')
    } catch (error) {
      ElMessage.error('复制失败')
    }
  }
}

/**
 * 生成随机密码（使用默认规则）
 */
const generatePassword = async () => {
  try {
    // 使用默认密码规则生成密码
    if (defaultPasswordRuleId.value) {
      const password = await window.go.app.App.GeneratePasswordByRule(defaultPasswordRuleId.value)
      formData.password = password
      ElMessage.success(t('account.passwordGeneratedWithDefaultRule'))
    } else {
      // 如果没有默认规则，使用简单的密码生成逻辑
      const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
      let password = ''
      for (let i = 0; i < 12; i++) {
        password += chars.charAt(Math.floor(Math.random() * chars.length))
      }
      formData.password = password
      ElMessage.success(t('account.passwordGenerated'))
    }
  } catch (error) {
    console.error('生成密码失败:', error)
    ElMessage.error(t('account.generatePasswordFailed'))
  }
}

/**
 * 处理密码命令
 */
const handlePasswordCommand = (command) => {
  switch (command) {
    case 'generate-default':
      generatePassword()
      break
    case 'select-rule':
      showPasswordRuleSelection()
      break
    default:
      console.warn('未知的密码命令:', command)
  }
}

/**
 * 显示密码规则选择对话框
 */
const showPasswordRuleSelection = async () => {
  try {
    passwordRulesLoading.value = true
    showPasswordRuleDialog.value = true

    // 加载密码规则列表
    const rules = await window.go.app.App.GetAllPasswordRules()
    passwordRules.value = rules || []

    // 设置默认选择的规则
    if (defaultPasswordRuleId.value) {
      selectedPasswordRuleId.value = defaultPasswordRuleId.value
    } else if (passwordRules.value.length > 0) {
      selectedPasswordRuleId.value = passwordRules.value[0].id
    }
  } catch (error) {
    console.error('加载密码规则失败:', error)
    ElMessage.error(t('account.loadPasswordRulesFailed'))
  } finally {
    passwordRulesLoading.value = false
  }
}

/**
 * 使用选择的规则生成密码
 */
const generatePasswordWithRule = async () => {
  if (!selectedPasswordRuleId.value) {
    ElMessage.warning(t('account.selectRuleFirst'))
    return
  }

  try {
    const password = await window.go.app.App.GeneratePasswordByRule(selectedPasswordRuleId.value)
    formData.password = password
    showPasswordRuleDialog.value = false
    ElMessage.success(t('account.passwordGeneratedWithRule'))
  } catch (error) {
    console.error('生成密码失败:', error)
    ElMessage.error(t('account.generatePasswordFailed'))
  }
}

/**
 * 关闭密码规则选择对话框
 */
const handlePasswordRuleDialogClose = () => {
  showPasswordRuleDialog.value = false
  selectedPasswordRuleId.value = ''
}



/**
 * 加载密码规则列表（用于生成密码）
 */
const loadPasswordRulesForGeneration = async () => {
  try {
    const rules = await window.go.app.App.GetAllPasswordRules()
    passwordRules.value = rules || []

    // 查找默认规则
    const defaultRule = passwordRules.value.find(rule => rule.is_default)
    if (defaultRule) {
      defaultPasswordRuleId.value = defaultRule.id
    }

    console.log('加载密码规则成功:', passwordRules.value)
  } catch (error) {
    console.error('加载密码规则失败:', error)
  }
}

/**
 * 使用默认规则生成密码
 * 20251017 陈凤庆 修改逻辑：当没有密码规则时，使用内置默认规则（大小写+数字+特殊字符，8位）
 */
const generatePasswordWithDefaultRule = async () => {
  try {
    let password = ''
    let ruleName = '默认规则'

    // 先检查是否有密码规则
    const rules = await window.go.app.App.GetAllPasswordRules()

    if (rules && rules.length > 0) {
      // 如果有密码规则，尝试使用默认规则
      const defaultRule = rules.find(rule => rule.is_default)
      if (defaultRule && defaultRule.id) {
        password = await window.go.app.App.GeneratePasswordByRule(defaultRule.id)
        ruleName = defaultRule.name || '默认规则'
      } else {
        // 如果没有默认规则，使用第一个规则
        const firstRule = rules[0]
        password = await window.go.app.App.GeneratePasswordByRule(firstRule.id)
        ruleName = firstRule.name || '规则'
      }
    } else {
      // 如果没有任何密码规则，使用内置默认规则生成密码
      console.log('没有找到密码规则，使用内置默认规则生成密码')
      password = generatePasswordWithBuiltinRule()
      ruleName = '内置默认规则'
    }

    if (!password) {
      ElMessage.error('生成密码失败：无法生成密码')
      return
    }

    formData.password = password
    ElMessage.success(`使用"${ruleName}"生成密码成功`)
  } catch (error) {
    console.error('生成密码失败:', error)
    // 如果API调用失败，使用内置默认规则
    console.log('API调用失败，使用内置默认规则生成密码')
    const password = generatePasswordWithBuiltinRule()
    formData.password = password
    ElMessage.success('使用内置默认规则生成密码成功')
  }
}

/**
 * 使用内置默认规则生成密码（大小写+数字+特殊字符，8位）
 * 20251017 陈凤庆 新增函数
 */
const generatePasswordWithBuiltinRule = () => {
  const uppercase = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const lowercase = 'abcdefghijklmnopqrstuvwxyz'
  const numbers = '0123456789'
  const specialChars = '!@#$%^&*'

  let password = ''

  // 确保至少包含一个大写字母、小写字母、数字和特殊字符
  password += uppercase[Math.floor(Math.random() * uppercase.length)]
  password += lowercase[Math.floor(Math.random() * lowercase.length)]
  password += numbers[Math.floor(Math.random() * numbers.length)]
  password += specialChars[Math.floor(Math.random() * specialChars.length)]

  // 剩余4位从所有字符中随机选择
  const allChars = uppercase + lowercase + numbers + specialChars
  for (let i = 4; i < 8; i++) {
    password += allChars[Math.floor(Math.random() * allChars.length)]
  }

  // 打乱密码字符顺序
  return password.split('').sort(() => Math.random() - 0.5).join('')
}

/**
 * 处理密码规则命令
 */
const handlePasswordRuleCommand = async (ruleId) => {
  try {
    if (!ruleId) {
      ElMessage.error('密码规则ID不能为空')
      return
    }

    // 重新加载密码规则列表，确保数据是最新的
    console.log('重新加载密码规则列表，确保数据最新')
    const rules = await window.go.app.App.GetAllPasswordRules()
    passwordRules.value = rules || []
    console.log('当前密码规则列表:', passwordRules.value)

    // 验证规则ID是否存在
    const rule = passwordRules.value.find(r => r.id === ruleId)
    if (!rule) {
      console.error('规则ID不存在:', ruleId)
      ElMessage.error('选择的密码规则不存在，请重新选择')
      return
    }

    console.log('使用规则生成密码:', rule.name, ruleId)
    const password = await window.go.app.App.GeneratePasswordByRule(ruleId)
    if (!password) {
      ElMessage.error('生成密码失败：返回值为空')
      return
    }

    formData.password = password
    ElMessage.success(`使用"${rule.name}"生成密码成功`)
  } catch (error) {
    console.error('生成密码失败:', error)
    const errorMsg = error?.message || error?.toString() || '未知错误'
    ElMessage.error('生成密码失败: ' + errorMsg)
  }
}

/**
 * 加载默认密码规则
 */
const loadDefaultPasswordRule = async () => {
  try {
    const rules = await window.go.app.App.GetAllPasswordRules()
    const defaultRule = rules.find(rule => rule.is_default)
    if (defaultRule) {
      defaultPasswordRuleId.value = defaultRule.id
      console.log(`[PasswordDetail] 加载默认密码规则成功: ${defaultRule.name}`)
    } else {
      console.log(`[PasswordDetail] 未找到默认密码规则`)
    }
  } catch (error) {
    console.error('加载默认密码规则失败:', error)
    // 不显示错误消息，因为这不是关键功能
  }
}

/**
 * 打开地址
 */
const openUrl = () => {
  if (formData.url) {
    ElMessage.info(`打开地址: ${formData.url}`)
    // 在实际应用中会调用系统默认浏览器
  } else {
    ElMessage.warning('地址为空')
  }
}

/**
 * 保存账号
 * @modify 20251002 陈凤庆 修改数据传递，确保type字段正确传递给后端
 */
const handleSave = async () => {
  if (!formRef.value) return

  try {
    // 20251017 陈凤庆 添加保存前的数据验证和调试信息
    console.log('[PasswordDetail] 开始保存，当前表单数据:', formData)
    console.log('[PasswordDetail] 表单验证规则:', formRules)

    // 验证表单
    await formRef.value.validate()
    console.log('[PasswordDetail] 表单验证通过')

    saving.value = true

    // 20251017 陈凤庆 检查必填字段
    if (!formData.title || !formData.username || !formData.password || !formData.typeid) {
      const missingFields = []
      if (!formData.title) missingFields.push('标题')
      if (!formData.username) missingFields.push('用户名')
      if (!formData.password) missingFields.push('密码')
      if (!formData.typeid) missingFields.push('类型')
      throw new Error(`缺少必填字段: ${missingFields.join(', ')}`)
    }

    // 20251002 陈凤庆 准备保存数据，确保typeid字段正确
    const saveData = {
      ...formData,
      type: formData.typeid || formData.type || '', // 保持type字段兼容性
      typeid: formData.typeid || formData.type || '' // 使用typeid字段
    }

    console.log('[PasswordDetail] 准备保存数据:', saveData)

    // 发送保存事件
    emit('save', saveData)

    // 20251017 陈凤庆 保存用户名到历史记录（不阻塞主流程）
    if (formData.username && formData.username.trim()) {
      // 异步保存，不等待结果，避免阻塞主流程
      saveUsernameToHistory(formData.username.trim())
        .then(() => {
          console.log('用户名历史记录保存完成')
        })
        .catch((error) => {
          console.error('保存用户名到历史记录失败:', error)
          // 即使保存历史记录失败，也不影响账号保存成功的流程
        })
    }

    //ElMessage.success(props.isNew ? '账号创建成功' : '账号更新成功')
    handleClose()

  } catch (error) {
    console.error('保存失败:', error)
    // 20251017 陈凤庆 改进错误处理，显示更详细的错误信息
    let errorMsg = '未知错误'

    if (typeof error === 'string') {
      errorMsg = error
    } else if (error?.message) {
      errorMsg = error.message
    } else if (error?.toString && typeof error.toString === 'function') {
      const errorStr = error.toString()
      if (errorStr !== '[object Object]') {
        errorMsg = errorStr
      }
    }

    // 如果错误是由于表单验证失败，显示更友好的消息
    if (error?.name === 'ValidationError' || errorMsg.includes('validation')) {
      errorMsg = '请检查必填字段是否已正确填写'
    }

    ElMessage.error(`保存失败: ${errorMsg}`)
  } finally {
    saving.value = false
  }
}

/**
 * 删除账号
 */
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm(
      t('account.confirmDeleteMessage', { title: formData.title }),
      t('account.confirmDeleteTitle'),
      {
        confirmButtonText: t('common.delete'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    emit('delete', formData.id)
    handleClose()

  } catch (error) {
    // 用户取消删除
  }
}

/**
 * 根据分组ID加载标签列表
 * @param {string} groupId 分组ID
 * @author 陈凤庆
 * @date 20250101
 * @modify 20251002 陈凤庆 修复API调用方法名，GetTabsByGroup改为GetTypesByGroup，增强日志记录
 */
const loadTabsByGroup = async (groupId) => {
  try {
    console.log(`[PasswordDetail] 开始加载分组${groupId}的类型列表...`)
    console.log(`[PasswordDetail] 参数验证 - groupId类型: ${typeof groupId}, 值: ${groupId}`)

    // 20251002 陈凤庆 参数验证
    if (!groupId) {
      console.error('[PasswordDetail] 分组ID为空，无法加载类型列表')
      ElMessage.error('分组ID为空，无法加载类型列表')
      return
    }

    tabListLoading.value = true
    console.log(`[PasswordDetail] 设置加载状态为true`)

    // 20251002 陈凤庆 调用正确的API方法GetTypesByGroup
    console.log(`[PasswordDetail] 调用后端API: GetTypesByGroup(${groupId})`)
    const tabs = await window.go.app.App.GetTypesByGroup(groupId)

    console.log(`[PasswordDetail] API调用成功，返回数据:`, tabs)
    console.log(`[PasswordDetail] 返回数据类型: ${typeof tabs}, 是否为数组: ${Array.isArray(tabs)}`)
    console.log(`[PasswordDetail] 获取到${tabs?.length || 0}个类型`)

    // 20251002 陈凤庆 验证返回数据格式
    if (!Array.isArray(tabs)) {
      console.error('[PasswordDetail] 后端返回数据不是数组格式:', tabs)
      ElMessage.error('后端返回数据格式错误')
      tabList.value = []
      return
    }

    // 更新标签列表
    tabList.value = tabs || []
    console.log(`[PasswordDetail] 类型列表更新完成，当前列表长度: ${tabList.value.length}`)

  } catch (error) {
    console.error('[PasswordDetail] 加载类型列表失败:', error)
    console.error('[PasswordDetail] 错误类型:', error.name)
    console.error('[PasswordDetail] 错误消息:', error.message)
    console.error('[PasswordDetail] 错误堆栈:', error.stack)

    // 20251002 陈凤庆 根据错误类型提供更具体的错误信息
    let errorMessage = '加载类型列表失败'
    if (error.message) {
      errorMessage += `: ${error.message}`
    }

    ElMessage.error(errorMessage)
    tabList.value = []
  } finally {
    tabListLoading.value = false
    console.log(`[PasswordDetail] 设置加载状态为false，加载过程结束`)
  }
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.left-footer-buttons {
  display: flex;
}

.right-footer-buttons {
  display: flex;
  gap: 10px;
}

:deep(.el-input-group__append) {
  padding: 0;
}

:deep(.el-input-group__append .el-button) {
  margin: 0;
  border-left: 1px solid #dcdfe6;
}

:deep(.el-button-group .el-button) {
  margin: 0;
}

/* 20251017 陈凤庆 确保用户名输入框有足够宽度，按钮左右排列 */
:deep(.el-form-item .el-input-group) {
  width: 100%;
}

:deep(.el-input-group__append .el-button-group) {
  display: flex;
  flex-direction: row;
}

:deep(.el-input-group__append .el-button-group .el-button) {
  flex-shrink: 0;
  white-space: nowrap;
}

/* 20251017 陈凤庆 用户名表单项相对定位，支持浮动建议列表 */
.username-form-item {
  position: relative;
}

/* 20251003 陈凤庆 备注字段容器样式 - 上下布局，按钮与编辑框边框对齐 */
.notes-field-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.notes-input-wrapper {
  width: 100%;
}

.notes-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  /* 与编辑框右边框对齐 */
  margin-right: 0;
  padding-right: 0;
}

.notes-actions .el-button-group {
  background: white;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* 确保备注编辑框与其他输入框对齐 */
.notes-field-container .el-textarea {
  width: 100%;
}

.notes-field-container .el-textarea__inner {
  width: 100%;
  box-sizing: border-box;
}

.default-rule-tag {
  color: #409eff;
  font-size: 12px;
  margin-left: 8px;
}

/* 20251017 陈凤庆 用户名历史记录下拉菜单样式 */
.username-history-dropdown :deep(.el-button) {
  border-left: 1px solid #dcdfe4;
  border-right: 1px solid #dcdfe4;
  border-top: 1px solid #dcdfe4;
  border-bottom: 1px solid #dcdfe4;
}

.username-history-dropdown :deep(.el-button:hover) {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #f0f9ff;
}

/* 用户名下拉菜单项选中状态 */
.username-history-dropdown :deep(.el-dropdown-menu__item.is-selected) {
  background-color: #f0f9ff;
  color: #409eff;
}

/* 用户名下拉菜单项中的删除按钮 */
.username-history-dropdown :deep(.delete-username-btn) {
  opacity: 0;
  transition: opacity 0.3s;
}

.username-history-dropdown :deep(.el-dropdown-menu__item:hover .delete-username-btn) {
  opacity: 1;
}

/* 20251017 陈凤庆 用户名输入框下方建议列表样式 - 改为浮动不撑开布局 */
.username-suggestions-list {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  border: 1px solid #dcdfe4;
  border-radius: 4px;
  background: #fff;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 1000;
  max-height: 240px;
  overflow-y: auto;
}

/* 20251017 陈凤庆 移除suggestions-header和suggestions-title样式，不再需要 */

.suggestions-content {
  padding: 4px 0;
}

.suggestion-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 14px;
}

.suggestion-item:hover,
.suggestion-item.is-selected {
  background-color: #f0f9ff;
  color: #409eff;
}

.suggestion-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.delete-suggestion-btn {
  opacity: 0;
  transition: opacity 0.3s;
  margin-left: 8px;
}

.suggestion-item:hover .delete-suggestion-btn {
  opacity: 1;
}

/* 修复密码规则下拉按钮的样式 */
.password-rule-dropdown :deep(.el-button) {
  border-left: 1px solid #dcdfe4;
  border-right: 1px solid #dcdfe4;
  border-top: 1px solid #dcdfe4;
  border-bottom: 1px solid #dcdfe4;
}

.password-rule-dropdown :deep(.el-button:hover) {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #f0f9ff;
}


</style>
