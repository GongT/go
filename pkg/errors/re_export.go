package errors

import stderrors "errors"

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

func Is(err error, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target any) bool {
	return stderrors.As(err, target)
}

func AsType[E error](err error) (E, bool) {
	return stderrors.AsType[E](err)
}

var ErrUnsupported = stderrors.ErrUnsupported

func NewStd(message string) error {
	return stderrors.New(message)
}
