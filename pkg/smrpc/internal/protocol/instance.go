package protocol

type StructVar struct {
	TypeId uint
	Values []any
}

type Channel struct {
}

type Callback struct {
}

type TypeId int

const (
	TypeString TypeId = iota
	TypeInt
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeFloat32
	TypeFloat64
	TypeBool
)

type BasicVar struct {
	TypeId TypeId
	Value  any
}

type Variable interface {
	*StructVar | *Channel | *Callback
}

type Instance struct {
	FuncId  uint
	Inputs  []StructVar
	Outputs []StructVar
}
