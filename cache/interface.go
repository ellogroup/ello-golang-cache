package cache

import "context"

type OnDemandFetcher[K comparable, V any] interface {
	FetchByKey(ctx context.Context, k K) (V, error)
}

type AsyncFetcher[K comparable, V any] interface {
	FetchAll(ctx context.Context) (map[K]V, error)
}

type KeylessFetcher[V any] interface {
	Fetch(ctx context.Context) (V, error)
}
