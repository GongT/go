package writercache

import (
	"io"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/interfaces"
)

var _ interfaces.ModernWriter = (*WriteEvent)(nil)

type WriteEvent struct {
	target interfaces.ModernWriter

	onWrite func()
	dirty   bool
}

func New(target interfaces.ModernWriter, onWrite func()) *WriteEvent {
	r := &WriteEvent{}
	r.RouteTo(target, onWrite)
	return r
}

// 重置缓存为干净状态
func (cw *WriteEvent) RouteTo(target interfaces.ModernWriter, onWrite func()) {
	if cw.target != nil || cw.onWrite != nil {
		panic("WriteEvent: RouteTo 只能调用一次")
	}
	myenv.AssertPtr("target", target, "onWrite", onWrite)
	cw.target = target
	cw.onWrite = onWrite
	cw.dirty = false
}

func (cw *WriteEvent) Close() {
	cw.target = nil
	cw.onWrite = nil
	cw.dirty = false
}

// 重置缓存为干净状态
func (cw *WriteEvent) Reset() {
	cw.dirty = false
}

// Write implements [interfaces.ModernWriter].
func (cw *WriteEvent) Write(p []byte) (n int, err error) {
	if cw.target == nil {
		return -1, io.EOF
	}
	cw.handleWrite()
	return cw.target.Write(p)
}

// WriteByte implements [interfaces.ModernWriter].
func (cw *WriteEvent) WriteByte(c byte) error {
	if cw.target == nil {
		return io.EOF
	}
	cw.handleWrite()
	return cw.target.WriteByte(c)
}

// WriteRune implements [interfaces.ModernWriter].
func (cw *WriteEvent) WriteRune(r rune) (n int, err error) {
	if cw.target == nil {
		return -1, io.EOF
	}
	cw.handleWrite()
	return cw.target.WriteRune(r)
}

// WriteString implements [interfaces.ModernWriter].
func (cw *WriteEvent) WriteString(s string) (n int, err error) {
	if cw.target == nil {
		return -1, io.EOF
	}
	cw.handleWrite()
	return cw.target.WriteString(s)
}

func (cw *WriteEvent) handleWrite() {
	if !cw.dirty {
		cw.dirty = true
		cw.onWrite()
	}
}
