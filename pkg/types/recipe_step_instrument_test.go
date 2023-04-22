package types

import (
	"context"
	"math"
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepInstrumentCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{
			InstrumentID:        pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Name:                fake.LoremIpsumSentence(exampleQuantity),
			RecipeStepProductID: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               fake.LoremIpsumSentence(exampleQuantity),
			PreferenceRank:      uint8(fake.Number(1, math.MaxUint8)),
			Optional:            fake.Bool(),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepInstrumentUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentUpdateRequestInput{
			InstrumentID:        pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Name:                pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToRecipeStep: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			RecipeStepProductID: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			PreferenceRank:      pointers.Pointer(uint8(fake.Number(1, math.MaxUint8))),
			Optional:            pointers.Pointer(fake.Bool()),
			MinimumQuantity:     pointers.Pointer(fake.Uint32()),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
