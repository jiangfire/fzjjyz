package main

import (
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
	"fmt"
)

func main() {
	i18n.Init("zh_CN")
	msg := i18n.T("progress.generating_kyber")
	fmt.Printf("Translation: %q\n", msg)
	fmt.Printf("Test: ")
	fmt.Printf("  "+msg+" ", 1, 4)
	fmt.Println()
}
