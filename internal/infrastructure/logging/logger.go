package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// LogLevel defines the severity of a log message.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of a LogLevel.
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides a simple, opinionated logging interface.
type Logger struct {
	mu        sync.Mutex
	stdLogger *log.Logger
	minLevel  LogLevel
}

// NewLogger creates a new Logger instance.
// output is the io.Writer where log messages will be written (e.g., os.Stdout, a file).
// minLevel sets the minimum level for messages to be logged.
func NewLogger(output io.Writer) *Logger {
	return &Logger{
		stdLogger: log.New(output, "", log.Ldate|log.Ltime|log.Lshortfile),
		minLevel:  LevelInfo, // Default minimum level
	}
}

// SetMinLevel sets the minimum log level for the logger.
func (l *Logger) SetMinLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// logf formats and writes a log message if its level meets the minimum level.
func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	prefix := fmt.Sprintf("[%s] ", level.String())
	l.stdLogger.Output(3, prefix+fmt.Sprintf(format, v...)) // 3 skips logf, log functions themselves
}

// Debugf logs a debug message.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(LevelDebug, format, v...)
}

// Infof logs an info message.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(LevelInfo, format, v...)
}

// Warnf logs a warning message.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(LevelWarn, format, v...)
}

// Errorf logs an error message.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(LevelError, format, v...)
}

// Fatalf logs a fatal message and then exits the application.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logf(LevelFatal, format, v...)
	os.Exit(1)
}
