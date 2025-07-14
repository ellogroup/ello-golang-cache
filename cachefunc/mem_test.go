package cachefunc

import (
	"errors"
	"fmt"
	"github.com/ellogroup/ello-golang-clock/clock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMemCache_Get(t *testing.T) {
	errFn := errors.New("fn error")

	past := time.Date(2025, time.January, 5, 12, 0, 0, 0, time.Local)
	now := time.Date(2025, time.January, 10, 12, 0, 0, 0, time.Local)
	future := time.Date(2025, time.January, 15, 12, 0, 0, 0, time.Local)
	expFuture := future.Sub(now)
	sysClock = clock.NewFixed(now)

	type args[K comparable, V any] struct {
		k   K
		fn  func() (V, error)
		exp time.Duration
	}
	type testCase[K comparable, V any] struct {
		name    string
		m       map[K]*Item[V]
		args    args[K, V]
		want    V
		wantM   map[K]*Item[V]
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[string, int]{
		{
			name: "value already in cache, returns cached value",
			m:    map[string]*Item[int]{"test": {val: 123}},
			args: args[string, int]{
				k:   "test",
				fn:  func() (int, error) { return 0, nil },
				exp: expFuture,
			},
			want:    123,
			wantM:   map[string]*Item[int]{"test": {val: 123}},
			wantErr: assert.NoError,
		},
		{
			name: "value already in cache with exp in future, returns cached value",
			m:    map[string]*Item[int]{"test": {val: 123, exp: future}},
			args: args[string, int]{
				k:   "test",
				fn:  func() (int, error) { return 0, nil },
				exp: expFuture,
			},
			want:    123,
			wantM:   map[string]*Item[int]{"test": {val: 123, exp: future}},
			wantErr: assert.NoError,
		},
		{
			name: "value in cache with exp in past, stores func value in cache, returns value",
			m:    map[string]*Item[int]{"test": {val: 123, exp: past}},
			args: args[string, int]{
				k:   "test",
				fn:  func() (int, error) { return 123, nil },
				exp: expFuture,
			},
			want:    123,
			wantM:   map[string]*Item[int]{"test": {val: 123, exp: future}},
			wantErr: assert.NoError,
		},
		{
			name: "value in cache with exp now, stores func value in cache, returns value",
			m:    map[string]*Item[int]{"test": {val: 123, exp: now}},
			args: args[string, int]{
				k:   "test",
				fn:  func() (int, error) { return 123, nil },
				exp: expFuture,
			},
			want:    123,
			wantM:   map[string]*Item[int]{"test": {val: 123, exp: future}},
			wantErr: assert.NoError,
		},
		{
			name: "value not in cache, stores func value in cache, returns value",
			m:    map[string]*Item[int]{"test": {val: 123}},
			args: args[string, int]{
				k:   "another-test",
				fn:  func() (int, error) { return 456, nil },
				exp: NoExpire,
			},
			want:    456,
			wantM:   map[string]*Item[int]{"test": {val: 123}, "another-test": {val: 456}},
			wantErr: assert.NoError,
		},
		{
			name: "value not in cache, stores func value with exp in cache, returns value",
			m:    map[string]*Item[int]{"test": {val: 123}},
			args: args[string, int]{
				k:   "another-test",
				fn:  func() (int, error) { return 456, nil },
				exp: expFuture,
			},
			want:    456,
			wantM:   map[string]*Item[int]{"test": {val: 123}, "another-test": {val: 456, exp: future}},
			wantErr: assert.NoError,
		},
		{
			name: "value not in cache, func returns error, returns func error",
			m:    map[string]*Item[int]{"test": {val: 123}},
			args: args[string, int]{
				k:   "another-test",
				fn:  func() (int, error) { return 0, errFn },
				exp: NoExpire,
			},
			want:  0,
			wantM: map[string]*Item[int]{"test": {val: 123}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, errFn, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := &Memory[string, int]{
				m: tt.m,
			}
			got, err := sut.Get(tt.args.k, tt.args.fn, tt.args.exp)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v, <func>, %v)", tt.args.k, tt.args.exp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v, <func>, %v)", tt.args.k, tt.args.exp)
			assert.Equalf(t, tt.wantM, sut.m, "Get(%v, <func>, %v)", tt.args.k, tt.args.exp)
		})
	}
}
