package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipe_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Recipe{}

		expected := &RecipeUpdateInput{
			Name:               "example",
			Source:             "example",
			Description:        "example",
			InspiredByRecipeID: func(x uint64) *uint64 { return &x }(1),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Source, i.Source)
		assert.Equal(t, expected.Description, i.Description)
		assert.Equal(t, expected.InspiredByRecipeID, i.InspiredByRecipeID)
	})
}
