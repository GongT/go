package struct_stream

import (
	"bytes"
	"io"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
	"github.com/gongt/go/pkg/interfaces"
)

var _ io.ReadWriter = (*structuredBuffer)(nil)

type StructuredBuffer = *structuredBuffer

// structuredBuffer 一个[长度+缓冲区]序列
type structuredBuffer struct {
	buffer    bytes.Buffer
	size_buff []byte

	errored bool
	closed  bool
}

func New() StructuredBuffer {
	return &structuredBuffer{
		size_buff: make([]byte, 4), // 长度4字节的int32表示，使用LittleEndian存储
	}
}

// Close 逻辑关闭，由Write方调用，供Read方检查使用
func (sb *structuredBuffer) Close() error {
	sb.closed = true
	if closable, ok := any(sb.buffer).(io.Closer); ok {
		return closable.Close()
	}
	return nil
}

// IsClosed 返回是否已被逻辑关闭
func (sb *structuredBuffer) IsClosed() bool {
	if closable, ok := any(sb.buffer).(interfaces.Closable); ok {
		return closable.IsClosed()
	}
	return sb.closed
}

func (sb *structuredBuffer) check() error {
	if sb.errored {
		return errors.EnsureTrace(sharederrors.ErrBrokenPipe)
	}
	return nil
}
