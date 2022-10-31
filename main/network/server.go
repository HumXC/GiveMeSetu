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
	"log"
	"mime"
	"net/http"
	"path"
)

//go:embed icon.png
var icon []byte

type SetuServer struct {
	serverMux *http.ServeMux
	lib       storage.ImgLib
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
	}

	s.serverMux.HandleFunc("/ping", s.ping)
	s.serverMux.HandleFunc("/favicon.ico", setIcon)
	s.serverMux.HandleFunc("/library/", s.library)

	return s
}

func (s *SetuServer) ping(w http.ResponseWriter, r *http.Request) {
	setContentType(w, s.lib.Setus[0])
	f, err := s.lib.GetStream(s.lib.Setus[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b := make([]byte, 128)
	for {
		_, err := f.Read(b)
		w.Write(b)
		if err != nil {
			break
		}
	}
}
func (s *SetuServer) library(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("library"))

}

func setContentType(w http.ResponseWriter, fileName string) {
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(fileName)))
}

func setIcon(w http.ResponseWriter, r *http.Request) {
	setContentType(w, ".png")
	w.Write(icon)
}
