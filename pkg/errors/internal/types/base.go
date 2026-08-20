package types

import (
	"maps"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/internal/tools"
	"github.com/gongt/go/pkg/errors/stacktrace"
)

var _ internal.EE = (*ErrorObjectBase)(nil)

type ErrorObjectBase struct {
	err error

	overrideMessage string
	stack           stacktrace.StackTraceArray
	details         map[string]any
}

func (e *ErrorObjectBase) StackTrace() stacktrace.StackTraceArray {
	return e.stack
}

func (e *ErrorObjectBase) Details() map[string]any {
	return e.details
}

func (e *ErrorObjectBase) DetailsCreate() map[string]any {
	if e.details == nil {
		e.details = make(map[string]any)
	}
	return e.details
}

func (e *ErrorObjectBase) WithDetails(detail_pairs ...any) internal.EE {
	if e == nil {
		return nil
	}
	d := e.DetailsCreate()
	tools.ApplyPairsToMap(d, detail_pairs)
	return e
}

func (e *ErrorObjectBase) SetDetails(details map[string]any) internal.EE {
	if e == nil {
		return nil
	}
	d := e.DetailsCreate()
	maps.Copy(d, details)
	return e
}

func (e *ErrorObjectBase) SetDetail(key string, value any) internal.EE {
	if e == nil {
		return nil
	}
	d := e.DetailsCreate()
	d[key] = value
	return e
}

func (e *ErrorObjectBase) Error() string {
	if e == nil {
		return "<nil error>"
	}
	if e.overrideMessage != "" {
		return e.overrideMessage
	}
	return e.err.Error()
}

func (e *ErrorObjectBase) Unwrap() error {
	return e.err
}

// 内部接口: 复写Error()的返回值
func (e *ErrorObjectBase) OverrideMessage(newMessage string) internal.EE {
	e.overrideMessage = newMessage
	return e
}

// 内部接口: 返回一个可修改的details map
func (e *ErrorObjectBase) InitDetails() map[string]any {
	if e.details == nil {
		e.details = make(map[string]any)
	}
	return e.details
}

// --------------------------------------------------------------------
// 静态创建方法

func Create(err error, add int) *ErrorObjectBase {
	myenv.Assert(err != nil, "Create: err is nil pointer")
	stack := stacktrace.CaptureStackTrace(add+1, 0)
	return &ErrorObjectBase{
		err:   err,
		stack: stack,
	}
}

func CreateEnsureStackTrace(err error, add int) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(internal.StackTrace); ok {
		return err
	}

	return Create(err, add+1)
}
