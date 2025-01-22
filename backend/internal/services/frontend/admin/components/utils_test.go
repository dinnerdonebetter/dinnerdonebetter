package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testExample struct {
	ArbitraryString string
	ArbitraryBool   bool
	ArbitraryFloat  float64
}

func TestGetFieldNames(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"ArbitraryBool",
			"ArbitraryFloat",
			"ArbitraryString",
		}
		actual := GetFieldNames[testExample]()

		assert.Equal(t, expected, actual)
	})

	T.Run("panics with a pointer", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() { GetFieldNames[*testExample]() })
	})
}

func TestGetFieldValues(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := testExample{
			ArbitraryBool:   true,
			ArbitraryFloat:  1.2,
			ArbitraryString: "hello world",
		}

		expected := []string{
			"true",
			"1.2",
			"hello world",
		}
		actual := GetFieldValues(x)

		assert.Equal(t, expected, actual)
	})

	T.Run("panics with a pointer", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() { GetFieldValues(&testExample{}) })
	})
}
