package logger

import (
	"log"
	"os"
	"runtime"
)

// Logger is the global logger instance
var Logger *log.Logger

func init() {
	// Create or open the log file
	logFile, err := os.OpenFile("logging.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Initialize the logger
	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

// LogWithTrace logs a message with the current execution trace
func LogWithTrace(msg string) {
	// Get the current caller information
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		Logger.Printf("ERROR: Unable to retrieve caller information\n")
		return
	}

	// Get the function name
	funcName := runtime.FuncForPC(pc).Name()

	// Log the message with trace details
	Logger.Printf("[%s:%d %s] %s\n", file, line, funcName, msg)
}

// LogErrorWithTrace logs an error with stack trace
func LogErrorWithTrace(err error) {
	if err == nil {
		return
	}

	// Get the current caller information
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		Logger.Printf("ERROR: Unable to retrieve caller information\n")
		return
	}

	// Get the function name
	funcName := runtime.FuncForPC(pc).Name()

	// Log the error with trace details
	Logger.Printf("[ERROR: %s:%d %s] %v\n", file, line, funcName, err)
}
