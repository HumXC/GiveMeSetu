package storage

import (
	"crypto/md5"
	"errors"
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
	Dir       string              // 库所的文件夹
	Name      string              // 文件夹的名称 path.Base(Dir)
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

// 在当前库创建一个子库
func (i *Library) CreateLib(name string) error {
	if _, ok := i.SubLib[name]; ok {
		return fmt.Errorf("folder \"%s\" is existed", name)
	}
	dir := path.Join(i.Dir, name)
	err := os.Mkdir(dir, 0775)
	if err != nil {
		return err
	}
	lib, err := newLib(dir)
	if err != nil {
		return err
	}
	i.SubLib[name] = lib
	return nil
}

func delInMap[T any | *Library](m map[string]T, name string) bool {
	_, ok := m[name]
	if ok {
		delete(m, name)
	}
	return ok
}

// 从 Setus 里删除记录，不会删除文件，如果文件还在原来的地方，下一次启动仍然会将他重新扫描到
// 如果删除成功，返回 true
// 如果返回 false 则表示 name 在 Setus 里不存在
func (i *Library) Rm(name string) bool {
	return delInMap(i.Setus, name)
}

// 删除一个库，跟 Rm() 一样，只是从 SubLib 中删除记录，不会删除磁盘上的文件
func (i *Library) RmLib(name string) bool {
	return delInMap(i.SubLib, name)
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
	if _, ok := i.Setus[name]; !ok {
		return nil, errors.New("can not find setu: " + name)
	}
	file, err := os.OpenFile(path.Join(i.Dir, name), os.O_RDONLY, 0775)
	return file, err
}

// 从根目录创建 Lib
func GetLib(rootLibDir string) (*Library, error) {
	return newLib(rootLibDir)
}

// 为 dir 创建一个库实例
func newLib(dir string) (*Library, error) {
	lib := Library{
		Dir:    dir,
		Name:   path.Base(dir),
		SubLib: make(map[string]*Library, 0),
		Setus:  make(map[string]any, 0),
	}
	entrys, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range entrys {
		n := f.Name()
		if f.IsDir() {
			// 创建子库
			subLib, err := newLib(path.Join(dir, n))
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
