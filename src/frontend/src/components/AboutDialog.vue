<template>
  <el-dialog
    v-model="visible"
    :title="$t('about.title')"
    width="500px"
    :before-close="handleClose"
  >
    <div class="about-content">
      <!-- 应用信息 -->
      <div class="app-info">
        <div class="app-logo">
          <el-icon class="logo-icon"><Key /></el-icon>
          <h2>{{ appInfo.name }}</h2>
        </div>
        <p class="app-description">{{ $t('about.description') }}</p>
        <p class="version">{{ $t('about.version') }}: v{{ appInfo.version }}</p>
        <p class="build-date">{{ $t('about.buildDate') }}: {{ appInfo.buildDate }}</p>
      </div>

      <!-- 项目链接 -->
      <el-divider />
      <div class="project-links">
        <h3>{{ $t('about.projectLinks') }}</h3>
        <div class="link-item">
          <el-icon><Link /></el-icon>
          <span>GitHub: </span>
          <el-link
            href="https://github.com/cfqlilei/wepass"
            target="_blank"
            type="primary"
          >
            https://github.com/cfqlilei/wepass
          </el-link>
        </div>
        <!-- <div class="link-item">
          <el-icon><Link /></el-icon>
          <span>Gitee: </span>
          <el-link
            href="https://gitee.com/cfqlilei/wepassword"
            target="_blank"
            type="primary"
          >
            https://gitee.com/cfqlilei/wepassword
          </el-link>
        </div> -->
      </div>

      <!-- 作者信息 -->
      <el-divider />
      <div class="author-info">
        <h3>{{ $t('about.authorInfo') }}</h3>
        <div class="info-item">
          <el-icon><User /></el-icon>
          <span>{{ $t('about.author') }}: 微易软件 陈凤庆</span>
        </div>
        <div class="info-item">
          <el-icon><Message /></el-icon>
          <span>{{ $t('about.email') }}: cfq@wesoftcn.com</span>
        </div>
      </div>

      <!-- 支持信息 -->
      <el-divider />
      <div class="support-info">
        <h3>{{ $t('about.technicalSupport') }}</h3>
        <div class="support-item">
          <el-icon><QuestionFilled /></el-icon>
          <span>{{ $t('about.contactInfo') }}</span>
        </div>
        <ul class="support-list">
          <!-- <li>GitHub Issues: 在项目仓库中提交问题</li> -->
          <li>{{ $t('about.emailContact') }}: cfq@wesoftcn.com</li>
          <li>{{ $t('about.qqGroup') }}: {{ $t('about.qqGroupTbd') }}</li>
        </ul>
      </div>

      <!-- 授权信息 -->
      <el-divider />
      <div class="license-info">
        <h3>{{ $t('about.licenseInfo') }}</h3>
        <div class="license-item">
          <el-icon><Document /></el-icon>
          <span>{{ $t('about.usageAgreement') }}: {{ $t('about.limitedFree') }}</span>
          <!-- <span>开源协议: MIT License</span> -->
        </div>
        <!-- <div class="license-item">
          <el-icon><Document /></el-icon>
          <span>开源协议: MIT License</span>
          <span>开源协议: MIT License</span>
        </div> -->
        <div class="license-item">
          <el-icon><Lock /></el-icon>
          <span>{{ $t('about.copyright') }}: © 2025 {{ $t('about.authorName') }}. All rights reserved.</span>
        </div>
        <!-- <p class="license-description">
          本软件基于MIT开源协议发布，您可以自由使用、修改和分发本软件。
          使用本软件时请遵守相关法律法规，作者不承担任何使用风险。
        </p> -->
      </div>

      <!-- 技术栈 -->
      <el-divider />
      <div class="tech-stack">
        <h3>{{ $t('about.techStack') }}</h3>
        <div class="tech-tags">
          <el-tag type="primary">Wails v2</el-tag>
          <el-tag type="success">Go 1.21+</el-tag>
          <el-tag type="info">Vue.js 3</el-tag>
          <el-tag type="warning">Element Plus</el-tag>
          <el-tag>SQLite</el-tag>
          <el-tag>AES-256-GCM</el-tag>
        </div>
      </div>

      <!-- 特别感谢 -->
      <el-divider />
      <div class="thanks">
        <h3>{{ $t('about.specialThanks') }}</h3>
        <p>{{ $t('about.thanksDescription') }}</p>
        <ul class="thanks-list">
          <li>Wails - {{ $t('about.wailsDescription') }}</li>
          <li>Vue.js - {{ $t('about.vueDescription') }}</li>
          <li>Element Plus - {{ $t('about.elementPlusDescription') }}</li>
          <li>Go - {{ $t('about.goDescription') }}</li>
          <li>SQLite - {{ $t('about.sqliteDescription') }}</li>
        </ul>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="handleClose">{{ $t('common.confirm') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Key, Link, User, Message, QuestionFilled,
  Document, Lock
} from '@element-plus/icons-vue'
import { apiService } from '../services/api'

const { t } = useI18n()

/**
 * 关于对话框组件
 * @author 陈凤庆
 * @description 显示应用信息、作者信息、授权信息等
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
const appInfo = ref({
  name: 'wepass',
  version: '1.0.0',
  buildDate: '2025-10-01'
})

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    loadAppInfo()
  }
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

/**
 * 加载应用信息
 */
const loadAppInfo = async () => {
  try {
    const info = await apiService.getAppInfo()
    appInfo.value = info
  } catch (error) {
    console.error('Failed to load app info:', error)
  }
}

/**
 * 关闭对话框
 */
const handleClose = () => {
  visible.value = false
}
</script>

<style scoped>
.about-content {
  max-height: 500px;
  overflow-y: auto;
  padding: 10px 0;
}

.app-info {
  text-align: center;
  margin-bottom: 20px;
}

.app-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 10px;
}

.logo-icon {
  font-size: 32px;
  color: #007bff;
}

.app-logo h2 {
  margin: 0;
  color: #007bff;
  font-size: 28px;
}

.app-description {
  color: #666;
  margin: 10px 0;
}

.version {
  color: #999;
  font-size: 14px;
}

.project-links h3,
.author-info h3,
.support-info h3,
.license-info h3,
.tech-stack h3,
.thanks h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 16px;
}

.link-item,
.info-item,
.support-item,
.license-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.link-item .el-icon,
.info-item .el-icon,
.support-item .el-icon,
.license-item .el-icon {
  font-size: 16px;
  color: #666;
}

.support-list,
.thanks-list {
  margin: 10px 0;
  padding-left: 20px;
}

.support-list li,
.thanks-list li {
  margin-bottom: 5px;
  color: #666;
}

.license-description {
  margin-top: 10px;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
  font-size: 12px;
  color: #666;
  line-height: 1.5;
}

.tech-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: center;
}
</style>
