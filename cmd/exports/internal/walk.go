package internal

import (
	"iter"

	sourcecode "github.com/gongt/go/pkg/reflection/source_code"
)

type WalkResult struct {
	AbsolutePath string

	PackageName string
	FQDN        string
}

func WalkFilesIn(dir string) iter.Seq[WalkResult] {
	root, err := sourcecode.DetectPackageName(dir)

	if err != nil {
		return iter.Empty[WalkResult]().WithError(err)
	}

	return walkDir(dir, root)
}

func walkDir(dir string) iter.Seq[WalkResult] {

}
