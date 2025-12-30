package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
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
		Short: i18n.T("keygen.short"),
		Long:  i18n.T("keygen.long"),
		RunE:  runKeygen,
	}

	cmd.Flags().StringVarP(&keygenOutputDir, "output-dir", "d", ".", i18n.T("keygen.flags.output-dir"))
	cmd.Flags().StringVarP(&keygenName, "name", "n", "", i18n.T("keygen.flags.name"))
	cmd.Flags().BoolVarP(&keygenForce, "force", "f", false, i18n.T("keygen.flags.force"))

	return cmd
}

func runKeygen(cmd *cobra.Command, args []string) error {
	// 生成默认名称
	if keygenName == "" {
		keygenName = fmt.Sprintf("key_%d", time.Now().Unix())
	}

	// 创建输出目录
	if err := os.MkdirAll(keygenOutputDir, 0755); err != nil {
		return fmt.Errorf(i18n.T("error.cannot_create_dir"), keygenOutputDir, err)
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
				return fmt.Errorf(i18n.T("error.output_file_exists"), p)
			}
		}
	}

	// 显示生成信息
	if verbose {
		fmt.Printf("%s\n", i18n.T("status.generating_keys"))
		fmt.Printf("  %s: %s\n", i18n.T("keygen.flags.output-dir"), keygenOutputDir)
		fmt.Printf("  %s: %s\n", i18n.T("keygen.flags.name"), keygenName)
	}

	// 1. 生成 Kyber 密钥对
	fmt.Printf("  [1/4] %s ", i18n.T("progress.generating_kyber"))
	kyberPub, kyberPriv, err := crypto.GenerateKyberKeys()
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.keygen_kyber_failed", err)
	}
	fmt.Println(i18n.T("status.done"))

	// 2. 生成 ECDH 密钥对
	fmt.Printf("  [2/4] %s ", i18n.T("progress.generating_ecdh"))
	ecdhPub, ecdhPriv, err := crypto.GenerateECDHKeys()
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.keygen_ecdh_failed", err)
	}
	fmt.Println(i18n.T("status.done"))

	// 3. 生成 Dilithium 密钥对
	fmt.Printf("  [3/4] %s ", i18n.T("progress.generating_dilithium"))
	dilithiumPub, dilithiumPriv, err := crypto.GenerateDilithiumKeys()
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.keygen_dilithium_failed", err)
	}
	fmt.Println(i18n.T("status.done"))

	// 4. 保存密钥文件
	fmt.Printf("  [4/4] %s ", i18n.T("progress.saving_keys"))

	// 保存 Kyber+ECDH 密钥对
	if err := crypto.SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.save_keys_failed", err)
	}

	// 保存 Dilithium 密钥对
	if err := crypto.SaveDilithiumKeys(dilithiumPub, dilithiumPriv, dilithiumPubPath, dilithiumPrivPath); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.save_dilithium_failed", err)
	}

	fmt.Println(i18n.T("status.done"))

	// 显示摘要
	fmt.Println("\n" + i18n.T("status.success_keygen"))
	fmt.Printf("\n"+i18n.T("keygen_info.files")+"\n",
		filepath.Base(pubPath),
		filepath.Base(privPath),
		filepath.Base(dilithiumPubPath),
		filepath.Base(dilithiumPrivPath))

	fmt.Println("\n" + i18n.T("security.warning"))
	fmt.Println(i18n.T("security.protect_keys"))
	fmt.Println(i18n.T("security.no_sharing"))
	fmt.Println(i18n.T("security.secure_storage"))

	return nil
}
