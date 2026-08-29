package strtools

import (
	"unicode"
	"unicode/utf8"
)

func UcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	first, size := utf8.DecodeRuneInString(s)
	uc := unicode.ToUpper(first)
	if first == uc {
		return s
	}
	return string(uc) + s[size:]
}

func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	first, size := utf8.DecodeRuneInString(s)
	lc := unicode.ToLower(first)
	if first == lc {
		return s
	}
	return string(lc) + s[size:]
}
