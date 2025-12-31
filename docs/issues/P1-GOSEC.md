# P1 - Gosec 问题修复清单

**优先级：🟡 高**
**数量：100个**
**风险：安全漏洞**
**状态：待修复**

---

## 📋 问题分类

### 1. G304 - 潜在文件包含漏洞 (58个)

**风险：** 攻击者可能通过恶意路径访问任意文件

#### 1.1 cmd/fzjjyz/decrypt.go (1个)

**第54行：**
```go
headerFile, err := os.Open(decryptInput)
```
**修复方案：**
```go
// 验证路径安全性
if !filepath.IsAbs(decryptInput) {
    return fmt.Errorf("path must be absolute")
}
if strings.Contains(decryptInput, "..") {
    return fmt.Errorf("path traversal detected")
}
headerFile, err := os.Open(decryptInput)
```

**状态：** ⬜ 待修复

---

#### 1.2 cmd/fzjjyz/decrypt_dir.go (2个)

**第69行：**
```go
headerFile, err := os.Open(decryptDirInput)
```

**第167行：**
```go
zipData, err := os.ReadFile(tempZipPath)
```

**修复方案：** 同上，添加路径验证

**状态：** ⬜ 待修复

---

#### 1.3 cmd/fzjjyz/info.go (1个)

**第38行：**
```go
data, err := os.ReadFile(infoInput)
```

**状态：** ⬜ 待修复

---

#### 1.4 internal/crypto/archive.go (2个)

**第115行：**
```go
file, err := os.Open(path)
```

**第199行：**
```go
dstFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, file.Mode())
```

**状态：** ⬜ 待修复

---

#### 1.5 internal/crypto/archive_test.go (2个)

**第110行：**
```go
content, err := os.ReadFile(extractedFile)
```

**第297行：**
```go
content, err := os.ReadFile(filepath.Join(extractDir, tt.path))
```

**状态：** ⬜ 待修复

---

#### 1.6 internal/crypto/benchmark_test.go (2个)

**第264行：**
```go
originalData, _ := os.ReadFile(inputPath)
```

**第265行：**
```go
decryptedData, _ := os.ReadFile(decryptedPath)
```

**状态：** ⬜ 待修复

---

#### 1.7 internal/crypto/hash_utils.go (1个)

**第14行：**
```go
file, err := os.Open(path)
```

**状态：** ⬜ 待修复

---

#### 1.8 internal/crypto/integration_test.go (9个)

**第68行：**
```go
f, err := os.Open(encryptedFile)
```

**第106行：**
```go
decryptedData, err := os.ReadFile(decryptedFile)
```

**第142行：**
```go
encryptedData, _ := os.ReadFile(encryptedFile)
```

**第223行：**
```go
decryptedData, _ := os.ReadFile(decryptedFile)
```

**第262行：**
```go
decryptedData, _ := os.ReadFile(decryptedFile)
```

**第310行：**
```go
decData, _ := os.ReadFile(decPath)
```

**第366行：**
```go
origData, _ := os.ReadFile(testFile)
```

**第367行：**
```go
decData, _ := os.ReadFile(decryptedFile)
```

**第411行：**
```go
decData, _ := os.ReadFile(decPath)
```

**状态：** ⬜ 待修复

---

#### 1.9 internal/crypto/keyfile.go (6个)

**第118行：**
```go
pubPEM, err := os.ReadFile(pubPath)
```

**第126行：**
```go
privPEM, err := os.ReadFile(privPath)
```

**第139行：**
```go
pubPEM, err := os.ReadFile(pubPath)
```

**第158行：**
```go
privPEM, err := os.ReadFile(privPath)
```

**第260行：**
```go
pubPEM, err := os.ReadFile(pubPath)
```

**第268行：**
```go
privPEM, err := os.ReadFile(privPath)
```

**状态：** ⬜ 待修复

---

#### 1.10 internal/crypto/keyfile_test.go (2个)

**第132行：**
```go
pubContent, _ := os.ReadFile(pubPath)
```

**第133行：**
```go
privContent, _ := os.ReadFile(privPath)
```

**状态：** ⬜ 待修复

---

#### 1.11 internal/crypto/operations_shared.go (2个)

**第127行：**
```go
encryptedData, err := os.ReadFile(inputPath)
```

**第222行：**
```go
plaintext, err := os.ReadFile(inputPath)
```

