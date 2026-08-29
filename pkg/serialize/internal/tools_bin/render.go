package tools_bin

import (
	"go/ast"
	"path"
	"strconv"
	"strings"

	sourcecode "github.com/gongt/go/pkg/source_code"
)

// renderNode 原样输出ast节点对应的源码片段
func renderNode(file sourcecode.FileInfo, node ast.Node) string {
	var sb strings.Builder
	file.CloneAt(&sb, node.Pos(), node.End())
	return sb.String()
}

// typeText 渲染一个类型表达式为可写入生成文件的源码文本，自动将外部包引用替换为生成文件中的导入别名
func typeText(buf sourcecode.GoFileBuffer, file sourcecode.FileInfo, expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		return "*" + typeText(buf, file, star.X)
	}

	rendered := renderNode(file, expr)

	var qualifier *ast.Ident
	ast.Inspect(expr, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		qualifier, _ = selector.X.(*ast.Ident)
		return false
	})
	if qualifier == nil || !strings.HasPrefix(rendered, qualifier.Name) {
		return rendered
	}

	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		localName := path.Base(importPath)
		if spec.Name != nil {
			localName = spec.Name.Name
		}
		if localName == qualifier.Name {
			return buf.AddImport(importPath) + rendered[len(qualifier.Name):]
		}
	}
	return rendered
}
