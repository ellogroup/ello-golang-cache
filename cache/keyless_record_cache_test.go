package cache

import (
	"context"
	"fmt"
	"github.com/ellogroup/ello-golang-cache/driver"
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"
)

func newActiveRecordCacheStub() *RecordCache[int, string] {
	c := driver.NewMemoryCache[int, RecordCacheItem[string]]()
	c.Set(context.Background(), 0, RecordCacheItem[string]{V: "active1", T: time.Now().Add(time.Hour * 1)})
	return &RecordCache[int, string]{
		log:       zap.NewNop(),
		cache:     c,
		recordTtl: 10 * time.Second,
	}
}

func newStaleRecordCacheStub() *RecordCache[int, string] {
	c := driver.NewMemoryCache[int, RecordCacheItem[string]]()
	c.Set(context.Background(), 0, RecordCacheItem[string]{V: "active1", T: time.Now().Add(-time.Hour * 1)})
	return &RecordCache[int, string]{
		log:             zap.NewNop(),
		cache:           c,
		onDemandFetcher: nil,
		recordTtl:       10 * time.Second,
	}
}

func TestKeylessRecordCache_Get(t *testing.T) {
	type testCase[V any] struct {
		name    string
		k       *KeylessRecordCache[V]
		want    string
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Successfully fetch without key",
			k:       &KeylessRecordCache[string]{newActiveRecordCacheStub()},
			want:    "active1",
			wantErr: false,
		},
		{
			name:    "return error when record Stale",
			k:       &KeylessRecordCache[string]{newStaleRecordCacheStub()},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.Get(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockFetcher struct {
}

func (m *mockFetcher) Fetch(_ context.Context) (string, error) {
	return "hello", nil
}

type mockFetcherError struct {
}

func (m *mockFetcherError) Fetch(_ context.Context) (string, error) {
	return "", fmt.Errorf("error")
}

func TestNewKeylessRecordCacheOnDemand(t *testing.T) {
	type args struct {
		driver driver.CacheDriver[int, RecordCacheItem[string]]
		f      KeylessFetcher[string]
		ttl    time.Duration
	}
	type testCase[V any] struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name: "successfully init KeylessRecordCacheOnDemand and get value",
			args: args{
				driver: driver.NewMemoryCache[int, RecordCacheItem[string]](),
				f:      &mockFetcher{},
				ttl:    10 * time.Second,
			},
			want:    "hello",
			wantErr: false,
		},
		{
			name: "successfully init KeylessRecordCacheOnDemand and get error",
			args: args{
				driver: driver.NewMemoryCache[int, RecordCacheItem[string]](),
				f:      &mockFetcherError{},
				ttl:    10 * time.Second,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeylessRecordCacheOnDemand(tt.args.driver, tt.args.f, tt.args.ttl)
			got, err := k.Get(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeylessRecordCacheOnDemand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKeylessRecordCacheAsync(t *testing.T) {
	type args struct {
		driver driver.CacheDriver[int, RecordCacheItem[string]]
		f      KeylessFetcher[string]
		ttl    time.Duration
	}
	type testCase[V any] struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name: "successfully init KeylessRecordCacheAsync and get value",
			args: args{
				driver: driver.NewMemoryCache[int, RecordCacheItem[string]](),
				f:      &mockFetcher{},
				ttl:    10000 * time.Second,
			},
			want:    "hello",
			wantErr: false,
		},
		{
			name: "successfully init KeylessRecordCacheAsync and get error",
			args: args{
				driver: driver.NewMemoryCache[int, RecordCacheItem[string]](),
				f:      &mockFetcherError{},
				ttl:    10 * time.Second,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeylessRecordCacheAsync(tt.args.driver, tt.args.f, tt.args.ttl)
			got, err := k.Get(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeylessRecordCacheAsync() = %v, want %v", got, tt.want)
			}
		})
	}
}
