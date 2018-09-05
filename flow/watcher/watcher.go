package watcher

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher -
type Watcher struct {
	*fsnotify.Watcher
	*Manager
	context context.Context
}

//NewWatcher creates watcher instance for given manager
func NewWatcher(m *Manager) *Watcher {
	w, _ := fsnotify.NewWatcher()
	return &Watcher{
		Watcher: w,
		Manager: m,
		context: m.context,
	}
}

// Start watching process for given watch configuration
func (w *Watcher) Start() {
	go func() {
		for {
			err := filepath.Walk(w.Watch, func(path string, info os.FileInfo, err error) error {
				if info == nil {
					w.cancelFunc()
					return errors.New("nil directory")
				}
				if info.IsDir() {
					if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") || w.IsIgnored(path) {
						return filepath.SkipDir
					}
				}
				if w.IsWatchedFile(path) {
					w.Add(path)
				}
				return nil
			})

			if err != nil {
				w.context.Done()
				break
			}
			// sweep for new files every 1 second

			time.Sleep(1 * time.Second)
		}
	}()
}

// IsIgnored check if path is ignored by config
func (w *Watcher) IsIgnored(path string) bool {

	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	for _, e := range w.Ignore {
		p := filepath.FromSlash(w.Watch + "/" + e)
		if strings.HasPrefix(p, "./") {
			p = strings.Replace(p, "./", "", 1)
		}
		if strings.HasPrefix(path, strings.TrimSpace(p)) {
			return true
		}
	}
	return false
}

// IsWatchedFile check if file is watched
func (w *Watcher) IsWatchedFile(file string) bool {
	ext := filepath.Ext(file)
	for _, e := range w.Extensions {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}
