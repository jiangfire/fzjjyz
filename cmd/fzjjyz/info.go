package main

import (
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"github.com/spf13/cobra"
)

var (
	infoInput string
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "查看加密文件信息",
		Long: `解析并显示加密文件的详细信息，包括：
  • 文件名和原始大小
  • 加密时间戳
  • 使用的算法
  • 签名状态
  • 完整性验证

示例：
  fzjjyz info -i encrypted.fzj
  fzjjyz info --input data.fzj`,
		RunE: runInfo,
	}

	cmd.Flags().StringVarP(&infoInput, "input", "i", "", "加密文件路径 (必需)")
	cmd.MarkFlagRequired("input")

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) error {
	// 验证输入文件
	if _, err := os.Stat(infoInput); err != nil {
		return fmt.Errorf("文件不存在: %s", infoInput)
	}

	// 读取文件
	data, err := os.ReadFile(infoInput)
	if err != nil {
		return fmt.Errorf("无法读取文件: %v", err)
	}

	// 解析文件头
	header, err := format.ParseFileHeaderFromBytes(data)
	if err != nil {
		return fmt.Errorf("文件头解析失败: %v", err)
	}

	// 验证文件头
	if err := header.Validate(); err != nil {
		return fmt.Errorf("文件头验证失败: %v", err)
	}

	// 获取文件信息
	fileInfo, _ := os.Stat(infoInput)

	// 显示信息
	fmt.Printf("📁 文件信息: %s\n\n", filepath.Base(infoInput))

	// 基本信息
	fmt.Println("基本信息:")
	fmt.Printf("  文件名:        %s\n", header.Filename)
	fmt.Printf("  原始大小:      %d bytes\n", header.FileSize)
	fmt.Printf("  加密大小:      %d bytes\n", fileInfo.Size())
	fmt.Printf("  压缩率:        %.1f%%\n", float64(fileInfo.Size())/float64(header.FileSize)*100)
	fmt.Printf("  时间戳:        %s\n", format.UnixTime(header.Timestamp))

	// 算法信息
	fmt.Println("\n加密信息:")
	algoName := "未知"
	if header.Algorithm == 0x02 {
		algoName = "Kyber768 + ECDH + AES-256-GCM"
	}
	fmt.Printf("  算法:          %s (0x%02x)\n", algoName, header.Algorithm)
	fmt.Printf("  版本:          0x%04x\n", header.Version)
	fmt.Printf("  魔数:          %c%c%c\\x%02x\n", header.Magic[0], header.Magic[1], header.Magic[2], header.Magic[3])

	// 密钥信息
	fmt.Println("\n密钥信息:")
	fmt.Printf("  Kyber封装:     %d bytes\n", header.KyberEncLen)
	fmt.Printf("  ECDH公钥:      %d bytes\n", header.ECDHLen)
	fmt.Printf("  IV/Nonce:      %d bytes\n", header.IVLen)
	fmt.Printf("  签名:          %d bytes\n", header.SigLen)

	// 完整性信息
	fmt.Println("\n完整性:")
	fmt.Printf("  SHA256哈希:    %x...\n", header.SHA256Hash[:8])

	// 验证状态
	fmt.Println("\n验证状态:")
	if header.SigLen > 0 {
		fmt.Println("  签名:          ✅ 存在")
	} else {
		fmt.Println("  签名:          ❌ 不存在")
	}

	// 检查数据完整性
	expectedSize := header.GetHeaderSize()
	if len(data) > expectedSize {
		fmt.Println("  数据完整性:   ✅ 完整")
	} else {
		fmt.Println("  数据完整性:   ❌ 不完整")
	}

	return nil
}
