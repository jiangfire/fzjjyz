package zjcrypto

import (
	"fmt"
	"io"
	"sync"
)

// StreamProcessor 流式处理器接口.
type StreamProcessor interface {
	Process(input io.Reader, output io.Writer) error
}

// MultiWriter 支持多目标写入的包装器.
type MultiWriter struct {
	writers []io.Writer
}

// NewMultiWriter 创建多目标写入器.
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

// Write 实现 io.Writer 接口，数据同时写入所有目标.
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, fmt.Errorf("write to multiwriter target: %w", err)
		}
	}
	return len(p), nil
}

// PipeStream 管道流处理器.
type PipeStream struct {
	reader io.Reader
	writer io.Writer
	errCh  chan error
	once   sync.Once
}

// NewPipeStream 创建管道流.
func NewPipeStream() *PipeStream {
	reader, writer := io.Pipe()
	return &PipeStream{
		reader: reader,
		writer: writer,
		errCh:  make(chan error, 1),
	}
}

// GetReader 获取读取端.
func (ps *PipeStream) GetReader() io.Reader {
	return ps.reader
}

// GetWriter 获取写入端.
func (ps *PipeStream) GetWriter() io.Writer {
	return ps.writer
}

// Close 关闭管道.
func (ps *PipeStream) Close() error {
	if closer, ok := ps.writer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			return fmt.Errorf("close pipe writer: %w", err)
		}
	}
	return nil
}

// SetError 设置错误（只设置一次）.
func (ps *PipeStream) SetError(err error) {
	ps.once.Do(func() {
		if err != nil {
			ps.errCh <- err
		}
		close(ps.errCh)
	})
}

// GetError 获取错误.
func (ps *PipeStream) GetError() error {
	select {
	case err := <-ps.errCh:
		return err
	default:
		return nil
	}
}

// CopyWithHash 同时复制数据并计算哈希.
func CopyWithHash(dst io.Writer, src io.Reader, hasher io.Writer) (written int64, hash [32]byte, err error) {
	// 使用 MultiWriter 同时写入目标和哈希计算器
	multiWriter := NewMultiWriter(dst, hasher)

	written, err = io.Copy(multiWriter, src)
	if err != nil {
		return written, hash, fmt.Errorf("copy with hash: %w", err)
	}

	// 如果是 StreamingHash，获取最终结果
	if sh, ok := hasher.(*StreamingHash); ok {
		hash = sh.Sum()
	}

	return written, hash, nil
}
