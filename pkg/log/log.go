package log

import (
	"strings"
	"sync"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"

	"go.uber.org/zap"
)

var (
	_ abstractions.ILogger                   = (*Logger)(nil)
	_ abstractions.ConfigSubscriber[*Config] = (*Logger)(nil)
)

type Config struct {
	Level string `mapstructure:"level"`
}

type Logger struct {
	zap   *zap.Logger
	level zap.AtomicLevel
	cfg   *Config
	mx    sync.Mutex
}

func (l *Logger) SectionKey() string {
	return "LOG"
}

func (l *Logger) UpdateConfig(newCfg *Config) error {
	l.mx.Lock()
	defer l.mx.Unlock()

	if newCfg.Level != l.cfg.Level {
		level := levelFromString(newCfg.Level)
		l.level.SetLevel(level.Level())
	}

	l.cfg = newCfg

	return nil
}

func NewLogger(local bool, level string) *Logger {
	var logger Logger

	atomicLevel := levelFromString(level)

	var cfg zap.Config
	if local {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.DisableStacktrace = true
	cfg.Level = atomicLevel
	cfg.OutputPaths = []string{"stdout"}
	logger.zap, _ = cfg.Build(zap.WithCaller(true))
	logger.level = atomicLevel

	return &logger
}

func (l *Logger) Close() error {
	if err := l.zap.Sync(); err != nil && !isStdoutSyncErr(err) {
		return err
	}

	return nil
}

func isStdoutSyncErr(err error) bool {
	return strings.Contains(err.Error(), "sync")
}
