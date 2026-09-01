// @exported

package packer

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/gongt/go/pkg/binaries/internal/endian"
	"github.com/gongt/go/pkg/errors"
)

type Unpacker = *unpacker
type unpacker struct {
	buff   bytes.Buffer
	endian endian.ByteOrder
}

// NewUnpack 创建一个新的unpacker实例，不应再使用data的原始引用，尤其是不能修改data的内容
func NewUnpack(byteOrder endian.ByteOrder, data []byte) Unpacker {
	if byteOrder.IsNull() {
		byteOrder = endian.LittleEndian
	}
	return &unpacker{
		buff:   *bytes.NewBuffer(data),
		endian: byteOrder,
	}
}

// Format 按指定格式输出unpacker
func (p *unpacker) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(f, "unpacker{endian: %s, buff: %v}", p.endian.String(), p.buff)
	case 's':
		fmt.Fprintf(f, "unpacker{endian: %s, len: %d}", p.endian.String(), p.buff.Len())
	case 'q':
		fmt.Fprintf(f, "%q", p.buff.String())
	}
}

// Len 返回缓冲区中剩余的数据长度
func (p *unpacker) Len() int {
	return p.buff.Len()
}

func (p *unpacker) Grow(n int) {
	p.buff.Grow(n)
}

// reader 从缓冲区读取或预览指定长度的数据
func reader[T any](peek bool, p *unpacker, size int, readFunc func([]byte, endian.ByteOrder) (T, error)) (T, error) {
	if p.buff.Len() < size {
		var zero T
		return zero, errors.EnsureTrace(io.EOF)
	}
	var data []byte
	if peek {
		data, _ = p.buff.Peek(size)
	} else {
		data = p.buff.Next(size)
	}
	return readFunc(data, p.endian)
}

// Peek 从缓冲区中预览指定长度的数据，但不移动读取指针
func (p *unpacker) Peek(length int) ([]byte, error) {
	return reader(true, p, length, _read_bytes)
}

// Next 从缓冲区中读取指定长度的数据
func (p *unpacker) Next(length int) ([]byte, error) {
	return reader(false, p, length, _read_bytes)
}

// PeekSafe 从缓冲区中预览指定长度的数据，但不移动读取指针
//
// 返回复制，可长期引用
func (p *unpacker) PeekSafe(length int) ([]byte, error) {
	return reader(true, p, length, _read_bytes_safe)
}

// NextSafe 从缓冲区中读取指定长度的数据
//
// 返回复制，可长期引用
func (p *unpacker) NextSafe(length int) ([]byte, error) {
	return reader(false, p, length, _read_bytes_safe)
}

// PeekString 从缓冲区中预览指定长度的字符串，但不移动读取指针
func (p *unpacker) PeekString(length int) (string, error) {
	return reader(true, p, length, _read_string)
}

// NextString 从缓冲区中读取指定长度的字符串
func (p *unpacker) NextString(length int) (string, error) {
	return reader(false, p, length, _read_string)
}

/* 无符号整数 */

// PeekUint64 从缓冲区中预览一个uint64值
func (p *unpacker) PeekUint64() (uint64, error) {
	return reader(true, p, 8, _read_uint64)
}

// PeekUint 从缓冲区中预览一个无符号整数
func (p *unpacker) PeekUint() (uint, error) {
	v, err := p.PeekUint64()
	return uint(v), err
}

// NextUint64 从缓冲区中读取一个uint64值
func (p *unpacker) NextUint64() (uint64, error) {
	return reader(false, p, 8, _read_uint64)
}

// NextUint 从缓冲区中读取一个无符号整数
func (p *unpacker) NextUint() (uint, error) {
	v, err := p.NextUint64()
	return uint(v), err
}

// PeekUint32 从缓冲区中预览一个uint32值
func (p *unpacker) PeekUint32() (uint32, error) {
	return reader(true, p, 4, _read_uint32)
}

// NextUint32 从缓冲区中读取一个uint32值
func (p *unpacker) NextUint32() (uint32, error) {
	return reader(false, p, 4, _read_uint32)
}

// PeekUint16 从缓冲区中预览一个uint16值
func (p *unpacker) PeekUint16() (uint16, error) {
	return reader(true, p, 2, _read_uint16)
}

// NextUint16 从缓冲区中读取一个uint16值
func (p *unpacker) NextUint16() (uint16, error) {
	return reader(false, p, 2, _read_uint16)
}

// PeekUint8 从缓冲区中预览一个uint8值
func (p *unpacker) PeekUint8() (uint8, error) {
	return reader(true, p, 1, _read_uint8)
}

// NextUint8 从缓冲区中读取一个uint8值
func (p *unpacker) NextUint8() (uint8, error) {
	return reader(false, p, 1, _read_uint8)
}

// 读取布尔值，1表示true，0表示false，其他返回错误
func (p *unpacker) PeekBool() (bool, error) {
	return reader(true, p, 1, _read_bool)
}

// NextBool 从缓冲区中读取一个布尔值
func (p *unpacker) NextBool() (bool, error) {
	return reader(false, p, 1, _read_bool)
}

/* 有符号整数 */

// PeekInt 从缓冲区中预览一个有符号整数
func (p *unpacker) PeekInt() (int, error) {
	v, err := p.PeekInt64()
	return int(v), err
}

