package iterator

import (
	"iter"

	"github.com/gongt/go/pkg/errors/internal"
)

// [底层操作] 循环遍历一个错误的Wrap()链
//   - 遇到Join()时: 继续递归它的第一个错误
//
// @exported
func IterWrapChain(err error) iter.Seq[error] {
	return func(yield func(error) bool) {
		for err != nil {
			if !yield(err) {
				return
			}

			if unwrapper, ok := err.(internal.UnWrap); ok {
				err = unwrapper.Unwrap()
			} else if errs := extract_joined(err); len(errs) > 0 {
				err = errs[0]
			} else {
				break
			}
		}
	}
}

// [底层操作] 循环遍历一个错误的reason链
func IterReasonChain(err error) iter.Seq[error] {
	return func(yield func(error) bool) {
		for err != nil {
			if !yield(err) {
				return
			}

			if r, found := GetReason(err); found {
				err = r
			} else {
				break
			}
		}
	}
}

// [底层操作] 循环遍历一个错误的Wrap()、Join()链，
//   - 第二个yield值为Join()层级
//
// @exported
func IterWrapTree(err error) iter.Seq2[error, uint] {
	var inner func(yield func(error, uint) bool, err error, level uint) bool

	inner = func(yield func(error, uint) bool, err error, level uint) bool {
		for err != nil {
			if !yield(err, level) {
				return false
			}

			if unwrapper, ok := err.(internal.UnWrap); ok {
				err = unwrapper.Unwrap()
			} else if errs := extract_joined(err); len(errs) > 0 {
				for _, each := range errs {
					if !inner(yield, each, level+1) {
						return false
					}
				}
				break
			} else {
				break
			}
		}
		return true
	}

	return func(yield func(error, uint) bool) {
		inner(yield, err, 0)
	}
}

func extract_joined(err error) []error {
	if err == nil {
		return nil
	}
	e, ok := err.(internal.UnJoin)
	if !ok {
		return nil
	}
	return e.Unwrap()
}

// [底层操作] 遍历 Wrap()、Join() 链，并在发现附加数据时yield
//
// @exported
func IterEveryDetail(err error) iter.Seq[map[string]any] {
	return func(yield func(map[string]any) bool) {
		for err := range IterWrapTree(err) {
			if detailer, ok := err.(internal.Detailer); ok {
				details := detailer.Details()
				if len(details) > 0 {
					if !yield(details) {
						return
					}
				}
			}
		}
	}
}
