package logger

import (
	"encoding/json"

	"github.com/rs/zerolog"
)

type maskingLevelWriter struct {
	underlying zerolog.LevelWriter
	fields     map[string]struct{}
	nested     map[string]map[string]struct{} // parent → set of child keys
}

func newMaskingLevelWriter(underlying zerolog.LevelWriter, cfg MaskingConfig) *maskingLevelWriter {
	fields := make(map[string]struct{}, len(cfg.Fields))
	for _, f := range cfg.Fields {
		fields[f] = struct{}{}
	}

	nested := make(map[string]map[string]struct{}, len(cfg.NestedFields))
	for parent, children := range cfg.NestedFields {
		childSet := make(map[string]struct{}, len(children))
		for _, c := range children {
			childSet[c] = struct{}{}
		}
		nested[parent] = childSet
	}

	return &maskingLevelWriter{
		underlying: underlying,
		fields:     fields,
		nested:     nested,
	}
}

// WriteLevel implements zerolog.LevelWriter.
// Debug logs bypass masking entirely.
func (w *maskingLevelWriter) WriteLevel(l zerolog.Level, p []byte) (int, error) {
	if l < zerolog.InfoLevel {
		return w.underlying.WriteLevel(l, p)
	}
	masked := w.maskJSON(p)
	n, err := w.underlying.WriteLevel(l, masked)
	// Always report the original length so zerolog doesn't treat a short write as an error.
	if err == nil {
		n = len(p)
	}
	return n, err
}

// Write satisfies io.Writer; delegates without masking because level is unknown at this point.
func (w *maskingLevelWriter) Write(p []byte) (int, error) {
	return w.underlying.Write(p)
}

// maskJSON parses one JSON log line, masks configured fields, and re-serialises.
// Returns the original bytes unchanged if parsing fails or nothing was masked.
func (w *maskingLevelWriter) maskJSON(p []byte) []byte {
	if len(w.fields) == 0 && len(w.nested) == 0 {
		return p
	}

	// Trim trailing newline so json.Unmarshal doesn't fail.
	trimmed := p
	hasNewline := len(p) > 0 && p[len(p)-1] == '\n'
	if hasNewline {
		trimmed = p[:len(p)-1]
	}

	// Use RawMessage to avoid deserialising values we don't need to touch.
	var entry map[string]json.RawMessage
	if err := json.Unmarshal(trimmed, &entry); err != nil {
		return p
	}

	changed := false

	// Mask top-level fields.
	for field := range w.fields {
		if raw, ok := entry[field]; ok {
			if masked, ok := maskRawString(raw); ok {
				entry[field] = masked
				changed = true
			}
		}
	}

	// Mask nested fields (e.g. metadata.callVariable8).
	for parent, children := range w.nested {
		raw, ok := entry[parent]
		if !ok {
			continue
		}
		var sub map[string]json.RawMessage
		if err := json.Unmarshal(raw, &sub); err != nil {
			continue
		}
		subChanged := false
		for child := range children {
			if childRaw, ok := sub[child]; ok {
				if masked, ok := maskRawString(childRaw); ok {
					sub[child] = masked
					subChanged = true
				}
			}
		}
		if subChanged {
			newRaw, err := json.Marshal(sub)
			if err != nil {
				continue
			}
			entry[parent] = newRaw
			changed = true
		}
	}

	if !changed {
		return p
	}

	result, err := json.Marshal(entry)
	if err != nil {
		return p
	}
	if hasNewline {
		result = append(result, '\n')
	}
	return result
}

// maskRawString takes a json.RawMessage that should be a JSON string,
// masks all but the last 3 characters, and returns the new raw value.
// Returns (nil, false) if the value is not a non-empty JSON string.
func maskRawString(raw json.RawMessage) (json.RawMessage, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil || s == "" {
		return nil, false
	}
	masked, err := json.Marshal(maskPhone(s))
	if err != nil {
		return nil, false
	}
	return masked, true
}

const maskPrefix = "******"

// maskPhone masks a phone string to "******XYZ" where XYZ are the last 3 runes.
// Strings with 3 or fewer runes are returned unchanged to avoid revealing more than the original.
// The prefix is always 6 asterisks regardless of original length, so log readers
// cannot infer the phone length from the masked output.
func maskPhone(s string) string {
	runes := []rune(s)
	if len(runes) <= 3 {
		return s
	}
	return maskPrefix + string(runes[len(runes)-3:])
}
