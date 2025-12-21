package utils

import (
	"fmt"
	"io"
	"sync"
)

// 日志器（安静原则：无用信息保持安静）
type Logger struct {
	writer  io.Writer
	silent  bool
	verbose bool
	mu      sync.Mutex
}

func NewLogger(w io.Writer, silent, verbose bool) *Logger {
	return &Logger{writer: w, silent: silent, verbose: verbose}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.silent {
		return
	}
	l.log("INFO", format, v...)
}

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
	fmt.Fprintf(l.writer, "[%s] %s\n", level, fmt.Sprintf(format, v...))
}
