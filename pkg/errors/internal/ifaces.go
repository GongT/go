// @export
package internal

import (
	"github.com/gongt/go/pkg/errors/stacktrace"
)

type StackTrace interface {
	// 返回错误栈
	//
	// 可能是 nil，这个接口比较底层，通常不应该用到
	StackTrace() stacktrace.StackTraceArray
}

type UnJoin interface {
	// 返回被Join()合并的所有错误
	//
	// 长度至少为2
	Unwrap() []error
}

func IsJoined(err error) bool {
	_, ok := err.(UnJoin)
	return ok
}

type UnWrap interface {
	Unwrap() error
}

type Detailer interface {
	// 返回当前错误附带的动态数据
	//
	// 可能是nil，这个接口比较底层，通常不应该用到
	//
	// 返回的map可以修改
	Details() map[string]any
}

type E interface {
	error
	Detailer
	StackTrace
}

type MessageOverrider interface {
	OverrideMessage(string) EE
}

type DetailerWriter interface {
	WithDetails(detail_pairs ...any) EE
	SetDetails(details map[string]any) EE
	SetDetail(key string, value any) EE

	// 类似[Detailer.Details]，但不会返回nil，如果没有数据，则会初始化
	//
	// 如果已有，则等于[Detailer.Details]
	DetailsCreate() map[string]any
}

type AlsoBe interface {
	// 让当前错误对象也可以被判定为指定类型
	AlsoBe(error) EE
}

type EE interface {
	E
	DetailerWriter
	MessageOverrider
	AlsoBe
}
