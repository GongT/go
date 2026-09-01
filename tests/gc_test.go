package tests

import (
	"runtime"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors/errfmt"
	"github.com/stretchr/testify/require"
)

type SomeType struct {
	value string
}

func Test_Gc(t *testing.T) {
	defer errfmt.Recover(t)
	myenv.T(t)

	var cleanupCalled bool

	someInstance := &SomeType{
		value: "test value",
	}

	runtime.AddCleanup(someInstance, func(_ *string) {
		cleanupCalled = true
	}, (*string)(nil))

	someInstance = nil
	runtime.GC()

	go func() {
		require.Equal(t, true, cleanupCalled)
	}()
}
