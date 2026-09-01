package closable

import (
	"bytes"
	"io"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/interfaces"
)

var _ interfaces.Closable = (*Buffer)(nil)

// 可以关闭的[bytes.Buffer]，关闭后所有写入操作都会返回错误，但仍然可以读取
type Buffer struct {
	buff bytes.Buffer

	closed bool
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{buff: *bytes.NewBuffer(buf)}
}
func NewBufferString(s string) *Buffer {
	return &Buffer{buff: *bytes.NewBufferString(s)}
}

func (b *Buffer) IsClosed() bool {
	return b.closed
}

func (b *Buffer) notClosed() error {
	if b.closed {
		return errors.EnsureTrace(io.ErrClosedPipe)
	}
	return nil
}

func (b *Buffer) Close() error {
	b.closed = true
	return nil
}

func (b *Buffer) Bytes() []byte                         { return b.buff.Bytes() }
func (b *Buffer) AvailableBuffer() []byte               { return b.buff.AvailableBuffer() }
func (b *Buffer) String() string                        { return b.buff.String() }
func (b *Buffer) Peek(n int) ([]byte, error)            { return b.buff.Peek(n) }
func (b *Buffer) Len() int                              { return b.buff.Len() }
func (b *Buffer) Cap() int                              { return b.buff.Cap() }
func (b *Buffer) Available() int                        { return b.buff.Available() }
func (b *Buffer) WriteTo(w io.Writer) (int64, error)    { return b.buff.WriteTo(w) }
func (b *Buffer) Read(p []byte) (int, error)            { return b.buff.Read(p) }
func (b *Buffer) Next(n int) []byte                     { return b.buff.Next(n) }
func (b *Buffer) UnreadRune() error                     { return b.buff.UnreadRune() }
func (b *Buffer) UnreadByte() error                     { return b.buff.UnreadByte() }
func (b *Buffer) ReadBytes(delim byte) ([]byte, error)  { return b.buff.ReadBytes(delim) }
func (b *Buffer) ReadString(delim byte) (string, error) { return b.buff.ReadString(delim) }
func (b *Buffer) ReadByte() (byte, error)               { return b.buff.ReadByte() }
func (b *Buffer) ReadRune() (rune, int, error)          { return b.buff.ReadRune() }

func (b *Buffer) Truncate(n int) {
	myenv.Must(b.notClosed())
	b.buff.Truncate(n)
}

func (b *Buffer) Reset() {
	myenv.Must(b.notClosed())
	b.buff.Reset()
}

func (b *Buffer) Grow(n int) {
	myenv.Must(b.notClosed())
	b.buff.Grow(n)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.Write(p)
}

func (b *Buffer) WriteString(s string) (n int, err error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.WriteString(s)
}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.ReadFrom(r)
}

func (b *Buffer) WriteByte(c byte) error {
	if err := b.notClosed(); err != nil {
		return err
	}
	return b.buff.WriteByte(c)
}

func (b *Buffer) WriteRune(r rune) (n int, err error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.WriteRune(r)
}
