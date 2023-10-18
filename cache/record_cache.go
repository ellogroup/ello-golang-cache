package cache

import (
	"context"
	"fmt"
	"github.com/ellogroup/ello-golang-cache/driver"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

// RecordCache for a detailed explanation of the below, RecordCacheItem, KeylessRecordCache and driver.Cache
// please see https://ellogroup.atlassian.net/wiki/spaces/EP/pages/12648450/Cache+Package
type RecordCache[K comparable, V any] struct {
	log             *zap.Logger
	onDemandFetcher OnDemandFetcher[K, V]
	asyncFetcher    AsyncFetcher[K, V]
	cache           driver.Cache[K, RecordCacheItem[V]]
	recordTtl       time.Duration
	allTtl          time.Duration
	lastUpdated     time.Time
	cron            *cron.Cron
}

func NewRecordCache[K comparable, V any](cache driver.Cache[K, RecordCacheItem[V]]) *RecordCache[K, V] {
	return &RecordCache[K, V]{
		cache: cache,
		log:   zap.NewNop(),
	}
}

func (r *RecordCache[K, V]) AddLogger(l *zap.Logger) *RecordCache[K, V] {
	r.log = l
	return r
}

func (r *RecordCache[K, V]) SetOnDemandFetcher(f OnDemandFetcher[K, V], ttl time.Duration) *RecordCache[K, V] {
	r.onDemandFetcher = f
	return r.refreshStaleRecordsEvery(ttl)
}

func (r *RecordCache[K, V]) SetAsyncFetcher(f AsyncFetcher[K, V], ttl time.Duration) *RecordCache[K, V] {
	r.asyncFetcher = f
	return r.refreshAllRecordsEvery(ttl)
}

func (r *RecordCache[K, V]) refreshStaleRecordsEvery(ttl time.Duration) *RecordCache[K, V] {
	r.recordTtl = ttl
	err := r.setSchedule()
	if err != nil {
		r.log.Error("Could not start scheduler", zap.Error(err))
	}
	r.log.Debug("Refresh Stale initialized", zap.String("Ttl", ttl.String()))
	return r
}

func (r *RecordCache[K, V]) refreshAllRecordsEvery(ttl time.Duration) *RecordCache[K, V] {
	r.allTtl = ttl
	err := r.setSchedule()
	if err != nil {
		r.log.Error("Could not start scheduler", zap.Error(err))
	}
	r.log.Debug("Refresh all set", zap.String("Ttl", ttl.String()))
	return r
}

func (r *RecordCache[K, V]) needRefreshing(k K) bool {
	v, ok := r.cache.Get(context.Background(), k)
	if r.onDemandFetcher == nil {
		return !ok || v.IsStale(r.allTtl)
	}
	return !ok || v.IsStale(r.recordTtl)
}

func (r *RecordCache[K, V]) Get(ctx context.Context, k K) (V, error) {
	if r.needRefreshing(k) {
		if err := r.refreshItem(k); err != nil {
			return *new(V), err
		}
	}
	record, ok := r.cache.Get(ctx, k)
	if !ok {
		return *new(V), fmt.Errorf("record not in cache")
	}
	return record.V, nil
}

func (r *RecordCache[K, V]) refreshAllRecords() {
	r.log.Info("Refreshing all records")
	if r.asyncFetcher == nil {
		return
	}
	latest, err := r.asyncFetcher.FetchAll(context.Background())
	if err != nil {
		return
	}
	if !r.cache.Clear(context.Background()) {
		r.log.Warn("could not empty cache when refreshing all")
	}
	now := time.Now()
	for k, v := range latest {
		r.cache.Set(context.Background(), k, RecordCacheItem[V]{V: v, T: now})
	}
	r.log.Info("Cache refreshed")
}

func (r *RecordCache[K, V]) removeStale() {
	for k, v := range r.cache.All(context.Background()) {
		if v.IsStale(r.recordTtl) {
			r.cache.Delete(context.Background(), k)
		}
	}
}

func (r *RecordCache[K, V]) refreshCache() {
	if r.onDemandFetcher != nil {
		r.removeStale()
	}
	if r.asyncFetcher != nil &&
		(r.lastUpdated.IsZero() || r.lastUpdated.Before(time.Now().Add(-1*r.allTtl))) {
		r.refreshAllRecords()
		r.lastUpdated = time.Now()
	}
}

func (r *RecordCache[K, V]) refreshItem(k K) error {
	if r.onDemandFetcher == nil {
		return fmt.Errorf("value not in cache and on demand fetcher is not initalised")
	}
	r.log.Info("Refreshing record", zap.Any("key", k))
	v, err := r.onDemandFetcher.FetchByKey(context.Background(), k)
	if err != nil {
		return err
	}
	r.cache.Set(context.Background(), k, RecordCacheItem[V]{V: v, T: time.Now()})
	r.log.Debug("Refreshing record", zap.Any("Key", k), zap.Any("Value", v))
	return nil
}

func (r *RecordCache[K, V]) setSchedule() error {
	if r.cron != nil {
		r.log.Debug("Cron already set")
		return nil
	}
	r.refreshCache()
	r.cron = cron.New()
	if _, err := r.cron.AddFunc("* * * * *", func() { r.refreshCache() }); err != nil {
		return err
	}
	r.cron.Start()
	return nil
}
