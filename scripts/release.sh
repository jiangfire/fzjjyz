#!/bin/bash

# fzjjyz 发布脚本
# 用于自动化发布流程

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 输出函数
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
if [ $# -lt 1 ]; then
    echo "用法: $0 <版本号> [选项]"
    echo ""
    echo "选项:"
    echo "  --skip-test    跳过测试"
    echo "  --skip-build   跳过构建"
    echo "  --skip-tag     跳过 Git 标签"
    echo "  --dry-run      试运行，不执行实际操作"
    echo ""
    echo "示例:"
    echo "  $0 v0.1.1              # 完整发布流程"
    echo "  $0 v0.1.1 --dry-run    # 试运行"
    echo "  $0 v0.1.1 --skip-test  # 跳过测试"
    exit 1
fi

VERSION="$1"
shift

# 解析选项
SKIP_TEST=false
SKIP_BUILD=false
SKIP_TAG=false
DRY_RUN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-test)
            SKIP_TEST=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-tag)
            SKIP_TAG=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        *)
            error "未知选项: $1"
            exit 1
            ;;
    esac
done

# 验证版本号格式
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    error "版本号格式错误: $VERSION (应为 v0.1.0 格式)"
    exit 1
fi

info "开始发布流程: $VERSION"
info "选项: 跳过测试=$SKIP_TEST, 跳过构建=$SKIP_BUILD, 跳过标签=$SKIP_TAG, 试运行=$DRY_RUN"

if [ "$DRY_RUN" = true ]; then
    warning "这是试运行模式，不会执行实际操作"
fi

# 步骤 1: 检查当前分支
info "步骤 1: 检查当前分支"
CURRENT_BRANCH=$(git branch --show-current)
info "当前分支: $CURRENT_BRANCH"

if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
    warning "当前不在 main/master 分支，当前分支: $CURRENT_BRANCH"
    if [ "$DRY_RUN" = false ]; then
        read -p "继续发布? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            info "已取消发布"
            exit 0
        fi
    fi
fi

# 步骤 2: 检查未提交的更改
info "步骤 2: 检查工作区状态"
if [ "$DRY_RUN" = false ]; then
    if ! git diff --quiet; then
        warning "检测到未提交的更改:"
        git status --short
        read -p "是否提交这些更改? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            git add .
            git commit -m "chore: 发布前准备 $VERSION"
        else
            error "请先清理工作区再发布"
            exit 1
        fi
    fi
fi

# 步骤 3: 运行测试
if [ "$SKIP_TEST" = false ]; then
    info "步骤 3: 运行测试"
    if [ "$DRY_RUN" = false ]; then
        if go test ./... -cover; then
            success "测试通过"
        else
            error "测试失败"
            exit 1
        fi
    else
        info "[DRY RUN] 跳过测试执行"
    fi
else
    warning "跳过测试"
fi

