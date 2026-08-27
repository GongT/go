package packer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Unpacker = *unpacker
type unpacker struct {
	buff   bytes.Buffer
	endian binary.ByteOrder
}

// NewUnpack 创建一个新的unpacker实例，不应再使用data的原始引用
func NewUnpack(endian binary.ByteOrder, data []byte) Unpacker {
	return &unpacker{
		buff:   *bytes.NewBuffer(data),
		endian: endian,
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

// reader 从缓冲区读取或预览指定长度的数据
func reader[T any](peek bool, p *unpacker, size int, readFunc func([]byte) (T, error)) (T, error) {
	if p.buff.Len() < size {
		var zero T
		return zero, io.EOF
	}
	var data []byte
	if peek {
		data, _ = p.buff.Peek(size)
	} else {
		data = p.buff.Next(size)
	}
	return readFunc(data)
}

// ReadWithLen 从缓冲区中读取带长度前缀的数据，如果数据不足则返回io.EOF，且不读取任何数据
//
// 注意此方法返回的buffer是原始缓冲区的切片，如果长期引用需要复制，否则其内容可能会被后续读取操作覆盖
// func (p *unpacker) ReadWithLen() ([]byte, error) {
// 	size, err := p.PeekUint32()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if p.buff.Len() < int(4+size) {
// 		return nil, io.EOF
// 	}
// 	_ = p.buff.Next(4) // 跳过长度前缀
// 	return p.Next(int(size))
// }

// ReadStringWithLen 从缓冲区中读取带长度前缀的字符串，如果数据不足则返回io.EOF，且不读取任何数据
// func (p *unpacker) ReadStringWithLen() (string, error) {
// 	data, err := p.ReadWithLen()
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(data), nil
// }

// Peek 从缓冲区中预览指定长度的数据，但不移动读取指针
func (p *unpacker) Peek(length int) ([]byte, error) {
	return reader(true, p, length, func(data []byte) ([]byte, error) {
		return data, nil
	})
}

// PeekSafe 从缓冲区中预览指定长度的数据，但不移动读取指针
//
// 返回复制，可长期引用
func (p *unpacker) PeekSafe(length int) ([]byte, error) {
	return reader(true, p, length, func(data []byte) ([]byte, error) {
		copied := make([]byte, len(data))
		copy(copied, data)
		return copied, nil
	})
}

// Next 从缓冲区中读取指定长度的数据
func (p *unpacker) Next(length int) ([]byte, error) {
	return reader(false, p, length, func(data []byte) ([]byte, error) {
		return data, nil
	})
}

// NextSafe 从缓冲区中读取指定长度的数据
//
// 返回复制，可长期引用
func (p *unpacker) NextSafe(length int) ([]byte, error) {
	return reader(false, p, length, func(data []byte) ([]byte, error) {
		copied := make([]byte, len(data))
		copy(copied, data)
		return copied, nil
	})
}

// PeekString 从缓冲区中预览指定长度的字符串，但不移动读取指针
func (p *unpacker) PeekString(length int) (string, error) {
	return reader(true, p, length, func(data []byte) (string, error) {
		return string(data), nil
	})
}

// NextString 从缓冲区中读取指定长度的字符串
func (p *unpacker) NextString(length int) (string, error) {
	return reader(false, p, length, func(data []byte) (string, error) {
		return string(data), nil
	})
}

/* 无符号整数 */

// PeekUint64 从缓冲区中预览一个uint64值
func (p *unpacker) PeekUint64() (uint64, error) {
	return reader(true, p, 8, func(data []byte) (uint64, error) {
		return p.endian.Uint64(data), nil
	})
}

// PeekUint 从缓冲区中预览一个无符号整数
func (p *unpacker) PeekUint() (uint64, error) {
	return p.PeekUint64()
}

// NextUint64 从缓冲区中读取一个uint64值
func (p *unpacker) NextUint64() (uint64, error) {
	return reader(false, p, 8, func(data []byte) (uint64, error) {
		return p.endian.Uint64(data), nil
	})
}

// NextUint 从缓冲区中读取一个无符号整数
func (p *unpacker) NextUint() (uint64, error) {
	return p.NextUint64()
}

// PeekUint32 从缓冲区中预览一个uint32值
func (p *unpacker) PeekUint32() (uint32, error) {
	return reader(true, p, 4, func(data []byte) (uint32, error) {
		return p.endian.Uint32(data), nil
	})
}

// NextUint32 从缓冲区中读取一个uint32值
func (p *unpacker) NextUint32() (uint32, error) {
	return reader(false, p, 4, func(data []byte) (uint32, error) {
		return p.endian.Uint32(data), nil
	})
}

// PeekUint16 从缓冲区中预览一个uint16值
func (p *unpacker) PeekUint16() (uint16, error) {
	return reader(true, p, 2, func(data []byte) (uint16, error) {
		return p.endian.Uint16(data), nil
	})
}

// NextUint16 从缓冲区中读取一个uint16值
func (p *unpacker) NextUint16() (uint16, error) {
	return reader(false, p, 2, func(data []byte) (uint16, error) {
		return p.endian.Uint16(data), nil
	})
}

// PeekUint8 从缓冲区中预览一个uint8值
func (p *unpacker) PeekUint8() (uint8, error) {
	return reader(true, p, 1, func(data []byte) (uint8, error) {
		return data[0], nil
	})
}

// PeekByte 从缓冲区中预览一个字节
func (p *unpacker) PeekByte() (byte, error) {
	return p.PeekUint8()
}

// 读取布尔值，1表示true，0表示false，其他返回错误
func (p *unpacker) PeekBool() (bool, error) {
	return reader(true, p, 1, convBool)
}

// NextUint8 从缓冲区中读取一个uint8值
func (p *unpacker) NextUint8() (uint8, error) {
	return reader(false, p, 1, func(data []byte) (uint8, error) {
		return data[0], nil
	})
}

// NextByte 从缓冲区中读取一个字节
func (p *unpacker) NextByte() (byte, error) {
	return p.NextUint8()
}

// NextBool 从缓冲区中读取一个布尔值
func (p *unpacker) NextBool() (bool, error) {
	return reader(false, p, 1, convBool)
}

/* 有符号整数 */

// PeekInt 从缓冲区中预览一个有符号整数
func (p *unpacker) PeekInt() (int64, error) {
	return p.PeekInt64()
}

// PeekInt64 从缓冲区中预览一个int64值
func (p *unpacker) PeekInt64() (int64, error) {
	return reader(true, p, 8, func(data []byte) (int64, error) {
		return int64(p.endian.Uint64(data)), nil
	})
}

// NextInt 从缓冲区中读取一个有符号整数
func (p *unpacker) NextInt() (int64, error) {
	return p.NextInt64()
}

// NextInt64 从缓冲区中读取一个int64值
func (p *unpacker) NextInt64() (int64, error) {
	return reader(false, p, 8, func(data []byte) (int64, error) {
		return int64(p.endian.Uint64(data)), nil
	})
}

// PeekInt32 从缓冲区中预览一个int32值
func (p *unpacker) PeekInt32() (int32, error) {
	return reader(true, p, 4, func(data []byte) (int32, error) {
		return int32(p.endian.Uint32(data)), nil
	})
}

// NextInt32 从缓冲区中读取一个int32值
func (p *unpacker) NextInt32() (int32, error) {
	return reader(false, p, 4, func(data []byte) (int32, error) {
		return int32(p.endian.Uint32(data)), nil
	})
}

// PeekInt16 从缓冲区中预览一个int16值
func (p *unpacker) PeekInt16() (int16, error) {
	return reader(true, p, 2, func(data []byte) (int16, error) {
		return int16(p.endian.Uint16(data)), nil
	})
}

// NextInt16 从缓冲区中读取一个int16值
func (p *unpacker) NextInt16() (int16, error) {
	return reader(false, p, 2, func(data []byte) (int16, error) {
		return int16(p.endian.Uint16(data)), nil
	})
}

// PeekInt8 从缓冲区中预览一个int8值
func (p *unpacker) PeekInt8() (int8, error) {
	return reader(true, p, 1, func(data []byte) (int8, error) {
		return int8(data[0]), nil
	})
}

// NextInt8 从缓冲区中读取一个int8值
func (p *unpacker) NextInt8() (int8, error) {
	return reader(false, p, 1, func(data []byte) (int8, error) {
		return int8(data[0]), nil
	})
}

/* 浮点数 */

// PeekFloat32 从缓冲区中预览一个float32值
func (p *unpacker) PeekFloat32() (float32, error) {
	return reader(true, p, 4, func(data []byte) (float32, error) {
		return math.Float32frombits(p.endian.Uint32(data)), nil
	})
}

// NextFloat32 从缓冲区中读取一个float32值
func (p *unpacker) NextFloat32() (float32, error) {
	return reader(false, p, 4, func(data []byte) (float32, error) {
		return math.Float32frombits(p.endian.Uint32(data)), nil
	})
}

// PeekFloat64 从缓冲区中预览一个float64值
func (p *unpacker) PeekFloat64() (float64, error) {
	return reader(true, p, 8, func(data []byte) (float64, error) {
		return math.Float64frombits(p.endian.Uint64(data)), nil
	})
}

// NextFloat64 从缓冲区中读取一个float64值
func (p *unpacker) NextFloat64() (float64, error) {
	return reader(false, p, 8, func(data []byte) (float64, error) {
		return math.Float64frombits(p.endian.Uint64(data)), nil
	})
}

/* 工具 */

// convBool 将字节转换为布尔值
func convBool(data []byte) (bool, error) {
	switch data[0] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %d", data[0])
	}
}
