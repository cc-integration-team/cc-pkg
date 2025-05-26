package logger

import "context"

var defaultLogger Logger = &nullLogger{}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Debugf(format string, v ...any) {
	defaultLogger.Debugf(format, v...)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Infof(format string, v ...any) {
	defaultLogger.Infof(format, v...)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Warnf(format string, v ...any) {
	defaultLogger.Warnf(format, v...)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(format string, v ...any) {
	defaultLogger.Errorf(format, v...)
}

func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

func Fatalf(format string, v ...any) {
	defaultLogger.Fatalf(format, v...)
}

func WithContext(ctx context.Context) Logger {
	return defaultLogger.WithContext(ctx)
}
