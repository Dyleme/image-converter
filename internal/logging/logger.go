package logging

import (
	"context"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type contextKey string

const loggerKey = contextKey("logger")

// NewLogger returns configurated logrus.Logger.
func NewLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
}

var (
	defaultLogger     *logrus.Logger
	defaultLoggerOnce sync.Once
)

// DefaultLogger returns logger.
// Logger can be created by this function only once.
func DefaultLogger() *logrus.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger()
	})

	return defaultLogger
}

// WithLogger function put logger into context.
func WithLogger(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext function get logger from context.
// or returns default logger if context is empty.
func FromContext(ctx context.Context) *logrus.Logger {
	if logger, ok := ctx.Value(loggerKey).(*logrus.Logger); ok {
		return logger
	}

	return DefaultLogger()
}
