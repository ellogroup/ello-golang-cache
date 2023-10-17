package cache

import "time"

type RecordCacheItem[V any] struct {
	V V
	T time.Time
}

func (rci *RecordCacheItem[V]) IsStale(ttl time.Duration) bool {
	return rci.T.Before(time.Now().Add(-1 * ttl))
}
