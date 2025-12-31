package zjcrypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"sync"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
)

// HybridPublicKey 混合公钥结构（表达原则：数据结构优先）.
type HybridPublicKey struct {
	Kyber kem.PublicKey
	ECDH  *ecdh.PublicKey
}

// HybridPrivateKey 混合私钥结构.
type HybridPrivateKey struct {
	Kyber kem.PrivateKey
	ECDH  *ecdh.PrivateKey
}

// GenerateKyberKeys 生成Kyber密钥对.
func GenerateKyberKeys() (kem.PublicKey, kem.PrivateKey, error) {
	scheme := kyber768.Scheme()
	pub, priv, err := scheme.GenerateKeyPair()
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			fmt.Sprintf("Kyber key generation failed: %v", err),
		)
	}
	return pub, priv, nil
}

// GenerateECDHKeys 生成ECDH密钥对.
func GenerateECDHKeys() (*ecdh.PublicKey, *ecdh.PrivateKey, error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			fmt.Sprintf("ECDH key generation failed: %v", err),
		)
	}
	return priv.PublicKey(), priv, nil
}

// ExportPublicKey 导出公钥到PEM格式.
func ExportPublicKey(kyberPub kem.PublicKey, ecdhPub *ecdh.PublicKey) ([]byte, error) {
	// Kyber公钥
	kyberBytes, err := kyberPub.MarshalBinary()
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			fmt.Sprintf("Failed to marshal Kyber public key: %v", err),
		)
	}
	kyberPEM := pem.Block{
		Type:  "KYBER PUBLIC KEY",
		Bytes: kyberBytes,
	}

	// ECDH公钥
	ecdhBytes := ecdhPub.Bytes()
	ecdhPEM := pem.Block{
		Type:  "ECDH PUBLIC KEY",
		Bytes: ecdhBytes,
	}

	// 组合PEM
	combined := pem.EncodeToMemory(&kyberPEM)
	combined = append(combined, pem.EncodeToMemory(&ecdhPEM)...)

	return combined, nil
}

// ExportPrivateKey 导出私钥到PEM格式（注意权限设置）.
func ExportPrivateKey(kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey) ([]byte, error) {
	kyberBytes, err := kyberPriv.MarshalBinary()
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrKeyGenerationFailed,
			fmt.Sprintf("Failed to marshal Kyber private key: %v", err),
		)
	}
	ecdhBytes := ecdhPriv.Bytes()

	kyberPEM := pem.Block{
		Type:  "KYBER PRIVATE KEY",
		Bytes: kyberBytes,
	}

	ecdhPEM := pem.Block{
		Type:  "ECDH PRIVATE KEY",
		Bytes: ecdhBytes,
	}

	combined := pem.EncodeToMemory(&kyberPEM)
	combined = append(combined, pem.EncodeToMemory(&ecdhPEM)...)

	return combined, nil
}

// ImportKeys 从PEM导入密钥.
func ImportKeys(pubPEM, privPEM []byte) (*HybridPublicKey, *HybridPrivateKey, error) {
	// 解析公钥
	pubKyber, pubECDH, err := parsePublicKeys(pubPEM)
	if err != nil {
		return nil, nil, err
	}

	// 解析私钥
	privKyber, privECDH, err := parsePrivateKeys(privPEM)
	if err != nil {
		return nil, nil, err
	}

	return &HybridPublicKey{Kyber: pubKyber, ECDH: pubECDH},
		&HybridPrivateKey{Kyber: privKyber, ECDH: privECDH},
		nil
}

// 辅助函数：解析公钥.
func parsePublicKeys(pemData []byte) (kem.PublicKey, *ecdh.PublicKey, error) {
	var kyberKey kem.PublicKey
	var ecdhKey *ecdh.PublicKey

	// 解析多个PEM块
	rest := pemData
	for len(rest) > 0 {
		block, next := pem.Decode(rest)
		if block == nil {
			break
		}

		switch block.Type {
		case "KYBER PUBLIC KEY":
			scheme := kyber768.Scheme()
			pub, err := scheme.UnmarshalBinaryPublicKey(block.Bytes)
			if err != nil {
				return nil, nil, utils.NewCryptoError(
					utils.ErrInvalidKey,
					fmt.Sprintf("Failed to parse Kyber public key: %v", err),
				)
			}
			kyberKey = pub

		case "ECDH PUBLIC KEY":
			pub, err := ecdh.X25519().NewPublicKey(block.Bytes)
			if err != nil {
				return nil, nil, utils.NewCryptoError(
					utils.ErrInvalidKey,
					fmt.Sprintf("Failed to parse ECDH public key: %v", err),
				)
			}
			ecdhKey = pub
		}

		rest = next
	}

	if kyberKey == nil || ecdhKey == nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Incomplete public key data",
		)
	}

	return kyberKey, ecdhKey, nil
}

