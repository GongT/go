// @exported

package internal

import (
	"iter"

	"github.com/gongt/go/assets/exports/internal/imported"
)

type privateType struct{}
type PublicType = *privateType

type SomeType struct{}

type SomeInterface interface {
	Method()
}

func TestingTypes(_ map[SomeType]map[uint64]SomeType, _ ****string, _ [][][][]string, _ [123]any) {
}

func TestingTypeStruct(_ *struct {
	// field1 是私有的，导出此函数无论如何都无法通过编译
	field1 SomeType
	Field2 string `json:"wow"`
}) {
}

func TestingTypeInterface(_ interface {
	Method1() string
	Method2(param SomeType) bool
}) {
}

func TestingTypeSignature[T SomeInterface](_ func(name T) (map[string]imported.Type, error)) {}

// func TestingInvalid(_ aaa) {}

func TestingChannels(_ chan SomeType, _ <-chan SomeType, _ chan<- SomeType) {
}

func Test123[T imported.Type](v any, t T, a, b string, s1 *Struct1, _ int) (string, error) {
	return "", nil
}

func TestIter[Y any]() iter.Seq2[string, Y] {
	return nil
}

// TestPrivate 有私有引用，编译时会出错，但不影响导出，用户必须自己保证相关类型已经导出
func TestPrivate(incorrect *privateType, correct PublicType) error {
	return nil
}
