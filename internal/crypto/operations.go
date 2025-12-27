package crypto

import (
	"crypto/ecdh"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
)

// EncryptFile 加密文件
// 输入: 原始文件路径, 输出文件路径, Kyber公钥, ECDH公钥, Dilithium私钥
// 返回: 错误
//
// 加密流程:
// 1. 读取原始文件
// 2. 混合密钥封装 (Kyber + ECDH)
// 3. AES-256-GCM 加密
// 4. 计算 SHA256 哈希
// 5. Dilithium3 签名
// 6. 构建并序列化文件头
// 7. 写入 [头部] + [密文]
func EncryptFile(inputPath, outputPath string, kyberPub kem.PublicKey, ecdhPub *ecdh.PublicKey, dilithiumPriv interface{}) error {
	// 调用核心加密逻辑
	header, ciphertext, err := EncryptFileCore(inputPath, kyberPub, ecdhPub, dilithiumPriv)
	if err != nil {
		return err
	}

	// 序列化头部
	headerBytes, err := serializeHeader(header)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrSerializationFailed,
			"Header serialization failed: "+err.Error(),
		)
	}

	// 写入加密文件
	return writeEncryptedFile(outputPath, headerBytes, ciphertext)
}

// DecryptFile 解密文件
// 输入: 加密文件路径, 输出文件路径, Kyber私钥, ECDH私钥, Dilithium公钥
// 返回: 错误
//
// 解密流程:
// 1. 读取并解析文件头
// 2. 验证文件头
// 3. 混合密钥解封装
// 4. AES-256-GCM 解密
// 5. 验证 SHA256 哈希
// 6. 验证 Dilithium3 签名
// 7. 写入解密文件
func DecryptFile(inputPath, outputPath string, kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, dilithiumPub interface{}) error {
	// 调用核心解密逻辑
	plaintext, err := DecryptFileCore(inputPath, kyberPriv, ecdhPriv, dilithiumPub)
	if err != nil {
		return err
	}

	// 写入解密文件
	return writeDecryptedFile(outputPath, plaintext)
}
