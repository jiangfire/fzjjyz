package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCLIIntegration CLI 集成测试.
//
//nolint:gocognit,gocyclo,funlen // 测试函数需要完整覆盖所有场景，复杂度和长度是必要的
func TestCLIIntegration(t *testing.T) {
	// 跳过在 CI 中的测试（如果需要）
	if testing.Short() {
		t.Skip("跳过 CLI 集成测试")
	}

	// 构建 CLI 可执行文件
	executable := buildCLI(t)
	defer func() {
		if err := os.Remove(executable); err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}()

	// 创建临时测试目录
	testDir, err := os.MkdirTemp("", "fzjjyz-test-*") // #nosec G301 - 测试环境使用临时目录
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}()

	// 测试数据
	testFile := filepath.Join(testDir, "test.txt")
	encryptedFile := filepath.Join(testDir, "test.txt.fzj")
	keyPrefix := "testkey"
	pubKey := filepath.Join(testDir, keyPrefix+"_public.pem")
	privKey := filepath.Join(testDir, keyPrefix+"_private.pem")
	dilithiumPubKey := filepath.Join(testDir, keyPrefix+"_dilithium_public.pem")
	dilithiumPrivKey := filepath.Join(testDir, keyPrefix+"_dilithium_private.pem")
	extractedPubKey := filepath.Join(testDir, "extracted_public.pem")

	// 测试步骤
	t.Run("1. 生成密钥对", func(t *testing.T) {
		// #nosec G204 - 测试环境执行命令
		cmd := exec.Command(executable, "keygen", "-d", testDir, "-n", keyPrefix)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("密钥生成失败: %v\n输出: %s", err, output)
		}

		// 验证密钥文件是否存在
		for _, file := range []string{pubKey, privKey, dilithiumPubKey, dilithiumPrivKey} {
			if _, err := os.Stat(file); os.IsNotExist(err) { // #nosec G304 - 测试环境使用临时文件路径
				t.Errorf("密钥文件未创建: %s", file)
			}
		}

		t.Log("✅ 密钥生成成功")
	})

	t.Run("2. 创建测试文件", func(t *testing.T) {
		content := "这是测试文件内容，用于测试加密和解密功能。\n"
		content += "时间戳: " + time.Now().Format(time.RFC3339) + "\n"
		content += "测试数据: " + strings.Repeat("ABC", 100) + "\n"

		// #nosec G306 - 测试环境使用标准权限
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("创建测试文件失败: %v", err)
		}

		t.Log("✅ 测试文件创建成功")
	})

	t.Run("3. 加密文件", func(t *testing.T) {
		cmd := exec.Command(executable, "encrypt",
			"-i", testFile,
			"-o", encryptedFile,
			"-p", pubKey,
			"-s", dilithiumPrivKey,
			"-v",
		) // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("加密失败: %v\n输出: %s", err, output)
		}

		// 验证加密文件是否存在
		if _, err := os.Stat(encryptedFile); os.IsNotExist(err) { // #nosec G304 - 测试环境使用临时文件路径
			t.Errorf("加密文件未创建: %s", encryptedFile)
		}

		t.Logf("✅ 加密成功\n输出: %s", output)
	})

	t.Run("4. 解密文件", func(t *testing.T) {
		decryptedFile := filepath.Join(testDir, "decrypted.txt")

		cmd := exec.Command(executable, "decrypt",
			"-i", encryptedFile,
			"-o", decryptedFile,
			"-p", privKey,
			"-s", dilithiumPubKey,
			"-v",
		) // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("解密失败: %v\n输出: %s", err, output)
		}

		// 验证解密文件是否存在
		if _, err := os.Stat(decryptedFile); os.IsNotExist(err) { // #nosec G304 - 测试环境使用临时文件路径
			t.Errorf("解密文件未创建: %s", decryptedFile)
		}

		// 验证解密内容是否与原文件一致
		original, err := os.ReadFile(testFile) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("读取原始文件失败: %v", err)
		}
		decrypted, err := os.ReadFile(decryptedFile) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("读取解密文件失败: %v", err)
		}
		if !bytes.Equal(original, decrypted) {
			t.Errorf("解密内容与原文件不一致")
		}

		t.Logf("✅ 解密成功\n输出: %s", output)
	})

	t.Run("5. 密钥管理 - 导出公钥", func(t *testing.T) {
		cmd := exec.Command(executable, "keymanage",
			"-a", "export",
			"-s", privKey,
			"-o", extractedPubKey,
		) // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("公钥导出失败: %v\n输出: %s", err, output)
		}

		// 验证导出的公钥文件是否存在
		if _, err := os.Stat(extractedPubKey); os.IsNotExist(err) { // #nosec G304 - 测试环境使用临时文件路径
			t.Errorf("导出的公钥文件未创建: %s", extractedPubKey)
		}

		// 验证导出的公钥与原公钥内容一致
		originalPub, err := os.ReadFile(pubKey) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("读取原始公钥失败: %v", err)
		}
		extractedPub, err := os.ReadFile(extractedPubKey) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("读取导出公钥失败: %v", err)
		}
		if !bytes.Equal(originalPub, extractedPub) {
			t.Errorf("导出的公钥与原公钥不一致")
		}

		t.Logf("✅ 公钥导出成功\n输出: %s", output)
	})

	t.Run("6. 密钥管理 - 验证密钥对", func(t *testing.T) {
		cmd := exec.Command(executable, "keymanage",
			"-a", "verify",
			"-p", pubKey,
			"-s", privKey,
		) // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("密钥验证失败: %v\n输出: %s", err, output)
		}

		if !strings.Contains(string(output), "✅") {
			t.Errorf("密钥验证未通过: %s", output)
		}

		t.Logf("✅ 密钥验证成功\n输出: %s", output)
	})

	t.Run("7. 密钥管理 - 导入密钥", func(t *testing.T) {
		importDir := filepath.Join(testDir, "imported")
		// #nosec G301 - 测试环境使用标准权限
		if err := os.MkdirAll(importDir, 0755); err != nil {
			t.Fatalf("创建导入目录失败: %v", err)
		}

		cmd := exec.Command(executable, "keymanage",
			"-a", "import",
			"-p", pubKey,
			"-s", privKey,
			"-d", importDir,
		) // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("密钥导入失败: %v\n输出: %s", err, output)
		}

		// 验证导入的文件是否存在
		importedPub := filepath.Join(importDir, filepath.Base(pubKey))
		importedPriv := filepath.Join(importDir, filepath.Base(privKey))
		if _, err := os.Stat(importedPub); os.IsNotExist(err) { // #nosec G304 - 测试环境使用临时文件路径
			t.Errorf("导入的公钥文件未创建: %s", importedPub)
		}
		if _, err := os.Stat(importedPriv); os.IsNotExist(err) { // #nosec G304 - 测试环境使用临时文件路径
			t.Errorf("导入的私钥文件未创建: %s", importedPriv)
		}

		t.Logf("✅ 密钥导入成功\n输出: %s", output)
	})

	t.Run("8. 版本信息", func(t *testing.T) {
		cmd := exec.Command(executable, "version") // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("获取版本信息失败: %v", err)
		}

		if !strings.Contains(string(output), "fzjjyz") {
			t.Errorf("版本输出不包含应用名称: %s", output)
		}

		t.Logf("✅ 版本信息获取成功\n输出: %s", output)
	})

	t.Run("9. 错误处理 - 错误密钥解密", func(t *testing.T) {
		// 创建错误的密钥
		wrongKeyDir, _ := os.MkdirTemp("", "wrong-key-*") // #nosec G301 - 测试环境使用临时目录
		defer func() {
			if err := os.RemoveAll(wrongKeyDir); err != nil {
				t.Logf("cleanup warning: %v", err)
			}
		}()

		wrongKeyPrefix := "wrongkey"
		cmd := exec.Command(executable, "keygen", "-d", wrongKeyDir, "-n", wrongKeyPrefix) // #nosec G204 - 测试环境执行命令
		if err := cmd.Run(); err != nil {
			t.Fatalf("生成错误密钥失败: %v", err)
		}

		wrongPrivKey := filepath.Join(wrongKeyDir, wrongKeyPrefix+"_private.pem")

		// 尝试用错误的密钥解密
		decryptedFile := filepath.Join(testDir, "wrong_decrypted.txt")
		cmd = exec.Command(executable, "decrypt",
			"-i", encryptedFile,
			"-o", decryptedFile,
			"-p", wrongPrivKey,
		) // #nosec G204 - 测试环境执行命令
		output, err := cmd.CombinedOutput()

		// 应该失败
		if err == nil {
			t.Errorf("使用错误密钥解密应该失败，但没有返回错误")
		}

		// 验证解密文件不存在（或内容错误）
		if _, err := os.Stat(decryptedFile); err == nil { // #nosec G304 - 测试环境使用临时文件路径
			// 文件存在，检查内容是否应该无效
			content, err := os.ReadFile(decryptedFile) // #nosec G304 - 测试环境使用临时文件路径
			if err != nil {
				t.Fatalf("读取解密文件失败: %v", err)
			}
			if len(content) > 0 && !bytes.Equal(content, []byte("这是测试文件内容")) {
				t.Logf("✅ 错误密钥解密失败（符合预期）: %s", output)
			}
		} else {
			t.Logf("✅ 错误密钥解密失败（符合预期）: %s", output)
		}
	})
}

