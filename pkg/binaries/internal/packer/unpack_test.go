package packer

import (
	"io"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/binaries/internal/endian"
	"github.com/stretchr/testify/require"
)

func TestUnpackMatchesNativeEndianBytes(t *testing.T) {
	myenv.T(t)

	data := []byte{
		0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01,
		0x18, 0x17, 0x16, 0x15, 0x14, 0x13, 0x12, 0x11,
		0x24, 0x23, 0x22, 0x21,
		0x32, 0x31,
		0x01, 0x42,
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xfc, 0xff, 0xff, 0xff,
		0xfb, 0xff,
		0xfa,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xc0,
		0x00, 0x00, 0xc0, 0x3f,
	}
	unpacker := NewUnpack(endian.LittleEndian, data)

	assertNext := func(name string, expected any, read func() (any, error)) {
		t.Helper()
		actual, err := read()
		require.NoErrorf(t, err, "%s returned an error", name)
		require.Equalf(t, expected, actual, "%s returned an unexpected value", name)
	}

	assertNext("NextUint", uint(0x0102030405060708), func() (any, error) { return unpacker.NextUint() })
	assertNext("NextUint64", uint64(0x1112131415161718), func() (any, error) { return unpacker.NextUint64() })
	assertNext("NextUint32", uint32(0x21222324), func() (any, error) { return unpacker.NextUint32() })
	assertNext("NextUint16", uint16(0x3132), func() (any, error) { return unpacker.NextUint16() })
	assertNext("NextBool", true, func() (any, error) { return unpacker.NextBool() })
	assertNext("NextUint8", uint8(0x42), func() (any, error) { return unpacker.NextUint8() })
	assertNext("NextInt", -2, func() (any, error) { return unpacker.NextInt() })
	assertNext("NextInt64", int64(-3), func() (any, error) { return unpacker.NextInt64() })
	assertNext("NextInt32", int32(-4), func() (any, error) { return unpacker.NextInt32() })
	assertNext("NextInt16", int16(-5), func() (any, error) { return unpacker.NextInt16() })
	assertNext("NextInt8", int8(-6), func() (any, error) { return unpacker.NextInt8() })
	assertNext("NextFloat64", -2.5, func() (any, error) { return unpacker.NextFloat64() })
	assertNext("NextFloat32", float32(1.5), func() (any, error) { return unpacker.NextFloat32() })

	require.Zero(t, unpacker.Len())

	peek := NewUnpack(endian.LittleEndian, []byte{0xaa, 0xbb, 0xcc})
	preview, err := peek.Peek(2)
	require.NoError(t, err)
	require.Equal(t, []byte{0xaa, 0xbb}, preview)
	require.Equal(t, 3, peek.Len())
	_, err = peek.Next(4)
	require.ErrorIs(t, err, io.EOF)
}
