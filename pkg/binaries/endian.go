package binaries

import "encoding/binary"

// LittleEndian 包装[binary.LittleEndian]添加了各种类型的函数
var LittleEndian = binary.LittleEndian

type littleEndian interface{}

func NewEndian() Endian {
	binary.LittleEndian
}
