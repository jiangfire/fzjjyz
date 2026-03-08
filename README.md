# fzj - 后量子文件加密工具

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Post-Quantum](https://img.shields.io/badge/Post%20Quantum-Crypto-purple.svg)]()
[![Security](https://img.shields.io/badge/Security-Audit--Ready-blue.svg)]()

**fzj** 是一个基于后量子密码学的文件加密工具，提供面向未来的安全保护。

> 🔔 **v1.0.3** - 热修复版本（修复跨平台 Windows 路径解析导致的 CI 失败）

## 📚 文档导航

所有文档已归档到 `docs/` 目录：

### 🎯 快速开始
- **[docs/README_MAIN.md](docs/README_MAIN.md)** - 完整项目介绍和特性
- **[docs/INSTALL.md](docs/INSTALL.md)** - 安装和构建指南
- **[docs/USAGE.md](docs/USAGE.md)** - 完整命令参考和示例

### 🔧 开发文档
- **[docs/CONTRIBUTING.md](docs/CONTRIBUTING.md)** - 贡献指南
- **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - 系统架构
- **[docs/CODE_QUALITY.md](docs/CODE_QUALITY.md)** - 代码质量报告

### 📦 其他文档
- **[docs/CHANGELOG.md](docs/CHANGELOG.md)** - 版本历史
- **[docs/INDEX.md](docs/INDEX.md)** - 完整文档索引

## 🚀 快速开始

### 1. 安装
```bash
# 克隆源码
git clone https://codeberg.org/jiangfire/fzj
cd fzj

# 构建二进制
go build -o fzj ./cmd/fzj

# 验证安装
./fzj version
```

### 2. 生成密钥对
```bash
fzj keygen -d ./keys -n mykey
```

### 3. 加密文件
```bash
echo "这是一个秘密消息" > secret.txt
fzj encrypt -i secret.txt -o secret.fzj \
  -p keys/mykey_public.pem \
  -s keys/mykey_dilithium_private.pem
```

### 4. 解密文件
```bash
fzj decrypt -i secret.fzj -o recovered.txt \
  -p keys/mykey_private.pem \
  -s keys/mykey_dilithium_public.pem
```

## ✨ 核心特性

- 🔐 **混合加密**: Kyber768 + ECDH 双重密钥封装
- 🔒 **认证加密**: AES-256-GCM 提供机密性和完整性
- 📝 **数字签名**: Dilithium3 签名验证
- ⚡ **高性能**: 1MB 文件加密 < 40ms，解密 < 50ms
- 🛡️ **安全优先**: 零信任架构，最小权限原则
- 🧰 **智能缓存**: 带 TTL 和大小限制的密钥缓存
- 📊 **性能监控**: 内置基准测试
- 🌍 **跨平台**: Windows/Linux/macOS 全支持
- 📦 **开箱即用**: 完整的 CLI 工具，8个核心命令
- 💡 **友好提示**: 详细的错误信息和解决方案
- 🌍 **国际化**: 自动检测 LANG 环境变量，支持中英文
- 🛡️ **路径防护**: 自动检测并阻止 ZIP 路径遍历攻击

## 🔧 技术架构

### 算法组合
| 算法 | 用途 | 标准 | 安全级别 |
|------|------|------|----------|
| **Kyber768** | 后量子密钥封装 | NIST PQC | AES-192 |
| **X25519 ECDH** | 传统密钥交换 | RFC 7748 | ~128位 |
| **AES-256-GCM** | 认证加密 | FIPS 197 | 256位 |
| **Dilithium3** | 数字签名 | NIST PQC | SHA384 |
| **SHA256** | 完整性校验 | FIPS 180-4 | 256位 |

### 性能指标
| 操作 | 文件大小 | 耗时 | 说明 |
|------|----------|------|------|
| 密钥生成 | - | ~450ms | Kyber + ECDH + Dilithium |
| 加密 | 1MB | ~35ms | 混合加密 + 签名 |
| 解密 | 1MB | ~40ms | 完整验证 |
| 信息查看 | 4.5KB | <10ms | 快速解析 |
| 缓存加载 | - | <1μs | 内存命中 |

## 🛠️ 命令概览

```bash
# 1. 密钥管理
fzj keygen -d ./keys -n mykey

# 2. 文件加密/解密
fzj encrypt -i input.txt -o output.fzj -p keys/public.pem -s keys/dilithium_priv.pem
fzj decrypt -i output.fzj -o recovered.txt -p keys/private.pem -s keys/dilithium_pub.pem

# 3. 目录加密/解密 (v0.2.0 新增)
fzj encrypt-dir -i ./myproject -o project.fzj -p keys/public.pem -s keys/dilithium_priv.pem
fzj decrypt-dir -i project.fzj -o restored -p keys/private.pem -s keys/dilithium_pub.pem

# 4. 信息查看
fzj info -i output.fzj

# 5. 密钥管理
fzj keymanage -a verify -p keys/public.pem -s keys/private.pem
fzj keymanage -a export -s keys/private.pem -o extracted_public.pem
fzj keymanage -a cache-info  # 查看缓存信息

# 6. 国际化 (v0.2.0 新增)
export LANG=en_US  # 切换到英文
export LANG=zh_CN  # 切换到中文
```

## 🎯 使用场景

### 1. 安全文件传输
```bash
# 加密敏感文件，通过不安全渠道传输
fzj encrypt -i sensitive.doc -o sensitive.fzj -p recipient_public.pem -s my_private.pem
# 发送 .fzj 文件，接收方使用私钥解密
```

### 2. 安全备份（目录）
```bash
# 直接加密目录 (v0.2.0 新增)
fzj encrypt-dir -i ./important_data -o backup.fzj -p backup_public.pem -s backup_private.pem
```

### 3. 机密文档共享
```bash
# 团队成员间共享加密文档
fzj encrypt -i project.docx -o project.fzj -p team_public.pem -s my_private.pem
```

## 🔒 安全警告

⚠️ **重要提示**:
- 这是一个研究性质的项目，虽然使用了行业标准加密算法
- 生产环境使用前请进行充分的安全评估
- 请妥善保管私钥文件，不要与他人分享
- 建议定期轮换密钥（3-6个月）

### v1.0.3 热修复摘要
- 修复跨平台路径解析：`decrypt` 默认输出名现在可正确处理 `C:\...` 这类 Windows 风格路径
- 在 Linux/macOS 环境中同样只取安全文件名，避免保留绝对路径
- 修复 GitHub CI 用例 `TestSafeDefaultOutputFromHeader/absolute_windows` 失败

### v1.0.2 安全修复摘要
- 修复“签名可剥离绕过”：提供 `--verify-key` 时，密文必须携带有效签名
- 修复默认输出路径风险：`decrypt` 默认输出仅使用头部文件名的 basename
- 加强明文落盘安全：解密输出权限调整为 `0600`
- 降低目录解密明文暴露窗口：`decrypt-dir` 改为随机临时文件并尽早删除
- 补齐 `keymanage -a cache-info` 实现，与文档保持一致

## 🤝 参与贡献

欢迎各种形式的贡献！请先阅读 [贡献指南](docs/CONTRIBUTING.md)。

### 贡献类型
- 🐛 报告 Bug
- 💡 提出新功能
- 📝 改进文档
- 🔧 提交代码
- ✅ 添加测试
- 📊 性能优化

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🔗 项目链接

- **项目主页**: https://codeberg.org/jiangfire/fzj
- **Issue 追踪**: https://codeberg.org/jiangfire/fzj/issues
- **讨论区**: https://codeberg.org/jiangfire/fzj/discussions

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

**当前版本**: v1.0.3
**最后更新**: 2026-03-08
**状态**: ✅ 生产就绪 (持续优化中)
