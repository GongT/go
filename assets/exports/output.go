// SOME HEADING TEXT
package test_pkg

import (
	pkg2 "iter"

	pkg0 "github.com/gongt/go/assets/exports/internal"
	pkg1 "github.com/gongt/go/assets/exports/internal/imported"
	pkg3 "github.com/gongt/go/pkg/errors/stacktrace"
)

// - ./assets/exports/internal/functions.go

type PublicType = pkg0.PublicType
type SomeType = pkg0.SomeType
type SomeInterface = pkg0.SomeInterface

func TestingTypes(arg1 map[pkg0.SomeType]map[uint64]pkg0.SomeType, arg2 ****string, arg3 [][][][]string, arg4 [123]any) {
	pkg0.TestingTypes(arg1, arg2, arg3, arg4)
}

func TestingTypeStruct(arg1 *struct {
	field1 pkg0.SomeType
	Field2 string "json:\"wow\""
}) {
	pkg0.TestingTypeStruct(arg1)
}

func TestingTypeInterface(arg1 interface {
	Method1() string
	Method2(param pkg0.SomeType) bool
}) {
	pkg0.TestingTypeInterface(arg1)
}

func TestingTypeSignature[T pkg0.SomeInterface](arg1 func(name T) (map[string]pkg1.Type, error)) {
	pkg0.TestingTypeSignature(arg1)
}

func TestingChannels(arg1 chan pkg0.SomeType, arg2 <-chan pkg0.SomeType, arg3 chan<- pkg0.SomeType) {
	pkg0.TestingChannels(arg1, arg2, arg3)
}

func Test123[T pkg1.Type](arg1 any, arg2 T, arg3 string, arg4 string, arg5 *pkg0.Struct1, arg6 int) (string, error) {
	return pkg0.Test123(arg1, arg2, arg3, arg4, arg5, arg6)
}

func TestIter[Y any]() pkg2.Seq2[string, Y] {
	return pkg0.TestIter[Y]()
}

// TestPrivate 有私有引用，编译时会出错，但不影响导出，用户必须自己保证相关类型已经导出
func TestPrivate(arg1 *pkg0.privateType, arg2 pkg0.PublicType) error {
	return pkg0.TestPrivate(arg1, arg2)
}

// - ./assets/exports/internal/imported/type.go

// - ./assets/exports/internal/input.go

type GenericType2[T comparable] = pkg0.GenericType2[T]
type Struct1 = pkg0.Struct1
type IFace1 = pkg0.IFace1

// Documentation for testX function
func TestX() (*pkg0.Struct1, error) {
	return pkg0.TestX()
}

func R() *pkg3.StackTraceArray {
	return pkg0.R()
}

// - ./assets/exports/internal/new.go

func NewSomeStruct() *pkg0.SomeStruct {
	return pkg0.NewSomeStruct()
}

// - ./assets/exports/internal/test.go

func Test2() {
	pkg0.Test2()
}
