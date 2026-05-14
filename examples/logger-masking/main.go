package main

import (
	"github.com/cc-integration-team/cc-pkg/v3/pkg/logger"
)

func main() {
	// -------------------------------------------------------------------------
	// Ví dụ 1: Masking top-level fields
	// Các field customerPhone, callerID, calleeID sẽ bị mask từ level Info trở lên.
	// -------------------------------------------------------------------------
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "example-service",
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  false,
		},
		Masking: logger.MaskingConfig{
			Enabled: true,
			Fields:  []string{"customerPhone", "callerID", "calleeID"},
		},
	}))

	logger.Info("--- Example 1: top-level field masking ---")

	// Info: customerPhone bị mask → "******567"
	// Output: {"level":"info","customerPhone":"******567","callerID":"******678","message":"inbound call"}
	logger.WithFields(logger.Fields{
		"customerPhone": "0901234567",
		"callerID":      "0912345678",
		"calleeID":      "100",
	}).Info("inbound call")

	// Debug: KHÔNG mask, giữ nguyên giá trị gốc để dễ troubleshoot
	// Output: {"level":"debug","customerPhone":"0901234567","callerID":"0912345678","message":"raw debug"}
	logger.WithFields(logger.Fields{
		"customerPhone": "0901234567",
		"callerID":      "0912345678",
	}).Debug("raw debug")

	// -------------------------------------------------------------------------
	// Ví dụ 2: Masking nested fields trong metadata map
	// callVariable8 bên trong metadata bị mask, các key khác giữ nguyên.
	// -------------------------------------------------------------------------
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "example-service",
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  false,
		},
		Masking: logger.MaskingConfig{
			Enabled: true,
			Fields:  []string{"customerPhone"},
			NestedFields: map[string][]string{
				"metadata": {"callVariable8"},
			},
		},
	}))

	logger.Info("--- Example 2: nested field masking ---")

	// Output:
	// {
	//   "level": "info",
	//   "customerPhone": "******567",          ← top-level: masked
	//   "metadata": {
	//     "callVariable8": "******909",         ← nested: masked
	//     "callVariable3": "other-value",       ← không trong config: giữ nguyên
	//     "queueName":     "queue_support"      ← không trong config: giữ nguyên
	//   },
	//   "message": "call with metadata"
	// }
	logger.WithFields(logger.Fields{
		"customerPhone": "0901234567",
		"metadata": map[string]any{
			"callVariable8": "0909090909",
			"callVariable3": "other-value",
			"queueName":     "queue_support",
		},
	}).Info("call with metadata")

	// -------------------------------------------------------------------------
	// Ví dụ 3: Masking tắt — toàn bộ giá trị ghi nguyên
	// -------------------------------------------------------------------------
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "example-service",
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  false,
		},
		Masking: logger.MaskingConfig{
			Enabled: false, // tắt masking
		},
	}))

	logger.Info("--- Example 3: masking disabled ---")

	// Output: {"level":"info","customerPhone":"0901234567","message":"no masking"}
	logger.WithFields(logger.Fields{
		"customerPhone": "0901234567",
	}).Info("no masking")

	// -------------------------------------------------------------------------
	// Ví dụ 4: Pretty console + masking (dùng cho môi trường dev)
	// -------------------------------------------------------------------------
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "example-service",
		Caller:  true,
		Console: logger.LoggerConsoleConfig{
			Level:   "info",
			Enabled: true,
			Pretty:  true,
		},
		Masking: logger.MaskingConfig{
			Enabled: true,
			Fields:  []string{"customerPhone", "callerID"},
		},
	}))

	logger.Info("--- Example 4: pretty console + masking ---")

	logger.WithFields(logger.Fields{
		"customerPhone": "0901234567",
		"callerID":      "0912345678",
		"agentID":       "agent_001",
		"direction":     "INBOUND",
	}).Info("call assigned")
}
