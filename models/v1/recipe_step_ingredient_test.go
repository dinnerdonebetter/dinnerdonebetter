package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepIngredient{}

		expected := &RecipeStepIngredientUpdateInput{
			ValidIngredientID:     fake.Uint64(),
			IngredientNotes:       fake.Word(),
			QuantityType:          fake.Word(),
			QuantityValue:         fake.Float32(),
			QuantityNotes:         fake.Word(),
			ProductOfRecipeStepID: func(x uint64) *uint64 { return &x }(fake.Uint64()),
		}

		i.Update(expected)
		assert.Equal(t, expected.ValidIngredientID, i.ValidIngredientID)
		assert.Equal(t, expected.IngredientNotes, i.IngredientNotes)
		assert.Equal(t, expected.QuantityType, i.QuantityType)
		assert.Equal(t, expected.QuantityValue, i.QuantityValue)
		assert.Equal(t, expected.QuantityNotes, i.QuantityNotes)
		assert.Equal(t, expected.ProductOfRecipeStepID, i.ProductOfRecipeStepID)
	})
}

func TestRecipeStepIngredient_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStepIngredient := &RecipeStepIngredient{
			ValidIngredientID:     uint64(fake.Uint32()),
			IngredientNotes:       fake.Word(),
			QuantityType:          fake.Word(),
			QuantityValue:         fake.Float32(),
			QuantityNotes:         fake.Word(),
			ProductOfRecipeStepID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		}

		expected := &RecipeStepIngredientUpdateInput{
			ValidIngredientID:     recipeStepIngredient.ValidIngredientID,
			IngredientNotes:       recipeStepIngredient.IngredientNotes,
			QuantityType:          recipeStepIngredient.QuantityType,
			QuantityValue:         recipeStepIngredient.QuantityValue,
			QuantityNotes:         recipeStepIngredient.QuantityNotes,
			ProductOfRecipeStepID: recipeStepIngredient.ProductOfRecipeStepID,
		}
		actual := recipeStepIngredient.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
