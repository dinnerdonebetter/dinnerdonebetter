package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipe_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Recipe{}

		expected := &RecipeUpdateInput{
			Name:               fake.Word(),
			Source:             fake.Word(),
			Description:        fake.Word(),
			InspiredByRecipeID: func(x uint64) *uint64 { return &x }(fake.Uint64()),
			Private:            fake.Bool(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Source, i.Source)
		assert.Equal(t, expected.Description, i.Description)
		assert.Equal(t, expected.InspiredByRecipeID, i.InspiredByRecipeID)
		assert.Equal(t, expected.Private, i.Private)
	})
}

func TestRecipe_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipe := &Recipe{
			Name:               fake.Word(),
			Source:             fake.Word(),
			Description:        fake.Word(),
			InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			Private:            fake.Bool(),
		}

		expected := &RecipeUpdateInput{
			Name:               recipe.Name,
			Source:             recipe.Source,
			Description:        recipe.Description,
			InspiredByRecipeID: recipe.InspiredByRecipeID,
			Private:            recipe.Private,
		}
		actual := recipe.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
