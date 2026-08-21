package iterator

import (
	"iter"

	"github.com/gongt/go/pkg/errors"
)

// 错误回调，如果返回 true，则继续迭代，返回 false，则停止迭代
type ErrorHandler = func(err error) bool

type Yield[T any] = func(T, error) bool
type iterFuncErr[T any] = func(Yield[T])

type Yield2[T1 any, T2 any] = func(T1, T2, error) bool
type iterFuncErr2[T1 any, T2 any] = func(Yield2[T1, T2])

func defaultErrorHandler(_ error) bool {
	return false
}

// CreateIterator 创建一个迭代器，遇到错误时调用error_handler处理
//
// 默认的error_handler会直接停止迭代，忽略发生的错误
func CreateIterator[T any](worker iterFuncErr[T], errorHandler ErrorHandler) iter.Seq[T] {
	if errorHandler == nil {
		errorHandler = defaultErrorHandler
	}
	return func(yield func(T) bool) {
		worker(func(v T, err error) bool {
			if err != nil {
				// 处理错误，如果返回 false，则停止迭代
				if !errorHandler(err) {
					return false
				}
			}
			if !yield(v) {
				return false
			}
			return true
		})
	}
}

func CreateIterator2[T1 any, T2 any](worker iterFuncErr2[T1, T2], errorHandler ErrorHandler) iter.Seq2[T1, T2] {
	if errorHandler == nil {
		errorHandler = defaultErrorHandler
	}
	return func(yield func(T1, T2) bool) {
		worker(func(v1 T1, v2 T2, err error) bool {
			if err != nil {
				// 处理错误，如果返回 false，则停止迭代
				if !errorHandler(err) {
					return false
				}
			}
			if !yield(v1, v2) {
				return false
			}
			return true
		})
	}
}

type errorHolder struct {
	error
}

func (e *errorHolder) Get() error {
	return e.error
}

// 遇到第一个错误时停止，记录此错误
func RecordFirstErrorBreak() (ErrorHandler, *errorHolder) {
	e := &errorHolder{}
	return func(err error) bool {
		e.error = err
		return false
	}, e
}

// 迭代器会继续迭代，只记录第一个遇到的错误
func RecordFirstErrorContinue() (ErrorHandler, *errorHolder) {
	e := &errorHolder{}
	return func(err error) bool {
		if e.error == nil {
			e.error = err
		}
		return true
	}, e
}

// 迭代器会继续迭代，记录最后一个遇到的错误
func RecordLastErrorContinue() (ErrorHandler, *errorHolder) {
	e := &errorHolder{}
	return func(err error) bool {
		e.error = err
		return true
	}, e
}

type errorsHolder struct {
	errors []error
}

func (e *errorsHolder) Get() error {
	return errors.Concat(e.errors, "迭代器产生了多个错误")
}

// 一直运行，并Join所有的错误
func RecordAllErrorsContinue() (ErrorHandler, *errorsHolder) {
	e := &errorsHolder{}
	return func(err error) bool {
		e.errors = append(e.errors, err)
		return true
	}, e
}
