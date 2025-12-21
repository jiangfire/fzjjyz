package format

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

// ParseFileHeader 从 Reader 中解析文件头
func ParseFileHeader(r io.Reader) (*FileHeader, error) {
	header := &FileHeader{}

	// 读取固定字段 (10字节)
	var magic [4]byte
	if _, err := io.ReadFull(r, magic[:]); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read magic: %v", err),
		)
	}
	header.Magic = magic

	// 验证 Magic
	if !IsValidMagic(header.Magic[:]) {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidMagic,
			"Invalid magic number",
		)
	}

	// 读取 Version (2字节)
	if err := binary.Read(r, binary.BigEndian, &header.Version); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read version: %v", err),
		)
	}

	// 验证 Version
	if !IsVersionSupported(header.Version) {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidVersion,
			fmt.Sprintf("Unsupported version: 0x%04x", header.Version),
		)
	}

	// 读取 Algorithm (1字节)
	if err := binary.Read(r, binary.BigEndian, &header.Algorithm); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read algorithm: %v", err),
		)
	}

	// 验证 Algorithm
	if header.Algorithm != 0x02 {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidAlgorithm,
			"Unsupported algorithm",
		)
	}

	// 读取 Flags (1字节)
	if err := binary.Read(r, binary.BigEndian, &header.Flags); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read flags: %v", err),
		)
	}

	// 读取 FilenameLen (1字节)
	if err := binary.Read(r, binary.BigEndian, &header.FilenameLen); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read filename length: %v", err),
		)
	}

	// 读取 Filename
	if header.FilenameLen > 0 {
		filenameBytes := make([]byte, header.FilenameLen)
		if _, err := io.ReadFull(r, filenameBytes); err != nil {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidFormat,
				fmt.Sprintf("Failed to read filename: %v", err),
			)
		}
		header.Filename = string(filenameBytes)
	}

	// 读取 FileSize (8字节)
	if err := binary.Read(r, binary.BigEndian, &header.FileSize); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read file size: %v", err),
		)
	}

	// 读取 Timestamp (4字节)
	if err := binary.Read(r, binary.BigEndian, &header.Timestamp); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read timestamp: %v", err),
		)
	}

	// 读取 KyberEncLen (2字节)
	if err := binary.Read(r, binary.BigEndian, &header.KyberEncLen); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read Kyber length: %v", err),
		)
	}

	// 读取 KyberEnc
	if header.KyberEncLen > 0 {
		header.KyberEnc = make([]byte, header.KyberEncLen)
		if _, err := io.ReadFull(r, header.KyberEnc); err != nil {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidFormat,
				fmt.Sprintf("Failed to read Kyber encapsulation: %v", err),
			)
		}
	}

	// 读取 ECDHLen (1字节)
	if err := binary.Read(r, binary.BigEndian, &header.ECDHLen); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read ECDH length: %v", err),
		)
	}

	// 读取 ECDHPub
	if header.ECDHLen > 0 {
		if _, err := io.ReadFull(r, header.ECDHPub[:]); err != nil {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidFormat,
				fmt.Sprintf("Failed to read ECDH public key: %v", err),
			)
		}
	}

	// 读取 IVLen (1字节)
	if err := binary.Read(r, binary.BigEndian, &header.IVLen); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read IV length: %v", err),
		)
	}

	// 读取 IV
	if header.IVLen > 0 {
		if _, err := io.ReadFull(r, header.IV[:]); err != nil {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidFormat,
				fmt.Sprintf("Failed to read IV: %v", err),
			)
		}
	}

	// 读取 SigLen (2字节)
	if err := binary.Read(r, binary.BigEndian, &header.SigLen); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read signature length: %v", err),
		)
	}

	// 读取 Signature
	if header.SigLen > 0 {
		header.Signature = make([]byte, header.SigLen)
		if _, err := io.ReadFull(r, header.Signature); err != nil {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidFormat,
				fmt.Sprintf("Failed to read signature: %v", err),
			)
		}
	}

	// 读取 SHA256Hash (32字节)
	if _, err := io.ReadFull(r, header.SHA256Hash[:]); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			fmt.Sprintf("Failed to read SHA256 hash: %v", err),
		)
	}

	return header, nil
}

