package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &ValidIngredientPreparation{}

		expected := &ValidIngredientPreparationUpdateInput{
			Notes: fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}

func TestValidIngredientPreparation_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		validIngredientPreparation := &ValidIngredientPreparation{
			Notes: fake.Word(),
		}

		expected := &ValidIngredientPreparationUpdateInput{
			Notes: validIngredientPreparation.Notes,
		}
		actual := validIngredientPreparation.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
