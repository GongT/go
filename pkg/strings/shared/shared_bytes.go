package shared

import (
	"bytes"
	"io"
	"sync"
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

// Bytes 返回缓冲区中的字节切片，并通过use回调提供给调用者
func (b *Buffer) Bytes[T any](use func([]byte) T) T {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return use(b.buff.Bytes())
}

// AvailableBuffer 返回缓冲区中可用的字节切片，并通过use回调提供给调用者
func (b *Buffer) AvailableBuffer[T any](use func([]byte) T) T {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return use(b.buff.AvailableBuffer())
}

// String 返回缓冲区的内容作为字符串
func (b *Buffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.String()
}

// Peek 检查缓冲区的前 n 个字节，如果不足则use仍会被调用，然后返回EOF
func (b *Buffer) Peek[T any](n int, use func([]byte) (T, error)) (T, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	buf, err := b.buff.Peek(n)
	ret, useErr := use(buf)
	if useErr == nil {
		return ret, err
	} else {
		return ret, useErr
	}
}

// Len 返回缓冲区中未读取的字节数
func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.Len()
}

// Cap 返回缓冲区的容量
func (b *Buffer) Cap() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.Cap()
}

// Available 返回缓冲区中可用的字节数
func (b *Buffer) Available() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buff.Available()
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

// Next 返回缓冲区的下 n 个字节，通过use回调提供给调用者
func (b *Buffer) Next[T any](n int, use func([]byte) T) T {
	b.mu.Lock()
	defer b.mu.Unlock()
	return use(b.buff.Next(n))
}

// UnreadRune 将缓冲区的最后一个读取的字符回退
func (b *Buffer) UnreadRune() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.UnreadRune()
}

// UnreadByte 将缓冲区的最后一个读取的字节回退
func (b *Buffer) UnreadByte() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.UnreadByte()
}

// ReadBytes 从缓冲区中读取直到遇到指定的分隔符，并返回包含分隔符的字节切片
func (b *Buffer) ReadBytes(delim byte) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadBytes(delim)
}

// ReadString 从缓冲区中读取直到遇到指定的分隔符，并返回包含分隔符的字符串
func (b *Buffer) ReadString(delim byte) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadString(delim)
}

// ReadByte 从缓冲区中读取一个字节
func (b *Buffer) ReadByte() (byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadByte()
}

// ReadRune 从缓冲区中读取一个 UTF-8 编码的字符，并返回该字符及其字节长度
func (b *Buffer) ReadRune() (rune, int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buff.ReadRune()
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

// Grow 增加缓冲区的容量，以确保至少可以容纳 n 个字节
func (b *Buffer) Grow(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buff.Grow(n)
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
