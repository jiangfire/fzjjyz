package main

import (
	"os"

	"github.com/spf13/cobra"
)

// 版本信息
const (
	Version     = "0.1.0"
	AppName     = "fzjjyz"
	Description = "后量子文件加密工具 - 使用 Kyber768 + ECDH + AES-256-GCM + Dilithium3"
)

// 根命令
var rootCmd = &cobra.Command{
	Use:   "fzjjyz",
	Short: Description,
	Long: `fzjjyz - 后量子文件加密工具

使用以下算法提供安全的文件加密：
  • Kyber768 - 后量子密钥封装
  • X25519 ECDH - 传统密钥交换
  • AES-256-GCM - 认证加密
  • Dilithium3 - 数字签名

快速开始：
  # 生成密钥对
  fzjjyz keygen -d ./keys -n mykey

  # 加密文件
  fzjjyz encrypt -i plaintext.txt -o encrypted.fzj -p keys/mykey_public.pem -s keys/mykey_dilithium_private.pem

  # 解密文件
  fzjjyz decrypt -i encrypted.fzj -o decrypted.txt -p keys/mykey_private.pem -s keys/mykey_dilithium_public.pem

  # 查看文件信息
  fzjjyz info -i encrypted.fzj

项目主页: https://codeberg.org/jiangfire/fzjjyz`,
	Version: Version,
}

// 全局标志
var (
	verbose bool
	force   bool
)

func init() {
	// 添加全局标志
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "启用详细输出")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "强制覆盖现有文件")

	// 添加子命令
	rootCmd.AddCommand(
		newEncryptCmd(),
		newDecryptCmd(),
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
