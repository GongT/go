package strtools

import (
	"iter"
)

// SplitSingle 用ch分割字符串，类似 [strings.Split](s, ch)，但多个ch连在一起时，视为不是分隔符
func SplitSingle(s string, ch byte) iter.Seq[string] {
	return func(yield func(string) bool) {
		var buf []byte
		i := 0
		for i < len(s) {
			if s[i] != ch {
				buf = append(buf, s[i])
				i++
				continue
			}

			// 统计连续的ch的个数
			j := i
			for j < len(s) && s[j] == ch {
				j++
			}
			dotCount := j - i

			if dotCount == 1 {
				// 单独的ch，视为分隔符
				if !yield(string(buf)) {
					return
				}
				buf = buf[:0]
			} else {
				// 连续多个ch，不是分隔符，原样保留
				buf = append(buf, s[i:j]...)
			}
			i = j
		}
		yield(string(buf))
	}
}
