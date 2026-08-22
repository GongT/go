package writercache

import (
	"io"

	"github.com/gongt/go/internal/myenv"
)

type AllWriter interface {
	io.Writer
	io.StringWriter
	io.ByteWriter

	WriteRune(r rune) (n int, err error)
}

var _ AllWriter = (*WriteEvent)(nil)

type WriteEvent struct {
	target AllWriter

	onWrite func()
	dirty   bool
}

func New(target AllWriter, onWrite func()) *WriteEvent {
	r := &WriteEvent{}
	r.RouteTo(target, onWrite)
	return r
}

// 重置缓存为干净状态
func (cw *WriteEvent) RouteTo(target AllWriter, onWrite func()) {
	if cw.target != nil || cw.onWrite != nil {
		panic("WriteEvent: RouteTo 只能调用一次")
	}
	myenv.AssertPtr("target", target, "onWrite", onWrite)
	cw.target = target
	cw.onWrite = onWrite
	cw.dirty = false
}

// 重置缓存为干净状态
func (cw *WriteEvent) Reset() {
	cw.dirty = false
}

// Write implements [AllWriter].
func (cw *WriteEvent) Write(p []byte) (n int, err error) {
	cw.handleWrite()
	return cw.target.Write(p)
}

// WriteByte implements [AllWriter].
func (cw *WriteEvent) WriteByte(c byte) error {
	cw.handleWrite()
	return cw.target.WriteByte(c)
}

// WriteRune implements [AllWriter].
func (cw *WriteEvent) WriteRune(r rune) (n int, err error) {
	cw.handleWrite()
	return cw.target.WriteRune(r)
}

// WriteString implements [AllWriter].
func (cw *WriteEvent) WriteString(s string) (n int, err error) {
	cw.handleWrite()
	return cw.target.WriteString(s)
}

func (cw *WriteEvent) handleWrite() {
	if !cw.dirty {
		cw.dirty = true
		cw.onWrite()
	}
}
