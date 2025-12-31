// Package main 提供文件加密解密命令行工具.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz/utils"
	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
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
		Short: i18n.T("decrypt.short"),
		Long:  i18n.T("decrypt.long"),
		RunE:  runDecrypt,
	}

	cmd.Flags().StringVarP(&decryptInput, "input", "i", "", i18n.T("decrypt.flags.input"))
	cmd.Flags().StringVarP(&decryptOutput, "output", "o", "", i18n.T("decrypt.flags.output"))
	cmd.Flags().StringVarP(&decryptPrivKey, "private-key", "p", "", i18n.T("decrypt.flags.private-key"))
	cmd.Flags().StringVarP(&decryptVerifyKey, "verify-key", "s", "", i18n.T("decrypt.flags.verify-key"))
	cmd.Flags().BoolVarP(&decryptForce, "force", "f", false, i18n.T("decrypt.flags.force"))
	cmd.Flags().IntVar(&decryptBufferSize, "buffer-size", 0, i18n.T("decrypt.flags.buffer-size"))
	cmd.Flags().BoolVar(&decryptStreaming, "streaming", true, i18n.T("decrypt.flags.streaming"))

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("private-key")

	return cmd
}

func runDecrypt(_ *cobra.Command, _ []string) error {
	return executeDecryptCommand()
}

func executeDecryptCommand() error {
	// 步骤1: 验证输入
	//nolint:wrapcheck
	if err := utils.ValidateInputFile(decryptInput); err != nil {
		return err
	}
	//nolint:wrapcheck
	if err := utils.CheckOutputConflict(decryptOutput, decryptForce); err != nil {
		return err
	}

	// 步骤2: 解析文件头
	header, err := parseDecryptHeader()
	if err != nil {
		return err
	}

	// 步骤3: 加载密钥
	reporter := utils.NewProgressReporter(3, verbose)
	hybridPriv, dilithiumPub, err := loadDecryptKeys(reporter)
	if err != nil {
		return err
	}

	// 步骤4: 执行解密
	if err := executeDecrypt(reporter, hybridPriv, dilithiumPub, header); err != nil {
		return err
	}

	// 步骤5: 显示结果
	return showDecryptResult(header)
}

func parseDecryptHeader() (*format.FileHeader, error) {
	headerFile := utils.MustOpen(decryptInput)
	defer func() {
		_ = headerFile.Close()
	}()

	header, err := format.ParseFileHeader(headerFile)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.parse_header_failed"), err)
	}

	// 设置默认输出路径
	if decryptOutput == "" {
		decryptOutput = header.Filename
	}

	return header, nil
}

func loadDecryptKeys(reporter *utils.ProgressReporter) (*crypto.HybridPrivateKey, interface{}, error) {
	reporter.Step("progress.loading_keys")

	// 加载私钥
	hybridPriv, err := utils.LoadHybridPrivateKey(decryptPrivKey)
	if err != nil {
		reporter.Failed()
		//nolint:wrapcheck
		return nil, nil, err
	}

	// 加载验证公钥
	dilithiumPub, err := utils.LoadDilithiumVerifyKey(decryptVerifyKey)
	if err != nil {
		reporter.Failed()
		//nolint:wrapcheck
		return nil, nil, err
	}

	// 显示警告
	if decryptVerifyKey == "" {
		reporter.Warning("status.warning_no_sign_verify")
	}

	reporter.Done()
	return hybridPriv, dilithiumPub, nil
}

func executeDecrypt(
	reporter *utils.ProgressReporter,
	hybridPriv *crypto.HybridPrivateKey,
	dilithiumPub interface{},
	header *format.FileHeader,
) error {
	// 显示详细信息
	reporter.InfoString("file_info.encrypted_file", decryptInput)
	reporter.InfoString("file_info.decrypted_file", decryptOutput)
	reporter.InfoString("status.private_key", decryptPrivKey)
	if decryptVerifyKey != "" {
		reporter.InfoString("status.verify_key", decryptVerifyKey)
	}
	reporter.InfoString("file_info.original_filename", header.Filename)
	reporter.InfoBool("status.streaming_mode", decryptStreaming)

	// 计算缓冲区大小
	bufSize := getDecryptBufferSize()
	reporter.Info("file_info.buffer_size", bufSize/1024)

	// 执行解密
	reporter.Step("progress.decrypting")
	decryptFunc := getDecryptFunction(hybridPriv, dilithiumPub, bufSize)
	if err := decryptFunc(); err != nil {
		reporter.Failed()
		return fmt.Errorf("decrypt failed: %w",
			i18n.TranslateError("error.decrypt_failed", err))
	}
	reporter.Done()

	// 验证步骤
	reporter.Step("progress.verifying")
	reporter.Done()

	return nil
}

func getDecryptBufferSize() int {
	if decryptBufferSize > 0 {
		return decryptBufferSize * 1024
	}
	size, _ := utils.GetFileSize(decryptInput)
	return crypto.OptimalBufferSize(size)
}

func getDecryptFunction(hybridPriv *crypto.HybridPrivateKey, dilithiumPub interface{}, bufSize int) func() error {
	if decryptStreaming {
		return func() error {
			return crypto.DecryptFileStreaming(
				decryptInput, decryptOutput,
				hybridPriv.Kyber, hybridPriv.ECDH,
				dilithiumPub,
				bufSize,
			)
		}
	}
	return func() error {
		return crypto.DecryptFile(
			decryptInput, decryptOutput,
			hybridPriv.Kyber, hybridPriv.ECDH,
			dilithiumPub,
		)
	}
}

func showDecryptResult(header *format.FileHeader) error {
	decryptedInfo, err := os.Stat(decryptOutput)
	if err != nil {
		fmt.Println("\n" + i18n.T("status.failed"))
		return nil //nolint:nilerr
	}

	encryptedInfo, _ := os.Stat(decryptInput)

	reporter := utils.NewProgressReporter(1, true)
	reporter.Summary("status.success_decrypt")

	summary := i18n.T("file_info.decrypt_summary")
	fmt.Printf("%s\\n",
		fmt.Sprintf(summary,
			filepath.Base(decryptInput), encryptedInfo.Size(),
			filepath.Base(decryptOutput), decryptedInfo.Size(),
			header.Filename,
			format.UnixTime(header.Timestamp)))

	return nil
}
