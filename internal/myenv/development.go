//go:build !release

package myenv

import (
	"fmt"
	"testing"
)

const IsDebug = true
const IsRelease = false

var IsTesting = testing.Testing()

func Assert(condition bool, message string, args ...any) {
	if !condition {
		panic(fmt.Sprintf(message, args...))
	}
}
