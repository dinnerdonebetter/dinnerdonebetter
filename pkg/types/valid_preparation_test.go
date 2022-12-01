package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidPreparationCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{
			Name:                  fake.LoremIpsumSentence(exampleQuantity),
			Description:           fake.LoremIpsumSentence(exampleQuantity),
			IconPath:              fake.LoremIpsumSentence(exampleQuantity),
			PastTense:             fake.LoremIpsumSentence(exampleQuantity),
			YieldsNothing:         fake.Bool(),
			RestrictToIngredients: fake.Bool(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{
			Name:                  stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Description:           stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			IconPath:              stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			PastTense:             stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			YieldsNothing:         boolPointer(fake.Bool()),
			RestrictToIngredients: boolPointer(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
