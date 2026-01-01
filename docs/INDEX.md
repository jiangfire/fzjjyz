# fzj 文档中心

**后量子文件加密工具 | v0.2.0**

---

## 📚 文档导航

### 🎯 用户文档（从零开始）

| 文档 | 说明 | 适合人群 |
|------|------|----------|
| **[README_MAIN.md](README_MAIN.md)** | 项目总览、特性、快速开始 | 所有用户 |
| **[INSTALL.md](INSTALL.md)** | 安装和构建指南 | 新用户 |
| **[USAGE.md](USAGE.md)** | 完整命令参考和示例 | 所有用户 |
| **[CHANGELOG.md](CHANGELOG.md)** | 版本历史和变更记录 | 所有用户 |

### 🔧 开发文档

| 文档 | 说明 | 适合人群 |
|------|------|----------|
| **[CONTRIBUTING.md](CONTRIBUTING.md)** | 贡献指南和开发流程 | 开发者 |
| **[ARCHITECTURE.md](ARCHITECTURE.md)** | 系统架构和模块设计 | 开发者 |
| **[CODE_QUALITY.md](CODE_QUALITY.md)** | 代码质量修复报告 | 开发者 |
| **[PERFORMANCE.md](PERFORMANCE.md)** | 性能基准和优化策略 | 开发者 |

### 📦 发布相关

| 文档 | 说明 | 适合人群 |
|------|------|----------|
| **[RELEASE_GUIDE.md](RELEASE_GUIDE.md)** | 发布流程和完整指南 | 维护者 |

### 🐛 问题分析（issues/）

| 文件 | 说明 | 优先级 |
|------|------|--------|
| **[P0-ERRCHECK.md](issues/P0-ERRCHECK.md)** | errcheck 问题 (100个) | P0 |
| **[P0-WRAPCHECK.md](issues/P0-WRAPCHECK.md)** | wrapcheck 问题 (75个) | P0 |
| **[P1-GOSEC.md](issues/P1-GOSEC.md)** | gosec 安全问题 (100个) | P1 |
| **[P1-STATICCHECK.md](issues/P1-STATICCHECK.md)** | staticcheck 问题 (45个) | P1 |
| **[P2-REVIVE.md](issues/P2-REVIVE.md)** | revive 规范问题 (57个) | P2 |
| **[P2-GODOT.md](issues/P2-GODOT.md)** | godot 注释问题 (50个) | P2 |
| **[P3-LOW.md](issues/P3-LOW.md)** | 低优先级问题 (37个) | P3 |

---

## 🚀 快速开始

### 首次使用？
1. 阅读 **[README_MAIN.md](README_MAIN.md)** 了解项目
2. 按照 **[INSTALL.md](INSTALL.md)** 安装工具
3. 查看 **[USAGE.md](USAGE.md)** 学习使用

### 想要贡献？
1. 阅读 **[CONTRIBUTING.md](CONTRIBUTING.md)** 了解流程
2. 查看 **[CODE_QUALITY.md](CODE_QUALITY.md)** 了解代码标准
3. 参考 **[ARCHITECTURE.md](ARCHITECTURE.md)** 理解系统设计
4. 运行 `golangci-lint run` 检查代码质量

### 遇到问题？
1. 查看 **[USAGE.md](USAGE.md)** 的错误处理章节
2. 搜索 **[CHANGELOG.md](CHANGELOG.md)** 看是否已知问题
3. 在 issues 中搜索或创建新 Issue

---

## 📊 项目状态概览

| 指标 | 状态 | 详情 |
|------|------|------|
| **版本** | ✅ v0.2.0 | 最新稳定版 |
| **代码质量** | ✅ 100% | 462/462 问题已修复 |
| **测试覆盖** | ✅ >80% | 关键模块 >80% |
| **构建状态** | ✅ 100% | 无编译错误 |
| **文档完整性** | ✅ 完整 | 14+ 文档文件 |

---

## 🔍 按主题查找

### 加密相关
- **算法说明**: README_MAIN.md#技术架构
- **使用示例**: USAGE.md#完整工作流示例
- **性能数据**: PERFORMANCE.md

### 开发相关
- **架构设计**: ARCHITECTURE.md
- **代码规范**: CODE_QUALITY.md
- **贡献流程**: CONTRIBUTING.md

### 发布相关
- **发布流程**: RELEASE_GUIDE.md
- **检查清单**: RELEASE_CHECKLIST.md
- **部署总结**: DEPLOYMENT_SUMMARY.md

---

## 📞 获取帮助

### 文档内搜索
```bash
# 在文档中搜索关键词
grep -r "关键词" docs/
```

### 命令帮助
```bash
# 查看所有命令
fzj --help

# 查看特定命令
fzj encrypt --help
```

### 详细输出
```bash
# 使用 verbose 模式
fzj encrypt -i input.txt -o output.fzj -p pub.pem -s priv.pem -v
```

---

## 🔄 文档更新日志

### 2025-01-01
- ✅ 文档归档完成
- ✅ 所有文件移至 docs/ 目录
- ✅ 创建 INDEX.md 导航
- ✅ 清理临时文件

### 2025-12-31
- ✅ 添加 DEPLOYMENT_SUMMARY.md
- ✅ 更新 RELEASE_GUIDE.md
- ✅ 完善问题分析文档

---

**维护者**: fzj 开发团队
**最后更新**: 2025-01-01
**版本**: v0.2.0