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
	_ "embed"
	"fmt"
	"give-me-setu/main/conf"
	"give-me-setu/util"
	"os"
	"path"
)

var DataDir string

func init() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		DataDir = os.Args[1]
	} else {
		dir, err := os.Executable()
		if err != nil {
			panic(err)
		}
		DataDir = path.Join(path.Dir(dir), "data")
	}

	util.InitDir(DataDir)
	c := conf.GetConfig(path.Join(DataDir, "config.yaml"))
	fmt.Println(c)
}
func main() {
}
