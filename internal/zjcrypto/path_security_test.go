package zjcrypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFilePathInDirectory(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "fzjjyz-base-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(baseDir)
	}()

	validPath := filepath.Join("nested", "file.txt")
	validAbs, err := validateFilePathInDirectory(validPath, baseDir)
	if err != nil {
		t.Fatalf("合法路径不应报错: %v", err)
	}
	if !filepath.IsAbs(validAbs) {
		t.Fatalf("返回路径应为绝对路径: %s", validAbs)
	}

	// ../ 逃逸必须被阻止
	if _, err := validateFilePathInDirectory(filepath.Join("..", "evil.txt"), baseDir); err == nil {
		t.Fatal("路径遍历未被拦截")
	}
}

func TestIsSubPath(t *testing.T) {
	baseDir, err := os.MkdirTemp("", "fzjjyz-base-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(baseDir)
	}()

	targetInside := filepath.Join(baseDir, "a", "b.txt")
	ok, err := isSubPath(baseDir, targetInside)
	if err != nil {
		t.Fatalf("isSubPath 返回错误: %v", err)
	}
	if !ok {
		t.Fatal("目录内路径应返回 true")
	}

	targetOutside := filepath.Clean(filepath.Join(baseDir, "..", "outside.txt"))
	ok, err = isSubPath(baseDir, targetOutside)
	if err != nil {
		t.Fatalf("isSubPath 返回错误: %v", err)
	}
	if ok {
		t.Fatal("目录外路径应返回 false")
	}
}
