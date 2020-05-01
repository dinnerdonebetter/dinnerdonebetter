package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestIngredientTagMapping_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &IngredientTagMapping{}

		expected := &IngredientTagMappingUpdateInput{
			ValidIngredientTagID: fake.Uint64(),
		}

		i.Update(expected)
		assert.Equal(t, expected.ValidIngredientTagID, i.ValidIngredientTagID)
	})
}

func TestIngredientTagMapping_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ingredientTagMapping := &IngredientTagMapping{
			ValidIngredientTagID: uint64(fake.Uint32()),
		}

		expected := &IngredientTagMappingUpdateInput{
			ValidIngredientTagID: ingredientTagMapping.ValidIngredientTagID,
		}
		actual := ingredientTagMapping.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
