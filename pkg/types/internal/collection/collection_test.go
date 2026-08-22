package collection

import (
	"testing"

	"github.com/gongt/go/internal/myenv"
	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	col := Collection[uint64]{}

	rm1 := col.Add(1)
	assert.Equal(t, 1, col.Len())

	rm2 := col.Add(2)
	assert.Equal(t, 2, col.Len())
	rm21 := col.Add(2)
	assert.Equal(t, 3, col.Len())

	rm3 := col.Add(3)
	assert.Equal(t, 4, col.Len())

	rm21()
	assert.Equal(t, 3, col.Len())

	for v := range col.Items() {
		assert.Contains(t, []uint64{1, 2, 3}, v)
	}

	rm2()

	for v := range col.Items() {
		assert.Contains(t, []uint64{1, 3}, v)
	}

	rm1()
	rm3()
	assert.Equal(t, 0, col.Len())
}
