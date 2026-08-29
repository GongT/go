package sbins

import (
	"bytes"
	"encoding/binary"
	"math"
	"uuid"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/interfaces"
	"github.com/gongt/go/pkg/serialize/kinds"
	"github.com/gongt/go/pkg/strings/packer"
)

type PacketRead = *packet_read
type packet_read struct {
	unpack  packer.Unpacker
	started bool
	errored bool
	ended   bool
}

func NewBodyRead(reader ) PacketRead {
	return &packet_read{
		unpack:  packer.NewUnpack(binary.NativeEndian, data),
		started: true,
	}
}

func (p *packet_read) err[T error](err T) T {
	p.errored = true
	p.unpack = nil
	return err
}

func (p *packet_read) err2[T error, VT any](v VT, err T) (VT, T) {
	if any(err) != nil {
		return v, p.err(err)
	}
	var no T
	return v, no
}

func (p *packet_read) next_size() (int, error) {
	if p.unpack.Len() == 0 || p.ended {
		err := p.err(errors.NewAnonymous("数据已结束"))
		return -1, err
	}
	length, err := p.unpack.NextInt32()
	if err != nil {
		return -1, p.err(err)
	}

	if myenv.IsDebug {
		if length <= 0 || length == math.MaxInt32 {
			return -1, p.err(errors.NewAnonymous("长度数据异常"))
		}
	}

	return int(length), nil
}

func (p *packet_read) next_kind(expected kinds.ValueType) error {
	if p.unpack.Len() == 0 || p.ended {
		return p.err(errors.NewAnonymous("数据已结束"))
	}
	kind, err := p.err2(p.unpack.NextUint16())
	if err != nil {
		return err
	}
	if kinds.ValueType(kind) != expected {
		return p.err(errors.NewAnonymous("字段类型数据异常")).WithDetails(
			"expected", expected.String(),
			"actual", kinds.ValueType(kind).String(),
		)
	}
	return nil
}

func (p *packet_read) ReadTypeRaw(id kinds.ValueType) ([]byte, error) {
	if err := p.next_kind(id); err != nil {
		return nil, err
	}
	if length, err := p.next_size(); err != nil {
		return nil, err
	} else {
		return p.err2(p.unpack.NextSafe(length))
	}
}

func (p *packet_read) Read() ([]byte, error) {
	if err := p.next_kind(kinds.TypeIdBytes); err != nil {
		return nil, err
	}
	if length, err := p.next_size(); err != nil {
		return nil, err
	} else {
		return p.err2(p.unpack.NextSafe(length))
	}
}

func (p *packet_read) ReadString() (string, error) {
	if err := p.next_kind(kinds.TypeIdString); err != nil {
		return "", err
	}
	if length, err := p.next_size(); err != nil {
		return "", err
	} else {
		return p.err2(p.unpack.NextString(length))
	}
}

func (p *packet_read) ReadUint() (uint, error) {
	if err := p.next_kind(kinds.TypeIdUint); err != nil {
		return 0, err
	}

	return p.err2(p.unpack.NextUint())
}
func (p *packet_read) ReadUint64() (uint64, error) {
	if err := p.next_kind(kinds.TypeIdUint64); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextUint64())
}

func (p *packet_read) ReadUint32() (uint32, error) {
	if err := p.next_kind(kinds.TypeIdUint32); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextUint32())
}

func (p *packet_read) ReadUint16() (uint16, error) {
	if err := p.next_kind(kinds.TypeIdUint16); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextUint16())
}

func (p *packet_read) ReadUint8() (uint8, error) {
	if err := p.next_kind(kinds.TypeIdUint8); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextUint8())
}

func (p *packet_read) ReadByte() (byte, error) {
	if err := p.next_kind(kinds.TypeIdByte); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextUint8())
}

func (p *packet_read) ReadBool() (bool, error) {
	if err := p.next_kind(kinds.TypeIdBool); err != nil {
		return false, err
	}
	return p.err2(p.unpack.NextBool())
}

func (p *packet_read) ReadInt() (int, error) {
	if err := p.next_kind(kinds.TypeIdInt); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextInt())
}

func (p *packet_read) ReadInt64() (int64, error) {
	if err := p.next_kind(kinds.TypeIdInt64); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextInt64())
}

func (p *packet_read) ReadInt32() (int32, error) {
	if err := p.next_kind(kinds.TypeIdInt32); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextInt32())
}

func (p *packet_read) ReadInt16() (int16, error) {
	if err := p.next_kind(kinds.TypeIdInt16); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextInt16())
}

func (p *packet_read) ReadInt8() (int8, error) {
	if err := p.next_kind(kinds.TypeIdInt8); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextInt8())
}

func (p *packet_read) ReadFloat64() (float64, error) {
	if err := p.next_kind(kinds.TypeIdFloat64); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextFloat64())
}

func (p *packet_read) ReadFloat32() (float32, error) {
	if err := p.next_kind(kinds.TypeIdFloat32); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextFloat32())
}

func (p *packet_read) ReadComplex128() (complex128, error) {
	if err := p.next_kind(kinds.TypeIdComplex128); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextComplex128())
}

func (p *packet_read) ReadComplex64() (complex64, error) {
	if err := p.next_kind(kinds.TypeIdComplex64); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextComplex64())
}

func (p *packet_read) ReadAnyId(id kinds.ValueType) (uint64, error) {
	if err := p.next_kind(id); err != nil {
		return 0, err
	}
	return p.err2(p.unpack.NextUint64())
}

func (p *packet_read) ReadAnyBytes(id kinds.ValueType) ([]byte, error) {
	if err := p.next_kind(id); err != nil {
		return nil, err
	}

	if length, err := p.next_size(); err != nil {
		return nil, err
	} else {
		return p.err2(p.unpack.Next(length))
	}
}

func (p *packet_read) ReadStruct(guid [16]byte) ([]byte, error) {
	body, err := p.ReadAnyBytes(kinds.TypeIdStruct)
	if err != nil {
		return nil, err
	}

	got := body[:16] // Extract the GUID from the body
	if !bytes.Equal(got, guid[:]) {
		return nil, errors.NewAnonymous("结构体GUID不匹配").WithDetails("expected", uuid.UUID(guid), "actual", uuid.UUID(got))
	}

	return body[16:], nil
}

func (p *packet_read) ReadArray[T interfaces.BuiltinLiterials](childType kinds.ValueType) ([]T, error) {
	if err := p.next_kind(kinds.TypeIdArray); err != nil {
		return nil, err
	}

	if size, err := p.next_size(); err != nil {
		return nil, err
	} else {
		data := make([]T, size)
		for i := range data {
			if data[i], err = p.unpack.NextAny[T](); err != nil {
				return nil, err
			}
		}
		return data, nil
	}
}

func (p *packet_read) ReadAnyHeader(id kinds.ValueType) (int, error) {
	if err := p.next_kind(id); err != nil {
		return 0, err
	}
	return p.next_size()
}
