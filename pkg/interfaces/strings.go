package interfaces

import "io"

// StringLike 是一个泛型接口，表示可以被视为字符串的类型。表达此对象要被用作字符串。
type StringLike = ByteSeq

// ByteSeq 是一个泛型接口，表示可以被视为字节序列的类型。表达此对象要被用作字节序列。
type ByteSeq interface {
	~string | ~[]byte
}

type ToBytes interface {
	Bytes() []byte
}

type BuiltinLiterials interface {
	string | bool | int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8 | float64 | float32 | complex128 | complex64
}

type CustomLiterials interface {
	~string | ~bool | ~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8 | ~float64 | ~float32 | ~complex128 | ~complex64
}

type ModernWriter interface {
	io.Writer
	io.StringWriter
	io.ByteWriter

	WriteRune(r rune) (n int, err error)
}

type wrapWriter struct {
	io.Writer
}

func ModernizeWriter(w io.Writer) ModernWriter {
	if mw, ok := w.(ModernWriter); ok {
		return mw
	}
	return &wrapWriter{w}
}

func (w *wrapWriter) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func (w *wrapWriter) WriteByte(c byte) error {
	_, err := w.Write([]byte{c})
	return err
}

func (w *wrapWriter) WriteRune(r rune) (n int, err error) {
	return w.Write([]byte(string(r)))
}
