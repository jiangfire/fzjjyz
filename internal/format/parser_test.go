package format

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

func TestParseFileHeader(t *testing.T) {
	// 创建一个完整的文件头
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x01,
		FilenameLen: 8,
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

	// 填充随机数据
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

	// 解析
	parsed, err := ParseFileHeader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseFileHeader failed: %v", err)
	}

	// 验证所有字段
	if !bytes.Equal(header.Magic[:], parsed.Magic[:]) {
		t.Error("Magic mismatch")
	}
	if header.Version != parsed.Version {
		t.Errorf("Version mismatch: %04x != %04x", header.Version, parsed.Version)
	}
	if header.Algorithm != parsed.Algorithm {
		t.Errorf("Algorithm mismatch: %d != %d", header.Algorithm, parsed.Algorithm)
	}
	if header.Flags != parsed.Flags {
		t.Errorf("Flags mismatch: %d != %d", header.Flags, parsed.Flags)
	}
	if header.FilenameLen != parsed.FilenameLen {
		t.Errorf("FilenameLen mismatch: %d != %d", header.FilenameLen, parsed.FilenameLen)
	}
	if header.Filename != parsed.Filename {
		t.Errorf("Filename mismatch: %s != %s", header.Filename, parsed.Filename)
	}
	if header.FileSize != parsed.FileSize {
		t.Errorf("FileSize mismatch: %d != %d", header.FileSize, parsed.FileSize)
	}
	if header.Timestamp != parsed.Timestamp {
		t.Errorf("Timestamp mismatch: %d != %d", header.Timestamp, parsed.Timestamp)
	}
	if header.KyberEncLen != parsed.KyberEncLen {
		t.Errorf("KyberEncLen mismatch: %d != %d", header.KyberEncLen, parsed.KyberEncLen)
	}
	if !bytes.Equal(header.KyberEnc, parsed.KyberEnc) {
		t.Error("KyberEnc mismatch")
	}
	if header.ECDHLen != parsed.ECDHLen {
		t.Errorf("ECDHLen mismatch: %d != %d", header.ECDHLen, parsed.ECDHLen)
	}
	if !bytes.Equal(header.ECDHPub[:], parsed.ECDHPub[:]) {
		t.Error("ECDHPub mismatch")
	}
	if header.IVLen != parsed.IVLen {
		t.Errorf("IVLen mismatch: %d != %d", header.IVLen, parsed.IVLen)
	}
	if !bytes.Equal(header.IV[:], parsed.IV[:]) {
		t.Error("IV mismatch")
	}
	if header.SigLen != parsed.SigLen {
		t.Errorf("SigLen mismatch: %d != %d", header.SigLen, parsed.SigLen)
	}
	if !bytes.Equal(header.Signature, parsed.Signature) {
		t.Error("Signature mismatch")
	}
	if !bytes.Equal(header.SHA256Hash[:], parsed.SHA256Hash[:]) {
		t.Error("SHA256Hash mismatch")
	}
}

func TestParseFileHeaderValidation(t *testing.T) {
	// 测试解析后验证
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x00,
		FilenameLen: 12,  // document.pdf = 12 bytes
		Filename:    "document.pdf",
		FileSize:    2048,
		Timestamp:   uint32(1734672000),
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
	parsed, err := ParseFileHeader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseFileHeader failed: %v", err)
	}

	// 验证
	if err := parsed.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}

func TestParseInvalidMagic(t *testing.T) {
	// 创建无效的 magic
	header := &FileHeader{
		Magic:       [4]byte{'X', 'Y', 'Z', 0x01}, // 无效
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

	data, _ := header.MarshalBinary()
	_, err := ParseFileHeader(bytes.NewReader(data))
	if err == nil {
		t.Error("Should fail on invalid magic")
	}
	if !utils.IsFormatError(err) {
		t.Errorf("Expected format error, got: %v", err)
	}
}

func TestParseUnsupportedVersion(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0200, // 未来版本
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

	data, _ := header.MarshalBinary()
	_, err := ParseFileHeader(bytes.NewReader(data))
	if err == nil {
		t.Error("Should fail on unsupported version")
	}
	if !utils.IsFormatError(err) {
		t.Errorf("Expected format error, got: %v", err)
	}
}

func TestParseEmptyFile(t *testing.T) {
	// 测试空文件
	_, err := ParseFileHeader(bytes.NewReader([]byte{}))
	if err == nil {
		t.Error("Should fail on empty data")
	}
}

func TestParseFileWithFile(t *testing.T) {
	// 测试从文件解析
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.fzj")

	// 创建完整的文件头数据
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x01,
		FilenameLen: 13,  // encrypted.bin = 13 bytes
		Filename:    "encrypted.bin",
		FileSize:    1024 * 1024,
		Timestamp:   uint32(1734672000),
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

	// 填充随机数据
	rand.Read(header.KyberEnc)
	rand.Read(header.ECDHPub[:])
	rand.Read(header.IV[:])
	rand.Read(header.Signature)
	rand.Read(header.SHA256Hash[:])

	data, _ := header.MarshalBinary()
	if err := os.WriteFile(testFile, data, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// 从文件解析
	f, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer f.Close()

	parsed, err := ParseFileHeader(f)
	if err != nil {
		t.Fatalf("ParseFileHeader from file failed: %v", err)
	}

	if parsed.Filename != "encrypted.bin" {
		t.Errorf("Filename mismatch: %s", parsed.Filename)
	}
	if parsed.FileSize != 1024*1024 {
		t.Errorf("FileSize mismatch: %d", parsed.FileSize)
	}
}

func TestParseWithCorruptedLengths(t *testing.T) {
	// 测试长度字段不匹配的情况
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x00,
		FilenameLen: 8,
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

	data, _ := header.MarshalBinary()

	// 修改长度字段使其远大于实际数据，导致读取时数据不足
	// KyberEncLen 在位置 29-30，修改为 5000 (远大于剩余数据)
	data[29] = 0x13
	data[30] = 0x88 // 5000 instead of 1088

	// 现在数据总长度是 3619，KyberEncLen=5000 会导致解析器尝试读取5000字节
	// 但数据只有3619字节，应该失败
	parsed, err := ParseFileHeader(bytes.NewReader(data))
	if err == nil {
		t.Errorf("Should fail on corrupted length (length > data). Got parsed with KyberEncLen=%d, KyberEnc length=%d",
			parsed.KyberEncLen, len(parsed.KyberEnc))
	} else {
		t.Logf("Correctly failed with error: %v", err)
	}
}

func TestParseWithTruncatedData(t *testing.T) {
	// 测试数据被截断的情况
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x00,
		FilenameLen: 8,
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

	data, _ := header.MarshalBinary()

	// 截断数据，移除部分 KyberEnc 数据
	truncated := data[:len(data)-100]
	_, err := ParseFileHeader(bytes.NewReader(truncated))
	if err == nil {
		t.Error("Should fail on truncated data")
	} else {
		t.Logf("Correctly failed with error: %v", err)
	}
}

func TestParseHeaderWithEmptyFields(t *testing.T) {
	// 测试解析空字段
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

	data, _ := header.MarshalBinary()
	parsed, err := ParseFileHeader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseFileHeader with empty fields failed: %v", err)
	}

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
