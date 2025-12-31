package zjcrypto

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// TestHybridEncryption 测试混合加密的封装和解封装.
func TestHybridEncryption(t *testing.T) {
	// 生成密钥对
	kyberPub, kyberPriv, err := GenerateKyberKeys()
	if err != nil {
		t.Fatalf("生成 Kyber 密钥失败: %v", err)
	}

	ecdhPub, ecdhPriv, err := GenerateECDHKeys()
	if err != nil {
		t.Fatalf("生成 ECDH 密钥失败: %v", err)
	}

	// 加密 - 封装
	encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
	encapsulated, ecdhTempPub, sharedSecret, err := encryptor.Encapsulate()
	if err != nil {
		t.Fatalf("封装失败: %v", err)
	}

	// 验证封装结果长度
	if len(encapsulated) != 1088 {
		t.Errorf("期望封装长度 1088，实际得到 %d", len(encapsulated))
	}

	if len(ecdhTempPub) != 32 {
		t.Errorf("期望临时 ECDH 公钥长度 32，实际得到 %d", len(ecdhTempPub))
	}

	if len(sharedSecret) != 32 {
		t.Errorf("期望共享密钥长度 32，实际得到 %d", len(sharedSecret))
	}

	// 解密 - 解封装
	decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)
	decryptedSecret, err := decryptor.Decapsulate(encapsulated, ecdhTempPub)
	if err != nil {
		t.Fatalf("解封装失败: %v", err)
	}

	// 验证共享密钥一致
	if !bytes.Equal(sharedSecret, decryptedSecret) {
		t.Error("共享密钥不匹配")
	}
}

// TestAESGCMEncryption 测试 AES-GCM 加密解密.
func TestAESGCMEncryption(t *testing.T) {
	sharedSecret := make([]byte, 32)
	if _, err := rand.Read(sharedSecret); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}
	plaintext := []byte("Test data for encryption")

	// 加密
	ciphertext, nonce, err := AESGCMEncrypt(sharedSecret, plaintext)
	if err != nil {
		t.Fatalf("AES-GCM 加密失败: %v", err)
	}

	// 验证 nonce 长度
	if len(nonce) != 12 {
		t.Errorf("期望 nonce 长度 12，实际得到 %d", len(nonce))
	}

	// 解密
	decrypted, err := AESGCMDecrypt(sharedSecret, ciphertext, nonce)
	if err != nil {
		t.Fatalf("AES-GCM 解密失败: %v", err)
	}

	// 验证解密结果
	if !bytes.Equal(plaintext, decrypted) {
		t.Error("解密数据与原文不匹配")
	}
}

// TestEncryptionIntegrity 测试加密完整性验证（防篡改）.
func TestEncryptionIntegrity(t *testing.T) {
	sharedSecret := make([]byte, 32)
	if _, err := rand.Read(sharedSecret); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}
	plaintext := []byte("Sensitive data")

	ciphertext, nonce, err := AESGCMEncrypt(sharedSecret, plaintext)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 篡改密文
	ciphertext[0] ^= 0xFF

	// 应该失败
	_, err = AESGCMDecrypt(sharedSecret, ciphertext, nonce)
	if err == nil {
		t.Error("篡改后的密文应该解密失败")
	}
}

// TestHybridWithDifferentData 测试不同数据的混合加密.
func TestHybridWithDifferentData(t *testing.T) {
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	testCases := []struct {
		name string
		data []byte
	}{
		{"空数据", []byte{}},
		{"单字节", []byte{0x00}},
		{"大块数据", make([]byte, 1024*1024)}, // 1MB
		{"二进制数据", []byte{0x00, 0xFF, 0x80, 0x7F}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
			encapsulated, ecdhTempPub, sharedSecret, err := encryptor.Encapsulate()
			if err != nil {
				t.Fatalf("封装失败: %v", err)
			}

			decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)
			decryptedSecret, err := decryptor.Decapsulate(encapsulated, ecdhTempPub)
			if err != nil {
				t.Fatalf("解封装失败: %v", err)
			}

			if !bytes.Equal(sharedSecret, decryptedSecret) {
				t.Error("共享密钥不匹配")
			}

			// 使用共享密钥加密数据
			ciphertext, nonce, err := AESGCMEncrypt(sharedSecret, tc.data)
			if err != nil {
				t.Fatalf("数据加密失败: %v", err)
			}

			decrypted, err := AESGCMDecrypt(sharedSecret, ciphertext, nonce)
			if err != nil {
				t.Fatalf("数据解密失败: %v", err)
			}

			if !bytes.Equal(tc.data, decrypted) {
				t.Errorf("数据不匹配: 原文长度 %d, 解密长度 %d", len(tc.data), len(decrypted))
			}
		})
	}
}

// TestHybridEncryptorConstructor 测试构造函数.
func TestHybridEncryptorConstructor(t *testing.T) {
	kyberPub, _, _ := GenerateKyberKeys()
	ecdhPub, _, _ := GenerateECDHKeys()

	encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
	if encryptor == nil {
		t.Fatal("构造函数返回 nil")
	}

	if encryptor.kyberPub != kyberPub {
		t.Error("Kyber 公钥未正确设置")
	}

	if encryptor.ecdhPub != ecdhPub {
		t.Error("ECDH 公钥未正确设置")
	}
}

