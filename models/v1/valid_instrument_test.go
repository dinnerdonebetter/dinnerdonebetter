package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &ValidInstrument{}

		expected := &ValidInstrumentUpdateInput{
			Name:        fake.Word(),
			Variant:     fake.Word(),
			Description: fake.Word(),
			Icon:        fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Variant, i.Variant)
		assert.Equal(t, expected.Description, i.Description)
		assert.Equal(t, expected.Icon, i.Icon)
	})
}

func TestValidInstrument_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		validInstrument := &ValidInstrument{
			Name:        fake.Word(),
			Variant:     fake.Word(),
			Description: fake.Word(),
			Icon:        fake.Word(),
		}

		expected := &ValidInstrumentUpdateInput{
			Name:        validInstrument.Name,
			Variant:     validInstrument.Variant,
			Description: validInstrument.Description,
			Icon:        validInstrument.Icon,
		}
		actual := validInstrument.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
