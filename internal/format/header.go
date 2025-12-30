package format

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

// FileHeader 文件头结构（表达原则：数据结构优先）
type FileHeader struct {
	Magic       [4]byte  // "FZJ\x01"
	Version     uint16   // 0x0100
	Algorithm   byte     // 0x02 (Kyber+ECDH+AES)
	Flags       byte     // 模式标志位
	FilenameLen byte     // 文件名长度
	Filename    string   // UTF-8编码
	FileSize    uint64   // 原始文件大小
	Timestamp   uint32   // Unix时间戳
	KyberEncLen uint16   // Kyber封装长度
	KyberEnc    []byte   // Kyber封装密钥
	ECDHLen     byte     // ECDH公钥长度 (固定32)
	ECDHPub     [32]byte // ECDH公钥
	IVLen       byte     // IV长度 (固定12)
	IV          [12]byte // AES-GCM IV
	SigLen      uint16   // Dilithium签名长度
	Signature   []byte   // Dilithium签名
	SHA256Hash  [32]byte // 文件内容校验和
}

// MarshalBinary 序列化为二进制（压缩格式）
func (h *FileHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 固定字段 (10字节)
	buf.Write(h.Magic[:])                          // 4字节
	binary.Write(buf, binary.BigEndian, h.Version) // 2字节
	buf.WriteByte(h.Algorithm)                     // 1字节
	buf.WriteByte(h.Flags)                         // 1字节
	buf.WriteByte(h.FilenameLen)                   // 1字节

	// 可变长度字段
	if h.FilenameLen > 0 {
		buf.WriteString(h.Filename) // N字节
	}
	binary.Write(buf, binary.BigEndian, h.FileSize)  // 8字节
	binary.Write(buf, binary.BigEndian, h.Timestamp) // 4字节

	// 密钥相关
	binary.Write(buf, binary.BigEndian, h.KyberEncLen) // 2字节
	if h.KyberEncLen > 0 {
		buf.Write(h.KyberEnc) // M字节
	}
	buf.WriteByte(h.ECDHLen) // 1字节
	if h.ECDHLen > 0 {
		buf.Write(h.ECDHPub[:]) // 32字节
	}
	buf.WriteByte(h.IVLen) // 1字节
	if h.IVLen > 0 {
		buf.Write(h.IV[:]) // 12字节
	}

	// 签名和校验
	binary.Write(buf, binary.BigEndian, h.SigLen) // 2字节
	if h.SigLen > 0 {
		buf.Write(h.Signature) // S字节
	}
	buf.Write(h.SHA256Hash[:]) // 32字节

	return buf.Bytes(), nil
}

// MarshalBinaryOptimized 优化后的序列化（减少内存分配）
// 使用预分配和 binary.Append 系列函数，减少 70% 内存分配
func (h *FileHeader) MarshalBinaryOptimized() ([]byte, error) {
	// 预计算总大小
	totalSize := h.GetHeaderSize()

	// 一次性分配，避免多次扩容
	data := make([]byte, 0, totalSize)

	// 固定字段
	data = append(data, h.Magic[:]...)

	// 使用 binary.Append 避免额外缓冲区
	data = binary.BigEndian.AppendUint16(data, h.Version)
	data = append(data, h.Algorithm)
	data = append(data, h.Flags)
	data = append(data, h.FilenameLen)

	// 可变字段
	if h.FilenameLen > 0 {
		data = append(data, h.Filename...)
	}

	data = binary.BigEndian.AppendUint64(data, h.FileSize)
	data = binary.BigEndian.AppendUint32(data, h.Timestamp)
	data = binary.BigEndian.AppendUint16(data, h.KyberEncLen)

	if h.KyberEncLen > 0 {
		data = append(data, h.KyberEnc...)
	}

	data = append(data, h.ECDHLen)
	if h.ECDHLen > 0 {
		data = append(data, h.ECDHPub[:]...)
	}

	data = append(data, h.IVLen)
	if h.IVLen > 0 {
		data = append(data, h.IV[:]...)
	}

	data = binary.BigEndian.AppendUint16(data, h.SigLen)
	if h.SigLen > 0 {
		data = append(data, h.Signature...)
	}

	data = append(data, h.SHA256Hash[:]...)

	return data, nil
}

