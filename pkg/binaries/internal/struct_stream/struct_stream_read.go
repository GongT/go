package struct_stream

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
	"github.com/gongt/go/pkg/strings/closable"
)

var _ io.ReadWriter = (*structuredBuffer)(nil)

type StructuredBuffer = *structuredBuffer

// structuredBuffer 一个[长度+缓冲区]序列
type structuredBuffer struct {
	buffer    closable.Buffer
	size_buff []byte

	errored bool
}

func New() StructuredBuffer {
	return &structuredBuffer{
		size_buff: make([]byte, 4), // 长度4字节的int32表示，使用LittleEndian存储
	}
}

// PeekSize 返回下一个结构化数据块的大小，没有数据返回0
func (sb *structuredBuffer) PeekSize() (int, error) {
	if err := sb.check(); err != nil {
		return -1, err
	}

	if sb.buffer.Len() < 4 {
		return 0, nil
	}

	buff, err := sb.buffer.Peek(4)
	if err != nil {
		sb.errored = true
		return -1, errors.EnsureTrace(err)
	}

	size := int(binary.LittleEndian.Uint32(buff))

	if size <= 0 {
		sb.errored = true
		return -1, errors.EnsureTrace(sharederrors.ErrDataCorrupted)
	}

	if sb.buffer.Len() < 4+size { // 虽然长度收到，但数据尚未完整
		return 0, nil
	}

	return size, nil
}

// Read 从结构化缓冲区读取数据，如果out不足以容纳数据，则返回错误
func (sb *structuredBuffer) Read(out []byte) (int, error) {
	size, err := sb.PeekSize()
	if err != nil {
		return 0, err
	}
	if size == 0 {
		return 0, nil
	}
	if len(out) < size {
		return 0, errors.EnsureTrace(sharederrors.ErrEntityTooLarge)
	}

	sb.buffer.Next(4) // 丢弃长度
	n, err := sb.buffer.Read(out[:size])

	if err != nil {
		sb.errored = true
		return n, errors.EnsureTrace(err)
	}
	return n, nil
}

// Write 将数据写入结构化缓冲区，前面会写入数据长度
func (sb *structuredBuffer) Write(data []byte) (int, error) {
	if err := sb.check(); err != nil {
		return 0, err
	}

	size := len(data)
	if size == 0 {
		return 0, nil
	}

	if size > math.MaxInt32 {
		return 0, errors.EnsureTrace(sharederrors.ErrEntityTooLarge)
	}

	binary.LittleEndian.PutUint32(sb.size_buff, uint32(size))
	sb.buffer.Write(sb.size_buff)
	sb.buffer.Write(data)

	return len(data), nil
}

func (sb *structuredBuffer) check() error {
	if sb.errored {
		return errors.EnsureTrace(sharederrors.ErrBrokenPipe)
	}
	return nil
}

// WriteTo 将Write进来的数据（添加长度后）写入到指定的io.Writer中
func (sb *structuredBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := sb.buffer.WriteTo(w)
	return n, err
}

// ReadFrom 从指定的io.Reader中读取，准备给Read使用
func (sb *structuredBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = sb.buffer.ReadFrom(r)
	return n, err
}
