package stl

import (
	"go/build"
	"os/exec"
	"strings"
)

var stdlibs []string

// IsStandardImport 检查给定的包路径是否属于标准库
func IsStandardImport(path string) bool {
	path = strings.Trim(path, "\"")

	// Special handling for the "C" pseudo-package
	if path == "C" {
		return true
	}

	// Internal packages, compiler tooling, and standard packages do not contain a "." in their first path segment.
	// This covers paths like "fmt", "net/http", or "crypto/rand".
	firstSegment := path
	if idx := strings.Index(path, "/"); idx != -1 {
		firstSegment = path[:idx]
	}

	if !strings.Contains(firstSegment, ".") {
		return true
	}

	// Fallback/Robust Check: Try to find the package using the go/build context.
	// Standard library packages are located under GOROOT, not GOPATH or modules.
	pkg, err := build.Import(path, "", build.FindOnly)
	if err == nil && pkg.Goroot {
		return true
	}

	return false
}

// ListStandardLibrary 运行go list std并返回标准库包列表
func ListStandardLibrary() []string {
	if stdlibs == nil {
		out, err := exec.Command("go", "list", "std").Output()
		if err != nil {
			return nil
		}
		for line := range strings.Lines(string(out)) {
			if line != "" {
				stdlibs = append(stdlibs, line)
			}
		}
	}
	return stdlibs
}
