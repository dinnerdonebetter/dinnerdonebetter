package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeStepInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepInstrument{}

		expected := &RecipeStepInstrumentUpdateInput{
			InstrumentID: func(x uint64) *uint64 { return &x }(1),
			RecipeStepID: 1,
			Notes:        "example",
		}

		i.Update(expected)
		assert.Equal(t, expected.InstrumentID, i.InstrumentID)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
		assert.Equal(t, expected.Notes, i.Notes)
	})
}
