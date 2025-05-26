package logger

import "context"

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
	WithContext(ctx context.Context) Logger
}
