package context

import "reflect"

type SerializeContext = *serializeContext
type serializeContext struct{}

func NewSerialize() SerializeContext {
	return &serializeContext{}
}

func (s SerializeContext) HandleChannel(direction reflect.ChanDir) []byte {
	// TODO: 记录此通道的方向信息，并返回（本对象）唯一ID
	return nil
}

func (s SerializeContext) HandleFunc(argumentCount int, returnCount int) []byte {
	// TODO: 记录此函数的参数和返回值数量，并返回（本对象）唯一ID
	return nil
}
