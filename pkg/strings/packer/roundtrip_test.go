package packer

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackThenUnpackNativeEndianValues(t *testing.T) {
	packer := NewPack(binary.NativeEndian)
	values := struct {
		uintValue    uint
		uint64Value  uint64
		uint32Value  uint32
		uint16Value  uint16
		byteValue    byte
		boolValue    bool
		uint8Value   uint8
		intValue     int
		int64Value   int64
		int32Value   int32
		int16Value   int16
		int8Value    int8
		float64Value float64
		float32Value float32
		text         string
		raw          []byte
	}{
		uintValue:    0x0102030405060708,
		uint64Value:  0x1112131415161718,
		uint32Value:  0x21222324,
		uint16Value:  0x3132,
		byteValue:    0x41,
		boolValue:    true,
		uint8Value:   0x42,
		intValue:     -2,
		int64Value:   -3,
		int32Value:   -4,
		int16Value:   -5,
		int8Value:    -6,
		float64Value: -2.5,
		float32Value: 1.5,
		text:         "native",
		raw:          []byte{0xde, 0xad, 0xbe, 0xef},
	}

	packer.WriteUint(values.uintValue)
	packer.WriteUint64(values.uint64Value)
	packer.WriteUint32(values.uint32Value)
	packer.WriteUint16(values.uint16Value)
	packer.WriteByte(values.byteValue)
	packer.WriteBool(values.boolValue)
	packer.WriteUint8(values.uint8Value)
	packer.WriteInt(values.intValue)
	packer.WriteInt64(values.int64Value)
	packer.WriteInt32(values.int32Value)
	packer.WriteInt16(values.int16Value)
	packer.WriteInt8(values.int8Value)
	packer.WriteFloat64(values.float64Value)
	packer.WriteFloat32(values.float32Value)
	packer.WriteString(values.text)
	packer.Write(values.raw)

	unpacker := NewUnpack(binary.NativeEndian, packer.Bytes())
	assert := func(name string, expected any, read func() (any, error)) {
		t.Helper()
		actual, err := read()
		require.NoErrorf(t, err, "%s returned an error", name)
		require.Equalf(t, expected, actual, "%s returned an unexpected value", name)
	}
	assert("NextUint", uint64(values.uintValue), func() (any, error) { return unpacker.NextUint() })
	assert("NextUint64", values.uint64Value, func() (any, error) { return unpacker.NextUint64() })
	assert("NextUint32", values.uint32Value, func() (any, error) { return unpacker.NextUint32() })
	assert("NextUint16", values.uint16Value, func() (any, error) { return unpacker.NextUint16() })
	assert("NextByte", values.byteValue, func() (any, error) { return unpacker.NextByte() })
	assert("NextBool", values.boolValue, func() (any, error) { return unpacker.NextBool() })
	assert("NextUint8", values.uint8Value, func() (any, error) { return unpacker.NextUint8() })
	assert("NextInt", int64(values.intValue), func() (any, error) { return unpacker.NextInt() })
	assert("NextInt64", values.int64Value, func() (any, error) { return unpacker.NextInt64() })
	assert("NextInt32", values.int32Value, func() (any, error) { return unpacker.NextInt32() })
	assert("NextInt16", values.int16Value, func() (any, error) { return unpacker.NextInt16() })
	assert("NextInt8", values.int8Value, func() (any, error) { return unpacker.NextInt8() })
	assert("NextFloat64", values.float64Value, func() (any, error) { return unpacker.NextFloat64() })
	assert("NextFloat32", values.float32Value, func() (any, error) { return unpacker.NextFloat32() })

	text, err := unpacker.NextString(len(values.text))
	require.NoError(t, err)
	require.Equal(t, values.text, text)
	raw, err := unpacker.Next(len(values.raw))
	require.NoError(t, err)
	require.Equal(t, values.raw, raw)
	require.Zero(t, unpacker.Len())
}
