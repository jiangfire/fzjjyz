# P3 - 低优先级问题修复清单

**优先级：🔵 低**
**数量：37个**
**风险：代码可读性、可维护性**
**状态：可选修复**

---

## 📋 问题分类

### 1. funlen - 函数过长 (10个)

**问题：** 函数超过 60 行（测试函数超过 40 行）

#### 1.1 cmd/fzjjyz/main_test.go

**第318行：**
```go
func TestCLIBenchmark(t *testing.T) {
    // 86 行
}
```
**建议：** 保持现状，测试函数复杂度可接受

**状态：** ⬜ 建议保持

---

#### 1.2 internal/crypto/archive.go

**第133行：**
```go
func ExtractZipToDirectory(zipData []byte, targetDir string) error {
    // 71 行
}
```
**建议：** 可拆分为辅助函数，但当前逻辑清晰

**状态：** ⬜ 建议保持

---

#### 1.3 internal/crypto/integration_test.go

**第15行：**
```go
func TestIntegrationEndToEnd(t *testing.T) {
    // 62 行
}
```
**建议：** 保持现状，集成测试需要完整流程

**状态：** ⬜ 建议保持

---

#### 1.4 internal/crypto/keygen.go

**第226行：**
```go
func GenerateKeyPairParallel() (
    // 66 行
)
```
**建议：** 可拆分，但并行生成逻辑需要整体控制

**状态：** ⬜ 建议保持

---

#### 1.5 internal/crypto/stream_test.go

**第171行：**
```go
func TestHeaderOptimized(t *testing.T) {
    // 76 行
}
```
**建议：** 保持现状

**状态：** ⬜ 建议保持

---

#### 1.6 internal/format/header.go

**第34行：**
```go
func (h *FileHeader) MarshalBinary() ([]byte, error) {
    // 41 行
}
```
**建议：** 已有优化版本 `MarshalBinaryOptimized`

**状态：** ⬜ 建议保持

---

**第161行：**
```go
func (h *FileHeader) UnmarshalBinary(data []byte) error {
    // 53 行
}
```
**建议：** 逻辑清晰，保持现状

**状态：** ⬜ 建议保持

---

#### 1.7 internal/format/header_test.go

**第10行：**
```go
func TestFileHeaderSerialization(t *testing.T) {
    // 67 行
}
```
**建议：** 保持现状

**状态：** ⬜ 建议保持

---

#### 1.8 internal/format/parser.go

**第13行：**
```go
func ParseFileHeader(r io.Reader) (*FileHeader, error) {
    // 53 行
}
```
**建议：** 逻辑清晰，保持现状

**状态：** ⬜ 建议保持

---

#### 1.9 internal/format/parser_test.go

**第13行：**
```go
func TestParseFileHeader(t *testing.T) {
    // 46 行
}
```
**建议：** 保持现状

**状态：** ⬜ 建议保持

---

### 2. lll - 行过长 (11个)

**问题：** 代码行超过 120 个字符

#### 2.1 cmd/fzjjyz/keymanage.go

