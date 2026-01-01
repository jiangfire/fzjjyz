# fzj - 后量子文件加密工具
# Makefile for building, testing, and releasing

# ==================== 变量定义 ====================

# 项目信息
PROJECT_NAME := fzj
MODULE_NAME := codeberg.org/jiangfire/fzj
MAIN_PACKAGE := cmd/fzjjyz

# 版本信息 (从 git 获取， fallback 到 main.go 中的版本)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || grep -oP 'Version\s*=\s*"\K[^"]+' cmd/fzjjyz/main.go)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建输出目录
DIST_DIR := dist
BUILD_DIR := build

# Go 构建参数
GO := go
GOFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
GOOS_LIST := linux darwin windows
GOARCH_LIST := amd64 arm64

# 测试相关
TEST_FLAGS := -v -race -coverprofile=coverage.out -covermode=atomic
TEST_COVERAGE_THRESHOLD := 80

# Lint 工具
GOLANGCI_LINT := golangci-lint

# ==================== 基础目标 ====================

.PHONY: help build build-dev test test-cover lint clean tidy vendor cross-build package checksum release ci

# 显示帮助信息
help:
	@echo "可用目标:"
	@echo "  build          - 构建当前平台的发布版本"
	@echo "  build-dev      - 构建开发版本（带调试信息）"
	@echo "  test           - 运行所有测试"
	@echo "  test-cover     - 运行测试并生成覆盖率报告"
	@echo "  lint           - 运行代码检查"
	@echo "  lint-fix       - 运行代码检查并自动修复"
	@echo "  clean          - 清理构建产物"
	@echo "  tidy           - 整理 go.mod"
	@echo "  vendor         - 生成 vendor 目录"
	@echo "  cross-build    - 跨平台构建（CI 用）"
	@echo "  package        - 打包所有构建产物"
	@echo "  checksum       - 生成校验和"
	@echo "  release        - 完整发布流程（构建 + 打包 + 校验）"
	@echo "  ci             - CI 流程（lint + test + build）"
	@echo "  version        - 显示版本信息"

# 构建当前平台版本
build:
	@echo "→ 构建 $(PROJECT_NAME) v$(VERSION)..."
	@mkdir -p $(DIST_DIR)
	$(GO) build $(GOFLAGS) -o $(DIST_DIR)/$(PROJECT_NAME) ./$(MAIN_PACKAGE)
	@echo "✓ 构建完成: $(DIST_DIR)/$(PROJECT_NAME)"

# 构建开发版本（带调试信息）
build-dev:
	@echo "→ 构建开发版本..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(PROJECT_NAME) ./$(MAIN_PACKAGE)
	@echo "✓ 开发版本完成: $(BUILD_DIR)/$(PROJECT_NAME)"

# 运行测试
test:
	@echo "→ 运行测试..."
	$(GO) test ./... $(TEST_FLAGS) -timeout 5m
	@echo "✓ 测试完成"

# 运行测试并生成覆盖率报告
test-cover:
	@echo "→ 运行测试并生成覆盖率报告..."
	$(GO) test ./... $(TEST_FLAGS) -timeout 5m
	@echo "→ 生成覆盖率报告..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ 覆盖率报告: coverage.html"
	@echo "→ 检查覆盖率阈值..."
	@$(GO) tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//' | awk '{if ($$1 < $(TEST_COVERAGE_THRESHOLD)) {print "✗ 覆盖率低于阈值: " $$1 "% < $(TEST_COVERAGE_THRESHOLD)%"; exit 1} else {print "✓ 覆盖率达标: " $$1 "% >= $(TEST_COVERAGE_THRESHOLD)%"}}'

# 代码检查
lint:
	@echo "→ 运行代码检查..."
	$(GOLANGCI_LINT) run --timeout 5m
	@echo "✓ 代码检查完成"

