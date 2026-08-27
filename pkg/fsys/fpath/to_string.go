package fpath

import (
	"os"
	"path/filepath"

	"github.com/gongt/go/pkg/fsys/fpath/internal"
)

type PathLikeRo interface {
	*IPath | IPath | ~string
}

type PathLike interface {
	*Path | Path | PathLikeRo
}

type fileLike interface {
	*File | File | ~string
}

type pathOrFileLike interface {
	fileLike | PathLike
}

func ToString[T pathOrFileLike](input T) (r string) {
	switch v := any(input).(type) {
	case *Path:
		r = v.Raw()
	case Path:
		r = v.Raw()
	case *IPath:
		r = v.Raw()
	case IPath:
		r = v.Raw()
	case *File:
		r = v.Name
	case File:
		r = v.Name
	case string:
		r = string(v)
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
	internal.AssertValidPath(r)
	return
}

func ToPath[T pathOrFileLike](input T) *Path {
	switch v := any(input).(type) {
	case *Path:
		return v
	case Path:
		return &v
	case *IPath:
		return v.Mutate()
	case IPath:
		return v.Mutate()
	case *File:
		return New(v.Name)
	case File:
		return New(v.Name)
	case string:
		return New(v)
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
}

func ToImmutable[T pathOrFileLike](input T) *IPath {
	switch v := any(input).(type) {
	case *Path:
		return v.Immutable()
	case Path:
		return v.Immutable()
	case *IPath:
		return v
	case IPath:
		return &v
	case *File:
		return INew(v.Name)
	case File:
		return INew(v.Name)
	case string:
		return INew(v)
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
}

func toStringAssertAbs[T pathOrFileLike](input T) string {
	t := ToString(input)
	if !filepath.IsAbs(t) {
		panic(PathErr.New("传入的参数必须是绝对路径"))
	}
	t = filepath.Clean(t)
	return t
}

func toStringMustAbs[T pathOrFileLike](input T, base string) string {
	t := ToString(input)
	if !filepath.IsAbs(t) {
		t = filepath.Join(base, t)
	} else {
		t = filepath.Clean(t)
	}
	return t
}

func toStringAbsCwd2[T1 pathOrFileLike, T2 pathOrFileLike](a T1, b T2) (*Path, *Path) {
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
