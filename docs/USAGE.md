# 使用文档

本文档提供 fzj 的完整使用指南，包括所有命令的详细说明、参数解释和实际示例。

## 📖 概述

fzj 是一个命令行工具，提供 8 个核心命令：

1. **keygen** - 生成密钥对（Kyber+ECDH+Dilithium）
2. **encrypt** - 加密文件（混合加密 + 签名）
3. **decrypt** - 解密文件（验证 + 恢复）
4. **encrypt-dir** - 加密文件夹（打包 + 加密）
5. **decrypt-dir** - 解密文件夹（解密 + 解包）
6. **info** - 查看加密文件信息
7. **keymanage** - 密钥管理（导出/导入/验证/缓存信息）
8. **version** - 版本信息

### 国际化支持

工具支持多语言（自动检测 `LANG` 环境变量）：
- **简体中文** (`zh_CN`) - 默认
- **English** (`en_US`)

```bash
# 使用中文
export LANG=zh_CN
fzj encrypt --help

# 使用英文
export LANG=en_US
fzj encrypt --help
```

### 全局标志

所有命令支持以下全局标志：

- `-v, --verbose`: 启用详细输出，显示更多信息
- `-f, --force`: 强制覆盖现有文件
- `--help`: 显示帮助信息

---

## 🔑 keygen - 生成密钥对

### 语法
```bash
fzj keygen -d <目录> -n <名称> [flags]
```

### 参数说明

| 参数 | 简写 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--dir` | `-d` | string | ✅ | - | 密钥保存目录 |
| `--name` | `-n` | string | ✅ | - | 密钥名称 |
| `--force` | `-f` | bool | ❌ | false | 覆盖现有文件 |

### 生成的文件

执行成功后，会在指定目录生成 4 个 PEM 文件：

```
{目录}/{名称}_public.pem              # Kyber+ECDH 公钥
{目录}/{名称}_private.pem             # Kyber+ECDH 私钥 (0600权限)
{目录}/{名称}_dilithium_public.pem    # Dilithium 公钥
{目录}/{名称}_dilithium_private.pem   # Dilithium 私钥 (0600权限)
```

### 使用示例

#### 基本使用
```bash
fzj keygen -d ./keys -n mykey
```

**输出**:
```
[1/4] 生成 Kyber768 密钥... 完成
[2/4] 生成 ECDH X25519 密钥... 完成
[3/4] 生成 Dilithium3 签名密钥... 完成
[4/4] 保存密钥文件... 完成

✅ 密钥对生成成功！

生成的文件:
  • mykey_public.pem (公钥)
  • mykey_private.pem (私钥 - 0600权限)
  • mykey_dilithium_public.pem (签名公钥)
  • mykey_dilithium_private.pem (签名私钥 - 0600权限)

⚠️  安全提示:
  • 请妥善保管私钥文件
  • 不要将私钥分享给他人
  • 建议使用安全的存储介质
```

#### 强制覆盖
```bash
fzj keygen -d ./keys -n mykey --force
```

#### 不同目录
```bash
# Linux/macOS
fzj keygen -d ~/.secure/keys -n production

# Windows
fzj keygen -d C:\Secure\Keys -n production
```

### 安全提示

⚠️ **重要安全建议**:
- 私钥文件会自动设置为 0600 权限（仅所有者可读写）
- 请将密钥存储在安全位置
- 定期轮换密钥（建议 3-6 个月）
- 不要将私钥提交到版本控制系统

---

## 🔐 encrypt - 加密文件

### 语法
```bash
fzj encrypt -i <输入文件> -o <输出文件> -p <公钥> -s <签名私钥> [flags]
```

### 参数说明

| 参数 | 简写 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--input` | `-i` | string | ✅ | - | 输入文件路径 |
| `--output` | `-o` | string | ❌ | `{输入}.fzj` | 输出文件路径 |
| `--public-key` | `-p` | string | ✅ | - | Kyber+ECDH 公钥文件 |
| `--sign-key` | `-s` | string | ✅ | - | Dilithium 签名私钥文件 |
| `--force` | `-f` | bool | ❌ | false | 覆盖输出文件 |
| `--verbose` | `-v` | bool | ❌ | false | 显示详细信息 |

### 加密流程

1. **读取输入文件**
2. **密钥封装**: Kyber768 + ECDH 混合密钥交换
3. **数据加密**: AES-256-GCM 认证加密
4. **数字签名**: Dilithium3 签名（可选）
5. **构建文件头**: 包含元数据和验证信息
6. **写入加密文件**

