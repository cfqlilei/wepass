<template>
  <div v-if="show" class="search-bar">
    <el-input
      v-model="keyword"
      :placeholder="$t('main.searchPlaceholder')"
      :prefix-icon="Search"
      clearable
      @input="handleInput"
      @clear="handleClear"
    />
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import { Search } from '@element-plus/icons-vue'

/**
 * 搜索栏组件
 * @author 陈凤庆
 * @date 20251001
 * @description 提供搜索输入框和搜索逻辑
 */
export default {
  name: 'SearchBar',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    modelValue: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue', 'search', 'clear'],
  setup(props, { emit }) {
    const keyword = ref(props.modelValue)

    // 监听 modelValue 变化
    watch(() => props.modelValue, (newValue) => {
      keyword.value = newValue
    })

    /**
     * 处理输入事件
     */
    const handleInput = () => {
      emit('update:modelValue', keyword.value)
      emit('search', keyword.value)
    }

    /**
     * 处理清除事件
     */
    const handleClear = () => {
      keyword.value = ''
      emit('update:modelValue', '')
      emit('clear')
    }

    return {
      keyword,
      Search,
      handleInput,
      handleClear
    }
  }
}
</script>

<style scoped>
.search-bar {
  padding: 8px 16px;
  border-bottom: 1px solid #dee2e6;
  background-color: #f8f9fa;
}
</style>

