//go:build !release

package myenv

import (
	"fmt"
	"testing"
)

const IsDebug = true
const IsRelease = false

var IsTesting = testing.Testing()

func Assert(condition bool, message string, args ...any) {
	if !condition {
		panic(fmt.Sprintf(message, args...))
	}
}

// AssertPtr 检查多个指针是否为 nil，如果为 nil 则 panic
//
//	AssertPtr("ptr1", ptr1, "另一个指针", nil) // panic("另一个指针: 指针不能为 nil")
func AssertPtr(pairs ...any) {
	for i := 0; i < len(pairs); i += 2 {
		name := pairs[i]
		ptr := pairs[i+1]
		Assert(ptr != nil, "%s: 指针不能为 nil", name)
	}
}
