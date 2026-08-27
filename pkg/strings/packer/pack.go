package packer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type Packer = *packer
type packer struct {
	buff   bytes.Buffer
	endian binary.ByteOrder
}

// NewPack 创建一个新的packer实例
func NewPack(endian binary.ByteOrder) Packer {
	if endian == nil {
		endian = binary.LittleEndian
	}
	return &packer{
		buff:   bytes.Buffer{},
		endian: endian,
	}
}

// Format 按指定格式输出packer
func (p *packer) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(f, "packer{endian: %s, buff: %v}", p.endian.String(), p.buff)
	case 's':
		fmt.Fprintf(f, "packer{endian: %s, len: %d}", p.endian.String(), p.buff.Len())
	case 'q':
		fmt.Fprintf(f, "%q", p.buff.String())
	}
}

// Bytes 返回缓冲区中的全部数据
func (p *packer) Bytes() []byte {
	return p.buff.Bytes()
}

// Len 返回缓冲区中的数据长度
func (p *packer) Len() int {
	return p.buff.Len()
}

// Reset 清空缓冲区中的数据
func (p *packer) Reset() {
	p.buff.Reset()
}

// Write 写入任意数据
func (p *packer) Write(data []byte) {
	p.buff.Write(data)
}

// WriteString 写入任意字符串
func (p *packer) WriteString(value string) {
	p.buff.WriteString(value)
}

// WriteWithLen 写入带长度前缀的数据
//
// 长度由uint32表示
// func (p *packer) WriteWithLen(data []byte) error {
// 	size := len(data)
// 	if size > math.MaxUint32 {
// 		return fmt.Errorf("data too large: %d bytes", size)
// 	}
// 	if err := p.WriteUint32(uint32(size)); err != nil {
// 		return err
// 	}
// 	_, err := p.buff.Write(data)
// 	return err
// }

// WriteStringWithLen 写入带长度前缀的字符串
// func (p *packer) WriteStringWithLen(value string) error {
// 	return p.WriteWithLen([]byte(value))
// }

/* 无符号整数 */

// WriteUint 写入无符号整数
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

// WriteByte 写入一个字节
func (p *packer) WriteByte(value byte) {
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

// WriteUint8 写入uint8值
func (p *packer) WriteUint8(value uint8) {
	p.buff.WriteByte(value)
}

/* 有符号整数 */

// WriteInt 写入有符号整数
func (p *packer) WriteInt(value int) {
	p.WriteInt64(int64(value))
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
