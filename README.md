# Ello Go cache packages

## Record Cache

`cache.RecordCache` can be used to easily cache any items. It requires a driver, a fetcher and a ttl.

### Drivers

Drivers are used to store the records. A custom driver can be provided that implements the `driver.Cache` interface, or 
an existing driver can be used:

#### Memory Driver

Memory driver stores items in a map.

```go
// Example of RecordCache to store key value pairs in a map (key = int, val = string)
c := cache.NewRecordCache[int, string](
    driver.NewMemoryCache[int, cache.RecordCacheItem[string]]()
)
```

#### Redis Driver

Redis driver stores items in Redis. This requires an instance of the redis client provided by 
`github.com/redis/go-redis/v9`.

```go
// Example of RecordCache to store key value pairs in Redis (key = int, val = string)
key := "redis_key_name"
client := redis.NewClient()

c := cache.NewRecordCache[int, string](
    driver.NewRedisCacheDriver[int, cache.RecordCacheItem[string]](key, client)
)
```

### Fetchers

Fetchers are used by the `cache.RecordCache` to fetch the data to be cached, and must implement one of the following 
interfaces: 

#### OnDemandFetcher interface

`cache.OnDemandFetcher` fetches an individual record as and when it is required. It will be requested if the item hasn't 
been fetched before, or if the timestamp of when it was last fetched exceeds the ttl. This allows only individual items 
to be fetched and only when required, but occasionally requests will have to wait for the record to be fetched.

```go
type ExampleOnDemandFetcher struct {}

// Implement cache.OnDemandFetcher for key value pairs of int: string
func (f ExampleOnDemandFetcher) FetchByKey(ctx context.Context, k int) (string, error) {
	// return string for the provided int key
}

// Create RecordCache and set the on demand fetcher with a ttl of 1 hour (if the requested item was previous cached 
// longer than an hour ago, it will be fetched again)
c := cache.NewRecordCache[int, string](driver).SetOnDemandFetcher(&ExampleOnDemandFetcher{}, 60 * time.Minute)
```

#### AsyncFetcher interface

`cache.AsyncFetcher` fetches all possible records and is run asynchronously according to the ttl. This allows all 
records to always be available in the cache, but means _all_ records need to be fetched which isn't always possible or 
feasible.

```go
type ExampleAsyncFetcher struct {}

// Implement cache.AsyncFetcher for key value pairs of int: string
func (f ExampleAsyncFetcher) FetchAll(ctx context.Context) (map[int]string, error) {
	// return map[int]string containing all values
}

// Create RecordCache and set the async fetcher with a ttl of 1 hour (the fetcher will be called every 1 hour)
c := cache.NewRecordCache[int, string](driver).SetAsyncFetcher(&ExampleAsyncFetcher{}, 60 * time.Minute)
```

## Keyless Record Cache

Keyless Record Cache is an implementation of `cache.RecordCache` that does not require a key. This is useful for when 
only a single value needs to be cached instead of key/value pairs. The fetcher needs to implement 
`cache.KeylessFetcher`.

There are both on demand and async options:

```go
// Example for storing a string

d := driver.NewMemoryCache[int, cache.RecordCacheItem[string]]()

// On demand
c := cache.NewKeylessRecordCacheOnDemand[string](
    d,
    f, // Implementation of cache.KeylessFetcher
    60 * time.Minute,
)
// Get value
val := c.Get(ctx)

// Asynchronous
c := cache.NewKeylessRecordCacheAsync[string](
    d,
    f, // Implementation of cache.KeylessFetcher
    60 * time.Minute,
)
// Get value
val := c.Get(ctx)
```