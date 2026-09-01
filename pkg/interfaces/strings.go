package interfaces

// StringLike 是一个泛型接口，表示可以被视为字符串的类型。表达此对象要被用作字符串。
type StringLike = ByteSeq

// ByteSeq 是一个泛型接口，表示可以被视为字节序列的类型。表达此对象要被用作字节序列。
type ByteSeq interface {
	~string | ~[]byte
}

type ToBytes interface {
	Bytes() []byte
}

type BuiltinLiterials interface {
	string | bool | int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8 | float64 | float32 | complex128 | complex64
}

type CustomLiterials interface {
	~string | ~bool | ~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8 | ~float64 | ~float32 | ~complex128 | ~complex64
}
