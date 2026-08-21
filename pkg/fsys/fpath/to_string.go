package fpath

import "path/filepath"

type PathLike interface {
	*Path | *IPath | string
}
type PathLikeRo interface {
	*IPath | string
}

func ToString[T PathLike](input T) string {
	switch v := any(input).(type) {
	case *Path:
		return v.String()
	case *IPath:
		return v.String()
	case string:
		return v
	default:
		panic(PathErr.New("传入的类型必须是 *Path 或 *IPath 或 string"))
	}
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
