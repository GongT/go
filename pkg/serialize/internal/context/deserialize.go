package context

import (
	"reflect"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/serialize/internal/sbins"
)

type DeserializeContext = *deserializeContext
type deserializeContext struct {
	funcs    map[uint64]reflect.Type
	channels map[uint64]reflect.Type
	metadata string

	read sbins.PacketRead

	implementFunction func(id idType, rType reflect.Type) any
	implementChannel  func(id idType, rType reflect.Type) any
}

func NewDeserialize() (DeserializeContext, error) {
	r := &deserializeContext{
		funcs:    make(map[uint64]reflect.Type),
		channels: make(map[uint64]reflect.Type),
	}

	return r, nil
}

func (d DeserializeContext) SetImplementFunction(f func(id idType, rType reflect.Type) any) {
	if d.implementFunction != nil {
		panic("implementFunction重复设置")
	}
	d.implementFunction = f
}

func (d DeserializeContext) SetImplementChannel(f func(id idType, rType reflect.Type) any) {
	if d.implementChannel != nil {
		panic("implementChannel重复设置")
	}
	d.implementChannel = f
}

func (d DeserializeContext) Metadata() string {
	return d.metadata
}

func (d DeserializeContext) Packer() sbins.PacketRead {
	return d.read
}

func (d DeserializeContext) Initialize(packet []byte) error {
	if d.read != nil {
		return errors.NewAnonymous("DeserializeContext不可重复使用")
	}
	d.read = sbins.NewPacketRead(packet)

	if meta, err := d.read.Starting(); err == nil {
		d.metadata = meta
	} else {
		return err
	}
	return nil
}

// Wrap 创建一个新的反序列化上下文，用于处理嵌套的数据体
func (d DeserializeContext) Wrap(data []byte) DeserializeContext {
	return &deserializeContext{
		funcs:    d.funcs,
		channels: d.channels,
		metadata: d.metadata,

		read: sbins.NewBodyRead(data),

		implementFunction: d.implementFunction,
		implementChannel:  d.implementChannel,
	}
}

// HandleChannel 记录反序列化时从数据流中读到的通道ID，并绑定其（由生成代码静态得知的）完整反射类型信息
func (d DeserializeContext) HandleChannel(id idType, t reflect.Type) (any, error) {
	d.channels[id] = t
	if d.implementChannel != nil {
		r := d.implementChannel(id, t)
		if myenv.IsDebug && !t.AssignableTo(reflect.TypeOf(r)) {
			panic(errors.NewAnonymous("通道实现类型不匹配").WithDetails("expected", t, "received", reflect.TypeOf(r)))
		}
		return r, nil
	}
	return nil, errors.NewAnonymous("未实现通道接收功能")
}

// HandleFunc 记录反序列化时从数据流中读到的函数ID，并绑定其（由生成代码静态得知的）完整反射类型信息
func (d DeserializeContext) HandleFunc(id idType, t reflect.Type) (any, error) {
	d.funcs[id] = t
	if d.implementFunction != nil {
		inst := d.implementFunction(id, t)
		if myenv.IsDebug && !t.AssignableTo(reflect.TypeOf(inst)) {
			panic(errors.NewAnonymous("函数实现类型不匹配").WithDetails("expected", t, "received", reflect.TypeOf(inst)))
		}
		return inst, nil
	}
	return nil, errors.NewAnonymous("未实现函数接收功能")
}

func (d DeserializeContext) Lookup(id idType) (reflect.Type, bool) {
	if t, ok := d.funcs[id]; ok {
		return t, ok
	}
	if t, ok := d.channels[id]; ok {
		return t, ok
	}
	return nil, false
}
