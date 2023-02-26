package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func TestRecipeStepVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			Name:                fake.LoremIpsumSentence(exampleQuantity),
			RecipeStepProductID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               fake.LoremIpsumSentence(exampleQuantity),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     pointers.Uint32(fake.Uint32()),
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
			Name:                pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToRecipeStep: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			RecipeStepProductID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:     pointers.Uint32(fake.Uint32()),
			MaximumQuantity:     pointers.Uint32(fake.Uint32()),
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