**状态：** ⬜ 待修复

---

#### 1.12 internal/crypto/operations_test.go (10个)

**第58行：**
```go
decryptedData, err := os.ReadFile(decryptedFile)
```

**第101行：**
```go
decryptedData, err := os.ReadFile(decryptedFile)
```

**第147行：**
```go
decryptedData, err := os.ReadFile(decryptedFile)
```

**第185行：**
```go
f, err := os.Open(encryptedFile)
```

**第249行：**
```go
encryptedData, err := os.ReadFile(encryptedFile)
```

**第351行：**
```go
decryptedData, err := os.ReadFile(decryptedFile)
```

**第394行：**
```go
decryptedData, err := os.ReadFile(decryptedFile)
```

**第452行：**
```go
decryptedData, _ := os.ReadFile(decryptedFile)
```

**第504行：**
```go
encData, _ := os.ReadFile(encryptedFile)
```

**状态：** ⬜ 待修复

---

#### 1.13 internal/crypto/signature.go (1个)

**第104行：**
```go
data, err := os.ReadFile(filePath)
```

**状态：** ⬜ 待修复

---

#### 1.14 internal/crypto/stream_test.go (3个)

**第300行：**
```go
encData, _ := os.ReadFile(encryptedFile)
```

**第855行：**
```go
encData, _ := os.ReadFile(encryptedFile)
```

**状态：** ⬜ 待修复

---

### 2. G306 - 文件权限过松 (25个)

**风险：** 敏感数据可能被其他用户读取

#### 2.1 cmd/fzjjyz/encrypt_dir.go (1个)

**第109行：**
```go
if err := os.WriteFile(tempZipPath, zipData, 0644); err != nil {
```
**修复方案：**
```go
if err := os.WriteFile(tempZipPath, zipData, 0600); err != nil {
```

**状态：** ⬜ 待修复

---

#### 2.2 cmd/fzjjyz/keymanage.go (1个)

**第82行：**
```go
if err := os.WriteFile(keymanageOutput, pubPEM, 0644); err != nil {
```
**修复方案：**
```go
if err := os.WriteFile(keymanageOutput, pubPEM, 0600); err != nil {
```

**状态：** ⬜ 待修复

---

#### 2.3 cmd/fzjjyz/main_test.go (3个)

**第64行：**
```go
if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
```

**第348行：**
```go
if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
```

**修复方案：** 使用 0600

**状态：** ⬜ 待修复

---

#### 2.4 internal/crypto/archive.go (4个)

**第135行：**
```go
if err := os.MkdirAll(targetDir, 0755); err != nil {
```

**第176行：**
```go
if err := os.MkdirAll(targetPath, 0755); err != nil {
```

**第183行：**
```go
if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
```

**修复方案：** 使用 0750（生产代码），测试代码可保持 0755

**状态：** ⬜ 待修复

---

#### 2.5 internal/crypto/integration_test.go (1个)

**第44行：**
```go
if err := os.WriteFile(originalFile, testData, 0644); err != nil {
```

**状态：** ⬜ 待修复

---

#### 2.6 internal/crypto/keyfile.go (2个)

**第87行：**
```go
if err := os.WriteFile(pubPath, pubPEM, 0644); err != nil {
```

**第520行：**
```go
if err := os.WriteFile(pubPath, keyPair.Public, 0644); err != nil {
```

**修复方案：** 使用 0600

**状态：** ⬜ 待修复

---

#### 2.7 internal/crypto/operations_shared.go (2个)

**第115行：**
```go
if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
```

**第173行：**
```go
if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
```

**修复方案：** 使用 0600

**状态：** ⬜ 待修复

---

#### 2.8 internal/crypto/operations_test.go (7个)

**第34行：**
```go
if err := os.WriteFile(originalFile, originalData, 0644); err != nil {
```

**第85行：**
```go
if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
```

**第131行：**
```go
if err := os.WriteFile(largeFile, largeData, 0644); err != nil {
```

**第174行：**
```go
if err := os.WriteFile(testFile, testData, 0644); err != nil {
```

**第238行：**
```go
if err := os.WriteFile(testFile, []byte("Test data"), 0644); err != nil {
```

**第261行：**
```go
if err := os.WriteFile(tamperedFile, tamperedData, 0644); err != nil {
```

**第295行：**
```go
if err := os.WriteFile(testFile, []byte("Secret data"), 0644); err != nil {
```

