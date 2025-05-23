package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Manager[C any] struct {
	v              *viper.Viper
	path           string
	cfg            *C
	subscribersMap map[string][]func(*C) error
	mx             sync.RWMutex
}

func NewManager[C any](path string) (*Manager[C], error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cfg := new(C)

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	var loader Manager[C]
	loader.v = v
	loader.path = path
	loader.cfg = cfg
	loader.subscribersMap = make(map[string][]func(*C) error)

	return &loader, nil
}

func (l *Manager[C]) Config() *C {
	return l.cfg
}

func (l *Manager[C]) Watch(logger *zap.Logger) {
	l.v.WatchConfig()
	l.v.OnConfigChange(func(e fsnotify.Event) {
		newCfg := new(C)

		if err := l.v.Unmarshal(newCfg); err != nil {
			logger.Error("failed to unmarshal config", zap.Error(err))
		}

		logger.Info("config changed", zap.String("path", l.path))

		l.mx.RLock()
		for section, subscribers := range l.subscribersMap {
			for _, subscriber := range subscribers {
				if err := subscriber(newCfg); err != nil {
					logger.Error("subscriber error", zap.String("section", section), zap.Error(err))
				}
			}
		}
		l.mx.RUnlock()

		l.cfg = newCfg
	})
}

func (l *Manager[C]) Subscribe(key string, updateFn func(cfg *C) error) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.subscribersMap[key] = append(l.subscribersMap[key], updateFn)
}
