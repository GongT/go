package myenv

import "errors"

// Must 如果 err 不为 nil，则触发 panic
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must1 如果 err 不为 nil，则触发 panic，否则返回 v
func Must1[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Must2 如果 err 不为 nil，则触发 panic，否则返回 v1, v2
func Must2[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

func NotImplemented() {
	panic(errors.New("not implemented"))
}
