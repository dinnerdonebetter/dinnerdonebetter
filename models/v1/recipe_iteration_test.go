package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeIteration_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeIteration{}

		expected := &RecipeIterationUpdateInput{
			RecipeID:            1,
			EndDifficultyRating: 1.23,
			EndComplexityRating: 1.23,
			EndTasteRating:      1.23,
			EndOverallRating:    1.23,
		}

		i.Update(expected)
		assert.Equal(t, expected.RecipeID, i.RecipeID)
		assert.Equal(t, expected.EndDifficultyRating, i.EndDifficultyRating)
		assert.Equal(t, expected.EndComplexityRating, i.EndComplexityRating)
		assert.Equal(t, expected.EndTasteRating, i.EndTasteRating)
		assert.Equal(t, expected.EndOverallRating, i.EndOverallRating)
	})
}
