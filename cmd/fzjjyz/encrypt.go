package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	cmd.MarkFlagRequired("input")
	cmd.MarkFlagRequired("public-key")
	cmd.MarkFlagRequired("sign-key")

	return cmd
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(encryptInput); err != nil {
		return fmt.Errorf(i18n.T("error.input_file_not_exists"), encryptInput)
	}

	// 设置默认输出路径
	if encryptOutput == "" {
		encryptOutput = encryptInput + ".fzj"
	}

	// 检查输出文件是否已存在
	if !encryptForce {
		if _, err := os.Stat(encryptOutput); err == nil {
			return fmt.Errorf(i18n.T("error.output_file_exists"), encryptOutput)
		}
	}

	// 显示进度
	fmt.Printf(i18n.T("status.encrypting_file")+"\n", filepath.Base(encryptInput))
	if verbose {
		fmt.Printf(i18n.T("file_info.original_file")+"\n", encryptInput)
		fmt.Printf(i18n.T("file_info.encrypted_file")+"\n", encryptOutput)
		fmt.Printf("  %s: %s\n", i18n.T("status.public_key"), encryptPubKey)
		fmt.Printf("  %s: %s\n", i18n.T("status.sign_key"), encryptSignKey)
		fmt.Printf("  %s: %v\n", i18n.T("status.streaming_mode"), encryptStreaming)
	}

	// 加载密钥（使用缓存）
	fmt.Printf("\n[1/3] %s ", i18n.T("progress.loading_keys"))
	hybridPub, err := crypto.LoadPublicKeyCached(encryptPubKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.load_public_key_failed", err, encryptPubKey)
	}

	dilithiumPriv, err := crypto.LoadDilithiumPrivateKeyCached(encryptSignKey)
	if err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.load_sign_key_failed", err, encryptSignKey)
	}
	fmt.Println(i18n.T("status.done"))

	// 确定缓冲区大小
	var bufSize int
	if encryptBufferSize > 0 {
		bufSize = encryptBufferSize * 1024
	} else {
		stat, _ := os.Stat(encryptInput)
		bufSize = crypto.OptimalBufferSize(stat.Size())
	}

	if verbose {
		fmt.Printf(i18n.T("file_info.buffer_size")+"\n", bufSize/1024)
	}

	// 执行加密
	fmt.Printf("[2/3] %s ", i18n.T("progress.encrypting"))
	var encryptFunc func() error
	if encryptStreaming {
		encryptFunc = func() error {
			return crypto.EncryptFileStreaming(
				encryptInput, encryptOutput,
				hybridPub.Kyber, hybridPub.ECDH,
				dilithiumPriv,
				bufSize,
			)
		}
	} else {
		encryptFunc = func() error {
			return crypto.EncryptFile(
				encryptInput, encryptOutput,
				hybridPub.Kyber, hybridPub.ECDH,
				dilithiumPriv,
			)
		}
	}

	if err := encryptFunc(); err != nil {
		fmt.Println(i18n.T("status.failed"))
		return i18n.TranslateError("error.encrypt_failed", err)
	}
	fmt.Println(i18n.T("status.done"))

	// 显示结果
	fmt.Printf("[3/3] %s ", i18n.T("progress.verifying"))
	encryptedInfo, _ := os.Stat(encryptOutput)
	originalInfo, _ := os.Stat(encryptInput)
	fmt.Println(i18n.T("status.done"))

	fmt.Printf("\n%s\n\n", i18n.T("status.success_encrypt"))
	summary := i18n.T("file_info.encrypt_summary")
	fmt.Printf("%s\n",
		fmt.Sprintf(summary,
			filepath.Base(encryptInput), originalInfo.Size(),
			filepath.Base(encryptOutput), encryptedInfo.Size(),
			float64(encryptedInfo.Size())/float64(originalInfo.Size())*100))

	return nil
}
