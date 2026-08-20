package iterator_test

import (
	"errors"
	"iter"
	"testing"

	"github.com/gongt/go/internal/myenv"
	. "github.com/gongt/go/pkg/iterator"
	"github.com/stretchr/testify/assert"
)

func TestIterateBreak(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	var result []int
	for i := range iterateOverValue() {
		result = append(result, i)
	}

	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, result)
}

func TestErrorCatcher(t *testing.T) {
	myenv.RedirectDebugTesting(t)
	
	var result []int
	handler, holder := CreateErrorStore(nil)


func iterateOverValue() iter.Seq[int] {
	return CreateIterator(func(yield Yield[int]) {
		// Example iteration logic
		for i := range 10 {
			if !yield(i, nil) {
				break
			}
			if i == 5 {
				// Simulate an error condition
				if !yield(0, errors.New("An error occurred at i=5")) {
					break
				}
			}
		}
	}, nil)
}
