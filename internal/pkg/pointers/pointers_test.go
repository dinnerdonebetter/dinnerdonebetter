package pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := false
		actual := Pointer(expected)

		require.NotNil(t, actual)
		assert.Equal(t, expected, *actual)
	})
}
