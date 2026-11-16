package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional_ZeroValue(t *testing.T) {
	var o Optional[int]
	assert.True(t, o.IsEmpty())
	assert.False(t, o.IsPresent())
	val, ok := o.Get()
	assert.Zero(t, val)
	assert.False(t, ok)
	assert.Zero(t, o.Val())
	assert.Equal(t, 42, o.OrElse(42))
	assert.Equal(t, "EmptyOptional[int]{0}", o.String())
	assert.PanicsWithValue(t, ErrNil, func() { o.MustVal() })

	calls := 0
	assert.Equal(t, 42, o.OrElseGet(func() int {
		calls++
		return 42
	}))
	assert.Equal(t, 1, calls)
	o.IfPresent(func(int) { t.Fatal("IfPresent called for zero-value Optional") })
	assert.True(t, o.Filter(func(int) bool {
		t.Fatal("Filter called for zero-value Optional")
		return true
	}).IsEmpty())

	var ptr Optional[*int]
	ptrVal, ptrOK := ptr.Get()
	assert.Nil(t, ptrVal)
	assert.False(t, ptrOK)
	var boolean Optional[bool]
	boolVal, boolOK := boolean.Get()
	assert.False(t, boolVal)
	assert.False(t, boolOK)
}

func TestOptional_ZeroAndFalseArePresent(t *testing.T) {
	for _, o := range []Optional[int]{Of(0), OfNullable(0), FromOK(0, true)} {
		assert.True(t, o.IsPresent())
		assert.False(t, o.IsEmpty())
		val, ok := o.Get()
		assert.Zero(t, val)
		assert.True(t, ok)
		assert.Zero(t, o.MustVal())
		assert.Zero(t, o.OrElse(42))
		assert.Zero(t, o.OrElseGet(func() int {
			t.Fatal("supplier called for present zero")
			return 42
		}))
		calls := 0
		o.IfPresent(func(val int) {
			calls++
			assert.Zero(t, val)
		})
		assert.Equal(t, 1, calls)
		assert.True(t, o.Filter(func(val int) bool { return val == 0 }).IsPresent())
		assert.True(t, o.Filter(func(int) bool { return false }).IsEmpty())
	}
	for _, o := range []Optional[bool]{Of(false), OfNullable(false), FromOK(false, true)} {
		val, ok := o.Get()
		assert.False(t, val)
		assert.True(t, ok)
		assert.True(t, o.IsPresent())
		assert.False(t, o.MustVal())
		assert.False(t, o.OrElse(true))
	}
}

func TestFromOK_DiscardAbsentValue(t *testing.T) {
	o := FromOK(42, false)
	assert.Equal(t, Empty[int](), o)
	val, ok := o.Get()
	assert.Zero(t, val)
	assert.False(t, ok)
	assert.Zero(t, o.Val())

	n := 42
	ptr := FromOK(&n, false)
	assert.Equal(t, Empty[*int](), ptr)
	assert.Nil(t, ptr.Val())
	boxed := FromOK[any](&n, false)
	boxedVal, boxedOK := boxed.Get()
	assert.Nil(t, boxedVal)
	assert.False(t, boxedOK)
	assert.Nil(t, boxed.Val())
}

func TestOptional_NilContracts(t *testing.T) {
	var pointer *int
	var slice []int
	var mapping map[string]int
	var channel chan int
	var function func()
	for _, tt := range []struct {
		name string
		val  any
	}{
		{name: "nil interface", val: nil},
		{name: "typed pointer", val: pointer},
		{name: "typed slice", val: slice},
		{name: "typed map", val: mapping},
		{name: "typed channel", val: channel},
		{name: "typed function", val: function},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.PanicsWithValue(t, ErrNil, func() { Of(tt.val) })
			nullable := OfNullable(tt.val)
			assert.True(t, nullable.IsEmpty())
			val, ok := nullable.Get()
			assert.Nil(t, val)
			assert.False(t, ok)
			assert.IsType(t, tt.val, nullable.Val())
			assert.Nil(t, nullable.Val())
			assert.PanicsWithValue(t, ErrNil, func() { nullable.MustVal() })

			present := FromOK(tt.val, true)
			assert.True(t, present.IsPresent())
			val, ok = present.Get()
			assert.IsType(t, tt.val, val)
			assert.Nil(t, val)
			assert.True(t, ok)
			assert.IsType(t, tt.val, present.MustVal())
			assert.Nil(t, present.MustVal())
		})
	}

	// Val preserves the legacy boxed typed nil; Get uses any's zero value instead.
	boxed := OfNullable[any](pointer)
	if boxed.Val() == nil {
		t.Fatal("Val lost the typed nil stored by OfNullable")
	}
	if val, ok := boxed.Get(); val != nil || ok {
		t.Fatalf("Get() = (%v, %v), want (nil, false)", val, ok)
	}
	present := FromOK[any](pointer, true)
	if val, ok := present.Get(); val != any(pointer) || !ok {
		t.Fatalf("Get() = (%v, %v), want (typed nil, true)", val, ok)
	}

	assert.PanicsWithValue(t, ErrNil, func() { Of(pointer) })
	assert.True(t, OfNullable(pointer).IsEmpty())
	ptrVal, ptrOK := FromOK(pointer, true).Get()
	assert.Nil(t, ptrVal)
	assert.True(t, ptrOK)
}

func TestFromOK_MapLookup(t *testing.T) {
	values := map[string]*int{"present nil": nil}
	val, ok := values["present nil"]
	present := FromOK(val, ok)
	assert.True(t, present.IsPresent())
	assert.Nil(t, present.MustVal())
	val, ok = values["missing"]
	assert.True(t, FromOK(val, ok).IsEmpty())
}
