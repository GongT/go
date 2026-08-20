package color_builder

import (
	"fmt"
	"strings"

	CSI "github.com/gongt/go/pkg/string_helper/csi"
	"github.com/gongt/go/pkg/string_helper/padding_builder"
)

// 是一个支持切换颜色开关的字符串构建器，用于命令行输出
//   - 一般和CSI包配合使用
type ColorBuilder struct {
	padding_builder.PaddingBuilder

	colorDisabled bool
}

func New(enable bool) *ColorBuilder {
	r := &ColorBuilder{
		PaddingBuilder: *padding_builder.New(),
		colorDisabled:  !enable,
	}

	r.SetColorEnabled(enable)

	return r
}

func (b *ColorBuilder) SetColorEnabled(enable bool) {
	if b.colorDisabled == !enable {
		return
	}

	b.colorDisabled = !enable
	if enable {
		b.SetFormatFunc(fmt.Sprintf)
	} else {
		b.SetFormatFunc(b.monochromeFormat)
	}
}

func (b *ColorBuilder) IsColorEnabled() bool {
	return !b.colorDisabled
}

func (b *ColorBuilder) WriteWrap(c CSI.St, s string) {
	b.Write("%s%s%s", c, s, c.ToReset())
}
func (b *ColorBuilder) WriteWrapF(c CSI.St, s string, args ...any) {
	b.Write("%s", c)
	b.Write(s, args...)
	b.Write("%s", c.ToReset())
}

func (b *ColorBuilder) monochromeFormat(s string, args ...any) string {
	for i, arg := range args {
		if str, ok := arg.(string); ok {
			if strings.HasPrefix(str, "\x1b[") && strings.HasSuffix(str, "m") {
				args[i] = ""
			}
		} else if _, ok := arg.(CSI.Seq); ok {
			args[i] = ""
		} else if _, ok := arg.(CSI.St); ok {
			args[i] = ""
		}
	}
	return fmt.Sprintf(s, args...)
}
