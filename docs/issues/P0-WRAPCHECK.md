# P0 - Wrapcheck 问题修复清单

**优先级：🔴 严重**
**数量：73个**
**风险：丢失错误堆栈信息，难以调试**
**状态：待修复**

---

## 📋 问题详情

### 1. cmd/fzjjyz/decrypt.go (3个)

#### 1.1 第99行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_private_key_failed", err, decryptPrivKey)
```
**修复方案：**
```go
return fmt.Errorf("load private key failed: %w",
    i18n.TranslateError("error.load_private_key_failed", err, decryptPrivKey))
```

#### 1.2 第107行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_verify_key_failed", err, decryptVerifyKey)
```
**修复方案：**
```go
return fmt.Errorf("load verify key failed: %w",
    i18n.TranslateError("error.load_verify_key_failed", err, decryptVerifyKey))
```

#### 1.3 第151行 - i18n.TranslateError
```go
return i18n.TranslateError("error.decrypt_failed", err)
```
**修复方案：**
```go
return fmt.Errorf("decrypt failed: %w",
    i18n.TranslateError("error.decrypt_failed", err))
```

**状态：** ⬜ 待修复
**预计时间：** 5分钟

---

### 2. cmd/fzjjyz/decrypt_dir.go (5个)

#### 2.1 第101行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_private_key_failed", err, decryptDirPrivKey)
```

#### 2.2 第109行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_verify_key_failed", err, decryptDirVerifyKey)
```

#### 2.3 第158行 - i18n.TranslateError
```go
return i18n.TranslateError("error.decrypt_failed", err)
```

#### 2.4 第170行 - i18n.TranslateError
```go
return i18n.TranslateError("error.cannot_read_data", err)
```

#### 2.5 第181行 - i18n.TranslateError
```go
return i18n.TranslateError("error.extract_failed", err)
```

**状态：** ⬜ 待修复
**预计时间：** 8分钟

---

### 3. cmd/fzjjyz/encrypt.go (3个)

#### 3.1 第79行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_public_key_failed", err, encryptPubKey)
```

#### 3.2 第85行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_sign_key_failed", err, encryptSignKey)
```

#### 3.3 第126行 - i18n.TranslateError
```go
return i18n.TranslateError("error.encrypt_failed", err)
```

**状态：** ⬜ 待修复
**预计时间：** 5分钟

---

### 4. cmd/fzjjyz/encrypt_dir.go (5个)

#### 4.1 第80行 - i18n.TranslateError
```go
return i18n.TranslateError("error.pack_failed", err)
```

#### 4.2 第94行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_public_key_failed", err, encryptDirPubKey)
```

#### 4.3 第100行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_sign_key_failed", err, encryptDirSignKey)
```

#### 4.4 第111行 - i18n.TranslateError
```go
return i18n.TranslateError("error.temp_file_failed", err)
```

#### 4.5 第149行 - i18n.TranslateError
```go
return i18n.TranslateError("error.encrypt_failed", err)
```

**状态：** ⬜ 待修复
**预计时间：** 8分钟

---

### 5. cmd/fzjjyz/keygen.go (5个)

#### 5.1 第74行 - i18n.TranslateError
```go
return i18n.TranslateError("error.keygen_kyber_failed", err)
```

#### 5.2 第83行 - i18n.TranslateError
```go
return i18n.TranslateError("error.keygen_ecdh_failed", err)
```

#### 5.3 第92行 - i18n.TranslateError
```go
return i18n.TranslateError("error.keygen_dilithium_failed", err)
```

#### 5.4 第102行 - i18n.TranslateError
```go
return i18n.TranslateError("error.save_keys_failed", err)
```

#### 5.5 第108行 - i18n.TranslateError
```go
return i18n.TranslateError("error.save_dilithium_failed", err)
```

**状态：** ⬜ 待修复
**预计时间：** 8分钟

---

### 6. cmd/fzjjyz/keymanage.go (8个)

#### 6.1 第68行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_private_key_failed", err, keymanagePrivKey)
```

#### 6.2 第78行 - i18n.TranslateError
```go
return i18n.TranslateError("error.export_key_failed", err)
```

#### 6.3 第83行 - i18n.TranslateError
```go
return i18n.TranslateError("error.save_export_failed", err)
```

#### 6.4 第110行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_public_key_failed", err, keymanagePubKey)
```

#### 6.5 第115行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_private_key_failed", err, keymanagePrivKey)
```

#### 6.6 第126行 - i18n.TranslateError
```go
return i18n.TranslateError("error.save_keys_failed", err)
```

#### 6.7 第148行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_public_key_failed", err, keymanagePubKey)
```

#### 6.8 第155行 - i18n.TranslateError
```go
return i18n.TranslateError("error.load_private_key_failed", err, keymanagePrivKey)
```

**状态：** ⬜ 待修复
**预计时间：** 12分钟

---

### 7. cmd/fzjjyz/utils/progress.go (2个)

#### 7.1 第134行 - Read
```go
return n, err
```
**修复方案：**
```go
return n, fmt.Errorf("read failed: %w", err)
```

#### 7.2 第162行 - Write
```go
return n, err
```
**修复方案：**
```go
return n, fmt.Errorf("write failed: %w", err)
```

