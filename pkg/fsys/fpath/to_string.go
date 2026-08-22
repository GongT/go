package fpath

import (
	"os"
	"path/filepath"

	"github.com/gongt/go/pkg/fsys/fpath/internal"
)

type PathLike interface {
	*Path | *IPath | string
}
type PathLikeRo interface {
	*IPath | string
}

func ToString[T PathLike](input T) (r string) {
	switch v := any(input).(type) {
	case *Path:
		r = v.String()
	case *IPath:
		r = v.String()
	case string:
		r = v
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
	internal.AssertValidPath(r)
	return
}

func ToPath[T PathLike](input T) *Path {
	switch v := any(input).(type) {
	case *Path:
		return v
	case *IPath:
		return v.Mutate()
	case string:
		return New(v)
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
}

func ToImmutable[T PathLike](input T) *IPath {
	switch v := any(input).(type) {
	case *Path:
		return v.Immutable()
	case *IPath:
		return v
	case string:
		return INew(v)
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
}

func toStringAssertAbs[T PathLike](input T) string {
	t := ToString(input)
	if !filepath.IsAbs(t) {
		panic(PathErr.New("传入的参数必须是绝对路径"))
	}
	t = filepath.Clean(t)
	return t
}

func toStringMustAbs[T PathLike](input T, base string) string {
	t := ToString(input)
	if !filepath.IsAbs(t) {
		t = filepath.Join(base, t)
	} else {
		t = filepath.Clean(t)
	}
	return t
}

func toStringAbsCwd2[T1 PathLike, T2 PathLike](a T1, b T2) (*Path, *Path) {
	aPath := ToPath(a)
	bPath := ToPath(b)

	aAbs := aPath.IsAbs()
	bAbs := bPath.IsAbs()

	if aAbs && bAbs {
		return aPath, bPath
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(PathErr.Wrap(err))
	}

	if aAbs && !bAbs {
		return aPath, New(cwd + "/" + bPath.Raw())
	}
	if !aAbs && bAbs {
		return New(cwd + "/" + aPath.Raw()), bPath
	}
	return New(cwd + "/" + aPath.Raw()), New(cwd + "/" + bPath.Raw())
}
