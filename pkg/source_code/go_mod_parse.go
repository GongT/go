package sourcecode

import (
	"go/ast"
	"go/parser"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/filesys"
	"golang.org/x/mod/modfile"
)

func FindGoMod(fromPath string) (string, error) {
	modFile, _ := filesys.FindUpUntilEntry(fromPath, "go.mod", ".git")
	if path.Base(modFile) != "go.mod" {
		return "", errors.NewAnonymous("无法找到go.mod文件").WithDetails("path", fromPath)
	}

	return modFile, nil
}

func CalculateImportPath(filePath string) (string, error) {
	modFile, err := FindGoMod(filePath)
	if err != nil {
		return "", err
	}

	goModData, _ := os.ReadFile(modFile)
	modName := modfile.ModulePath(goModData)

	relativePath, _ := filepath.Rel(path.Dir(modFile), filePath)

	return path.Join(modName, relativePath), nil
}

func DetectPackageName(folderPath string) (string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return "", errors.Extend(err, "无法读取文件夹").WithDetails("path", folderPath)
	}
	var goFiles []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "exports.go" {
			continue
		}
		if path.Ext(entry.Name()) == ".go" {
			goFiles = append(goFiles, entry)
		}
	}

	if len(goFiles) == 0 {
		return normalize(path.Base(folderPath)), nil
	}

	var foundPackageName string

	for _, file := range goFiles {
		filePath := path.Join(folderPath, file.Name())

		var ast *ast.File
		ast, err = ReadSourceFile(filePath, parser.PackageClauseOnly)
		if err != nil {
			return "", err
		}

		name := ast.Name.Name
		name = strings.TrimSuffix(name, "_test")

		if foundPackageName == "" {
			foundPackageName = name
		} else if foundPackageName != name {
			return "", errors.NewAnonymous("目录中包名不一致").WithDetails("path", filePath, "foundPackageName", foundPackageName, "currentPackageName", name)
		}
	}

	return foundPackageName, nil
}

func normalize(name string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9]+`)
	return strings.Join(re.FindAllString(name, -1), "")
}
