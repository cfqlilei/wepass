<template>
  <el-dialog
    v-model="visible"
    :title="$t('passwordRuleSettings.title')"
    :width="dialogWidth"
    :top="dialogTop"
    :before-close="handleClose"
    :modal="true"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    :destroy-on-close="false"
    :append-to-body="true"
    :lock-scroll="true"
    class="password-rule-settings-dialog"
    draggable
  >
    <div class="rule-settings">
      <!-- 密码规则选择 -->
      <el-card class="rule-selector-card">
        <template #header>
          <div class="card-header">
            <span>选择密码规则</span>
          </div>
        </template>

        <el-form-item label="密码规则">
          <div style="display: flex; gap: 10px; align-items: center; width: 100%;">
            <el-select
              v-model="selectedRuleId"
              placeholder="请选择密码规则"
              style="flex: 1;"
              @change="handleRuleChange"
            >
              <!-- 20251017 陈凤庆 删除硬编码的通用规则选项，只从数据库加载规则 -->
              <el-option
                v-for="rule in savedRules"
                :key="rule.id"
                :label="rule.name"
                :value="rule.id"
              />
            </el-select>
            <!-- 删除按钮 - 20251017 陈凤庆 新增删除功能 -->
            <el-button
              v-if="selectedRuleId"
              type="danger"
              :icon="Delete"
              @click="handleDeleteRule"
              :disabled="deleting"
              title="删除密码规则"
            >
              删除
            </el-button>
            <!-- 默认规则开关 - 在下拉框右侧 -->
            <!-- 20251017 陈凤庆 修复重复通用规则问题，只要选择了规则就显示默认开关 -->
            <div v-if="selectedRuleId" style="display: flex; align-items: center; gap: 8px; white-space: nowrap;">
              <span style="font-size: 14px;">默认规则</span>
              <el-switch
                v-model="isCurrentRuleDefault"
                @change="handleDefaultRuleToggle"
                :disabled="settingDefault"
              />
            </div>
          </div>
        </el-form-item>
      </el-card>

      <!-- 密码规则配置 -->
      <el-card class="rule-config-card">
        <template #header>
          <div class="card-header">
            <span>密码规则配置</span>
          </div>
        </template>

        <!-- 配置模式选择 -->
        <el-form-item label="配置模式">
          <el-radio-group v-model="configMode" @change="handleConfigModeChange">
            <el-radio label="general">默认方式</el-radio>
            <el-radio label="custom">自定义规则</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 默认方式配置 -->
        <div v-if="configMode === 'general'" class="general-config">
          <el-form :model="generalRule" label-width="140px">
            <el-form-item>
              <template #label>
                <div class="form-label-with-switch">
                  <span>包含大写字母</span>
                  <el-switch v-model="generalRule.includeUppercase" />
                </div>
              </template>
              <div v-if="generalRule.includeUppercase" class="min-count-input">
                <span>最小位数：</span>
                <el-input-number
                  v-model="generalRule.minUppercase"
                  :min="1"
                  :max="20"
                  size="small"
                  style="width: 80px"
                />
                <span>位</span>
              </div>
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="form-label-with-switch">
                  <span>包含小写字母</span>
                  <el-switch v-model="generalRule.includeLowercase" />
                </div>
              </template>
              <div v-if="generalRule.includeLowercase" class="min-count-input">
                <span>最小位数：</span>
                <el-input-number
                  v-model="generalRule.minLowercase"
                  :min="1"
                  :max="20"
                  size="small"
                  style="width: 80px"
                />
                <span>位</span>
              </div>
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="form-label-with-switch">
                  <span>包含数字</span>
                  <el-switch v-model="generalRule.includeNumbers" />
                </div>
              </template>
              <div v-if="generalRule.includeNumbers" class="min-count-input">
                <span>最小位数：</span>
                <el-input-number
                  v-model="generalRule.minNumbers"
                  :min="1"
                  :max="20"
                  size="small"
                  style="width: 80px"
                />
                <span>位</span>
              </div>
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="form-label-with-switch">
                  <span>包含特殊字符</span>
                  <el-switch v-model="generalRule.includeSpecialChars" />
                </div>
              </template>
              <div v-if="generalRule.includeSpecialChars" class="min-count-input">
                <span>最小位数：</span>
                <el-input-number
                  v-model="generalRule.minSpecialChars"
                  :min="1"
                  :max="20"
                  size="small"
                  style="width: 80px"
                />
                <span>位</span>
              </div>
            </el-form-item>

            <el-form-item>
              <template #label>
                <div class="form-label-with-switch">
                  <span>自定义字符</span>
                  <el-switch v-model="generalRule.includeCustomChars" />
                </div>
              </template>
              <div v-if="generalRule.includeCustomChars">
                <div class="min-count-input">
                  <span>最小位数：</span>
                  <el-input-number
                    v-model="generalRule.minCustomChars"
                    :min="1"
                    :max="20"
                    size="small"
                    style="width: 80px"
                  />
                  <span>位</span>
                </div>
                <div class="custom-chars-input">
                  <el-input
                    v-model="generalRule.customSpecialChars"
                    placeholder="请输入自定义字符，如：@#$abc123"
                    style="margin-top: 8px"
                  />
                </div>
              </div>
            </el-form-item>

            <el-form-item label="密码长度">
              <el-input-number
                v-model="generalRule.length"
                :min="getMinPasswordLength()"
                :max="128"
              />
              <span class="length-hint">（最小长度：{{ getMinPasswordLength() }}）</span>
            </el-form-item>
          </el-form>

          <!-- 生成密码测试 -->
          <div class="password-test">
            <el-button type="primary" @click="generateTestPassword">生成密码</el-button>
            <el-input
              v-model="generatedPassword"
              placeholder="生成的密码将显示在这里"
              readonly
              style="margin-left: 10px; flex: 1"
            >
              <template #append>
                <el-button @click="copyPassword" :disabled="!generatedPassword">
                  复制
                </el-button>
              </template>
            </el-input>
          </div>
        </div>

        <!-- 自定义规则配置 -->
        <div v-if="configMode === 'custom'" class="custom-config">
          <el-form :model="customRule" label-width="140px">
            <el-form-item label="自定义规则模式">
              <el-input
                v-model="customRule.pattern"
                type="textarea"
                :rows="3"
                placeholder="请输入自定义规则模式，例如：a{5}A{3}U{2}d{4}"
              />
            </el-form-item>
          </el-form>

          <!-- 生成密码测试 -->
          <div class="password-test">
            <el-button type="primary" @click="generateTestPassword">生成密码</el-button>
            <el-input
              v-model="generatedPassword"
              placeholder="生成的密码将显示在这里"
              readonly
              style="margin-left: 10px; flex: 1"
            >
              <template #append>
                <el-button @click="copyPassword" :disabled="!generatedPassword">
                  复制
                </el-button>
              </template>
            </el-input>
          </div>
        </div>
      </el-card>

      <!-- 自定义密码规则说明 -->
      <el-card class="rule-card">
        <template #header>
          <span>自定义密码规则说明</span>
        </template>
        
        <div class="custom-rules-help">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="a">Lower-Case Alphanumeric: abcdefghijklmnopqrstuvwxyz 0123456789</el-descriptions-item>
            <el-descriptions-item label="A">Mixed-Case Alphanumeric: ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz 0123456789</el-descriptions-item>
            <el-descriptions-item label="U">Upper-Case Alphanumeric: ABCDEFGHIJKLMNOPQRSTUVWXYZ 0123456789</el-descriptions-item>
            <el-descriptions-item label="d">Digit: 0123456789</el-descriptions-item>
            <el-descriptions-item label="h">Lower-Case Hex Character: 0123456789 abcdef</el-descriptions-item>
            <el-descriptions-item label="H">Upper-Case Hex Character: 0123456789 ABCDEF</el-descriptions-item>
            <el-descriptions-item label="l">Lower-Case Letter: abcdefghijklmnopqrstuvwxyz</el-descriptions-item>
            <el-descriptions-item label="L">Mixed-Case Letter: ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz</el-descriptions-item>
            <el-descriptions-item label="u">Upper-Case Letter: ABCDEFGHIJKLMNOPQRSTUVWXYZ</el-descriptions-item>
            <el-descriptions-item label="v">Lower-Case Vowel: aeiou</el-descriptions-item>
            <el-descriptions-item label="V">Mixed-Case Vowel: AEIOU aeiou</el-descriptions-item>
            <el-descriptions-item label="Z">Upper-Case Vowel: AEIOU</el-descriptions-item>
            <el-descriptions-item label="c">Lower-Case Consonant: bcdfghjklmnpqrstvwxyz</el-descriptions-item>
            <el-descriptions-item label="C">Mixed-Case Consonant: BCDFGHJKLMNPQRSTVWXYZ bcdfghjklmnpqrstvwxyz</el-descriptions-item>
            <el-descriptions-item label="z">Upper-Case Consonant: BCDFGHJKLMNPQRSTVWXYZ</el-descriptions-item>
            <el-descriptions-item label="p">Punctuation: ,.;:</el-descriptions-item>
            <el-descriptions-item label="b">Bracket: ()[]{}&lt;&gt;</el-descriptions-item>
            <el-descriptions-item label="s">Printable 7-Bit Special Character: !&quot;#$%&amp;'()*+,-./:;&lt;=&gt;?@[\]^_`{|}~</el-descriptions-item>
            <el-descriptions-item label="S">Printable 7-Bit ASCII: A-Z, a-z, 0-9, !&quot;#$%&amp;'()*+,-./:;&lt;=&gt;?@[\]^_`{|}~</el-descriptions-item>
            <el-descriptions-item label="x">Latin-1 Supplement: Range [U+00A1, U+00FF] except U+00AD</el-descriptions-item>
            <el-descriptions-item label="\">\">Escape (Fixed Char): Use following character as is.</el-descriptions-item>
            <el-descriptions-item label="{n}">Escape (Repeat): Repeat the previous placeholder n times.</el-descriptions-item>
            <el-descriptions-item label="[...]">Custom Char Set: Define a custom character set.</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-card>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="success" @click="showSaveRuleDialog">保存为密码规则</el-button>
      </div>
    </template>

    <!-- 保存密码规则对话框 -->
    <el-dialog
      v-model="showSaveDialog"
      title="保存密码规则"
      width="400px"
    >
      <el-form :model="saveForm" label-width="100px" ref="saveFormRef">
        <el-form-item
          label="规则名称"
          required
          :rules="[{ required: true, message: '请输入规则名称', trigger: 'blur' }]"
          prop="name"
        >
          <el-select
            v-model="saveForm.name"
            placeholder="请选择或输入规则名称"
            filterable
            allow-create
            default-first-option
            @change="checkRuleNameExists"
            style="width: 100%"
          >
            <el-option
              v-for="rule in savedRules"
              :key="rule.id"
              :label="rule.name"
              :value="rule.name"
            />
          </el-select>
          <div v-if="nameExistsWarning" class="name-warning">
            <el-text type="warning" size="small">
              规则名称已存在，保存将覆盖现有规则
            </el-text>
          </div>
        </el-form-item>
        <el-form-item label="规则描述">
          <el-input
            v-model="saveForm.description"
            placeholder="请输入规则描述"
            type="textarea"
            :rows="2"
          />
        </el-form-item>
        <el-form-item label="设为默认规则">
          <el-switch
            v-model="saveForm.isDefault"
            :disabled="saving"
          />
          <span class="default-rule-hint">设为默认规则后，创建账号时将优先使用此规则生成密码</span>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showSaveDialog = false">取消</el-button>
        <el-button type="primary" @click="savePasswordRule" :loading="saving">
          {{ nameExistsWarning ? '覆盖保存' : '保存' }}
        </el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'

