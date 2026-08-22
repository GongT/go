package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
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

	env := myenv.Must(codegen.CreateEnvironment("github.com/gongt/go/cmd/inject-logger", MAGIC_STRING))

	opts := &Options{}
	myenv.MustNil(env.ParseArgs(opts))

	if len(env.Args) > 0 {
		panic(errors.NewAnonymous("unexpected arguments").WithDetails("args", env.Args))
	}

	outputDir := path.Dir(env.InputFile)
	outputBuffer := sourcecode.NewGoFileBuffer()

	fileWriter := myenv.Must(env.NewOutput(env.InputFile, outputBuffer))

	goMod := myenv.Must(env.FindGoMod())
	if !fpath.IsLocal(outputDir, goMod.Dir()) {
		panic(errors.NewAnonymous("输出文件不在项目内").WithDetails("outputDir", outputDir, "projectDir", goMod.Dir()))
	}

	pkgName := myenv.Must(sourcecode.DetectPackageName(outputDir))
	outputBuffer.SetPackageName(pkgName)

	var gArgs []string
	if opts.TagPrefix != "" {
		gArgs = append(gArgs, fmt.Sprintf("--prefix=%q", opts.TagPrefix))
	}
	codegen.WriteGeneratorComment(outputBuffer.Heading(), env.GeneratorFullName, gArgs)

	// 查找父目录，通过package名称组装日志tag
	ignoreSelf := func(ent string) bool {
		if sourcecode.SimulateGolangBuild(ent) {
			return true
		}
		if ent == path.Base(env.InputFile) {
			return true
		}
		return false
	}
	names := []string{}
	for dir := range fsys.ClimbingPath(outputDir) {
		var ignore func(string) bool = nil
		if dir.Raw() == outputDir {
			ignore = ignoreSelf
		}
		if fpath.IsEquals(dir, goMod.Dir()) {
			break
		}

		package_declarations := myenv.Must(sourcecode.DetectPackageNames(dir.Raw(), ignore))
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

	myenv.MustNil(fileWriter.WriteFile())

	signals.AppQuit.Set(0)
}
