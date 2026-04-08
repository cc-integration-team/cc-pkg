package logger

import (
	"fmt"
	"log/slog"
)

type slogLogAdapter struct {
	slog *slog.Logger
}

func newSlogLogAdapter() Logger {
	return &slogLogAdapter{
		slog: slog.Default(),
	}
}

func (l *slogLogAdapter) Debug(msg string) {
	l.slog.Debug(msg)
}

func (l *slogLogAdapter) Info(msg string) {
	l.slog.Info(msg)
}

func (l *slogLogAdapter) Infof(format string, v ...any) {
	l.slog.Info(fmt.Sprintf(format, v...))
}

func (l *slogLogAdapter) Warn(msg string) {
	l.slog.Warn(msg)
}

func (l *slogLogAdapter) Error(msg string) {
	l.slog.Error(msg)
}

func (l *slogLogAdapter) Fatal(msg string) {
	l.slog.Error("unsupported log level: fatal")
}

func (l *slogLogAdapter) WithFields(fields Fields) Logger {
	newSlog := l.slog
	for k, v := range fields {
		newSlog = newSlog.With(k, v)
	}
	return &slogLogAdapter{
		slog: newSlog,
	}
}
