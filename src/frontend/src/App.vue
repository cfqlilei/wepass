<template>
  <div id="app" class="app-container">
    <router-view />
  </div>
</template>

<script>
/**
 * 应用根组件
 * @author 陈凤庆
 * @description Vue.js 应用的根组件，包含路由视图
 */
export default {
  name: 'App',
}
</script>

<style>
/* 全局样式 */
.app-container {
  width: 100%;
  /* 20250928 陈凤庆 修改高度为自适应内容，避免底部空白 */
  min-height: 100vh;
  height: 100vh;
  /* 20250127 陈凤庆 修改为hidden，防止不必要的滚动条 */
  overflow: hidden;
  font-family: "PingFang SC", "Microsoft YaHei", "Noto Sans CJK SC", sans-serif;
  background-color: #ffffff;
  color: #212529;
  font-size: 12px;
  line-height: 1.5;
  /* 20250127 陈凤庆 使用flex布局，确保内容填满整个容器 */
  display: flex;
  flex-direction: column;
  justify-content: stretch;
  align-items: stretch;
  /* 20250127 陈凤庆 添加box-sizing确保正确计算尺寸 */
  box-sizing: border-box;
  /* 20250127 陈凤庆 确保最大宽度不超出视口 */
  max-width: 100vw;
}

/* Element Plus 样式覆盖 */
.el-button {
  font-family: inherit;
}

.el-input__inner {
  font-family: inherit;
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 禁用文本选择 */
.no-select {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
</style>