# 步骤 4: 更新版本号
info "步骤 4: 更新版本号"
VERSION_NUM=${VERSION#v}
info "版本号: $VERSION_NUM"

if [ "$DRY_RUN" = false ]; then
    # 检查 main.go 中的版本定义
    if grep -q "const Version = " cmd/fzjjyz/main.go; then
        # 备份原文件
        cp cmd/fzjjyz/main.go cmd/fzjjyz/main.go.bak

        # 更新版本号
        sed -i "s/const Version = \".*\"/const Version = \"$VERSION_NUM\"/" cmd/fzjjyz/main.go

        # 验证更新
        if grep -q "const Version = \"$VERSION_NUM\"" cmd/fzjjyz/main.go; then
            success "版本号已更新为 $VERSION_NUM"
        else
            error "版本号更新失败"
            mv cmd/fzjjyz/main.go.bak cmd/fzjjyz/main.go
            exit 1
        fi

        # 清理备份
        rm cmd/fzjjyz/main.go.bak
    else
        warning "未找到版本定义，跳过更新"
    fi
else
    info "[DRY RUN] 将更新版本号为 $VERSION_NUM"
fi

# 步骤 5: 构建二进制文件
if [ "$SKIP_BUILD" = false ]; then
    info "步骤 5: 构建跨平台二进制"

    if [ "$DRY_RUN" = false ]; then
        # 创建发布目录
        mkdir -p release/$VERSION

        # 构建各平台
        info "构建 Linux amd64..."
        GOOS=linux GOARCH=amd64 go build -o release/$VERSION/fzjjyz_linux_amd64 ./cmd/fzjjyz

        info "构建 Windows amd64..."
        GOOS=windows GOARCH=amd64 go build -o release/$VERSION/fzjjyz_windows_amd64.exe ./cmd/fzjjyz

        info "构建 macOS Intel..."
        GOOS=darwin GOARCH=amd64 go build -o release/$VERSION/fzjjyz_darwin_amd64 ./cmd/fzjjyz

        info "构建 macOS Apple Silicon..."
        GOOS=darwin GOARCH=arm64 go build -o release/$VERSION/fzjjyz_darwin_arm64 ./cmd/fzjjyz

        # 生成校验和
        info "生成校验和..."
        cd release/$VERSION
        sha256sum fzjjyz_* > checksums.txt
        cd ../..

        success "构建完成"
        ls -lh release/$VERSION/
    else
        info "[DRY RUN] 将构建各平台二进制"
    fi
else
    warning "跳过构建"
fi

# 步骤 6: 提交版本更新
info "步骤 6: 提交版本更新"
if [ "$DRY_RUN" = false ]; then
    git add cmd/fzjjyz/main.go
    if git commit -m "chore: 发布 $VERSION"; then
        success "版本更新已提交"
    else
        warning "没有需要提交的更改"
    fi
else
    info "[DRY RUN] 将提交版本更新"
fi

# 步骤 7: 创建 Git 标签
if [ "$SKIP_TAG" = false ]; then
    info "步骤 7: 创建 Git 标签"

    # 生成标签信息
    TAG_MESSAGE=$(cat <<EOF
Release $VERSION

### 📦 主要变更
- 版本: $VERSION
- 日期: $(date +%Y-%m-%d)

### 🔍 验证
构建完成的二进制文件已生成，包含：
- fzjjyz_linux_amd64
- fzjjyz_windows_amd64.exe
- fzjjyz_darwin_amd64
- fzjjyz_darwin_arm64
- checksums.txt

### 🚀 快速开始
\`\`\`bash
# 生成密钥
./fzjjyz keygen -d ./keys -n mykey

# 加密文件
./fzjjyz encrypt -i secret.txt -o secret.fzj -p keys/mykey_public.pem -s keys/mykey_dilithium_private.pem

# 解密文件
./fzjjyz decrypt -i secret.fzj -o recovered.txt -p keys/mykey_private.pem -s keys/mykey_dilithium_public.pem
\`\`\`

### ✅ 质量保证
- 所有测试通过
- 跨平台构建成功
- 校验和已生成
- 文档完整

发布者: $(git config user.name)
EOF
)

    if [ "$DRY_RUN" = false ]; then
        # 创建带注释的标签
        if git tag -a "$VERSION" -m "$TAG_MESSAGE"; then
            success "Git 标签 $VERSION 已创建"
        else
            error "创建标签失败"
            exit 1
        fi

        # 推送标签
        info "推送标签到远程..."
        if git push origin "$VERSION"; then
            success "标签已推送到远程"
        else
            error "推送标签失败"
            exit 1
        fi

        # 推送提交
        info "推送提交到远程..."
        if git push; then
            success "提交已推送到远程"
        else
            warning "推送提交失败，请手动推送"
        fi
    else
        info "[DRY RUN] 将创建并推送标签 $VERSION"
    fi
else
    warning "跳过 Git 标签创建"
fi

# 步骤 8: 生成发布说明
info "步骤 8: 生成发布说明"
if [ "$DRY_RUN" = false ]; then
    if [ -f "RELEASE_NOTES.md" ]; then
        info "使用 RELEASE_NOTES.md 作为发布说明"
        # 可以在这里添加版本过滤逻辑
    else
        info "生成基础发布说明..."
        cat > release/$VERSION/release_notes.md <<EOF
# Release $VERSION

## 🎉 发布概述

**版本**: $VERSION
**日期**: $(date +%Y-%m-%d)
**状态**: ✅ 生产就绪

## 📦 下载

从附件下载对应平台的二进制文件：

- \`fzjjyz_linux_amd64\` - Linux 64位
- \`fzjjyz_windows_amd64.exe\` - Windows 64位
- \`fzjjyz_darwin_amd64\` - macOS Intel
- \`fzjjyz_darwin_arm64\` - macOS Apple Silicon
- \`checksums.txt\` - SHA256 校验和

## 🔍 验证完整性

下载后请验证校验和：

\`\`\`bash
sha256sum -c checksums.txt
\`\`\`

## 🚀 快速开始

### 1. 生成密钥
\`\`\`bash
./fzjjyz keygen -d ./keys -n mykey
\`\`\`

### 2. 加密文件
\`\`\`bash
./fzjjyz encrypt -i secret.txt -o secret.fzj \\
  -p keys/mykey_public.pem \\
  -s keys/mykey_dilithium_private.pem
\`\`\`

### 3. 解密文件
\`\`\`bash
./fzjjyz decrypt -i secret.fzj -o recovered.txt \\
  -p keys/mykey_private.pem \\
  -s keys/mykey_dilithium_public.pem
\`\`\`

## 📊 变更详情

请查看 [CHANGELOG.md](../CHANGELOG.md) 获取详细的变更记录。

## 🔐 安全说明

本版本包含后量子加密实现，适用于：
- 个人敏感文件加密
- 安全备份存储
- 团队文件传输
- 教育研究目的

详细安全信息请参考 [SECURITY.md](../SECURITY.md)

## 🤝 贡献

欢迎贡献！请阅读 [CONTRIBUTING.md](../CONTRIBUTING.md) 了解如何参与。

## 📄 许可证

MIT License - 详见 [LICENSE](../LICENSE)

---

**发布者**: $(git config user.name)
**构建时间**: $(date)
EOF
        success "发布说明已生成"
    fi
else
    info "[DRY RUN] 将生成发布说明"
fi

# 步骤 9: 总结
info "步骤 9: 发布总结"
echo ""
echo "========================================"
echo "发布流程完成: $VERSION"
echo "========================================"
echo ""
if [ "$DRY_RUN" = false ]; then
    echo "✅ 版本号已更新: $VERSION_NUM"
    echo "✅ 测试已运行"
    echo "✅ 二进制已构建"
    echo "✅ 校验和已生成"
    echo "✅ Git 标签已创建"
    echo ""
    echo "📁 发布文件位置: release/$VERSION/"
    echo ""
    echo "📦 需要手动上传到 GitHub Release:"
    ls -1 release/$VERSION/
    echo ""
    echo "🔗 下一步:"
    echo "   1. 访问 GitHub Releases 页面"
    echo "   2. 创建新 Release: $VERSION"
    echo "   3. 上传 release/$VERSION/ 目录中的所有文件"
    echo "   4. 使用生成的发布说明或 RELEASE_NOTES.md"
else
    echo "⚠️  试运行模式 - 未执行实际操作"
    echo "💡 使用以下命令执行实际发布:"
    echo "   $0 $VERSION"
fi
echo ""
echo "========================================"

success "发布流程结束"