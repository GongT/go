package endian

import "encoding/binary"

// @exported
type impl interface {
	binary.AppendByteOrder
	binary.ByteOrder
}

var _ impl = LittleEndian
var _ impl = BigEndian
var _ impl = NativeEndian

// LittleEndian 包装[binary.LittleEndian]添加了各种类型的函数
//
// @exported
var LittleEndian = create(binary.LittleEndian, 1)

// BigEndian 包装[binary.BigEndian]添加了各种类型的函数
//
// @exported
var BigEndian = create(binary.BigEndian, 2)

// NativeEndian 包装本地字节序添加了各种类型的函数
//
// @exported
var NativeEndian = create(binary.NativeEndian, 3)

type ByteOrder struct {
	impl

	_type int
}

func create(order impl, typ int) ByteOrder {
	return ByteOrder{impl: order, _type: typ}
}

func (b ByteOrder) IsNull() bool {
	return b._type == 0
}
