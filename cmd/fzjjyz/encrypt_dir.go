// Package main 提供文件加密解密命令行工具.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz/utils"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"codeberg.org/jiangfire/fzjjyz/internal/zjcrypto"
	"github.com/spf13/cobra"
)

var (
	encryptDirInput      string
	encryptDirOutput     string
	encryptDirPubKey     string
	encryptDirSignKey    string
	encryptDirForce      bool
	encryptDirBufferSize int
	encryptDirStreaming  bool
)

func newEncryptDirCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt-dir",
		Short: i18n.T("encrypt-dir.short"),
		Long:  i18n.T("encrypt-dir.long"),
		RunE:  runEncryptDir,
	}

	cmd.Flags().StringVarP(&encryptDirInput, "input", "i", "", i18n.T("encrypt-dir.flags.input"))
	cmd.Flags().StringVarP(&encryptDirOutput, "output", "o", "", i18n.T("encrypt-dir.flags.output"))
	cmd.Flags().StringVarP(&encryptDirPubKey, "public-key", "p", "", i18n.T("encrypt-dir.flags.public-key"))
	cmd.Flags().StringVarP(&encryptDirSignKey, "sign-key", "s", "", i18n.T("encrypt-dir.flags.sign-key"))
	cmd.Flags().BoolVarP(&encryptDirForce, "force", "f", false, i18n.T("encrypt-dir.flags.force"))
	cmd.Flags().IntVar(&encryptDirBufferSize, "buffer-size", 0, i18n.T("encrypt-dir.flags.buffer-size"))
	cmd.Flags().BoolVar(&encryptDirStreaming, "streaming", true, i18n.T("encrypt-dir.flags.streaming"))

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagRequired("public-key")
	_ = cmd.MarkFlagRequired("sign-key")

	return cmd
}

//nolint:funlen
func runEncryptDir(_ *cobra.Command, _ []string) error {
	// 验证源目录
	//nolint:wrapcheck
	if err := utils.ValidateInputDir(encryptDirInput); err != nil {
		return err
	}

	// 检查输出文件是否已存在
	//nolint:wrapcheck
	if err := utils.CheckOutputConflict(encryptDirOutput, encryptDirForce); err != nil {
		return err
	}

	// 显示进度
	fmt.Printf(i18n.T("status.encrypting_dir")+"\n", filepath.Base(encryptDirInput))
	if verbose {
		fmt.Printf("  %s: %s\n", i18n.T("file_info.source_dir"), encryptDirInput)
		fmt.Printf("  %s: %s\n", i18n.T("file_info.encrypted_file"), encryptDirOutput)
		fmt.Printf("  %s: %s\n", i18n.T("status.public_key"), encryptDirPubKey)
		fmt.Printf("  %s: %s\n", i18n.T("status.sign_key"), encryptDirSignKey)
	}

	// [1/4] 打包成ZIP
	fmt.Printf("\n")
	fmt.Printf("[1/4] %s ", i18n.T("progress.packing"))
	var zipBuffer bytes.Buffer
	if err := zjcrypto.CreateZipFromDirectory(encryptDirInput, &zipBuffer, zjcrypto.DefaultArchiveOptions); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("pack failed: %w",
			i18n.TranslateError("error.pack_failed", err))
	}
	zipData := zipBuffer.Bytes()

	// 获取ZIP信息
	zipSize := len(zipData)
	fileCount, _ := zjcrypto.CountZipFiles(zipData)
	fmt.Printf(i18n.T("archive.packed")+"\n", zipSize, fileCount)

	// [2/4] 加载密钥
	fmt.Printf("[2/4] %s ", i18n.T("progress.loading_keys"))
	hybridPub, err := utils.LoadHybridPublicKey(encryptDirPubKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		//nolint:wrapcheck
		return err
	}

	dilithiumPriv, err := utils.LoadDilithiumPrivateKey(encryptDirSignKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		//nolint:wrapcheck
		return err
	}
	fmt.Println(i18n.T("status.done"))

	// [3/4] 加密ZIP数据
	fmt.Printf("[3/4] %s ", i18n.T("progress.encrypting"))

	// 临时保存ZIP到文件（以便复用现有加密流程）
	tempZipFile, err := os.CreateTemp(filepath.Dir(encryptDirOutput), "fzjjyz-encrypt-*.zip")
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("temp file failed: %w",
			i18n.TranslateError("error.temp_file_failed", err))
	}
	tempZipPath := tempZipFile.Name()
	if closeErr := tempZipFile.Close(); closeErr != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("temp file failed: %w",
			i18n.TranslateError("error.temp_file_failed", closeErr))
	}
	defer func() {
		_ = os.Remove(tempZipPath)
	}()

	if err := os.WriteFile(tempZipPath, zipData, 0600); err != nil { // #nosec G703
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("temp file failed: %w",
			i18n.TranslateError("error.temp_file_failed", err))
	}

	// 确定缓冲区大小
	bufSize := calculateBufferSizeFromLength(int64(zipSize), encryptDirBufferSize)

	if verbose {
		fmt.Printf(i18n.T("file_info.buffer_size")+"\n", bufSize/1024)
	}

	if err := runEncryptWithMode(
		tempZipPath,
		encryptDirOutput,
		hybridPub,
		dilithiumPriv,
		encryptDirStreaming,
		bufSize,
	); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("encrypt failed: %w",
			i18n.TranslateError("error.encrypt_failed", err))
	}
	fmt.Println(i18n.T("status.done"))

	// [4/4] 验证结果
	fmt.Printf("[4/4] %s ", i18n.T("progress.verifying"))
	encryptedInfo, _ := os.Stat(encryptDirOutput)
	fmt.Println(i18n.T("status.done"))

	// 显示结果
	fmt.Printf("\n%s\n\n", i18n.T("status.success_encrypt"))
	summary := i18n.T("dir_info.encrypt_summary")
	fmt.Printf("%s\n",
		fmt.Sprintf(summary,
			encryptDirInput, fileCount,
			zipSize,
			filepath.Base(encryptDirOutput), encryptedInfo.Size(),
			float64(encryptedInfo.Size())/float64(zipSize)*100))

	return nil
}
