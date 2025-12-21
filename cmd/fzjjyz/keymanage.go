package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	keymanageAction    string
	keymanagePubKey    string
	keymanagePrivKey   string
	keymanageOutput    string
	keymanageOutputDir string
)

func newKeymanageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keymanage",
		Short: "密钥管理工具",
		Long: `管理加密密钥，支持导出、导入和验证操作。

可用操作:
  export    从私钥文件中提取并导出公钥
  import    导入密钥文件到指定目录
  verify    验证密钥对是否匹配

示例:
  # 导出公钥
  fzjjyz keymanage export --private-key private.pem --output public_extracted.pem

  # 验证密钥对
  fzjjyz keymanage verify --public-key public.pem --private-key private.pem

  # 导入密钥
  fzjjyz keymanage import --public-key pub.pem --private-key priv.pem --output-dir ./keys`,
		RunE: runKeymanage,
	}

	cmd.Flags().StringVarP(&keymanageAction, "action", "a", "", "操作类型: export/import/verify (必需)")
	cmd.Flags().StringVarP(&keymanagePubKey, "public-key", "p", "", "公钥文件路径")
	cmd.Flags().StringVarP(&keymanagePrivKey, "private-key", "s", "", "私钥文件路径")
	cmd.Flags().StringVarP(&keymanageOutput, "output", "o", "", "输出文件路径 (用于export)")
	cmd.Flags().StringVarP(&keymanageOutputDir, "output-dir", "d", ".", "输出目录 (用于import)")

	cmd.MarkFlagRequired("action")

	return cmd
}

func runKeymanage(cmd *cobra.Command, args []string) error {
	switch keymanageAction {
	case "export":
		return runExport()
	case "import":
		return runImport()
	case "verify":
		return runVerify()
	default:
		return fmt.Errorf("未知操作: %s (支持: export, import, verify)", keymanageAction)
	}
}

// export: 从私钥文件中提取并导出公钥
func runExport() error {
	if keymanagePrivKey == "" {
		return fmt.Errorf("必须提供 --private-key")
	}
	if keymanageOutput == "" {
		return fmt.Errorf("必须提供 --output")
	}

	fmt.Println("导出公钥...")

	// 加载私钥文件
	hybridPriv, err := crypto.LoadPrivateKey(keymanagePrivKey)
	if err != nil {
		return fmt.Errorf("加载私钥失败: %v", err)
	}

	// 从私钥中提取公钥
	kyberPub := hybridPriv.Kyber.Public()
	ecdhPub := hybridPriv.ECDH.PublicKey()

	// 导出公钥到 PEM
	pubPEM, err := crypto.ExportPublicKey(kyberPub, ecdhPub)
	if err != nil {
		return fmt.Errorf("导出公钥失败: %v", err)
	}

	// 保存公钥文件
	if err := os.WriteFile(keymanageOutput, pubPEM, 0644); err != nil {
		return fmt.Errorf("保存公钥文件失败: %v", err)
	}

	fmt.Printf("✅ 公钥已导出到: %s\n", keymanageOutput)
	return nil
}

// import: 导入密钥到指定目录
func runImport() error {
	if keymanagePubKey == "" || keymanagePrivKey == "" {
		return fmt.Errorf("必须提供 --public-key 和 --private-key")
	}

	if keymanageOutputDir == "" {
		keymanageOutputDir = "."
	}

	fmt.Println("导入密钥...")

	// 创建输出目录
	if err := os.MkdirAll(keymanageOutputDir, 0755); err != nil {
		return fmt.Errorf("无法创建目录: %v", err)
	}

	// 加载密钥
	hybridPub, err := crypto.LoadPublicKey(keymanagePubKey)
	if err != nil {
		return fmt.Errorf("加载公钥失败: %v", err)
	}

	hybridPriv, err := crypto.LoadPrivateKey(keymanagePrivKey)
	if err != nil {
		return fmt.Errorf("加载私钥失败: %v", err)
	}

	// 生成新路径
	basePub := filepath.Base(keymanagePubKey)
	basePriv := filepath.Base(keymanagePrivKey)
	newPubPath := filepath.Join(keymanageOutputDir, basePub)
	newPrivPath := filepath.Join(keymanageOutputDir, basePriv)

	// 保存到新位置
	if err := crypto.SaveKeyFiles(hybridPub.Kyber, hybridPub.ECDH, hybridPriv.Kyber, hybridPriv.ECDH, newPubPath, newPrivPath); err != nil {
		return fmt.Errorf("保存密钥失败: %v", err)
	}

	fmt.Printf("✅ 密钥已导入到: %s\n", keymanageOutputDir)
	fmt.Printf("  公钥: %s\n", basePub)
	fmt.Printf("  私钥: %s\n", basePriv)

	return nil
}

// verify: 验证密钥对是否匹配
func runVerify() error {
	if keymanagePubKey == "" || keymanagePrivKey == "" {
		return fmt.Errorf("必须提供 --public-key 和 --private-key")
	}

	fmt.Println("验证密钥对...")

	// 加载公钥
	hybridPub, err := crypto.LoadPublicKey(keymanagePubKey)
	if err != nil {
		fmt.Println("❌ 公钥加载失败")
		return fmt.Errorf("加载公钥失败: %v", err)
	}

	// 加载私钥
	hybridPriv, err := crypto.LoadPrivateKey(keymanagePrivKey)
	if err != nil {
		fmt.Println("❌ 私钥加载失败")
		return fmt.Errorf("加载私钥失败: %v", err)
	}

	// 验证密钥对是否匹配
	// 比较公钥字节
	kyberPubBytes, _ := hybridPub.Kyber.MarshalBinary()
	kyberPrivPubBytes, _ := hybridPriv.Kyber.Public().MarshalBinary()

	ecdhPubBytes := hybridPub.ECDH.Bytes()
	ecdhPrivPubBytes := hybridPriv.ECDH.PublicKey().Bytes()

	kyberMatch := bytes.Equal(kyberPubBytes, kyberPrivPubBytes)
	ecdhMatch := bytes.Equal(ecdhPubBytes, ecdhPrivPubBytes)

	if kyberMatch && ecdhMatch {
		fmt.Println("✅ 密钥对验证通过")
		fmt.Println("  Kyber:  ✅ 匹配")
		fmt.Println("  ECDH:   ✅ 匹配")
	} else {
		fmt.Println("❌ 密钥对不匹配")
		if !kyberMatch {
			fmt.Println("  Kyber:  ❌ 不匹配")
		}
		if !ecdhMatch {
			fmt.Println("  ECDH:   ❌ 不匹配")
		}
		return fmt.Errorf("密钥对不匹配")
	}

	return nil
}

