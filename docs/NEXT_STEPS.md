# 下一步行动计划

本文档整合了当前项目状态、CI 问题分析和开发计划，提供完整的行动指南。

## 📋 当前项目状态

### ✅ 已完成功能
- 核心加密/解密功能 (Kyber768 + ECDH + AES-256-GCM + Dilithium3)
- 目录加密/解密功能
- 国际化支持 (中英文)
- 密钥缓存安全增强 (TTL + 大小限制)
- 性能基准测试

### ⚠️ 当前问题
- **编译失败**: 9 处非恒定格式字符串错误
- **CI 问题**: 448 个 golangci-lint 问题
- **Git 状态**: 多个文件已修改但未提交

---

## 🚨 P0 - 立即行动 (今天必须完成)

### 1. 修复编译错误

**问题**: 项目无法编译，需要修复 9 处格式字符串错误

**执行**:
```bash
# 1. 查看具体错误
go vet ./...

# 2. 修复以下文件中的格式字符串问题
# internal/i18n/cobra.go:16,20,28
# cmd/fzjjyz/decrypt.go:164
# cmd/fzjjyz/decrypt_dir.go:182
# cmd/fzjjyz/encrypt.go:136
# cmd/fzjjyz/encrypt_dir.go:160
# cmd/fzjjyz/keygen.go:64
# cmd/fzjjyz/keymanage.go:93,139

# 3. 验证修复
go build ./cmd/fzjjyz
```

**修复示例**:
```go
// ❌ 错误
fmt.Printf(i18n.T("status.success_encrypt") + "\n")

// ✅ 正确
fmt.Printf("%s\n", i18n.T("status.success_encrypt"))
```

---

### 2. 运行测试验证功能

**目标**: 确保所有功能正常工作

**执行**:
```bash
# 运行所有测试
go test ./...

# 运行性能基准测试
go test -bench=. ./internal/crypto/

# 查看测试覆盖率
go test -cover ./...
```

**预期结果**:
- 所有测试通过
- 基准测试正常运行
- 无编译错误

---

### 3. 提交当前变更

**目标**: 保存当前开发进度

**执行**:
```bash
# 1. 查看当前状态
git status

# 2. 添加所有变更
git add .

# 3. 提交代码
git commit -m "feat: 添加目录加密/解密功能、国际化支持、密钥缓存安全增强

主要变更:
- 新增 encrypt-dir 和 decrypt-dir 命令
- 实现国际化系统 (zh_CN, en_US)
- 密钥缓存安全增强 (TTL + 大小限制)
- 性能基准测试套件
- 文件头序列化优化

技术细节:
- 使用 archive/zip 处理目录打包
- i18n 系统自动检测 LANG 环境变量
- 缓存最大 100 个密钥，TTL 1 小时
- 并行密钥生成优化性能"

# 4. 推送到远程
git push origin master
```

---

### 4. 构建验证

**目标**: 确保二进制文件正常工作

**执行**:
```bash
# 构建
go build ./cmd/fzjjyz

# 验证版本信息
./fzjjyz version

# 测试基本命令
./fzjjyz --help
./fzjjyz encrypt --help
./fzjjyz decrypt --help
./fzjjyz keygen --help
```

---

## 🟠 P1 - 短期优化 (1-2周内)

### 5. 完善测试覆盖

**目标**: 测试覆盖率 > 80%

**重点测试**:
```bash
# 1. 目录操作测试
go test -v ./cmd/fzjjyz/... -run "Dir"

# 2. 国际化测试
go test -v ./internal/i18n/...

# 3. 集成测试
go test -v ./cmd/fzjjyz/...

# 4. 边界测试
# - 大文件 (>100MB)
# - 空文件
# - 特殊字符文件名
# - 错误场景
```

**需要添加的测试**:
- `internal/i18n/` - 翻译准确性测试
- `cmd/fzjjyz/decrypt_dir.go` - 目录解密测试
- `cmd/fzjjyz/encrypt_dir.go` - 目录加密测试
- `internal/crypto/archive.go` - 归档处理测试

---

### 6. 文档更新

**目标**: 文档与代码同步

**需要更新的文件**:

#### 6.1 USAGE.md
```bash
# 添加目录加密/解密示例
# 更新命令列表
# 添加国际化说明
# 更新性能数据
```

#### 6.2 CHANGELOG.md
```bash
# 记录 v0.1.1 变更
# 添加日期和版本号
# 分类记录功能/修复/优化
```

#### 6.3 README.md (可选)
```bash
# 更新快速开始示例
# 添加目录操作说明
# 更新性能指标
```

---

### 7. 代码清理

**目标**: 保持代码库整洁

