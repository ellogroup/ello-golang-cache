# Ello Go cache packages

## CacheFunc

`cachefunc` allows the result of a function to be cached, and will only call the function when the cached item does not 
exist or needs to be refreshed.

### Usage

```go
package main

import (
	"errors"
	"fmt"
	"github.com/ellogroup/ello-golang-cache/v2/cachefunc"
	"time"
)

func main() {
	// Will "fetch" user 1
	user1, err := getUser(1)
	fmt.Printf("getUser(1) => %+v, %+v\n\n", user1, err)

	// Will return the now cached user 1
	user1, err = getUser(1)
	fmt.Printf("getUser(1) => %+v, %+v\n\n", user1, err)

	// Will return not found error
	user99, err := getUser(99)
	fmt.Printf("getUser(99) => %+v, %+v\n\n", user99, err)
}

var cache = cachefunc.New[int, user]()

type user struct {
	id   int
	name string
}

var users = map[int]user{
	1: {1, "user1"},
	2: {2, "user2"},
}

func getUser(id int) (user, error) {
	return cache.Get(id, func() (user, error) {
		// This function is only called when the cache needs to set/refresh this item

		fmt.Printf("fetching user %q...\n\n", id)

		if u, ok := users[id]; ok {
			// The item is only cached if a nil error returned
			return u, nil
		}
		// Any errors returned by this func will be returned by the cachefunc.Get method
		// Errors are not cached
		return user{}, errors.New("user not found")
	}, time.Hour) // Sets duration of cached items to 1 hour
}
```

### Interface

The `CacheFunc` interface contains the method to cache the result of func `fn` for key `k`.

```go
type CacheFunc[K comparable, V any] interface {
    Get(k K, fn func() (V, error), ttl time.Duration) (V, error)
}
```

### Implementations

#### Default

The default implementation is Memory. See below.

```go
cache := cachefunc.New[int, string]()
```

#### Memory

Returns an in-memory implementation, which stores items in a map.

```go
memCache := cachefunc.NewMemory[int, string]()
```

#### Nop

Returns a no-operation implementation, which always calls the provided func. Useful for tests.

```go
nopCache := cachefunc.NewNop[int, string]()
```