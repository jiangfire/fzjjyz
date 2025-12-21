// 遵循清晰原则：简单明了的错误定义
package utils

import "fmt"

// 错误代码枚举（表达原则：数据结构优先）
type ErrorCode int

const (
	// 系统级错误
	ErrSystem ErrorCode = iota
	ErrIOError

	// 格式错误
	ErrInvalidMagic
	ErrInvalidFormat
	ErrVersionMismatch
	ErrInvalidVersion
	ErrInvalidAlgorithm

	// 加密错误
	ErrInvalidKey
	ErrKeyGenerationFailed
	ErrInvalidData
	ErrEncryptionFailed
	ErrDecryptionFailed
	ErrSerializationFailed

	// 验证错误
	ErrAuthFailed
	ErrSignatureVerification
	ErrHashMismatch
	ErrSigningFailed
	ErrVerificationFailed

	// 用户错误
	ErrInvalidParameter
	ErrFileNotFound
)

// 自定义错误结构（透明原则：清晰状态）
type CryptoError struct {
	Code    ErrorCode
	Message string
}

func (e *CryptoError) Error() string {
	return e.Message
}

// 错误上下文（模块原则：可组合）
type ErrorContext struct {
	Operation string
	Position  int64
	File      string
}

func (ctx *ErrorContext) Wrap(code ErrorCode, msg string) error {
	return &CryptoError{
		Code:    code,
		Message: fmt.Sprintf("%s at position %d in %s: %s",
			ctx.Operation, ctx.Position, ctx.File, msg),
	}
}

// 工厂函数（修复原则：及早抛出明确异常）
func NewCryptoError(code ErrorCode, msg string) *CryptoError {
	return &CryptoError{Code: code, Message: msg}
}

// 错误分类函数
func IsFormatError(err error) bool {
	if ce, ok := err.(*CryptoError); ok {
		return ce.Code >= ErrInvalidMagic && ce.Code <= ErrInvalidAlgorithm
	}
	return false
}

func IsSecurityError(err error) bool {
	if ce, ok := err.(*CryptoError); ok {
		return ce.Code >= ErrInvalidKey && ce.Code <= ErrHashMismatch
	}
	return false
}
