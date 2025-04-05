package log

import (
	"strings"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"

	"go.uber.org/zap"
)

var _ abstractions.ILogger = (*Logger)(nil)

type Logger struct {
	zap *zap.Logger
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
