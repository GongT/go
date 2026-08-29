package sbins

import (
	"encoding/binary"
	"math"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/interfaces"
	"github.com/gongt/go/pkg/serialize/kinds"
	"github.com/gongt/go/pkg/strings/packer"
)

type PacketWrite = *packet_body_write
type packet_body_write struct {
	pack packer.Packer

	started bool
	done    bool
}

func NewPacketBody() *packet_body_write {
	return &packet_body_write{
		pack: packer.NewPack(binary.NativeEndian),
	}
}

func (p *packet_body_write) Bytes() []byte {
	if !p.done {
		panic(errors.NewAnonymous("packet not done"))
	}
	return p.pack.Bytes()
}

func (p *packet_body_write) Starting() error {
	if p.started {
		return errors.NewAnonymous("packet already started")
	}
	if p.done {
		panic(errors.NewAnonymous("packet already done"))
	}
	p.started = true
	return nil
}

func (p *packet_body_write) size(length int) {
	if !p.started {
		panic(errors.NewAnonymous("packet not started"))
	}
	if p.done {
		panic(errors.NewAnonymous("packet already done"))
	}
	if myenv.IsDebug {
		if length <= 0 || length >= math.MaxInt32 {
			panic(errors.NewAnonymous("invalid length"))
		}
	}
	p.pack.WriteInt32(int32(length))
}

func (p *packet_body_write) kind(kind kinds.ValueType) {
	if p.done {
		panic(errors.NewAnonymous("packet already done"))
	}
	p.pack.WriteUint16(uint16(kind))
}

func (p *packet_body_write) Write(data []byte) {
	p.kind(kinds.TypeIdBytes)
	p.size(len(data))
	p.pack.Write(data)
}

func (p *packet_body_write) WriteString(value string) {
	p.kind(kinds.TypeIdString)
	p.size(len(value))
	p.pack.Write([]byte(value))
}

func (p *packet_body_write) WriteUint(value uint) {
	p.kind(kinds.TypeIdUint)
	p.pack.WriteUint(value)
}

func (p *packet_body_write) WriteUint64(value uint64) {
	p.kind(kinds.TypeIdUint64)
	p.pack.WriteUint64(value)
}

func (p *packet_body_write) WriteUint32(value uint32) {
	p.kind(kinds.TypeIdUint32)
	p.pack.WriteUint32(value)
}

func (p *packet_body_write) WriteUint16(value uint16) {
	p.kind(kinds.TypeIdUint16)
	p.pack.WriteUint16(value)
}

func (p *packet_body_write) WriteUint8(value uint8) {
	p.kind(kinds.TypeIdUint8)
	p.pack.WriteUint8(value)
}

func (p *packet_body_write) WriteByte(value byte) {
	p.kind(kinds.TypeIdByte)
	p.pack.WriteUint8(value)
}

func (p *packet_body_write) WriteBool(value bool) {
	p.kind(kinds.TypeIdBool)
	p.pack.WriteBool(value)
}

func (p *packet_body_write) WriteInt(value int) {
	p.kind(kinds.TypeIdInt)
	p.pack.WriteInt(value)
}

func (p *packet_body_write) WriteInt64(value int64) {
	p.kind(kinds.TypeIdInt64)
	p.pack.WriteInt64(value)
}

func (p *packet_body_write) WriteInt32(value int32) {
	p.kind(kinds.TypeIdInt32)
	p.pack.WriteInt32(value)
}

func (p *packet_body_write) WriteInt16(value int16) {
	p.kind(kinds.TypeIdInt16)
	p.pack.WriteInt16(value)
}

func (p *packet_body_write) WriteInt8(value int8) {
	p.kind(kinds.TypeIdInt8)
	p.pack.WriteInt8(value)
}

func (p *packet_body_write) WriteFloat64(value float64) {
	p.kind(kinds.TypeIdFloat64)
	p.pack.WriteFloat64(value)
}

func (p *packet_body_write) WriteFloat32(value float32) {
	p.kind(kinds.TypeIdFloat32)
	p.pack.WriteFloat32(value)
}

func (p *packet_body_write) WriteComplex128(value complex128) {
	p.kind(kinds.TypeIdComplex128)
	p.pack.WriteComplex128(value)
}

func (p *packet_body_write) WriteComplex64(value complex64) {
	p.kind(kinds.TypeIdComplex64)
	p.pack.WriteComplex64(value)
}

func (p *packet_body_write) WriteAnyId(id kinds.ValueType, value uint64) {
	p.kind(id)
	p.pack.WriteUint64(value)
}

func (p *packet_body_write) WriteAnyBytes(id kinds.ValueType, data []byte) {
	p.kind(id)
	p.size(len(data))
	p.pack.Write(data)
}

func (p *packet_body_write) WriteArray[T interfaces.BuiltinLiterials](childType kinds.ValueType, data []T) {
	p.kind(kinds.TypeIdArray)
	p.size(len(data))
	before := p.pack.Len()
	for _, v := range data {
		p.pack.WriteAny(v)
	}
	after := p.pack.Len()
	if after-before != len(data)*kinds.SizeOfType(childType) {
		panic("mismatched array size")
	}
}

func (p *packet_body_write) WriteStruct(guid [16]byte, data []byte) {
	p.kind(kinds.TypeIdStruct)
	p.size(len(data) + len(guid))
	p.pack.Write(guid[:])
	p.pack.Write(data)
}

func (p *packet_body_write) WriteAnyHeader(id kinds.ValueType, size int) {
	p.kind(id)
	p.size(size)
}
