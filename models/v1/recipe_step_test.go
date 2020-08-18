package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStep{}

		expected := &RecipeStepUpdateInput{
			Index:                     uint(fake.Uint32()),
			PreparationID:             fake.Uint64(),
			PrerequisiteStep:          fake.Uint64(),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
			Notes:                     fake.Word(),
			RecipeID:                  fake.Uint64(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Index, i.Index)
		assert.Equal(t, expected.PreparationID, i.PreparationID)
		assert.Equal(t, expected.PrerequisiteStep, i.PrerequisiteStep)
		assert.Equal(t, expected.MinEstimatedTimeInSeconds, i.MinEstimatedTimeInSeconds)
		assert.Equal(t, expected.MaxEstimatedTimeInSeconds, i.MaxEstimatedTimeInSeconds)
		assert.Equal(t, expected.TemperatureInCelsius, i.TemperatureInCelsius)
		assert.Equal(t, expected.Notes, i.Notes)
		assert.Equal(t, expected.RecipeID, i.RecipeID)
	})
}

func TestRecipeStep_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStep := &RecipeStep{
			Index:                     uint(fake.Uint32()),
			PreparationID:             uint64(fake.Uint32()),
			PrerequisiteStep:          uint64(fake.Uint32()),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
			Notes:                     fake.Word(),
			RecipeID:                  uint64(fake.Uint32()),
		}

		expected := &RecipeStepUpdateInput{
			Index:                     recipeStep.Index,
			PreparationID:             recipeStep.PreparationID,
			PrerequisiteStep:          recipeStep.PrerequisiteStep,
			MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
			MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
			TemperatureInCelsius:      recipeStep.TemperatureInCelsius,
			Notes:                     recipeStep.Notes,
			RecipeID:                  recipeStep.RecipeID,
		}
		actual := recipeStep.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
