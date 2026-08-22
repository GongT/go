package fsys

import (
	"iter"
	"os"

	"github.com/gongt/go/pkg/fsys/fpath"
)

// 从startDir开始，向上展开每个目录，直到满足predicate的条件
func FindUpUntil[T fpath.PathLike](startDir T, predicate func(dir *fpath.IPath) bool) (*fpath.IPath, bool) {
	for dir := range ClimbingPathPhysical(startDir) {
		if predicate(dir) {
			return dir, true
		}
	}

	return nil, false
}

// 向上查找，直到找到其中一个文件，返回此文件路径
func FindUpUntilEntry[T fpath.PathLike](startDir T, fileNames ...string) (*fpath.IPath, bool) {
	var found *fpath.IPath
	FindUpUntil(startDir, func(dir *fpath.IPath) bool {
		entities, err := os.ReadDir(dir.Raw())
		if err != nil {
			return false
		}
		for _, f := range fileNames {
			for _, entry := range entities {
				if entry.Name() == f {
					found = dir.Join(f)
					return true
				}
			}
		}
		return false
	})

	if found == nil {
		return nil, false
	}

	return found, true
}

// 先realpath然后向上展开
func ClimbingPathPhysical[T fpath.PathLike](dir T) iter.Seq[*fpath.IPath] {
	iter := fpath.ToPath(dir)
	if err := iter.RealpathMissing(); err != nil {
		iter.Normalize()
	}
	return ClimbingPath(iter)
}

// 向上逐层展开文件夹路径，包括传入的目录本身，直到根目录
//
// 不要求任何目录存在。无规范化，纯字符串操作
func ClimbingPath[T fpath.PathLike](dir T) iter.Seq[*fpath.IPath] {
	iter := fpath.ToPath(dir)
	return func(yield func(*fpath.IPath) bool) {
		current := iter.Raw()
		for {
			if !yield(iter.Immutable()) {
				return
			}

			iter.Dir()

			next := iter.Raw()
			if next == current { // 调用LogicalDir()后，路径不变，说明已经到达根目录
				return
			}
			current = next
		}
	}
}
