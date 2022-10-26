/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-26
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
	Path      string             // 库的位置
	Name      string             // 对外显示的名字
	Items     []string           // 所包含的媒体
}

func (i *ImgLib) Add(name string) {
	i.Items = append(i.Items, name)
}
func NewImgLib(dir string) *ImgLib {
	lib := ImgLib{
		Path:   dir,
		Name:   path.Base(dir),
		SubLib: make(map[string]*ImgLib),
	}
	entrys, err := os.ReadDir(dir)
	if err != nil {
		log.Panicf("Can not create ImgLib: %v", err)
	}
	for _, f := range entrys {
		name := f.Name()
		fullName := path.Join(dir, name)
		if f.IsDir() {
			// 创建子库
			subLib := NewImgLib(fullName)
			lib.SubLib[name] = subLib
		} else {
			t := mime.TypeByExtension(path.Ext(fullName))
			if strings.Contains(t, "image/") {
				lib.Items = append(lib.Items, name)
			}
		}
	}
	return &lib
}
