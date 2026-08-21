package fpath

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys/fpath/internal"
)

const maxSymlinkDepth = 50

var ErrLoop = errors.NewTemplate("Too many symbolic links were encountered")

func resolveRealpath(name string, allowMissingFinal bool, allowBrokenLink bool, depth int) (string, error) {
	if depth > maxSymlinkDepth {
		return "", ErrLoop.New()
	}

	absolute := name
	volume := filepath.VolumeName(absolute)
	root := volume + string(filepath.Separator)
	parts := splitPath(absolute[len(volume):])
	resolved := root

	for i := range parts {
		part := parts[i]
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			resolved = filepath.Dir(resolved)
			continue
		}
		candidate := filepath.Join(resolved, part)
		info, err := os.Lstat(candidate)
		if err != nil {
			if os.IsNotExist(err) && allowMissingFinal {
				if i+1 < len(parts) {
					candidate += string(filepath.Separator) + strings.Join(parts[i+1:], string(filepath.Separator))
				}
				return internal.Clean(candidate), nil
			}
			return "", PathErr.Wrap(err)
		}

		if info.Mode()&os.ModeSymlink == 0 {
			resolved = candidate
			continue
		}

		target, err := os.Readlink(candidate)
		if err != nil {
			return "", PathErr.Wrap(err)
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(candidate), target)
		}
		remaining := target
		if len(parts) > i+1 {
			remaining += string(filepath.Separator) + strings.Join(parts[i+1:], string(filepath.Separator))
		}
		allowTargetMissing := allowMissingFinal || (allowBrokenLink && i == len(parts)-1)
		return resolveRealpath(remaining, allowTargetMissing, allowBrokenLink, depth+1)
	}

	return internal.Clean(resolved), nil
}

func splitPath(name string) []string {
	return strings.Split(name, string(filepath.Separator))
}

// 跟踪符号链接转换成真实路径，解析后路径一定存在，但可能是损坏的符号链接
func (p *Path) Realpath() error {
	if err := p.MaybeConvertAbsolute(); err != nil {
		return err
	}
	value, err := resolveRealpath(p.value, false, true, 0)
	if err != nil {
		return PathErr.Wrap(err).WithDetails("path", p.value, "mode", "normal")
	}
	p.value = value
	p.canonicalizeCache = value
	return nil
}

func (p *Path) MustRealpath() {
	if err := p.Realpath(); err != nil {
		panic(err)
	}
}

// 跟踪符号链接转换成真实路径，并且返回的一定是非符号链接文件（目录），任意路径不存在，都返回错误
//
// 理应和[filepath.EvalSymlinks]行为一致
func (p *Path) RealpathExisting() error {
	if err := p.MaybeConvertAbsolute(); err != nil {
		return err
	}
	value, err := resolveRealpath(p.value, false, false, 0)
	if err != nil {
		return PathErr.Wrap(err).WithDetails("path", p.value, "mode", "enforce")
	}
	if _, err = os.Stat(value); err != nil {
		return PathErr.Wrap(err).WithDetails("path", p.value, "mode", "enforce")
	}
	p.value = value
	p.canonicalizeCache = value
	return nil
}

func (p *Path) MustRealpathExisting() {
	if err := p.RealpathExisting(); err != nil {
		panic(err)
	}
}

// 跟踪符号链接转换成真实路径，解析后路径可以不存在
func (p *Path) RealpathMissing() error {
	if err := p.MaybeConvertAbsolute(); err != nil {
		return err
	}
	value, err := resolveRealpath(p.value, true, true, 0)
	if err != nil {
		return PathErr.Wrap(err).WithDetails("path", p.value, "mode", "relax")
	}
	p.value = value
	p.canonicalizeCache = value
	return nil
}

func (p *Path) MustRealpathMissing() {
	if err := p.RealpathMissing(); err != nil {
		panic(err)
	}
}
