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
		Fields:  []string{"customerPhone"}, // customerPhone always masked via auto-masking
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

	// Create a single Masker instance from the same config.
	// Inject this wherever conditional masking is needed in the app.
	masker := logger.NewMasker(cfg)

	logger.Info("--- Masker: masking disabled check ---")
	logger.WithFields(logger.Fields{
		"isEnabled": masker.IsEnabled(),
	}).Info("masker state")

	// -------------------------------------------------------------------------
	// Ví dụ 1: INBOUND call
	// Customer là caller  → mask callerID
	// Agent là callee     → KHÔNG mask calleeID (là extension nội bộ)
	// -------------------------------------------------------------------------
	inbound := callEvent{
		agentID:       "agent_001",
		eventName:     "ALERTING",
		callID:        "26714589",
		direction:     "Callin",
		customerPhone: "0901234567",
		callerID:      "0901234567", // customer
		calleeID:      "74501",      // agent extension
	}

	logger.Info("--- Masker: INBOUND call ---")
	logEvent(inbound, masker)

	// -------------------------------------------------------------------------
	// Ví dụ 2: OUTBOUND call
	// Agent là caller     → KHÔNG mask callerID (là extension nội bộ)
	// Customer là callee  → mask calleeID
	// -------------------------------------------------------------------------
	outbound := callEvent{
		agentID:       "agent_001",
		eventName:     "ALERTING",
		callID:        "26714590",
		direction:     "Callout",
		customerPhone: "0987654321",
		callerID:      "74501",      // agent extension
		calleeID:      "0987654321", // customer
	}

	logger.Info("--- Masker: OUTBOUND call ---")
	logEvent(outbound, masker)

	// -------------------------------------------------------------------------
	// Ví dụ 3: Masking tắt — tất cả giữ nguyên
	// -------------------------------------------------------------------------
	disabledMasker := logger.NewMasker(logger.MaskingConfig{Enabled: false})

	logger.Info("--- Masker: disabled ---")
	logger.WithFields(logger.Fields{
		"customerPhone": disabledMasker.Mask("0901234567"),         // unchanged
		"callerID":      disabledMasker.MaskIf("0901234567", true), // unchanged
	}).Info("no masking applied")

	// -------------------------------------------------------------------------
	// Ví dụ 4: MaskIf với điều kiện động
	// -------------------------------------------------------------------------
	logger.Info("--- Masker: MaskIf with dynamic condition ---")

	for _, direction := range []string{"Callin", "Callout"} {
		logger.WithFields(logger.Fields{
			"direction": direction,
			"callerID":  masker.MaskIf("0901234567", direction == "Callin"),
			"calleeID":  masker.MaskIf("0987654321", direction == "Callout"),
		}).Info("direction-based masking")
	}
}

// logEvent logs a call event using Masker for direction-aware masking.
// customerPhone is always masked (handled by MaskingWriter via config).
// callerID/calleeID are masked conditionally based on direction.
func logEvent(e callEvent, masker *logger.Masker) {
	logger.WithFields(logger.Fields{
		"agentID":       e.agentID,
		"eventName":     e.eventName,
		"callID":        e.callID,
		"direction":     e.direction,
		"customerPhone": e.customerPhone, // auto-masked by MaskingWriter (field in config)
		"callerID":      masker.MaskIf(e.callerID, e.direction == "Callin"),
		"calleeID":      masker.MaskIf(e.calleeID, e.direction == "Callout"),
	}).Info("call event")
}
