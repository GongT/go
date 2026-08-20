package errors

import (
	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/stacktrace"
)

// 遍历错误链，找到第一个实现了StackTrace接口（且返回不是nil）的错误，返回其堆栈信息
func GetStackTrace(err error) stacktrace.StackTraceArray {
	itr := err
	for itr != nil {
		// 首先检查join
		if unjoin, ok := itr.(internal.UnJoin); ok {
			errs := unjoin.Unwrap()
			if len(errs) > 0 {
				// 递归检查第一个错误
				found_root_cause := GetStackTrace(errs[0])
				if len(found_root_cause) > 0 {
					return found_root_cause
				}
				// 找不到则使用join自身的stacktrace
			}
		}

		if tracer, ok := itr.(internal.StackTrace); ok {
			// 找到！
			if childStack := tracer.StackTrace(); len(childStack) > 0 {
				return childStack
			}
		}

		// 没找到，继续检查
		if unwrapper, ok := itr.(internal.UnWrap); ok {
			itr = unwrapper.Unwrap()
		} else {
			break // 整个链都没有
		}
	}
	return nil
}

func Unjoin(err error) []error {
	if e, ok := err.(internal.UnJoin); ok {
		return e.Unwrap()
	}
	return nil
}

func IterateErrorStack(err error) stacktrace.Iter {
	return stacktrace.IterateStack(GetStackTrace(err))
}