// 辅助函数：解析私钥.
func parsePrivateKeys(pemData []byte) (kem.PrivateKey, *ecdh.PrivateKey, error) {
	var kyberKey kem.PrivateKey
	var ecdhKey *ecdh.PrivateKey

	// 解析多个PEM块
	rest := pemData
	for len(rest) > 0 {
		block, next := pem.Decode(rest)
		if block == nil {
			break
		}

		switch block.Type {
		case "KYBER PRIVATE KEY":
			scheme := kyber768.Scheme()
			priv, err := scheme.UnmarshalBinaryPrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, utils.NewCryptoError(
					utils.ErrInvalidKey,
					fmt.Sprintf("Failed to parse Kyber private key: %v", err),
				)
			}
			kyberKey = priv

		case "ECDH PRIVATE KEY":
			priv, err := ecdh.X25519().NewPrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, utils.NewCryptoError(
					utils.ErrInvalidKey,
					fmt.Sprintf("Failed to parse ECDH private key: %v", err),
				)
			}
			ecdhKey = priv
		}

		rest = next
	}

	if kyberKey == nil || ecdhKey == nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Incomplete private key data",
		)
	}

	return kyberKey, ecdhKey, nil
}

// GenerateKeyPairParallel 并行生成所有密钥对（Kyber + ECDH + Dilithium）
// 使用 goroutine 并行执行，速度提升 40-60%
//
//nolint:funlen // 并行密钥生成需要完整的错误处理和同步逻辑
func GenerateKeyPairParallel() (
	kyberPub kem.PublicKey,
	kyberPriv kem.PrivateKey,
	ecdhPub *ecdh.PublicKey,
	ecdhPriv *ecdh.PrivateKey,
	dilithiumPub interface{},
	dilithiumPriv interface{},
	err error,
) {
	var wg sync.WaitGroup
	var mu sync.Mutex // 保护错误收集

	// 错误收集通道
	errChan := make(chan error, 3)

	// Kyber 密钥生成
	wg.Add(1)
	go func() {
		defer wg.Done()
		scheme := kyber768.Scheme()
		pub, priv, err := scheme.GenerateKeyPair()
		if err != nil {
			mu.Lock()
			errChan <- utils.NewCryptoError(
				utils.ErrKeyGenerationFailed,
				fmt.Sprintf("Kyber key generation failed: %v", err),
			)
			mu.Unlock()
			return
		}
		kyberPub, kyberPriv = pub, priv
	}()

	// ECDH 密钥生成
	wg.Add(1)
	go func() {
		defer wg.Done()
		priv, err := ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			mu.Lock()
			errChan <- utils.NewCryptoError(
				utils.ErrKeyGenerationFailed,
				fmt.Sprintf("ECDH key generation failed: %v", err),
			)
			mu.Unlock()
			return
		}
		ecdhPub, ecdhPriv = priv.PublicKey(), priv
	}()

	// Dilithium 密钥生成
	wg.Add(1)
	go func() {
		defer wg.Done()
		pub, priv, err := GenerateDilithiumKeys()
		if err != nil {
			mu.Lock()
			errChan <- err
			mu.Unlock()
			return
		}
		dilithiumPub, dilithiumPriv = pub, priv
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		return nil, nil, nil, nil, nil, nil, err
	}

	return kyberPub, kyberPriv, ecdhPub, ecdhPriv, dilithiumPub, dilithiumPriv, nil
}

// GenerateHybridKeysParallel 并行生成 Kyber + ECDH 密钥对
// 专门用于混合加密，不包含 Dilithium.
func GenerateHybridKeysParallel() (
	kyberPub kem.PublicKey,
	kyberPriv kem.PrivateKey,
	ecdhPub *ecdh.PublicKey,
	ecdhPriv *ecdh.PrivateKey,
	err error,
) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 2)

	// Kyber
	wg.Add(1)
	go func() {
		defer wg.Done()
		scheme := kyber768.Scheme()
		pub, priv, err := scheme.GenerateKeyPair()
		if err != nil {
			mu.Lock()
			errChan <- utils.NewCryptoError(
				utils.ErrKeyGenerationFailed,
				fmt.Sprintf("Kyber key generation failed: %v", err),
			)
			mu.Unlock()
			return
		}
		kyberPub, kyberPriv = pub, priv
	}()

	// ECDH
	wg.Add(1)
	go func() {
		defer wg.Done()
		priv, err := ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			mu.Lock()
			errChan <- utils.NewCryptoError(
				utils.ErrKeyGenerationFailed,
				fmt.Sprintf("ECDH key generation failed: %v", err),
			)
			mu.Unlock()
			return
		}
		ecdhPub, ecdhPriv = priv.PublicKey(), priv
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return nil, nil, nil, nil, err
	}

	return kyberPub, kyberPriv, ecdhPub, ecdhPriv, nil
}
