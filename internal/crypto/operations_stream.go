package crypto

import (
	"crypto/ecdh"
	"os"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
)

// EncryptFileStreaming 流式加密文件
// 这是流式处理的入口函数，提供与原 EncryptFile 兼容的接口
//
// 注意：由于 AES-GCM 需要完整数据，实际实现仍然需要读取整个文件到内存。
// 但提供了缓冲区池优化和统一的接口设计，便于未来扩展真正的流式处理。
//
// 参数：
//   inputPath: 输入文件路径
//   outputPath: 输出文件路径
//   kyberPub: Kyber 公钥
//   ecdhPub: ECDH 公钥
//   dilithiumPriv: Dilithium 私钥（可选，传 nil 跳过签名）
//   bufferSize: 缓冲区大小（建议使用 OptimalBufferSize 自动选择）
//
// 返回：
//   error: 错误信息
func EncryptFileStreaming(
	inputPath, outputPath string,
	kyberPub kem.PublicKey,
	ecdhPub *ecdh.PublicKey,
	dilithiumPriv interface{},
	bufferSize int,
) error {
	encryptor, err := NewStreamingEncryptor(kyberPub, ecdhPub, dilithiumPriv, bufferSize)
	if err != nil {
		return err
	}
	return encryptor.EncryptFile(inputPath, outputPath)
}

// DecryptFileStreaming 流式解密文件
// 这是流式处理的入口函数，提供与原 DecryptFile 兼容的接口
//
// 注意：由于 AES-GCM 需要完整数据，实际实现仍然需要读取整个文件到内存。
// 但提供了缓冲区池优化和统一的接口设计，便于未来扩展真正的流式处理。
//
// 参数：
//   inputPath: 加密文件路径
//   outputPath: 输出文件路径
//   kyberPriv: Kyber 私钥
//   ecdhPriv: ECDH 私钥
//   dilithiumPub: Dilithium 公钥（可选，传 nil 跳过签名验证）
//   bufferSize: 缓冲区大小（建议使用 OptimalBufferSize 自动选择）
//
// 返回：
//   error: 错误信息
func DecryptFileStreaming(
	inputPath, outputPath string,
	kyberPriv kem.PrivateKey,
	ecdhPriv *ecdh.PrivateKey,
	dilithiumPub interface{},
	bufferSize int,
) error {
	decryptor, err := NewStreamingDecryptor(kyberPriv, ecdhPriv, dilithiumPub, bufferSize)
	if err != nil {
		return err
	}
	return decryptor.DecryptFile(inputPath, outputPath)
}

// EncryptFileStreamingAuto 自动选择缓冲区大小的流式加密
// 这是推荐的使用方式，自动根据文件大小选择最优缓冲区
func EncryptFileStreamingAuto(
	inputPath, outputPath string,
	kyberPub kem.PublicKey,
	ecdhPub *ecdh.PublicKey,
	dilithiumPriv interface{},
) error {
	// 获取文件大小
	info, err := os.Stat(inputPath)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to get file info: "+err.Error(),
		)
	}

	bufferSize := OptimalBufferSize(info.Size())
	return EncryptFileStreaming(inputPath, outputPath, kyberPub, ecdhPub, dilithiumPriv, bufferSize)
}

// DecryptFileStreamingAuto 自动选择缓冲区大小的流式解密
// 这是推荐的使用方式，自动根据文件大小选择最优缓冲区
func DecryptFileStreamingAuto(
	inputPath, outputPath string,
	kyberPriv kem.PrivateKey,
	ecdhPriv *ecdh.PrivateKey,
	dilithiumPub interface{},
) error {
	// 获取文件大小
	info, err := os.Stat(inputPath)
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to get file info: "+err.Error(),
		)
	}

	bufferSize := OptimalBufferSize(info.Size())
	return DecryptFileStreaming(inputPath, outputPath, kyberPriv, ecdhPriv, dilithiumPub, bufferSize)
}
