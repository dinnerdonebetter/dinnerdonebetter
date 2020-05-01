package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestIterationMedia_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &IterationMedia{}

		expected := &IterationMediaUpdateInput{
			Source:   fake.Word(),
			Mimetype: fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Source, i.Source)
		assert.Equal(t, expected.Mimetype, i.Mimetype)
	})
}

func TestIterationMedia_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		iterationMedia := &IterationMedia{
			Source:   fake.Word(),
			Mimetype: fake.Word(),
		}

		expected := &IterationMediaUpdateInput{
			Source:   iterationMedia.Source,
			Mimetype: iterationMedia.Mimetype,
		}
		actual := iterationMedia.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
