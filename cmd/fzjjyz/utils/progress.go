// Package utils 提供工具函数（LOD原则：隔离输出细节）
package utils

import (
	"fmt"
	"io"
	"sync"
	"time"

	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
)

// ProgressBar 进度条结构.
type ProgressBar struct {
	Total      int64
	current    int64
	startTime  time.Time
	lastUpdate time.Time
	lock       sync.Mutex
	prefix     string
}

// NewProgressBar 创建新的进度条.
func NewProgressBar(total int64, prefix string) *ProgressBar {
	return &ProgressBar{
		Total:     total,
		startTime: time.Now(),
		prefix:    prefix,
	}
}

// Add 增加进度.
func (pb *ProgressBar) Add(n int64) {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	pb.current += n
	pb.lastUpdate = time.Now()
	pb.render()
}

// Set 设置当前进度.
func (pb *ProgressBar) Set(current int64) {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	pb.current = current
	pb.lastUpdate = time.Now()
	pb.render()
}

// Complete 完成进度条.
func (pb *ProgressBar) Complete() {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	pb.current = pb.Total
	pb.render()
	fmt.Println() // 换行
}

// render 渲染进度条.
func (pb *ProgressBar) render() {
	if pb.Total == 0 {
		return
	}

	percent := float64(pb.current) / float64(pb.Total) * 100

	// 进度条宽度
	barWidth := 40
	filled := int(float64(barWidth) * percent / 100)

	// 创建进度条字符串
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	// 计算速度
	elapsed := time.Since(pb.startTime).Seconds()
	var speed float64
	if elapsed > 0 {
		speed = float64(pb.current) / elapsed
	}

	// 估算剩余时间
	var eta string
	if speed > 0 && pb.current < pb.Total {
		remaining := float64(pb.Total - pb.current)
		etaSeconds := remaining / speed
		if etaSeconds < 60 {
			eta = fmt.Sprintf("%.0fs", etaSeconds)
		} else {
			eta = fmt.Sprintf("%.0fm", etaSeconds/60)
		}
	} else {
		eta = "--"
	}

	// 格式化输出
	fmt.Printf("\r%s [%s] %.1f%% (%d/%d) 速度: %.1f KB/s  ETA: %s",
		pb.prefix,
		bar,
		percent,
		pb.current,
		pb.Total,
		speed/1024,
		eta,
	)
}

// ProgressReader 包装 Reader 以显示进度.
type ProgressReader struct {
	reader io.Reader
	bar    *ProgressBar
}

// NewProgressReader 创建进度读取器.
func NewProgressReader(reader io.Reader, total int64, prefix string) *ProgressReader {
	return &ProgressReader{
		reader: reader,
		bar:    NewProgressBar(total, prefix),
	}
}

// Read 实现 io.Reader 接口.
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	if n > 0 {
		pr.bar.Add(int64(n))
	}
	return n, err //nolint:wrapcheck
}

// Close 完成进度.
func (pr *ProgressReader) Close() {
	pr.bar.Complete()
}

// ProgressWriter 包装 Writer 以显示进度.
type ProgressWriter struct {
	writer io.Writer
	bar    *ProgressBar
}

// NewProgressWriter 创建进度写入器.
func NewProgressWriter(writer io.Writer, total int64, prefix string) *ProgressWriter {
	return &ProgressWriter{
		writer: writer,
		bar:    NewProgressBar(total, prefix),
	}
}

// Write 实现 io.Writer 接口.
func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.writer.Write(p)
	if n > 0 {
		pw.bar.Add(int64(n))
	}
	return n, err //nolint:wrapcheck
}

// Close 完成进度.
func (pw *ProgressWriter) Close() {
	pw.bar.Complete()
}

// ProgressReporter reports progress with step-by-step updates.
type ProgressReporter struct {
	verbose bool
	step    int
	total   int
}

// NewProgressReporter creates a new progress reporter.
func NewProgressReporter(total int, verbose bool) *ProgressReporter {
	return &ProgressReporter{verbose: verbose, total: total}
}

// Step reports a step (eliminates 20+ repetitions).
func (p *ProgressReporter) Step(key string, _ ...interface{}) {
	p.step++
	fmt.Printf("[%d/%d] %s ", p.step, p.total, i18n.T(key))
}

// Done reports step completion.
func (p *ProgressReporter) Done() {
	fmt.Println(i18n.T("status.done"))
}

// Failed reports step failure.
func (p *ProgressReporter) Failed() {
	fmt.Println(i18n.T("status.failed"))
}

// Info displays information (eliminates verbose repetition).
func (p *ProgressReporter) Info(key string, value interface{}) {
	if p.verbose {
		fmt.Printf("  %s: %v\n", i18n.T(key), value)
	}
}

// InfoString displays string information.
func (p *ProgressReporter) InfoString(key string, value string) {
	if p.verbose {
		fmt.Printf("  %s: %s\n", i18n.T(key), value)
	}
}

// InfoBool displays boolean information.
func (p *ProgressReporter) InfoBool(key string, value bool) {
	if p.verbose {
		fmt.Printf("  %s: %v\n", i18n.T(key), value)
	}
}

// Warning displays a warning.
func (p *ProgressReporter) Warning(key string) {
	fmt.Println(i18n.T(key))
}

// Summary displays a summary.
func (p *ProgressReporter) Summary(title string, args ...interface{}) {
	fmt.Printf("\n%s\n\n", fmt.Sprintf(i18n.T(title), args...))
}

// PrintStatus prints status directly (compatible with existing code).
func PrintStatus(key string, args ...interface{}) {
	fmt.Printf(i18n.T(key)+"\n", args...)
}

// PrintVerbose prints in verbose mode (requires external verbose flag).
func PrintVerbose(verbose bool, key string, args ...interface{}) {
	if verbose {
		fmt.Printf(i18n.T(key)+"\n", args...)
	}
}
