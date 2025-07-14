package cachefunc

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNop_Get(t *testing.T) {
	errFn := errors.New("fn error")
	type args[K comparable, V any] struct {
		key K
		fn  func() (V, error)
		exp time.Duration
	}
	type testCase[K comparable, V any] struct {
		name    string
		n       Nop[K, V]
		args    args[K, V]
		want    V
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[string, int]{
		{
			name: "func returns value, returns value",
			n:    Nop[string, int]{},
			args: args[string, int]{"key", func() (int, error) {
				return 123, nil
			}, NoExpire},
			want:    123,
			wantErr: assert.NoError,
		},
		{
			name: "func returns error, returns error",
			n:    Nop[string, int]{},
			args: args[string, int]{"key", func() (int, error) {
				return 0, errFn
			}, NoExpire},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, errFn, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Get(tt.args.key, tt.args.fn, tt.args.exp)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v, <func>, %v)", tt.args.key, tt.args.exp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v, <func>, %v)", tt.args.key, tt.args.exp)
		})
	}
}
