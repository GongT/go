package fpath

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/assert"
)

func Test_IsLocal(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	assert.True(t, IsLocal("/a/b/c/d", "/a/b"))
	assert.True(t, IsLocal("c/d", "/a/b"))
	assert.False(t, IsLocal("../c/d", "/a/b"))
	assert.False(t, IsLocal("c/../../d", "/a/b"))
}

func Test_ToRelative(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	assert.Equal(t, "./c/d", MustRelative("/a/b/c/d", "/a/b"))
	assert.Equal(t, "../..", MustRelative("/a/b", "/a/b/c/d"))
	assert.Equal(t, ".", MustRelative("/a/b", "/a/b"))
	assert.Equal(t, "../../a/b/c/d", MustRelative("/a/b/c/d", "/x/y"))
}
