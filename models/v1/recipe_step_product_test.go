package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepProduct{}

		expected := &RecipeStepProductUpdateInput{
			Name:         fake.Word(),
			RecipeStepID: fake.Uint64(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}

func TestRecipeStepProduct_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStepProduct := &RecipeStepProduct{
			Name:         fake.Word(),
			RecipeStepID: uint64(fake.Uint32()),
		}

		expected := &RecipeStepProductUpdateInput{
			Name:         recipeStepProduct.Name,
			RecipeStepID: recipeStepProduct.RecipeStepID,
		}
		actual := recipeStepProduct.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
