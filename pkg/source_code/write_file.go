package sourcecode

import "github.com/gongt/go/pkg/source_code/internal/writer"

type GoFileBuffer = writer.GoFileBuffer

func NewGoFileBuffer() GoFileBuffer {
	return writer.NewGoFileBuffer()
}
