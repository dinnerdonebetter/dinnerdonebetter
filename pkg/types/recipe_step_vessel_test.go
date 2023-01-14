package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			Name:                fake.LoremIpsumSentence(exampleQuantity),
			RecipeStepProductID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               fake.LoremIpsumSentence(exampleQuantity),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     fake.Uint32(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepVesselUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselUpdateRequestInput{
			Name:                stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToRecipeStep: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			RecipeStepProductID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:     uint32Pointer(fake.Uint32()),
			MaximumQuantity:     uint32Pointer(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