### 使用示例

#### 基本加密
```bash
fzj encrypt -i secret.txt -o secret.fzj \
  -p keys/mykey_public.pem \
  -s keys/mykey_dilithium_private.pem
```

**输出**:
```
📖 读取输入文件: secret.txt (1.2 KB)
🔑 加密公钥: keys/mykey_public.pem
📝 签名私钥: keys/mykey_dilithium_private.pem

[████████████████████] 100% - 加密完成

✅ 加密成功！
📁 输出文件: secret.fzj (4.5 KB)
📊 加密耗时: 3ms
🔐 文件格式: fzj v1.0
```

#### 不签名的加密
```bash
fzj encrypt -i data.txt -o data.fzj -p keys/mykey_public.pem
```

#### 强制覆盖
```bash
fzj encrypt -i secret.txt -o secret.fzj \
  -p keys/mykey_public.pem \
  -s keys/mykey_dilithium_private.pem \
  --force
```

#### 详细输出
```bash
fzj encrypt -i secret.txt -o secret.fzj \
  -p keys/mykey_public.pem \
  -s keys/mykey_dilithium_private.pem \
  -v
```

**详细输出**:
```
加密文件: secret.txt
  输入: secret.txt
  输出: secret.fzj
  公钥: keys/mykey_public.pem
  签名密钥: keys/mykey_dilithium_private.pem

[1/3] 加载密钥... 完成
[2/3] 加密文件... 完成
[3/3] 验证... 完成

✅ 加密成功！

文件信息:
  原始文件: secret.txt (1234 bytes)
  加密文件: secret.fzj (4576 bytes)
  压缩率: 370.8%
```

#### 默认输出文件名
```bash
# 不指定 -o，会自动使用 input.fzj
fzj encrypt -i secret.txt -p keys/public.pem -s keys/private.pem
# 输出: secret.fzj
```

### 性能提示

- **大文件支持**: 工具使用流式处理，支持任意大小文件
- **内存占用**: 内存占用与文件大小无关
- **进度显示**: 大文件会自动显示进度条
- **密钥缓存**: 自动缓存密钥，后续操作速度提升 1000x+
- **并行优化**: 密钥生成使用多核 CPU 并行处理

### 缓存机制

工具内置智能密钥缓存系统：

```bash
# 第一次加载密钥（约 1-2ms）
fzj encrypt -i file.txt -o file.fzj -p pub.pem -s priv.pem

# 第二次加载（<1μs，几乎瞬间）
fzj encrypt -i file2.txt -o file2.fzj -p pub.pem -s priv.pem
```

**缓存特性**:
- ✅ **自动启用**: 无需配置
- ✅ **TTL 过期**: 1 小时后自动失效
- ✅ **大小限制**: 最多缓存 100 个密钥
- ✅ **后台清理**: 每 5 分钟清理过期条目
- ✅ **线程安全**: 支持并发访问

**缓存状态查看**:
```bash
# 查看缓存信息（需要启用 verbose 模式）
fzj keymanage -a cache-info
```

---

## 🔓 decrypt - 解密文件

### 语法
```bash
fzj decrypt -i <输入文件> -o <输出文件> -p <私钥> -s <验证公钥> [flags]
```

### 参数说明

| 参数 | 简写 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--input` | `-i` | string | ✅ | - | 加密文件路径 |
| `--output` | `-o` | string | ❌ | 原文件名 | 输出文件路径 |
| `--private-key` | `-p` | string | ✅ | - | Kyber+ECDH 私钥文件 |
| `--verify-key` | `-s` | string | ❌ | - | Dilithium 验证公钥文件 |
| `--force` | `-f` | bool | ❌ | false | 覆盖输出文件 |
| `--verbose` | `-v` | bool | ❌ | false | 显示详细信息 |

### 解密流程

1. **读取加密文件**
2. **解析文件头**: 验证魔数、版本、格式
3. **密钥解封装**: Kyber768 + ECDH 密钥恢复
4. **数据解密**: AES-256-GCM 解密（自动验证认证标签）
5. **哈希验证**: SHA256 完整性检查
6. **签名验证**: Dilithium3 签名验证（如果提供 `--verify-key` 则强制执行）
7. **恢复文件**: 使用原始文件名或指定路径

### 使用示例

#### 完整解密（带签名验证）
```bash
fzj decrypt -i secret.fzj -o recovered.txt \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem
```

**输出**:
```
📖 读取加密文件: secret.fzj (4.5 KB)
🔑 解密私钥: keys/mykey_private.pem
🔍 验证公钥: keys/mykey_dilithium_public.pem

