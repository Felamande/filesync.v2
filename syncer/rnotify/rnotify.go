package rnotify

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

type Watcher struct {
	SkipDirs []string
	watcher  *fsnotify.Watcher
	root     string
	Events   chan fsnotify.Event
}

func NewWatcher(root string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = w.Add(root)

	return &Watcher{
		watcher: w,
		root:    root,
		Events:  make(chan fsnotify.Event),
	}, err

}

func (w *Watcher) Add(path string) error {
	return w.watcher.Add(path)
}

func (w *Watcher) Skip(dirs ...string) *Watcher {
	w.SkipDirs = append(w.SkipDirs, dirs...)
	return w
}
func (w *Watcher) Start() (chan fsnotify.Event, chan error, error) {
	fmt.Println("root", w.root)
	err := filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}
		for _, dir := range w.SkipDirs {
			if filepath.Base(path) == dir {
				fmt.Println("skip", dir)
				return filepath.SkipDir
			}
		}
		err = w.watcher.Add(path)
		if err != nil {
			w.watcher.Errors <- err
		}
		fmt.Println("add watch", path)
		return nil
	})
	go func() {
		for {
			select {
			case e := <-w.watcher.Events:
				for _, dir := range w.SkipDirs {
					if filepath.Base(e.Name) == dir {
						goto end
					}
				}
				if e.Op == fsnotify.Create && isDir(e.Name) {
					w.watcher.Add(e.Name)
					fmt.Println("add watch", e.Name)
				}
				w.Events <- e
			end:
			}
		}
	}()
	return w.Events, w.watcher.Errors, err
}

func isDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
