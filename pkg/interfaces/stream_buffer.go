package interfaces

import "io"

type BufferedWriter interface {
	io.Writer

	Available() int
	Grow(n int)
	// Cap() int
	// Len() int
	AvailableBuffer() []byte
}
