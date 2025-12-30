package crypto

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
)

// TestBufferPool 测试缓冲区池
func TestBufferPool(t *testing.T) {
	t.Run("基本功能", func(t *testing.T) {
		pool := NewBufferPool(64 * 1024)
		buf := pool.Get()
		if len(buf) != 64*1024 {
			t.Errorf("期望缓冲区大小 64KB, 实际 %d", len(buf))
		}
		pool.Put(buf)
	})

	t.Run("自动调整大小", func(t *testing.T) {
		// 小于最小值
		pool := NewBufferPool(1024)
		if pool.Get() == nil {
			t.Error("应该返回有效缓冲区")
		}

		// 大于最大值
		pool = NewBufferPool(10 * 1024 * 1024)
		buf := pool.Get()
		if cap(buf) > MaxBufferSize {
			t.Errorf("缓冲区大小超出限制: %d", cap(buf))
		}
	})

	t.Run("OptimalBufferSize", func(t *testing.T) {
		tests := []struct {
			size     int64
			expected int
		}{
			{5 * 1024 * 1024, 64 * 1024},          // 5MB -> 64KB
			{50 * 1024 * 1024, 256 * 1024},        // 50MB -> 256KB
			{500 * 1024 * 1024, 512 * 1024},       // 500MB -> 512KB
			{5 * 1024 * 1024 * 1024, 1024 * 1024}, // 5GB -> 1MB
		}

		for _, tt := range tests {
			result := OptimalBufferSize(tt.size)
			if result != tt.expected {
				t.Errorf("OptimalBufferSize(%d) = %d, 期望 %d", tt.size, result, tt.expected)
			}
		}
	})
}

// TestHashUtils 测试流式哈希工具
func TestHashUtils(t *testing.T) {
	t.Run("HashFile", func(t *testing.T) {
		// 创建临时文件
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")
		testData := []byte("Hello, World! This is a test file for streaming hash.")
		if err := os.WriteFile(testFile, testData, 0644); err != nil {
			t.Fatal(err)
		}

		hash, err := HashFile(testFile)
		if err != nil {
			t.Fatalf("HashFile failed: %v", err)
		}

		// 验证哈希不为空
		if hash == [32]byte{} {
			t.Error("Hash should not be empty")
		}
	})

	t.Run("HashReader", func(t *testing.T) {
		testData := []byte("Test data for reader hash")
		reader := bytes.NewReader(testData)
		hash, err := HashReader(reader)
		if err != nil {
			t.Fatalf("HashReader failed: %v", err)
		}

		if hash == [32]byte{} {
			t.Error("Hash should not be empty")
		}
	})

	t.Run("StreamingHash", func(t *testing.T) {
		hasher := NewStreamingHash()
		data := []byte("Test data")
		hasher.Write(data)
		result := hasher.Sum()

		if result == [32]byte{} {
			t.Error("StreamingHash should produce a result")
		}
	})
}

// TestKeyGenParallel 测试并行密钥生成
func TestKeyGenParallel(t *testing.T) {
	t.Run("GenerateHybridKeysParallel", func(t *testing.T) {
		kyberPub, kyberPriv, ecdhPub, ecdhPriv, err := GenerateHybridKeysParallel()
		if err != nil {
			t.Fatalf("GenerateHybridKeysParallel failed: %v", err)
		}

		if kyberPub == nil || kyberPriv == nil {
			t.Error("Kyber keys should not be nil")
		}
		if ecdhPub == nil || ecdhPriv == nil {
			t.Error("ECDH keys should not be nil")
		}
	})
}

