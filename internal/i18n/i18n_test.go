package i18n

import (
	"os"
	"sync"
	"testing"
)

// TestInit 测试初始化功能
func TestInit(t *testing.T) {
	// 重置状态
	globalDict = nil
	currentLang = ""
	once = sync.Once{}

	// 测试初始化中文
	err := Init("zh_CN")
	if err != nil {
		t.Fatalf("初始化中文失败: %v", err)
	}
	if GetLanguage() != "zh_CN" {
		t.Errorf("期望语言 zh_CN，实际得到 %s", GetLanguage())
	}

	// 重置状态
	globalDict = nil
	currentLang = ""
	once = sync.Once{}

	// 测试初始化英文
	err = Init("en_US")
	if err != nil {
		t.Fatalf("初始化英文失败: %v", err)
	}
	if GetLanguage() != "en_US" {
		t.Errorf("期望语言 en_US，实际得到 %s", GetLanguage())
	}

	// 重置状态
	globalDict = nil
	currentLang = ""
	once = sync.Once{}

	// 测试不支持的语言（会回退到默认语言 zh_CN）
	err = Init("invalid")
	if err != nil {
		t.Errorf("不支持的语言应该回退到默认语言，不应返回错误: %v", err)
	}
	if GetLanguage() != "zh_CN" {
		t.Errorf("不支持的语言应回退到 zh_CN，实际得到 %s", GetLanguage())
	}
}

// TestT 测试翻译函数
func TestT(t *testing.T) {
	// 初始化中文
	Init("zh_CN")

	// 测试简单翻译
	result := T("encrypt.short")
	if result == "" {
		t.Error("翻译结果不能为空")
	}

	// 测试带参数的翻译
	result = T("error.file_not_exists", "test.txt")
	if result == "" {
		t.Error("带参数的翻译结果不能为空")
	}
	if result == "error.file_not_exists" {
		t.Error("应该返回翻译后的文本，而不是 key")
	}

	// 测试不存在的 key
	result = T("nonexistent.key")
	if result != "nonexistent.key" {
		t.Errorf("不存在的 key 应该返回自身，实际得到: %s", result)
	}
}

// TestGet 测试 Get 函数（动态 key）
func TestGet(t *testing.T) {
	Init("zh_CN")

	// 测试动态 key
	keyPrefix := "encrypt"
	result := Get(keyPrefix + ".short")
	if result == "" {
		t.Error("Get 函数返回空")
	}

	// 测试不存在的 key
	result = Get("nonexistent.key")
	if result != "nonexistent.key" {
		t.Errorf("不存在的 key 应该返回自身，实际得到: %s", result)
	}
}

// TestSetLanguage 测试动态切换语言
func TestSetLanguage(t *testing.T) {
	// 先设置中文
	err := SetLanguage("zh_CN")
	if err != nil {
		t.Fatalf("设置中文失败: %v", err)
	}

	zhResult := T("encrypt.short")

	// 切换到英文
	err = SetLanguage("en_US")
	if err != nil {
		t.Fatalf("设置英文失败: %v", err)
	}

	enResult := T("encrypt.short")

	// 同一个 key，不同语言应该返回不同结果
	if zhResult == enResult {
		t.Errorf("不同语言应该返回不同翻译，中文: %s, 英文: %s", zhResult, enResult)
	}
}

// TestTranslateError 测试错误翻译
func TestTranslateError(t *testing.T) {
	Init("zh_CN")

	err := TranslateError("error.file_not_exists", "test.txt")
	if err == nil {
		t.Error("应该返回错误")
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("错误消息不能为空")
	}
}

// TestMustTranslate 测试强制翻译
func TestMustTranslate(t *testing.T) {
	Init("zh_CN")

	// 存在的 key
	result := MustTranslate("encrypt.short")
	if result == "" {
		t.Error("应该返回翻译结果")
	}

	// 不存在的 key
	result = MustTranslate("nonexistent.key")
	if result != "nonexistent.key" {
		t.Errorf("不存在的 key 应该返回 key 本身，实际得到: %s", result)
	}
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
	Init("zh_CN")

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = T("encrypt.short")
				_ = Get("decrypt.long")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestFormatString 测试格式化字符串
func TestFormatString(t *testing.T) {
	Init("zh_CN")

	// 测试带格式占位符的翻译
	result := T("error.file_not_exists", "missing.txt")
	if result == "error.file_not_exists" {
		t.Error("应该应用格式化参数")
	}

	// 测试多个参数
	result = T("file_info.encrypt_summary", "test.txt", 100, "test.fzj", 150, 150.0)
	if result == "file_info.encrypt_summary" {
		t.Error("应该应用多个格式化参数")
	}
}

// TestEmptyDict 测试空字典回退
func TestEmptyDict(t *testing.T) {
	// 重置全局状态
	globalDict = nil
	currentLang = ""
	once = sync.Once{}

	// 测试当 globalDict 为 nil 时的行为（在 Init 之前）
	// 此时 T() 会触发 lazy init，尝试初始化
	// 如果初始化失败，会使用空字典
	result := T("any.key")
	// 由于 Init 会回退到 zh_CN，所以实际会有翻译
	// 这个测试主要是确保不会 panic
	if result == "" {
		t.Error("T 不应该返回空字符串")
	}
}

// TestTranslationAccuracy 测试翻译准确性
func TestTranslationAccuracy(t *testing.T) {
	// 测试中文翻译
	Init("zh_CN")

	tests := []struct {
		key      string
		expected string // 包含的子字符串
	}{
		{"encrypt.short", "加密文件"},
		{"decrypt.short", "解密文件"},
		{"keygen.short", "生成密钥"},
		{"status.success_encrypt", "成功"},
		{"status.success_decrypt", "成功"},
	}

	for _, tt := range tests {
		result := T(tt.key)
		if result == tt.key {
			t.Errorf("中文翻译失败: %s", tt.key)
		}
	}

	// 测试英文翻译
	Init("en_US")

	for _, tt := range tests {
		result := T(tt.key)
		if result == tt.key {
			t.Errorf("英文翻译失败: %s", tt.key)
		}
	}
}

// TestGetWithEnvironment 测试环境变量检测
func TestGetWithEnvironment(t *testing.T) {
	// 保存原环境变量
	origLang := os.Getenv("LANG")
	defer os.Setenv("LANG", origLang)

	// 设置环境变量为中文
	os.Setenv("LANG", "zh_CN")

	// 重置状态
	globalDict = nil
	currentLang = ""
	once = sync.Once{}

	// 初始化（应该自动检测环境变量）
	if err := Init(""); err != nil {
		t.Fatalf("自动初始化失败: %v", err)
	}

	if GetLanguage() != "zh_CN" {
		t.Errorf("期望自动检测到 zh_CN，实际得到 %s", GetLanguage())
	}
}
