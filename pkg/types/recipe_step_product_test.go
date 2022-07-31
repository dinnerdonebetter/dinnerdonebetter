package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:                 fake.LoremIpsumSentence(exampleQuantity),
			Type:                 RecipeStepProductIngredientType,
			QuantityType:         fake.LoremIpsumSentence(exampleQuantity),
			MinimumQuantityValue: fake.Float32(),
			QuantityNotes:        fake.LoremIpsumSentence(exampleQuantity),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{
			Name:                 stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Type:                 stringPointer(RecipeStepProductIngredientType),
			QuantityType:         stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantityValue: float32Pointer(fake.Float32()),
			MaximumQuantityValue: float32Pointer(fake.Float32()),
			QuantityNotes:        stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
