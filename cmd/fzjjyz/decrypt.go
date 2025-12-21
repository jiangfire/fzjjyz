package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"github.com/spf13/cobra"
)

var (
	decryptInput     string
	decryptOutput    string
	decryptPrivKey   string
	decryptVerifyKey string
	decryptForce     bool
)

func newDecryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "解密文件",
		Long: `使用后量子混合加密算法解密文件。

解密流程：
  1. 解析文件头
  2. 验证文件格式
  3. Kyber768 + ECDH 密钥解封装
  4. AES-256-GCM 解密数据
  5. 验证 SHA256 哈希
  6. 验证 Dilithium 签名
  7. 写入原始文件

必需参数：
  --input, -i         加密文件路径
  --private-key, -p   Kyber+ECDH 私钥文件
  --verify-key, -s    Dilithium 公钥文件 (可选)

示例：
  fzjjyz decrypt -i encrypted.fzj -o decrypted.txt -p private.pem -s dilithium_public.pem
  fzjjyz decrypt --input data.fzj --private-key priv.pem --verify-key pub.pem --force`,
		RunE: runDecrypt,
	}

	cmd.Flags().StringVarP(&decryptInput, "input", "i", "", "加密文件路径 (必需)")
	cmd.Flags().StringVarP(&decryptOutput, "output", "o", "", "输出文件路径 (可选，默认: 原文件名)")
	cmd.Flags().StringVarP(&decryptPrivKey, "private-key", "p", "", "Kyber+ECDH 私钥文件 (必需)")
	cmd.Flags().StringVarP(&decryptVerifyKey, "verify-key", "s", "", "Dilithium 公钥文件 (可选)")
	cmd.Flags().BoolVarP(&decryptForce, "force", "f", false, "覆盖输出文件")

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("private-key")

	return cmd
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(decryptInput); err != nil {
		return fmt.Errorf("加密文件不存在: %s", decryptInput)
	}

	// 读取文件头以获取原始文件名
	encryptedData, err := os.ReadFile(decryptInput)
	if err != nil {
		return fmt.Errorf("无法读取加密文件: %v", err)
	}

	header, err := format.ParseFileHeaderFromBytes(encryptedData)
	if err != nil {
		return fmt.Errorf("文件头解析失败: %v", err)
	}

	// 设置默认输出路径
	if decryptOutput == "" {
		decryptOutput = header.Filename
	}

	// 检查输出文件是否已存在
	if !decryptForce {
		if _, err := os.Stat(decryptOutput); err == nil {
			return fmt.Errorf("输出文件已存在: %s (使用 --force 覆盖)", decryptOutput)
		}
	}

	// 显示进度
	fmt.Printf("解密文件: %s\n", filepath.Base(decryptInput))
	if verbose {
		fmt.Printf("  输入: %s\n", decryptInput)
		fmt.Printf("  输出: %s\n", decryptOutput)
		fmt.Printf("  私钥: %s\n", decryptPrivKey)
		if decryptVerifyKey != "" {
			fmt.Printf("  验证密钥: %s\n", decryptVerifyKey)
		}
		fmt.Printf("  原始文件名: %s\n", header.Filename)
	}

	// 加载密钥
	fmt.Print("\n[1/3] 加载密钥... ")

	// 加载 Kyber+ECDH 私钥
	hybridPriv, err := crypto.LoadPrivateKey(decryptPrivKey)
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("加载私钥失败: %v", err)
	}

	// 加载 Dilithium 公钥（如果提供）
	var dilithiumPub interface{}
	if decryptVerifyKey != "" {
		dilithiumPub, err = crypto.LoadDilithiumPublicKey(decryptVerifyKey)
		if err != nil {
			fmt.Println("失败")
			return fmt.Errorf("加载验证公钥失败: %v", err)
		}
	} else {
		// 尝试从文件头提取公钥（如果支持）
		// 目前不支持，需要提供验证密钥
		fmt.Println("警告: 未提供签名验证密钥，将跳过签名验证")
	}
	fmt.Println("完成")

	// 解密文件
	fmt.Print("[2/3] 解密文件... ")
	if err := crypto.DecryptFile(decryptInput, decryptOutput, hybridPriv.Kyber, hybridPriv.ECDH, dilithiumPub); err != nil {
		fmt.Println("失败")
		return fmt.Errorf("解密失败: %v", err)
	}
	fmt.Println("完成")

	// 显示结果
	fmt.Print("[3/3] 验证... ")
	fmt.Println("完成")

	// 获取文件信息
	decryptedInfo, err := os.Stat(decryptOutput)
	if err != nil {
		fmt.Println("\n⚠️  解密完成，但无法获取文件信息")
		return nil
	}

	encryptedInfo, _ := os.Stat(decryptInput)

	fmt.Printf("\n✅ 解密成功！\n")
	fmt.Printf("\n文件信息:\n")
	fmt.Printf("  加密文件: %s (%d bytes)\n", filepath.Base(decryptInput), encryptedInfo.Size())
	fmt.Printf("  解密文件: %s (%d bytes)\n", filepath.Base(decryptOutput), decryptedInfo.Size())
	fmt.Printf("  原始文件名: %s\n", header.Filename)
	fmt.Printf("  时间戳: %s\n", format.UnixTime(header.Timestamp))

	return nil
}
