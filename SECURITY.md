# 安全文档

本文档详细说明 fzjjyz 的安全架构、算法选择、最佳实践和已知限制。

## 📋 目录

- [安全声明](#安全声明)
- [加密架构](#加密架构)
- [算法详解](#算法详解)
- [安全特性](#安全特性)
- [安全最佳实践](#安全最佳实践)
- [已知限制与缓解](#已知限制与缓解)
- [安全审计建议](#安全审计建议)
- [事件响应](#事件响应)
- [合规性说明](#合规性说明)
- [未来改进计划](#未来改进计划)

---

## 安全声明

### ⚠️ 重要警告

**fzjjyz 是一个研究性质的项目**，虽然使用了行业标准的加密算法，但在生产环境使用前请：

1. **充分理解安全风险**
2. **进行独立的安全评估**
3. **考虑咨询安全专家**
4. **定期更新软件版本**

### 适用场景

✅ **推荐使用场景**:
- 个人文件加密
- 小型团队内部文件传输
- 安全备份
- 教育和研究目的
- CTF 竞赛

❌ **不推荐使用场景**:
- 高价值商业机密（未经专业审计）
- 国家级敏感数据
- 实时通信系统
- 缺乏安全维护的环境

---

## 加密架构

### 整体架构图

```
原始文件
    ↓
┌─────────────────────────────────────────┐
│ 1. 密钥封装层 (Key Encapsulation)      │
│    - Kyber768 (后量子)                  │
│    - X25519 ECDH (传统)                 │
│    输出: 32字节共享密钥                 │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│ 2. 数据加密层 (Data Encryption)        │
│    - AES-256-GCM                        │
│    输出: 密文 + 认证标签                │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│ 3. 签名层 (Digital Signature) 可选     │
│    - Dilithium3                         │
│    输出: 数字签名                       │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│ 4. 文件封装层 (File Format)            │
│    - 自定义二进制格式                   │
│    - 包含元数据和验证信息               │
└─────────────────────────────────────────┘
    ↓
加密文件 (.fzj)
```

### 安全边界

```
┌─────────────────────────────────────────────────────────────┐
│                    安全边界                                 │
│                                                             │
│  外部攻击者 ──→  传输/存储  ──→  加密文件  ──→  内部威胁    │
│                                                             │
│  攻击面:                                                   │
│  • 文件窃取        ✓ 已防护 (加密)                          │
│  • 篡改检测        ✓ 已防护 (GCM认证 + 哈希)               │
│  • 量子攻击        ✓ 已防护 (Kyber768)                     │
│  • 密钥泄露        ⚠️ 依赖用户保管                         │
│  • 侧信道攻击      ⚠️ 依赖标准库实现                       │
│  • 实现漏洞        ⚠️ 需要代码审计                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 算法详解

### 1. Kyber768 - 后量子密钥封装

**用途**: 抵抗量子计算机攻击的密钥交换

**参数**:
- 类型: ML-KEM (Module-Lattice-Based Key-Encapsulation Mechanism)
- 安全级别: NIST PQC Level 3 (相当于 AES-192)
- 密钥大小: 公钥 1184 字节，私钥 2400 字节
- 密文大小: 1088 字节
- 共享密钥: 32 字节

**选择理由**:
- ✅ NIST 后量子密码学标准化候选
- ✅ 经过广泛密码学分析
- ✅ 性能优秀
- ✅ Cloudflare CIRCL 实现成熟

**潜在风险**:
- ⚠️ 相对较新，长期安全性仍在研究中
- ⚠️ 可能存在未知的实现漏洞

**缓解措施**:
- 使用经过审计的实现 (Cloudflare CIRCL)
- 保持依赖更新
- 监控密码学社区的最新研究

### 2. X25519 ECDH - 传统密钥交换

**用途**: 提供传统安全层的冗余保护

**参数**:
- 曲线: Curve25519
- 安全级别: ~128 位
- 密钥大小: 32 字节
- 性能: 极快

**选择理由**:
- ✅ 成熟、广泛部署
- ✅ 抗侧信道攻击设计
- ✅ 快速性能
- ✅ 双重保护 (如果 Kyber 被破解，仍有 ECDH 保护)

**潜在风险**:
- ⚠️ 量子计算机可破解 (Shor 算法)

**缓解措施**:
- 与 Kyber 结合使用 (混合模式)
- 未来可轻松替换或移除

### 3. AES-256-GCM - 认证加密

**用途**: 数据机密性和完整性

**参数**:
- 模式: GCM (Galois/Counter Mode)
- 密钥大小: 256 位 (32 字节)
- IV 大小: 12 字节 (96 位)
- 认证标签: 16 字节
- 性能: 硬件加速

**选择理由**:
- ✅ NIST 标准，广泛使用
- ✅ 认证加密 (AEAD)，同时提供加密和完整性
- ✅ 抗填充预言攻击
- ✅ 并行化友好

**潜在风险**:
- ⚠️ IV 重用会导致密钥流重用 (灾难性)
- ⚠️ 实现不当可能导致时序攻击

**缓解措施**:
- 使用加密安全随机数生成器 (crypto/rand)
- 每次加密使用新 IV
- 使用标准库实现 (crypto/aes, crypto/cipher)

### 4. Dilithium3 - 数字签名

**用途**: 文件来源认证和不可否认性

**参数**:
- 类型: ML-DSA (Module-Lattice-Based Digital Signature Algorithm)
- 安全级别: NIST PQC Level 3
- 公钥大小: 1952 字节
- 私钥大小: 4000 字节
- 签名大小: 3293 字节
- 消息哈希: SHA256

**选择理由**:
- ✅ NIST 后量子密码学标准化候选
- ✅ 强安全性保证
- ✅ 签名大小合理
- ✅ Cloudflare CIRCL 实现

**潜在风险**:
- ⚠️ 签名较大，增加文件大小
- ⚠️ 新算法，长期安全性仍在研究

**缓解措施**:
- 可选使用 (不签名的文件更小)
- 保持实现更新

### 5. SHA256 - 完整性校验

**用途**: 文件完整性验证

**参数**:
- 输出: 32 字节哈希值
- 抗碰撞性: 2^128

**选择理由**:
- ✅ 标准且广泛使用
- ✅ 性能优秀
- ✅ 足够的抗碰撞性

---

## 安全特性

### 1. 后量子安全 (Post-Quantum Security)

```
传统加密 (RSA, ECDH) ──→ 量子计算机 ──→ 破解
                          ↓
Kyber768 ──→ 量子计算机 ──→ 仍安全 (预计 2030+ 年)
```

**实现**:
- Kyber768 提供后量子安全
- ECDH 提供传统安全冗余
- 混合模式确保双重保护

**时间线**:
- 2025: 当前状态
- 2030+: 预计量子计算机威胁传统加密
- 2035+: 后量子密码学成为标准

### 2. 前向保密 (Forward Secrecy)

```
每次加密流程:
1. 生成新的临时密钥对
2. 使用接收方公钥封装
3. 加密数据
4. 丢弃临时私钥
```

**优势**:
- 即使长期私钥泄露，历史消息仍安全
- 每次加密使用独立密钥

### 3. 认证加密 (Authenticated Encryption)

```
AES-256-GCM 提供:
├─ 机密性: 数据不可读
├─ 完整性: 检测篡改
└─ 真实性: 确保来源
```

**攻击防护**:
- ✅ 密文篡改检测
- ✅ 填充攻击防护
- ✅ 重放攻击 (通过时间戳)

### 4. 权限最小化

**私钥文件权限**:
- Unix/Linux: 0600 (仅所有者可读写)
- Windows: 设置只读属性

**实现**:
```go
// internal/crypto/keyfile.go
func SetSecurePermissions(path string) error {
    // Unix: chmod 600
    // Windows: 设置文件属性
}
```

### 5. 错误隔离

**设计原则**:
- 不泄露内部细节
- 用户友好错误消息
- 详细日志 (verbose 模式)

**示例**:
```
用户看到: "解密失败: 密钥不匹配"
内部日志: "AES-GCM decryption failed: Authentication failed"
```

---

## 安全最佳实践

### 1. 密钥管理

#### 密钥生成
```bash
# ✅ 正确做法
fzjjyz keygen -d ~/.secure/keys -n production

# ❌ 错误做法
fzjjyz keygen -d ./public_html -n keys  # 目录公开
fzjjyz keygen -d /tmp -n mykey         # 临时目录
```

#### 密钥存储
```bash
# ✅ 推荐方案
# 1. 加密存储
gpg --symmetric --cipher-algo AES256 keys/

# 2. 硬件安全模块 (HSM)
# 3. 密钥管理服务 (KMS)

# ❌ 避免
# 1. 提交到 Git
# 2. 存储在云盘未加密
# 3. 通过邮件发送
```

#### 密钥轮换
```bash
# 建议周期: 3-6 个月
# 或者: 重要事件后 (人员离职、系统迁移)

# 1. 生成新密钥
fzjjyz keygen -d ./keys -n newkey

# 2. 验证新密钥
fzjjyz keymanage -a verify -p keys/newkey_public.pem -s keys/newkey_private.pem

# 3. 重新加密重要文件
for file in important/*.fzj; do
    fzjjyz decrypt -i "$file" -p keys/oldkey_private.pem -o temp.txt
    fzjjyz encrypt -i temp.txt -o "$file" -p keys/newkey_public.pem -s keys/newkey_dilithium_private.pem
done

# 4. 备份并安全删除旧密钥
tar -czf oldkey_backup.tar.gz keys/oldkey*
# 安全删除...
```

### 2. 文件加密最佳实践

#### 加密前
```bash
# ✅ 检查文件完整性
sha256sum sensitive.txt

# ✅ 备份原始文件
cp sensitive.txt sensitive.txt.backup

# ✅ 验证密钥对
fzjjyz keymanage -a verify -p public.pem -s private.pem

# ❌ 不要加密已损坏的文件
# ❌ 不要在多用户系统上存储明文密钥
```

#### 加密过程
```bash
# ✅ 使用签名验证来源
fzjjyz encrypt -i data.txt -o data.fzj \
  -p recipient_public.pem \
  -s my_dilithium_private.pem

# ✅ 详细模式检查
fzjjyz encrypt -i data.txt -o data.fzj \
  -p recipient_public.pem \
  -s my_dilithium_private.pem \
  -v
```

#### 加密后
```bash
# ✅ 验证加密文件
fzjjyz info -i data.fzj

# ✅ 测试解密 (在安全环境)
fzjjyz decrypt -i data.fzj -o test.txt \
  -p my_private.pem \
  -s recipient_public.pem
diff data.txt test.txt

# ✅ 安全删除原始文件 (如果需要)
# Linux/macOS: shred -u sensitive.txt
# Windows: 使用安全删除工具

# ✅ 记录元数据
echo "文件: data.fzj" >> encryption_log.txt
echo "时间: $(date)" >> encryption_log.txt
echo "接收方: recipient@domain.com" >> encryption_log.txt
```

### 3. 环境安全

#### 系统环境
```bash
# ✅ 保持系统更新
# Windows: Windows Update
# Linux: sudo apt update && sudo apt upgrade
# macOS: 系统偏好设置 → 软件更新

# ✅ 使用防火墙
# Linux: sudo ufw enable
# Windows: 启用 Windows Defender 防火墙

# ✅ 防病毒软件
# 保持病毒库更新
```

#### 运行环境
```bash
# ✅ 在受信任的环境中操作
# - 私人网络
# - VPN
# - 物理安全的计算机

# ❌ 避免
# - 公共 Wi-Fi
# - 共享计算机
# - 未打补丁的系统
```

#### 临时文件
```bash
# ✅ 使用安全临时目录
export TMPDIR=/tmp/secure_$$
mkdir -p $TMPDIR
chmod 700 $TMPDIR

# ✅ 操作完成后清理
trap "rm -rf $TMPDIR" EXIT
```

### 4. 传输安全

#### 传输前
```bash
# ✅ 压缩和加密
tar -czf data.tar.gz /path/to/data
fzjjyz encrypt -i data.tar.gz -o data.tar.gz.fzj \
  -p recipient_public.pem \
  -s my_dilithium_private.pem

# ✅ 生成校验和
sha256sum data.tar.gz.fzj > data.tar.gz.fzj.sha256
```

#### 传输中
```bash
# ✅ 使用安全传输协议
# - HTTPS
# - SFTP/SCP
# - VPN

# ❌ 避免
# - HTTP
# - FTP
# - 未加密的邮件
```

#### 传输后
```bash
# ✅ 接收方验证
# 1. 检查校验和
sha256sum -c data.tar.gz.fzj.sha256

# 2. 查看文件信息
fzjjyz info -i data.tar.gz.fzj

# 3. 解密测试
fzjjyz decrypt -i data.tar.gz.fzj -o data.tar.gz \
  -p recipient_private.pem \
  -s sender_public.pem

# 4. 验证完整性
tar -tzf data.tar.gz
```

### 5. 备份策略

#### 3-2-1 备份规则
```
3 份副本:
  ├─ 原始文件
  ├─ 本地备份
  └─ 远程备份

2 种不同介质:
  ├─ 本地磁盘
  └─ 外部存储/云存储

1 份异地备份:
  └─ 物理分离的位置
```

#### 备份流程
```bash
# 1. 创建加密备份
fzjjyz encrypt -i backup.tar.gz -o backup.tar.gz.fzj \
  -p backup_public.pem \
  -s backup_dilithium_private.pem

# 2. 多位置存储
cp backup.tar.gz.fzj /mnt/external/
cp backup.tar.gz.fzj /remote/backup/

# 3. 定期验证
fzjjyz info -i /mnt/external/backup.tar.gz.fzj
fzjjyz info -i /remote/backup/backup.tar.gz.fzj

# 4. 测试恢复
mkdir -p /tmp/restore_test
fzjjyz decrypt -i /mnt/external/backup.tar.gz.fzj \
  -o /tmp/restore_test/backup.tar.gz \
  -p backup_private.pem \
  -s backup_dilithium_public.pem
```

---

## 已知限制与缓解

### 1. 实现成熟度

**限制**:
- 新项目，未经过广泛审计
- 可能存在未知漏洞

**缓解**:
- ✅ 使用经过验证的加密库 (Cloudflare CIRCL)
- ✅ 遵循最佳实践
- ✅ 代码审查和测试
- ⚠️ 建议进行独立安全审计

### 2. 侧信道攻击

**限制**:
- 标准库可能有侧信道漏洞
- 时间攻击可能泄露信息

**缓解**:
- ✅ 使用常数时间操作
- ✅ 依赖标准库的安全实现
- ⚠️ 高安全场景建议使用硬件安全模块

### 3. 密钥管理依赖用户

**限制**:
- 密钥安全完全依赖用户
- 无密钥恢复机制

**缓解**:
- ✅ 提供详细的安全指南
- ✅ 自动设置安全权限
- ✅ 强制备份建议
- ⚠️ 考虑集成密钥管理系统

### 4. 文件格式兼容性

**限制**:
- 自定义二进制格式
- 不同版本间可能不兼容

**缓解**:
- ✅ 明确的版本号
- ✅ 向后兼容设计
- ✅ 详细的格式文档
- ⚠️ 重大变更时提供迁移工具

### 5. 性能与安全权衡

**限制**:
- 后量子算法较慢
- 签名增加文件大小

**缓解**:
- ✅ 可选签名
- ✅ 流式处理支持大文件
- ✅ 性能优化 (基准测试 1MB < 40ms)

---

## 安全审计建议

### 审计清单

#### 1. 密码学审查
- [ ] 验证算法选择是否合适
- [ ] 检查密钥生成随机性
- [ ] 验证 IV 生成和使用
- [ ] 检查密钥派生过程
- [ ] 验证签名验证逻辑

#### 2. 实现审查
- [ ] 检查缓冲区处理
- [ ] 验证错误处理
- [ ] 检查内存管理
- [ ] 验证文件权限设置
- [ ] 检查临时文件清理

#### 3. 协议审查
- [ ] 验证文件格式解析
- [ ] 检查长度验证
- [ ] 验证魔数检查
- [ ] 检查时间戳验证
- [ ] 验证哈希比较 (常数时间)

#### 4. 环境审查
- [ ] 检查依赖版本
- [ ] 验证构建过程
- [ ] 检查部署安全
- [ ] 验证备份策略
- [ ] 检查日志记录

### 推荐审计工具

```bash
# 1. 静态分析
go vet ./...
gosec ./...

# 2. 依赖扫描
govulncheck ./...
npm audit  # 如果有 JS 依赖

# 3. 代码质量
golangci-lint run

# 4. 测试覆盖
go test -cover ./...
go tool cover -html=coverage.out
```

### 第三方审计

**建议**:
- 高价值数据: 聘请专业密码学公司审计
- 开源项目: 社区同行评审
- 学术研究: 发表论文接受学术审查

**审计机构**:
- NCC Group
- Trail of Bits
- Cure53
- 密码学研究机构

---

## 事件响应

### 1. 密钥泄露

**检测**:
- 未授权的文件访问
- 异常的解密尝试
- 系统入侵迹象

**响应流程**:
```bash
# 1. 立即停止使用泄露密钥
# 2. 评估泄露范围
# 3. 生成新密钥对
fzjjyz keygen -d ./keys -n newkey

# 4. 重新加密所有相关文件
# 5. 撤销旧密钥
# 6. 调查泄露原因
# 7. 更新安全措施
```

### 2. 文件篡改

**检测**:
- 解密失败
- 哈希验证失败
- 签名验证失败

**响应流程**:
```bash
# 1. 隔离可疑文件
mv suspicious.fzj quarantine/

# 2. 检查文件信息
fzjjyz info -i quarantine/suspicious.fzj

# 3. 尝试解密 (查看篡改位置)
fzjjyz decrypt -i quarantine/suspicious.fzj -o /dev/null -p key.pem 2>&1

# 4. 从备份恢复
# 5. 检查系统安全性
```

### 3. 实现漏洞发现

**报告流程**:
1. **保密**: 不要公开披露
2. **联系**: security@jiangfire.com
3. **提供**: 详细复现步骤
4. **等待**: 修复发布
5. **验证**: 修复有效性

**响应时间**:
- 高危漏洞: 24 小时内响应
- 中危漏洞: 72 小时内响应
- 低危漏洞: 1 周内响应

---

## 合规性说明

### 加密出口管制

**注意**:
- Kyber, Dilithium: 无出口限制 (公开算法)
- AES: 受美国出口管制，但开源软件通常豁免
- 建议: 了解所在国家/地区的相关法规

### 数据保护法规

**GDPR (欧盟)**:
- ✅ 加密是推荐的技术措施
- ⚠️ 需要记录处理活动
- ⚠️ 数据主体权利 (访问、删除)

**CCPA (加州)**:
- ✅ 加密帮助保护个人信息
- ⚠️ 需要通知数据泄露

**其他法规**:
- 根据所在地区咨询法律专家

---

## 未来改进计划

### 短期 (2026)

#### 1. 硬件安全模块支持
```
目标: 支持 HSM/KMS 存储私钥
实现: PKCS#11 接口
影响: 提高密钥安全性
```

#### 2. 密钥派生函数
```
目标: 支持从密码派生密钥
算法: Argon2id
影响: 更方便的密钥管理
```

#### 3. 多重签名
```
目标: 支持多方签名
场景: 团队文件批准
影响: 增强协作安全性
```

### 中期 (2026-2027)

#### 1. 更换算法
```
目标: 跟随 NIST 标准更新
候选: ML-KEM, ML-DSA (标准化版本)
影响: 长期安全性保证
```

#### 2. 协议升级
```
目标: 支持流式加密
场景: 超大文件 (>10GB)
影响: 内存效率提升
```

#### 3. 审计和认证
```
目标: 第三方安全审计
标准: FIPS 140-3
影响: 生产环境可用性
```

### 长期 (2027+)

#### 1. 量子安全网络
```
目标: 集成量子密钥分发 (QKD)
场景: 最高安全级别
影响: 理论上无条件安全
```

#### 2. 零知识证明
```
目标: 无需解密的文件验证
技术: ZK-SNARKs
影响: 隐私保护增强
```

#### 3. 去中心化存储
```
目标: IPFS/区块链集成
场景: 分布式加密存储
影响: 抗审查和高可用性
```

---

## 总结

### 安全要点

1. **算法选择**: Kyber768 + ECDH + AES-256-GCM + Dilithium3
2. **安全级别**: 后量子安全 (NIST Level 3)
3. **主要优势**: 双重保护、认证加密、来源验证
4. **主要风险**: 密钥管理、实现成熟度、用户操作

### 使用建议

**适合**:
- ✅ 个人敏感文件
- ✅ 小型团队协作
- ✅ 安全备份
- ✅ 教育研究

**不适合**:
- ❌ 高价值商业机密 (未经审计)
- ❌ 国家级敏感数据
- ❌ 缺乏维护的环境

### 安全原则

```
1. 不要信任，始终验证
2. 密钥管理是关键
3. 备份至关重要
4. 保持软件更新
5. 定期安全审计
```

---

## 获取帮助

### 安全报告

**发现安全问题？请立即报告**:
- **邮箱**: security@jiangfire.com
- **PGP**: [下载公钥](https://codeberg.org/jiangfire/fzjjyz/security/pgp)
- **响应时间**: 24 小时内

### 安全咨询

**一般安全问题**:
- 查阅本文档
- 查看 [USAGE.md](USAGE.md) 中的安全提示
- 在项目讨论区提问

### 社区资源

- **安全公告**: 项目 Releases 页面
- **已知问题**: GitHub Issues (安全标签)
- **最佳实践**: 社区 Wiki

---

**版本**: v0.1.0
**最后更新**: 2025-12-21
**维护者**: fzjjyz 安全团队
**状态**: 🟡 审计中