package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// bufWriter is an in-memory zerolog.LevelWriter used in tests.
type bufWriter struct{ bytes.Buffer }

func (b *bufWriter) WriteLevel(_ zerolog.Level, p []byte) (int, error) {
	return b.Write(p)
}

func newTestWriter(cfg MaskingConfig) (*maskingLevelWriter, *bufWriter) {
	buf := &bufWriter{}
	return newMaskingLevelWriter(buf, cfg), buf
}

// parseLog unmarshals a single JSON log line from the buffer.
func parseLog(t *testing.T, buf *bufWriter) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(bytes.TrimRight(buf.Bytes(), "\n"), &m); err != nil {
		t.Fatalf("failed to parse log output: %v\nraw: %s", err, buf.String())
	}
	return m
}

// --- maskPhone unit tests ---

func TestMaskPhone(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"0901234567", "******567"},    // normal 10-digit
		{"08123456789", "******789"},   // 11-digit: still 6 stars
		{"1234", "******234"},          // 4-digit: still 6 stars
		{"123", "123"},                 // ≤3 chars: unchanged
		{"12", "12"},                   // ≤3 chars: unchanged
		{"", ""},                       // empty: unchanged
	}
	for _, tc := range cases {
		if got := maskPhone(tc.input); got != tc.want {
			t.Errorf("maskPhone(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// --- maskJSON unit tests ---

func TestMaskJSON_TopLevelField(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone", "callerID"},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","customerPhone":"0901234567","callerID":"0912345678","message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	if got := m["customerPhone"]; got != "******567" {
		t.Errorf("customerPhone = %q, want ******567", got)
	}
	if got := m["callerID"]; got != "******678" {
		t.Errorf("callerID = %q, want ******678", got)
	}
	if got := m["message"]; got != "test" {
		t.Errorf("message should be unchanged, got %q", got)
	}
}

func TestMaskJSON_NestedField(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		NestedFields: map[string][]string{
			"metadata": {"callVariable8"},
		},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","metadata":{"callVariable8":"0909090909","callVariable3":"other","queueName":"support"},"message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	meta, ok := m["metadata"].(map[string]any)
	if !ok {
		t.Fatal("metadata should be a map")
	}
	if got := meta["callVariable8"]; got != "******909" {
		t.Errorf("callVariable8 = %q, want ******909", got)
	}
	// Non-configured keys must remain untouched.
	if got := meta["callVariable3"]; got != "other" {
		t.Errorf("callVariable3 = %q, want 'other'", got)
	}
	if got := meta["queueName"]; got != "support" {
		t.Errorf("queueName = %q, want 'support'", got)
	}
}

func TestMaskJSON_DebugBypassesMasking(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"debug","customerPhone":"0901234567","message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.DebugLevel, []byte(line))

	m := parseLog(t, buf)
	// Debug level: value must NOT be masked.
	if got := m["customerPhone"]; got != "0901234567" {
		t.Errorf("debug log should not be masked, got customerPhone=%q", got)
	}
}

func TestMaskJSON_DisabledPassesThrough(t *testing.T) {
	// When masking is disabled, newMaskingLevelWriter is never called in the adapter,
	// but we test the writer directly with an empty field set to ensure safety.
	cfg := MaskingConfig{Enabled: false}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","customerPhone":"0901234567","message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	if got := m["customerPhone"]; got != "0901234567" {
		t.Errorf("no fields configured, value should be unchanged, got %q", got)
	}
}

func TestMaskJSON_EmptyPhoneUnchanged(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","customerPhone":"","message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	if got := m["customerPhone"]; got != "" {
		t.Errorf("empty string should be unchanged, got %q", got)
	}
}

func TestMaskJSON_NonStringFieldSkipped(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"count"}, // numeric field — should not be touched
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","count":42,"message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	if got := m["count"]; got != float64(42) {
		t.Errorf("numeric field should be unchanged, got %v", got)
	}
}

func TestMaskJSON_InvalidJSONPassesThrough(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"},
	}
	w, buf := newTestWriter(cfg)

	line := `not valid json` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	if got := strings.TrimRight(buf.String(), "\n"); got != "not valid json" {
		t.Errorf("invalid JSON should pass through unchanged, got %q", got)
	}
}

func TestMaskJSON_FieldNotPresent(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","message":"no phone here"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	if _, exists := m["customerPhone"]; exists {
		t.Error("customerPhone should not be injected into log entry")
	}
	if got := m["message"]; got != "no phone here" {
		t.Errorf("message unchanged, got %q", got)
	}
}

func TestMaskJSON_NewlinePreserved(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","customerPhone":"0901234567","message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("output should end with newline")
	}
}

func TestMaskJSON_TopLevelAndNestedCombined(t *testing.T) {
	cfg := MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"},
		NestedFields: map[string][]string{
			"metadata": {"callVariable8"},
		},
	}
	w, buf := newTestWriter(cfg)

	line := `{"level":"info","customerPhone":"0901234567","metadata":{"callVariable8":"0909090909","queueName":"q1"},"message":"test"}` + "\n"
	_, _ = w.WriteLevel(zerolog.InfoLevel, []byte(line))

	m := parseLog(t, buf)
	if got := m["customerPhone"]; got != "******567" {
		t.Errorf("customerPhone = %q, want ******567", got)
	}
	meta := m["metadata"].(map[string]any)
	if got := meta["callVariable8"]; got != "******909" {
		t.Errorf("callVariable8 = %q, want ******909", got)
	}
	if got := meta["queueName"]; got != "q1" {
		t.Errorf("queueName should be unchanged, got %q", got)
	}
}
