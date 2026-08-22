package main

import (
	"os"

	"github.com/gongt/go/pkg/string_helper/strtools"
)

func main() {
	text, err := strtools.TinyHashStream(os.Stdin)
	if err != nil {
		panic(err)
	}
	os.Stdout.WriteString(text)
}
