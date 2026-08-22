package codegen

import (
	"os"
	"path/filepath"
	"unsafe"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys/fpath"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/internal"
	"github.com/gongt/go/pkg/source_code/internal/writer"
	"github.com/jessevdk/go-flags"
)

var _ internal.CodegenEnvironmentInterface = (*GeneratorEnvironment)(nil)

type GeneratorEnvironment struct {
	InputFile         string `long:"input" description:"Path to the input file" env:"GOFILE"`
	Cwd               string
	Args              []string
	GeneratorFullName string // 工具执行方式，例如 "github.com/gongt/go/cmd/exports"
	Magic             []byte // GUID
	goMod             *sourcecode.GoModFile
}

func CreateEnvironment(myName string, magic string) (*GeneratorEnvironment, error) {
	var ret = &GeneratorEnvironment{
		GeneratorFullName: myName,
		Magic:             unsafe.Slice(unsafe.StringData(magic), len(magic)),
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Extend(err, "无法获取当前工作目录")
	}
	ret.Cwd = cwd

	args, err := flags.NewParser(ret, flags.IgnoreUnknown).Parse()
	if err != nil {
		return nil, errors.Extend(err, "无法解析命令行参数")
	}

	ret.Args = args

	if ret.InputFile == "" {
		return nil, errors.NewAnonymous("必须通过--input提供输入文件路径")
	}

	f := fpath.New(ret.Cwd).Join(ret.InputFile)
	err = f.RealpathExisting()
	if err != nil {
		return nil, err
	}

	ret.InputFile = f.Raw()
	if stat, err := os.Stat(ret.InputFile); err != nil || stat.IsDir() {
		if os.IsNotExist(err) {
			// 允许不存在
		} else {
			return nil, errors.NewAnonymous("输入文件路径无效").WithDetails("path", ret.InputFile)
		}
	}

	return ret, nil
}

func (env *GeneratorEnvironment) ReadAsString() (string, error) {
	data, err := os.ReadFile(env.InputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", errors.Extend(err, "无法读取输入文件").WithDetails("path", env.InputFile)
	}
	return string(data), nil
}

func (env *GeneratorEnvironment) ParseArgs(output any) error {
	args, err := flags.NewParser(output, flags.IgnoreUnknown).ParseArgs(env.Args)
	if err != nil {
		return errors.Extend(err, "无法解析命令行参数")
	}
	env.Args = args
	return nil
}

func (env *GeneratorEnvironment) NewOutput(path string, content toBytes) (*SafeTextWriter, error) {
	if filepath.IsAbs(path) {
		mod, err := env.FindGoMod()
		if err != nil {
			return nil, errors.Extend(err, "无法查找go.mod文件")
		}
		if !fpath.IsLocal(path, mod.Dir()) {
			return nil, errors.NewAnonymous("拒绝向项目外输出文件").WithDetails("path", path, "projectDir", mod.Dir())
		}
	} else {
		path = filepath.Join(env.Cwd, path)
	}
	f := fpath.New(path)

	r := NewTextWriter(f.Raw(), []byte(env.Magic), content)

	if buffer, ok := content.(*writer.GoFileBuffer); ok {
		buffer.Heading().WriteString(env.DoNotModify())
	}

	if err := r.Prepare(); err != nil {
		return nil, errors.Extend(err, "无法准备输出文件").WithDetails("path", r.Path())
	}

	return r, nil
}

func (env *GeneratorEnvironment) DoNotModify() string {
	return internal.CreateSafeGuardComment(env.Magic, env.GeneratorFullName)
}

func (env *GeneratorEnvironment) FindGoMod() (*sourcecode.GoModFile, error) {
	if env.goMod != nil {
		return env.goMod, nil
	}
	modFile, err := sourcecode.FindGoMod(env.Cwd)
	if err != nil {
		return nil, err
	}
	env.goMod = modFile
	return modFile, nil
}