**状态：** ⬜ 待修复
**预计时间：** 3分钟

---

### 8. internal/crypto/archive.go (15个)

#### 8.1 第58行 - filepath.Abs
```go
return err
```

#### 8.2 第62行 - filepath.Walk
```go
return filepath.Walk(absSource, func(path string, info os.FileInfo, walkErr error) error {
    // ...
})
```

#### 8.3 第80行 - os.Readlink
```go
return err
```

#### 8.4 第85行 - os.Stat
```go
return err
```

#### 8.5 第92行 - filepath.Rel
```go
return err
```

#### 8.6 第105行 - Writer.Create
```go
return err
```

#### 8.7 第111行 - Writer.Create
```go
return err
```

#### 8.8 第117行 - os.Open
```go
return err
```

#### 8.9 第126行 - io.Copy
```go
return err
```

#### 8.10 第177行 - os.MkdirAll
```go
return err
```

#### 8.11 第184行 - os.MkdirAll
```go
return err
```

#### 8.12 第190行 - File.Open
```go
return err
```

#### 8.13 第201行 - os.OpenFile
```go
return err
```

#### 8.14 第211行 - io.Copy
```go
return err
```

#### 8.15 第222行 - zip.NewReader
```go
return 0, err
```

**状态：** ⬜ 待修复
**预计时间：** 20分钟

---

### 9. internal/crypto/hash_utils.go (4个)

#### 9.1 第16行 - os.Open
```go
return result, err
```

#### 9.2 第26行 - io.Copy
```go
return result, err
```

#### 9.3 第40行 - io.Copy
```go
return result, err
```

#### 9.4 第64行 - hash.Write
```go
return sh.hash.Write(p)
```

**状态：** ⬜ 待修复
**预计时间：** 5分钟

---

### 10. internal/crypto/stream_utils.go (3个)

#### 10.1 第30行 - Write
```go
return n, err
```

#### 10.2 第67行 - Close
```go
return closer.Close()
```

#### 10.3 第104行 - io.Copy
```go
return written, hash, err
```

**状态：** ⬜ 待修复
**预计时间：** 5分钟

---

### 11. internal/format/header.go (18个)

#### 11.1 第39行 - Buffer.Write
```go
return nil, err
```

#### 11.2 第42行 - binary.Write
```go
return nil, err
```

#### 11.3 第45行 - Buffer.WriteByte
```go
return nil, err
```

#### 11.4 第48行 - Buffer.WriteByte
```go
return nil, err
```

#### 11.5 第51行 - Buffer.WriteByte
```go
return nil, err
```

#### 11.6 第57行 - Buffer.WriteString
```go
return nil, err
```

#### 11.7 第61行 - binary.Write
```go
return nil, err
```

#### 11.8 第64行 - binary.Write
```go
return nil, err
```

#### 11.9 第69行 - binary.Write
```go
return nil, err
```

#### 11.10 第73行 - Buffer.Write
```go
return nil, err
```

#### 11.11 第77行 - Buffer.WriteByte
```go
return nil, err
```

#### 11.12 第81行 - Buffer.Write
```go
return nil, err
```

#### 11.13 第85行 - Buffer.WriteByte
```go
return nil, err
```

#### 11.14 第89行 - Buffer.Write
```go
return nil, err
```

#### 11.15 第95行 - binary.Write
```go
return nil, err
```

#### 11.16 第99行 - Buffer.Write
```go
return nil, err
```

#### 11.17 第103行 - Buffer.Write
```go
return nil, err
```

**状态：** ⬜ 待修复
**预计时间：** 25分钟

---

### 12. internal/crypto/operations_shared.go (2个)

#### 12.1 第105行 - MarshalBinaryOptimized
```go
return header.MarshalBinaryOptimized()
```

#### 12.2 第127行 - os.ReadFile
```go
return 0, err
```

**状态：** ⬜ 待修复
**预计时间：** 3分钟

---

## 📊 统计信息

| 文件类型 | 数量 | 预计时间 |
|---------|------|---------|
| cmd/fzjjyz/ | 26个 | 44分钟 |
| internal/crypto/ | 24个 | 38分钟 |
| internal/format/ | 18个 | 25分钟 |
| cmd/fzjjyz/utils/ | 2个 | 3分钟 |
| **总计** | **73个** | **110分钟 (1.8小时)** |

---

## 🔧 修复模板

### 模板1：i18n.TranslateError
```go
// 原代码
return i18n.TranslateError("error.key", err, arg)

// 修复后
return fmt.Errorf("operation failed: %w",
    i18n.TranslateError("error.key", err, arg))
```

### 模板2：标准库错误
```go
// 原代码
return err

// 修复后
return fmt.Errorf("operation failed: %w", err)
```

### 模板3：接口方法
```go
// 原代码
return n, err

// 修复后
return n, fmt.Errorf("operation failed: %w", err)
```

### 模板4：多返回值
```go
// 原代码
return written, hash, err

// 修复后
return written, hash, fmt.Errorf("operation failed: %w", err)
```

---

## ✅ 验证标准

修复后运行：
```bash
golangci-lint run --disable-all --enable=wrapcheck
```

应输出：`0 issues`

---

**创建时间：** 2025-12-30
**预计完成：** 2025-12-30
**负责人：** 待分配
