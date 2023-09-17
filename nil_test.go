package optional

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	type args[T any] struct {
		val T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[any]{
		{
			name: "test int 0",
			args: args[any]{val: 0},
			want: false,
		},
		{
			name: "test pointer nil",
			args: args[any]{val: func() *int {
				return nil
			}()},
			want: true,
		},
		{
			name: "test pointer",
			args: args[any]{val: func() *int {
				a := 0
				return &a
			}()},
			want: false,
		},
		{
			name: "test nil",
			args: args[any]{nil},
			want: true,
		},
		{
			name: "test nil map",
			args: args[any]{func() map[string]string {
				var m map[string]string
				return m
			}()},
			want: true,
		},
		{
			name: "test empty map",
			args: args[any]{make(map[string]string)},
			want: false,
		},
		{
			name: "test nil slice",
			args: args[any]{func() []uint8 {
				var s []uint8
				return s
			}()},
			want: true,
		},
		{
			name: "test empty slice",
			args: args[any]{make([]uint8, 0)},
			want: false,
		},
		{
			name: "test struct pointer nil",
			args: args[any]{func() *time.Time {
				var t *time.Time
				return t
			}()},
			want: true,
		},
		{
			name: "test struct pointer by new",
			args: args[any]{new(time.Time)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.args.val); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNilWithInt(t *testing.T) {
	type args struct {
		val int
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{
		{
			name: "test 0",
			args: args{0},
			want: false,
		},
		{
			name: "test func",
			args: args{val: func() int {
				return 0
			}()},
			want: false,
		},
		{
			name: "test 1",
			args: args{1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.args.val); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNilWithFloat(t *testing.T) {
	type args struct {
		val float64
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{
		{
			name: "test 0",
			args: args{0},
			want: false,
		},
		{
			name: "test NaN by func",
			args: args{val: func() float64 {
				return math.NaN()
			}()},
			want: false,
		},
		{
			name: "test -Inf",
			args: args{math.Inf(-1)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.args.val); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequiredNonNull(t *testing.T) {
	t.Run("test nil", func(t *testing.T) {
		assert.Panics(t, func() {
			RequireNonNull[any](nil)
		})
	})
	t.Run("test not nil", func(t *testing.T) {
		assert.Equal(t, RequireNonNull[int](1), 1)
	})
}
