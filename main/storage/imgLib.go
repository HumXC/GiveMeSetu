package storage

import (
	"give-me-setu/util"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type ImgLib struct {
	ParentLib *ImgLib            // 上一级文件夹 (库)
	SubLib    map[string]*ImgLib // 子文件夹 (库)
	Dir       string             // 库所在的位置
	Name      string             // 显示的名称
	Setus     map[string]any     // 所包含的媒体
}

// 添加图片到库
func (i *ImgLib) Add(srcName string, name string) error {
	f, err := os.Create(path.Join(i.Dir, name))
	if err != nil {
		return err
	}
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(f, src)
	if err != nil {
		return err
	}
	i.Setus[name] = nil
	return nil
}

// 返回指定路径的 lib，第二个返回值是无法进入的路径
func (i *ImgLib) Go(libName string) (*ImgLib, []string) {
	names := make([]string, 0)
	for _, v := range strings.Split(libName, "/") {
		if v != "" {
			names = append(names, v)
		}
	}
	var lib *ImgLib = i
	for index, name := range names {
		if v, ok := i.SubLib[name]; ok {
			lib = v
		} else {

			return lib, names[index:]
		}
	}
	return lib, nil
}

func GetLib(rootLibDir string) *ImgLib {
	return newLib(path.Dir(rootLibDir), path.Base(rootLibDir))
}

func (i *ImgLib) GetFile(setuName string) (io.ReadCloser, error) {
	file, err := os.OpenFile(path.Join(i.Dir, setuName), os.O_RDONLY, 0775)
	return file, err
}

func newLib(dir string, name string) *ImgLib {
	fullName := path.Join(dir, name)
	lib := ImgLib{
		Dir:    fullName,
		Name:   name,
		SubLib: make(map[string]*ImgLib),
		Setus:  make(map[string]any),
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
			if util.IsMIMEType(n, "image") {
				lib.Setus[n] = nil
			}
		}
	}
	return &lib
}
