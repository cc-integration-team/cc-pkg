package logger

import "net/url"

// Masker provides explicit, condition-aware value masking for use at the call
// site. It complements the automatic field-name masking done by MaskingWriter:
//
//   - Use MaskingConfig.Fields for fields that are always sensitive (e.g. customerPhone).
//   - Use Masker when the masking decision depends on runtime context (e.g. only
//     mask callerID on inbound calls, calleeID on outbound calls).
//
// A zero-value Masker is safe to use and acts as a no-op (masking disabled).
type Masker struct {
	enabled bool
}

// NewMasker creates a Masker whose enabled state mirrors MaskingConfig.Enabled.
// Inject the single Masker instance wherever conditional masking is needed.
func NewMasker(cfg MaskingConfig) *Masker {
	return &Masker{enabled: cfg.Enabled}
}

// Mask returns the masked form of s ("******XYZ").
// Returns s unchanged when masking is disabled or s has 3 or fewer runes.
func (m *Masker) Mask(s string) string {
	if !m.enabled {
		return s
	}
	return maskPhone(s)
}

// MaskIf masks s only when condition is true, otherwise returns s as-is.
// Useful when the same field should be masked only in certain call directions
// or event types.
//
// Example:
//
//	"callerID": masker.MaskIf(callerID, direction == "Callin")
//	"calleeID": masker.MaskIf(calleeID, direction == "Callout")
func (m *Masker) MaskIf(s string, condition bool) string {
	if !condition {
		return s
	}
	return m.Mask(s)
}

// MaskURLParams returns rawURL with the specified query parameters replaced by "***".
// Returns rawURL unchanged when masking is disabled, params is empty, or the URL cannot be parsed.
//
// Example:
//
//	"url": masker.MaskURLParams(rawURL, "callerid", "token")
func (m *Masker) MaskURLParams(rawURL string, params ...string) string {
	if !m.enabled || len(params) == 0 {
		return rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	for _, p := range params {
		if q.Has(p) {
			q.Set(p, "***")
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// IsEnabled reports whether masking is active.
func (m *Masker) IsEnabled() bool {
	return m.enabled
}
