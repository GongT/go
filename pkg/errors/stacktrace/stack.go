package stacktrace

import (
	"iter"
	"runtime"
	_ "unsafe"
)

type StackTraceArray = []uintptr

const MAX_STACK_DEPTH = 32

// 抓取当前调用栈
//
// skip: 跳过的帧数
//
// cap: 最大深度，如果<=0，则使用默认值 [MAX_STACK_DEPTH]
func CaptureStackTrace(skip int, cap int) StackTraceArray {
	if cap <= 0 {
		cap = MAX_STACK_DEPTH
	}
	pcs := make(StackTraceArray, cap)
	n := runtime.Callers(skip+2, pcs)
	return pcs[:n]
}

type Iter = iter.Seq2[runtime.Frame, uint]

// 迭代每一帧
//
// 根据文档，迭代过程中禁止修改stack（虽然这种情况几乎不可能出现）
func IterateStack(stack StackTraceArray) Iter {
	return IterateFrames(runtime.CallersFrames(stack))
}

// 把runtime.CallersFrames的返回值转换成迭代器，第二个值为帧索引
func IterateFrames(frames *runtime.Frames) Iter {
	return func(yield func(runtime.Frame, uint) bool) {
		var index uint
		for {
			frame, more := frames.Next()
			if frame.PC == 0 {
				// frames.Next()返回的frame.PC为0表示没有更多帧
				// 正常来讲应该不可能发生
				return
			}
			if !yield(frame, index) || !more {
				return
			}
			index++
		}
	}
}

// var panicFunctionPC uintptr
//
// func loadPanicPc() {
// 	recover()
// 	pc, _, _, ok := runtime.Caller(1)
// 	if ok {
// 		if myenv.IsDebug {
// 			myenv.Assert(runtime.FuncForPC(pc).Name() == "runtime.gopanic", "栈状态错误")
// 		}
// 		fmt.Printf("Detected gopanic PC is %x\n", pc)
// 		panicFunctionPC = pc
// 	}
// 	if panicFunctionPC == 0 {
// 		println("Failed to detect gopanic PC, returning 0")
// 	}
// }

// func init() {
// 	defer loadPanicPc()
// 	panic("trigger")
// }

func DetectPanicPath(stack StackTraceArray) bool {
	// if panicFunctionPC > 0 && slices.Contains(stack, panicFunctionPC+1) {
	// gopanic+1 大概就是stack中的值，但完全无法保证
	// return true
	// } else {
	for frame := range IterateStack(stack) {
		if frame.Func != nil && frame.Func.Name() == "runtime.gopanic" {
			println("Detected panic path in stack trace", frame.PC)
			return true
		}
	}
	return false
	// }
}
