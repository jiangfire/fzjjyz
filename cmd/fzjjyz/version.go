// Package main 提供文件加密解密命令行工具.
package main

import (
	"fmt"

	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: i18n.T("version.short"),
		Long:  i18n.T("version.long"),
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("%s - %s\n", i18n.T("app.name"), i18n.T("app.description"))
			fmt.Printf("%s %s\n", i18n.T("version.label"), Version)
			fmt.Printf("%s %s\n", i18n.T("version.app_name"), AppName)
			fmt.Printf("%s %s\n", i18n.T("version.description"), Description)
		},
	}

	return cmd
}
