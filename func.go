package optional

import "context"

// Map If a value is present, returns an Optional describing (as if by OfNullable) the result of applying the given mapping function to the value, otherwise returns an empty Optional.
// If the mapping function returns a null result then this method returns an empty Optional.
func Map[T, E any](o Optional[T], f func(T) E) Optional[E] {
	if o.IsPresent() {
		return OfNullable(f(o.Val()))
	}
	return Empty[E]()
}

// FlatMap If a value is present, returns the result of applying the given Optional-bearing mapping function to the value, otherwise returns an empty Optional.
// This method is similar to Map(Function), but the mapping function is one whose result is already an Optional, and if invoked, flatMap does not wrap it within an additional Optional.
func FlatMap[T, E any](o Optional[T], f func(Optional[T]) E) Optional[E] {
	if o.IsPresent() {
		return OfNullable(f(o))
	}
	return Empty[E]()
}

// MapWithCtx If a value is present, returns an Optional describing (as if by OfNullable) the result of applying the given mapping function to the value, otherwise returns an empty Optional.
// If the mapping function returns a null result then this method returns an empty Optional.
func MapWithCtx[T, E any](ctx context.Context, o Optional[T], f func(context.Context, T) E) Optional[E] {
	if o.IsPresent() {
		return OfNullable(f(ctx, o.Val()))
	}
	return Empty[E]()
}

// FlatMapWithCtx If a value is present, returns the result of applying the given Optional-bearing mapping function to the value, otherwise returns an empty Optional.
// This method is similar to MapWithCtx(Function), but the mapping function is one whose result is already an Optional, and if invoked, flatMap does not wrap it within an additional Optional.
func FlatMapWithCtx[T, E any](ctx context.Context, o Optional[T], f func(context.Context, Optional[T]) E) Optional[E] {
	if o.IsPresent() {
		return OfNullable(f(ctx, o))
	}
	return Empty[E]()
}

// MapWithErr If the given err is not null, returns an empty Optional and this err.
// If a value is present, returns an Optional describing (as if by OfNullable) the result of applying the given mapping function to the value, otherwise returns an empty Optional.
// If the mapping function returns a null result then this method returns an empty Optional.
func MapWithErr[T, E any](o Optional[T], f func(T) (E, error), err error) (Optional[E], error) {
	if err != nil {
		return Empty[E](), err
	}
	if o.IsPresent() {
		if r, err := f(o.Val()); err != nil {
			return Empty[E](), err
		} else {
			return OfNullable(r), nil
		}
	}
	return Empty[E](), nil
}

// FlatMapWithErr If the given err is not null, returns an empty Optional and this err.
// If a value is present, returns the result of applying the given Optional-bearing mapping function to the value, otherwise returns an empty Optional.
// This method is similar to MapWithErr(Function), but the mapping function is one whose result is already an Optional, and if invoked, flatMap does not wrap it within an additional Optional.
func FlatMapWithErr[T, E any](o Optional[T], f func(Optional[T]) (E, error), err error) (Optional[E], error) {
	if err != nil {
		return Empty[E](), err
	}
	if o.IsPresent() {
		if r, err := f(o); err != nil {
			return Empty[E](), err
		} else {
			return OfNullable(r), nil
		}
	}
	return Empty[E](), nil
}

// MapWithCtxAndErr If the given err is not null, returns an empty Optional and this err.
// If a value is present, returns an Optional describing (as if by OfNullable) the result of applying the given mapping function to the value, otherwise returns an empty Optional.
// If the mapping function returns a null result then this method returns an empty Optional.
func MapWithCtxAndErr[T, E any](ctx context.Context, o Optional[T], f func(context.Context, T) (E, error), err error) (Optional[E], error) {
	if err != nil {
		return Empty[E](), err
	}
	if o.IsPresent() {
		if r, err := f(ctx, o.Val()); err != nil {
			return Empty[E](), err
		} else {
			return OfNullable(r), nil
		}
	}
	return Empty[E](), nil
}

// FlatMapWithCtxAndErr If the given err is not null, returns an empty Optional and this err.
// If a value is present, returns the result of applying the given Optional-bearing mapping function to the value, otherwise returns an empty Optional.
// This method is similar to MapWithCtxAndErr(Function), but the mapping function is one whose result is already an Optional, and if invoked, flatMap does not wrap it within an additional Optional.
func FlatMapWithCtxAndErr[T, E any](ctx context.Context, o Optional[T], f func(context.Context, Optional[T]) (E, error), err error) (Optional[E], error) {
	if err != nil {
		return Empty[E](), err
	}
	if o.IsPresent() {
		if r, err := f(ctx, o); err != nil {
			return Empty[E](), err
		} else {
			return OfNullable(r), nil
		}
	}
	return Empty[E](), nil
}
