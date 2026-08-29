package streaming

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/require"
)

func testWrite(r *chunkReader, text string) <-chan error {
	writeDone := make(chan error, 1)
	go func() {
		_, err := r.Write([]byte(text))
		writeDone <- err
	}()
	return writeDone
}

func testClose(r *chunkReader) <-chan error {
	closeDone := make(chan error, 1)
	go func() {
		closeDone <- r.Close()
	}()
	return closeDone
}

func expect(t *testing.T, r *chunkReader, expected string) {
	chunk, err := r.Next()
	require.NoError(t, err)
	require.Equal(t, expected, string(chunk))
}

func TestChunkReaderWrite(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	// 普通
	reader := NewChunkReader([]byte("|sep|"))
	writeDone := testWrite(reader, "left|sep|right|sep|")

	expect(t, reader, "left|sep|")
	expect(t, reader, "right|sep|")
	require.NoError(t, <-writeDone)

	// 缺少最终分隔符
	reader.Write([]byte("tail"))
	closeDone := testClose(reader)
	expect(t, reader, "tail")
	require.NoError(t, <-closeDone)

	// 空白
	reader = NewChunkReader([]byte("|"))
	writeDone = testWrite(reader, "|||")
	expect(t, reader, "|")
	expect(t, reader, "|")
	expect(t, reader, "|")
	require.NoError(t, <-writeDone)
	require.NoError(t, <-testClose(reader))

	// 不输出分隔符
	reader = NewChunkReader([]byte("|sep|"))
	reader.IncludingSep = false
	writeDone = testWrite(reader, "left|sep|right|sep|")
	expect(t, reader, "left")
	expect(t, reader, "right")
	require.NoError(t, <-writeDone)
	require.NoError(t, <-testClose(reader))
}
