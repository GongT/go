package fpath

import (
	"strings"

	"github.com/gongt/go/internal/myenv"
)

// File 表示一个文件名（不包括路径），它是一个不可变对象，不要修改Name字段的值
type File struct {
	Name string
}

func NewFile(name string) *File {
	if myenv.IsDebug {
		if strings.Contains(name, "/") || strings.Contains(name, "\\") {
			panic(PathErr.New("文件名不能包含路径分隔符"))
		}
	}
	return &File{
		Name: name,
	}
}

func (f *File) String() string {
	return f.Name
}

// 将文件名拆分为主干和扩展名
//
// 其中扩展名包含点号
func (f *File) SplitExt() (string, string) {
	i := strings.LastIndexByte(f.Name, '.')
	if i == -1 { // 没有扩展名
		return f.Name, ""
	}
	return f.Name[:i], f.Name[i:]
}

// 返回文件名的扩展名部分
func (f *File) Ext() string {
	_, ex := f.SplitExt()
	return ex
}

// 返回文件名的主干部分，即去掉扩展名的部分
func (f *File) Stem() string {
	ba, _ := f.SplitExt()
	return ba
}