// NextInt 从缓冲区中读取一个有符号整数
func (p *unpacker) NextInt() (int, error) {
	v, err := p.NextInt64()
	return int(v), err
}

// [unpacker.NextInt]
func (p *unpacker) NextSize() (int, error) {
	if v, err := p.NextInt(); err == nil {
		if v <= 0 || v >= math.MaxInt32 {
			return -1, errors.NewAnonymous("长度值异常")
		}
		return v, nil
	} else {
		return -1, err
	}
}

// PeekInt64 从缓冲区中预览一个int64值
func (p *unpacker) PeekInt64() (int64, error) {
	return reader(true, p, 8, _read_int64)
}

// NextInt64 从缓冲区中读取一个int64值
func (p *unpacker) NextInt64() (int64, error) {
	return reader(false, p, 8, _read_int64)
}

// PeekInt32 从缓冲区中预览一个int32值
func (p *unpacker) PeekInt32() (int32, error) {
	return reader(true, p, 4, _read_int32)
}

// NextInt32 从缓冲区中读取一个int32值
func (p *unpacker) NextInt32() (int32, error) {
	return reader(false, p, 4, _read_int32)
}

// PeekInt16 从缓冲区中预览一个int16值
func (p *unpacker) PeekInt16() (int16, error) {
	return reader(true, p, 2, _read_int16)
}

// NextInt16 从缓冲区中读取一个int16值
func (p *unpacker) NextInt16() (int16, error) {
	return reader(false, p, 2, _read_int16)
}

// PeekInt8 从缓冲区中预览一个int8值
func (p *unpacker) PeekInt8() (int8, error) {
	return reader(true, p, 1, _read_int8)
}

// NextInt8 从缓冲区中读取一个int8值
func (p *unpacker) NextInt8() (int8, error) {
	return reader(false, p, 1, _read_int8)
}

/* 浮点数 */

// PeekFloat32 从缓冲区中预览一个float32值
func (p *unpacker) PeekFloat32() (float32, error) {
	return reader(true, p, 4, _read_float32)
}

// NextFloat32 从缓冲区中读取一个float32值
func (p *unpacker) NextFloat32() (float32, error) {
	return reader(false, p, 4, _read_float32)
}

// PeekFloat64 从缓冲区中预览一个float64值
func (p *unpacker) PeekFloat64() (float64, error) {
	return reader(true, p, 8, _read_float64)
}

// NextFloat64 从缓冲区中读取一个float64值
func (p *unpacker) NextFloat64() (float64, error) {
	return reader(false, p, 8, _read_float64)
}

/* 复数 */

// PeekComplex64 从缓冲区中预览一个complex64值
func (p *unpacker) PeekComplex64() (complex64, error) {
	return reader(true, p, 8, _read_complex64)
}

// NextComplex64 从缓冲区中读取一个complex64值
func (p *unpacker) NextComplex64() (complex64, error) {
	return reader(false, p, 8, _read_complex64)
}

// PeekComplex128 从缓冲区中预览一个complex128值
func (p *unpacker) PeekComplex128() (complex128, error) {
	return reader(true, p, 16, _read_complex128)
}

// NextComplex128 从缓冲区中读取一个complex128值
func (p *unpacker) NextComplex128() (complex128, error) {
	return reader(false, p, 16, _read_complex128)
}

/* 工具 */

func _read_bytes(data []byte, _ endian.ByteOrder) ([]byte, error) {
	return data, nil
}

func _read_bytes_safe(data []byte, _ endian.ByteOrder) ([]byte, error) {
	copied := make([]byte, len(data))
	copy(copied, data)
	return copied, nil
}

func _read_string(data []byte, _ endian.ByteOrder) (string, error) {
	return string(data), nil
}

func _read_bool(data []byte, _ endian.ByteOrder) (bool, error) {
	switch data[0] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %d", data[0])
	}
}

func _read_uint64(data []byte, endian endian.ByteOrder) (uint64, error) {
	return endian.Uint64(data), nil
}
func _read_uint32(data []byte, endian endian.ByteOrder) (uint32, error) {
	return endian.Uint32(data), nil
}
func _read_uint16(data []byte, endian endian.ByteOrder) (uint16, error) {
	return endian.Uint16(data), nil
}
func _read_uint8(data []byte, endian endian.ByteOrder) (uint8, error) {
	return data[0], nil
}

func _read_int64(data []byte, endian endian.ByteOrder) (int64, error) {
	return endian.Int64(data), nil
}
func _read_int32(data []byte, endian endian.ByteOrder) (int32, error) {
	return endian.Int32(data), nil
}
func _read_int16(data []byte, endian endian.ByteOrder) (int16, error) {
	return endian.Int16(data), nil
}
func _read_int8(data []byte, endian endian.ByteOrder) (int8, error) {
	return int8(data[0]), nil
}

func _read_float32(data []byte, endian endian.ByteOrder) (float32, error) {
	return endian.Float32(data), nil
}
func _read_float64(data []byte, endian endian.ByteOrder) (float64, error) {
	return endian.Float64(data), nil
}

func _read_complex64(data []byte, endian endian.ByteOrder) (complex64, error) {
	return endian.Complex64(data), nil
}

func _read_complex128(data []byte, endian endian.ByteOrder) (complex128, error) {
	return endian.Complex128(data), nil
}
