package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (l *Logger) Zap() *zap.Logger {
	return l.zap
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
