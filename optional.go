package optional

import (
	"fmt"
	"reflect"
)

// Optional holds a value that may be absent. Its zero value is empty.
type Optional[T any] struct {
	val     T
	present bool
}

func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

// Of Returns an Optional describing the given non-null value.
// panic if value is nil.
func Of[T any](val T) Optional[T] {
	return Optional[T]{
		val:     RequireNonNull(val),
		present: true,
	}
}

// OfNullable Returns an Optional describing the given value, if non-null, otherwise returns an empty Optional.
func OfNullable[T any](val T) Optional[T] {
	return Optional[T]{
		val:     val,
		present: !IsNil(val),
	}
}

// FromOK returns a present Optional when ok is true, even if val is nil.
// When ok is false, it returns an empty Optional and discards val.
func FromOK[T any](val T, ok bool) Optional[T] {
	if !ok {
		return Empty[T]()
	}
	return Optional[T]{val: val, present: true}
}

// Get returns the value and true when present, or T's zero value and false when empty.
func (o Optional[T]) Get() (T, bool) {
	if !o.present {
		var zero T
		return zero, false
	}
	return o.val, true
}

// Val returns the stored value, even when empty. Unlike Get, it does not normalize
// an empty OfNullable's typed nil to T's zero value.
func (o Optional[T]) Val() T {
	return o.val
}

// MustVal returns the value is present, else panic
func (o Optional[T]) MustVal() T {
	if o.IsPresent() {
		return o.Val()
	}
	panic(ErrNil)
}

// IsPresent If a value is present, returns true, otherwise false.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty If a value is not present, returns true, otherwise false.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// OrElse If a value is present, returns the value, otherwise returns defaultVal.
func (o Optional[T]) OrElse(defaultVal T) T {
	if o.IsEmpty() {
		return defaultVal
	}
	return o.Val()
}

// OrElseGet If a value is present, returns the value, otherwise returns the result produced by the supplying function.
func (o Optional[T]) OrElseGet(defaultValSupplier func() T) T {
	if o.IsPresent() {
		return o.Val()
	}
	return defaultValSupplier()
}

// IfPresent If a value is present, performs the given action with the value, otherwise does nothing.
// f: the function to be performed, if a value is present
func (o Optional[T]) IfPresent(f func(T)) {
	if !o.IsEmpty() {
		f(o.Val())
	}
}

// Filter If a value is present, and the value matches the given predicate, returns an Optional describing the value, otherwise returns an empty Optional.
func (o Optional[T]) Filter(f func(T) bool) Optional[T] {
	if o.IsEmpty() {
		return o
	}
	if f(o.Val()) {
		return o
	}
	return Empty[T]()
}

func (o Optional[T]) String() string {
	if o.IsPresent() {
		return fmt.Sprintf("Optional[%s]{%v}", reflect.TypeOf(o.Val()), o.Val())
	}
	return fmt.Sprintf("EmptyOptional[%s]{%v}", reflect.TypeOf(o.Val()), o.Val())
}
