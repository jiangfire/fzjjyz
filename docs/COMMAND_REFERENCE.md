# 常用命令参考文档

这个文档记录了我在使用 golangci-lint 和其他工具时经常遇到的命令错误，避免重复踩坑。

## 🐛 golangci-lint 命令错误

### ❌ 错误的命令

```bash
# 错误 1: 使用 --disable-all 标志
golangci-lint run --disable-all --enable=errcheck

# 错误 2: 使用 -D all
golangci-lint run -D all -E errcheck

# 错误 3: 使用 --path-pattern
golangci-lint run --disable-all --enable=errcheck --path-pattern=internal/i18n/i18n_test.go
```

### ✅ 正确的命令

```bash
# ✅ 正确 1: 启用特定 linter
golangci-lint run --enable=errcheck

# ✅ 正确 2: 禁用特定 linter（不是 all）
golangci-lint run --disable=staticcheck --enable=errcheck

# ✅ 正确 3: 检查特定目录
golangci-lint run --enable=errcheck ./internal/i18n/

# ✅ 正确 4: 检查特定文件
golangci-lint run --enable=errcheck internal/i18n/i18n_test.go
```

## 📋 常用 golangci-lint 命令

### 按 linter 类型检查

```bash
# 只检查 errcheck
golangci-lint run --enable=errcheck

# 只检查 wrapcheck
golangci-lint run --enable=wrapcheck

# 只检查 gosec
golangci-lint run --enable=gosec

# 只检查 staticcheck
golangci-lint run --enable=staticcheck

# 只检查 revive
golangci-lint run --enable=revive

# 只检查 godot
golangci-lint run --enable=godot
```

### 检查特定目录

```bash
# 检查 cmd 目录
golangci-lint run ./cmd/fzjjyz/

# 检查 internal/crypto 目录
golangci-lint run ./internal/crypto/

# 检查所有 Go 文件
golangci-lint run ./...
```

### 自动修复

```bash
# 自动修复可修复的问题
golangci-lint run --fix

# 只修复特定 linter
golangci-lint run --enable=godot --fix
```

### 查看支持的 linter

```bash
# 查看所有可用 linter
golangci-lint help linters

# 查看运行状态
golangci-lint run --help
```

## 🧪 Go 测试命令

### 基本测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/i18n/...

# 运行特定包的详细测试
go test ./internal/i18n/... -v

# 运行特定测试函数
go test ./internal/i18n/... -v -run TestInit

# 显示测试覆盖率
go test ./... -cover

# 显示详细覆盖率
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试特定文件

```bash
# 测试单个文件
go test ./internal/i18n/i18n_test.go ./internal/i18n/i18n.go

# 或者进入目录测试
cd internal/i18n && go test -v
```

## 🔨 Go 构建命令

### 构建

```bash
# 构建所有
go build ./...

# 构建特定包
go build ./cmd/fzjjyz

# 构建并安装
go install ./cmd/fzjjyz
```

### 模块管理

```bash
# 清理依赖
go mod tidy

# 查看依赖
go list -m all

# 更新依赖
go get -u ./...
```

## 📊 Git 常用命令

### 查看状态

```bash
# 查看修改状态
git status

# 查看修改统计
git diff --stat

# 查看具体修改
git diff

# 查看某个文件的修改
git diff internal/i18n/i18n_test.go
```

### 提交代码

```bash
# 添加文件
git add .

# 提交
git commit -m "fix: 修复 i18n_test.go 的 errcheck 问题"

# 查看最近提交
git log --oneline -5
```

## 🎯 修复流程常用命令

### 1. 检查问题

```bash
# 查看当前目录的所有问题
golangci-lint run

# 只看 errcheck 问题
golangci-lint run --enable=errcheck

# 只看 wrapcheck 问题
golangci-lint run --enable=wrapcheck
```

### 2. 修复后验证

```bash
# 运行测试
go test ./... -v

# 构建验证
go build ./...

# 再次检查 linter
golangci-lint run --enable=errcheck
```

### 3. 查看改动

```bash
# 查看修改统计
git diff --stat

# 查看具体代码改动
git diff
```

## ⚠️ 常见错误总结

| 错误命令 | 正确命令 | 原因 |
|---------|---------|------|
| `--disable-all` | `--enable=xxx` | golangci-lint 没有 `--disable-all` 标志 |
| `-D all` | `-D flagname` | `all` 不是有效的 linter 名称 |
| `--path-pattern` | 直接指定路径 | 没有这个标志，直接在命令后加路径 |

## 💡 提示

1. **始终使用 `--enable=xxx` 来启用特定 linter**
2. **检查特定目录时，使用 `./path/to/dir/` 格式**
3. **修复后务必运行测试验证**
4. **使用 `go test ./... -v` 查看详细测试输出**

---

**创建时间：** 2025-12-31
**最后更新：** 2025-12-31
**维护者：** Claude Code
