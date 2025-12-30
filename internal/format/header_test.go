package format

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func TestFileHeaderSerialization(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x01,
		FilenameLen: 8, // 修正：test.txt 是 8 字节
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
		FilenameLen: 12, // document.pdf = 12 bytes
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

// TestValidateInvalidECDHLength 测试无效ECDH长度
func TestValidateInvalidECDHLength(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		FilenameLen: 0,
		Filename:    "",
		KyberEncLen: 0,
		KyberEnc:    nil,
		ECDHLen:     16, // 无效长度
		ECDHPub:     [32]byte{},
		IVLen:       12,
		IV:          [12]byte{},
		SigLen:      0,
		Signature:   nil,
		SHA256Hash:  [32]byte{},
	}

	err := header.Validate()
	if err == nil {
		t.Error("无效ECDH长度应该返回错误")
	}
}

// TestValidateInvalidIVLength 测试无效IV长度
func TestValidateInvalidIVLength(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		FilenameLen: 0,
		Filename:    "",
		KyberEncLen: 0,
		KyberEnc:    nil,
		ECDHLen:     32,
		ECDHPub:     [32]byte{},
		IVLen:       8, // 无效长度
		IV:          [12]byte{},
		SigLen:      0,
		Signature:   nil,
		SHA256Hash:  [32]byte{},
	}

	err := header.Validate()
	if err == nil {
		t.Error("无效IV长度应该返回错误")
	}
}

// TestUnmarshalBinaryInvalidData 测试反序列化无效数据
func TestUnmarshalBinaryInvalidData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"too short", []byte{1, 2, 3}},
		{"invalid magic", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"truncated after magic", []byte{'F', 'Z', 'J', 0x01, 0x01, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var header FileHeader
			err := header.UnmarshalBinary(tt.data)
			if err == nil {
				t.Errorf("无效数据 %s 应该返回错误", tt.name)
			}
		})
	}
}

// TestMarshalBinaryOptimizedConsistency 测试优化序列化一致性
func TestMarshalBinaryOptimizedConsistency(t *testing.T) {
	// 测试不同大小的数据
	testCases := []struct {
		name     string
		filename string
		kyberLen int
		sigLen   int
	}{
		{"small", "a.txt", 100, 100},
		{"medium", "medium.txt", 1088, 2420},
		{"large", "large.bin", 1088, 2700},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			header := &FileHeader{
				Magic:       [4]byte{'F', 'Z', 'J', 0x01},
				Version:     0x0100,
				Algorithm:   0x02,
				FilenameLen: byte(len(tc.filename)),
				Filename:    tc.filename,
				FileSize:    1024,
				Timestamp:   uint32(time.Now().Unix()),
				KyberEncLen: uint16(tc.kyberLen),
				KyberEnc:    make([]byte, tc.kyberLen),
				ECDHLen:     32,
				ECDHPub:     [32]byte{},
				IVLen:       12,
				IV:          [12]byte{},
				SigLen:      uint16(tc.sigLen),
				Signature:   make([]byte, tc.sigLen),
				SHA256Hash:  [32]byte{},
			}

			standard, _ := header.MarshalBinary()
			optimized, _ := header.MarshalBinaryOptimized()

			if !bytes.Equal(standard, optimized) {
				t.Errorf("序列化结果不一致: standard=%d, optimized=%d", len(standard), len(optimized))
			}
		})
	}
}

// TestGetHeaderSizeConsistency 测试头部大小计算一致性
func TestGetHeaderSizeConsistency(t *testing.T) {
	testCases := []struct {
		name     string
		header   *FileHeader
	}{
		{
			"empty",
			&FileHeader{
				Magic:       [4]byte{'F', 'Z', 'J', 0x01},
				Version:     0x0100,
				Algorithm:   0x02,
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
			},
		},
		{
			"full",
			&FileHeader{
				Magic:       [4]byte{'F', 'Z', 'J', 0x01},
				Version:     0x0100,
				Algorithm:   0x02,
				FilenameLen: 12,
				Filename:    "document.pdf",
				FileSize:    1048576,
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
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calculated := tc.header.GetHeaderSize()
			actual, _ := tc.header.MarshalBinary()

			if calculated != len(actual) {
				t.Errorf("大小计算不一致: 计算=%d, 实际=%d", calculated, len(actual))
			}
		})
	}
}