[████████████████████] 100% - 解密完成

✅ 解密成功！
📁 输出文件: recovered.txt (1.2 KB)
📊 解密耗时: 4ms
✅ 哈希验证: 通过
✅ 签名验证: 通过
```

#### 不验证签名的解密
```bash
fzj decrypt -i data.fzj -o data.txt -p keys/mykey_private.pem
```

**输出**:
```
📖 读取加密文件: data.fzj (4.5 KB)
🔑 解密私钥: keys/mykey_private.pem
⚠️  警告: 未提供签名验证密钥，将跳过签名验证
完成

[████████████████████] 100% - 解密完成

✅ 解密成功！
📁 输出文件: data.txt (1.2 KB)
📊 解密耗时: 3ms
✅ 哈希验证: 通过
⚠️  签名验证: 跳过
```

#### 使用原始文件名
```bash
# 不指定 -o，会自动使用文件头中的原始文件名（已做 basename 安全清洗）
fzj decrypt -i secret.fzj -p keys/mykey_private.pem -s keys/public.pem
# 输出: secret.txt (原始文件名)
```

#### 强制覆盖
```bash
fzj decrypt -i secret.fzj -o recovered.txt \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem \
  --force
```

#### 详细输出
```bash
fzj decrypt -i secret.fzj -o recovered.txt \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem \
  -v
```

### 错误处理

#### 密钥不匹配
```
❌ 解密失败: AES-GCM decryption failed: Authentication failed

可能原因:
  1. 密钥不匹配（使用了错误的私钥）
  2. 文件已损坏或被篡改
  3. 文件格式不正确

安全提示:
  - 如果提示哈希不匹配，文件可能已被篡改，请勿使用
  - 始终提供签名验证密钥以确保文件完整性
```

#### 文件损坏
```
❌ 文件头解析失败: Failed to read Kyber encapsulation: unexpected EOF

可能原因:
  1. 文件传输过程中损坏
  2. 文件被截断
  3. 不是有效的 fzj 加密文件

建议:
  - 重新传输文件
  - 使用 info 命令验证文件格式
  - 检查文件大小是否完整
```

#### 签名验证失败
```
❌ 签名验证失败: Dilithium signature verification failed

可能原因:
  1. 使用了错误的验证公钥
  2. 文件被伪造或篡改
  3. 签名密钥不匹配

安全警告:
  ⚠️  文件可能已被篡改！请勿使用解密后的数据
  - 检查公钥是否与加密时使用的签名私钥匹配
  - 确认文件来源可信
```

#### 密钥加载失败
```
❌ 加载公钥失败: open keys/public.pem: no such file or directory

提示:
  1. 请检查公钥文件路径是否正确: keys/public.pem
  2. 确保公钥文件格式正确（PEM 格式）
  3. 检查文件权限（需可读）
  4. 如果是首次使用，请先生成密钥对: fzj keygen
```

---

## 🔐 encrypt-dir - 加密文件夹

### 语法
```bash
fzj encrypt-dir -i <输入目录> -o <输出文件> -p <公钥> -s <签名私钥> [flags]
```

### 参数说明

| 参数 | 简写 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--input` | `-i` | string | ✅ | - | 输入目录路径 |
| `--output` | `-o` | string | ❌ | `{输入}.fzj` | 输出加密文件路径 |
| `--public-key` | `-p` | string | ✅ | - | Kyber+ECDH 公钥文件 |
| `--sign-key` | `-s` | string | ✅ | - | Dilithium 签名私钥文件 |
| `--force` | `-f` | bool | ❌ | false | 覆盖输出文件 |
| `--verbose` | `-v` | bool | ❌ | false | 显示详细信息 |

### 加密流程

1. **扫描目录**: 递归扫描所有文件和子目录
2. **打包 ZIP**: 将整个目录打包成 ZIP 归档
3. **密钥封装**: Kyber768 + ECDH 混合密钥交换
4. **数据加密**: AES-256-GCM 认证加密
5. **数字签名**: Dilithium3 签名（可选）
6. **写入加密文件**: 保存为单个 .fzj 文件

### 使用示例

