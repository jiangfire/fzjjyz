package zjcrypto

import (
	"crypto/ecdh"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
)

const (
	// encryptedFilePerm 加密文件权限。
	encryptedFilePerm = 0644
	// decryptedFilePerm 解密文件权限。
	decryptedFilePerm = 0644
)

// EncryptionData 包含加密所需的所有数据.
type EncryptionData struct {
	Plaintext    []byte
	Filename     string
	FileSize     uint64
	Encapsulated []byte
	ECDHTempPub  []byte
	SharedSecret []byte
	IV           []byte
	Hash         [32]byte
	Signature    []byte
}

// DecryptionData 包含解密所需的所有数据.
type DecryptionData struct {
	Ciphertext   []byte
	Header       *format.FileHeader
	SharedSecret []byte
	Plaintext    []byte
	Hash         [32]byte
}

// prepareEncryptionKeys 执行混合密钥封装
// 返回: Kyber密文, 临时ECDH公钥, 组合共享密钥.
func prepareEncryptionKeys(kyberPub kem.PublicKey, ecdhPub *ecdh.PublicKey) ([]byte, []byte, []byte, error) {
	encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
	return encryptor.Encapsulate()
}

// encryptAESGCM 使用AES-256-GCM加密数据.
func encryptAESGCM(sharedSecret []byte, plaintext []byte) (ciphertext []byte, iv []byte, err error) {
	// AESGCMEncrypt 已经生成随机 nonce
	ciphertext, iv, err = AESGCMEncrypt(sharedSecret, plaintext)
	if err != nil {
		return nil, nil, err
	}

	return ciphertext, iv, nil
}

// decryptAESGCM 使用AES-256-GCM解密数据.
func decryptAESGCM(sharedSecret []byte, ciphertext []byte, iv []byte) ([]byte, error) {
	return AESGCMDecrypt(sharedSecret, ciphertext, iv)
}

// calculateHash 计算数据的SHA256哈希.
func calculateHash(data []byte) [32]byte {
	return sha256.Sum256(data)
}

// signHash 对哈希进行Dilithium签名.
func signHash(hash []byte, dilithiumPriv interface{}) ([]byte, error) {
	return SignHash(hash, dilithiumPriv)
}

// verifyHashSignature 验证哈希签名.
func verifyHashSignature(hash []byte, signature []byte, dilithiumPub interface{}) (bool, error) {
	return VerifyHashSignature(hash, signature, dilithiumPub)
}

// buildFileHeader 构建文件头.
func buildFileHeader(
	filename string,
	fileSize uint64,
	encapsulated []byte,
	ecdhTempPub []byte,
	iv []byte,
	signature []byte,
	hash [32]byte,
) (*format.FileHeader, error) {
	// 转换类型
	var ecdhPubArray [32]byte
	copy(ecdhPubArray[:], ecdhTempPub)

	var ivArray [12]byte
	copy(ivArray[:], iv)

	return format.NewFileHeader(
		filename,
		fileSize,
		encapsulated,
		ecdhPubArray,
		ivArray,
		signature,
		hash,
	), nil
}

// serializeHeader 序列化文件头.
func serializeHeader(header *format.FileHeader) ([]byte, error) {
	// 优先使用优化的序列化方法
	data, err := header.MarshalBinaryOptimized()
	if err != nil {
		return nil, fmt.Errorf("marshal header: %w", err)
	}
	return data, nil
}

// writeEncryptedFile 写入加密文件（头部 + 密文）.
func writeEncryptedFile(outputPath string, headerBytes []byte, ciphertext []byte) error {
	// 预分配内存
	outputData := make([]byte, 0, len(headerBytes)+len(ciphertext))
	outputData = append(outputData, headerBytes...)
	outputData = append(outputData, ciphertext...)

	if err := os.WriteFile(outputPath, outputData, encryptedFilePerm); err != nil {
		return fmt.Errorf("write encrypted file: %w", err)
	}
	return nil
}

// parseEncryptedFile 读取并解析加密文件.
func parseEncryptedFile(inputPath string) (header *format.FileHeader, ciphertext []byte, err error) {
	// G304: inputPath 应由调用方验证
	encryptedData, err := os.ReadFile(inputPath) //nolint:gosec
	if err != nil {
		return nil, nil, fmt.Errorf("read encrypted file: %w", err)
	}

	// 解析文件头
	header, err = format.ParseFileHeaderFromBytes(encryptedData)
	if err != nil {
		return nil, nil, fmt.Errorf("parse file header: %w", err)
	}

	// 验证头部
	if err := header.Validate(); err != nil {
		return nil, nil, fmt.Errorf("header validation failed: %w", err)
	}

	// 提取密文
	headerSize := header.GetHeaderSize()
	if len(encryptedData) <= headerSize {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"File too short - no ciphertext after header",
		)
	}
	ciphertext = encryptedData[headerSize:]

	return header, ciphertext, nil
}

