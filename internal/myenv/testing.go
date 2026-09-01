package myenv

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func RedirectDebugOutput(output io.Writer) {
	if !testing.Testing() {
		panic("RedirectDebugOutput仅用于go test测试")
	}
	log.SetOutput(output)
}

func RedirectDebugTesting(t testing.TB) *Tester {
	if !testing.Testing() {
		panic("RedirectDebugTesting仅用于go test测试")
	}
	RedirectDebugOutput(t.Output())
	return &Tester{TB: t}
}

func T(t testing.TB) *Tester {
	return RedirectDebugTesting(t)
}

type Tester struct {
	testing.TB
}

func (tester *Tester) Must[T any](value T, err error) T {
	require.NoError(tester.TB, err)
	return value
}
