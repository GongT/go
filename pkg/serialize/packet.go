package serialize

import (
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/strings/packer"
)

const (
	startFlag = 0x12345678ABCDEF00
	endFlag   = 0x00FEDCBA87654321
)

// packet 一个简易的数据包
//   - 8 起始标志，用于校验
//   - 8 帧长度
//   - 8 元数据长度
//   - n 元数据
//   - 8 数据长度
//   - n 数据
//   - 8 结束标志，用于校验
type packet struct {
	pack packer.Packer

	done bool
}

// NewPacket 创建一个新的[packet]
func CreatePacket() *packet {
	r := &packet{
		pack: packer.NewPack(nil),
	}
	r.pack.WriteUint64(startFlag)
	return r
}

func (p *packet) Pack() []byte {
	if !p.done {
		p.pack.WriteUint64(endFlag)
		p.done = true
	}
	return p.pack.Bytes()
}

func (p *packet) Output() packer.Packer {
	return p.pack
}

func (p *packet) size(length int) {
	if p.done {
		panic(errors.NewAnonymous("packet already done"))
	}
	p.pack.WriteUint64(uint64(length))
}

func (p *packet) Write(data []byte) {
	p.size(len(data))
	p.pack.Write(data)
}

func (p *packet) WriteString(value string) {
	p.pack.Write([]byte(value))
}

func (p *packet) WriteUint(value uint) {
	p.size(8)
	p.pack.WriteUint(value)
}

func (p *packet) WriteUint64(value uint64) {
	p.size(8)
	p.pack.WriteUint64(value)
}

func (p *packet) WriteUint32(value uint32) {
	p.size(4)
	p.pack.WriteUint32(value)
}

func (p *packet) WriteUint16(value uint16) {
	p.size(2)
	p.pack.WriteUint16(value)
}

func (p *packet) WriteByte(value byte) {
	p.size(1)
	p.pack.WriteByte(value)
}

func (p *packet) WriteBool(value bool) {
	p.size(1)
	p.pack.WriteBool(value)
}

func (p *packet) WriteUint8(value uint8) {
	p.size(1)
	p.pack.WriteUint8(value)
}

func (p *packet) WriteInt(value int) {
	p.size(8)
	p.pack.WriteInt(value)
}

func (p *packet) WriteInt64(value int64) {
	p.size(8)
	p.pack.WriteInt64(value)
}

func (p *packet) WriteInt32(value int32) {
	p.size(4)
	p.pack.WriteInt32(value)
}

func (p *packet) WriteInt16(value int16) {
	p.size(2)
	p.pack.WriteInt16(value)
}

func (p *packet) WriteInt8(value int8) {
	p.size(1)
	p.pack.WriteInt8(value)
}

func (p *packet) WriteFloat64(value float64) {
	p.size(8)
	p.pack.WriteFloat64(value)
}

func (p *packet) WriteFloat32(value float32) {
	p.size(4)
	p.pack.WriteFloat32(value)
}
