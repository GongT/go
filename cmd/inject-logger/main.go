package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/signals"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/codegen"
)

const MAGIC_STRING = "9ed0de76-4d48-4068-b392-a152cd6d3ab1"

type Options struct {
	TagPrefix string `long:"prefix" description:"Prefix to add to the log tag" default:""`
}

func main() {
	defer signals.AppQuit.Finalize()

	env, err := codegen.CreateEnvironment("inject-logger", MAGIC_STRING)
	if err != nil {
		panic(err)
	}

	opts := &Options{}
	if err := env.ParseArgs(opts); err != nil {
		panic(err)
	}
	if len(env.Args) > 0 {
		panic(errors.NewAnonymous("unexpected arguments").WithDetails("args", env.Args))
	}

	original, err := env.ReadAsString()
	if err != nil {
		panic(err)
	}

	file, err := env.NewOutput(env.InputFile, "")
	if err != nil {
		panic(err)
	}

	sb := file.Output()

	gen_found := false
	for line := range strings.Lines(original) {
		if strings.HasPrefix(line, "//go:generate ") {
			gen_found = true
			sb.WriteString(line)
		}
	}

	goMod, err := env.FindGoMod()
	if err != nil {
		panic(err)
	}
	rootDir := path.Dir(goMod)
	finalDir := path.Dir(env.InputFile)

	if !strings.HasPrefix(finalDir, rootDir) || finalDir == rootDir {
		panic(errors.NewAnonymous("finalDir is not a subdirectory of rootDir").WithDetails("finalDir", finalDir, "rootDir", rootDir))
	}

	if !gen_found {
		sb.WriteString("//go:generate github.com/gongt/go/cmd/inject-logger")
		if opts.TagPrefix != "" {
			fmt.Fprintf(sb, " --prefix=%q", opts.TagPrefix)
		}
		sb.WriteString("\n")
	}

	names := []string{}
	itr := finalDir
	for itr != rootDir {
		pkgName, err := sourcecode.DetectPackageName(itr)
		if err != nil {
			panic(err)
		}
		names = append([]string{pkgName}, names...)
		itr = path.Dir(itr)
	}

	sb.WriteString(`
import (
	"github.com/gongt/go/pkg/logger"
)

`)

	if opts.TagPrefix != "" {
		names = append([]string{opts.TagPrefix}, names...)
	}
	fmt.Fprintf(sb, "const _dbg_tag = %q\n", strings.Join(names, ":"))

	sb.WriteString(`
func debug(msg string) {
	logger.DLog(_dbg_tag, msg)
}
var _ = debug
func debugf(fmt string, args ...any) {
	logger.DLogF(_dbg_tag, fmt, args...)
}
var _ = debugf

func print(fmt string) {
	logger.Log(_dbg_tag, fmt)
}
var _ = print
func printf(fmt string, args ...any) {
	logger.LogF(_dbg_tag, fmt, args...)
}
var _ = printf
`)

	err = file.WriteFile()
	if err != nil {
		panic(err)
	}

	signals.AppQuit.Set(0)
}
