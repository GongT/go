package bin_rpc

import (
	"log"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/source_code/codegen"
)

func RunMain(env codegen.GeneratorEnvironment) {
	myenv.Must(env.NoMoreArgs())

	log.Println("hello world", env)
}
