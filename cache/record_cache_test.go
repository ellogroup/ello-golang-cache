package cache

import (
	"context"
	"fmt"
	"github.com/ellogroup/ello-golang-cache/driver"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"reflect"
	"testing"
	"time"
)

func newCacheStub() driver.Cache[string, RecordCacheItem[int]] {
	c := driver.NewMemoryCache[string, RecordCacheItem[int]]()

	c.Set(context.Background(), "active1", RecordCacheItem[int]{V: 1, T: time.Now().Add(time.Hour * 1)})
	c.Set(context.Background(), "active2", RecordCacheItem[int]{V: 2, T: time.Now().Add(time.Hour * 1)})
	c.Set(context.Background(), "stale1", RecordCacheItem[int]{V: 1, T: time.Now().Add(-time.Hour * 1)})
	c.Set(context.Background(), "stale2", RecordCacheItem[int]{V: 2, T: time.Now().Add(-time.Hour * 1)})
	return c
}

type onDemandFetcherMock[K string, v int] struct{}

func newOnDemandFetcherMock[K string, V int]() onDemandFetcherMock[K, V] {
	return onDemandFetcherMock[K, V]{}
}

func (sf onDemandFetcherMock[K, V]) FetchByKey(_ context.Context, _ string) (int, error) {
	return 10, nil
}

type asyncFetcherMock[s string, i int] struct{}

func newAsyncFetcherMock() asyncFetcherMock[string, int] {
	return asyncFetcherMock[string, int]{}
}

func (a asyncFetcherMock[K, V]) FetchAll(_ context.Context) (map[string]int, error) {
	return map[string]int{
		"active1": 1,
		"active2": 2,
		"stale1":  1,
		"stale2":  2,
	}, nil
}

type FetcherError[K string, v int] struct{}

func newFetcherError[K string, V int]() FetcherError[K, V] {
	return FetcherError[K, V]{}
}

func (sf FetcherError[K, v]) FetchByKey(_ context.Context, _ string) (int, error) {
	return 1000, fmt.Errorf("error")
}

func TestRecordCache_Get(t *testing.T) {
	type fields struct {
		onDemandFetcher OnDemandFetcher[string, int]
		asyncFetcher    AsyncFetcher[string, int]
		recordTtl       time.Duration
		allTtl          time.Duration
		lastUpdated     time.Time
		cron            *cron.Cron
	}
	tests := []struct {
		name    string
		fields  fields
		args    string
		want    int
		wantErr bool
	}{
		{
			name: "Successfully fetch active record",
			fields: fields{
				onDemandFetcher: nil,
				asyncFetcher:    nil,
				recordTtl:       100 * time.Second,
				allTtl:          100 * time.Second,
			},
			args:    "active1",
			want:    1,
			wantErr: false,
		},
		{
			name: "Successfully refreshItem and return Stale record",
			fields: fields{
				onDemandFetcher: newOnDemandFetcherMock(),
				asyncFetcher:    nil,
				recordTtl:       100 * time.Second,
				allTtl:          100 * time.Second,
			},
			args:    "stale1",
			want:    10,
			wantErr: false,
		},
		{
			name: "Error with fetching Stale record",
			fields: fields{
				onDemandFetcher: newFetcherError(),
				asyncFetcher:    nil,
				recordTtl:       100 * time.Second,
				allTtl:          100 * time.Second,
			},
			args:    "stale1",
			want:    0,
			wantErr: true,
		},
		{
			name: "Error as onDemandFetcher is nil",
			fields: fields{
				onDemandFetcher: nil,
				asyncFetcher:    nil,
				recordTtl:       100 * time.Second,
				allTtl:          100 * time.Second,
			},
			args:    "stale1",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RecordCache[string, int]{
				log:             zap.NewNop(),
				onDemandFetcher: tt.fields.onDemandFetcher,
				asyncFetcher:    tt.fields.asyncFetcher,
				cache:           newCacheStub(),
				recordTtl:       tt.fields.recordTtl,
				allTtl:          tt.fields.allTtl,
				lastUpdated:     tt.fields.lastUpdated,
				cron:            tt.fields.cron,
			}
			got, err := r.Get(context.Background(), tt.args)
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

func TestRecordCache_SetAsyncFetcher(t *testing.T) {
	type fields struct {
		recordTtl time.Duration
	}
	type args struct {
		ttl time.Duration
		k   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     int
		wantErr  bool
		nextCron time.Time
	}{
		{
			name: "successfully fetch all and return cached active value",
			fields: fields{
				recordTtl: 100 * time.Second,
			},
			args: args{
				ttl: 100 * time.Second,
				k:   "active2",
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "error when records added to cache is Stale",
			args: args{
				ttl: -70 * time.Second,
				k:   "active2",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRecordCache[string, int](driver.NewMemoryCache[string, RecordCacheItem[int]]())
			rc := r.SetAsyncFetcher(newAsyncFetcherMock(), tt.args.ttl)
			rc.recordTtl = tt.fields.recordTtl
			got, err := rc.Get(context.Background(), tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refreshAllRecordsEvery() = %v, want %v", got, tt.want)
			}
			timeRound := time.Now().Add(time.Second * 30).Round(time.Minute)
			next := rc.cron.Entries()[0].Next
			if !timeRound.Equal(next) {
				fmt.Printf("Last refreshItem: %v, round = %v", rc.cron.Entries()[0].Next.Add(-time.Minute*1), timeRound)
				t.Errorf("next scheuled refreshItem all = %v, want= %v", rc.cron.Entries()[0].Next, timeRound)
			}
		})
	}
}

func TestRecordCache_RefreshStaleRecordsEvery(t *testing.T) {
	type fields struct {
		sf OnDemandFetcher[string, int]
	}
	type args struct {
		ttl time.Duration
		k   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "successfully return cached active value",
			fields: fields{sf: newOnDemandFetcherMock()},
			args: args{
				ttl: 100 * time.Second,
				k:   "active2",
			},
			want:    2,
			wantErr: false,
		},
		{
			name:   "successfully refreshItem and return Stale value",
			fields: fields{sf: newOnDemandFetcherMock()},
			args: args{
				ttl: 100 * time.Second,
				k:   "stale2",
			},
			want:    10,
			wantErr: false,
		},
		{
			name:   "error  onDemandFetcher not fixed",
			fields: fields{sf: nil},
			args: args{
				ttl: 100 * time.Second,
				k:   "stale2",
			},
			want:    0,
			wantErr: true,
		},
		{
			name:   "error  onDemandFetcher returns error ",
			fields: fields{sf: newFetcherError()},
			args: args{
				ttl: 100 * time.Second,
				k:   "stale2",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RecordCache[string, int]{
				log:             zap.NewNop(),
				cache:           newCacheStub(),
				onDemandFetcher: tt.fields.sf,
			}
			got, err := r.refreshStaleRecordsEvery(tt.args.ttl).Get(context.Background(), tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("refreshAllRecordsEvery() = %v, want %v", got, tt.want)
			}
		})
	}
}
