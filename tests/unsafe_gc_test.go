package tests

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/require"
)

// 无法优化字符串，确保测试时字符串一定在堆上
var someString string = " "

func returnUnsafeHeap() []byte {
	var stringVal string
	stringVal = "the quick brown fox" + someString + "jumps over the lazy dog"
	noop(stringVal)
	return unsafe.Slice(unsafe.StringData(stringVal), len(stringVal))
}

func noop(s string) {}

// go对unsafeptr的垃圾回收行为进行测试
//
// 验证在调用runtime.GC()后，unsafe.Slice获取的字节切片仍然保持有效
func Test_UnsafeStringSliceAlive(t *testing.T) {
	myenv.T(t)

	val := returnUnsafeHeap()

	require.Equal(t, []byte("the quick brown fox jumps over the lazy dog"), val)

	runtime.GC()

	require.Equal(t, []byte("the quick brown fox jumps over the lazy dog"), val)
}

// 验证通过unsafe.Slice修改字符串内容后，字符串的值是否发生变化
func Test_UnsafeModifyString(t *testing.T) {
	myenv.T(t)

	var str string
	str = "the quick brown fox" + someString + "jumps over the lazy dog"

	noop(str)

	val := unsafe.Slice(unsafe.StringData(str), len(str))

	val[0] = 'T'

	require.Equal(t, "The quick brown fox jumps over the lazy dog", str)
}
