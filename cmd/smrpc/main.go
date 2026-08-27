package main

import (
	_ "unsafe"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/signals"
	_ "github.com/gongt/go/pkg/smrpc/generate"
	"github.com/gongt/go/pkg/source_code/codegen"
)

// Link our local declaration to the target private function
// Format: //go:linkname [local_alias] [module_path]/[package].[private_func_name]
//
//go:linkname runMain github.com/gongt/go/pkg/smrpc/generate/internal/bin_rpc.RunMain
func runMain(codegen.GeneratorEnvironment)

const MAGIC_STRING = "caa81d42-87f1-4462-9e42-3ef9590748cf"

func main() {
	defer signals.AppQuit.Finalize()

	env := myenv.Must1(codegen.CreateEnvironment("github.com/gongt/go/cmd/smrpc", MAGIC_STRING))

	runMain(env)

	signals.AppQuit.Set(0)
}
