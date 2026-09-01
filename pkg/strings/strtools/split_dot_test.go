package strtools

import (
	"slices"
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/assert"
)

func TestSplitDot(t *testing.T) {
	myenv.T(t)

	str := "a.b..c...d.e"

	r := slices.Collect(SplitSingle(str, '.'))

	assert.Equal(t, []string{"a", "b..c...d", "e"}, r)
}
