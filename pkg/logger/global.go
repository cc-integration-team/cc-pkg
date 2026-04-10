package logger

var defaultLogger Logger = newSlogLogAdapter()

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Debugf(msg string, args ...any) {
	defaultLogger.Debugf(msg, args...)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Infof(msg string, args ...any) {
	defaultLogger.Infof(msg, args...)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Warnf(msg string, args ...any) {
	defaultLogger.Warnf(msg, args...)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(msg string, args ...any) {
	defaultLogger.Errorf(msg, args...)
}

func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

func Fatalf(msg string, args ...any) {
	defaultLogger.Fatalf(msg, args...)
}

func WithFields(fields Fields) Logger {
	return defaultLogger.WithFields(fields)
}
