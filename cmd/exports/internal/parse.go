package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/iterator"
	sourcecode "github.com/gongt/go/pkg/source_code"
)

func ParseFiles(dirPath *fpath.IPath) (sourcecode.GoFileBuffer, error) {
	rdr := sourcecode.NewPackageReader()
	rdr.Recursive = true

	infoMap, err := rdr.ReadFilesType(dirPath.Raw())
	if err != nil {
		return nil, err
	}

	fsb := sourcecode.NewGoFileBuffer()

	for _, content := range iterator.SortedMap(infoMap) {
		ParseFile(fsb, content)
	}

	return fsb, nil
}

type ExportType int

const (
	NoExport ExportType = iota
	Export
	ForceExport
)

func ParseFile(sb sourcecode.GoFileBuffer, fileInfo *sourcecode.FileInfo) {
	if len(fileInfo.Decls) == 0 {
		log.Printf("No declarations found in file: %s", fileInfo.Filename)
		return
	}

	log.Printf("Processing file: %s", fileInfo.Filename)
	// log.Println(content.Package.Path())

	defaultsExport := checkDefault(fileInfo)
	log.Printf("  - defaults export? %v", defaultsExport)

	for _, decl := range fileInfo.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			if decl.Recv != nil {
				// Method - skip
				continue
			}
			exportType := isExported(decl.Doc, defaultsExport)
			if exportType == NoExport {
				log.Printf("  - private symbol %s", decl.Name.Name)
				continue
			}

			emitFunc(sb, exportType, decl, fileInfo)
		case *ast.GenDecl: // Type declarations are represented as GenDecl in the AST
			if decl.Tok == token.IMPORT {
				// import - skip
				continue
			}
			exportType := isExported(decl.Doc, defaultsExport)
			if exportType == NoExport {
				log.Printf("  - private symbol %v", decl.Specs)
				continue
			}

			doc := decl.Doc

			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					emitType(sb, exportType, spec, fileInfo, doc)
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						emitVar(sb, exportType, decl.Tok == token.VAR, name.Name, fileInfo, doc)
					}
				default:
					log.Printf("Other spec type: %T", spec)
				}
			}
		default:
			log.Printf("Other declaration type: %T", decl)
		}
	}
}

func emitVar(sb sourcecode.GoFileBuffer, exportType ExportType, isVar bool, symbol string, file *sourcecode.FileInfo, doc *ast.CommentGroup) {
	exported := exportedName(symbol, exportType)
	if exported == "" {
		return
	}
	log.Printf("  - export variable %s", symbol)

	copyDoc(sb, doc)

	if isVar {
		sb.WriteString("var ")
	} else {
		sb.WriteString("const ")
	}

	sb.WriteString(exported)
	sb.WriteString(" = ")
	sb.WriteString(symbol)
	sb.WriteByte('\n')
}

func emitFunc(sb sourcecode.GoFileBuffer, exportType ExportType, decl *ast.FuncDecl, file *sourcecode.FileInfo) {
	exported := exportedName(decl.Name.Name, exportType)
	if exported == "" {
		return
	}
	log.Printf("  - export function %s", decl.Name.Name)

	copyDoc(sb, decl.Doc)
	sb.WriteString("func ")
	sb.WriteString(exported)

	// 泛型参数部分
	if decl.Type.TypeParams != nil {
		params := make([]string, 0, len(decl.Type.TypeParams.List))
		for _, field := range decl.Type.TypeParams.List {
			constraint := renderType(sb, file, field.Type)
			if ident, ok := field.Type.(*ast.Ident); ok {
				if ident.Name != "any" && ident.Name != "comparable" {
					constraint = sb.AddImport(file.Container().Path()) + "." + ident.Name
				}
			}
			for _, name := range field.Names {
				params = append(params, name.Name+" "+constraint)
			}
		}
		fmt.Fprintf(sb, "[%s]", strings.Join(params, ", "))
	}

	// 参数部分
	recallArgList := []string{}
	inputs := []string{}
	if decl.Type.Params != nil {
		num := 0
		for _, field := range decl.Type.Params.List {
			names := field.Names
			if len(names) == 0 {
				names = []*ast.Ident{{Name: "_"}}
			}
			for range names {
				num++
				variableName := fmt.Sprintf("arg%d", num)
				inputs = append(inputs, variableName+" "+renderType(sb, file, field.Type))
				callArgument := variableName
				if _, ok := field.Type.(*ast.Ellipsis); ok {
					callArgument += "..."
				}
				recallArgList = append(recallArgList, callArgument)
			}
		}
	}
	fmt.Fprintf(sb, "(%s)", strings.Join(inputs, ", "))

	// 返回部分
	hasReturn := false
	if decl.Type.Results != nil && len(decl.Type.Results.List) > 0 {
		hasReturn = true
		outputs := []string{}

		for _, field := range decl.Type.Results.List {
			count := len(field.Names)
			if count == 0 {
				count = 1
			}
			for range count {
				outputs = append(outputs, renderType(sb, file, field.Type))
			}
		}

		fmt.Fprintf(sb, "(%s)", strings.Join(outputs, ", "))
	}

	// 函数体
	sb.WriteString(" {\n\t")
	if hasReturn {
		sb.WriteString("return ")
	}
	call := sb.AddImport(file.Container().Path()) + "." + decl.Name.Name
	if decl.Type.TypeParams != nil && decl.Type.Params != nil && len(decl.Type.Params.List) == 0 {
		arguments := make([]string, 0)
		for _, field := range decl.Type.TypeParams.List {
			for _, name := range field.Names {
				arguments = append(arguments, name.Name)
			}
		}
		call += "[" + strings.Join(arguments, ", ") + "]"
	}
	fmt.Fprintf(sb, "%s(%s)", call, strings.Join(recallArgList, ", "))
	sb.WriteString("\n}\n\n")
}

