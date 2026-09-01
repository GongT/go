package struct_stream

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
	"github.com/gongt/go/pkg/strings/strtools"
)

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

func (sb *structuredBuffer) WriteString(s string) (int, error) {
	return sb.Write(strtools.UnsafeBytes(s))
}

// ReadFrom 从指定的io.Reader中读取，准备给Read使用
func (sb *structuredBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	n, err = sb.buffer.ReadFrom(r)
	return n, err
}
