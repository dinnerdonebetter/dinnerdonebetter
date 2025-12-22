package mealplanning

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredient{}
		input := &RecipeStepIngredientUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.RecipeStepProductRecipeID = pointer.To(t.Name())
		input.Quantity.Max = pointer.To(fake.Float32())
		input.Optional = pointer.To(true)
		input.VesselIndex = pointer.To(fake.Uint16())
		input.ProductPercentageToUse = pointer.To(fake.Float32())

		x.Update(input)
	})
}

func TestRecipeStepIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			ValidIngredientPreparationID:     pointer.To(t.Name()),
			ValidIngredientMeasurementUnitID: pointer.To(t.Name()),
			Quantity:                         types.Float32RangeWithOptionalMax{Min: fake.Float32()},
			QuantityNotes:                    t.Name(),
			IngredientNotes:                  t.Name(),
			Optional:                         fake.Bool(),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("recipe step product does not require bridge IDs", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			ProductOfRecipeStepIndex:        pointer.To(uint64(0)),
			ProductOfRecipeStepProductIndex: pointer.To(uint64(0)),
			Quantity:                        types.Float32RangeWithOptionalMax{Min: fake.Float32()},
			QuantityNotes:                   t.Name(),
			IngredientNotes:                 t.Name(),
			Optional:                        fake.Bool(),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepIngredientDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientDatabaseCreationInput{
			ID:                t.Name(),
			MeasurementUnitID: t.Name(),
			Quantity:          types.Float32RangeWithOptionalMax{Min: fake.Float32()},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientDatabaseCreationInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepIngredientUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateRequestInput{
			IngredientID:      pointer.To(t.Name()),
			MeasurementUnitID: pointer.To(t.Name()),
			Quantity:          types.Float32RangeWithOptionalMaxUpdateRequestInput{Min: pointer.To(fake.Float32())},
			QuantityNotes:     pointer.To(t.Name()),
			IngredientNotes:   pointer.To(t.Name()),
			Optional:          pointer.To(fake.Bool()),
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
