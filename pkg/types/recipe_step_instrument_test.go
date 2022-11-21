package types

import (
	"context"
	"math"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepInstrumentCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{
			InstrumentID:        stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Name:                fake.LoremIpsumSentence(exampleQuantity),
			ProductOfRecipeStep: fake.Bool(),
			RecipeStepProductID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               fake.LoremIpsumSentence(exampleQuantity),
			PreferenceRank:      uint8(fake.Number(1, math.MaxUint8)),
			Optional:            fake.Bool(),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     fake.Uint32(),
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
			InstrumentID:        stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Name:                stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToRecipeStep: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			ProductOfRecipeStep: boolPointer(fake.Bool()),
			RecipeStepProductID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:               stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			PreferenceRank:      uint8Pointer(uint8(fake.Number(1, math.MaxUint8))),
			Optional:            boolPointer(fake.Bool()),
			MinimumQuantity:     uint32Pointer(fake.Uint32()),
			MaximumQuantity:     uint32Pointer(fake.Uint32()),
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
