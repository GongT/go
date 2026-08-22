package writer

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/gongt/go/pkg/string_helper/strtools"
	writercache "github.com/gongt/go/pkg/string_helper/writer_cache"
)

var _ writercache.AllWriter = (*GoFileBuffer)(nil)

type GoFileBuffer struct {
	writercache.WriteEvent

	buff bytes.Buffer

	header      bytes.Buffer
	packageName string
	imports     map[string]string

	formatted []byte
}

// NewGoFileBuffer 通过拼接字符串创建一个Go语言文件
func NewGoFileBuffer() *GoFileBuffer {
	r := &GoFileBuffer{
		imports: make(map[string]string),
	}
	r.buff.Grow(10240)

	r.RouteTo(&r.buff, r.invalidateCache)

	return r
}

func (c *GoFileBuffer) AddImport(pkgPath string) string {
	alias, exists := c.imports[pkgPath]
	if !exists {
		c.imports[pkgPath] = strtools.TinyHash(pkgPath)
		alias = c.imports[pkgPath]
	}
	return alias
}

// SetHeading 设置文件头部注释，通常用于添加版权信息或生成标记
func (c *GoFileBuffer) Heading() writercache.AllWriter {
	return writercache.New(&c.header, c.invalidateCache)
}

// SetPackageName 设置Go文件的包名，不设置会出错
func (c *GoFileBuffer) SetPackageName(packageName string) {
	c.packageName = packageName
}

func (c *GoFileBuffer) Body() writercache.AllWriter {
	return c
}

func (c *GoFileBuffer) finalize() {
	buff := &bytes.Buffer{}

	// heading
	buff.Write(c.header.Bytes())
	buff.WriteByte('\n')

	// package
	if c.packageName == "" {
		panic("没有设置 packageName")
	}
	buff.WriteString("package ")
	buff.WriteString(c.packageName)
	buff.WriteString("\n\n")

	// import
	buff.WriteString("import (\n")
	paths := make([]string, 0, len(c.imports))
	for pkgPath := range c.imports {
		paths = append(paths, pkgPath)
	}
	sort.Strings(paths)
	standardPaths := make([]string, 0, len(paths))
	externalPaths := make([]string, 0, len(paths))
	for _, pkgPath := range paths {
		if strings.Contains(pkgPath, ".") {
			externalPaths = append(externalPaths, pkgPath)
		} else {
			standardPaths = append(standardPaths, pkgPath)
		}
	}
	paths = append(standardPaths, externalPaths...)
	for index, pkgPath := range paths {
		if index == len(standardPaths) && len(standardPaths) > 0 && len(externalPaths) > 0 {
			buff.WriteByte('\n')
		}
		alias := c.imports[pkgPath]
		buff.WriteByte('\t')
		buff.WriteString(alias)
		buff.WriteString(" \"")
		buff.WriteString(pkgPath)
		buff.WriteString("\"\n")
	}
	buff.WriteString(")\n")

	// body
	body := c.buff.Bytes()
	buff.Write(body)

	src, err := format.Source(buff.Bytes())
	if err != nil {
		fmt.Printf("Error formatting code: %v\n", err)
		return
	}

	c.Reset()
	c.formatted = src
}

func (c *GoFileBuffer) Bytes() []byte {
	if c.formatted == nil {
		c.finalize()
	}
	return c.formatted
}

func (c *GoFileBuffer) String() string {
	return string(c.Bytes())
}

func (c *GoFileBuffer) invalidateCache() {
	c.formatted = nil
}
