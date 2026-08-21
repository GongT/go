package internal

import (
	"path/filepath"
	"strings"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/errors"
)

var PathErr = errors.NewTemplate("路径错误")

// 清理路径，去掉多余的斜杠和点
func Clean(p string) string {
	return filepath.ToSlash(filepath.Clean(p))
}

func AssertValidPath(p string) {
	if err := CheckValidPath(p); err != nil {
		panic(err)
	}
}

func AssertValidPaths(paths []string) {
	if err := CheckValidPaths(paths); err != nil {
		panic(err)
	}
}

func CheckValidPath(p string) error {
	if myenv.IsDebug {
		if p == "" {
			return PathErr.New("参数包含空值")
		}
		if strings.Contains(p, "\\") {
			return PathErr.New("参数包含反斜杠: %s", p)
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
