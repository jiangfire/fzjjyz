package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestSilentMode(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, true, false) // silent=true, verbose=false

	logger.Info("should not appear")
	logger.Error("should appear")

	output := buf.String()
	if strings.Contains(output, "should not appear") {
		t.Error("Silent mode should suppress info")
	}
	if !strings.Contains(output, "should appear") {
		t.Error("Error should always appear")
	}
}

func TestVerboseMode(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, false, true) // silent=false, verbose=true

	logger.Info("info message")
	logger.Debug("debug message")

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Error("Verbose mode should show info")
	}
	if !strings.Contains(output, "debug message") {
		t.Error("Verbose mode should show debug")
	}
}

func TestNormalMode(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, false, false) // silent=false, verbose=false

	logger.Info("info message")
	logger.Debug("debug message")
	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Error("Normal mode should show info")
	}
	if strings.Contains(output, "debug message") {
		t.Error("Normal mode should not show debug")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Normal mode should show error")
	}
}

func TestLoggerConcurrency(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, false, false)

	// 测试并发安全性
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message %d", id)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有消息都已写入
	output := buf.String()
	count := strings.Count(output, "concurrent message")
	if count != 10 {
		t.Errorf("Expected 10 concurrent messages, got %d", count)
	}
}

func TestLoggerFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, false, false)

	logger.Info("test %s %d", "message", 123)

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Error("Output should contain level prefix")
	}
	if !strings.Contains(output, "test message 123") {
		t.Error("Format should be applied correctly")
	}
}

func TestLoggerEmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, false, false)

	logger.Info("")
	logger.Error("")

	output := buf.String()
	// Should not panic, just output empty messages
	if !strings.Contains(output, "[INFO]") || !strings.Contains(output, "[ERROR]") {
		t.Error("Empty messages should still produce output")
	}
}
