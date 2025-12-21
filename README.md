# fzjjyz - 后量子文件加密工具

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-100%25-passing-brightgreen.svg)]()
[![Post-Quantum](https://img.shields.io/badge/Post%20Quantum-Crypto-purple.svg)]()

**fzjjyz** 是一个基于后量子密码学的文件加密工具，提供面向未来的安全保护。

## ✨ 核心特性

- 🔐 **混合加密**: Kyber768 + ECDH 双重密钥封装，结合后量子和传统安全性
- 🔒 **认证加密**: AES-256-GCM 提供机密性和完整性保护
- 📝 **数字签名**: Dilithium3 签名验证，确保文件来源可信
- ⚡ **高性能**: 1MB 文件加密 < 40ms，解密 < 50ms
- 🛡️ **安全优先**: 零信任架构，最小权限原则，私钥自动设置 0600 权限
- 🌍 **跨平台**: Windows/Linux/macOS 全支持
- 📦 **开箱即用**: 完整的 CLI 工具，6个核心命令

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

## 📊 性能指标

| 操作 | 文件大小 | 耗时 | 说明 |
|------|----------|------|------|
| 密钥生成 | - | ~450ms | Kyber + ECDH + Dilithium |
| 加密 | 1MB | ~35ms | 混合加密 + 签名 |
| 解密 | 1MB | ~40ms | 完整验证 |
| 信息查看 | 4.5KB | <10ms | 快速解析 |

**测试环境**: Windows 11, Go 1.25.4, AMD Ryzen 7

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
│   │   ├── hybrid.go        # 混合加密
│   │   ├── signature.go     # 签名系统
│   │   ├── operations.go    # 文件操作
│   │   ├── keygen.go        # 密钥生成
│   │   └── keyfile.go       # 密钥管理
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
# 密钥管理
fzjjyz keygen -d ./keys -n mykey

# 文件加密/解密
fzjjyz encrypt -i input.txt -o output.fzj -p keys/public.pem -s keys/dilithium_priv.pem
fzjjyz decrypt -i output.fzj -o recovered.txt -p keys/private.pem -s keys/dilithium_pub.pem

# 信息查看
fzjjyz info -i output.fzj

# 密钥管理
fzjjyz keymanage -a verify -p keys/public.pem -s keys/private.pem
fzjjyz keymanage -a export -s keys/private.pem -o extracted_public.pem
fzjjyz keymanage -a import -p keys/public.pem -s keys/private.pem -d ./backup

# 版本信息
fzjjyz version
```

## 🎯 使用场景

### 1. 安全文件传输
```bash
# 加密敏感文件，通过不安全渠道传输
fzjjyz encrypt -i sensitive.doc -o sensitive.fzj -p recipient_public.pem -s my_private.pem
# 发送 .fzj 文件，接收方使用私钥解密
```

### 2. 安全备份
```bash
# 加密备份文件
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

详细安全信息请查看 [SECURITY.md](SECURITY.md)。

## 🤝 参与贡献

欢迎各种形式的贡献！请先阅读 [贡献指南](CONTRIBUTING.md)。

### 贡献类型
- 🐛 报告 Bug
- 💡 提出新功能
- 📝 改进文档
- 🔧 提交代码
- ✅ 添加测试

### 快速开始
```bash
# 1. Fork 项目
# 2. 创建特性分支
git checkout -b feature/amazing-feature

# 3. 开发和测试
go test ./...
go build ./cmd/fzjjyz

# 4. 提交 PR
git push origin feature/amazing-feature
```

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

**当前版本**: v0.1.0
**最后更新**: 2025-12-21
**状态**: ✅ 生产就绪