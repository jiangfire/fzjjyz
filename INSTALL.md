# 安装指南

本指南将帮助您在不同平台上安装和配置 fzjjyz。

## 📋 系统要求

### 最低要求
- **Go**: 1.25.4 或更高版本
- **内存**: 256 MB
- **磁盘空间**: 10 MB
- **操作系统**: Windows 10+, Linux, macOS 10.15+

### 推荐配置
- **Go**: 1.26+
- **内存**: 512 MB
- **磁盘空间**: 50 MB（包含测试数据）
- **存储**: SSD（提高文件操作速度）

### 检查 Go 版本
```bash
go version
```

如果未安装 Go 或版本过低，请访问 [Go 官网](https://go.dev/dl/) 下载安装。

---

## 📥 安装方式

### 方式 1: 从源码构建（推荐）

这是最灵活的安装方式，适合开发者和高级用户。

#### 步骤 1: 获取源码

```bash
# 使用 Git 克隆（推荐）
git clone https://codeberg.org/jiangfire/fzjjyz
cd fzjjyz

# 或者下载发布包
# 访问 https://codeberg.org/jiangfire/fzjjyz/releases
# 下载并解压
```

#### 步骤 2: 构建二进制

```bash
# Linux / macOS
go build -o fzjjyz ./cmd/fzjjyz

# Windows
go build -o fzjjyz.exe ./cmd/fzjjyz
```

**构建选项**:
```bash
# 优化构建（减小体积）
go build -ldflags="-s -w" -o fzjjyz ./cmd/fzjjyz

# 调试构建（包含调试信息）
go build -gcflags="all=-N -l" -o fzjjyz_debug ./cmd/fzjjyz

# 跨平台构建
GOOS=linux GOARCH=amd64 go build -o fzjjyz_linux ./cmd/fzjjyz
GOOS=windows GOARCH=amd64 go build -o fzjjyz_windows.exe ./cmd/fzjjyz
GOOS=darwin GOARCH=amd64 go build -o fzjjyz_macos ./cmd/fzjjyz
```

#### 步骤 3: 验证安装

```bash
# Linux / macOS
./fzjjyz version

# Windows
.\fzjjyz.exe version
```

**预期输出**:
```
fzjjyz - 后量子文件加密工具
版本: 0.1.1
应用名称: fzjjyz
描述: 后量子文件加密工具 - 使用 Kyber768 + ECDH + AES-256-GCM + Dilithium3
```

**v0.1.1 新增特性**:
- ✅ 智能密钥缓存（TTL + 大小限制）
- ✅ 并行密钥生成（3x 加速）
- ✅ 性能基准测试
- ✅ 代码重构（消除重复代码）
- ✅ 优化的错误提示

#### 步骤 4: 全局安装（可选）

**Linux / macOS**:
```bash
# 复制到系统路径
sudo cp fzjjyz /usr/local/bin/

# 验证
which fzjjyz
fzjjyz version
```

**Windows (以管理员身份运行)**:
```cmd
# 复制到系统路径
copy fzjjyz.exe C:\Windows\System32\

# 验证
fzjjyz version
```

---

### 方式 2: 使用 Go 安装

如果您已经配置好 Go 环境，可以直接安装。

```bash
# 安装到 GOPATH/bin
go install codeberg.org/jiangfire/fzjjyz/cmd/fzjjyz@latest

# 验证安装
fzjjyz version
```

**注意**: 确保 `$GOPATH/bin` 在您的 `PATH` 环境变量中。

---

### 方式 3: 预编译二进制

适合快速部署，无需编译。

#### 下载地址
访问 [Releases 页面](https://codeberg.org/jiangfire/fzjjyz/releases) 下载对应平台的预编译二进制：

| 平台 | 文件名 | 架构 |
|------|--------|------|
| Windows | `fzjjyz-windows-amd64.exe` | x86-64 |
| Linux | `fzjjyz-linux-amd64` | x86-64 |
| macOS | `fzjjyz-darwin-amd64` | x86-64 |
| macOS (Apple Silicon) | `fzjjyz-darwin-arm64` | ARM64 |

#### 安装步骤

**Linux / macOS**:
```bash
# 1. 下载
wget https://codeberg.org/jiangfire/fzjjyz/releases/download/v0.1.0/fzjjyz-linux-amd64

# 2. 添加执行权限
chmod +x fzjjyz-linux-amd64

# 3. 重命名（可选）
mv fzjjyz-linux-amd64 fzjjyz

# 4. 移动到系统路径（可选）
sudo mv fzjjyz /usr/local/bin/

# 5. 验证
fzjjyz version
```

**Windows**:
```powershell
# 1. 下载（使用浏览器或 PowerShell）
Invoke-WebRequest -Uri "https://codeberg.org/jiangfire/fzjjyz/releases/download/v0.1.0/fzjjyz-windows-amd64.exe" -OutFile "fzjjyz.exe"

# 2. 验证
.\fzjjyz.exe version

# 3. 全局安装（可选，以管理员身份运行）
copy fzjjyz.exe C:\Windows\System32\
```

---

## 🔧 平台特定说明

### Windows

#### PowerShell 示例
```powershell
# 生成密钥
.\fzjjyz.exe keygen -d .\keys -n mykey

# 加密文件
.\fzjjyz.exe encrypt -i secret.txt -o secret.fzj `
  -p .\keys\mykey_public.pem `
  -s .\keys\mykey_dilithium_private.pem

# 解密文件
.\fzjjyz.exe decrypt -i secret.fzj -o recovered.txt `
  -p .\keys\mykey_private.pem `
  -s .\keys\mykey_dilithium_public.pem
```

#### 常见问题

**问题**: PowerShell 执行策略错误
```powershell
# 错误信息: 无法加载脚本，因为在此系统上禁止运行脚本
# 解决方案:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**问题**: Windows Defender 警告
```
这是正常的，因为是新发布的命令行工具。
选择"更多信息" -> "仍要运行"
```

### Linux

#### Bash 示例
```bash
# 添加执行权限
chmod +x fzjjyz

# 添加到 PATH（临时）
export PATH=$PATH:$(pwd)

# 添加到 PATH（永久）
echo 'export PATH=$PATH:/path/to/fzjjyz' >> ~/.bashrc
source ~/.bashrc

# 或者复制到系统路径
sudo cp fzjjyz /usr/local/bin/
```

#### 权限设置
```bash
# 如果遇到权限问题
chmod +x fzjjyz

# 验证
ls -la fzjjyz
# 应该显示: -rwxr-xr-x (755)
```

#### 常见发行版

**Ubuntu / Debian**:
```bash
# 安装 Go (如果需要)
sudo apt update
sudo apt install golang-go

# 构建
go build -o fzjjyz ./cmd/fzjjyz
```

**CentOS / RHEL**:
```bash
# 安装 Go (如果需要)
sudo yum install golang

# 构建
go build -o fzjjyz ./cmd/fzjjyz
```

### macOS

#### 基本使用
```bash
# 添加执行权限
chmod +x fzjjyz

# 如果遇到"无法打开"错误
xattr -d com.apple.quarantine fzjjyz

# 或者在系统偏好设置 -> 安全性与隐私中允许运行
```

#### Apple Silicon (M1/M2)
```bash
# 如果下载的是 Intel 版本，需要 Rosetta 2
# 或者下载 arm64 版本

# 检查架构
file fzjjyz
# 应该显示: Mach-O 64-bit executable arm64
```

#### 常见问题

**问题**: "无法打开，因为无法验证开发者"
```
解决方案 1:
- 系统偏好设置 -> 安全性与隐私
- 点击"仍要打开"

解决方案 2:
xattr -d com.apple.quarantine fzjjyz
```

---

## ✅ 验证安装

运行以下命令验证安装是否成功：

### 1. 版本检查
```bash
fzjjyz version
```
**预期输出**: 显示版本号和描述

### 2. 帮助信息
```bash
fzjjyz --help
```
**预期输出**: 显示所有可用命令

### 3. 快速测试
```bash
# 创建临时目录
mkdir -p /tmp/fzjjyz_test
cd /tmp/fzjjyz_test

# 生成测试密钥
fzjjyz keygen -d ./keys -n test

# 创建测试文件
echo "Hello, fzjjyz!" > hello.txt

# 加密
fzjjyz encrypt -i hello.txt -o hello.fzj \
  -p keys/test_public.pem \
  -s keys/test_dilithium_private.pem

# 解密
fzjjyz decrypt -i hello.fzj -o hello_decrypted.txt \
  -p keys/test_private.pem \
  -s keys/test_dilithium_public.pem

# 验证
cat hello_decrypted.txt
# 应该显示: Hello, fzjjyz!

# 清理
cd ~ && rm -rf /tmp/fzjjyz_test
```

---

## 🔍 依赖说明

### 自动管理

Go Modules 会自动管理所有依赖。主要依赖包括：

```bash
# 查看依赖
go list -m all
```

**核心依赖**:
- `github.com/cloudflare/circl v1.6.1` - 后量子密码学库
- `github.com/spf13/cobra v1.10.2` - CLI 框架
- `github.com/schollz/progressbar/v3 v3.18.0` - 进度条

### 手动更新依赖
```bash
# 更新所有依赖
go get -u ./...

# 清理未使用的依赖
go mod tidy

# 验证依赖
go mod verify
```

---

## ❓ 常见问题

### Q: 构建失败，提示 "module not found"
**A**:
```bash
# 确保使用 Go 1.25.4+
go version

# 下载依赖
go mod download

# 清理并重新构建
go clean -cache
go build ./cmd/fzjjyz
```

### Q: 运行时提示 "permission denied"
**A**:
```bash
# Linux/macOS: 添加执行权限
chmod +x fzjjyz

# Windows: 以管理员身份运行
# 或检查文件是否被锁定
```

### Q: Windows 防火墙警告
**A**: 这是正常的，因为是命令行工具。选择"允许访问"。

### Q: 如何卸载？
**A**:

**Linux/macOS**:
```bash
# 删除二进制
sudo rm /usr/local/bin/fzjjyz

# 删除配置文件（如果存在）
rm -rf ~/.config/fzjjyz/
```

**Windows**:
```cmd
# 删除二进制
del C:\Windows\System32\fzjjyz.exe

# 删除工作目录（如果存在）
rmdir /s C:\Users\YourUsername\fzjjyz
```

### Q: 如何更新到最新版本？
**A**:
```bash
# 方式 1: 重新构建
cd fzjjyz
git pull origin main
go build -o fzjjyz ./cmd/fzjjyz

# 方式 2: 重新下载
# 访问 Releases 页面下载最新版本
```

### Q: 遇到 "out of memory" 错误
**A**:
- 确保系统至少有 256MB 可用内存
- 对于超大文件（>1GB），建议使用 64 位系统
- 关闭其他占用内存的程序

### Q: 加密/解密速度慢
**A**:
- 使用 SSD 存储
- 确保 Go 版本 >= 1.25
- 检查 CPU 使用率
- 对于大文件，工具会自动优化

---

## 📚 下一步

安装完成后，您可以：

1. **阅读使用文档**: 查看 [USAGE.md](USAGE.md) 学习如何使用所有命令
2. **快速开始**: 按照 README.md 的快速开始指南操作
3. **了解安全**: 阅读 [SECURITY.md](SECURITY.md) 了解安全最佳实践
4. **开始开发**: 如果您想贡献代码，查看 [DEVELOPMENT.md](DEVELOPMENT.md)

---

## 🆘 获取帮助

如果在安装过程中遇到问题：

1. 查看本文档的常见问题部分
2. 阅读 [SECURITY.md](SECURITY.md) 了解已知限制
3. 在项目 [Issues](https://codeberg.org/jiangfire/fzjjyz/issues) 页面搜索或创建 Issue
4. 加入项目讨论区寻求帮助

---

**版本**: v0.1.1
**最后更新**: 2025-12-26
**维护者**: fzjjyz 开发团队

---

## 🆕 v0.1.1 新增验证步骤

安装完成后，验证 v0.1.1 新特性：

### 1. 验证缓存系统

```bash
# 生成测试密钥
fzjjyz keygen -d ./test_keys -n test

# 创建测试文件
echo "测试缓存性能" > test.txt

# 第一次加密（加载密钥）
fzjjyz encrypt -i test.txt -o test.fzj \
  -p test_keys/test_public.pem \
  -s test_keys/test_dilithium_private.pem

# 第二次加密（使用缓存，应显著更快）
fzjjyz encrypt -i test.txt -o test2.fzj \
  -p test_keys/test_public.pem \
  -s test_keys/test_dilithium_private.pem --verbose
```

**预期结果**: 第二次操作明显更快（<1ms 密钥加载）

### 2. 验证并行密钥生成

```bash
# 生成密钥并观察速度
time fzjjyz keygen -d ./test_keys2 -n test2
```

**预期结果**: 约 150ms（v0.1.0 约 450ms）

### 3. 运行性能基准测试

```bash
# 运行完整基准测试
go test -bench=. -benchmem ./internal/crypto/
```

**预期输出**:
```
BenchmarkEncryptFile/1MB-8          100    ~1.9ms    ~3.2MB/op    29 allocs/op
BenchmarkDecryptFile/1MB-8          100    ~1.6ms    ~2.1MB/op    36 allocs/op
BenchmarkKeyGeneration/Parallel-8   10000   ~117ms    ~90KB/op     26 allocs/op
BenchmarkCachePerformance/Cached-8  10000   ~76ns     192 B/op     2 allocs/op
```

### 4. 验证错误提示优化

```bash
# 测试错误提示（使用错误的密钥路径）
fzjjyz encrypt -i test.txt -o test.fzj \
  -p wrong_path.pem \
  -s test_keys/test_dilithium_private.pem
```

**预期结果**: 详细的错误信息和解决方案建议

### 5. 清理测试文件

```bash
# 清理测试数据
rm -rf test_keys test_keys2 test.txt test*.fzj
```

---

## 📋 版本对比

| 特性 | v0.1.0 | v0.1.1 | 改进 |
|------|--------|--------|------|
| 密钥生成 | ~450ms | ~150ms | 3x 加速 |
| 密钥加载 | ~1-2ms | <1μs | 1000x+ 加速 |
| 代码重复 | ~600行 | 0行 | 100% 消除 |
| 错误提示 | 基础 | 详细+解决方案 | 用户友好 |
| 性能测试 | 无 | 完整基准套件 | 新增 |

**升级建议**: 所有用户建议升级到 v0.1.1，获得性能提升和安全增强。