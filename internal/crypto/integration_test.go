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

// TestIntegrationEndToEnd 完整的端到端加密解密流程测试
//
//nolint:funlen // 端到端测试需要完整覆盖所有流程
func TestIntegrationEndToEnd(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-integration-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	// 步骤1: 生成密钥对
	t.Log("步骤1: 生成密钥对...")
	kyberPubRaw, kyberPrivRaw, err := GenerateKyberKeys()
	if err != nil {
		t.Fatalf("生成Kyber密钥失败: %v", err)
	}
	ecdhPub, ecdhPriv, err := GenerateECDHKeys()
	if err != nil {
		t.Fatalf("生成ECDH密钥失败: %v", err)
	}
	_, dilithiumPriv, err := GenerateDilithiumKeys()
	if err != nil {
		t.Fatalf("生成Dilithium密钥失败: %v", err)
	}

	kyberPub := kyberPubRaw
	kyberPriv := kyberPrivRaw

	// 步骤2: 创建测试文件
	t.Log("步骤2: 创建测试文件...")
	testData := []byte("这是一个完整的端到端测试文件，包含中文、English、1234567890!@#$%^&*()_+{}|:<>?[]\\;',./")
	originalFile := filepath.Join(tmpDir, "original.txt")
	if err := os.WriteFile(originalFile, testData, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 步骤3: 加密文件
	t.Log("步骤3: 加密文件...")
	encryptedFile := filepath.Join(tmpDir, "encrypted.fzj")
	err = EncryptFile(originalFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		t.Fatalf("文件加密失败: %v", err)
	}

	// 步骤4: 验证加密文件格式
	t.Log("步骤4: 验证加密文件格式...")
	encInfo, err := os.Stat(encryptedFile)
	if err != nil {
		t.Fatalf("获取加密文件信息失败: %v", err)
	}
	if encInfo.Size() <= int64(len(testData)) {
		t.Errorf("加密文件太小: %d bytes (原始 %d bytes)", encInfo.Size(), len(testData))
	}

	// 步骤5: 解析加密文件头
	t.Log("步骤5: 解析加密文件头...")
	f, err := os.Open(encryptedFile)
	if err != nil {
		t.Fatalf("打开加密文件失败: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Logf("警告: 关闭文件失败: %v", err)
		}
	}()

	header, err := format.ParseFileHeader(f)
	if err != nil {
		t.Fatalf("解析文件头失败: %v", err)
	}

	// 验证头部字段
	if !format.IsValidMagic(header.Magic[:]) {
		t.Error("魔数验证失败")
	}
	if !format.IsVersionSupported(header.Version) {
		t.Error("版本验证失败")
	}
	if header.Filename != "original.txt" {
		t.Errorf("文件名错误: 期望 'original.txt', 得到 '%s'", header.Filename)
	}
	if header.FileSize != uint64(len(testData)) {
		t.Errorf("文件大小错误: 期望 %d, 得到 %d", len(testData), header.FileSize)
	}
	if header.Algorithm != 0x02 {
		t.Errorf("算法错误: 期望 0x02, 得到 0x%02x", header.Algorithm)
	}

	// 步骤6: 解密文件
	t.Log("步骤6: 解密文件...")
	decryptedFile := filepath.Join(tmpDir, "decrypted.txt")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err != nil {
		t.Fatalf("文件解密失败: %v", err)
	}

	// 步骤7: 验证解密数据
	t.Log("步骤7: 验证解密数据...")
	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}

	if !bytes.Equal(testData, decryptedData) {
		t.Errorf("解密数据不匹配\n原始: %s\n解密: %s", testData, decryptedData)
	}

	t.Log("✅ 端到端测试成功完成！")
}

