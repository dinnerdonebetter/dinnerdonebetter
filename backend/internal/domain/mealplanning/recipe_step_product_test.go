package mealplanning

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/types"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProduct{
			MeasurementUnit:        &ValidMeasurementUnit{},
			ContainedInVesselIndex: new(uint16(3)),
		}
		input := &RecipeStepProductUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Compostable = new(true)
		input.ContainedInVesselIndex = new(uint16(1))
		input.IsLiquid = new(true)
		input.IsWaste = new(true)

		x.Update(input)
	})
}

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:              t.Name(),
			Type:              RecipeStepProductIngredientType,
			MeasurementUnitID: new(t.Name()),
			MeasurementQuantity: types.OptionalFloat32Range{
				Max: nil,
				Min: new(fake.Float32()),
			},
			QuantityNotes:            t.Name(),
			Compostable:              fake.Bool(),
			StorageDurationInSeconds: types.OptionalUint32Range{Max: new(fake.Uint32())},
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: new(fake.Float32()),
				Min: new(fake.Float32()),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{
			Name:                        new(t.Name()),
			Type:                        new(RecipeStepProductIngredientType),
			MeasurementUnitID:           new(t.Name()),
			MeasurementQuantity:         types.OptionalFloat32Range{Max: new(fake.Float32()), Min: new(fake.Float32())},
			QuantityNotes:               new(t.Name()),
			Compostable:                 new(fake.Bool()),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{Max: new(fake.Float32()), Min: new(fake.Float32())},
			StorageDurationInSeconds:    types.OptionalUint32Range{Max: new(fake.Uint32()), Min: new(fake.Uint32())},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
