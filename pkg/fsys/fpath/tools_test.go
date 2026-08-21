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
