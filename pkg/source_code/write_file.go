package sourcecode

import "github.com/gongt/go/pkg/source_code/internal/writer"

// buffer

type GoFileBuffer = writer.GoFileBuffer

func NewGoFileBuffer() GoFileBuffer {
	return writer.NewGoFileBuffer()
}

var OriginalName = writer.OriginalName
var IndexName = writer.IndexName
var HashName = writer.HashName

// type_resolver

var ErrNotName = writer.ErrNotName
var ErrBasic = writer.ErrBasic

type TypeResolver = writer.TypeResolver

func NewTypeResolver(file GoFileBuffer) TypeResolver {
	return writer.NewTypeResolver(file)
}
