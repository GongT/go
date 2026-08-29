// @exported

package streaming

import (
	"io"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/strings/shared"
)

var _ io.ReadWriteCloser = (*streamBuffer)(nil)

type StreamBuffer = *streamBuffer

// 阻塞Read操作，直到传入的[]byte能够被填满，或流关闭
type streamBuffer struct {
	buffer shared.Buffer
	closed bool

	// WaterMark 设置缓冲区的水位线，当缓冲区长度超过该值时，写入操作将阻塞，直到有数据被读出为止，默认4096字节
	WaterMark int

	// 有数据被写入
	written chan struct{}

	// 有数据被读出
	readden chan struct{}
}

func NewStreamBuffer() StreamBuffer {
	return &streamBuffer{
		written:   make(chan struct{}),
		readden:   make(chan struct{}),
		WaterMark: 4096,
	}
}

func notify(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func (sb *streamBuffer) emit(p []byte) (int, error) {
	lp := len(p)

	if sb.closed {
		if sb.buffer.Len() == lp {
			sb.buffer.Read(p)
			return lp, nil
		}
		return 0, io.EOF
	} else {
		// it's ensured the buffer has enough data to fill p

		n, _ := sb.buffer.Read(p)
		myenv.Assert(n == lp, "数据异常")

		notify(sb.readden)
	}
	return lp, nil
}

// Read 从缓冲区中读取数据，缓冲区总是可以被填满（返回1 == len(p)），除非流关闭并返回EOF
func (sb *streamBuffer) Read(p []byte) (int, error) {
	lp := len(p)

	if sb.WaterMark > 0 {
		// 如果当前读取的长度超过水位线的一半，则将水位线调整为当前读取长度的两倍
		if lp*2 > sb.WaterMark {
			sb.WaterMark = lp * 2
			notify(sb.readden)
		}
	}

	for sb.buffer.Len() < lp {
		<-sb.written

		if sb.closed {
			return sb.emit(p)
		}
	}

	return sb.emit(p)
}

func (sb *streamBuffer) Write(p []byte) (int, error) {
	if sb.closed {
		return 0, io.ErrClosedPipe
	}

	r, err := sb.buffer.Write(p)
	if err != nil {
		notify(sb.written)
	}

	if sb.WaterMark > 0 {
		for sb.buffer.Len() > sb.WaterMark {
			<-sb.readden
		}
	}

	return r, err
}

func (sb *streamBuffer) Close() error {
	if sb.closed {
		panic(errors.NewAnonymous("重复释放"))
	}

	sb.closed = true
	close(sb.written)
	sb.written = nil
	close(sb.readden)
	sb.readden = nil
	return nil
}

// func (b *streamBuffer) Bytes() []byte                             { return b.buffer.Bytes() }
// func (b *streamBuffer) AvailableBuffer() []byte                   { return b.buffer.AvailableBuffer() }
// func (b *streamBuffer) Peek(n int) ([]byte, error)                { return b.buffer.Peek(n) }
// func (b *streamBuffer) Len() int                                  { return b.buffer.Len() }
// func (b *streamBuffer) Cap() int                                  { return b.buffer.Cap() }
// func (b *streamBuffer) Available() int                            { return b.buffer.Available() }
// func (b *streamBuffer) WriteTo(w io.Writer) (int64, error)        { return b.buffer.WriteTo(w) }
// func (b *streamBuffer) ReadFrom(r io.Reader) (n int64, err error) { return b.buffer.ReadFrom(r) }
// func (b *streamBuffer) Next(n int) []byte                         { return b.buffer.Next(n) }
// func (b *streamBuffer) UnreadRune() error                         { return b.buffer.UnreadRune() }
// func (b *streamBuffer) UnreadByte() error                         { return b.buffer.UnreadByte() }
// func (b *streamBuffer) ReadBytes(delim byte) ([]byte, error)      { return b.buffer.ReadBytes(delim) }
// func (b *streamBuffer) ReadString(delim byte) (string, error)     { return b.buffer.ReadString(delim) }
// func (b *streamBuffer) ReadByte() (byte, error)                   { return b.buffer.ReadByte() }
// func (b *streamBuffer) ReadRune() (rune, int, error)              { return b.buffer.ReadRune() }
// func (b *streamBuffer) Reset()                                    {}
// func (b *streamBuffer) Grow(n int)                                {}
// func (b *streamBuffer) WriteString(s string) (n int, err error)   {}
// func (b *streamBuffer) WriteByte(c byte) error                    {}
// func (b *streamBuffer) WriteRune(r rune) (n int, err error)       {}
