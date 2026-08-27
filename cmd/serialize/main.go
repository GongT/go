package main

import (
	_ "unsafe"

	"github.com/gongt/go/internal/myenv"
	_ "github.com/gongt/go/pkg/serialize"
	"github.com/gongt/go/pkg/signals"
	"github.com/gongt/go/pkg/source_code/codegen"
)

// Link our local declaration to the target private function
// Format: //go:linkname [local_alias] [module_path]/[package].[private_func_name]
//
//go:linkname runMain github.com/gongt/go/pkg/serialize/internal/tools_bin/main.SerializeGeneratorMain
func runMain(codegen.GeneratorEnvironment)

const MAGIC_STRING = "4a18ae49-5f26-43e3-b409-7b306174dd8c"

func main() {
	defer signals.AppQuit.Finalize()

	env := myenv.Must1(codegen.CreateEnvironment("github.com/gongt/go/cmd/serialize", MAGIC_STRING))

	runMain(env)

	signals.AppQuit.Set(0)
}
