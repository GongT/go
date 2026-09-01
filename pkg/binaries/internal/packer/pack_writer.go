package packer

import (
	"fmt"
	"io"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
)

func (p *packer_writer_impl) Format(f fmt.State, verb rune) {
	if buff, ok := p.buff.(fmt.Stringer); verb == 's' && ok {
		fmt.Fprintf(f, "%q", buff.String())
	} else {
		fmt.Fprintf(f, "packer{endian: %s, buff: %v}", p.endian.String(), p.buff)
	}
}

func (p *packer_writer_impl) Close() error {
	if p.closed {
		return errors.EnsureTrace(sharederrors.ErrDuplicateCall)
	}
	p.closed = true
	if closer, ok := p.buff.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func (p *packer_writer_impl) IsClosed() bool {
	return p.closed
}

func (p *packer_writer_impl) IsBuffered() bool {
	return false
}

/* 字符串 */

func (p *packer_writer_impl) WriteBytes(data []byte) {
	p.buff.Write(data)
}

func (p *packer_writer_impl) WriteString(value string) {
	p.buff.WriteString(value)
}

/* 无符号整数 */

func (p *packer_writer_impl) WriteUint(value uint) {
	p.WriteUint64(uint64(value))
}

func (p *packer_writer_impl) WriteUint64(value uint64) {
	p.endian.PutUint64(p.mem[:], value)
	p.buff.Write(p.mem[:8])
}

func (p *packer_writer_impl) WriteUint32(value uint32) {
	p.endian.PutUint32(p.mem[:], value)
	p.buff.Write(p.mem[:4])
}

func (p *packer_writer_impl) WriteUint16(value uint16) {
	p.endian.PutUint16(p.mem[:], value)
	p.buff.Write(p.mem[:2])
}

func (p *packer_writer_impl) WriteUint8(value uint8) {
	p.buff.WriteByte(value)
}

func (p *packer_writer_impl) WriteByte(value byte) { p.WriteUint8(value) }

func (p *packer_writer_impl) WriteBool(value bool) {
	if value {
		p.WriteUint8(1)
	} else {
		p.WriteUint8(0)
	}
}

/* 有符号整数 */

func (p *packer_writer_impl) WriteSize(value int) {
	p.endian.PutSize(p.mem[:], value)
	p.buff.Write(p.mem[:8])
}

func (p *packer_writer_impl) WriteInt(value int) {
	p.WriteInt64(int64(value))
}

func (p *packer_writer_impl) WriteInt64(value int64) {
	p.endian.PutInt64(p.mem[:], value)
	p.buff.Write(p.mem[:8])
}

func (p *packer_writer_impl) WriteInt32(value int32) {
	p.endian.PutInt32(p.mem[:], value)
	p.buff.Write(p.mem[:4])
}

func (p *packer_writer_impl) WriteRune(value rune) { p.WriteInt32(value) }

func (p *packer_writer_impl) WriteInt16(value int16) {
	p.endian.PutInt16(p.mem[:], value)
	p.buff.Write(p.mem[:2])
}

func (p *packer_writer_impl) WriteInt8(value int8) {
	p.buff.WriteByte(byte(value))
}

/* 浮点数 */

func (p *packer_writer_impl) WriteFloat64(value float64) {
	p.endian.PutFloat64(p.mem[:], value)
	p.buff.Write(p.mem[:8])
}

func (p *packer_writer_impl) WriteFloat32(value float32) {
	p.endian.PutFloat32(p.mem[:], value)
	p.buff.Write(p.mem[:4])
}

/* 复数 */

func (p *packer_writer_impl) WriteComplex64(value complex64) {
	p.endian.PutComplex64(p.mem[:8], value)
	p.buff.Write(p.mem[:8])
}

func (p *packer_writer_impl) WriteComplex128(value complex128) {
	p.endian.PutComplex128(p.mem[:16], value)
	p.buff.Write(p.mem[:16])
}
