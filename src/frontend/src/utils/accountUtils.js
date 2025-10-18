/**
 * 账号相关工具函数
 * @author 陈凤庆
 * @date 20251003
 * @description 提供账号操作的通用工具函数
 */

import { ElMessage } from "element-plus";
import { apiService } from "../services/api.js";
import i18n from "@/i18n";

/**
 * 复制账号密码到剪贴板（标准函数，可复用）
 * @param {string} accountId 账号ID
 * @param {string} accountTitle 账号标题（用于提示信息）
 * @returns {Promise<void>}
 * @author 20251003 陈凤庆 创建标准复制密码函数
 */
export async function copyAccountPassword(
  accountId,
  accountTitle = i18n.global.t("accountUtils.defaultAccountTitle")
) {
  try {
    console.log(
      `[accountUtils.copyAccountPassword] ${i18n.global.t("accountUtils.startCopyingPassword", { accountId: accountId })}`
    );

    // 调用后端安全复制方法，10秒后自动清理剪贴板
    await apiService.copyAccountPassword(accountId);

    ElMessage.success(i18n.global.t("accountUtils.passwordCopiedSuccess"));
    console.log(
      `[accountUtils.copyAccountPassword] ${i18n.global.t("accountUtils.passwordCopySuccess", { accountId: accountId })}`
    );
  } catch (error) {
    console.error(
      `[accountUtils.copyAccountPassword] ${i18n.global.t("accountUtils.passwordCopyFailed")}:`,
      error
    );
    ElMessage.error(
      `${i18n.global.t("accountUtils.passwordCopyFailed")}: ${error.message || i18n.global.t("common.unknownError")}`
    );
    throw error;
  }
}

/**
 * 获取账号密码并显示到界面
 * @param {string} accountId 账号ID
 * @returns {Promise<string>} 解密后的密码
 * @author 20251003 陈凤庆 创建获取密码函数
 */
export async function getAccountPassword(accountId) {
  try {
    console.log(
      `[accountUtils.getAccountPassword] ${i18n.global.t("accountUtils.startGettingPassword", { accountId: accountId })}`
    );

    // 调用后端接口获取解密后的密码
    const password = await apiService.getAccountPassword(accountId);

    console.log(
      `[accountUtils.getAccountPassword] ${i18n.global.t("accountUtils.passwordGetSuccess", { accountId: accountId })}`
    );
    return password;
  } catch (error) {
    console.error(
      `[accountUtils.getAccountPassword] ${i18n.global.t("accountUtils.passwordGetFailed")}:`,
      error
    );
    ElMessage.error(
      `${i18n.global.t("accountUtils.passwordGetFailed")}: ${error.message || i18n.global.t("common.unknownError")}`
    );
    throw error;
  }
}

/**
 * 获取账号详情
 * @param {string} accountId 账号ID
 * @returns {Promise<Object>} 账号详情
 * @author 20251003 陈凤庆 创建获取账号详情函数
 */
export async function getAccountDetail(accountId) {
  try {
    console.log(
      `[accountUtils.getAccountDetail] ${i18n.global.t("accountUtils.startGettingAccountDetail", { accountId: accountId })}`
    );

    // 调用后端接口获取账号详情
    const accountDetail = await apiService.getAccountDetail(accountId);

    console.log(
      `[accountUtils.getAccountDetail] ${i18n.global.t("accountUtils.accountDetailGetSuccess", { accountId: accountId })}`
    );
    return accountDetail;
  } catch (error) {
    console.error(
      `[accountUtils.getAccountDetail] ${i18n.global.t("accountUtils.accountDetailGetFailed")}:`,
      error
    );
    ElMessage.error(
      `${i18n.global.t("accountUtils.accountDetailGetFailed")}: ${error.message || i18n.global.t("common.unknownError")}`
    );
    throw error;
  }
}
