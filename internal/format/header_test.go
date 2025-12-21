package format

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestFileHeaderSerialization(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x01,
		FilenameLen: 8,  // 修正：test.txt 是 8 字节
		Filename:    "test.txt",
		FileSize:    1024,
		Timestamp:   1734672000,
		KyberEncLen: 1088,
		KyberEnc:    make([]byte, 1088),
		ECDHLen:     32,
		ECDHPub:     [32]byte{},
		IVLen:       12,
		IV:          [12]byte{},
		SigLen:      2420,
		Signature:   make([]byte, 2420),
		SHA256Hash:  [32]byte{},
	}

	// 填充测试数据
	rand.Read(header.KyberEnc)
	rand.Read(header.ECDHPub[:])
	rand.Read(header.IV[:])
	rand.Read(header.Signature)
	rand.Read(header.SHA256Hash[:])

	// 序列化
	data, err := header.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	// 验证大小（包含Kyber和Dilithium签名，应在4KB以内）
	if len(data) > 4096 {
		t.Errorf("Header too large: %d bytes (should be < 4096)", len(data))
	}

	// 反序列化
	parsed := &FileHeader{}
	err = parsed.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}

	// 验证一致性
	if !bytes.Equal(header.Magic[:], parsed.Magic[:]) {
		t.Error("Magic mismatch")
	}
	if header.FileSize != parsed.FileSize {
		t.Errorf("FileSize mismatch: %d != %d", header.FileSize, parsed.FileSize)
	}
	if header.Filename != parsed.Filename {
		t.Errorf("Filename mismatch: %s != %s", header.Filename, parsed.Filename)
	}
	if header.Timestamp != parsed.Timestamp {
		t.Errorf("Timestamp mismatch: %d != %d", header.Timestamp, parsed.Timestamp)
	}
	if !bytes.Equal(header.KyberEnc, parsed.KyberEnc) {
		t.Error("KyberEnc mismatch")
	}
	if !bytes.Equal(header.ECDHPub[:], parsed.ECDHPub[:]) {
		t.Error("ECDHPub mismatch")
	}
	if !bytes.Equal(header.IV[:], parsed.IV[:]) {
		t.Error("IV mismatch")
	}
	if !bytes.Equal(header.Signature, parsed.Signature) {
		t.Error("Signature mismatch")
	}
	if !bytes.Equal(header.SHA256Hash[:], parsed.SHA256Hash[:]) {
		t.Error("SHA256Hash mismatch")
	}
}

func TestMagicValidation(t *testing.T) {
	validMagic := []byte{'F', 'Z', 'J', 0x01}
	if !IsValidMagic(validMagic) {
		t.Error("Valid magic should pass")
	}

	invalidMagic := []byte{'F', 'X', 'J', 0x01}
	if IsValidMagic(invalidMagic) {
		t.Error("Invalid magic should fail")
	}

	// 测试长度不足
	shortMagic := []byte{'F', 'Z', 'J'}
	if IsValidMagic(shortMagic) {
		t.Error("Short magic should fail")
	}
}

func TestVersionCompatibility(t *testing.T) {
	// 测试版本兼容性
	if !IsVersionSupported(0x0100) {
		t.Error("Version 0x0100 should be supported")
	}

	if IsVersionSupported(0x0200) {
		t.Error("Future version should not be supported yet")
	}

	if IsVersionSupported(0x0001) {
		t.Error("Old version should not be supported")
	}
}

func TestNewFileHeader(t *testing.T) {
	filename := "document.pdf"
	fileSize := uint64(2048)
	kyberEnc := make([]byte, 1088)
	ecdhPub := [32]byte{}
	iv := [12]byte{}
	signature := make([]byte, 2420)
	hash := [32]byte{}

	// 填充随机数据
	rand.Read(kyberEnc)
	rand.Read(ecdhPub[:])
	rand.Read(iv[:])
	rand.Read(signature)
	rand.Read(hash[:])

	header := NewFileHeader(filename, fileSize, kyberEnc, ecdhPub, iv, signature, hash)

	// 验证固定字段
	if header.Magic != [4]byte{'F', 'Z', 'J', 0x01} {
		t.Error("Magic incorrect")
	}
	if header.Version != 0x0100 {
		t.Error("Version incorrect")
	}
	if header.Algorithm != 0x02 {
		t.Error("Algorithm incorrect")
	}
	if header.Filename != filename {
		t.Errorf("Filename incorrect: %s", header.Filename)
	}
	if header.FileSize != fileSize {
		t.Errorf("FileSize incorrect: %d", header.FileSize)
	}
	if header.ECDHLen != 32 {
		t.Error("ECDHLen should be 32")
	}
	if header.IVLen != 12 {
		t.Error("IVLen should be 12")
	}

	// 验证时间戳（应该在合理范围内）
	if header.Timestamp == 0 {
		t.Error("Timestamp should not be zero")
	}
}

func TestHeaderWithEmptyFields(t *testing.T) {
	// 测试空字段的序列化
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x00,
		FilenameLen: 0,
		Filename:    "",
		FileSize:    0,
		Timestamp:   0,
		KyberEncLen: 0,
		KyberEnc:    nil,
		ECDHLen:     0,
		ECDHPub:     [32]byte{},
		IVLen:       0,
		IV:          [12]byte{},
		SigLen:      0,
		Signature:   nil,
		SHA256Hash:  [32]byte{},
	}

	data, err := header.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary with empty fields failed: %v", err)
	}

	// 反序列化
	parsed := &FileHeader{}
	err = parsed.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary with empty fields failed: %v", err)
	}

	// 验证
	if parsed.Filename != "" {
		t.Error("Filename should be empty")
	}
	if parsed.KyberEncLen != 0 {
		t.Error("KyberEncLen should be 0")
	}
	if parsed.KyberEnc != nil && len(parsed.KyberEnc) > 0 {
		t.Error("KyberEnc should be empty")
	}
}

func TestHeaderSizeEstimation(t *testing.T) {
	// 测试典型大小的头部
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x01,
		FilenameLen: 12,  // document.pdf = 12 bytes
		Filename:    "document.pdf",
		FileSize:    1024 * 1024, // 1MB
		Timestamp:   1734672000,
		KyberEncLen: 1088,
		KyberEnc:    make([]byte, 1088),
		ECDHLen:     32,
		ECDHPub:     [32]byte{},
		IVLen:       12,
		IV:          [12]byte{},
		SigLen:      2420,
		Signature:   make([]byte, 2420),
		SHA256Hash:  [32]byte{},
	}

	data, _ := header.MarshalBinary()
	t.Logf("Typical header size: %d bytes", len(data))

	// 验证头部大小在合理范围内（3KB - 4KB）
	if len(data) < 3000 || len(data) > 4096 {
		t.Errorf("Header size %d not in expected range 3000-4096", len(data))
	}
}
