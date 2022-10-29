/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-27
 * @FilePath: /give-me-setu/main/storage/imgLib.go
 * @Description: 图库
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package storage

import (
	"log"
	"mime"
	"os"
	"path"
	"strings"
)

type ImgLib struct {
	ParentLib *ImgLib            // 上一级文件夹 (库)
	SubLib    map[string]*ImgLib // 子文件夹 (库)
	Dir       string             // 库所在的位置（父目录）
	Name      string             // 库相对与根库的路径
	Items     []string           // 所包含的媒体
}

func (i *ImgLib) Add(name string) {
	i.Items = append(i.Items, name)
}
func Get(rootLibDir string) *ImgLib {
	return newLib(path.Dir(rootLibDir), path.Base(rootLibDir))
}

func newLib(dir string, name string) *ImgLib {
	lib := ImgLib{
		Dir:    dir,
		Name:   name,
		SubLib: make(map[string]*ImgLib),
	}
	fullName := path.Join(dir, name)
	entrys, err := os.ReadDir(fullName)
	if err != nil {
		log.Panicf("Can not create ImgLib: %v", err)
	}
	for _, f := range entrys {
		n := f.Name()
		if f.IsDir() {
			// 创建子库
			subLib := newLib(fullName, path.Join(name, n))
			lib.SubLib[n] = subLib
		} else {
			t := mime.TypeByExtension(path.Ext(fullName))
			if strings.Contains(t, "image/") {
				lib.Items = append(lib.Items, n)
			}
		}
	}
	return &lib
}
