package logger

import (
	"context"
	"testing"
)

func TestContextKey(t *testing.T) {
	ctx := context.Background()

	// Default logger should be null logger
	logger := FromContext(ctx)
	if logger == nil {
		t.Error("Expected default logger, got nil")
	}
	if logger != defaultLogger {
		t.Error("Expected default logger, got different logger")
	}

	// inject a custom logger
	nullLogger := newNullLogAdapter()
	t.Logf("defaultLogger addr: %p", defaultLogger)
	t.Logf("nullLogger addr: %p", nullLogger)
	if nullLogger == defaultLogger {
		t.Error("Expected different logger instances")
	}

	ctx = WithContext(ctx, nullLogger)
	logger1 := FromContext(ctx)
	if logger1 != nullLogger {
		t.Error("Expected injected logger, got different logger")
	}
	if logger1 == defaultLogger {
		t.Error("Expected injected logger, got default logger")
	}
}
