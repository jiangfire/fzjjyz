// Package main 提供文件加密解密命令行工具.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz/utils"
	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"codeberg.org/jiangfire/fzjjyz/internal/zjcrypto"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
	"github.com/spf13/cobra"
)

var (
	decryptDirInput      string
	decryptDirOutput     string
	decryptDirPrivKey    string
	decryptDirVerifyKey  string
	decryptDirForce      bool
	decryptDirBufferSize int
	decryptDirStreaming  bool
)

func newDecryptDirCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt-dir",
		Short: i18n.T("decrypt-dir.short"),
		Long:  i18n.T("decrypt-dir.long"),
		RunE:  runDecryptDir,
	}

	cmd.Flags().StringVarP(&decryptDirInput, "input", "i", "", i18n.T("decrypt-dir.flags.input"))
	cmd.Flags().StringVarP(&decryptDirOutput, "output", "o", "", i18n.T("decrypt-dir.flags.output"))
	cmd.Flags().StringVarP(&decryptDirPrivKey, "private-key", "p", "", i18n.T("decrypt-dir.flags.private-key"))
	cmd.Flags().StringVarP(&decryptDirVerifyKey, "verify-key", "s", "", i18n.T("decrypt-dir.flags.verify-key"))
	cmd.Flags().BoolVarP(&decryptDirForce, "force", "f", false, i18n.T("decrypt-dir.flags.force"))
	cmd.Flags().IntVar(&decryptDirBufferSize, "buffer-size", 0, i18n.T("decrypt-dir.flags.buffer-size"))
	cmd.Flags().BoolVar(&decryptDirStreaming, "streaming", true, i18n.T("decrypt-dir.flags.streaming"))

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagRequired("private-key")

	return cmd
}

//nolint:gocognit,funlen
func runDecryptDir(_ *cobra.Command, _ []string) error {
	// 验证输入文件
	//nolint:wrapcheck
	if err := utils.ValidateInputFile(decryptDirInput); err != nil {
		return err
	}

	// 验证输出目录（如果已存在）
	outputInfo, err := os.Stat(decryptDirOutput)
	if err == nil {
		if !outputInfo.IsDir() {
			return fmt.Errorf(i18n.T("error.output_not_dir"), decryptDirOutput)
		}
		if !decryptDirForce {
			// 检查目录是否为空
			entries, _ := os.ReadDir(decryptDirOutput)
			if len(entries) > 0 {
				return fmt.Errorf(i18n.T("error.output_dir_not_empty"), decryptDirOutput)
			}
		}
	}

	// 读取文件头以获取信息
	headerFile, err := os.Open(decryptDirInput) // #nosec G304 - 文件路径来自用户输入，已通过参数验证
	if err != nil {
		return fmt.Errorf(i18n.T("error.cannot_open_file"), err)
	}
	defer func() {
		_ = headerFile.Close()
	}()

	header, err := format.ParseFileHeader(headerFile)
	if err != nil {
		return fmt.Errorf(i18n.T("error.parse_header_failed"), err)
	}

	// 显示进度
	fmt.Printf(i18n.T("status.decrypting_dir")+"\n", filepath.Base(decryptDirInput))
	if verbose {
		fmt.Printf("  %s: %s\n", i18n.T("file_info.encrypted_file"), decryptDirInput)
		fmt.Printf("  %s: %s\n", i18n.T("file_info.output_dir"), decryptDirOutput)
		fmt.Printf("  %s: %s\n", i18n.T("status.private_key"), decryptDirPrivKey)
		if decryptDirVerifyKey != "" {
			fmt.Printf("  %s: %s\n", i18n.T("status.verify_key"), decryptDirVerifyKey)
		}
		fmt.Printf("  %s: %s\n", i18n.T("file_info.original_filename"), header.Filename)
	}

	// [1/4] 加载密钥
	fmt.Printf("\n[1/4] %s ", i18n.T("progress.loading_keys"))
	hybridPriv, err := utils.LoadHybridPrivateKey(decryptDirPrivKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		//nolint:wrapcheck
		return err
	}

	var dilithiumPub = (*mode3.PublicKey)(nil)
	if decryptDirVerifyKey != "" {
		dilithiumPub, err = utils.LoadDilithiumVerifyKey(decryptDirVerifyKey)
		if err != nil {
			fmt.Println(i18n.T("status.failed"))
			//nolint:wrapcheck
			return err
		}
	} else {
		fmt.Println(i18n.T("status.warning_no_sign_verify"))
	}
	fmt.Println(i18n.T("status.done"))

	// [2/4] 解密数据
	fmt.Printf("[2/4] %s ", i18n.T("progress.decrypting"))

	// 确定缓冲区大小
	bufSize := calculateBufferSizeFromFile(decryptDirInput, decryptDirBufferSize)

	if verbose {
		fmt.Printf(i18n.T("file_info.buffer_size")+"\n", bufSize/1024)
	}

	// 使用随机临时文件，避免固定路径带来的覆盖和竞争风险
	tempZipFile, err := os.CreateTemp("", "fzjjyz-decrypt-*.zip")
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf(i18n.T("error.cannot_open_temp"), err)
	}
	tempZipPath := tempZipFile.Name()
	if closeErr := tempZipFile.Close(); closeErr != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf(i18n.T("error.cannot_open_temp"), closeErr)
	}

	if err := runDecryptWithMode(
		decryptDirInput,
		tempZipPath,
		hybridPriv,
		dilithiumPub,
		decryptDirStreaming,
		bufSize,
	); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("decrypt failed: %w",
			i18n.TranslateError("error.decrypt_failed", err))
	}
	defer func() {
		if removeErr := os.Remove(tempZipPath); removeErr != nil {
			_ = removeErr // 忽略清理错误，不影响主流程
		}
	}() // 清理临时文件

	// 读取解密的ZIP数据
	zipData, err := os.ReadFile(tempZipPath) // #nosec G304 - 临时文件路径由程序生成，安全可控
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("cannot read data: %w",
			i18n.TranslateError("error.cannot_read_data", err))
	}
	// 读取后尽快删除明文临时文件，减少明文驻留时间窗口
	if removeErr := os.Remove(tempZipPath); removeErr == nil {
		tempZipPath = ""
	}

	zipSize := len(zipData)
	fileCount, _ := zjcrypto.CountZipFiles(zipData)
	fmt.Printf(i18n.T("archive.decrypted")+"\n", zipSize)

	// [3/4] 解压ZIP
	fmt.Printf("[3/4] %s ", i18n.T("progress.extracting"))
	if err := zjcrypto.ExtractZipToDirectory(zipData, decryptDirOutput); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return fmt.Errorf("extract failed: %w",
			i18n.TranslateError("error.extract_failed", err))
	}
	fmt.Println(i18n.T("status.done"))

	// [4/4] 验证结果
	fmt.Printf("[4/4] %s ", i18n.T("progress.verifying"))
	fmt.Println(i18n.T("status.done"))

	// 显示结果
	fmt.Printf("\n%s\n\n", i18n.T("status.success_decrypt"))
	summary := i18n.T("dir_info.decrypt_summary")
	fmt.Printf("%s\n",
		fmt.Sprintf(summary,
			filepath.Base(decryptDirInput), zipSize,
			zipSize,
			fileCount,
			decryptDirOutput,
			header.Filename,
			format.UnixTime(header.Timestamp)))

	return nil
}
