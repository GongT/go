package tools_bin

import (
	"go/types"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/i18n/type_name"
)

type fieldKind int

const (
	kindFallback fieldKind = iota
	kindScalar             // 布尔、整数、浮点数、字符串
	kindSlice
	kindMap
	kindChan
	kindFunc
)

// classifyType 根据类型的底层结构对字段进行分类，命名类型（例如 type X int）会被展开为其底层类型判断
func classifyType(t types.Type) fieldKind {
	switch u := t.Underlying().(type) {
	case *types.Basic:
		if _, ok := basicWriteInfo(u.Kind()); ok {
			return kindScalar
		} else {
			panic(errors.NewAnonymous("无法处理的标量类型: %s", type_name.TranslateBasicType(u)))
		}
	case *types.Slice:
		return kindSlice
	case *types.Map:
		return kindMap
	case *types.Chan:
		return kindChan
	case *types.Signature:
		return kindFunc
	default:
		// 结构体、指针、interface（含error）、数组等，一律回退
		return kindFallback
	}
}

// isInterfaceType 判断类型的底层是否为interface（包括 error 和 any）或任何类型的指针
func isInterfaceType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Interface)
	if !ok {
		_, ok = t.Underlying().(*types.Pointer)
	}
	return ok
}

// basicWriteInfo 返回标量类型对应的packer读写方法名后缀、写入时需要转换到的go内建类型名
func basicWriteInfo(kind types.BasicKind) (goType string, ok bool) {
	switch kind {
	case types.Bool:
		return "bool", true
	case types.String:
		return "string", true
	case types.Int:
		return "int", true
	case types.Int8:
		return "int8", true
	case types.Int16:
		return "int16", true
	case types.Int32:
		return "int32", true
	case types.Int64:
		return "int64", true
	case types.Uint:
		return "uint", true
	case types.Uint8:
		return "uint8", true
	case types.Uint16:
		return "uint16", true
	case types.Uint32:
		return "uint32", true
	case types.Uint64:
		return "uint64", true
	case types.Float32:
		return "float32", true
	case types.Float64:
		return "float64", true
	case types.Complex128:
		return "complex128", true
	case types.Complex64:
		return "complex64", true
	default:
		// uintptr、unsafe.Pointer等，packer不支持
		return "", false
	}
}

func basicBytes(kind types.BasicKind) (size int, ok bool) {
	switch kind {
	case types.Bool:
		return 1, true
	case types.Int8, types.Uint8:
		return 1, true
	case types.Int16, types.Uint16:
		return 2, true
	case types.Int32, types.Uint32, types.Float32:
		return 4, true
	case types.Int64, types.Uint64, types.Float64, types.Complex64:
		return 8, true
	case types.Complex128:
		return 16, true
	default:
		return -1, false
	}
}
