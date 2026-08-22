// SOME HEADING TEXT
package test_pkg

import (
	udzBalWrrTt "iter"

	bZpItAMbgIOW "github.com/gongt/go/assets/exports/internal"
	bCUOMTOFUEed "github.com/gongt/go/pkg/errors/stacktrace"
)

type GenericType2[T comparable] = bZpItAMbgIOW.GenericType2[T]
type Struct1 = bZpItAMbgIOW.Struct1
type IFace1 = bZpItAMbgIOW.IFace1

func Test123[T bZpItAMbgIOW.GenericType](arg1 any, arg2 T, arg3 string, arg4 string, arg5 *Struct1, arg6 int) (string, error) {
	return bZpItAMbgIOW.Test123(arg1, arg2, arg3, arg4, arg5, arg6)
}

// Documentation for testX function
func TestX() (*Struct1, error) {
	return bZpItAMbgIOW.TestX()
}

func R() *bCUOMTOFUEed.StackTraceArray {
	return bZpItAMbgIOW.R()
}

func TestIter[Y any]() udzBalWrrTt.Seq2[string, Y] {
	return bZpItAMbgIOW.TestIter[Y]()
}

func Test2() {
	bZpItAMbgIOW.Test2()
}
