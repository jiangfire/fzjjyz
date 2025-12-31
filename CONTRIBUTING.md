# 贡献指南

欢迎贡献 fzjjyz 项目！本文档将帮助您了解如何参与项目开发、提交代码和改进文档。

## 📋 目录

- [快速开始](#快速开始)
- [贡献类型](#贡献类型)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [代码审查](#代码审查)
- [社区准则](#社区准则)
- [获得帮助](#获得帮助)

---

## 快速开始

### 1. 选择贡献方式

**首次贡献？** 从这些开始：
- 📝 改进文档
- 🐛 报告 Bug
- 💡 提出新功能建议
- ✅ 添加测试用例

**有开发经验？** 直接参与代码：
- 🔧 修复 Bug
- 🚀 实现新特性
- ⚡ 性能优化
- 🔒 安全改进

### 2. 环境准备

```bash
# 1. Fork 项目
# 访问 https://codeberg.org/jiangfire/fzjjyz
# 点击 Fork 按钮

# 2. 克隆你的 Fork
git clone https://codeberg.org/your-username/fzjjyz
cd fzjjyz

# 3. 添加上游仓库
git remote add upstream https://codeberg.org/jiangfire/fzjjyz

# 4. 安装依赖
go mod download
go mod tidy

# 5. 验证构建
go build -o fzjjyz ./cmd/fzjjyz
./fzjjyz version
```

### 3. 你的第一个贡献

```bash
# 1. 创建分支
git checkout -b docs/fix-typo

# 2. 修改文件
# 编辑 README.md 或其他文档

# 3. 提交
git add .
git commit -m "docs: 修复 README 中的拼写错误"

# 4. 推送
git push origin docs/fix-typo

# 5. 创建 Pull Request
# 在 Codeberg 上点击 "Create Pull Request"
```

---

## 贡献类型

### 1. 文档改进 📝

**适合**: 所有贡献者，特别是新手

**示例**:
- 修复拼写错误和语法
- 补充使用示例
- 翻译文档
- 改进 README 结构
- 添加 FAQ

**提交规范**:
```
docs: 修复 USAGE.md 中的命令示例
docs: 添加密钥轮换的最佳实践
docs: 翻译 SECURITY.md 为中文
```

### 2. Bug 报告 🐛

**适合**: 所有用户

**好的 Bug 报告包含**:
- 清晰的问题描述
- 复现步骤
- 预期行为 vs 实际行为
- 环境信息 (OS, Go 版本)
- 错误日志或截图

**模板**:
```markdown
## 问题描述
[清晰描述问题]

## 复现步骤
1. [步骤 1]
2. [步骤 2]
3. [步骤 3]

## 环境
- OS: [例如: Windows 11]
- Go 版本: [例如: 1.26.0]
- fzjjyz 版本: [例如: v0.1.0]

## 预期行为
[期望的结果]

## 实际行为
[实际的结果]

## 错误日志
[如果有]
```

### 3. 功能建议 💡

**适合**: 有创新想法的贡献者

**建议内容**:
- 解决的问题
- 实现方案
- 使用场景
- 潜在影响

**模板**:
```markdown
## 功能描述
[功能名称]

## 解决的问题
[当前的痛点]

## 实现方案
[技术细节]

## 使用场景
[何时使用]

## 备注
[其他信息]
```

### 4. 代码贡献 🔧

**适合**: 有 Go 开发经验

**代码类型**:
- Bug 修复
- 新特性
- 性能优化
- 测试覆盖
- 重构改进

**要求**:
- 通过所有测试
- 遵循代码规范
- 添加相关测试
- 更新文档

### 5. 测试贡献 ✅

**适合**: 注重质量的贡献者

**测试类型**:
- 单元测试
- 集成测试
- 性能测试
- 边界测试
- 错误处理测试

**目标**: 保持 95%+ 覆盖率

---

## 开发流程

### 1. 分支策略

```
main (稳定分支)
    ↑
develop (开发分支)
    ↑
feature/your-feature (特性分支)
    ↑
hotfix/urgent-fix (紧急修复)
```

**分支命名**:
- `feature/description` - 新特性
- `bugfix/description` - Bug 修复
- `docs/description` - 文档改进
- `hotfix/description` - 紧急修复
- `refactor/description` - 重构

### 2. 完整工作流

```bash
# 1. 同步最新代码
git checkout main
git pull upstream main

# 2. 创建特性分支
git checkout -b feature/your-feature

# 3. 进行开发
# 编辑代码...

# 4. 运行测试
go test ./...
go test -cover ./...

# 5. 格式化代码
go fmt ./...
go vet ./...

# 6. 提交代码
git add .
git commit -m "feat: 添加你的特性"

# 7. 推送分支
git push origin feature/your-feature

# 8. 创建 Pull Request
# 在 Codeberg 上创建 PR

# 9. 等待审查
# 根据反馈修改...

# 10. 合并后清理
git checkout main
git branch -d feature/your-feature
git pull upstream main
```

### 3. Pull Request 流程

**PR 标题格式**:
```
类型(范围): 简短描述

例如:
feat(crypto): 添加 Kyber768 密钥生成
fix(cli): 修复 Windows 权限问题
docs: 更新安装指南
```

**PR 描述模板**:
```markdown
## 描述
[清晰描述变更]

## 类型
- [ ] Bug 修复
- [ ] 新特性
- [ ] 文档改进
- [ ] 性能优化
- [ ] 重构
- [ ] 其他

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 添加了测试
- [ ] 更新了文档
- [ ] 所有测试通过
- [ ] 无编译错误

## 测试
[描述如何测试]

## 截图/日志
[如果有]
```

---

## 代码规范

### 1. Go 代码风格

```go
// ✅ 推荐
package main

import (
    "fmt"
    "os"
)

// 驼峰命名，首字母大写导出
type HybridPublicKey struct {
    Kyber *kyber768.PublicKey
    ECDH  *ecdh.PublicKey
}

// 明确的错误处理
func Example() error {
    data, err := os.ReadFile("file.txt")
    if err != nil {
        return fmt.Errorf("读取文件失败: %w", err)
    }

    if len(data) == 0 {
        return errors.New("文件为空")
    }

    return nil
}

// 充分的注释
// GenerateKeys 生成混合密钥对
// 返回: 私钥、公钥、错误
func GenerateKeys() (*HybridPrivateKey, *HybridPublicKey, error) {
    // 实现...
}
```

### 2. 格式化工具

```bash
# 自动格式化
go fmt ./...

# 静态分析
go vet ./...

# 代码检查 (如果安装)
golangci-lint run

# 导入排序
goimports -w .
```

**⚠️ 编辑代码时的注意事项**:
- 使用 Edit 工具时，确保 `old_string` 完全匹配文件中的内容（包括特殊字符）
- 编辑后立即验证结果，避免转义错误传播

### 3. 测试规范

```go
// internal/crypto/example_test.go
package crypto

import (
    "testing"
)

// 表驱动测试
func TestExample(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        {
            name:    "empty input",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ExampleFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ExampleFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("ExampleFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}

// 基准测试
func BenchmarkExample(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ExampleFunction("test")
    }
}
```

### 4. 提交信息规范

**格式**:
```
类型(范围): 简短描述 [不超过 50 字符]

详细描述 [可选，72 字符/行]

[空行]

[可选: 脚注、Breaking Changes 等]
```

**类型**:
- `feat`: 新特性
- `fix`: Bug 修复
- `docs`: 文档
- `style`: 代码格式 (不影响功能)
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试
- `chore`: 构建/工具

**示例**:
```
feat(crypto): 添加 Dilithium3 签名支持

- 实现签名生成和验证
- 集成到加密流程
- 添加相关测试

Closes #123
```

```
fix(cli): 修复 Windows 权限设置问题

Windows 不支持 chmod 600，使用文件属性替代

相关: #45
```

---

## 代码审查

### 1. 审查清单

**提交前自检**:
- [ ] 代码可读性
- [ ] 测试覆盖率 > 90%
- [ ] 文档已更新
- [ ] 无编译警告
- [ ] 通过 lint 检查
- [ ] 性能考虑
- [ ] 安全考虑

**审查者检查**:
- [ ] 功能正确性
- [ ] 边界情况处理
- [ ] 错误处理完整
- [ ] 代码风格一致
- [ ] 测试充分
- [ ] 文档准确
- [ ] 无安全漏洞

### 2. 审查反馈处理

```bash
# 1. 获取 PR 分支
git fetch upstream
git checkout -b pr-branch upstream/pr-branch

# 2. 应用反馈修改
# 编辑代码...

# 3. 提交修改
git add .
git commit -m "fix: 根据审查反馈修改"

# 4. 推送更新
git push origin pr-branch
```

### 3. 常见审查意见

**性能**:
```
建议: 使用缓冲区减少内存分配
建议: 考虑并发处理
```

**安全**:
```
建议: 添加输入验证
建议: 使用常数时间比较
```

**代码质量**:
```
建议: 提取重复代码为函数
建议: 添加更多注释
建议: 简化复杂逻辑
```

---

## 社区准则

### 1. 行为准则

我们遵循 [Contributor Covenant 2.0](https://www.contributor-covenant.org/version/2/0/code_of_conduct/)。

**核心原则**:
- ✅ 尊重他人，包容友好
- ✅ 建设性反馈
- ✅ 专注技术讨论
- ✅ 欢迎新手和多样性

**禁止行为**:
- ❌ 侮辱、骚扰、歧视
- ❌ 恶意代码或破坏
- ❌ 个人攻击
- ❌ 不适当内容

### 2. 沟通规范

**Issue 讨论**:
- 清晰描述问题
- 提供复现步骤
- 保持礼貌
- 及时回应

**PR 讨论**:
- 专注代码质量
- 提供建设性意见
- 解释修改原因
- 感谢贡献者

**社区交流**:
- 使用中性语言
- 避免情绪化
- 寻求共识
- 尊重不同观点

### 3. 决策流程

**小问题**:
- 维护者直接决定
- 快速推进

**大问题**:
- 社区讨论
- 收集反馈
- 维护者决策
- 公开结果

**争议解决**:
- 冷静讨论
- 寻求妥协
- 必要时投票
- 维护者仲裁

---

## 获得帮助

### 1. 文档资源

**必读文档**:
- [README.md](README.md) - 项目概览
- [INSTALL.md](INSTALL.md) - 安装指南
- [USAGE.md](USAGE.md) - 使用文档
- [DEVELOPMENT.md](DEVELOPMENT.md) - 开发指南
- [SECURITY.md](SECURITY.md) - 安全说明

### 2. 问题求助

**遇到问题？按顺序尝试**:

1. **搜索文档**
   ```bash
   grep -r "你的问题" *.md
   ```

2. **搜索 Issue**
   - 在 Codeberg Issues 搜索
   - 查看已关闭的 Issue

3. **社区讨论**
   - 创建 Issue (标记为 "question")
   - 在讨论区提问

4. **直接联系**
   - 邮件: contact@jiangfire.com
   - 安全问题: security@jiangfire.com

### 3. 学习资源

**Go 语言**:
- [Go 官方文档](https://go.dev/doc/)
- [Go 有效编程](https://go.dev/doc/effective_go)
- [Go 语言规范](https://go.dev/ref/spec)

**密码学**:
- [NIST PQC](https://csrc.nist.gov/projects/post-quantum-cryptography)
- [Cloudflare CIRCL](https://github.com/cloudflare/circl)
- [AES-GCM 说明](https://en.wikipedia.org/wiki/Galois/Counter_Mode)

**工具**:
- [Git 教程](https://git-scm.com/book/zh/v2)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

---

## 贡献者列表

### 感谢所有贡献者！

**核心维护者**:
- @jiangfire - 项目发起人

**贡献者**:
- [添加你的名字...]

### 如何加入

**成为贡献者**:
1. 提交 3+ 个被接受的 PR
2. 积极参与社区讨论
3. 遵守行为准则

**成为维护者**:
1. 深入理解项目架构
2. 持续贡献高质量代码
3. 帮助审查他人 PR
4. 获得现有维护者认可

---

## 快速参考

### 常用命令

```bash
# 开发
go build ./cmd/fzjjyz
go test ./...
go fmt ./...
go vet ./...

# 测试
go test -cover ./...
go test -v -run TestEncrypt ./cmd/fzjjyz/

# 构建发布
go build -ldflags="-s -w" -o fzjjyz ./cmd/fzjjyz

# 跨平台
GOOS=windows GOARCH=amd64 go build -o fzjjyz.exe ./cmd/fzjjyz
```

### 提交模板

```bash
# 配置 Git 模板
git config --global commit.template ~/.gitmessage

# ~/.gitmessage 内容:
# 类型(范围): 简短描述
#
# 详细描述...
#
# Breaking Changes: [如果有]
# Related Issues: #123, #456
```

### PR 检查清单

```markdown
## PR 检查清单

### 基础要求
- [ ] 遵循代码规范
- [ ] 通过所有测试
- [ ] 无编译错误
- [ ] 代码格式化

### 功能相关
- [ ] 功能完整实现
- [ ] 边界情况处理
- [ ] 错误处理完善
- [ ] 性能考虑

### 测试相关
- [ ] 单元测试覆盖
- [ ] 集成测试通过
- [ ] 边界测试存在
- [ ] 错误测试完整

### 文档相关
- [ ] README 更新 (如果需要)
- [ ] 代码有注释
- [ ] API 文档完整
- [ ] 示例代码正确

### 其他
- [ ] 无安全漏洞
- [ ] 无性能退化
- [ ] 向后兼容
- [ ] PR 描述清晰
```

---

## 总结

### 贡献流程图

```
1. 选择贡献类型
   ↓
2. 准备环境
   ↓
3. 创建分支
   ↓
4. 进行开发
   ↓
5. 测试验证
   ↓
6. 提交 PR
   ↓
7. 等待审查
   ↓
8. 根据反馈修改
   ↓
9. 合并完成
   ↓
10. 庆祝！🎉
```

### 成功秘诀

1. **从小开始**: 从文档改进开始
2. **保持沟通**: 遇到问题及时提问
3. **学习规范**: 遵循项目代码风格
4. **充分测试**: 测试是质量保证
5. **耐心等待**: 审查需要时间
6. **持续学习**: 从反馈中成长

### 下一步

**现在就开始**:
1. 📖 阅读 [DEVELOPMENT.md](DEVELOPMENT.md) 搭建环境
2. 🔍 浏览 [Issues](https://codeberg.org/jiangfire/fzjjyz/issues) 找感兴趣的任务
3. 💬 加入讨论，提出你的想法
4. 🚀 提交你的第一个 PR！

---

**我们期待你的贡献！** 🎉

**问题？** 查看 [FAQ](https://codeberg.org/jiangfire/fzjjyz/discussions) 或创建 Issue

**版本**: v0.1.0
**最后更新**: 2025-12-21
**维护者**: fzjjyz 社区