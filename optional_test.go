package optional

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional_String(t *testing.T) {
	type s interface {
		String() string
	}
	type testCase struct {
		name string
		args s
		want string
	}

	tests := []testCase{
		{
			name: "empty slice string",
			args: OfNullable(make([]string, 0)),
			want: "Optional[[]string]{[]}",
		},
		{
			name: "slice string Ofnullable",
			args: OfNullable([]string{"a", "b", "c"}),
			want: fmt.Sprintf("Optional[[]string]{%v}", []string{"a", "b", "c"}),
		},
		{
			name: "slice string",
			args: Of([]string{"a", "b", "c"}),
			want: fmt.Sprintf("Optional[[]string]{%v}", []string{"a", "b", "c"}),
		},
		{
			name: "empty slice",
			args: OfNullable(func() any {
				var s []string
				return s
			}()),
			want: "EmptyOptional[[]string]{[]}",
		},
		{
			name: "empty int",
			args: Empty[int](),
			want: "EmptyOptional[int]{0}",
		},
		{
			name: "nil string pointer",
			args: OfNullable[*string](nil),
			want: "EmptyOptional[*string]{<nil>}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptional_MustVal(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Panics(t, func() {
			Empty[int]().MustVal()
		})
	})
	t.Run("not nil", func(t *testing.T) {
		assert.Equal(t, Of(1).MustVal(), 1)
	})
}

func TestOptional_IsEmpty(t *testing.T) {
	type isEmptyable interface {
		IsEmpty() bool
	}
	type testCase struct {
		name string
		args isEmptyable
		want bool
	}

	tests := []testCase{
		{
			name: "nil empty",
			args: OfNullable[any](nil),
			want: true,
		},
		{
			name: "not empty",
			args: Of(1),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptional_OrElse(t *testing.T) {
	type args[T any] struct {
		o Optional[T]
		d T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}

	tests := []testCase[any]{
		{
			name: "int else",
			args: args[any]{
				o: Of[any](2),
				d: 1,
			},
			want: 2,
		},
		{
			name: "zero else",
			args: args[any]{
				o: Of[any](0),
				d: 1,
			},
			want: 0,
		},
		{
			name: "nil else",
			args: args[any]{
				o: Empty[any](),
				d: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.o.OrElse(tt.args.d); got != tt.want {
				t.Errorf("OrElse(%v) = %v, want %v", tt.args.d, got, tt.want)
			}
		})
	}
}
