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

	// 不允许绝对路径，防止绕过基础目录
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute path is not allowed: %s", path)
	}

	cleanPath := filepath.Clean(path)
	if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute base: %w", err)
	}

	absPath, err := filepath.Abs(filepath.Join(absBase, cleanPath))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	inBase, err := isSubPath(absBase, absPath)
	if err != nil {
		return "", fmt.Errorf("failed to validate path scope: %w", err)
	}
	if !inBase {
		return "", fmt.Errorf("path escapes base directory: %s", path)
	}

	return absPath, nil
}

// validateAndExtractPath 安全地验证解压路径.
func validateAndExtractPath(path, targetDir string) (string, error) {
	return validateFilePathInDirectory(path, targetDir)
}

// isSubPath 检查 target 是否位于 base 目录内（含 base 本身）.
func isSubPath(base, target string) (bool, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false, err
	}
	rel = filepath.Clean(rel)
	if rel == "." {
		return true, nil
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false, nil
	}
	return true, nil
}
