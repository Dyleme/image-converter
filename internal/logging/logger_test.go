package logging_test

import (
	"context"
	"testing"

	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	logger := logging.NewLogger("")
	assert.NotNil(t, logger)
}

func TestDefaultLogger(t *testing.T) {
	logger1 := logging.DefaultLogger()
	assert.NotNil(t, logger1)

	logger2 := logging.DefaultLogger()
	assert.NotNil(t, logger2)

	assert.Equal(t, logger1, logger2)
}

func TestContext(t *testing.T) {
	ctx := context.Background()

	defaultLogger := logging.DefaultLogger()

	logger := logging.FromContext(ctx)
	assert.Equal(t, logger, defaultLogger, "should be equal")

	newLogger := logging.NewLogger(logrus.DebugLevel.String())

	ctx = logging.WithLogger(ctx, newLogger)

	logger = logging.FromContext(ctx)
	assert.Equal(t, logger, newLogger)
}
