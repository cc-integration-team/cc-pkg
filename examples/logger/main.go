package main

import (
	"errors"

	"github.com/cc-integration-team/cc-pkg/v2/pkg/logger"
)

func main() {
	// l := logger.NewZerologAdapter(logger.LoggerConfig{
	// 	Console: logger.LoggerConsoleConfig{
	// 		Level:   "debug",
	// 		Pretty:  true,
	// 		Enabled: true,
	// 	},
	// })

	// logger.SetDefaultLogger(l)

	// logger.WithFields(logger.Fields{
	// 	"module": "main",
	// 	"action": "test",
	// }).Info("This is an info message with fields")

	// reusedLogger := logger.WithFields(logger.Fields{
	// 	"module": "main",
	// 	"action": "reuse",
	// })
	// reusedLogger.Debug("This is a debug message from reused logger")
	// reusedLogger.Error("This is an error message from reused logger")
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		CallerDebug: true,
		File:        logger.LoggerFileConfig{},
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  true,
		},
	}))
	err := errors.New("e ")
	logger.WithFields(logger.Fields{
		"action": "test",
		"module": "main",
		"error":  err,
	}).Info("This is an info message with fields")
}
