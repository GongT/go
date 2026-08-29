package type_name

import (
	"go/types"
	"reflect"
	"strconv"
	"strings"
)

func PointerLevel(t reflect.Type) int {
	level := 0
	for t.Kind() == reflect.Pointer {
		level++
		t = t.Elem()
	}
	return level
}

func TranslateInterfaceType(target any) string {
	reflectType := reflect.TypeOf(target)
	return TranslateType(reflectType)
}
func TranslateBasicType(targetType *types.Basic) string {
	switch targetType.Kind() {
	case types.Bool:
		return "布尔值"
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		return name_number(type_bits(targetType), true, 0)
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		return name_number(type_bits(targetType), false, 0)
	case types.Float32, types.Float64:
		return name_number(type_bits(targetType), true, 1)
	case types.Complex64, types.Complex128:
		return name_number(type_bits(targetType), true, 2)
	case types.String:
		return "字符串"
	case types.UnsafePointer:
		return "任意指针"
	case types.Invalid:
		return "无效类型"
	case types.UntypedBool, types.UntypedInt, types.UntypedRune, types.UntypedFloat, types.UntypedComplex, types.UntypedString, types.UntypedNil:
		return "无法推断类型"
	default:
		return "(" + targetType.String() + ")"
	}
}

func type_bits(t *types.Basic) int {
	switch t.Kind() {
	case types.Int:
		return 0
	case types.Int8:
		return 8
	case types.Int16:
		return 16
	case types.Int32:
		return 32
	case types.Int64:
		return 64
	case types.Uint:
		return 0
	case types.Uint8:
		return 8
	case types.Uint16:
		return 16
	case types.Uint32:
		return 32
	case types.Uint64:
		return 64
	case types.Float32:
		return 32
	case types.Float64:
		return 64
	case types.Complex64:
		return 64
	case types.Complex128:
		return 128
	default:
		return -1
	}
}

func name_number(bits int, signed bool, t byte) string {
	sb := strings.Builder{}
	switch bits {
	case 0:
		// pass
	case 8:
		sb.WriteByte('8')
	case 16:
		sb.WriteString("16")
	case 32:
		sb.WriteString("32")
	case 64:
		sb.WriteString("64")
	case 128:
		sb.WriteString("128")
	default:
		sb.WriteString(strconv.Itoa(bits))
	}
	sb.WriteRune('位')
	if !signed {
		sb.WriteString("无符号")
	}
	switch t {
	case 0:
		sb.WriteRune('整')
	case 1:
		sb.WriteString("浮点")
	case 2:
		sb.WriteRune('复')
	}
	sb.WriteString("数")
	return sb.String()
}

func TranslateType(targetType reflect.Type) string {
	switch targetType.Kind() {
	case reflect.Bool:
		return "布尔值"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return name_number(targetType.Bits(), true, 0)
	case reflect.Uint8:
		return "字节"
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return name_number(targetType.Bits(), false, 0)
	case reflect.Float32, reflect.Float64:
		return name_number(targetType.Bits(), true, 1)
	case reflect.Complex64, reflect.Complex128:
		return name_number(targetType.Bits(), true, 2)
	case reflect.String:
		return "字符串"
	case reflect.Slice:
		return TranslateType(targetType.Elem()) + "数组"
	case reflect.Array:
		return "定长" + TranslateType(targetType.Elem()) + "数组(" + strconv.Itoa(targetType.Len()) + ")"
	case reflect.Map:
		return "(" + TranslateType(targetType.Key()) + "->" + TranslateType(targetType.Elem()) + ")映射"
	case reflect.Interface:
		return "任意接口"
	case reflect.Struct:
		return "结构体"
	case reflect.Uintptr:
		return "C指针"
	case reflect.UnsafePointer:
		return "任意指针"
	case reflect.Pointer:
		level := PointerLevel(targetType)
		final := indirectType(targetType)
		if level == 1 {
			return TranslateType(final) + "指针"
		}
		return strconv.Itoa(level) + "级" + TranslateType(final) + "指针"
	case reflect.Func:
		return "函数指针"
	case reflect.Chan:
		switch targetType.ChanDir() {
		case reflect.SendDir:
			return "发送" + TranslateType(targetType.Elem()) + "通道"
		case reflect.RecvDir:
			return "接收" + TranslateType(targetType.Elem()) + "通道"
		default:
			return "双向" + TranslateType(targetType.Elem()) + "通道"
		}
	default:
		return "(" + targetType.String() + ")"
	}
}

func indirectType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}
