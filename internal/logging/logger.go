package logging

import (
	"context"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type contextKey string

const loggerKey = contextKey("logger")

func NewLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
}

var (
	defaultLogger     *logrus.Logger
	defaultLoggerOnce sync.Once
)

func DefaultLogger() *logrus.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger()
	})

	return defaultLogger
}

func WithLogger(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) *logrus.Logger {
	if logger, ok := ctx.Value(loggerKey).(*logrus.Logger); ok {
		return logger
	}

	return DefaultLogger()
}
