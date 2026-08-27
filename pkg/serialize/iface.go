package serialize

import "github.com/gongt/go/pkg/serialize/internal/context"

type SerializeContext = context.SerializeContext
type DeserializeContext = context.DeserializeContext
type ValueType int

const (
	ValueTypeValue ValueType = iota
	ValueTypeChannel
	ValueTypeFunc
)

type Serializer = func(ctx SerializeContext, value any) error
type Deserializer = func(ctx DeserializeContext) (any, error)

type Marshaller interface {
	Marshal(ctx SerializeContext) ([]byte, error)
	Unmarshal(ctx SerializeContext, data []byte) error
}

type JsonMarshaller interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// Marshal 尝试调用Marshal和MarshalJSON方法进行序列化
func Marshal(ctx SerializeContext, value any) ([]byte, error) {
	if m, ok := value.(Marshaller); ok {
		return m.Marshal(ctx)
	}
	if m, ok := value.(JsonMarshaller); ok {
		return m.MarshalJSON()
	}
	return nil, nil
}

// Unmarshal 尝试调用Unmarshal和UnmarshalJSON方法进行反序列化
func Unmarshal(ctx SerializeContext, data []byte, value any) error {
	if m, ok := value.(Marshaller); ok {
		return m.Unmarshal(ctx, data)
	}
	if m, ok := value.(JsonMarshaller); ok {
		return m.UnmarshalJSON(data)
	}
	return nil
}
