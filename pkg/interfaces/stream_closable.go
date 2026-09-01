package interfaces

import "io"

type Closable interface {
	io.Closer

	IsClosed() bool
}

type ReadClosable interface {
	io.Reader

	Closable
}

type WriteClosable interface {
	io.Writer

	Closable
}

type ReadWriteClosable interface {
	io.Reader
	io.Writer

	Closable
}
