package zjcrypto

import (
	"crypto/ecdh"
	"encoding/pem"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"codeberg.org/jiangfire/fzjjyz/internal/utils"
	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/sign/dilithium/mode3"
)

// 缓存配置常量.
const (
	// MaxCacheSize 缓存最大容量（密钥数量）.
	MaxCacheSize = 100
	// DefaultCacheTTL 默认缓存过期时间（1小时）.
	DefaultCacheTTL = 1 * time.Hour
	// CacheCleanupInterval 缓存清理间隔.
	CacheCleanupInterval = 5 * time.Minute
	// windowsOS 操作系统名称常量，避免重复字符串。
	windowsOS = "windows"
	// keyFilePerm 私钥文件权限。
	keyFilePerm = 0600
	// pubKeyFilePerm 公钥文件权限。
	pubKeyFilePerm = 0644
)

// keyCacheEntry 缓存条目，包含创建时间和TTL.
type keyCacheEntry struct {
	key       interface{}
	createdAt time.Time
	ttl       time.Duration
}

// keyCache 密钥缓存，使用 sync.Map 实现线程安全.
var keyCache sync.Map

// init 初始化缓存清理定时器
//
//nolint:gochecknoinits // init函数用于必要的包级初始化：启动后台缓存清理任务
func init() {
	startCacheCleanup()
}

// startCacheCleanup 启动后台缓存清理任务.
func startCacheCleanup() {
	// 使用匿名函数避免全局变量，定时器会自动重新调度
	_ = time.AfterFunc(CacheCleanupInterval, func() {
		cleanupExpiredKeys()
		// 重新调度下一次清理
		startCacheCleanup()
	})
}

// cleanupExpiredKeys 清理过期的缓存条目.
func cleanupExpiredKeys() {
	now := time.Now()
	keyCache.Range(func(key, value interface{}) bool {
		entry := value.(*keyCacheEntry)
		if now.Sub(entry.createdAt) >= entry.ttl {
			keyCache.Delete(key)
		}
		return true
	})
}

// SaveKeyFiles 保存密钥文件（遵循安全原则）.
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
	if err := os.WriteFile(pubPath, pubPEM, pubKeyFilePerm); err != nil {
		return fmt.Errorf("save public key: %w", err)
	}

	// 保存私钥（严格权限0600）
	// 在 Windows 上，权限设置与 Unix 不同，需要特殊处理
	if err := os.WriteFile(privPath, privPEM, keyFilePerm); err != nil {
		return fmt.Errorf("save private key: %w", err)
	}

	// 在 Unix/Linux 系统上，确保私钥权限正确
	if runtime.GOOS != windowsOS {
		if err := os.Chmod(privPath, keyFilePerm); err != nil {
			return fmt.Errorf("set private key permissions: %w", err)
		}
	}

	return nil
}

// LoadKeyFiles 加载密钥文件.
func LoadKeyFiles(pubPath, privPath string) (*HybridPublicKey, *HybridPrivateKey, error) {
	// G304: 调用方应验证路径安全性
	pubPEM, err := os.ReadFile(pubPath) //nolint:gosec
	if err != nil {
		return nil, nil, fmt.Errorf("read public key file: %w", err)
	}

	// G304: 调用方应验证路径安全性
	privPEM, err := os.ReadFile(privPath) //nolint:gosec
	if err != nil {
		return nil, nil, fmt.Errorf("read private key file: %w", err)
	}

	return ImportKeys(pubPEM, privPEM)
}

// LoadPublicKey 只加载公钥文件.
func LoadPublicKey(pubPath string) (*HybridPublicKey, error) {
	// G304: 调用方应验证路径安全性
	pubPEM, err := os.ReadFile(pubPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("read public key file: %w", err)
	}

	// 只解析公钥部分
	pubKyber, pubECDH, err := parsePublicKeys(pubPEM)
	if err != nil {
		return nil, err
	}

	return &HybridPublicKey{Kyber: pubKyber, ECDH: pubECDH}, nil
}

