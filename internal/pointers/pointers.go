package pointers

import (
	"time"
)

// String returns a pointer to a string.
func String(x string) *string {
	return &x
}

// Uint8 returns a pointer to the provided uint32.
func Uint8(x uint8) *uint8 {
	return &x
}

// Uint32 returns a pointer to a uint32.
func Uint32(x uint32) *uint32 {
	return &x
}

// Uint64 returns a pointer to the provided uint64.
func Uint64(x uint64) *uint64 {
	return &x
}

// Bool returns a pointer to the provided bool.
func Bool(x bool) *bool {
	return &x
}

// Float32 returns a pointer to the provided float32.
func Float32(x float32) *float32 {
	return &x
}

// Float64 returns a pointer to the provided float64.
func Float64(x float64) *float64 {
	return &x
}

// Time returns a pointer to the provided float64.
func Time(x time.Time) *time.Time {
	return &x
}
