package log

import (
	"context"

	"github.com/DavidMovas/gopherbox/pkg/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap    *zap.Logger
	closer *closer.Closer
}

func NewLogger(local bool, level string) (*Logger, error) {
	logger := &Logger{}

	c := closer.NewCloser()
	logger.closer = c

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

	c.Push(logger.zap.Sync)

	return logger, nil
}

func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

func (l *Logger) Stop() error {
	return l.closer.Close(context.Background())
}

func levelFromString(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}
