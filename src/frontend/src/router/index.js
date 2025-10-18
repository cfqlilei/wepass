import {
  createRouter,
  createWebHistory,
  createWebHashHistory,
} from "vue-router";
import { useVaultStore } from "@/stores/vault";
import { apiService } from "@/services/api";

/**
 * 路由配置
 * @author 陈凤庆
 * @description 定义应用的路由规则和导航守卫
 * @modify 20251003 陈凤庆 根据平台自动选择路由模式，解决跨平台兼容性问题
 */

/**
 * 检测当前运行平台
 * @returns {string} 平台名称：'windows', 'macos', 'linux', 'unknown'
 */
function detectPlatform() {
  // 20251003 陈凤庆 增强平台检测逻辑，解决Windows平台路由问题
  if (typeof window !== "undefined") {
    const userAgent = window.navigator.userAgent.toLowerCase();
    const platform = window.navigator.platform.toLowerCase();

    console.log(`[平台检测] UserAgent: ${userAgent}`);
    console.log(`[平台检测] Platform: ${platform}`);

    // 多种方式检测Windows平台
    if (
      userAgent.includes("win") ||
      platform.includes("win") ||
      userAgent.includes("windows") ||
      platform.includes("windows")
    ) {
      console.log(`[平台检测] 确认为Windows平台`);
      return "windows";
    }

    if (userAgent.includes("mac") || platform.includes("mac")) {
      console.log(`[平台检测] 确认为macOS平台`);
      return "macos";
    }

    if (userAgent.includes("linux") || platform.includes("linux")) {
      console.log(`[平台检测] 确认为Linux平台`);
      return "linux";
    }
  }

  // 20251003 陈凤庆 如果检测失败，默认使用Hash模式（更安全）
  console.log(`[平台检测] 无法确定平台，默认使用Hash模式`);
  return "unknown";
}

/**
 * 根据平台选择路由模式
 * @returns {History} 路由历史对象
 */
function createRouterHistory() {
  const platform = detectPlatform();
  console.log(`[路由] 检测到平台: ${platform}`);

  // 20251003 陈凤庆 Windows和未知平台都使用Hash模式，确保兼容性
  if (platform === "windows" || platform === "unknown") {
    console.log(`[路由] ${platform}平台，使用Hash模式避免404问题`);
    return createWebHashHistory();
  } else {
    // 仅macOS和Linux平台使用History模式
    console.log(`[路由] ${platform}平台，使用History模式`);
    return createWebHistory();
  }
}

const routes = [
  {
    path: "/",
    name: "Home",
    redirect: "/login",
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/LoginView.vue"),
    meta: {
      title: "login.title",
      requiresAuth: false,
    },
  },
  {
    path: "/main",
    name: "Main",
    component: () => import("@/views/MainView.vue"),
    meta: {
      title: "main.title",
      requiresAuth: true,
    },
  },
  {
    path: "/test",
    name: "Test",
    component: () => import("@/test/PasswordDetailTest.vue"),
    meta: {
      title: "test.passwordDetailTest",
      requiresAuth: false,
    },
  },
];

