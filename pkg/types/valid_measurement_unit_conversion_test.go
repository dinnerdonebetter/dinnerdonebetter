package types

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
)

func TestValidMeasurementUnitConversion_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitConversion{
			OnlyForIngredient: &ValidIngredient{},
		}
		input := &ValidMeasurementUnitConversionUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}
