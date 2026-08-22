package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gongt/go/cmd/exports/internal"
	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/signals"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

const MAGIC_STRING = "33084b74-3121-47cd-9093-db68da5893bc"

type Options struct {
	VerboseMode bool `long:"verbose" description:"Enable verbose mode for debugging"`
}

func main() {
	defer signals.AppQuit.Finalize()

	env := myenv.Must(codegen.CreateEnvironment("github.com/gongt/go/cmd/exports", MAGIC_STRING))

	if path.Base(env.InputFile) != "exports.go" {
		panic("This generator is only supported inside exports.go")
	}

	var opts Options
	myenv.MustNil(env.ParseArgs(&opts))

	internalDir := fpath.New(env.InputFile, "..", "internal").Immutable().MustRealpath()

	if stat, err := os.Stat(internalDir.Raw()); err != nil || !stat.IsDir() {
		panic(fmt.Sprintf("The provided folder %s does not contain an internal directory", internalDir))
	}

	if !opts.VerboseMode {
		log.SetOutput(io.Discard)
	}
	bs := myenv.Must(internal.ParseFiles(internalDir))
	if !opts.VerboseMode {
		log.SetOutput(os.Stderr)
	}

	pkgName := myenv.Must(sourcecode.DetectPackageName(filepath.Dir(env.InputFile)))
	bs.SetPackageName(pkgName)

	writer := myenv.Must(env.NewOutput(env.InputFile, bs))

	codegen.WriteGeneratorComment(bs.Heading(), env.GeneratorFullName, nil)

	myenv.MustNil(writer.WriteFile())

	signals.AppQuit.Set(0)
}
