package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func TestRecipeCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{
			Name:               fake.LoremIpsumSentence(exampleQuantity),
			Source:             fake.LoremIpsumSentence(exampleQuantity),
			Description:        fake.LoremIpsumSentence(exampleQuantity),
			InspiredByRecipeID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
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
			Name:               pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Source:             pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Description:        pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			InspiredByRecipeID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			SealOfApproval:     pointers.Bool(fake.Bool()),
			YieldsPortions:     pointers.Uint8(fake.Uint8()),
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
