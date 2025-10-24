/**
 * Wails API 服务
 * @author 陈凤庆
 * @description 封装与 Go 后端的通信接口
 */

/**
 * API 服务类
 */
class ApiService {
  /**
   * 20250127 陈凤庆 动态检查Wails API是否可用
   * @returns {Object|null} Wails API对象或null
   */
  getWailsAPI() {
    if (
      typeof window !== "undefined" &&
      window.go &&
      window.go.app &&
      window.go.app.App
    ) {
      return window.go.app.App;
    }
    return null;
  }
  /**
   * 检查密码库是否存在
   * @param {string} vaultPath 密码库路径
   * @returns {Promise<boolean>} 是否存在
   */
  async checkVaultExists(vaultPath) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，使用模拟数据");
      return false;
    }
    return await wailsAPI.CheckVaultExists(vaultPath);
  }

  /**
   * 创建新密码库
   * @param {string} vaultName 密码库名称（无需后缀）
   * @param {string} password 登录密码
   * @param {string} language 语言代码（可选，默认为zh-CN）
   * @returns {Promise<string>} 创建的密码库完整路径
   */
  async createVault(vaultName, password, language = "zh-CN") {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return "";
    }
    return await wailsAPI.CreateVault(vaultName, password, language);
  }

  /**
   * 打开密码库
   * @param {string} vaultPath 密码库路径
   * @param {string} password 登录密码
   * @returns {Promise<void>}
   */
  async openVault(vaultPath, password) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.OpenVault(vaultPath, password);
  }

  /**
   * 检查密码库是否已打开
   * @returns {Promise<boolean>} 是否已打开
   * @author 陈凤庆
   * @description 20251003 用于前端检查登录状态，避免重复登录
   */
  async isVaultOpened() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return false;
    }
    return await wailsAPI.IsVaultOpened();
  }

  /**
   * 关闭密码库
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @description 20251003 关闭密码库并清理状态，用于退出登录
   */
  async closeVault() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.CloseVault();
  }

  /**
   * 获取当前密码库路径
   * @returns {Promise<string>} 当前密码库路径
   * @author 陈凤庆
   * @date 2025-10-17
   * @description 获取当前打开的密码库文件路径
   */
  async getCurrentVaultPath() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return "";
    }
    return await wailsAPI.GetCurrentVaultPath();
  }

  /**
   * 打开密码库所在目录
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @description 20251002 打开当前密码库文件所在的目录
   */
  async openVaultDirectory() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API not initialized");
    }
    return await wailsAPI.OpenVaultDirectory();
  }

  /**
   * 获取最近使用的密码库列表
   * @returns {Promise<Array<string>>} 密码库路径列表
   */
  async getRecentVaults() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [];
    }
    return await wailsAPI.GetRecentVaults();
  }

  /**
   * 检查最近使用的密码库状态
   * @returns {Promise<Object>} 包含存在的文件路径和简化模式状态
   * @author 陈凤庆
   * @description 20250928 检查最近使用的密码库文件是否存在，决定是否使用简化模式
   */
  async checkRecentVaultStatus() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return {
        hasValidVault: false,
        vaultPath: "",
        isSimplified: false,
      };
    }
    return await wailsAPI.CheckRecentVaultStatus();
  }

  /**
   * 获取所有分组
   * @returns {Promise<Array>} 分组列表
   * @modify 20251002 陈凤庆 后端返回字符串ID，避免JavaScript精度丢失
   */
  async getGroups() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [
        {
          id: "1",
          name: "默认",
          parent_id: "0",
          icon: "fa-folder-open",
          sort_order: 0,
        },
      ];
    }
    const groups = await wailsAPI.GetGroups();
    // 后端已返回字符串ID，直接使用
    return groups;
  }

  /**
   * 创建新分组
   * @param {string} name 分组名称
   * @returns {Promise<Object>} 创建的分组
   * @modify 20251002 陈凤庆 使用字符串ID，避免JavaScript精度丢失
   * @modify 20251002 陈凤庆 删除parentID参数，后端不需要层级结构
   */
  async createGroup(name) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return null;
    }
    const group = await wailsAPI.CreateGroup(name);
    // 后端已返回字符串ID，直接使用
    return group;
  }

  /**
   * 重命名分组
   * @param {string} id 分组ID
   * @param {string} newName 新的分组名称
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增分组重命名前端API
   */
  async renameGroup(id, newName) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.RenameGroup(id, newName);
  }

  /**
   * 删除分组
   * @param {string} id 分组ID
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增分组删除前端API
   */
  async deleteGroup(id) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.DeleteGroup(id);
  }

  /**
   * 将分组向左移动一位
   * @param {string} id 分组ID
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增分组左移前端API
   */
  async moveGroupLeft(id) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.MoveGroupLeft(id);
  }

  /**
   * 将分组向右移动一位
   * @param {string} id 分组ID
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增分组右移前端API
   */
  async moveGroupRight(id) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.MoveGroupRight(id);
  }

  /**
   * 根据分组ID获取页签列表
   * @param {string} groupID 分组ID
   * @returns {Promise<Array>} 页签列表
   * @modify 20251002 陈凤庆 使用字符串ID，避免JavaScript精度丢失
   * @modify 20251002 陈凤庆 GetTabsByGroup改名为GetTypesByGroup，适配后端方法名变更
   */
  async getTabsByGroup(groupID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [
        {
          id: "1",
          name: "网站账号",
          icon: "fa-globe",
          type: "default",
          group_id: groupID,
          sort_order: 0,
        },
      ];
    }
    // 20251002 陈凤庆 后端方法名已改为GetTypesByGroup
    const tabs = await wailsAPI.GetTypesByGroup(groupID);
    // 后端返回GUID字符串，直接使用
    return tabs;
  }

  /**
   * 获取所有页签
   * @returns {Promise<Array>} 页签列表
   */
  async getAllTabs() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [];
    }
    return await wailsAPI.GetAllTabs();
  }

  /**
   * 根据分组ID获取类型列表（别名方法）
   * @param {string} groupId 分组ID
   * @returns {Promise<Array>} 类型列表
   * @author 20251005 陈凤庆 为ChangeGroupDialog组件提供的别名方法
   */
  async getTypesByGroupId(groupId) {
    return this.getTabsByGroup(groupId);
  }

  /**
   * 更新账号分组
   * @param {string} accountId 账号ID
   * @param {string} typeId 新的类型ID
   * @returns {Promise<void>}
   * @author 20251005 陈凤庆 新增更新账号分组的API方法
   */
  async updateAccountGroup(accountId, typeId) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateAccountGroup(accountId, typeId);
  }

  /**
   * 获取所有类型/标签
   * @returns {Promise<Array>} 类型列表
   * @author 20251004 陈凤庆 新增获取所有类型的API，用于导出功能
   */
  async getTypes() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [];
    }
    return await wailsAPI.GetAllTypes();
  }

  /**
   * 创建新页签
   * @param {string} name 页签名称
   * @param {string} icon 图标
   * @param {number} groupID 所属分组ID
   * @returns {Promise<Object>} 创建的页签
   * @modify 20251016 陈凤庆 修复方法名，后端已改为CreateType
   */
  async createTab(name, icon, groupID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return null;
    }
    return await wailsAPI.CreateType(name, groupID, icon);
  }

  /**
   * 更新页签
   * @param {Object} tab 页签信息
   * @returns {Promise<void>}
   * @modify 20251016 陈凤庆 修复方法名，后端已改为UpdateType
   */
  async updateTab(tab) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateType(tab);
  }

  /**
   * 删除页签
   * @param {number} id 页签ID
   * @returns {Promise<void>}
   */
  async deleteTab(id) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.DeleteType(id);
  }

  /**
   * 标签上移
   * @param {string} id 标签ID
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增标签上移API调用
   */
  async moveTabUp(id) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.MoveTypeUp(id);
  }

  /**
   * 标签下移
   * @param {string} id 标签ID
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增标签下移API调用
   */
  async moveTabDown(id) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.MoveTypeDown(id);
  }

  /**
   * 更新标签排序
   * @param {string} typeID 标签ID
   * @param {number} newSortOrder 新的排序号
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增标签排序更新API调用，用于拖拽排序
   */
  async updateTabSortOrder(typeID, newSortOrder) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateTypeSortOrder(typeID, newSortOrder);
  }

  /**
   * 在指定标签后插入新标签
   * @param {string} name 新标签名称
   * @param {string} groupID 分组ID
   * @param {string} icon 图标
   * @param {string} afterTypeID 在此标签后插入，如果为空则插入到最后
   * @returns {Promise<Object>} 创建的标签
   * @author 20251002 陈凤庆 新增在指定标签后插入新标签API调用
   */
  async insertTabAfter(name, groupID, icon, afterTypeID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return null;
    }
    return await wailsAPI.InsertTypeAfter(name, groupID, icon, afterTypeID);
  }

  /**
   * 更新分组排序
   * @param {string} groupID 分组ID
   * @param {number} newSortOrder 新的排序号
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增分组排序更新API调用，用于拖拽排序
   */
  async updateGroupSortOrder(groupID, newSortOrder) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateGroupSortOrder(groupID, newSortOrder);
  }

  /**
   * 根据分组获取账号
   * @param {number} groupID 分组ID
   * @returns {Promise<Array>} 账号列表
   * @modify 20251001 陈凤庆 添加详细日志记录和错误处理
   */
  async getPasswordsByGroup(groupID) {
    return this.getAccountsByConditions({ group_id: groupID });
  }

  /**
   * 根据查询条件获取账号列表
   * @param {Object} conditions 查询条件对象，如 {group_id: "xxx", type_id: "xxx"}
   * @returns {Promise<Array>} 账号列表
   * @author 20251003 陈凤庆 新增统一的账号查询前端API
   */
  async getAccountsByConditions(conditions) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [];
    }

    try {
      const conditionsJson = JSON.stringify(conditions);
      console.log(
        "[API] 调用后端 GetAccountsByConditions，参数:",
        conditionsJson
      );
      const result = await wailsAPI.GetAccountsByConditions(conditionsJson);
      console.log("[API] 后端调用成功，返回结果:", result);

      // 验证返回数据格式
      if (!Array.isArray(result)) {
        console.error("[API] 后端返回数据不是数组格式:", result);
        throw new Error("Backend returned invalid data format, expected array");
      }

      return result;
    } catch (error) {
      console.error("[API] getAccountsByConditions 调用失败:", error);
      throw error;
    }
  }

  /**
   * 根据标签ID获取账号列表
   * @param {string} tabID 标签ID
   * @returns {Promise<Array>} 账号列表
   * @author 20251002 陈凤庆 新增按标签ID查询账号的前端API
   * @modify 20251003 陈凤庆 重构为调用统一的getAccountsByConditions方法
   */
  async getAccountsByTab(tabID) {
    return this.getAccountsByConditions({ type_id: tabID });
  }

  /**
   * 创建新账号
   * @param {Object} item 账号数据
   * @returns {Promise<Object>} 创建的账号
   * @modify 20251002 陈凤庆 修改API调用，使用CreateAccount接口，传递单独参数而非对象
   */
  async createPasswordItem(item) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return null;
    }

    // 20251019 陈凤庆 修复问题 004：增强API调用的调试日志
    console.log("[API] ========== 创建账号API调用 ==========");
    console.log("[API] 📥 原始传入参数:", item);

    // 准备API调用参数
    const apiParams = {
      title: item.title || "",
      username: item.username || "",
      password: item.password || "",
      url: item.url || "",
      typeID: item.type || item.typeid || "",
      notes: item.notes || "",
      inputMethod: item.input_method || 1,
    };

    console.log("[API] 📤 准备调用后端API，参数详情:");
    console.log(
      "  - title:",
      `"${apiParams.title}" (长度: ${apiParams.title.length})`
    );
    console.log(
      "  - username:",
      `"${apiParams.username}" (长度: ${apiParams.username.length})`
    );
    console.log(
      "  - password:",
      apiParams.password ? "***已设置***" : "未设置",
      `(长度: ${apiParams.password.length})`
    );
    console.log(
      "  - url:",
      `"${apiParams.url}" (长度: ${apiParams.url.length})`
    );
    console.log(
      "  - typeID:",
      `"${apiParams.typeID}" (长度: ${apiParams.typeID.length})`
    );
    console.log(
      "  - notes:",
      `"${apiParams.notes}" (长度: ${apiParams.notes.length})`
    );
    console.log("  - inputMethod:", apiParams.inputMethod);

    // 20251019 陈凤庆 新增：关键字段验证
    if (!apiParams.typeID) {
      console.error("[API] ❌ 关键字段typeID为空，这会导致后端验证失败");
      console.error("[API] 字段来源分析:");
      console.error("  - item.type:", item.type);
      console.error("  - item.typeid:", item.typeid);
      throw new Error("typeID字段为空，无法创建账号");
    }

    if (!apiParams.title) {
      console.error("[API] ❌ 关键字段title为空，这会导致后端验证失败");
      throw new Error("title字段为空，无法创建账号");
    }

    try {
      console.log("[API] 🚀 开始调用后端CreateAccount接口...");
      const result = await wailsAPI.CreateAccount(
        apiParams.title,
        apiParams.username,
        apiParams.password,
        apiParams.url,
        apiParams.typeID,
        apiParams.notes,
        apiParams.inputMethod
      );

      console.log("[API] ✅ 后端API调用成功，返回结果:", result);
      console.log("[API] ========== 创建账号API调用完成 ==========");
      return result;
    } catch (error) {
      console.error("[API] ❌ 后端API调用失败:", error);
      console.error("[API] 失败时的参数:", apiParams);
      console.error("[API] ========== 创建账号API调用失败 ==========");
      throw error;
    }
  }

  /**
   * 创建新账号（新方法名）
   * @param {Object} item 账号数据
   * @returns {Promise<Object>} 创建的账号
   */
  async createAccount(item) {
    return this.createPasswordItem(item);
  }

  /**
   * 更新账号项
   * @param {Object} item 账号项数据
   * @returns {Promise<void>}
   * @author 20251002 陈凤庆 新增更新账号项的前端API
   */
  async updatePasswordItem(item) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }

    try {
      console.log("[API] 调用后端 UpdateAccount，参数:", item);

      // 20251002 陈凤庆 调用UpdateAccount接口，传递完整的账号对象
      await wailsAPI.UpdateAccount(item);

      console.log("[API] 后端调用成功，账号更新完成");
    } catch (error) {
      console.error("[API] updatePasswordItem 调用失败:", error);
      console.error("[API] 错误类型:", error.name);
      console.error("[API] 错误消息:", error.message);
      console.error("[API] 错误堆栈:", error.stack);

      // 重新抛出错误，保持原始错误信息
      throw error;
    }
  }

  /**
   * 删除账号
   * @param {string} accountId 账号ID
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增删除账号的前端API
   */
  async deleteAccount(accountId) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }

    try {
      console.log("[API] 调用后端 DeleteAccount，参数:", accountId);
      await wailsAPI.DeleteAccount(accountId);
      console.log("[API] 后端调用成功，账号删除完成");
    } catch (error) {
      console.error("[API] deleteAccount 调用失败:", error);
      throw error;
    }
  }

  /**
   * 搜索账号
   * @param {string} keyword 搜索关键词
   * @returns {Promise<Array>} 搜索结果
   */
  async searchPasswords(keyword) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [];
    }
    return await wailsAPI.SearchAccounts(keyword);
  }

  /**
   * 更新账号使用次数
   * @param {string} accountId 账号ID
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增使用次数统计功能
   */
  async updateAccountUsage(accountId) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateAccountUsage(accountId);
  }

  /**
   * 获取所有账号
   * @returns {Promise<Array>} 所有账号列表
   * @author 20251003 陈凤庆 新增使用次数统计功能
   */
  async getAllAccounts() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [];
    }
    return await wailsAPI.GetAllAccounts();
  }

  /**
   * 更新应用使用统计
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增使用天数统计功能
   */
  async updateAppUsage() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateAppUsage();
  }

  /**
   * 根据账号ID获取账号详情（用于详情页面显示）
   * @param {string} accountID 账号ID
   * @returns {Promise<Object>} 账号详情
   * @author 20251003 陈凤庆 新增账号详情查询接口
   */
  async getAccountDetail(accountID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return null;
    }

    try {
      console.log("[API] 调用后端 GetAccountDetail，参数:", accountID);
      const result = await wailsAPI.GetAccountDetail(accountID);
      console.log("[API] GetAccountDetail 调用成功，返回结果:", result);
      return result;
    } catch (error) {
      console.error("[API] getAccountDetail 调用失败:", error);
      throw error;
    }
  }

  /**
   * 根据账号ID获取解密后的密码
   * @param {string} accountID 账号ID
   * @returns {Promise<string>} 解密后的密码
   * @author 20251003 陈凤庆 新增密码查询接口
   */
  async getAccountPassword(accountID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return "";
    }

    try {
      console.log("[API] 调用后端 GetAccountPassword，参数:", accountID);
      const result = await wailsAPI.GetAccountPassword(accountID);
      console.log("[API] GetAccountPassword 调用成功");
      return result;
    } catch (error) {
      console.error("[API] getAccountPassword 调用失败:", error);
      throw error;
    }
  }

  /**
   * 根据账号ID获取解密后的备注
   * @param {string} accountID 账号ID
   * @returns {Promise<string>} 解密后的备注
   * @author 20251003 陈凤庆 新增备注查询接口
   */
  async getAccountNotes(accountID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return "";
    }

    try {
      console.log("[API] 调用后端 GetAccountNotes，参数:", accountID);
      const result = await wailsAPI.GetAccountNotes(accountID);
      console.log("[API] GetAccountNotes 调用成功");
      return result;
    } catch (error) {
      console.error("[API] getAccountNotes 调用失败:", error);
      throw error;
    }
  }

  /**
   * 复制账号备注到剪贴板
   * @param {string} accountID 账号ID
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增备注复制接口
   */
  async copyAccountNotes(accountID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API 未初始化");
    }

    try {
      console.log("[API] 调用后端 CopyAccountNotes，参数:", accountID);
      await wailsAPI.CopyAccountNotes(accountID);
      console.log("[API] CopyAccountNotes 调用成功");
    } catch (error) {
      console.error("[API] copyAccountNotes 调用失败:", error);
      throw error;
    }
  }

  /**
   * 复制账号密码到剪贴板（标准函数，可复用）
   * @param {string} accountID 账号ID
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增标准复制密码函数
   */
  async copyAccountPassword(accountID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }

    try {
      console.log("[API] 调用后端 CopyAccountPassword，参数:", accountID);
      await wailsAPI.CopyAccountPassword(accountID);
      console.log("[API] CopyAccountPassword 调用成功");
    } catch (error) {
      console.error("[API] copyAccountPassword 调用失败:", error);
      throw error;
    }
  }

  /**
   * 根据账号ID获取完整账号数据（用于编辑）
   * @param {string} accountID 账号ID
   * @returns {Promise<Object>} 完整账号数据
   * @author 20251003 陈凤庆 新增获取完整账号数据接口
   */
  async getAccountByID(accountID) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return null;
    }

    try {
      console.log("[API] 调用后端 GetAccountByID，参数:", accountID);
      const result = await wailsAPI.GetAccountByID(accountID);
      console.log("[API] GetAccountByID 调用成功，返回结果:", result);
      return result;
    } catch (error) {
      console.error("[API] getAccountByID 调用失败:", error);
      throw error;
    }
  }

  /**
   * 获取使用天数
   * @returns {Promise<number>} 使用天数
   * @author 20251003 陈凤庆 新增使用天数统计功能
   */
  async getUsageDays() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return 30; // 返回模拟数据
    }
    return await wailsAPI.GetUsageDays();
  }

  /**
   * 获取应用信息
   * @returns {Promise<Object>} 应用信息
   */
  async getAppInfo() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return {
        name: "密码管理器",
        version: "1.0.0",
        author: "陈凤庆",
        buildDate: "2025-10-01",
        time: new Date().toLocaleString("zh-CN"),
      };
    }
    return await wailsAPI.GetAppInfo();
  }

  /**
   * 获取更新日志
   * @returns {Promise<Array>} 更新日志列表
   */
  async getChangeLog() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回模拟数据");
      return [
        {
          Version: "1.0.0",
          Date: "2025-10-01",
          Changes: [
            "[新增] 密码库创建界面优化",
            "[新增] 自动填充功能",
            "[优化] 界面布局和交互体验",
          ],
        },
      ];
    }
    return await wailsAPI.GetChangeLog();
  }

  /**
   * 获取日志配置
   * @returns {Promise<Object>} 日志配置
   * @author 陈凤庆
   * @date 2025-10-03
   */
  async getLogConfig() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回默认配置");
      return {
        enable_info_log: true,
        enable_debug_log: true,
      };
    }
    return await wailsAPI.GetLogConfig();
  }

  /**
   * 设置日志配置
   * @param {Object} config 日志配置
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @date 2025-10-03
   */
  async setLogConfig(config) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.SetLogConfig(config);
  }

  /**
   * 获取锁定配置
   * @returns {Promise<Object>} 锁定配置
   * @author 陈凤庆
   * @date 2025-10-04
   */
  async getLockConfig() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回默认锁定配置");
      return {
        enable_auto_lock: false,
        enable_timer_lock: false,
        enable_minimize_lock: false,
        lock_time_minutes: 10,
        enable_system_lock: true,
        system_lock_minutes: 120,
      };
    }
    return await wailsAPI.GetLockConfig();
  }

  /**
   * 设置锁定配置
   * @param {Object} config 锁定配置
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @date 2025-10-04
   */
  async setLockConfig(config) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.SetLockConfig(config);
  }

  /**
   * 获取快捷键配置
   * @returns {Promise<Object>} 快捷键配置
   * @author 陈凤庆
   * @date 2025-10-14
   */
  async getHotkeyConfig() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回默认快捷键配置");
      return {
        enable_global_hotkey: true,
        show_hide_hotkey: "Ctrl+Alt+H",
      };
    }
    return await wailsAPI.GetHotkeyConfig();
  }

  /**
   * 设置快捷键配置
   * @param {Object} config 快捷键配置
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @date 2025-10-14
   */
  async setHotkeyConfig(config) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.SetHotkeyConfig(config);
  }

  /**
   * 手动触发锁定
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @date 2025-10-04
   */
  async triggerLock() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.TriggerLock();
  }

  /**
   * 更新用户活动时间
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @date 2025-10-04
   */
  async updateUserActivity() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.UpdateUserActivity();
  }

  /**
   * 获取应用配置
   * @returns {Promise<Object>} 应用配置
   * @author 陈凤庆
   * @date 2025-10-03
   */
  async getAppConfig() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化，返回默认配置");
      return {
        theme: "light",
        language: "zh-CN",
      };
    }
    return await wailsAPI.GetAppConfig();
  }

  /**
   * 设置应用配置
   * @param {Object} config 应用配置
   * @returns {Promise<void>}
   * @author 陈凤庆
   * @date 2025-10-03
   */
  async setAppConfig(config) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return;
    }
    return await wailsAPI.SetAppConfig(config);
  }

  /**
   * 验证旧登录密码
   * @param {string} oldPassword 旧登录密码
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增验证旧登录密码的前端API
   */
  async verifyOldPassword(oldPassword) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API 未初始化");
    }

    try {
      console.log("[API] 调用后端 VerifyOldPassword");
      await wailsAPI.VerifyOldPassword(oldPassword);
      console.log("[API] VerifyOldPassword 调用成功");
    } catch (error) {
      console.error("[API] verifyOldPassword 调用失败:", error);
      throw error;
    }
  }

  /**
   * 修改登录密码
   * @param {string} oldPassword 旧登录密码
   * @param {string} newPassword 新登录密码
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增修改登录密码的前端API
   */
  async changeLoginPassword(oldPassword, newPassword) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API 未初始化");
    }

    try {
      console.log("[API] 调用后端 ChangeLoginPassword");
      await wailsAPI.ChangeLoginPassword(oldPassword, newPassword);
      console.log("[API] ChangeLoginPassword 调用成功");
    } catch (error) {
      console.error("[API] changeLoginPassword 调用失败:", error);
      throw error;
    }
  }

  /**
   * 验证旧登录密码
   * @param {string} oldPassword 旧登录密码
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增验证旧登录密码的前端API
   */
  async verifyOldPassword(oldPassword) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API 未初始化");
    }

    try {
      console.log("[API] 调用后端 VerifyOldPassword");
      await wailsAPI.VerifyOldPassword(oldPassword);
      console.log("[API] VerifyOldPassword 调用成功");
    } catch (error) {
      console.error("[API] verifyOldPassword 调用失败:", error);
      throw error;
    }
  }

  /**
   * 选择导出路径
   * @returns {Promise<string>} 选择的导出路径
   * @author 20251003 陈凤庆 新增选择导出路径的前端API
   */
  async selectExportPath() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return "";
    }

    try {
      console.log("[API] 调用后端 SelectExportPath");
      const result = await wailsAPI.SelectExportPath();
      console.log("[API] SelectExportPath 调用成功");
      return result;
    } catch (error) {
      console.error("[API] selectExportPath 调用失败:", error);
      throw error;
    }
  }

  /**
   * 选择导入文件
   * @returns {Promise<string>} 选择的导入文件路径
   * @author 20251003 陈凤庆 新增选择导入文件的前端API
   */
  async selectImportFile() {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      return "";
    }

    try {
      console.log("[API] 调用后端 SelectImportFile");
      const result = await wailsAPI.SelectImportFile();
      console.log("[API] SelectImportFile 调用成功");
      return result;
    } catch (error) {
      console.error("[API] selectImportFile 调用失败:", error);
      throw error;
    }
  }

  /**
   * 导出密码库
   * @param {string} loginPassword 登录密码
   * @param {string} backupPassword 备份密码
   * @param {string} exportPath 导出路径
   * @param {string[]} accountIDs 要导出的账号ID列表（手动选择模式）
   * @param {string[]} groupIDs 要导出的分组ID列表（按分组导出模式）
   * @param {string[]} typeIDs 要导出的类别ID列表（按类别导出模式）
   * @param {boolean} exportAll 是否导出所有账号
   * @returns {Promise<void>}
   * @author 20251003 陈凤庆 新增导出密码库的前端API
   * @modify 20251003 陈凤庆 支持按分组和按类别导出
   */
  async exportVault(
    loginPassword,
    backupPassword,
    exportPath,
    accountIDs,
    groupIDs,
    typeIDs,
    exportAll
  ) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API 未初始化");
    }

    try {
      console.log("[API] 调用后端 ExportVault");
      console.log("[API] 导出参数:", {
        accountIDs: accountIDs || [],
        groupIDs: groupIDs || [],
        typeIDs: typeIDs || [],
        exportAll: exportAll || false,
      });

      await wailsAPI.ExportVault(
        loginPassword,
        backupPassword,
        exportPath,
        accountIDs || [],
        groupIDs || [],
        typeIDs || [],
        exportAll || false
      );
      console.log("[API] ExportVault 调用成功");
    } catch (error) {
      console.error("[API] exportVault 调用失败:", error);
      throw error;
    }
  }

  /**
   * 导入密码库
   * @param {string} importPath 导入文件路径
   * @param {string} backupPassword 备份密码（解压密码）
   * @returns {Promise<Object>} 导入结果
   * @author 20251003 陈凤庆 新增导入密码库的前端API
   */
  async importVault(importPath, backupPassword) {
    const wailsAPI = this.getWailsAPI();
    if (!wailsAPI) {
      console.warn("Wails API 未初始化");
      throw new Error("Wails API 未初始化");
    }

    try {
      console.log("[API] 调用后端 ImportVault");
      const result = await wailsAPI.ImportVault(importPath, backupPassword);
      console.log("[API] ImportVault 调用成功");
      return result;
    } catch (error) {
      console.error("[API] importVault 调用失败:", error);
      throw error;
    }
  }
}

// 导出单例实例
export const apiService = new ApiService();
