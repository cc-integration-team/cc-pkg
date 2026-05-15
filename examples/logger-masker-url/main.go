package main

import (
	"github.com/cc-integration-team/cc-pkg/v3/pkg/logger"
)

// client simulates an outbound HTTP client that dials a CTI endpoint.
type client struct {
	masker    *logger.Masker
	direction string // "Callin" | "Callout"
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
		"url":       c.masker.MaskURLParams(rawURL, c.paramsToMask()...),
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

	masker := logger.NewMasker(cfg)

	// -------------------------------------------------------------------------
	// Ví dụ 1: INBOUND call — mask "callerid" (số khách hàng gọi vào)
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: INBOUND call ---")
	inbound := &client{masker: masker, direction: "Callin"}
	inbound.connect("https://cti.example.com/dial?callerid=0901234567&calledid=74501&callid=100")
	// url logged: https://cti.example.com/dial?calledid=74501&callerid=***&callid=100

	// -------------------------------------------------------------------------
	// Ví dụ 2: OUTBOUND call — mask "calledid" (số khách hàng được gọi ra)
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: OUTBOUND call ---")
	outbound := &client{masker: masker, direction: "Callout"}
	outbound.connect("https://cti.example.com/dial?callerid=74501&calledid=0987654321&callid=101")
	// url logged: https://cti.example.com/dial?calledid=***&callerid=74501&callid=101

	// -------------------------------------------------------------------------
	// Ví dụ 3: Direction không xác định — không mask gì cả
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: unknown direction ---")
	unknown := &client{masker: masker, direction: "Unknown"}
	unknown.connect("https://cti.example.com/dial?callerid=0901234567&calledid=0987654321&callid=102")
	// url logged: nguyên vẹn (paramsToMask trả nil)

	// -------------------------------------------------------------------------
	// Ví dụ 4: Masking tắt — URL luôn nguyên vẹn dù params có giá trị
	// -------------------------------------------------------------------------
	logger.Info("--- MaskURLParams: masking disabled ---")
	disabledMasker := logger.NewMasker(logger.MaskingConfig{Enabled: false})
	disabledClient := &client{masker: disabledMasker, direction: "Callin"}
	disabledClient.connect("https://cti.example.com/dial?callerid=0901234567&calledid=74501&callid=103")
	// url logged: nguyên vẹn vì masker.enabled = false
}
