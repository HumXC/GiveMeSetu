package storage

import (
	"crypto/md5"
	"fmt"
	"give-me-setu/util"
	"io"
	"os"
	"path"
	"strings"
)

type Setu struct {
	Name string
	Md5  string
}
type Library struct {
	ParentLib *Library            // 上一级文件夹 (库)
	SubLib    map[string]*Library // 子文件夹 (库)
	Dir       string              // 库所在的位置
	Name      string              // 文件夹的名称
	Setus     map[string]any      // 所包含的媒体
}

// 添加图片到库
// srcName 是需要添加的文件, extName 是文件扩展名
// 如果成功添加将返回该文件的 md5 校验和
func (i *Library) Add(srcName, extName string) (string, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return "", err
	}
	defer src.Close()
	buf, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", md5.Sum(buf))
	f, err := os.Create(path.Join(i.Dir, sum+extName))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, src)
	if err != nil {
		return "", err
	}

	i.Setus[sum] = sum
	return sum, nil
}

// 返回指定路径的 lib，第二个返回值是无法进入的路径(如果有)
func (i *Library) Go(libName string) (*Library, []string) {
	names := make([]string, 0)
	for _, v := range strings.Split(libName, "/") {
		if v != "" {
			names = append(names, v)
		}
	}
	var lib *Library = i
	for index, name := range names {
		if v, ok := i.SubLib[name]; ok {
			lib = v
		} else {
			return lib, names[index:]
		}
	}
	return lib, []string{}
}

// 获取文件
func (i *Library) GetFile(setuName string) (io.ReadCloser, error) {
	file, err := os.OpenFile(path.Join(i.Dir, setuName), os.O_RDONLY, 0775)
	return file, err
}

// 从根目录创建 Lib
func GetLib(rootLibDir string) (*Library, error) {
	return newLib(path.Dir(rootLibDir), path.Base(rootLibDir))
}

// 创建一个库，dir 是库的文件夹
func newLib(dir string, name string) (*Library, error) {
	fullName := path.Join(dir, name)
	lib := Library{
		Dir:    fullName,
		Name:   name,
		SubLib: make(map[string]*Library),
		Setus:  make(map[string]any),
	}
	entrys, err := os.ReadDir(fullName)
	if err != nil {
		return nil, err
	}
	for _, f := range entrys {
		n := f.Name()
		if f.IsDir() {
			// 创建子库
			subLib, err := newLib(fullName, n)
			if err != nil {
				return nil, err
			}
			lib.SubLib[n] = subLib
		} else {
			if util.IsMIMEType(n, "image") {
				lib.Setus[n] = nil
			}
		}
	}
	return &lib, nil
}
