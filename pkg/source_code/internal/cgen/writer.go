package cgen

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/gongt/go/pkg/interfaces"
	sourcecode "github.com/gongt/go/pkg/source_code"
	"github.com/gongt/go/pkg/source_code/internal/writer"
	"github.com/gongt/go/pkg/terminal"
)

type SafeTextWriter struct {
	path    *fpath.IPath
	magic   []byte
	content interfaces.ToBytes

	// 判断回调，当文件存在且不包含magic时，调用此函数判断是否安全覆盖，默认直接返回false
	//
	// 如果返回false，执行提示用户是否覆盖逻辑
	IsSafe func(string) bool
	Mkdir  bool // 是否创建目录

	prepared bool
}

func NewTextWriter[T fpath.PathLike](path T, magicBytes []byte, content interfaces.ToBytes) *SafeTextWriter {
	r := &SafeTextWriter{
		path:    fpath.ToImmutable(path),
		magic:   magicBytes,
		content: content,
	}
	return r
}

func (w *SafeTextWriter) IsGoFile() bool {
	_, ok := w.content.(writer.GoFileBuffer)
	return ok
}

func (w *SafeTextWriter) LearnPackageName() error {
	if buff, ok := w.content.(writer.GoFileBuffer); ok {
		pkgName, err := sourcecode.DetectPackageNameExcept(w.path.Dir(), w.path.Base().Name)
		if err != nil {
			return err
		}
		buff.SetPackageName(pkgName)
	} else {
		return errors.NewAnonymous("内容不是Go文件缓冲区")
	}
	return nil
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

	bytes := w.content.Bytes()

	name := "." + w.path.Base().String() + ".tmp"
	temp := w.path.SetBase(name)
	if err := os.WriteFile(temp.Raw(), bytes, 0644); err != nil {
		return errors.Extend(err, "文件写入错误").WithDetails("path", w.path.Raw())
	}

	if err := os.Rename(temp.Raw(), w.path.Raw()); err != nil {
		return errors.Extend(err, "文件重命名错误").WithDetails("from", temp.Raw(), "to", w.path.Raw())
	}

	log.Printf("写入%d字节，文件: %s", len(bytes), w.path.Raw())

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

	if len(content) == 0 {
		return
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
	fmt.Printf("  文件 %s 已存在，且不是由此工具生成的。\n", w.path.Raw())
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
