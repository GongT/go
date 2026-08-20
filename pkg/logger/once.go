package logger

import "github.com/gongt/go/pkg/logger/internal/once"

func DebugLogOnce(format string, v ...interface{}) {
	once.DebugLogOnce(1, format, v...)
}

func AlertOnce(format string, v ...interface{}) {
	once.AlertOnce(1, format, v...)
}
