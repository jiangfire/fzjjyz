package zjcrypto

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestCreateZipFromDirectory 测试目录打包成ZIP.
func TestCreateZipFromDirectory(t *testing.T) {
	// 创建临时测试目录
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	// 创建测试文件结构
	testDir := filepath.Join(tmpDir, "source")
	//nolint:gosec // 测试目录使用宽松权限
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	//nolint:gosec // 测试文件使用标准权限
	if err := os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	//nolint:gosec // 测试文件使用标准权限
	if err := os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}
	subDir := filepath.Join(testDir, "sub")
	//nolint:gosec // 测试目录使用宽松权限
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	//nolint:gosec // 测试文件使用标准权限
	if err := os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("content3"), 0644); err != nil {
		t.Fatalf("Failed to write file3: %v", err)
	}

	// 打包
	var buf bytes.Buffer
	err = CreateZipFromDirectory(testDir, &buf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包失败: %v", err)
	}

	// 验证ZIP数据不为空
	if buf.Len() == 0 {
		t.Error("ZIP数据为空")
	}

	// 验证可以读取ZIP
	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("读取ZIP失败: %v", err)
	}

	// 验证文件数量（应该有3个文件 + 1个子目录）
	if len(reader.File) < 3 {
		t.Errorf("期望至少3个文件，实际得到 %d", len(reader.File))
	}
}

// TestCreateZipFromDirectoryNotFound 测试源目录不存在.
func TestCreateZipFromDirectoryNotFound(t *testing.T) {
	var buf bytes.Buffer
	err := CreateZipFromDirectory("/nonexistent/path", &buf, DefaultArchiveOptions)
	if err == nil {
		t.Error("源目录不存在应该返回错误")
	}
}

// TestCreateZipFromDirectoryNotDir 测试源路径不是目录.
func TestCreateZipFromDirectoryNotDir(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "testfile_*")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("Warning: failed to remove temp file: %v", err)
		}
	}()
	if _, err := tmpFile.WriteString("test"); err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}

	var buf bytes.Buffer
	err = CreateZipFromDirectory(tmpFile.Name(), &buf, DefaultArchiveOptions)
	if err == nil {
		t.Error("源路径是文件应该返回错误")
	}
}

// TestExtractZipToDirectory 测试ZIP解压到目录.
func TestExtractZipToDirectory(t *testing.T) {
	// 先创建一个ZIP
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	// 创建源目录
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0750); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test content"), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// 打包
	var zipBuf bytes.Buffer
	err = CreateZipFromDirectory(sourceDir, &zipBuf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包失败: %v", err)
	}

	// 解压到新目录
	extractDir := filepath.Join(tmpDir, "extracted")
	err = ExtractZipToDirectory(zipBuf.Bytes(), extractDir)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证解压结果
	extractedFile := filepath.Join(extractDir, "test.txt")
	content, err := os.ReadFile(extractedFile) //nolint:gosec
	if err != nil {
		t.Fatalf("读取解压文件失败: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("内容不匹配，期望 'test content'，实际得到 '%s'", string(content))
	}
}

// TestExtractZipToDirectoryPathTraversal 测试路径遍历攻击防护.
func TestExtractZipToDirectoryPathTraversal(t *testing.T) {
	// 创建恶意ZIP（包含 .. 路径）
	var maliciousZip bytes.Buffer
	writer := zip.NewWriter(&maliciousZip)

	// 尝试写入到父目录
	header, err := writer.Create("../../etc/passwd")
	if err != nil {
		t.Fatalf("Failed to create header: %v", err)
	}
	if _, err := header.Write([]byte("malicious")); err != nil {
		t.Fatalf("Failed to write header: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// 尝试解压
	targetDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(targetDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	err = ExtractZipToDirectory(maliciousZip.Bytes(), filepath.Join(targetDir, "output"))
	if err == nil {
		t.Error("路径遍历攻击应该被阻止")
	}
}

// TestGetZipSize 测试计算ZIP大小.
func TestGetZipSize(t *testing.T) {
	// 创建测试目录
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	testDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(testDir, 0750); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "file.txt"), []byte("12345"), 0600); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	var buf bytes.Buffer
	err = CreateZipFromDirectory(testDir, &buf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包失败: %v", err)
	}

	size, err := GetZipSize(buf.Bytes())
	if err != nil {
		t.Fatalf("获取ZIP大小失败: %v", err)
	}

	if size <= 0 {
		t.Errorf("期望正数大小，实际得到 %d", size)
	}
}

// TestCountZipFiles 测试统计ZIP文件数量.
func TestCountZipFiles(t *testing.T) {
	// 创建测试目录
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	testDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(testDir, 0750); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("1"), 0600); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("2"), 0600); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(testDir, "sub"), 0750); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "sub", "file3.txt"), []byte("3"), 0600); err != nil {
		t.Fatalf("Failed to write file3: %v", err)
	}

	var buf bytes.Buffer
	err = CreateZipFromDirectory(testDir, &buf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包失败: %v", err)
	}

	count, err := CountZipFiles(buf.Bytes())
	if err != nil {
		t.Fatalf("统计文件数量失败: %v", err)
	}

	// 至少应该有3个文件
	if count < 3 {
		t.Errorf("期望至少3个文件，实际得到 %d", count)
	}
}

