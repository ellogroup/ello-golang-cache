package cachefunc

import (
	"sync"
	"time"
)

type Memory[K comparable, V any] struct {
	m     map[K]*Item[V]
	mutex sync.RWMutex
}

func NewMemory[K comparable, V any]() CacheFunc[K, V] {
	c := &Memory[K, V]{
		m: make(map[K]*Item[V]),
	}
	return c
}

func (c *Memory[K, V]) Get(k K, fn func() (V, error), exp time.Duration) (V, error) {
	if i, ok := c.read(k); ok && !i.Expired() {
		// found in cache
		return i.val, nil
	}
	// fetch value
	v, err := fn()
	if err == nil {
		// only cache if no error
		c.write(k, NewItemExpiresIn(v, exp))
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
