package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationRequestInput{
			IngredientID:      pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
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
			IngredientID:      pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MeasurementUnitID: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:   pointers.Pointer(fake.Float32()),
			QuantityNotes:     pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			IngredientNotes:   pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Optional:          pointers.Pointer(fake.Bool()),
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