/**
 * 密码规则设置组件
 * @author 陈凤庆
 * @description 设置和管理密码生成规则
 * @modify 20251003 陈凤庆 添加响应式布局支持，自适应屏幕大小
 */

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
const showSaveDialog = ref(false)
const saving = ref(false)
const nameExistsWarning = ref(false)
const deleting = ref(false) // 20251017 陈凤庆 新增删除状态

// 窗口宽度响应式变量
const windowWidth = ref(window.innerWidth)

// 响应式对话框宽度
const dialogWidth = computed(() => {
  // 最大宽度750px，最小宽度500px，确保不超出主界面宽度
  const maxWidth = 750
  const minWidth = 500

  // 左右边距各80px，总共160px，增加边距确保不被遮挡
  const margin = 160

  // 计算可用宽度
  const availableWidth = windowWidth.value - margin

  // 返回合适的宽度，确保不超过最大宽度且不小于最小宽度，保留边距
  const finalWidth = Math.min(maxWidth, Math.max(minWidth, availableWidth))

  return `${finalWidth}px`
})

// 响应式对话框高度
const dialogHeight = computed(() => {
  // 根据窗口高度动态计算对话框高度
  const windowHeight = window.innerHeight
  const bottomMargin = 120 // 底部边距100px + 20px缓冲
  const maxHeight = Math.min(800, windowHeight - bottomMargin) // 最大高度800px或窗口高度减去底部边距
  const minHeight = 600 // 20251017 陈凤庆 增加最小高度从500px到600px

  return `${Math.max(minHeight, maxHeight)}px`
})

