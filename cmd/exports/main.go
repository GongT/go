package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

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

var allowFileName = regexp.MustCompile(`^exports(?:_(.+))?\.go$`)

func main() {
	defer signals.AppQuit.Finalize()

	env := myenv.Must(codegen.CreateEnvironment("github.com/gongt/go/cmd/exports", MAGIC_STRING))

	if !allowFileName.MatchString(path.Base(env.InputFile)) {
		panic(errors.New("此生成器仅支持在 exports_xxx.go 中使用"))
	}

	var opts Options
	myenv.MustNil(env.ParseArgs(&opts))

	internalDir, err := fpath.New(env.InputFile).SetFilename("internal").Immutable().RealpathExisting()

	if err != nil {
		panic(fmt.Errorf("无法找到 internal 目录: %w", err))
	}

	if stat, err := os.Stat(internalDir.Raw()); err != nil || !stat.IsDir() {
		panic(fmt.Errorf("%s必须是目录", internalDir))
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
