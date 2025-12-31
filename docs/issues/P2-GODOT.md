# P2 - Godot 问题修复清单

**优先级：🟢 中**
**数量：50个**
**风险：文档规范**
**状态：待修复**

---

## 📋 问题详情

**问题：** 所有中文注释都需要以句号结尾

### 1. cmd/fzjjyz/main.go (3个)

**第10行：**
```go
// 版本信息
```
**修复：**
```go
// 版本信息。
```

**第17行：**
```go
// 根命令（文本将在 init 中通过 i18n 翻译）
```
**修复：**
```go
// 根命令（文本将在 init 中通过 i18n 翻译）。
```

**第20行：**
```go
// 全局标志
```
**修复：**
```go
// 全局标志。
```

**状态：** ⬜ 待修复

---

### 2. cmd/fzjjyz/main_test.go (4个)

**第13行：**
```go
// TestCLIIntegration CLI 集成测试
```
**修复：**
```go
// TestCLIIntegration CLI 集成测试。
```

**第250行：**
```go
// buildCLI 构建 CLI 可执行文件
```
**修复：**
```go
// buildCLI 构建 CLI 可执行文件。
```

**第282行：**
```go
// TestCLIHelp 测试帮助信息
```
**修复：**
```go
// TestCLIHelp 测试帮助信息。
```

**第317行：**
```go
// TestCLIBenchmark 简单的性能测试
```
**修复：**
```go
// TestCLIBenchmark 简单的性能测试。
```

**第417行：**
```go
// 辅助函数
```
**修复：**
```go
// 辅助函数。
```

**状态：** ⬜ 待修复

---

### 3. cmd/fzjjyz/utils/errors.go (11个)

**第8行：**
```go
// ErrorType 定义错误类型
```

**第25行：**
```go
// UserError 用户友好的错误包装
```

**第43行：**
```go
// NewUserError 创建用户友好错误
```

**第52行：**
```go
// WrapError 包装原始错误为用户友好错误
```

**第99行：**
```go
// PrintError 打印用户友好的错误信息
```

**第118行：**
```go
// PrintWarning 打印警告信息
```

**第123行：**
```go
// PrintSuccess 打印成功信息
```

**第128行：**
```go
// PrintInfo 打印信息
```

**第133行：**
```go
// ConfirmPrompt 确认提示
```

**第157行：**
```go
// SelectPrompt 选择提示
```

**状态：** ⬜ 待修复

---

### 4. cmd/fzjjyz/utils/progress.go (10个)

**第10行：**
```go
// ProgressBar 进度条结构
```

**第20行：**
```go
// NewProgressBar 创建新的进度条
```

**第29行：**
```go
// Add 增加进度
```

**第39行：**
```go
// Set 设置当前进度
```

**第49行：**
```go
// Complete 完成进度条
```

**第59行：**
```go
// render 渲染进度条
```

**第114行：**
```go
// ProgressReader 包装 Reader 以显示进度
```

**第120行：**
```go
// NewProgressReader 创建进度读取器
```

**第128行：**
```go
// Read 实现 io.Reader 接口
```

**第137行：**
```go
// Close 完成进度
```

**第142行：**
```go
// ProgressWriter 包装 Writer 以显示进度
```

**第148行：**
```go
// NewProgressWriter 创建进度写入器
```

**第156行：**
```go
// Write 实现 io.Writer 接口
```

**第165行：**
```go
// Close 完成进度
```

**状态：** ⬜ 待修复

---

### 5. internal/crypto/archive.go (5个)

**第14行：**
```go
// ArchiveOptions 打包选项
```

**第21行：**
```go
// DefaultArchiveOptions 默认打包选项
```

**第30行：**
```go
// 返回: 错误
```

**第132行：**
```go
// 返回: 错误
```

**第218行：**
```go
// GetZipSize 计算ZIP数据的总大小（用于进度条）
```

**第233行：**
```go
// CountZipFiles 统计ZIP中的文件数量
```

**状态：** ⬜ 待修复

---

### 6. internal/crypto/buffer_pool.go (5个)

**第6行：**
```go
// DefaultBufferSize 默认缓冲区大小：64KB
```

**第9行：**
```go
// MaxBufferSize 最大缓冲区大小：1MB
```

**第12行：**
```go
// MinBufferSize 最小缓冲区大小：4KB
```

**第16行：**
```go
// BufferPool 缓冲区池，用于减少 GC 压力
```

**第21行：**
```go
// NewBufferPool 创建新的缓冲区池
```

**第39行：**
```go
// Get 从池中获取一个缓冲区
```

**第44行：**
```go
// Put 将缓冲区归还到池中
```

**第54行：**
```go
// OptimalBufferSize 根据文件大小推荐缓冲区大小
```

**状态：** ⬜ 待修复

---

### 7. internal/crypto/hash_utils.go (4个)

**第10行：**
```go
// 使用 io.Copy 避免一次性读取整个文件到内存
```

**第34行：**
```go
// HashReader 流式计算 Reader 的 SHA256 哈希值
```

**第48行：**
```go
// StreamingHash 流式哈希计算器，支持 Write 接口
```

**第55行：**
```go
// NewStreamingHash 创建新的流式哈希计算器
```

**状态：** ⬜ 待修复

---

### 8. internal/crypto/keygen.go (1个)

**第26行：**
```go
// 生成Kyber密钥对
```

**状态：** ⬜ 待修复

---

### 9. internal/crypto/operations_shared.go (1个)

**第166行：**
```go
func decapsulateKeys(kyberPriv kem.PrivateKey, ecdhPriv *ecdh.PrivateKey, encapsulated []byte, ecdhPub []byte) ([]byte, error) {
```
**缺少注释**

**状态：** ⬜ 待修复

---

### 10. internal/format/header.go (3个)

**第12行：**
```go
// FileHeader 文件头结构（表达原则：数据结构优先）
```

**第33行：**
```go
// MarshalBinary 序列化为二进制（压缩格式）
```

**第109行：**
```go
// MarshalBinaryOptimized 优化后的序列化（减少内存分配）
```

**状态：** ⬜ 待修复

---

### 11. internal/i18n/i18n.go (1个)

**第46行：**
```go
// 无法加载默认语言 %s: %v
```

**状态：** ⬜ 待修复

---

## 📊 统计信息

| 文件类型 | 数量 | 预计时间 |
|---------|------|---------|
| cmd/fzjjyz/ | 18个 | 10分钟 |
| internal/crypto/ | 12个 | 8分钟 |
| internal/format/ | 3个 | 2分钟 |
| internal/i18n/ | 1个 | 1分钟 |
| 其他 | 16个 | 10分钟 |
| **总计** | **50个** | **31分钟** |

---

## 🔧 修复方法

### 方法1：自动修复（推荐）
```bash
golangci-lint run --fix
```
此命令会自动修复所有 godot 问题。

### 方法2：手动修复
在所有中文注释末尾添加句号 `。`

**示例：**
```go
// ❌ 原代码
// 版本信息
// 用户友好的错误包装

// ✅ 修复后
// 版本信息。
// 用户友好的错误包装。
```

---

## ✅ 验证标准

修复后运行：
```bash
golangci-lint run --disable-all --enable=godot
```

应输出：`0 issues`

---

## 📝 注意事项

1. **英文注释不需要句号** - 只有中文注释需要
2. **行内注释** - 如果是完整句子，也需要句号
3. **代码块注释** - 简短说明可不加句号
4. **自动修复** - 使用 `--fix` 参数可自动处理大部分问题

---

**创建时间：** 2025-12-30
**预计完成：** 2025-12-30
**负责人：** 待分配
