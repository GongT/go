// @exported

package streaming

import (
	"io"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/strings/shared"
	"github.com/gongt/go/pkg/strings/strtools"
)

var _ io.ReadWriteCloser = (*streamBuffer)(nil)

type StreamBuffer = *streamBuffer

// 阻塞Read操作，直到传入的[]byte能够被填满，或流关闭
//
// 此对象用于实现 “不断产生数据，消费者不断读取” 的场景，作为最初数据源使用，不适合放在管道中间作为缓冲
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

func notify(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

func (sb *streamBuffer) emit(p []byte) (int, error) {
	lp := len(p)

	if sb.closed {
		if sb.buffer.Len() >= lp {
			sb.buffer.Read(p)
			return lp, nil
		}
		return 0, io.EOF
	} else {
		// 调用前已经确保缓冲区中有足够的数据填满p
		n, _ := sb.buffer.Read(p)
		myenv.Assert(n == lp, "数据异常")

		notify(sb.readden)
	}
	return lp, nil
}

// Wait 阻塞直到缓冲区中有lp个字节
func (sb *streamBuffer) Wait(lp int) {
	if sb.closed {
		return
	}

	if sb.WaterMark > 0 {
		// 如果当前读取的长度超过水位线的一半，则将水位线调整为当前读取长度的两倍
		if lp*2 > sb.WaterMark {
			sb.WaterMark = lp * 2
		}
	}

	for sb.buffer.Len() < lp { // 一直等待到写入的数据足够填满p
		<-sb.written
		if sb.closed { // 等待时发生关闭
			return
		}
	}
}

// Read 从缓冲区中读取数据，缓冲区总是可以被填满（返回1 == len(p)），除非流关闭并返回EOF
func (sb *streamBuffer) Read(p []byte) (int, error) {
	lp := len(p)
	sb.Wait(lp)

	return sb.emit(p)
}

// Write 向缓冲区中写入数据，如果缓冲区超过水位线则阻塞，直到有数据被读出为止
func (sb *streamBuffer) Write(p []byte) (int, error) {
	if sb.closed {
		return 0, io.ErrClosedPipe
	}

	r, err := sb.buffer.Write(p)
	if err != nil {
		notify(sb.written)
	}

	if sb.WaterMark > 0 {
		// 如果本次写入的数据超过了水位线，则暂不返回
		for sb.buffer.Len() > sb.WaterMark {
			<-sb.readden
		}
	}

	return r, err
}

// WriteString 向缓冲区中写入字符串
func (b *streamBuffer) WriteString(s string) (n int, err error) {
	return b.Write(strtools.UnsafeBytes(s))
}

/* 免内存分配复制接口 */

func (b *streamBuffer) WriteTo(w io.Writer) (n int64, err error) {
	n, err = b.buffer.WriteTo(w)
	notify(b.readden)
	return
}
func (b *streamBuffer) ReadFrom(r io.Reader) (rn int64, err error) {
	var n int
	for {
		n, err = b.buffer.Appender(512, func(buff []byte) (int, error) {
			return r.Read(buff)
		})

		rn += int64(n)
		notify(b.written)

		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return
		}
	}
}

/* 缓冲区操作 */

func (b *streamBuffer) Reset() {
	b.buffer.Reset()
	notify(b.readden)
}

func (b *streamBuffer) Cap() int { return b.buffer.Cap() }

// func (b *streamBuffer) Grow(n int)     { b.buffer.Grow(n) }
// func (b *streamBuffer) Available() int { return b.buffer.Available() }
// func (b *streamBuffer) Len() int       { return b.buffer.Len() }

/* 可直接通过bufio实现 */

// func (b *streamBuffer) ReadBytes(delim byte) ([]byte, error)  { return b.buffer.ReadBytes(delim) }
// func (b *streamBuffer) ReadString(delim byte) (string, error) { return b.buffer.ReadString(delim) }
// func (b *streamBuffer) Peek(n int) ([]byte, error) { return b.buffer.Peek(n) }
// func (b *streamBuffer) WriteByte(c byte) error {}
// func (b *streamBuffer) WriteRune(r rune) (n int, err error) {}
// func (b *streamBuffer) ReadByte() (byte, error) {}
// func (b *streamBuffer) ReadRune() (rune, int, error) {}
// func (b *streamBuffer) UnreadRune() error {}
// func (b *streamBuffer) UnreadByte() error {}

/* 难以实现 */

// func (b *streamBuffer) Next(n int) []byte                         { return b.buffer.Next(n) }
// func (b *streamBuffer) Bytes() []byte                             { return b.buffer.Bytes() }
// func (b *streamBuffer) AvailableBuffer() []byte                   { return b.buffer.AvailableBuffer() }
