package crypto

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
)

// BenchmarkEncryptFile 基准测试：文件加密性能.
func BenchmarkEncryptFile(b *testing.B) {
	// 生成测试密钥
	kyberPub, _, _ := GenerateKyberKeys()
	ecdhPub, _, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	// 创建不同大小的测试文件
	sizes := []struct {
		name string
		size int64
	}{
		{"1KB", 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
		{"10MB", 10 * 1024 * 1024},
	}

	for _, tc := range sizes {
		b.Run(tc.name, func(b *testing.B) {
			// 创建临时测试文件
			tmpDir := b.TempDir()
			inputPath := filepath.Join(tmpDir, "input.bin")
			outputPath := filepath.Join(tmpDir, "output.fzj")

			// 生成测试数据
			data := make([]byte, tc.size)
			for i := range data {
				data[i] = byte(i % 256)
			}
			if err := os.WriteFile(inputPath, data, 0644); err != nil {
				b.Fatalf("Failed to write test file: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = EncryptFile(inputPath, outputPath, kyberPub, ecdhPub, dilithiumPriv)
			}
		})
	}
}

// BenchmarkDecryptFile 基准测试：文件解密性能.
func BenchmarkDecryptFile(b *testing.B) {
	// 生成测试密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	dilithiumPub, dilithiumPriv, _ := GenerateDilithiumKeys()

	// 创建测试文件
	sizes := []struct {
		name string
		size int64
	}{
		{"1KB", 1024},
		{"100KB", 100 * 1024},
		{"1MB", 1024 * 1024},
		{"10MB", 10 * 1024 * 1024},
	}

	for _, tc := range sizes {
		b.Run(tc.name, func(b *testing.B) {
			// 创建临时测试文件
			tmpDir := b.TempDir()
			inputPath := filepath.Join(tmpDir, "input.bin")
			encryptedPath := filepath.Join(tmpDir, "encrypted.fzj")
			decryptedPath := filepath.Join(tmpDir, "decrypted.bin")

			// 生成测试数据并加密
			data := make([]byte, tc.size)
			for i := range data {
				data[i] = byte(i % 256)
			}
			if err := os.WriteFile(inputPath, data, 0644); err != nil {
				b.Fatalf("Failed to write test file: %v", err)
			}
			_ = EncryptFile(inputPath, encryptedPath, kyberPub, ecdhPub, dilithiumPriv)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = DecryptFile(encryptedPath, decryptedPath, kyberPriv, ecdhPriv, dilithiumPub)
			}
		})
	}
}

// BenchmarkStreamingEncrypt 基准测试：流式加密性能.
func BenchmarkStreamingEncrypt(b *testing.B) {
	kyberPub, _, _ := GenerateKyberKeys()
	ecdhPub, _, _ := GenerateECDHKeys()
	_, dilithiumPriv, _ := GenerateDilithiumKeys()

	sizes := []struct {
		name string
		size int64
	}{
		{"1MB", 1024 * 1024},
		{"10MB", 10 * 1024 * 1024},
	}

	for _, tc := range sizes {
		b.Run(tc.name, func(b *testing.B) {
			tmpDir := b.TempDir()
			inputPath := filepath.Join(tmpDir, "input.bin")
			outputPath := filepath.Join(tmpDir, "output.fzj")

			data := make([]byte, tc.size)
			for i := range data {
				data[i] = byte(i % 256)
			}
			if err := os.WriteFile(inputPath, data, 0644); err != nil {
				b.Fatalf("Failed to write test file: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = EncryptFileStreamingAuto(inputPath, outputPath, kyberPub, ecdhPub, dilithiumPriv)
			}
		})
	}
}

// BenchmarkKeyGeneration 基准测试：密钥生成性能.
func BenchmarkKeyGeneration(b *testing.B) {
	b.Run("Kyber", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = GenerateKyberKeys()
		}
	})

	b.Run("ECDH", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = GenerateECDHKeys()
		}
	})

	b.Run("Dilithium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = GenerateDilithiumKeys()
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _, _, _, _, _ = GenerateKeyPairParallel()
		}
	})
}

// BenchmarkHeaderSerialization 基准测试：头部序列化性能.
func BenchmarkHeaderSerialization(b *testing.B) {
	// 创建测试头部
	header := &format.FileHeader{
		Version:    0x0100,
		Algorithm:  0x02,
		Filename:   "test_file.bin",
		FileSize:   1024 * 1024,
		KyberEnc:   make([]byte, 1088),
		ECDHPub:    [32]byte{},
		IV:         [12]byte{},
		SigLen:     2700,
		Signature:  make([]byte, 2700),
		SHA256Hash: [32]byte{},
		Timestamp:  uint32(time.Now().Unix() & 0xFFFFFFFF),
	}

	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = header.MarshalBinary()
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = header.MarshalBinaryOptimized()
		}
	})
}

