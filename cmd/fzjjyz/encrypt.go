// Package main 提供文件加密解密命令行工具.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz/utils"
	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
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
		Short: i18n.T("encrypt.short"),
		Long:  i18n.T("encrypt.long"),
		RunE:  runEncrypt,
	}

	cmd.Flags().StringVarP(&encryptInput, "input", "i", "", i18n.T("encrypt.flags.input"))
	cmd.Flags().StringVarP(&encryptOutput, "output", "o", "", i18n.T("encrypt.flags.output"))
	cmd.Flags().StringVarP(&encryptPubKey, "public-key", "p", "", i18n.T("encrypt.flags.public-key"))
	cmd.Flags().StringVarP(&encryptSignKey, "sign-key", "s", "", i18n.T("encrypt.flags.sign-key"))
	cmd.Flags().BoolVarP(&encryptForce, "force", "f", false, i18n.T("encrypt.flags.force"))
	cmd.Flags().IntVar(&encryptBufferSize, "buffer-size", 0, i18n.T("encrypt.flags.buffer-size"))
	cmd.Flags().BoolVar(&encryptStreaming, "streaming", true, i18n.T("encrypt.flags.streaming"))

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("public-key")
	_ = cmd.MarkFlagRequired("sign-key")

	return cmd
}

func runEncrypt(_ *cobra.Command, _ []string) error {
	return executeEncryptCommand()
}

func executeEncryptCommand() error {
	// 步骤1: 验证输入
	//nolint:wrapcheck
	if err := utils.ValidateInputFile(encryptInput); err != nil {
		return err
	}
	//nolint:wrapcheck
	if err := utils.CheckOutputConflict(encryptOutput, encryptForce); err != nil {
		return err
	}

	// 步骤2: 准备输出路径
	prepareEncryptOutput()

	// 步骤3: 加载密钥
	reporter := utils.NewProgressReporter(3, verbose)
	hybridPub, dilithiumPriv, err := loadEncryptKeys(reporter)
	if err != nil {
		return err
	}

	// 步骤4: 执行加密
	if err := executeEncrypt(reporter, hybridPub, dilithiumPriv); err != nil {
		return err
	}

	// 步骤5: 显示结果
	return showEncryptResult()
}

func prepareEncryptOutput() {
	if encryptOutput == "" {
		encryptOutput = encryptInput + ".fzj"
	}
}

func loadEncryptKeys(reporter *utils.ProgressReporter) (*crypto.HybridPublicKey, interface{}, error) {
	reporter.Step("progress.loading_keys")

	// 加载公钥
	hybridPub, err := utils.LoadHybridPublicKey(encryptPubKey)
	if err != nil {
		reporter.Failed()
		//nolint:wrapcheck
		return nil, nil, err
	}

	// 加载签名私钥
	dilithiumPriv, err := utils.LoadDilithiumPrivateKey(encryptSignKey)
	if err != nil {
		reporter.Failed()
		//nolint:wrapcheck
		return nil, nil, err
	}

	reporter.Done()
	return hybridPub, dilithiumPriv, nil
}

func executeEncrypt(
	reporter *utils.ProgressReporter,
	hybridPub *crypto.HybridPublicKey,
	dilithiumPriv interface{},
) error {
	// 显示详细信息
	reporter.InfoString("file_info.original_file", encryptInput)
	reporter.InfoString("file_info.encrypted_file", encryptOutput)
	reporter.InfoString("status.public_key", encryptPubKey)
	reporter.InfoString("status.sign_key", encryptSignKey)
	reporter.InfoBool("status.streaming_mode", encryptStreaming)

	// 计算缓冲区大小
	bufSize := getEncryptBufferSize()
	reporter.Info("file_info.buffer_size", bufSize/1024)

	// 执行加密
	reporter.Step("progress.encrypting")
	encryptFunc := getEncryptFunction(hybridPub, dilithiumPriv, bufSize)
	if err := encryptFunc(); err != nil {
		reporter.Failed()
		return fmt.Errorf("encrypt failed: %w",
			i18n.TranslateError("error.encrypt_failed", err))
	}
	reporter.Done()

	// 验证步骤
	reporter.Step("progress.verifying")
	reporter.Done()

	return nil
}

func getEncryptBufferSize() int {
	if encryptBufferSize > 0 {
		return encryptBufferSize * 1024
	}
	size, _ := utils.GetFileSize(encryptInput)
	return crypto.OptimalBufferSize(size)
}

func getEncryptFunction(hybridPub *crypto.HybridPublicKey, dilithiumPriv interface{}, bufSize int) func() error {
	if encryptStreaming {
		return func() error {
			return crypto.EncryptFileStreaming(
				encryptInput, encryptOutput,
				hybridPub.Kyber, hybridPub.ECDH,
				dilithiumPriv,
				bufSize,
			)
		}
	}
	return func() error {
		return crypto.EncryptFile(
			encryptInput, encryptOutput,
			hybridPub.Kyber, hybridPub.ECDH,
			dilithiumPriv,
		)
	}
}

func showEncryptResult() error {
	encryptedInfo, _ := os.Stat(encryptOutput)
	originalInfo, _ := os.Stat(encryptInput)

	reporter := utils.NewProgressReporter(1, true)
	reporter.Summary("status.success_encrypt")

	summary := i18n.T("file_info.encrypt_summary")
	fmt.Printf("%s\\n",
		fmt.Sprintf(summary,
			filepath.Base(encryptInput), originalInfo.Size(),
			filepath.Base(encryptOutput), encryptedInfo.Size(),
			float64(encryptedInfo.Size())/float64(originalInfo.Size())*100))

	return nil
}
