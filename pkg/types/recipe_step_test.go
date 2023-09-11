package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func buildValidRecipeStepCreationRequestInput() *RecipeStepCreationRequestInput {
	return &RecipeStepCreationRequestInput{
		PreparationID:                 "PreparationID",
		MinimumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
		MaximumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
		MinimumTemperatureInCelsius:   pointers.Pointer(float32(123.45)),
		Notes:                         "Notes",
		ExplicitInstructions:          "ExplicitInstructions",
		Instruments: []*RecipeStepInstrumentCreationRequestInput{
			{
				InstrumentID:    pointers.Pointer("InstrumentID"),
				Name:            "Name",
				MinimumQuantity: fake.Uint32(),
			},
		},
		Products: []*RecipeStepProductCreationRequestInput{
			{
				Name:              "Name",
				Type:              RecipeStepProductIngredientType,
				MeasurementUnitID: pointers.Pointer("MeasurementUnitID"),
				MinimumQuantity:   pointers.Pointer(fake.Float32()),
				QuantityNotes:     "QuantityNotes",
			},
		},
		Ingredients: []*RecipeStepIngredientCreationRequestInput{
			{
				IngredientID:      func(s string) *string { return &s }("IngredientID"),
				MeasurementUnitID: "MeasurementUnitID",
				QuantityNotes:     "QuantityNotes",
				IngredientNotes:   "IngredientNotes",
				MinimumQuantity:   1,
			},
		},
	}
}

func TestRecipeStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStep{
			MinimumTemperatureInCelsius: pointers.Pointer(float32(123.45)),
		}
		input := &RecipeStepUpdateRequestInput{}

		fake.Struct(&input)
		input.Optional = pointers.Pointer(true)
		input.StartTimerAutomatically = pointers.Pointer(true)
		input.MinimumTemperatureInCelsius = pointers.Pointer(float32(543.21))
		input.MaximumTemperatureInCelsius = pointers.Pointer(float32(123.45))

		x.Update(input)
	})
}

func TestRecipeStepCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := buildValidRecipeStepCreationRequestInput()

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})

	T.Run("with too many ingredients", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationRequestInput{
			PreparationID:                 t.Name(),
			MinimumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Pointer(float32(123.45)),
			Notes:                         t.Name(),
			ExplicitInstructions:          t.Name(),
			Products: []*RecipeStepProductCreationRequestInput{
				{
					Name: t.Name(),
				},
			},
			Ingredients: []*RecipeStepIngredientCreationRequestInput{},
		}

		for i := 0; i < maxIngredientsPerStep*2; i++ {
			x.Ingredients = append(x.Ingredients, &RecipeStepIngredientCreationRequestInput{
				IngredientID:      func(s string) *string { return &s }(t.Name()),
				MeasurementUnitID: t.Name(),
				QuantityNotes:     t.Name(),
				IngredientNotes:   t.Name(),
				MinimumQuantity:   1,
			})
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{
			Index:                         pointers.Pointer(fake.Uint32()),
			Preparation:                   &ValidPreparation{},
			MinimumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Pointer(float32(123.45)),
			Notes:                         pointers.Pointer(t.Name()),
			ExplicitInstructions:          pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
