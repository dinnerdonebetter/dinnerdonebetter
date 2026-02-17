package mealplanning

import (
	"testing"

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
		input.RecipeStepProductRecipeID = new(t.Name())
		input.Quantity.Max = new(fake.Float32())
		input.Optional = new(true)
		input.VesselIndex = new(fake.Uint16())
		input.ProductPercentageToUse = new(fake.Float32())

		x.Update(input)
	})
}

func TestRecipeStepIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			ValidIngredientPreparationID:     new(t.Name()),
			ValidIngredientMeasurementUnitID: new(t.Name()),
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
			ProductOfRecipeStepIndex:        new(uint64(0)),
			ProductOfRecipeStepProductIndex: new(uint64(0)),
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
			IngredientID:      new(t.Name()),
			MeasurementUnitID: new(t.Name()),
			Quantity:          types.Float32RangeWithOptionalMaxUpdateRequestInput{Min: new(fake.Float32())},
			QuantityNotes:     new(t.Name()),
			IngredientNotes:   new(t.Name()),
			Optional:          new(fake.Bool()),
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
