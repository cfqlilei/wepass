/**
 * 语言切换工具函数
 * @author 陈凤庆
 * @date 2025-10-05
 * @description 处理语言切换和Element Plus语言包更新
 */

import { ElConfigProvider } from 'element-plus'
import { setLocale, getElementPlusLocale } from '@/i18n'

// 导入Element Plus语言包
import zhCn from "element-plus/dist/locale/zh-cn.mjs";
import zhTw from "element-plus/dist/locale/zh-tw.mjs";
import en from "element-plus/dist/locale/en.mjs";
import ru from "element-plus/dist/locale/ru.mjs";
import es from "element-plus/dist/locale/es.mjs";
import fr from "element-plus/dist/locale/fr.mjs";
import de from "element-plus/dist/locale/de.mjs";

// Element Plus语言包映射
const elementPlusLocales = {
  zhCn,
  zhTw,
  en,
  ru,
  es,
  fr,
  de
}

/**
 * 切换应用语言
 * @param {string} locale 语言代码
 * @param {Function} reloadCallback 重新加载回调函数
 */
export function switchLanguage(locale, reloadCallback) {
  try {
    // 更新i18n语言
    setLocale(locale)
    
    // 获取对应的Element Plus语言包
    const elementPlusLocaleName = getElementPlusLocale(locale)
    const elementPlusLocale = elementPlusLocales[elementPlusLocaleName] || zhCn
    
    // 更新Element Plus语言包
    updateElementPlusLocale(elementPlusLocale)
    
    // 更新HTML lang属性
    document.documentElement.lang = locale
    
    // 如果提供了重新加载回调，执行它
    if (typeof reloadCallback === 'function') {
      reloadCallback()
    }
    
    return true
  } catch (error) {
    console.error('切换语言失败:', error)
    return false
  }
}

/**
 * 更新Element Plus语言包
 * @param {Object} locale Element Plus语言包对象
 */
function updateElementPlusLocale(locale) {
  // 这里需要通过全局配置来更新Element Plus的语言包
  // 由于Element Plus在初始化后不容易动态更改语言包，
  // 我们将在设置保存后重新加载页面来确保语言切换生效
}

/**
 * 获取Element Plus语言包
 * @param {string} locale 语言代码
 * @returns {Object} Element Plus语言包对象
 */
export function getElementPlusLocaleObject(locale) {
  const elementPlusLocaleName = getElementPlusLocale(locale)
  return elementPlusLocales[elementPlusLocaleName] || zhCn
}
