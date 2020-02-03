package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredPreparationInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RequiredPreparationInstrument{}

		expected := &RequiredPreparationInstrumentUpdateInput{
			InstrumentID:  1,
			PreparationID: 1,
			Notes:         "example",
		}

		i.Update(expected)
		assert.Equal(t, expected.InstrumentID, i.InstrumentID)
		assert.Equal(t, expected.PreparationID, i.PreparationID)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}
