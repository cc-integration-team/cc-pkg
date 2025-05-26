package logger

import "context"

type nullLogger struct{}

func (l *nullLogger) Debug(msg string)                       {}
func (l *nullLogger) Debugf(format string, v ...any)         {}
func (l *nullLogger) Info(msg string)                        {}
func (l *nullLogger) Infof(format string, v ...any)          {}
func (l *nullLogger) Warn(msg string)                        {}
func (l *nullLogger) Warnf(format string, v ...any)          {}
func (l *nullLogger) Error(msg string)                       {}
func (l *nullLogger) Errorf(format string, v ...any)         {}
func (l *nullLogger) Fatal(msg string)                       {}
func (l *nullLogger) Fatalf(format string, v ...any)         {}
func (l *nullLogger) WithContext(ctx context.Context) Logger { return l }
