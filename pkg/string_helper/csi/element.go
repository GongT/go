package CSI

import (
	"fmt"
	"math/bits"
	"strings"
)

// St 是一个位掩码，表示一个或多个 ANSI 转义码
//
// 其中低 32 位表示颜色值（其中的最高位表示颜色位数），
// 高 32 位表示样式位
type St uint64

const (
	Reset    St = 1 << 63                       // 重置其他bit
	ResetAll St = Reset | 0xffff_ffff_0000_0000 // 重置所有样式和颜色

	Bold   St = 1 << (0 + 32)     // 粗体、高亮
	Dim    St = 1 << (1 + 32)     // 暗淡
	Italic St = 1 << (2 + 32)     // 斜体
	Under  St = 1 << (3 + 32)     // 下划线
	Under2 St = Under | 1<<(4+32) // 双线
	Under3 St = Under | 1<<(5+32) // 波浪线
	Under4 St = Under | 1<<(6+32) // 点线
	Under5 St = Under | 1<<(7+32) // 长点线
	Blink  St = 1 << (8 + 32)     // 闪烁
	Swap   St = 1 << (9 + 32)     // 交换前景、背景色
	Hide   St = 1 << (10 + 32)    // 隐藏
	Strike St = 1 << (11 + 32)    // 删除线
	Over   St = 1 << (12 + 32)    // 上划线
	Fore   St = 1 << (13 + 32)    // 低位是前景色
	Back   St = 1 << (14 + 32)    // 低位是背景色
)

// 设置如何解释低32位的颜色值，不支持古董8、16色
const (
	Color8  St = 0       // 8-bit 索引颜色
	Color24 St = 1 << 31 // 24-bit 真彩色
)

const COLOR_MASK St = 0x7fffffff            // 低31位1
const STYLE_MASK St = 0x7fff_ffff_0000_0000 // 高32位除了最高位都是1

// Fore + c
func Fg(c St) St {
	return Fore | (c & COLOR_MASK)
}

// Back + c
func Bg(c St) St {
	return Back | (c & COLOR_MASK)
}

// 真彩色值
func TC(r, g, b uint8) St {
	return Color24 + St(uint64(r)<<16+uint64(g)<<8+uint64(b))
}

// 完整字符串，包含CSR
func (e St) String() string {
	s := e.Sequence()
	if s == "" {
		return ""
	}
	return "\x1b[" + s + "m"
}

// 不包含CSR的字符串表示
func (e St) Sequence() string {
	if e == 0 {
		return ""
	} else if e&ResetAll == ResetAll {
		return "0"
	}

	bs := strings.Builder{}
	bs.Grow(bits.OnesCount32(uint32(e)) * 3)

	if e&Reset == 0 { // set
		if e&Bold != 0 {
			bs.WriteString("1;")
		} else if e&Dim != 0 {
			bs.WriteString("2;")
		}
		if e&Italic != 0 {
			bs.WriteString("3;")
		}
		if e&Under != 0 {
			bs.WriteString("4")
			if e&Under2 != 0 {
				bs.WriteString(":2")
			} else if e&Under3 != 0 {
				bs.WriteString(":3")
			} else if e&Under4 != 0 {
				bs.WriteString(":4")
			} else if e&Under5 != 0 {
				bs.WriteString(":5")
			}
			bs.WriteString(";")
		}
		if e&Blink != 0 {
			bs.WriteString("5;")
		}
		if e&Swap != 0 {
			bs.WriteString("7;")
		}
		if e&Hide != 0 {
			bs.WriteString("8;")
		}
		if e&Strike != 0 {
			bs.WriteString("9;")
		}
		if e&Over != 0 {
			bs.WriteString("53;")
		}
		if (e&Fore != 0 || e&Back != 0) && e.HasColor() {
			if e&Fore != 0 {
				bs.WriteString("38;")
				e.EscapeColor(&bs)
			}
			if e&Back != 0 {
				bs.WriteString("48;")
				e.EscapeColor(&bs)
			}
		}

	} else { // reset
		if e&Bold != 0 || e&Dim != 0 {
			bs.WriteString("22;")
		}
		if e&Italic != 0 {
			bs.WriteString("23;")
		}
		if e&Under != 0 {
			bs.WriteString("24;")
		}
		if e&Blink != 0 {
			bs.WriteString("25;")
		}
		if e&Swap != 0 {
			bs.WriteString("27;")
		}
		if e&Hide != 0 {
			bs.WriteString("28;")
		}
		if e&Strike != 0 {
			bs.WriteString("29;")
		}
		if e&Over != 0 {
			bs.WriteString("55;")
		}
		if e&Fore != 0 {
			bs.WriteString("39;")
		}
		if e&Back != 0 {
			bs.WriteString("49;")
		}
	}
	if bs.Len() == 0 {
		return ""
	}
	return bs.String()[:bs.Len()-1]
}

func (e St) HasColor() bool {
	if e&Color24 != 0 { // TrueColor
		return e&0xffffff != 0
	} else {
		return e&0xff != 0
	}
}

// 将颜色值转换为ANSI转义码的参数部分
func (e St) EscapeColor(bs *strings.Builder) {
	if e&Color24 != 0 { // TrueColor
		fmt.Fprintf(bs, "2;%d;%d;%d;", (e>>16)&0xff, (e>>8)&0xff, e&0xff)
	} else {
		fmt.Fprintf(bs, "5;%d;", e&0xff)
	}
}

// e & ^m
func (e *St) Unset(m St) {
	*e &^= m
}

// e | m
func (e *St) Set(m St) {
	*e |= m
}

// e | Reset
func (e St) ToReset() St {
	return e | Reset
}

// e & ^(Fore | Back | 0xffffff)
//
// 移除颜色位，保留样式位
func (e St) noColorBits() St {
	return e &^ (Fore | Back | 0xffffff)
}
