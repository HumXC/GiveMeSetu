/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-26
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-26
 * @FilePath: /give-me-setu/main/conf/config.go
 * @Description: 从环境变量和文件获取配置
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package conf

import (
	_ "embed"
	"give-me-setu/util"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed config.yaml
var configFile []byte

type Config struct {
	Library  string   `yaml:"library"`
	Database Database `yaml:"database"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func GetConfig(path string) *Config {
	var c Config

	if util.IsExit(path) {
		file, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		configFile = file
	} else {
		err := os.WriteFile(path, configFile, 0775)
		if err != nil {
			panic(err)
		}
	}
	err := yaml.Unmarshal(configFile, &c)
	if err != nil {
		panic(err)
	}
	return &c
}
