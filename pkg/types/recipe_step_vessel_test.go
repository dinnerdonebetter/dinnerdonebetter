package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			Name:                fake.LoremIpsumSentence(exampleQuantity),
			RecipeStepProductID: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               fake.LoremIpsumSentence(exampleQuantity),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
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
			Name:                pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToRecipeStep: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			RecipeStepProductID: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:     pointers.Pointer(fake.Uint32()),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
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
