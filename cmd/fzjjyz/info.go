package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	infoInput string
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: i18n.T("info.short"),
		Long:  i18n.T("info.long"),
		RunE:  runInfo,
	}

	cmd.Flags().StringVarP(&infoInput, "input", "i", "", i18n.T("info.flags.input"))
	cmd.MarkFlagRequired("input")

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(infoInput); err != nil {
		return fmt.Errorf(i18n.T("error.file_not_exists"), infoInput)
	}

	// 读取文件
	data, err := os.ReadFile(infoInput)
	if err != nil {
		return fmt.Errorf(i18n.T("error.cannot_read_file"), err)
	}

	// 解析文件头
	header, err := format.ParseFileHeaderFromBytes(data)
	if err != nil {
		return fmt.Errorf(i18n.T("error.parse_header_failed"), err)
	}

	// 验证文件头
	if err := header.Validate(); err != nil {
		return fmt.Errorf(i18n.T("error.validate_header_failed"), err)
	}

	// 获取文件信息
	fileInfo, _ := os.Stat(infoInput)

	// 显示信息
	fmt.Printf(i18n.T("file_info.header")+"\n\n", filepath.Base(infoInput))

	// 基本信息
	fmt.Println(i18n.T("file_info.basic"))
	fmt.Printf("  "+i18n.T("file_info.original_filename")+"\n", header.Filename)
	fmt.Printf("  "+i18n.T("file_info.original_file")+"\n", "", header.FileSize)
	fmt.Printf("  "+i18n.T("file_info.encrypted_file")+"\n", "", fileInfo.Size())
	fmt.Printf("  "+i18n.T("file_info.compressed_rate")+"\n", float64(fileInfo.Size())/float64(header.FileSize)*100)
	fmt.Printf("  "+i18n.T("file_info.timestamp")+"\n", format.UnixTime(header.Timestamp))

	// 算法信息
	fmt.Println("\n" + i18n.T("file_info.encryption"))
	algoName := "未知"
	if header.Algorithm == 0x02 {
		algoName = "Kyber768 + ECDH + AES-256-GCM"
	}
	fmt.Printf("  "+i18n.T("file_info.algorithm")+"\n", algoName, header.Algorithm)
	fmt.Printf("  "+i18n.T("file_info.version")+"\n", header.Version)
	fmt.Printf("  "+i18n.T("file_info.magic")+"\n", header.Magic[0], header.Magic[1], header.Magic[2], header.Magic[3])

	// 密钥信息
	fmt.Println("\n" + i18n.T("file_info.keys"))
	fmt.Printf("  "+i18n.T("file_info.kyber")+"\n", header.KyberEncLen)
	fmt.Printf("  "+i18n.T("file_info.ecdh")+"\n", header.ECDHLen)
	fmt.Printf("  "+i18n.T("file_info.iv")+"\n", header.IVLen)
	fmt.Printf("  "+i18n.T("file_info.signature")+"\n", header.SigLen)

	// 完整性信息
	fmt.Println("\n" + i18n.T("file_info.integrity"))
	fmt.Printf("  "+i18n.T("file_info.hash")+"\n", header.SHA256Hash[:8])

	// 验证状态
	fmt.Println("\n" + i18n.T("file_info.verification"))
	if header.SigLen > 0 {
		fmt.Printf("  "+i18n.T("file_info.signature_status")+" %s\n", i18n.T("file_info.exists"))
	} else {
		fmt.Printf("  "+i18n.T("file_info.signature_status")+" %s\n", i18n.T("file_info.not_exists"))
	}

	// 检查数据完整性
	expectedSize := header.GetHeaderSize()
	if len(data) > expectedSize {
		fmt.Printf("  "+i18n.T("file_info.data_integrity")+" %s\n", i18n.T("file_info.complete"))
	} else {
		fmt.Printf("  "+i18n.T("file_info.data_integrity")+" %s\n", i18n.T("file_info.incomplete"))
	}

	return nil
}
