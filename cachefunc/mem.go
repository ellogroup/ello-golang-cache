package cachefunc

import (
	"sync"
	"time"
)

// Memory is a memory based implementation of CacheFunc. Items are stored in a map.
type Memory[K comparable, V any] struct {
	m     map[K]*Item[V]
	mutex sync.RWMutex
}

// NewMemory returns a new instance of Memory
func NewMemory[K comparable, V any]() CacheFunc[K, V] {
	c := &Memory[K, V]{
		m: make(map[K]*Item[V]),
	}
	return c
}

func (c *Memory[K, V]) Get(k K, fn func() (V, error), ttl time.Duration) (V, error) {
	if i, ok := c.read(k); ok && !i.Expired() {
		// found in cache
		return i.val, nil
	}
	// fetch value
	v, err := fn()
	if err == nil {
		// only cache if no error
		c.write(k, NewItemWithTTL(v, ttl))
	}
	// return value
	return v, err
}

func (c *Memory[K, V]) read(k K) (*Item[V], bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	v, ok := c.m[k]
	return v, ok && v != nil
}

func (c *Memory[K, V]) write(k K, v *Item[V]) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.m[k] = v
}
