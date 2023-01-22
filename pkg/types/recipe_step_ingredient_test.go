package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func TestRecipeStepIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			IngredientID:      pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
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
			IngredientID:      pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MeasurementUnitID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:   pointers.Float32(fake.Float32()),
			QuantityNotes:     pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			IngredientNotes:   pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Optional:          pointers.Bool(fake.Bool()),
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
