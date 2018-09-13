package dbe

import (
	"github.com/sedind/flow/logger"
)

var (
	logLevel = "debug"
	// Logger for dbe action
	Logger = logger.New(logLevel)
)

// LogLevel sets logging level
func LogLevel(level string) {
	logLevel = level
	Logger = logger.New(logLevel)
}
