/**
 * 锁定事件服务
 * @author 陈凤庆
 * @date 2025-10-04
 * @description 处理自动锁定事件和状态监听
 */

import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";
import { useVaultStore } from "@/stores/vault";

// 安全的i18n访问函数
function safeT(key, fallback = key) {
  try {
    // 尝试从全局window对象获取i18n实例
    if (window.i18n && window.i18n.global && window.i18n.global.t) {
      return window.i18n.global.t(key);
    }
    // 如果没有i18n实例，返回fallback
    return fallback;
  } catch (error) {
    console.warn(`i18n translation failed for key: ${key}`, error);
    return fallback;
  }
}

class LockEventService {
  constructor() {
    this.isMonitoring = false;
    this.checkInterval = null;
    this.router = null;
    this.vaultStore = null;
    this.lastCheckTime = Date.now();
    this.minimizeTime = null; // 记录窗口最小化时间
    this.lockConfig = null; // 20251004 陈凤庆 存储锁定配置
  }

  /**
   * 初始化锁定事件服务
   * @param {Object} router Vue Router实例
   * @param {Object} vaultStore Vault Store实例
   */
  init(router, vaultStore) {
    this.router = router;
    this.vaultStore = vaultStore;
    console.log("[锁定事件服务] 初始化完成");
  }

  /**
   * 设置锁定配置
   * @param {Object} lockConfig 锁定配置对象
   */
  setLockConfig(lockConfig) {
    this.lockConfig = lockConfig;
    console.log("[锁定事件服务] 锁定配置已设置:", lockConfig);
  }

  /**
   * 开始监听锁定事件
   */
  startMonitoring() {
    if (this.isMonitoring) {
      console.log("[锁定事件服务] 已在监听中，跳过重复启动");
      return;
    }

    this.isMonitoring = true;
    console.log("[锁定事件服务] 开始监听锁定事件");

    // 每30秒检查一次锁定状态（减少频率避免过于敏感）
    this.checkInterval = setInterval(() => {
      this.checkLockStatus();
    }, 30000);

    // 监听窗口事件
    this.setupWindowEventListeners();
  }

  /**
   * 停止监听锁定事件
   */
  stopMonitoring() {
    if (!this.isMonitoring) {
      return;
    }

    this.isMonitoring = false;
    console.log("[锁定事件服务] 停止监听锁定事件");

    if (this.checkInterval) {
      clearInterval(this.checkInterval);
      this.checkInterval = null;
    }

    // 移除窗口事件监听器
    this.removeWindowEventListeners();
  }

  /**
   * 检查锁定状态（简化版本）
   */
  async checkLockStatus() {
    try {
      // 20251004 陈凤庆 检查后端是否触发了锁定
      const isLockTriggered = await window.go?.app?.App?.IsLockTriggered?.();

      console.log("[锁定事件服务] 检查锁定触发状态:", isLockTriggered);

      if (isLockTriggered) {
        console.log("[锁定事件服务] 检测到后端触发锁定，准备跳转到登录页面");
        await this.handleLockEvent();
      } else {
        console.log("[锁定事件服务] 锁定状态正常，继续监听");
      }
    } catch (error) {
      // 如果API调用失败，记录错误但不触发锁定
      console.log("[锁定事件服务] API调用失败:", error.message);
      console.error("[锁定事件服务] 详细错误信息:", error);
    }
  }

  /**
   * 处理锁定事件
   */
  async handleLockEvent() {
    if (!this.isMonitoring) {
      return;
    }

    console.log("[锁定事件服务] 处理锁定事件");

    // 停止监听
    this.stopMonitoring();

    // 更新前端状态
    if (this.vaultStore) {
      this.vaultStore.setVaultOpened(false, "");
      console.log(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.frontendStateUpdated", "前端状态已更新")}`
      );
    }

    // 清理敏感数据
    this.clearSensitiveData();

    // 跳转到登录页面
    if (this.router && this.router.currentRoute.value.path !== "/login") {
      console.log(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.redirectToLogin", "跳转到登录页面")}`
      );
      await this.router.push("/login");

      // 显示锁定消息
      ElMessage.info({
        message: safeT(
          "lockEventService.vaultAutoLocked",
          "密码库已自动锁定，请重新登录"
        ),
        duration: 3000,
        showClose: true,
      });
    }
  }

  /**
   * 清理敏感数据
   */
  clearSensitiveData() {
    try {
      // 清理localStorage中的敏感数据
      const sensitiveKeys = [
        "vault_password",
        "current_accounts",
        "decrypted_data",
      ];

      sensitiveKeys.forEach((key) => {
        if (localStorage.getItem(key)) {
          localStorage.removeItem(key);
        }
      });

      // 20251017 陈凤庆 清理sessionStorage中的敏感数据
      // 注意：登录密码现在存储在后端内存中，会在密码库关闭时自动清理
      sessionStorage.clear();

      console.log(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.sensitiveDataCleared", "敏感数据已清理")}`
      );
    } catch (error) {
      console.error(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.clearSensitiveDataFailed", "清理敏感数据失败")}`,
        error
      );
    }
  }

