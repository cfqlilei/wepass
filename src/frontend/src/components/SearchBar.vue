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
 * @description 提供搜索输入框和搜索逻辑，支持防抖机制
 * @modify 20251018 陈凤庆 添加防抖机制，避免快速输入时发送过多请求
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
    },
    debounceDelay: {
      type: Number,
      default: 300 // 防抖延迟时间，单位毫秒
    }
  },
  emits: ['update:modelValue', 'search', 'clear'],
  setup(props, { emit }) {
    const keyword = ref(props.modelValue)
    let debounceTimer = null

    // 监听 modelValue 变化
    watch(() => props.modelValue, (newValue) => {
      keyword.value = newValue
    })

    /**
     * 处理输入事件 - 带防抖
     */
    const handleInput = () => {
      emit('update:modelValue', keyword.value)

      // 清除之前的防抖计时器
      if (debounceTimer) {
        clearTimeout(debounceTimer)
      }

      // 设置新的防抖计时器
      debounceTimer = setTimeout(() => {
        emit('search', keyword.value)
      }, props.debounceDelay)
    }

    /**
     * 处理清除事件
     */
    const handleClear = () => {
      keyword.value = ''
      emit('update:modelValue', '')

      // 清除防抖计时器
      if (debounceTimer) {
        clearTimeout(debounceTimer)
      }

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

