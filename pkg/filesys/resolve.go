package filesys

import (
	"os"
	"path"
	"slices"

	"github.com/gongt/go/pkg/errors"
)

func ResolvePath(segments ...string) (string, error) {
	rooted := false
	for i, seg := range slices.Backward(segments) {
		if seg == "" {
			return "", errors.NewAnonymous("Resolve的第%d个参数为空", i)
		}
		if path.IsAbs(seg) {
			rooted = true
			segments = segments[i:]
			break
		}
	}
	if !rooted {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		segments[0] = wd + "/" + segments[0]
	}

	return path.Join(segments...), nil
}

func MustResolvePath(segments ...string) string {
	ret, err := ResolvePath(segments...)
	if err != nil {
		panic(err)
	}
	return ret
}
