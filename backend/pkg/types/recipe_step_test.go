package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func buildValidRecipeStepCreationRequestInput() *RecipeStepCreationRequestInput {
	return &RecipeStepCreationRequestInput{
		PreparationID: "PreparationID",
		EstimatedTimeInSeconds: OptionalUint32Range{
			Max: pointer.To(fake.Uint32()),
			Min: pointer.To(fake.Uint32()),
		},
		TemperatureInCelsius: OptionalFloat32Range{
			Max: nil,
			Min: pointer.To(float32(123.45)),
		},
		Notes:                "Notes",
		ExplicitInstructions: "ExplicitInstructions",
		Instruments: []*RecipeStepInstrumentCreationRequestInput{
			{
				InstrumentID: pointer.To("InstrumentID"),
				Name:         "Name",
				Quantity:     Uint32RangeWithOptionalMax{Min: fake.Uint32()},
			},
		},
		Products: []*RecipeStepProductCreationRequestInput{
			{
				Name:              "Name",
				Type:              RecipeStepProductIngredientType,
				MeasurementUnitID: pointer.To("MeasurementUnitID"),
				Quantity:          OptionalFloat32Range{Min: pointer.To(float32(1))},
				QuantityNotes:     "QuantityNotes",
			},
		},
		Ingredients: []*RecipeStepIngredientCreationRequestInput{
			{
				IngredientID:      func(s string) *string { return &s }("IngredientID"),
				MeasurementUnitID: "MeasurementUnitID",
				QuantityNotes:     "QuantityNotes",
				IngredientNotes:   "IngredientNotes",
				Quantity:          Float32RangeWithOptionalMax{Min: 1},
			},
		},
	}
}

func TestRecipeStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStep{
			TemperatureInCelsius: OptionalFloat32Range{Min: pointer.To(float32(123.45))},
		}

		input := &RecipeStepUpdateRequestInput{}
		assert.NoError(t, fake.Struct(&input))
		input.Optional = pointer.To(true)
		input.StartTimerAutomatically = pointer.To(true)
		input.TemperatureInCelsius.Min = pointer.To(float32(543.21))
		input.TemperatureInCelsius.Max = pointer.To(float32(123.45))

		x.Update(input)
	})
}

func TestRecipeStepCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := buildValidRecipeStepCreationRequestInput()

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
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
			PreparationID: t.Name(),
			EstimatedTimeInSeconds: OptionalUint32Range{
				Max: pointer.To(fake.Uint32()),
				Min: pointer.To(fake.Uint32()),
			},
			TemperatureInCelsius: OptionalFloat32Range{
				Max: nil,
				Min: pointer.To(float32(123.45)),
			},
			Notes:                t.Name(),
			ExplicitInstructions: t.Name(),
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
				Quantity:          Float32RangeWithOptionalMax{Min: 1},
			})
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepDatabaseCreationInput{
			ID:            t.Name(),
			PreparationID: t.Name(),
			Products: []*RecipeStepProductDatabaseCreationInput{
				{
					ID:                  t.Name(),
					Name:                t.Name(),
					Type:                RecipeStepProductIngredientType,
					BelongsToRecipeStep: t.Name(),
				},
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{
			Index:       pointer.To(fake.Uint32()),
			Preparation: &ValidPreparation{},
			EstimatedTimeInSeconds: OptionalUint32Range{
				Max: pointer.To(fake.Uint32()),
				Min: pointer.To(fake.Uint32()),
			},
			TemperatureInCelsius: OptionalFloat32Range{
				Max: nil,
				Min: pointer.To(float32(123.45)),
			},
			Notes:                pointer.To(t.Name()),
			ExplicitInstructions: pointer.To(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
