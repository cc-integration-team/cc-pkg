package main

import (
	"github.com/cc-integration-team/cc-pkg/v3/pkg/logger"
)

// callEvent simulates a CTI call event from Cisco Finesse.
type callEvent struct {
	agentID       string
	eventName     string
	callID        string
	direction     string // "Callin" | "Callout"
	customerPhone string
	callerID      string
	calleeID      string
}

func main() {
	cfg := logger.MaskingConfig{
		Enabled: true,
		Fields:  []string{"customerPhone"}, // customerPhone always masked via MaskingWriter
	}

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

	logger.WithFields(logger.Fields{
		"isEnabled": logger.IsMaskingEnabled(),
	}).Info("masker state")

	// -------------------------------------------------------------------------
	// Ví dụ 1: INBOUND call
	// Customer là caller  → mask callerID
	// Agent là callee     → KHÔNG mask calleeID (là extension nội bộ)
	// -------------------------------------------------------------------------
	logger.Info("--- INBOUND call ---")
	logEvent(callEvent{
		agentID:       "agent_001",
		eventName:     "ALERTING",
		callID:        "26714589",
		direction:     "Callin",
		customerPhone: "0901234567",
		callerID:      "0901234567",
		calleeID:      "74501",
	})

	// -------------------------------------------------------------------------
	// Ví dụ 2: OUTBOUND call
	// Agent là caller     → KHÔNG mask callerID (là extension nội bộ)
	// Customer là callee  → mask calleeID
	// -------------------------------------------------------------------------
	logger.Info("--- OUTBOUND call ---")
	logEvent(callEvent{
		agentID:       "agent_001",
		eventName:     "ALERTING",
		callID:        "26714590",
		direction:     "Callout",
		customerPhone: "0987654321",
		callerID:      "74501",
		calleeID:      "0987654321",
	})

	// -------------------------------------------------------------------------
	// Ví dụ 3: Masking tắt — kiểm tra IsMaskingEnabled trước khi log
	// -------------------------------------------------------------------------
	logger.SetDefaultMasker(logger.MaskingConfig{Enabled: false})
	logger.Info("--- masking disabled ---")
	logger.WithFields(logger.Fields{
		"isEnabled":     logger.IsMaskingEnabled(),
		"customerPhone": logger.Mask("0901234567"),         // unchanged
		"callerID":      logger.MaskIf("0901234567", true), // unchanged
	}).Info("no masking applied")
}

// logEvent logs a call event using direction-aware masking via global helpers.
// customerPhone is always masked by MaskingWriter (field in config).
// callerID/calleeID are masked conditionally based on direction.
func logEvent(e callEvent) {
	logger.WithFields(logger.Fields{
		"agentID":       e.agentID,
		"eventName":     e.eventName,
		"callID":        e.callID,
		"direction":     e.direction,
		"customerPhone": e.customerPhone,
		"callerID":      logger.MaskIf(e.callerID, e.direction == "Callin"),
		"calleeID":      logger.MaskIf(e.calleeID, e.direction == "Callout"),
	}).Info("call event")
}