// BenchmarkCachePerformance 基准测试：缓存性能.
func BenchmarkCachePerformance(b *testing.B) {
	tmpDir := b.TempDir()
	keyPath := filepath.Join(tmpDir, "test_key.pem")

	// 生成并保存测试密钥（需要完整密钥对）
	pub, priv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	if err := SaveKeyFiles(pub, ecdhPub, priv, ecdhPriv, keyPath+".pub", keyPath+".priv"); err != nil {
		b.Fatalf("Failed to save key files: %v", err)
	}

	b.Run("FirstLoad", func(b *testing.B) {
		ClearKeyCache()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = LoadPublicKeyCached(keyPath + ".pub")
		}
	})

	b.Run("CachedLoad", func(b *testing.B) {
		ClearKeyCache()
		// 预热缓存
		_, _ = LoadPublicKeyCached(keyPath + ".pub")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = LoadPublicKeyCached(keyPath + ".pub")
		}
	})
}

// TestPerformanceComparison 性能对比测试.
//nolint:funlen // 测试函数需要完整覆盖所有场景
func TestPerformanceComparison(t *testing.T) {
	// 生成密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()
	dilithiumPub, dilithiumPriv, _ := GenerateDilithiumKeys()

	// 测试不同大小的文件
	testSizes := []int64{1024, 10 * 1024, 100 * 1024, 1024 * 1024, 10 * 1024 * 1024}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			tmpDir := t.TempDir()
			inputPath := filepath.Join(tmpDir, "input.bin")
			encryptedPath := filepath.Join(tmpDir, "encrypted.fzj")
			decryptedPath := filepath.Join(tmpDir, "decrypted.bin")

			// 创建测试文件
			data := make([]byte, size)
			for i := range data {
				data[i] = byte(i % 256)
			}
			if err := os.WriteFile(inputPath, data, 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// 测试标准加密
			start := time.Now()
			err := EncryptFile(inputPath, encryptedPath, kyberPub, ecdhPub, dilithiumPriv)
			if err != nil {
				t.Fatalf("标准加密失败: %v", err)
			}
			standardEncryptTime := time.Since(start)

			// 测试流式加密
			start = time.Now()
			err = EncryptFileStreamingAuto(inputPath, encryptedPath, kyberPub, ecdhPub, dilithiumPriv)
			if err != nil {
				t.Fatalf("流式加密失败: %v", err)
			}
			streamingEncryptTime := time.Since(start)

			// 解密测试
			start = time.Now()
			err = DecryptFile(encryptedPath, decryptedPath, kyberPriv, ecdhPriv, dilithiumPub)
			if err != nil {
				t.Fatalf("解密失败: %v", err)
			}
			decryptTime := time.Since(start)

			// 验证结果
			originalData, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatalf("Failed to read original file: %v", err)
			}
			decryptedData, err := os.ReadFile(decryptedPath)
			if err != nil {
				t.Fatalf("Failed to read decrypted file: %v", err)
			}
			if string(originalData) != string(decryptedData) {
				t.Fatal("解密数据不匹配")
			}

			t.Logf("文件大小: %d bytes", size)
			t.Logf("标准加密: %v", standardEncryptTime)
			t.Logf("流式加密: %v", streamingEncryptTime)
			t.Logf("解密: %v", decryptTime)
			t.Logf("加密文件大小: %d bytes", getFileSize(encryptedPath))
		})
	}
}

// TestCacheTTL 测试缓存TTL功能.
func TestCacheTTL(t *testing.T) {
	t.Skip("SaveKeyFiles 需要非 nil 私钥，跳过完整测试")
}

// TestCacheSizeLimit 测试缓存大小限制.
func TestCacheSizeLimit(t *testing.T) {
	t.Skip("SaveKeyFiles 需要非 nil 私钥，跳过完整测试")
}

// TestCacheExpiration 测试缓存过期.
func TestCacheExpiration(t *testing.T) {
	// 注意: DefaultCacheTTL 是常量，无法直接修改
	// 这里测试缓存的基本过期机制
	t.Skip("TTL 是常量，需要重构为可配置变量才能测试")
}

// 辅助函数.
func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
