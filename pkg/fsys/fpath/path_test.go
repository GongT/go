package fpath

import (
	"os"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	t.Run("相对路径", func(t *testing.T) {
		p := New("a/.b/c")

		assert.Equal(t, false, p.IsAbs())
		assert.Equal(t, false, p.IsRoot())

		assert.Equal(t, "a/.b/c", p.String())

		p.Resolve("../d")
		assert.Equal(t, "a/.b/d", p.String())

		p.Join("x", "...", "z.test.txt")
		assert.Equal(t, "a/.b/d/x/.../z.test.txt", p.String())

		assert.Equal(t, "z.test.txt", p.Base().Name)
		assert.Equal(t, "z.test", p.Base().Stem())
		assert.Equal(t, ".txt", p.Base().Ext())

		p = New("a/..")
		assert.Equal(t, false, p.IsRoot())
		assert.Equal(t, ".", p.String())
		p.Join("../..")
		assert.Equal(t, "../..", p.String())
	})

	t.Run("绝对路径", func(t *testing.T) {
		pp := New("a/.b/!d/x/.../z.test.txt")
		p := New("/root")

		assert.Equal(t, true, p.IsAbs())

		p.ResolveWith(pp)
		assert.Equal(t, "/root/a/.b/!d/x/.../z.test.txt", p.String())

		p.Resolve("/")
		assert.Equal(t, "/", p.String()) // TODO: windows?

		p.Resolve("..")
		assert.Equal(t, "/", p.String())

		p.Resolve("a/b/c")
		assert.Equal(t, "/a/b/c", p.String())

		p.SetFilename("new.txt")
		assert.Equal(t, "/a/b/new.txt", p.String())

		r := New("/")
		assert.Equal(t, true, r.IsRoot())
		r.Dir()
		assert.Equal(t, "/", r.String())
	})

	t.Run("其他测试", func(t *testing.T) {
		p := New("/a/../../..")
		assert.Equal(t, true, p.IsRoot())

		assert.True(t, New("/a/b/./c").NeedsNormalize())
		assert.False(t, New("/a/b/c").NeedsNormalize())

		p = New("../../a")
		p.SetBase("b")
		assert.Equal(t, "../../b", p.String())
		p.SetDir("/x/y")

		assert.Equal(t, "/x/y/b", p.String())
	})

	t.Run("读写测试", func(t *testing.T) {
		tmp := t.TempDir()
		p := New(tmp)
		p.Join("a", "..", "b", "..", "c.txt")
		assert.Equal(t, tmp+"/a/../b/../c.txt", p.Raw())

		os.WriteFile(p.Raw(), []byte("aaaaaaa"), 0644)
	})
}