// 20251003 陈凤庆 根据平台自动选择路由模式
const router = createRouter({
  history: createRouterHistory(),
  routes,
});

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const platform = detectPlatform();

  // 20251003 陈凤庆 详细的路由和URL检查，帮助调试Windows平台问题
  console.log(
    `[路由守卫] ==================== 路由守卫开始 ====================`
  );
  console.log(`[路由守卫] 平台: ${platform}`);
  console.log(`[路由守卫] 当前URL: ${window.location.href}`);
  console.log(
    `[路由守卫] 路由模式: ${router.options.history.base ? "History" : "Hash"}`
  );
  console.log(`[路由守卫] 从 ${from.path} 跳转到 ${to.path}`);
  console.log(`[路由守卫] to.fullPath: ${to.fullPath}`);
  console.log(`[路由守卫] from.fullPath: ${from.fullPath}`);

  const vaultStore = useVaultStore();
  console.log(`[路由守卫] 当前 isVaultOpened = ${vaultStore.isVaultOpened}`);

  // 动态设置页面标题
  try {
    const appInfo = await apiService.getAppInfo();
    const appTitle = `${appInfo.name} v${appInfo.version}`;

    if (to.path === "/login") {
      document.title = appTitle;
    } else if (to.path === "/main") {
      document.title = appTitle;
    } else if (to.meta.title) {
      document.title = to.meta.title;
    } else {
      document.title = appTitle;
    }
  } catch (error) {
    console.error("获取应用信息失败，使用默认标题:", error);
    if (to.meta.title) {
      document.title = to.meta.title;
    } else {
      document.title = "wepass";
    }
  }

  // 20251004 陈凤庆 优先检查锁定触发状态（简化锁定功能）
  console.log(`[路由守卫] 🔒 检查锁定触发状态...`);
  let isLockTriggered = false;
  try {
    if (
      window.go &&
      window.go.app &&
      window.go.app.App &&
      window.go.app.App.IsLockTriggered
    ) {
      isLockTriggered = await window.go.app.App.IsLockTriggered();
      console.log(`[路由守卫] 锁定触发状态: ${isLockTriggered}`);
    }
  } catch (error) {
    console.log(`[路由守卫] 检查锁定触发状态失败:`, error);
  }

  // 如果检测到锁定触发，直接处理锁定逻辑
  if (isLockTriggered) {
    console.log(`[路由守卫] 🔒 检测到锁定触发，强制跳转到登录页面`);
    vaultStore.closeVault(); // 更新前端状态
    if (to.path !== "/login") {
      next("/login");
      return;
    } else {
      next();
      return;
    }
  }

  // 20251003 陈凤庆 检查后端的密码库打开状态，而不是前端store
  // 这样可以避免前端状态丢失导致的重复登录问题
  console.log(`[路由守卫] 🔍 开始检查后端密码库状态...`);
  console.log(`[路由守卫] 路由跳转: ${from.path} → ${to.path}`);
  console.log(
    `[路由守卫] 前端store当前状态: isVaultOpened = ${vaultStore.isVaultOpened}`
  );

  let isVaultOpened = false;
  let stateCheckAttempts = 0;
  const maxAttempts = 3;

  // 20251003 陈凤庆 添加重试机制，解决可能的时序问题
  while (stateCheckAttempts < maxAttempts) {
    stateCheckAttempts++;
    console.log(`[路由守卫] 第${stateCheckAttempts}次状态检查...`);

    try {
      isVaultOpened = await apiService.isVaultOpened();
      console.log(
        `[路由守卫] 🔍 后端状态检查结果: isVaultOpened = ${isVaultOpened}`
      );

      // 如果状态明确，跳出重试循环
      break;
    } catch (error) {
      console.error(
        `[路由守卫] ❌ 第${stateCheckAttempts}次状态检查失败:`,
        error
      );

      if (stateCheckAttempts < maxAttempts) {
        console.log(`[路由守卫] 等待100ms后重试...`);
        await new Promise((resolve) => setTimeout(resolve, 100));
      } else {
        console.error("[路由守卫] ❌ 所有状态检查尝试都失败，假设未登录");
        isVaultOpened = false;
      }
    }
  }

  // 20251003 陈凤庆 详细的状态同步逻辑
  console.log(`[路由守卫] 状态同步检查:`);
  console.log(`[路由守卫] - 后端状态: ${isVaultOpened}`);
  console.log(`[路由守卫] - 前端状态: ${vaultStore.isVaultOpened}`);

  if (isVaultOpened && !vaultStore.isVaultOpened) {
    console.log(
      `[路由守卫] 🔄 状态不同步，更新前端状态: 设置 isVaultOpened = true`
    );
    vaultStore.setVaultOpened(true, "");
    console.log(`[路由守卫] ✅ 前端状态已同步为: ${vaultStore.isVaultOpened}`);
  } else if (!isVaultOpened && vaultStore.isVaultOpened) {
    console.log(
      `[路由守卫] 🔄 状态不同步，更新前端状态: 设置 isVaultOpened = false`
    );
    vaultStore.closeVault();
    console.log(`[路由守卫] ✅ 前端状态已同步为: ${vaultStore.isVaultOpened}`);
  } else {
    console.log(`[路由守卫] ✅ 前后端状态一致，无需同步`);
  }

  // 检查是否需要认证
  if (to.meta.requiresAuth && !isVaultOpened) {
    console.log(`[路由守卫] ${platform}: 需要认证但未登录，跳转到登录页`);
    next("/login");
  } else if (to.path === "/login" && isVaultOpened) {
    console.log(`[路由守卫] ${platform}: 已登录但访问登录页，跳转到主页`);
    next("/main");
  } else {
    console.log(`[路由守卫] ${platform}: 允许跳转`);
    next();
  }
});

export default router;
