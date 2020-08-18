package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepInstrument{}

		expected := &RecipeStepInstrumentUpdateInput{
			InstrumentID: func(x uint64) *uint64 { return &x }(fake.Uint64()),
			RecipeStepID: fake.Uint64(),
			Notes:        fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.InstrumentID, i.InstrumentID)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}

func TestRecipeStepInstrument_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStepInstrument := &RecipeStepInstrument{
			InstrumentID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			RecipeStepID: uint64(fake.Uint32()),
			Notes:        fake.Word(),
		}

		expected := &RecipeStepInstrumentUpdateInput{
			InstrumentID: recipeStepInstrument.InstrumentID,
			RecipeStepID: recipeStepInstrument.RecipeStepID,
			Notes:        recipeStepInstrument.Notes,
		}
		actual := recipeStepInstrument.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
