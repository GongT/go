package streaming

import (
	"testing"
	"time"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/require"
)

func TestStreamBufferReadAfterWrite(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	sb := NewStreamBuffer()

	writeDone := make(chan error, 1)
	go func() {
		_, err := sb.Write([]byte("hello"))
		writeDone <- err
	}()

	select {
	case err := <-writeDone:
		require.NoError(t, err)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("写入被卡住，说明写入时没有唤醒等待的读取操作")
	}

	readDone := make(chan struct{})
	var (
		buf []byte
		n   int
		err error
	)
	go func() {
		buf = make([]byte, 5)
		n, err = sb.Read(buf)
		close(readDone)
	}()

	select {
	case <-readDone:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("读取没有在写入后恢复，说明写入信号未正确触发")
	}

	require.NoError(t, err)
	require.Equal(t, 5, n)
	require.Equal(t, "hello", string(buf[:n]))
}

func TestStreamBuffer_WaterMark(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	sb := NewStreamBuffer()
	sb.WaterMark = 10

	writeDone := make(chan error, 1)
	go func() {
		_, err := sb.Write([]byte("hellohellohello"))
		writeDone <- err
	}()

	select {
	case <-writeDone:
		t.Fatal("写入返回，水位线阻塞未生效")
	case <-time.After(80 * time.Millisecond):
		// ok
	}
}

func TestStreamBuffer_ReadLarger(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	sb := NewStreamBuffer()
	go sb.Write([]byte("hello"))

	returned := make(chan struct{})
	go func() {
		buf := make([]byte, 20)
		sb.Read(buf)
		close(returned)
	}()

	select {
	case <-returned:
		t.Fatal("读取过早返回")
	case <-time.After(80 * time.Millisecond):
		// ok
	}
}
