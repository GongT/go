package myenv

import (
	"os"

	"golang.org/x/term"
)

var StderrIsTerminal bool = term.IsTerminal(int(os.Stderr.Fd()))
var StdoutIsTerminal bool = term.IsTerminal(int(os.Stdout.Fd()))
