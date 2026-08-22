// @exported
package config_loader

import (
	"regexp"
)

// Regexp 是一个包装了 regexp.Regexp 的结构体，用于从配置文件中读取
type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) FromString(s string) error {
	re, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	r.Regexp = re
	return nil
}
