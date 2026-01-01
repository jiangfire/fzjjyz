package zjcrypto

import (
	"fmt"
	"path/filepath"
	"strings"
)

// validateFilePathInDirectory 验证路径是否在指定目录内.
// 用于解压等场景，确保不会逃逸到父目录.
func validateFilePathInDirectory(path, baseDir string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if baseDir == "" {
		return "", fmt.Errorf("base directory cannot be empty")
	}

	// 构建完整路径
	fullPath := filepath.Join(baseDir, path)

	// 转换为绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute base: %w", err)
	}

	// 验证路径在基础目录内
	if !strings.HasPrefix(absPath, absBase) {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}

	// 检查路径遍历字符
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}

	return absPath, nil
}

// validateAndExtractPath 安全地验证解压路径.
func validateAndExtractPath(path, targetDir string) (string, error) {
	return validateFilePathInDirectory(path, targetDir)
}
