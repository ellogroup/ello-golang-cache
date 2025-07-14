package cachefunc

import (
	"github.com/ellogroup/ello-golang-clock/clock"
	"time"
)

const (
	NoCache  = time.Duration(0)
	NoExpire = time.Duration(-1)
)

var sysClock = clock.NewSystem()

// CacheFunc allows the result of a function to be cached, and will only call the function when the cached item does not
// exist or needs to be refreshed.
type CacheFunc[K comparable, V any] interface {
	Get(k K, fn func() (V, error), ttl time.Duration) (V, error)
}

// New returns a new instance of CacheFunc. By default, the Memory implementation is used.
func New[K comparable, V any]() CacheFunc[K, V] {
	return NewMemory[K, V]()
}

// Item represents a cached item with an expiry time
type Item[V any] struct {
	val V
	exp time.Time
}

// NewItem returns a new Item
func NewItem[V any](val V) *Item[V] {
	return &Item[V]{val: val}
}

// NewItemWithTTL returns a new Item with a expiry time calculated from ttl
func NewItemWithTTL[V any](val V, ttl time.Duration) *Item[V] {
	if ttl == NoExpire {
		return &Item[V]{val: val}
	}
	return &Item[V]{val: val, exp: sysClock.Now().Add(ttl)}
}

// Expired returns whether the Item is currently expired or not
func (i *Item[V]) Expired() bool {
	return !i.exp.IsZero() && !i.exp.After(sysClock.Now())
}
