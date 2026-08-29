// @exported

package packer

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/gongt/go/pkg/errors"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
	"github.com/gongt/go/pkg/interfaces"
)

type Packer interface {
	Close() error
	WriteBytes(data []byte)
	WriteString(value string)
	WriteUint(value uint)
	WriteUint64(value uint64)
	WriteUint32(value uint32)
	WriteUint16(value uint16)
	WriteUint8(value uint8)
	WriteBool(value bool)
	WriteInt(value int)
	WriteSize(value int)
	WriteInt64(value int64)
	WriteInt32(value int32)
	WriteInt16(value int16)
	WriteInt8(value int8)
	WriteFloat64(value float64)
	WriteFloat32(value float32)
	WriteComplex64(value complex64)
	WriteComplex128(value complex128)
}

var _ Packer = (*packer)(nil)

type packer struct {
	buff   interfaces.ModernWriter
	endian binary.ByteOrder

	closed bool
}

// NewPack 创建一个新的packer实例
func NewPack(buff io.Writer, endian binary.ByteOrder) Packer {
	if endian == nil {
		endian = binary.LittleEndian
	}

	r := &packer{
		buff:   interfaces.ModernizeWriter(buff),
		endian: endian,
	}
	return r
}

// Format 按指定格式输出packer
func (p *packer) Format(f fmt.State, verb rune) {
	if buff, ok := p.buff.(fmt.Stringer); verb == 's' && ok {
		fmt.Fprintf(f, "%q", buff.String())
	} else {
		fmt.Fprintf(f, "packer{endian: %s, buff: %v}", p.endian.String(), p.buff)
	}
}

func (p *packer) Close() error {
	if p.closed {
		return errors.EnsureTrace(sharederrors.ErrDuplicateCall)
	}
	p.closed = true
	if closer, ok := p.buff.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// WriteBytes 写入任意数据
func (p *packer) WriteBytes(data []byte) {
	p.buff.Write(data)
}

// WriteString 写入任意字符串
func (p *packer) WriteString(value string) {
	p.buff.WriteString(value)
}

/* 无符号整数 */

// WriteUint 写入无符号整数（一律写入为uint64）
func (p *packer) WriteUint(value uint) {
	p.WriteUint64(uint64(value))
}

// WriteUint64 写入uint64值
func (p *packer) WriteUint64(value uint64) {
	var buf [8]byte
	p.endian.PutUint64(buf[:], value)
	p.buff.Write(buf[:])
}

// WriteUint32 写入uint32值
func (p *packer) WriteUint32(value uint32) {
	var buf [4]byte
	p.endian.PutUint32(buf[:], value)
	p.buff.Write(buf[:])
}

// WriteUint16 写入uint16值
func (p *packer) WriteUint16(value uint16) {
	var buf [2]byte
	p.endian.PutUint16(buf[:], value)
	p.buff.Write(buf[:])
}

// WriteUint8 写入uint8值
func (p *packer) WriteUint8(value uint8) {
	p.buff.WriteByte(value)
}

// 写入布尔值，true写入1，false写入0
func (p *packer) WriteBool(value bool) {
	if value {
		p.buff.WriteByte(1)
	} else {
		p.buff.WriteByte(0)
	}
}

/* 有符号整数 */

// WriteInt 写入有符号整数，8字节
func (p *packer) WriteInt(value int) {
	p.WriteInt64(int64(value))
}

// [packer.WriteInt] 用于写入长度值
//
// 虽然go中长度是int(64)表示，但我将其限制为不能超过 math.MaxInt32 （当然也不会小于0）
func (p *packer) WriteSize(value int) {
	if value <= 0 || value >= math.MaxInt32 {
		panic(errors.NewAnonymous("长度值异常"))
	}
	p.WriteInt(value)
}

// WriteInt64 写入int64值
func (p *packer) WriteInt64(value int64) {
	var buf [8]byte
	p.endian.PutUint64(buf[:], uint64(value))
	p.buff.Write(buf[:])
}

// WriteInt32 写入int32值
func (p *packer) WriteInt32(value int32) {
	var buf [4]byte
	p.endian.PutUint32(buf[:], uint32(value))
	p.buff.Write(buf[:])
}

// WriteInt16 写入int16值
func (p *packer) WriteInt16(value int16) {
	var buf [2]byte
	p.endian.PutUint16(buf[:], uint16(value))
	p.buff.Write(buf[:])
}

// WriteInt8 写入int8值
func (p *packer) WriteInt8(value int8) {
	p.buff.WriteByte(byte(value))
}

/* 浮点数 */

// WriteFloat64 写入float64值
func (p *packer) WriteFloat64(value float64) {
	var buf [8]byte
	p.endian.PutUint64(buf[:], math.Float64bits(value))
	p.buff.Write(buf[:])
}

// WriteFloat32 写入float32值
func (p *packer) WriteFloat32(value float32) {
	var buf [4]byte
	p.endian.PutUint32(buf[:], math.Float32bits(value))
	p.buff.Write(buf[:])
}

/* 复数 */

// WriteComplex64 写入complex64值
func (p *packer) WriteComplex64(value complex64) {
	var buf [8]byte
	p.endian.PutUint32(buf[:4], math.Float32bits(real(value)))
	p.endian.PutUint32(buf[4:], math.Float32bits(imag(value)))
	p.buff.Write(buf[:])
}

// WriteComplex128 写入complex128值
func (p *packer) WriteComplex128(value complex128) {
	var buf [16]byte
	p.endian.PutUint64(buf[:8], math.Float64bits(real(value)))
	p.endian.PutUint64(buf[8:], math.Float64bits(imag(value)))

	p.buff.Write(buf[:])
}
