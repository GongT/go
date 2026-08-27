package main

import (
	"os"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/strings/strtools"
)

func main() {
	var result string

	switch len(os.Args) {
	case 1:
		result = myenv.Must1(strtools.TinyHashStream(os.Stdin))
	case 2:
		result = strtools.TinyHash([]byte(os.Args[1]))
	default:
		panic("invalid number of arguments")
	}

	os.Stdout.WriteString(result)

	if myenv.StdoutIsTerminal {
		os.Stdout.WriteString("\n")
	}
}
