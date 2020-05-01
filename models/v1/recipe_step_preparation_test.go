package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepPreparation{}

		expected := &RecipeStepPreparationUpdateInput{
			ValidPreparationID: fake.Uint64(),
			Notes:              fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.ValidPreparationID, i.ValidPreparationID)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}

func TestRecipeStepPreparation_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStepPreparation := &RecipeStepPreparation{
			ValidPreparationID: uint64(fake.Uint32()),
			Notes:              fake.Word(),
		}

		expected := &RecipeStepPreparationUpdateInput{
			ValidPreparationID: recipeStepPreparation.ValidPreparationID,
			Notes:              recipeStepPreparation.Notes,
		}
		actual := recipeStepPreparation.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
