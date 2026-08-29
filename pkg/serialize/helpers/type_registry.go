package helpers

import (
	"reflect"
)

type uuid = [16]byte

var typeRegistry map[uuid]reflect.Type

func RegisterType(guid uuid, ptrType reflect.Type) {
	if typeRegistry == nil {
		typeRegistry = make(map[uuid]reflect.Type)
	}
	typeRegistry[guid] = ptrType
}

func GetMarshaller(guid uuid) (reflect.Type, bool) {
	if typeRegistry == nil {
		return nil, false
	}
	m, ok := typeRegistry[guid]
	return m, ok
}
