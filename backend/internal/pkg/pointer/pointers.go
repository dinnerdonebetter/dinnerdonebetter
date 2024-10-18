package pointer

// To returns a pointer to a value.
func To[T any](x T) *T {
	return &x
}

// Dereference returns the value of a pointer.
func Dereference[T any](x *T) T {
	if x == nil {
		var zero T
		return zero
	}
	return *x
}
