package utils_test

import (
	"testing"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

// 测试1: 自定义错误类型创建.
func TestCustomErrorType(t *testing.T) {
	err := utils.NewCryptoError(utils.ErrInvalidKey, "key validation failed")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Code != utils.ErrInvalidKey {
		t.Errorf("Expected code %v, got %v", utils.ErrInvalidKey, err.Code)
	}
}

// 测试2: 错误上下文.
func TestErrorContext(t *testing.T) {
	ctx := utils.ErrorContext{
		Operation: "encrypt",
		Position:  1024,
		File:      "test.txt",
	}
	err := ctx.Wrap(utils.ErrIOError, "read failed")
	expected := "encrypt at position 1024 in test.txt: read failed"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

// 测试3: 错误分类.
func TestErrorClassification(t *testing.T) {
	err := utils.NewCryptoError(utils.ErrInvalidMagic, "bad magic")
	if !utils.IsFormatError(err) {
		t.Error("Should be format error")
	}

	err = utils.NewCryptoError(utils.ErrAuthFailed, "auth failed")
	if !utils.IsSecurityError(err) {
		t.Error("Should be security error")
	}
}

// 测试4: 错误代码枚举值.
func TestErrorCodes(t *testing.T) {
	// 验证错误代码在预期范围内
	codes := []utils.ErrorCode{
		utils.ErrSystem,
		utils.ErrIOError,
		utils.ErrInvalidMagic,
		utils.ErrInvalidFormat,
		utils.ErrVersionMismatch,
		utils.ErrInvalidKey,
		utils.ErrKeyGenerationFailed,
		utils.ErrAuthFailed,
		utils.ErrSignatureVerification,
		utils.ErrHashMismatch,
		utils.ErrInvalidParameter,
		utils.ErrFileNotFound,
	}

	for _, code := range codes {
		if code < utils.ErrSystem || code > utils.ErrFileNotFound {
			t.Errorf("Error code %v out of expected range", code)
		}
	}
}

// 测试5: 错误消息格式.
func TestErrorMessageFormat(t *testing.T) {
	err := utils.NewCryptoError(utils.ErrInvalidKey, "test message")
	expected := "test message"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

// 测试6: 空错误上下文.
func TestEmptyErrorContext(t *testing.T) {
	ctx := utils.ErrorContext{}
	err := ctx.Wrap(utils.ErrSystem, "system error")
	expected := " at position 0 in : system error"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

// 测试7: 错误分类边界.
func TestErrorClassificationBoundaries(t *testing.T) {
	// 格式错误边界
	formatErr := utils.NewCryptoError(utils.ErrVersionMismatch, "version")
	if !utils.IsFormatError(formatErr) {
		t.Error("Version mismatch should be format error")
	}

	// 安全错误边界
	securityErr := utils.NewCryptoError(utils.ErrHashMismatch, "hash")
	if !utils.IsSecurityError(securityErr) {
		t.Error("Hash mismatch should be security error")
	}

	// 非错误类型
	if utils.IsFormatError(nil) {
		t.Error("nil should not be format error")
	}
}
