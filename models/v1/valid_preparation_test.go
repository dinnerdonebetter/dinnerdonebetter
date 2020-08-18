package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &ValidPreparation{}

		expected := &ValidPreparationUpdateInput{
			Name:                       fake.Word(),
			Description:                fake.Word(),
			Icon:                       fake.Word(),
			ApplicableToAllIngredients: fake.Bool(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Description, i.Description)
		assert.Equal(t, expected.Icon, i.Icon)
		assert.Equal(t, expected.ApplicableToAllIngredients, i.ApplicableToAllIngredients)
	})
}

func TestValidPreparation_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		validPreparation := &ValidPreparation{
			Name:                       fake.Word(),
			Description:                fake.Word(),
			Icon:                       fake.Word(),
			ApplicableToAllIngredients: fake.Bool(),
		}

		expected := &ValidPreparationUpdateInput{
			Name:                       validPreparation.Name,
			Description:                validPreparation.Description,
			Icon:                       validPreparation.Icon,
			ApplicableToAllIngredients: validPreparation.ApplicableToAllIngredients,
		}
		actual := validPreparation.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
