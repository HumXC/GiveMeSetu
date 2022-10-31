/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-31
 * @FilePath: /give-me-setu/main/main.go
 * @Description: main
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package main

import (
	"give-me-setu/main/conf"
	"give-me-setu/main/database"
	"give-me-setu/main/network"
	"give-me-setu/util"
	"os"
	"path"
)

var Cfg conf.Config

func init() {
	dataDir := ""
	if len(os.Args) > 1 && os.Args[1] != "" {
		dataDir = os.Args[1]
	} else if os.Getenv("DATADIR") != "" {
		dataDir = os.Getenv("DATADIR")
	} else {
		dir, err := os.Executable()
		if err != nil {
			panic(err)
		}
		dataDir = path.Join(path.Dir(dir), "data")
	}
	util.InitDir(dataDir)
	Cfg = *conf.Get(path.Join(dataDir, "config.yaml"))
	Cfg.DataDir = dataDir
	Cfg.Library = path.Join(dataDir, "library")
	util.InitDir(Cfg.Library)
	_ = database.Get(Cfg)

}
func main() {

	network.NewServer(Cfg.Library).Run("12345")

}
