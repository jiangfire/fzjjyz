### Security
- 修复签名可剥离绕过：提供验签公钥时，密文必须包含有效签名。
- 收敛 `decrypt` 默认输出路径：仅使用头部文件名的 basename，防止路径穿越写入。
- 解密输出文件权限从 `0644` 调整为 `0600`。
- `decrypt-dir` 使用随机临时 ZIP 文件并在读取后尽快删除，降低明文驻留风险。

### Fixed
- 实现 `keymanage -a cache-info`，与 README 和使用文档保持一致。
- 新增签名剥离攻击、默认输出路径清洗、CLI cache-info 的回归测试。

### Notes
- 常规测试 `go test ./...` 已通过。
- 本地 `go test -race` 在当前工具链环境异常（`runtime/race` 包不可用），需在完整 race 工具链环境下执行。
