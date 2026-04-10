package logger

import "context"

type contextKey struct{}

var (
	loggerKey = contextKey{}
)

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return defaultLogger
	}
	return logger
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Logger defines the interface that all loggers must implement
type Logger interface {
	Debug(msg string)
	Debugf(msg string, args ...any)
	Info(msg string)
	Infof(msg string, args ...any)
	Warn(msg string)
	Warnf(msg string, args ...any)
	Error(msg string)
	Errorf(msg string, args ...any)
	Fatal(msg string)
	Fatalf(msg string, args ...any)
	WithFields(fields Fields) Logger
}

type Fields map[string]any
