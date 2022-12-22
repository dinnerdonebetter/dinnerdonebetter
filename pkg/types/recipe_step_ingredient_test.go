package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			IngredientID:      stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MeasurementUnitID: fake.LoremIpsumSentence(exampleQuantity),
			MinimumQuantity:   fake.Float32(),
			QuantityNotes:     fake.LoremIpsumSentence(exampleQuantity),
			IngredientNotes:   fake.LoremIpsumSentence(exampleQuantity),
			Optional:          fake.Bool(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepIngredientUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateRequestInput{
			IngredientID:      stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MeasurementUnitID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:   float32Pointer(fake.Float32()),
			QuantityNotes:     stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			IngredientNotes:   stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Optional:          boolPointer(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
