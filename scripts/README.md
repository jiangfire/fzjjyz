# 发布自动化脚本

这些脚本帮助你自动化 fzjjyz 项目的发布流程。

## 目录结构

```
scripts/
├── release.sh          # Linux/macOS 发布脚本
├── release.bat         # Windows 发布脚本
└── README.md          # 本文件
```

## 快速开始

### Linux/macOS

```bash
# 1. 添加执行权限
chmod +x scripts/release.sh

# 2. 运行发布
./scripts/release.sh v0.1.1

# 3. 试运行
./scripts/release.sh v0.1.1 --dry-run
```

### Windows

```cmd
# 1. 运行发布
scripts\release.bat v0.1.1

# 2. 试运行
scripts\release.bat v0.1.1 --dry-run
```

## 功能特性

### 自动化流程

1. **环境检查**
   - 验证当前分支
   - 检查未提交的更改
   - 验证版本号格式

2. **质量保证**
   - 运行完整测试套件
   - 自动更新版本号
   - 跨平台构建

3. **发布准备**
   - 生成 SHA256 校验和
   - 创建 Git 标签
   - 推送到远程仓库

4. **发布说明**
   - 自动生成基础发布说明
   - 支持使用现有 RELEASE_NOTES.md

### 选项说明

| 选项 | 说明 | 使用场景 |
|------|------|----------|
| `--skip-test` | 跳过测试 | 快速测试发布流程 |
| `--skip-build` | 跳过构建 | 只更新版本号和标签 |
| `--dry-run` | 试运行 | 验证流程，不执行实际操作 |

## GitHub Actions

项目包含三个 GitHub Actions workflow：

### 1. Test and Build (test.yml)

**触发条件**:
- 推送到 main/master 分支
- Pull Request
- 手动触发

**功能**:
- 运行测试套件
- 代码格式检查
- 跨平台构建验证
- 安全审计 (govulncheck, gosec)

### 2. Release (release.yml)

**触发条件**:
- 推送版本标签 (v*)
- 手动触发

**功能**:
- 自动测试
- 多平台构建 (Linux/Windows/macOS)
- 生成校验和
- 创建 GitHub Release
- 上传所有二进制文件
- 自动清理临时工件

### 3. Quick Release (quick-release.yml)

**触发条件**:
- 手动触发 (需要输入版本号)

**功能**:
- 快速发布流程
- 自动创建标签
- 生成发布说明
- 上传二进制文件

## 使用示例

### 场景 1: 完整自动化发布

```bash
# 1. 创建版本标签
git tag -a v0.1.1 -m "Release v0.1.1"

# 2. 推送标签
git push origin v0.1.1

# 3. GitHub Actions 自动执行
# - 运行测试
# - 构建所有平台
# - 创建 Release
# - 上传文件
```

### 场景 2: 手动触发快速发布

在 GitHub 页面:
1. 进入 Actions 标签页
2. 选择 "Quick Release"
3. 点击 "Run workflow"
4. 输入版本号 (如: v0.1.1)
5. 点击 "Run workflow"

### 场景 3: 本地脚本发布

```bash
# Linux/macOS
./scripts/release.sh v0.1.1

# Windows
scripts\release.bat v0.1.1
```

## 发布清单

### 发布前检查

- [ ] 所有测试通过
- [ ] 代码已格式化
- [ ] 文档已更新
- [ ] CHANGELOG.md 已更新
- [ ] 版本号已确定

### 发布后验证

- [ ] GitHub Release 已创建
- [ ] 所有二进制文件已上传
- [ ] 校验和文件已上传
- [ ] 发布说明已填写
- [ ] Git 标签已推送

## 文件说明

### 生成的文件

发布过程中会生成以下文件：

```
release/
└── v0.1.1/
    ├── fzjjyz_linux_amd64
    ├── fzjjyz_windows_amd64.exe
    ├── fzjjyz_darwin_amd64
    ├── fzjjyz_darwin_arm64
    ├── checksums.txt
    └── release_notes.md (可选)
```

### GitHub Release 内容

推荐上传到 GitHub Release 的文件：
- 所有平台的二进制文件
- checksums.txt
- RELEASE_NOTES.md (如果存在)

## 故障排除

### 常见问题

**Q: 脚本执行失败，提示权限不足**
```bash
chmod +x scripts/release.sh
```

**Q: Git 标签已存在**
```bash
# 删除旧标签
git tag -d v0.1.1
git push origin --delete v0.1.1

# 重新创建
./scripts/release.sh v0.1.1
```

**Q: 构建失败**
- 确保 Go 环境正确安装
- 检查依赖: `go mod tidy`
- 清理缓存: `go clean -cache`

**Q: GitHub Actions 失败**
- 查看 Actions 日志
- 检查 Secrets 配置
- 验证权限设置

## 最佳实践

### 版本号管理

遵循 [语义化版本](https://semver.org/lang/zh-CN/):

- `0.1.0` → `0.1.1`: Bug 修复
- `0.1.0` → `0.2.0`: 新功能
- `0.1.0` → `1.0.0`: 重大变更

### 发布频率

- **小版本**: 修复 bug 时发布
- **中版本**: 新功能完成时发布
- **大版本**: 稳定性和兼容性保证后发布

### 发布说明

保持发布说明的完整性：
- 清晰的变更描述
- 升级指南
- 已知问题
- 安全说明

## 高级用法

### 自定义发布流程

编辑 `release.sh` 或 `release.bat` 来：
- 添加自定义构建步骤
- 集成其他工具
- 自定义发布说明格式

### CI/CD 集成

可以与其他工具集成：
- Docker 镜像构建
- 自动部署
- 通知系统 (Slack/Email)

## 相关文档

- [CHANGELOG.md](../CHANGELOG.md) - 版本变更记录
- [RELEASE_CHECKLIST.md](../RELEASE_CHECKLIST.md) - 发布清单
- [RELEASE_NOTES.md](../RELEASE_NOTES.md) - 发布说明模板
- [CONTRIBUTING.md](../CONTRIBUTING.md) - 贡献指南

## 支持

如有问题或建议：
- 提交 Issue
- 创建 Pull Request
- 联系维护者

---

**版本**: 1.0
**最后更新**: 2025-12-27
**维护者**: fzjjyz 开发团队