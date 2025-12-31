package zjcrypto

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSaveKeyFiles(t *testing.T) {
	// 生成密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	// 保存到临时目录
	tempDir := t.TempDir()
	pubPath := filepath.Join(tempDir, "test_public.pem")
	privPath := filepath.Join(tempDir, "test_private.pem")

	if err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		t.Fatalf("SaveKeyFiles failed: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(pubPath); os.IsNotExist(err) {
		t.Error("Public key file not created")
	}
	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		t.Error("Private key file not created")
	}

	// 验证私钥权限（仅在非 Windows 系统上）
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
}

func TestLoadKeyFiles(t *testing.T) {
	// 生成并保存密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	tempDir := t.TempDir()
	pubPath := filepath.Join(tempDir, "pub.pem")
	privPath := filepath.Join(tempDir, "priv.pem")

	if err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		t.Fatalf("SaveKeyFiles failed: %v", err)
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
		t.Error("Kyber public key mismatch")
	}

	kyberPrivBytes, _ := kyberPriv.MarshalBinary()
	loadedKyberPrivBytes, _ := loadedPriv.Kyber.MarshalBinary()
	if !bytes.Equal(kyberPrivBytes, loadedKyberPrivBytes) {
		t.Error("Kyber private key mismatch")
	}

	if !bytes.Equal(ecdhPub.Bytes(), loadedPub.ECDH.Bytes()) {
		t.Error("ECDH public key mismatch")
	}
	if !bytes.Equal(ecdhPriv.Bytes(), loadedPriv.ECDH.Bytes()) {
		t.Error("ECDH private key mismatch")
	}
}

func TestLoadMissingFiles(t *testing.T) {
	tempDir := t.TempDir()

	// 测试缺失公钥文件
	_, _, err := LoadKeyFiles(
		filepath.Join(tempDir, "missing_pub.pem"),
		filepath.Join(tempDir, "missing_priv.pem"),
	)
	if err == nil {
		t.Error("Should fail on missing public key file")
	}
	// 错误类型是 ErrInvalidParameter，不是 ErrFormatError
	// 只要返回错误即可，无需额外处理
}

func TestSaveToNonExistentDirectory(t *testing.T) {
	// 生成密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	// 尝试保存到不存在的目录
	nonExistentDir := "/tmp/nonexistent_test_dir_12345"
	pubPath := filepath.Join(nonExistentDir, "pub.pem")
	privPath := filepath.Join(nonExistentDir, "priv.pem")

	err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath)
	if err == nil {
		t.Error("Should fail when directory doesn't exist")
	}
}

func TestSaveKeyFilesContent(t *testing.T) {
	// 生成密钥
	kyberPub, kyberPriv, _ := GenerateKyberKeys()
	ecdhPub, ecdhPriv, _ := GenerateECDHKeys()

	tempDir := t.TempDir()
	pubPath := filepath.Join(tempDir, "test.pem")
	privPath := filepath.Join(tempDir, "test_priv.pem")

	if err := SaveKeyFiles(kyberPub, ecdhPub, kyberPriv, ecdhPriv, pubPath, privPath); err != nil {
		t.Fatalf("SaveKeyFiles failed: %v", err)
	}

	// 读取文件内容验证格式
	pubContent, err := os.ReadFile(pubPath)
	if err != nil {
		t.Fatalf("读取公钥文件失败: %v", err)
	}
	privContent, err := os.ReadFile(privPath)
	if err != nil {
		t.Fatalf("读取私钥文件失败: %v", err)
	}

	// 验证PEM格式
	if len(pubContent) == 0 {
		t.Error("Public key file is empty")
	}
	if len(privContent) == 0 {
		t.Error("Private key file is empty")
	}

	// 验证包含正确的PEM头
	expectedPubHeader := "-----BEGIN KYBER PUBLIC KEY-----"
	expectedPrivHeader := "-----BEGIN KYBER PRIVATE KEY-----"

	if !bytesContains(pubContent, []byte(expectedPubHeader)) {
		t.Errorf("Public key missing expected header: %s", expectedPubHeader)
	}
	if !bytesContains(privContent, []byte(expectedPrivHeader)) {
		t.Errorf("Private key missing expected header: %s", expectedPrivHeader)
	}
}

// 辅助函数.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func bytesContains(s, subs []byte) bool {
	if len(subs) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(subs); i++ {
		if bytesEqual(s[i:i+len(subs)], subs) {
			return true
		}
	}
	return false
}
