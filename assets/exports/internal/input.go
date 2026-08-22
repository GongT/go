// @exported
package internal

import (
	"iter"

	altername "github.com/gongt/go/pkg/errors/stacktrace"
)

// @private
type GenericType interface {
	string | int
}

type GenericType2[T comparable] struct {
	v T
}

type Struct1 struct {
	value string
}
type IFace1 interface {
	Method1()
}

func Test123[T GenericType](v any, t T, a, b string, s1 *Struct1, _ int) (string, error) {
	return "", nil
}

// Documentation for testX function
//
// @exported
func TestX() (r *Struct1, err error) {
	return
}

func R() *altername.StackTraceArray {
	return nil
}

// @private
const PrivateConst = 123

func TestIter[Y any]() iter.Seq2[string, Y] {
	return nil
}
