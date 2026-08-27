package slices

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/assert"
)

func Test_Intersect(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	a := []int{1, 3, 2, 3, 4, 5, 6}
	b := []int{1, 2, 3, 6}
	IntersectInplace(&a, b)

	assert.Equal(t, []int{1, 3, 2, 6}, a)
}
