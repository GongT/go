package signals

import (
	"fmt"
	"log"
	"os"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/errors/errfmt"

	"golang.org/x/term"
)

var AppQuit *quitHandler = New()
var StderrIsTerminal bool = term.IsTerminal(int(os.Stderr.Fd()))

type quitHandler struct {
	ch chan struct{}

	codeSet  bool
	exitCode int
}

func New() *quitHandler {
	return &quitHandler{
		ch:       make(chan struct{}),
		exitCode: -1,
		codeSet:  false,
	}
}

func (h *quitHandler) Wait() {
	<-h.ch
}

func (h *quitHandler) MainFinish() {
	close(h.ch)

	if h.codeSet {
		if StderrIsTerminal {
			fmt.Fprintln(os.Stderr, "bye, bye~")
		}

		os.Exit(h.exitCode)
	} else {
		log.Println("退出时没有设置退出码，默认使用 1")
		os.Exit(1)
	}
}

// 在main()函数中使用defer AppQuit.Finalize()，可以捕获panic并打印错误信息
func (h *quitHandler) Finalize() {
	if r := recover(); r != nil {
		AppQuit.Fatal(r)
	} else {
		AppQuit.MainFinish()
	}
}

func (h *quitHandler) Fatal(e any) {
	close(h.ch)
	var code int

	r := errfmt.FormatError(e, StderrIsTerminal)
	os.Stderr.WriteString(r)

	if err, ok := e.(error); ok {
		code, _ = errors.GetCode(err)
	}

	if code == 0 {
		code = -1
	}
	h.Overwrite(code)

	// log.Printf("最终退出码: %d\n", h.exitCode)
	os.Exit(h.exitCode)
}

// 设置退出码，强制覆盖之前的设置
func (h *quitHandler) Overwrite(code int) {
	h.codeSet = true
	h.exitCode = code
}

// 设置退出码，如果已经设置过，则不会覆盖
//
// 除非原来的是 0 且新的不是 0
func (h *quitHandler) Set(code int) {
	if !h.codeSet {
		h.codeSet = true
		h.exitCode = code
	} else if h.exitCode == 0 && code != 0 {
		h.exitCode = code
	}
}
