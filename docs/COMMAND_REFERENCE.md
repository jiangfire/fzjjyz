# 常用命令参考文档

## 📋 golangci-lint 命令

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

## 💡 提示

1. **始终使用 `--enable=xxx` 来启用特定 linter**
2. **检查特定目录时，使用 `./path/to/dir/` 格式**
3. **修复后务必运行测试验证**
4. **使用 `go test ./... -v` 查看详细测试输出**

---

**创建时间：** 2025-12-31
**最后更新：** 2025-12-31
**维护者：** fzjjyz 开发团队
