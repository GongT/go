package sourcecode

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys"
	"github.com/gongt/go/pkg/fsys/fpath"
	"golang.org/x/mod/modfile"
)

type GoModFile struct {
	Path       *fpath.IPath
	content    []byte
	moduleName string
}

func FindGoMod[T fpath.PathLike](fromPath T) (*GoModFile, error) {
	found, _ := fsys.FindUpUntilEntry(fromPath, "go.mod", ".git")
	if found.Base().Name != "go.mod" {
		return nil, errors.NewAnonymous("无法找到go.mod文件").WithDetails("path", fromPath)
	}

	return OpenGoMod(found), nil
}

func OpenGoMod[T fpath.PathLike](filePath T) *GoModFile {
	return &GoModFile{
		Path: fpath.ToImmutable(filePath),
	}
}

func (g *GoModFile) Dir() *fpath.IPath {
	return g.Path.Dir()
}
func (g *GoModFile) LoadFile() error {
	if g.content != nil {
		return nil
	}

	content, err := os.ReadFile(g.Path.Raw())
	if err != nil {
		return errors.Extend(err, "无法读取go.mod文件").WithDetails("path", g.Path)
	}

	g.content = content
	return nil
}

func (g *GoModFile) GetModuleName() (string, error) {
	if g.moduleName == "" {
		if err := g.LoadFile(); err != nil {
			return "", err
		}
		modName := modfile.ModulePath(g.content)
		if modName == "" {
			return "", errors.NewAnonymous("go.mod文件中缺少模块名").WithDetails("path", g.Path)
		}
		g.moduleName = modName
	}
	return g.moduleName, nil
}

// 给定一个文件路径，计算它的绝对导入路径（import path）
func (g *GoModFile) CalculateImportPath(filePath string) (string, error) {
	relativePath, err := fpath.ToRelative(filePath, g.Dir())
	if err != nil {
		return "", errors.Extend(err, "无法计算相对路径").WithDetails("base", path.Dir(g.Path.Raw()), "target", filePath)
	}

	var moduleName string
	if moduleName, err = g.GetModuleName(); err != nil {
		return "", err
	}

	if relativePath == "." {
		return moduleName, nil
	}
	if strings.HasPrefix(relativePath, "..") {
		return "", errors.NewAnonymous("文件路径不在go.mod模块内").WithDetails("filePath", filePath, "goModPath", g.Path)
	}

	// 去掉文件名部分，删掉relativePath开头的"./"
	relativePath = filepath.Dir(relativePath)

	return moduleName + "/" + relativePath, nil
}
