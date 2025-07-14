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

type CacheFunc[K comparable, V any] interface {
	Get(k K, fn func() (V, error), exp time.Duration) (V, error)
}

func New[K comparable, V any]() CacheFunc[K, V] {
	return NewMemory[K, V]()
}

type Item[V any] struct {
	val V
	exp time.Time
}

func NewItem[V any](val V) *Item[V] {
	return &Item[V]{val: val}
}

func NewItemExpiresIn[V any](val V, d time.Duration) *Item[V] {
	if d == NoExpire {
		return &Item[V]{val: val}
	}
	return &Item[V]{val: val, exp: sysClock.Now().Add(d)}
}

func (i *Item[V]) Expired() bool {
	return !i.exp.IsZero() && !i.exp.After(sysClock.Now())
}
