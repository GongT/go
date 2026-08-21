package fpath

import (
	"path/filepath"
	"strings"

	"github.com/gongt/go/internal/myenv"
)

// 判断what是否在container目录下
//   - 虽然what支持相对路径，但最好不要这样，因为基于container而非cwd，后续的路径操作可能会产生意外结果
//   - container必须是绝对路径
func IsLocal[T PathLike](what T, container T) bool {
	switch v := any(what).(type) {
	case *Path:
		return isLocal1(v.Raw(), container)
	case *IPath:
		return isLocal1(v.Raw(), container)
	case string:
		return isLocal1(v, container)
	default:
		panic("unsupported type")
	}
}

func isLocal1[T PathLike](what string, _container T) bool {
	var container string
	switch v := any(_container).(type) {
	case *Path:
		container = v.Raw()
	case *IPath:
		container = v.Raw()
	case string:
		container = v
	default:
		panic("unsupported type")
	}

	if !filepath.IsAbs(container) {
		panic("container path must be absolute path")
	}

	container = filepath.Clean(container)

	if !filepath.IsAbs(what) {
		what = filepath.Join(container, what)
	} else {
		what = filepath.Clean(what)
	}

	if myenv.IsWindows {
		what = strings.ToLower(what)
		container = strings.ToLower(container)
	}

	if what == container {
		return true
	}

	return strings.HasPrefix(what, container+"/")
}
