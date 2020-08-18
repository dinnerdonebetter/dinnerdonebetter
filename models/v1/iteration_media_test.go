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
			Path:              fake.Word(),
			Mimetype:          fake.Word(),
			RecipeIterationID: fake.Uint64(),
			RecipeStepID:      func(x uint64) *uint64 { return &x }(fake.Uint64()),
		}

		i.Update(expected)
		assert.Equal(t, expected.Path, i.Path)
		assert.Equal(t, expected.Mimetype, i.Mimetype)
		assert.Equal(t, expected.RecipeIterationID, i.RecipeIterationID)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}

func TestIterationMedia_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		iterationMedia := &IterationMedia{
			Path:              fake.Word(),
			Mimetype:          fake.Word(),
			RecipeIterationID: uint64(fake.Uint32()),
			RecipeStepID:      func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		}

		expected := &IterationMediaUpdateInput{
			Path:              iterationMedia.Path,
			Mimetype:          iterationMedia.Mimetype,
			RecipeIterationID: iterationMedia.RecipeIterationID,
			RecipeStepID:      iterationMedia.RecipeStepID,
		}
		actual := iterationMedia.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
