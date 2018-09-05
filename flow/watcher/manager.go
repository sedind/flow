package watcher

import (
	"context"
	"crypto/md5"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/sedind/flow/flow/config"
	"github.com/sedind/flow/flow/logger"
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
	go m.build(fsnotify.Event{Name: ":start:"})
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
	m.once.Do(func() {
		defer func() {
			m.once = &sync.Once{}
		}()

		m.buildTransaction(func() error {
			now := time.Now()
			m.Logger.Print("Rebuild on : %s", event.Name)
			if m.ChangeCommand == "" {
				m.Logger.Print("`change_command` not provided in flow.yml")
				return nil
			}

			cmd := exec.Command(m.ChangeCommand, m.ChangeArgs...)

			err := m.runAndListen(cmd)
			if err != nil {
				return err
			}
			tt := time.Since(now)
			m.Logger.Success("Change command Completed (PID: %d) (TIME: %s)", cmd.Process.Pid, tt)
			m.Restart <- true
			return nil
		})
	})
}

func (m *Manager) buildTransaction(fn func() error) {
	err := fn()
	if err != nil {
		m.Logger.Error("Error!")
		m.Logger.Error(err)
	}
}

// ID generates md5 has of current working directory
func ID(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}
