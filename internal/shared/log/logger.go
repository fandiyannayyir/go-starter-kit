package log

import (
	"context"
	"io"
)

type Logger interface {
	WithError(err any) Logger
	WithField(key string, val any) Logger
	WithFields(fields map[string]any) Logger
	WithContext(ctx context.Context) Logger

	Debug(message string)
	Error(message string)
	Fatal(message string)
	Info(message string)
	Warn(message string)

	GetWriter() io.Writer
	Close() error
}
