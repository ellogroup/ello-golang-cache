package cachefunc

import (
	"github.com/ellogroup/ello-golang-clock/clock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestItem_Expired(t *testing.T) {
	now := time.Date(2025, time.January, 10, 12, 0, 0, 0, time.Local)
	sysClock = clock.NewFixed(now)
	past := now.Add(-1)
	future := now.Add(1)

	type testCase[V any] struct {
		name string
		i    Item[V]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "zero exp, returns false",
			i:    Item[int]{exp: time.Time{}},
			want: false,
		},
		{
			name: "exp in past, returns true",
			i:    Item[int]{exp: past},
			want: true,
		},
		{
			name: "exp now, returns true",
			i:    Item[int]{exp: now},
			want: true,
		},
		{
			name: "exp in future, returns false",
			i:    Item[int]{exp: future},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.i.Expired(), "Expired()")
		})
	}
}
