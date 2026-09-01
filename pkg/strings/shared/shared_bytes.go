package shared

import (
	"bytes"
	"io"
	"sync"

	"github.com/gongt/go/internal/myenv"
)

// 可以并发访问的[bytes.Buffer]，读写不能同时进行
type Buffer struct {
	buff bytes.Buffer

	mu sync.RWMutex
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{buff: *bytes.NewBuffer(buf)}
}
func NewBufferString(s string) *Buffer {
	return &Buffer{buff: *bytes.NewBufferString(s)}
}

// Lock 锁定并运行指定的临界区函数
func (b *Buffer) Lock(critical func(buff *bytes.Buffer)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	critical(&b.buff)
}

// Lock1 锁定并运行指定的临界区函数
func (b *Buffer) Lock1[T1 any](critical func(*bytes.Buffer) T1) T1 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return critical(&b.buff)
}

// Lock2 锁定并运行指定的临界区函数
func (b *Buffer) Lock2[T1 any, T2 any](critical func(*bytes.Buffer) (T1, T2)) (T1, T2) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return critical(&b.buff)
}

// RLock 锁定并运行指定的临界区函数
func (b *Buffer) RLock(critical func(buff *bytes.Buffer)) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	critical(&b.buff)
}

// RLock1 锁定并运行指定的临界区函数
func (b *Buffer) RLock1[T1 any](critical func(*bytes.Buffer) T1) T1 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return critical(&b.buff)
}

// RLock2 锁定并运行指定的临界区函数
func (b *Buffer) RLock2[T1 any, T2 any](critical func(*bytes.Buffer) (T1, T2)) (T1, T2) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return critical(&b.buff)
}

// Bytes 返回缓冲区中的字节切片，复制整个buffer
func (b *Buffer) Bytes() []byte {
	b.mu.RLock()
	defer b.mu.RUnlock()
	r := make([]byte, b.buff.Len())
	copy(r, b.buff.Bytes())
	return r
}

// Appender 扩展缓冲区的可用字节切片，并通过use回调提供给调用者
//
// use 接收的切片长度恰好=cap
//
// 返回实际写入的字节数和可能的错误
func (b *Buffer) Appender(cap int, use func([]byte) (int, error)) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buff.Grow(cap)
	buf := b.buff.AvailableBuffer()
	n, err := use(buf)
	if n > 0 {
		b.buff.Write(buf[:n])
		if n != cap && myenv.IsDebug {
			for i := n; i < cap; i++ {
				buf[i] = 0xEE // 填充未使用的字节以便调试
			}
		}
	} else if n < 0 {
		panic("shared.Buffer.Appender: Callback returned negative count")
	}
	return n, err
}

// String 返回缓冲区的内容作为字符串
func (b *Buffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.String()
}

// Cap 返回缓冲区的容量
func (b *Buffer) Cap() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.Cap()
}

// Len 返回缓冲区中未读取的字节数
func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.Len()
}

// WriteTo 将缓冲区的内容写入到指定的io.Writer中
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.WriteTo(w)
}

// Read 从缓冲区中复制数据到提供的字节切片中
func (b *Buffer) Read(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.Read(p)
}

// ReadBytes 从缓冲区中读取直到遇到指定的分隔符，并返回包含分隔符的字节切片
//
// 如果在缓冲区中未找到分隔符，则返回缓冲区的所有数据和EOF。
func (b *Buffer) ReadBytes(delim byte) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadBytes(delim)
}

// ReadString 从缓冲区中读取直到遇到指定的分隔符，并返回包含分隔符的字符串
//
// 如果在缓冲区中未找到分隔符，则返回缓冲区的所有数据和EOF。
func (b *Buffer) ReadString(delim byte) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadString(delim)
}

// Truncate 将缓冲区的长度截断为 n
func (b *Buffer) Truncate(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buff.Truncate(n)
}

// Reset 重置缓冲区，使其长度为 0
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buff.Reset()
}

// Write 将提供的字节切片写入缓冲区
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.Write(p)
}

// WriteString 将提供的字符串写入缓冲区
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.WriteString(s)
}

// ReadFrom 从提供的 io.Reader 中读取数据并写入缓冲区
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadFrom(r)
}

// WriteByte 将一个字节写入缓冲区
func (b *Buffer) WriteByte(c byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.WriteByte(c)
}

// WriteRune 将一个 UTF-8 编码的字符写入缓冲区
func (b *Buffer) WriteRune(r rune) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.WriteRune(r)
}
