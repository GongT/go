//go:generate go run golang.org/x/tools/cmd/stringer -type=ValueType -output=value_kind_generate.go

package kinds

import "go/types"

type ValueType uint16

// 基本类型

const (
	TypeIdUnknown ValueType = iota // 异常
	TypeIdBytes
	TypeIdString
	TypeIdUint
	TypeIdUint64
	TypeIdUint32
	TypeIdUint16
	TypeIdUint8
	TypeIdInt
	TypeIdInt64
	TypeIdInt32
	TypeIdInt16
	TypeIdInt8
	TypeIdByte
	TypeIdBool
	TypeIdFloat64
	TypeIdFloat32
	TypeIdComplex128
	TypeIdComplex64

	TypeIdJson
	TypeIdFunc
	TypeIdChannel
	TypeIdStruct

	TypeIdArray
	TypeIdMap
)

func IdOfType(typ types.BasicKind) string {
	switch typ {
	case types.Bool:
		return "TypeIdBool"
	case types.Int:
		return "TypeIdInt"
	case types.Int8:
		return "TypeIdInt8"
	case types.Int16:
		return "TypeIdInt16"
	case types.Int32:
		return "TypeIdInt32"
	case types.Int64:
		return "TypeIdInt64"
	case types.Uint:
		return "TypeIdUint"
	case types.Uint8:
		return "TypeIdUint8"
	case types.Uint16:
		return "TypeIdUint16"
	case types.Uint32:
		return "TypeIdUint32"
	case types.Uint64:
		return "TypeIdUint64"
	case types.Float32:
		return "TypeIdFloat32"
	case types.Float64:
		return "TypeIdFloat64"
	case types.Complex64:
		return "TypeIdComplex64"
	case types.Complex128:
		return "TypeIdComplex128"
	case types.String:
		return "TypeIdString"
	default:
		return "TypeIdUnknown"
	}
}

func SizeOfType(typ ValueType) int {
	switch typ {
	case TypeIdBool:
		return 1
	case TypeIdInt8, TypeIdUint8, TypeIdByte:
		return 1
	case TypeIdInt16, TypeIdUint16:
		return 2
	case TypeIdInt32, TypeIdUint32, TypeIdFloat32:
		return 4
	case TypeIdInt, TypeIdUint, TypeIdInt64, TypeIdUint64, TypeIdFloat64, TypeIdComplex64:
		return 8
	case TypeIdComplex128:
		return 16
	default:
		return -1
	}
}
