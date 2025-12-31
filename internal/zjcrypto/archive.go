// Package crypto 提供文件加密、解密和归档处理功能
package zjcrypto

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

const (
	// maxTotalSize 限制解压总大小，防止解压缩炸弹（1GB）。
	maxTotalSize = 1024 * 1024 * 1024
	// dirPerm 使用更安全的目录权限。
	dirPerm = 0750
)

// ArchiveOptions 打包选项.
type ArchiveOptions struct {
	IncludePatterns []string // 包含的文件模式（glob）
	ExcludePatterns []string // 排除的文件模式（glob）
	FollowSymlinks  bool     // 是否跟随符号链接
}

// DefaultArchiveOptions 默认打包选项.
var DefaultArchiveOptions = ArchiveOptions{
	IncludePatterns: []string{"**/*"},
	ExcludePatterns: []string{},
	FollowSymlinks:  false,
}

// CreateZipFromDirectory 将目录打包成ZIP
// 输入: 源目录路径, 输出缓冲区, 打包选项
// 返回: 错误
//
//nolint:gocognit,funlen // 目录遍历和ZIP打包逻辑复杂且需要完整处理所有文件类型和路径转换
func CreateZipFromDirectory(sourceDir string, output io.Writer, opts ArchiveOptions) error {
	// 确保源目录存在
	info, err := os.Stat(sourceDir)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Source directory not found: "+err.Error(),
		)
	}
	if !info.IsDir() {
		return utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Source path is not a directory",
		)
	}

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(output)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// 获取源目录的绝对路径（用于计算相对路径）
	absSource, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 递归遍历目录
	//nolint:wrapcheck // filepath.Walk 的错误已在回调中包装，此处直接返回
	return filepath.Walk(absSource, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk error at %s: %w", path, walkErr)
		}

		// 跳过目录本身（只处理内容）
		if path == absSource {
			return nil
		}

		// 处理符号链接
		if info.Mode()&os.ModeSymlink != 0 {
			info, err = handleSymlink(path, absSource, opts.FollowSymlinks)
			if err != nil {
				return err
			}
			if info == nil {
				return nil // 跳过符号链接
			}
		}

		// 计算在ZIP中的相对路径
		relPath, err := filepath.Rel(absSource, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		// 转换为ZIP路径格式（使用正斜杠）
		zipPath := filepath.ToSlash(relPath)

		// 处理目录
		if info.IsDir() {
			// ZIP中目录以斜杠结尾
			if !strings.HasSuffix(zipPath, "/") {
				zipPath += "/"
			}
			_, err := zipWriter.Create(zipPath)
			if err != nil {
				return fmt.Errorf("create zip dir entry: %w", err)
			}
			return nil
		}

		// 处理文件
		header, err := zipWriter.Create(zipPath)
		if err != nil {
			return fmt.Errorf("create zip file entry: %w", err)
		}

		// 复制文件内容
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open file %s: %w", path, err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		_, err = io.Copy(header, file)
		if err != nil {
			return fmt.Errorf("copy file content: %w", err)
		}
		return nil
	})
}