// buildCLI 构建 CLI 可执行文件.
func buildCLI(t *testing.T) string {
	// 创建临时可执行文件路径
	tmpFile, err := os.CreateTemp("", "fzjjyz-test-*.exe")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	executable := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("close temp file failed: %v", err)
	}

	// 获取项目根目录（当前工作目录的父目录）
	projectRoot := "C:\\Users\\yimo\\Codes\\fzjjyz"
	if cwd, err := os.Getwd(); err == nil {
		// 如果测试在 cmd/fzjjyz 目录下运行，需要回到项目根目录
		if filepath.Base(cwd) == "fzjjyz" && filepath.Base(filepath.Dir(cwd)) == "cmd" {
			projectRoot = filepath.Dir(filepath.Dir(cwd))
		}
	}

	// 构建命令
	cmd := exec.Command("go", "build", "-o", executable, "./cmd/fzjjyz") // #nosec G204 - 测试环境执行命令
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("构建 CLI 失败: %v\n输出: %s", err, output)
	}

	t.Logf("CLI 构建成功: %s", executable)
	return executable
}

// TestCLIHelp 测试帮助信息.
func TestCLIHelp(t *testing.T) {
	executable := buildCLI(t)
	defer func() {
		if err := os.Remove(executable); err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}()

	tests := []struct {
		name    string
		command []string
	}{
		{"根命令帮助", []string{"--help"}},
		{"加密帮助", []string{"encrypt", "--help"}},
		{"解密帮助", []string{"decrypt", "--help"}},
		{"密钥生成帮助", []string{"keygen", "--help"}},
		{"密钥管理帮助", []string{"keymanage", "--help"}},
		{"版本信息", []string{"version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(executable, tt.command...) // #nosec G204 - 测试环境执行命令
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("执行 %v 失败: %v", tt.command, err)
			}

			outputStr := string(output)
			if len(outputStr) < 10 {
				t.Errorf("输出太短: %s", outputStr)
			}

			t.Logf("命令 %v 输出: %s", tt.command, outputStr[:minInt(200, len(outputStr))])
		})
	}
}

