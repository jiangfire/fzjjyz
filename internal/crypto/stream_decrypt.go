package crypto

import (
	"crypto/ecdh"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// StreamingDecryptor 流式解密器
// 支持大文件解密，内存占用仅与缓冲区大小相关
//
// 注意：虽然名为"流式"，但由于 AES-GCM 需要完整数据才能验证认证标签，
// 实际实现仍然需要一次性读取整个文件到内存。
// 真正的流式解密需要使用 AES-CTR + HMAC 等替代方案。
type StreamingDecryptor struct {
	kyberPriv    kem.PrivateKey
	ecdhPriv     *ecdh.PrivateKey
	dilithiumPub *mode3.PublicKey
	bufferSize   int
	pool         *BufferPool
}

// NewStreamingDecryptor 创建流式解密器
func NewStreamingDecryptor(
	kyberPriv kem.PrivateKey,
	ecdhPriv *ecdh.PrivateKey,
	dilithiumPub interface{},
	bufferSize int,
) (*StreamingDecryptor, error) {
	if bufferSize < MinBufferSize || bufferSize > MaxBufferSize {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Buffer size out of range",
		)
	}

	var pub *mode3.PublicKey
	if dilithiumPub != nil {
		var ok bool
		pub, ok = dilithiumPub.(*mode3.PublicKey)
		if !ok {
			return nil, utils.NewCryptoError(
				utils.ErrInvalidKey,
				"Invalid Dilithium3 public key",
			)
		}
	}

	return &StreamingDecryptor{
		kyberPriv:    kyberPriv,
		ecdhPriv:     ecdhPriv,
		dilithiumPub: pub,
		bufferSize:   bufferSize,
		pool:         NewBufferPool(bufferSize),
	}, nil
}

// DecryptFile 流式解密文件
// 使用核心解密逻辑，支持缓冲区池优化
func (sd *StreamingDecryptor) DecryptFile(inputPath, outputPath string) error {
	// 调用核心解密逻辑
	plaintext, err := DecryptFileCore(inputPath, sd.kyberPriv, sd.ecdhPriv, sd.dilithiumPub)
	if err != nil {
		return err
	}

	// 写入解密文件
	return writeDecryptedFile(outputPath, plaintext)
}
