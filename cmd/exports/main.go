package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/gongt/go/cmd/exports/internal"
	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/signals"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

const MAGIC_STRING = "33084b74-3121-47cd-9093-db68da5893bc"

type Options struct {
	VerboseMode bool `long:"verbose" short:"v" description:"Enable verbose mode for debugging" env:"VERBOSE"`
}

var allowFileName = regexp.MustCompile(`^exports(?:_(.+))?\.go$`)

func main() {
	defer signals.AppQuit.Finalize()

	env := myenv.Must1(codegen.CreateEnvironment("github.com/gongt/go/cmd/exports", MAGIC_STRING))

	if !allowFileName.MatchString(env.InputPath().Base().Name) {
		panic(errors.New("此生成器仅支持在 exports_xxx.go 中使用"))
	}

	var opts Options
	myenv.Must(env.ParseArgs(&opts))
	myenv.Must(env.NoMoreArgs())

	internalDir, err := env.InputPath().SetFilename("internal").RealpathExisting()

	if err != nil {
		panic(fmt.Errorf("无法找到 internal 目录: %w", err))
	}

	if stat, err := os.Stat(internalDir.Raw()); err != nil || !stat.IsDir() {
		panic(fmt.Errorf("%s必须是目录", internalDir))
	}

	if !opts.VerboseMode {
		log.SetOutput(io.Discard)
	}

	bs := myenv.Must1(internal.ParseFiles(env, internalDir)) // main
	if !opts.VerboseMode {
		log.SetOutput(os.Stderr)
	}

	pkgName := myenv.Must1(sourcecode.DetectPackageName(env.InputPath().Dir()))
	bs.SetPackageName(pkgName)

	writer := myenv.Must1(env.NewOutput(env.InputPath(), bs))

	codegen.WriteGeneratorComment(bs.Heading(), env.GeneratorName(), nil)

	myenv.Must(writer.WriteFile())

	signals.AppQuit.Set(0)
}
