/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-31
 * @FilePath: /give-me-setu/main/network/server.go
 * @Description: 服务端
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package network

import (
	_ "embed"
	"give-me-setu/main/storage"
	"give-me-setu/util"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

//go:embed icon.png
var icon []byte

type SetuServer struct {
	serverMux *http.ServeMux
	lib       storage.ImgLib
	client    http.Client
}

func (s *SetuServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("A request from %s: %s", r.Host, r.URL)
	s.serverMux.ServeHTTP(w, r)
}
func (s *SetuServer) Run(port string) error {
	return http.ListenAndServe(":"+port, s)
}

func NewServer(rootLibDir string) *SetuServer {
	s := &SetuServer{
		serverMux: http.NewServeMux(),
		lib:       *storage.GetLib(rootLibDir),
		client:    *http.DefaultClient,
	}

	s.serverMux.HandleFunc("/ping", s.ping)
	s.serverMux.HandleFunc("/favicon.ico", setIcon)
	s.serverMux.HandleFunc("/library/add", s.libraryAdd)

	return s
}

func (s *SetuServer) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func (s *SetuServer) libraryAdd(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	imgUrl := query.Get("url")
	libName := query.Get("lib")
	name := query.Get("name")

	switch {
	case name == "":
		// TODO: 从文件内容判断文件类型，而不是从文件名
		name_ := path.Base(imgUrl)
		if util.IsMIMEType(name_, "image") {
			name = name_
		} else {
			w.Write([]byte("\"name\" is not a image"))
			return
		}

	case imgUrl == "":
		// TODO: 支持上传文件，而不是使用url
		w.Write([]byte("\"url\" is empty"))
		return

	}

	lib, err := s.lib.Go(libName)
	if err != nil {
		log.Printf("Can not find lib: %v", err)
		w.Write([]byte(err.Error()))
		return
	}
	if _, ok := lib.Setus[name]; ok {
		w.Write([]byte("File \"" + name + "\" is exist"))
		return
	}
	os.TempDir()
	resp, err := s.client.Get(imgUrl)
	if err != nil {
		log.Printf("Failed to get web image: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file, err := os.CreateTemp("", "give-me-setu")
	if err != nil {
		log.Printf("Failed to create temp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	defer os.Remove(path.Join(os.TempDir(), file.Name()))
	b := make([]byte, 120)
	for {
		_, err := resp.Body.Read(b)
		file.Write(b)

		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				log.Printf("Failed to get web image: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	}
	lib.Add(file, name)
	w.Write([]byte("ok"))

}

func setContentType(w http.ResponseWriter, fileName string) {
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(fileName)))
}

func setIcon(w http.ResponseWriter, r *http.Request) {
	setContentType(w, ".png")
	w.Write(icon)
}
