package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepProduct{}

		expected := &RecipeStepProductUpdateInput{
			Name:         "example",
			RecipeStepID: 1,
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}
