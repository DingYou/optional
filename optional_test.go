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

func TestMust(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Panics(t, func() {
			Empty[int]().MustVal()
		})
	})
	t.Run("not nil", func(t *testing.T) {
		assert.Equal(t, Of(1).MustVal(), 1)
	})
}
