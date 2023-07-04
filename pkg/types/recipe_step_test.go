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
		PreparationID:                 fake.LoremIpsumSentence(exampleQuantity),
		MinimumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
		MaximumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
		MinimumTemperatureInCelsius:   pointers.Pointer(float32(123.45)),
		Notes:                         fake.LoremIpsumSentence(exampleQuantity),
		ExplicitInstructions:          fake.LoremIpsumSentence(exampleQuantity),
		Instruments: []*RecipeStepInstrumentCreationRequestInput{
			{
				InstrumentID:    pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
				Name:            fake.LoremIpsumSentence(exampleQuantity),
				MinimumQuantity: fake.Uint32(),
			},
		},
		Products: []*RecipeStepProductCreationRequestInput{
			{
				Name:              fake.LoremIpsumSentence(exampleQuantity),
				Type:              RecipeStepProductIngredientType,
				MeasurementUnitID: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
				MinimumQuantity:   pointers.Pointer(fake.Float32()),
				QuantityNotes:     fake.LoremIpsumSentence(exampleQuantity),
			},
		},
		Ingredients: []*RecipeStepIngredientCreationRequestInput{
			{
				IngredientID:      func(s string) *string { return &s }(fake.LoremIpsumSentence(exampleQuantity)),
				MeasurementUnitID: fake.LoremIpsumSentence(exampleQuantity),
				QuantityNotes:     fake.LoremIpsumSentence(exampleQuantity),
				IngredientNotes:   fake.LoremIpsumSentence(exampleQuantity),
				MinimumQuantity:   1,
			},
		},
	}
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
			PreparationID:                 fake.LoremIpsumSentence(exampleQuantity),
			MinimumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: pointers.Pointer(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Pointer(float32(123.45)),
			Notes:                         fake.LoremIpsumSentence(exampleQuantity),
			ExplicitInstructions:          fake.LoremIpsumSentence(exampleQuantity),
			Products: []*RecipeStepProductCreationRequestInput{
				{
					Name: fake.LoremIpsumSentence(exampleQuantity),
				},
			},
			Ingredients: []*RecipeStepIngredientCreationRequestInput{},
		}

		for i := 0; i < maxIngredientsPerStep*2; i++ {
			x.Ingredients = append(x.Ingredients, &RecipeStepIngredientCreationRequestInput{
				IngredientID:      func(s string) *string { return &s }(fake.LoremIpsumSentence(exampleQuantity)),
				MeasurementUnitID: fake.LoremIpsumSentence(exampleQuantity),
				QuantityNotes:     fake.LoremIpsumSentence(exampleQuantity),
				IngredientNotes:   fake.LoremIpsumSentence(exampleQuantity),
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
			Notes:                         pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			ExplicitInstructions:          pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
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
