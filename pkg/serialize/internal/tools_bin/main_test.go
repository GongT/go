package tools_bin

import (
	"log"
	"runtime"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/fsys/fpath"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/stretchr/testify/require"
)

// 定位仓库中的 assets/serialize 目录，不依赖测试执行时的工作目录
func assetsDir(t *testing.T) string {
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "无法获取当前测试文件路径")

	mod, err := sourcecode.FindGoMod(fpath.INew(thisFile))
	require.NoError(t, err, "无法定位go.mod")

	return mod.Dir().Join("assets", "serialize").Raw()
}

// TestGenerateFixtures 对assets/serialize下的真实示例结构体运行代码生成核心逻辑，
// 校验生成出的Marshal/Unmarshal代码在语法上是合法的Go源码
func TestGenerateFixtures(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	rdr := sourcecode.NewPackageReader()
	infos, err := rdr.ReadFilesType(assetsDir(t))
	require.NoError(t, err)

	generated := 0
	for path, file := range infos {
		structs := collectStructs(file)
		if len(structs) == 0 {
			continue
		}

		buf := sourcecode.NewGoFileBuffer()
		buf.SetPackageName(file.Container().Name())

		g := &genCtx{
			buf:        buf,
			file:       file,
			info:       file.Container().TypesInfo(),
			helpersPkg: buf.AddImport("github.com/gongt/go/pkg/serialize"),
			errorsPkg:  buf.AddImport("github.com/gongt/go/pkg/errors"),
		}

		for _, st := range structs {
			g.emitMarshalFunc(st)
			g.emitUnmarshalFunc(st)
		}

		src := buf.Bytes()
		require.NotEmpty(t, src, "文件%s生成的代码格式化失败", path)

		log.Printf("文件[%s]生成的代码长度: %d", path, len(src))

		generated++
	}

	require.Positive(t, generated, "assets/serialize下应至少有一个包含结构体的文件参与生成")
}
