// Package main 提供文件加密解密命令行工具.
package main

import (
	"crypto/ecdh"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz/utils"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"codeberg.org/jiangfire/fzjjyz/internal/zjcrypto"
	"github.com/cloudflare/circl/kem"
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

func runKeygen(_ *cobra.Command, _ []string) error {
	return executeKeygenCommand()
}

func executeKeygenCommand() error {
	// 步骤1: 准备
	paths, err := prepareKeygen()
	if err != nil {
		return err
	}

	// 步骤2: 生成密钥
	reporter := utils.NewProgressReporter(4, verbose)
	keys, err := generateKeys(reporter)
	if err != nil {
		return err
	}

	// 步骤3: 保存密钥
	if err := saveKeys(reporter, keys, paths); err != nil {
		return err
	}

	// 步骤4: 显示结果
	return showKeygenResult(paths)
}

func prepareKeygen() ([]string, error) {
	// 生成默认名称
	if keygenName == "" {
		keygenName = fmt.Sprintf("key_%d", time.Now().Unix())
	}

	// 创建输出目录
	if err := os.MkdirAll(keygenOutputDir, 0750); err != nil {
		return nil, fmt.Errorf(i18n.T("error.cannot_create_dir"), keygenOutputDir, err)
	}

	// 构建文件路径
	pubPath := filepath.Join(keygenOutputDir, keygenName+"_public.pem")
	privPath := filepath.Join(keygenOutputDir, keygenName+"_private.pem")
	dilithiumPubPath := filepath.Join(keygenOutputDir, keygenName+"_dilithium_public.pem")
	dilithiumPrivPath := filepath.Join(keygenOutputDir, keygenName+"_dilithium_private.pem")

	paths := []string{pubPath, privPath, dilithiumPubPath, dilithiumPrivPath}

	// 检查冲突
	if !keygenForce {
		for _, p := range paths {
			if utils.FileExists(p) {
				return nil, fmt.Errorf(i18n.T("error.output_file_exists"), p)
			}
		}
	}

	// 显示详细信息
	if verbose {
		fmt.Printf("%s\\n", i18n.T("status.generating_keys"))
		fmt.Printf("  %s: %s\\n", i18n.T("keygen.flags.output-dir"), keygenOutputDir)
		fmt.Printf("  %s: %s\\n", i18n.T("keygen.flags.name"), keygenName)
	}

	return paths, nil
}

type keyPair struct {
	kyberPub      interface{}
	kyberPriv     interface{}
	ecdhPub       interface{}
	ecdhPriv      interface{}
	dilithiumPub  interface{}
	dilithiumPriv interface{}
}

func generateKeys(reporter *utils.ProgressReporter) (*keyPair, error) {
	// 1. Kyber
	reporter.Step("progress.generating_kyber")
	kyberPub, kyberPriv, err := zjcrypto.GenerateKyberKeys()
	if err != nil {
		reporter.Failed()
		return nil, fmt.Errorf("kyber key generation failed: %w",
			i18n.TranslateError("error.keygen_kyber_failed", err))
	}
	reporter.Done()

	// 2. ECDH
	reporter.Step("progress.generating_ecdh")
	ecdhPub, ecdhPriv, err := zjcrypto.GenerateECDHKeys()
	if err != nil {
		reporter.Failed()
		return nil, fmt.Errorf("ecdh key generation failed: %w",
			i18n.TranslateError("error.keygen_ecdh_failed", err))
	}
	reporter.Done()

	// 3. Dilithium
	reporter.Step("progress.generating_dilithium")
	dilithiumPub, dilithiumPriv, err := zjcrypto.GenerateDilithiumKeys()
	if err != nil {
		reporter.Failed()
		return nil, fmt.Errorf("dilithium key generation failed: %w",
			i18n.TranslateError("error.keygen_dilithium_failed", err))
	}
	reporter.Done()

	return &keyPair{
		kyberPub: kyberPub, kyberPriv: kyberPriv,
		ecdhPub: ecdhPub, ecdhPriv: ecdhPriv,
		dilithiumPub: dilithiumPub, dilithiumPriv: dilithiumPriv,
	}, nil
}

func saveKeys(reporter *utils.ProgressReporter, keys *keyPair, paths []string) error {
	reporter.Step("progress.saving_keys")

	pubPath, privPath, dilithiumPubPath, dilithiumPrivPath := paths[0], paths[1], paths[2], paths[3]

	// 保存 Kyber+ECDH - 需要类型转换
	kyberPub := keys.kyberPub.(kem.PublicKey)
	kyberPriv := keys.kyberPriv.(kem.PrivateKey)
	ecdhPub := keys.ecdhPub.(*ecdh.PublicKey)
	ecdhPriv := keys.ecdhPriv.(*ecdh.PrivateKey)

	if err := zjcrypto.SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		reporter.Failed()
		return fmt.Errorf("save keys failed: %w",
			i18n.TranslateError("error.save_keys_failed", err))
	}

	// 保存 Dilithium
	if err := zjcrypto.SaveDilithiumKeys(
		keys.dilithiumPub, keys.dilithiumPriv,
		dilithiumPubPath, dilithiumPrivPath,
	); err != nil {
		reporter.Failed()
		return fmt.Errorf("save dilithium keys failed: %w",
			i18n.TranslateError("error.save_dilithium_failed", err))
	}

	reporter.Done()
	return nil
}

func showKeygenResult(paths []string) error {
	pubPath, privPath, dilithiumPubPath, dilithiumPrivPath := paths[0], paths[1], paths[2], paths[3]

	fmt.Println("\n" + i18n.T("status.success_keygen"))
	fmt.Printf("\n"+i18n.T("keygen_info.files")+"\\n",
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
