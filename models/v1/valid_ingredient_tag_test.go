package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientTag_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &ValidIngredientTag{}

		expected := &ValidIngredientTagUpdateInput{
			Name: fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
	})
}

func TestValidIngredientTag_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		validIngredientTag := &ValidIngredientTag{
			Name: fake.Word(),
		}

		expected := &ValidIngredientTagUpdateInput{
			Name: validIngredientTag.Name,
		}
		actual := validIngredientTag.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
