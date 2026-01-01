# 发布指南

本文档说明如何使用 Makefile 和 GitHub Actions 进行项目发布。

## 📋 目录

- [本地发布](#本地发布)
- [GitHub Actions 发布](#github-actions-发布)
- [发布流程说明](#发布流程说明)
- [故障排除](#故障排除)

---

## 本地发布

### 前置要求

- Go 1.25+
- Make (Linux/macOS) 或 Git Bash (Windows)
- golangci-lint (可选，用于代码检查)
- git (用于版本管理)

### 发布步骤

#### 1. 完整发布流程

```bash
# 执行完整发布（清理 + 测试 + 检查 + 跨平台构建 + 校验和）
make release

# 或者使用本地发布（包含打包）
make local-release
```

#### 2. 分步发布

```bash
# 1. 清理旧产物
make clean

# 2. 运行测试
make test

# 3. 代码检查
make lint

# 4. 跨平台构建
make cross-build

# 5. 生成校验和
make checksum

# 6. (可选) 打包
make package
```

#### 3. CI 发布（跳过 lint）

```bash
make ci-release
```

### 常用命令

```bash
# 查看帮助
make help

# 构建当前平台版本
make build

# 构建开发版本（带调试）
make build-dev

# 运行测试并生成覆盖率报告
make test-cover

# 显示版本信息
make version

# 整理依赖
make tidy
```

---

## GitHub Actions 发布

### 触发方式

#### 1. 自动触发（推荐）

推送到标签时自动触发：

```bash
# 创建标签
git tag v0.1.0

# 推送标签
git push origin v0.1.0
```

#### 2. 手动触发

在 GitHub 网页端：
1. 进入 Actions 标签页
2. 选择 "Release" workflow
3. 点击 "Run workflow"
4. 输入版本号（可选）
5. 点击 "Run workflow"

### GitHub Actions 工作流

#### `test.yml` - 测试和构建
- **触发**: push 到 main/master 分支，PR，手动触发
- **功能**:
  - 运行单元测试
  - 代码格式检查
  - 跨平台构建测试
  - 功能测试
  - 安全审计

#### `release.yml` - 正式发布
- **触发**: push 到 v* 标签，手动触发
- **功能**:
  - 完整测试和 lint
  - 跨平台构建 (Linux/Windows/macOS, amd64/arm64)
  - 生成校验和
  - 创建 GitHub Release
  - 上传所有二进制文件

#### `quick-release.yml` - 快速发布
- **触发**: 手动触发
- **功能**:
  - 快速构建和发布
  - 自动创建标签
  - 适合紧急发布

---

## 发布流程说明

### 完整发布流程

```
1. 代码提交到主分支
   ↓
2. 创建版本标签 (vX.Y.Z)
   ↓
3. 推送标签到 GitHub
   ↓
4. GitHub Actions 自动触发
   ↓
5. 质量检查 (测试 + Lint)
   ↓
6. 跨平台构建 (6个平台)
   ↓
7. 生成校验和
   ↓
8. 创建 Release 页面
   ↓
9. 上传所有文件
   ↓
10. 发布完成 ✅
```

### 产物说明

发布完成后，GitHub Release 包含以下文件：

| 文件 | 说明 |
|------|------|
| `fzj-linux-amd64` | Linux 64位 |
| `fzj-linux-arm64` | Linux ARM64 |
| `fzj-darwin-amd64` | macOS Intel |
| `fzj-darwin-arm64` | macOS Apple Silicon |
| `fzj-windows-amd64.exe` | Windows 64位 |
| `fzj-windows-arm64.exe` | Windows ARM64 |
| `checksums.sha256` | SHA256 校验和 |

---

## 故障排除

### 常见问题

#### 1. Make 命令在 Windows 上不可用

**解决方案**: 使用 Git Bash 或 WSL

```bash
# Git Bash (安装 Git 时自带)
make release

# 或者直接使用 Go 命令
go build -ldflags "-s -w" -o dist/fzj ./cmd/fzj
```

#### 2. 跨平台构建失败

**检查点**:
- 确保系统支持目标平台
- 检查 Go 版本是否支持目标平台
- 确保有足够的磁盘空间

```bash
# 检查可用平台
go tool dist list

# 单独测试某个平台
GOOS=linux GOARCH=amd64 go build -o /dev/null ./cmd/fzj
```

#### 3. GitHub Actions 构建失败

**常见原因**:
- 依赖安装失败
- 测试失败
- 代码检查失败

**解决方案**:
1. 查看 Actions 日志
2. 本地运行相同命令测试
3. 修复问题后重新推送标签

```bash
# 本地重现 CI 流程
make ci
```

#### 4. 校验和验证失败

**验证方法**:

```bash
# Linux/macOS
sha256sum -c checksums.sha256

# Windows (PowerShell)
Get-FileHash -Algorithm SHA256 fzj-windows-amd64.exe

# 或者使用校验和文件
sha256sum fzj-windows-amd64.exe
# 然后与 checksums.sha256 中的值对比
```

### 调试 CI 问题

#### 查看详细日志

```bash
# 在 GitHub Actions 中启用详细输出
# 修改 workflow 文件，添加：
- name: Debug info
  run: |
    set -x
    echo "Go version: $(go version)"
    echo "Working dir: $(pwd)"
    ls -la
```

#### 手动触发测试

```bash
# 推送到测试分支
git checkout -b test-release
git add .
git commit -m "test: release workflow"
git push origin test-release

# 在 Actions 中手动触发，选择 test-release 分支
```

---

## 最佳实践

### 版本号规范

使用 [语义化版本](https://semver.org/lang/zh-CN/):

- `v1.0.0` - 正式版本
- `v0.1.0` - 开发版本
- `v1.0.1` - 修复版本
- `v1.1.0` - 功能更新

### 发布前检查清单

- [ ] 所有测试通过
- [ ] 代码检查通过
- [ ] 文档已更新
- [ ] CHANGELOG.md 已更新（如果有）
- [ ] 版本号已更新
- [ ] 本地构建测试通过

### 安全建议

1. **验证二进制**: 始终验证 SHA256 校验和
2. **来源验证**: 确保从官方仓库下载
3. **权限控制**: GitHub token 最小权限原则

---

## 相关文件

- `Makefile` - 构建和发布脚本
- `.github/workflows/release.yml` - 正式发布流程
- `.github/workflows/test.yml` - 测试和质量检查
- `.github/workflows/quick-release.yml` - 快速发布
- `docs/RELEASE_GUIDE.md` - 本指南

---

## 联系与支持

如有问题，请：
1. 查看 GitHub Issues
2. 阅读项目文档
3. 提交 Issue 或 PR