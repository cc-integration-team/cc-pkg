package logger

type nulllogAdapter struct{}

func NewNulllogAdapter() *nulllogAdapter {
	return &nulllogAdapter{}
}

func (l *nulllogAdapter) Debug(msg string)                       {}
func (l *nulllogAdapter) Debugf(format string, v ...any)         {}
func (l *nulllogAdapter) Info(msg string)                        {}
func (l *nulllogAdapter) Infof(format string, v ...any)          {}
func (l *nulllogAdapter) Warn(msg string)                        {}
func (l *nulllogAdapter) Warnf(format string, v ...any)          {}
func (l *nulllogAdapter) Error(msg string)                       {}
func (l *nulllogAdapter) Errorf(format string, v ...any)         {}
func (l *nulllogAdapter) Fatal(msg string)                       {}
func (l *nulllogAdapter) Fatalf(format string, v ...any)         {}
