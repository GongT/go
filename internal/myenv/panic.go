package myenv

// MustNil 如果 err 不为 nil，则触发 panic
func MustNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Must 如果 err 不为 nil，则触发 panic，否则返回 v
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
