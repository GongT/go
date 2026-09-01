package closable

import (
	"io"
	"strings"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/interfaces"
)

var _ interfaces.Closable = (*Builder)(nil)

// 可以关闭的[strings.Builder]，关闭后所有写入操作都会返回错误，但仍然可以读取
type Builder struct {
	buff strings.Builder

	closed bool
}

func (b *Builder) IsClosed() bool {
	return b.closed
}

func (b *Builder) notClosed() error {
	if b.closed {
		return errors.EnsureTrace(io.ErrClosedPipe)
	}
	return nil
}

func (b *Builder) Close() error {
	b.closed = true
	return nil
}

func (b *Builder) String() string { return b.buff.String() }
func (b *Builder) Len() int       { return b.buff.Len() }
func (b *Builder) Cap() int       { return b.buff.Cap() }

func (b *Builder) Reset() {
	myenv.Must(b.notClosed())
	b.buff.Reset()
}

func (b *Builder) Grow(n int) {
	myenv.Must(b.notClosed())
	b.buff.Grow(n)
}

func (b *Builder) Write(p []byte) (int, error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.Write(p)
}

func (b *Builder) WriteByte(c byte) error {
	if err := b.notClosed(); err != nil {
		return err
	}
	return b.buff.WriteByte(c)
}

func (b *Builder) WriteRune(r rune) (int, error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.WriteRune(r)
}

func (b *Builder) WriteString(s string) (int, error) {
	if err := b.notClosed(); err != nil {
		return 0, err
	}
	return b.buff.WriteString(s)
}
