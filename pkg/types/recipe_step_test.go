package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func buildValidRecipeStepCreationRequestInput() *RecipeStepCreationRequestInput {
	return &RecipeStepCreationRequestInput{
		Index:                         fake.Uint32(),
		PreparationID:                 fake.LoremIpsumSentence(exampleQuantity),
		MinimumEstimatedTimeInSeconds: pointers.Uint32(fake.Uint32()),
		MaximumEstimatedTimeInSeconds: pointers.Uint32(fake.Uint32()),
		MinimumTemperatureInCelsius:   pointers.Float32(float32(123.45)),
		Notes:                         fake.LoremIpsumSentence(exampleQuantity),
		ExplicitInstructions:          fake.LoremIpsumSentence(exampleQuantity),
		Instruments: []*RecipeStepInstrumentCreationRequestInput{
			{
				Name:            fake.LoremIpsumSentence(exampleQuantity),
				MinimumQuantity: fake.Uint32(),
			},
		},
		Products: []*RecipeStepProductCreationRequestInput{
			{
				Name:              fake.LoremIpsumSentence(exampleQuantity),
				Type:              RecipeStepProductIngredientType,
				MeasurementUnitID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
				MinimumQuantity:   pointers.Float32(fake.Float32()),
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
			Index:                         fake.Uint32(),
			PreparationID:                 fake.LoremIpsumSentence(exampleQuantity),
			MinimumEstimatedTimeInSeconds: pointers.Uint32(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: pointers.Uint32(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Float32(float32(123.45)),
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
			Index:                         pointers.Uint32(fake.Uint32()),
			Preparation:                   &ValidPreparation{},
			MinimumEstimatedTimeInSeconds: pointers.Uint32(fake.Uint32()),
			MaximumEstimatedTimeInSeconds: pointers.Uint32(fake.Uint32()),
			MinimumTemperatureInCelsius:   pointers.Float32(float32(123.45)),
			Notes:                         pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			ExplicitInstructions:          pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
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
