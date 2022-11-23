package main

import (
	"give-me-setu/conf"
	"give-me-setu/network"
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
}
func main() {
	s, err := network.NewServer(Cfg.Library)
	if err != nil {
		panic(err)
	}
	err = s.Run("12345")
	if err != nil {
		panic(err)
	}
}
