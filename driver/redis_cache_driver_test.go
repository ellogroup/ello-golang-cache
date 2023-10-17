package driver

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"reflect"
	"testing"
)

func TestRedisDriver_All(t *testing.T) {
	type testCase[K Key, V Value] struct {
		name string
		c    RedisCache[K, V]
		want map[K]V
	}
	tests := []testCase[Key, Value]{
		{
			name: "Get All",
			c: RedisCache[Key, Value]{
				c:   getRedisMockGetAll(),
				key: "test",
			},
			want: map[Key]Value{
				Key{Id: "key_1"}: {
					Value: "value_1",
					Inner: Inner{"inner_value_1"},
				},
				Key{Id: "key_2"}: {
					Value: "value_2",
					Inner: Inner{"inner_value_2"},
				},
			},
		},
		{
			name: "Get All Empty",
			c: RedisCache[Key, Value]{
				c:   getRedisMockGetAllEmpty(),
				key: "test",
			},
			want: map[Key]Value{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.All(context.Background()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDriver_Clear(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		c    RedisCache[K, V]
		want bool
	}
	tests := []testCase[Key, Value]{
		{
			name: "successful clear",
			c: RedisCache[Key, Value]{

				c:   getRedisMockClear(),
				key: "test",
			},
			want: true,
		},
		{
			name: "unsuccessful clear",
			c: RedisCache[Key, Value]{

				c:   getRedisMockClearError(),
				key: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Clear(context.Background()); got != tt.want {
				t.Errorf("Clear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDriver_Delete(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		c    RedisCache[K, V]
		arg  K
		want bool
	}
	tests := []testCase[Key, Value]{
		{
			name: "successful delete",
			c: RedisCache[Key, Value]{

				c:   getRedisMockDelete(),
				key: "test",
			},
			arg:  key1(),
			want: true,
		},
		{
			name: "key not in cache, returns true",
			c: RedisCache[Key, Value]{

				c:   getRedisMockDelete(),
				key: "test",
			},
			arg:  key2(),
			want: true,
		},
		{
			name: "unsuccessful delete",
			c: RedisCache[Key, Value]{

				c:   getRedisMockDeleteError(),
				key: "test",
			},
			arg:  key1(),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Delete(context.Background(), tt.arg); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDriver_Get(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name  string
		c     RedisCache[K, V]
		arg   K
		want  V
		want1 bool
	}
	tests := []testCase[Key, Value]{
		{
			name: "successful get",
			c: RedisCache[Key, Value]{

				c:   getRedisMockGet(),
				key: "test",
			},
			arg:   key1(),
			want:  value1(),
			want1: true,
		},
		{
			name: "get no value set (mocked for key)",
			c: RedisCache[Key, Value]{

				c:   getRedisMockGetNil(),
				key: "test",
			},
			arg:   key1(),
			want1: false,
		},
		{
			name: "get no value set (not mocked for key)",
			c: RedisCache[Key, Value]{

				c:   getRedisMockGetNil(),
				key: "test",
			},
			arg:   key2(),
			want1: false,
		},
		{
			name: "get returns error",
			c: RedisCache[Key, Value]{

				c:   getRedisMockGetError(),
				key: "test",
			},
			arg:   key2(),
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(context.Background(), tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRedisDriver_Has(t *testing.T) {
	type testCase[K comparable, V any] struct {
		name string
		c    RedisCache[K, V]
		arg  K
		want bool
	}
	tests := []testCase[Key, Value]{
		{
			name: "has true",
			c: RedisCache[Key, Value]{

				c:   getRedisMockHasTrue(),
				key: "test",
			},
			arg:  key1(),
			want: true,
		},
		{
			name: "has false (mocked for key)",
			c: RedisCache[Key, Value]{

				c:   getRedisMockHasFalse(),
				key: "test",
			},
			arg:  key1(),
			want: false,
		},
		{
			name: "has false (not mocked for key)",
			c: RedisCache[Key, Value]{

				c:   getRedisMockHasTrue(),
				key: "test",
			},
			arg:  key2(),
			want: false,
		},
		{
			name: "get returns error",
			c: RedisCache[Key, Value]{

				c:   getRedisMockHasError(),
				key: "test",
			},
			arg:  key1(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Has(context.Background(), tt.arg); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisDriver_Set(t *testing.T) {
	type args struct {
		key   Key
		value Value
	}
	type testCase[K comparable, V any] struct {
		name string
		c    RedisCache[K, V]
		args args
		want bool
	}
	tests := []testCase[Key, Value]{
		{
			name: "successful set",
			c: RedisCache[Key, Value]{

				c:   getRedisMockSet(),
				key: "test",
			},
			args: args{key: key1(), value: value1()},
			want: true,
		},
		{
			name: "set error returned",
			c: RedisCache[Key, Value]{

				c:   getRedisMockSetError(),
				key: "test",
			},
			args: args{key: key1(), value: value1()},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Set(context.Background(), tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

type Key struct {
	Id string
}

func (k Key) toGob() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(k)
	if err != nil {
		return ""
	}
	return b.String()
}

type Inner struct {
	InnerValue string
}

type Value struct {
	Value string
	Inner Inner
}

func (v Value) toGob() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)
	if err != nil {
		return ""
	}
	return b.String()
}

func getRedisMockGetAll() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHGetAll("test").SetVal(
		map[string]string{key1().toGob(): value1().toGob(), key2().toGob(): value2().toGob()})
	return r
}

func getRedisMockGetAllEmpty() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHGetAll("test").SetVal(
		map[string]string{})
	return r
}

func getRedisMockGet() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHGet("test", key1().toGob()).SetVal(value1().toGob())
	return r
}

func getRedisMockGetError() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHGet("test", key1().toGob()).SetErr(fmt.Errorf("error"))
	return r
}

func getRedisMockGetNil() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHGet("test", key1().toGob()).RedisNil()
	return r
}

func getRedisMockSet() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHSet("test", key1().toGob(), value1().toGob()).SetVal(1)
	return r
}

func getRedisMockSetError() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHSet("test", key1().toGob(), value1().toGob()).SetErr(fmt.Errorf("error"))
	return r
}

func getRedisMockHasTrue() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHExists("test", key1().toGob()).SetVal(true)
	return r
}

func getRedisMockHasFalse() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHExists("test", key1().toGob()).SetVal(false)
	return r
}

func getRedisMockHasError() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHExists("test", key1().toGob()).SetErr(fmt.Errorf("error"))
	return r
}

func getRedisMockDelete() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHDel("test", key1().toGob()).SetVal(1)
	return r
}

func getRedisMockDeleteError() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectHDel("test", key1().toGob()).SetErr(fmt.Errorf("error"))
	return r
}

func getRedisMockClear() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectDel("test").SetVal(1)
	return r
}

func getRedisMockClearError() *redis.Client {
	r, mock := redismock.NewClientMock()
	mock.ExpectDel("test").SetErr(fmt.Errorf("error"))
	return r
}

func value1() Value {
	return Value{Value: "value_1", Inner: Inner{InnerValue: "inner_value_1"}}
}

func value2() Value {
	return Value{Value: "value_2", Inner: Inner{InnerValue: "inner_value_2"}}
}

func key2() Key {
	return Key{Id: "key_2"}
}

func key1() Key {
	return Key{Id: "key_1"}
}
