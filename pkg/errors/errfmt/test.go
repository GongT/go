package errfmt

import (
	"os"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
)

func TestNoError(t *testing.T, err error) {
	if err == nil {
		return
	}

	os.Stderr.WriteString(FormatError(err, true))
	t.FailNow()
}

// PanicBox 包裹一个错误，如果它实现了 [internal.StackTrace] 接口，则Error()会打印堆栈信息，如没有实现则不会被包裹。
// 这会在最终的错误输出中显示堆栈信息，
// 无法避免panic本身也打印堆栈信息
//
// 可使用 errors.Unwrap() 获取原始错误
type panicBox struct {
	error
}

func (p *panicBox) Error() string {
	return FormatError(p.error, myenv.StderrIsTerminal)
}

func (p *panicBox) Unwrap() error {
	return p.error
}

// Panic 如果err不为nil，以 [panicBox] 方式抛出异常
func Panic(err error) {
	if err == nil {
		return
	}

	if strace := errors.GetStackTrace(err); len(strace) > 0 {
		panic(&panicBox{err})
	} else {
		panic(err)
	}
}
