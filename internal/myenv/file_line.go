package myenv

import (
	"errors"
	"runtime"
)

type lowError struct {
	error

	trace []uintptr
}

func NewErr(message string) *lowError {
	trace := make([]uintptr, 32)
	n := runtime.Callers(1, trace)
	trace = trace[:n]

	return &lowError{
		error: errors.New(message),
		trace: trace,
	}
}

func (e *lowError) StackTrace() []uintptr {
	return e.trace
}

func CurrentFileLine() (string, int) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic(NewErr("无法在运行时获取当前文件路径和行号"))
	}
	return file, line
}
