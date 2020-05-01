package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeTag_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeTag{}

		expected := &RecipeTagUpdateInput{
			Name: fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
	})
}

func TestRecipeTag_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeTag := &RecipeTag{
			Name: fake.Word(),
		}

		expected := &RecipeTagUpdateInput{
			Name: recipeTag.Name,
		}
		actual := recipeTag.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