// TestKeyCache 测试密钥缓存
func TestKeyCache(t *testing.T) {
	t.Run("缓存基本功能", func(t *testing.T) {
		// 先清空缓存
		ClearKeyCache()

		// 创建临时密钥文件
		tmpDir := t.TempDir()
		kyberPub, kyberPriv, ecdhPub, ecdhPriv, _ := GenerateHybridKeysParallel()
		pubPath := filepath.Join(tmpDir, "test_public.pem")
		privPath := filepath.Join(tmpDir, "test_private.pem")

		// 保存密钥
		if err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
			t.Fatal(err)
		}

		// 第一次加载（应该从文件读取）
		key1, err := LoadPublicKeyCached(pubPath)
		if err != nil {
			t.Fatalf("First load failed: %v", err)
		}

		// 第二次加载（应该从缓存返回）
		key2, err := LoadPublicKeyCached(pubPath)
		if err != nil {
			t.Fatalf("Second load failed: %v", err)
		}

		// 验证缓存命中
		if GetCacheSize() != 1 {
			t.Errorf("Expected cache size 1, got %d", GetCacheSize())
		}

		// 验证两个指针相同（缓存命中）
		if key1 != key2 {
			t.Error("Cache should return same pointer")
		}

		// 清空缓存
		ClearKeyCache()
		if GetCacheSize() != 0 {
			t.Errorf("Cache should be empty after clear, got %d", GetCacheSize())
		}
	})
}

// TestHeaderOptimized 测试优化的头部序列化
func TestHeaderOptimized(t *testing.T) {
	t.Run("MarshalBinaryOptimized", func(t *testing.T) {
		// 创建测试头部
		header := &format.FileHeader{
			Magic:       [4]byte{'F', 'Z', 'J', 0x01},
			Version:     0x0100,
			Algorithm:   0x02,
			Flags:       0x00,
			FilenameLen: 8,
			Filename:    "test.txt",
			FileSize:    1024,
			Timestamp:   1234567890,
			KyberEncLen: 1088,
			KyberEnc:    make([]byte, 1088),
			ECDHLen:     32,
			ECDHPub:     [32]byte{},
			IVLen:       12,
			IV:          [12]byte{},
			SigLen:      3293,
			Signature:   make([]byte, 3293),
			SHA256Hash:  [32]byte{},
		}

		// 填充随机数据
		rand.Read(header.KyberEnc)
		rand.Read(header.ECDHPub[:])
		rand.Read(header.IV[:])
		rand.Read(header.Signature)
		rand.Read(header.SHA256Hash[:])

		// 使用原方法
		original, err := header.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary failed: %v", err)
		}

		// 使用优化方法
		optimized, err := header.MarshalBinaryOptimized()
		if err != nil {
			t.Fatalf("MarshalBinaryOptimized failed: %v", err)
		}

		// 验证结果相同
		if len(original) != len(optimized) {
			t.Errorf("Length mismatch: original %d, optimized %d", len(original), len(optimized))
		}

		for i := range original {
			if original[i] != optimized[i] {
				t.Errorf("Byte %d mismatch: original %d, optimized %d", i, original[i], optimized[i])
				break
			}
		}
	})

	t.Run("优化效果验证", func(t *testing.T) {
		// 验证优化后的序列化使用预分配
		header := &format.FileHeader{
			Magic:       [4]byte{'F', 'Z', 'J', 0x01},
			Version:     0x0100,
			Algorithm:   0x02,
			FilenameLen: 8,
			Filename:    "test.txt",
			FileSize:    1024,
			Timestamp:   1234567890,
			KyberEncLen: 1088,
			KyberEnc:    make([]byte, 1088),
			ECDHLen:     32,
			ECDHPub:     [32]byte{},
			IVLen:       12,
			IV:          [12]byte{},
			SigLen:      3293,
			Signature:   make([]byte, 3293),
			SHA256Hash:  [32]byte{},
		}

		optimized, _ := header.MarshalBinaryOptimized()
		expectedSize := header.GetHeaderSize()

		if len(optimized) != expectedSize {
			t.Errorf("Optimized size %d != expected %d", len(optimized), expectedSize)
		}
	})
}

