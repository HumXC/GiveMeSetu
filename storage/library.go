package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type Library struct {
	LibDB     *SetuDBs            // 数据库表操作
	ParentLib *Library            // 上一级文件夹 (库)
	SubLib    map[string]*Library // 子文件夹 (库)
	Dir       string              // 库所在的位置
	Name      string              // 文件夹的名称
	Setus     map[string]any      // 所包含的媒体
}

// 添加文件添加到库文件夹，并 md5 字符串添加到 Setus
// srcName 是需要添加的文件, extName 是文件扩展名, origin 是来源，可以是一个 url
// 如果成功添加将返回该文件的 md5 校验和，也是文件名
func (i *Library) Add(srcName string) (string, error) {
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
	f, err := os.Create(path.Join(i.Dir, sum))
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.Write(buf)
	i.Setus[sum] = nil
	return sum, nil
}

// 从 Setus 里删除记录，不会删除文件
// 如果删除成功，返回 true
// 如果返回 false 则表示 name 在 Setus 里不存在
func (i *Library) Del(name string) bool {
	_, ok := i.Setus[name]
	if ok {
		delete(i.Setus, name)
	}
	return ok
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
func (i *Library) GetFile(name string) (io.ReadCloser, error) {
	file, err := os.OpenFile(path.Join(i.Dir, name), os.O_RDONLY, 0775)
	return file, err
}

// 从根目录创建 Lib
func GetLib(rootLibDir string) (*Library, error) {
	return newLib(path.Dir(rootLibDir), path.Base(rootLibDir))
}

// 创建一个库，dir 是库的文件夹
func newLib(dir, name string) (*Library, error) {
	fullName := path.Join(dir, name)
	lib := Library{
		Dir:    fullName,
		Name:   name,
		SubLib: make(map[string]*Library, 0),
		Setus:  make(map[string]any, 0),
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
			lib.Setus[n] = nil
		}
	}
	return &lib, nil
}
