import { createI18n } from 'vue-i18n'

// 导入语言文件
import zhCN from './locales/zh-CN.js'
import zhTW from './locales/zh-TW.js'
import enUS from './locales/en-US.js'
import ruRU from './locales/ru-RU.js'
import esES from './locales/es-ES.js'
import frFR from './locales/fr-FR.js'
import deDE from './locales/de-DE.js'

/**
 * 国际化配置
 * @author 陈凤庆
 * @date 2025-10-05
 * @description 配置vue-i18n多语言支持
 */

// 支持的语言列表
export const supportedLocales = [
  { code: 'zh-CN', name: '简体中文', elementPlusLocale: 'zhCn' },
  { code: 'zh-TW', name: '繁體中文', elementPlusLocale: 'zhTw' },
  { code: 'en-US', name: 'English', elementPlusLocale: 'en' },
  { code: 'ru-RU', name: 'Русский', elementPlusLocale: 'ru' },
  { code: 'es-ES', name: 'Español', elementPlusLocale: 'es' },
  { code: 'fr-FR', name: 'Français', elementPlusLocale: 'fr' },
  { code: 'de-DE', name: 'Deutsch', elementPlusLocale: 'de' }
]

// 获取默认语言
function getDefaultLocale() {
  // 优先从localStorage获取用户设置的语言
  const savedLocale = localStorage.getItem('app-language')
  if (savedLocale && supportedLocales.find(locale => locale.code === savedLocale)) {
    return savedLocale
  }
  
  // 其次从浏览器语言获取
  const browserLocale = navigator.language || navigator.userLanguage
  const matchedLocale = supportedLocales.find(locale => 
    locale.code === browserLocale || locale.code.split('-')[0] === browserLocale.split('-')[0]
  )
  
  if (matchedLocale) {
    return matchedLocale.code
  }
  
  // 默认返回简体中文
  return 'zh-CN'
}

// 创建i18n实例
const i18n = createI18n({
  legacy: false, // 使用Composition API模式
  locale: getDefaultLocale(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'zh-TW': zhTW,
    'en-US': enUS,
    'ru-RU': ruRU,
    'es-ES': esES,
    'fr-FR': frFR,
    'de-DE': deDE
  }
})

/**
 * 切换语言
 * @param {string} locale 语言代码
 */
export function setLocale(locale) {
  if (supportedLocales.find(l => l.code === locale)) {
    i18n.global.locale.value = locale
    localStorage.setItem('app-language', locale)
    document.documentElement.lang = locale
  }
}

/**
 * 获取当前语言
 * @returns {string} 当前语言代码
 */
export function getCurrentLocale() {
  return i18n.global.locale.value
}

/**
 * 获取语言对应的Element Plus语言包名称
 * @param {string} locale 语言代码
 * @returns {string} Element Plus语言包名称
 */
export function getElementPlusLocale(locale) {
  const supportedLocale = supportedLocales.find(l => l.code === locale)
  return supportedLocale ? supportedLocale.elementPlusLocale : 'zhCn'
}

export default i18n
