package main

import (
	"github.com/cc-integration-team/cc-pkg/v3/pkg/logger"
)

func main() {
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "ncxo-ah-media-adapter:v1.0.1",
		File:    logger.LoggerFileConfig{},
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  false,
		},
	}))
	logger.WithFields(logger.Fields{
		"total": 300,
	}).Info("example message")
	logger.WithFields(logger.Fields{
		"total": 100,
	}).Debug("example message")
	logger.SetDefaultLogger(logger.NewZerologAdapter(logger.LoggerConfig{
		Service: "ncxo-ah-media-adapter:v1.0.1",
		File:    logger.LoggerFileConfig{},
		Caller:  true,
		Console: logger.LoggerConsoleConfig{
			Level:   "debug",
			Enabled: true,
			Pretty:  true,
		},
	}))
	logger.WithFields(logger.Fields{
		"abs": 123,
	}).WithFields(logger.Fields{
		"123": 123,
	}).Debug("example message")
}