// TestUnixTimeEdgeCases 测试UnixTime边界情况
func TestUnixTimeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		timestamp uint32
	}{
		{"zero", 0},
		{"epoch", 1735689600}, // 2025-01-01
		{"max", 0xFFFFFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnixTime(tt.timestamp)
			if result == "" {
				t.Error("UnixTime 不应该返回空字符串")
			}
			t.Logf("%s: %s", tt.name, result)
		})
	}
}

// TestValidateSignatureLengthMismatch 测试签名长度不匹配
func TestValidateSignatureLengthMismatch(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		FilenameLen: 0,
		Filename:    "",
		KyberEncLen: 0,
		KyberEnc:    nil,
		ECDHLen:     32,
		ECDHPub:     [32]byte{},
		IVLen:       12,
		IV:          [12]byte{},
		SigLen:      100,
		Signature:   make([]byte, 200), // 长度不匹配
		SHA256Hash:  [32]byte{},
	}

	err := header.Validate()
	if err == nil {
		t.Error("签名长度不匹配应该返回错误")
	}
}

// TestValidateKyberEncLengthMismatch 测试Kyber封装长度不匹配
func TestValidateKyberEncLengthMismatch(t *testing.T) {
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		FilenameLen: 0,
		Filename:    "",
		KyberEncLen: 100,
		KyberEnc:    make([]byte, 200), // 长度不匹配
		ECDHLen:     32,
		ECDHPub:     [32]byte{},
		IVLen:       12,
		IV:          [12]byte{},
		SigLen:      0,
		Signature:   nil,
		SHA256Hash:  [32]byte{},
	}

	err := header.Validate()
	if err == nil {
		t.Error("Kyber封装长度不匹配应该返回错误")
	}
}

// TestUnmarshalBinaryReadErrors 测试反序列化读取错误
func TestUnmarshalBinaryReadErrors(t *testing.T) {
	// 创建不完整的头部数据
	incompleteData := make([]byte, 15) // 小于最小要求

	var header FileHeader
	err := header.UnmarshalBinary(incompleteData)
	if err == nil {
		t.Error("不完整数据应该返回错误")
	}
}

// TestMarshalBinaryWithRealData 测试使用真实数据的序列化
func TestMarshalBinaryWithRealData(t *testing.T) {
	// 使用真实大小的数据
	filename := "large_file.bin"
	header := &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x00,
		FilenameLen: byte(len(filename)),
		FileSize:    10485760,
		Timestamp:   uint32(time.Now().Unix()),
		KyberEncLen: 1088,
		ECDHLen:     32,
		IVLen:       12,
		SigLen:      2700,
	}
	header.Filename = filename
	header.KyberEnc = make([]byte, 1088)
	header.ECDHPub = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	header.IV = [12]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	header.Signature = make([]byte, 2700)
	header.SHA256Hash = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	// 标准序列化
	data1, err1 := header.MarshalBinary()
	if err1 != nil {
		t.Fatalf("标准序列化失败: %v", err1)
	}

	// 优化序列化
	data2, err2 := header.MarshalBinaryOptimized()
	if err2 != nil {
		t.Fatalf("优化序列化失败: %v", err2)
	}

	// 验证结果一致
	if !bytes.Equal(data1, data2) {
		t.Error("两种序列化方法结果不一致")
	}

	// 验证可以反序列化
	var decoded FileHeader
	if err := decoded.UnmarshalBinary(data1); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 验证关键字段
	if decoded.Filename != header.Filename {
		t.Errorf("文件名不匹配: %s vs %s", decoded.Filename, header.Filename)
	}
	if decoded.FileSize != header.FileSize {
		t.Errorf("文件大小不匹配: %d vs %d", decoded.FileSize, header.FileSize)
	}
}

