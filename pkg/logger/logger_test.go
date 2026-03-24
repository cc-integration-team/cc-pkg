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
	slogAdapter := newSlogLogAdapter()
	t.Logf("defaultLogger addr: %p", defaultLogger)
	t.Logf("slogAdapter addr: %p", slogAdapter)
	if slogAdapter == defaultLogger {
		t.Error("Expected different logger instances")
	}

	ctx = WithContext(ctx, slogAdapter)
	logger1 := FromContext(ctx)
	if logger1 != slogAdapter {
		t.Error("Expected injected logger, got different logger")
	}
	if logger1 == defaultLogger {
		t.Error("Expected injected logger, got default logger")
	}
}
