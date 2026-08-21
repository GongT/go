package sourcecode

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys"
	"golang.org/x/mod/modfile"
)

type GoModFile struct {
	Path       string
	content    []byte
	moduleName string
}

func FindGoMod(fromPath string) (*GoModFile, error) {
	found, _ := fsys.FindUpUntilEntry(fromPath, "go.mod", ".git")
	if path.Base(found) != "go.mod" {
		return nil, errors.NewAnonymous("无法找到go.mod文件").WithDetails("path", fromPath)
	}

	return OpenGoMod(found), nil
}

func OpenGoMod(filePath string) *GoModFile {
	return &GoModFile{
		Path: filePath,
	}
}

func (g *GoModFile) Dir() string {
	return path.Dir(g.Path)
}
func (g *GoModFile) LoadFile() error {
	if g.content != nil {
		return nil
	}

	content, err := os.ReadFile(g.Path)
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
	relativePath, err := filepath.Rel(path.Dir(g.Path), filePath)
	if err != nil {
		return "", errors.Extend(err, "无法计算相对路径").WithDetails("base", path.Dir(g.Path), "target", filePath)
	}

	relativePath = filepath.Clean(relativePath)

	if moduleName, err := g.GetModuleName(); err == nil {
		if relativePath == "." {
			return moduleName, nil
		}
		if strings.HasPrefix(relativePath, "..") {
			return "", errors.NewAnonymous("文件路径不在go.mod模块内").WithDetails("filePath", filePath, "goModPath", g.Path)
		}
		return moduleName + "/" + relativePath, nil
	} else {
		return "", err
	}
}
