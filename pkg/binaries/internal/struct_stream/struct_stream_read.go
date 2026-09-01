package struct_stream

import (
	"encoding/binary"
	"io"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
)

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

func (sb *structuredBuffer) ReadString() (string, error) {
	size, err := sb.PeekSize()
	if err != nil {
		return "", err
	}
	if size == 0 {
		return "", nil
	}

	buf := make([]byte, size)
	n, err := sb.Read(buf)
	if err != nil {
		return "", err
	}
	if n != size {
		return "", errors.EnsureTrace(sharederrors.ErrDataCorrupted)
	}
	return string(buf), nil
}


// WriteTo 将Write进来的数据（添加长度后）写入到指定的io.Writer中
func (sb *structuredBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := sb.buffer.WriteTo(w)
	return n, err
}