// TestIntegrationTamperDetection 篡改检测集成测试.
func TestIntegrationTamperDetection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-tamper-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	// 生成密钥
	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	// 创建并加密文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Tamper test data"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "test.enc")
	if err := EncryptFile(testFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv); err != nil {
		t.Fatalf("文件加密失败: %v", err)
	}

	// 篡改密文
	encryptedData, err := os.ReadFile(encryptedFile)
	if err != nil {
		t.Fatalf("读取加密文件失败: %v", err)
	}
	tamperedData := make([]byte, len(encryptedData))
	copy(tamperedData, encryptedData)
	tamperedData[200] ^= 0xFF // 篡改数据

	tamperedFile := filepath.Join(tmpDir, "tampered.enc")
	if err := os.WriteFile(tamperedFile, tamperedData, 0644); err != nil {
		t.Fatalf("创建篡改文件失败: %v", err)
	}

	// 尝试解密 - 应该失败
	decryptedFile := filepath.Join(tmpDir, "decrypted.txt")
	err = DecryptFile(tamperedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv))
	if err == nil {
		t.Error("篡改检测失败: 应该检测到篡改")
	}
	t.Logf("✅ 篡改检测成功: %v", err)
}

// TestIntegrationWrongKey 使用错误密钥解密测试.
func TestIntegrationWrongKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-wrongkey-*")
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

	// 创建并加密
	testFile := filepath.Join(tmpDir, "secret.txt")
	if err := os.WriteFile(testFile, []byte("Secret content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "secret.enc")
	if err := EncryptFile(testFile, encryptedFile, kyberPub1Typed, ecdhPub1, dilithiumPriv1); err != nil {
		t.Fatalf("文件加密失败: %v", err)
	}

	// 使用错误密钥解密
	decryptedFile := filepath.Join(tmpDir, "wrong.txt")
	err = DecryptFile(encryptedFile, decryptedFile, kyberPriv2, ecdhPriv2, DilithiumGetPublicKey(dilithiumPriv1))
	if err == nil {
		t.Error("使用错误密钥应该解密失败")
	}
	t.Logf("✅ 错误密钥检测成功: %v", err)
}

// TestIntegrationEmptyFile 空文件集成测试.
func TestIntegrationEmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-empty-*")
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

	// 空文件
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "empty.enc")
	if err := EncryptFile(emptyFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv); err != nil {
		t.Fatalf("空文件加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "empty_decrypted.txt")
	dilithiumPub := DilithiumGetPublicKey(dilithiumPriv)
	if err := DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, dilithiumPub); err != nil {
		t.Fatalf("空文件解密失败: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}
	if len(decryptedData) != 0 {
		t.Errorf("空文件解密后应为0字节，实际: %d", len(decryptedData))
	}
	t.Log("✅ 空文件测试成功")
}

// TestIntegrationLargeFile 大文件集成测试.
func TestIntegrationLargeFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-large-*")
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

	// 创建 500KB 随机数据
	largeData := make([]byte, 500*1024)
	if _, err := rand.Read(largeData); err != nil {
		t.Fatalf("生成随机数据失败: %v", err)
	}

	largeFile := filepath.Join(tmpDir, "large.bin")
	if err := os.WriteFile(largeFile, largeData, 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "large.enc")
	if err := EncryptFile(largeFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv); err != nil {
		t.Fatalf("大文件加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "large_decrypted.bin")
	dilithiumPub := DilithiumGetPublicKey(dilithiumPriv)
	if err := DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, dilithiumPub); err != nil {
		t.Fatalf("大文件解密失败: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}
	if !bytes.Equal(largeData, decryptedData) {
		t.Error("大文件解密数据不匹配")
	}
	t.Log("✅ 大文件测试成功")
}

// TestIntegrationMultipleFiles 多文件并发集成测试.
func TestIntegrationMultipleFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-multi-*")
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

	// 创建多个文件
	files := []struct {
		name string
		data []byte
	}{
		{"file1.txt", []byte("First file content")},
		{"file2.bin", []byte{0x00, 0x01, 0x02, 0xFF}},
		{"file3.txt", []byte("Third file with 中文")},
		{"file4.txt", []byte("")}, // 空文件
	}

	// 加密所有 files
	for _, f := range files {
		origPath := filepath.Join(tmpDir, f.name)
		if err := os.WriteFile(origPath, f.data, 0644); err != nil {
			t.Fatalf("创建文件 %s 失败: %v", f.name, err)
		}

		encPath := filepath.Join(tmpDir, f.name+".enc")
		if err := EncryptFile(origPath, encPath, kyberPub, ecdhPub, dilithiumPriv); err != nil {
			t.Fatalf("文件 %s 加密失败: %v", f.name, err)
		}

		decPath := filepath.Join(tmpDir, f.name+".dec")
		if err := DecryptFile(encPath, decPath, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv)); err != nil {
			t.Fatalf("文件 %s 解密失败: %v", f.name, err)
		}

		decData, err := os.ReadFile(decPath)
		if err != nil {
			t.Fatalf("读取解密文件 %s 失败: %v", f.name, err)
		}
		if !bytes.Equal(f.data, decData) {
			t.Errorf("文件 %s 数据不匹配", f.name)
		}
	}
	t.Log("✅ 多文件测试成功")
}

// TestIntegrationKeyPersistence 密钥持久化集成测试.
func TestIntegrationKeyPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-keys-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	// 生成并保存密钥
	kyberPubRaw, kyberPrivRaw, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	kyberPub := kyberPubRaw.(kem.PublicKey)
	kyberPriv := kyberPrivRaw.(kem.PrivateKey)

	pubPath := filepath.Join(tmpDir, "pub.pem")
	privPath := filepath.Join(tmpDir, "priv.pem")

	// 保存密钥
	if err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		t.Fatalf("保存密钥失败: %v", err)
	}

	// 加载密钥
	hybridPub, hybridPriv, err := LoadKeyFiles(pubPath, privPath)
	if err != nil {
		t.Fatalf("加载密钥失败: %v", err)
	}

	// 生成签名密钥（不保存）
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	// 使用加载的密钥进行加密解密
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Key persistence test"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	encryptedFile := filepath.Join(tmpDir, "test.enc")
	// 使用 HybridPublicKey.Kyber 和 HybridPublicKey.ECDH
	if err := EncryptFile(testFile, encryptedFile, hybridPub.Kyber, hybridPub.ECDH, dilithiumPriv); err != nil {
		t.Fatalf("使用加载的密钥加密失败: %v", err)
	}

	decryptedFile := filepath.Join(tmpDir, "test.dec")
	// 使用 HybridPrivateKey.Kyber 和 HybridPrivateKey.ECDH
	dilithiumPub := DilithiumGetPublicKey(dilithiumPriv)
	if err := DecryptFile(encryptedFile, decryptedFile, hybridPriv.Kyber, hybridPriv.ECDH, dilithiumPub); err != nil {
		t.Fatalf("使用加载的密钥解密失败: %v", err)
	}

	origData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("读取原始文件失败: %v", err)
	}
	decData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("读取解密文件失败: %v", err)
	}
	if !bytes.Equal(origData, decData) {
		t.Error("密钥持久化后数据不匹配")
	}
	t.Log("✅ 密钥持久化测试成功")
}

