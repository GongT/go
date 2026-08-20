package types

import (
	stderrors "errors"
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
	alsoBeTypes     []error
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

// Is被[stderrors.Is]调用
//   - 在 e == target == false 判断之后
//   - 在 Unwrap() 之前
func (e *ErrorObjectBase) Is(target error) bool {
	if e == nil || target == nil {
		return false
	}
	for _, t := range e.alsoBeTypes {
		if stderrors.Is(target, t) {
			// 几乎完全等同于 target==t 但不能完全保证，所以还是调用一下 Is
			// 此处target和t的先后顺序没有合理性，无法真正确定
			// 因为理论上target和t都是裸的error对象，类似os.NotExist这种
			return true
		}
	}
	return false
}

func (e *ErrorObjectBase) AlsoBe(t error) internal.EE {
	e.alsoBeTypes = append(e.alsoBeTypes, t)
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
