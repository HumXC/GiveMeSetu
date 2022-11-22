package storage

import (
	"log"

	"github.com/howeyc/fsnotify"
)

func NewWatcher(lib *Library) *LibWatcher {
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
	Library *Library
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
