package common

import (
	"log"
	"os"
)

var logDebugEnabled = true
var logInfo = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
var logError = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
var logDebug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)
var logNull *log.Logger

func init() {
	devNull, _ := os.Open(os.DevNull)
	logNull = log.New(devNull, "[NULL]", log.LstdFlags)
}

// SetDebugEnabled sets/unsets global debug logs
func SetDebugEnabled(val bool) {
	logInfo.Print("Setting debug enabled: ", val)
	logDebugEnabled = val
}

// LogInfo returns an info-level logger
func LogInfo() *log.Logger {
	return logInfo
}

// LogError returns an error-level logger
func LogError() *log.Logger {
	return logError
}

// LogDebug returns a debug-level logger (or /dev/null if debug is disabled)
func LogDebug() *log.Logger {
	if logDebugEnabled {
		return logDebug
	}
	return logNull
}
