package main

import (
	"fmt"
	"strings"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/fsys"
	"github.com/gongt/go/pkg/fsys/fpath"
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

	env := myenv.Must1(codegen.CreateEnvironment("github.com/gongt/go/cmd/inject-logger", MAGIC_STRING))

	opts := &Options{}
	myenv.Must(env.ParseArgs(opts))

	myenv.Must(env.NoMoreArgs())

	outputDir := env.InputPath().Dir()
	outputBuffer := sourcecode.NewGoFileBuffer()

	fileWriter := myenv.Must1(env.NewOutput(env.InputPath().Raw(), outputBuffer))

	pkgName := myenv.Must1(sourcecode.DetectPackageName(outputDir))
	outputBuffer.SetPackageName(pkgName)

	var gArgs []string
	if opts.TagPrefix != "" {
		gArgs = append(gArgs, fmt.Sprintf("--prefix=%q", opts.TagPrefix))
	}
	codegen.WriteGeneratorComment(outputBuffer.Heading(), env.GeneratorName(), gArgs)

	// 查找父目录，通过package名称组装日志tag
	ignoreSelf := func(ent string) bool {
		if sourcecode.SimulateGolangBuild(ent) {
			return true
		}
		if ent == env.InputPath().Base().Name {
			return true
		}
		return false
	}
	names := []string{}
	for dir := range fsys.ClimbingPath(outputDir) {
		var ignore func(string) bool = nil
		if fpath.IsEquals(dir, outputDir) {
			ignore = ignoreSelf
		}
		if fpath.IsEquals(dir, env.WorkspaceRoot()) {
			break
		}

		package_declarations := myenv.Must1(sourcecode.DetectPackageNames(dir.Raw(), ignore))
		names = append([]string{package_declarations[0]}, names...)
	}

	// 生成代码
	loggerPkg := outputBuffer.AddImport("github.com/gongt/go/pkg/logger")

	if opts.TagPrefix != "" {
		names = append([]string{opts.TagPrefix}, names...)
	}
	fmt.Fprintf(outputBuffer.Body(), "const _dbg_tag = %q\n", strings.Join(names, ":"))

	outputBuffer.Body().WriteString(fmt.Sprintf(`
func debug(msg string) {
	%s.DLog(_dbg_tag, msg)
}
var _ = debug
func debugf(fmt string, args ...any) {
	%s.DLogF(_dbg_tag, fmt, args...)
}
var _ = debugf

func print(fmt string) {
	%s.Log(_dbg_tag, fmt)
}
var _ = print
func printf(fmt string, args ...any) {
	%s.LogF(_dbg_tag, fmt, args...)
}
var _ = printf
`, loggerPkg, loggerPkg, loggerPkg, loggerPkg))

	myenv.Must(fileWriter.WriteFile())

	signals.AppQuit.Set(0)
}
