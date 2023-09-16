package optional

// Map If a value is present, returns an Optional describing (as if by OfNullable) the result of applying the given mapping function to the value, otherwise returns an empty Optional.
// If the mapping function returns a null result then this method returns an empty Optional.
func Map[T, E any](o Optional[T], f func(T) E) Optional[E] {
	if o.IsPresent() {
		return OfNullable(f(o.Val()))
	}
	return Empty[E]()
}

func FlatMap[T, E any](o Optional[T], f func(Optional[T]) E) Optional[E] {
	return OfNullable(f(o))
}
