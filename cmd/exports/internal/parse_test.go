package internal

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors/errfmt"
	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/source_code/codegen"
	"github.com/stretchr/testify/require"
)

func Test_ParseFiles(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	me, _ := myenv.CurrentFileLine()

	f := fpath.New(me).Join("../../../../assets/exports/internal")
	errfmt.TestNoError(t, f.RealpathExisting())
	log.Printf("输入: %s", f.Raw())

	buff, err := ParseFiles(codegen.TestingEnvironment(t), f.Immutable())

	errfmt.TestNoError(t, err)

	buff.Heading().WriteString("// SOME HEADING TEXT")
	buff.SetPackageName("test_pkg")

	expectFile := f.Immutable().Join("../output.go").Raw()
	content, err := os.ReadFile(expectFile)
	errfmt.TestNoError(t, err)

	defer func() {
		if t.Failed() {
			debugFile := f.Immutable().Join("../.debug.go")
			log.Printf("测试失败，调试文件: %s", debugFile)
			os.WriteFile(debugFile.Raw(), buff.Bytes(), 0644)
		}
	}()

	require.NoError(t, buff.CheckSyntax())
	require.Equal(t, strings.TrimSpace(string(content)), strings.TrimSpace(buff.String()))
}