**第125行：**
```go
if err := crypto.SaveKeyFiles(hybridPub.Kyber, hybridPub.ECDH, hybridPriv.Kyber, hybridPriv.ECDH, newPubPath, newPrivPath); err != nil {
```
**修复方案：**
```go
if err := crypto.SaveKeyFiles(
    hybridPub.Kyber,
    hybridPub.ECDH,
    hybridPriv.Kyber,
    hybridPriv.ECDH,
    newPubPath,
    newPrivPath,
); err != nil {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 2.2 internal/crypto/integration_test.go (3个)

**第219行：**
```go
if err := DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv)); err != nil {
```

**第258行：**
```go
if err := DecryptFile(encryptedFile, decryptedFile, kyberPriv, ecdhPriv, DilithiumGetPublicKey(dilithiumPriv)); err != nil {
```

**第362行：**
```go
if err := DecryptFile(encryptedFile, decryptedFile, hybridPriv.Kyber, hybridPriv.ECDH, DilithiumGetPublicKey(dilithiumPriv)); err != nil {
```

**修复方案：** 拆分参数到多行

**状态：** ⬜ 待修复
**预计时间：** 3分钟

---

#### 2.3 internal/crypto/operations.go (2个)

**第22行：**
```go
func EncryptFile(inputPath, outputPath string, kyberPub kem.PublicKey, ecdhPub *ecdh.PublicKey, dilithiumPriv interface{}) error {
```

**第54行：**
```go
func DecryptFile(inputPath, outputPath string, kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, dilithiumPub interface{}) error {
```

**修复方案：**
```go
func EncryptFile(
    inputPath, outputPath string,
    kyberPub kem.PublicKey,
    ecdhPub *ecdh.PublicKey,
    dilithiumPriv interface{},
) error {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 2.4 internal/crypto/operations_shared.go (1个)

**第166行：**
```go
func decapsulateKeys(kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, encapsulated []byte, ecdhPub []byte) ([]byte, error) {
```

**修复方案：** 同上

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 2.5 internal/format/header.go (2个)

**第267行：**
```go
func NewFileHeader(filename string, fileSize uint64, kyberEnc []byte, ecdhPub [32]byte, iv [12]byte, signature []byte, hash [32]byte) *FileHeader {
```

**修复方案：**
```go
func NewFileHeader(
    filename string,
    fileSize uint64,
    kyberEnc []byte,
    ecdhPub [32]byte,
    iv [12]byte,
    signature []byte,
    hash [32]byte,
) *FileHeader {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 2.6 internal/format/header_test.go (2个)

**第525行：**
```go
header.ECDHPub = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
```

**第528行：**
```go
header.SHA256Hash = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
```

**修复方案：**
```go
header.ECDHPub = [32]byte{
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
    17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
}
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 2.7 internal/i18n/en_US.go (1个)

**第246行：**
```go
"status.warning_no_sign_verify": "⚠️  Warning: No signature verification key provided, skipping signature verification",
```

**修复方案：**
```go
"status.warning_no_sign_verify": "⚠️  Warning: No signature verification key " +
    "provided, skipping signature verification",
```

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

### 3. gocognit - 认知复杂度高 (4个)

**问题：** 认知复杂度超过 30

#### 3.1 cmd/fzjjyz/decrypt_dir.go

**第47行：**
```go
func runDecryptDir(cmd *cobra.Command, args []string) error {
    // 认知复杂度: 32
}
```
**建议：** 保持现状，逻辑虽然复杂但清晰

**状态：** ⬜ 建议保持

---

#### 3.2 cmd/fzjjyz/main_test.go

**第14行：**
```go
func TestCLIIntegration(t *testing.T) {
    // 认知复杂度: 54
}
```
**建议：** 保持现状，集成测试需要完整覆盖

**状态：** ⬜ 建议保持

---

#### 3.3 internal/crypto/archive.go

**第31行：**
```go
func CreateZipFromDirectory(sourceDir string, output io.Writer, opts ArchiveOptions) error {
    // 认知复杂度: 36
}
```
**建议：** 保持现状，归档逻辑复杂但必要

**状态：** ⬜ 建议保持

---

#### 3.4 internal/crypto/stream_test.go

**第257行：**
```go
func TestStreamingEncryption(t *testing.T) {
    // 认知复杂度: 31
}
```
**建议：** 保持现状

**状态：** ⬜ 建议保持

---

### 4. goconst - 字符串常量重复 (3个)

**问题：** 相同字符串多次出现，应定义为常量

#### 4.1 internal/crypto/keyfile.go

**第104行：**
```go
if runtime.GOOS != "windows" {
    // "windows" 出现 4 次
}
```
**修复方案：**
```go
const osWindows = "windows"

if runtime.GOOS != osWindows {
```

**状态：** ⬜ 待修复
**预计时间：** 2分钟

---

#### 4.2 internal/i18n/i18n_test.go (2个)

**第21行：**
```go
if GetLanguage() != "zh_CN" {
    // "zh_CN" 出现 4 次
}
```
**修复方案：**
```go
const testLang = "zh_CN"

if GetLanguage() != testLang {
```

**第76行：**
```go
if result != "nonexistent.key" {
    // "nonexistent.key" 出现 3 次
}
```
**修复方案：**
```go
const nonexistentKey = "nonexistent.key"

if result != nonexistentKey {
```

**状态：** ⬜ 待修复
**预计时间：** 3分钟

---

### 5. gochecknoinits - init函数 (2个)

**问题：** 不应使用 init 函数

#### 5.1 cmd/fzjjyz/main.go

**第26行：**
```go
func init() {
    i18n.Init("")
}
```
**建议：** 保持现状，init 用于国际化初始化是可接受的模式

**状态：** ⬜ 建议保持

---

#### 5.2 internal/crypto/keyfile.go

**第41行：**
```go
func init() {
    // 缓存初始化
}
```
**建议：** 保持现状

**状态：** ⬜ 建议保持

---

### 6. unused - 未使用代码 (2个)

**问题：** 定义但未使用的变量/函数

#### 6.1 internal/crypto/keyfile.go

**第38行：**
```go
var cacheCleanupTimer *time.Timer
```
**修复方案：** 直接删除

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

#### 6.2 internal/i18n/i18n.go

**第130行：**
```go
func formatString(format string, args ...interface{}) string {
    return fmt.Sprintf(format, args...)
}
```
**修复方案：** 直接删除

**状态：** ⬜ 待修复
**预计时间：** 1分钟

---

## 📊 统计信息

| 问题类型 | 数量 | 建议操作 | 预计时间 |
|---------|------|---------|---------|
| funlen | 10 | 保持 | 0分钟 |
| lll | 11 | 修复 | 13分钟 |
| gocognit | 4 | 保持 | 0分钟 |
| goconst | 3 | 修复 | 5分钟 |
| gochecknoinits | 2 | 保持 | 0分钟 |
| unused | 2 | 删除 | 2分钟 |
| **总计** | **32个** | - | **20分钟** |

---

## 🔧 修复模板

### 模板1：长函数签名
```go
// 原代码
func myFunc(a int, b string, c []byte, d interface{}, e error, f bool) error {

// 修复后
func myFunc(
    a int,
    b string,
    c []byte,
    d interface{},
    e error,
    f bool,
) error {
```

### 模板2：长数组字面量
```go
// 原代码
arr := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

// 修复后
arr := [32]byte{
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
    17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
}
```

### 模板3：字符串常量
```go
// 原代码
if runtime.GOOS != "windows" {
    if runtime.GOOS != "windows" {
    }
}

// 修复后
const osWindows = "windows"

if runtime.GOOS != osWindows {
    if runtime.GOOS != osWindows {
    }
}
```

### 模板4：删除未使用代码
```go
// 删除
var cacheCleanupTimer *time.Timer

// 删除
func unusedFunc() {
}
```

---

## ✅ 验证标准

修复后运行：
```bash
golangci-lint run --disable-all --enable=lll,goconst,unused
```

应输出：`0 issues` 或仅 funlen/gocognit/gochecknoinits（可忽略）

---

**创建时间：** 2025-12-30
**预计完成：** 2025-12-30（可选）
**负责人：** 待分配
**优先级：** 低（可选修复）
