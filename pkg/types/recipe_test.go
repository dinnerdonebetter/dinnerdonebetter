package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{
			Name:               fake.LoremIpsumSentence(exampleQuantity),
			Source:             fake.LoremIpsumSentence(exampleQuantity),
			Description:        fake.LoremIpsumSentence(exampleQuantity),
			InspiredByRecipeID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Steps: []*RecipeStepCreationRequestInput{
				buildValidRecipeStepCreationRequestInput(),
			},
			SealOfApproval: fake.Bool(),
			YieldsPortions: fake.Uint8(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateRequestInput{
			Name:               stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Source:             stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Description:        stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			InspiredByRecipeID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			SealOfApproval:     boolPointer(fake.Bool()),
			YieldsPortions:     uint8Pointer(fake.Uint8()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
