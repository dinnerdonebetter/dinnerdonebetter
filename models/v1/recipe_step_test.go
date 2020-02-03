package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStep{}

		expected := &RecipeStepUpdateInput{
			Index:                     1,
			PreparationID:             1,
			PrerequisiteStep:          1,
			MinEstimatedTimeInSeconds: 1,
			MaxEstimatedTimeInSeconds: 1,
			TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(1),
			Notes:                     "example",
			RecipeID:                  1,
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
