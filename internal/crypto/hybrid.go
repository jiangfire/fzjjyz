package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
)

// HybridEncryptor 混合加密器
// 结合后量子密码（Kyber）和传统密码（ECDH）提供双重保护.
type HybridEncryptor struct {
	kyberPub kem.PublicKey
	ecdhPub  *ecdh.PublicKey
}

// HybridDecryptor 混合解密器.
type HybridDecryptor struct {
	kyberPriv kem.PrivateKey
	ecdhPriv  *ecdh.PrivateKey
}

// NewHybridEncryptor 创建混合加密器.
func NewHybridEncryptor(kyberPub kem.PublicKey, ecdhPub *ecdh.PublicKey) *HybridEncryptor {
	return &HybridEncryptor{
		kyberPub: kyberPub,
		ecdhPub:  ecdhPub,
	}
}

// NewHybridDecryptor 创建混合解密器.
func NewHybridDecryptor(kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey) *HybridDecryptor {
	return &HybridDecryptor{
		kyberPriv: kyberPriv,
		ecdhPriv:  ecdhPriv,
	}
}

// Encapsulate 执行混合密钥封装
// 返回: Kyber 密文 (1088B), 临时 ECDH 公钥 (32B), 组合共享密钥 (32B)
//
// 加密流程:
// 1. Kyber768 封装 → 1088B 密文 + 32B Kyber 共享密钥
// 2. 生成临时 ECDH 密钥对 → 32B 临时公钥 + 32B ECDH 共享密钥
// 3. SHA256(Kyber密钥 + ECDH密钥) → 32B 最终密钥.
func (e *HybridEncryptor) Encapsulate() (encapsulated []byte, ecdhPub []byte, sharedSecret []byte, err error) {
	// 步骤1: Kyber 封装
	kyberScheme := kyber768.Scheme()
	encapsulated, kyberSecret, err := kyberScheme.Encapsulate(e.kyberPub)
	if err != nil {
		return nil, nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"Kyber encapsulation failed",
		)
	}

	// 步骤2: ECDH 临时密钥对生成和交换
	ecdhKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"ECDH key generation failed",
		)
	}
	ecdhSecret, err := ecdhKey.ECDH(e.ecdhPub)
	if err != nil {
		return nil, nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"ECDH key exchange failed",
		)
	}

	// 返回临时 ECDH 公钥（需要存储在文件头中）
	ecdhPubBytes := ecdhKey.PublicKey().Bytes()

	// 步骤3: 组合共享密钥
	combined := sha256.Sum256(append(kyberSecret, ecdhSecret...))
	return encapsulated, ecdhPubBytes, combined[:], nil
}

// Decapsulate 执行混合密钥解封装
// 输入: Kyber 密文 (1088B), 临时 ECDH 公钥 (32B)
// 返回: 组合共享密钥 (32B)
//
// 解密流程:
// 1. Kyber768 解封装 → 32B Kyber 共享密钥
// 2. ECDH X25519 密钥交换（使用临时公钥）→ 32B ECDH 共享密钥
// 3. SHA256(Kyber密钥 + ECDH密钥) → 32B 最终密钥.
func (d *HybridDecryptor) Decapsulate(encapsulated []byte, ecdhPub []byte) ([]byte, error) {
	// 步骤1: Kyber 解封装
	kyberScheme := kyber768.Scheme()
	kyberSecret, err := kyberScheme.Decapsulate(d.kyberPriv, encapsulated)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrAuthFailed,
			"Kyber decapsulation failed",
		)
	}

	// 步骤2: ECDH 密钥交换
	// 解析加密器的临时 ECDH 公钥
	tempECDHPub, err := ecdh.X25519().NewPublicKey(ecdhPub)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid temporary ECDH public key",
		)
	}
	ecdhSecret, err := d.ecdhPriv.ECDH(tempECDHPub)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrAuthFailed,
			"ECDH key exchange failed",
		)
	}

	// 步骤3: 组合共享密钥
	combined := sha256.Sum256(append(kyberSecret, ecdhSecret...))
	return combined[:], nil
}

// AESGCMEncrypt 使用 AES-256-GCM 加密数据
// 输入: 32B 密钥, 明文数据
// 返回: 密文, 12B Nonce, 错误
//
// AES-GCM 特性:
// - 机密性: AES-256 加密
// - 完整性: GCM 认证标签
// - 防重放: 随机 Nonce.
func AESGCMEncrypt(key []byte, plaintext []byte) (ciphertext []byte, nonce []byte, err error) {
	// 验证密钥长度 (必须是 32B 用于 AES-256)
	if len(key) != 32 {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid AES key length, expected 32 bytes",
		)
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid AES key",
		)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"GCM mode failed",
		)
	}

	// 生成随机 Nonce (12B)
	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"Nonce generation failed",
		)
	}

	// 加密
	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// AESGCMDecrypt 使用 AES-256-GCM 解密数据
// 输入: 32B 密钥, 密文, 12B Nonce
// 返回: 明文, 错误
//
// 自动验证数据完整性和真实性.
func AESGCMDecrypt(key []byte, ciphertext []byte, nonce []byte) ([]byte, error) {
	// 验证密钥长度
	if len(key) != 32 {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid AES key length, expected 32 bytes",
		)
	}

	// 验证 nonce 长度
	if len(nonce) != 12 {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidData,
			"Invalid nonce length, expected 12 bytes",
		)
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid AES key",
		)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"GCM mode failed",
		)
	}

	// 解密（自动验证认证标签）
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrAuthFailed,
			"Authentication failed - data may be tampered or invalid key",
		)
	}

	return plaintext, nil
}
