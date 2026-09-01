package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/iterator"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

func ParseFiles(env codegen.GeneratorEnvironment, dirPath *fpath.IPath) (sourcecode.GoFileBuffer, error) {
	rdr := sourcecode.NewPackageReader()
	rdr.Recursive = true

	infoMap, err := rdr.ReadFilesType(dirPath.Raw())
	if err != nil {
		return nil, err
	}

	fsb := sourcecode.NewGoFileBuffer()
	fsb.NamePackage = sourcecode.OriginalName

	for _, content := range iterator.SortedMap(infoMap) {
		if len(content.Decls) == 0 {
			log.Printf("No declarations found in file: %s", content.Filename)
			continue
		}
		fmt.Fprintf(fsb, "// - %s\n\n", fpath.MustRelative(content.Filename, env.WorkspaceRoot()))
		cnt := ParseFile(fsb, content)

		if cnt == 0 {
			fsb.WriteString("// No exported symbols\n")
		}

		fsb.WriteString("\n\n")
	}

	return fsb, nil
}

type ExportType int

const (
	NoExport ExportType = iota
	Export
	ForceExport
)

func ParseFile(sb sourcecode.GoFileBuffer, fileInfo sourcecode.FileInfo) (exportedCount int) {
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

			exportedCount += emitFunc(sb, exportType, decl, fileInfo)
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
					exportedCount += emitType(sb, exportType, spec, fileInfo, doc)
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						exportedCount += emitVar(sb, exportType, decl.Tok == token.VAR, name.Name, fileInfo, doc)
					}
				default:
					log.Printf("Other spec type: %T", spec)
				}
			}
		default:
			log.Printf("Other declaration type: %T", decl)
		}
	}
	return
}

func emitVar(sb sourcecode.GoFileBuffer, exportType ExportType, isVar bool, symbol string, file sourcecode.FileInfo, doc *ast.CommentGroup) int {
	exported := exportedName(symbol, exportType)
	if exported == "" {
		return 0
	}

	log.Printf("  - export variable %s", symbol)
	if sb.AddExport(exported, nil) {
		log.Printf("    ! duplicate, skip")
		return 0
	}

	copyDoc(sb, doc)

	if isVar {
		sb.WriteString("var ")
	} else {
		sb.WriteString("const ")
	}

	sb.WriteString(exported)
	sb.WriteString(" = ")
	sb.WriteString(sb.AddImport(file.Container().Path()))
	sb.WriteString(".")
	sb.WriteString(symbol)
	sb.WriteByte('\n')
	return 1
}

func typeParams(sb sourcecode.GoFileBuffer, tp *ast.FieldList, file sourcecode.FileInfo) string {
	if tp == nil {
		return ""
	}

	params := make([]string, 0, len(tp.List))
	for _, field := range tp.List {
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
	printList(sb, params, '[')

	arguments := make([]string, 0)
	for _, field := range tp.List {
		for _, name := range field.Names {
			arguments = append(arguments, name.Name)
		}
	}
	return "[" + strings.Join(arguments, ", ") + "]"
}

func emitFunc(sb sourcecode.GoFileBuffer, exportType ExportType, decl *ast.FuncDecl, file sourcecode.FileInfo) int {
	exported := exportedName(decl.Name.Name, exportType)
	if exported == "" {
		return 0
	}
	log.Printf("  - export function %s", decl.Name.Name)
	if sb.AddExport(exported, nil) {
		log.Printf("    ! duplicate, skip")
		return 0
	}

	copyDoc(sb, decl.Doc)
	sb.WriteString("func ")
	sb.WriteString(exported)

	// 泛型参数部分
	genericCall := typeParams(sb, decl.Type.TypeParams, file)

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
	printList(sb, inputs, '(')

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

		printList(sb, outputs, '(')
	}

	// 函数体
	sb.WriteString(" {\n\t")
	if hasReturn {
		sb.WriteString("return ")
	}
	pkgIdent := sb.AddImport(file.Container().Path())
	fmt.Fprintf(sb, "%s.%s", pkgIdent, decl.Name.Name)

	if decl.Type.TypeParams != nil && decl.Type.Params != nil && len(decl.Type.Params.List) == 0 {
		// 如果函数有泛型参数，但没有输入参数，则需要在调用时传递泛型参数，否则编译器无法推断泛型类型
		sb.WriteString(genericCall)
	}

	printList(sb, recallArgList, '(')
	sb.WriteString("\n}\n\n")
	return 1
}

func printList(sb sourcecode.GoFileBuffer, list []string, quote byte) {
	if quote != 0 {
		sb.WriteByte(quote)
	}
	for i, item := range list {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(item)
	}
	if quote != 0 {
		switch quote {
		case '(':
			sb.WriteByte(')')
		case '[':
			sb.WriteByte(']')
		case '{':
			sb.WriteByte('}')
		default:
			panic(fmt.Sprintf("Unsupported quote character: %c", quote))
		}
	}
}

func renderType(sb sourcecode.GoFileBuffer, file sourcecode.FileInfo, expr ast.Expr) string {
	typ := file.Container().TypesInfo().TypeOf(expr)
	tr := sourcecode.NewTypeResolver(sb)

	return tr.QualifyType(typ)
}

func emitType(sb sourcecode.GoFileBuffer, exportType ExportType, spec *ast.TypeSpec, file sourcecode.FileInfo, doc *ast.CommentGroup) int {
	exported := exportedName(spec.Name.Name, exportType)
	if exported == "" {
		return 0
	}

	log.Printf("  - export type %s", spec.Name.Name)
	if sb.AddExport(exported, nil) {
		log.Printf("    ! duplicate, skip")
		return 0
	}

	id := sb.AddImport(file.Container().Path())

	copyDoc(sb, doc)
	sb.WriteString("type ")
	sb.WriteString(exported)

	genericCall := typeParams(sb, spec.TypeParams, file)

	sb.WriteString(" = ")
	sb.WriteString(id)
	sb.WriteByte('.')
	sb.WriteString(spec.Name.Name)
	sb.WriteString(genericCall)

	sb.WriteByte('\n')
	return 1
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

	return ""
	// u := unicode.ToUpper(r)
	//
	//	if u == r {
	//		panic(fmt.Sprintf("无大小写的符号被标为导出: %q", symbol))
	//	}
	//
	// return string(unicode.ToUpper(r)) + symbol[utf8.RuneLen(r):]
}

var exportedRegex = regexp.MustCompile(`^//\s*@(public|exported)`)
var privateRegex = regexp.MustCompile(`^//\s*@(private|internal|unexported)`)

func checkDefault(file sourcecode.FileInfo) ExportType {
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
