// Package zjcrypto 提供缓冲池功能
package zjcrypto

import "sync"

const (
	// DefaultBufferSize 默认缓冲区大小：64KB.
	DefaultBufferSize = 64 * 1024

	// MaxBufferSize 最大缓冲区大小：1MB.
	MaxBufferSize = 1024 * 1024

	// MinBufferSize 最小缓冲区大小：4KB.
	MinBufferSize = 4 * 1024
)

// BufferPool 缓冲区池，用于减少 GC 压力.
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool 创建新的缓冲区池.
func NewBufferPool(bufferSize int) *BufferPool {
	if bufferSize < MinBufferSize {
		bufferSize = MinBufferSize
	}
	if bufferSize > MaxBufferSize {
		bufferSize = MaxBufferSize
	}

	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, bufferSize)
			},
		},
	}
}

// Get 从池中获取一个缓冲区.
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put 将缓冲区归还到池中.
func (bp *BufferPool) Put(b []byte) {
	// 只有容量大于 0 的缓冲区才放回池中
	if cap(b) > 0 {
		// 重置切片长度，方便复用
		b = b[:cap(b)]
		//nolint:staticcheck // SA6002: sync.Pool.Put 接受 interface{}，[]byte 在 Go 1.18+ 可用，用于减少 GC 压力
		bp.pool.Put(b)
	}
}

// OptimalBufferSize 根据文件大小推荐缓冲区大小.
func OptimalBufferSize(fileSize int64) int {
	switch {
	case fileSize < 10*1024*1024: // < 10MB
		return 64 * 1024
	case fileSize < 100*1024*1024: // < 100MB
		return 256 * 1024
	case fileSize < 1024*1024*1024: // < 1GB
		return 512 * 1024
	default: // >= 1GB
		return 1024 * 1024
	}
}
