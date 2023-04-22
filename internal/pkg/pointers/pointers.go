package pointers

// Pointer returns a pointer to a value.
func Pointer[T any](x T) *T {
	return &x
}
