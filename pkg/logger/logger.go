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
	Debugf(format string, v ...any)
	Info(msg string)
	Infof(format string, v ...any)
	Warn(msg string)
	Warnf(format string, v ...any)
	Error(msg string)
	Errorf(format string, v ...any)
	Fatal(msg string)
	Fatalf(format string, v ...any)
	WithFields(fields Fields) Logger
}

type Fields map[string]any
