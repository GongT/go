package internal

import (
	"path/filepath"
	"slices"
)

// 从后向前找第一个绝对路径，然后在此切断
func ShrinkArguments(segments []string) ([]string, bool) {
	for i, seg := range slices.Backward(segments) {
		if filepath.IsAbs(seg) {
			return segments[i:], true
		}
	}
	return segments, false
}
