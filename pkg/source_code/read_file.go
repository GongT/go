package sourcecode

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/gongt/go/pkg/errors"
)

type Mode = parser.Mode

func ReadSourceFile(filePath string, mode parser.Mode) (*ast.File, error) {
	if mode == 0 {
		mode = parser.SkipObjectResolution + parser.ParseComments
	}
	fSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fSet, filePath, nil, mode)
	if err != nil {
		return nil, errors.Extend(err, "文件解析失败").WithDetails("filePath", filePath)
	}
	return astFile, nil
}
