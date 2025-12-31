// Package utils provides error handling and logging utilities.
//
//nolint:revive
package utils

import (
	"errors"
	"fmt"
)

// ErrorCode defines error types for the crypto system.
type ErrorCode int

const (
	// ErrSystem represents system-level errors.
	ErrSystem ErrorCode = iota
	// ErrIOError represents input/output errors.
	ErrIOError

	// ErrInvalidMagic represents invalid magic number errors.
	ErrInvalidMagic
	// ErrInvalidFormat represents invalid format errors.
	ErrInvalidFormat
	ErrVersionMismatch
	ErrInvalidVersion
	ErrInvalidAlgorithm

	// ErrInvalidKey represents encryption-related errors.
	ErrInvalidKey
	ErrKeyGenerationFailed
	ErrInvalidData
	ErrEncryptionFailed
	ErrDecryptionFailed
	ErrSerializationFailed

	// ErrAuthFailed represents verification errors.
	ErrAuthFailed
	ErrSignatureVerification
	ErrHashMismatch
	ErrSigningFailed
	ErrVerificationFailed

	// ErrInvalidParameter represents user errors.
	ErrInvalidParameter
	ErrFileNotFound
)

// CryptoError represents a custom error with code and message.
type CryptoError struct {
	Code    ErrorCode
	Message string
}

func (e *CryptoError) Error() string {
	return e.Message
}

// ErrorContext provides context for errors.
type ErrorContext struct {
	Operation string
	Position  int64
	File      string
}

// Wrap wraps an error with context.
func (ctx *ErrorContext) Wrap(code ErrorCode, msg string) error {
	return &CryptoError{
		Code: code,
		Message: fmt.Sprintf("%s at position %d in %s: %s",
			ctx.Operation, ctx.Position, ctx.File, msg),
	}
}

// NewCryptoError creates a new crypto error.
func NewCryptoError(code ErrorCode, msg string) *CryptoError {
	return &CryptoError{Code: code, Message: msg}
}

// IsFormatError checks if error is a format error.
func IsFormatError(err error) bool {
	ce := &CryptoError{}
	if errors.As(err, &ce) {
		return ce.Code >= ErrInvalidMagic && ce.Code <= ErrInvalidAlgorithm
	}
	return false
}

// IsSecurityError checks if error is a security error.
func IsSecurityError(err error) bool {
	ce := &CryptoError{}
	if errors.As(err, &ce) {
		return ce.Code >= ErrInvalidKey && ce.Code <= ErrHashMismatch
	}
	return false
}
