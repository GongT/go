package fpath

import (
	"os"
)

// 如果路径不是绝对路径，使用工作目录作为基准转换成绝对路径
//
// 不常用，通常用 [Path.ToAbs]
func (p *Path) MaybeConvertAbsolute() error {
	if !p.IsAbs() {
		workingDirectory, err := os.Getwd()
		if err != nil {
			return PathErr.Wrap(err)
		}
		p.value = workingDirectory + "/" + p.value
	}
	return nil
}

// 如果路径不是绝对路径，使用工作目录作为基准转换成绝对路径
func (p *Path) ToAbs() {
	if err := p.MaybeConvertAbsolute(); err != nil {
		panic(err)
	}
}
