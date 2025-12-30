# fzjjyz - 后量子文件加密工具

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-100%25-passing-brightgreen.svg)]()
[![Post-Quantum](https://img.shields.io/badge/Post%20Quantum-Crypto-purple.svg)]()
[![Security](https://img.shields.io/badge/Security-Audit--Ready-blue.svg)]()
[![Performance](https://img.shields.io/badge/Performance-Optimized-green.svg)]()

**fzjjyz** 是一个基于后量子密码学的文件加密工具，提供面向未来的安全保护。

> 🔔 **2025-12-30 更新**: v0.2.0 - 目录加密/解密、国际化支持（中英文）、路径遍历防护

## ✨ 核心特性

- 🔐 **混合加密**: Kyber768 + ECDH 双重密钥封装，结合后量子和传统安全性
- 🔒 **认证加密**: AES-256-GCM 提供机密性和完整性保护
- 📝 **数字签名**: Dilithium3 签名验证，确保文件来源可信
- ⚡ **高性能**: 1MB 文件加密 < 40ms，解密 < 50ms
- 🛡️ **安全优先**: 零信任架构，最小权限原则，私钥自动设置 0600 权限
- 🧰 **智能缓存**: 带 TTL 和大小限制的密钥缓存，防止内存泄漏
- 📊 **性能监控**: 内置基准测试，支持性能分析
- 🌍 **跨平台**: Windows/Linux/macOS 全支持
- 📦 **开箱即用**: 完整的 CLI 工具，8个核心命令（含目录加密）
- 💡 **友好提示**: 详细的错误信息和解决方案建议
- 🌍 **国际化**: 自动检测 LANG 环境变量，支持中英文
- 🛡️ **路径防护**: 自动检测并阻止 ZIP 路径遍历攻击

## 🚀 快速开始

### 1. 安装

```bash
# 克隆源码
git clone https://codeberg.org/jiangfire/fzjjyz
cd fzjjyz

# 构建二进制
go build -o fzjjyz ./cmd/fzjjyz

# 验证安装
./fzjjyz version
```

### 2. 生成密钥对

```bash
fzjjyz keygen -d ./keys -n mykey
```

**生成的文件:**
- `keys/mykey_public.pem` - Kyber+ECDH 公钥
- `keys/mykey_private.pem` - Kyber+ECDH 私钥 (0600)
- `keys/mykey_dilithium_public.pem` - Dilithium 公钥
- `keys/mykey_dilithium_private.pem` - Dilithium 私钥 (0600)

### 3. 加密文件

```bash
# 创建测试文件
echo "这是一个秘密消息" > secret.txt

# 加密
fzjjyz encrypt -i secret.txt -o secret.fzj \
  -p keys/mykey_public.pem \
  -s keys/mykey_dilithium_private.pem
```

### 4. 解密文件

```bash
# 解密
fzjjyz decrypt -i secret.fzj -o recovered.txt \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem

# 验证
diff secret.txt recovered.txt && echo "✅ 解密成功！"
```

## 🔧 技术架构

### 加密流程

```
原始文件
    ↓
[1] 密钥封装: Kyber768 + ECDH
    ↓ 生成: 32字节共享密钥
[2] 数据加密: AES-256-GCM
    ↓ 生成: 加密数据 + 认证标签
[3] 数字签名: Dilithium3 (可选)
    ↓ 生成: 签名
[4] 文件封装: 自定义二进制格式
    ↓ 输出: .fzj 文件
```

### 算法组合

| 算法 | 用途 | 标准 | 安全级别 |
|------|------|------|----------|
| **Kyber768** | 后量子密钥封装 | NIST PQC | AES-192 |
| **X25519 ECDH** | 传统密钥交换 | RFC 7748 | ~128位 |
| **AES-256-GCM** | 认证加密 | FIPS 197 | 256位 |
| **Dilithium3** | 数字签名 | NIST PQC | SHA384 |
| **SHA256** | 完整性校验 | FIPS 180-4 | 256位 |

### 安全特性

- ✅ **后量子安全**: Kyber 抵抗量子计算机攻击
- ✅ **前向保密**: 每次加密使用新临时密钥
- ✅ **双重保护**: Kyber + ECDH 双重密钥封装
- ✅ **认证加密**: AES-GCM 防止密文篡改
- ✅ **来源认证**: Dilithium3 签名验证
- ✅ **完整性校验**: SHA256 哈希验证
- ✅ **缓存安全**: TTL 过期 + 大小限制，防止内存耗尽
- ✅ **权限控制**: 私钥文件自动设置 0600 权限

## 📊 性能指标

| 操作 | 文件大小 | 耗时 | 说明 |
|------|----------|------|------|
| 密钥生成 | - | ~450ms | Kyber + ECDH + Dilithium (并行) |
| 加密 | 1MB | ~35ms | 混合加密 + 签名 |
| 解密 | 1MB | ~40ms | 完整验证 |
| 信息查看 | 4.5KB | <10ms | 快速解析 |
| 缓存加载 | - | <1μs | 内存命中 |

**测试环境**: Windows 11, Go 1.25.4, AMD Ryzen 7

### 基准测试

运行内置基准测试：
```bash
go test -bench=. -benchmem ./internal/crypto/
```

**性能特点**:
- ✅ 缓存加速: 后续加载速度提升 1000x+
- ✅ 内存优化: 缓存上限 100 个密钥，自动清理过期条目
- ✅ 并行优化: 密钥生成使用多核 CPU
- ✅ 流式处理: 支持大文件，缓冲区自动优化

## 📁 项目结构

```
fzjjyz/
├── cmd/fzjjyz/              # CLI 工具
│   ├── main.go              # 主入口
│   ├── encrypt.go           # 加密命令
│   ├── decrypt.go           # 解密命令
│   ├── keygen.go            # 密钥生成
│   ├── keymanage.go         # 密钥管理
│   ├── info.go              # 信息查看
│   ├── version.go           # 版本信息
│   ├── main_test.go         # 集成测试
│   └── utils/               # 工具模块
│       ├── progress.go      # 进度条
│       └── errors.go        # 错误处理
│
├── internal/                # 内部模块
│   ├── crypto/              # 密码学核心
│   │   ├── hybrid.go        # 混合加密 (Kyber + ECDH)
│   │   ├── signature.go     # 签名系统 (Dilithium3)
│   │   ├── operations.go    # 核心文件操作
│   │   ├── operations_shared.go  # 共享函数库
│   │   ├── operations_stream.go  # 流式操作接口
│   │   ├── stream_encrypt.go    # 流式加密器
│   │   ├── stream_decrypt.go    # 流式解密器
│   │   ├── keygen.go        # 密钥生成 (支持并行)
│   │   ├── keyfile.go       # 密钥管理 (带TTL缓存)
│   │   ├── buffer_pool.go   # 缓冲区池优化
│   │   ├── benchmark_test.go # 性能基准测试
│   │   └── *_test.go        # 单元测试
│   │
│   ├── format/              # 文件格式
│   │   ├── header.go        # 文件头结构
│   │   └── parser.go        # 解析器
│   │
│   └── utils/               # 工具函数
│       ├── errors.go        # 错误系统
│       └── logger.go        # 日志系统
│
├── test_cli/                # 测试数据
├── go.mod                   # 依赖管理
├── README.md                # 项目说明 (本文件)
├── INSTALL.md               # 安装指南
├── USAGE.md                 # 使用文档
├── DEVELOPMENT.md           # 开发指南
├── SECURITY.md              # 安全文档
├── CONTRIBUTING.md          # 贡献指南
└── CHANGELOG.md             # 变更记录
```

## 📚 文档导航

### 用户文档
- 📖 [INSTALL.md](INSTALL.md) - 安装和构建指南
- 📝 [USAGE.md](USAGE.md) - 完整命令参考和示例
- 🔒 [SECURITY.md](SECURITY.md) - 安全策略和最佳实践

### 开发文档
- 👨‍💻 [DEVELOPMENT.md](DEVELOPMENT.md) - 开发环境和指南
- 🤝 [CONTRIBUTING.md](CONTRIBUTING.md) - 贡献流程
- 📊 [CHANGELOG.md](CHANGELOG.md) - 版本历史

## 🛠️ 命令概览

```bash
# 1. 密钥管理
fzjjyz keygen -d ./keys -n mykey

# 2. 文件加密/解密
fzjjyz encrypt -i input.txt -o output.fzj -p keys/public.pem -s keys/dilithium_priv.pem
fzjjyz decrypt -i output.fzj -o recovered.txt -p keys/private.pem -s keys/dilithium_pub.pem

# 3. 目录加密/解密 (v0.2.0 新增)
fzjjyz encrypt-dir -i ./myproject -o project.fzj -p keys/public.pem -s keys/dilithium_priv.pem
fzjjyz decrypt-dir -i project.fzj -o restored -p keys/private.pem -s keys/dilithium_pub.pem

# 4. 信息查看
fzjjyz info -i output.fzj

# 5. 密钥管理
fzjjyz keymanage -a verify -p keys/public.pem -s keys/private.pem
fzjjyz keymanage -a export -s keys/private.pem -o extracted_public.pem
fzjjyz keymanage -a cache-info  # 查看缓存信息

# 6. 国际化 (v0.2.0 新增)
export LANG=en_US  # 切换到英文
export LANG=zh_CN  # 切换到中文

# 7. 高级选项
# 指定缓冲区大小（KB）
fzjjyz encrypt -i large.bin -o large.fzj -p pub.pem -s priv.pem --buffer-size 1024
# 强制覆盖
fzjjyz decrypt -i file.fzj -o out.txt -p priv.pem -s pub.pem --force
# 详细输出
fzjjyz encrypt -i file.txt -o file.fzj -p pub.pem -s priv.pem --verbose
```

### 错误提示示例

工具提供详细的错误信息和解决方案：

```
❌ 加载公钥失败: open keys/public.pem: no such file or directory

提示:
  1. 请检查公钥文件路径是否正确: keys/public.pem
  2. 确保公钥文件格式正确（PEM 格式）
  3. 检查文件权限（需可读）
  4. 如果是首次使用，请先生成密钥对: fzjjyz keygen
```

## 🎯 使用场景

### 1. 安全文件传输
```bash
# 加密敏感文件，通过不安全渠道传输
fzjjyz encrypt -i sensitive.doc -o sensitive.fzj -p recipient_public.pem -s my_private.pem
# 发送 .fzj 文件，接收方使用私钥解密
```

### 2. 安全备份（目录）
```bash
# 方式1: 直接加密目录 (v0.2.0 新增)
fzjjyz encrypt-dir -i ./important_data -o backup.fzj -p backup_public.pem -s backup_private.pem

# 方式2: 先打包再加密
tar -czf backup.tar.gz ./important_data/
fzjjyz encrypt -i backup.tar.gz -o backup.fzj -p backup_public.pem -s backup_private.pem

# 存储到云端或外部存储
```

### 3. 机密文档共享
```bash
# 团队成员间共享加密文档
fzjjyz encrypt -i project.docx -o project.fzj -p team_public.pem -s my_private.pem
# 团队成员使用各自私钥解密
```

## 🔒 安全警告

⚠️ **重要提示**:
- 这是一个研究性质的项目，虽然使用了行业标准加密算法
- 生产环境使用前请进行充分的安全评估
- 请妥善保管私钥文件，不要与他人分享
- 建议定期轮换密钥（3-6个月）

### 安全最佳实践

1. **密钥管理**
   - 私钥文件权限设置为 0600（仅所有者可读写）
   - 使用密钥缓存功能，避免频繁读取磁盘
   - 定期清理缓存：`fzjjyz keymanage -a clear-cache`

2. **文件完整性**
   - 始终提供签名验证密钥进行完整性检查
   - 解密时验证哈希和签名，防止篡改

3. **内存安全**
   - 缓存自动过期（默认1小时）
   - 缓存大小限制（最多100个密钥）
   - 后台自动清理过期条目

详细安全信息请查看 [SECURITY.md](SECURITY.md)。

## 🤝 参与贡献

欢迎各种形式的贡献！请先阅读 [贡献指南](CONTRIBUTING.md)。

### 贡献类型
- 🐛 报告 Bug
- 💡 提出新功能
- 📝 改进文档
- 🔧 提交代码
- ✅ 添加测试
- 📊 性能优化

### 开发工作流
```bash
# 1. 克隆并设置
git clone https://codeberg.org/jiangfire/fzjjyz
cd fzjjyz

# 2. 创建特性分支
git checkout -b feature/amazing-feature

# 3. 开发和测试
go test ./...                    # 运行所有测试
go test -bench=. ./internal/crypto/  # 运行基准测试
go build ./cmd/fzjjyz            # 构建二进制

# 4. 代码质量检查
go vet ./...
go fmt ./...

# 5. 提交 PR
git add .
git commit -m "feat: 添加新功能"
git push origin feature/amazing-feature
```

### 项目状态
- ✅ **P0 核心修复完成**: 并行密钥生成、代码去重、测试覆盖、Git清理
- ✅ **P1 安全增强**: 缓存TTL、大小限制、错误提示优化
- ✅ **P2 性能优化**: 基准测试、文档更新
- 🔄 **P3 进行中**: 高级功能开发

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🔗 项目链接

- **项目主页**: https://codeberg.org/jiangfire/fzjjyz
- **Issue 追踪**: https://codeberg.org/jiangfire/fzjjyz/issues
- **讨论区**: https://codeberg.org/jiangfire/fzjjyz/discussions

## 📞 联系方式

### 安全报告
发现安全问题？请发送邮件至: **security@jiangfire.com**

### 一般咨询
- 项目主页讨论区
- GitHub Issues
- 邮件联系

## 🙏 致谢

- [Cloudflare CIRCL](https://github.com/cloudflare/circl) - 后量子密码学库
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- Go 社区 - 优秀的标准库和工具链

---

**注意**: 这是一个后量子密码学研究项目，旨在探索和演示后量子加密技术。请在理解安全风险的前提下使用。

**当前版本**: v0.2.0
**最后更新**: 2025-12-30
**状态**: ✅ 生产就绪 (持续优化中)

---

## 📋 更新日志

### 2025-12-30 (v0.2.0) ✨
- ✅ **目录加密**: `encrypt-dir` 和 `decrypt-dir` 命令
- ✅ **国际化**: 自动检测 LANG，支持中英文双语
- ✅ **路径防护**: 自动检测并阻止 ZIP 路径遍历攻击
- ✅ **缓存信息**: `keymanage -a cache-info` 查看统计
- ✅ **测试覆盖**: 关键模块 >80% 测试覆盖率
- ✅ **文档更新**: 完整的使用文档和架构说明

### 2025-12-26 (v0.1.1)
- ✅ **安全增强**: 密钥缓存支持 TTL 和大小限制
- ✅ **错误优化**: 详细的错误提示和解决方案
- ✅ **性能测试**: 新增基准测试套件
- ✅ **代码质量**: 消除重复代码，提升可维护性
- ✅ **文档更新**: 完善使用指南和最佳实践

### 2025-12-21 (v0.1.0)
- ✅ 初始版本发布
- ✅ 核心加密功能
- ✅ 流式处理支持
- ✅ 完整的 CLI 工具