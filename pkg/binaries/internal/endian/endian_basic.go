package endian

import (
	"log"
	"math"

	"github.com/gongt/go/pkg/errors"
)

// uint8

func (s ByteOrder) Uint8(b []byte) uint8                 { return b[0] }
func (s ByteOrder) AppendUint8(b []byte, v uint8) []byte { return append(b, v) }
func (s ByteOrder) PutUint8(b []byte, v uint8)           { b[0] = byte(v) }

// byte (uint8)

func (s ByteOrder) Byte(b []byte) byte                 { return s.Uint8(b) }
func (s ByteOrder) AppendByte(b []byte, v byte) []byte { return s.AppendUint8(b, v) }
func (s ByteOrder) PutByte(b []byte, v byte)           { s.PutUint8(b, v) }

// size (int保存为uint32)
func (s ByteOrder) Size(b []byte) int {
	return int(s.impl.Uint32(b))
}

// size (int保存为uint32)
func (s ByteOrder) PutSize(b []byte, v int) {
	if v <= 0 || v >= math.MaxInt32 {
		panic(errors.NewAnonymous("长度值异常"))
	}
	s.impl.PutUint32(b, uint32(v))
}

// size (int保存为uint32)
func (s ByteOrder) AppendSize(b []byte, v int) []byte {
	if v <= 0 || v >= math.MaxInt32 {
		panic(errors.NewAnonymous("长度值异常"))
	}
	return s.AppendUint32(b, uint32(v))
}

// bool

func (s ByteOrder) Bool(b []byte) bool                 { return fromByte(b[0]) }
func (s ByteOrder) AppendBool(b []byte, v bool) []byte { return s.AppendUint8(b, toByte(v)) }
func (s ByteOrder) PutBool(b []byte, v bool)           { s.PutUint8(b, toByte(v)) }

func toByte(v bool) uint8 {
	if v {
		return 1
	}
	return 0
}
func fromByte(b uint8) bool {
	switch b {
	case 1:
		return true
	case 0:
		return false
	default:
		log.Printf("Invalid boolean value: %v\n", b)
		return false
	}
}

// int (signed)

func (s ByteOrder) Int(b []byte) int                 { return int(s.Uint64(b)) }
func (s ByteOrder) AppendInt(b []byte, v int) []byte { return s.AppendUint64(b, uint64(v)) }
func (s ByteOrder) PutInt(b []byte, v int)           { s.PutUint64(b, uint64(v)) }

func (s ByteOrder) Int64(b []byte) int64                 { return int64(s.Uint64(b)) }
func (s ByteOrder) AppendInt64(b []byte, v int64) []byte { return s.AppendUint64(b, uint64(v)) }
func (s ByteOrder) PutInt64(b []byte, v int64)           { s.PutUint64(b, uint64(v)) }

func (s ByteOrder) Int32(b []byte) int32                 { return int32(s.Uint32(b)) }
func (s ByteOrder) AppendInt32(b []byte, v int32) []byte { return s.AppendUint32(b, uint32(v)) }
func (s ByteOrder) PutInt32(b []byte, v int32)           { s.PutUint32(b, uint32(v)) }

func (s ByteOrder) Int16(b []byte) int16                 { return int16(s.Uint16(b)) }
func (s ByteOrder) AppendInt16(b []byte, v int16) []byte { return s.AppendUint16(b, uint16(v)) }
func (s ByteOrder) PutInt16(b []byte, v int16)           { s.PutUint16(b, uint16(v)) }

func (s ByteOrder) Int8(b []byte) int8                 { return int8(s.Uint8(b)) }
func (s ByteOrder) AppendInt8(b []byte, v int8) []byte { return s.AppendUint8(b, uint8(v)) }
func (s ByteOrder) PutInt8(b []byte, v int8)           { s.PutUint8(b, uint8(v)) }

func (s ByteOrder) Rune(b []byte) rune                 { return s.Int32(b) }
func (s ByteOrder) AppendRune(b []byte, v rune) []byte { return s.AppendInt32(b, v) }
func (s ByteOrder) PutRune(b []byte, v rune)           { s.PutInt32(b, v) }

// float

func (s ByteOrder) Float32(b []byte) float32 {
	return math.Float32frombits(s.Uint32(b))
}
func (s ByteOrder) AppendFloat32(b []byte, v float32) []byte {
	return s.AppendUint32(b, math.Float32bits(v))
}
func (s ByteOrder) PutFloat32(b []byte, v float32) { s.PutUint32(b, math.Float32bits(v)) }

func (s ByteOrder) Float64(b []byte) float64 {
	return math.Float64frombits(s.Uint64(b))
}
func (s ByteOrder) AppendFloat64(b []byte, v float64) []byte {
	return s.AppendUint64(b, math.Float64bits(v))
}
func (s ByteOrder) PutFloat64(b []byte, v float64) {
	s.PutUint64(b, math.Float64bits(v))
}

// complex

func (s ByteOrder) Complex64(b []byte) complex64 {
	return complex64(complex(s.Float32(b), s.Float32(b[4:])))
}
func (s ByteOrder) AppendComplex64(b []byte, v complex64) []byte {
	b = s.AppendFloat32(b, real(v))
	return s.AppendFloat32(b, imag(v))
}
func (s ByteOrder) PutComplex64(b []byte, v complex64) {
	s.PutFloat32(b, real(v))
	s.PutFloat32(b[4:], imag(v))
}

func (s ByteOrder) Complex128(b []byte) complex128 {
	return complex128(complex(s.Float64(b), s.Float64(b[8:])))
}
func (s ByteOrder) AppendComplex128(b []byte, v complex128) []byte {
	b = s.AppendFloat64(b, real(v))
	return s.AppendFloat64(b, imag(v))
}
func (s ByteOrder) PutComplex128(b []byte, v complex128) {
	s.PutFloat64(b, real(v))
	s.PutFloat64(b[8:], imag(v))
}
