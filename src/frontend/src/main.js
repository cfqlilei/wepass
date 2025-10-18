import { createApp } from "vue";
import { createPinia } from "pinia";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

// 导入Element Plus语言包
import zhCn from "element-plus/dist/locale/zh-cn.mjs";
import zhTw from "element-plus/dist/locale/zh-tw.mjs";
import en from "element-plus/dist/locale/en.mjs";
import ru from "element-plus/dist/locale/ru.mjs";
import es from "element-plus/dist/locale/es.mjs";
import fr from "element-plus/dist/locale/fr.mjs";
import de from "element-plus/dist/locale/de.mjs";

import App from "./App.vue";
import router from "./router";
import { apiService } from "./services/api";
import i18n, { getCurrentLocale, getElementPlusLocale } from "./i18n";

/**
 * Vue.js 应用主入口
 * @author 陈凤庆
 * @description 初始化 Vue 应用和相关插件，包括国际化支持
 * @modify 2025-10-05 陈凤庆 添加多语言支持
 */

// Element Plus语言包映射
const elementPlusLocales = {
  zhCn,
  zhTw,
  en,
  ru,
  es,
  fr,
  de,
};

// 获取当前语言对应的Element Plus语言包
function getCurrentElementPlusLocale() {
  const currentLocale = getCurrentLocale();
  const elementPlusLocaleName = getElementPlusLocale(currentLocale);
  return elementPlusLocales[elementPlusLocaleName] || zhCn;
}

const app = createApp(App);

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

// 使用插件
app.use(createPinia());
app.use(router);
app.use(i18n);
app.use(ElementPlus, {
  // Element Plus 配置 - 动态语言包
  locale: getCurrentElementPlusLocale(),
});

// 20250928 陈凤庆 全局注册API服务，解决组件中无法访问apiService的问题
app.config.globalProperties.$apiService = apiService;
app.provide("apiService", apiService);

// 20251005 陈凤庆 将i18n实例挂载到window对象，供其他模块使用
window.i18n = i18n;

// 挂载应用
app.mount("#app");
