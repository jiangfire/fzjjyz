package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/pem"
	"fmt"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"codeberg.org/jiangfire/fzjjyz/internal/utils"
)

// 密钥对结构（表达原则：数据结构优先）
type HybridPublicKey struct {
	Kyber kem.PublicKey
	ECDH  *ecdh.PublicKey
}

type HybridPrivateKey struct {
	Kyber kem.PrivateKey
	ECDH  *ecdh.PrivateKey
}

// 生成Kyber密钥对
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

// 生成ECDH密钥对
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

// 导出公钥到PEM格式
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

// 导出私钥到PEM格式（注意权限设置）
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

// 从PEM导入密钥
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

// 辅助函数：解析公钥
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

// 辅助函数：解析私钥
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
