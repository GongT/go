package context

import (
	"reflect"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/serialize/internal/sbins"
)

type idType = uint64

type SerializeContext = *serializeContext
type serializeContext struct {
	counter uint64
	types   map[uint64]reflect.Type

	sendFunction func(id idType, value any, rType reflect.Type) any

	sendChannel func(id idType, value any, rType reflect.Type) any

	write sbins.PacketWrite
}

func NewSerialize() (SerializeContext, error) {
	r := &serializeContext{
		types: make(map[uint64]reflect.Type),
	}

	return r, nil
}

func (s SerializeContext) Initialize(metadata string) error {
	s.write = sbins.NewPacket(metadata)
	return nil
}

func (s SerializeContext) Wrap() SerializeContext {
	write := sbins.NewBodyField()

	return &serializeContext{
		counter:      s.counter,
		types:        s.types,
		sendFunction: s.sendFunction,
		sendChannel:  s.sendChannel,
		write:        write,
	}
}

func (s SerializeContext) SetSendFunction(f func(id idType, value any, rType reflect.Type) any) {
	if s.sendFunction != nil {
		panic("sendFunction重复设置")
	}
	s.sendFunction = f
}

func (s SerializeContext) SetSendChannel(f func(id idType, value any, rType reflect.Type) any) {
	if s.sendChannel != nil {
		panic("sendChannel重复设置")
	}
	s.sendChannel = f
}

func (s SerializeContext) Packer() sbins.PacketWrite {
	return s.write
}

// HandleChannel 为遇到的通道字段分配唯一ID，并保存其完整反射类型信息，供调用者通过[serializeContext.GetChannelType]取回
func (s SerializeContext) HandleChannel(value any, t reflect.Type) (idType, error) {
	id := s.register(t)
	if s.sendChannel != nil {
		s.sendChannel(id, value, t)
	} else {
		return id, errors.NewAnonymous("未实现通道发送功能")
	}
	return id, nil
}

// HandleFunc 为遇到的函数字段分配唯一ID，并保存其完整反射类型信息，供调用者通过[serializeContext.GetFuncType]取回
func (s SerializeContext) HandleFunc(value any, t reflect.Type) (idType, error) {
	id := s.register(t)
	if s.sendFunction != nil {
		s.sendFunction(id, value, t)
	} else {
		return id, errors.NewAnonymous("未实现函数发送功能")
	}
	return id, nil
}

func (s SerializeContext) register(t reflect.Type) idType {
	s.counter++
	id := s.counter
	s.types[id] = t

	return id
}

func (s SerializeContext) Lookup(id idType) (reflect.Type, bool) {
	t, ok := s.types[id]
	return t, ok
}
