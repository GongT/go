package fpath

type PathLike interface {
	*Path | *IPath | string
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
		panic("fpath.ToString: 传入的类型必须是 *Path 或 *IPath 或 string")
	}
}
