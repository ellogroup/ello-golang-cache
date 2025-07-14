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
