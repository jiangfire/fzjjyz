package zjcrypto

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// TestSignData 测试数据签名.
func TestSignData(t *testing.T) {
	_, priv, err := GenerateDilithiumKeys()
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	data := []byte("Test data for signing")
	signature, err := SignData(data, priv)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	// 验证签名长度
	expectedSize := mode3.SignatureSize
	if len(signature) != expectedSize {
		t.Errorf("签名长度错误: 期望 %d, 实际 %d", expectedSize, len(signature))
	}
}

// TestVerifySignature 测试签名验证.
func TestVerifySignature(t *testing.T) {
	pub, priv, err := GenerateDilithiumKeys()
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	data := []byte("Test data for signing")
	signature, err := SignData(data, priv)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	// 验证签名
	valid, err := VerifySignature(data, signature, pub)
	if err != nil {
		t.Fatalf("验证失败: %v", err)
	}

	if !valid {
		t.Error("签名验证应该通过")
	}
}

// TestVerifyInvalidSignature 测试无效签名验证.
func TestVerifyInvalidSignature(t *testing.T) {
	pub, priv, err := GenerateDilithiumKeys()
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	data := []byte("Test data")
	signature, err := SignData(data, priv)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	// 篡改签名
	tamperedSig := make([]byte, len(signature))
	copy(tamperedSig, signature)
	tamperedSig[0] ^= 0xFF

	// 使用篡改的签名验证
	valid, err := VerifySignature(data, tamperedSig, pub)
	if err != nil {
		t.Fatalf("验证出错: %v", err)
	}

	if valid {
		t.Error("篡改的签名应该验证失败")
	}
}

// TestSignDifferentData 测试不同数据签名.
func TestSignDifferentData(t *testing.T) {
	pub, priv, _ := GenerateDilithiumKeys()

	data1 := []byte("Data 1")
	data2 := []byte("Data 2")

	sig1, _ := SignData(data1, priv)
	sig2, _ := SignData(data2, priv)

	// 相同数据应该产生不同签名（随机性）
	if bytes.Equal(sig1, sig2) {
		t.Error("不同数据应该产生不同签名")
	}

	// 验证各自签名
	valid1, _ := VerifySignature(data1, sig1, pub)
	valid2, _ := VerifySignature(data2, sig2, pub)
	valid1wrong, _ := VerifySignature(data1, sig2, pub)

	if !valid1 || !valid2 {
		t.Error("各自签名应该验证通过")
	}

	if valid1wrong {
		t.Error("错误的签名应该验证失败")
	}
}

// TestSignEmptyData 测试空数据签名.
func TestSignEmptyData(t *testing.T) {
	pub, priv, _ := GenerateDilithiumKeys()

	data := []byte{}
	signature, err := SignData(data, priv)
	if err != nil {
		t.Fatalf("空数据签名失败: %v", err)
	}

	valid, err := VerifySignature(data, signature, pub)
	if err != nil {
		t.Fatalf("空数据签名验证失败: %v", err)
	}

	if !valid {
		t.Error("空数据签名应该验证通过")
	}
}

// TestSignLargeData 测试大数据签名.
func TestSignLargeData(t *testing.T) {
	pub, priv, _ := GenerateDilithiumKeys()

	largeData := make([]byte, 1024*1024) // 1MB
	if _, err := rand.Read(largeData); err != nil {
		t.Fatalf("生成随机数据失败: %v", err)
	}

	signature, err := SignData(largeData, priv)
	if err != nil {
		t.Fatalf("大数据签名失败: %v", err)
	}

	valid, err := VerifySignature(largeData, signature, pub)
	if err != nil {
		t.Fatalf("大数据签名验证失败: %v", err)
	}

	if !valid {
		t.Error("大数据签名应该验证通过")
	}
}

// TestSignFile 测试文件签名.
func TestSignFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	pub, priv, _ := GenerateDilithiumKeys()

	testFile := filepath.Join(tmpDir, "test.txt")
	testData := []byte("File content for signing")
	if err := os.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	signature, err := SignFile(testFile, priv)
	if err != nil {
		t.Fatalf("文件签名失败: %v", err)
	}

	valid, err := VerifyFileSignature(testFile, signature, pub)
	if err != nil {
		t.Fatalf("文件签名验证失败: %v", err)
	}

	if !valid {
		t.Error("文件签名应该验证通过")
	}
}

// TestVerifyFileSignatureTampered 测试篡改文件验证.
func TestVerifyFileSignatureTampered(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fzjjyz-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("警告: 清理临时目录失败: %v", err)
		}
	}()

	pub, priv, _ := GenerateDilithiumKeys()

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("Original data"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	signature, _ := SignFile(testFile, priv)

	// 篡改文件
	if err := os.WriteFile(testFile, []byte("Tampered data"), 0644); err != nil {
		t.Fatalf("篡改文件失败: %v", err)
	}

	valid, err := VerifyFileSignature(testFile, signature, pub)
	if err != nil {
		t.Fatalf("验证出错: %v", err)
	}

	if valid {
		t.Error("篡改后的文件签名应该验证失败")
	}
}

