/*
 * @Author: HumXC Hum-XC@outlook.com
 * @Date: 2022-10-25
 * @LastEditors: HumXC Hum-XC@outlook.com
 * @LastEditTime: 2022-10-26
 * @FilePath: /give-me-setu/main/storage/watcher.go
 * @Description: 文件夹的监听器, 图库的管理
 *
 * Copyright (c) 2022 by HumXC Hum-XC@outlook.com, All Rights Reserved.
 */
package storage

import (
	"log"

	"github.com/howeyc/fsnotify"
)

func NewWatcher(lib *ImgLib) *LibWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic("The watcher failed to initialize: " + err.Error())
	}

	return &LibWatcher{
		watcher: w,
		Library: lib,
	}
}

type LibWatcher struct {
	watcher *fsnotify.Watcher
	Library *ImgLib
}

func (w *LibWatcher) Watch() {
	err := w.watcher.Watch(w.Library.Dir)
	if err != nil {
		panic("The watcher failed to start watch: " + err.Error())
	}
	for {
		select {
		case e := <-w.watcher.Event:
			switch {
			case e.IsCreate():
				log.Printf("created: %s", e.Name)
				log.Printf("lib: %v", w.Library)
			}
		case e := <-w.watcher.Error:
			log.Panicf("watcher has a error: %v", e)
		}
	}
}