# 代码检查并自动修复
lint-fix:
	@echo "→ 运行代码检查并自动修复..."
	$(GOLANGCI_LINT) run --fix --timeout 5m
	@echo "✓ 代码检查和修复完成"

# 清理
clean:
	@echo "→ 清理构建产物..."
	rm -rf $(DIST_DIR) $(BUILD_DIR) coverage.out coverage.html
	@echo "✓ 清理完成"

# 整理 go.mod
tidy:
	@echo "→ 整理 go.mod..."
	$(GO) mod tidy
	@echo "✓ go.mod 已整理"

# 生成 vendor
vendor:
	@echo "→ 生成 vendor..."
	$(GO) mod vendor
	@echo "✓ vendor 目录已生成"

# ==================== 发布相关目标 ====================

# 跨平台构建（GitHub Actions 版本）
cross-build:
	@echo "→ 跨平台构建 (版本: $(VERSION))..."
	@mkdir -p $(DIST_DIR)
	@for os in $(GOOS_LIST); do \
		for arch in $(GOARCH_LIST); do \
			echo "  → 构建 $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch $(GO) build $(GOFLAGS) -o $(DIST_DIR)/$(PROJECT_NAME)-$$os-$$arch ./$(MAIN_PACKAGE); \
			if [ "$$os" = "windows" ]; then \
				mv $(DIST_DIR)/$(PROJECT_NAME)-$$os-$$arch $(DIST_DIR)/$(PROJECT_NAME)-$$os-$$arch.exe; \
			fi; \
		done; \
	done
	@echo "✓ 跨平台构建完成: $(DIST_DIR)/"

# 打包（仅用于本地，GitHub Actions 使用内置打包）
package:
	@echo "→ 打包构建产物..."
	@cd $(DIST_DIR) && for file in $(PROJECT_NAME)-*; do \
		if [ -f "$$file" ]; then \
			echo "  → 打包 $$file..."; \
			tar -czf "$$file.tar.gz" "$$file" 2>/dev/null || zip -q "$$file.zip" "$$file"; \
		fi; \
	done
	@echo "✓ 打包完成"

# 生成校验和
checksum:
	@echo "→ 生成校验和..."
	@cd $(DIST_DIR) && \
		if command -v sha256sum >/dev/null 2>&1; then \
			sha256sum $(PROJECT_NAME)-* > checksums.sha256; \
			echo "✓ 校验和已生成: checksums.sha256"; \
			cat checksums.sha256; \
		elif command -v shasum >/dev/null 2>&1; then \
			shasum -a 256 $(PROJECT_NAME)-* > checksums.sha256; \
			echo "✓ 校验和已生成: checksums.sha256"; \
			cat checksums.sha256; \
		else \
			echo "⚠  未找到 sha256sum/shasum 工具，跳过校验和生成"; \
		fi

# 完整发布流程（本地使用）
release: clean lint test cross-build checksum
	@echo ""
	@echo "========== 发布完成 =========="
	@echo "版本: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "产物目录: $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/
	@echo "============================"

# ==================== CI/CD 目标 ====================

# CI 流程（快速检查）
ci: lint test build
	@echo "✓ CI 流程完成"

# CI 发布流程（跳过 lint，加速构建）
ci-release: test cross-build checksum
	@echo ""
	@echo "========== CI 发布完成 =========="
	@echo "版本: $(VERSION)"
	@echo "产物: $(DIST_DIR)/"
	@echo "================================"

# ==================== 开发辅助目标 ====================

# 快速开发构建和测试
dev: build-dev test
	@echo "✓ 开发构建和测试完成"

# 运行本地测试
run-test:
	@echo "→ 运行快速测试..."
	$(GO) test ./internal/zjcrypto -v -run TestEncryptFile
	@echo "✓ 快速测试完成"

# 检查依赖更新
outdated:
	@echo "→ 检查依赖更新..."
	$(GO) list -u -m all | grep -v '\[.*\]'

# 显示版本信息
version:
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"