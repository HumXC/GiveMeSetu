/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-26
 * @FilePath: /give-me-setu/util/util.go
 * @Description: 定义一些常用的无处安放的函数
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package util

import "os"

/**
 * @description: 判断文件是否存在
 * @param {string} path
 * @return {*}
 */
func IsExit(path string) bool {
	_, err := os.Stat(path)
	if os.IsExist(err) {
		return true
	} else {
		return false
	}

}

/**
 * @description: 如果文件夹不存在则创建
 * @param {string} path
 * @return {bool}
 */
func InitDir(path string) {
	if IsExit(path) {
		os.MkdirAll(path, 0755)
	}
}
