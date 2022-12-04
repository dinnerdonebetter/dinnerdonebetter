package pointers

import (
	"time"
)

// Pointer returns a pointer to a value.
func Pointer[T any](x T) *T {
	return &x
}

// String returns a pointer to a string.
func String(x string) *string {
	return Pointer[string](x)
}

// Uint8 returns a pointer to the provided uint32.
func Uint8(x uint8) *uint8 {
	return Pointer[uint8](x)
}

// Uint16 returns a pointer to a uint16.
func Uint16(x uint16) *uint16 {
	return Pointer[uint16](x)
}

// Uint32 returns a pointer to a uint32.
func Uint32(x uint32) *uint32 {
	return Pointer[uint32](x)
}

// Uint64 returns a pointer to the provided uint64.
func Uint64(x uint64) *uint64 {
	return Pointer[uint64](x)
}

// Int32 returns a pointer to a uint32.
func Int32(x int32) *int32 {
	return Pointer[int32](x)
}

// Bool returns a pointer to the provided bool.
func Bool(x bool) *bool {
	return Pointer[bool](x)
}

// Float32 returns a pointer to the provided float32.
func Float32(x float32) *float32 {
	return Pointer[float32](x)
}

// Float64 returns a pointer to the provided float64.
func Float64(x float64) *float64 {
	return Pointer[float64](x)
}

// Time returns a pointer to the provided float64.
func Time(x time.Time) *time.Time {
	return Pointer[time.Time](x)
}
