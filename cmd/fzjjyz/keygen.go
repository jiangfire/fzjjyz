package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	keygenOutputDir string
	keygenName      string
	keygenForce     bool
)

func newKeygenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "生成后量子密钥对",
		Long: `生成完整的密钥对组合，包括：
  • Kyber768 + ECDH 密钥对 (用于加密/解密)
  • Dilithium3 密钥对 (用于签名/验证)

生成的文件：
  {name}_public.pem          - Kyber+ECDH 公钥
  {name}_private.pem         - Kyber+ECDH 私钥 (0600权限)
  {name}_dilithium_public.pem  - Dilithium 公钥
  {name}_dilithium_private.pem - Dilithium 私钥 (0600权限)

示例：
  fzjjyz keygen -d ./keys -n mykey
  fzjjyz keygen --output-dir ./keys --name mykey --force`,
		RunE: runKeygen,
	}

	cmd.Flags().StringVarP(&keygenOutputDir, "output-dir", "d", ".", "输出目录")
	cmd.Flags().StringVarP(&keygenName, "name", "n", "", "密钥名称前缀 (默认: 时间戳)")
	cmd.Flags().BoolVarP(&keygenForce, "force", "f", false, "覆盖现有文件")

	return cmd
}

func runKeygen(cmd *cobra.Command, args []string) error {
	// 生成默认名称
	if keygenName == "" {
		keygenName = fmt.Sprintf("key_%d", time.Now().Unix())
	}

	// 创建输出目录
	if err := os.MkdirAll(keygenOutputDir, 0755); err != nil {
		return fmt.Errorf("无法创建目录 %s: %v", keygenOutputDir, err)
	}

	// 构建文件路径
	pubPath := filepath.Join(keygenOutputDir, keygenName+"_public.pem")
	privPath := filepath.Join(keygenOutputDir, keygenName+"_private.pem")
	dilithiumPubPath := filepath.Join(keygenOutputDir, keygenName+"_dilithium_public.pem")
	dilithiumPrivPath := filepath.Join(keygenOutputDir, keygenName+"_dilithium_private.pem")

	// 检查文件是否已存在
	paths := []string{pubPath, privPath, dilithiumPubPath, dilithiumPrivPath}
	if !keygenForce {
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return fmt.Errorf("文件已存在: %s (使用 --force 覆盖)", p)
			}
		}
	}

	// 显示生成信息
	if verbose {
		fmt.Printf("生成密钥对...\n")
		fmt.Printf("  输出目录: %s\n", keygenOutputDir)
		fmt.Printf("  密钥名称: %s\n", keygenName)
	}

	// 1. 生成 Kyber 密钥对
	fmt.Print("  [1/4] 生成 Kyber768 密钥... ")
	kyberPub, kyberPriv, err := crypto.GenerateKyberKeys()
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("Kyber密钥生成失败: %v", err)
	}
	fmt.Println("完成")

	// 2. 生成 ECDH 密钥对
	fmt.Print("  [2/4] 生成 ECDH X25519 密钥... ")
	ecdhPub, ecdhPriv, err := crypto.GenerateECDHKeys()
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("ECDH密钥生成失败: %v", err)
	}
	fmt.Println("完成")

	// 3. 生成 Dilithium 密钥对
	fmt.Print("  [3/4] 生成 Dilithium3 签名密钥... ")
	dilithiumPub, dilithiumPriv, err := crypto.GenerateDilithiumKeys()
	if err != nil {
		fmt.Println("失败")
		return fmt.Errorf("Dilithium密钥生成失败: %v", err)
	}
	fmt.Println("完成")

	// 4. 保存密钥文件
	fmt.Print("  [4/4] 保存密钥文件... ")

	// 保存 Kyber+ECDH 密钥对
	if err := crypto.SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		fmt.Println("失败")
		return fmt.Errorf("保存密钥文件失败: %v", err)
	}

	// 保存 Dilithium 密钥对
	if err := crypto.SaveDilithiumKeys(dilithiumPub, dilithiumPriv, dilithiumPubPath, dilithiumPrivPath); err != nil {
		fmt.Println("失败")
		return fmt.Errorf("保存Dilithium密钥失败: %v", err)
	}

	fmt.Println("完成")

	// 显示摘要
	fmt.Println("\n✅ 密钥对生成成功！")
	fmt.Println("\n生成的文件:")
	fmt.Printf("  • %s (公钥)\n", filepath.Base(pubPath))
	fmt.Printf("  • %s (私钥 - 0600权限)\n", filepath.Base(privPath))
	fmt.Printf("  • %s (签名公钥)\n", filepath.Base(dilithiumPubPath))
	fmt.Printf("  • %s (签名私钥 - 0600权限)\n", filepath.Base(dilithiumPrivPath))

	fmt.Println("\n⚠️  安全提示:")
	fmt.Println("  • 请妥善保管私钥文件")
	fmt.Println("  • 不要将私钥分享给他人")
	fmt.Println("  • 建议使用安全的存储介质")

	return nil
}