// 响应式对话框顶部位置
// 20251017 陈凤庆 调整对话框位置，底部边距减少为100px
const dialogTop = computed(() => {
  const windowHeight = window.innerHeight
  const bottomMargin = 100 // 底部边距100px
  const dialogHeightValue = parseInt(dialogHeight.value) // 获取对话框高度数值

  // 计算顶部位置，确保底部边距为100px
  const topPosition = Math.max(20, windowHeight - dialogHeightValue - bottomMargin)

  return `${topPosition}px`
})

// 界面状态
const selectedRuleId = ref('')
const configMode = ref('general')
const generatedPassword = ref('')
const isCurrentRuleDefault = ref(false)
const settingDefault = ref(false)

// 通用规则设置
const generalRule = reactive({
  includeUppercase: true,
  includeLowercase: true,
  includeNumbers: true,
  includeSpecialChars: false,
  includeCustomChars: false,
  minUppercase: 1,
  minLowercase: 1,
  minNumbers: 1,
  minSpecialChars: 1,
  minCustomChars: 1,
  length: 12,
  customSpecialChars: ''
})

// 自定义规则设置
const customRule = reactive({
  pattern: '',
  description: ''
})

// 已保存的规则
const savedRules = ref([])

// 保存表单
const saveForm = reactive({
  name: '',
  description: '',
  isDefault: false
})

