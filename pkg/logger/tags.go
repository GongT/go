package logger

import (
	"github.com/gongt/go/pkg/logger/internal/tags"
)

type DebugTag string

func (t DebugTag) IsEnabled() bool {
	return tags.CheckEnabled(tags.DebugTag(t))
}

func (t DebugTag) Enable() bool {
	return tags.Enable(tags.DebugTag(t))
}

func (t DebugTag) Disable() bool {
	return tags.Disable(tags.DebugTag(t))
}

func IsEnabled(tag string) bool {
	return tags.CheckEnabled(tags.DebugTag(tag))
}

func Enable(setting string) bool {
	return tags.Enable(tags.DebugTag(setting))
}

func Disable(setting string) bool {
	return tags.Disable(tags.DebugTag(setting))
}