// TestExtractEmptyZip 测试解压空ZIP.
func TestExtractEmptyZip(t *testing.T) {
	// 创建空目录并打包
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	emptyDir := filepath.Join(tmpDir, "empty")
	if err := os.MkdirAll(emptyDir, 0750); err != nil {
		t.Fatalf("Failed to create empty dir: %v", err)
	}

	var buf bytes.Buffer
	err = CreateZipFromDirectory(emptyDir, &buf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包空目录失败: %v", err)
	}

	// 解压
	extractDir := filepath.Join(tmpDir, "extracted")
	err = ExtractZipToDirectory(buf.Bytes(), extractDir)
	if err != nil {
		t.Fatalf("解压空ZIP失败: %v", err)
	}

	// 验证目录存在
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		t.Error("解压后目录应该存在")
	}
}

// TestCreateZipWithSubdirectories 测试包含子目录的打包.
func TestCreateZipWithSubdirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	testDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(filepath.Join(testDir, "a", "b", "c"), 0750); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "a", "b", "c", "deep.txt"), []byte("deep"), 0600); err != nil {
		t.Fatalf("Failed to write deep file: %v", err)
	}

	var buf bytes.Buffer
	err = CreateZipFromDirectory(testDir, &buf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包失败: %v", err)
	}

	// 验证ZIP结构
	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("读取ZIP失败: %v", err)
	}

	// 查找深层文件
	found := false
	for _, f := range reader.File {
		if f.Name == "a/b/c/deep.txt" {
			found = true
			break
		}
	}

	if !found {
		t.Error("未找到深层文件")
	}
}

// TestExtractZipWithDirectories 测试解压包含目录的ZIP.
func TestExtractZipWithDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzj_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	// 创建带目录结构的源
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(filepath.Join(sourceDir, "dir1", "dir2"), 0750); err != nil {
		t.Fatalf("Failed to create source dirs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "root.txt"), []byte("root"), 0600); err != nil {
		t.Fatalf("Failed to write root.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "dir1", "file1.txt"), []byte("file1"), 0600); err != nil {
		t.Fatalf("Failed to write file1.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "dir1", "dir2", "file2.txt"), []byte("file2"), 0600); err != nil {
		t.Fatalf("Failed to write file2.txt: %v", err)
	}

	// 打包
	var zipBuf bytes.Buffer
	err = CreateZipFromDirectory(sourceDir, &zipBuf, DefaultArchiveOptions)
	if err != nil {
		t.Fatalf("打包失败: %v", err)
	}

	// 解压
	extractDir := filepath.Join(tmpDir, "extracted")
	err = ExtractZipToDirectory(zipBuf.Bytes(), extractDir)
	if err != nil {
		t.Fatalf("解压失败: %v", err)
	}

	// 验证所有文件
	tests := []struct {
		path     string
		expected string
	}{
		{"root.txt", "root"},
		{"dir1/file1.txt", "file1"},
		{"dir1/dir2/file2.txt", "file2"},
	}

	for _, tt := range tests {
		content, err := os.ReadFile(filepath.Join(extractDir, tt.path)) //nolint:gosec
		if err != nil {
			t.Errorf("读取 %s 失败: %v", tt.path, err)
			continue
		}
		if string(content) != tt.expected {
			t.Errorf("%s 内容不匹配，期望 '%s'，实际得到 '%s'", tt.path, tt.expected, string(content))
		}
	}
}
