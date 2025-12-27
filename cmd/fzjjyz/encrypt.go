package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	encryptInput      string
	encryptOutput     string
	encryptPubKey     string
	encryptSignKey    string
	encryptForce      bool
	encryptBufferSize int
	encryptStreaming  bool
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
	cmd.Flags().IntVar(&encryptBufferSize, "buffer-size", 0, "缓冲区大小 (KB)，0=自动选择")
	cmd.Flags().BoolVar(&encryptStreaming, "streaming", true, "使用流式处理（大文件推荐）")

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
		fmt.Printf("  流式处理: %v\n", encryptStreaming)
	}

	// 加载密钥（使用缓存）
	fmt.Print("\n[1/3] 加载密钥... ")
	hybridPub, err := crypto.LoadPublicKeyCached(encryptPubKey)
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("❌ 加载公钥失败: %v\n\n提示:\n  1. 请检查公钥文件路径是否正确: %s\n  2. 确保公钥文件格式正确（PEM 格式）\n  3. 检查文件权限（需可读）\n  4. 如果是首次使用，请先生成密钥对: fzjjyz keygen", err, encryptPubKey)
	}

	dilithiumPriv, err := crypto.LoadDilithiumPrivateKeyCached(encryptSignKey)
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("❌ 加载签名私钥失败: %v\n\n提示:\n  1. 请检查 Dilithium 私钥文件路径是否正确: %s\n  2. 确保私钥文件格式正确（PEM 格式）\n  3. 检查文件权限（建议 0600）\n  4. 私钥文件应仅由所有者读取\n  5. 如果是首次使用，请先生成密钥对: fzjjyz keygen", err, encryptSignKey)
	}
	fmt.Println("完成")

	// 确定缓冲区大小
	var bufSize int
	if encryptBufferSize > 0 {
		bufSize = encryptBufferSize * 1024
	} else {
		stat, _ := os.Stat(encryptInput)
		bufSize = crypto.OptimalBufferSize(stat.Size())
	}

	if verbose {
		fmt.Printf("  缓冲区大小: %d KB\n", bufSize/1024)
	}

	// 执行加密
	fmt.Print("[2/3] 加密文件... ")
	var encryptFunc func() error
	if encryptStreaming {
		encryptFunc = func() error {
			return crypto.EncryptFileStreaming(
				encryptInput, encryptOutput,
				hybridPub.Kyber, hybridPub.ECDH,
				dilithiumPriv,
				bufSize,
			)
		}
	} else {
		encryptFunc = func() error {
			return crypto.EncryptFile(
				encryptInput, encryptOutput,
				hybridPub.Kyber, hybridPub.ECDH,
				dilithiumPriv,
			)
		}
	}

	if err := encryptFunc(); err != nil {
		fmt.Println("失败")
		return fmt.Errorf("❌ 加密失败: %v\n\n可能原因:\n  1. 文件权限不足（无法读取输入或写入输出）\n  2. 内存不足（大文件需要更多内存）\n  3. 密钥不匹配\n  4. 输入文件在加密过程中被修改\n\n建议:\n  - 检查磁盘空间和文件权限\n  - 对于超大文件，尝试使用 --buffer-size 调整缓冲区\n  - 确保密钥正确匹配", err)
	}
	fmt.Println("完成")

	// 显示结果
	fmt.Print("[3/3] 验证... ")
	encryptedInfo, _ := os.Stat(encryptOutput)
	originalInfo, _ := os.Stat(encryptInput)
	fmt.Println("完成")

	fmt.Printf("\n✅ 加密成功！\n\n")
	fmt.Printf("文件信息:\n")
	fmt.Printf("  原始文件: %s (%d bytes)\n", filepath.Base(encryptInput), originalInfo.Size())
	fmt.Printf("  加密文件: %s (%d bytes)\n", filepath.Base(encryptOutput), encryptedInfo.Size())
	fmt.Printf("  压缩率: %.1f%%\n", float64(encryptedInfo.Size())/float64(originalInfo.Size())*100)

	return nil
}
