//go:build windows
// +build windows

package services

/*
#include <windows.h>
*/
import "C"

/**
 * Windows平台辅助函数
 * @author 陈凤庆
 * @description 20251003 Windows平台通用的辅助函数
 */

/**
 * cBoolToBool C语言BOOL转Go语言bool
 * @param cBool C语言BOOL值
 * @return bool Go语言bool值
 * @description 20251003 陈凤庆 辅助函数，用于C和Go之间的布尔值转换
 */
func cBoolToBool(cBool C.BOOL) bool {
	return cBool != 0
}

