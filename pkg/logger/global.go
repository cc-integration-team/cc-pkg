package logger

var defaultLogger Logger = newSlogLogAdapter()

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func Debug(msg string) {
	defaultLogger.Debug(msg)
}

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Warn(msg string) {
	defaultLogger.Warn(msg)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

func WithFields(fields Fields) Logger {
	return defaultLogger.WithFields(fields)
}
