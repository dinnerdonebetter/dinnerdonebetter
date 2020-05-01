package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeIteration_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeIteration{}

		expected := &RecipeIterationUpdateInput{
			EndDifficultyRating: fake.Float32(),
			EndComplexityRating: fake.Float32(),
			EndTasteRating:      fake.Float32(),
			EndOverallRating:    fake.Float32(),
		}

		i.Update(expected)
		assert.Equal(t, expected.EndDifficultyRating, i.EndDifficultyRating)
		assert.Equal(t, expected.EndComplexityRating, i.EndComplexityRating)
		assert.Equal(t, expected.EndTasteRating, i.EndTasteRating)
		assert.Equal(t, expected.EndOverallRating, i.EndOverallRating)
	})
}

func TestRecipeIteration_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeIteration := &RecipeIteration{
			EndDifficultyRating: fake.Float32(),
			EndComplexityRating: fake.Float32(),
			EndTasteRating:      fake.Float32(),
			EndOverallRating:    fake.Float32(),
		}

		expected := &RecipeIterationUpdateInput{
			EndDifficultyRating: recipeIteration.EndDifficultyRating,
			EndComplexityRating: recipeIteration.EndComplexityRating,
			EndTasteRating:      recipeIteration.EndTasteRating,
			EndOverallRating:    recipeIteration.EndOverallRating,
		}
		actual := recipeIteration.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
