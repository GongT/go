package terminal

import (
	"os"

	"github.com/gongt/go/internal/myenv"
	"golang.org/x/term"
)

const tty_device_name = "/dev/tty"

func IsTTY(stream *os.File) bool {
	return term.IsTerminal(int(stream.Fd()))
}

func HasTTY() bool {
	if myenv.IsWindows {
		return false
	}

	if stat, err := os.Stat(tty_device_name); err != nil {
		return false
	} else {
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
}

func OpenTTY() (*os.File, error) {
	return os.OpenFile(tty_device_name, os.O_RDWR, 0)
}
