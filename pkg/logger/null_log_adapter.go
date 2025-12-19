package logger

type nullLogAdapter struct {
	_ byte // prevent empty struct to compare equal
}

func newNullLogAdapter() Logger {
	return &nullLogAdapter{}
}

func (l *nullLogAdapter) Debug(msg string)                {}
func (l *nullLogAdapter) Debugf(format string, v ...any)  {}
func (l *nullLogAdapter) Info(msg string)                 {}
func (l *nullLogAdapter) Infof(format string, v ...any)   {}
func (l *nullLogAdapter) Warn(msg string)                 {}
func (l *nullLogAdapter) Warnf(format string, v ...any)   {}
func (l *nullLogAdapter) Error(msg string)                {}
func (l *nullLogAdapter) Errorf(format string, v ...any)  {}
func (l *nullLogAdapter) Fatal(msg string)                {}
func (l *nullLogAdapter) Fatalf(format string, v ...any)  {}
func (l *nullLogAdapter) WithFields(fields Fields) Logger { return l }
