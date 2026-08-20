package filesys

import (
	"os"
	"path"
)

// 从startDir开始，向上展开每个目录，直到满足predicate的条件
func FindUpUntil(startDir string, predicate func(dir string) bool) (string, bool) {
	dir := startDir
	for {
		if predicate(dir) {
			return dir, true
		}

		parent := path.Dir(dir)
		if parent == dir { // 到顶了
			return "", false
		}
		dir = parent
	}
}

// 向上查找，直到找到其中一个文件，返回此文件路径
func FindUpUntilEntry(startDir string, fileNames ...string) (string, bool) {
	var idx int = -1
	dir, _ := FindUpUntil(startDir, func(dir string) bool {
		for i, f := range fileNames {
			if _, err := os.Stat(path.Join(dir, f)); err == nil {
				idx = i
				return true
			}
		}
		return false
	})

	if idx == -1 {
		return "", false
	}

	return path.Join(dir, fileNames[idx]), true
}