#### 基本加密
```bash
fzj encrypt-dir -i ./my_project -o project_backup.fzj \
  -p keys/mykey_public.pem \
  -s keys/mykey_dilithium_private.pem
```

**输出**:
```
📁 扫描目录: ./my_project
📦 打包文件: 15 个文件 (总计 2.3 MB)
🔑 加密公钥: keys/mykey_public.pem
📝 签名私钥: keys/mykey_dilithium_private.pem

[████████████████████] 100% - 加密完成

✅ 加密成功！
📁 输出文件: project_backup.fzj (3.1 MB)
📊 加密耗时: 120ms
🔐 文件格式: fzj v1.0
```

#### 详细输出
```bash
fzj encrypt-dir -i ./data -o data.fzj \
  -p keys/public.pem \
  -s keys/private.pem \
  -v
```

**详细输出**:
```
加密目录: ./data
  输入: ./data
  输出: data.fzj
  公钥: keys/public.pem
  签名密钥: keys/private.pem

[1/4] 扫描目录... 完成 (15 文件)
[2/4] 打包 ZIP... 完成 (2.3 MB)
[3/4] 加密... 完成
[4/4] 验证... 完成

✅ 加密成功！

文件信息:
  原始大小: 2.3 MB
  加密大小: 3.1 MB
  文件数量: 15
```

---

## 🔓 decrypt-dir - 解密文件夹

### 语法
```bash
fzj decrypt-dir -i <输入文件> -o <输出目录> -p <私钥> -s <验证公钥> [flags]
```

### 参数说明

| 参数 | 简写 | 类型 | 必需 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--input` | `-i` | string | ✅ | - | 加密文件路径 |
| `--output` | `-o` | string | ❌ | 原目录名 | 输出目录路径 |
| `--private-key` | `-p` | string | ✅ | - | Kyber+ECDH 私钥文件 |
| `--verify-key` | `-s` | string | ❌ | - | Dilithium 验证公钥文件 |
| `--force` | `-f` | bool | ❌ | false | 覆盖现有文件 |
| `--verbose` | `-v` | bool | ❌ | false | 显示详细信息 |

### 解密流程

1. **读取加密文件**
2. **解析文件头**: 验证魔数、版本、格式
3. **密钥解封装**: Kyber768 + ECDH 密钥恢复
4. **数据解密**: AES-256-GCM 解密
5. **哈希验证**: SHA256 完整性检查
6. **签名验证**: Dilithium3 签名验证（如果提供）
7. **解压 ZIP**: 解压到指定目录
8. **恢复目录结构**: 保持原始目录层级

### 使用示例

#### 完整解密
```bash
fzj decrypt-dir -i project_backup.fzj -o restored_project \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem
```

**输出**:
```
📖 读取加密文件: project_backup.fzj (3.1 MB)
🔑 解密私钥: keys/mykey_private.pem
🔍 验证公钥: keys/mykey_dilithium_public.pem

[████████████████████] 100% - 解密完成
[████████████████████] 100% - 解压完成

✅ 解密成功！
📁 输出目录: restored_project
📊 解密耗时: 150ms
✅ 哈希验证: 通过
✅ 签名验证: 通过
📦 文件数量: 15
```

#### 使用原始目录名
```bash
# 不指定 -o，会自动使用加密文件中的原始目录名
fzj decrypt-dir -i project_backup.fzj \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem
# 输出: my_project (原始目录名)
```

#### 不验证签名
```bash
fzj decrypt-dir -i data.fzj -o restored \
  -p keys/mykey_private.pem
```

**输出**:
```
📖 读取加密文件: data.fzj (3.1 MB)
🔑 解密私钥: keys/mykey_private.pem
⚠️  警告: 未提供签名验证密钥，将跳过签名验证

[████████████████████] 100% - 解密完成
[████████████████████] 100% - 解压完成

✅ 解密成功！
📁 输出目录: restored
📊 解密耗时: 145ms
✅ 哈希验证: 通过
⚠️  签名验证: 跳过
```

### 安全特性

#### 路径遍历防护
工具会自动检测并阻止恶意 ZIP 文件尝试逃逸到父目录：

```
❌ 安全警告: 检测到路径遍历攻击，已阻止
  恶意路径: ../../etc/passwd
  已自动清理并拒绝解压
