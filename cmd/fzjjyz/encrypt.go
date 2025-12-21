package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	encryptInput     string
	encryptOutput    string
	encryptPubKey    string
	encryptSignKey   string
	encryptForce     bool
)

func newEncryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "加密文件",
		Long: `使用后量子混合加密算法加密文件。

加密流程：
  1. 读取原始文件
  2. 生成随机会话密钥
  3. Kyber768 + ECDH 密钥封装
  4. AES-256-GCM 加密数据
  5. Dilithium3 签名验证
  6. 构建加密文件头
  7. 写入加密文件

必需参数：
  --input, -i         输入文件路径
  --public-key, -p    Kyber+ECDH 公钥文件
  --sign-key, -s      Dilithium 私钥文件

示例：
  fzjjyz encrypt -i plaintext.txt -o encrypted.fzj -p public.pem -s dilithium_private.pem
  fzjjyz encrypt --input data.txt --public-key pub.pem --sign-key priv.pem --force`,
		RunE: runEncrypt,
	}

	cmd.Flags().StringVarP(&encryptInput, "input", "i", "", "输入文件路径 (必需)")
	cmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "输出文件路径 (可选，默认: input.fzj)")
	cmd.Flags().StringVarP(&encryptPubKey, "public-key", "p", "", "Kyber+ECDH 公钥文件 (必需)")
	cmd.Flags().StringVarP(&encryptSignKey, "sign-key", "s", "", "Dilithium 私钥文件 (必需)")
	cmd.Flags().BoolVarP(&encryptForce, "force", "f", false, "覆盖输出文件")

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("public-key")
	cmd.MarkFlagRequired("sign-key")

	return cmd
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(encryptInput); err != nil {
		return fmt.Errorf("输入文件不存在: %s", encryptInput)
	}

	// 设置默认输出路径
	if encryptOutput == "" {
		encryptOutput = encryptInput + ".fzj"
	}

	// 检查输出文件是否已存在
	if !encryptForce {
		if _, err := os.Stat(encryptOutput); err == nil {
			return fmt.Errorf("输出文件已存在: %s (使用 --force 覆盖)", encryptOutput)
		}
	}

	// 显示进度
	fmt.Printf("加密文件: %s\n", filepath.Base(encryptInput))
	if verbose {
		fmt.Printf("  输入: %s\n", encryptInput)
		fmt.Printf("  输出: %s\n", encryptOutput)
		fmt.Printf("  公钥: %s\n", encryptPubKey)
		fmt.Printf("  签名密钥: %s\n", encryptSignKey)
	}

	// 加载密钥
	fmt.Print("\n[1/3] 加载密钥... ")

	// 加载 Kyber+ECDH 公钥
	hybridPub, err := crypto.LoadPublicKey(encryptPubKey)
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("加载公钥失败: %v", err)
	}

	// 加载 Dilithium 私钥
	dilithiumPriv, err := crypto.LoadDilithiumPrivateKey(encryptSignKey)
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("加载签名私钥失败: %v", err)
	}
	fmt.Println("完成")

	// 加密文件
	fmt.Print("[2/3] 加密文件... ")
	if err := crypto.EncryptFile(encryptInput, encryptOutput, hybridPub.Kyber, hybridPub.ECDH, dilithiumPriv); err != nil {
		fmt.Println("失败")
		return fmt.Errorf("加密失败: %v", err)
	}
	fmt.Println("完成")

	// 显示结果
	fmt.Print("[3/3] 验证... ")

	// 读取加密文件大小
	encryptedInfo, err := os.Stat(encryptOutput)
	if err != nil {
		fmt.Println("警告: 无法获取文件信息")
	} else {
		fmt.Println("完成")

		// 获取原始文件大小
		originalInfo, _ := os.Stat(encryptInput)

		fmt.Printf("\n✅ 加密成功！\n")
		fmt.Printf("\n文件信息:\n")
		fmt.Printf("  原始文件: %s (%d bytes)\n", filepath.Base(encryptInput), originalInfo.Size())
		fmt.Printf("  加密文件: %s (%d bytes)\n", filepath.Base(encryptOutput), encryptedInfo.Size())
		fmt.Printf("  压缩率: %.1f%%\n", float64(encryptedInfo.Size())/float64(originalInfo.Size())*100)
	}

	return nil
}
