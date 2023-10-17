package driver

import "context"

type CacheDriver[K comparable, V any] interface {
	Has(ctx context.Context, key K) bool
	Get(ctx context.Context, k K) (V, bool)
	All(ctx context.Context) map[K]V
	Set(ctx context.Context, key K, value V) bool
	Delete(ctx context.Context, key K) bool
	Clear(ctx context.Context) bool
}
