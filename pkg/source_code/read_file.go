package sourcecode

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/goforj/godump"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/source_code/internal/reader"
	"golang.org/x/tools/go/packages"
)

type Mode = parser.Mode
type FileInfo = reader.FileInfo
type PackageInfo = reader.PackageInfo

type PackageReader struct {
	Mode              parser.Mode
	Recursive         bool
	Testing           bool
	BuildFlags        []string
	IgnoreSyntaxError bool // 是否忽略语法错误，默认false
}

func NewPackageReader() *PackageReader {
	return &PackageReader{
		Mode:       parser.SkipObjectResolution + parser.ParseComments,
		Recursive:  true,
		Testing:    false,
		BuildFlags: nil,
	}
}

// 获取一个目录的模块名，如果没有源码，则返回（处理过的）文件夹名
func (p *PackageReader) GetModuleName(dir string) (string, error) {
	mode := packages.NeedName
	pkgs, _, err := p.callLoad(dir, mode)
	if err != nil {
		return "", err
	}
	for _, pkg := range pkgs {
		return pkg.Name, nil
	}
	return normalize(filepath.Base(dir)), nil
}

// GetPackageName 从文件中读出 “package xxx” 语句
//
// 如果文件不存在，返回错误。
// 如果有语法错误等问题导致无法解析，返回空字符串而非错误。
func (p *PackageReader) GetPackageName(filePath string) (string, error) {
	fSet := token.NewFileSet()
	buff, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Extend(err, "无法读取文件").WithDetails("path", filePath)
	}
	astFile, _ := parser.ParseFile(fSet, filePath, buff, parser.PackageClauseOnly)
	if astFile == nil || astFile.Name == nil {
		return "", nil
	}
	return astFile.Name.Name, nil
}

// ReadFile 读取一个文件，这不需要packages包，也不做test等判断，除了mode以外的参数无效
func (p *PackageReader) ParseFileAst(filePath string) (*ast.File, error) {
	fSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fSet, filePath, nil, p.Mode)
	if err != nil {
		return nil, errors.Extend(err, "文件解析失败").WithDetails("path", filePath)
	}
	return astFile, nil
}

// ReadFiles 读取一个目录下的所有go文件，返回文件信息，只有ast
func (p *PackageReader) ReadFilesAst(dir string) (map[string]FileInfo, error) {
	mode := packages.LoadFiles | packages.NeedSyntax
	_, files, err := p.callLoad(dir, mode)
	return files, err
}

// ReadFiles 读取一个目录下的所有go文件，返回文件信息，包括类型
func (p *PackageReader) ReadFilesType(dir string) (map[string]FileInfo, error) {
	_, files, err := p.callLoad(dir, packages.LoadSyntax)
	return files, err
}

func (p *PackageReader) callLoad(dir string, mode packages.LoadMode) ([]*packages.Package, map[string]FileInfo, error) {
	if p.Testing {
		mode |= packages.NeedForTest
	}
	var mu sync.Mutex

	contentMap := make(map[string][]byte)
	fset := token.NewFileSet()
	cfg := &packages.Config{
		Mode:       mode,
		Dir:        dir,
		BuildFlags: p.BuildFlags,
		Fset:       fset,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			mu.Lock()
			defer mu.Unlock()

			contentMap[filename] = src
			ast, err := parser.ParseFile(fset, filename, src, p.Mode)
			if err != nil && p.IgnoreSyntaxError {
				err = nil
			}

			// 根据 golang.org/x/tools/go/packages/packages.go 中的parseFiles函数，ast可以为nil
			return ast, err
		},
	}

	filter := "."
	if p.Recursive {
		filter = "./..."
	}

	pkgs, err := packages.Load(cfg, filter)
	if err != nil {
		return nil, nil, errors.Extend(err, "包解析失败").WithDetails("dir", dir)
	}

	result := make(map[string]FileInfo)
	for pkg := range packages.Postorder(pkgs) {
		if len(pkg.Errors) > 0 {
			log.Printf("包%s解析失败:\n", pkg.ID)
			for index, suberr := range pkg.Errors {
				log.Printf("  - %d: %v", index, suberr)
			}
		}
		if pkg.Module != nil {
			if pkg.Module.Error != nil {
				return pkgs, nil, errors.NewAnonymous("包%s所属模块解析失败", pkg.ID).WithDetails("error", pkg.Module.Error.Err)
			}
		}

		pkgRead := reader.NewPackageInfo(pkg)

		for _, ast := range pkg.Syntax {
			path := fset.Position(ast.FileStart).Filename
			if path == "" {
				log.Println("文件没有对应的文件名")
				godump.Dump(ast)
				continue
			}

			content, ok := contentMap[path]
			if !ok {
				log.Printf("文件%s没有对应的内容\n", path)
			}

			result[path] = reader.NewFileInfo(ast, fset, path, content, pkgRead)
		}
	}

	return pkgs, result, nil
}
