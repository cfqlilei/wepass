import { defineStore } from "pinia";
import { ref, computed } from "vue";

/**
 * 应用状态管理
 * @author 陈凤庆
 * @description 管理应用的全局状态，包括分组、页签、搜索等
 */
export const useAppStore = defineStore("app", () => {
  // 状态
  const groups = ref([]);
  const currentGroupId = ref(null); // 20250101 陈凤庆 初始值改为null，避免找不到分组
  const tabs = ref([]);
  const currentTabId = ref(null); // 20250101 陈凤庆 初始值改为null，避免找不到标签
  const passwords = ref([]);
  const searchKeyword = ref("");
  const isSearching = ref(false);
  const searchResults = ref([]);

  // 20251019 陈凤庆 新增：分组页签记忆功能
  // 存储每个分组最后选中的页签ID，格式：{ groupId: tabId }
  const groupTabMemory = ref({});

  // 计算属性
  const currentGroup = computed(() => {
    return groups.value.find((g) => g.id === currentGroupId.value);
  });

  const currentTab = computed(() => {
    // 20250101 陈凤庆 如果currentTabId为null，返回null
    if (currentTabId.value === null || currentTabId.value === undefined) {
      return null;
    }
    return tabs.value.find((t) => t.id === currentTabId.value);
  });

  /**
   * 根据当前标签筛选账号项
   * 20250101 陈凤庆 添加标签筛选逻辑
   * 20250101 陈凤庆 确保始终返回数组，避免null.length错误
   * 20251002 陈凤庆 修改为按标签ID筛选，而不是使用filter规则
   */
  const filteredPasswords = computed(() => {
    // 如果正在搜索,返回搜索结果
    if (isSearching.value) {
      return searchResults.value || [];
    }

    // 如果没有选中标签,返回所有账号项
    if (!currentTab.value) {
      return passwords.value || [];
    }

    // 20251002 陈凤庆 根据标签ID筛选账号项
    const tab = currentTab.value;
    let filtered = passwords.value || [];

    // 按标签ID筛选账号项
    if (tab.id) {
      filtered = filtered.filter(
        (p) => p.typeid === tab.id || p.tab_id === tab.id
      );
    }

    return filtered;
  });

  // 动作
  /**
   * 设置分组列表
   * @param {Array} groupList 分组列表
   */
  const setGroups = (groupList) => {
    groups.value = groupList;
  };

  /**
   * 设置当前分组
   * @param {string} groupId 分组ID
   * @modify 20251002 陈凤庆 使用字符串ID，避免JavaScript精度丢失
   */
  const setCurrentGroup = (groupId) => {
    currentGroupId.value = groupId;
  };

  /**
   * 设置页签列表
   * @param {Array} tabList 页签列表
   */
  const setTabs = (tabList) => {
    tabs.value = tabList;
  };

  /**
   * 设置当前页签
   * @param {number} tabId 页签ID
   */
  const setCurrentTab = (tabId) => {
    currentTabId.value = tabId;

    // 20251019 陈凤庆 记录当前分组的页签选择
    if (currentGroupId.value && tabId) {
      rememberGroupTab(currentGroupId.value, tabId);
    }
  };

  /**
   * 设置账号列表
   * @param {Array} passwordList 账号列表
   * 20250101 陈凤庆 确保不会设置null值，避免null.length错误
   */
  const setPasswords = (passwordList) => {
    passwords.value = passwordList || [];
  };

  /**
   * 添加账号
   * @param {Object} password 账号
   */
  const addPassword = (password) => {
    passwords.value.push(password);
  };

  /**
   * 更新账号
   * @param {Object} password 账号
   */
  const updatePassword = (password) => {
    const index = passwords.value.findIndex((p) => p.id === password.id);
    if (index !== -1) {
      passwords.value[index] = password;
    }
  };

  /**
   * 删除账号
   * @param {number} passwordId 账号ID
   */
  const removePassword = (passwordId) => {
    const index = passwords.value.findIndex((p) => p.id === passwordId);
    if (index !== -1) {
      passwords.value.splice(index, 1);
    }
  };

  /**
   * 设置搜索关键词
   * @param {string} keyword 搜索关键词
   */
  const setSearchKeyword = (keyword) => {
    searchKeyword.value = keyword;
    isSearching.value = keyword.length > 0;
  };

  /**
   * 设置搜索结果
   * @param {Array} results 搜索结果
   */
  const setSearchResults = (results) => {
    searchResults.value = results;
  };

  /**
   * 清除搜索
   */
  const clearSearch = () => {
    searchKeyword.value = "";
    isSearching.value = false;
    searchResults.value = [];
  };

  /**
   * 记录分组的页签选择
   * @param {string} groupId 分组ID
   * @param {string} tabId 页签ID
   * @author 20251019 陈凤庆 新增分组页签记忆功能
   */
  const rememberGroupTab = (groupId, tabId) => {
    if (!groupId || !tabId) return;

    groupTabMemory.value[groupId] = tabId;

    // 持久化到localStorage
    try {
      localStorage.setItem(
        "groupTabMemory",
        JSON.stringify(groupTabMemory.value)
      );
      console.log(`[Store] 记录分组${groupId}的页签选择: ${tabId}`);
    } catch (error) {
      console.error("[Store] 保存页签记忆失败:", error);
    }
  };

  /**
   * 获取分组记忆的页签ID
   * @param {string} groupId 分组ID
   * @returns {string|null} 页签ID或null
   * @author 20251019 陈凤庆 新增分组页签记忆功能
   */
  const getRememberedTab = (groupId) => {
    if (!groupId) return null;

    const rememberedTabId = groupTabMemory.value[groupId];
    console.log(`[Store] 获取分组${groupId}的记忆页签: ${rememberedTabId}`);
    return rememberedTabId || null;
  };

  /**
   * 从localStorage加载页签记忆
   * @author 20251019 陈凤庆 新增分组页签记忆功能
   */
  const loadGroupTabMemory = () => {
    try {
      const saved = localStorage.getItem("groupTabMemory");
      if (saved) {
        groupTabMemory.value = JSON.parse(saved);
        console.log("[Store] 已加载页签记忆:", groupTabMemory.value);
      }
    } catch (error) {
      console.error("[Store] 加载页签记忆失败:", error);
      groupTabMemory.value = {};
    }
  };

  /**
   * 清除分组的页签记忆
   * @param {string} groupId 分组ID（可选，不传则清除所有）
   * @author 20251019 陈凤庆 新增分组页签记忆功能
   */
  const clearGroupTabMemory = (groupId = null) => {
    if (groupId) {
      delete groupTabMemory.value[groupId];
      console.log(`[Store] 已清除分组${groupId}的页签记忆`);
    } else {
      groupTabMemory.value = {};
      console.log("[Store] 已清除所有页签记忆");
    }

    // 更新localStorage
    try {
      localStorage.setItem(
        "groupTabMemory",
        JSON.stringify(groupTabMemory.value)
      );
    } catch (error) {
      console.error("[Store] 保存页签记忆失败:", error);
    }
  };

  return {
    // 状态
    groups,
    currentGroupId,
    tabs,
    currentTabId,
    passwords,
    searchKeyword,
    isSearching,
    searchResults,
    groupTabMemory,

    // 计算属性
    currentGroup,
    currentTab,
    filteredPasswords,

    // 动作
    setGroups,
    setCurrentGroup,
    setTabs,
    setCurrentTab,
    setPasswords,
    addPassword,
    updatePassword,
    removePassword,
    setSearchKeyword,
    setSearchResults,
    clearSearch,

    // 20251019 陈凤庆 新增：页签记忆功能
    rememberGroupTab,
    getRememberedTab,
    loadGroupTabMemory,
    clearGroupTabMemory,
  };
});
