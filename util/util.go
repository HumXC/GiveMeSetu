/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-27
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

func Replace(str string, target string, args ...any) string {
	str = strings.Replace(str, target, "%s", len(args))
	return fmt.Sprintf(str, args...)
}

/**
 * @description: 判断文件是否存在
 * @param {string} path
 * @return {*}
 */
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}

}

/**
 * @description: 如果文件夹不存在则创建
 * @param {string} path
 * @return {bool}
 */
func InitDir(path string) {
	if IsExist(path) {
		os.MkdirAll(path, 0755)
	}
}