```

---

## ℹ️ info - 查看文件信息

### 语法
```bash
fzj info -i <加密文件>
```

### 参数说明

| 参数 | 简写 | 类型 | 必需 | 说明 |
|------|------|------|------|------|
| `--input` | `-i` | string | ✅ | 加密文件路径 |

### 显示的信息

- **基础信息**: 文件名、大小、时间戳、格式版本
- **算法信息**: 密钥封装、数据加密、数字签名算法
- **密钥信息**: Kyber 密文长度、ECDH 公钥、IV、签名长度
- **完整性**: SHA256 哈希、签名状态、验证结果

### 使用示例

#### 查看信息
```bash
fzj info -i secret.fzj
```

**输出**:
```
📁 文件信息: secret.fzj

基础信息:
  文件名:        secret.fzj
  原始大小:      1234 bytes
  加密大小:      4576 bytes
  压缩率:        370.8%
  时间戳:        2025-12-21 19:33:55
  格式版本:      fzj v1.0

算法信息:
  算法:          Kyber768 + ECDH + AES-256-GCM (0x02)
  版本:          0x0100
  魔数:          FZJ\x01

密钥信息:
  Kyber封装:     1088 bytes
  ECDH公钥:      32 bytes
  IV/Nonce:      12 bytes
  签名:          3293 bytes

完整性:
  SHA256哈希:    95fe0f8f44d17d26...
  签名状态:      ✅ 存在

验证状态:
  签名:          ✅ 存在
  数据完整性:   ✅ 完整
```

#### 查看未签名的文件
```bash
fzj info -i data.fzj
```

**输出**:
```
📁 文件信息: data.fzj

基础信息:
  文件名:        data.fzj
  原始大小:      1024 bytes
  加密大小:      4352 bytes
  时间戳:        2025-12-21 19:35:10

算法信息:
  算法:          Kyber768 + ECDH + AES-256-GCM (0x02)

完整性:
  SHA256哈希:    a1b2c3d4e5f6...

验证状态:
  签名:          ❌ 未签名
  数据完整性:   ✅ 完整
```

### 信息用途

- **验证文件**: 确认文件是否为有效的 fzj 加密文件
- **检查完整性**: 验证文件未被篡改
- **获取元数据**: 了解加密参数和原始文件信息
- **调试问题**: 排查加密/解密问题

---

## 🔐 keymanage - 密钥管理

### 语法
```bash
fzj keymanage -a <动作> [参数]
```

### 动作类型

| 动作 | 说明 | 必需参数 |
|------|------|----------|
| `export` | 从私钥导出公钥 | `-s` 私钥, `-o` 输出 |
| `verify` | 验证密钥对匹配 | `-p` 公钥, `-s` 私钥 |
| `import` | 导入密钥到目录 | `-p` 公钥, `-s` 私钥, `-d` 目录 |
| `cache-info` | 查看缓存信息 | 无 |

### 1. export - 导出公钥

从私钥文件中提取并导出对应的公钥。

**语法**:
```bash
fzj keymanage -a export -s <私钥文件> -o <输出公钥文件>
```

**示例**:
```bash
fzj keymanage -a export \
  -s keys/mykey_private.pem \
  -o keys/extracted_public.pem
```

**输出**:
```
导出公钥...
✅ 公钥已导出到: keys/extracted_public.pem
```

**用途**:
- 从私钥备份中恢复公钥
- 分享公钥给他人

### 2. verify - 验证密钥对

验证公钥和私钥是否匹配。

**语法**:
```bash
fzj keymanage -a verify -p <公钥文件> -s <私钥文件>
```

**示例**:
```bash
fzj keymanage -a verify \
  -p keys/mykey_public.pem \
  -s keys/mykey_private.pem
```

**输出** (匹配):
```
验证密钥对...
✅ 密钥对验证通过
  Kyber:  ✅ 匹配
  ECDH:   ✅ 匹配
```

**输出** (不匹配):
```
验证密钥对...
❌ 密钥对不匹配
  Kyber:  ❌ 不匹配
  ECDH:   ❌ 不匹配
```

**用途**:
- 确认密钥对是否正确配对
- 检查密钥文件是否损坏
- 验证备份的完整性

### 3. import - 导入密钥

将密钥文件导入到指定目录，自动处理权限。

**语法**:
```bash
fzj keymanage -a import \
  -p <公钥文件> \
  -s <私钥文件> \
  -d <目标目录>
```

**示例**:
```bash
fzj keymanage -a import \
  -p keys/mykey_public.pem \
  -s keys/mykey_private.pem \
  -d ./backup/keys
