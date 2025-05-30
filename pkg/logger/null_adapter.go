package logger

import "context"

type nullAdapter struct{}

func (l *nullAdapter) Debug(msg string)                       {}
func (l *nullAdapter) Debugf(format string, v ...any)         {}
func (l *nullAdapter) Info(msg string)                        {}
func (l *nullAdapter) Infof(format string, v ...any)          {}
func (l *nullAdapter) Warn(msg string)                        {}
func (l *nullAdapter) Warnf(format string, v ...any)          {}
func (l *nullAdapter) Error(msg string)                       {}
func (l *nullAdapter) Errorf(format string, v ...any)         {}
func (l *nullAdapter) Fatal(msg string)                       {}
func (l *nullAdapter) Fatalf(format string, v ...any)         {}
func (l *nullAdapter) WithContext(ctx context.Context) Logger { return l }
