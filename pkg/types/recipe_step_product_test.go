package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProduct{
			MeasurementUnit:        &ValidMeasurementUnit{},
			ContainedInVesselIndex: pointers.Pointer(uint16(3)),
		}
		input := &RecipeStepProductUpdateRequestInput{}

		fake.Struct(&input)
		input.Compostable = pointers.Pointer(true)
		input.ContainedInVesselIndex = pointers.Pointer(uint16(1))
		input.IsLiquid = pointers.Pointer(true)
		input.IsWaste = pointers.Pointer(true)

		x.Update(input)
	})
}

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:                               t.Name(),
			Type:                               RecipeStepProductIngredientType,
			MeasurementUnitID:                  pointers.Pointer(t.Name()),
			MinimumQuantity:                    pointers.Pointer(fake.Float32()),
			QuantityNotes:                      t.Name(),
			Compostable:                        fake.Bool(),
			MaximumStorageDurationInSeconds:    pointers.Pointer(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{
			Name:                               pointers.Pointer(t.Name()),
			Type:                               pointers.Pointer(RecipeStepProductIngredientType),
			MeasurementUnitID:                  pointers.Pointer(t.Name()),
			MinimumQuantity:                    pointers.Pointer(fake.Float32()),
			MaximumQuantity:                    pointers.Pointer(fake.Float32()),
			QuantityNotes:                      pointers.Pointer(t.Name()),
			Compostable:                        pointers.Pointer(fake.Bool()),
			MaximumStorageDurationInSeconds:    pointers.Pointer(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
