// Package audit provides standardized logging for accountability
package audit

import (
	"log"
	"os"
)

// Logger provides structured audit logging
type Logger struct {
	logger *log.Logger
}

// NewLogger creates a new audit logger
func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "[AUDIT] ", log.LstdFlags|log.Lmicroseconds),
	}
}

// Info logs an informational message
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.log("INFO", msg, keyvals...)
}

// Error logs an error message
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.log("ERROR", msg, keyvals...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.log("WARN", msg, keyvals...)
}

func (l *Logger) log(level string, msg string, keyvals ...interface{}) {
	l.logger.Printf("[%s] %s %v", level, msg, keyvals)
}