// ExtractZipToDirectory 将ZIP解压到目录
// 输入: ZIP数据, 目标目录路径
// 返回: 错误
//
//nolint:funlen,gocognit // 解压逻辑需要完整处理路径验证、目录创建和文件写入，复杂度较高
func ExtractZipToDirectory(zipData []byte, targetDir string) error {
	// 检查ZIP数据大小，防止解压缩炸弹
	if int64(len(zipData)) > maxTotalSize {
		return utils.NewCryptoError(
			utils.ErrInvalidParameter,
			fmt.Sprintf("ZIP data too large: %d bytes (max: %d)", len(zipData), maxTotalSize),
		)
	}

	// 创建目标目录 - 使用更安全的权限
	if err := os.MkdirAll(targetDir, dirPerm); err != nil {
		return fmt.Errorf("create target directory: %w", err)
	}

	// 创建ZIP读取器
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("invalid ZIP format: %w", err)
	}

	// 计算总大小用于防炸弹
	var totalSize int64

	// 遍历ZIP中的所有文件
	for _, file := range reader.File {
		// 防止路径遍历攻击 - 检查原始ZIP路径
		if strings.Contains(file.Name, "..") || strings.HasPrefix(file.Name, "/") || strings.HasPrefix(file.Name, "\\") {
			return utils.NewCryptoError(
				utils.ErrInvalidParameter,
				"Invalid file path in ZIP: "+file.Name,
			)
		}

		// 构建完整的目标路径
		// G305: 已在前面检查 ..、/、\ 并验证路径前缀，防止路径遍历
		targetPath := filepath.Join(targetDir, file.Name)

		// 再次验证最终路径 - 防止路径遍历
		absTarget, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("get absolute target path: %w", err)
		}
		absDir, err := filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("get absolute dir path: %w", err)
		}
		if !strings.HasPrefix(absTarget, absDir) {
			return utils.NewCryptoError(
				utils.ErrInvalidParameter,
				"Path traversal detected: "+file.Name,
			)
		}

		// 累加大小检查（防止解压缩炸弹）
		// G115: uint64 -> int64 转换，但有 maxTotalSize=1GB 限制，不会溢出
		size := file.UncompressedSize64
		if size > uint64(maxTotalSize) {
			return utils.NewCryptoError(
				utils.ErrInvalidParameter,
				fmt.Sprintf("File size exceeds limit: %d bytes", size),
			)
		}
		totalSize += int64(size)
		if totalSize > maxTotalSize {
			return utils.NewCryptoError(
				utils.ErrInvalidParameter,
				fmt.Sprintf("Total uncompressed size exceeds limit: %d bytes", totalSize),
			)
		}

		// 处理目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, dirPerm); err != nil {
				return fmt.Errorf("create directory %s: %w", targetPath, err)
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(targetPath), dirPerm); err != nil {
			return fmt.Errorf("create parent dir for %s: %w", targetPath, err)
		}

		// 打开ZIP中的文件
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("open zip entry %s: %w", file.Name, err)
		}
		defer func() {
			if closeErr := srcFile.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		// 创建目标文件 - 使用文件原始权限
		// G304: targetPath 已在前面验证在 targetDir 内
		dstFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, file.Mode())
		if err != nil {
			return fmt.Errorf("create output file %s: %w", targetPath, err)
		}
		defer func() {
			if closeErr := dstFile.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		// 复制内容
		// G110: 已通过 totalSize 检查限制总大小（maxTotalSize=1GB），防止解压缩炸弹
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return fmt.Errorf("copy content to %s: %w", targetPath, err)
		}
	}

	return nil
}

// GetZipSize 计算ZIP数据的总大小（用于进度条）.
func GetZipSize(zipData []byte) (int64, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return 0, fmt.Errorf("create ZIP reader: %w", err)
	}

	var totalSize int64
	for _, file := range reader.File {
		// G115: 显式转换，注意潜在的溢出
		// file.UncompressedSize64 是 uint64，转换为 int64
		// 检查是否超过 maxTotalSize (1GB)，防止溢出
		if file.UncompressedSize64 > uint64(maxTotalSize) {
			return 0, fmt.Errorf("file size exceeds maximum: %d bytes", file.UncompressedSize64)
		}
		totalSize += int64(file.UncompressedSize64)
	}

	return totalSize, nil
}

// CountZipFiles 统计ZIP中的文件数量.
func CountZipFiles(zipData []byte) (int, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return 0, fmt.Errorf("create ZIP reader: %w", err)
	}
	return len(reader.File), nil
}

// handleSymlink 处理符号链接，降低嵌套复杂度.
// 返回: 更新后的FileInfo，如果跳过则返回nil，错误表示验证失败.
func handleSymlink(path, absSource string, followSymlinks bool) (os.FileInfo, error) {
	if !followSymlinks {
		return nil, nil // 跳过符号链接
	}

	// 读取链接目标
	target, err := os.Readlink(path)
	if err != nil {
		return nil, fmt.Errorf("readlink %s: %w", path, err)
	}

	// 验证目标在源目录内，防止符号链接指向外部
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("abs path %s: %w", target, err)
	}
	if !strings.HasPrefix(absTarget, absSource) {
		return nil, fmt.Errorf("symlink target outside source directory: %s", target)
	}

	// 获取目标信息
	info, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("stat symlink target %s: %w", target, err)
	}

	return info, nil
}
