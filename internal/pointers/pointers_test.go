package pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoolPointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := false
		actual := Bool(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func TestFloat32Pointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var expected float32 = 123.0
		actual := Float32(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func TestFloat64Pointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var expected float64 = 123
		actual := Float64(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func TestStringPointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := t.Name()
		actual := String(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func TestUint32Pointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var expected uint32 = 123
		actual := Uint32(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func TestUint64Pointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var expected uint64 = 123
		actual := Uint64(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}

func TestUint8Pointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var expected uint8 = 123
		actual := Uint8(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}
