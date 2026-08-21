package sourcecode

import (
	"go/ast"
	"go/parser"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/hashicorp/go-set/v3"
)

func SimulateGolangBuild(f string) bool {
	// 试图模拟go build判断是否要编译的逻辑
	// 实际只是跳过test和隐藏文件
	if strings.HasSuffix(f, "_test.go") {
		return true
	}
	if strings.HasPrefix(f, ".") || strings.HasPrefix(f, "_") {
		return true
	}
	return false
}

// 检测指定目录的包名
func DetectPackageName(folderPath string) (string, error) {
	entries, err := DetectPackageNames(folderPath, SimulateGolangBuild)
	if err != nil {
		return "", err
	}
	if len(entries) == 1 {
		return entries[0], nil
	}
	return "", errors.NewAnonymous("目录中包名不一致").WithDetails("path", folderPath, "packageNames", entries)
}

// 遍历指定目录下的所有Go文件（不递归），收集它们的包名，返回列表。
func DetectPackageNames(folderPath string, ignore func(string) bool) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, errors.Extend(err, "无法读取文件夹").WithDetails("path", folderPath)
	}
	var goFiles []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) != ".go" {
			continue
		}
		if ignore != nil && ignore(entry.Name()) {
			continue
		}
		goFiles = append(goFiles, entry)
	}

	if len(goFiles) == 0 {
		return []string{normalize(path.Base(folderPath))}, nil
	}

	var results = set.New[string](2)

	for _, file := range goFiles {
		filePath := path.Join(folderPath, file.Name())

		var ast *ast.File
		ast, err = ReadSourceFile(filePath, parser.PackageClauseOnly)
		if err != nil {
			return nil, err
		}

		name := ast.Name.Name

		results.Insert(name)
	}

	return results.Slice(), nil
}

func normalize(name string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9]+`)
	return strings.Join(re.FindAllString(name, -1), "")
}
