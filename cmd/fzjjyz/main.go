package main

import (
	"os"

	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"github.com/spf13/cobra"
)

// 版本信息
const (
	Version     = "0.1.0"
	AppName     = "fzjjyz"
	Description = "后量子文件加密工具 - 使用 Kyber768 + ECDH + AES-256-GCM + Dilithium3"
)

// 根命令（文本将在 init 中通过 i18n 翻译）
var rootCmd *cobra.Command

// 全局标志
var (
	verbose bool
	force   bool
)

func init() {
	// 初始化国际化（从环境变量 LANG 读取）
	if err := i18n.Init(""); err != nil {
		// 如果初始化失败，使用默认语言（中文）
		i18n.Init("zh_CN")
	}

	// 创建根命令
	rootCmd = &cobra.Command{
		Use:     "fzjjyz",
		Short:   i18n.T("app.description"),
		Long:    i18n.T("app.long"),
		Version: Version,
	}

	// 添加全局标志
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, i18n.T("flags.verbose"))
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, i18n.T("flags.force"))

	// 添加子命令
	rootCmd.AddCommand(
		newEncryptCmd(),
		newDecryptCmd(),
		newEncryptDirCmd(),
		newDecryptDirCmd(),
		newKeygenCmd(),
		newKeymanageCmd(),
		newInfoCmd(),
		newVersionCmd(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
