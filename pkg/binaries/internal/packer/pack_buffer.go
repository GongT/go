package packer

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
)

func (p *packer_buffer_impl) Format(f fmt.State, verb rune) {
	if buff, ok := p.buff.(fmt.Stringer); verb == 's' && ok {
		fmt.Fprintf(f, "%q", buff.String())
	} else {
		fmt.Fprintf(f, "packer{endian: %s, buff: %v}", p.endian.String(), p.buff)
	}
}

func (p *packer_buffer_impl) g(n int) []byte {
	if p.buff.Available() < n {
		p.buff.Grow(n)
	}
	return p.buff.AvailableBuffer()[:n]
}

func (p *packer_buffer_impl) Close() error {
	if p.closed {
		return errors.EnsureTrace(sharederrors.ErrDuplicateCall)
	}
	p.closed = true
	if closer, ok := p.buff.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (p *packer_buffer_impl) IsClosed() bool {
	return p.closed
}

func (p *packer_buffer_impl) IsBuffered() bool {
	return true
}

/* 字符串 */

func (p *packer_buffer_impl) WriteBytes(data []byte) { p.buff.Write(data) }

func (p *packer_buffer_impl) WriteString(value string) {
	s := unsafe.StringData(value)
	p.buff.Write(unsafe.Slice(s, len(value)))
}

/* 无符号整数 */

func (p *packer_buffer_impl) WriteUint(value uint) { p.WriteUint64(uint64(value)) }

func (p *packer_buffer_impl) WriteUint64(value uint64) {
	b := p.g(8)
	p.endian.PutUint64(b, value)
	p.buff.Write(b)
}

func (p *packer_buffer_impl) WriteUint32(value uint32) {
	b := p.g(4)
	p.endian.PutUint32(b, value)
	p.buff.Write(b)
}

func (p *packer_buffer_impl) WriteUint16(value uint16) {
	b := p.g(2)
	p.endian.PutUint16(b, value)
	p.buff.Write(b)
}

func (p *packer_buffer_impl) WriteUint8(value uint8) {
	b := p.g(1)
	p.endian.PutUint8(b, value)
	p.buff.Write(b)
}

func (p *packer_buffer_impl) WriteByte(value byte) { p.WriteUint8(value) }

func (p *packer_buffer_impl) WriteBool(value bool) {
	if value {
		p.WriteUint8(1)
	} else {
		p.WriteUint8(0)
	}
}

/* 有符号整数 */

func (p *packer_buffer_impl) WriteSize(value int) {
	b := p.g(8)
	p.endian.PutSize(b, value)
	p.buff.Write(b)
}

func (p *packer_buffer_impl) WriteInt(value int) { p.WriteInt64(int64(value)) }

func (p *packer_buffer_impl) WriteInt64(value int64) {
	b := p.g(8)
	p.endian.PutInt64(b, value)
	p.buff.Write(b)
}

// WriteInt32 写入int32值
func (p *packer_buffer_impl) WriteInt32(value int32) {
	b := p.g(4)
	p.endian.PutInt32(b, value)
	p.buff.Write(b)
}

func (p *packer_buffer_impl) WriteRune(value rune) { p.WriteInt32(value) }

// WriteInt16 写入int16值
func (p *packer_buffer_impl) WriteInt16(value int16) {
	b := p.g(2)
	p.endian.PutInt16(b, value)
	p.buff.Write(b)
}

// WriteInt8 写入int8值
func (p *packer_buffer_impl) WriteInt8(value int8) {
	b := p.g(1)
	p.endian.PutInt8(b, value)
	p.buff.Write(b)
}

/* 浮点数 */

// WriteFloat64 写入float64值
func (p *packer_buffer_impl) WriteFloat64(value float64) {
	b := p.g(8)
	p.endian.PutFloat64(b, value)
	p.buff.Write(b)
}

// WriteFloat32 写入float32值
func (p *packer_buffer_impl) WriteFloat32(value float32) {
	b := p.g(4)
	p.endian.PutFloat32(b, value)
	p.buff.Write(b)
}

/* 复数 */

// WriteComplex64 写入complex64值
func (p *packer_buffer_impl) WriteComplex64(value complex64) {
	b := p.g(8)
	p.endian.PutComplex64(b, value)
	p.buff.Write(b)
}

// WriteComplex128 写入complex128值
func (p *packer_buffer_impl) WriteComplex128(value complex128) {
	b := p.g(16)
	p.endian.PutComplex128(b, value)
	p.buff.Write(b)
}