```

**输出**:
```
导入密钥...
✅ 密钥已导入到: ./backup/keys
  公钥: mykey_public.pem
  私钥: mykey_private.pem
```

**用途**:
- 备份密钥到安全位置
- 迁移密钥到新环境
- 整理密钥文件

### 4. cache-info - 查看缓存信息

查看当前密钥缓存的状态和统计信息。

**语法**:
```bash
fzj keymanage -a cache-info
```

**示例**:
```
缓存信息:
  总条目: 3
  已过期: 1
  估算大小: 300 bytes
```

---

## 📊 version - 版本信息

### 语法
```bash
fzj version
```

### 示例输出
```
fzj - 后量子文件加密工具
版本: 0.1.0
应用名称: fzj
描述: 后量子文件加密工具 - 使用 Kyber768 + ECDH + AES-256-GCM + Dilithium3
```

---

## 🔄 完整工作流示例

### 场景: 安全文件传输

#### 步骤 1: 发送方生成密钥对
```bash
# 创建工作目录
mkdir -p ~/secure_transfer
cd ~/secure_transfer

# 生成密钥对
fzj keygen -d ./keys -n sender
```

#### 步骤 2: 发送方加密文件
```bash
# 准备要发送的文件
echo "机密信息: 项目预算 100万" > budget.txt

# 加密（使用接收方的公钥）
# 假设接收方已经提供了 public.pem
fzj encrypt -i budget.txt -o budget.fzj \
  -p receiver_public.pem \
  -s keys/sender_dilithium_private.pem
```

#### 步骤 3: 传输加密文件
```bash
# 通过任何渠道发送 budget.fzj
# 邮件、网盘、即时通讯等
```

#### 步骤 4: 接收方解密
```bash
# 接收方使用自己的私钥解密
fzj decrypt -i budget.fzj -o budget_decrypted.txt \
  -p receiver_private.pem \
  -s sender_public.pem

# 查看解密内容
cat budget_decrypted.txt
```

#### 步骤 5: 验证
```bash
# 查看加密文件信息
fzj info -i budget.fzj

# 验证密钥对
fzj keymanage -a verify \
  -p receiver_public.pem \
  -s receiver_private.pem
```

### 场景: 安全备份（目录）

```bash
# 1. 生成专用备份密钥
fzj keygen -d ./backup_keys -n backup_2025

# 2. 加密整个目录
fzj encrypt-dir -i /important/data -o data_backup.fzj \
  -p backup_keys/backup_2025_public.pem \
  -s backup_keys/backup_2025_dilithium_private.pem

# 3. 上传到云端
aws s3 cp data_backup.fzj s3://my-backup/

# 4. 定期验证
fzj info -i data_backup.fzj

# 5. 恢复时解密
fzj decrypt-dir -i data_backup.fzj -o /restore/data \
  -p backup_keys/backup_2025_private.pem \
  -s backup_keys/backup_2025_dilithium_public.pem
```

### 场景: 安全备份（先打包再加密）

```bash
# 1. 生成专用备份密钥
fzj keygen -d ./backup_keys -n backup_2025

# 2. 打包备份文件
tar -czf data_backup.tar.gz /important/data/

# 3. 加密备份包
fzj encrypt -i data_backup.tar.gz -o data_backup.tar.gz.fzj \
  -p backup_keys/backup_2025_public.pem \
  -s backup_keys/backup_2025_dilithium_private.pem

# 4. 上传到云端
aws s3 cp data_backup.tar.gz.fzj s3://my-backup/

# 5. 定期验证
fzj info -i data_backup.tar.gz.fzj

# 6. 恢复时解密
fzj decrypt -i data_backup.tar.gz.fzj -o data_backup.tar.gz \
  -p backup_keys/backup_2025_private.pem \
  -s backup_keys/backup_2025_dilithium_public.pem
