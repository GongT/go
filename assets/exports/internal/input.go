// @exported
package internal

import (
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
