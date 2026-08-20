package types

import (
	"iter"
	"slices"
	"sync"
)

// Collection是一个任意类型指针的可重复的集合
type Collection[T comparable] []T

type RemoveFunc func()

func (c *Collection[T]) Add(item T) RemoveFunc {
	deleted := false
	*c = append(*c, item)

	return func() {
		if deleted {
			panic("Collection::remove重复调用")
		}
		deleted = true

		for i, v := range *c {
			if v == item {
				*c = append((*c)[:i], (*c)[i+1:]...)
				return
			}
		}
		panic("Collection::remove未找到元素")
	}
}

func (c *Collection[T]) Len() int {
	return len(*c)
}

func (c *Collection[T]) Items() iter.Seq[T] {
	return slices.Values(*c)
}

type SharedCollection[T comparable] struct {
	c  Collection[T]
	mu sync.RWMutex
}

func (sc *SharedCollection[T]) Add(item T) RemoveFunc {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	fn := sc.c.Add(item)

	return func() {
		sc.mu.Lock()
		defer sc.mu.Unlock()

		fn()
	}
}

func (sc *SharedCollection[T]) Len() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.c.Len()
}

func (sc *SharedCollection[T]) Items() iter.Seq[T] {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	c := slices.Clone(sc.c)
	return slices.Values(c)
}

func (sc *SharedCollection[T]) Exec(f func(itr iter.Seq[T])) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	f(sc.c.Items())
}