// TestHybridDecryptorConstructor 测试解密器构造函数.
func TestHybridDecryptorConstructor(t *testing.T) {
	_, kyberPriv, _ := GenerateKyberKeys()
	_, ecdhPriv, _ := GenerateECDHKeys()

	decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)
	if decryptor == nil {
		t.Fatal("构造函数返回 nil")
	}

	if decryptor.kyberPriv != kyberPriv {
		t.Error("Kyber 私钥未正确设置")
	}

	if decryptor.ecdhPriv != ecdhPriv {
		t.Error("ECDH 私钥未正确设置")
	}
}

// TestSharedSecretLength 测试共享密钥长度一致性.
func TestSharedSecretLength(t *testing.T) {
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
	decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)

	// 多次测试确保长度一致
	for i := 0; i < 10; i++ {
		encapsulated, ecdhTempPub, sharedSecret1, err := encryptor.Encapsulate()
		if err != nil {
			t.Fatalf("第 %d 次封装失败: %v", i, err)
		}

		sharedSecret2, err := decryptor.Decapsulate(encapsulated, ecdhTempPub)
		if err != nil {
			t.Fatalf("第 %d 次解封装失败: %v", i, err)
		}

		if len(sharedSecret1) != 32 || len(sharedSecret2) != 32 {
			t.Errorf("第 %d 次: 密钥长度错误 (1:%d, 2:%d)", i, len(sharedSecret1), len(sharedSecret2))
		}

		if !bytes.Equal(sharedSecret1, sharedSecret2) {
			t.Errorf("第 %d 次: 密钥不匹配", i)
		}
	}
}

// TestAESGCMNonceUniqueness 测试 Nonce 唯一性.
func TestAESGCMNonceUniqueness(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}
	data := []byte("test")

	nonces := make(map[string]bool)
	for i := 0; i < 100; i++ {
		_, nonce, err := AESGCMEncrypt(key, data)
		if err != nil {
			t.Fatalf("加密失败: %v", err)
		}

		nonceStr := string(nonce)
		if nonces[nonceStr] {
			t.Errorf("发现重复的 nonce: %x", nonce)
		}
		nonces[nonceStr] = true
	}
}

// TestAESGCMEmptyData 测试空数据加密.
func TestAESGCMEmptyData(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}

	ciphertext, nonce, err := AESGCMEncrypt(key, []byte{})
	if err != nil {
		t.Fatalf("空数据加密失败: %v", err)
	}

	decrypted, err := AESGCMDecrypt(key, ciphertext, nonce)
	if err != nil {
		t.Fatalf("空数据解密失败: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("期望空数据，实际得到长度 %d", len(decrypted))
	}
}

// TestAESGCMInvalidKey 测试无效密钥.
func TestAESGCMInvalidKey(t *testing.T) {
	// 密钥长度错误
	invalidKeys := [][]byte{
		make([]byte, 16), // AES-128 (不支持)
		make([]byte, 24), // AES-192 (不支持)
		make([]byte, 0),  // 空密钥
		make([]byte, 64), // 过长
	}

	data := []byte("test")
	for _, key := range invalidKeys {
		_, _, err := AESGCMEncrypt(key, data)
		if err == nil {
			t.Errorf("长度 %d 的密钥应该失败", len(key))
		}
	}
}

// TestAESGCMTamperedCiphertext 测试密文篡改检测.
func TestAESGCMTamperedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}
	data := []byte("sensitive data")

	ciphertext, nonce, err := AESGCMEncrypt(key, data)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 测试不同位置的篡改
	positions := []int{0, len(ciphertext) / 2, len(ciphertext) - 1}

	for _, pos := range positions {
		tampered := make([]byte, len(ciphertext))
		copy(tampered, ciphertext)
		tampered[pos] ^= 0xFF

		_, err := AESGCMDecrypt(key, tampered, nonce)
		if err == nil {
			t.Errorf("位置 %d 的篡改应该被检测到", pos)
		}
	}
}

// TestAESGCMTamperedNonce 测试 Nonce 篡改检测.
func TestAESGCMTamperedNonce(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}
	data := []byte("test")

	ciphertext, nonce, err := AESGCMEncrypt(key, data)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 篡改 nonce
	tamperedNonce := make([]byte, len(nonce))
	copy(tamperedNonce, nonce)
	tamperedNonce[0] ^= 0xFF

	_, err = AESGCMDecrypt(key, ciphertext, tamperedNonce)
	if err == nil {
		t.Error("篡改的 nonce 应该导致解密失败")
	}
}

// TestHybridConsistency 测试混合加密一致性.
func TestHybridConsistency(t *testing.T) {
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
	decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)

	// 相同输入应该产生不同的封装（因为随机种子）
	encapsulated1, ecdhTempPub1, secret1, _ := encryptor.Encapsulate()
	encapsulated2, ecdhTempPub2, secret2, _ := encryptor.Encapsulate()

	// 封装应该不同
	if bytes.Equal(encapsulated1, encapsulated2) {
		t.Error("相同输入的两次封装应该不同（随机性）")
	}

	// 临时 ECDH 公钥也应该不同
	if bytes.Equal(ecdhTempPub1, ecdhTempPub2) {
		t.Error("临时 ECDH 公钥应该不同")
	}

	// 共享密钥应该不同（因为随机种子不同）
	if bytes.Equal(secret1, secret2) {
		t.Error("不同随机种子应该产生不同共享密钥")
	}

	// 解封装应该正确恢复
	decrypted1, _ := decryptor.Decapsulate(encapsulated1, ecdhTempPub1)
	decrypted2, _ := decryptor.Decapsulate(encapsulated2, ecdhTempPub2)

	if !bytes.Equal(secret1, decrypted1) || !bytes.Equal(secret2, decrypted2) {
		t.Error("解封装未能正确恢复共享密钥")
	}
}
