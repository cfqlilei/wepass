<template>
  <div
    class="resizer"
    @mousedown="handleMouseDown"
  ></div>
</template>

<script>
/**
 * 可拖动分隔栏组件
 * @author 陈凤庆
 * @date 20251001
 * @description 处理侧边栏拖动调整大小的逻辑
 */
export default {
  name: 'Resizer',
  emits: ['start-resize'],
  setup(_, { emit }) {
    /**
     * 处理鼠标按下事件
     */
    const handleMouseDown = (event) => {
      emit('start-resize', event)
    }

    return {
      handleMouseDown
    }
  }
}
</script>

<style scoped>
.resizer {
  width: 4px;
  /* 20251003 陈凤庆 隐藏分隔线，设置为透明 */
  background-color: transparent;
  cursor: col-resize;
  position: relative;
  transition: background-color 0.2s;
}

.resizer:hover {
  /* 20251003 陈凤庆 悬停时显示分隔线 */
  background-color: #007bff;
}
</style>