**执行**:
```bash
# 1. 检查 .gitignore
cat .gitignore

# 2. 移除不必要的测试文件
# 检查 test_i18n.go 是否需要保留
# 检查 testdir/ 和 testfile.txt

# 3. 确认不提交的文件
git status

# 4. 清理临时文件
rm -f test_i18n.go  # 如果不需要
rm -rf testdir/ testfile.txt  # 如果是测试数据

# 5. 更新 .gitignore (如果需要)
echo "testdir/" >> .gitignore
echo "testfile.txt" >> .gitignore
echo "test_i18n.go" >> .gitignore
```

---

## 🟡 P2 - CI/CD 集成 (2-3周内)

### 8. 修复 CI 问题

**目标**: 通过 golangci-lint 检查

**优先级修复**:

#### 8.1 关键问题 (必须修复)
```bash
# 1. errcheck - 关键位置 (20处)
# - defer Close() 错误处理
# - MarkFlagRequired() 错误检查
# - os.Remove() 错误处理

# 2. wrapcheck - 核心代码 (30处)
# - 外部包错误包装
# - 添加错误上下文
```

#### 8.2 配置调整
```bash
# 如果某些规则过于严格，调整 .golangci.yml
# 例如：测试文件放宽要求
```

---

### 9. 配置 GitHub Actions

**创建文件**: `.github/workflows/go-ci.yml`

```yaml
name: Go CI

on:
  push:
    branches: [ master ]
  pull_request:
    branches: [ master ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Cache Go modules
        uses: actions/cache@v4
        with:
          path: |
            ~/go/pkg/mod
            ~/.cache/go-build
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

      - name: Install dependencies
        run: go mod download

      - name: Run go vet
        run: go vet ./...

      - name: Run tests
        run: go test -v ./... -coverprofile=coverage.out

      - name: Build
        run: go build ./cmd/fzjjyz

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
```

---

## 📊 完整执行时间表

| 阶段 | 任务 | 预计时间 | 优先级 |
|------|------|----------|--------|
| **Day 1** | 修复编译错误 | 30分钟 | P0 |
| **Day 1** | 运行测试验证 | 15分钟 | P0 |
| **Day 1** | 提交代码 | 10分钟 | P0 |
| **Day 1** | 构建验证 | 5分钟 | P0 |
| **Day 2-3** | 完善测试覆盖 | 4-6小时 | P1 |
| **Day 4** | 文档更新 | 2-3小时 | P1 |
| **Day 5** | 代码清理 | 1小时 | P1 |
| **Week 2** | CI 问题修复 | 6-8小时 | P2 |
| **Week 2** | GitHub Actions | 2小时 | P2 |

---

## ✅ 验证清单

### P0 完成标准
- [ ] `go build ./cmd/fzjjyz` 成功
- [ ] `go test ./...` 全部通过
- [ ] `go test -bench=. ./internal/crypto/` 正常
- [ ] `git commit` 成功提交
- [ ] `./fzjjyz version` 显示正确版本

### P1 完成标准
- [ ] 测试覆盖率 > 80%
- [ ] USAGE.md 更新完成
- [ ] CHANGELOG.md 更新完成
- [ ] .gitignore 正确配置
- [ ] 无未跟踪的测试文件

### P2 完成标准
- [ ] `golangci-lint run` 无错误
- [ ] GitHub Actions 配置完成
- [ ] CI 流程正常运行
- [ ] 代码质量检查通过

---

## 🎯 里程碑

### 里程碑 1: 代码提交 (今天)
- 修复编译错误
- 验证功能正常
- 提交并推送代码

### 里程碑 2: 测试完善 (本周)
- 测试覆盖率达标
- 文档更新完成
- 代码库整洁

### 里程碑 3: CI 就绪 (下周)
- 通过 lint 检查
- CI/CD 配置完成
- 自动化测试

---

## 📝 执行脚本

### 一键执行脚本 (P0)

```bash
#!/bin/bash
echo "=== fzjjyz P0 验证脚本 ==="

echo "1. 检查编译..."
go build ./cmd/fzjjyz
if [ $? -ne 0 ]; then
    echo "❌ 编译失败，请先修复"
    exit 1
fi
echo "✅ 编译成功"

echo "2. 运行测试..."
go test ./...
if [ $? -ne 0 ]; then
    echo "❌ 测试失败"
    exit 1
fi
echo "✅ 测试通过"

echo "3. 运行基准测试..."
go test -bench=. ./internal/crypto/
echo "✅ 基准测试完成"

echo "4. 验证版本..."
./fzjjyz version
echo "✅ 版本验证完成"

echo "=== 所有 P0 检查通过 ==="
```

---

## 🔖 参考文档

- [CI 问题分析](./CI_ISSUES.md) - 详细的 lint 问题分析
- [开发计划](./DEVELOPMENT_PLAN.md) - 长期开发规划
- [架构文档](./ARCHITECTURE.md) - 系统架构说明
- [性能文档](./PERFORMANCE.md) - 性能基准和优化

---

**文档版本**: 1.0
**创建日期**: 2025-12-30
**状态**: 待执行
**下一步**: 修复编译错误 → 运行测试 → 提交代码
