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
			InstrumentID:  fake.Uint64(),
			PreparationID: fake.Uint64(),
			Notes:         fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.InstrumentID, i.InstrumentID)
		assert.Equal(t, expected.PreparationID, i.PreparationID)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}

func TestRequiredPreparationInstrument_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		requiredPreparationInstrument := &RequiredPreparationInstrument{
			InstrumentID:  uint64(fake.Uint32()),
			PreparationID: uint64(fake.Uint32()),
			Notes:         fake.Word(),
		}

		expected := &RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  requiredPreparationInstrument.InstrumentID,
			PreparationID: requiredPreparationInstrument.PreparationID,
			Notes:         requiredPreparationInstrument.Notes,
		}
		actual := requiredPreparationInstrument.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
