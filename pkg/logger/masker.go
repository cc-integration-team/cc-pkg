package logger

import "net/url"

type masker struct {
	enabled bool
}

func newMasker(cfg MaskingConfig) *masker {
	return &masker{enabled: cfg.Enabled}
}

func (m *masker) mask(s string) string {
	if !m.enabled {
		return s
	}
	return maskPhone(s)
}

func (m *masker) maskIf(s string, condition bool) string {
	if !condition {
		return s
	}
	return m.mask(s)
}

// maskURLParams returns rawURL with the specified query parameters masked.
// visibleSuffix controls how many trailing characters remain visible (0 = hide all → "***").
// Returns rawURL unchanged when masking is disabled, params is empty, or the URL cannot be parsed.
func (m *masker) maskURLParams(rawURL string, visibleSuffix int, params ...string) string {
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
			q.Set(p, maskSuffix(q.Get(p), visibleSuffix))
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// maskSuffix masks s, keeping the last visibleSuffix runes visible.
// visibleSuffix <= 0 or >= len(s) always returns "***".
func maskSuffix(s string, visibleSuffix int) string {
	runes := []rune(s)
	if visibleSuffix <= 0 || visibleSuffix >= len(runes) {
		return "***"
	}
	return "***" + string(runes[len(runes)-visibleSuffix:])
}

func (m *masker) isEnabled() bool {
	return m.enabled
}
