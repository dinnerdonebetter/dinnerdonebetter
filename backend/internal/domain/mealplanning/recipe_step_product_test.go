package mealplanning

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProduct{
			MeasurementUnit:        &recipeenums.ValidMeasurementUnit{},
			ContainedInVesselIndex: pointer.To(uint16(3)),
		}
		input := &RecipeStepProductUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Compostable = pointer.To(true)
		input.ContainedInVesselIndex = pointer.To(uint16(1))
		input.IsLiquid = pointer.To(true)
		input.IsWaste = pointer.To(true)

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
			MeasurementUnitID: pointer.To(t.Name()),
			Quantity: types.OptionalFloat32Range{
				Max: nil,
				Min: pointer.To(fake.Float32()),
			},
			QuantityNotes:            t.Name(),
			Compostable:              fake.Bool(),
			StorageDurationInSeconds: types.OptionalUint32Range{Max: pointer.To(fake.Uint32())},
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: pointer.To(fake.Float32()),
				Min: pointer.To(fake.Float32()),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
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
			Name:                        pointer.To(t.Name()),
			Type:                        pointer.To(RecipeStepProductIngredientType),
			MeasurementUnitID:           pointer.To(t.Name()),
			Quantity:                    types.OptionalFloat32Range{Max: pointer.To(fake.Float32()), Min: pointer.To(fake.Float32())},
			QuantityNotes:               pointer.To(t.Name()),
			Compostable:                 pointer.To(fake.Bool()),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{Max: pointer.To(fake.Float32()), Min: pointer.To(fake.Float32())},
			StorageDurationInSeconds:    types.OptionalUint32Range{Max: pointer.To(fake.Uint32()), Min: pointer.To(fake.Uint32())},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
