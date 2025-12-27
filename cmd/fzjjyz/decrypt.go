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
	decryptInput      string
	decryptOutput     string
	decryptPrivKey    string
	decryptVerifyKey  string
	decryptForce      bool
	decryptBufferSize int
	decryptStreaming  bool
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
	cmd.Flags().IntVar(&decryptBufferSize, "buffer-size", 0, "缓冲区大小 (KB)，0=自动选择")
	cmd.Flags().BoolVar(&decryptStreaming, "streaming", true, "使用流式处理（大文件推荐）")

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("private-key")

	return cmd
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(decryptInput); err != nil {
		return fmt.Errorf("加密文件不存在: %s", decryptInput)
	}

	// 读取文件头以获取原始文件名（使用缓存读取，避免大文件问题）
	// 对于流式解密，我们只需要读取头部
	headerFile, err := os.Open(decryptInput)
	if err != nil {
		return fmt.Errorf("无法打开加密文件: %v", err)
	}
	defer headerFile.Close()

	header, err := format.ParseFileHeader(headerFile)
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
		fmt.Printf("  流式处理: %v\n", decryptStreaming)
	}

	// 加载密钥（使用缓存）
	fmt.Print("\n[1/3] 加载密钥... ")
	hybridPriv, err := crypto.LoadPrivateKeyCached(decryptPrivKey)
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("❌ 加载私钥失败: %v\n\n提示:\n  1. 请检查私钥文件路径是否正确: %s\n  2. 确保私钥文件格式正确（PEM 格式）\n  3. 检查文件权限（建议 0600）\n  4. 私钥文件应仅由所有者读取\n  5. 确保使用与加密时匹配的私钥", err, decryptPrivKey)
	}

	var dilithiumPub interface{}
	if decryptVerifyKey != "" {
		dilithiumPub, err = crypto.LoadDilithiumPublicKeyCached(decryptVerifyKey)
		if err != nil {
			fmt.Println("失败")
			return fmt.Errorf("❌ 加载验证公钥失败: %v\n\n提示:\n  1. 请检查 Dilithium 公钥文件路径是否正确: %s\n  2. 确保公钥文件格式正确（PEM 格式）\n  3. 检查文件权限（需可读）\n  4. 确保使用与加密时匹配的公钥\n  5. 如果未提供签名密钥，可省略此参数（但无法验证签名）", err, decryptVerifyKey)
		}
	} else {
		fmt.Println("⚠️  警告: 未提供签名验证密钥，将跳过签名验证")
	}
	fmt.Println("完成")

	// 确定缓冲区大小
	var bufSize int
	if decryptBufferSize > 0 {
		bufSize = decryptBufferSize * 1024
	} else {
		stat, _ := os.Stat(decryptInput)
		bufSize = crypto.OptimalBufferSize(stat.Size())
	}

	if verbose {
		fmt.Printf("  缓冲区大小: %d KB\n", bufSize/1024)
	}

	// 执行解密
	fmt.Print("[2/3] 解密文件... ")
	var decryptFunc func() error
	if decryptStreaming {
		decryptFunc = func() error {
			return crypto.DecryptFileStreaming(
				decryptInput, decryptOutput,
				hybridPriv.Kyber, hybridPriv.ECDH,
				dilithiumPub,
				bufSize,
			)
		}
	} else {
		decryptFunc = func() error {
			return crypto.DecryptFile(
				decryptInput, decryptOutput,
				hybridPriv.Kyber, hybridPriv.ECDH,
				dilithiumPub,
			)
		}
	}

	if err := decryptFunc(); err != nil {
		fmt.Println("失败")
		return fmt.Errorf("❌ 解密失败: %v\n\n可能原因:\n  1. 密钥不匹配（使用了错误的私钥）\n  2. 文件已损坏或被篡改\n  3. 文件格式不正确（不是 fzjjyz 加密文件）\n  4. 签名验证失败（文件可能被篡改）\n  5. 文件权限不足\n\n安全提示:\n  - 如果提示哈希不匹配，文件可能已被篡改，请勿使用\n  - 如果提示签名无效，密钥可能不匹配或文件被修改\n  - 建议始终提供签名验证密钥以确保文件完整性", err)
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