tar -xzf data_backup.tar.gz
```

---

## ⚠️ 错误处理

### 常见错误及解决方案

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| `file not found` | 文件不存在 | 检查文件路径和名称 |
| `permission denied` | 权限不足 | 检查文件权限，添加执行权限 |
| `authentication failed` | 密钥不匹配 | 使用正确的密钥对 |
| `invalid file format` | 文件损坏或不是加密文件 | 检查文件完整性 |
| `signature verification failed` | 签名验证失败 | 检查签名密钥 |
| `output file exists` | 输出文件已存在 | 使用 `--force` 或指定新路径 |
| `input/output error` | 文件读写失败 | 检查磁盘空间和文件权限 |

### 错误处理改进 (v0.2.0+)

工具在 v0.2.0 中增强了错误处理机制：

#### 1. 文件操作错误处理
- **改进点**: `archive.go` 中的 defer 错误处理
- **效果**: 确保文件关闭错误能正确返回，提升数据完整性
- **场景**: 目录加密/解密时的文件操作

#### 2. 用户输入错误处理
- **改进点**: `errors.go` 中的 ConfirmPrompt 增强
- **效果**: 处理输入失败情况，返回安全的默认值
- **场景**: 交互式确认提示时的异常处理

#### 3. 错误信息优化
- **详细错误分类**: 8 种错误类型，清晰分类
- **解决方案建议**: 每个错误提供具体解决步骤
- **多语言支持**: 错误信息支持中英文切换

### 详细错误信息

使用 `--verbose` 标志获取详细错误信息：

```bash
fzj decrypt -i secret.fzj -o output.txt -p key.pem --verbose
```

---

## 💡 最佳实践

### 1. 密钥管理
- ✅ 使用强密码保护存储密钥的目录
- ✅ 定期轮换密钥（3-6个月）
- ✅ 备份私钥到安全的离线存储
- ❌ 不要通过不安全渠道传输私钥
- ❌ 不要将密钥提交到版本控制

### 2. 文件加密
- ✅ 加密前备份原始文件
- ✅ 使用签名验证文件来源
- ✅ 验证解密结果
- ❌ 不要加密已损坏的文件
- ❌ 不要在多用户系统上存储明文密钥

### 3. 目录加密
- ✅ 加密前检查目录大小
- ✅ 使用 encrypt-dir 处理多个文件
- ✅ 解密后验证目录结构
- ❌ 不要加密包含敏感临时文件的目录
- ❌ 不要加密系统目录

### 4. 性能优化
- ✅ 使用 SSD 存储大文件
- ✅ 关闭不必要的程序释放内存
- ✅ 使用 64 位系统处理大文件
- ✅ 定期清理临时文件

### 5. 安全考虑
- ✅ 在受信任的环境中操作
- ✅ 保持系统和依赖更新
- ✅ 使用防火墙和防病毒软件
- ✅ 定期审计加密文件

---

## 🔧 高级用法

### 批量加密脚本

```bash
#!/bin/bash
# batch_encrypt.sh

KEYS_DIR="./keys"
PUBLIC_KEY="$KEYS_DIR/mykey_public.pem"
SIGN_KEY="$KEYS_DIR/mykey_dilithium_private.pem"
OUTPUT_DIR="./encrypted"

mkdir -p "$OUTPUT_DIR"

for file in sensitive/*; do
  if [ -f "$file" ]; then
    filename=$(basename "$file")
    output="$OUTPUT_DIR/${filename}.fzj"
    echo "加密: $filename"
    fzj encrypt -i "$file" -o "$output" -p "$PUBLIC_KEY" -s "$SIGN_KEY"
  fi
done

echo "批量加密完成！"
```

### Windows PowerShell 批量加密

```powershell
# batch_encrypt.ps1
$keysDir = ".\keys"
$publicKey = "$keysDir\mykey_public.pem"
$signKey = "$keysDir\mykey_dilithium_private.pem"
$outputDir = ".\encrypted"

New-Item -ItemType Directory -Force -Path $outputDir

Get-ChildItem -Path ".\sensitive" -File | ForEach-Object {
  $inputFile = $_.FullName
  $outputFile = "$outputDir\$($_.Name).fzj"

  Write-Host "加密: $($_.Name)"
  fzj encrypt -i $inputFile -o $outputFile -p $publicKey -s $signKey
}

Write-Host "批量加密完成！"
```

### 密钥轮换脚本

```bash
#!/bin/bash
# rotate_keys.sh

OLD_NAME="old_key"
NEW_NAME="new_key"
KEYS_DIR="./keys"

# 生成新密钥
echo "生成新密钥..."
fzj keygen -d "$KEYS_DIR" -n "$NEW_NAME"

# 验证新密钥
echo "验证新密钥..."
fzj keymanage -a verify \
  -p "$KEYS_DIR/${NEW_NAME}_public.pem" \
  -s "$KEYS_DIR/${NEW_NAME}_private.pem"

# 备份旧密钥（安全存储）
echo "备份旧密钥..."
mkdir -p ./key_backup
cp "$KEYS_DIR/${OLD_NAME}"* ./key_backup/

echo "密钥轮换完成！"
echo "新公钥: $KEYS_DIR/${NEW_NAME}_public.pem"
echo "请更新所有使用旧公钥的系统"
```

### 目录批量加密脚本

```bash
#!/bin/bash
# batch_encrypt_dir.sh

KEYS_DIR="./keys"
PUBLIC_KEY="$KEYS_DIR/mykey_public.pem"
SIGN_KEY="$KEYS_DIR/mykey_dilithium_private.pem"
BASE_DIR="./projects"

# 遍历所有子目录
for dir in "$BASE_DIR"/*/; do
  if [ -d "$dir" ]; then
    dirname=$(basename "$dir")
    output="./encrypted/${dirname}.fzj"

    echo "加密目录: $dirname"
    fzj encrypt-dir -i "$dir" -o "$output" \
      -p "$PUBLIC_KEY" -s "$SIGN_KEY"
  fi