// decapsulateKeys 解封装密钥.
func decapsulateKeys(
	kyberPriv kem.PrivateKey,
	ecdhPriv *ecdh.PrivateKey,
	encapsulated []byte,
	ecdhPub []byte,
) ([]byte, error) {
	decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)
	sharedSecret, err := decryptor.Decapsulate(encapsulated, ecdhPub)
	if err != nil {
		return nil, fmt.Errorf("decapsulate: %w", err)
	}
	return sharedSecret, nil
}

// writeDecryptedFile 写入解密文件.
func writeDecryptedFile(outputPath string, plaintext []byte) error {
	if err := os.WriteFile(outputPath, plaintext, decryptedFilePerm); err != nil {
		return fmt.Errorf("write decrypted file: %w", err)
	}
	return nil
}

// verifyDecryptionIntegrity 验证解密数据的完整性和签名.
func verifyDecryptionIntegrity(plaintext []byte, header *format.FileHeader, dilithiumPub interface{}) error {
	// 验证哈希
	hash := calculateHash(plaintext)
	if hash != header.SHA256Hash {
		return utils.NewCryptoError(
			utils.ErrHashMismatch,
			"SHA256 hash mismatch - file may be corrupted",
		)
	}

	// 验证签名（如果提供）
	if dilithiumPub != nil && header.SigLen > 0 {
		valid, err := verifyHashSignature(hash[:], header.Signature, dilithiumPub)
		if err != nil {
			return utils.NewCryptoError(
				utils.ErrVerificationFailed,
				"Signature verification error: "+err.Error(),
			)
		}
		if !valid {
			return utils.NewCryptoError(
				utils.ErrVerificationFailed,
				"Invalid signature - file may be tampered",
			)
		}
	}

	return nil
}

// EncryptFileCore 加密文件的核心逻辑
// 这个函数被 EncryptFile 和 StreamingEncryptor 共用.
func EncryptFileCore(
	inputPath string,
	kyberPub kem.PublicKey,
	ecdhPub *ecdh.PublicKey,
	dilithiumPriv interface{},
) (header *format.FileHeader, ciphertext []byte, err error) {
	// G304: inputPath 应由调用方验证
	plaintext, err := os.ReadFile(inputPath) //nolint:gosec
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to read input file: "+err.Error(),
		)
	}

	// 1. 混合密钥封装
	encapsulated, ecdhTempPub, sharedSecret, err := prepareEncryptionKeys(kyberPub, ecdhPub)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"Hybrid encapsulation failed: "+err.Error(),
		)
	}

	// 2. AES-GCM 加密
	ciphertext, iv, err := encryptAESGCM(sharedSecret, plaintext)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrEncryptionFailed,
			"AES-GCM encryption failed: "+err.Error(),
		)
	}

	// 3. 计算哈希
	hash := calculateHash(plaintext)

	// 4. 签名
	signature, err := signHash(hash[:], dilithiumPriv)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrSigningFailed,
			"Hash signing failed: "+err.Error(),
		)
	}

	// 5. 构建头部
	filename := filepath.Base(inputPath)
	header, err = buildFileHeader(
		filename,
		uint64(len(plaintext)),
		encapsulated,
		ecdhTempPub,
		iv,
		signature,
		hash,
	)
	if err != nil {
		return nil, nil, err
	}

	return header, ciphertext, nil
}

// DecryptFileCore 解密文件的核心逻辑
// 这个函数被 DecryptFile 和 StreamingDecryptor 共用.
func DecryptFileCore(
	inputPath string,
	kyberPriv kem.PrivateKey,
	ecdhPriv *ecdh.PrivateKey,
	dilithiumPub interface{},
) (plaintext []byte, err error) {
	// 1. 读取并解析文件
	header, ciphertext, err := parseEncryptedFile(inputPath)
	if err != nil {
		return nil, err
	}

	// 2. 密钥解封装
	sharedSecret, err := decapsulateKeys(kyberPriv, ecdhPriv, header.KyberEnc, header.ECDHPub[:])
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrAuthFailed,
			"Hybrid decapsulation failed: "+err.Error(),
		)
	}

	// 3. AES-GCM 解密
	plaintext, err = decryptAESGCM(sharedSecret, ciphertext, header.IV[:])
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrDecryptionFailed,
			"AES-GCM decryption failed: "+err.Error(),
		)
	}

	// 4. 验证完整性和签名
	if err := verifyDecryptionIntegrity(plaintext, header, dilithiumPub); err != nil {
		return nil, err
	}

	return plaintext, nil
}
