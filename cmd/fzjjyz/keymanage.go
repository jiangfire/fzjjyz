// Package main 提供文件加密解密命令行工具.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"codeberg.org/jiangfire/fzjjyz/internal/zjcrypto"
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
		Short: i18n.T("keymanage.short"),
		Long:  i18n.T("keymanage.long"),
		RunE:  runKeymanage,
	}

	cmd.Flags().StringVarP(&keymanageAction, "action", "a", "", i18n.T("keymanage.flags.action"))
	cmd.Flags().StringVarP(&keymanagePubKey, "public-key", "p", "", i18n.T("keymanage.flags.public-key"))
	cmd.Flags().StringVarP(&keymanagePrivKey, "private-key", "s", "", i18n.T("keymanage.flags.private-key"))
	cmd.Flags().StringVarP(&keymanageOutput, "output", "o", "", i18n.T("keymanage.flags.output"))
	cmd.Flags().StringVarP(&keymanageOutputDir, "output-dir", "d", ".", i18n.T("keymanage.flags.output-dir"))

	_ = cmd.MarkFlagRequired("action")

	return cmd
}

func runKeymanage(_ *cobra.Command, _ []string) error {
	switch keymanageAction {
	case "export":
		return runExport()
	case "import":
		return runImport()
	case "verify":
		return runVerify()
	default:
		return fmt.Errorf(i18n.T("error.unknown_action"), keymanageAction)
	}
}

// export: 从私钥文件中提取并导出公钥
func runExport() error {
	if keymanagePrivKey == "" {
		return fmt.Errorf(i18n.T("error.missing_required_flags"), "--private-key")
	}
	if keymanageOutput == "" {
		return fmt.Errorf(i18n.T("error.missing_required_flags"), "--output")
	}

	fmt.Println(i18n.T("status.public_key") + "...")

	// 加载私钥文件
	hybridPriv, err := zjcrypto.LoadPrivateKey(keymanagePrivKey)
	if err != nil {
		return fmt.Errorf("load private key failed: %w",
			i18n.TranslateError("error.load_private_key_failed", err, keymanagePrivKey))
	}

	// 从私钥中提取公钥
	kyberPub := hybridPriv.Kyber.Public()
	ecdhPub := hybridPriv.ECDH.PublicKey()

	// 导出公钥到 PEM
	pubPEM, err := zjcrypto.ExportPublicKey(kyberPub, ecdhPub)
	if err != nil {
		return fmt.Errorf("export key failed: %w",
			i18n.TranslateError("error.export_key_failed", err))
	}

	// 保存公钥文件
	if err := os.WriteFile(keymanageOutput, pubPEM, 0600); err != nil {
		return fmt.Errorf("save export failed: %w",
			i18n.TranslateError("error.save_export_failed", err))
	}

	fmt.Printf(i18n.T("status.success_export")+"\n", keymanageOutput)
	return nil
}

// import: 导入密钥到指定目录
func runImport() error {
	if keymanagePubKey == "" || keymanagePrivKey == "" {
		return fmt.Errorf("%s", i18n.T("error.missing_both_keys"))
	}

	if keymanageOutputDir == "" {
		keymanageOutputDir = "."
	}

	fmt.Println(i18n.T("status.success_import") + "...")

	// 创建输出目录
	if err := os.MkdirAll(keymanageOutputDir, 0750); err != nil {
		return fmt.Errorf(i18n.T("error.cannot_create_dir"), keymanageOutputDir, err)
	}

	// 加载密钥
	hybridPub, err := zjcrypto.LoadPublicKey(keymanagePubKey)
	if err != nil {
		return fmt.Errorf("load public key failed: %w",
			i18n.TranslateError("error.load_public_key_failed", err, keymanagePubKey))
	}

	hybridPriv, err := zjcrypto.LoadPrivateKey(keymanagePrivKey)
	if err != nil {
		return fmt.Errorf("load private key failed: %w",
			i18n.TranslateError("error.load_private_key_failed", err, keymanagePrivKey))
	}

	// 生成新路径
	basePub := filepath.Base(keymanagePubKey)
	basePriv := filepath.Base(keymanagePrivKey)
	newPubPath := filepath.Join(keymanageOutputDir, basePub)
	newPrivPath := filepath.Join(keymanageOutputDir, basePriv)

	// 保存到新位置
	if err := zjcrypto.SaveKeyFiles(
		hybridPub.Kyber,
		hybridPub.ECDH,
		hybridPriv.Kyber,
		hybridPriv.ECDH,
		newPubPath,
		newPrivPath,
	); err != nil {
		return fmt.Errorf("save keys failed: %w",
			i18n.TranslateError("error.save_keys_failed", err))
	}

	fmt.Printf(i18n.T("status.success_import")+"\n", keymanageOutputDir)
	fmt.Printf("  %s: %s\n", i18n.T("status.public_key"), basePub)
	fmt.Printf("  %s: %s\n", i18n.T("status.private_key"), basePriv)

	return nil
}

// verify: 验证密钥对是否匹配
func runVerify() error {
	if keymanagePubKey == "" || keymanagePrivKey == "" {
		return fmt.Errorf("%s", i18n.T("error.missing_both_keys"))
	}

	fmt.Println(i18n.T("status.success_verify") + "...")

	// 加载公钥
	hybridPub, err := zjcrypto.LoadPublicKey(keymanagePubKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("load public key failed: %w",
			i18n.TranslateError("error.load_public_key_failed", err, keymanagePubKey))
	}

	// 加载私钥
	hybridPriv, err := zjcrypto.LoadPrivateKey(keymanagePrivKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("load private key failed: %w",
			i18n.TranslateError("error.load_private_key_failed", err, keymanagePrivKey))
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
		fmt.Println(i18n.T("status.success_verify"))
		fmt.Printf(i18n.T("keymanage_verify.kyber")+"\n", "✅")
		fmt.Printf(i18n.T("keymanage_verify.ecdh")+"\n", "✅")
	} else {
		fmt.Println(i18n.T("status.failed_verify"))
		if !kyberMatch {
			fmt.Printf(i18n.T("keymanage_verify.kyber")+"\n", "❌")
		}
		if !ecdhMatch {
			fmt.Printf(i18n.T("keymanage_verify.ecdh")+"\n", "❌")
		}
		return fmt.Errorf("verify keys failed: %w",
			i18n.TranslateError("error.verify_keys_failed"))
	}

	return nil
}
