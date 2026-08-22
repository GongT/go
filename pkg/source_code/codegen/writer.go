package codegen

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/terminal"
)

type toBytes interface {
	Bytes() []byte
}

type SafeTextWriter struct {
	path    *fpath.IPath
	magic   []byte
	content toBytes

	// 判断回调，当文件存在且不包含magic时，调用此函数判断是否安全覆盖，默认直接返回false
	//
	// 如果返回false，执行提示用户是否覆盖逻辑
	IsSafe func(string) bool
	Mkdir  bool // 是否创建目录

	prepared bool
}

func NewTextWriter(path string, magicBytes []byte, content toBytes) *SafeTextWriter {
	r := &SafeTextWriter{
		path:    fpath.INew(path),
		magic:   magicBytes,
		content: content,
	}
	return r
}

func (w *SafeTextWriter) Path() string {
	return w.path.Raw()
}

func (w *SafeTextWriter) SetPath(path string) {
	w.path = fpath.INew(path)
}

func (w *SafeTextWriter) Prepare() error {
	if w.prepared {
		return nil
	}

	dir := w.path.Dir()
	if w.Mkdir {
		if _, err := os.Stat(dir.Raw()); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(dir.Raw(), 0755); err != nil {
					return errors.Extend(err, "无法创建目录").WithDetails("path", dir)
				}
			} else {
				return errors.Extend(err, "无法访问目录").WithDetails("path", dir)
			}
		}
		w.Mkdir = false
	}

	w.prepared = true

	return nil
}

func (w *SafeTextWriter) WriteFile() error {
	w.safeGuard()
	return w.ForceOverride()
}

func (w *SafeTextWriter) ForceOverride() error {
	if !w.prepared {
		panic("ForceOverride()/WriteFile()之前必须先调用Prepare()")
	}

	fw, err := os.Create(w.path.Raw())
	if err != nil {
		return errors.Extend(err, "无法打开文件").WithDetails("path", w.path.Raw())
	}
	defer fw.Close()

	bytes := w.content.Bytes()
	if _, err := fw.Write(bytes); err != nil {
		return errors.Extend(err, "文件写入错误").WithDetails("path", w.path.Raw())
	}

	log.Println("写入文件: ", w.path.Raw())

	return nil
}

func (w *SafeTextWriter) safeGuard() {
	content, err := os.ReadFile(w.path.Raw())
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	if bytes.Contains(content, w.magic) {
		return
	}

	if w.IsSafe != nil && w.IsSafe(w.path.Raw()) {
		return
	}

	if terminal.HasTTY() {
		w.ask()
	} else {
		panic(fmt.Sprintf(`  文件 %s 不为空，且不是由此工具生成的。你需要选择其中一种方式来允许覆盖行为:
  - 删除它
  - 在其中添加字符串 "%s" （不包括引号）`, w.path.Raw(), string(w.magic)))
	}
}

func (w *SafeTextWriter) ask() {
	fmt.Printf("  文件 %s 已存在，且不是由此工具生成的。", w.path.Raw())
	fmt.Println("  如果继续，文件将被覆盖，任何更改都将丢失。")
	fmt.Print("  是否继续？ (y/n): ")

	var response string

	f, err := terminal.OpenTTY()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = fmt.Fscanln(f, &response)
	if err != nil {
		response = "n"
	}

	rsp := strings.ToLower(response)
	if rsp != "y" && rsp != "yes" && rsp != "是" {
		fmt.Print("中止操作...\n\n")
		os.Exit(1)
	}
}
