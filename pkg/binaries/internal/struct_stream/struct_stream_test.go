package struct_stream

import (
	"encoding/binary"
	"testing"

	"github.com/gongt/go/internal/myenv"
	sharederrors "github.com/gongt/go/pkg/errors/shared"
	"github.com/stretchr/testify/require"
)

func TestChunkReaderWrite(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	stream := New()

	s := "test data"
	stream.WriteString(s)

	bt := make([]byte, 1024)
	binary.LittleEndian.PutUint32(bt, uint32(len(s)))
	copy(bt[4:], []byte(s))

	l := 4 + len(s)
	require.Equal(t, bt[:l], stream.buffer.Bytes())

	s = "hello~world~"
	stream.WriteString(s)

	binary.LittleEndian.PutUint32(bt[l:], uint32(len(s)))
	copy(bt[l+4:], []byte(s))

	l += 4 + len(s)
	require.Equal(t, bt[:l], stream.buffer.Bytes())

	// read

	s1, err := stream.ReadString()
	require.NoError(t, err)
	require.Equal(t, "test data", s1)

	_, err = stream.Read(make([]byte, 2))
	require.ErrorIs(t, err, sharederrors.ErrEntityTooLarge)

	s2, err := stream.ReadString()
	require.NoError(t, err)
	require.Equal(t, "hello~world~", s2)
}
