package types

import (
	"errors"
	"strings"

	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/stacktrace"
)

var _ internal.EE = (*ErrorObjectJoined)(nil)
var _ internal.UnJoin = (*ErrorObjectJoined)(nil)

type ErrorObjectJoined struct {
	self ErrorObjectBase

	errs []error
}

func (e *ErrorObjectJoined) Unwrap() []error {
	return e.errs
}

func (e *ErrorObjectJoined) Error() string {
	sb := strings.Builder{}
	sb.WriteString(e.self.Error())
	sb.WriteByte(':')
	for _, err := range e.errs {
		sb.WriteString("\n  - ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (e *ErrorObjectJoined) OverrideMessage(message string) internal.EE {
	e.self.OverrideMessage(message)
	return e
}
func (e *ErrorObjectJoined) StackTrace() stacktrace.StackTraceArray {
	return e.self.StackTrace()
}
func (e *ErrorObjectJoined) Details() map[string]any {
	return e.self.Details()
}
func (e *ErrorObjectJoined) DetailsCreate() map[string]any {
	return e.self.DetailsCreate()
}
func (e *ErrorObjectJoined) WithDetails(detail_pairs ...any) internal.EE {
	e.self.WithDetails(detail_pairs...)
	return e
}
func (e *ErrorObjectJoined) SetDetails(details map[string]any) internal.EE {
	e.self.SetDetails(details)
	return e
}
func (e *ErrorObjectJoined) SetDetail(key string, value any) internal.EE {
	e.self.SetDetail(key, value)
	return e
}

var default_title = "发生多个错误"

// 合并多个同级别的错误
//
// errs中的nil会被忽略，如果至少有2个非nil，返回一个ErrorObjectJoined，否则返回errs中唯一的非nil错误，全nil则返回nil
func Join(add int, errs []error, force_message bool) internal.EE {
	return JoinMessage(add+1, default_title, errs, force_message)
}

func JoinMessage(add int, msg string, errs []error, force_message bool) internal.EE {
	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	} else if len(nonNilErrs) == 1 {
		if force_message {
			return Wrap(add+1, nonNilErrs[0], false, msg, nil)
		} else if r, ok := nonNilErrs[0].(internal.EE); ok {
			return r
		} else {
			return Create(nonNilErrs[0], add+1)
		}
	}

	e := &ErrorObjectJoined{
		self: ErrorObjectBase{
			err:   errors.New(msg),
			stack: stacktrace.CaptureStackTrace(add+1, 0),
		},
		errs: nonNilErrs,
	}

	return e
}
