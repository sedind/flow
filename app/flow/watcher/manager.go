package watcher

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/sedind/flow/app/flow/config"
	"github.com/sedind/flow/app/flow/logger"
)

// Manager watcher manager object
type Manager struct {
	*config.WatcherConfig
	Logger     *logger.Logger
	Restart    chan bool
	once       *sync.Once
	ID         string
	context    context.Context
	cancelFunc context.CancelFunc
}

// New creates new watcher manager instance
func New(c *config.WatcherConfig) *Manager {
	return NewWithContext(context.Background(), c)
}

// NewWithContext creates new watcher manager instance
func NewWithContext(ctx context.Context, c *config.WatcherConfig) *Manager {
	ctx, cancelFunc := context.WithCancel(ctx)
	m := &Manager{
		WatcherConfig: c,
		Logger:        logger.New(c.Name, true),
		Restart:       make(chan bool),
		once:          &sync.Once{},
		ID:            ID(c.Name),
		context:       ctx,
		cancelFunc:    cancelFunc,
	}
	return m
}

// Start watcher
func (m *Manager) Start() error {
	m.Logger.Success(fmt.Sprintf("Start watcher %s", m.Name))
	w := NewWatcher(m)
	w.Start()
	go func() {
		for {
			select {
			case event := <-w.Events:
				if event.Op != fsnotify.Chmod {
					go m.build(event)
				}
				w.Remove(event.Name)
				w.Add(event.Name)
			case err := <-w.Errors:
				m.Logger.Error(err)
			case <-m.context.Done():
				break
			}
		}
	}()

	m.runner()
	return nil
}

func (m *Manager) build(event fsnotify.Event) {
	m.Logger.Success(event.Name)
}

// ID generates md5 has of current working directory
func ID(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}
