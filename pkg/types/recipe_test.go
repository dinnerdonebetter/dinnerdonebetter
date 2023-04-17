package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"

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
			InspiredByRecipeID: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Steps: []*RecipeStepCreationRequestInput{
				buildValidRecipeStepCreationRequestInput(),
				buildValidRecipeStepCreationRequestInput(),
			},
			SealOfApproval:           fake.Bool(),
			MinimumEstimatedPortions: fake.Float32(),
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
			Name:                     pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Source:                   pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Description:              pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			InspiredByRecipeID:       pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			SealOfApproval:           pointers.Bool(fake.Bool()),
			MinimumEstimatedPortions: pointers.Float32(fake.Float32()),
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
