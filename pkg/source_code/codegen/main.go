package codegen

import (
	"os"
	"testing"
	"unsafe"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys"
	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/interfaces"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/internal/cgen"
	"github.com/gongt/go/pkg/source_code/internal/writer"
	"github.com/gongt/go/pkg/types"
	"github.com/jessevdk/go-flags"
)

type globalArgs struct {
	InputFile string `long:"input" description:"Path to the input file" env:"GOFILE"`
}

type GeneratorEnvironment = *generatorEnvironment
type generatorEnvironment struct {
	InputFile    *fpath.IPath
	contentCache []byte

	Cwd           *fpath.IPath // 启动时的工作目录
	workspaceRoot *fpath.IPath // 工作区根目录，go.mod所在目录

	InitialArgs []string // 初始命令行参数
	unusedArgs  []string // 从未使用的命令行参数

	GeneratorFullName string // 工具执行方式，例如 "github.com/gongt/go/cmd/exports"
	magicBytes        []byte // GUID
	goMod             *sourcecode.GoModFile
}

var CurrentGenerator GeneratorEnvironment

func GetEnvironment() GeneratorEnvironment {
	if CurrentGenerator == nil {
		panic(errors.NewAnonymous("当前没有生成器环境"))
	}
	return CurrentGenerator
}

func TestingEnvironment(t *testing.T) GeneratorEnvironment {
	cwd := fpath.INew(myenv.Must1(os.Getwd()))
	ret := &generatorEnvironment{
		GeneratorFullName: "testing",
		magicBytes:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Cwd:               cwd,
		InitialArgs:       []string{},
		unusedArgs:        []string{},
		InputFile:         cwd.Join("_virtual_input_file.go"),
		contentCache:      []byte{},
	}

	var err error

	ret.goMod, err = sourcecode.FindGoMod(cwd)
	if err != nil {
		panic(errors.Extend(err, "无法找到go.mod文件"))
	}
	ret.workspaceRoot = ret.goMod.Dir()

	return ret
}

func CreateEnvironment(myName string, magic string) (GeneratorEnvironment, error) {
	var ret = &generatorEnvironment{
		GeneratorFullName: myName,
		magicBytes:        unsafe.Slice(unsafe.StringData(magic), len(magic)),
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Extend(err, "无法获取当前工作目录")
	}
	ret.Cwd = fpath.INew(cwd)

	opts := &globalArgs{}
	args, err := flags.NewParser(opts, flags.IgnoreUnknown).Parse()
	if err != nil {
		return nil, errors.Extend(err, "无法解析命令行参数")
	}

	ret.InitialArgs = args

	if opts.InputFile == "" {
		return nil, errors.NewAnonymous("未发现GOFILE环境，必须通过--input提供输入文件路径")
	}

	inFile := ret.Cwd.Resolve(opts.InputFile)
	inFile, err = inFile.RealpathMissing()
	if err != nil {
		return nil, err
	}
	ret.InputFile = inFile

	ret.goMod, err = sourcecode.FindGoMod(ret.InputFile)
	if err != nil {
		return nil, errors.Extend(err, "无法找到go.mod文件")
	}
	ret.workspaceRoot = ret.goMod.Dir()

	data, err := fsys.ReadFileOrEmpty(ret.InputFile)
	if err != nil {
		return nil, errors.Extend(err, "无法读取输入文件").WithDetails("path", ret.InputFile)
	}
	ret.contentCache = data

	CurrentGenerator = ret
	return CurrentGenerator, nil
}

func (env GeneratorEnvironment) GeneratorName() string       { return env.GeneratorFullName }
func (env GeneratorEnvironment) InputPath() *fpath.IPath     { return env.InputFile }
func (env GeneratorEnvironment) InputPathRaw() string        { return env.InputFile.Raw() }
func (env GeneratorEnvironment) InputContent() []byte        { return env.contentCache }
func (env GeneratorEnvironment) Magic() []byte               { return env.magicBytes }
func (env GeneratorEnvironment) WorkspaceRoot() *fpath.IPath { return env.workspaceRoot }
func (env GeneratorEnvironment) UnusedArgs() []string        { return env.unusedArgs }

func (env GeneratorEnvironment) NoMoreArgs() error {
	un := env.UnusedArgs()
	if len(un) > 0 {
		panic(errors.NewAnonymous("未知命令行参数: %s", un[0]).WithDetails("args", un))
	}
	return nil
}

func (env GeneratorEnvironment) ParseArgs(output any) error {
	unused, err := flags.NewParser(output, flags.IgnoreUnknown).ParseArgs(env.InitialArgs)
	if err != nil {
		return errors.Extend(err, "无法解析命令行参数")
	}

	types.IntersectInplace(&env.unusedArgs, unused)

	return nil
}

func (env *generatorEnvironment) NewOutput[T fpath.PathLike](ipath T, content interfaces.ToBytes) (*cgen.SafeTextWriter, error) {
	path := fpath.ToImmutable(ipath)

	if path.IsAbs() {
		if !fpath.IsLocal(path, env.workspaceRoot) {
			return nil, errors.NewAnonymous("拒绝向项目外输出文件").WithDetails("path", path, "projectDir", env.workspaceRoot)
		}
	} else {
		path = env.Cwd.Join(path)
	}

	r := cgen.NewTextWriter(path, env.magicBytes, content)

	if buffer, ok := content.(writer.GoFileBuffer); ok {
		buffer.Heading().WriteString(env.DoNotModify())
	}

	if err := r.Prepare(); err != nil {
		return nil, errors.Extend(err, "无法准备输出文件").WithDetails("path", r.Path())
	}

	return r, nil
}

func (env GeneratorEnvironment) DoNotModify() string {
	return cgen.CreateSafeGuardComment(env.magicBytes, env.GeneratorFullName)
}

func (env GeneratorEnvironment) GoMod() *sourcecode.GoModFile {
	if env.goMod != nil {
		return env.goMod
	}
	modFile := sourcecode.OpenGoMod(env.workspaceRoot.Join("go.mod"))
	env.goMod = modFile
	return modFile
}
