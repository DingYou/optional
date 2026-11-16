package optional

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_CrossType(t *testing.T) {
	calls := 0
	mapped := Map(Of(42), func(val int) string {
		calls++
		return strconv.Itoa(val)
	})
	val, ok := mapped.Get()
	assert.Equal(t, "42", val)
	assert.True(t, ok)
	assert.Equal(t, 1, calls)
	assert.True(t, Map(Of(42), func(int) bool { return false }).IsPresent())
	assert.True(t, Map(Of(42), func(int) *int { return nil }).IsEmpty())
	boxed := Map(Of(42), func(int) any {
		var pointer *int
		return pointer
	})
	assert.True(t, boxed.IsEmpty())
	if val, ok := boxed.Get(); val != nil || ok {
		t.Fatalf("Get() = (%v, %v), want (nil, false)", val, ok)
	}
}

func TestFlatMap_CrossTypeAndCallbackResult(t *testing.T) {
	tests := []struct {
		name   string
		result Optional[string]
	}{
		{name: "present", result: Of("42")},
		{name: "present zero", result: Of("")},
		{name: "empty", result: Empty[string]()},
		{name: "zero value", result: Optional[string]{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			// The explicit result type also checks that FlatMap does not nest Optional.
			var mapped Optional[string] = FlatMap(Of(42), func(val int) Optional[string] {
				calls++
				assert.Equal(t, 42, val)
				return tt.result
			})
			assert.Equal(t, tt.result, mapped)
			assert.Equal(t, 1, calls)
		})
	}

	var pointer *int
	for _, tt := range []struct {
		name   string
		result Optional[any]
	}{
		{name: "present typed nil", result: FromOK[any](pointer, true)},
		{name: "empty typed nil", result: OfNullable[any](pointer)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mapped := FlatMap(Of(42), func(int) Optional[any] { return tt.result })
			assert.Equal(t, tt.result.IsPresent(), mapped.IsPresent())
			if mapped.Val() != any(pointer) {
				t.Fatal("FlatMap did not preserve the callback's stored typed nil")
			}
		})
	}
}

func TestMapAndFlatMap_EmptyShortCircuit(t *testing.T) {
	var zero Optional[int]
	for _, o := range []Optional[int]{zero, Empty[int]()} {
		mapped := Map(o, func(int) string {
			t.Fatal("Map called for an empty Optional")
			return "unexpected"
		})
		assert.True(t, mapped.IsEmpty())
		flatMapped := FlatMap(o, func(int) Optional[string] {
			t.Fatal("FlatMap called for an empty Optional")
			return Of("unexpected")
		})
		assert.True(t, flatMapped.IsEmpty())
	}
	assert.True(t, Map(Empty[int](), (func(int) string)(nil)).IsEmpty())
	assert.True(t, FlatMap(Empty[int](), (func(int) Optional[string])(nil)).IsEmpty())
}

func TestContextAndErrorVariants_PreserveContract(t *testing.T) {
	ctx := context.Background()
	inputErr := errors.New("input error")
	callbackErr := errors.New("callback error")
	tests := []struct {
		name string
		run  func(Optional[int], error, error, *int) (Optional[string], error)
	}{
		{
			name: "MapWithCtx",
			run: func(o Optional[int], _, _ error, calls *int) (Optional[string], error) {
				return MapWithCtx(ctx, o, func(gotCtx context.Context, val int) string {
					*calls += 1
					assert.Equal(t, ctx, gotCtx)
					return strconv.Itoa(val)
				}), nil
			},
		},
		{
			name: "FlatMapWithCtx",
			run: func(o Optional[int], _, _ error, calls *int) (Optional[string], error) {
				return FlatMapWithCtx(ctx, o, func(gotCtx context.Context, val Optional[int]) string {
					*calls += 1
					assert.Equal(t, ctx, gotCtx)
					return strconv.Itoa(val.MustVal())
				}), nil
			},
		},
		{
			name: "MapWithErr",
			run: func(o Optional[int], initial, callback error, calls *int) (Optional[string], error) {
				return MapWithErr(o, func(val int) (string, error) {
					*calls += 1
					return strconv.Itoa(val), callback
				}, initial)
			},
		},
		{
			name: "FlatMapWithErr",
			run: func(o Optional[int], initial, callback error, calls *int) (Optional[string], error) {
				return FlatMapWithErr(o, func(val Optional[int]) (string, error) {
					*calls += 1
					return strconv.Itoa(val.MustVal()), callback
				}, initial)
			},
		},
		{
			name: "MapWithCtxAndErr",
			run: func(o Optional[int], initial, callback error, calls *int) (Optional[string], error) {
				return MapWithCtxAndErr(ctx, o, func(gotCtx context.Context, val int) (string, error) {
					*calls += 1
					assert.Equal(t, ctx, gotCtx)
					return strconv.Itoa(val), callback
				}, initial)
			},
		},
		{
			name: "FlatMapWithCtxAndErr",
			run: func(o Optional[int], initial, callback error, calls *int) (Optional[string], error) {
				return FlatMapWithCtxAndErr(ctx, o, func(gotCtx context.Context, val Optional[int]) (string, error) {
					*calls += 1
					assert.Equal(t, ctx, gotCtx)
					return strconv.Itoa(val.MustVal()), callback
				}, initial)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			present, err := tt.run(Of(42), nil, nil, &calls)
			assert.NoError(t, err)
			assert.Equal(t, "42", present.MustVal())
			assert.Equal(t, 1, calls)
			var zero Optional[int]
			for _, o := range []Optional[int]{zero, Empty[int]()} {
				empty, err := tt.run(o, nil, nil, &calls)
				assert.NoError(t, err)
				assert.True(t, empty.IsEmpty())
				assert.Equal(t, 1, calls)
			}
			if tt.name == "MapWithCtx" || tt.name == "FlatMapWithCtx" {
				return
			}
			empty, err := tt.run(Of(42), inputErr, nil, &calls)
			assert.ErrorIs(t, err, inputErr)
			assert.True(t, empty.IsEmpty())
			assert.Equal(t, 1, calls)
			empty, err = tt.run(Of(42), nil, callbackErr, &calls)
			assert.ErrorIs(t, err, callbackErr)
			assert.True(t, empty.IsEmpty())
			assert.Equal(t, 2, calls)
		})
	}
}

func TestFlatMapVariants_StillWrapCallbackResult(t *testing.T) {
	ctx := context.Background()
	ctxResult := FlatMapWithCtx(ctx, Of(42), func(context.Context, Optional[int]) Optional[string] {
		return Empty[string]()
	})
	assert.True(t, ctxResult.IsPresent())
	assert.True(t, ctxResult.MustVal().IsEmpty())
	errResult, err := FlatMapWithErr(Of(42), func(Optional[int]) (Optional[string], error) {
		return Empty[string](), nil
	}, nil)
	assert.NoError(t, err)
	assert.True(t, errResult.IsPresent())
	assert.True(t, errResult.MustVal().IsEmpty())

	ctxErrResult, err := FlatMapWithCtxAndErr(ctx, Of(42), func(context.Context, Optional[int]) (Optional[string], error) {
		return Empty[string](), nil
	}, nil)
	assert.NoError(t, err)
	assert.True(t, ctxErrResult.IsPresent())
	assert.True(t, ctxErrResult.MustVal().IsEmpty())
}
