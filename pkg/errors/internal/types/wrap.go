package types

import (
	"fmt"

	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/stacktrace"
)

var _ internal.EE = (*ErrorObjectWrapped)(nil)
var _ internal.UnWrap = (*ErrorObjectWrapped)(nil)

type ErrorObjectWrapped struct {
	base ErrorObjectBase

	context string
}

func (e *ErrorObjectWrapped) Unwrap() error {
	return &e.base
}

func (e *ErrorObjectWrapped) Error() string {
	if e.context == "" {
		return e.base.Error()
	}
	return fmt.Sprintf("%s: %s", e.context, e.base.Error())
}

func (e *ErrorObjectWrapped) OverrideMessage(message string) internal.EE {
	e.base.OverrideMessage(message)
	return e
}
func (e *ErrorObjectWrapped) StackTrace() stacktrace.StackTraceArray {
	return e.base.StackTrace()
}
func (e *ErrorObjectWrapped) Details() map[string]any {
	return e.base.Details()
}
func (e *ErrorObjectWrapped) DetailsCreate() map[string]any {
	return e.base.DetailsCreate()
}
func (e *ErrorObjectWrapped) WithDetails(detail_pairs ...any) internal.EE {
	e.base.WithDetails(detail_pairs...)
	return e
}
func (e *ErrorObjectWrapped) SetDetails(details map[string]any) internal.EE {
	e.base.SetDetails(details)
	return e
}
func (e *ErrorObjectWrapped) SetDetail(key string, value any) internal.EE {
	e.base.SetDetail(key, value)
	return e
}

func (e *ErrorObjectWrapped) Is(target error) bool {
	return e.base.Is(target)
}

func (e *ErrorObjectWrapped) AlsoBe(t error) internal.EE {
	e.base.AlsoBe(t)
	return e
}

func Wrap(stackAdd int, err error, forceStack bool, msg string, args []any) *ErrorObjectWrapped {
	if err == nil {
		return nil
	}

	var stack stacktrace.StackTraceArray
	if forceStack {
		stack = stacktrace.CaptureStackTrace(1+stackAdd, 0)
	} else if tracer, ok := err.(internal.StackTrace); ok {
		stack = tracer.StackTrace()
	} else {
		stack = stacktrace.CaptureStackTrace(1+stackAdd, 0)
	}

	if len(args) != 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return &ErrorObjectWrapped{
		base: ErrorObjectBase{
			err:   err,
			stack: stack,
		},
		context: msg,
	}
}
