<template>
  <el-dialog
    v-model="visible"
    :title="$t('passwordGenerator.title')"
    width="500px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      label-width="120px"
    >
      <el-form-item :label="$t('passwordGenerator.selectRule')">
        <el-select
          v-model="formData.ruleType"
          :placeholder="$t('passwordGenerator.selectRulePlaceholder')"
          style="width: 100%"
          @change="handleRuleTypeChange"
        >
          <el-option
            v-for="rule in savedRules"
            :key="rule.id"
            :label="rule.name"
            :value="rule.id"
          />
          <el-option
            :label="$t('passwordGenerator.generalRule')"
            value="general"
          />
          <el-option
            :label="$t('passwordGenerator.customRule')"
            value="custom"
          />
        </el-select>
      </el-form-item>

      <!-- 通用密码规则 -->
      <template v-if="formData.ruleType === 'general'">
        <el-form-item :label="t('passwordGenerator.includeUppercase')">
          <el-switch v-model="formData.includeUppercase" />
        </el-form-item>
        <el-form-item :label="t('passwordGenerator.includeLowercase')">
          <el-switch v-model="formData.includeLowercase" />
        </el-form-item>
        <el-form-item :label="t('passwordGenerator.includeNumbers')">
          <el-switch v-model="formData.includeNumbers" />
        </el-form-item>
        <el-form-item :label="t('passwordGenerator.includeSpecialChars')">
          <el-switch v-model="formData.includeSpecialChars" />
        </el-form-item>
        <el-form-item :label="t('passwordGenerator.passwordLength')">
          <el-input-number
            v-model="formData.length"
            :min="4"
            :max="128"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="t('passwordGenerator.customSpecialChars')">
          <el-input
            v-model="formData.customSpecialChars"
            :placeholder="t('passwordGenerator.defaultSpecialChars')"
          />
        </el-form-item>
      </template>

      <!-- 自定义密码规则 -->
      <template v-if="formData.ruleType === 'custom'">
        <el-form-item label="密码模式">
          <el-input
            v-model="formData.customPattern"
            type="textarea"
            :rows="3"
            placeholder="请输入密码模式，例如：Aaa111"
          />
        </el-form-item>
        <div class="pattern-help">
          <el-text size="small" type="info">
            模式说明：a=小写字母, A=大小写字母, U=大写字母, d=数字, h=小写十六进制, H=大写十六进制, 
            l=小写字母, L=大小写字母, u=大写字母, v=小写元音, V=大小写元音, Z=大写元音, 
            c=小写辅音, C=大小写辅音, z=大写辅音, p=标点符号, b=括号, s=特殊字符, S=可打印ASCII, 
            x=Latin-1补充, \=转义字符, {n}=重复n次, [...]=自定义字符集
          </el-text>
        </div>
      </template>

      <!-- 生成的密码 -->
      <el-form-item :label="$t('passwordGenerator.generatedPassword')">
        <el-input
          v-model="generatedPassword"
          readonly
          :placeholder="$t('passwordGenerator.clickToGenerate')"
        >
          <template #append>
            <el-button
              :icon="RefreshRight"
              @click="generatePassword"
              :title="$t('passwordGenerator.generatePassword')"
            />
            <el-button
              :icon="CopyDocument"
              @click="copyPassword"
              :title="$t('passwordGenerator.copyPassword')"
            />
          </template>
        </el-input>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">{{ t('common.close') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { CopyDocument, RefreshRight } from '@element-plus/icons-vue'

const { t } = useI18n()

/**
 * 密码生成器组件
 * @author 陈凤庆
 * @description 根据规则生成密码的对话框组件
 */

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'generated'])

// 响应式数据
const formRef = ref()
const visible = ref(false)
const generatedPassword = ref('')

// 表单数据
const formData = reactive({
  ruleType: 'general',
  includeUppercase: true,
  includeLowercase: true,
  includeNumbers: true,
  includeSpecialChars: false,
  length: 12,
  customSpecialChars: '',
  customPattern: ''
})

// 保存的密码规则
const savedRules = ref([])

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    // 对话框打开时加载最新的密码规则
    loadPasswordRules()
  }
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

/**
 * 加载密码规则
 */
