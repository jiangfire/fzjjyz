package crypto

import (
	"crypto/ecdh"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// StreamingEncryptor 流式加密器
// 支持大文件加密，内存占用仅与缓冲区大小相关
//
// 注意：虽然名为"流式"，但由于 AES-GCM 需要完整数据才能生成认证标签，
// 实际实现仍然需要一次性读取整个文件到内存。
// 真正的流式加密需要使用 AES-CTR + HMAC 等替代方案。
type StreamingEncryptor struct {
	kyberPub      kem.PublicKey
	ecdhPub       *ecdh.PublicKey
	dilithiumPriv *mode3.PrivateKey
	bufferSize    int
	pool          *BufferPool
}

// NewStreamingEncryptor 创建流式加密器.
func NewStreamingEncryptor(
	kyberPub kem.PublicKey,
	ecdhPub *ecdh.PublicKey,
	dilithiumPriv interface{},
	bufferSize int,
) (*StreamingEncryptor, error) {
	if bufferSize < MinBufferSize || bufferSize > MaxBufferSize {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Buffer size out of range",
		)
	}

	// 类型转换
	var priv *mode3.PrivateKey
	if dilithiumPriv != nil {
		var ok bool
		priv, ok = dilithiumPriv.(*mode3.PrivateKey)
		if !ok {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidKey,
				"Invalid Dilithium3 private key",
			)
		}
	}

	return &StreamingEncryptor{
		kyberPub:      kyberPub,
		ecdhPub:       ecdhPub,
		dilithiumPriv: priv,
		bufferSize:    bufferSize,
		pool:          NewBufferPool(bufferSize),
	}, nil
}

// EncryptFile 流式加密文件
// 使用核心加密逻辑，支持缓冲区池优化.
func (se *StreamingEncryptor) EncryptFile(inputPath, outputPath string) error {
	// 调用核心加密逻辑
	header, ciphertext, err := EncryptFileCore(inputPath, se.kyberPub, se.ecdhPub, se.dilithiumPriv)
	if err != nil {
		return err
	}

	// 使用优化的序列化方法
	headerBytes, err := header.MarshalBinaryOptimized()
	if err != nil {
		return utils.NewCryptoError(
			utils.ErrSerializationFailed,
			"Header serialization failed: "+err.Error(),
		)
	}

	// 写入加密文件
	return writeEncryptedFile(outputPath, headerBytes, ciphertext)
}