done

echo "所有目录加密完成！"
```

---

## 📈 性能参考

### 测试数据（单文件）

| 文件大小 | 加密耗时 | 解密耗时 | 加密后大小 |
|----------|----------|----------|------------|
| 1 KB | <1ms | <1ms | 4.5 KB |
| 100 KB | 3ms | 4ms | 400 KB |
| 1 MB | 35ms | 40ms | 984 KB |
| 10 MB | 350ms | 400ms | 9.8 MB |
| 100 MB | 3.5s | 4.0s | 98 MB |

### 测试数据（目录加密）

| 目录大小 | 文件数量 | 加密耗时 | 解密耗时 |
|----------|----------|----------|----------|
| 10 MB | 50 | 120ms | 150ms |
| 100 MB | 500 | 1.2s | 1.5s |
| 1 GB | 5000 | 12s | 15s |

**说明**:
- 时间包括打包、密钥封装、加密、签名
- 目录加密会先打包成 ZIP，再加密
- 性能与 CPU 性能和磁盘 I/O 相关
- **缓存加速**: 后续操作速度提升 1000x+

### 运行基准测试

```bash
# 完整基准测试
go test -bench=. -benchmem ./internal/crypto/

# 仅加密基准
go test -bench=Encrypt -benchmem ./internal/crypto/

# 仅缓存基准
go test -bench=Cache -benchmem ./internal/crypto/

# 目录打包基准
go test -bench=Archive -benchmem ./internal/crypto/
```

---

## 🆘 获取帮助

### 命令帮助
```bash
# 查看所有命令
fzj --help

# 查看特定命令帮助
fzj encrypt --help
fzj decrypt --help
fzj encrypt-dir --help
fzj decrypt-dir --help
fzj keygen --help
```

### 详细输出
```bash
# 使用 -v 查看详细过程
fzj encrypt -i input.txt -o output.fzj -p pub.pem -s priv.pem -v
```

### 错误诊断
```bash
# 详细错误信息
fzj decrypt -i file.fzj -p key.pem --verbose 2>&1 | tee error.log
```

---

**版本**: v1.0.4
**最后更新**: 2026-03-08
**维护者**: fzj 开发团队

---

## 📋 更新日志

### 2025-12-30 更新
#### 新增功能
- ✅ **目录加密命令**: `encrypt-dir` - 支持整个目录打包加密
- ✅ **目录解密命令**: `decrypt-dir` - 支持解密并恢复目录结构
- ✅ **路径遍历防护**: 自动检测并阻止恶意 ZIP 文件
- ✅ **缓存信息查看**: `keymanage -a cache-info` 查看缓存统计

#### 改进
- ✅ **国际化支持**: 完整的中英文双语支持
- ✅ **文档更新**: 添加目录加密完整示例
- ✅ **错误提示**: 更友好的错误信息和解决方案

### 2025-12-26 更新
#### 新增功能
- ✅ **智能密钥缓存**: 自动缓存 + TTL 过期 + 大小限制
- ✅ **详细错误提示**: 包含解决方案和安全建议
- ✅ **性能基准测试**: 完整的测试套件和性能分析

#### 改进
- ✅ **错误信息**: 更友好的提示和解决方案
- ✅ **文档更新**: 添加缓存机制说明和性能数据
- ✅ **代码质量**: 消除重复代码，提升可维护性

#### 性能提升
- **缓存加速**: 1000x+ 后续加载速度
- **并行优化**: 多核 CPU 利用
- **序列化优化**: 5x 头部序列化速度
