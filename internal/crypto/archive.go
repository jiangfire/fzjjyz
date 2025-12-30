package crypto

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

// ArchiveOptions 打包选项
type ArchiveOptions struct {
	IncludePatterns []string // 包含的文件模式（glob）
	ExcludePatterns []string // 排除的文件模式（glob）
	FollowSymlinks  bool     // 是否跟随符号链接
}

// DefaultArchiveOptions 默认打包选项
var DefaultArchiveOptions = ArchiveOptions{
	IncludePatterns: []string{"**/*"},
	ExcludePatterns: []string{},
	FollowSymlinks:  false,
}

// CreateZipFromDirectory 将目录打包成ZIP
// 输入: 源目录路径, 输出缓冲区, 打包选项
// 返回: 错误
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
		return err
	}

	// 递归遍历目录
	return filepath.Walk(absSource, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// 跳过目录本身（只处理内容）
		if path == absSource {
			return nil
		}

		// 检查是否为符号链接
		if info.Mode()&os.ModeSymlink != 0 {
			if !opts.FollowSymlinks {
				return nil // 跳过符号链接
			}
			// 读取链接目标
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			// 获取目标信息
			info, err = os.Stat(target)
			if err != nil {
				return err
			}
		}

		// 计算在ZIP中的相对路径
		relPath, err := filepath.Rel(absSource, path)
		if err != nil {
			return err
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
			return err
		}

		// 处理文件
		header, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		// 复制文件内容
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		_, err = io.Copy(header, file)
		return err
	})
}

// ExtractZipToDirectory 将ZIP解压到目录
// 输入: ZIP数据, 目标目录路径
// 返回: 错误
func ExtractZipToDirectory(zipData []byte, targetDir string) error {
	// 创建目标目录
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to create target directory: "+err.Error(),
		)
	}

	// 创建ZIP读取器
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Invalid ZIP format: "+err.Error(),
		)
	}

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
		targetPath := filepath.Join(targetDir, file.Name)

		// 再次验证最终路径
		absTarget, _ := filepath.Abs(targetPath)
		absDir, _ := filepath.Abs(targetDir)
		if !strings.HasPrefix(absTarget, absDir) {
			return utils.NewCryptoError(
				utils.ErrInvalidParameter,
				"Path traversal detected: "+file.Name,
			)
		}

		// 处理目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// 打开ZIP中的文件
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := srcFile.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		// 创建目标文件
		dstFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, file.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := dstFile.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()

		// 复制内容
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}

	return nil
}

// GetZipSize 计算ZIP数据的总大小（用于进度条）
func GetZipSize(zipData []byte) (int64, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return 0, err
	}

	var totalSize int64
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	return totalSize, nil
}

// CountZipFiles 统计ZIP中的文件数量
func CountZipFiles(zipData []byte) (int, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return 0, err
	}
	return len(reader.File), nil
}
