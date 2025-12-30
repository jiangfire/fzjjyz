package i18n

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// TranslateCommand 自动翻译 Cobra 命令的所有文本
func TranslateCommand(cmd *cobra.Command, keyPrefix string) {
	if cmd == nil {
		return
	}

	// 翻译命令描述
	if cmd.Short != "" {
		cmd.Short = Get(keyPrefix + ".short")
	}

	if cmd.Long != "" {
		cmd.Long = Get(keyPrefix + ".long")
	}

	// 翻译标志描述
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Usage != "" {
			// 构建标志的翻译键：keyPrefix.flags.flagName
			flagKey := keyPrefix + ".flags." + flag.Name
			translated := Get(flagKey)
			if translated != "" {
				flag.Usage = translated
			}
		}
	})

	// 递归处理子命令
	for _, subCmd := range cmd.Commands() {
		subKeyPrefix := keyPrefix + "." + subCmd.Name()
		TranslateCommand(subCmd, subKeyPrefix)
	}
}

// TranslateError 创建翻译后的错误
func TranslateError(key string, args ...interface{}) error {
	return &TranslatedError{
		key:  key,
		args: args,
	}
}

// TranslatedError 翻译后的错误类型
type TranslatedError struct {
	key  string
	args []interface{}
}

func (e *TranslatedError) Error() string {
	return T(e.key, e.args...)
}

// MustTranslate 强制翻译（用于测试或确保翻译存在）
func MustTranslate(key string, args ...interface{}) string {
	result := T(key, args...)
	if result == "" {
		return key
	}
	return result
}
