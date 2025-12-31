package i18n

// zhCN 中文翻译字典.
type zhCN struct{}

func (z *zhCN) Get(key string) string {
	translation, exists := zhTranslations[key]
	if !exists {
		return ""
	}
	return translation
}

// zhTranslations 中文翻译映射.
var zhTranslations = map[string]string{
	// 根命令和应用信息
	"app.name":        "fzjjyz",
	"app.description": "后量子文件加密工具 - 使用 Kyber768 + ECDH + AES-256-GCM + Dilithium3",
	"app.long": `fzjjyz - 后量子文件加密工具

使用以下算法提供安全的文件加密：
  • Kyber768 - 后量子密钥封装
  • X25519 ECDH - 传统密钥交换
  • AES-256-GCM - 认证加密
  • Dilithium3 - 数字签名

快速开始：
  # 生成密钥对
  fzjjyz keygen -d ./keys -n mykey

  # 加密文件
  fzjjyz encrypt -i plaintext.txt -o encrypted.fzj -p keys/mykey_public.pem -s keys/mykey_dilithium_private.pem

  # 解密文件
  fzjjyz decrypt -i encrypted.fzj -o decrypted.txt -p keys/mykey_private.pem -s keys/mykey_dilithium_public.pem

  # 查看文件信息
  fzjjyz info -i encrypted.fzj

项目主页: https://codeberg.org/jiangfire/fzjjyz`,

	// 全局标志
	"flags.verbose": "启用详细输出",
	"flags.force":   "强制覆盖现有文件",

	// encrypt 命令
	"encrypt.short": "加密文件",
	"encrypt.long": `使用后量子混合加密算法加密文件。

加密流程：
  1. 读取原始文件
  2. 生成随机会话密钥
  3. Kyber768 + ECDH 密钥封装
  4. AES-256-GCM 加密数据
  5. Dilithium3 签名验证
  6. 构建加密文件头
  7. 写入加密文件

必需参数：
  --input, -i         输入文件路径
  --public-key, -p    Kyber+ECDH 公钥文件
  --sign-key, -s      Dilithium 私钥文件

示例：
  fzjjyz encrypt -i plaintext.txt -o encrypted.fzj -p public.pem -s dilithium_private.pem
  fzjjyz encrypt --input data.txt --public-key pub.pem --sign-key priv.pem --force`,
	"encrypt.flags.input":       "输入文件路径 (必需)",
	"encrypt.flags.output":      "输出文件路径 (可选，默认: input.fzj)",
	"encrypt.flags.public-key":  "Kyber+ECDH 公钥文件 (必需)",
	"encrypt.flags.sign-key":    "Dilithium 私钥文件 (必需)",
	"encrypt.flags.force":       "覆盖输出文件",
	"encrypt.flags.buffer-size": "缓冲区大小 (KB)，0=自动选择",
	"encrypt.flags.streaming":   "使用流式处理（大文件推荐）",

	// decrypt 命令
	"decrypt.short": "解密文件",
	"decrypt.long": `使用后量子混合加密算法解密文件。

解密流程：
  1. 解析文件头
  2. 验证文件格式
  3. Kyber768 + ECDH 密钥解封装
  4. AES-256-GCM 解密数据
  5. 验证 SHA256 哈希
  6. 验证 Dilithium 签名
  7. 写入原始文件

必需参数：
  --input, -i         加密文件路径
  --private-key, -p   Kyber+ECDH 私钥文件
  --verify-key, -s    Dilithium 公钥文件 (可选)

示例：
  fzjjyz decrypt -i encrypted.fzj -o decrypted.txt -p private.pem -s dilithium_public.pem
  fzjjyz decrypt --input data.fzj --private-key priv.pem --verify-key pub.pem --force`,
	"decrypt.flags.input":       "加密文件路径 (必需)",
	"decrypt.flags.output":      "输出文件路径 (可选，默认: 原文件名)",
	"decrypt.flags.private-key": "Kyber+ECDH 私钥文件 (必需)",
	"decrypt.flags.verify-key":  "Dilithium 公钥文件 (可选)",
	"decrypt.flags.force":       "覆盖输出文件",
	"decrypt.flags.buffer-size": "缓冲区大小 (KB)，0=自动选择",
	"decrypt.flags.streaming":   "使用流式处理（大文件推荐）",

	// encrypt-dir 命令
	"encrypt-dir.short": "加密文件夹",
	"encrypt-dir.long": `将整个文件夹打包成ZIP，然后使用后量子混合加密算法加密。

加密流程：
  1. 扫描源目录，递归获取所有文件
  2. 将目录结构打包成ZIP格式
  3. 读取ZIP数据到内存
  4. Kyber768 + ECDH 密钥封装
  5. AES-256-GCM 加密ZIP数据
  6. Dilithium3 签名验证
  7. 构建加密文件头
  8. 写入加密文件 (.fzj)

必需参数：
  --input, -i         源目录路径
  --output, -o        输出加密文件路径
  --public-key, -p    Kyber+ECDH 公钥文件
  --sign-key, -s      Dilithium 私钥文件

示例：
  fzjjyz encrypt-dir -i ./sensitive_data -o secure.fzj -p public.pem -s dilithium_private.pem
  fzjjyz encrypt-dir --input ./confidential --output backup.fzj --public-key pub.pem --sign-key priv.pem --force`,
	"encrypt-dir.flags.input":       "源目录路径 (必需)",
	"encrypt-dir.flags.output":      "输出加密文件路径 (必需)",
	"encrypt-dir.flags.public-key":  "Kyber+ECDH 公钥文件 (必需)",
	"encrypt-dir.flags.sign-key":    "Dilithium 私钥文件 (必需)",
	"encrypt-dir.flags.force":       "覆盖输出文件",
	"encrypt-dir.flags.buffer-size": "缓冲区大小 (KB)，0=自动选择",
	"encrypt-dir.flags.streaming":   "使用流式处理",

	// decrypt-dir 命令
	"decrypt-dir.short": "解密文件夹",
	"decrypt-dir.long": `解密加密的文件夹存档，并恢复原始目录结构。

解密流程：
  1. 解析加密文件头
  2. 验证文件格式
  3. Kyber768 + ECDH 密钥解封装
  4. AES-256-GCM 解密数据
  5. 验证 SHA256 哈希
  6. 验证 Dilithium 签名
  7. 解压ZIP到目标目录
  8. 恢复原始目录结构

必需参数：
  --input, -i         加密文件路径 (.fzj)
  --output, -o        输出目录路径
  --private-key, -p   Kyber+ECDH 私钥文件
  --verify-key, -s    Dilithium 公钥文件 (可选)

示例：
  fzjjyz decrypt-dir -i secure.fzj -o ./restored -p private.pem -s dilithium_public.pem
  fzjjyz decrypt-dir --input backup.fzj --output ./recovered --private-key priv.pem --verify-key pub.pem --force`,
	"decrypt-dir.flags.input":       "加密文件路径 (必需)",
	"decrypt-dir.flags.output":      "输出目录路径 (必需)",
	"decrypt-dir.flags.private-key": "Kyber+ECDH 私钥文件 (必需)",
	"decrypt-dir.flags.verify-key":  "Dilithium 公钥文件 (可选)",
	"decrypt-dir.flags.force":       "覆盖输出目录中的现有文件",
	"decrypt-dir.flags.buffer-size": "缓冲区大小 (KB)，0=自动选择",
	"decrypt-dir.flags.streaming":   "使用流式处理",

	// keygen 命令
	"keygen.short": "生成后量子密钥对",
	"keygen.long": `生成完整的密钥对组合，包括：
  • Kyber768 + ECDH 密钥对 (用于加密/解密)
  • Dilithium3 密钥对 (用于签名/验证)

生成的文件：
  {name}_public.pem          - Kyber+ECDH 公钥
  {name}_private.pem         - Kyber+ECDH 私钥 (0600权限)
  {name}_dilithium_public.pem  - Dilithium 公钥
  {name}_dilithium_private.pem - Dilithium 私钥 (0600权限)

示例：
  fzjjyz keygen -d ./keys -n mykey
  fzjjyz keygen --output-dir ./keys --name mykey --force`,
	"keygen.flags.output-dir": "输出目录",
	"keygen.flags.name":       "密钥名称前缀 (默认: 时间戳)",
	"keygen.flags.force":      "覆盖现有文件",

	// keymanage 命令
	"keymanage.short": "密钥管理工具",
	"keymanage.long": `管理加密密钥，支持导出、导入和验证操作。

可用操作:
  export    从私钥文件中提取并导出公钥
  import    导入密钥文件到指定目录
  verify    验证密钥对是否匹配

示例:
  # 导出公钥
  fzjjyz keymanage export --private-key private.pem --output public_extracted.pem

  # 验证密钥对
  fzjjyz keymanage verify --public-key public.pem --private-key private.pem

  # 导入密钥
  fzjjyz keymanage import --public-key pub.pem --private-key priv.pem --output-dir ./keys`,
	"keymanage.flags.action":      "操作类型: export/import/verify (必需)",
	"keymanage.flags.public-key":  "公钥文件路径",
	"keymanage.flags.private-key": "私钥文件路径",
	"keymanage.flags.output":      "输出文件路径 (用于export)",
	"keymanage.flags.output-dir":  "输出目录 (用于import)",

	// info 命令
	"info.short": "查看加密文件信息",
	"info.long": `解析并显示加密文件的详细信息，包括：
  • 文件名和原始大小
  • 加密时间戳
  • 使用的算法
  • 签名状态
  • 完整性验证

示例：
  fzjjyz info -i encrypted.fzj
  fzjjyz info --input data.fzj`,
	"info.flags.input": "加密文件路径 (必需)",

	// version 命令
	"version.short":       "显示版本信息",
	"version.long":        "显示 fzjjyz 的版本信息和构建详情",
	"version.info":        "版本信息",
	"version.label":       "版本:",
	"version.app_name":    "应用名称:",
	"version.description": "描述:",

	// 进度提示
	"progress.loading_keys":         "加载密钥...",
	"progress.encrypting":           "加密文件...",
	"progress.verifying":            "验证...",
	"progress.decrypting":           "解密文件...",
	"progress.packing":              "打包文件夹...",
	"progress.extracting":           "解压文件...",
	"progress.generating_kyber":     "生成 Kyber768 密钥...",
	"progress.generating_ecdh":      "生成 ECDH X25519 密钥...",
	"progress.generating_dilithium": "生成 Dilithium3 签名密钥...",
	"progress.saving_keys":          "保存密钥文件...",

	// 状态消息
	"status.done":                   "完成",
	"status.failed":                 "失败",
	"status.warning_no_sign_verify": "⚠️  警告: 未提供签名验证密钥，将跳过签名验证",
	"status.success_encrypt":        "✅ 加密成功！",
	"status.success_decrypt":        "✅ 解密成功！",
	"status.success_keygen":         "✅ 密钥对生成成功！",
	"status.success_export":         "✅ 公钥已导出到: %s",
	"status.success_import":         "✅ 密钥已导入到: %s",
	"status.success_verify":         "✅ 密钥对验证通过",
	"status.failed_verify":          "❌ 密钥对不匹配",
	"status.encrypting_file":        "加密文件: %s",
	"status.decrypting_file":        "解密文件: %s",
	"status.encrypting_dir":         "加密文件夹: %s",
	"status.decrypting_dir":         "解密文件夹: %s",
	"status.generating_keys":        "生成密钥对...",
	"status.public_key":             "公钥",
	"status.sign_key":               "签名密钥",
	"status.streaming_mode":         "流式处理",

	// 文件信息输出
	"file_info.header":            "📁 文件信息: %s",
	"file_info.basic":             "基本信息:",
	"file_info.encryption":        "加密信息:",
	"file_info.keys":              "密钥信息:",
	"file_info.integrity":         "完整性:",
	"file_info.verification":      "验证状态:",
	"file_info.original_file":     "原始文件: %s (%d bytes)",
	"file_info.encrypted_file":    "加密文件: %s (%d bytes)",
	"file_info.decrypted_file":    "解密文件: %s (%d bytes)",
	"file_info.compressed_rate":   "压缩率: %.1f%%",
	"file_info.timestamp":         "时间戳: %s",
	"file_info.algorithm":         "算法: %s (0x%02x)",
	"file_info.version":           "版本: 0x%04x",
	"file_info.magic":             "魔数: %c%c%c\\x%02x",
	"file_info.kyber":             "Kyber封装: %d bytes",
	"file_info.ecdh":              "ECDH公钥: %d bytes",
	"file_info.iv":                "IV/Nonce: %d bytes",
	"file_info.signature":         "签名: %d bytes",
	"file_info.hash":              "SHA256哈希: %x...",
	"file_info.signature_status":  "签名:",
	"file_info.data_integrity":    "数据完整性:",
	"file_info.exists":            "存在",
	"file_info.not_exists":        "不存在",
	"file_info.complete":          "完整",
	"file_info.incomplete":        "不完整",
	"file_info.original_filename": "原始文件名: %s",
	"file_info.file_count":        "文件数量: %d 个",
	"file_info.source_dir":        "源目录: %s",
	"file_info.output_dir":        "输出目录: %s",
	"file_info.zip_size":          "ZIP大小: %d bytes",
	"file_info.decrypted_size":    "解密大小: %d bytes",
	"file_info.buffer_size":       "缓冲区大小: %d KB",

	// 文件夹加密/解密信息
	"dir_info.encrypt_summary": `文件信息:
  源目录: %s
  文件数量: %d 个
  ZIP大小: %d bytes
  加密文件: %s (%d bytes)
  压缩率: %.1f%%`,
	"dir_info.decrypt_summary": `文件信息:
  加密文件: %s (%d bytes)
  解密大小: %d bytes
  文件数量: %d 个
  输出目录: %s
  原始文件名: %s
  时间戳: %s`,

	// 单文件加密/解密信息
	"file_info.encrypt_summary": `文件信息:
  原始文件: %s (%d bytes)
  加密文件: %s (%d bytes)
  压缩率: %.1f%%`,
	"file_info.decrypt_summary": `文件信息:
  加密文件: %s (%d bytes)
  解密文件: %s (%d bytes)
  原始文件名: %s
  时间戳: %s`,

	// 密钥生成信息
	"keygen_info.files": `生成的文件:
  • %s (公钥)
  • %s (私钥 - 0600权限)
  • %s (签名公钥)
  • %s (签名私钥 - 0600权限)`,

	// 密钥验证信息
	"keymanage_verify.kyber": "  Kyber:  %s",
	"keymanage_verify.ecdh":  "  ECDH:   %s",

	// 安全提示
	"security.warning":        "⚠️  安全提示:",
	"security.protect_keys":   "• 请妥善保管私钥文件",
	"security.no_sharing":     "• 不要将私钥分享给他人",
	"security.secure_storage": "• 建议使用安全的存储介质",

	// 打包/解压信息
	"archive.packed":    "完成 (大小: %d bytes, 文件数: %d)",
	"archive.decrypted": "完成 (大小: %d bytes)",

	// 错误信息 - 文件相关
	"error.file_not_exists":           "文件不存在: %s",
	"error.input_file_not_exists":     "输入文件不存在: %s",
	"error.encrypted_file_not_exists": "加密文件不存在: %s",
	"error.source_dir_not_exists":     "源目录不存在: %s",
	"error.input_not_dir":             "输入路径不是目录: %s",
	"error.output_not_dir":            "输出路径不是目录: %s",
	"error.output_file_exists":        "输出文件已存在: %s (使用 --force 覆盖)",
	"error.output_dir_not_empty":      "输出目录非空: %s (使用 --force 覆盖)",
	"error.cannot_create_dir":         "无法创建目录 %s: %v",
	"error.cannot_open_file":          "无法打开加密文件: %v",
	"error.cannot_read_file":          "无法读取文件: %v",
	"error.cannot_read_data":          "无法读取解密数据: %v",
	"error.cannot_open_temp":          "无法打开临时文件: %v",

	// 错误信息 - 密钥相关
	"error.key_not_found": "密钥文件不存在: %s",
	"error.key_invalid":   "密钥格式无效",
	"error.load_public_key_failed": `❌ 加载公钥失败: %v

提示:
  1. 请检查公钥文件路径是否正确: %s
  2. 确保公钥文件格式正确（PEM 格式）
  3. 检查文件权限（需可读）
  4. 如果是首次使用，请先生成密钥对: fzjjyz keygen`,
	"error.load_private_key_failed": `❌ 加载私钥失败: %v

提示:
  1. 请检查私钥文件路径是否正确: %s
  2. 确保私钥文件格式正确（PEM 格式）
  3. 检查文件权限（建议 0600）
  4. 私钥文件应仅由所有者读取
  5. 确保使用与加密时匹配的私钥`,
	"error.load_sign_key_failed": `❌ 加载签名私钥失败: %v

提示:
  1. 请检查 Dilithium 私钥文件路径是否正确: %s
  2. 确保私钥文件格式正确（PEM 格式）
  3. 检查文件权限（建议 0600）
  4. 私钥文件应仅由所有者读取
  5. 如果是首次使用，请先生成密钥对: fzjjyz keygen`,
	"error.load_verify_key_failed": `❌ 加载验证公钥失败: %v

提示:
  1. 请检查 Dilithium 公钥文件路径是否正确: %s
  2. 确保公钥文件格式正确（PEM 格式）
  3. 检查文件权限（需可读）
  4. 确保使用与加密时匹配的公钥
  5. 如果未提供签名密钥，可省略此参数（但无法验证签名）`,
	"error.keygen_kyber_failed":     "Kyber密钥生成失败: %v",
	"error.keygen_ecdh_failed":      "ECDH密钥生成失败: %v",
	"error.keygen_dilithium_failed": "Dilithium密钥生成失败: %v",
	"error.save_keys_failed":        "保存密钥文件失败: %v",
	"error.save_dilithium_failed":   "保存Dilithium密钥失败: %v",
	"error.export_key_failed":       "导出公钥失败: %v",
	"error.save_export_failed":      "保存公钥文件失败: %v",
	"error.import_keys_failed":      "导入密钥失败: %v",
	"error.verify_keys_failed":      "密钥对不匹配",

	// 错误信息 - 加密/解密相关
	"error.encrypt_failed": `❌ 加密失败: %v

可能原因:
  1. 文件权限不足（无法读取输入或写入输出）
  2. 内存不足（大文件需要更多内存）
  3. 密钥不匹配
  4. 输入文件在加密过程中被修改

建议:
  - 检查磁盘空间和文件权限
  - 对于超大文件，尝试使用 --buffer-size 调整缓冲区
  - 确保密钥正确匹配`,
	"error.decrypt_failed": `❌ 解密失败: %v

可能原因:
  1. 密钥不匹配（使用了错误的私钥）
  2. 文件已损坏或被篡改
  3. 文件格式不正确（不是 fzjjyz 加密文件）
  4. 签名验证失败（文件可能被篡改）
  5. 文件权限不足

安全提示:
  - 如果提示哈希不匹配，文件可能已被篡改，请勿使用
  - 如果提示签名无效，密钥可能不匹配或文件被修改
  - 建议始终提供签名验证密钥以确保文件完整性`,
	"error.pack_failed": `❌ 打包失败: %v

可能原因:
  1. 目录权限不足
  2. 包含不支持的文件类型（如符号链接）
  3. 磁盘空间不足`,
	"error.extract_failed": `❌ 解压失败: %v

可能原因:
  1. 输出目录权限不足
  2. 磁盘空间不足
  3. ZIP文件损坏`,
	"error.temp_file_failed":       "❌ 临时文件创建失败: %v",
	"error.parse_header_failed":    "文件头解析失败: %v",
	"error.validate_header_failed": "文件头验证失败: %v",

	// 错误信息 - 其他
	"error.unknown_action":         "未知操作: %s (支持: export, import, verify)",
	"error.missing_required_flags": "必须提供 %s",
	"error.missing_both_keys":      "必须提供 --public-key 和 --private-key",
	"error.nothing_to_do":          "没有可执行的操作",
}