// 表单引用
const saveFormRef = ref(null)

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
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
  if (newVal) {
    // 对话框打开时加载密码规则
    loadPasswordRules()
  }
})

/**
 * 加载所有密码规则
 */
const loadPasswordRules = async () => {
  try {
    const rules = await window.go.app.App.GetAllPasswordRules()
    savedRules.value = rules || []
  } catch (error) {
    console.error('加载密码规则失败:', error)
    ElMessage.error('加载密码规则失败')
  }
}

/**
 * 处理规则选择变化
 * @author 陈凤庆
 * @date 20251017
 * @description 修复重复通用规则问题，只从数据库加载规则
 */
const handleRuleChange = async (ruleId) => {
  if (!ruleId) {
    return
  }

  // 加载选中的规则配置
  try {
    const rule = await window.go.app.App.GetPasswordRuleByID(ruleId)
    loadRuleConfig(rule)
    isCurrentRuleDefault.value = rule.is_default || false
  } catch (error) {
    console.error('加载规则配置失败:', error)
    ElMessage.error('加载规则配置失败')
    isCurrentRuleDefault.value = false
  }
}

/**
 * 处理配置模式变化
 */
const handleConfigModeChange = (mode) => {
  if (mode === 'general') {
    resetGeneralRule()
  } else {
    resetCustomRule()
  }
}

/**
 * 处理默认规则开关切换
 */
