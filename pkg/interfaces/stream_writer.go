package interfaces

import "io"

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
	return WrapWriter(w)
}

func WrapWriter(w io.Writer) ModernWriter {
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
