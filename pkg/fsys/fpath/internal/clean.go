package internal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
)

// 清理路径，去掉多余的斜杠和点
func Clean(p string) string {
	return filepath.ToSlash(filepath.Clean(p))
}

func AssertValidPath(p string) {
	if myenv.IsDebug {
		if p == "" {
			panic("路径包含空值")
		}
		if strings.Contains(p, "\\") {
			panic(fmt.Sprintf("路径中包含反斜杠: %s", p))
		}
	}
}

func AssertValidPaths(paths []string) {
	if myenv.IsDebug {
		for _, p := range paths {
			AssertValidPath(p)
		}
	}
}

func CheckValidPath(p string) error {
	if myenv.IsDebug {
		if p == "" {
			return errors.NewAnonymous("路径包含空值")
		}
		if strings.Contains(p, "\\") {
			return errors.NewAnonymous("路径中包含反斜杠: %s", p)
		}
	}
	return nil
}

func CheckValidPaths(paths []string) error {
	if myenv.IsDebug {
		for _, p := range paths {
			if err := CheckValidPath(p); err != nil {
				return err
			}
		}
	}
	return nil
}