const handleDefaultRuleToggle = async (isDefault) => {
  if (!selectedRuleId.value) {
    return
  }

  settingDefault.value = true

  try {
    // 使用正确的后端API函数名
    await window.go.app.App.SetPasswordRuleAsDefault(selectedRuleId.value, isDefault)

    if (isDefault) {
      ElMessage.success('已设置为默认规则')
    } else {
      ElMessage.success('已取消默认规则')
    }

    // 重新加载规则列表以更新状态
    await loadPasswordRules()
  } catch (error) {
    console.error('设置默认规则失败:', error)
    ElMessage.error('设置默认规则失败: ' + error.message)
    // 恢复开关状态
    isCurrentRuleDefault.value = !isDefault
  } finally {
    settingDefault.value = false
  }
}

/**
 * 处理默认规则变化
 */
const handleDefaultRuleChange = async (isDefault) => {
  if (!selectedRuleId.value) {
    return
  }

  try {
    saving.value = true
    await window.go.app.App.SetPasswordRuleAsDefault(selectedRuleId.value, isDefault)

    if (isDefault) {
      ElMessage.success('已设置为默认密码规则')
      // 重新加载规则列表以更新其他规则的默认状态
      await loadPasswordRules()
    } else {
      ElMessage.success('已取消默认密码规则设置')
    }
  } catch (error) {
    console.error('设置默认规则失败:', error)
    ElMessage.error('设置默认规则失败')
    // 恢复开关状态
    isCurrentRuleDefault.value = !isDefault
  } finally {
    saving.value = false
  }
}

/**
 * 重置通用规则为默认值
 */
const resetGeneralRule = () => {
  Object.assign(generalRule, {
    includeUppercase: true,
    includeLowercase: true,
    includeNumbers: true,
    includeSpecialChars: false,
    includeCustomChars: false,
    minUppercase: 1,
    minLowercase: 1,
    minNumbers: 1,
    minSpecialChars: 1,
    minCustomChars: 1,
    length: 12,
    customSpecialChars: ''
  })
}

/**
 * 重置自定义规则为默认值
 */
const resetCustomRule = () => {
  Object.assign(customRule, {
    pattern: '',
    description: ''
  })
}

/**
 * 加载规则配置到界面
 */
const loadRuleConfig = (rule) => {
  try {
    const config = JSON.parse(rule.config)

    if (rule.rule_type === 'general') {
      configMode.value = 'general'
      Object.assign(generalRule, config)
    } else {
      configMode.value = 'custom'
      Object.assign(customRule, config)
    }
  } catch (error) {
    console.error('解析规则配置失败:', error)
    ElMessage.error('规则配置格式错误')
  }
}

/**
 * 计算最小密码长度
 */
const getMinPasswordLength = () => {
  let minLength = 0
  if (generalRule.includeUppercase) minLength += generalRule.minUppercase
  if (generalRule.includeLowercase) minLength += generalRule.minLowercase
  if (generalRule.includeNumbers) minLength += generalRule.minNumbers
  if (generalRule.includeSpecialChars) minLength += generalRule.minSpecialChars
  if (generalRule.includeCustomChars) minLength += generalRule.minCustomChars
  return Math.max(minLength, 4)
}

/**
 * 转换前端数据结构为后端期望的格式
 */
const convertGeneralRuleToBackend = (frontendRule) => {
  return {
    include_uppercase: frontendRule.includeUppercase,
    include_lowercase: frontendRule.includeLowercase,
    include_numbers: frontendRule.includeNumbers,
    include_special_chars: frontendRule.includeSpecialChars,
    include_custom_chars: frontendRule.includeCustomChars,
    min_uppercase: frontendRule.minUppercase,
    min_lowercase: frontendRule.minLowercase,
    min_numbers: frontendRule.minNumbers,
    min_special_chars: frontendRule.minSpecialChars,
    min_custom_chars: frontendRule.minCustomChars,
    length: frontendRule.length,
    custom_special_chars: frontendRule.customSpecialChars
  }
}

/**
 * 转换自定义规则数据结构为后端期望的格式
 */
const convertCustomRuleToBackend = (frontendRule) => {
  return {
    pattern: frontendRule.pattern,
    length: frontendRule.length
  }
}

/**
 * 生成测试密码
 */
