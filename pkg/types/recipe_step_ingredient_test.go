package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

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
		input.RecipeStepProductRecipeID = pointers.Pointer(t.Name())
		input.MaximumQuantity = pointers.Pointer(fake.Float32())
		input.Optional = pointers.Pointer(true)
		input.VesselIndex = pointers.Pointer(fake.Uint16())
		input.ProductPercentageToUse = pointers.Pointer(fake.Float32())

		x.Update(input)
	})
}

func TestRecipeStepIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			IngredientID:      pointers.Pointer(t.Name()),
			MeasurementUnitID: t.Name(),
			MinimumQuantity:   fake.Float32(),
			QuantityNotes:     t.Name(),
			IngredientNotes:   t.Name(),
			Optional:          fake.Bool(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
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
			MinimumQuantity:   fake.Float32(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepIngredientUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateRequestInput{
			IngredientID:      pointers.Pointer(t.Name()),
			MeasurementUnitID: pointers.Pointer(t.Name()),
			MinimumQuantity:   pointers.Pointer(fake.Float32()),
			QuantityNotes:     pointers.Pointer(t.Name()),
			IngredientNotes:   pointers.Pointer(t.Name()),
			Optional:          pointers.Pointer(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