// TestIntegrationSpecialFilenames 特殊文件名集成测试.
func TestIntegrationSpecialFilenames(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-special-*")
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

	specialNames := []string{
		"test-file_v1.2.txt",
		"file with spaces.txt",
		"file-with-dashes.txt",
		"file_with_underscores.txt",
	}

	for _, name := range specialNames {
		origPath := filepath.Join(tmpDir, name)
		testData := []byte("Test content for " + name)
		if err := os.WriteFile(origPath, testData, 0644); err != nil {
			t.Fatalf("创建文件 %s 失败: %v", name, err)
		}

		encPath := filepath.Join(tmpDir, name+".enc")
		if err := EncryptFile(origPath, encPath, kyberPub, ecdhPub, dilithiumPriv); err != nil {
			t.Fatalf("文件 %s 加密失败: %v", name, err)
		}

		decPath := filepath.Join(tmpDir, name+".dec")
		if err := DecryptFile(encPath, decPath, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv)); err != nil {
			t.Fatalf("文件 %s 解密失败: %v", name, err)
		}

		decData, err := os.ReadFile(decPath)
		if err != nil {
			t.Fatalf("读取解密文件 %s 失败: %v", name, err)
		}
		if !bytes.Equal(testData, decData) {
			t.Errorf("特殊文件名 %s 数据不匹配", name)
		}
	}
	t.Log("✅ 特殊文件名测试成功")
}