**状态：** ⬜ 待修复

---

#### 2.9 internal/crypto/operations_test.go (续)

**第335行：**
```go
if err := os.WriteFile(binaryFile, binaryData, 0644); err != nil {
```

**第378行：**
```go
if err := os.WriteFile(specialFile, []byte("Special chars test"), 0644); err != nil {
```

**第426行：**
```go
if err := os.WriteFile(file, data, 0644); err != nil {
```

**第476行：**
```go
if err := os.WriteFile(testFile, testData, 0644); err != nil {
```

**状态：** ⬜ 待修复

---

#### 2.10 internal/crypto/signature_test.go (3个)

**第169行：**
```go
if err := os.WriteFile(testFile, testData, 0644); err != nil {
```

**第199行：**
```go
if err := os.WriteFile(testFile, []byte("Original data"), 0644); err != nil {
```

**第206行：**
```go
if err := os.WriteFile(testFile, []byte("Tampered data"), 0644); err != nil {
```

**状态：** ⬜ 待修复

---

#### 2.11 internal/crypto/stream_test.go (2个)

**第66行：**
```go
if err := os.WriteFile(testFile, testData, 0644); err != nil {
```

**第282行：**
```go
if err := os.WriteFile(originalFile, testData, 0644); err != nil {
```

**第344行：**
```go
if err := os.WriteFile(originalFile, testData, 0644); err != nil {
```

**状态：** ⬜ 待修复

---

### 3. G301 - 目录权限过松 (7个)

**风险：** 目录可能被其他用户访问

#### 3.1 cmd/fzjjyz/keygen.go (1个)

**第42行：**
```go
if err := os.MkdirAll(keygenOutputDir, 0755); err != nil {
```
**修复方案：**
```go
if err := os.MkdirAll(keygenOutputDir, 0750); err != nil {
```

**状态：** ⬜ 待修复

---

#### 3.2 cmd/fzjjyz/keymanage.go (1个)

**第103行：**
```go
if err := os.MkdirAll(keymanageOutputDir, 0755); err != nil {
```
**修复方案：**
```go
if err := os.MkdirAll(keymanageOutputDir, 0750); err != nil {
```

**状态：** ⬜ 待修复

---

#### 3.3 cmd/fzjjyz/main_test.go (1个)

**第168行：**
```go
if err := os.MkdirAll(importDir, 0755); err != nil {
```

**状态：** ⬜ 待修复

---

#### 3.4 internal/crypto/archive.go (3个)

**第135行：**
```go
if err := os.MkdirAll(targetDir, 0755); err != nil {
```

**第176行：**
```go
if err := os.MkdirAll(targetPath, 0755); err != nil {
```

**第183行：**
```go
if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
```

**状态：** ⬜ 待修复

---

#### 3.5 internal/crypto/keyfile.go (1个)

**第42行：**
```go
if err := os.MkdirAll(keygenOutputDir, 0755); err != nil {
```

**状态：** ⬜ 待修复

---

### 4. G204 - 子进程变量注入 (10个)

**风险：** 命令注入（在测试环境中可接受）

#### 4.1 cmd/fzjjyz/main_test.go (10个)

**第43行：**
```go
cmd := exec.Command(executable, "keygen", "-d", testDir, "-n", keyPrefix)
```

**第72行：**
```go
cmd := exec.Command(executable, "encrypt", "-i", testFile, "-o", encryptedFile, ...)
```

**第95行：**
```go
cmd := exec.Command(executable, "decrypt", "-i", encryptedFile, "-o", decryptedFile, ...)
```

**第123行：**
```go
cmd := exec.Command(executable, "keymanage", "-a", "export", ...)
```

**第149行：**
```go
cmd := exec.Command(executable, "keymanage", "-a", "verify", ...)
```

**第172行：**
```go
cmd := exec.Command(executable, "keymanage", "-a", "import", ...)
```

**第197行：**
```go
cmd := exec.Command(executable, "version")
```

**第216行：**
```go
cmd := exec.Command(executable, "keygen", "-d", wrongKeyDir, "-n", wrongKeyPrefix)
```

**第225行：**
```go
cmd = exec.Command(executable, "decrypt", "-i", encryptedFile, ...)
```

**第270行：**
```go
cmd := exec.Command("go", "build", "-o", executable, "./cmd/fzjjyz")
```

