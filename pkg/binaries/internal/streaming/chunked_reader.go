// @exported

package streaming

import (
	"bytes"
	"io"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/interfaces"
)

var _ io.ReadWriteCloser = (*chunkReader)(nil)

type ChunkReader = *chunkReader

// 输入数据，根据一个序列进行分割，每次Next返回一段，直到没有数据为止
type chunkReader struct {
	buffer bytes.Buffer
	sep    []byte

	done chan struct{}
	ch   chan []byte

	// IncludingSep 如果为true，则返回的数据块包含分隔符，否则不包含分隔符，默认为true
	IncludingSep bool
}

// NewChunkReader 创建一个新的chunkReader实例，sep为数据块的分隔符
func NewChunkReader[T interfaces.ByteSeq](sep T) ChunkReader {
	if len(sep) == 0 {
		panic(errors.NewAnonymous("sep不能为空"))
	}

	var arr []byte = make([]byte, len(sep))
	copy(arr, sep)
	return &chunkReader{
		sep:  arr,
		ch:   make(chan []byte),
		done: make(chan struct{}),

		IncludingSep: true,
	}
}

// Read 读取数据块，输入的p为存放数据的缓冲区，返回实际读取的字节数和可能的错误
//
// 注意: p必须能够容纳一个完整的数据块，否则返回错误，
// 除非有特殊需要，否则应该调用 [chunkReader.Next]
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	select {
	case data := <-cr.ch:
		if len(p) < len(data) {
			return 0, ErrInsufficientBuffer
		}
		copy(p, data)
		return len(data), nil
	case <-cr.done:
		return 0, io.EOF
	}
}

// Next 返回下一个完整的数据块，如果没有数据则等待
func (cr *chunkReader) Next() ([]byte, error) {
	select {
	case data := <-cr.ch:
		return data, nil
	case <-cr.done:
		return nil, io.EOF
	}
}

// Write 向缓冲区写入数据流，如果缓冲了一个完整的数据块，但还没有Read，则会阻塞
func (cr *chunkReader) Write(p []byte) (n int, err error) {
	if cr.done == nil {
		return 0, io.ErrClosedPipe
	}

	if _, err := cr.buffer.Write(p); err != nil {
		return 0, err
	}

	for {
		index := bytes.Index(cr.buffer.Bytes(), cr.sep)
		if index < 0 {
			break
		}

		chunk := cr.buffer.Next(index + len(cr.sep))
		if !cr.IncludingSep {
			chunk = chunk[:len(chunk)-len(cr.sep)]
		}
		if err := cr.emit(chunk); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (cr *chunkReader) emit(data []byte) error {
	select {
	case cr.ch <- data:
		return nil
	case <-cr.done:
		return io.EOF
	}
}

func (cr *chunkReader) Close() error {
	if cr.done == nil {
		panic(errors.NewAnonymous("重复释放"))
	}

	if cr.buffer.Len() > 0 {
		_ = cr.emit(cr.buffer.Bytes())
		cr.buffer.Reset()
	}

	close(cr.done)
	cr.done = nil

	close(cr.ch)
	cr.ch = nil

	return nil
}