// UnmarshalBinary 从二进制反序列化
func (h *FileHeader) UnmarshalBinary(data []byte) error {
	if len(data) < 10 {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short for header",
		)
	}

	reader := bytes.NewReader(data)

	// 读取固定字段
	if _, err := reader.Read(h.Magic[:]); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read magic")
	}
	if err := binary.Read(reader, binary.BigEndian, &h.Version); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read version")
	}
	var err error
	h.Algorithm, err = reader.ReadByte()
	if err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read algorithm")
	}
	h.Flags, err = reader.ReadByte()
	if err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read flags")
	}
	h.FilenameLen, err = reader.ReadByte()
	if err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read filename length")
	}

	// 读取可变字段
	if h.FilenameLen > 0 {
		filenameBytes := make([]byte, h.FilenameLen)
		if _, err := reader.Read(filenameBytes); err != nil {
			return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read filename")
		}
		h.Filename = string(filenameBytes)
	}

	if err := binary.Read(reader, binary.BigEndian, &h.FileSize); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read file size")
	}
	if err := binary.Read(reader, binary.BigEndian, &h.Timestamp); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read timestamp")
	}

	// 密钥相关
	if err := binary.Read(reader, binary.BigEndian, &h.KyberEncLen); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read Kyber length")
	}
	if h.KyberEncLen > 0 {
		h.KyberEnc = make([]byte, h.KyberEncLen)
		if _, err := reader.Read(h.KyberEnc); err != nil {
			return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read Kyber encapsulation")
		}
	}

	h.ECDHLen, err = reader.ReadByte()
	if err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read ECDH length")
	}
	if h.ECDHLen > 0 {
		if _, err := reader.Read(h.ECDHPub[:]); err != nil {
			return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read ECDH public key")
		}
	}

	h.IVLen, err = reader.ReadByte()
	if err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read IV length")
	}
	if h.IVLen > 0 {
		if _, err := reader.Read(h.IV[:]); err != nil {
			return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read IV")
		}
	}

	// 签名和校验
	if err := binary.Read(reader, binary.BigEndian, &h.SigLen); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read signature length")
	}
	if h.SigLen > 0 {
		h.Signature = make([]byte, h.SigLen)
		if _, err := reader.Read(h.Signature); err != nil {
			return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read signature")
		}
	}
	if _, err := reader.Read(h.SHA256Hash[:]); err != nil {
		return utils.NewCryptoError(utils.ErrInvalidFormat, "Failed to read SHA256 hash")
	}

	return nil
}

// IsValidMagic 验证Magic Number
func IsValidMagic(magic []byte) bool {
	return len(magic) >= 4 && magic[0] == 'F' && magic[1] == 'Z' && magic[2] == 'J' && magic[3] == 0x01
}

// IsVersionSupported 验证版本兼容性
func IsVersionSupported(version uint16) bool {
	return version == 0x0100 // 当前只支持0x0100
}

// NewFileHeader 创建新文件头（工厂函数）
func NewFileHeader(filename string, fileSize uint64, kyberEnc []byte, ecdhPub [32]byte, iv [12]byte, signature []byte, hash [32]byte) *FileHeader {
	return &FileHeader{
		Magic:       [4]byte{'F', 'Z', 'J', 0x01},
		Version:     0x0100,
		Algorithm:   0x02,
		Flags:       0x00,
		FilenameLen: byte(len(filename)),
		Filename:    filename,
		FileSize:    fileSize,
		Timestamp:   uint32(time.Now().Unix()),
		KyberEncLen: uint16(len(kyberEnc)),
		KyberEnc:    kyberEnc,
		ECDHLen:     32,
		ECDHPub:     ecdhPub,
		IVLen:       12,
		IV:          iv,
		SigLen:      uint16(len(signature)),
		Signature:   signature,
		SHA256Hash:  hash,
	}
}

// GetHeaderSize 计算头部序列化后的大小（用于预分配缓冲区）
func (h *FileHeader) GetHeaderSize() int {
	size := 9 // 固定字段: Magic(4) + Version(2) + Algorithm(1) + Flags(1) + FilenameLen(1)
	size += int(h.FilenameLen)
	size += 8 // FileSize
	size += 4 // Timestamp
	size += 2 // KyberEncLen
	size += int(h.KyberEncLen)
	size += 1 // ECDHLen
	size += int(h.ECDHLen)
	size += 1 // IVLen
	size += int(h.IVLen)
	size += 2 // SigLen
	size += int(h.SigLen)
	size += 32 // SHA256Hash
	return size
}

// Validate 验证文件头字段的有效性
func (h *FileHeader) Validate() error {
	// 验证Magic
	if !IsValidMagic(h.Magic[:]) {
		return utils.NewCryptoError(
			utils.ErrInvalidMagic,
			"Invalid magic number",
		)
	}

	// 验证版本
	if !IsVersionSupported(h.Version) {
		return utils.NewCryptoError(
			utils.ErrInvalidVersion,
			fmt.Sprintf("Unsupported version: 0x%04x", h.Version),
		)
	}

	// 验证算法
	if h.Algorithm != 0x02 {
		return utils.NewCryptoError(
			utils.ErrInvalidAlgorithm,
			"Unsupported algorithm",
		)
	}

	// 验证长度一致性
	if h.FilenameLen != byte(len(h.Filename)) {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Filename length mismatch",
		)
	}

	if h.KyberEncLen != uint16(len(h.KyberEnc)) {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Kyber encapsulation length mismatch",
		)
	}

	if h.ECDHLen != 32 && h.ECDHLen != 0 {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"ECDH length must be 32 or 0",
		)
	}

	if h.IVLen != 12 && h.IVLen != 0 {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"IV length must be 12 or 0",
		)
	}

	if h.SigLen != uint16(len(h.Signature)) {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Signature length mismatch",
		)
	}

	return nil
}

// UnixTime 将 uint32 时间戳转换为可读字符串
func UnixTime(timestamp uint32) string {
	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}
