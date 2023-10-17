package driver

import (
	"context"
	"reflect"
	"testing"
)

func TestMemoryCache_All(t *testing.T) {
	type fields struct {
		c map[int]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[int]string
		wantLen int
	}{
		{
			name: "cache full returns non-empty map",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			want:    map[int]string{1: "one", 2: "two"},
			wantLen: 2,
		},
		{
			name: "cache empty returns empty map",
			fields: fields{
				c: make(map[int]string),
			},
			want:    make(map[int]string),
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryCache[int, string]{
				c: tt.fields.c,
			}
			if got := m.All(context.Background()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
			if got := len(m.All(context.Background())); !reflect.DeepEqual(got, tt.wantLen) {
				t.Errorf("All() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	type fields struct {
		c map[int]string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "cache full returns true",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			want: true,
		},
		{
			name: "cache empty returns true",
			fields: fields{
				c: map[int]string{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryCache[int, string]{
				c: tt.fields.c,
			}
			if got := m.Clear(context.Background()); got != tt.want {
				t.Errorf("Clear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	type fields struct {
		c map[int]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    int
		want    bool
		wantLen int
	}{
		{
			name: "cache has k removes and returns true",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			args:    1,
			want:    true,
			wantLen: 1,
		},
		{
			name: "cache empty returns true",
			fields: fields{
				c: map[int]string{},
			},
			args:    1,
			want:    true,
			wantLen: 0,
		},
		{
			name: "cache doe not have k, returns true",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			args:    9,
			want:    true,
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryCache[int, string]{
				c: tt.fields.c,
			}
			got := m.Delete(context.Background(), tt.args)
			if got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(len(m.c), tt.wantLen) {
				t.Errorf("cache size = %v, wanted %v", len(m.c), tt.wantLen)
			}
		})
	}
}

func TestMemoryCache_Get(t *testing.T) {
	type fields struct {
		c map[int]string
	}
	tests := []struct {
		name   string
		fields fields
		args   int
		want   string
		want1  bool
	}{
		{
			name: "cache has k   returns k and true",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			args:  1,
			want:  "one",
			want1: true,
		},
		{
			name: "empty cache   returns empty k return and false",
			fields: fields{
				c: map[int]string{},
			},
			args:  1,
			want:  "",
			want1: false,
		},
		{
			name: "cache does not have k  empty k return and false",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			args:  10,
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryCache[int, string]{
				c: tt.fields.c,
			}
			got, got1 := m.Get(context.Background(), tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemoryCache_Has(t *testing.T) {
	type fields struct {
		c map[int]string
	}
	tests := []struct {
		name   string
		fields fields
		args   int
		want   bool
	}{
		{
			name: "cache has k   returns true",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			args: 1,
			want: true,
		},
		{
			name: "cache does not hav k   returns false",
			fields: fields{
				c: map[int]string{1: "one", 2: "two"},
			},
			args: 10,
			want: false,
		},
		{
			name: "empty cache   returns false",
			fields: fields{
				c: map[int]string{},
			},
			args: 1,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryCache[int, string]{
				c: tt.fields.c,
			}
			if got := m.Has(context.Background(), tt.args); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryCache_Set(t *testing.T) {
	type fields struct {
		c map[int]string
	}
	type args struct {
		k int
		v string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantLen int
	}{
		{
			name: "add to empty cache  returns true",
			fields: fields{
				c: map[int]string{},
			},
			args: args{
				k: 1,
				v: "one",
			},
			want:    true,
			wantLen: 1,
		},
		{
			name: "add new key to init cache  returns true",
			fields: fields{
				c: map[int]string{1: "one"},
			},
			args: args{
				k: 2,
				v: "two",
			},
			want:    true,
			wantLen: 2,
		},
		{
			name: "add existing key to init cache  overwrites value",
			fields: fields{
				c: map[int]string{1: "one"},
			},
			args: args{
				k: 1,
				v: "two",
			},
			want:    true,
			wantLen: 1,
		},
		{
			name: "add existing key with nil value  adds and returns true",
			fields: fields{
				c: map[int]string{1: "one"},
			},
			args: args{
				k: 3,
			},
			want:    true,
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemoryCache[int, string]{
				c: tt.fields.c,
			}
			got := m.Set(context.Background(), tt.args.k, tt.args.v)
			if got != tt.want {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(len(m.c), tt.wantLen) {
				t.Errorf("cache size = %v, wanted %v", len(m.c), tt.wantLen)
			}
		})
	}
}

func TestNewMemoryCache(t *testing.T) {
	tests := []struct {
		name string
		want map[int]string
	}{{
		name: "init memory cache successfully",
		want: map[int]string{},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryCache[int, string](); !reflect.DeepEqual(got.c, tt.want) {
				t.Errorf("NewMemoryCache() = %v, want %v", got.c, tt.want)
			}
		})
	}
}
