package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequiredPreparationInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RequiredPreparationInstrument{}

		expected := &RequiredPreparationInstrumentUpdateInput{
			ValidInstrumentID: fake.Uint64(),
			Notes:             fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.ValidInstrumentID, i.ValidInstrumentID)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}

func TestRequiredPreparationInstrument_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		requiredPreparationInstrument := &RequiredPreparationInstrument{
			ValidInstrumentID: uint64(fake.Uint32()),
			Notes:             fake.Word(),
		}

		expected := &RequiredPreparationInstrumentUpdateInput{
			ValidInstrumentID: requiredPreparationInstrument.ValidInstrumentID,
			Notes:             requiredPreparationInstrument.Notes,
		}
		actual := requiredPreparationInstrument.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
