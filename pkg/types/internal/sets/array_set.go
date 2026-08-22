// @exported
package sets

import (
	"slices"
	"sync"
)

// Set 集合的数组简易实现，适合只有几个元素的情况
type Set[T comparable] []T

func (s Set[T]) Has(item T) bool {
	return slices.Contains(s, item)
}

func (s *Set[T]) Add(item T) {
	if !s.Has(item) {
		*s = append(*s, item)
	}
}

func (s *Set[T]) Delete(item T) {
	for i, v := range *s {
		if v == item {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}

// SharedSet 线程安全的 [Set]
type SharedSet[T comparable] struct {
	mu  sync.RWMutex
	set Set[T]
}

func (s *SharedSet[T]) Has(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set.Has(item)
}

func (s *SharedSet[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set.Add(item)
}

func (s *SharedSet[T]) Delete(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set.Delete(item)
}
