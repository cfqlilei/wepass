<!--
/**
 * 账号详情组件测试页面
 * @author 陈凤庆
 * @date 2025-01-27
 * @modify 20251002 陈凤庆 密码详情改为账号详情，保持命名一致性
 * @description 用于测试账号详情组件的修改效果
 */
-->
<template>
  <div class="test-container">
    <h2>{{ $t('passwordDetailTest.title') }}</h2>
    <el-button type="primary" @click="showDetail">{{ $t('passwordDetailTest.showDetail') }}</el-button>

    <AccountDetail
      v-model="showDialog"
      :account-data="testData"
      :is-new="false"
      @save="handleSave"
      @delete="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import AccountDetail from '../components/AccountDetail.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

/**
 * 测试数据和方法
 * @modify 20251002 陈凤庆 密码详情改为账号详情，保持命名一致性
 */
const showDialog = ref(false)

const testData = ref({
  id: 1,
  title: '测试网站',
  username: 'testuser@example.com',
  password: 'TestPassword123!',
  url: 'https://example.com',
  type: '网站',
  notes: '这是一个测试备注，用于验证文本选择和复制功能是否正常工作。',
  group_id: 1,
  icon: ''
})

const showDetail = () => {
  showDialog.value = true
}

const handleSave = (data) => {
  console.log('保存数据:', data)
}

const handleDelete = (id) => {
  console.log('删除ID:', id)
}
</script>

<style scoped>
.test-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

h2 {
  margin-bottom: 20px;
  color: #333;
}
</style>
