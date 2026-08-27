package tools_bin

import (
	"maps"
	"slices"

	"github.com/gongt/go/internal/myenv"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

func SerializeGeneratorMain(env codegen.GeneratorEnvironment) {
	myenv.Must(env.NoMoreArgs())

	rdr := sourcecode.NewPackageReader()

	infos := myenv.Must1(rdr.ReadFilesType(env.WorkspaceRoot().Raw()))

	for info := range maps.Values(infos) {
		for decl := range slices.Values(info.Decls) {
			println(decl) // TODO
		}
	}
}
