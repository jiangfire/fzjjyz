# P2 - Revive 问题修复清单

**优先级：🟢 中**
**数量：57个**
**风险：代码规范、可维护性**
**状态：待修复**

---

## 📋 问题分类

### 1. unused-parameter - 未使用参数 (15个)

**问题：** 函数参数声明但未使用

#### 1.1 cmd/fzjjyz/decrypt.go (1个)

**第46行：**
```go
func runDecrypt(cmd *cobra.Command, args []string) error {
    // cmd 未使用
}
```
**修复方案：**
```go
func runDecrypt(_ *cobra.Command, args []string) error {
    // 使用 _ 忽略参数
}
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.2 cmd/fzjjyz/decrypt_dir.go (1个)

**第47行：**
```go
func runDecryptDir(cmd *cobra.Command, args []string) error {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.3 cmd/fzjjyz/encrypt.go (1个)

**第46行：**
```go
func runEncrypt(cmd *cobra.Command, args []string) error {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.4 cmd/fzjjyz/encrypt_dir.go (1个)

**第48行：**
```go
func runEncryptDir(cmd *cobra.Command, args []string) error {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.5 cmd/fzjjyz/info.go (1个)

**第31行：**
```go
func runInfo(cmd *cobra.Command, args []string) error {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.6 cmd/fzjjyz/keygen.go (1个)

**第35行：**
```go
func runKeygen(cmd *cobra.Command, args []string) error {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.7 cmd/fzjjyz/keymanage.go (1个)

**第41行：**
```go
func runKeymanage(cmd *cobra.Command, args []string) error {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.8 cmd/fzjjyz/version.go (1个)

**第15行：**
```go
Run: func(cmd *cobra.Command, args []string) {
```

**修复方案：**
```go
Run: func(_ *cobra.Command, _ []string) {
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.9 internal/crypto/keyfile.go (3个)

**第340行：**
```go
keyCache.Range(func(key, value interface{}) bool {
    // key 未使用
})
```

**修复方案：**
```go
keyCache.Range(func(_, value interface{}) bool {
})
```

**第479行：**
```go
keyCache.Range(func(key, value interface{}) bool {
    // value 未使用
})
```

**修复方案：**
```go
keyCache.Range(func(key, _ interface{}) bool {
})
```

**第488行：**
```go
keyCache.Range(func(key, value interface{}) bool {
    // key 未使用
})
```

**修复方案：**
```go
keyCache.Range(func(_, value interface{}) bool {
})
```

**状态：** ⬜ 待修复
**预计时间：** 3分钟

---

#### 1.10 internal/crypto/stream_encrypt.go (1个)

**第86行：**
```go
func (se *StreamingEncryptor) encryptData(input io.Reader, output io.Writer, sharedSecret []byte, nonce []byte) error {
    // input 未使用
}
```

**修复方案：**
```go
func (se *StreamingEncryptor) encryptData(_ io.Reader, output io.Writer, sharedSecret []byte, nonce []byte) error {
}
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.11 internal/i18n/i18n.go (1个)

**第153行：**
```go
func (e *emptyDict) Get(key string) string {
    // key 未使用
}
```

**修复方案：**
```go
func (e *emptyDict) Get(_ string) string {
    return ""
}
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 1.12 internal/i18n/i18n_test.go (1个)

**第156行：**
```go
func TestConcurrentAccess(t *testing.T) {
    // t 未使用
}
```

**修复方案：**
```go
func TestConcurrentAccess(_ *testing.T) {
}
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

### 2. package-comments - 缺少包注释 (6个)

**问题：** 包缺少文档注释

#### 2.1 cmd/fzjjyz/decrypt.go

**第1行：**
```go
package main
```
**修复方案：**
```go
// Package main implements the decrypt command for fzjjyz.
package main
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.2 internal/crypto/archive.go

**第1行：**
```go
package crypto
```
**修复方案：**
```go
// Package crypto provides post-quantum encryption functionality.
package crypto
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.3 internal/crypto/keyfile.go

**第1行：**
```go
package crypto
```
**修复方案：** 已有包注释，跳过

**状态：** ⬜ 已修复

---

#### 2.4 internal/crypto/stream_decrypt.go

**第1行：**
```go
package crypto
```
**修复方案：**
```go
// Package crypto provides post-quantum encryption functionality.
package crypto
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.5 internal/format/header.go

**第1行：**
```go
package format
```
**修复方案：**
```go
// Package format handles file header serialization and parsing.
package format
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.6 internal/format/parser.go

**第1行：**
```go
package format
```
**修复方案：** 已有包注释，跳过

**状态：** ⬜ 已修复

---

#### 2.7 internal/i18n/cobra.go

**第1行：**
```go
package i18n
```
**修复方案：**
```go
// Package i18n provides internationalization support.
package i18n
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.8 internal/utils/errors.go

**第1行：**
```go
package utils
```
**修复方案：**
```go
// Package utils provides error handling utilities.
package utils
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.9 internal/utils/logging.go

**第1行：**
```go
package utils
```
**修复方案：**
```go
// Package utils provides logging utilities.
package utils
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

### 3. exported - 缺少导出项注释 (20个)

**问题：** 导出的函数、类型、常量缺少文档

#### 3.1 internal/crypto/keyfile.go (2个)

**第66行：**
```go
// 保存密钥文件（遵循安全原则）
func SaveKeyFiles(...) error {
```
**修复方案：**
```go
// SaveKeyFiles saves key files following security principles.
func SaveKeyFiles(...) error {
```

**第116行：**
```go
// 加载密钥文件
func LoadKeyFiles(...) {
```
**修复方案：**
```go
// LoadKeyFiles loads key files from specified paths.
func LoadKeyFiles(...) {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 3.2 internal/crypto/keygen.go (8个)

**第15行：**
```go
// 密钥对结构（表达原则：数据结构优先）
type HybridPublicKey struct {
```
**修复方案：**
```go
// HybridPublicKey represents a hybrid public key (Kyber + ECDH).
type HybridPublicKey struct {
```

**第21行：**
```go
type HybridPrivateKey struct {
```
**修复方案：**
```go
// HybridPrivateKey represents a hybrid private key (Kyber + ECDH).
type HybridPrivateKey struct {
```

**第26行：**
```go
// 生成Kyber密钥对
func GenerateKyberKeys() (kem.PublicKey, kem.PrivateKey, error) {
```
**修复方案：**
```go
// GenerateKyberKeys generates a new Kyber key pair.
func GenerateKyberKeys() (kem.PublicKey, kem.PrivateKey, error) {
```

**第39行：**
```go
// 生成ECDH密钥对
func GenerateECDHKeys() (*ecdh.PublicKey, *ecdh.PrivateKey, error) {
```
**修复方案：**
```go
// GenerateECDHKeys generates a new ECDH key pair.
func GenerateECDHKeys() (*ecdh.PublicKey, *ecdh.PrivateKey, error) {
```

**第51行：**
```go
// 导出公钥到PEM格式
func ExportPublicKey(pub interface{}) (string, error) {
```
**修复方案：**
```go
// ExportPublicKey exports a public key to PEM format.
func ExportPublicKey(pub interface{}) (string, error) {
```

**第80行：**
```go
// 导出私钥到PEM格式（注意权限设置）
func ExportPrivateKey(priv interface{}) (string, error) {
```
**修复方案：**
```go
// ExportPrivateKey exports a private key to PEM format.
func ExportPrivateKey(priv interface{}) (string, error) {
```

**第107行：**
```go
// 从PEM导入密钥
func ImportKeys(pubPEM, privPEM string) (interface{}, interface{}, error) {
```
**修复方案：**
```go
// ImportKeys imports keys from PEM format.
func ImportKeys(pubPEM, privPEM string) (interface{}, interface{}, error) {
```

**状态：** ⬜ 待修复
**预计时间：** 8分钟

---

#### 3.3 internal/crypto/operations_shared.go (1个)

**第166行：**
```go
func decapsulateKeys(kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, encapsulated []byte, ecdhPub []byte) ([]byte, error) {
```
**修复方案：**
```go
// decapsulateKeys decapsulates shared secret from encrypted keys.
func decapsulateKeys(kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, encapsulated []byte, ecdhPub []byte) ([]byte, error) {
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 3.4 internal/crypto/stream_utils.go (1个)

**第1行：**
```go
// 流式处理工具函数
```
**修复方案：**
```go
// Package crypto provides post-quantum encryption functionality.
package crypto
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 3.5 internal/format/header.go (2个)

**第267行：**
```go
func NewFileHeader(filename string, fileSize uint64, kyberEnc []byte, ecdhPub [32]byte, iv [12]byte, signature []byte, hash [32]byte) *FileHeader {
```
**修复方案：**
```go
// NewFileHeader creates a new file header with the specified parameters.
func NewFileHeader(filename string, fileSize uint64, kyberEnc []byte, ecdhPub [32]byte, iv [12]byte, signature []byte, hash [32]byte) *FileHeader {
```

**第312行：**
```go
func GetHeaderInfo(header *FileHeader) *HeaderInfo {
```
**修复方案：**
```go
// GetHeaderInfo extracts basic information from a file header.
func GetHeaderInfo(header *FileHeader) *HeaderInfo {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 3.6 internal/format/parser.go (1个)

**第300行：**
```go
// GetHeaderInfo 从文件头提取基本信息（用于快速预览）
type HeaderInfo struct {
```
**修复方案：**
```go
// HeaderInfo contains basic file header information.
type HeaderInfo struct {
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 3.7 internal/utils/errors.go (5个)

**第6行：**
```go
// 错误代码枚举（表达原则：数据结构优先）
type ErrorCode int
```
**修复方案：**
```go
// ErrorCode represents different types of errors.
type ErrorCode int
```

**第41行：**
```go
// 自定义错误结构（透明原则：清晰状态）
type CryptoError struct {
```
**修复方案：**
```go
// CryptoError represents a cryptographic error with context.
type CryptoError struct {
```

**第51行：**
```go
// 错误上下文（模块原则：可组合）
type ErrorContext struct {
```
**修复方案：**
```go
// ErrorContext provides error context and wrapping.
type ErrorContext struct {
```

**第66行：**
```go
// 工厂函数（修复原则：及早抛出明确异常）
func NewCryptoError(code ErrorCode, msg string) error {
```
**修复方案：**
```go
// NewCryptoError creates a new CryptoError.
func NewCryptoError(code ErrorCode, msg string) error {
```

**第71行：**
```go
// 错误分类函数
func IsFormatError(err error) bool {
```
**修复方案：**
```go
// IsFormatError checks if error is a format error.
func IsFormatError(err error) bool {
```

**状态：** ⬜ 待修复
**预计时间：** 5分钟

---

#### 3.8 internal/utils/logging.go (2个)

**第9行：**
```go
// 日志器（安静原则：无用信息保持安静）
type Logger struct {
```
**修复方案：**
```go
// Logger provides thread-safe logging with verbosity control.
type Logger struct {
```

**第17行：**
```go
func NewLogger(w io.Writer, silent, verbose bool) *Logger {
```
**修复方案：**
```go
// NewLogger creates a new logger with specified settings.
func NewLogger(w io.Writer, silent, verbose bool) *Logger {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

### 4. empty-block - 空代码块 (2个)

**问题：** 空的 if 块或函数体

#### 4.1 cmd/fzjjyz/decrypt_dir.go (1个)

**第161行：**
```go
if removeErr := os.Remove(tempZipPath); removeErr != nil {
    // 忽略清理错误，不影响主流程
}
```
**修复方案：**
```go
if removeErr := os.Remove(tempZipPath); removeErr != nil {
    // 忽略清理错误，不影响主流程
    return fmt.Errorf("cleanup failed: %w", removeErr)
}
```
或
```go
_ = os.Remove(tempZipPath) // 忽略清理错误
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 4.2 cmd/fzjjyz/encrypt_dir.go (1个)

**第169行：**
```go
if removeErr := os.Remove(tempZipPath); removeErr != nil {
    // 忽略清理错误
}
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

### 5. var-naming - 包名问题 (2个)

**问题：** 包名与标准库冲突或无意义

#### 5.1 cmd/fzjjyz/utils/errors.go

**第1行：**
```go
package utils
```
**修复方案：** 保持现状，包名可接受

**状态：** ⬜ 已忽略

---

#### 5.2 internal/utils/errors.go

**第1行：**
```go
package utils
```
**修复方案：** 保持现状

**状态：** ⬜ 已忽略

---

### 6. redefines-builtin-id - 重定义内置函数 (1个)

**问题：** 重定义了 Go 内置函数

#### 6.1 cmd/fzjjyz/main_test.go

**第418行：**
```go
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```
**修复方案：**
```go
func minInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

### 7. increment-decrement - 递增风格 (2个)

**问题：** 应使用 `++` 而非 `+= 1`

#### 7.1 internal/format/header.go (2个)

**第294行：**
```go
size += 1 // ECDHLen
```
**修复方案：**
```go
size++ // ECDHLen
```

**第296行：**
```go
size += 1 // IVLen
```
**修复方案：**
```go
size++ // IVLen
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

## 📊 统计信息

| 问题类型 | 数量 | 修复难度 | 预计时间 |
|---------|------|---------|---------|
| unused-parameter | 15 | 简单 | 15分钟 |
| package-comments | 6 | 简单 | 6分钟 |
| exported | 20 | 中等 | 20分钟 |
| empty-block | 2 | 简单 | 2分钟 |
| redefines-builtin-id | 1 | 简单 | 1分钟 |
| increment-decrement | 2 | 简单 | 1分钟 |
| var-naming | 2 | 忽略 | 0分钟 |
| **总计** | **48个** | - | **45分钟** |

---

## 🔧 修复模板

### 模板1：忽略未使用参数
```go
func myFunc(_ *cobra.Command, args []string) error {
    // 只使用 args
}
```

### 模板2：包注释
```go
// Package <name> provides <description>.
package <name>
```

### 模板3：导出项注释
```go
// <Name> <description in sentence case>.
func <Name>() {
```

### 模板4：空代码块
```go
// 添加实际逻辑或注释
if err := cleanup(); err != nil {
    log.Printf("cleanup warning: %v", err)
}
```

### 模板5：避免内置函数名
```go
// 原：min
// 改：minInt, minVal, minimum
func minInt(a, b int) int {
```

### 模板6：递增风格
```go
// 原：size += 1
// 改：size++
size++
```

---

## ✅ 验证标准

修复后运行：
```bash
golangci-lint run --disable-all --enable=revive
```

应输出：`0 issues` 或仅可忽略的问题

---

**创建时间：** 2025-12-30
**预计完成：** 2025-12-31
**负责人：** 待分配
