package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Instrument{}

		expected := &InstrumentUpdateInput{
			Name:        "example",
			Variant:     "example",
			Description: "example",
			Icon:        "example",
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Variant, i.Variant)
		assert.Equal(t, expected.Description, i.Description)
		assert.Equal(t, expected.Icon, i.Icon)
	})
}
