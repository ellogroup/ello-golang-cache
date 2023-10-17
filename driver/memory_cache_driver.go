package driver

import "context"

type MemoryCache[K comparable, V any] struct {
	c map[K]V
}

func NewMemoryCache[K comparable, V any]() MemoryCache[K, V] {
	return MemoryCache[K, V]{
		c: make(map[K]V),
	}
}

func (m MemoryCache[K, V]) Has(_ context.Context, key K) bool {
	_, ok := m.c[key]
	return ok
}

func (m MemoryCache[K, V]) Get(_ context.Context, key K) (V, bool) {
	v, ok := m.c[key]
	return v, ok
}

func (m MemoryCache[K, V]) All(_ context.Context) map[K]V {
	return m.c
}

func (m MemoryCache[K, V]) Set(ctx context.Context, key K, value V) bool {
	m.c[key] = value
	return m.Has(ctx, key)
}

func (m MemoryCache[K, V]) Delete(ctx context.Context, key K) bool {
	delete(m.c, key)
	return !m.Has(ctx, key)
}

func (m MemoryCache[K, V]) Clear(_ context.Context) bool {
	m.c = make(map[K]V)
	return len(m.c) == 0
}
