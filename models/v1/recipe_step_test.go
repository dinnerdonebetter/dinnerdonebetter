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
			ValidPreparationID:        fake.Uint64(),
			PrerequisiteStepID:        func(x uint64) *uint64 { return &x }(fake.Uint64()),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			YieldsProductName:         fake.Word(),
			YieldsQuantity:            uint(fake.Uint32()),
			Notes:                     fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Index, i.Index)
		assert.Equal(t, expected.ValidPreparationID, i.ValidPreparationID)
		assert.Equal(t, expected.PrerequisiteStepID, i.PrerequisiteStepID)
		assert.Equal(t, expected.MinEstimatedTimeInSeconds, i.MinEstimatedTimeInSeconds)
		assert.Equal(t, expected.MaxEstimatedTimeInSeconds, i.MaxEstimatedTimeInSeconds)
		assert.Equal(t, expected.YieldsProductName, i.YieldsProductName)
		assert.Equal(t, expected.YieldsQuantity, i.YieldsQuantity)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}

func TestRecipeStep_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStep := &RecipeStep{
			Index:                     uint(fake.Uint32()),
			ValidPreparationID:        uint64(fake.Uint32()),
			PrerequisiteStepID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			YieldsProductName:         fake.Word(),
			YieldsQuantity:            uint(fake.Uint32()),
			Notes:                     fake.Word(),
		}

		expected := &RecipeStepUpdateInput{
			Index:                     recipeStep.Index,
			ValidPreparationID:        recipeStep.ValidPreparationID,
			PrerequisiteStepID:        recipeStep.PrerequisiteStepID,
			MinEstimatedTimeInSeconds: recipeStep.MinEstimatedTimeInSeconds,
			MaxEstimatedTimeInSeconds: recipeStep.MaxEstimatedTimeInSeconds,
			YieldsProductName:         recipeStep.YieldsProductName,
			YieldsQuantity:            recipeStep.YieldsQuantity,
			Notes:                     recipeStep.Notes,
		}
		actual := recipeStep.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
