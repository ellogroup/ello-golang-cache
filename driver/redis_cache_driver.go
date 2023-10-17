package driver

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/redis/go-redis/v9"
)

type RedisCache[K comparable, V any] struct {
	c   *redis.Client
	key string
}

func (r *RedisCache[K, V]) Has(ctx context.Context, key K) bool {
	var k bytes.Buffer
	err := gob.NewEncoder(&k).Encode(key)
	if err != nil {
		return false
	}
	has := r.c.HExists(ctx, r.key, k.String())
	if has.Err() != nil {
		return false
	}
	return has.Val()
}

func (r *RedisCache[K, V]) Get(ctx context.Context, key K) (V, bool) {
	var k bytes.Buffer
	err := gob.NewEncoder(&k).Encode(key)
	if err != nil {
		return *new(V), false
	}
	b, err := r.c.HGet(ctx, r.key, k.String()).Bytes()
	if err != nil {
		return *new(V), false
	}
	var v V
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&v)
	return v, err == nil
}

func (r *RedisCache[K, V]) All(ctx context.Context) map[K]V {
	m := map[K]V{}
	all := r.c.HGetAll(ctx, r.key).Val()
	for key, value := range all {
		var k K
		err := gob.NewDecoder(bytes.NewReader([]byte(key))).Decode(&k)
		if err != nil {
			continue
		}
		var v V
		err = gob.NewDecoder(bytes.NewReader([]byte(value))).Decode(&v)
		if err != nil {
			continue
		}
		m[k] = v
	}
	return m
}

func (r *RedisCache[K, V]) Set(ctx context.Context, key K, value V) bool {
	var k bytes.Buffer
	err := gob.NewEncoder(&k).Encode(key)
	if err != nil {
		return false
	}
	var v bytes.Buffer
	err = gob.NewEncoder(&v).Encode(value)
	if err != nil {
		return false
	}
	err = r.c.HSet(ctx, r.key, k.String(), v.String()).Err()
	return err == nil
}

func (r *RedisCache[K, V]) Delete(ctx context.Context, key K) bool {
	var k bytes.Buffer
	err := gob.NewEncoder(&k).Encode(key)
	r.c.HDel(ctx, r.key, k.String())
	return err == nil
}

func (r *RedisCache[K, V]) Clear(ctx context.Context) bool {
	err := r.c.Del(ctx, r.key).Err()
	return err == nil
}

func NewRedisCacheDriver[K comparable, V any](redisKey string, redisClient *redis.Client) *RedisCache[K, V] {
	return &RedisCache[K, V]{
		c:   redisClient,
		key: redisKey,
	}
}
