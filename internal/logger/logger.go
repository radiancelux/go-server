package logger

import (
	"log"
	"os"
)

// ServerLogger implements the Logger interface
type ServerLogger struct {
	logger *log.Logger
}

// NewServerLogger creates a new server logger
func NewServerLogger() *ServerLogger {
	return &ServerLogger{
		logger: log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile),
	}
}

// Info logs an info message
func (l *ServerLogger) Info(msg string, args ...any) {
	l.logger.Printf("[INFO] "+msg, args...)
}

// Error logs an error message
func (l *ServerLogger) Error(msg string, args ...any) {
	l.logger.Printf("[ERROR] "+msg, args...)
}

// Debug logs a debug message
func (l *ServerLogger) Debug(msg string, args ...any) {
	l.logger.Printf("[DEBUG] "+msg, args...)
}