func renderNode(file *sourcecode.FileInfo, node ast.Node) string {
	var rendered strings.Builder
	file.CloneAt(&rendered, node.Pos(), node.End())
	return rendered.String()
}

func renderType(sb sourcecode.GoFileBuffer, file *sourcecode.FileInfo, expr ast.Expr) string {
	if pointer, ok := expr.(*ast.StarExpr); ok {
		return "*" + renderType(sb, file, pointer.X)
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
			return sb.AddImport(importPath) + rendered[len(qualifier.Name):]
		}
	}
	return rendered
}

func emitType(sb sourcecode.GoFileBuffer, exportType ExportType, spec *ast.TypeSpec, file *sourcecode.FileInfo, doc *ast.CommentGroup) {
	exported := exportedName(spec.Name.Name, exportType)
	if exported == "" {
		return
	}

	log.Printf("  - export type %s", spec.Name.Name)

	id := sb.AddImport(file.Container().Path())

	copyDoc(sb, doc)
	sb.WriteString("type ")
	sb.WriteString(exported)
	sb.WriteString(" = ")
	sb.WriteString(id)
	sb.WriteByte('.')
	sb.WriteString(spec.Name.Name)

	sb.WriteByte('\n')
}

func exportedName(symbol string, exportType ExportType) string {
	if exportType == NoExport {
		panic("调用路径错误")
	}
	if symbol[0] == '_' {
		if exportType == ForceExport {
			symbol = symbol[1:]
		} else {
			return ""
		}
	}

	if symbol == "" {
		panic("空符号被标为导出")
	}

	r, _ := utf8.DecodeRuneInString(symbol)
	if unicode.IsUpper(r) {
		return symbol
	}

	u := unicode.ToUpper(r)
	if u == r {
		panic(fmt.Sprintf("无大小写的符号被标为导出: %q", symbol))
	}

	return string(unicode.ToUpper(r)) + symbol[utf8.RuneLen(r):]
}

var exportedRegex = regexp.MustCompile(`^//\s*@(public|exported)`)
var privateRegex = regexp.MustCompile(`^//\s*@(private|internal|unexported)`)

func checkDefault(file *sourcecode.FileInfo) ExportType {
	for _, group := range file.Comments {
		if group.Pos() > file.Package {
			break
		}

		for _, comment := range group.List {
			if exportedRegex.MatchString(comment.Text) {
				return Export
			}
		}
	}
	return NoExport
}

func isExported(doc *ast.CommentGroup, defaultsExport ExportType) ExportType {
	if doc == nil {
		return defaultsExport
	}

	for _, comment := range doc.List {
		if exportedRegex.MatchString(comment.Text) {
			return ForceExport
		}
		if privateRegex.MatchString(comment.Text) {
			return NoExport
		}
	}
	return defaultsExport
}

func copyDoc(sb sourcecode.GoFileBuffer, doc *ast.CommentGroup) {
	if doc == nil {
		return
	}
	for _, comment := range doc.List {
		if exportedRegex.MatchString(comment.Text) || privateRegex.MatchString(comment.Text) {
			continue
		}
		sb.WriteString(comment.Text)
		sb.WriteByte('\n')
	}
}
