package checker

import (
	"reflect"

	"github.com/gongt/go/internal/myenv"
)

func IsPointer(v any) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Pointer, reflect.Func, reflect.Invalid:
		return true
	default:
		return false
	}
}

func DevelAssertPointer(v any, what ...string) {
	if myenv.IsDebug && !IsPointer(v) {
		if len(what) > 0 {
			panic("DevelAssertPointer: 断言失败: 不是指针类型: " + what[0])
		} else {
			panic("DevelAssertPointer: 断言失败: 不是指针类型")
		}
	}
}
