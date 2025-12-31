package i18n

import (
	"fmt"
	"os"
	"sync"
)

// TranslationDict 翻译字典接口.
type TranslationDict interface {
	Get(key string) string
}

// 全局变量.
var (
	globalDict  TranslationDict
	currentLang string
	once        sync.Once
	mu          sync.RWMutex
)

// 默认语言.
const defaultLang = "zh_CN"

// Init 初始化国际化系统
// lang: 语言代码（如 "zh_CN", "en_US"），空字符串表示从环境变量读取
func Init(lang string) error {
	mu.Lock()
	defer mu.Unlock()

	// 如果未指定语言，从环境变量读取
	if lang == "" {
		lang = os.Getenv("LANG")
		if lang == "" {
			lang = defaultLang
		}
	}

	// 加载对应的语言字典
	dict, err := loadLanguage(lang)
	if err != nil {
		// 回退到默认语言
		if lang != defaultLang {
			dict, err = loadLanguage(defaultLang)
			if err != nil {
				return fmt.Errorf("无法加载默认语言 %s: %w", defaultLang, err)
			}
			currentLang = defaultLang
			globalDict = dict
			return nil
		}
		return err
	}

	currentLang = lang
	globalDict = dict
	return nil
}

// T 翻译主函数
// 用法: i18n.T("encrypt.short") 或 i18n.T("error.file_not_exists", "test.txt").
func T(key string, args ...interface{}) string {
	once.Do(func() {
		if globalDict == nil {
			if err := Init(""); err != nil {
				// 初始化失败，使用空字典
				globalDict = &emptyDict{}
			}
		}
	})

	mu.RLock()
	defer mu.RUnlock()

	// 先查找翻译
	var translation string
	if globalDict != nil {
		translation = globalDict.Get(key)
	}

	if translation == "" {
		translation = key
	}

	// 如果没有参数，直接返回翻译结果
	if len(args) == 0 {
		return translation
	}

	// 有参数时，使用翻译结果作为格式字符串
	return fmt.Sprintf(translation, args...)
}

// Get 翻译函数（无格式化参数），用于动态 key
// 用法: i18n.Get(keyPrefix + ".short").
func Get(key string) string {
	once.Do(func() {
		if globalDict == nil {
			if err := Init(""); err != nil {
				// 初始化失败，使用空字典
				globalDict = &emptyDict{}
			}
		}
	})

	mu.RLock()
	defer mu.RUnlock()

	if globalDict != nil {
		if translation := globalDict.Get(key); translation != "" {
			return translation
		}
	}
	return key
}

// SetLanguage 动态切换语言.
func SetLanguage(lang string) error {
	return Init(lang)
}

// GetLanguage 获取当前语言.
func GetLanguage() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// loadLanguage 加载指定语言的字典.
func loadLanguage(lang string) (TranslationDict, error) {
	switch lang {
	case "zh_CN", "zh-CN", "zh":
		return &zhCN{}, nil
	case "en_US", "en-US", "en":
		return &enUS{}, nil
	default:
		// 尝试加载，如果失败则返回错误
		return nil, fmt.Errorf("不支持的语言: %s", lang)
	}
}

// emptyDict 空字典，用于错误回退.
type emptyDict struct{}

func (e *emptyDict) Get(_ string) string {
	return ""
}