// LoadPrivateKey 只加载私钥文件.
func LoadPrivateKey(privPath string) (*HybridPrivateKey, error) {
	// G304: 调用方应验证路径安全性
	privPEM, err := os.ReadFile(privPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("read private key file: %w", err)
	}

	// 只解析私钥部分
	privKyber, privECDH, err := parsePrivateKeys(privPEM)
	if err != nil {
		return nil, err
	}

	return &HybridPrivateKey{Kyber: privKyber, ECDH: privECDH}, nil
}

// DilithiumKeyPair 包含 Dilithium3 密钥对的 PEM 格式.
type DilithiumKeyPair struct {
	Public  []byte
	Private []byte
}

// ExportDilithiumKeys 导出 Dilithium3 密钥对到 PEM 格式.
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

// ImportDilithiumKeys 从 PEM 格式导入 Dilithium3 密钥对.
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

// LoadDilithiumKeys 从文件加载 Dilithium3 密钥对.
func LoadDilithiumKeys(pubPath, privPath string) (interface{}, interface{}, error) {
	// G304: 调用方应验证路径安全性
	pubPEM, err := os.ReadFile(pubPath) //nolint:gosec
	if err != nil {
		return nil, nil, fmt.Errorf("read Dilithium public key file: %w", err)
	}

	// G304: 调用方应验证路径安全性
	privPEM, err := os.ReadFile(privPath) //nolint:gosec
	if err != nil {
		return nil, nil, fmt.Errorf("read Dilithium private key file: %w", err)
	}

	return ImportDilithiumKeys(pubPEM, privPEM)
}