  /**
   * 设置窗口事件监听器
   */
  setupWindowEventListeners() {
    // 监听窗口最小化事件
    document.addEventListener(
      "visibilitychange",
      this.handleVisibilityChange.bind(this)
    );

    // 监听窗口焦点事件
    window.addEventListener("blur", this.handleWindowBlur.bind(this));
    window.addEventListener("focus", this.handleWindowFocus.bind(this));

    console.log(
      `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.windowEventListenersSet", "窗口事件监听器已设置")}`
    );
  }

  /**
   * 移除窗口事件监听器
   */
  removeWindowEventListeners() {
    document.removeEventListener(
      "visibilitychange",
      this.handleVisibilityChange.bind(this)
    );
    window.removeEventListener("blur", this.handleWindowBlur.bind(this));
    window.removeEventListener("focus", this.handleWindowFocus.bind(this));

    console.log("[锁定事件服务] 窗口事件监听器已移除");
  }

  /**
   * 处理页面可见性变化
   */
  handleVisibilityChange() {
    if (document.hidden) {
      console.log("[锁定事件服务] 页面隐藏 - 可能是窗口最小化");
      // 记录最小化时间，用于后续检查
      this.minimizeTime = Date.now();
      // 通知后端窗口最小化
      this.notifyWindowMinimize();

      // 20251004 陈凤庆 根据需求描述，最小化锁定应该立即返回登录状态
      // 检查是否启用了最小化锁定，如果是则立即触发锁定检查
      if (
        this.lockConfig &&
        this.lockConfig.enable_auto_lock &&
        this.lockConfig.enable_minimize_lock
      ) {
        console.log(
          "[锁定事件服务] 检测到最小化且启用了最小化锁定，立即检查锁定状态"
        );
        // 立即检查锁定状态，不等待窗口恢复
        setTimeout(() => {
          this.checkLockStatus();
        }, 100); // 短暂延迟确保后端处理完成
      }
    } else {
      console.log("[锁定事件服务] 页面显示 - 窗口恢复");
      // 检查是否从最小化状态恢复
      if (this.minimizeTime) {
        const hiddenDuration = Date.now() - this.minimizeTime;
        console.log(`[锁定事件服务] 窗口隐藏时长: ${hiddenDuration}ms`);
        this.minimizeTime = null;

        // 20251005 陈凤庆 只在启用了最小化锁定时才检查状态，避免过于频繁的检查
        if (
          this.lockConfig &&
          this.lockConfig.enable_auto_lock &&
          this.lockConfig.enable_minimize_lock
        ) {
          setTimeout(() => {
            console.log(
              "[锁定事件服务] 窗口恢复后检查锁定状态（最小化锁定已启用）"
            );
            this.checkLockStatus();
          }, 200); // 适中延迟
        } else {
          console.log(
            "[锁定事件服务] 窗口恢复，但未启用最小化锁定，跳过锁定检查"
          );
        }
      }

      // 通知后端窗口获得焦点
      this.notifyWindowFocus();
    }
  }

  /**
   * 处理窗口失去焦点
   */
  handleWindowBlur() {
    console.log(
      `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.windowLostFocus", "窗口失去焦点")}`
    );
    // 可以在这里添加额外的处理逻辑
  }

  /**
   * 处理窗口获得焦点
   */
  handleWindowFocus() {
    console.log(
      `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.windowGainedFocus", "窗口获得焦点")}`
    );
    // 20251005 陈凤庆 移除窗口焦点时的锁定检查，避免过于敏感
    this.notifyWindowFocus();
  }

  /**
   * 通知后端窗口最小化
   */
  async notifyWindowMinimize() {
    try {
      await window.go?.app?.App?.OnWindowMinimize?.();
      console.log(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.backendNotifiedMinimize", "后端已通知窗口最小化")}`
      );
    } catch (error) {
      console.error(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.notifyMinimizeFailed", "通知最小化失败")}`,
        error
      );
    }
  }

  /**
   * 通知后端窗口获得焦点
   */
  async notifyWindowFocus() {
    try {
      await window.go?.app?.App?.OnWindowFocus?.();
      console.log(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.backendNotifiedFocus", "后端已通知窗口获得焦点")}`
      );
    } catch (error) {
      console.error(
        `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.notifyFocusFailed", "通知焦点失败")}`,
        error
      );
    }
  }

  /**
   * 手动触发锁定检查
   */
  async triggerLockCheck() {
    console.log(
      `[${safeT("lockEventService.logPrefix", "锁定事件服务")}] ${safeT("lockEventService.manualLockCheckTriggered", "手动锁定检查已触发")}`
    );
    await this.checkLockStatus();
  }
}

// 创建单例实例
const lockEventService = new LockEventService();

export default lockEventService;
