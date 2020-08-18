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
			IngredientID:    func(x uint64) *uint64 { return &x }(fake.Uint64()),
			QuantityType:    fake.Word(),
			QuantityValue:   fake.Float32(),
			QuantityNotes:   fake.Word(),
			ProductOfRecipe: fake.Bool(),
			IngredientNotes: fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.IngredientID, i.IngredientID)
		assert.Equal(t, expected.QuantityType, i.QuantityType)
		assert.Equal(t, expected.QuantityValue, i.QuantityValue)
		assert.Equal(t, expected.QuantityNotes, i.QuantityNotes)
		assert.Equal(t, expected.ProductOfRecipe, i.ProductOfRecipe)
		assert.Equal(t, expected.IngredientNotes, i.IngredientNotes)
	})
}

func TestRecipeStepIngredient_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStepIngredient := &RecipeStepIngredient{
			IngredientID:    func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			QuantityType:    fake.Word(),
			QuantityValue:   fake.Float32(),
			QuantityNotes:   fake.Word(),
			ProductOfRecipe: fake.Bool(),
			IngredientNotes: fake.Word(),
		}

		expected := &RecipeStepIngredientUpdateInput{
			IngredientID:    recipeStepIngredient.IngredientID,
			QuantityType:    recipeStepIngredient.QuantityType,
			QuantityValue:   recipeStepIngredient.QuantityValue,
			QuantityNotes:   recipeStepIngredient.QuantityNotes,
			ProductOfRecipe: recipeStepIngredient.ProductOfRecipe,
			IngredientNotes: recipeStepIngredient.IngredientNotes,
		}
		actual := recipeStepIngredient.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
