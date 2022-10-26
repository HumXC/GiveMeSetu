/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-26
 * @FilePath: /give-me-setu/main/main.go
 * @Description: main
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package main

import (
	"fmt"
	"give-me-setu/main/storage"
	"give-me-setu/util"
)

var imgLib string = "../data/img"

func init() {
	util.InitDir(imgLib)

}
func main() {
	lib := storage.NewImgLib(imgLib)
	fmt.Println(lib)
}
