package crypto

import (
	"crypto/ecdh"
	"encoding/pem"
	"fmt"
	"os"
	"runtime"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// 保存密钥文件（遵循安全原则）
func SaveKeyFiles(
	kyberPub kem.PublicKey,
	ecdhPub *ecdh.PublicKey,
	kyberPriv kem.PrivateKey,
	ecdhPriv *ecdh.PrivateKey,
	pubPath, privPath string,
) error {
	// 导出公钥
	pubPEM, err := ExportPublicKey(kyberPub, ecdhPub)
	if err != nil {
		return err
	}

	// 导出私钥
	privPEM, err := ExportPrivateKey(kyberPriv, ecdhPriv)
	if err != nil {
		return err
	}

	// 保存公钥（默认权限）
	if err := os.WriteFile(pubPath, pubPEM, 0644); err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to save public key",
		)
	}

	// 保存私钥（严格权限0600）
	// 在 Windows 上，权限设置与 Unix 不同，需要特殊处理
	if err := os.WriteFile(privPath, privPEM, 0600); err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to save private key",
		)
	}

	// 在 Unix/Linux 系统上，确保私钥权限正确
	if runtime.GOOS != "windows" {
		if err := os.Chmod(privPath, 0600); err != nil {
			return utils.NewCryptoError(
				utils.ErrIOError,
				"Failed to set private key permissions",
			)
		}
	}

	return nil
}

// 加载密钥文件
func LoadKeyFiles(pubPath, privPath string) (*HybridPublicKey, *HybridPrivateKey, error) {
	pubPEM, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Public key file not found or unreadable",
		)
	}

	privPEM, err := os.ReadFile(privPath)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Private key file not found or unreadable",
		)
	}

	return ImportKeys(pubPEM, privPEM)
}

// LoadPublicKey 只加载公钥文件
func LoadPublicKey(pubPath string) (*HybridPublicKey, error) {
	pubPEM, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Public key file not found or unreadable",
		)
	}

	// 只解析公钥部分
	pubKyber, pubECDH, err := parsePublicKeys(pubPEM)
	if err != nil {
		return nil, err
	}

	return &HybridPublicKey{Kyber: pubKyber, ECDH: pubECDH}, nil
}

// LoadPrivateKey 只加载私钥文件
func LoadPrivateKey(privPath string) (*HybridPrivateKey, error) {
	privPEM, err := os.ReadFile(privPath)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Private key file not found or unreadable",
		)
	}

	// 只解析私钥部分
	privKyber, privECDH, err := parsePrivateKeys(privPEM)
	if err != nil {
		return nil, err
	}

	return &HybridPrivateKey{Kyber: privKyber, ECDH: privECDH}, nil
}

// DilithiumKeyPair 包含 Dilithium3 密钥对的 PEM 格式
type DilithiumKeyPair struct {
	Public  []byte
	Private []byte
}

// ExportDilithiumKeys 导出 Dilithium3 密钥对到 PEM 格式
func ExportDilithiumKeys(pub interface{}, priv interface{}) (*DilithiumKeyPair, error) {
	// 确保公钥是 Dilithium3 类型
	pubKey, ok := pub.(*mode3.PublicKey)
	if !ok {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 public key type",
		)
	}

	// 确保私钥是 Dilithium3 类型
	privKey, ok := priv.(*mode3.PrivateKey)
	if !ok {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 private key type",
		)
	}

	// 导出公钥
	pubBytes := pubKey.Bytes()
	pubPEM := pem.Block{
		Type:  "DILITHIUM3 PUBLIC KEY",
		Bytes: pubBytes,
	}

	// 导出私钥
	privBytes := privKey.Bytes()
	privPEM := pem.Block{
		Type:  "DILITHIUM3 PRIVATE KEY",
		Bytes: privBytes,
	}

	return &DilithiumKeyPair{
		Public:  pem.EncodeToMemory(&pubPEM),
		Private: pem.EncodeToMemory(&privPEM),
	}, nil
}

