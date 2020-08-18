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
			Notes:              fake.Word(),
			ValidPreparationID: fake.Uint64(),
			ValidIngredientID:  fake.Uint64(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Notes, i.Notes)
		assert.Equal(t, expected.ValidPreparationID, i.ValidPreparationID)
		assert.Equal(t, expected.ValidIngredientID, i.ValidIngredientID)
	})
}

func TestValidIngredientPreparation_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		validIngredientPreparation := &ValidIngredientPreparation{
			Notes:              fake.Word(),
			ValidPreparationID: uint64(fake.Uint32()),
			ValidIngredientID:  uint64(fake.Uint32()),
		}

		expected := &ValidIngredientPreparationUpdateInput{
			Notes:              validIngredientPreparation.Notes,
			ValidPreparationID: validIngredientPreparation.ValidPreparationID,
			ValidIngredientID:  validIngredientPreparation.ValidIngredientID,
		}
		actual := validIngredientPreparation.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
