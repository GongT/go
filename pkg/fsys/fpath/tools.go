package fpath

import (
	"path/filepath"
	"strings"

	"github.com/gongt/go/internal/myenv"
)

func Clean[T PathLike](p T) *Path {
	return New(filepath.ToSlash(filepath.Clean(ToString(p))))
}

// 判断what是否在container目录下
//   - 虽然what支持相对路径，但最好不要这样，因为基于container而非cwd，后续的路径操作可能会产生意外结果
//   - container必须是绝对路径
func IsLocal[T1 PathLike, T2 PathLike](what_a T1, container_a T2) bool {
	container := toStringAssertAbs(container_a)
	what := toStringMustAbs(what_a, container)

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

// 将what相对于container的路径返回
//
//   - 虽然what支持相对路径，但最好不要这样，因为基于container而非cwd，后续的路径操作可能会产生意外结果
//
//   - container必须是绝对路径
//
//     ToRelative("/a/b/c/d", "/a/b") => "./c/d"
//     ToRelative("/a/b", "/a/b/c/d") => "../.."
//     ToRelative("/a/b", "/a/b") => "."
//     ToRelative("/a/b/c/d", "/x/y") => "../../a/b/c/d"
func ToRelative[T1 PathLike, T2 PathLike](what_a T1, container_a T2) (string, error) {
	container := toStringAssertAbs(container_a)
	what := toStringMustAbs(what_a, container)

	rel, err := filepath.Rel(container, what)
	if err != nil {
		return "", PathErr.Wrap(err).WithDetails("what", what, "container", container)
	}

	if !strings.HasPrefix(rel, ".") {
		rel = "./" + rel
	}

	return rel, nil
}

func MustRelative[T1 PathLike, T2 PathLike](what T1, container T2) string {
	if s, err := ToRelative(what, container); err != nil {
		panic(PathErr.Wrap(err).WithDetails("what", what, "container", container))
	} else {
		return s
	}
}

// 判断两个路径是否逻辑相等
//
// 自动基于cwd转换为绝对路径
func IsEquals[T1 PathLike, T2 PathLike](a_ T1, b_ T2) bool {
	a, b := toStringAbsCwd2(a_, b_)

	if a.value == b.value {
		return true
	}

	if myenv.IsWindows {
		a.value = strings.ToLower(a.value)
		b.value = strings.ToLower(b.value)
	}

	return a.Normalize().Raw() == b.Normalize().Raw()
}

// 判断两个路径是否实际相同
//
// 自动基于cwd转换为绝对路径
func IsSameEntity[T1 PathLike, T2 PathLike](a_ T1, b_ T2) bool {
	a, b := toStringAbsCwd2(a_, b_)

	if a.value == b.value {
		return true
	}

	a.RealpathMissing()
	b.RealpathMissing()

	return a.Raw() == b.Raw()
}
