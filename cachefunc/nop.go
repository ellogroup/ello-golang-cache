package cachefunc

import (
	"time"
)

// Nop is a a no-operation implementation of CacheFunc
type Nop[K comparable, V any] struct{}

// NewNop returns a new instance of Nop
func NewNop[K comparable, V any]() CacheFunc[K, V] {
	return &Nop[K, V]{}
}

func (n Nop[K, V]) Get(_ K, fn func() (V, error), _ time.Duration) (V, error) {
	return fn()
}