const generateTestPassword = async () => {
  try {
    let password = ''

    if (configMode.value === 'general') {
      // 转换数据结构并使用通用规则生成密码
      const backendRule = convertGeneralRuleToBackend(generalRule)
      console.log('发送到后端的通用规则数据:', backendRule)
      password = await window.go.app.App.GeneratePasswordByGeneralConfig(backendRule)
    } else {
      // 转换数据结构并使用自定义规则生成密码
      const backendRule = convertCustomRuleToBackend(customRule)
      console.log('发送到后端的自定义规则数据:', backendRule)
      password = await window.go.app.App.GeneratePasswordByCustomConfig(backendRule)
    }

    generatedPassword.value = password
    ElMessage.success('密码生成成功')
  } catch (error) {
    console.error('生成密码失败:', error)
    ElMessage.error('生成密码失败: ' + (error.message || error))
  }
}

/**
 * 复制密码到剪贴板
 */
const copyPassword = async () => {
  if (!generatedPassword.value) return

  try {
    await navigator.clipboard.writeText(generatedPassword.value)
    ElMessage.success('密码已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error('复制失败')
  }
}

/**
 * 显示保存规则对话框
 */
const showSaveRuleDialog = () => {
  // 默认显示当前选择的密码规则名称
  const currentRule = savedRules.value.find(rule => rule.id === selectedRuleId.value)
  saveForm.name = currentRule ? currentRule.name : ''
  saveForm.description = ''
  saveForm.isDefault = false
  nameExistsWarning.value = false
  showSaveDialog.value = true
}

/**
 * 检查规则名称是否存在
 */
const checkRuleNameExists = () => {
  if (!saveForm.name.trim()) {
    nameExistsWarning.value = false
    return
  }

  nameExistsWarning.value = savedRules.value.some(rule => rule.name === saveForm.name.trim())
}

/**
 * 保存密码规则
 */
const savePasswordRule = async () => {
  if (!saveForm.name.trim()) {
    ElMessage.error('请输入规则名称')
    return
  }

  saving.value = true

  try {
    let result

    if (configMode.value === 'general') {
      // 转换数据结构并保存通用规则
      const backendRule = convertGeneralRuleToBackend(generalRule)
      console.log('保存通用规则数据:', backendRule)
      result = await window.go.app.App.CreateGeneralPasswordRule(
        saveForm.name.trim(),
        saveForm.description.trim(),
        backendRule
      )
    } else {
      // 转换数据结构并保存自定义规则
      const backendRule = convertCustomRuleToBackend(customRule)
      console.log('保存自定义规则数据:', backendRule)
      result = await window.go.app.App.CreateCustomPasswordRule(
        saveForm.name.trim(),
        saveForm.description.trim(),
        backendRule
      )
    }

    // 如果设置为默认规则，调用设置默认规则的API
    if (saveForm.isDefault && result.id) {
      try {
        await window.go.app.App.SetPasswordRuleAsDefault(result.id, true)
        console.log('设置默认规则成功:', result.id)
      } catch (error) {
        console.error('设置默认规则失败:', error)
        ElMessage.warning('密码规则保存成功，但设置为默认规则失败')
      }
    }

    ElMessage.success('密码规则保存成功')
    showSaveDialog.value = false

    // 重新加载规则列表
    await loadPasswordRules()

    // 选择刚保存的规则并加载其配置
    selectedRuleId.value = result.id
    if (result.id) {
      await handleRuleChange(result.id)
    }

  } catch (error) {
    console.error('保存密码规则失败:', error)
    ElMessage.error('保存密码规则失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

/**
 * 删除密码规则
 * 20251017 陈凤庆 新增删除功能
 */
const handleDeleteRule = async () => {
  if (!selectedRuleId.value) {
    ElMessage.warning('请先选择要删除的密码规则')
    return
  }

  // 获取当前选中的规则信息
  const currentRule = savedRules.value.find(rule => rule.id === selectedRuleId.value)
  if (!currentRule) {
    ElMessage.error('未找到选中的密码规则')
    return
  }

  try {
    // 弹出确认对话框
    await ElMessageBox.confirm(
      `确定要删除密码规则"${currentRule.name}"吗？删除后无法恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    // 用户确认删除，开始删除操作
    deleting.value = true

    // 调用后端API删除密码规则
    await window.go.app.App.DeletePasswordRule(selectedRuleId.value)

    ElMessage.success(`密码规则"${currentRule.name}"删除成功`)

    // 重新加载规则列表
    await loadPasswordRules()

    // 清空选中的规则
    selectedRuleId.value = ''
    isCurrentRuleDefault.value = false

    // 重置配置模式为默认方式
    configMode.value = 'general'
    resetGeneralRule()

  } catch (error) {
    if (error === 'cancel') {
      // 用户取消删除
      return
    }

    console.error('删除密码规则失败:', error)
    const errorMsg = error?.message || error?.toString() || '未知错误'
    ElMessage.error('删除密码规则失败: ' + errorMsg)
  } finally {
    deleting.value = false
  }
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  visible.value = false
}

// 组件挂载时初始化
onMounted(async () => {
  window.addEventListener('resize', handleResize)
  await loadPasswordRules()
})
</script>

<style scoped>
/* 对话框样式 */
:deep(.el-dialog) {
  --el-dialog-padding-primary: 0;
  border-radius: 8px;
  box-shadow: 0 12px 32px 4px rgba(0, 0, 0, 0.36), 0 8px 20px rgba(0, 0, 0, 0.72);
  z-index: 9999 !important;
  min-height: 600px; /* 20251017 陈凤庆 增加最小高度从500px到600px */
  max-height: calc(100vh - 120px); /* 20251017 陈凤庆 调整最大高度，确保底部边距100px + 20px缓冲 */
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
  padding: 20px;
}

:deep(.el-dialog__headerbtn) {
  top: 15px;
  right: 15px;
}

.rule-settings {
  max-height: 65vh;
  min-height: 400px;
  overflow-y: auto;
}

.rule-selector-card,
.rule-config-card {
  margin-bottom: 16px; /* 减少底部边距 */
}

.form-label-with-switch {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.min-count-input {
  display: flex;
  align-items: center;
  justify-content: flex-start; /* 左对齐 */
  gap: 8px;
  margin-top: 6px; /* 减少顶部边距 */
  padding-left: 0; /* 确保左对齐 */
}

.custom-chars-input {
  margin-top: 8px;
}

.length-hint {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}

.password-test {
  display: flex;
  align-items: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

.general-config,
.custom-config {
  margin-top: 16px; /* 减少顶部边距 */
}

/* 调整表单项间距，使界面更紧凑 */
:deep(.el-form-item) {
  margin-bottom: 16px; /* 减少表单项底部边距 */
}

:deep(.el-form-item:last-child) {
  margin-bottom: 0; /* 最后一个表单项无底部边距 */
}

/* 调整卡片内边距 */
:deep(.el-card__body) {
  padding: 16px; /* 减少卡片内边距 */
}

/* 调整密码规则选择卡片的底部边距 */
.rule-selector-card {
  margin-bottom: 12px; /* 进一步减少底部边距 */
}

.name-warning {
  margin-top: 4px;
}

.default-rule-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: nowrap;
  gap: 10px;
}

.custom-rules-help {
  font-size: 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  flex-wrap: nowrap;
}

/* 表格操作按钮不换行 */
.table-actions {
  display: flex;
  gap: 8px;
  white-space: nowrap;
}

:deep(.el-table .el-table__cell) {
  white-space: nowrap;
}

:deep(.el-table .el-button + .el-button) {
  margin-left: 0;
}

/* 小屏幕适配 */
@media (max-width: 768px) {
  .rule-settings {
    max-height: 60vh;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  :deep(.el-dialog__body) {
    padding: 16px;
  }

  .custom-rules-help {
    font-size: 11px;
  }
}

@media (max-width: 480px) {
  .rule-settings {
    max-height: 55vh;
  }

  :deep(.el-dialog__body) {
    padding: 12px;
  }

  .dialog-footer {
    flex-direction: column;
    gap: 8px;
  }

  .dialog-footer .el-button {
    width: 100%;
  }
}

.default-rule-hint {
  margin-left: 10px;
  font-size: 12px;
  color: #909399;
}
</style>
