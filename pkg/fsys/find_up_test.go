package fsys

import (
	"iter"
	"slices"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/gongt/go/pkg/fsys/fpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mapfrom(paths []*fpath.IPath) []string {
	var result []string
	for _, p := range paths {
		result = append(result, p.Raw())
	}
	return result
}

func TestClimbingPath(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	check := func(expected []string, actual iter.Seq[*fpath.IPath]) {
		val := slices.Collect(actual)
		assert.Equal(t, len(expected), len(val))
		assert.Equal(t, expected, mapfrom(val))
	}

	check([]string{"/a/b/c/d", "/a/b/c", "/a/b", "/a", "/"}, ClimbingPath("/a/b/c/d"))

	check([]string{"/x/.././y", "/x/../.", "/x/..", "/x", "/"}, ClimbingPath("/x/.././y"))

	check([]string{"a/b/c/d", "a/b/c", "a/b", "a", "."}, ClimbingPath("a/b/c/d"))
}

func TestFindUp(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	me, _ := myenv.CurrentFileLine()

	r, ok := FindUpUntilEntry(me, "go.mod", "go.sum")
	require.True(t, ok)
	assert.Equal(t, "go.mod", r.Base().Name)
}
