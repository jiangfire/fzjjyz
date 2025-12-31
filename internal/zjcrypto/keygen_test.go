package zjcrypto

import (
	"bytes"
	"crypto/ecdh"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestKyberKeyGeneration(t *testing.T) {
	pub, priv, err := GenerateKyberKeys()
	if err != nil {
		t.Fatalf("Key generation failed: %v", err)
	}

	if pub == nil || priv == nil {
		t.Fatal("Generated keys should not be nil")
	}

	// 验证密钥可以序列化
	pubBytes, err := pub.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal public key: %v", err)
	}
	if len(pubBytes) != 1184 {
		t.Errorf("Expected Kyber768 public key length 1184, got %d", len(pubBytes))
	}
}

func TestECDHKeyGeneration(t *testing.T) {
	pub, priv, err := GenerateECDHKeys()
	if err != nil {
		t.Fatalf("ECDH key generation failed: %v", err)
	}

	if pub == nil || priv == nil {
		t.Fatal("Generated ECDH keys should not be nil")
	}

	// 验证密钥类型
	if pub.Curve() != ecdh.X25519() {
		t.Error("ECDH key should use X25519 curve")
	}
}

func TestHybridKeyExportImport(t *testing.T) {
	// 生成密钥对
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	// 导出到PEM
	pubPEM, err := ExportPublicKey(kyberPub, ecdhPub)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	privPEM, err := ExportPrivateKey(kyberPriv, ecdhPriv)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// 验证PEM格式
	if !bytes.Contains(pubPEM, []byte("KYBER PUBLIC KEY")) {
		t.Error("Public key PEM should contain KYBER PUBLIC KEY header")
	}
	if !bytes.Contains(privPEM, []byte("KYBER PRIVATE KEY")) {
		t.Error("Private key PEM should contain KYBER PRIVATE KEY header")
	}

	// 导入并验证
	importedPub, importedPriv, err := ImportKeys(pubPEM, privPEM)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// 验证密钥一致性 - 使用 MarshalBinary
	kyberPubBytes, _ := kyberPub.MarshalBinary()
	importedKyberPubBytes, _ := importedPub.Kyber.MarshalBinary()
	if !bytes.Equal(kyberPubBytes, importedKyberPubBytes) {
		t.Error("Kyber public key mismatch after import")
	}

	kyberPrivBytes, _ := kyberPriv.MarshalBinary()
	importedKyberPrivBytes, _ := importedPriv.Kyber.MarshalBinary()
	if !bytes.Equal(kyberPrivBytes, importedKyberPrivBytes) {
		t.Error("Kyber private key mismatch after import")
	}

	if !bytes.Equal(ecdhPub.Bytes(), importedPub.ECDH.Bytes()) {
		t.Error("ECDH public key mismatch after import")
	}
	if !bytes.Equal(ecdhPriv.Bytes(), importedPriv.ECDH.Bytes()) {
		t.Error("ECDH private key mismatch after import")
	}
}

func TestKeyValidation(t *testing.T) {
	// 测试无效PEM
	_, _, err := ImportKeys([]byte("invalid"), []byte("invalid"))
	if err == nil {
		t.Error("Should fail on invalid PEM")
	}

	// 测试空PEM
	_, _, err = ImportKeys([]byte{}, []byte{})
	if err == nil {
		t.Error("Should fail on empty PEM")
	}
}

func TestKeyFileSaveLoad(t *testing.T) {
	// 生成密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	// 保存到临时文件
	tempDir := t.TempDir()
	pubPath := filepath.Join(tempDir, "test_public.pem")
	privPath := filepath.Join(tempDir, "test_private.pem")

	err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath)
	if err != nil {
		t.Fatalf("SaveKeyFiles failed: %v", err)
	}

	// 验证私钥文件权限（仅在非 Windows 系统上）
	// Windows 的文件权限模型与 Unix 不同
	if runtime.GOOS != "windows" {
		info, err := os.Stat(privPath)
		if err != nil {
			t.Fatal("Cannot stat private key file")
		}
		perm := info.Mode().Perm()
		if perm != 0600 {
			t.Errorf("Private key permissions should be 0600, got %o", perm)
		}
	}

	// 加载密钥
	loadedPub, loadedPriv, err := LoadKeyFiles(pubPath, privPath)
	if err != nil {
		t.Fatalf("LoadKeyFiles failed: %v", err)
	}

	// 验证一致性
	kyberPubBytes, _ := kyberPub.MarshalBinary()
	loadedKyberPubBytes, _ := loadedPub.Kyber.MarshalBinary()
	if !bytes.Equal(kyberPubBytes, loadedKyberPubBytes) {
		t.Error("Loaded Kyber public key doesn't match original")
	}

	kyberPrivBytes, _ := kyberPriv.MarshalBinary()
	loadedKyberPrivBytes, _ := loadedPriv.Kyber.MarshalBinary()
	if !bytes.Equal(kyberPrivBytes, loadedKyberPrivBytes) {
		t.Error("Loaded Kyber private key doesn't match original")
	}

	if !bytes.Equal(ecdhPub.Bytes(), loadedPub.ECDH.Bytes()) {
		t.Error("Loaded ECDH public key doesn't match original")
	}
	if !bytes.Equal(ecdhPriv.Bytes(), loadedPriv.ECDH.Bytes()) {
		t.Error("Loaded ECDH private key doesn't match original")
	}
}

func TestKeyFileCorruption(t *testing.T) {
	// 测试损坏的密钥文件
	tempDir := t.TempDir()
	corruptPath := filepath.Join(tempDir, "corrupt.pem")
	if err := os.WriteFile(corruptPath, []byte("not a valid pem"), 0600); err != nil {
		t.Fatalf("创建损坏密钥文件失败: %v", err)
	}

	_, _, err := LoadKeyFiles(corruptPath, corruptPath)
	if err == nil {
		t.Error("Should fail on corrupted key file")
	}
}
