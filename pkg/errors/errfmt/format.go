package errfmt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/errors/internal"
	"github.com/gongt/go/pkg/errors/internal/iterator"
	"github.com/gongt/go/pkg/errors/stacktrace"
	"github.com/gongt/go/pkg/string_helper/color_builder"
	CSI "github.com/gongt/go/pkg/string_helper/csi"
)

type SB = *color_builder.ColorBuilder

func FormatError(e any, color bool) string {
	sb := color_builder.New(color)
	switch err := e.(type) {
	case *panicBox:
		formatError(sb, err.Unwrap())
	case error:
		formatError(sb, err)
	default:
		formatError(sb, fmt.Errorf("FormatError输入了非error的类型: %v", err))
	}
	return sb.String()
}

func formatError(sb SB, err error) {
	c := CSI.Swap | CSI.Fore | CSI.Red
	defer sb.Padding("%s %s", c, c.ToReset())()
	var level uint
	for err := range iterator.IterReasonChain(err) {
		formatErrorOne(sb, err, level)
		level++
	}
}

func check_join(sb SB, err error) {
	if _, ok := err.(internal.UnJoin); ok {
		sb.Write("%s<合并错误>%s ", CSI.Dim, CSI.Dim.ToReset())
	}
}

func formatErrorOne(sb SB, err error, level uint) {
	if level == 0 {
		c := CSI.Swap | CSI.Fore | CSI.Red
		sb.Write("%s 错误 %s ", c|CSI.Bold, c.ToReset())

		check_join(sb, err)

		sb.WriteLine("%s%s", err.Error(), CSI.Bold.ToReset())

		if _, ok := err.(detailer); ok {
			formatDetails(sb, err, level) // details 会遍历整个树，所以只需要在开头输出一次
		}
	} else {
		c := CSI.Swap | CSI.Fore | CSI.Yellow | CSI.Bold
		sb.Write("%s 以上错误发生于处理另一个错误时 %s ", c, c.ToReset())

		check_join(sb, err)

		// r := sb.Padding("%s %s ", c, c.ToReset())
		sb.WriteRawLine(err.Error())
		// r()
	}

	if strace := errors.GetStackTrace(err); len(strace) > 0 {
		formatStack(sb, strace, level)
	} else {
		sb.WriteLine("  - %s缺少栈信息%s", CSI.Fore|CSI.Yellow, CSI.Fore.ToReset())
		formatStack(sb, stacktrace.CaptureStackTrace(0, 0), level)
	}
	sb.NewLine()
}

type detailer interface {
	Details() map[string]any
}

func formatDetails(sb SB, err error, _ uint) {
	details := make(map[string]any)
	for detail := range iterator.IterEveryDetail(err) {
		for k, v := range detail {
			if k == "reason" { // reason字段是原因链
				continue
			}
			if _, exists := details[k]; !exists {
				details[k] = v
			}
		}
	}

	if len(details) == 0 {
		return
	}

	sb.WriteLine("» %s附加数据%s", CSI.Fore|CSI.Blue, CSI.Fore.ToReset())
	for k, v := range details {
		sb.WriteLine(" ◈ %s: %s%v%s", k, CSI.Dim, v, CSI.Dim.ToReset())
	}

}

func formatStack(sb *color_builder.ColorBuilder, stack stacktrace.StackTraceArray, _ uint) {
	size := uint(len(stack))
	num_w := len(fmt.Sprintf("%d", size-1))

	for frame, index := range stacktrace.IterateStack(stack) {
		is_stl := false
		sb.Write(" %*d: ", num_w, size-index-1)

		if frame.Function == "" {
			sb.WriteWrapF(CSI.Italic|CSI.Dim, "<未知函数 +%x>", frame.PC)
		} else {
			var funcName, anon string

			parts := strings.Split(frame.Function, ".") // Split返回不可能为空，至少有一个元素
			last, parts, _ := pop(parts)

			if matched, _ := regexp.MatchString(`^func\d+$`, *last); matched { // 匿名函数
				anon = *last
				last, parts, _ = pop(parts)
				if last == nil {
					sb.WriteWrap(CSI.Fg(CSI.Red), frame.Function)
					return
				} else {
					funcName = *last
				}
			} else {
				funcName = *last
			}

			pkg := strings.Join(parts, ".")
			is_stl = is_standard_library(pkg)
			if is_stl {
				sb.WriteWrap(CSI.Dim|CSI.Italic, frame.Function)
			} else {
				sb.WriteWrap(CSI.Dim, pkg+".")
				sb.WriteWrap(CSI.Green|CSI.Fore|CSI.Bold, funcName)
				if anon != "" {
					sb.WriteRaw(".")
					sb.WriteWrap(CSI.Dim|CSI.Italic, anon)
				}
			}
		}

		sb.NewLine()

		if frame.File != "" {
			if is_stl {
				sb.WriteWrapF(CSI.Dim, "        %s:%d\n", frame.File, frame.Line)
			} else {
				sb.WriteLine("        %s:%d", frame.File, frame.Line)
			}
		}
	}
}

func pop[T any](s []T) (*T, []T, bool) {
	if len(s) == 0 {
		return nil, nil, false
	}
	return &s[len(s)-1], s[:len(s)-1], true
}

func is_standard_library(pkg string) bool {
	// 标准库包名不包含点号 (?)
	if !strings.Contains(pkg, ".") {
		if pkg == "main" {
			return false
		}
		return true
	}
	return false
}