// ParseFileHeaderFromBytes 从字节切片解析文件头
func ParseFileHeaderFromBytes(data []byte) (*FileHeader, error) {
	return ParseFileHeader(bytes.NewReader(data))
}

// ExtractHeaderSize 从加密文件中提取头部大小
// 这在需要只读取头部而不解析完整数据时很有用
func ExtractHeaderSize(data []byte) (int, error) {
	if len(data) < 31 {
		return 0, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short to contain header size",
		)
	}

	// 固定字段位置
	// Magic: 0-3
	// Version: 4-5
	// Algorithm: 6
	// Flags: 7
	// FilenameLen: 8
	// Filename: 9 to 9+FilenameLen
	// FileSize: 9+FilenameLen to 9+FilenameLen+7
	// Timestamp: 9+FilenameLen+8 to 9+FilenameLen+11
	// KyberEncLen: 9+FilenameLen+12 to 9+FilenameLen+13

	filenameLen := data[8]
	pos := 9 + int(filenameLen) + 8 + 4 // Fixed + Filename + FileSize + Timestamp

	if len(data) < pos+2 {
		return 0, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short for KyberEncLen",
		)
	}

	kyberEncLen := int(data[pos])<<8 | int(data[pos+1])
	pos += 2 + kyberEncLen // KyberEncLen + KyberEnc

	if len(data) < pos+1 {
		return 0, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short for ECDHLen",
		)
	}

	ecdhLen := int(data[pos])
	pos += 1 + ecdhLen // ECDHLen + ECDHPub

	if len(data) < pos+1 {
		return 0, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short for IVLen",
		)
	}

	ivLen := int(data[pos])
	pos += 1 + ivLen // IVLen + IV

	if len(data) < pos+2 {
		return 0, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short for SigLen",
		)
	}

	sigLen := int(data[pos])<<8 | int(data[pos+1])
	pos += 2 + sigLen // SigLen + Signature

	if len(data) < pos+32 {
		return 0, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Data too short for SHA256Hash",
		)
	}

	pos += 32 // SHA256Hash
	return pos, nil
}

// IsValidEncryptedFile 检查数据是否为有效的加密文件格式
func IsValidEncryptedFile(data []byte) bool {
	if len(data) < 10 {
		return false
	}

	// 检查 Magic
	if !IsValidMagic(data[0:4]) {
		return false
	}

	// 检查 Version
	version := uint16(data[4])<<8 | uint16(data[5])
	if !IsVersionSupported(version) {
		return false
	}

	// 检查 Algorithm
	if data[6] != 0x02 {
		return false
	}

	return true
}

// GetHeaderInfo 从文件头提取基本信息（用于快速预览）
type HeaderInfo struct {
	Filename  string
	FileSize  uint64
	Timestamp uint32
	Algorithm string
	HasKyber bool
	HasECDH  bool
	HasIV    bool
	HasSig   bool
}

func GetHeaderInfo(header *FileHeader) *HeaderInfo {
	algo := "Kyber768+ECDH+AES256-GCM"
	if header.Algorithm != 0x02 {
		algo = fmt.Sprintf("Unknown(0x%02x)", header.Algorithm)
	}

	return &HeaderInfo{
		Filename:  header.Filename,
		FileSize:  header.FileSize,
		Timestamp: header.Timestamp,
		Algorithm: algo,
		HasKyber:  header.KyberEncLen > 0,
		HasECDH:   header.ECDHLen > 0,
		HasIV:     header.IVLen > 0,
		HasSig:    header.SigLen > 0,
	}
}
