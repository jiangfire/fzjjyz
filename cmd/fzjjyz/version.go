package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Long:  `显示 fzjjyz 的版本信息和构建详情`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("fzjjyz - 后量子文件加密工具\n")
			fmt.Printf("版本: %s\n", Version)
			fmt.Printf("应用名称: %s\n", AppName)
			fmt.Printf("描述: %s\n", Description)
		},
	}

	return cmd
}
