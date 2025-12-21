package crypto

import (
	"crypto/ecdh"
	"crypto/sha256"
	"os"
	"path/filepath"

	"codeberg.org/jiangfire/fzjjyz/internal/format"
	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
)

// EncryptFile 加密文件
// 输入: 原始文件路径, 输出文件路径, Kyber公钥, ECDH公钥, Dilithium私钥
// 返回: 错误
func EncryptFile(inputPath, outputPath string, kyberPub kem.PublicKey, ecdhPub *ecdh.PublicKey, dilithiumPriv interface{}) error {
	// 读取原始文件
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to read input file: "+err.Error(),
		)
	}

	// 获取文件名（用于头部）
	filename := filepath.Base(inputPath)

	// 1. 混合密钥封装
	encryptor := NewHybridEncryptor(kyberPub, ecdhPub)
	encapsulated, ecdhTempPub, sharedSecret, err := encryptor.Encapsulate()
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"Hybrid encapsulation failed: "+err.Error(),
		)
	}

	// 2. AES-GCM 加密数据
	ciphertext, nonce, err := AESGCMEncrypt(sharedSecret, plaintext)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrEncryptionFailed,
			"AES-GCM encryption failed: "+err.Error(),
		)
	}

	// 3. 计算 SHA256 哈希（用于完整性验证）
	hash := sha256.Sum256(plaintext)

	// 4. 对哈希进行签名
	signature, err := SignHash(hash[:], dilithiumPriv)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrSigningFailed,
			"Hash signing failed: "+err.Error(),
		)
	}

	// 5. 构建文件头
	// 将 ECDH 临时公钥转换为 [32]byte
	var ecdhPubArray [32]byte
	copy(ecdhPubArray[:], ecdhTempPub)

	// 将 IV 转换为 [12]byte
	var ivArray [12]byte
	copy(ivArray[:], nonce)

	header := format.NewFileHeader(
		filename,
		uint64(len(plaintext)),
		encapsulated,
		ecdhPubArray,
		ivArray,
		signature,
		hash,
	)

	// 6. 序列化头部
	headerBytes, err := header.MarshalBinary()
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrSerializationFailed,
			"Header serialization failed: "+err.Error(),
		)
	}

	// 7. 写入加密文件
	// 格式: [头部] + [密文]
	outputData := make([]byte, 0, len(headerBytes)+len(ciphertext))
	outputData = append(outputData, headerBytes...)
	outputData = append(outputData, ciphertext...)

	err = os.WriteFile(outputPath, outputData, 0644)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to write output file: "+err.Error(),
		)
	}

	return nil
}

// DecryptFile 解密文件
// 输入: 加密文件路径, 输出文件路径, Kyber私钥, ECDH私钥, Dilithium公钥
// 返回: 错误
func DecryptFile(inputPath, outputPath string, kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, dilithiumPub interface{}) error {
	// 读取加密文件
	encryptedData, err := os.ReadFile(inputPath)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to read encrypted file: "+err.Error(),
		)
	}

	// 1. 解析文件头
	header, err := format.ParseFileHeaderFromBytes(encryptedData)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Failed to parse file header: "+err.Error(),
		)
	}

	// 2. 验证文件头
	if err := header.Validate(); err != nil {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"Header validation failed: "+err.Error(),
		)
	}

	// 3. 提取密文（跳过头部）
	headerSize := header.GetHeaderSize()
	if len(encryptedData) <= headerSize {
		return utils.NewCryptoError(
			utils.ErrInvalidFormat,
			"File too short - no ciphertext after header",
		)
	}
	ciphertext := encryptedData[headerSize:]

	// 4. 混合密钥解封装
	decryptor := NewHybridDecryptor(kyberPriv, ecdhPriv)
	sharedSecret, err := decryptor.Decapsulate(header.KyberEnc, header.ECDHPub[:])
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrAuthFailed,
			"Hybrid decapsulation failed: "+err.Error(),
		)
	}

	// 5. AES-GCM 解密数据
	plaintext, err := AESGCMDecrypt(sharedSecret, ciphertext, header.IV[:])
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrDecryptionFailed,
			"AES-GCM decryption failed: "+err.Error(),
		)
	}

	// 6. 验证 SHA256 哈希
	hash := sha256.Sum256(plaintext)
	if hash != header.SHA256Hash {
		return utils.NewCryptoError(
			utils.ErrHashMismatch,
			"SHA256 hash mismatch - file may be corrupted",
		)
	}

	// 7. 验证签名
	valid, err := VerifyHashSignature(hash[:], header.Signature, dilithiumPub)
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

	// 8. 写入解密文件
	err = os.WriteFile(outputPath, plaintext, 0644)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to write output file: "+err.Error(),
		)
	}

	return nil
}
