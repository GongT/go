//go:build release

package myenv

const IsDebug = false
const IsRelease = true
const IsTesting = false

func Assert(condition bool, message string, args ...any) {}
