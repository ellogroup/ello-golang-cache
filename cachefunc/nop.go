package cachefunc

import (
	"time"
)

type Nop[K comparable, V any] struct{}

func NewNop[K comparable, V any]() CacheFunc[K, V] {
	return &Nop[K, V]{}
}

func (n Nop[K, V]) Get(_ K, fn func() (V, error), _ time.Duration) (V, error) {
	return fn()
}