// TestCLIBenchmark 简单的性能测试.
//
//nolint:funlen
func TestCLIBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}

	executable := buildCLI(t)
	defer func() {
		if err := os.Remove(executable); err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}()

	// 创建临时目录
	testDir, err := os.MkdirTemp("", "fzjjyz-bench-*") // #nosec G301 - 测试环境使用临时目录
	if err != nil {
		t.Fatalf("创建测试目录失败: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}()

	// 生成密钥
	keyPrefix := "bench"
	// #nosec G204 - 测试环境执行命令
	cmd := exec.Command(executable, "keygen", "-d", testDir, "-n", keyPrefix)
	if err := cmd.Run(); err != nil {
		t.Fatalf("密钥生成失败: %v", err)
	}

	pubKey := filepath.Join(testDir, keyPrefix+"_public.pem")
	privKey := filepath.Join(testDir, keyPrefix+"_private.pem")
	dilithiumPubKey := filepath.Join(testDir, keyPrefix+"_dilithium_public.pem")
	dilithiumPrivKey := filepath.Join(testDir, keyPrefix+"_dilithium_private.pem")

	// 创建测试文件（1MB）
	testFile := filepath.Join(testDir, "large.txt")
	largeContent := strings.Repeat("LARGE_CONTENT_", 70000) // ~1MB
	// #nosec G306 - 测试环境使用标准权限
	if err := os.WriteFile(testFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("创建大文件失败: %v", err)
	}

	// 测量加密时间
	t.Run("加密性能", func(t *testing.T) {
		encryptedFile := filepath.Join(testDir, "large.fzj")
		start := time.Now()

		cmd := exec.Command(executable, "encrypt",
			"-i", testFile,
			"-o", encryptedFile,
			"-p", pubKey,
			"-s", dilithiumPrivKey,
		) // #nosec G204 - 测试环境执行命令
		if err := cmd.Run(); err != nil {
			t.Fatalf("加密失败: %v", err)
		}

		elapsed := time.Since(start)
		t.Logf("加密 1MB 文件耗时: %v", elapsed)

		// 验证加密文件大小
		info, err := os.Stat(encryptedFile) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("获取加密文件信息失败: %v", err)
		}
		t.Logf("加密后文件大小: %d bytes", info.Size())
	})

	// 测量解密时间
	t.Run("解密性能", func(t *testing.T) {
		encryptedFile := filepath.Join(testDir, "large.fzj")
		decryptedFile := filepath.Join(testDir, "large_decrypted.txt")

		// 先加密（使用 --force 覆盖现有文件）
		cmd := exec.Command(executable, "encrypt",
			"-i", testFile,
			"-o", encryptedFile,
			"-p", pubKey,
			"-s", dilithiumPrivKey,
			"--force",
		) // #nosec G204 - 测试环境执行命令
		if err := cmd.Run(); err != nil {
			t.Fatalf("加密失败: %v", err)
		}

		// 测量解密
		start := time.Now()
		cmd = exec.Command(executable, "decrypt",
			"-i", encryptedFile,
			"-o", decryptedFile,
			"-p", privKey,
			"-s", dilithiumPubKey,
			"--force",
		) // #nosec G204 - 测试环境执行命令
		if err := cmd.Run(); err != nil {
			t.Fatalf("解密失败: %v", err)
		}

		elapsed := time.Since(start)
		t.Logf("解密 1MB 文件耗时: %v", elapsed)

		// 验证解密内容
		original, err := os.ReadFile(testFile) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("读取原始文件失败: %v", err)
		}
		decrypted, err := os.ReadFile(decryptedFile) // #nosec G304 - 测试环境使用临时文件路径
		if err != nil {
			t.Fatalf("读取解密文件失败: %v", err)
		}
		if !bytes.Equal(original, decrypted) {
			t.Errorf("解密内容不匹配")
		}
	})
}

// 辅助函数.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
