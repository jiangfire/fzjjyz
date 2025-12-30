package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("private-key")

	return cmd
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(decryptInput); err != nil {
		return fmt.Errorf(i18n.T("error.encrypted_file_not_exists"), decryptInput)
	}

	// 读取文件头以获取原始文件名（使用缓存读取，避免大文件问题）
	// 对于流式解密，我们只需要读取头部
	headerFile, err := os.Open(decryptInput)
	if err != nil {
		return fmt.Errorf(i18n.T("error.cannot_open_file"), err)
	}
	defer headerFile.Close()

	header, err := format.ParseFileHeader(headerFile)
	if err != nil {
		return fmt.Errorf(i18n.T("error.parse_header_failed"), err)
	}

	// 设置默认输出路径
	if decryptOutput == "" {
		decryptOutput = header.Filename
	}

	// 检查输出文件是否已存在
	if !decryptForce {
		if _, err := os.Stat(decryptOutput); err == nil {
			return fmt.Errorf(i18n.T("error.output_file_exists"), decryptOutput)
		}
	}

	// 显示进度
	fmt.Printf(i18n.T("status.decrypting_file")+"\n", filepath.Base(decryptInput))
	if verbose {
		fmt.Printf("  %s: %s\n", i18n.T("file_info.encrypted_file"), decryptInput)
		fmt.Printf("  %s: %s\n", i18n.T("file_info.decrypted_file"), decryptOutput)
		fmt.Printf("  %s: %s\n", i18n.T("status.private_key"), decryptPrivKey)
		if decryptVerifyKey != "" {
			fmt.Printf("  %s: %s\n", i18n.T("status.verify_key"), decryptVerifyKey)
		}
		fmt.Printf("  %s: %s\n", i18n.T("file_info.original_filename"), header.Filename)
		fmt.Printf("  %s: %v\n", i18n.T("status.streaming_mode"), decryptStreaming)
	}

	// 加载密钥（使用缓存）
	fmt.Printf("\n[1/3] %s ", i18n.T("progress.loading_keys"))
	hybridPriv, err := crypto.LoadPrivateKeyCached(decryptPrivKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.load_private_key_failed", err, decryptPrivKey)
	}

	var dilithiumPub interface{}
	if decryptVerifyKey != "" {
		dilithiumPub, err = crypto.LoadDilithiumPublicKeyCached(decryptVerifyKey)
		if err != nil {
			fmt.Println(i18n.T("status.failed"))
			return i18n.TranslateError("error.load_verify_key_failed", err, decryptVerifyKey)
		}
	} else {
		fmt.Println(i18n.T("status.warning_no_sign_verify"))
	}
	fmt.Println(i18n.T("status.done"))

	// 确定缓冲区大小
	var bufSize int
	if decryptBufferSize > 0 {
		bufSize = decryptBufferSize * 1024
	} else {
		stat, _ := os.Stat(decryptInput)
		bufSize = crypto.OptimalBufferSize(stat.Size())
	}

	if verbose {
		fmt.Printf(i18n.T("file_info.buffer_size")+"\n", bufSize/1024)
	}

	// 执行解密
	fmt.Printf("[2/3] %s ", i18n.T("progress.decrypting"))
	var decryptFunc func() error
	if decryptStreaming {
		decryptFunc = func() error {
			return crypto.DecryptFileStreaming(
				decryptInput, decryptOutput,
				hybridPriv.Kyber, hybridPriv.ECDH,
				dilithiumPub,
				bufSize,
			)
		}
	} else {
		decryptFunc = func() error {
			return crypto.DecryptFile(
				decryptInput, decryptOutput,
				hybridPriv.Kyber, hybridPriv.ECDH,
				dilithiumPub,
			)
		}
	}

	if err := decryptFunc(); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.decrypt_failed", err)
	}
	fmt.Println(i18n.T("status.done"))

	// 显示结果
	fmt.Printf("[3/3] %s ", i18n.T("progress.verifying"))
	fmt.Println(i18n.T("status.done"))

	// 获取文件信息
	decryptedInfo, err := os.Stat(decryptOutput)
	if err != nil {
		fmt.Println("\n" + i18n.T("status.failed"))
		return nil
	}

	encryptedInfo, _ := os.Stat(decryptInput)

	fmt.Printf("\n%s\n\n", i18n.T("status.success_decrypt"))
	summary := i18n.T("file_info.decrypt_summary")
	fmt.Printf("%s\n",
		fmt.Sprintf(summary,
			filepath.Base(decryptInput), encryptedInfo.Size(),
			filepath.Base(decryptOutput), decryptedInfo.Size(),
			header.Filename,
			format.UnixTime(header.Timestamp)))

	return nil
}
