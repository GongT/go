// @exported

package packer

import (
	"io"

	"github.com/gongt/go/pkg/binaries/internal/endian"
	"github.com/gongt/go/pkg/interfaces"
)

type Packer interface {
	// Close 调用buff的Close()方法，如果不支持[io.Closer]，则什么都不做
	Close() error

	// WriteBytes 写入任意数据
	WriteBytes(data []byte)

	// WriteString 写入任意字符串
	WriteString(value string)

	// WriteUint 写入无符号整数（一律写入为uint64）
	WriteUint(value uint)

	// WriteUint64 写入uint64值
	WriteUint64(value uint64)

	// WriteUint32 写入uint32值
	WriteUint32(value uint32)

	// WriteUint16 写入uint16值
	WriteUint16(value uint16)

	// WriteUint8 写入uint8值
	WriteUint8(value uint8)

	// WriteByte 写入单个字节
	WriteByte(value byte)

	// WriteBool 写入布尔值，true写入1，false写入0
	WriteBool(value bool)

	// WriteInt 写入有符号整数，8字节
	WriteInt(value int)

	// [Packer.WriteInt] 用于写入长度值
	//
	// 虽然go中长度是int(64)表示，但我将其限制为不能超过 math.MaxInt32 （当然也不会小于0）
	WriteSize(value int)

	// WriteInt64 写入int64值
	WriteInt64(value int64)

	// WriteInt32 写入int32值
	WriteInt32(value int32)

	// WriteRune 写入rune值
	WriteRune(value rune)

	// WriteInt16 写入int16值
	WriteInt16(value int16)

	// WriteInt8 写入int8值
	WriteInt8(value int8)

	// WriteFloat64 写入float64值
	WriteFloat64(value float64)

	// WriteFloat32 写入float32值
	WriteFloat32(value float32)

	// WriteComplex64 写入complex64值
	WriteComplex64(value complex64)

	// WriteComplex128 写入complex128值
	WriteComplex128(value complex128)

	// IsClosed 返回packer是否已经关闭，是否有意义取决于具体实现
	IsClosed() bool
	// IsBuffered 返回packer是否是基于缓冲区的实现
	IsBuffered() bool
}

var _ Packer = (*packer_writer_impl)(nil)
var _ Packer = (*packer_buffer_impl)(nil)

type packer_writer_impl struct {
	buff   interfaces.ModernWriter
	endian endian.ByteOrder
	mem    [16]byte

	closed bool
}

type packer_buffer_impl struct {
	buff   interfaces.BufferedWriter
	endian endian.ByteOrder

	closed bool
}

// NewPack 创建一个新的packer实例
func NewPack[T io.Writer](buff T, byteOrder endian.ByteOrder) Packer {
	return NewPackEndian(buff, endian.LittleEndian)
}

func NewPackEndian[T io.Writer](buff T, byteOrder endian.ByteOrder) Packer {
	switch t := any(buff).(type) {
	case interfaces.BufferedWriter:
		return &packer_buffer_impl{
			buff:   t,
			endian: byteOrder,
		}
	default:
		return &packer_writer_impl{
			buff:   interfaces.ModernizeWriter(buff),
			endian: byteOrder,
		}
	}
}
