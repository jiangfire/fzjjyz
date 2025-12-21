package utils

import (
	"fmt"
	"strings"
)

// ErrorType 定义错误类型
type ErrorType string

const (
	ErrorInvalidInput    ErrorType = "输入错误"
	ErrorKeyNotFound     ErrorType = "密钥文件不存在"
	ErrorKeyInvalid      ErrorType = "密钥格式无效"
	ErrorFileNotFound    ErrorType = "文件不存在"
	ErrorFileExists      ErrorType = "文件已存在"
	ErrorPermission      ErrorType = "权限不足"
	ErrorEncryption      ErrorType = "加密失败"
	ErrorDecryption      ErrorType = "解密失败"
	ErrorSignature       ErrorType = "签名验证失败"
	ErrorIO              ErrorType = "读写错误"
	ErrorUnknown         ErrorType = "未知错误"
)

// UserError 用户友好的错误包装
type UserError struct {
	Type    ErrorType
	Message string
	Inner   error
}

func (e *UserError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %s (细节: %v)", e.Type, e.Message, e.Inner)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *UserError) Unwrap() error {
	return e.Inner
}

// NewUserError 创建用户友好错误
func NewUserError(errType ErrorType, message string, inner error) error {
	return &UserError{
		Type:    errType,
		Message: message,
		Inner:   inner,
	}
}

// WrapError 包装原始错误为用户友好错误
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}

	// 如果已经是 UserError，直接返回
	if _, ok := err.(*UserError); ok {
		return err
	}

	// 根据错误内容推断错误类型
	errStr := strings.ToLower(err.Error())
	var errType ErrorType

	switch {
	case strings.Contains(errStr, "permission") || strings.Contains(errStr, "denied"):
		errType = ErrorPermission
	case strings.Contains(errStr, "not found") || strings.Contains(errStr, "no such file"):
		errType = ErrorFileNotFound
	case strings.Contains(errStr, "exists") || strings.Contains(errStr, "already"):
		errType = ErrorFileExists
	case strings.Contains(errStr, "key") || strings.Contains(errStr, "decrypt") || strings.Contains(errStr, "encrypt"):
		if strings.Contains(errStr, "invalid") || strings.Contains(errStr, "format") {
			errType = ErrorKeyInvalid
		} else if strings.Contains(errStr, "decrypt") {
			errType = ErrorDecryption
		} else if strings.Contains(errStr, "encrypt") {
			errType = ErrorEncryption
		} else if strings.Contains(errStr, "signature") {
			errType = ErrorSignature
		} else {
			errType = ErrorKeyNotFound
		}
	case strings.Contains(errStr, "read") || strings.Contains(errStr, "write") || strings.Contains(errStr, "io"):
		errType = ErrorIO
	default:
		errType = ErrorUnknown
	}

	return &UserError{
		Type:    errType,
		Message: context,
		Inner:   err,
	}
}

// PrintError 打印用户友好的错误信息
func PrintError(err error) {
	if err == nil {
		return
	}

	// 检查是否是 UserError
	if userErr, ok := err.(*UserError); ok {
		fmt.Printf("❌ %s\n", userErr.Message)
		if userErr.Inner != nil && userErr.Type != ErrorUnknown {
			// 在 verbose 模式下显示内部错误细节
			fmt.Printf("   详情: %v\n", userErr.Inner)
		}
	} else {
		// 普通错误
		fmt.Printf("❌ 错误: %v\n", err)
	}
}

// PrintWarning 打印警告信息
func PrintWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// PrintSuccess 打印成功信息
func PrintSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// PrintInfo 打印信息
func PrintInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

// ConfirmPrompt 确认提示
func ConfirmPrompt(message string, defaultYes bool) bool {
	var response string

	if defaultYes {
		fmt.Printf("%s [Y/n]: ", message)
	} else {
		fmt.Printf("%s [y/N]: ", message)
	}

	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}

// SelectPrompt 选择提示
func SelectPrompt(message string, options []string, defaultIndex int) int {
	fmt.Printf("%s:\n", message)
	for i, opt := range options {
		marker := "  "
		if i == defaultIndex {
			marker = "→ "
		}
		fmt.Printf("  %d. %s%s\n", i+1, marker, opt)
	}

	var choice int
	fmt.Printf("\n请选择 [1-%d]: ", len(options))
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(options) {
		return defaultIndex
	}

	return choice - 1
}
