package tags

import (
	"os"

	"github.com/gongt/go/internal/myenv"
)

func init() {
	if myenv.IsTesting {
		Enable("*")
	} else {
		envDebug := os.Getenv("DEBUG")
		Enable(envDebug)
	}
}
