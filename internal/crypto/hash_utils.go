package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// HashFile 流式计算文件的 SHA256 哈希值
// 使用 io.Copy 避免一次性读取整个文件到内存.
func HashFile(path string) ([32]byte, error) {
	var result [32]byte

	file, err := os.Open(path)
	if err != nil {
		return result, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return result, fmt.Errorf("hash file: %w", err)
	}

	sum := hasher.Sum(nil)
	copy(result[:], sum)
	return result, nil
}

// HashReader 流式计算 Reader 的 SHA256 哈希值.
func HashReader(r io.Reader) ([32]byte, error) {
	var result [32]byte

	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return result, fmt.Errorf("hash reader: %w", err)
	}

	sum := hasher.Sum(nil)
	copy(result[:], sum)
	return result, nil
}

// StreamingHash 流式哈希计算器，支持 Write 接口.
type StreamingHash struct {
	hash   io.Writer
	result [32]byte
	summed bool
}

// NewStreamingHash 创建新的流式哈希计算器.
func NewStreamingHash() *StreamingHash {
	return &StreamingHash{
		hash: sha256.New(),
	}
}

// Write 实现 io.Writer 接口.
func (sh *StreamingHash) Write(p []byte) (n int, err error) {
	n, err = sh.hash.Write(p)
	if err != nil {
		return n, fmt.Errorf("hash write failed: %w", err)
	}
	return n, nil
}

// Sum 获取最终的哈希值（只计算一次）.
func (sh *StreamingHash) Sum() [32]byte {
	if !sh.summed {
		if h, ok := sh.hash.(interface{ Sum([]byte) []byte }); ok {
			sum := h.Sum(nil)
			copy(sh.result[:], sum)
			sh.summed = true
		}
	}
	return sh.result
}
