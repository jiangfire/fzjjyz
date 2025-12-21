package utils

import (
	"testing"
)

// 测试1: 自定义错误类型创建
func TestCustomErrorType(t *testing.T) {
	err := NewCryptoError(ErrInvalidKey, "key validation failed")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Code != ErrInvalidKey {
		t.Errorf("Expected code %v, got %v", ErrInvalidKey, err.Code)
	}
}

// 测试2: 错误上下文
func TestErrorContext(t *testing.T) {
	ctx := ErrorContext{
		Operation: "encrypt",
		Position:  1024,
		File:      "test.txt",
	}
	err := ctx.Wrap(ErrIOError, "read failed")
	expected := "encrypt at position 1024 in test.txt: read failed"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

// 测试3: 错误分类
func TestErrorClassification(t *testing.T) {
	err := NewCryptoError(ErrInvalidMagic, "bad magic")
	if !IsFormatError(err) {
		t.Error("Should be format error")
	}

	err = NewCryptoError(ErrAuthFailed, "auth failed")
	if !IsSecurityError(err) {
		t.Error("Should be security error")
	}
}

// 测试4: 错误代码枚举值
func TestErrorCodes(t *testing.T) {
	// 验证错误代码在预期范围内
	codes := []ErrorCode{
		ErrSystem,
		ErrIOError,
		ErrInvalidMagic,
		ErrInvalidFormat,
		ErrVersionMismatch,
		ErrInvalidKey,
		ErrKeyGenerationFailed,
		ErrAuthFailed,
		ErrSignatureVerification,
		ErrHashMismatch,
		ErrInvalidParameter,
		ErrFileNotFound,
	}

	for _, code := range codes {
		if code < ErrSystem || code > ErrFileNotFound {
			t.Errorf("Error code %v out of expected range", code)
		}
	}
}

// 测试5: 错误消息格式
func TestErrorMessageFormat(t *testing.T) {
	err := NewCryptoError(ErrInvalidKey, "test message")
	expected := "test message"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

// 测试6: 空错误上下文
func TestEmptyErrorContext(t *testing.T) {
	ctx := ErrorContext{}
	err := ctx.Wrap(ErrSystem, "system error")
	expected := " at position 0 in : system error"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

// 测试7: 错误分类边界
func TestErrorClassificationBoundaries(t *testing.T) {
	// 格式错误边界
	formatErr := NewCryptoError(ErrVersionMismatch, "version")
	if !IsFormatError(formatErr) {
		t.Error("Version mismatch should be format error")
	}

	// 安全错误边界
	securityErr := NewCryptoError(ErrHashMismatch, "hash")
	if !IsSecurityError(securityErr) {
		t.Error("Hash mismatch should be security error")
	}

	// 非错误类型
	if IsFormatError(nil) {
		t.Error("nil should not be format error")
	}
}
