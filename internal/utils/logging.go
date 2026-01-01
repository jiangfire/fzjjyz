// Package utils provides utility functions.
//
//nolint:revive // utils 在 internal 上下文中含义清晰
package utils

import (
	"fmt"
	"io"
	"sync"
)

// Logger provides logging capabilities with configurable verbosity.
type Logger struct {
	writer  io.Writer
	silent  bool
	verbose bool
	mu      sync.Mutex
}

// NewLogger creates a new logger instance.
func NewLogger(w io.Writer, silent, verbose bool) *Logger {
	return &Logger{writer: w, silent: silent, verbose: verbose}
}

// Info logs an informational message.
func (l *Logger) Info(format string, v ...interface{}) {
	if l.silent {
		return
	}
	l.log("INFO", format, v...)
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.silent || !l.verbose {
		return
	}
	l.log("DEBUG", format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log("ERROR", format, v...)
}

func (l *Logger) log(level, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _ = fmt.Fprintf(l.writer, "[%s] %s\n", level, fmt.Sprintf(format, v...))
}
