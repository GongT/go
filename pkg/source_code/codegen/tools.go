package codegen

import (
	"fmt"
	"io"
)

// WriteGeneratorComment 把当前生成器的go:generate注释写入到writer中
//
// 一般只会用在自我生成的代码文件中
func WriteGeneratorComment(writer io.Writer, generatorName string, arguments []string) {
	fmt.Fprintf(writer, "//go:generate go run %s", generatorName)
	for _, arg := range arguments {
		fmt.Fprintf(writer, " %s", arg)
	}
	fmt.Fprintln(writer)
}
