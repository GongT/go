package helpers

import (
	"reflect"

	"github.com/gongt/go/pkg/serialize/internal/context"
	"github.com/gongt/go/pkg/serialize/internal/sbins"
	"github.com/gongt/go/pkg/serialize/kinds"

	. "github.com/gongt/go/pkg/serialize"
)

// Marshal 尝试调用Marshal和MarshalJSON方法进行序列化
//
// 如果value既未实现[Marshaller]也未实现[JsonMarshaller]，返回(nil, nil)，调用者需自行判断data==nil并报错
//
// ctx 可以为nil
func HelperMarshal(value any, ctx SerializeContext) error {
	if m, ok := value.(Marshaller); ok {
		if data, err := m.Marshal(ctx); err != nil {
			return err
		} else {
			ser(ctx).Packer().WriteStruct(m.TypeId(), data)
		}
	} else if m, ok := value.(JsonMarshaller); ok {
		if data, err := m.MarshalJSON(); err != nil {
			return err
		} else {
			ser(ctx).Packer().WriteAnyBytes(kinds.TypeIdJson, data)
		}
	}

	return nil
}

// Unmarshal 尝试调用Unmarshal和UnmarshalJSON方法进行反序列化
//
// 如果receiver既未实现[Marshaller]也未实现[JsonMarshaller]，返回[ErrUnsupportedType]
//
// ctx 可以为nil
func HelperUnmarshal(receiver any, ctx DeserializeContext) error {
	if m, ok := receiver.(Marshaller); ok {
		bs, err := de(ctx).Packer().ReadStruct(m.TypeId())
		if err != nil {
			return err
		}
		return m.Unmarshal(bs, ctx)
	} else if m, ok := receiver.(JsonMarshaller); ok {
		if data, err := de(ctx).Packer().ReadAnyBytes(kinds.TypeIdJson); err != nil {
			return err
		} else {
			return m.UnmarshalJSON(data)
		}
	}
	return ErrUnsupportedType.New(receiver)
}

func EnsureUnmarshalContext(data []byte, ctx DeserializeContext) (DeserializeContext, error) {
	if ctx == nil {
		var err error
		if ctx, err = context.NewDeserialize(); err != nil {
			return nil, err
		}
		if err := de(ctx).Initialize(data); err != nil {
			return nil, err
		}
	} else {
		return de(ctx).Wrap(data), nil
	}
	return ctx, nil
}

func EnsureMarshalContext(ctx SerializeContext) (SerializeContext, error) {
	if ctx == nil {
		var err error
		if ctx, err = context.NewSerialize(); err != nil {
			return nil, err
		}
		if err := ser(ctx).Initialize(""); err != nil {
			return nil, err
		}
	} else {
		return ser(ctx).Wrap(), nil
	}
	return ctx, nil
}

func de(ctx DeserializeContext) context.DeserializeContext {
	return ctx.(context.DeserializeContext)
}
func ser(ctx SerializeContext) context.SerializeContext {
	return ctx.(context.SerializeContext)
}

func UseD(ctx DeserializeContext) sbins.PacketRead {
	return de(ctx).Packer()
}

func UseS(ctx SerializeContext) sbins.PacketWrite {
	return ser(ctx).Packer()
}

type idType = uint64

func ReceiveChannel(d DeserializeContext, id idType, t reflect.Type) (any, error) {
	return de(d).HandleChannel(id, t)
}

func ReceiveFunc(d DeserializeContext, id idType, t reflect.Type) (any, error) {
	return de(d).HandleFunc(id, t)
}

func SendChannel(d SerializeContext, value any, t reflect.Type) (idType, error) {
	return ser(d).HandleChannel(value, t)
}

func SendFunc(d SerializeContext, value any, t reflect.Type) (idType, error) {
	return ser(d).HandleFunc(value, t)
}