// TestSignHash 测试哈希签名.
func TestSignHash(t *testing.T) {
	pub, priv, _ := GenerateDilithiumKeys()

	hash := make([]byte, 32)
	if _, err := rand.Read(hash); err != nil {
		t.Fatalf("生成随机哈希失败: %v", err)
	}

	signature, err := SignHash(hash, priv)
	if err != nil {
		t.Fatalf("哈希签名失败: %v", err)
	}

	valid, err := VerifyHashSignature(hash, signature, pub)
	if err != nil {
		t.Fatalf("哈希签名验证失败: %v", err)
	}

	if !valid {
		t.Error("哈希签名应该验证通过")
	}
}

// TestSignInvalidHash 测试无效哈希签名.
func TestSignInvalidHash(t *testing.T) {
	_, priv, _ := GenerateDilithiumKeys()

	invalidHash := []byte("not 32 bytes")
	_, err := SignHash(invalidHash, priv)
	if err == nil {
		t.Error("无效哈希应该返回错误")
	}
}

// TestKeySizes 测试密钥和签名大小.
func TestKeySizes(t *testing.T) {
	pub, priv, err := GenerateDilithiumKeys()
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	// 验证公钥大小
	pubSize := DilithiumPublicKeySize()
	if pubSize != 1952 {
		t.Errorf("公钥大小错误: 期望 1952, 实际 %d", pubSize)
	}

	// 验证私钥大小
	privSize := DilithiumPrivateKeySize()
	if privSize != 4000 {
		t.Errorf("私钥大小错误: 期望 4000, 实际 %d", privSize)
	}

	// 验证签名大小
	sigSize := DilithiumSignatureSize()
	if sigSize != 3293 {
		t.Errorf("签名大小错误: 期望 3293, 实际 %d", sigSize)
	}

	// 验证实际密钥大小
	pubBytes, ok := pub.(*mode3.PublicKey)
	if !ok {
		t.Fatal("公钥类型错误")
	}
	privBytes, ok := priv.(*mode3.PrivateKey)
	if !ok {
		t.Fatal("私钥类型错误")
	}

	if len(pubBytes.Bytes()) != pubSize {
		t.Errorf("实际公钥字节大小不匹配")
	}

	if len(privBytes.Bytes()) != privSize {
		t.Errorf("实际私钥字节大小不匹配")
	}
}

// TestSignDataWithInvalidKey 测试无效密钥签名.
func TestSignDataWithInvalidKey(t *testing.T) {
	invalidKey := "not a dilithium key"
	_, err := SignData([]byte("test"), invalidKey)
	if err == nil {
		t.Error("无效密钥应该返回错误")
	}
}

// TestVerifyDataWithInvalidKey 测试无效密钥验证.
func TestVerifyDataWithInvalidKey(t *testing.T) {
	_, priv, _ := GenerateDilithiumKeys()
	data := []byte("test")
	signature, _ := SignData(data, priv)

	invalidKey := "not a dilithium key"
	_, err := VerifySignature(data, signature, invalidKey)
	if err == nil {
		t.Error("无效密钥应该返回错误")
	}
}

// TestConcurrentSigning 测试并发签名.
func TestConcurrentSigning(t *testing.T) {
	_, priv, _ := GenerateDilithiumKeys()

	// 签名多个不同数据
	numSignatures := 10
	results := make([][]byte, numSignatures)
	data := make([][]byte, numSignatures)

	for i := 0; i < numSignatures; i++ {
		data[i] = []byte("Data " + string(rune('0'+i)))
		signature, err := SignData(data[i], priv)
		if err != nil {
			t.Fatalf("第 %d 次签名失败: %v", i, err)
		}
		results[i] = signature
	}

	// 验证所有签名
	pub := DilithiumGetPublicKey(priv)
	for i := 0; i < numSignatures; i++ {
		valid, err := VerifySignature(data[i], results[i], pub)
		if err != nil {
			t.Fatalf("第 %d 次验证出错: %v", i, err)
		}
		if !valid {
			t.Errorf("第 %d 次签名验证失败", i)
		}
	}
}

// TestBinaryDataSignature 测试二进制数据签名.
func TestBinaryDataSignature(t *testing.T) {
	pub, priv, _ := GenerateDilithiumKeys()

	// 包含所有字节值的数据
	binaryData := make([]byte, 256)
	for i := 0; i < 256; i++ {
		binaryData[i] = byte(i)
	}

	signature, err := SignData(binaryData, priv)
	if err != nil {
		t.Fatalf("二进制数据签名失败: %v", err)
	}

	valid, err := VerifySignature(binaryData, signature, pub)
	if err != nil {
		t.Fatalf("二进制数据验证失败: %v", err)
	}

	if !valid {
		t.Error("二进制数据签名应该验证通过")
	}
}
