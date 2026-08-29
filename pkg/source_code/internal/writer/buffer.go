package writer

import (
	"bytes"
	"fmt"
	"go/format"
	"go/scanner"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/interfaces"
	"github.com/gongt/go/pkg/source_code/stl"
	CSI "github.com/gongt/go/pkg/strings/csi"
	"github.com/gongt/go/pkg/strings/strtools"
	writercache "github.com/gongt/go/pkg/strings/writer_cache"
)

var _ interfaces.ModernWriter = (GoFileBuffer)(nil)

// GoFileBuffer 通过拼接字符串创建一个Go语言文件
type GoFileBuffer = *goFileBuffer

type goFileBuffer struct {
	writercache.WriteEvent

	buff      bytes.Buffer
	formatted []byte
	checked   bool

	header      bytes.Buffer
	packageName string
	packagePath string

	// 路径 -> 别名映射
	imports map[string]string
	exports map[string]any

	NamePackage func(packageName string, buff GoFileBuffer) string
}

func NewGoFileBuffer() GoFileBuffer {
	r := &goFileBuffer{
		imports:     make(map[string]string),
		exports:     make(map[string]any),
		NamePackage: OriginalName,
	}
	r.buff.Grow(10240)

	r.RouteTo(&r.buff, r.invalidateCache)

	return r
}

func OriginalName(pkgPath string, _ GoFileBuffer) string {
	return filepath.Base(pkgPath)
}

func IndexName(pkgPath string, buff GoFileBuffer) string {
	return fmt.Sprintf("pkg%d", len(buff.imports))
}

func HashName(pkgPath string, _ GoFileBuffer) string {
	return strtools.TinyHash(pkgPath)
}

// AddImport 添加一个导入包，返回该包的别名，可能返回空字符串
func (c GoFileBuffer) AddImport(pkgPath string) string {
	if c.packagePath == pkgPath {
		return ""
	}
	alias, exists := c.imports[pkgPath]
	if !exists {
		alias = c.NamePackage(pkgPath, c)

		base := alias
		i := 1
		for {
			for _, existingAlias := range c.imports {
				if existingAlias == alias {
					alias = fmt.Sprintf("%s%d", base, i)
					i++
					continue
				}
			}
			break
		}

		c.imports[pkgPath] = alias
	}
	return alias
}

func (c GoFileBuffer) QualifyTypeName(pkgPath string, typeName string) string {
	alias := c.AddImport(pkgPath)
	if alias == "" {
		return typeName
	}
	return fmt.Sprintf("%s.%s", alias, typeName)
}

func (c GoFileBuffer) Qualifier(pkg *types.Package) string {
	return c.AddImport(pkg.Path())
}

// AddExport 将符号添加到导出列表中，如果符号已经存在则返回 true，否则返回 false
//
// 纯记录，不会对代码做任何操作
func (c GoFileBuffer) AddExport(symbol string, info any) bool {
	if _, exists := c.exports[symbol]; !exists {
		c.exports[symbol] = info
		return false
	} else {
		return true
	}
}

// SetHeading 设置文件头部注释，通常用于添加版权信息或生成标记
func (c GoFileBuffer) Heading() interfaces.ModernWriter {
	return writercache.New(&c.header, c.invalidateCache)
}

// SetPackageName 设置Go文件的包名，不设置会出错
func (c GoFileBuffer) SetPackageName(packageName string) {
	c.packageName = packageName
}

// SetPackagePath 设置Go文件的包路径，设置后AddImport此包返回空字符串
func (c GoFileBuffer) SetPackagePath(packagePath string) {
	c.packagePath = packagePath
}

func (c GoFileBuffer) Body() interfaces.ModernWriter {
	return c
}

// CheckSyntax 结束生成Go文件的过程，并输出格式化错误到日志
func (c GoFileBuffer) CheckSyntax() error {
	c.checked = true

	err := c.finalize()
	if err == nil {
		return nil
	}
	content := c.composit()

	fmt.Printf("%v\n", err)
	if e, ok := errors.AsType[scanner.ErrorList](err); ok {
		fmt.Printf("--------------------------\n%s\n--------------------------\n", markLine(content, e[0].Pos))
	} else {
		fmt.Printf("--------------------------\n%s\n--------------------------\n", content)
	}

	return err
}

func (c GoFileBuffer) composit() string {
	buff := &strings.Builder{}

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

	fmtImportLine := func(pkgPath string) {
		alias := c.imports[pkgPath]
		buff.WriteByte('\t')
		if filepath.Base(pkgPath) != alias {
			buff.WriteString(alias)
			buff.WriteByte(' ')
		}
		buff.WriteByte('"')
		buff.WriteString(pkgPath)
		buff.WriteString("\"\n")
	}

	for pkgPath := range c.imports {
		if !stl.IsStandardImport(pkgPath) {
			continue
		}
		fmtImportLine(pkgPath)
	}
	buff.WriteString("\n")

	for pkgPath := range c.imports {
		if stl.IsStandardImport(pkgPath) {
			continue
		}
		fmtImportLine(pkgPath)
	}

	buff.WriteString(")\n")

	// body
	body := c.buff.Bytes()
	buff.Write(body)
	return strings.TrimSpace(buff.String())
}

func (c GoFileBuffer) finalize() error {
	if c.formatted != nil {
		return nil
	}

	code := c.composit()

	if len(code) == 0 {
		return errors.NewAnonymous("生成的代码为空")
	}
	log.Printf("生成了%d个字节", len([]byte(code)))

	src, err := format.Source([]byte(code))
	if err != nil {
		log.Println("代码格式化错误")
		return errors.Extend(err, "代码格式化错误")
	}

	c.formatted = src
	log.Println("格式化结果成功")

	c.Close()
	return nil
}

// Bytes 返回最终生成的Go文件内容的字节切片
//
// 生成的代码已经格式化，因此语法上应该是合法的
func (c GoFileBuffer) Bytes() []byte {
	err := c.finalize()
	if err != nil {
		if !c.checked {
			c.CheckSyntax()
			c.checked = false
			panic(err)
		}
	}
	return c.formatted
}

func (c GoFileBuffer) String() string {
	return string(c.Bytes())
}

func (c GoFileBuffer) invalidateCache() {
	c.formatted = nil
}

func markLine(content string, pos token.Position) string {
	sb := strings.Builder{}
	lineNum := 1
	for line := range strings.Lines(content) {
		if lineNum == pos.Line {
			c := CSI.Fore | CSI.Red
			sb.WriteString(c.String())
			sb.WriteString(line)
			sb.WriteString(c.ToReset().String())
			if pos.Column >= 0 && pos.Column <= len(line) {
				sb.WriteString(CSI.Hide.String())
				sb.WriteString(line[0 : pos.Column-1])
				sb.WriteString(CSI.ResetAll.String())

				sb.WriteString(c.String())
				sb.WriteString("^\n\n")
				sb.WriteString(c.ToReset().String())
			} else {
				fmt.Fprintf(&sb, "error at col %d (out of range)\n", pos.Column)
			}
		} else {
			c := CSI.Dim
			sb.WriteString(c.String())
			sb.WriteString(line)
			sb.WriteString(c.ToReset().String())
		}
		lineNum++
	}
	return sb.String()
}
