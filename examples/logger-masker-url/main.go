package main

import (
	"github.com/cc-integration-team/cc-pkg/v3/pkg/logger"
)

// client simulates an outbound HTTP client that dials a CTI endpoint.
type client struct {
	direction     string // "Callin" | "Callout"
	visibleSuffix int    // trailing chars to keep visible; 0 = hide all
}

func (c *client) paramsToMask() []string {
	switch c.direction {
	case "Callin":
		return []string{"callerid"}
	case "Callout":
		return []string{"calledid"}
	default:
		return nil
	}
}

func (c *client) connect(rawURL string) {
	logger.WithFields(logger.Fields{
		"direction": c.direction,
		"url":       logger.MaskURLParams(rawURL, c.visibleSuffix, c.paramsToMask()...),
	}).Info("connecting to endpoint")
}

func main() {
	cfg := logger.MaskingConfig{Enabled: true}

	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "cti-adapter",
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  false,
		},
		Masking: cfg,
	}))
	logger.SetDefaultMasker(cfg)

	// -------------------------------------------------------------------------
	// Ví dụ 1: INBOUND — ẩn hoàn toàn callerid (visibleSuffix=0)
	// "0901234567" → "***"
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: INBOUND, hide all ---")
	(&client{direction: "Callin", visibleSuffix: 0}).connect(
		"https://cti.example.com/dial?callerid=0901234567&calledid=74501&callid=100",
	)

	// -------------------------------------------------------------------------
	// Ví dụ 2: INBOUND — hiện 2 ký tự cuối callerid (visibleSuffix=2)
	// "0901234567" → "***67"
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: INBOUND, show last 2 chars ---")
	(&client{direction: "Callin", visibleSuffix: 2}).connect(
		"https://cti.example.com/dial?callerid=0901234567&calledid=74501&callid=101",
	)

	// -------------------------------------------------------------------------
	// Ví dụ 3: OUTBOUND — ẩn hoàn toàn calledid
	// "0987654321" → "***"
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: OUTBOUND, hide all ---")
	(&client{direction: "Callout", visibleSuffix: 0}).connect(
		"https://cti.example.com/dial?callerid=74501&calledid=0987654321&callid=102",
	)

	// -------------------------------------------------------------------------
	// Ví dụ 4: Direction không xác định — không mask gì cả
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: unknown direction ---")
	(&client{direction: "Unknown", visibleSuffix: 0}).connect(
		"https://cti.example.com/dial?callerid=0901234567&calledid=0987654321&callid=103",
	)

	// -------------------------------------------------------------------------
	// Ví dụ 5: Masking tắt — URL luôn nguyên vẹn
	// -------------------------------------------------------------------------
	logger.SetDefaultMasker(logger.MaskingConfig{Enabled: false})
	logger.Info("--- MaskURLParams: masking disabled ---")
	(&client{direction: "Callin", visibleSuffix: 2}).connect(
		"https://cti.example.com/dial?callerid=0901234567&calledid=74501&callid=104",
	)
}
