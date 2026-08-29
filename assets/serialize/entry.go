//go:generate go run ../../cmd/serialize/main.go

package serialize

import (
	"github.com/gongt/go/assets/serialize/pkg1"
	"github.com/gongt/go/assets/serialize/pkg2"
)

type EmptyStruct struct{}

type ChildType interface {
	SomeMethod()
}

type privateType struct {
	Field1 pkg1.PublicType1
	Field2 pkg2.PointerToType2

	private1 uint32
}

type SomeType struct {
	Field1 pkg1.PublicType1
	Field2 *pkg2.PublicType2
	Field3 privateType
}

type TypeTest struct {
	typeInt   int
	typeInt8  int8
	typeInt16 int16
	typeInt32 int32
	typeInt64 int64

	typeBool   bool
	typeByte   byte
	typeUint   uint
	typeUint8  uint8
	typeUint16 uint16
	typeUint32 uint32
	typeUint64 uint64

	typeFloat32 float32
	typeFloat64 float64

	typeComplex64  complex64
	typeComplex128 complex128

	typeArray []int
	typeMap   map[string]int

	typeVArray []ChildType
	typeVMap   map[ChildType]ChildType

	typeString string
	typeBytes  []byte

	typeFunc    func(value int) string
	typeChannel <-chan int
}
