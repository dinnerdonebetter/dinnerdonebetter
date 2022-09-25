package pointers

// StringPointer returns a pointer to a string.
func StringPointer(s string) *string {
	return &s
}

// Uint8Pointer returns a pointer to the provided uint32.
func Uint8Pointer(x uint8) *uint8 {
	return &x
}

// Uint32Pointer returns a pointer to a uint32.
func Uint32Pointer(s uint32) *uint32 {
	return &s
}

// Uint64Pointer returns a pointer to the provided uint64.
func Uint64Pointer(x uint64) *uint64 {
	return &x
}

// BoolPointer returns a pointer to the provided bool.
func BoolPointer(x bool) *bool {
	return &x
}

// Float32Pointer returns a pointer to the provided float32.
func Float32Pointer(s float32) *float32 {
	return &s
}

// Float64Pointer returns a pointer to the provided float64.
func Float64Pointer(s float64) *float64 {
	return &s
}
