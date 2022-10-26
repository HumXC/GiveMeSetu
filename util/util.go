/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-25
 * @FilePath: /give-me-setu/util/util.go
 * @Description: 定义一些常用的无处安放的函数
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package util

import "os"

/**
 * @description: 如果文件夹不存在则创建
 * @param {string} path
 * @return {*}
 */
func InitDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}