// TestStreamingEncryption 测试流式加密（完整流程）
func TestStreamingEncryption(t *testing.T) {
	// 跳过长耗时测试
	if testing.Short() {
		t.Skip("跳过流式加密长耗时测试")
	}

	t.Run("小文件流式加密解密", func(t *testing.T) {
		// 生成密钥
		kyberPub, kyberPriv, ecdhPub, ecdhPriv, err := GenerateHybridKeysParallel()
		if err != nil {
			t.Fatal(err)
		}

		dilithiumPub, dilithiumPriv, err := GenerateDilithiumKeys()
		if err != nil {
			t.Fatal(err)
		}

		// 创建临时文件
		tmpDir := t.TempDir()
		originalFile := filepath.Join(tmpDir, "original.txt")
		encryptedFile := filepath.Join(tmpDir, "encrypted.fzj")
		decryptedFile := filepath.Join(tmpDir, "decrypted.txt")

		testData := []byte("这是一个测试文件，用于验证流式加密功能。This is a test file for streaming encryption verification.")
		if err := os.WriteFile(originalFile, testData, 0644); err != nil {
			t.Fatal(err)
		}

		// 流式加密
		t.Logf("开始加密...")
		err = EncryptFileStreaming(originalFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv, 64*1024)
		if err != nil {
			t.Fatalf("Streaming encryption failed: %v", err)
		}
		t.Logf("加密完成")

		// 验证加密文件存在
		if _, err := os.Stat(encryptedFile); err != nil {
			t.Fatalf("Encrypted file not created: %v", err)
		}

		// 读取加密文件头部用于调试
		encData, _ := os.ReadFile(encryptedFile)
		t.Logf("加密文件大小: %d", len(encData))

		// 流式解密
		t.Logf("开始解密...")
		err = DecryptFileStreaming(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, dilithiumPub, 64*1024)
		if err != nil {
			t.Fatalf("Streaming decryption failed: %v", err)
		}
		t.Logf("解密完成")

		// 验证解密内容
		decryptedData, err := os.ReadFile(decryptedFile)
		if err != nil {
			t.Fatal(err)
		}

		if string(decryptedData) != string(testData) {
			t.Errorf("Decrypted data mismatch: got %d bytes, want %d bytes", len(decryptedData), len(testData))
			t.Logf("Original: %s", testData)
			t.Logf("Decrypted: %s", decryptedData)
		}
	})

	t.Run("自动缓冲区选择", func(t *testing.T) {
		// 生成密钥
		kyberPub, kyberPriv, ecdhPub, ecdhPriv, err := GenerateHybridKeysParallel()
		if err != nil {
			t.Fatal(err)
		}

		dilithiumPub, dilithiumPriv, err := GenerateDilithiumKeys()
		if err != nil {
			t.Fatal(err)
		}

		// 创建临时文件
		tmpDir := t.TempDir()
		originalFile := filepath.Join(tmpDir, "original.txt")
		encryptedFile := filepath.Join(tmpDir, "encrypted.fzj")
		decryptedFile := filepath.Join(tmpDir, "decrypted.txt")

		testData := make([]byte, 100*1024) // 100KB
		rand.Read(testData)
		if err := os.WriteFile(originalFile, testData, 0644); err != nil {
			t.Fatal(err)
		}

		// 使用自动缓冲区选择
		err = EncryptFileStreamingAuto(originalFile, encryptedFile, kyberPub, ecdhPub, dilithiumPriv)
		if err != nil {
			t.Fatalf("Auto streaming encryption failed: %v", err)
		}

		err = DecryptFileStreamingAuto(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, dilithiumPub)
		if err != nil {
			t.Fatalf("Auto streaming decryption failed: %v", err)
		}

		// 验证
		decryptedData, err := os.ReadFile(decryptedFile)
		if err != nil {
			t.Fatal(err)
		}

		if len(decryptedData) != len(testData) {
			t.Errorf("Length mismatch: %d vs %d", len(decryptedData), len(testData))
		}
	})
}

// TestStreamingUtils 测试流式工具
func TestStreamingUtils(t *testing.T) {
	t.Run("MultiWriter", func(t *testing.T) {
		// 创建两个 buffer
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		// 使用 MultiWriter
		mw := NewMultiWriter(buf1, buf2)
		data := []byte("Hello")
		mw.Write(data)

		// 验证两个 buffer 都写入了数据
		if buf1.String() != "Hello" || buf2.String() != "Hello" {
			t.Errorf("MultiWriter failed: buf1=%s, buf2=%s", buf1.String(), buf2.String())
		}
	})
}

// Benchmark 测试性能对比 (已迁移到 benchmark_test.go)
// 此函数保留用于兼容性，建议使用 benchmark_test.go 中的版本

func BenchmarkKeyGen(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = GenerateKyberKeys()
			_, _, _ = GenerateECDHKeys()
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _, _, _, _, _ = GenerateKeyPairParallel()
		}
	})
}
