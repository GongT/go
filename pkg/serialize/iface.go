package serialize

import (
	"reflect"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/serialize/internal/context"
)

// NewSerializeContext 创建一个新的序列化上下文
//
// 一次性，不可重复使用
func NewSerializeContext(metadata string) (SerializeContext, error) {
	r, err := context.NewSerialize()
	if err != nil {
		return nil, err
	}
	if err := r.Initialize(metadata); err != nil {
		return nil, err
	}
	return r, nil
}

// NewDeserializeContext 创建一个新的反序列化上下文
//
// 一次性，不可重复使用
func NewDeserializeContext() (DeserializeContext, error) {
	return context.NewDeserialize()
}

type Marshaller interface {
	TypeId() [16]byte // UUID
	Marshal(ctx SerializeContext) ([]byte, error)
	Unmarshal(data []byte, ctx DeserializeContext) error
}

type JsonMarshaller interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// ErrUnsupportedType 表示遇到了既未实现[Marshaller]也未实现[JsonMarshaller]的类型，生成代码必须将其视为运行时错误
var ErrUnsupportedType = errors.NewTemplate("类型不支持序列化，类型=%T")

type idType = uint64

type SerializeContext interface {
	// SetSendFunction 处理函数，如果为nil则出错。其中value是 “*rType” 类型的实例
	SetSendFunction(f func(id idType, value any, rType reflect.Type) any)
	// SetSendChannel 处理通道，如果为nil则出错。其中value是 “*rType” 类型的实例
	SetSendChannel(f func(id idType, value any, rType reflect.Type) any)

	// 查询ID对应的反射类型
	Lookup(id idType) (reflect.Type, bool)
}

type DeserializeContext interface {
	// SetImplementFunction 处理函数，如果为nil则出错。其中value是 “*rType” 类型的实例
	//
	// 应返回函数指针
	SetImplementFunction(f func(id idType, rType reflect.Type) any)
	// SetImplementChannel 处理通道，如果为nil则出错。其中value是 “*rType” 类型的实例
	//
	// 应返回函数指针
	SetImplementChannel(f func(id idType, rType reflect.Type) any)

	// 查询ID对应的反射类型
	Lookup(id idType) (reflect.Type, bool)
}
