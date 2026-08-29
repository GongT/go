package tools_bin

import (
	"go/ast"
	"go/token"
	"log"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

func SerializeGeneratorMain(env codegen.GeneratorEnvironment) {
	myenv.Must(env.NoMoreArgs())

	rdr := sourcecode.NewPackageReader()
	infos := myenv.Must1(rdr.ReadFilesType(env.WorkspaceRoot().Raw()))

	target := env.InputPath().Raw()
	file, ok := infos[target]
	if !ok {
		panic(errors.NewAnonymous("未在类型分析结果中找到输入文件，路径=%s", target))
	}

	structs := collectStructs(file)
	if len(structs) == 0 {
		log.Printf("文件%s中没有可生成序列化代码的结构体，跳过", target)
		return
	}

	buf := sourcecode.NewGoFileBuffer()
	buf.SetPackageName(file.Container().Name())

	g := Ctx(file, buf)

	for _, st := range structs {
		g.emitHelpers(st)
		g.emitMarshalFunc(st)
		g.emitUnmarshalFunc(st)
	}

	stem, ext := env.InputPath().Base().SplitExt()
	outPath := env.InputPath().SetBase(stem + "_serialize_generate" + ext)

	writer := myenv.Must1(env.NewOutput(outPath, buf))

	myenv.Must(buf.CheckSyntax())

	myenv.Must(writer.WriteFile())
}

// collectStructs 仅收集输入文件顶层声明的结构体类型，保留字段声明顺序
func collectStructs(file sourcecode.FileInfo) []*structInfo {
	var out []*structInfo
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			out = append(out, &structInfo{
				name:   typeSpec.Name.Name,
				fields: collectStructFields(structType.Fields.List),
			})
		}
	}
	return out
}