// ImportDilithiumKeys 从 PEM 格式导入 Dilithium3 密钥对
func ImportDilithiumKeys(pubPEM, privPEM []byte) (interface{}, interface{}, error) {
	// 解析公钥
	pubBlock, _ := pem.Decode(pubPEM)
	if pubBlock == nil || pubBlock.Type != "DILITHIUM3 PUBLIC KEY" {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 public key PEM",
		)
	}
	var pubKey mode3.PublicKey
	if err := pubKey.UnmarshalBinary(pubBlock.Bytes); err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			fmt.Sprintf("Failed to parse Dilithium3 public key: %v", err),
		)
	}

	// 解析私钥
	privBlock, _ := pem.Decode(privPEM)
	if privBlock == nil || privBlock.Type != "DILITHIUM3 PRIVATE KEY" {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 private key PEM",
		)
	}
	var privKey mode3.PrivateKey
	if err := privKey.UnmarshalBinary(privBlock.Bytes); err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			fmt.Sprintf("Failed to parse Dilithium3 private key: %v", err),
		)
	}

	return &pubKey, &privKey, nil
}

// LoadDilithiumKeys 从文件加载 Dilithium3 密钥对
func LoadDilithiumKeys(pubPath, privPath string) (interface{}, interface{}, error) {
	pubPEM, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Dilithium public key file not found or unreadable",
		)
	}

	privPEM, err := os.ReadFile(privPath)
	if err != nil {
		return nil, nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Dilithium private key file not found or unreadable",
		)
	}

	return ImportDilithiumKeys(pubPEM, privPEM)
}

// LoadDilithiumPublicKey 只加载 Dilithium 公钥
func LoadDilithiumPublicKey(pubPath string) (interface{}, error) {
	pubPEM, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Dilithium public key file not found or unreadable",
		)
	}

	// 解析公钥
	pubBlock, _ := pem.Decode(pubPEM)
	if pubBlock == nil || pubBlock.Type != "DILITHIUM3 PUBLIC KEY" {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 public key PEM",
		)
	}
	var pubKey mode3.PublicKey
	if err := pubKey.UnmarshalBinary(pubBlock.Bytes); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			fmt.Sprintf("Failed to parse Dilithium3 public key: %v", err),
		)
	}

	return &pubKey, nil
}

// LoadDilithiumPrivateKey 只加载 Dilithium 私钥
func LoadDilithiumPrivateKey(privPath string) (interface{}, error) {
	privPEM, err := os.ReadFile(privPath)
	if err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidParameter,
			"Dilithium private key file not found or unreadable",
		)
	}

	// 解析私钥
	privBlock, _ := pem.Decode(privPEM)
	if privBlock == nil || privBlock.Type != "DILITHIUM3 PRIVATE KEY" {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			"Invalid Dilithium3 private key PEM",
		)
	}
	var privKey mode3.PrivateKey
	if err := privKey.UnmarshalBinary(privBlock.Bytes); err != nil {
		return nil, utils.NewCryptoError(
			utils.ErrInvalidKey,
			fmt.Sprintf("Failed to parse Dilithium3 private key: %v", err),
		)
	}

	return &privKey, nil
}

// SaveDilithiumKeys 保存 Dilithium3 密钥对到文件
func SaveDilithiumKeys(pub interface{}, priv interface{}, pubPath, privPath string) error {
	keyPair, err := ExportDilithiumKeys(pub, priv)
	if err != nil {
		return err
	}

	// 保存公钥
	if err := os.WriteFile(pubPath, keyPair.Public, 0644); err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to save Dilithium public key",
		)
	}

	// 保存私钥
	if err := os.WriteFile(privPath, keyPair.Private, 0600); err != nil {
		return utils.NewCryptoError(
			utils.ErrIOError,
			"Failed to save Dilithium private key",
		)
	}

	// 在 Unix/Linux 系统上，确保私钥权限正确
	if runtime.GOOS != "windows" {
		if err := os.Chmod(privPath, 0600); err != nil {
			return utils.NewCryptoError(
				utils.ErrIOError,
				"Failed to set Dilithium private key permissions",
			)
		}
	}

	return nil
}
