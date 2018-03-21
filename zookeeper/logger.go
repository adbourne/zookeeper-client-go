package zookeeper

import (
	"fmt"
	golog "log"
)

const (
	// levelKey is the key used by the logger to denote the logging level e.g info, debug etc
	levelKey = "level"

	// messageKey is the key used by the logger to denote the log message
	messageKey = "message"

	// urlKey is the key used by the logger to denote the URL used
	urlKey = "url"

	// requestKey is the key used by the logger to denote the request sent
	requestKey = "request"

	// errorKey is the key used by the logger to denote the error
	errorKey = "error"

	// statusKey is the key used by the logger to denote the status
	statusKey = "status"

	// LogLevelDebug is the value for the DEBUG log level
	levelDebug = "debug"

	// levelWarn is the value for the WARN log level
	levelWarn = "warn"
)

// LoggerFunc is an abstraction over the logger.
// It allows for different loggers to be used.
//
// The client log using key value pairs e.g:
//     Log(
//         "message", "Some Message",
//         "level", "info",
//         "foo", "bar",
//     )
//
// If your chosen logger does not accept key value pairs, then the implementation of this function should
// should convert the key value pairs to the appropriate string.
type LoggerFunc interface {
	Log(string, ...interface{})
}

// StdOutLogger is a simple implementation of LoggerFunc which writes to StdOut.
type StdOutLogger struct {
}

// Log redirects to the standard Go logger.
func (logger *StdOutLogger) Log(s string, args ...interface{}) {
	keyValues := []string{fmt.Sprintf("%s=", s)}
	isKey := false
	for _, arg := range args {
		delimiter := "%s  "
		if isKey {
			delimiter = "%s="
		}
		keyValues = append(keyValues, fmt.Sprintf(delimiter, arg))
		isKey = !isKey
	}

	logLine := ""
	for _, kv := range keyValues {
		logLine = fmt.Sprintf("%s%s", logLine, kv)
	}

	golog.Print(logLine)
}

// NewStdOutLogger creates a new StdOutLogger
func NewStdOutLogger() LoggerFunc {
	logger := &StdOutLogger{}
	logger.Log("message", "Starting logger...")
	return logger
}