const loadPasswordRules = async () => {
  try {
    const rules = await window.go.app.App.GetAllPasswordRules()
    savedRules.value = rules || []
    console.log('加载密码规则成功:', rules)
  } catch (error) {
    console.error('加载密码规则失败:', error)
    ElMessage.error('加载密码规则失败')
  }
}

/**
 * 处理规则类型变化
 */
const handleRuleTypeChange = () => {
  generatedPassword.value = ''
}

/**
 * 生成密码
 */
const generatePassword = async () => {
  try {
    let password = ''

    if (formData.ruleType === 'general') {
      password = generateGeneralPassword()
    } else if (formData.ruleType === 'custom') {
      password = generateCustomPassword()
    } else {
      // 使用保存的规则
      password = await generateSavedRulePassword()
    }

    if (!password) {
      ElMessage.error(t('passwordGenerator.generateFailed'))
      return
    }

    generatedPassword.value = password
    ElMessage.success(t('passwordGenerator.generateSuccess'))
  } catch (error) {
    console.error('生成密码失败:', error)
    const errorMsg = error?.message || error?.toString() || '未知错误'
    ElMessage.error('生成密码失败: ' + errorMsg)
  }
}

/**
 * 生成通用密码
 */
const generateGeneralPassword = () => {
  let charset = ''
  
  if (formData.includeLowercase) {
    charset += 'abcdefghijklmnopqrstuvwxyz'
  }
  if (formData.includeUppercase) {
    charset += 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  }
  if (formData.includeNumbers) {
    charset += '0123456789'
  }
  if (formData.includeSpecialChars) {
    charset += formData.customSpecialChars || '!@#$%^&*()_+-=[]{}|;:,.<>?'
  }
  
  if (!charset) {
    ElMessage.error(t('passwordGenerator.selectCharType'))
    return ''
  }
  
  let password = ''
  for (let i = 0; i < formData.length; i++) {
    password += charset.charAt(Math.floor(Math.random() * charset.length))
  }
  
  return password
}

/**
 * 生成自定义模式密码
 */
const generateCustomPassword = () => {
  if (!formData.customPattern) {
    ElMessage.error(t('passwordGenerator.enterPattern'))
    return ''
  }
  
  // 简化的自定义模式实现
  const pattern = formData.customPattern
  let password = ''
  
  for (let i = 0; i < pattern.length; i++) {
    const char = pattern[i]
    switch (char) {
      case 'a':
        password += String.fromCharCode(97 + Math.floor(Math.random() * 26))
        break
      case 'A':
        password += Math.random() > 0.5 
          ? String.fromCharCode(65 + Math.floor(Math.random() * 26))
          : String.fromCharCode(97 + Math.floor(Math.random() * 26))
        break
      case 'U':
        password += String.fromCharCode(65 + Math.floor(Math.random() * 26))
        break
      case 'd':
        password += Math.floor(Math.random() * 10).toString()
        break
      default:
        password += char
    }
  }
  
  return password
}

/**
 * 生成保存规则密码
 */
const generateSavedRulePassword = async () => {
  try {
    // 使用后端API生成密码
    const password = await window.go.app.App.GeneratePasswordByRule(formData.ruleType)
    return password
  } catch (error) {
    console.error('使用保存规则生成密码失败:', error)
    // 如果失败，使用默认规则
    return generateGeneralPassword()
  }
}

/**
 * 复制密码
 */
const copyPassword = async () => {
  if (!generatedPassword.value) {
    ElMessage.error(t('passwordGenerator.generateFirst'))
    return
  }
  
  try {
    await navigator.clipboard.writeText(generatedPassword.value)
    ElMessage.success(t('success.passwordCopied'))
  } catch (error) {
    ElMessage.error(t('error.copyFailed'))
  }
}

/**
 * 使用此密码
 */
const usePassword = () => {
  if (!generatedPassword.value) {
    ElMessage.error(t('passwordGenerator.generateFirst'))
    return
  }
  
  emit('generated', generatedPassword.value)
  handleClose()
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  visible.value = false
  generatedPassword.value = ''
}
</script>

<style scoped>
.pattern-help {
  margin-top: -10px;
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
