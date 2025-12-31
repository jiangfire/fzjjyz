package crypto

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"github.com/cloudflare/circl/kem"
)

// TestEncryptFile 测试完整文件加密流程.
func TestEncryptFile(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	// 生成密钥对
	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建测试文件
	originalFile := filepath.Join(tmpDir, "test.txt")
	originalData := []byte("Hello, this is a test file for encryption!")
	if err := os.WriteFile(originalFile, originalData, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 加密
	encryptedFile := filepath.Join(tmpDir, "test.txt.enc")
	err = EncryptFile(originalFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("文件加密失败: %v", err)
	}

	// 验证加密文件存在
	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		t.Fatal("加密文件未创建")
	}

	// 解密
	decryptedFile := filepath.Join(tmpDir, "test_decrypted.txt")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err != nil {
		t.Fatalf("文件解密失败: %v", err)
	}

	// 验证解密数据
	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}

	if !bytes.Equal(originalData, decryptedData) {
		t.Errorf("解密数据不匹配\n原始: %s\n解密: %s", originalData, decryptedData)
	}
}

// TestEncryptEmptyFile 测试空文件加密.
func TestEncryptEmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建空文件
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "empty.txt.enc")
	err = EncryptFile(emptyFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("空文件加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "empty_decrypted.txt")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err != nil {
		t.Fatalf("空文件解密失败: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}

	if len(decryptedData) != 0 {
		t.Errorf("期望空文件，实际得到长度 %d", len(decryptedData))
	}
}

// TestEncryptLargeFile 测试大文件加密.
func TestEncryptLargeFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建 100KB 测试数据
	largeData := make([]byte, 100*1024)
	if _, err := rand.Read(largeData); err != nil {
		t.Fatalf("生成随机数据失败: %v", err)
	}

	largeFile := filepath.Join(tmpDir, "large.bin")
	if err := os.WriteFile(largeFile, largeData, 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "large.bin.enc")
	err = EncryptFile(largeFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("大文件加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "large_decrypted.bin")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err != nil {
		t.Fatalf("大文件解密失败: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}

	if !bytes.Equal(largeData, decryptedData) {
		t.Error("大文件解密数据不匹配")
	}
}

// TestEncryptedFileFormat 测试加密文件格式.
//nolint:funlen // 测试函数需要完整验证所有文件头字段
func TestEncryptedFileFormat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, _, _ := GenerateKyberKeys()
	ecdhPub, _, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	testData := []byte("Format test")
	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "test.txt.enc")
	err = EncryptFile(testFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 读取加密文件并验证格式
	f, err := os.Open(encryptedFile)
	if err != nil {
		t.Fatalf("打开加密文件失败: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Logf("警告: 关闭文件失败: %v", err)
		}
	}()

	// 解析文件头
	header, err := format.ParseFileHeader(f)
	if err != nil {
		t.Fatalf("解析文件头失败: %v", err)
	}

	// 验证文件头字段
	if !format.IsValidMagic(header.Magic[:]) {
		t.Error("魔数验证失败")
	}

	if !format.IsVersionSupported(header.Version) {
		t.Error("版本验证失败")
	}

	if header.Filename != "test.txt" {
		t.Errorf("文件名错误: 期望 'test.txt', 得到 '%s'", header.Filename)
	}

	if header.FileSize != uint64(len(testData)) {
		t.Errorf("文件大小错误: 期望 %d, 得到 %d", len(testData), header.FileSize)
	}

	// 验证头部大小
	expectedSize := 9 + len(header.Filename) + 8 + 4 + 2 + 1088 + 1 + 32 + 1 + 12 + 2 + int(header.SigLen) + 32
	if header.GetHeaderSize() != expectedSize {
		t.Errorf("头部大小计算错误: 期望 %d, 得到 %d", expectedSize, header.GetHeaderSize())
	}
}

// TestTamperedEncryptedFile 测试篡改加密文件检测.
func TestTamperedEncryptedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建并加密文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Test data"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "test.txt.enc")
	err = EncryptFile(testFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 读取加密文件
	encryptedData, err := os.ReadFile(encryptedFile)
	if err != nil {
		t.Fatalf("读取加密文件失败: %v", err)
	}

	// 篡改密文部分（跳过头部，篡改数据）
	if len(encryptedData) > 100 {
		tamperedData := make([]byte, len(encryptedData))
		copy(tamperedData, encryptedData)
		tamperedData[100] ^= 0xFF // 篡改数据部分

		tamperedFile := filepath.Join(tmpDir, "tampered.enc")
		if err := os.WriteFile(tamperedFile, tamperedData, 0644); err != nil {
			t.Fatalf("创建篡改文件失败: %v", err)
		}

		// 尝试解密篡改文件
		decryptedFile := filepath.Join(tmpDir, "tampered_decrypted.txt")
		err = DecryptFile(tamperedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
		if err == nil {
			t.Error("应该检测到篡改并失败")
		}
	}
}

// TestInvalidKeyDecrypt 测试使用错误密钥解密.
func TestInvalidKeyDecrypt(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	// 加密密钥
	kyberPub1, _, _ := GenerateKyberKeys()
	ecdhPub1, _, _ := GenerateECDHKeys()
	_, dilithiumPriv1, _ := GenerateDilithiumKeys()

	// 解密密钥（不同）
	_, kyberPriv2, _ := GenerateKyberKeys()
	_, ecdhPriv2, _ := GenerateECDHKeys()

	kyberPub1Typed := kyberPub1.(kem.PublicKey)

	// 创建并加密文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Secret data"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "test.txt.enc")
	err = EncryptFile(testFile, encryptedFile, kyberPub1Typed, ecdhPub1, dilithiumPriv1)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 使用错误密钥尝试解密
	decryptedFile := filepath.Join(tmpDir, "wrong_decrypted.txt")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv2, ecdhPriv2, DilithiumGetPublicKey(dilithiumPriv1))
	if err == nil {
		t.Error("使用错误密钥应该解密失败")
	}
}

// TestBinaryFile 测试二进制文件.
func TestBinaryFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建二进制数据（包含所有可能的字节值）
	binaryData := make([]byte, 256)
	for i := 0; i < 256; i++ {
		binaryData[i] = byte(i)
	}

	binaryFile := filepath.Join(tmpDir, "binary.bin")
	if err := os.WriteFile(binaryFile, binaryData, 0644); err != nil {
		t.Fatalf("创建二进制文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "binary.bin.enc")
	err = EncryptFile(binaryFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("二进制文件加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "binary_decrypted.bin")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err != nil {
		t.Fatalf("二进制文件解密失败: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}

	if !bytes.Equal(binaryData, decryptedData) {
		t.Error("二进制文件解密数据不匹配")
	}
}

// TestFileWithSpecialChars 测试特殊字符文件名.
func TestFileWithSpecialChars(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 使用特殊字符文件名
	specialFile := filepath.Join(tmpDir, "test-file_v1.2.txt")
	if err := os.WriteFile(specialFile, []byte("Special chars test"), 0644); err != nil {
		t.Fatalf("创建特殊字符文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "test-file_v1.2.txt.enc")
	err = EncryptFile(specialFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("特殊字符文件加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "test-file_v1.2_decrypted.txt")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err != nil {
		t.Fatalf("特殊字符文件解密失败: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}

	expectedData := []byte("Special chars test")
	if !bytes.Equal(expectedData, decryptedData) {
		t.Error("特殊字符文件解密数据不匹配")
	}
}

// TestConcurrentEncrypt 测试并发加密.
func TestConcurrentEncrypt(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建多个测试文件
	numFiles := 5
	files := make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		file := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".txt")
		data := []byte("Test data " + string(rune('0'+i)))
		if err := os.WriteFile(file, data, 0644); err != nil {
			t.Fatalf("创建测试文件 %d 失败: %v", i, err)
		}
		files[i] = file
	}

	// 串行加密（并发测试需要更复杂的同步）
	encryptedFiles := make([]string, numFiles)
	for i, file := range files {
		encryptedFile := file + ".enc"
		err := EncryptFile(file, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
		if err != nil {
			t.Fatalf("文件 %d 加密失败: %v", i, err)
		}
		encryptedFiles[i] = encryptedFile
	}

	// 解密并验证
	for i, encFile := range encryptedFiles {
		decryptedFile := encFile + ".dec"
		err := DecryptFile(encFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
		if err != nil {
			t.Fatalf("文件 %d 解密失败: %v", i, err)
		}

		originalData, _ := os.ReadFile(files[i])
		decryptedData, _ := os.ReadFile(decryptedFile)
		if !bytes.Equal(originalData, decryptedData) {
			t.Errorf("文件 %d 数据不匹配", i)
		}
	}
}

// TestEncryptFileMetadata 测试加密后文件元数据.
func TestEncryptFileMetadata(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	kyberPubRaw, _, _ := GenerateKyberKeys()
	ecdhPub, _, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "metadata.txt")
	testData := []byte("Metadata test")
	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 获取原始文件信息
	origInfo, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("获取原始文件信息失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "metadata.txt.enc")
	err = EncryptFile(testFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 验证加密文件大小（应该大于原始文件 + 头部）
	encInfo, err := os.Stat(encryptedFile)
	if err != nil {
		t.Fatalf("获取加密文件信息失败: %v", err)
	}

	expectedMinSize := origInfo.Size() + 3500 // 头部约 3.6KB
	if encInfo.Size() < expectedMinSize {
		t.Errorf("加密文件太小: 期望至少 %d, 实际 %d", expectedMinSize, encInfo.Size())
	}

	// 验证加密文件不可读为文本
	encData, _ := os.ReadFile(encryptedFile)
	if bytes.Contains(encData, []byte("Metadata test")) {
		t.Error("加密文件包含明文")
	}
}
