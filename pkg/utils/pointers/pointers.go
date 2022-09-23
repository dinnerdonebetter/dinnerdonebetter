package pointers

// StringPointer returns a pointer to a string.
func StringPointer(s string) *string {
	return &s
}

// Uint32Pointer returns a pointer to a uint32.
func Uint32Pointer(s uint32) *uint32 {
	return &s
}

// Float32Pointer returns a pointer to a float32.
func Float32Pointer(s float32) *float32 {
	return &s
}
