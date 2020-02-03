package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeStepIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepIngredient{}

		expected := &RecipeStepIngredientUpdateInput{
			IngredientID:    func(x uint64) *uint64 { return &x }(1),
			QuantityType:    "example",
			QuantityValue:   1.23,
			QuantityNotes:   "example",
			ProductOfRecipe: false,
			IngredientNotes: "example",
			RecipeStepID:    1,
		}

		i.Update(expected)
		assert.Equal(t, expected.IngredientID, i.IngredientID)
		assert.Equal(t, expected.QuantityType, i.QuantityType)
		assert.Equal(t, expected.QuantityValue, i.QuantityValue)
		assert.Equal(t, expected.QuantityNotes, i.QuantityNotes)
		assert.Equal(t, expected.ProductOfRecipe, i.ProductOfRecipe)
		assert.Equal(t, expected.IngredientNotes, i.IngredientNotes)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}
