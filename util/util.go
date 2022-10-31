/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-31
 * @FilePath: /give-me-setu/util/util.go
 * @Description: 定义一些常用的无处安放的函数
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package util

import (
	"fmt"
	"os"
	"strings"
)

// 判断字符串是否以 target 开头
func IsStartWith(str string, target string) bool {
	return strings.Index(str, target) == 0
}

// 将字符串里的 old 替换成 args
func Replace(str string, old string, args ...any) string {
	str = strings.Replace(str, old, "%s", len(args))
	return fmt.Sprintf(str, args...)
}

// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}

}

// 如果文件夹不存在则创建
func InitDir(path string) {
	if !IsExist(path) {
		os.MkdirAll(path, 0755)
	}
}
