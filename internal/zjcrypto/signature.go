package zjcrypto

import (
	"crypto/rand"
	"fmt"
	"os"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// SignData 使用 Dilithium3 (mode3) 对数据进行签名
// 输入: 数据, Dilithium3 私钥
// 返回: 签名 (3293B), 错误.
func SignData(data []byte, privKey interface{}) (signature []byte, err error) {
	// 确保私钥是 Dilithium3 类型
	priv, ok := privKey.(*mode3.PrivateKey)
	if !ok {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 private key type",
		)
	}

	// 生成签名
	signature = make([]byte, mode3.SignatureSize)
	mode3.SignTo(priv, data, signature)
	return signature, nil
}

// VerifySignature 验证 Dilithium3 签名
// 输入: 数据, 签名 (2420B), Dilithium3 公钥
// 返回: bool (true = 验证通过), 错误.
func VerifySignature(data []byte, signature []byte, pubKey interface{}) (bool, error) {
	// 确保公钥是 Dilithium3 类型
	pub, ok := pubKey.(*mode3.PublicKey)
	if !ok {
		return false, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 public key type",
		)
	}

	// 验证签名
	valid := mode3.Verify(pub, data, signature)
	return valid, nil
}

// SignFile 对文件数据进行签名
// 输入: 文件路径, Dilithium3 私钥
// 返回: 签名 (2420B), 错误.
func SignFile(filePath string, privKey interface{}) (signature []byte, err error) {
	// 读取文件
	data, err := readFileData(filePath)
	if err != nil {
		return nil, err
	}

	return SignData(data, privKey)
}

// VerifyFileSignature 验证文件签名
// 输入: 文件路径, 签名, Dilithium3 公钥
// 返回: bool, 错误.
func VerifyFileSignature(filePath string, signature []byte, pubKey interface{}) (bool, error) {
	// 读取文件
	data, err := readFileData(filePath)
	if err != nil {
		return false, err
	}

	return VerifySignature(data, signature, pubKey)
}

// SignHash 对哈希值进行签名（用于文件加密流程）
// 输入: 32B 哈希值, Dilithium3 私钥
// 返回: 签名 (2420B), 错误.
func SignHash(hash []byte, privKey interface{}) (signature []byte, err error) {
	if len(hash) != 32 {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Hash must be 32 bytes",
		)
	}

	return SignData(hash, privKey)
}

// VerifyHashSignature 验证哈希签名
// 输入: 32B 哈希值, 签名, Dilithium3 公钥
// 返回: bool, 错误.
func VerifyHashSignature(hash []byte, signature []byte, pubKey interface{}) (bool, error) {
	if len(hash) != 32 {
		return false, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Hash must be 32 bytes",
		)
	}

	return VerifySignature(hash, signature, pubKey)
}

// 辅助函数：读取文件数据.
func readFileData(filePath string) ([]byte, error) {
	// G304: filePath 应由调用方验证
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return data, nil
}

// GenerateDilithiumKeys 生成 Dilithium3 密钥对
// 返回: 公钥, 私钥, 错误.
func GenerateDilithiumKeys() (interface{}, interface{}, error) {
	pub, priv, err := mode3.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			"Dilithium3 key generation failed",
		)
	}
	return pub, priv, nil
}

// DilithiumSignatureSize 返回 Dilithium3 签名大小.
func DilithiumSignatureSize() int {
	return mode3.SignatureSize
}

// DilithiumPublicKeySize 返回 Dilithium3 公钥大小.
func DilithiumPublicKeySize() int {
	return mode3.PublicKeySize
}

// DilithiumPrivateKeySize 返回 Dilithium3 私钥大小.
func DilithiumPrivateKeySize() int {
	return mode3.PrivateKeySize
}

// DilithiumGetPublicKey 从私钥获取公钥
// 输入: Dilithium3 私钥
// 返回: Dilithium3 公钥.
func DilithiumGetPublicKey(privKey interface{}) interface{} {
	priv, ok := privKey.(*mode3.PrivateKey)
	if !ok {
		return nil
	}
	return priv.Public()
}