**第301行：**
```go
cmd := exec.Command(executable, tt.command...)
```

**第335行：**
```go
cmd := exec.Command(executable, "keygen", "-d", testDir, "-n", keyPrefix)
```

**第357行：**
```go
cmd := exec.Command(executable, "encrypt", "-i", testFile, ...)
```

**第381行：**
```go
cmd := exec.Command(executable, "encrypt", "-i", testFile, ..., "--force")
```

**第394行：**
```go
cmd = exec.Command(executable, "decrypt", "-i", encryptedFile, ..., "--force")
```

**修复方案：** 在测试环境中，这些是可接受的。如果担心，可以添加注释：
```go
// #nosec G204 - 测试环境，变量来源可信
cmd := exec.Command(executable, ...)
```

**状态：** ⬜ 待修复（或标记为已忽略）

---

### 5. G110 - 解压缩炸弹 (1个)

**风险：** 大量解压导致 DoS

#### 5.1 internal/crypto/archive.go (1个)

**第210行：**
```go
if _, err := io.Copy(dstFile, srcFile); err != nil {
```

**修复方案：**
```go
// 添加大小限制
const maxExtractSize = 100 * 1024 * 1024 // 100MB
if file.UncompressedSize64 > maxExtractSize {
    return fmt.Errorf("file too large: %d bytes", file.UncompressedSize64)
}
if _, err := io.Copy(dstFile, srcFile); err != nil {
```

**状态：** ⬜ 待修复

---

### 6. G115 - 整数溢出 (1个)

**风险：** uint64 转换为 int64 可能溢出

#### 6.1 internal/crypto/archive.go (1个)

**第227行：**
```go
totalSize += int64(file.UncompressedSize64)
```

**修复方案：**
```go
if file.UncompressedSize64 > uint64(math.MaxInt64) {
    return 0, fmt.Errorf("file size overflow")
}
totalSize += int64(file.UncompressedSize64)
```

**状态：** ⬜ 待修复

---

## 📊 统计信息

| 问题类型 | 数量 | 优先级 | 预计时间 |
|---------|------|--------|---------|
| G304 - 文件包含 | 58 | 高 | 2小时 |
| G306 - 文件权限 | 25 | 中 | 1小时 |
| G301 - 目录权限 | 7 | 中 | 20分钟 |
| G204 - 子进程 | 10 | 低 | 15分钟 |
| G110 - 解压炸弹 | 1 | 高 | 15分钟 |
| G115 - 整数溢出 | 1 | 高 | 10分钟 |
| **总计** | **102个** | - | **4.5小时** |

---

## 🔧 修复模板

### 模板1：文件路径验证
```go
func validatePath(path string) error {
    if !filepath.IsAbs(path) {
        return fmt.Errorf("path must be absolute: %s", path)
    }
    if strings.Contains(path, "..") {
        return fmt.Errorf("path traversal detected: %s", path)
    }
    return nil
}

// 使用
if err := validatePath(userInput); err != nil {
    return err
}
data, err := os.ReadFile(userInput)
```

### 模板2：文件权限
```go
// 生产代码
os.WriteFile(path, data, 0600)  // 仅所有者可读写
os.MkdirAll(path, 0750)         // 仅所有者可读写执行

// 测试代码（可放宽）
os.WriteFile(path, data, 0644)
os.MkdirAll(path, 0755)
```

### 模板3：解压大小限制
```go
const maxExtractSize = 100 * 1024 * 1024 // 100MB

if file.UncompressedSize64 > maxExtractSize {
    return fmt.Errorf("file too large: %d > %d",
        file.UncompressedSize64, maxExtractSize)
}
```

### 模板4：整数溢出检查
```go
import "math"

if file.UncompressedSize64 > uint64(math.MaxInt64) {
    return 0, fmt.Errorf("size overflow")
}
totalSize += int64(file.UncompressedSize64)
```

### 模板5：子进程（测试）
```go
// #nosec G204 - 测试环境，变量来源可信
cmd := exec.Command(executable, "keygen", "-d", testDir, "-n", keyPrefix)
```

---

## ✅ 验证标准

修复后运行：
```bash
golangci-lint run --disable-all --enable=gosec
```

应输出：`0 issues` 或仅低优先级警告

---

**创建时间：** 2025-12-30
**预计完成：** 2025-12-31
**负责人：** 待分配
