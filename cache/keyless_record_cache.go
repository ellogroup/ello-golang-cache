package cache

import (
	"context"
	"github.com/ellogroup/ello-golang-cache/driver"
	"go.uber.org/zap"
	"time"
)

type onDemandFetcher[V any] struct {
	f KeylessFetcher[V]
}

func (o *onDemandFetcher[V]) FetchByKey(ctx context.Context, _ int) (V, error) {
	return o.f.Fetch(ctx)
}

func newOnDemandFetcher[V any](f KeylessFetcher[V]) *onDemandFetcher[V] {
	return &onDemandFetcher[V]{f: f}
}

type asyncFetcher[V any] struct {
	f KeylessFetcher[V]
}

func (a asyncFetcher[V]) FetchAll(ctx context.Context) (map[int]V, error) {
	v, err := a.f.Fetch(ctx)
	return map[int]V{0: v}, err
}

func newAsyncFetcher[V any](f KeylessFetcher[V]) asyncFetcher[V] {
	return asyncFetcher[V]{f: f}
}

// KeylessRecordCache to allow rotation of a cache with only one value
// a prime use case being a token cache for a service like salesforce
type KeylessRecordCache[V any] struct {
	*RecordCache[int, V]
}

func NewKeylessRecordCacheOnDemand[V any](driver driver.Cache[int, RecordCacheItem[V]], f KeylessFetcher[V], ttl time.Duration) *KeylessRecordCache[V] {
	return &KeylessRecordCache[V]{NewRecordCache[int, V](driver).SetOnDemandFetcher(newOnDemandFetcher(f), ttl)}
}

func NewKeylessRecordCacheAsync[V any](driver driver.Cache[int, RecordCacheItem[V]], f KeylessFetcher[V], ttl time.Duration) *KeylessRecordCache[V] {
	return &KeylessRecordCache[V]{
		NewRecordCache[int, V](driver).
			SetAsyncFetcher(newAsyncFetcher(f), ttl),
	}
}

func NewKeylessRecordCacheAsyncWithLogger[V any](driver driver.Cache[int, RecordCacheItem[V]], f KeylessFetcher[V], ttl time.Duration, log *zap.Logger) *KeylessRecordCache[V] {
	return &KeylessRecordCache[V]{
		NewRecordCache[int, V](driver).
			AddLogger(log).
			SetAsyncFetcher(newAsyncFetcher(f), ttl),
	}
}

func (k *KeylessRecordCache[V]) Get(ctx context.Context) (V, error) {
	return k.RecordCache.Get(ctx, 0)
}
