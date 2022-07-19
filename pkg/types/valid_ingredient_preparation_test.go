package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientPreparationCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationCreationRequestInput{
			Notes:              fake.LoremIpsumSentence(exampleQuantity),
			ValidPreparationID: fake.LoremIpsumSentence(exampleQuantity),
			ValidIngredientID:  fake.LoremIpsumSentence(exampleQuantity),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientPreparationUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationUpdateRequestInput{
			Notes:              stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			ValidPreparationID: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			ValidIngredientID:  stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
