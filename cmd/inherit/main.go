package main

import (
	"fmt"
	_ "unsafe"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/signals"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

const MAGIC_STRING = "b3d0022a-26a7-4d6c-8799-0f0a6ba1bcd4"

func main() {
	defer signals.AppQuit.Finalize()

	env := myenv.Must1(codegen.CreateEnvironment("github.com/gongt/go/cmd/inherit", MAGIC_STRING))

	output_file := env.InputPath().SetFilename(fmt.Sprintf("%s_inherit_generate.go", env.InputPath().Base().Stem()))

	fmt.Println("Output file:", output_file)

	buff := sourcecode.NewGoFileBuffer()
	writer := myenv.Must1(env.NewOutput(output_file, buff))

	myenv.Must(writer.LearnPackageName())

	myenv.NotImplemented()

	myenv.Must(writer.WriteFile())

	signals.AppQuit.Set(0)
}