// LoadDilithiumPublicKey 只加载 Dilithium 公钥.
func LoadDilithiumPublicKey(pubPath string) (interface{}, error) {
	// G304: 调用方应验证路径安全性
	pubPEM, err := os.ReadFile(pubPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("read Dilithium public key file: %w", err)
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

// LoadDilithiumPrivateKey 只加载 Dilithium 私钥.
func LoadDilithiumPrivateKey(privPath string) (interface{}, error) {
	// G304: 调用方应验证路径安全性
	privPEM, err := os.ReadFile(privPath) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("read Dilithium private key file: %w", err)
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

// checkCacheSize 检查缓存大小，如果超过限制则清理最旧的条目.
func checkCacheSize() {
	count := 0
	keyCache.Range(func(_, _ interface{}) bool {
		count++
		return true
	})

	if count >= MaxCacheSize {
		// 清理最旧的条目（简化策略：清理20%）
		toDelete := count / 5
		if toDelete == 0 {
			toDelete = 1
		}

		deleted := 0
		now := time.Now()
		keyCache.Range(func(key, value interface{}) bool {
			if deleted >= toDelete {
				return false
			}
			entry := value.(*keyCacheEntry)
			// 优先清理过期的，然后清理最旧的
			if now.Sub(entry.createdAt) >= entry.ttl || deleted < toDelete {
				keyCache.Delete(key)
				deleted++
			}
			return true
		})
	}
}

// LoadPublicKeyCached 带缓存的公钥加载（支持TTL和大小限制）
// 第一次加载时从文件读取并缓存，后续直接返回缓存结果.
func LoadPublicKeyCached(path string) (*HybridPublicKey, error) {
	cacheKey := "pub:" + path
	if cached, ok := keyCache.Load(cacheKey); ok {
		entry := cached.(*keyCacheEntry)
		if time.Since(entry.createdAt) < entry.ttl {
			return entry.key.(*HybridPublicKey), nil
		}
		// 已过期，删除
		keyCache.Delete(cacheKey)
	}

	// 检查缓存大小
	checkCacheSize()

	key, err := LoadPublicKey(path)
	if err != nil {
		return nil, err
	}

	keyCache.Store(cacheKey, &keyCacheEntry{
		key:       key,
		createdAt: time.Now(),
		ttl:       DefaultCacheTTL,
	})
	return key, nil
}

// LoadPrivateKeyCached 带缓存的私钥加载（支持TTL和大小限制）.
func LoadPrivateKeyCached(path string) (*HybridPrivateKey, error) {
	cacheKey := "priv:" + path
	if cached, ok := keyCache.Load(cacheKey); ok {
		entry := cached.(*keyCacheEntry)
		if time.Since(entry.createdAt) < entry.ttl {
			return entry.key.(*HybridPrivateKey), nil
		}
		keyCache.Delete(cacheKey)
	}

	checkCacheSize()

	key, err := LoadPrivateKey(path)
	if err != nil {
		return nil, err
	}

	keyCache.Store(cacheKey, &keyCacheEntry{
		key:       key,
		createdAt: time.Now(),
		ttl:       DefaultCacheTTL,
	})
	return key, nil
}

// LoadDilithiumPublicKeyCached 带缓存的 Dilithium 公钥加载（支持TTL和大小限制）.
func LoadDilithiumPublicKeyCached(path string) (interface{}, error) {
	cacheKey := "dilithium_pub:" + path
	if cached, ok := keyCache.Load(cacheKey); ok {
		entry := cached.(*keyCacheEntry)
		if time.Since(entry.createdAt) < entry.ttl {
			return entry.key, nil
		}
		keyCache.Delete(cacheKey)
	}

	checkCacheSize()

	key, err := LoadDilithiumPublicKey(path)
	if err != nil {
		return nil, err
	}

	keyCache.Store(cacheKey, &keyCacheEntry{
		key:       key,
		createdAt: time.Now(),
		ttl:       DefaultCacheTTL,
	})
	return key, nil
}

// LoadDilithiumPrivateKeyCached 带缓存的 Dilithium 私钥加载（支持TTL和大小限制）.
func LoadDilithiumPrivateKeyCached(path string) (interface{}, error) {
	cacheKey := "dilithium_priv:" + path
	if cached, ok := keyCache.Load(cacheKey); ok {
		entry := cached.(*keyCacheEntry)
		if time.Since(entry.createdAt) < entry.ttl {
			return entry.key, nil
		}
		keyCache.Delete(cacheKey)
	}

	checkCacheSize()

	key, err := LoadDilithiumPrivateKey(path)
	if err != nil {
		return nil, err
	}

	keyCache.Store(cacheKey, &keyCacheEntry{
		key:       key,
		createdAt: time.Now(),
		ttl:       DefaultCacheTTL,
	})
	return key, nil
}

// ClearKeyCache 清空密钥缓存
// 用于测试或手动清理缓存.
func ClearKeyCache() {
	keyCache.Range(func(key, _ interface{}) bool {
		keyCache.Delete(key)
		return true
	})
}

// GetCacheSize 获取当前缓存的密钥数量.
func GetCacheSize() int {
	count := 0
	keyCache.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// GetCacheInfo 获取缓存详细信息
// 返回：总条目数、已过期条目数、总大小（字节估算）.
func GetCacheInfo() (total int, expired int, estimatedSize int) {
	now := time.Now()
	keyCache.Range(func(_, value interface{}) bool {
		total++
		entry := value.(*keyCacheEntry)
		if now.Sub(entry.createdAt) >= entry.ttl {
			expired++
		}
		// 估算大小：每个条目约 100 字节 + 密钥大小
		estimatedSize += 100
		return true
	})
	return total, expired, estimatedSize
}

// SaveDilithiumKeys 保存 Dilithium3 密钥对到文件.
func SaveDilithiumKeys(pub interface{}, priv interface{}, pubPath, privPath string) error {
	keyPair, err := ExportDilithiumKeys(pub, priv)
	if err != nil {
		return err
	}

	// 保存公钥
	if err := os.WriteFile(pubPath, keyPair.Public, pubKeyFilePerm); err != nil {
		return fmt.Errorf("save Dilithium public key: %w", err)
	}

	// 保存私钥
	if err := os.WriteFile(privPath, keyPair.Private, keyFilePerm); err != nil {
		return fmt.Errorf("save Dilithium private key: %w", err)
	}

	// 在 Unix/Linux 系统上，确保私钥权限正确
	if runtime.GOOS != windowsOS {
		if err := os.Chmod(privPath, keyFilePerm); err != nil {
			return fmt.Errorf("set Dilithium private key permissions: %w", err)
		}
	}

	return nil
}
