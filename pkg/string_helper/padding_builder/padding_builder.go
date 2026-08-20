package padding_builder

import (
	"fmt"
	"strings"
)

type PaddingBuilder struct {
	sb strings.Builder

	indentLevel int

	paddingString  string
	cursorReturned bool

	fmt func(s string, args ...any) string
}

func New() *PaddingBuilder {
	return &PaddingBuilder{
		fmt:            fmt.Sprintf,
		cursorReturned: true,
	}
}

/* 标准 */

func (h *PaddingBuilder) String() string {
	return h.sb.String()
}

func (h *PaddingBuilder) Reset() {
	h.sb.Reset()
	h.indentLevel = 0
	h.paddingString = ""
	h.cursorReturned = true
}

func (h *PaddingBuilder) Len() int {
	return h.sb.Len()
}

func (h *PaddingBuilder) Grow(n int) {
	h.sb.Grow(n)
}

/* 写入 */

func (h *PaddingBuilder) Write(s string, args ...any) {
	h.indent_write(h.fmt(s, args...))
}

func (h *PaddingBuilder) WriteLine(s string, args ...any) {
	h.Write(s, args...)
	h.Write("\n")
}

func (h *PaddingBuilder) WriteRaw(s string) {
	h.indent_write(s)
}

func (h *PaddingBuilder) WriteRawLine(s string) {
	h.WriteRaw(s)
	h.WriteRaw("\n")
}

func (h *PaddingBuilder) NewLine() {
	h.WriteRaw("\n")
}

/* 缩进 */

func (h *PaddingBuilder) Indent() func() {
	return h.Padding(strings.Repeat(" ", 4))
}

func (h *PaddingBuilder) Padding(s string, args ...any) func() {
	s = h.fmt(s, args...)

	selfIndex := len(h.paddingString)
	h.paddingString += s

	return func() {
		if selfIndex != len(h.paddingString)-len(s) {
			panic("Padding() 调用顺序错误，必须按顺序依次调用")
		}
		h.paddingString = h.paddingString[:selfIndex]
	}
}

/* 扩展 */
func (h *PaddingBuilder) SetFormatFunc(f func(s string, args ...any) string) {
	h.fmt = f
}

/* 私有方法 */

func (h *PaddingBuilder) indent_write(s string) {
	if s == "" {
		return
	}
	if h.paddingString == "" {
		h.sb.WriteString(s)
		h.cursorReturned = nl_test(s)
		return
	}

	line_no := 0
	for line := range strings.Lines(s) { // Lines的每个元素都带有换行符
		if h.cursorReturned {
			h.sb.WriteString(h.paddingString)
		}
		h.sb.WriteString(line)
		line_no++
		h.cursorReturned = nl_test(line)
	}
}

func nl_test(s string) bool {
	return strings.HasSuffix(s, "\n") || strings.HasSuffix(s, "\r")
}
