package logger

var defaultLogger Logger = newSlogLogAdapter()
var defaultMasker = &masker{}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

// SetDefaultMasker initialises the package-level masker from cfg.
// Call once at boot alongside SetDefaultLogger.
func SetDefaultMasker(cfg MaskingConfig) {
	defaultMasker = newMasker(cfg)
}

// Mask returns the masked form of s ("******XYZ").
// Returns s unchanged when masking is disabled or s has 3 or fewer runes.
func Mask(s string) string {
	return defaultMasker.mask(s)
}

// MaskIf masks s only when condition is true, otherwise returns s as-is.
func MaskIf(s string, condition bool) string {
	return defaultMasker.maskIf(s, condition)
}

// MaskURLParams masks the specified query parameters in rawURL.
// visibleSuffix controls how many trailing characters remain visible (0 = hide all → "***").
// Returns rawURL unchanged when masking is disabled or no params are given.
func MaskURLParams(rawURL string, visibleSuffix int, params ...string) string {
	return defaultMasker.maskURLParams(rawURL, visibleSuffix, params...)
}

// IsMaskingEnabled reports whether the package-level masker is active.
func IsMaskingEnabled() bool {
	return defaultMasker.isEnabled()
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
