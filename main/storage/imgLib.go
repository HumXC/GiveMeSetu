/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-31
 * @FilePath: /give-me-setu/main/storage/imgLib.go
 * @Description: 图库
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package storage

import (
	"errors"
	"io"
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
	Name      string             // 库相对于根库的路径
	Setus     []string           // 所包含的媒体
}

func (i *ImgLib) Add(name string) {
	i.Setus = append(i.Setus, name)
}

// 返回指定路径的 lib
func (i *ImgLib) Go(libName string) (*ImgLib, error) {
	names := strings.Split(libName, "/")
	var lib *ImgLib = i
	for _, name := range names {
		if v, ok := i.SubLib[name]; ok {
			lib = v
		} else {
			return nil, errors.New("Can not found library: " + name)
		}
	}
	return lib, nil
}

func GetLib(rootLibDir string) *ImgLib {
	return newLib(path.Dir(rootLibDir), path.Base(rootLibDir))
}

func (i *ImgLib) GetStream(setuName string) (io.Reader, error) {
	file, err := os.OpenFile(path.Join(i.Dir, setuName), os.O_RDONLY, 0775)
	return file, err
}
func newLib(dir string, name string) *ImgLib {
	fullName := path.Join(dir, name)
	lib := ImgLib{
		Dir:    fullName,
		Name:   name,
		SubLib: make(map[string]*ImgLib),
	}

	entrys, err := os.ReadDir(fullName)
	if err != nil {
		log.Fatalf("Can not create ImgLib: %v", err)
	}
	for _, f := range entrys {
		n := f.Name()
		if f.IsDir() {
			// 创建子库
			subLib := newLib(fullName, path.Join(name, n))
			lib.SubLib[n] = subLib
		} else {
			t := mime.TypeByExtension(path.Ext(n))
			if strings.Contains(t, "image/") {
				lib.Setus = append(lib.Setus, n)
			}
		}
	}
	return &lib
}
